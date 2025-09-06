import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Search, Mail, CheckCircle, XCircle, Loader2, Eye, ArrowUpDown, Plus, Trash2, PencilLine, Filter, Clock, AlertCircle, Zap, Download, RefreshCw } from "lucide-react";
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
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";

const ITEMS_PER_PAGE = 5;

const mockTickets = [
  { id: 1, subject: "Order not received", customer: "John Doe", email: "john.doe@example.com", date: "2025-05-05", status: "Open", priority: "high", category: "order", lastActivity: "2 hours ago" },
  { id: 2, subject: "Refund request", customer: "Jane Smith", email: "jane.smith@example.com", date: "2025-05-04", status: "Pending", priority: "medium", category: "payment", lastActivity: "1 day ago" },
  { id: 3, subject: "Product damaged", customer: "Emily Davis", email: "emily.davis@example.com", date: "2025-05-03", status: "Closed", priority: "urgent", category: "product", lastActivity: "3 days ago" },
  { id: 4, subject: "Payment issue", customer: "Michael Brown", email: "michael.b@example.com", date: "2025-05-02", status: "Open", priority: "low", category: "payment", lastActivity: "5 hours ago" },
  { id: 5, subject: "Login issues", customer: "Sarah Wilson", email: "sarah.w@example.com", date: "2025-05-01", status: "Open", priority: "medium", category: "technical", lastActivity: "30 minutes ago" },
  { id: 6, subject: "Shipping delay inquiry", customer: "Robert Taylor", email: "rob.taylor@example.com", date: "2025-04-30", status: "Pending", priority: "low", category: "shipping", lastActivity: "4 hours ago" },
];

export default function Support() {
  const [search, setSearch] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState("all");
  const [priorityFilter, setPriorityFilter] = useState("all");
  const [categoryFilter, setCategoryFilter] = useState("all");
  const [sortBy, setSortBy] = useState("date");
  const [sortOrder, setSortOrder] = useState("desc");

  const { data: tickets = [], refetch } = useQuery({
    queryKey: ["support-tickets"],
    queryFn: async () => mockTickets,
  });
  const navigate = useNavigate();

  const getPriorityIcon = (priority) => {
    switch (priority) {
      case 'low': return <Clock className="w-4 h-4 text-blue-500" />;
      case 'medium': return <AlertCircle className="w-4 h-4 text-yellow-500" />;
      case 'high': return <AlertCircle className="w-4 h-4 text-orange-500" />;
      case 'urgent': return <Zap className="w-4 h-4 text-red-500" />;
      default: return null;
    }
  };

  const getPriorityColor = (priority) => {
    switch (priority) {
      case 'low': return 'text-blue-600 bg-blue-50 border-blue-200';
      case 'medium': return 'text-yellow-600 bg-yellow-50 border-yellow-200';
      case 'high': return 'text-orange-600 bg-orange-50 border-orange-200';
      case 'urgent': return 'text-red-600 bg-red-50 border-red-200';
      default: return 'text-gray-600 bg-gray-50 border-gray-200';
    }
  };

  const handleBulkAction = (action) => {
    toast.success(`Bulk action "${action}" applied to selected tickets`);
  };

  const handleExport = () => {
    toast.success("Tickets exported successfully");
  };

  let filteredTickets = tickets.filter(ticket => {
    const matchesSearch = search === "" || 
      ticket.subject.toLowerCase().includes(search.toLowerCase()) ||
      ticket.customer.toLowerCase().includes(search.toLowerCase()) ||
      ticket.email.toLowerCase().includes(search.toLowerCase());
    
    const matchesStatus = statusFilter === "all" || ticket.status.toLowerCase() === statusFilter.toLowerCase();
    const matchesPriority = priorityFilter === "all" || ticket.priority === priorityFilter;
    const matchesCategory = categoryFilter === "all" || ticket.category === categoryFilter;
    
    return matchesSearch && matchesStatus && matchesPriority && matchesCategory;
  });

  // Sort tickets
  filteredTickets = [...filteredTickets].sort((a, b) => {
    let comparison = 0;
    switch (sortBy) {
      case 'subject':
        comparison = a.subject.localeCompare(b.subject);
        break;
      case 'customer':
        comparison = a.customer.localeCompare(b.customer);
        break;
      case 'date':
        comparison = new Date(a.date) - new Date(b.date);
        break;
      case 'priority': {
        const priorityOrder = { low: 1, medium: 2, high: 3, urgent: 4 };
        comparison = priorityOrder[a.priority] - priorityOrder[b.priority];
        break;
      }
      default:
        comparison = 0;
    }
    return sortOrder === 'asc' ? comparison : -comparison;
  });

  const totalPages = Math.ceil(filteredTickets.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const paginatedTickets = filteredTickets.slice(startIndex, startIndex + ITEMS_PER_PAGE);

  const handleSort = (column) => {
    if (sortBy === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(column);
      setSortOrder('asc');
    }
  };

  return (
    <div className="space-y-6 p-6">
      {/* Header Stats */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="flex items-center p-6">
            <div className="flex items-center space-x-2">
              <CheckCircle className="h-8 w-8 text-green-600" />
              <div>
                <p className="text-2xl font-bold">{tickets.filter(t => t.status === 'Open').length}</p>
                <p className="text-xs text-muted-foreground">Open Tickets</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="flex items-center p-6">
            <div className="flex items-center space-x-2">
              <Loader2 className="h-8 w-8 text-yellow-600" />
              <div>
                <p className="text-2xl font-bold">{tickets.filter(t => t.status === 'Pending').length}</p>
                <p className="text-xs text-muted-foreground">Pending</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="flex items-center p-6">
            <div className="flex items-center space-x-2">
              <XCircle className="h-8 w-8 text-gray-600" />
              <div>
                <p className="text-2xl font-bold">{tickets.filter(t => t.status === 'Closed').length}</p>
                <p className="text-xs text-muted-foreground">Closed</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="flex items-center p-6">
            <div className="flex items-center space-x-2">
              <Zap className="h-8 w-8 text-red-600" />
              <div>
                <p className="text-2xl font-bold">{tickets.filter(t => t.priority === 'urgent').length}</p>
                <p className="text-xs text-muted-foreground">Urgent</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Support Tickets</CardTitle>
            <CardDescription>View and manage all your support tickets ({filteredTickets.length} total)</CardDescription>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" onClick={handleExport}>
              <Download className="mr-2 h-4 w-4" />
              Export
            </Button>
            <Button variant="outline" size="sm" onClick={() => refetch()}>
              <RefreshCw className="mr-2 h-4 w-4" />
              Refresh
            </Button>
            <Button onClick={() => navigate("/support/create")}>
              <Plus className="mr-2 h-4 w-4" /> New Ticket
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {/* Search and Filters */}
          <div className="flex flex-col gap-4 mb-6">
            <div className="flex flex-col gap-4 lg:flex-row lg:items-center">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Search tickets..."
                  value={search}
                  onChange={e => setSearch(e.target.value)}
                  className="pl-10"
                />
              </div>
              <div className="flex gap-2">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="outline" size="sm">
                      <Filter className="mr-2 h-4 w-4" />
                      Actions
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent>
                    <DropdownMenuItem onClick={() => handleBulkAction("mark-resolved")}>
                      Mark as Resolved
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => handleBulkAction("assign-agent")}>
                      Assign Agent
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onClick={() => handleBulkAction("delete")}>
                      Delete Selected
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
            
            <div className="flex flex-wrap gap-2">
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="open">Open</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="closed">Closed</SelectItem>
                </SelectContent>
              </Select>
              
              <Select value={priorityFilter} onValueChange={setPriorityFilter}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Priority" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Priority</SelectItem>
                  <SelectItem value="low">Low</SelectItem>
                  <SelectItem value="medium">Medium</SelectItem>
                  <SelectItem value="high">High</SelectItem>
                  <SelectItem value="urgent">Urgent</SelectItem>
                </SelectContent>
              </Select>
              
              <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                <SelectTrigger className="w-36">
                  <SelectValue placeholder="Category" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Categories</SelectItem>
                  <SelectItem value="general">General</SelectItem>
                  <SelectItem value="order">Order</SelectItem>
                  <SelectItem value="payment">Payment</SelectItem>
                  <SelectItem value="product">Product</SelectItem>
                  <SelectItem value="shipping">Shipping</SelectItem>
                  <SelectItem value="technical">Technical</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Support Tickets Table */}
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-12">
                    <input type="checkbox" className="rounded" />
                  </TableHead>
                  <TableHead>
                    <Button variant="ghost" onClick={() => handleSort('subject')}>
                      Subject
                      <ArrowUpDown className="ml-2 h-4 w-4" />
                    </Button>
                  </TableHead>
                  <TableHead>
                    <Button variant="ghost" onClick={() => handleSort('customer')}>
                      Customer
                      <ArrowUpDown className="ml-2 h-4 w-4" />
                    </Button>
                  </TableHead>
                  <TableHead>Priority</TableHead>
                  <TableHead>Category</TableHead>
                  <TableHead>
                    <Button variant="ghost" onClick={() => handleSort('date')}>
                      Created
                      <ArrowUpDown className="ml-2 h-4 w-4" />
                    </Button>
                  </TableHead>
                  <TableHead>Last Activity</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {paginatedTickets.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={9} className="text-center text-muted-foreground h-32">
                      No tickets found
                    </TableCell>
                  </TableRow>
                ) : (
                  paginatedTickets.map((ticket) => (
                    <TableRow key={ticket.id} className="hover:bg-muted/50">
                      <TableCell>
                        <input type="checkbox" className="rounded" />
                      </TableCell>
                      <TableCell className="font-medium">
                        <div className="cursor-pointer" onClick={() => navigate(`/support/details/${ticket.id}`)}>
                          {ticket.subject}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <div>
                            <div className="font-medium">{ticket.customer}</div>
                            <div className="text-sm text-muted-foreground flex items-center gap-1">
                              <Mail className="w-3 h-3" />
                              {ticket.email}
                            </div>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className={`inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border ${getPriorityColor(ticket.priority)}`}>
                          {getPriorityIcon(ticket.priority)}
                          {ticket.priority}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline" className="capitalize">
                          {ticket.category}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-muted-foreground">{ticket.date}</TableCell>
                      <TableCell className="text-sm text-muted-foreground">{ticket.lastActivity}</TableCell>
                      <TableCell>
                        <Badge variant={
                          ticket.status === "Open" ? "default" :
                          ticket.status === "Pending" ? "secondary" :
                          "destructive"
                        }>
                          {ticket.status === "Open" && <CheckCircle className="w-4 h-4 mr-1 inline" />}
                          {ticket.status === "Pending" && <Loader2 className="w-4 h-4 mr-1 inline animate-spin" />}
                          {ticket.status === "Closed" && <XCircle className="w-4 h-4 mr-1 inline" />}
                          {ticket.status}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button variant="ghost" size="icon" title="View Details" onClick={() => navigate(`/support/details/${ticket.id}`)}>
                            <Eye className="h-4 w-4" />
                          </Button>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="icon">
                                <PencilLine className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem>Edit Ticket</DropdownMenuItem>
                              <DropdownMenuItem>Assign Agent</DropdownMenuItem>
                              <DropdownMenuItem>Change Priority</DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem className="text-destructive">
                                <Trash2 className="mr-2 h-4 w-4" />
                                Delete
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
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
                    <PaginationLink isActive>{currentPage}</PaginationLink>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLink disabled className="pointer-events-none">
                      of {totalPages}
                    </PaginationLink>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationNext
                      onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                      disabled={currentPage === totalPages}
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
}