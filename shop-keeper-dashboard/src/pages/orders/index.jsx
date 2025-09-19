import { useState, useMemo } from "react";
import { toast } from "sonner";
import { Search, Eye, Truck, CheckCircle, ArrowUpDown, Loader2, Plus, ShoppingCart, Package2, Download } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { shopApi } from '@/lib/api';
import { generateInvoicePDF } from "@/utils/pdfGenerator";

const ITEMS_PER_PAGE = 10;

const STATUS_OPTIONS = {
  ALL: "all",
  PENDING: "pending",
  PROCESSING: "processing",
  COMPLETED: "completed",
  CANCELLED: "cancelled"
};

const Orders = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState(STATUS_OPTIONS.ALL);
  const [currentPage, setCurrentPage] = useState(1);
  const [sortBy, setSortBy] = useState("created_at");
  const [sortDirection, setSortDirection] = useState("desc");

  // Query parameters
  const queryParams = useMemo(() => {
    const params = {
      sort_by: sortBy,
      sort_direction: sortDirection,
      page: currentPage,
      per_page: ITEMS_PER_PAGE,
    };

    if (searchQuery?.trim()) {
      params.search = searchQuery.trim();
    }

    if (statusFilter !== STATUS_OPTIONS.ALL) {
      params.status = statusFilter;
    }

    return params;
  }, [searchQuery, statusFilter, sortBy, sortDirection, currentPage]);

  const queryClient = useQueryClient();

  // API calls
  const { 
    data: ordersResponse, 
    isLoading, 
    isError, 
    error 
  } = useQuery({
    queryKey: ['orders', queryParams],
    queryFn: () => shopApi.getOrders(queryParams),
  });

  const updateOrderStatusMutation = useMutation({
    mutationFn: ({ id, status }) => shopApi.updateOrder(id, { status }),
    onSuccess: () => {
      queryClient.invalidateQueries(['orders']);
    },
  });

  // Extract data
  const orders = Array.isArray(ordersResponse?.data) ? ordersResponse.data : [];
  const ordersMeta = ordersResponse?.meta || {};
  
  // Calculate today's orders
  const todayOrders = useMemo(() => {
    const today = new Date().toDateString();
    return orders.filter(order => {
      const orderDate = new Date(order.created_at).toDateString();
      return orderDate === today;
    }).length;
  }, [orders]);

  const handleSort = (key) => {
    if (sortBy === key) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortBy(key);
      setSortDirection("asc");
    }
    setCurrentPage(1);
  };

  const handleUpdateStatus = async (orderId, newStatus) => {
    try {
      await updateOrderStatusMutation.mutateAsync({ id: orderId, status: newStatus });
      toast.success(`Order status updated to ${newStatus}`);
    } catch (error) {
      toast.error(error.message || "Failed to update order status");
    }
  };

  const handleDownloadInvoice = async (order) => {
    try {
      await generateInvoicePDF(order);
      toast.success("Invoice downloaded successfully!");
    } catch (error) {
      console.error("PDF generation error:", error);
      toast.error("Failed to generate invoice. Please try again.");
    }
  };

  const getSortIcon = (column) => {
    if (sortBy !== column) return <ArrowUpDown className="ml-1 h-3 w-3 sm:ml-2 sm:h-4 sm:w-4" />;
    return sortDirection === "asc" 
      ? <ArrowUpDown className="ml-1 h-3 w-3 sm:ml-2 sm:h-4 sm:w-4 rotate-180" />
      : <ArrowUpDown className="ml-1 h-3 w-3 sm:ml-2 sm:h-4 sm:w-4" />;
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const getStatusVariant = (status) => {
    switch (status?.toLowerCase()) {
      case 'completed':
        return 'default';
      case 'processing':
        return 'secondary';
      case 'pending':
        return 'warning';
      case 'cancelled':
        return 'destructive';
      default:
        return 'secondary';
    }
  };


  if (isError) return (
    <Card>
      <CardContent className="flex items-center justify-center h-[400px] text-destructive">
        Error loading orders: {error?.message || 'Please try again later'}
      </CardContent>
    </Card>
  );

  return (
    <div className="space-y-4 sm:space-y-6 p-3 sm:p-6">
      <Card>
        <CardHeader className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <CardTitle className="text-lg sm:text-xl">Orders</CardTitle>
            <CardDescription className="text-sm">
              View and manage customer orders
            </CardDescription>
          </div>
          <Button onClick={() => navigate("/orders/create")} className="w-full sm:w-auto">
            <Plus className="mr-2 h-4 w-4" /> 
            <span className="hidden sm:inline">Create Custom Order</span>
            <span className="sm:hidden">Create Order</span>
          </Button>
        </CardHeader>
        <CardContent>
          {/* Search and Filter */}
          <div className="flex flex-col gap-4 md:flex-row md:items-center mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search orders..."
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value);
                  setCurrentPage(1);
                }}
                className="pl-10"
              />
            </div>
            <Select value={statusFilter} onValueChange={(value) => {
              setStatusFilter(value);
              setCurrentPage(1);
            }}>
              <SelectTrigger className="w-full md:w-[180px]">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value={STATUS_OPTIONS.ALL}>All Orders</SelectItem>
                  <SelectItem value={STATUS_OPTIONS.PENDING}>Pending</SelectItem>
                  <SelectItem value={STATUS_OPTIONS.PROCESSING}>Processing</SelectItem>
                  <SelectItem value={STATUS_OPTIONS.COMPLETED}>Completed</SelectItem>
                  <SelectItem value={STATUS_OPTIONS.CANCELLED}>Cancelled</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>

          {/* Orders Table */}
          <div className="rounded-md border overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="min-w-[100px]">
                    <Button variant="ghost" onClick={() => handleSort("order_number")} className="text-xs sm:text-sm">
                      Order ID
                      {getSortIcon("order_number")}
                    </Button>
                  </TableHead>
                  <TableHead className="min-w-[150px]">
                    <Button variant="ghost" onClick={() => handleSort("user")} className="text-xs sm:text-sm">
                      Customer
                      {getSortIcon("user")}
                    </Button>
                  </TableHead>
                  <TableHead className="hidden md:table-cell">
                    <Button variant="ghost" onClick={() => handleSort("created_at")} className="text-xs sm:text-sm">
                      Date
                      {getSortIcon("created_at")}
                    </Button>
                  </TableHead>
                  <TableHead className="hidden sm:table-cell min-w-[60px]">Items</TableHead>
                  <TableHead className="min-w-[80px]">
                    <Button variant="ghost" onClick={() => handleSort("total_amount")} className="text-xs sm:text-sm">
                      Amount
                      {getSortIcon("total_amount")}
                    </Button>
                  </TableHead>
                  <TableHead className="min-w-[90px]">Status</TableHead>
                  <TableHead className="text-right min-w-[100px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center py-8">
                      <div className="flex items-center justify-center space-x-2">
                        <Loader2 className="h-5 w-5 animate-spin" />
                        <span>Loading orders...</span>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : orders.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center py-16">
                      {searchQuery?.trim() || statusFilter !== STATUS_OPTIONS.ALL ? (
                        <div className="text-muted-foreground">
                          <Search className="h-12 w-12 mx-auto mb-4 opacity-50" />
                          <p className="text-lg font-medium mb-2">No orders match your search</p>
                          <p className="text-sm">Try adjusting your search criteria or filters</p>
                        </div>
                      ) : (
                        <div className="text-center">
                          <div className="mx-auto w-24 h-24 bg-muted rounded-full flex items-center justify-center mb-6">
                            <ShoppingCart className="h-12 w-12 text-muted-foreground" />
                          </div>
                          <h3 className="text-xl font-semibold text-foreground mb-3">No orders yet</h3>
                          <p className="text-muted-foreground mb-6 max-w-md mx-auto">
                            Get started by creating your first custom order manually, or wait for customers to place orders through your store.
                          </p>
                          <div className="flex flex-col sm:flex-row gap-3 justify-center">
                            <Button onClick={() => navigate("/orders/create")}>
                              <Plus className="mr-2 h-4 w-4" />
                              Create Your First Order
                            </Button>
                            <Button
                              variant="outline"
                              onClick={() => navigate("/products")}
                            >
                              <Package2 className="mr-2 h-4 w-4" />
                              View Products
                            </Button>
                          </div>
                        </div>
                      )}
                    </TableCell>
                  </TableRow>
                ) : (
                  orders.map((order) => (
                    <TableRow key={order.id}>
                      <TableCell className="font-medium text-xs sm:text-sm">
                        #{order.order_number || order.id}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <Avatar className="h-6 w-6 sm:h-8 sm:w-8 flex-shrink-0">
                            <AvatarFallback className="text-xs">
                              {order.user?.name?.[0] || order.customer_name?.[0] || 'U'}
                            </AvatarFallback>
                          </Avatar>
                          <div className="min-w-0">
                            <span className="text-xs sm:text-sm truncate block max-w-[100px] sm:max-w-none">
                              {order.user?.name || order.customer_name || 'Unknown Customer'}
                            </span>
                            <span className="text-xs text-muted-foreground md:hidden block">
                              {formatDate(order.created_at)}
                            </span>
                            <span className="text-xs text-muted-foreground sm:hidden block">
                              {order.items?.length || order.total_items || 0} items
                            </span>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="text-muted-foreground text-xs sm:text-sm hidden md:table-cell">
                        {formatDate(order.created_at)}
                      </TableCell>
                      <TableCell className="text-xs sm:text-sm hidden sm:table-cell">
                        {order.items?.length || order.total_items || 0}
                      </TableCell>
                      <TableCell className="text-xs sm:text-sm font-medium">
                        ${Number(order.total_amount || order.total || 0).toFixed(2)}
                      </TableCell>
                      <TableCell>
                        <Badge 
                          variant={getStatusVariant(order.status)}
                          className="text-xs px-2 py-1"
                        >
                          <span className="hidden sm:inline capitalize">{order.status}</span>
                          <span className="sm:hidden capitalize">
                            {order.status === "completed" ? "Done" :
                             order.status === "processing" ? "Proc" :
                             order.status === "pending" ? "Pend" : "Canc"}
                          </span>
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-1">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => navigate(`/orders/details/${order.id}`)}
                            className="p-1 sm:p-2"
                            title="View Details"
                          >
                            <Eye className="h-3 w-3 sm:h-4 sm:w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDownloadInvoice(order)}
                            className="p-1 sm:p-2"
                            title="Download Invoice"
                          >
                            <Download className="h-3 w-3 sm:h-4 sm:w-4" />
                          </Button>
                          {order.status === "pending" && (
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleUpdateStatus(order.id, "processing")}
                              className="p-1 sm:p-2"
                              disabled={updateOrderStatusMutation.isPending}
                              title="Mark as Processing"
                            >
                              <Truck className="h-3 w-3 sm:h-4 sm:w-4" />
                            </Button>
                          )}
                          {order.status === "processing" && (
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleUpdateStatus(order.id, "completed")}
                              className="p-1 sm:p-2"
                              disabled={updateOrderStatusMutation.isPending}
                              title="Mark as Completed"
                            >
                              <CheckCircle className="h-3 w-3 sm:h-4 sm:w-4" />
                            </Button>
                          )}
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {ordersMeta.last_page > 1 && (
            <div className="flex justify-center mt-6">
              <Pagination>
                <PaginationContent className="flex items-center gap-2">
                  <PaginationItem>
                    <PaginationPrevious
                      onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                      disabled={currentPage === 1}
                    />
                  </PaginationItem>

                  <PaginationItem>
                    <PaginationLink disabled className="pointer-events-none">
                      Page {currentPage} of {ordersMeta.last_page}
                    </PaginationLink>
                  </PaginationItem>

                  <PaginationItem>
                    <PaginationNext
                      onClick={() => setCurrentPage(p => Math.min(ordersMeta.last_page, p + 1))}
                      disabled={currentPage === ordersMeta.last_page}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default Orders;