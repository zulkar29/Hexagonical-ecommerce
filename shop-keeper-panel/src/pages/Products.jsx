import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Plus,
  Edit,
  Trash2,
  Package,
  Eye,
  LayoutGrid,
  List,
  Search,
  Filter,
  Download,
  Upload,
  MoreVertical,
  AlertTriangle,
  TrendingUp,
  TrendingDown,
  CheckCircle,
  XCircle,
  BarChart3,
  Settings,
  Copy,
  Archive,
  RefreshCw,
  ScanLine,
  Tag,
  DollarSign,
  Package2,
  Users,
  Calendar,
  Layers,
  SortAsc,
  SortDesc,
  FileText,
  Printer
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuCheckboxItem,
} from '@/components/ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from '@/components/ui/dialog';
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';

const EnhancedProductsPage = () => {
  const navigate = useNavigate();
  const [search, setSearch] = useState('');
  const [view, setView] = useState('table');
  const [selectedProducts, setSelectedProducts] = useState([]);
  const [sortBy, setSortBy] = useState('name');
  const [sortOrder, setSortOrder] = useState('asc');
  const [categoryFilter, setCategoryFilter] = useState('all');
  const [stockFilter, setStockFilter] = useState('all');
  const [statusFilter, setStatusFilter] = useState('all');
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showBulkActions, setShowBulkActions] = useState(false);
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const [currentPage, setCurrentPage] = useState(1);

  // Enhanced mock products with more data
  const mockProducts = [
    {
      id: 1,
      name: 'Basmati Rice 5kg',
      category: 'Rice & Grains',
      subcategory: 'Premium Rice',
      price: 500,
      costPrice: 400,
      stock: 50,
      minStock: 10,
      maxStock: 200,
      supplier: 'Rahman Traders',
      brand: 'Premium Foods',
      sku: 'RICE-BSM-5KG-001',
      barcode: '1234567890123',
      status: 'active',
      visibility: true,
      totalSold: 245,
      revenue: 122500,
      lastSoldAt: '2024-07-23',
      createdAt: '2024-01-15',
      updatedAt: '2024-07-20',
      image: '/api/placeholder/40/40',
      tags: ['premium', 'basmati', 'rice']
    },
    {
      id: 2,
      name: 'Onion 1kg',
      category: 'Vegetables',
      subcategory: 'Fresh Vegetables',
      price: 100,
      costPrice: 70,
      stock: 5,
      minStock: 10,
      maxStock: 100,
      supplier: 'Green Agro',
      brand: 'FreshVeg',
      sku: 'VEG-ONI-1KG-002',
      barcode: '1234567890124',
      status: 'active',
      visibility: true,
      totalSold: 180,
      revenue: 18000,
      lastSoldAt: '2024-07-22',
      createdAt: '2024-02-10',
      updatedAt: '2024-07-18',
      image: '/api/placeholder/40/40',
      tags: ['fresh', 'vegetable', 'onion']
    },
    {
      id: 3,
      name: 'Milk 1L',
      category: 'Dairy',
      subcategory: 'Fresh Milk',
      price: 100,
      costPrice: 85,
      stock: 30,
      minStock: 15,
      maxStock: 80,
      supplier: 'Milk Fresh Ltd.',
      brand: 'FreshMilk',
      sku: 'DAI-MIL-1L-003',
      barcode: '1234567890125',
      status: 'active',
      visibility: false,
      totalSold: 320,
      revenue: 32000,
      lastSoldAt: '2024-07-21',
      createdAt: '2024-01-20',
      updatedAt: '2024-07-15',
      image: '/api/placeholder/40/40',
      tags: ['fresh', 'dairy', 'milk']
    },
    {
      id: 4,
      name: 'Potato 1kg',
      category: 'Vegetables',
      subcategory: 'Root Vegetables',
      price: 80,
      costPrice: 55,
      stock: 0,
      minStock: 10,
      maxStock: 150,
      supplier: 'Potato House',
      brand: 'FarmFresh',
      sku: 'VEG-POT-1KG-004',
      barcode: '1234567890126',
      status: 'inactive',
      visibility: true,
      totalSold: 200,
      revenue: 16000,
      lastSoldAt: '2024-07-10',
      createdAt: '2024-03-05',
      updatedAt: '2024-07-12',
      image: '/api/placeholder/40/40',
      tags: ['fresh', 'vegetable', 'potato']
    },
    {
      id: 5,
      name: 'Cooking Oil 1L',
      category: 'Cooking',
      subcategory: 'Oils & Fats',
      price: 150,
      costPrice: 120,
      stock: 20,
      minStock: 8,
      maxStock: 60,
      supplier: 'Oil Mart',
      brand: 'PureCook',
      sku: 'COO-OIL-1L-005',
      barcode: '1234567890127',
      status: 'active',
      visibility: true,
      totalSold: 150,
      revenue: 22500,
      lastSoldAt: '2024-07-20',
      createdAt: '2024-02-15',
      updatedAt: '2024-07-19',
      image: '/api/placeholder/40/40',
      tags: ['cooking', 'oil', 'kitchen']
    },
    {
      id: 6,
      name: 'Sugar 1kg',
      category: 'Pantry',
      subcategory: 'Sweeteners',
      price: 120,
      costPrice: 95,
      stock: 8,
      minStock: 10,
      maxStock: 100,
      supplier: 'Sweet Supply Co.',
      brand: 'SweetLife',
      sku: 'PAN-SUG-1KG-006',
      barcode: '1234567890128',
      status: 'active',
      visibility: true,
      totalSold: 220,
      revenue: 26400,
      lastSoldAt: '2024-07-19',
      createdAt: '2024-01-25',
      updatedAt: '2024-07-17',
      image: '/api/placeholder/40/40',
      tags: ['sweet', 'sugar', 'pantry']
    }
  ];

  // Get unique categories for filter
  const categories = [...new Set(mockProducts.map(p => p.category))];

  // Filter and sort products
  const filteredAndSortedProducts = useMemo(() => {
    let filtered = mockProducts.filter(product => {
      const matchesSearch = product.name.toLowerCase().includes(search.toLowerCase()) ||
                           product.category.toLowerCase().includes(search.toLowerCase()) ||
                           product.supplier.toLowerCase().includes(search.toLowerCase()) ||
                           product.sku.toLowerCase().includes(search.toLowerCase()) ||
                           product.barcode.includes(search);

      const matchesCategory = categoryFilter === 'all' || product.category === categoryFilter;
      
      const matchesStock = stockFilter === 'all' ||
                          (stockFilter === 'in-stock' && product.stock > product.minStock) ||
                          (stockFilter === 'low-stock' && product.stock <= product.minStock && product.stock > 0) ||
                          (stockFilter === 'out-of-stock' && product.stock === 0);
      
      const matchesStatus = statusFilter === 'all' || product.status === statusFilter;

      return matchesSearch && matchesCategory && matchesStock && matchesStatus;
    });

    // Sort products
    filtered.sort((a, b) => {
      let aValue = a[sortBy];
      let bValue = b[sortBy];

      if (typeof aValue === 'string') {
        aValue = aValue.toLowerCase();
        bValue = bValue.toLowerCase();
      }

      if (sortOrder === 'asc') {
        return aValue < bValue ? -1 : aValue > bValue ? 1 : 0;
      } else {
        return aValue > bValue ? -1 : aValue < bValue ? 1 : 0;
      }
    });

    return filtered;
  }, [mockProducts, search, categoryFilter, stockFilter, statusFilter, sortBy, sortOrder]);

  // Pagination
  const totalPages = Math.ceil(filteredAndSortedProducts.length / itemsPerPage);
  const paginatedProducts = filteredAndSortedProducts.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  // Stats calculations
  const stats = useMemo(() => {
    const totalProducts = mockProducts.length;
    const activeProducts = mockProducts.filter(p => p.status === 'active').length;
    const lowStockProducts = mockProducts.filter(p => p.stock <= p.minStock && p.stock > 0).length;
    const outOfStockProducts = mockProducts.filter(p => p.stock === 0).length;
    const totalValue = mockProducts.reduce((sum, p) => sum + (p.stock * p.costPrice), 0);
    const totalRevenue = mockProducts.reduce((sum, p) => sum + p.revenue, 0);

    return {
      totalProducts,
      activeProducts,
      lowStockProducts,
      outOfStockProducts,
      totalValue,
      totalRevenue
    };
  }, [mockProducts]);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  const getStockStatus = (product) => {
    if (product.stock <= 0) return { status: 'out-of-stock', color: 'destructive', label: 'Out of Stock' };
    if (product.stock <= product.minStock) return { status: 'low-stock', color: 'warning', label: 'Low Stock' };
    return { status: 'in-stock', color: 'default', label: 'In Stock' };
  };

  const handleSelectProduct = (productId) => {
    setSelectedProducts(prev =>
      prev.includes(productId)
        ? prev.filter(id => id !== productId)
        : [...prev, productId]
    );
  };

  const handleSelectAll = () => {
    if (selectedProducts.length === paginatedProducts.length) {
      setSelectedProducts([]);
    } else {
      setSelectedProducts(paginatedProducts.map(p => p.id));
    }
  };

  const handleBulkAction = (action) => {
    console.log(`Bulk action: ${action} on products:`, selectedProducts);
    setSelectedProducts([]);
    setShowBulkActions(false);
  };

  const handleSort = (field) => {
    if (sortBy === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };

  const SortIcon = ({ field }) => {
    if (sortBy !== field) return <SortAsc className="h-4 w-4 opacity-50" />;
    return sortOrder === 'asc' ? <SortAsc className="h-4 w-4" /> : <SortDesc className="h-4 w-4" />;
  };

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto p-6 space-y-6">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold">Products</h1>
            <p className="text-muted-foreground mt-1">
              Manage your product catalog and inventory
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm">
              <Upload className="h-4 w-4 mr-2" />
              Import
            </Button>
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
            <Button onClick={() => navigate('/inventory/products/create')}>
              <Plus className="h-4 w-4 mr-2" />
              Add Product
            </Button>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Total Products</p>
                  <p className="text-2xl font-bold">{stats.totalProducts}</p>
                </div>
                <Package className="h-8 w-8 text-primary" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Active</p>
                  <p className="text-2xl font-bold text-green-600">{stats.activeProducts}</p>
                </div>
                <CheckCircle className="h-8 w-8 text-green-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Low Stock</p>
                  <p className="text-2xl font-bold text-yellow-600">{stats.lowStockProducts}</p>
                </div>
                <AlertTriangle className="h-8 w-8 text-yellow-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Out of Stock</p>
                  <p className="text-2xl font-bold text-red-600">{stats.outOfStockProducts}</p>
                </div>
                <XCircle className="h-8 w-8 text-red-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Inventory Value</p>
                  <p className="text-lg font-bold">{formatCurrency(stats.totalValue)}</p>
                </div>
                <DollarSign className="h-8 w-8 text-blue-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Total Revenue</p>
                  <p className="text-lg font-bold">{formatCurrency(stats.totalRevenue)}</p>
                </div>
                <BarChart3 className="h-8 w-8 text-purple-600" />
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Filters and Controls */}
        <Card>
          <CardContent className="p-6">
            <div className="flex flex-col lg:flex-row gap-4 mb-6">
              {/* Search */}
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  placeholder="Search products by name, SKU, category, or supplier..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-10"
                />
              </div>

              {/* Filters */}
              <div className="flex flex-wrap gap-2">
                <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                  <SelectTrigger className="w-40">
                    <SelectValue placeholder="Category" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Categories</SelectItem>
                    {categories.map(category => (
                      <SelectItem key={category} value={category}>{category}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                <Select value={stockFilter} onValueChange={setStockFilter}>
                  <SelectTrigger className="w-40">
                    <SelectValue placeholder="Stock Status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Stock</SelectItem>
                    <SelectItem value="in-stock">In Stock</SelectItem>
                    <SelectItem value="low-stock">Low Stock</SelectItem>
                    <SelectItem value="out-of-stock">Out of Stock</SelectItem>
                  </SelectContent>
                </Select>

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

              {/* View Toggle */}
              <div className="flex items-center gap-2">
                <Button
                  size="sm"
                  variant={view === 'table' ? 'default' : 'outline'}
                  onClick={() => setView('table')}
                >
                  <List className="h-4 w-4 mr-1" />
                  Table
                </Button>
                <Button
                  size="sm"
                  variant={view === 'grid' ? 'default' : 'outline'}
                  onClick={() => setView('grid')}
                >
                  <LayoutGrid className="h-4 w-4 mr-1" />
                  Grid
                </Button>
              </div>
            </div>

            {/* Bulk Actions */}
            {selectedProducts.length > 0 && (
              <div className="flex items-center justify-between mb-4 p-3 bg-primary/10 rounded-lg">
                <span className="text-sm font-medium">
                  {selectedProducts.length} product{selectedProducts.length > 1 ? 's' : ''} selected
                </span>
                <div className="flex items-center gap-2">
                  <Button size="sm" variant="outline" onClick={() => handleBulkAction('export')}>
                    <Download className="h-4 w-4 mr-1" />
                    Export
                  </Button>
                  <Button size="sm" variant="outline" onClick={() => handleBulkAction('archive')}>
                    <Archive className="h-4 w-4 mr-1" />
                    Archive
                  </Button>
                  <Button size="sm" variant="outline" onClick={() => setShowDeleteDialog(true)}>
                    <Trash2 className="h-4 w-4 mr-1" />
                    Delete
                  </Button>
                  <Button size="sm" variant="ghost" onClick={() => setSelectedProducts([])}>
                    <XCircle className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            )}

            {/* Table View */}
            {view === 'table' && (
              <div className="space-y-4">
                <div className="overflow-x-auto">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead className="w-12">
                          <Checkbox
                            checked={selectedProducts.length === paginatedProducts.length}
                            onCheckedChange={handleSelectAll}
                          />
                        </TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('name')}>
                          <div className="flex items-center">
                            Product <SortIcon field="name" />
                          </div>
                        </TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('category')}>
                          <div className="flex items-center">
                            Category <SortIcon field="category" />
                          </div>
                        </TableHead>
                        <TableHead>SKU</TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('price')}>
                          <div className="flex items-center">
                            Price <SortIcon field="price" />
                          </div>
                        </TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('stock')}>
                          <div className="flex items-center">
                            Stock <SortIcon field="stock" />
                          </div>
                        </TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('status')}>
                          <div className="flex items-center">
                            Status <SortIcon field="status" />
                          </div>
                        </TableHead>
                        <TableHead>Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {paginatedProducts.map((product) => {
                        const stockStatus = getStockStatus(product);
                        return (
                          <TableRow key={product.id} className="hover:bg-muted/50">
                            <TableCell>
                              <Checkbox
                                checked={selectedProducts.includes(product.id)}
                                onCheckedChange={() => handleSelectProduct(product.id)}
                              />
                            </TableCell>
                            <TableCell>
                              <div className="flex items-center space-x-3">
                                <Avatar className="h-10 w-10">
                                  <AvatarImage src={product.image} alt={product.name} />
                                  <AvatarFallback>
                                    <Package className="h-4 w-4" />
                                  </AvatarFallback>
                                </Avatar>
                                <div>
                                  <p
                                    className="font-medium text-primary cursor-pointer hover:underline"
                                    onClick={() => navigate(`/inventory/products/${product.id}`)}
                                  >
                                    {product.name}
                                  </p>
                                  <p className="text-xs text-muted-foreground">
                                    {product.brand} • {product.supplier}
                                  </p>
                                </div>
                              </div>
                            </TableCell>
                            <TableCell>
                              <Badge variant="secondary">{product.category}</Badge>
                            </TableCell>
                            <TableCell>
                              <span className="font-mono text-sm">{product.sku}</span>
                            </TableCell>
                            <TableCell>
                              <div>
                                <p className="font-medium">{formatCurrency(product.price)}</p>
                                <p className="text-xs text-muted-foreground">
                                  Cost: {formatCurrency(product.costPrice)}
                                </p>
                              </div>
                            </TableCell>
                            <TableCell>
                              <div className="space-y-1">
                                <Badge variant={stockStatus.color}>
                                  {product.stock}
                                </Badge>
                                <div className="w-16">
                                  <Progress 
                                    value={Math.min((product.stock / product.maxStock) * 100, 100)} 
                                    className="h-1"
                                  />
                                </div>
                              </div>
                            </TableCell>
                            <TableCell>
                              <div className="flex items-center space-x-2">
                                <Badge variant={product.status === 'active' ? 'default' : 'secondary'}>
                                  {product.status}
                                </Badge>
                                {!product.visibility && (
                                  <Badge variant="outline" className="text-xs">Hidden</Badge>
                                )}
                              </div>
                            </TableCell>
                            <TableCell>
                              <div className="flex items-center space-x-1">
                                <Button 
                                  size="sm" 
                                  variant="ghost"
                                  onClick={() => navigate(`/inventory/products/${product.id}`)}
                                >
                                  <Eye className="h-4 w-4" />
                                </Button>
                                <Button 
                                  size="sm" 
                                  variant="ghost"
                                  onClick={() => navigate(`/inventory/products/${product.id}?edit=true`)}
                                >
                                  <Edit className="h-4 w-4" />
                                </Button>
                                <DropdownMenu>
                                  <DropdownMenuTrigger asChild>
                                    <Button size="sm" variant="ghost">
                                      <MoreVertical className="h-4 w-4" />
                                    </Button>
                                  </DropdownMenuTrigger>
                                  <DropdownMenuContent align="end">
                                    <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                    <DropdownMenuItem>
                                      <Copy className="h-4 w-4 mr-2" />
                                      Duplicate
                                    </DropdownMenuItem>
                                    <DropdownMenuItem>
                                      <ScanLine className="h-4 w-4 mr-2" />
                                      Print Barcode
                                    </DropdownMenuItem>
                                    <DropdownMenuSeparator />
                                    <DropdownMenuItem>
                                      <Archive className="h-4 w-4 mr-2" />
                                      Archive
                                    </DropdownMenuItem>
                                    <DropdownMenuItem className="text-destructive">
                                      <Trash2 className="h-4 w-4 mr-2" />
                                      Delete
                                    </DropdownMenuItem>
                                  </DropdownMenuContent>
                                </DropdownMenu>
                              </div>
                            </TableCell>
                          </TableRow>
                        );
                      })}
                    </TableBody>
                  </Table>
                </div>

                {/* Pagination */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <span className="text-sm text-muted-foreground">Show</span>
                    <Select value={itemsPerPage.toString()} onValueChange={(value) => setItemsPerPage(parseInt(value))}>
                      <SelectTrigger className="w-20">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="5">5</SelectItem>
                        <SelectItem value="10">10</SelectItem>
                        <SelectItem value="20">20</SelectItem>
                        <SelectItem value="50">50</SelectItem>
                      </SelectContent>
                    </Select>
                    <span className="text-sm text-muted-foreground">
                      of {filteredAndSortedProducts.length} products
                    </span>
                  </div>

                  <div className="flex items-center space-x-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                      disabled={currentPage === 1}
                    >
                      Previous
                    </Button>
                    <span className="text-sm">
                      Page {currentPage} of {totalPages}
                    </span>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                      disabled={currentPage === totalPages}
                    >
                      Next
                    </Button>
                  </div>
                </div>
              </div>
            )}

            {/* Grid View */}
            {view === 'grid' && (
              <div className="space-y-4">
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
                  {paginatedProducts.map((product) => {
                    const stockStatus = getStockStatus(product);
                    return (
                      <Card 
                        key={product.id} 
                        className="hover:shadow-md transition-shadow relative"
                      >
                        <CardContent className="p-4">
                          {/* Selection Checkbox */}
                          <div className="absolute top-2 left-2">
                            <Checkbox
                              checked={selectedProducts.includes(product.id)}
                              onCheckedChange={() => handleSelectProduct(product.id)}
                            />
                          </div>

                          {/* Status Badges */}
                          <div className="absolute top-2 right-2 flex flex-col gap-1">
                            {product.status === 'inactive' && (
                              <Badge variant="secondary" className="text-xs">Inactive</Badge>
                            )}
                            {!product.visibility && (
                              <Badge variant="outline" className="text-xs">Hidden</Badge>
                            )}
                          </div>

                          <div className="flex flex-col items-center space-y-3 mt-6">
                            {/* Product Image */}
                            <Avatar className="h-16 w-16">
                              <AvatarImage src={product.image} alt={product.name} />
                              <AvatarFallback>
                                <Package className="h-6 w-6" />
                              </AvatarFallback>
                            </Avatar>

                            {/* Product Info */}
                            <div className="text-center space-y-2 w-full">
                              <h3 
                                className="font-semibold text-sm line-clamp-2 cursor-pointer hover:text-primary"
                                onClick={() => navigate(`/inventory/products/${product.id}`)}
                              >
                                {product.name}
                              </h3>
                              
                              <div className="space-y-1">
                                <Badge variant="secondary" className="text-xs">
                                  {product.category}
                                </Badge>
                                
                                <p className="text-xs text-muted-foreground">
                                  SKU: {product.sku}
                                </p>
                                
                                <p className="text-xs text-muted-foreground">
                                  {product.brand} • {product.supplier}
                                </p>
                              </div>

                              {/* Pricing */}
                              <div className="space-y-1">
                                <div className="flex justify-center items-center space-x-2">
                                  <span className="font-semibold text-primary">
                                    {formatCurrency(product.price)}
                                  </span>
                                  <span className="text-xs text-muted-foreground line-through">
                                    {formatCurrency(product.costPrice)}
                                  </span>
                                </div>
                                <p className="text-xs text-green-600">
                                  {(((product.price - product.costPrice) / product.price) * 100).toFixed(1)}% margin
                                </p>
                              </div>

                              {/* Stock Status */}
                              <div className="space-y-2">
                                <div className="flex justify-center items-center space-x-2">
                                  <Badge variant={stockStatus.color} className="text-xs">
                                    {product.stock} in stock
                                  </Badge>
                                </div>
                                
                                <div className="w-full">
                                  <Progress 
                                    value={Math.min((product.stock / product.maxStock) * 100, 100)} 
                                    className="h-1"
                                  />
                                  <div className="flex justify-between text-xs text-muted-foreground mt-1">
                                    <span>Min: {product.minStock}</span>
                                    <span>Max: {product.maxStock}</span>
                                  </div>
                                </div>
                              </div>

                              {/* Sales Info */}
                              <div className="grid grid-cols-2 gap-2 text-xs">
                                <div className="text-center">
                                  <p className="font-medium">{product.totalSold}</p>
                                  <p className="text-muted-foreground">Sold</p>
                                </div>
                                <div className="text-center">
                                  <p className="font-medium">{formatCurrency(product.revenue)}</p>
                                  <p className="text-muted-foreground">Revenue</p>
                                </div>
                              </div>
                            </div>

                            {/* Action Buttons */}
                            <div className="flex justify-center space-x-1 w-full">
                              <Button 
                                size="sm" 
                                variant="outline"
                                onClick={() => navigate(`/inventory/products/${product.id}`)}
                                className="flex-1"
                              >
                                <Eye className="h-3 w-3 mr-1" />
                                View
                              </Button>
                              <Button 
                                size="sm" 
                                variant="outline"
                                onClick={() => navigate(`/inventory/products/${product.id}?edit=true`)}
                              >
                                <Edit className="h-3 w-3" />
                              </Button>
                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button size="sm" variant="outline">
                                    <MoreVertical className="h-3 w-3" />
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                  <DropdownMenuItem>
                                    <Copy className="h-4 w-4 mr-2" />
                                    Duplicate
                                  </DropdownMenuItem>
                                  <DropdownMenuItem>
                                    <ScanLine className="h-4 w-4 mr-2" />
                                    Print Barcode
                                  </DropdownMenuItem>
                                  <DropdownMenuSeparator />
                                  <DropdownMenuItem>
                                    <Archive className="h-4 w-4 mr-2" />
                                    Archive
                                  </DropdownMenuItem>
                                  <DropdownMenuItem className="text-destructive">
                                    <Trash2 className="h-4 w-4 mr-2" />
                                    Delete
                                  </DropdownMenuItem>
                                </DropdownMenuContent>
                              </DropdownMenu>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    );
                  })}
                </div>

                {/* Grid Pagination */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <span className="text-sm text-muted-foreground">Show</span>
                    <Select value={itemsPerPage.toString()} onValueChange={(value) => setItemsPerPage(parseInt(value))}>
                      <SelectTrigger className="w-20">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="8">8</SelectItem>
                        <SelectItem value="16">16</SelectItem>
                        <SelectItem value="24">24</SelectItem>
                        <SelectItem value="32">32</SelectItem>
                      </SelectContent>
                    </Select>
                    <span className="text-sm text-muted-foreground">
                      of {filteredAndSortedProducts.length} products
                    </span>
                  </div>

                  <div className="flex items-center space-x-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                      disabled={currentPage === 1}
                    >
                      Previous
                    </Button>
                    <span className="text-sm">
                      Page {currentPage} of {totalPages}
                    </span>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                      disabled={currentPage === totalPages}
                    >
                      Next
                    </Button>
                  </div>
                </div>
              </div>
            )}

            {/* No Results */}
            {filteredAndSortedProducts.length === 0 && (
              <div className="text-center py-12">
                <Package className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-semibold mb-2">No products found</h3>
                <p className="text-muted-foreground mb-4">
                  {search || categoryFilter !== 'all' || stockFilter !== 'all' || statusFilter !== 'all'
                    ? 'Try adjusting your search or filters'
                    : 'Get started by adding your first product'
                  }
                </p>
                {!search && categoryFilter === 'all' && stockFilter === 'all' && statusFilter === 'all' && (
                  <Button onClick={() => navigate('/inventory/products/create')}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Your First Product
                  </Button>
                )}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Bulk Delete Dialog */}
        <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle className="flex items-center text-destructive">
                <AlertTriangle className="h-5 w-5 mr-2" />
                Delete Products
              </DialogTitle>
              <DialogDescription>
                Are you sure you want to delete {selectedProducts.length} product{selectedProducts.length > 1 ? 's' : ''}? 
                This action cannot be undone.
              </DialogDescription>
            </DialogHeader>
            
            <div className="bg-destructive/10 p-4 rounded-lg">
              <p className="text-sm text-destructive font-medium mb-2">This will permanently:</p>
              <ul className="text-sm text-destructive space-y-1">
                <li>• Remove the product{selectedProducts.length > 1 ? 's' : ''} from inventory</li>
                <li>• Delete all associated data</li>
                <li>• Remove from all sales reports</li>
              </ul>
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={() => setShowDeleteDialog(false)}>
                Cancel
              </Button>
              <Button 
                variant="destructive" 
                onClick={() => {
                  handleBulkAction('delete');
                  setShowDeleteDialog(false);
                }}
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Delete {selectedProducts.length} Product{selectedProducts.length > 1 ? 's' : ''}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
};

export default EnhancedProductsPage;