import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import {
  Search,
  Filter,
  Plus,
  Download,
  Upload,
  MoreHorizontal,
  Eye,
  Edit,
  Trash2,
  Mail,
  Phone,
  MapPin,
  Package,
  Calendar,
  DollarSign,
  Users,
  Store,
  Activity,
  Clock,
  CheckCircle,
  AlertTriangle,
  XCircle,
  TrendingUp,
  TrendingDown,
  RefreshCw,
  Settings,
  BarChart3,
  CreditCard,
  Shield,
  UserPlus,
  Building,
  Globe,
  Star,
  AlertCircle,
  Info,
  ChevronLeft,
  ChevronRight,
  ArrowUpDown,
  SortAsc,
  SortDesc
} from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import { Checkbox } from '@/components/ui/checkbox';

const TenantsPage = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedStatus, setSelectedStatus] = useState('all');
  const [selectedPlan, setSelectedPlan] = useState('all');
  const [sortBy, setSortBy] = useState('name');
  const [sortOrder, setSortOrder] = useState('asc');
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const [selectedTenants, setSelectedTenants] = useState([]);
  const [showAddDialog, setShowAddDialog] = useState(false);
  const [showBulkActionsDialog, setShowBulkActionsDialog] = useState(false);
  const [activeTab, setActiveTab] = useState('all');

  // Mock tenant data
  const [tenants, setTenants] = useState([
    {
      id: 1,
      name: 'Rahman Electronics',
      owner: 'Abdul Rahman',
      email: 'rahman@email.com',
      phone: '01712345678',
      plan: 'Business',
      status: 'active',
      joinDate: '2024-01-15',
      lastLogin: '2024-07-24T10:30:00',
      monthlyRevenue: 5000,
      totalRevenue: 35000,
      location: 'Dhanmondi, Dhaka',
      productsCount: 245,
      ordersCount: 1250,
      customersCount: 89,
      address: 'House 15, Road 7, Dhanmondi, Dhaka 1205',
      businessType: 'Electronics',
      website: 'https://rahmanelectronics.com',
      rating: 4.8,
      subscriptionEnd: '2024-08-15',
      notes: 'Premium customer with excellent payment history'
    },
    {
      id: 2,
      name: 'Fatima Fashion',
      owner: 'Fatima Khatun',
      email: 'fatima@email.com',
      phone: '01798765432',
      plan: 'Starter',
      status: 'trial',
      joinDate: '2024-07-10',
      lastLogin: '2024-07-23T15:45:00',
      monthlyRevenue: 2000,
      totalRevenue: 2000,
      location: 'Gulshan, Dhaka',
      productsCount: 89,
      ordersCount: 145,
      customersCount: 34,
      address: 'Shop 23, Gulshan Avenue, Dhaka 1212',
      businessType: 'Fashion & Apparel',
      website: 'https://fatimafashion.com',
      rating: 4.2,
      subscriptionEnd: '2024-08-10',
      notes: 'New trial customer, showing good engagement'
    },
    {
      id: 3,
      name: 'Modern Pharmacy',
      owner: 'Dr. Ahmed Ali',
      email: 'ahmed@modernpharmacy.com',
      phone: '01656789012',
      plan: 'Enterprise',
      status: 'active',
      joinDate: '2023-11-20',
      lastLogin: '2024-07-24T09:15:00',
      monthlyRevenue: 10000,
      totalRevenue: 80000,
      location: 'Uttara, Dhaka',
      productsCount: 567,
      ordersCount: 2340,
      customersCount: 156,
      address: 'Sector 7, Uttara, Dhaka 1230',
      businessType: 'Healthcare & Pharmacy',
      website: 'https://modernpharmacy.com',
      rating: 4.9,
      subscriptionEnd: '2024-11-20',
      notes: 'High-value enterprise client with multiple locations'
    },
    {
      id: 4,
      name: 'Green Grocers',
      owner: 'Mohammad Hasan',
      email: 'hasan@greengrocers.com',
      phone: '01534567890',
      plan: 'Business',
      status: 'suspended',
      joinDate: '2024-02-28',
      lastLogin: '2024-07-20T14:20:00',
      monthlyRevenue: 5000,
      totalRevenue: 25000,
      location: 'Mirpur, Dhaka',
      productsCount: 156,
      ordersCount: 680,
      customersCount: 67,
      address: 'Block C, Mirpur 10, Dhaka 1216',
      businessType: 'Grocery & Food',
      website: null,
      rating: 3.8,
      subscriptionEnd: '2024-08-28',
      notes: 'Suspended due to payment issues, attempting recovery'
    },
    {
      id: 5,
      name: 'Tech Solutions',
      owner: 'Rashida Begum',
      email: 'rashida@techsolutions.com',
      phone: '01445678901',
      plan: 'Enterprise',
      status: 'active',
      joinDate: '2023-09-15',
      lastLogin: '2024-07-24T11:00:00',
      monthlyRevenue: 10000,
      totalRevenue: 100000,
      location: 'Banani, Dhaka',
      productsCount: 423,
      ordersCount: 1890,
      customersCount: 234,
      address: 'Road 11, Banani, Dhaka 1213',
      businessType: 'Technology & Services',
      website: 'https://techsolutions.bd',
      rating: 4.7,
      subscriptionEnd: '2024-09-15',
      notes: 'Long-term enterprise client with API integration'
    },
    {
      id: 6,
      name: 'Dhaka Books',
      owner: 'Nasir Ahmed',
      email: 'nasir@dhakabooks.com',
      phone: '01356789012',
      plan: 'Starter',
      status: 'active',
      joinDate: '2024-06-01',
      lastLogin: '2024-07-22T16:30:00',
      monthlyRevenue: 2000,
      totalRevenue: 4000,
      location: 'New Market, Dhaka',
      productsCount: 345,
      ordersCount: 234,
      customersCount: 78,
      address: 'Shop 45, New Market, Dhaka 1205',
      businessType: 'Books & Education',
      website: 'https://dhakabooks.com',
      rating: 4.3,
      subscriptionEnd: '2024-08-01',
      notes: 'Bookstore with good customer base'
    }
  ]);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

  const formatLastLogin = (dateString) => {
    const now = new Date();
    const loginDate = new Date(dateString);
    const diffInHours = Math.floor((now - loginDate) / (1000 * 60 * 60));
    
    if (diffInHours < 1) return 'Online now';
    if (diffInHours < 24) return `${diffInHours}h ago`;
    const diffInDays = Math.floor(diffInHours / 24);
    return `${diffInDays}d ago`;
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'active': return 'bg-green-100 text-green-700 border-green-200';
      case 'trial': return 'bg-blue-100 text-blue-700 border-blue-200';
      case 'suspended': return 'bg-red-100 text-red-700 border-red-200';
      case 'expired': return 'bg-gray-100 text-gray-700 border-gray-200';
      default: return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  const getPlanColor = (plan) => {
    switch (plan) {
      case 'Starter': return 'bg-purple-100 text-purple-700 border-purple-200';
      case 'Business': return 'bg-blue-100 text-blue-700 border-blue-200';
      case 'Enterprise': return 'bg-orange-100 text-orange-700 border-orange-200';
      default: return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  const getRatingStars = (rating) => {
    return Array.from({ length: 5 }, (_, i) => (
      <Star
        key={i}
        className={`h-3 w-3 ${
          i < Math.floor(rating) 
            ? 'text-yellow-400 fill-current' 
            : 'text-gray-300'
        }`}
      />
    ));
  };

  // Filter and sort tenants
  const filteredTenants = tenants.filter(tenant => {
    const matchesSearch = 
      tenant.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      tenant.owner.toLowerCase().includes(searchQuery.toLowerCase()) ||
      tenant.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
      tenant.location.toLowerCase().includes(searchQuery.toLowerCase());
    
    const matchesStatus = selectedStatus === 'all' || tenant.status === selectedStatus;
    const matchesPlan = selectedPlan === 'all' || tenant.plan === selectedPlan;
    const matchesTab = activeTab === 'all' || tenant.status === activeTab;
    
    return matchesSearch && matchesStatus && matchesPlan && matchesTab;
  }).sort((a, b) => {
    let aValue = a[sortBy];
    let bValue = b[sortBy];
    
    if (sortBy === 'joinDate' || sortBy === 'lastLogin') {
      aValue = new Date(aValue);
      bValue = new Date(bValue);
    }
    
    if (typeof aValue === 'string') {
      aValue = aValue.toLowerCase();
      bValue = bValue.toLowerCase();
    }
    
    if (sortOrder === 'asc') {
      return aValue > bValue ? 1 : -1;
    } else {
      return aValue < bValue ? 1 : -1;
    }
  });

  // Pagination
  const totalPages = Math.ceil(filteredTenants.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const paginatedTenants = filteredTenants.slice(startIndex, startIndex + itemsPerPage);

  // Stats calculation
  const stats = {
    total: tenants.length,
    active: tenants.filter(t => t.status === 'active').length,
    trial: tenants.filter(t => t.status === 'trial').length,
    suspended: tenants.filter(t => t.status === 'suspended').length,
    totalRevenue: tenants.reduce((sum, t) => sum + t.totalRevenue, 0),
    avgRevenue: tenants.reduce((sum, t) => sum + t.monthlyRevenue, 0) / tenants.length
  };

  const handleSort = (column) => {
    if (sortBy === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(column);
      setSortOrder('asc');
    }
  };

  const handleSelectTenant = (tenantId) => {
    setSelectedTenants(prev => 
      prev.includes(tenantId) 
        ? prev.filter(id => id !== tenantId)
        : [...prev, tenantId]
    );
  };

  const handleSelectAll = () => {
    if (selectedTenants.length === paginatedTenants.length) {
      setSelectedTenants([]);
    } else {
      setSelectedTenants(paginatedTenants.map(t => t.id));
    }
  };

  const handleBulkAction = (action) => {
    console.log(`Bulk action: ${action} on tenants:`, selectedTenants);
    // Implement bulk actions here
    setSelectedTenants([]);
    setShowBulkActionsDialog(false);
  };

  const handleDeleteTenant = (tenantId) => {
    setTenants(prev => prev.filter(t => t.id !== tenantId));
  };

  const SortButton = ({ column, children }) => (
    <Button
      variant="ghost"
      size="sm"
      className="h-8 p-0 font-medium"
      onClick={() => handleSort(column)}
    >
      {children}
      {sortBy === column ? (
        sortOrder === 'asc' ? (
          <SortAsc className="ml-1 h-3 w-3" />
        ) : (
          <SortDesc className="ml-1 h-3 w-3" />
        )
      ) : (
        <ArrowUpDown className="ml-1 h-3 w-3" />
      )}
    </Button>
  );

  return (
    <div className="space-y-6 p-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Tenant Management</h1>
          <p className="text-muted-foreground">
            Manage shop registrations, subscriptions, and tenant accounts
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="outline" size="sm">
            <Upload className="h-4 w-4 mr-2" />
            Import
          </Button>
          <Button variant="outline" size="sm">
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Button onClick={() => setShowAddDialog(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Add Tenant
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Tenants</p>
                <p className="text-2xl font-bold">{stats.total}</p>
                <div className="flex items-center space-x-1 mt-1">
                  <TrendingUp className="h-3 w-3 text-green-500" />
                  <span className="text-xs text-green-500">+12% from last month</span>
                </div>
              </div>
              <Store className="h-8 w-8 text-primary" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active Tenants</p>
                <p className="text-2xl font-bold">{stats.active}</p>
                <div className="flex items-center space-x-1 mt-1">
                  <CheckCircle className="h-3 w-3 text-green-500" />
                  <span className="text-xs text-muted-foreground">{Math.round((stats.active / stats.total) * 100)}% active</span>
                </div>
              </div>
              <Activity className="h-8 w-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Revenue</p>
                <p className="text-2xl font-bold">{formatCurrency(stats.totalRevenue)}</p>
                <div className="flex items-center space-x-1 mt-1">
                  <DollarSign className="h-3 w-3 text-green-500" />
                  <span className="text-xs text-green-500">+18% from last month</span>
                </div>
              </div>
              <TrendingUp className="h-8 w-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Avg Revenue</p>
                <p className="text-2xl font-bold">{formatCurrency(Math.round(stats.avgRevenue))}</p>
                <div className="flex items-center space-x-1 mt-1">
                  <BarChart3 className="h-3 w-3 text-blue-500" />
                  <span className="text-xs text-muted-foreground">per tenant/month</span>
                </div>
              </div>
              <BarChart3 className="h-8 w-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters and Search */}
      <Card>
        <CardContent className="p-6">
          <div className="flex flex-col lg:flex-row gap-4">
            {/* Search */}
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  placeholder="Search tenants by name, owner, email, or location..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>

            {/* Filters */}
            <div className="flex flex-wrap gap-2">
              <Select value={selectedStatus} onValueChange={setSelectedStatus}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="trial">Trial</SelectItem>
                  <SelectItem value="suspended">Suspended</SelectItem>
                  <SelectItem value="expired">Expired</SelectItem>
                </SelectContent>
              </Select>

              <Select value={selectedPlan} onValueChange={setSelectedPlan}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Plan" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Plans</SelectItem>
                  <SelectItem value="Starter">Starter</SelectItem>
                  <SelectItem value="Business">Business</SelectItem>
                  <SelectItem value="Enterprise">Enterprise</SelectItem>
                </SelectContent>
              </Select>

              <Select value={itemsPerPage.toString()} onValueChange={(value) => setItemsPerPage(parseInt(value))}>
                <SelectTrigger className="w-20">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="10">10</SelectItem>
                  <SelectItem value="25">25</SelectItem>
                  <SelectItem value="50">50</SelectItem>
                  <SelectItem value="100">100</SelectItem>
                </SelectContent>
              </Select>

              <Button variant="outline" size="sm">
                <Filter className="h-4 w-4 mr-2" />
                More Filters
              </Button>

              <Button variant="outline" size="sm">
                <RefreshCw className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="all">All ({stats.total})</TabsTrigger>
          <TabsTrigger value="active">Active ({stats.active})</TabsTrigger>
          <TabsTrigger value="trial">Trial ({stats.trial})</TabsTrigger>
          <TabsTrigger value="suspended">Suspended ({stats.suspended})</TabsTrigger>
          <TabsTrigger value="expired">Expired (0)</TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="space-y-4">
          {/* Bulk Actions */}
          {selectedTenants.length > 0 && (
            <Card className="border-primary/20 bg-primary/5">
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Checkbox 
                      checked={selectedTenants.length === paginatedTenants.length}
                      onCheckedChange={handleSelectAll}
                    />
                    <span className="text-sm font-medium">
                      {selectedTenants.length} tenant{selectedTenants.length !== 1 ? 's' : ''} selected
                    </span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Button variant="outline" size="sm" onClick={() => setShowBulkActionsDialog(true)}>
                      Bulk Actions
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => setSelectedTenants([])}>
                      Clear Selection
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Tenants Table */}
          <Card>
            <CardContent className="p-0">
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-12">
                        <Checkbox
                          checked={selectedTenants.length === paginatedTenants.length && paginatedTenants.length > 0}
                          onCheckedChange={handleSelectAll}
                        />
                      </TableHead>
                      <TableHead>
                        <SortButton column="name">Tenant Details</SortButton>
                      </TableHead>
                      <TableHead>
                        <SortButton column="owner">Owner Info</SortButton>
                      </TableHead>
                      <TableHead>
                        <SortButton column="plan">Plan & Status</SortButton>
                      </TableHead>
                      <TableHead>
                        <SortButton column="monthlyRevenue">Performance</SortButton>
                      </TableHead>
                      <TableHead>
                        <SortButton column="joinDate">Dates</SortButton>
                      </TableHead>
                      <TableHead>
                        <SortButton column="lastLogin">Activity</SortButton>
                      </TableHead>
                      <TableHead className="w-12">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {paginatedTenants.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={8} className="text-center py-8">
                          <div className="flex flex-col items-center space-y-2">
                            <Store className="h-8 w-8 text-muted-foreground" />
                            <p className="text-muted-foreground">No tenants found</p>
                            <Button variant="outline" size="sm" onClick={() => setSearchQuery('')}>
                              Clear filters
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ) : (
                      paginatedTenants.map((tenant) => (
                        <TableRow key={tenant.id} className="hover:bg-muted/50">
                          <TableCell>
                            <Checkbox
                              checked={selectedTenants.includes(tenant.id)}
                              onCheckedChange={() => handleSelectTenant(tenant.id)}
                            />
                          </TableCell>
                          
                          {/* Tenant Details */}
                          <TableCell>
                            <div className="space-y-1">
                              <div className="flex items-center space-x-2">
                                <Avatar className="h-8 w-8">
                                  <AvatarFallback className="text-xs">
                                    {tenant.name.split(' ').map(n => n[0]).join('')}
                                  </AvatarFallback>
                                </Avatar>
                                <div>
                                  <p className="font-medium text-sm"><Link to={`/tenants/${tenant.id}`} className="text-primary underline hover:opacity-80">{tenant.name}</Link></p>
                                  <p className="text-xs text-muted-foreground">{tenant.businessType}</p>
                                </div>
                              </div>
                              <div className="flex items-center space-x-2 text-xs text-muted-foreground">
                                <MapPin className="h-3 w-3" />
                                <span>{tenant.location}</span>
                              </div>
                              <div className="flex items-center space-x-1">
                                {getRatingStars(tenant.rating)}
                                <span className="text-xs text-muted-foreground ml-1">{tenant.rating}</span>
                              </div>
                            </div>
                          </TableCell>

                          {/* Owner Info */}
                          <TableCell>
                            <div className="space-y-1">
                              <p className="font-medium text-sm">{tenant.owner}</p>
                              <div className="flex items-center space-x-1 text-xs text-muted-foreground">
                                <Mail className="h-3 w-3" />
                                <span>{tenant.email}</span>
                              </div>
                              <div className="flex items-center space-x-1 text-xs text-muted-foreground">
                                <Phone className="h-3 w-3" />
                                <span>{tenant.phone}</span>
                              </div>
                            </div>
                          </TableCell>

                          {/* Plan & Status */}
                          <TableCell>
                            <div className="space-y-2">
                              <Badge className={getPlanColor(tenant.plan) + ' text-xs'}>
                                {tenant.plan}
                              </Badge>
                              <Badge className={getStatusColor(tenant.status) + ' text-xs block w-fit'}>
                                {tenant.status}
                              </Badge>
                            </div>
                          </TableCell>

                          {/* Performance */}
                          <TableCell>
                            <div className="space-y-1">
                              <p className="font-semibold text-sm">
                                {formatCurrency(tenant.monthlyRevenue)}<span className="text-xs text-muted-foreground">/mo</span>
                              </p>
                              <p className="text-xs text-muted-foreground">
                                Total: {formatCurrency(tenant.totalRevenue)}
                              </p>
                              <div className="flex items-center space-x-3 text-xs text-muted-foreground">
                                <div className="flex items-center space-x-1">
                                  <Package className="h-3 w-3" />
                                  <span>{tenant.productsCount}</span>
                                </div>
                                <div className="flex items-center space-x-1">
                                  <Users className="h-3 w-3" />
                                  <span>{tenant.customersCount}</span>
                                </div>
                              </div>
                            </div>
                          </TableCell>

                          {/* Dates */}
                          <TableCell>
                            <div className="space-y-1">
                              <div className="text-xs">
                                <span className="text-muted-foreground">Joined: </span>
                                <span className="font-medium">{formatDate(tenant.joinDate)}</span>
                              </div>
                              <div className="text-xs">
                                <span className="text-muted-foreground">Expires: </span>
                                <span className="font-medium">{formatDate(tenant.subscriptionEnd)}</span>
                              </div>
                            </div>
                          </TableCell>

                          {/* Activity */}
                          <TableCell>
                            <div className="space-y-1">
                              <div className="flex items-center space-x-2">
                                <div className={`w-2 h-2 rounded-full ${
                                  formatLastLogin(tenant.lastLogin) === 'Online now' ? 'bg-green-500' :
                                  formatLastLogin(tenant.lastLogin).includes('h ago') ? 'bg-yellow-500' : 'bg-gray-400'
                                }`} />
                                <span className="text-xs">{formatLastLogin(tenant.lastLogin)}</span>
                              </div>
                              <p className="text-xs text-muted-foreground">
                                {tenant.ordersCount} orders total
                              </p>
                            </div>
                          </TableCell>

                          {/* Actions */}
                          <TableCell>
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                                  <MoreHorizontal className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end" className="w-48">
                                <DropdownMenuItem onClick={() => {}}>
                                  <Eye className="h-4 w-4 mr-2" />
                                  View Details
                                </DropdownMenuItem>
                                <DropdownMenuItem>
                                  <Edit className="h-4 w-4 mr-2" />
                                  Edit Tenant
                                </DropdownMenuItem>
                                <DropdownMenuItem>
                                  <BarChart3 className="h-4 w-4 mr-2" />
                                  View Analytics
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem>
                                  <Mail className="h-4 w-4 mr-2" />
                                  Send Message
                                </DropdownMenuItem>
                                <DropdownMenuItem>
                                  <CreditCard className="h-4 w-4 mr-2" />
                                  Billing History
                                </DropdownMenuItem>
                                <DropdownMenuItem>
                                  <Settings className="h-4 w-4 mr-2" />
                                  Account Settings
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                {tenant.status === 'active' && (
                                  <DropdownMenuItem className="text-orange-600">
                                    <Clock className="h-4 w-4 mr-2" />
                                    Suspend Account
                                  </DropdownMenuItem>
                                )}
                                {tenant.status === 'suspended' && (
                                  <DropdownMenuItem className="text-green-600">
                                    <CheckCircle className="h-4 w-4 mr-2" />
                                    Reactivate Account
                                  </DropdownMenuItem>
                                )}
                                <AlertDialog>
                                  <AlertDialogTrigger asChild>
                                    <DropdownMenuItem 
                                      className="text-red-600"
                                      onSelect={(e) => e.preventDefault()}
                                    >
                                      <Trash2 className="h-4 w-4 mr-2" />
                                      Delete Account
                                    </DropdownMenuItem>
                                  </AlertDialogTrigger>
                                  <AlertDialogContent>
                                    <AlertDialogHeader>
                                      <AlertDialogTitle>Delete Tenant Account</AlertDialogTitle>
                                      <AlertDialogDescription>
                                        Are you sure you want to delete {tenant.name}? This action cannot be undone and will permanently remove all tenant data.
                                      </AlertDialogDescription>
                                    </AlertDialogHeader>
                                    <AlertDialogFooter>
                                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                                      <AlertDialogAction 
                                        className="bg-red-600 hover:bg-red-700"
                                        onClick={() => handleDeleteTenant(tenant.id)}
                                      >
                                        Delete Account
                                      </AlertDialogAction>
                                    </AlertDialogFooter>
                                  </AlertDialogContent>
                                </AlertDialog>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>

          {/* Pagination */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <p className="text-sm text-muted-foreground">
                Showing {startIndex + 1} to {Math.min(startIndex + itemsPerPage, filteredTenants.length)} of {filteredTenants.length} tenants
              </p>
            </div>
            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                disabled={currentPage === 1}
              >
                <ChevronLeft className="h-4 w-4" />
                Previous
              </Button>
              <div className="flex items-center space-x-1">
                {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                  const page = i + 1;
                  return (
                    <Button
                      key={page}
                      variant={currentPage === page ? "default" : "outline"}
                      size="sm"
                      className="w-8"
                      onClick={() => setCurrentPage(page)}
                    >
                      {page}
                    </Button>
                  );
                })}
                {totalPages > 5 && (
                  <>
                    <span className="text-muted-foreground">...</span>
                    <Button
                      variant={currentPage === totalPages ? "default" : "outline"}
                      size="sm"
                      className="w-8"
                      onClick={() => setCurrentPage(totalPages)}
                    >
                      {totalPages}
                    </Button>
                  </>
                )}
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                disabled={currentPage === totalPages}
              >
                Next
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </TabsContent>
      </Tabs>

      {/* Add Tenant Dialog */}
      <Dialog open={showAddDialog} onOpenChange={setShowAddDialog}>
        <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Add New Tenant</DialogTitle>
            <DialogDescription>
              Register a new shop to the platform with their business details
            </DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="shopName">Shop Name *</Label>
              <Input id="shopName" placeholder="Enter shop name" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="businessType">Business Type</Label>
              <Select>
                <SelectTrigger>
                  <SelectValue placeholder="Select business type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="electronics">Electronics</SelectItem>
                  <SelectItem value="fashion">Fashion & Apparel</SelectItem>
                  <SelectItem value="grocery">Grocery & Food</SelectItem>
                  <SelectItem value="pharmacy">Healthcare & Pharmacy</SelectItem>
                  <SelectItem value="books">Books & Education</SelectItem>
                  <SelectItem value="technology">Technology & Services</SelectItem>
                  <SelectItem value="other">Other</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="ownerName">Owner Name *</Label>
              <Input id="ownerName" placeholder="Enter owner name" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email Address *</Label>
              <Input id="email" type="email" placeholder="Enter email address" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="phone">Phone Number *</Label>
              <Input id="phone" placeholder="Enter phone number" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="website">Website (Optional)</Label>
              <Input id="website" placeholder="https://example.com" />
            </div>
            <div className="md:col-span-2 space-y-2">
              <Label htmlFor="address">Address *</Label>
              <Input id="address" placeholder="Enter full business address" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="plan">Subscription Plan *</Label>
              <Select>
                <SelectTrigger>
                  <SelectValue placeholder="Select subscription plan" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="starter">Starter - ৳2,000/month</SelectItem>
                  <SelectItem value="business">Business - ৳5,000/month</SelectItem>
                  <SelectItem value="enterprise">Enterprise - ৳10,000/month</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="status">Initial Status</Label>
              <Select defaultValue="trial">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="trial">Trial (14 days)</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="md:col-span-2 space-y-2">
              <Label htmlFor="notes">Notes (Optional)</Label>
              <Textarea 
                id="notes" 
                placeholder="Additional notes about the tenant or business" 
                rows={3}
              />
            </div>
            <div className="md:col-span-2 flex items-center space-x-2">
              <Checkbox id="sendWelcome" defaultChecked />
              <Label htmlFor="sendWelcome" className="text-sm">
                Send welcome email with account setup instructions
              </Label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowAddDialog(false)}>
              Cancel
            </Button>
            <Button onClick={() => setShowAddDialog(false)}>
              <UserPlus className="h-4 w-4 mr-2" />
              Create Tenant
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Bulk Actions Dialog */}
      <Dialog open={showBulkActionsDialog} onOpenChange={setShowBulkActionsDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Bulk Actions</DialogTitle>
            <DialogDescription>
              Apply actions to {selectedTenants.length} selected tenant{selectedTenants.length !== 1 ? 's' : ''}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-2">
              <Button 
                variant="outline" 
                onClick={() => handleBulkAction('suspend')}
                className="justify-start"
              >
                <Clock className="h-4 w-4 mr-2" />
                Suspend Accounts
              </Button>
              <Button 
                variant="outline" 
                onClick={() => handleBulkAction('activate')}
                className="justify-start"
              >
                <CheckCircle className="h-4 w-4 mr-2" />
                Activate Accounts
              </Button>
              <Button 
                variant="outline" 
                onClick={() => handleBulkAction('message')}
                className="justify-start"
              >
                <Mail className="h-4 w-4 mr-2" />
                Send Message
              </Button>
              <Button 
                variant="outline" 
                onClick={() => handleBulkAction('export')}
                className="justify-start"
              >
                <Download className="h-4 w-4 mr-2" />
                Export Data
              </Button>
            </div>
            <Separator />
            <Button 
              variant="destructive" 
              onClick={() => handleBulkAction('delete')}
              className="w-full justify-start"
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Delete Selected Tenants
            </Button>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowBulkActionsDialog(false)}>
              Cancel
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default TenantsPage;