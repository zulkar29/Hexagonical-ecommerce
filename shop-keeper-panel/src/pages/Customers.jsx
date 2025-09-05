import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Search,
  Plus,
  Filter,
  Download,
  Upload,
  Users,
  Phone,
  MapPin,
  Calendar,
  ShoppingBag,
  TrendingUp,
  Edit,
  Trash2,
  Eye,
  Mail,
  Star,
  Gift,
  CreditCard,
  Clock,
  User,
  Building,
  MoreHorizontal,
  X,
  Check,
  AlertCircle,
  FileText,
  MessageSquare
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
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
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Progress } from '@/components/ui/progress';

const CustomerPage = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCustomer, setSelectedCustomer] = useState(null);
  const [showAddDialog, setShowAddDialog] = useState(false);
  const [showDetailsDialog, setShowDetailsDialog] = useState(false);
  const [filterType, setFilterType] = useState('all');
  const [sortBy, setSortBy] = useState('name');
  const [viewMode, setViewMode] = useState('table'); // 'table' or 'cards'

  // Mock customer data
  const customers = [
    {
      id: 1,
      name: 'Fatima Khan',
      phone: '01712345678',
      email: 'fatima.khan@gmail.com',
      address: 'House 45, Road 12, Dhanmondi, Dhaka-1205',
      area: 'Dhanmondi',
      joinDate: '2024-01-15',
      lastPurchase: '2025-01-20',
      totalPurchases: 45780,
      totalOrders: 67,
      averageOrderValue: 683,
      loyaltyPoints: 458,
      status: 'VIP',
      category: 'Regular',
      notes: 'Prefers organic products. Usually shops on weekends.',
      birthday: '1985-03-15',
      recentOrders: [
        { id: '#1234', date: '2025-01-20', amount: 2450, status: 'completed' },
        { id: '#1189', date: '2025-01-15', amount: 1890, status: 'completed' },
        { id: '#1156', date: '2025-01-08', amount: 3200, status: 'completed' }
      ]
    },
    {
      id: 2,
      name: 'Ahmed Ali',
      phone: '01798765432',
      email: 'ahmed.ali@yahoo.com',
      address: 'Flat 8B, Gulshan Avenue, Gulshan-1, Dhaka-1212',
      area: 'Gulshan',
      joinDate: '2024-03-22',
      lastPurchase: '2025-01-18',
      totalPurchases: 32450,
      totalOrders: 45,
      averageOrderValue: 721,
      loyaltyPoints: 324,
      status: 'Gold',
      category: 'Regular',
      notes: 'Bulk buyer. Owns a small restaurant.',
      birthday: '1978-09-28',
      recentOrders: [
        { id: '#1223', date: '2025-01-18', amount: 5670, status: 'completed' },
        { id: '#1198', date: '2025-01-12', amount: 3450, status: 'completed' }
      ]
    },
    {
      id: 3,
      name: 'Rashida Begum',
      phone: '01656789012',
      email: 'rashida.begum@hotmail.com',
      address: 'Block C, House 23, Uttara Sector 7, Dhaka-1230',
      area: 'Uttara',
      joinDate: '2023-11-08',
      lastPurchase: '2025-01-21',
      totalPurchases: 67890,
      totalOrders: 89,
      averageOrderValue: 763,
      loyaltyPoints: 679,
      status: 'VIP',
      category: 'Premium',
      notes: 'Long-time customer. Very loyal and refers others.',
      birthday: '1970-12-05',
      recentOrders: [
        { id: '#1245', date: '2025-01-21', amount: 1890, status: 'completed' },
        { id: '#1201', date: '2025-01-14', amount: 2340, status: 'completed' }
      ]
    },
    {
      id: 4,
      name: 'Mohammad Hasan',
      phone: '01534567890',
      email: 'hasan.mohammad@gmail.com',
      address: 'Road 5, Block B, Mirpur-10, Dhaka-1216',
      area: 'Mirpur',
      joinDate: '2024-06-10',
      lastPurchase: '2025-01-19',
      totalPurchases: 18670,
      totalOrders: 34,
      averageOrderValue: 549,
      loyaltyPoints: 187,
      status: 'Silver',
      category: 'Regular',
      notes: 'Young family. Price-conscious buyer.',
      birthday: '1990-07-18',
      recentOrders: [
        { id: '#1235', date: '2025-01-19', amount: 1230, status: 'completed' }
      ]
    },
    {
      id: 5,
      name: 'Nasir Ahmed',
      phone: '01623456789',
      email: 'nasir.ahmed@outlook.com',
      address: 'House 78, New Eskaton Road, Ramna, Dhaka-1000',
      area: 'Ramna',
      joinDate: '2024-09-15',
      lastPurchase: '2025-01-16',
      totalPurchases: 24500,
      totalOrders: 28,
      averageOrderValue: 875,
      loyaltyPoints: 245,
      status: 'Gold',
      category: 'Regular',
      notes: 'Prefers premium brands. Good payment history.',
      birthday: '1982-04-22',
      recentOrders: [
        { id: '#1212', date: '2025-01-16', amount: 2780, status: 'completed' }
      ]
    },
    {
      id: 6,
      name: 'Salma Khatun',
      phone: '01745632189',
      email: 'salma.khatun@gmail.com',
      address: 'Lane 3, Wari, Old Dhaka-1203',
      area: 'Wari',
      joinDate: '2024-12-01',
      lastPurchase: '2025-01-17',
      totalPurchases: 8450,
      totalOrders: 12,
      averageOrderValue: 704,
      loyaltyPoints: 85,
      status: 'Bronze',
      category: 'New',
      notes: 'New customer. Showing good potential.',
      birthday: '1995-11-30',
      recentOrders: [
        { id: '#1209', date: '2025-01-17', amount: 890, status: 'completed' }
      ]
    }
  ];

  // Filter and sort customers
  const filteredCustomers = useMemo(() => {
    let filtered = customers.filter(customer => {
      const matchesSearch = customer.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                           customer.phone.includes(searchQuery) ||
                           customer.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
                           customer.area.toLowerCase().includes(searchQuery.toLowerCase());
      
      const matchesFilter = filterType === 'all' || 
                           (filterType === 'vip' && customer.status === 'VIP') ||
                           (filterType === 'new' && customer.category === 'New') ||
                           (filterType === 'inactive' && 
                            new Date() - new Date(customer.lastPurchase) > 30 * 24 * 60 * 60 * 1000);
      
      return matchesSearch && matchesFilter;
    });

    // Sort customers
    filtered.sort((a, b) => {
      switch (sortBy) {
        case 'name':
          return a.name.localeCompare(b.name);
        case 'totalPurchases':
          return b.totalPurchases - a.totalPurchases;
        case 'lastPurchase':
          return new Date(b.lastPurchase) - new Date(a.lastPurchase);
        case 'joinDate':
          return new Date(b.joinDate) - new Date(a.joinDate);
        default:
          return 0;
      }
    });

    return filtered;
  }, [customers, searchQuery, filterType, sortBy]);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

  const getStatusColor = (status) => {
    switch (status) {
      case 'VIP': return 'bg-purple-100 text-purple-700';
      case 'Gold': return 'bg-yellow-100 text-yellow-700';
      case 'Silver': return 'bg-gray-100 text-gray-700';
      case 'Bronze': return 'bg-orange-100 text-orange-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'VIP': return <Star className="h-3 w-3" />;
      case 'Gold': return <Gift className="h-3 w-3" />;
      case 'Silver': return <CreditCard className="h-3 w-3" />;
      case 'Bronze': return <User className="h-3 w-3" />;
      default: return <User className="h-3 w-3" />;
    }
  };

  // Stats calculations
  const stats = {
    total: customers.length,
    new: customers.filter(c => c.category === 'New').length,
    vip: customers.filter(c => c.status === 'VIP').length,
    totalRevenue: customers.reduce((sum, c) => sum + c.totalPurchases, 0),
    averageValue: customers.reduce((sum, c) => sum + c.averageOrderValue, 0) / customers.length
  };

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Customer Management</h1>
          <p className="text-muted-foreground mt-1">
            Manage your customer relationships and track purchase history
          </p>
        </div>
        <div className="flex space-x-3">
          <Button variant="outline">
            <Upload className="h-4 w-4 mr-2" />
            Import
          </Button>
          <Button variant="outline">
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Dialog open={showAddDialog} onOpenChange={setShowAddDialog}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="h-4 w-4 mr-2" />
                Add Customer
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-md">
              <DialogHeader>
                <DialogTitle>Add New Customer</DialogTitle>
                <DialogDescription>
                  Create a new customer profile for your business
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="firstName">First Name</Label>
                    <Input id="firstName" placeholder="Mohammad" />
                  </div>
                  <div>
                    <Label htmlFor="lastName">Last Name</Label>
                    <Input id="lastName" placeholder="Rahman" />
                  </div>
                </div>
                <div>
                  <Label htmlFor="phone">Phone Number</Label>
                  <Input id="phone" placeholder="01712345678" />
                </div>
                <div>
                  <Label htmlFor="email">Email (Optional)</Label>
                  <Input id="email" type="email" placeholder="customer@example.com" />
                </div>
                <div>
                  <Label htmlFor="address">Address</Label>
                  <Textarea id="address" placeholder="House/Flat, Road, Area, Dhaka" />
                </div>
                <div>
                  <Label htmlFor="notes">Notes (Optional)</Label>
                  <Textarea id="notes" placeholder="Any special notes about this customer..." />
                </div>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setShowAddDialog(false)}>
                  Cancel
                </Button>
                <Button onClick={() => setShowAddDialog(false)}>
                  Add Customer
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-5 gap-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Customers</p>
                <p className="text-2xl font-bold">{stats.total}</p>
              </div>
              <Users className="h-8 w-8 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">New This Month</p>
                <p className="text-2xl font-bold">{stats.new}</p>
              </div>
              <TrendingUp className="h-8 w-8 text-green-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">VIP Customers</p>
                <p className="text-2xl font-bold">{stats.vip}</p>
              </div>
              <Star className="h-8 w-8 text-purple-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Revenue</p>
                <p className="text-2xl font-bold">{formatCurrency(stats.totalRevenue)}</p>
              </div>
              <ShoppingBag className="h-8 w-8 text-orange-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Avg Order Value</p>
                <p className="text-2xl font-bold">{formatCurrency(stats.averageValue)}</p>
              </div>
              <CreditCard className="h-8 w-8 text-indigo-600" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters and Search */}
      <Card>
        <CardContent className="p-6">
          <div className="flex flex-col sm:flex-row gap-4 items-center justify-between">
            <div className="flex flex-1 items-center space-x-4">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  placeholder="Search customers..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Select value={filterType} onValueChange={setFilterType}>
                <SelectTrigger className="w-40">
                  <Filter className="h-4 w-4 mr-2" />
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Customers</SelectItem>
                  <SelectItem value="vip">VIP Only</SelectItem>
                  <SelectItem value="new">New Customers</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                </SelectContent>
              </Select>
              <Select value={sortBy} onValueChange={setSortBy}>
                <SelectTrigger className="w-40">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="name">Sort by Name</SelectItem>
                  <SelectItem value="totalPurchases">Sort by Revenue</SelectItem>
                  <SelectItem value="lastPurchase">Sort by Last Purchase</SelectItem>
                  <SelectItem value="joinDate">Sort by Join Date</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex items-center space-x-2">
              <Button
                variant={viewMode === 'table' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setViewMode('table')}
              >
                Table
              </Button>
              <Button
                variant={viewMode === 'cards' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setViewMode('cards')}
              >
                Cards
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Customer List */}
      {viewMode === 'table' ? (
        <Card>
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Customer</TableHead>
                  <TableHead>Contact</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Total Purchases</TableHead>
                  <TableHead>Orders</TableHead>
                  <TableHead>Last Purchase</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredCustomers.map((customer) => (
                  <TableRow key={customer.id}>
                    <TableCell>
                      <div className="flex items-center space-x-3">
                        <Avatar className="h-10 w-10">
                          <AvatarFallback>
                            {customer.name.split(' ').map(n => n[0]).join('')}
                          </AvatarFallback>
                        </Avatar>
                        <div>
                          <p className="font-medium">{customer.name}</p>
                          <p className="text-sm text-muted-foreground">{customer.area}</p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div>
                        <p className="text-sm">{customer.phone}</p>
                        <p className="text-sm text-muted-foreground">{customer.email}</p>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge className={cn("text-xs", getStatusColor(customer.status))}>
                        {getStatusIcon(customer.status)}
                        <span className="ml-1">{customer.status}</span>
                      </Badge>
                    </TableCell>
                    <TableCell className="font-medium">
                      {formatCurrency(customer.totalPurchases)}
                    </TableCell>
                    <TableCell>{customer.totalOrders}</TableCell>
                    <TableCell className="text-sm">
                      {formatDate(customer.lastPurchase)}
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem
                            onClick={() => navigate(`/customer-details/${customer.id}`)}
                          >
                            <Eye className="h-4 w-4 mr-2" />
                            View Details
                          </DropdownMenuItem>
                          <DropdownMenuItem>
                            <Edit className="h-4 w-4 mr-2" />
                            Edit Customer
                          </DropdownMenuItem>
                          <DropdownMenuItem>
                            <MessageSquare className="h-4 w-4 mr-2" />
                            Send Message
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem className="text-destructive">
                            <Trash2 className="h-4 w-4 mr-2" />
                            Delete Customer
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredCustomers.map((customer) => (
            <Card key={customer.id} className="hover:shadow-md transition-shadow">
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <Avatar className="h-12 w-12">
                      <AvatarFallback>
                        {customer.name.split(' ').map(n => n[0]).join('')}
                      </AvatarFallback>
                    </Avatar>
                    <div>
                      <CardTitle className="text-lg">{customer.name}</CardTitle>
                      <CardDescription>{customer.area}</CardDescription>
                    </div>
                  </div>
                  <Badge className={cn("text-xs", getStatusColor(customer.status))}>
                    {getStatusIcon(customer.status)}
                    <span className="ml-1">{customer.status}</span>
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <p className="text-muted-foreground">Total Purchases</p>
                    <p className="font-semibold">{formatCurrency(customer.totalPurchases)}</p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">Orders</p>
                    <p className="font-semibold">{customer.totalOrders}</p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">Avg Order</p>
                    <p className="font-semibold">{formatCurrency(customer.averageOrderValue)}</p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">Loyalty Points</p>
                    <p className="font-semibold">{customer.loyaltyPoints}</p>
                  </div>
                </div>
                <div className="space-y-2">
                  <div className="flex items-center text-sm text-muted-foreground">
                    <Phone className="h-4 w-4 mr-2" />
                    {customer.phone}
                  </div>
                  <div className="flex items-center text-sm text-muted-foreground">
                    <Clock className="h-4 w-4 mr-2" />
                    Last purchase: {formatDate(customer.lastPurchase)}
                  </div>
                </div>
                <div className="flex space-x-2">
                  <Button
                    size="sm"
                    variant="outline"
                    className="flex-1"
                    onClick={() => navigate(`/customer-details/${customer.id}`)}
                  >
                    <Eye className="h-4 w-4 mr-2" />
                    View
                  </Button>
                  <Button size="sm" variant="outline" className="flex-1">
                    <Edit className="h-4 w-4 mr-2" />
                    Edit
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Customer Details Dialog */}
      <Dialog open={showDetailsDialog} onOpenChange={setShowDetailsDialog}>
        <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
          {selectedCustomer && (
            <>
              <DialogHeader>
                <div className="flex items-center space-x-4">
                  <Avatar className="h-16 w-16">
                    <AvatarFallback className="text-lg">
                      {selectedCustomer.name.split(' ').map(n => n[0]).join('')}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <DialogTitle className="text-2xl">{selectedCustomer.name}</DialogTitle>
                    <DialogDescription>
                      Customer since {formatDate(selectedCustomer.joinDate)}
                    </DialogDescription>
                    <Badge className={cn("mt-2", getStatusColor(selectedCustomer.status))}>
                      {getStatusIcon(selectedCustomer.status)}
                      <span className="ml-1">{selectedCustomer.status} Customer</span>
                    </Badge>
                  </div>
                </div>
              </DialogHeader>

              <div className="space-y-6">
                {/* Quick Stats */}
                <div className="grid grid-cols-3 gap-4">
                  <div className="text-center p-4 bg-muted rounded-lg">
                    <p className="text-2xl font-bold text-primary">
                      {formatCurrency(selectedCustomer.totalPurchases)}
                    </p>
                    <p className="text-sm text-muted-foreground">Total Purchases</p>
                  </div>
                  <div className="text-center p-4 bg-muted rounded-lg">
                    <p className="text-2xl font-bold text-primary">{selectedCustomer.totalOrders}</p>
                    <p className="text-sm text-muted-foreground">Total Orders</p>
                  </div>
                  <div className="text-center p-4 bg-muted rounded-lg">
                    <p className="text-2xl font-bold text-primary">{selectedCustomer.loyaltyPoints}</p>
                    <p className="text-sm text-muted-foreground">Loyalty Points</p>
                  </div>
                </div>

                {/* Contact Information */}
                <div>
                  <h3 className="font-semibold mb-3">Contact Information</h3>
                  <div className="space-y-3">
                    <div className="flex items-center space-x-3">
                      <Phone className="h-4 w-4 text-muted-foreground" />
                      <span>{selectedCustomer.phone}</span>
                    </div>
                    <div className="flex items-center space-x-3">
                      <Mail className="h-4 w-4 text-muted-foreground" />
                      <span>{selectedCustomer.email}</span>
                    </div>
                    <div className="flex items-start space-x-3">
                      <MapPin className="h-4 w-4 text-muted-foreground mt-1" />
                      <span>{selectedCustomer.address}</span>
                    </div>
                    <div className="flex items-center space-x-3">
                      <Calendar className="h-4 w-4 text-muted-foreground" />
                      <span>Birthday: {formatDate(selectedCustomer.birthday)}</span>
                    </div>
                  </div>
                </div>

                {/* Purchase History */}
                <div>
                  <h3 className="font-semibold mb-3">Recent Orders</h3>
                  <div className="space-y-2">
                    {selectedCustomer.recentOrders.map((order) => (
                      <div key={order.id} className="flex items-center justify-between p-3 border rounded-lg">
                        <div>
                          <p className="font-medium">{order.id}</p>
                          <p className="text-sm text-muted-foreground">{formatDate(order.date)}</p>
                        </div>
                        <div className="text-right">
                          <p className="font-medium">{formatCurrency(order.amount)}</p>
                          <Badge variant="secondary" className="text-xs">
                            {order.status}
                          </Badge>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Customer Notes */}
                <div>
                  <h3 className="font-semibold mb-3">Notes</h3>
                  <div className="p-4 bg-muted rounded-lg">
                    <p className="text-sm">{selectedCustomer.notes}</p>
                  </div>
                </div>

                {/* Loyalty Progress */}
                <div>
                  <h3 className="font-semibold mb-3">Loyalty Program</h3>
                  <div className="space-y-3">
                    <div className="flex justify-between text-sm">
                      <span>Current Points: {selectedCustomer.loyaltyPoints}</span>
                      <span>Next Reward: 500 points</span>
                    </div>
                    <Progress 
                      value={(selectedCustomer.loyaltyPoints % 500) / 500 * 100} 
                      className="h-2"
                    />
                    <p className="text-xs text-muted-foreground">
                      {500 - (selectedCustomer.loyaltyPoints % 500)} points to next reward
                    </p>
                  </div>
                </div>
              </div>

              <DialogFooter className="flex gap-2">
                <Button variant="outline">
                  <Edit className="h-4 w-4 mr-2" />
                  Edit Customer
                </Button>
                <Button variant="outline">
                  <MessageSquare className="h-4 w-4 mr-2" />
                  Send Message
                </Button>
                <Button variant="outline">
                  <FileText className="h-4 w-4 mr-2" />
                  View All Orders
                </Button>
                <Button onClick={() => setShowDetailsDialog(false)}>
                  Close
                </Button>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>

      {/* Empty State */}
      {filteredCustomers.length === 0 && (
        <Card>
          <CardContent className="py-16">
            <div className="text-center">
              {searchQuery || filterType !== 'all' ? (
                <>
                  <Search className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-medium mb-2">No customers found</h3>
                  <p className="text-muted-foreground mb-4">
                    Try adjusting your search or filter criteria
                  </p>
                  <Button
                    variant="outline"
                    onClick={() => {
                      setSearchQuery('');
                      setFilterType('all');
                    }}
                  >
                    Clear Filters
                  </Button>
                </>
              ) : (
                <>
                  <Users className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-medium mb-2">No customers yet</h3>
                  <p className="text-muted-foreground mb-4">
                    Start building your customer base by adding your first customer
                  </p>
                  <Button onClick={() => setShowAddDialog(true)}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add First Customer
                  </Button>
                </>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
};

export default CustomerPage;