import { useState, useMemo } from "react";
import { toast } from "sonner";
import { 
  Plus, 
  Search, 
  ArrowUpDown,
  Package,
  Loader2,
  Eye,
  PencilLine,
  Trash,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { DeleteConfirmationDialog } from "@/components/ui/delete-confirmation-dialog";
import { useNavigate } from "react-router-dom";
import { useProducts, useDeleteProduct, useCategories } from "@/hooks/useApi";

const ITEMS_PER_PAGE = 10;

const Products = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("");
  const [sortBy, setSortBy] = useState("created_at");
  const [sortDirection, setSortDirection] = useState("desc");
  const [currentPage, setCurrentPage] = useState(1);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [productToDelete, setProductToDelete] = useState(null);

  const navigate = useNavigate();

  // Query parameters - TanStack Query will handle debouncing via its caching
  const queryParams = useMemo(() => {
    const params = {
      sort_by: sortBy,
      sort_direction: sortDirection,
      page: currentPage,
      per_page: ITEMS_PER_PAGE,
    };
    
    // Only add search if it has a value
    if (searchQuery?.trim()) {
      params.search = searchQuery.trim();
    }
    
    // Only add category_id if it has a value
    if (selectedCategory) {
      params.category_id = selectedCategory;
    }
    
    return params;
  }, [searchQuery, selectedCategory, sortBy, sortDirection, currentPage]);

  // API calls with debouncing
  const { 
    data: productsResponse, 
    isLoading, 
    isError, 
    error 
  } = useProducts(queryParams);

  const { data: categoriesResponse } = useCategories();

  const deleteProductMutation = useDeleteProduct();

  // Extract data
  const products = Array.isArray(productsResponse?.data) ? productsResponse.data : [];
  const productsMeta = productsResponse?.meta || {};
  const categories = Array.isArray(categoriesResponse?.data) ? categoriesResponse.data : [];

  const handleSort = (key) => {
    if (sortBy === key) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortBy(key);
      setSortDirection("asc");
    }
    setCurrentPage(1);
  };

  const handleDelete = (product) => {
    setProductToDelete(product);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = async () => {
    if (!productToDelete) return;

    try {
      await deleteProductMutation.mutateAsync(productToDelete.id);
      toast.success("Product deleted successfully");
      setDeleteDialogOpen(false);
      setProductToDelete(null);
    } catch (error) {
      toast.error(error.message || "Failed to delete product");
    }
  };

  const handleCategoryChange = (value) => {
    setSelectedCategory(value === "all" ? "" : value);
    setCurrentPage(1);
  };

  const getSortIcon = (column) => {
    if (sortBy !== column) return <ArrowUpDown className="ml-1 h-3 w-3 sm:ml-2 sm:h-4 sm:w-4" />;
    return sortDirection === "asc" 
      ? <ArrowUpDown className="ml-1 h-3 w-3 sm:ml-2 sm:h-4 sm:w-4 rotate-180" />
      : <ArrowUpDown className="ml-1 h-3 w-3 sm:ml-2 sm:h-4 sm:w-4" />;
  };

  // Error handling
  if (isError) {
    return (
      <Card>
        <CardContent className="flex items-center justify-center h-[400px] text-destructive">
          Error loading products: {error?.message || 'Please try again later'}
        </CardContent>
      </Card>
    );
  }


  return (
    <div className="space-y-4 sm:space-y-6 p-3 sm:p-6">
      <Card>
        <CardHeader className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <CardTitle className="text-lg sm:text-xl">Product Inventory</CardTitle>
            <CardDescription className="text-sm">
              View and manage all your products in one place
            </CardDescription>
          </div>
          <Button onClick={() => navigate("/products/create")} className="w-full sm:w-auto">
            <Plus className="mr-2 h-4 w-4" /> 
            <span className="hidden xs:inline">Add Product</span>
            <span className="xs:hidden">Add</span>
          </Button>
        </CardHeader>
        <CardContent>
          {/* Search and filter section */}
          <div className="flex flex-col gap-4 md:flex-row md:items-center mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search products..."
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value);
                  setCurrentPage(1); // Reset to first page when searching
                }}
                className="pl-10"
              />
            </div>
            <Select value={selectedCategory || "all"} onValueChange={handleCategoryChange}>
              <SelectTrigger className="w-full md:w-[180px]">
                <SelectValue placeholder="All Categories" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Categories</SelectLabel>
                  <SelectItem value="all">All Categories</SelectItem>
                  {categories.map((category) => (
                    <SelectItem key={category.id} value={category.id.toString()}>
                      {category.name}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>

          {/* Products Table */}
          <div className="rounded-md border overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="min-w-[200px]">Product</TableHead>
                  <TableHead className="hidden md:table-cell">
                    <Button variant="ghost" onClick={() => handleSort("category")} className="text-xs sm:text-sm">
                      Category
                      {getSortIcon("category")}
                    </Button>
                  </TableHead>
                  <TableHead className="min-w-[80px]">
                    <Button variant="ghost" onClick={() => handleSort("price")} className="text-xs sm:text-sm">
                      Price
                      {getSortIcon("price")}
                    </Button>
                  </TableHead>
                  <TableHead className="hidden sm:table-cell min-w-[70px]">
                    <Button variant="ghost" onClick={() => handleSort("stock")} className="text-xs sm:text-sm">
                      Stock
                      {getSortIcon("stock")}
                    </Button>
                  </TableHead>
                  <TableHead className="min-w-[90px]">Status</TableHead>
                  <TableHead className="text-right min-w-[120px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8">
                      <div className="flex items-center justify-center space-x-2">
                        <Loader2 className="h-5 w-5 animate-spin" />
                        <span>Loading products...</span>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : products.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8">
                      <div className="flex flex-col items-center justify-center">
                        <Package className="h-8 w-8 text-muted-foreground mb-2" />
                        <span className="text-muted-foreground">
                          {searchQuery?.trim() || selectedCategory 
                            ? "No products match your search" 
                            : "No products found"
                          }
                        </span>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  products.map((product) => (
                  <TableRow key={product.id}>
                    <TableCell>
                      <div className="flex items-center gap-2 sm:gap-3">
                        <div className="h-8 w-8 sm:h-12 sm:w-12 rounded-lg border bg-muted flex items-center justify-center flex-shrink-0">
                          {product.images && product.images[0] ? (
                            <img 
                              src={product.images[0]} 
                              alt={product.name}
                              className="h-full w-full object-cover rounded-lg"
                            />
                          ) : (
                            <Package className="h-4 w-4 sm:h-6 sm:w-6 text-muted-foreground" />
                          )}
                        </div>
                        <div className="flex flex-col min-w-0">
                          <span className="font-medium text-xs sm:text-sm truncate">{product.name}</span>
                          <span className="text-xs text-muted-foreground">#{product.sku || product.id}</span>
                          <span className="text-xs text-muted-foreground md:hidden">
                            {product.category?.name || 'No Category'}
                          </span>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className="hidden md:table-cell text-xs sm:text-sm">
                      {product.category?.name || 'No Category'}
                    </TableCell>
                    <TableCell className="text-xs sm:text-sm font-medium">
                      ${Number(product.price).toFixed(2)}
                    </TableCell>
                    <TableCell className="hidden sm:table-cell text-xs sm:text-sm">{product.stock}</TableCell>
                    <TableCell>
                      <Badge 
                        variant={product.stock > 20 ? "default" : product.stock > 0 ? "secondary" : "destructive"}
                        className="text-xs px-2 py-1"
                      >
                        <span className="hidden sm:inline">
                          {product.stock === 0 ? "Out of Stock" : product.stock <= 10 ? "Low Stock" : "In Stock"}
                        </span>
                        <span className="sm:hidden">
                          {product.stock === 0 ? "Out" : product.stock <= 10 ? "Low" : "In"}
                        </span>
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => navigate(`/products/${product.id}`)}
                          className="p-1 sm:p-2"
                        >
                          <Eye className="h-3 w-3 sm:h-4 sm:w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => navigate(`/products/${product.id}/edit`)}
                          className="p-1 sm:p-2"
                        >
                          <PencilLine className="h-3 w-3 sm:h-4 sm:w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleDelete(product)}
                          className="text-destructive hover:text-destructive/90 p-1 sm:p-2"
                          disabled={deleteProductMutation.isPending}
                        >
                          <Trash className="h-3 w-3 sm:h-4 sm:w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {productsMeta.total > 0 && (
            <div className="flex flex-col sm:flex-row items-center justify-between mt-6 gap-4">
              <div className="text-sm text-muted-foreground">
                Showing {((productsMeta.current_page - 1) * productsMeta.per_page) + 1} to{' '}
                {Math.min(productsMeta.current_page * productsMeta.per_page, productsMeta.total)} of{' '}
                {productsMeta.total} products
              </div>
              
              {productsMeta.last_page > 1 && (
                <Pagination>
                  <PaginationContent className="flex items-center gap-2">
                    <PaginationItem>
                      <PaginationPrevious
                        onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                        disabled={currentPage === 1 || isLoading}
                        className="cursor-pointer"
                      />
                    </PaginationItem>

                    <PaginationItem>
                      <PaginationLink isActive={true}>
                        {currentPage}
                      </PaginationLink>
                    </PaginationItem>

                    <PaginationItem>
                      <PaginationLink disabled className="pointer-events-none">
                        of {productsMeta.last_page}
                      </PaginationLink>
                    </PaginationItem>

                    <PaginationItem>
                      <PaginationNext
                        onClick={() => setCurrentPage(Math.min(productsMeta.last_page, currentPage + 1))}
                        disabled={currentPage === productsMeta.last_page || isLoading}
                        className="cursor-pointer"
                      />
                    </PaginationItem>
                  </PaginationContent>
                </Pagination>
              )}
            </div>
          )}
          
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <DeleteConfirmationDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onConfirm={confirmDelete}
        itemName={productToDelete?.name}
      />
    </div>
  );
};

export default Products;