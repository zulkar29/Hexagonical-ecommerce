import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Search,
  MoreHorizontal,
  Eye,
  Edit,
  Trash2,
  Calendar,
  DollarSign,
  Store,
  Clock,
  CheckCircle,
  XCircle,
  CreditCard,
  Mail,
  Phone,
  MapPin,
  Package,
  Users,
  Star,
  Settings,
  BarChart3
} from 'lucide-react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

const TenantsPage = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedStatus, setSelectedStatus] = useState('all');
  const [selectedPlan, setSelectedPlan] = useState('all');

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
    
    return matchesSearch && matchesStatus && matchesPlan;
  });

  // Stats calculation
  const stats = {
    total: tenants.length,
    active: tenants.filter(t => t.status === 'active').length,
    trial: tenants.filter(t => t.status === 'trial').length,
    suspended: tenants.filter(t => t.status === 'suspended').length,
    totalRevenue: tenants.reduce((sum, t) => sum + t.totalRevenue, 0),
    avgRevenue: tenants.reduce((sum, t) => sum + t.monthlyRevenue, 0) / tenants.length
  };


  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center px-6">
          <div className="flex items-center gap-4">
            <Store className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Tenant Management</h1>
            <Badge variant="outline" className="ml-2">Admin Panel</Badge>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Stats Row */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4 w-full">
          <div className="flex items-center gap-3 bg-card p-3 rounded-lg border">
            <div className="p-2 rounded-md text-blue-600 bg-opacity-10 shrink-0">
              <Store className="h-5 w-5 text-blue-600" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Total Tenants</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{stats.total}</p>
                <span className="text-xs text-green-600 font-medium">+12%</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3 bg-card p-3 rounded-lg border">
            <div className="p-2 rounded-md text-green-600 bg-opacity-10 shrink-0">
              <CheckCircle className="h-5 w-5 text-green-600" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Active</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{stats.active}</p>
                <span className="text-xs text-green-600 font-medium">+8%</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3 bg-card p-3 rounded-lg border">
            <div className="p-2 rounded-md text-yellow-600 bg-opacity-10 shrink-0">
              <Clock className="h-5 w-5 text-yellow-600" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Trial</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{stats.trial}</p>
                <span className="text-xs text-yellow-600 font-medium">+15%</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3 bg-card p-3 rounded-lg border">
            <div className="p-2 rounded-md text-red-600 bg-opacity-10 shrink-0">
              <XCircle className="h-5 w-5 text-red-600" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Suspended</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{stats.suspended}</p>
                <span className="text-xs text-red-600 font-medium">-5%</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3 bg-card p-3 rounded-lg border">
            <div className="p-2 rounded-md text-emerald-600 bg-opacity-10 shrink-0">
              <DollarSign className="h-5 w-5 text-emerald-600" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Total Revenue</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{formatCurrency(stats.totalRevenue)}</p>
                <span className="text-xs text-emerald-600 font-medium">+18.2%</span>
              </div>
            </div>
          </div>
        </div>

        {/* Tenants Table */}
        <Card>
          <CardHeader className="py-4">
            {/* Filters and Search */}
            <div className="flex flex-wrap items-center gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  type="text"
                  placeholder="Search tenants by name, owner, email, or location..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 w-full"
                />
              </div>
              <Select value={selectedStatus} onValueChange={setSelectedStatus}>
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Filter by Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Statuses</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="trial">Trial</SelectItem>
                  <SelectItem value="suspended">Suspended</SelectItem>
                  <SelectItem value="expired">Expired</SelectItem>
                </SelectContent>
              </Select>
              <Select value={selectedPlan} onValueChange={setSelectedPlan}>
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Filter by Plan" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Plans</SelectItem>
                  <SelectItem value="Starter">Starter</SelectItem>
                  <SelectItem value="Business">Business</SelectItem>
                  <SelectItem value="Enterprise">Enterprise</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardHeader>
          <CardContent>
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Tenant Details</TableHead>
                    <TableHead>Owner Info</TableHead>
                    <TableHead>Plan & Status</TableHead>
                    <TableHead>Performance</TableHead>
                    <TableHead>Dates</TableHead>
                    <TableHead>Activity</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredTenants.map((tenant) => (
                    <TableRow key={tenant.id} className="hover:bg-muted/50">
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
                              <p className="font-medium text-sm">
                                <Link to={`/tenants/${tenant.id}`} className="text-primary hover:underline">
                                  {tenant.name}
                                </Link>
                              </p>
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
                            <DropdownMenuItem onClick={() => navigate(`/tenants/${tenant.id}`)}>
                              <Eye className="h-4 w-4 mr-2" />
                              View Details
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => navigate(`/tenants/${tenant.id}/manage`)}>
                              <Edit className="h-4 w-4 mr-2" />
                              Edit Tenant
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => navigate('/analytics/tenants')}>
                              <BarChart3 className="h-4 w-4 mr-2" />
                              View Analytics
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem>
                              <Mail className="h-4 w-4 mr-2" />
                              Send Message
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => navigate('/billing')}>
                              <CreditCard className="h-4 w-4 mr-2" />
                              Billing History
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => navigate(`/tenants/${tenant.id}/manage`)}>
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
                            <DropdownMenuItem className="text-destructive">
                              <Trash2 className="h-4 w-4 mr-2" />
                              Delete Account
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default TenantsPage;