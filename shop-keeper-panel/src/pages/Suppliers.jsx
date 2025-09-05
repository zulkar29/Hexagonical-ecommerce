import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Plus,
  Edit,
  Trash2,
  Building,
  Eye,
  LayoutGrid,
  List,
  Package,
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
  Phone,
  Mail,
  MapPin,
  Calendar,
  DollarSign,
  Users,
  Star,
  StarOff,
  Copy,
  Archive,
  RefreshCw,
  FileText,
  Printer,
  SortAsc,
  SortDesc,
  Clock,
  CreditCard,
  Truck
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

const EnhancedSuppliersPage = () => {
  const navigate = useNavigate();
  const [search, setSearch] = useState('');
  const [view, setView] = useState('table');
  const [selectedSuppliers, setSelectedSuppliers] = useState([]);
  const [sortBy, setSortBy] = useState('name');
  const [sortOrder, setSortOrder] = useState('asc');
  const [statusFilter, setStatusFilter] = useState('all');
  const [ratingFilter, setRatingFilter] = useState('all');
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const [currentPage, setCurrentPage] = useState(1);

  // Enhanced mock suppliers with comprehensive data
  const mockSuppliers = [
    {
      id: 1,
      name: 'Rahman Traders',
      email: 'contact@rahmantraders.com',
      phone: '01712345678',
      address: 'House 12, Road 5, Dhanmondi, Dhaka-1205',
      city: 'Dhaka',
      country: 'Bangladesh',
      productsCount: 12,
      activeProducts: 10,
      totalOrders: 45,
      totalValue: 250000,
      lastOrderDate: '2024-07-20',
      paymentTerms: '30 days',
      rating: 4.5,
      status: 'active',
      isPreferred: true,
      joinedDate: '2023-01-15',
      lastContact: '2024-07-18',
      image: '/api/placeholder/40/40',
      contactPerson: 'Abdul Rahman',
      designation: 'Sales Manager',
      businessType: 'Wholesale',
      taxId: 'BIN-123456789',
      bankDetails: 'Dutch Bangla Bank - 123456789',
      leadTime: '5-7 days',
      minimumOrder: 50000
    },
    {
      id: 2,
      name: 'Green Agro Ltd.',
      email: 'info@greenagro.com',
      phone: '01798765432',
      address: 'Plot 25, Gulshan Avenue, Gulshan-2, Dhaka-1212',
      city: 'Dhaka',
      country: 'Bangladesh',
      productsCount: 8,
      activeProducts: 8,
      totalOrders: 32,
      totalValue: 180000,
      lastOrderDate: '2024-07-22',
      paymentTerms: '15 days',
      rating: 4.2,
      status: 'active',
      isPreferred: false,
      joinedDate: '2023-03-20',
      lastContact: '2024-07-21',
      image: '/api/placeholder/40/40',
      contactPerson: 'Fatima Khan',
      designation: 'Business Development',
      businessType: 'Retail',
      taxId: 'BIN-987654321',
      bankDetails: 'BRAC Bank - 987654321',
      leadTime: '3-5 days',
      minimumOrder: 25000
    },
    {
      id: 3,
      name: 'Milk Fresh Ltd.',
      email: 'orders@milkfresh.bd',
      phone: '01656789012',
      address: 'Sector 10, Uttara Model Town, Dhaka-1230',
      city: 'Dhaka',
      country: 'Bangladesh',
      productsCount: 6,
      activeProducts: 5,
      totalOrders: 28,
      totalValue: 120000,
      lastOrderDate: '2024-07-15',
      paymentTerms: '7 days',
      rating: 3.8,
      status: 'active',
      isPreferred: false,
      joinedDate: '2023-06-10',
      lastContact: '2024-07-12',
      image: '/api/placeholder/40/40',
      contactPerson: 'Mohammad Ali',
      designation: 'Supply Manager',
      businessType: 'Manufacturer',
      taxId: 'BIN-456789123',
      bankDetails: 'Islami Bank - 456789123',
      leadTime: '2-4 days',
      minimumOrder: 15000
    },
    {
      id: 4,
      name: 'Potato House',
      email: 'sales@potatohouse.com',
      phone: '01534567890',
      address: 'Block D, Section 11, Mirpur, Dhaka-1216',
      city: 'Dhaka',
      country: 'Bangladesh',
      productsCount: 4,
      activeProducts: 2,
      totalOrders: 15,
      totalValue: 80000,
      lastOrderDate: '2024-06-25',
      paymentTerms: '45 days',
      rating: 3.2,
      status: 'inactive',
      isPreferred: false,
      joinedDate: '2023-08-05',
      lastContact: '2024-06-20',
      image: '/api/placeholder/40/40',
      contactPerson: 'Rashida Begum',
      designation: 'Owner',
      businessType: 'Wholesale',
      taxId: 'BIN-789123456',
      bankDetails: 'City Bank - 789123456',
      leadTime: '7-10 days',
      minimumOrder: 30000
    },
    {
      id: 5,
      name: 'Oil Mart International',
      email: 'procurement@oilmart.bd',
      phone: '01812345678',
      address: 'House 45, Road 12, Banani, Dhaka-1213',
      city: 'Dhaka',
      country: 'Bangladesh',
      productsCount: 15,
      activeProducts: 14,
      totalOrders: 60,
      totalValue: 400000,
      lastOrderDate: '2024-07-23',
      paymentTerms: '21 days',
      rating: 4.8,
      status: 'active',
      isPreferred: true,
      joinedDate: '2022-11-12',
      lastContact: '2024-07-22',
      image: '/api/placeholder/40/40',
      contactPerson: 'Ahmed Hassan',
      designation: 'Regional Manager',
      businessType: 'Distributor',
      taxId: 'BIN-321654987',
      bankDetails: 'Standard Chartered - 321654987',
      leadTime: '1-3 days',
      minimumOrder: 75000
    },
    {
      id: 6,
      name: 'Fresh Valley Suppliers',
      email: 'orders@freshvalley.net',
      phone: '01723456789',
      address: 'Warehouse 8, Savar Export Processing Zone, Dhaka',
      city: 'Savar',
      country: 'Bangladesh',
      productsCount: 20,
      activeProducts: 18,
      totalOrders: 38,
      totalValue: 290000,
      lastOrderDate: '2024-07-19',
      paymentTerms: '14 days',
      rating: 4.1,
      status: 'active',
      isPreferred: false,
      joinedDate: '2023-02-28',
      lastContact: '2024-07-16',
      image: '/api/placeholder/40/40',
      contactPerson: 'Nasir Ahmed',
      designation: 'Operations Head',
      businessType: 'Wholesale',
      taxId: 'BIN-147258369',
      bankDetails: 'Prime Bank - 147258369',
      leadTime: '4-6 days',
      minimumOrder: 40000
    }
  ];

  // Filter and sort suppliers
  const filteredAndSortedSuppliers = useMemo(() => {
    let filtered = mockSuppliers.filter(supplier => {
      const matchesSearch = supplier.name.toLowerCase().includes(search.toLowerCase()) ||
                           supplier.email.toLowerCase().includes(search.toLowerCase()) ||
                           supplier.phone.includes(search) ||
                           supplier.address.toLowerCase().includes(search.toLowerCase()) ||
                           supplier.contactPerson.toLowerCase().includes(search.toLowerCase());

      const matchesStatus = statusFilter === 'all' || supplier.status === statusFilter;
      
      const matchesRating = ratingFilter === 'all' ||
                          (ratingFilter === 'high' && supplier.rating >= 4.0) ||
                          (ratingFilter === 'medium' && supplier.rating >= 3.0 && supplier.rating < 4.0) ||
                          (ratingFilter === 'low' && supplier.rating < 3.0);

      return matchesSearch && matchesStatus && matchesRating;
    });

    // Sort suppliers
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
  }, [mockSuppliers, search, statusFilter, ratingFilter, sortBy, sortOrder]);

  // Pagination
  const totalPages = Math.ceil(filteredAndSortedSuppliers.length / itemsPerPage);
  const paginatedSuppliers = filteredAndSortedSuppliers.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  // Stats calculations
  const stats = useMemo(() => {
    const totalSuppliers = mockSuppliers.length;
    const activeSuppliers = mockSuppliers.filter(s => s.status === 'active').length;
    const preferredSuppliers = mockSuppliers.filter(s => s.isPreferred).length;
    const totalProducts = mockSuppliers.reduce((sum, s) => sum + s.productsCount, 0);
    const totalOrders = mockSuppliers.reduce((sum, s) => sum + s.totalOrders, 0);
    const totalValue = mockSuppliers.reduce((sum, s) => sum + s.totalValue, 0);
    const avgRating = mockSuppliers.reduce((sum, s) => sum + s.rating, 0) / totalSuppliers;

    return {
      totalSuppliers,
      activeSuppliers,
      preferredSuppliers,
      totalProducts,
      totalOrders,
      totalValue,
      avgRating
    };
  }, [mockSuppliers]);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

  const getRatingColor = (rating) => {
    if (rating >= 4.0) return 'text-green-600';
    if (rating >= 3.0) return 'text-yellow-600';
    return 'text-red-600';
  };

  const handleSelectSupplier = (supplierId) => {
    setSelectedSuppliers(prev =>
      prev.includes(supplierId)
        ? prev.filter(id => id !== supplierId)
        : [...prev, supplierId]
    );
  };

  const handleSelectAll = () => {
    if (selectedSuppliers.length === paginatedSuppliers.length) {
      setSelectedSuppliers([]);
    } else {
      setSelectedSuppliers(paginatedSuppliers.map(s => s.id));
    }
  };

  const handleBulkAction = (action) => {
    console.log(`Bulk action: ${action} on suppliers:`, selectedSuppliers);
    setSelectedSuppliers([]);
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

  const renderStars = (rating) => {
    return Array.from({ length: 5 }, (_, i) => (
      <Star
        key={i}
        className={cn(
          "h-3 w-3",
          i < Math.floor(rating) ? "fill-yellow-400 text-yellow-400" : "text-gray-300"
        )}
      />
    ));
  };

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto p-6 space-y-6">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold">Suppliers</h1>
            <p className="text-muted-foreground mt-1">
              Manage your supplier network and relationships
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
            <Button onClick={() => navigate('/inventory/suppliers/new')}>
              <Plus className="h-4 w-4 mr-2" />
              Add Supplier
            </Button>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-7 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Total Suppliers</p>
                  <p className="text-2xl font-bold">{stats.totalSuppliers}</p>
                </div>
                <Building className="h-8 w-8 text-primary" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Active</p>
                  <p className="text-2xl font-bold text-green-600">{stats.activeSuppliers}</p>
                </div>
                <CheckCircle className="h-8 w-8 text-green-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Preferred</p>
                  <p className="text-2xl font-bold text-yellow-600">{stats.preferredSuppliers}</p>
                </div>
                <Star className="h-8 w-8 text-yellow-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Products</p>
                  <p className="text-2xl font-bold">{stats.totalProducts}</p>
                </div>
                <Package className="h-8 w-8 text-blue-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Total Orders</p>
                  <p className="text-2xl font-bold">{stats.totalOrders}</p>
                </div>
                <Truck className="h-8 w-8 text-purple-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Total Value</p>
                  <p className="text-lg font-bold">{formatCurrency(stats.totalValue)}</p>
                </div>
                <DollarSign className="h-8 w-8 text-green-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Avg Rating</p>
                  <p className="text-2xl font-bold">{stats.avgRating.toFixed(1)}</p>
                </div>
                <BarChart3 className="h-8 w-8 text-orange-600" />
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
                  placeholder="Search suppliers by name, email, phone, or contact person..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-10"
                />
              </div>

              {/* Filters */}
              <div className="flex flex-wrap gap-2">
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

                <Select value={ratingFilter} onValueChange={setRatingFilter}>
                  <SelectTrigger className="w-32">
                    <SelectValue placeholder="Rating" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Ratings</SelectItem>
                    <SelectItem value="high">4+ Stars</SelectItem>
                    <SelectItem value="medium">3-4 Stars</SelectItem>
                    <SelectItem value="low">Below 3</SelectItem>
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
            {selectedSuppliers.length > 0 && (
              <div className="flex items-center justify-between mb-4 p-3 bg-primary/10 rounded-lg">
                <span className="text-sm font-medium">
                  {selectedSuppliers.length} supplier{selectedSuppliers.length > 1 ? 's' : ''} selected
                </span>
                <div className="flex items-center gap-2">
                  <Button size="sm" variant="outline" onClick={() => handleBulkAction('export')}>
                    <Download className="h-4 w-4 mr-1" />
                    Export
                  </Button>
                  <Button size="sm" variant="outline" onClick={() => handleBulkAction('email')}>
                    <Mail className="h-4 w-4 mr-1" />
                    Email
                  </Button>
                  <Button size="sm" variant="outline" onClick={() => setShowDeleteDialog(true)}>
                    <Trash2 className="h-4 w-4 mr-1" />
                    Delete
                  </Button>
                  <Button size="sm" variant="ghost" onClick={() => setSelectedSuppliers([])}>
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
                            checked={selectedSuppliers.length === paginatedSuppliers.length}
                            onCheckedChange={handleSelectAll}
                          />
                        </TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('name')}>
                          <div className="flex items-center">
                            Supplier <SortIcon field="name" />
                          </div>
                        </TableHead>
                        <TableHead>Contact Info</TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('productsCount')}>
                          <div className="flex items-center">
                            Products <SortIcon field="productsCount" />
                          </div>
                        </TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('totalValue')}>
                          <div className="flex items-center">
                            Total Value <SortIcon field="totalValue" />
                          </div>
                        </TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('rating')}>
                          <div className="flex items-center">
                            Rating <SortIcon field="rating" />
                          </div>
                        </TableHead>
                        <TableHead className="cursor-pointer" onClick={() => handleSort('lastOrderDate')}>
                          <div className="flex items-center">
                            Last Order <SortIcon field="lastOrderDate" />
                          </div>
                        </TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {paginatedSuppliers.map((supplier) => (
                        <TableRow key={supplier.id} className="hover:bg-muted/50">
                          <TableCell>
                            <Checkbox
                              checked={selectedSuppliers.includes(supplier.id)}
                              onCheckedChange={() => handleSelectSupplier(supplier.id)}
                            />
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center space-x-3">
                              <Avatar className="h-10 w-10">
                                <AvatarImage src={supplier.image} alt={supplier.name} />
                                <AvatarFallback>
                                  <Building className="h-4 w-4" />
                                </AvatarFallback>
                              </Avatar>
                              <div>
                                <p
                                  className="font-medium text-primary cursor-pointer hover:underline"
                                  onClick={() => navigate(`/inventory/suppliers/${supplier.id}`)}
                                >
                                  {supplier.name}
                                  {supplier.isPreferred && (
                                    <Star className="h-3 w-3 inline ml-1 fill-yellow-400 text-yellow-400" />
                                  )}
                                </p>
                                <p className="text-xs text-muted-foreground">
                                  {supplier.contactPerson} • {supplier.businessType}
                                </p>
                              </div>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="space-y-1">
                              <div className="flex items-center space-x-1">
                                <Phone className="h-3 w-3 text-muted-foreground" />
                                <span className="text-sm">{supplier.phone}</span>
                              </div>
                              <div className="flex items-center space-x-1">
                                <Mail className="h-3 w-3 text-muted-foreground" />
                                <span className="text-sm">{supplier.email}</span>
                              </div>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="text-center">
                              <Badge variant="outline">
                                {supplier.productsCount}
                              </Badge>
                              <p className="text-xs text-muted-foreground mt-1">
                                {supplier.activeProducts} active
                              </p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div>
                              <p className="font-medium">{formatCurrency(supplier.totalValue)}</p>
                              <p className="text-xs text-muted-foreground">
                                {supplier.totalOrders} orders
                              </p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center space-x-1">
                              <div className="flex">
                                {renderStars(supplier.rating)}
                              </div>
                              <span className={cn("text-sm font-medium", getRatingColor(supplier.rating))}>
                                {supplier.rating}
                              </span>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div>
                              <p className="text-sm">{formatDate(supplier.lastOrderDate)}</p>
                              <p className="text-xs text-muted-foreground">
                                {supplier.paymentTerms}
                              </p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <Badge variant={supplier.status === 'active' ? 'default' : 'secondary'}>
                              {supplier.status}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center space-x-1">
                              <Button 
                                size="sm" 
                                variant="ghost"
                                onClick={() => navigate(`/inventory/suppliers/${supplier.id}`)}
                              >
                                <Eye className="h-4 w-4" />
                              </Button>
                              <Button size="sm" variant="ghost">
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
                                    <Mail className="h-4 w-4 mr-2" />
                                    Send Email
                                  </DropdownMenuItem>
                                  <DropdownMenuItem>
                                    <FileText className="h-4 w-4 mr-2" />
                                    Generate Report
                                  </DropdownMenuItem>
                                  <DropdownMenuSeparator />
                                  <DropdownMenuItem>
                                    {supplier.isPreferred ? (
                                      <>
                                        <StarOff className="h-4 w-4 mr-2" />
                                        Remove Preferred
                                      </>
                                    ) : (
                                      <>
                                        <Star className="h-4 w-4 mr-2" />
                                        Mark Preferred
                                      </>
                                    )}
                                  </DropdownMenuItem>
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
                      ))}
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
                      of {filteredAndSortedSuppliers.length} suppliers
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
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                  {paginatedSuppliers.map((supplier) => (
                    <Card 
                      key={supplier.id} 
                      className="hover:shadow-md transition-shadow relative"
                    >
                      <CardContent className="p-6">
                        {/* Selection Checkbox */}
                        <div className="absolute top-3 left-3">
                          <Checkbox
                            checked={selectedSuppliers.includes(supplier.id)}
                            onCheckedChange={() => handleSelectSupplier(supplier.id)}
                          />
                        </div>

                        {/* Status Badges */}
                        <div className="absolute top-3 right-3 flex flex-col gap-1">
                          {supplier.isPreferred && (
                            <Badge variant="secondary" className="text-xs">
                              <Star className="h-3 w-3 mr-1 fill-yellow-400 text-yellow-400" />
                              Preferred
                            </Badge>
                          )}
                          {supplier.status === 'inactive' && (
                            <Badge variant="outline" className="text-xs">Inactive</Badge>
                          )}
                        </div>

                        <div className="flex flex-col items-center space-y-4 mt-6">
                          {/* Supplier Avatar */}
                          <Avatar className="h-16 w-16">
                            <AvatarImage src={supplier.image} alt={supplier.name} />
                            <AvatarFallback>
                              <Building className="h-6 w-6" />
                            </AvatarFallback>
                          </Avatar>

                          {/* Supplier Info */}
                          <div className="text-center space-y-3 w-full">
                            <div>
                              <h3 
                                className="font-semibold text-lg cursor-pointer hover:text-primary"
                                onClick={() => navigate(`/inventory/suppliers/${supplier.id}`)}
                              >
                                {supplier.name}
                              </h3>
                              <p className="text-sm text-muted-foreground">
                                {supplier.contactPerson}
                              </p>
                              <p className="text-xs text-muted-foreground">
                                {supplier.designation}
                              </p>
                            </div>

                            {/* Contact Info */}
                            <div className="space-y-2">
                              <div className="flex items-center justify-center space-x-1">
                                <Phone className="h-3 w-3 text-muted-foreground" />
                                <span className="text-sm">{supplier.phone}</span>
                              </div>
                              <div className="flex items-center justify-center space-x-1">
                                <MapPin className="h-3 w-3 text-muted-foreground" />
                                <span className="text-sm">{supplier.city}</span>
                              </div>
                            </div>

                            {/* Rating */}
                            <div className="flex items-center justify-center space-x-2">
                              <div className="flex">
                                {renderStars(supplier.rating)}
                              </div>
                              <span className={cn("text-sm font-medium", getRatingColor(supplier.rating))}>
                                {supplier.rating}
                              </span>
                            </div>

                            {/* Business Metrics */}
                            <div className="grid grid-cols-2 gap-4 text-sm">
                              <div className="text-center">
                                <p className="font-semibold">{supplier.productsCount}</p>
                                <p className="text-muted-foreground">Products</p>
                              </div>
                              <div className="text-center">
                                <p className="font-semibold">{supplier.totalOrders}</p>
                                <p className="text-muted-foreground">Orders</p>
                              </div>
                            </div>

                            {/* Financial Info */}
                            <div className="space-y-2">
                              <div className="text-center">
                                <p className="font-semibold text-primary text-lg">
                                  {formatCurrency(supplier.totalValue)}
                                </p>
                                <p className="text-xs text-muted-foreground">Total Business</p>
                              </div>
                              
                              <div className="flex items-center justify-between text-xs text-muted-foreground">
                                <span>Payment Terms:</span>
                                <span className="font-medium">{supplier.paymentTerms}</span>
                              </div>
                              
                              <div className="flex items-center justify-between text-xs text-muted-foreground">
                                <span>Lead Time:</span>
                                <span className="font-medium">{supplier.leadTime}</span>
                              </div>
                            </div>

                            {/* Last Activity */}
                            <div className="space-y-1">
                              <div className="flex items-center justify-center space-x-1">
                                <Calendar className="h-3 w-3 text-muted-foreground" />
                                <span className="text-xs text-muted-foreground">
                                  Last order: {formatDate(supplier.lastOrderDate)}
                                </span>
                              </div>
                              <div className="flex items-center justify-center space-x-1">
                                <Clock className="h-3 w-3 text-muted-foreground" />
                                <span className="text-xs text-muted-foreground">
                                  Contact: {formatDate(supplier.lastContact)}
                                </span>
                              </div>
                            </div>

                            {/* Business Type & Status */}
                            <div className="flex items-center justify-center space-x-2">
                              <Badge variant="outline" className="text-xs">
                                {supplier.businessType}
                              </Badge>
                              <Badge variant={supplier.status === 'active' ? 'default' : 'secondary'} className="text-xs">
                                {supplier.status}
                              </Badge>
                            </div>
                          </div>

                          {/* Action Buttons */}
                          <div className="flex justify-center space-x-2 w-full">
                            <Button 
                              size="sm" 
                              variant="outline"
                              onClick={() => navigate(`/inventory/suppliers/${supplier.id}`)}
                              className="flex-1"
                            >
                              <Eye className="h-3 w-3 mr-1" />
                              View
                            </Button>
                            <Button size="sm" variant="outline">
                              <Edit className="h-3 w-3" />
                            </Button>
                            <Button size="sm" variant="outline">
                              <Mail className="h-3 w-3" />
                            </Button>
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button size="sm" variant="outline">
                                  <MoreVertical className="h-3 w-3" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end">
                                <DropdownMenuItem>
                                  <Phone className="h-4 w-4 mr-2" />
                                  Call Supplier
                                </DropdownMenuItem>
                                <DropdownMenuItem>
                                  <FileText className="h-4 w-4 mr-2" />
                                  Generate Report
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem>
                                  {supplier.isPreferred ? (
                                    <>
                                      <StarOff className="h-4 w-4 mr-2" />
                                      Remove Preferred
                                    </>
                                  ) : (
                                    <>
                                      <Star className="h-4 w-4 mr-2" />
                                      Mark Preferred
                                    </>
                                  )}
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
                  ))}
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
                        <SelectItem value="6">6</SelectItem>
                        <SelectItem value="12">12</SelectItem>
                        <SelectItem value="18">18</SelectItem>
                        <SelectItem value="24">24</SelectItem>
                      </SelectContent>
                    </Select>
                    <span className="text-sm text-muted-foreground">
                      of {filteredAndSortedSuppliers.length} suppliers
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
            {filteredAndSortedSuppliers.length === 0 && (
              <div className="text-center py-12">
                <Building className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-semibold mb-2">No suppliers found</h3>
                <p className="text-muted-foreground mb-4">
                  {search || statusFilter !== 'all' || ratingFilter !== 'all'
                    ? 'Try adjusting your search or filters'
                    : 'Get started by adding your first supplier'
                  }
                </p>
                {!search && statusFilter === 'all' && ratingFilter === 'all' && (
                  <Button onClick={() => navigate('/inventory/suppliers/new')}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Your First Supplier
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
                Delete Suppliers
              </DialogTitle>
              <DialogDescription>
                Are you sure you want to delete {selectedSuppliers.length} supplier{selectedSuppliers.length > 1 ? 's' : ''}? 
                This action cannot be undone.
              </DialogDescription>
            </DialogHeader>
            
            <div className="bg-destructive/10 p-4 rounded-lg">
              <p className="text-sm text-destructive font-medium mb-2">This will permanently:</p>
              <ul className="text-sm text-destructive space-y-1">
                <li>• Remove the supplier{selectedSuppliers.length > 1 ? 's' : ''} from your database</li>
                <li>• Delete all contact information and history</li>
                <li>• Remove supplier associations from products</li>
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
                Delete {selectedSuppliers.length} Supplier{selectedSuppliers.length > 1 ? 's' : ''}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
};

export default EnhancedSuppliersPage;