import { useState, useMemo } from "react";
import { Edit, Trash2, Search, Plus, Loader2, ArrowUpDown, ArrowUp, ArrowDown, FolderTree, ChevronRight, ChevronDown } from "lucide-react";
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
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { DeleteConfirmationDialog } from "@/components/ui/delete-confirmation-dialog";
import { toast } from "sonner";

const ITEMS_PER_PAGE = 15;

// Enhanced mock data with hierarchical categories
const generateCategoriesData = () => {
  const categories = [
    // Electronics (parent)
    { 
      id: 1, 
      name: "Electronics", 
      slug: "electronics",
      description: "Electronic devices, gadgets, and accessories", 
      parent_id: null,
      level: 0,
      products_count: 156,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-15",
      children_count: 4
    },
    // Electronics children
    { 
      id: 2, 
      name: "Smartphones", 
      slug: "smartphones",
      description: "Mobile phones and accessories", 
      parent_id: 1,
      level: 1,
      products_count: 45,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-16",
      children_count: 3
    },
    { 
      id: 3, 
      name: "iPhone", 
      slug: "iphone",
      description: "Apple iPhone devices", 
      parent_id: 2,
      level: 2,
      products_count: 15,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-17",
      children_count: 0
    },
    { 
      id: 4, 
      name: "Android", 
      slug: "android",
      description: "Android devices", 
      parent_id: 2,
      level: 2,
      products_count: 25,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-18",
      children_count: 0
    },
    { 
      id: 5, 
      name: "Accessories", 
      slug: "phone-accessories",
      description: "Phone cases, chargers, etc.", 
      parent_id: 2,
      level: 2,
      products_count: 5,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-19",
      children_count: 0
    },
    { 
      id: 6, 
      name: "Laptops", 
      slug: "laptops",
      description: "Portable computers", 
      parent_id: 1,
      level: 1,
      products_count: 67,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-20",
      children_count: 0
    },
    { 
      id: 7, 
      name: "Headphones", 
      slug: "headphones",
      description: "Audio devices", 
      parent_id: 1,
      level: 1,
      products_count: 34,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-21",
      children_count: 0
    },
    { 
      id: 8, 
      name: "Gaming", 
      slug: "gaming",
      description: "Gaming consoles and accessories", 
      parent_id: 1,
      level: 1,
      products_count: 10,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-22",
      children_count: 0
    },

    // Fashion (parent)
    { 
      id: 9, 
      name: "Fashion", 
      slug: "fashion",
      description: "Clothing, shoes, and fashion accessories", 
      parent_id: null,
      level: 0,
      products_count: 234,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-23",
      children_count: 3
    },
    // Fashion children
    { 
      id: 10, 
      name: "Men's Clothing", 
      slug: "mens-clothing",
      description: "Clothing for men", 
      parent_id: 9,
      level: 1,
      products_count: 89,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-24",
      children_count: 0
    },
    { 
      id: 11, 
      name: "Women's Clothing", 
      slug: "womens-clothing",
      description: "Clothing for women", 
      parent_id: 9,
      level: 1,
      products_count: 112,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-25",
      children_count: 0
    },
    { 
      id: 12, 
      name: "Shoes", 
      slug: "shoes",
      description: "Footwear for all", 
      parent_id: 9,
      level: 1,
      products_count: 33,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-26",
      children_count: 0
    },

    // Home & Garden (parent)
    { 
      id: 13, 
      name: "Home & Garden", 
      slug: "home-garden",
      description: "Home improvement and garden supplies", 
      parent_id: null,
      level: 0,
      products_count: 78,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-27",
      children_count: 2
    },
    { 
      id: 14, 
      name: "Kitchen", 
      slug: "kitchen",
      description: "Kitchen appliances and tools", 
      parent_id: 13,
      level: 1,
      products_count: 45,
      is_active: true,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-28",
      children_count: 0
    },
    { 
      id: 15, 
      name: "Garden Tools", 
      slug: "garden-tools",
      description: "Tools for gardening", 
      parent_id: 13,
      level: 1,
      products_count: 33,
      is_active: false,
      image: "/api/placeholder/100/100",
      created_at: "2024-01-29",
      children_count: 0
    },
  ];

  return categories;
};

const Category = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [sortBy, setSortBy] = useState('created_at');
  const [sortDirection, setSortDirection] = useState('desc');
  const [viewMode, setViewMode] = useState('tree'); // 'table' or 'tree'
  const [expandedCategories, setExpandedCategories] = useState(new Set([1, 9, 13])); // Initially expand parent categories
  const [statusFilter, setStatusFilter] = useState('all');
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [categoryToDelete, setCategoryToDelete] = useState(null);
  const [isLoading] = useState(false);

  // Mock data
  const allCategories = generateCategoriesData();

  // Build hierarchy
  const buildCategoryTree = (categories, parentId = null) => {
    return categories
      .filter(cat => cat.parent_id === parentId)
      .map(cat => ({
        ...cat,
        children: buildCategoryTree(categories, cat.id)
      }));
  };

  const categoryTree = buildCategoryTree(allCategories);

  // Filter and search categories
  const filteredCategories = useMemo(() => {
    let filtered = allCategories;

    // Apply status filter
    if (statusFilter !== 'all') {
      filtered = filtered.filter(cat => 
        statusFilter === 'active' ? cat.is_active : !cat.is_active
      );
    }

    // Apply search filter
    if (searchQuery?.trim()) {
      const query = searchQuery.trim().toLowerCase();
      filtered = filtered.filter(cat => 
        cat.name.toLowerCase().includes(query) ||
        cat.description.toLowerCase().includes(query) ||
        cat.slug.toLowerCase().includes(query)
      );
    }

    return filtered;
  }, [searchQuery, statusFilter, allCategories]);

  // Pagination
  const totalPages = Math.ceil(filteredCategories.length / ITEMS_PER_PAGE);
  const paginatedCategories = filteredCategories.slice(
    (currentPage - 1) * ITEMS_PER_PAGE,
    currentPage * ITEMS_PER_PAGE
  );

  // Helper functions
  const handleSort = (column) => {
    if (sortBy === column) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(column);
      setSortDirection('asc');
    }
    setCurrentPage(1);
  };

  const getSortIcon = (column) => {
    if (sortBy !== column) return <ArrowUpDown className="h-4 w-4" />;
    return sortDirection === 'asc' ? <ArrowUp className="h-4 w-4" /> : <ArrowDown className="h-4 w-4" />;
  };

  const toggleExpanded = (categoryId) => {
    const newExpanded = new Set(expandedCategories);
    if (newExpanded.has(categoryId)) {
      newExpanded.delete(categoryId);
    } else {
      newExpanded.add(categoryId);
    }
    setExpandedCategories(newExpanded);
  };

  const handleDelete = (category) => {
    setCategoryToDelete(category);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = async () => {
    if (categoryToDelete) {
      try {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 1000));
        toast.success('Category deleted successfully');
        setDeleteDialogOpen(false);
        setCategoryToDelete(null);
      } catch {
        toast.error('Failed to delete category');
      }
    }
  };

  const renderCategoryRow = (category) => (
    <TableRow key={category.id}>
      <TableCell>
        <div className="flex items-center gap-2" style={{ paddingLeft: `${category.level * 24}px` }}>
          {category.children_count > 0 && (
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={() => toggleExpanded(category.id)}
            >
              {expandedCategories.has(category.id) ? 
                <ChevronDown className="h-4 w-4" /> : 
                <ChevronRight className="h-4 w-4" />
              }
            </Button>
          )}
          {category.children_count === 0 && <div className="w-6" />}
          
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gray-100 rounded-lg flex items-center justify-center">
              <FolderTree className="h-5 w-5 text-gray-600" />
            </div>
            <div>
              <div className="font-medium flex items-center gap-2">
                {category.name}
                {category.level > 0 && (
                  <Badge variant="outline" className="text-xs">
                    L{category.level}
                  </Badge>
                )}
              </div>
              <div className="text-xs text-muted-foreground">
                {category.slug}
              </div>
            </div>
          </div>
        </div>
      </TableCell>
      <TableCell className="text-muted-foreground">
        <div className="truncate" title={category.description}>
          {category.description || "No description"}
        </div>
      </TableCell>
      <TableCell>
        <Badge variant="outline">
          {category.products_count} products
        </Badge>
      </TableCell>
      <TableCell>
        {category.children_count > 0 && (
          <Badge variant="secondary" className="mr-2">
            {category.children_count} subcategories
          </Badge>
        )}
        <Badge variant={category.is_active ? "default" : "secondary"}>
          {category.is_active ? "Active" : "Inactive"}
        </Badge>
      </TableCell>
      <TableCell className="text-right">
        <div className="flex items-center justify-end gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate(`/category/${category.id}/edit`)}
            title="Edit Category"
          >
            <Edit className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => handleDelete(category)}
            className="text-destructive hover:text-destructive/90"
            title="Delete Category"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </TableCell>
    </TableRow>
  );

  const renderTreeView = (categories, level = 0) => {
    return categories.map(category => {
      const isExpanded = expandedCategories.has(category.id);
      const hasChildren = category.children && category.children.length > 0;
      
      return (
        <div key={category.id}>
          {renderCategoryRow({ ...category, level })}
          {hasChildren && isExpanded && renderTreeView(category.children, level + 1)}
        </div>
      );
    });
  };

  return (
    <div className="space-y-6 p-6">
      <Card className="w-full">
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <FolderTree className="h-5 w-5" />
              Categories ({allCategories.length})
            </CardTitle>
            <CardDescription>
              Manage your product categories with hierarchical organization
            </CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => setViewMode(viewMode === 'tree' ? 'table' : 'tree')}
            >
              {viewMode === 'tree' ? 'Table View' : 'Tree View'}
            </Button>
            <Button onClick={() => navigate("/category/create")}>
              <Plus className="mr-2 h-4 w-4" /> Add Category
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {/* Search and Filter section */}
          <div className="flex flex-col gap-4 md:flex-row md:items-center mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search by name, description, or slug..."
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value);
                  setCurrentPage(1);
                }}
                className="pl-10"
              />
            </div>
            <div className="flex gap-2">
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Categories Display */}
          <div className="w-full overflow-x-auto">
            <div className="rounded-md border">
              <Table className="w-full min-w-full">
              <TableHeader>
                <TableRow>
                  <TableHead 
                    className="cursor-pointer hover:bg-muted/50"
                    onClick={() => handleSort('name')}
                  >
                    <div className="flex items-center space-x-2">
                      <span>Category</span>
                      {getSortIcon('name')}
                    </div>
                  </TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead 
                    className="cursor-pointer hover:bg-muted/50"
                    onClick={() => handleSort('products_count')}
                  >
                    <div className="flex items-center space-x-2">
                      <span>Products</span>
                      {getSortIcon('products_count')}
                    </div>
                  </TableHead>
                  <TableHead 
                    className="cursor-pointer hover:bg-muted/50"
                    onClick={() => handleSort('status')}
                  >
                    <div className="flex items-center space-x-2">
                      <span>Status</span>
                      {getSortIcon('status')}
                    </div>
                  </TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center py-8">
                      <div className="flex items-center justify-center space-x-2">
                        <Loader2 className="h-5 w-5 animate-spin" />
                        <span>Loading categories...</span>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : viewMode === 'tree' ? (
                  // Tree View
                  categoryTree.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                        No categories found
                      </TableCell>
                    </TableRow>
                  ) : (
                    renderTreeView(categoryTree)
                  )
                ) : (
                  // Table View
                  paginatedCategories.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                        {searchQuery?.trim() ? `No categories found matching "${searchQuery.trim()}"` : "No categories found"}
                      </TableCell>
                    </TableRow>
                  ) : (
                    paginatedCategories.map((category) => renderCategoryRow(category))
                  )
                )}
              </TableBody>
            </Table>
            </div>
          </div>

          {/* Pagination - only show in table view */}
          {viewMode === 'table' && totalPages > 1 && (
            <div className="flex items-center justify-between mt-6">
              <div className="text-sm text-muted-foreground">
                Showing {((currentPage - 1) * ITEMS_PER_PAGE) + 1} to {Math.min(currentPage * ITEMS_PER_PAGE, filteredCategories.length)} of {filteredCategories.length} categories
              </div>
              <Pagination>
                <PaginationContent className="flex items-center gap-2">
                  <PaginationItem>
                    <PaginationPrevious
                      onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                      className={currentPage === 1 ? "pointer-events-none opacity-50" : "cursor-pointer"}
                    />
                  </PaginationItem>

                  {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                    const pageNumber = Math.max(1, Math.min(totalPages - 4, currentPage - 2)) + i;
                    if (pageNumber > totalPages) return null;
                    
                    return (
                      <PaginationItem key={pageNumber}>
                        <PaginationLink
                          onClick={() => setCurrentPage(pageNumber)}
                          isActive={currentPage === pageNumber}
                          className="cursor-pointer"
                        >
                          {pageNumber}
                        </PaginationLink>
                      </PaginationItem>
                    );
                  })}

                  <PaginationItem>
                    <PaginationNext
                      onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                      className={currentPage === totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <DeleteConfirmationDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onConfirm={confirmDelete}
        itemName={categoryToDelete?.name}
      />
    </div>
  );
};

export default Category;