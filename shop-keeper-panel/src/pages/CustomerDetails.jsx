import React, { useState, useEffect } from 'react';
import {
  ArrowLeft,
  Edit,
  Phone,
  Mail,
  MapPin,
  Calendar,
  Gift,
  Star,
  TrendingUp,
  TrendingDown,
  ShoppingBag,
  CreditCard,
  Clock,
  MessageSquare,
  FileText,
  Download,
  Share,
  AlertCircle,
  Plus,
  Eye,
  Package,
  Receipt,
  Users,
  Heart,
  Target,
  Award,
  Activity,
  DollarSign,
  MoreHorizontal,
  Filter,
  Search,
  RefreshCw
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';
import { Textarea } from '@/components/ui/textarea';
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
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
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell
} from 'recharts';
import { Label } from '@/components/ui/label';

const CustomerDetailsPage = ({ customerId = 1 }) => {
  const [customer, setCustomer] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [showNoteDialog, setShowNoteDialog] = useState(false);
  const [newNote, setNewNote] = useState('');
  const [orderFilter, setOrderFilter] = useState('all');
  const [orderSort, setOrderSort] = useState('date');

  // Mock customer data
  const customerData = {
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
    birthday: '1985-03-15',
    notes: [
      {
        id: 1,
        text: 'Prefers organic products. Usually shops on weekends.',
        date: '2024-12-15',
        author: 'Mohammad Rahman'
      },
      {
        id: 2,
        text: 'Called to inquire about bulk pricing for rice.',
        date: '2024-11-22',
        author: 'Shop Assistant'
      }
    ],
    orders: [
      {
        id: '#1234',
        date: '2025-01-20',
        amount: 2450,
        items: 8,
        status: 'completed',
        paymentMethod: 'bKash',
        products: [
          { name: 'Basmati Rice 5kg', quantity: 2, price: 500 },
          { name: 'Cooking Oil 1L', quantity: 3, price: 150 },
          { name: 'Onion 1kg', quantity: 5, price: 100 }
        ]
      },
      {
        id: '#1189',
        date: '2025-01-15',
        amount: 1890,
        items: 6,
        status: 'completed',
        paymentMethod: 'Cash',
        products: [
          { name: 'Dal 1kg', quantity: 2, price: 180 },
          { name: 'Sugar 1kg', quantity: 3, price: 120 },
          { name: 'Tea Leaves 500g', quantity: 2, price: 200 }
        ]
      },
      {
        id: '#1156',
        date: '2025-01-08',
        amount: 3200,
        items: 12,
        status: 'completed',
        paymentMethod: 'Nagad',
        products: [
          { name: 'Flour 2kg', quantity: 4, price: 160 },
          { name: 'Milk 1L', quantity: 6, price: 100 },
          { name: 'Salt 1kg', quantity: 8, price: 50 }
        ]
      },
      {
        id: '#1098',
        date: '2024-12-28',
        amount: 1650,
        items: 5,
        status: 'completed',
        paymentMethod: 'bKash',
        products: [
          { name: 'Potato 1kg', quantity: 10, price: 80 },
          { name: 'Onion 1kg', quantity: 5, price: 100 }
        ]
      }
    ],
    purchaseHistory: [
      { month: 'Jul', amount: 3200, orders: 8 },
      { month: 'Aug', amount: 4100, orders: 10 },
      { month: 'Sep', amount: 3800, orders: 9 },
      { month: 'Oct', amount: 5200, orders: 12 },
      { month: 'Nov', amount: 4800, orders: 11 },
      { month: 'Dec', amount: 6100, orders: 14 },
      { month: 'Jan', amount: 4500, orders: 10 }
    ],
    favoriteProducts: [
      { name: 'Basmati Rice 5kg', purchases: 15, lastBought: '2025-01-20' },
      { name: 'Cooking Oil 1L', purchases: 12, lastBought: '2025-01-18' },
      { name: 'Onion 1kg', purchases: 18, lastBought: '2025-01-15' },
      { name: 'Dal 1kg', purchases: 8, lastBought: '2025-01-10' }
    ],
    categorySpending: [
      { name: 'Rice & Grains', amount: 15200, percentage: 33 },
      { name: 'Vegetables', amount: 8900, percentage: 19 },
      { name: 'Cooking Essentials', amount: 7800, percentage: 17 },
      { name: 'Dairy Products', amount: 6200, percentage: 14 },
      { name: 'Spices & Condiments', amount: 4800, percentage: 10 },
      { name: 'Others', amount: 2880, percentage: 7 }
    ]
  };

  useEffect(() => {
    // Simulate API call to fetch customer data
    setCustomer(customerData);
  }, [customerId]);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

  const getStatusColor = (status) => {
    switch (status) {
      case 'VIP': return 'bg-purple-100 text-purple-700 border-purple-200';
      case 'Gold': return 'bg-yellow-100 text-yellow-700 border-yellow-200';
      case 'Silver': return 'bg-gray-100 text-gray-700 border-gray-200';
      case 'Bronze': return 'bg-orange-100 text-orange-700 border-orange-200';
      default: return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'VIP': return <Star className="h-4 w-4" />;
      case 'Gold': return <Award className="h-4 w-4" />;
      case 'Silver': return <Target className="h-4 w-4" />;
      case 'Bronze': return <Activity className="h-4 w-4" />;
      default: return <Users className="h-4 w-4" />;
    }
  };

  const addNote = () => {
    if (newNote.trim()) {
      const note = {
        id: customer.notes.length + 1,
        text: newNote,
        date: new Date().toISOString().split('T')[0],
        author: 'Mohammad Rahman'
      };
      setCustomer({
        ...customer,
        notes: [note, ...customer.notes]
      });
      setNewNote('');
      setShowNoteDialog(false);
    }
  };

  const filteredOrders = customer?.orders.filter(order => {
    if (orderFilter === 'all') return true;
    if (orderFilter === 'completed') return order.status === 'completed';
    if (orderFilter === 'recent') {
      const orderDate = new Date(order.date);
      const thirtyDaysAgo = new Date();
      thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);
      return orderDate >= thirtyDaysAgo;
    }
    return true;
  });

  if (!customer) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <RefreshCw className="h-8 w-8 animate-spin mx-auto mb-4 text-muted-foreground" />
          <p>Loading customer details...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Button variant="outline" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Customers
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-foreground">Customer Details</h1>
            <p className="text-muted-foreground">
              Complete profile and purchase history
            </p>
          </div>
        </div>
        <div className="flex space-x-3">
          <Button variant="outline">
            <Share className="h-4 w-4 mr-2" />
            Share
          </Button>
          <Button variant="outline">
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline">
                <MoreHorizontal className="h-4 w-4 mr-2" />
                Actions
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem>
                <MessageSquare className="h-4 w-4 mr-2" />
                Send Message
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Gift className="h-4 w-4 mr-2" />
                Add Loyalty Points
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Receipt className="h-4 w-4 mr-2" />
                Create Invoice
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="text-destructive">
                <AlertCircle className="h-4 w-4 mr-2" />
                Deactivate Customer
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {/* Customer Profile Card */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-start justify-between">
            <div className="flex items-center space-x-6">
              <Avatar className="h-20 w-20">
                <AvatarFallback className="text-2xl bg-primary text-primary-foreground">
                  {customer.name.split(' ').map(n => n[0]).join('')}
                </AvatarFallback>
              </Avatar>
              <div className="space-y-2">
                <div className="flex items-center space-x-3">
                  <h2 className="text-2xl font-bold">{customer.name}</h2>
                  <Badge className={cn("border", getStatusColor(customer.status))}>
                    {getStatusIcon(customer.status)}
                    <span className="ml-2">{customer.status} Customer</span>
                  </Badge>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                  <div className="flex items-center space-x-2">
                    <Phone className="h-4 w-4 text-muted-foreground" />
                    <span>{customer.phone}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Mail className="h-4 w-4 text-muted-foreground" />
                    <span>{customer.email}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <span>{customer.area}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <span>Joined {formatDate(customer.joinDate)}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Clock className="h-4 w-4 text-muted-foreground" />
                    <span>Last purchase {formatDate(customer.lastPurchase)}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Gift className="h-4 w-4 text-muted-foreground" />
                    <span>Birthday {formatDate(customer.birthday)}</span>
                  </div>
                </div>
              </div>
            </div>
            <Button onClick={() => setIsEditing(true)}>
              <Edit className="h-4 w-4 mr-2" />
              Edit Profile
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Purchases</p>
                <p className="text-2xl font-bold text-green-600">
                  {formatCurrency(customer.totalPurchases)}
                </p>
                <div className="flex items-center space-x-1 mt-1">
                  <TrendingUp className="h-3 w-3 text-green-500" />
                  <span className="text-xs text-green-500">+12.5% this month</span>
                </div>
              </div>
              <DollarSign className="h-8 w-8 text-green-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Orders</p>
                <p className="text-2xl font-bold text-blue-600">{customer.totalOrders}</p>
                <div className="flex items-center space-x-1 mt-1">
                  <TrendingUp className="h-3 w-3 text-blue-500" />
                  <span className="text-xs text-blue-500">+8% this month</span>
                </div>
              </div>
              <ShoppingBag className="h-8 w-8 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Average Order</p>
                <p className="text-2xl font-bold text-purple-600">
                  {formatCurrency(customer.averageOrderValue)}
                </p>
                <div className="flex items-center space-x-1 mt-1">
                  <TrendingUp className="h-3 w-3 text-purple-500" />
                  <span className="text-xs text-purple-500">+5.2% vs avg</span>
                </div>
              </div>
              <CreditCard className="h-8 w-8 text-purple-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Loyalty Points</p>
                <p className="text-2xl font-bold text-orange-600">{customer.loyaltyPoints}</p>
                <div className="flex items-center space-x-1 mt-1">
                  <Star className="h-3 w-3 text-orange-500" />
                  <span className="text-xs text-orange-500">VIP Status</span>
                </div>
              </div>
              <Award className="h-8 w-8 text-orange-600" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Tabs */}
      <Card>
        <CardContent className="p-0">
          <Tabs defaultValue="orders" className="w-full">
            <div className="border-b">
              <TabsList className="grid w-full grid-cols-5 bg-transparent">
                <TabsTrigger value="orders" className="flex items-center space-x-2">
                  <Receipt className="h-4 w-4" />
                  <span>Orders</span>
                </TabsTrigger>
                <TabsTrigger value="analytics" className="flex items-center space-x-2">
                  <TrendingUp className="h-4 w-4" />
                  <span>Analytics</span>
                </TabsTrigger>
                <TabsTrigger value="products" className="flex items-center space-x-2">
                  <Package className="h-4 w-4" />
                  <span>Products</span>
                </TabsTrigger>
                <TabsTrigger value="notes" className="flex items-center space-x-2">
                  <FileText className="h-4 w-4" />
                  <span>Notes</span>
                </TabsTrigger>
                <TabsTrigger value="loyalty" className="flex items-center space-x-2">
                  <Heart className="h-4 w-4" />
                  <span>Loyalty</span>
                </TabsTrigger>
              </TabsList>
            </div>

            {/* Orders Tab */}
            <TabsContent value="orders" className="p-6 space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <h3 className="text-lg font-semibold">Order History</h3>
                  <Badge variant="secondary">{customer.orders.length} total orders</Badge>
                </div>
                <div className="flex items-center space-x-2">
                  <Select value={orderFilter} onValueChange={setOrderFilter}>
                    <SelectTrigger className="w-40">
                      <Filter className="h-4 w-4 mr-2" />
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Orders</SelectItem>
                      <SelectItem value="completed">Completed</SelectItem>
                      <SelectItem value="recent">Last 30 Days</SelectItem>
                    </SelectContent>
                  </Select>
                  <Button variant="outline" size="sm">
                    <Plus className="h-4 w-4 mr-2" />
                    New Order
                  </Button>
                </div>
              </div>

              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Order ID</TableHead>
                    <TableHead>Date</TableHead>
                    <TableHead>Items</TableHead>
                    <TableHead>Amount</TableHead>
                    <TableHead>Payment</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredOrders?.map((order) => (
                    <TableRow key={order.id}>
                      <TableCell className="font-medium">{order.id}</TableCell>
                      <TableCell>{formatDate(order.date)}</TableCell>
                      <TableCell>{order.items} items</TableCell>
                      <TableCell className="font-semibold">
                        {formatCurrency(order.amount)}
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">{order.paymentMethod}</Badge>
                      </TableCell>
                      <TableCell>
                        <Badge variant="default">{order.status}</Badge>
                      </TableCell>
                      <TableCell>
                        <Button variant="ghost" size="sm">
                          <Eye className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TabsContent>

            {/* Analytics Tab */}
            <TabsContent value="analytics" className="p-6 space-y-6">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Purchase Trend */}
                <Card>
                  <CardHeader>
                    <CardTitle>Purchase Trend</CardTitle>
                    <CardDescription>Monthly spending pattern</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <ResponsiveContainer width="100%" height={300}>
                      <AreaChart data={customer.purchaseHistory}>
                        <defs>
                          <linearGradient id="amountGradient" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3}/>
                            <stop offset="95%" stopColor="#22c55e" stopOpacity={0}/>
                          </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                        <XAxis dataKey="month" />
                        <YAxis />
                        <Tooltip 
                          formatter={(value, name) => [
                            name === 'amount' ? formatCurrency(value) : value,
                            name === 'amount' ? 'Amount' : 'Orders'
                          ]}
                        />
                        <Area 
                          type="monotone" 
                          dataKey="amount" 
                          stroke="#22c55e" 
                          fillOpacity={1} 
                          fill="url(#amountGradient)" 
                        />
                      </AreaChart>
                    </ResponsiveContainer>
                  </CardContent>
                </Card>

                {/* Category Spending */}
                <Card>
                  <CardHeader>
                    <CardTitle>Category Spending</CardTitle>
                    <CardDescription>Spending distribution by category</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      {customer.categorySpending.map((category, index) => (
                        <div key={index} className="space-y-2">
                          <div className="flex justify-between text-sm">
                            <span className="font-medium">{category.name}</span>
                            <span>{formatCurrency(category.amount)} ({category.percentage}%)</span>
                          </div>
                          <Progress value={category.percentage} className="h-2" />
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            {/* Favorite Products Tab */}
            <TabsContent value="products" className="p-6 space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold">Favorite Products</h3>
                <Badge variant="secondary">Top purchases</Badge>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {customer.favoriteProducts.map((product, index) => (
                  <Card key={index}>
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <h4 className="font-medium">{product.name}</h4>
                          <p className="text-sm text-muted-foreground">
                            Purchased {product.purchases} times
                          </p>
                          <p className="text-xs text-muted-foreground">
                            Last bought: {formatDate(product.lastBought)}
                          </p>
                        </div>
                        <div className="text-right">
                          <Badge variant="secondary">{product.purchases}x</Badge>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </TabsContent>

            {/* Notes Tab */}
            <TabsContent value="notes" className="p-6 space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold">Customer Notes</h3>
                <Dialog open={showNoteDialog} onOpenChange={setShowNoteDialog}>
                  <DialogTrigger asChild>
                    <Button>
                      <Plus className="h-4 w-4 mr-2" />
                      Add Note
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Add Customer Note</DialogTitle>
                      <DialogDescription>
                        Record important information about this customer
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                      <div>
                        <Label htmlFor="note">Note</Label>
                        <Textarea
                          id="note"
                          value={newNote}
                          onChange={(e) => setNewNote(e.target.value)}
                          placeholder="Enter your note here..."
                          rows={4}
                        />
                      </div>
                    </div>
                    <DialogFooter>
                      <Button variant="outline" onClick={() => setShowNoteDialog(false)}>
                        Cancel
                      </Button>
                      <Button onClick={addNote}>Add Note</Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
              <div className="space-y-4">
                {customer.notes.map((note) => (
                  <Card key={note.id}>
                    <CardContent className="p-4">
                      <div className="space-y-2">
                        <p className="text-sm">{note.text}</p>
                        <div className="flex items-center justify-between text-xs text-muted-foreground">
                          <span>By {note.author}</span>
                          <span>{formatDate(note.date)}</span>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </TabsContent>

            {/* Loyalty Program Tab */}
            <TabsContent value="loyalty" className="p-6 space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Loyalty Status</CardTitle>
                    <CardDescription>Current tier and benefits</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="text-center">
                      <div className="inline-flex items-center justify-center w-16 h-16 bg-purple-100 rounded-full mb-4">
                        <Star className="h-8 w-8 text-purple-600" />
                      </div>
                      <h3 className="text-2xl font-bold text-purple-600">VIP Member</h3>
                      <p className="text-sm text-muted-foreground">Highest tier achieved</p>
                    </div>
                    <div className="space-y-3">
                      <div className="flex items-center justify-between">
                        <span className="text-sm">Current Points</span>
                        <span className="font-semibold">500 points</span>
                      </div>
                      <Progress 
                        value={(customer.loyaltyPoints % 500) / 500 * 100} 
                        className="h-3"
                      />
                      <p className="text-xs text-muted-foreground text-center">
                        {500 - (customer.loyaltyPoints % 500)} points to next reward
                      </p>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>VIP Benefits</CardTitle>
                    <CardDescription>Exclusive perks and rewards</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      <div className="flex items-center space-x-3">
                        <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                        <span className="text-sm">10% discount on all purchases</span>
                      </div>
                      <div className="flex items-center space-x-3">
                        <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                        <span className="text-sm">Free home delivery</span>
                      </div>
                      <div className="flex items-center space-x-3">
                        <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                        <span className="text-sm">Priority customer support</span>
                      </div>
                      <div className="flex items-center space-x-3">
                        <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                        <span className="text-sm">Exclusive product access</span>
                      </div>
                      <div className="flex items-center space-x-3">
                        <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                        <span className="text-sm">Birthday special offers</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Loyalty History */}
              <Card>
                <CardHeader>
                  <CardTitle>Points History</CardTitle>
                  <CardDescription>Recent loyalty point transactions</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                          <Plus className="h-4 w-4 text-green-600" />
                        </div>
                        <div>
                          <p className="font-medium text-sm">Purchase Reward</p>
                          <p className="text-xs text-muted-foreground">Order #1234 - Jan 20, 2025</p>
                        </div>
                      </div>
                      <span className="font-semibold text-green-600">+24 points</span>
                    </div>
                    <div className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
                          <Gift className="h-4 w-4 text-red-600" />
                        </div>
                        <div>
                          <p className="font-medium text-sm">Reward Redeemed</p>
                          <p className="text-xs text-muted-foreground">Free delivery - Jan 18, 2025</p>
                        </div>
                      </div>
                      <span className="font-semibold text-red-600">-50 points</span>
                    </div>
                    <div className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                          <Plus className="h-4 w-4 text-green-600" />
                        </div>
                        <div>
                          <p className="font-medium text-sm">Purchase Reward</p>
                          <p className="text-xs text-muted-foreground">Order #1189 - Jan 15, 2025</p>
                        </div>
                      </div>
                      <span className="font-semibold text-green-600">+18 points</span>
                    </div>
                    <div className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                          <Star className="h-4 w-4 text-blue-600" />
                        </div>
                        <div>
                          <p className="font-medium text-sm">VIP Tier Achieved</p>
                          <p className="text-xs text-muted-foreground">Tier upgrade bonus - Jan 10, 2025</p>
                        </div>
                      </div>
                      <span className="font-semibold text-blue-600">+100 points</span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>

      {/* Edit Profile Dialog */}
      <Dialog open={isEditing} onOpenChange={setIsEditing}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Edit Customer Profile</DialogTitle>
            <DialogDescription>
              Update customer information and preferences
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="editFirstName">First Name</Label>
                <Input id="editFirstName" defaultValue="Fatima" />
              </div>
              <div>
                <Label htmlFor="editLastName">Last Name</Label>
                <Input id="editLastName" defaultValue="Khan" />
              </div>
            </div>
            <div>
              <Label htmlFor="editPhone">Phone Number</Label>
              <Input id="editPhone" defaultValue={customer.phone} />
            </div>
            <div>
              <Label htmlFor="editEmail">Email Address</Label>
              <Input id="editEmail" type="email" defaultValue={customer.email} />
            </div>
            <div>
              <Label htmlFor="editAddress">Address</Label>
              <Textarea id="editAddress" defaultValue={customer.address} />
            </div>
            <div>
              <Label htmlFor="editBirthday">Birthday</Label>
              <Input id="editBirthday" type="date" defaultValue={customer.birthday} />
            </div>
            <div>
              <Label htmlFor="editStatus">Customer Status</Label>
              <Select defaultValue={customer.status.toLowerCase()}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bronze">Bronze</SelectItem>
                  <SelectItem value="silver">Silver</SelectItem>
                  <SelectItem value="gold">Gold</SelectItem>
                  <SelectItem value="vip">VIP</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditing(false)}>
              Cancel
            </Button>
            <Button onClick={() => setIsEditing(false)}>
              Save Changes
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Quick Actions Footer */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-semibold text-lg">Quick Actions</h3>
              <p className="text-sm text-muted-foreground">
                Common actions for {customer.name}
              </p>
            </div>
            <div className="flex space-x-3">
              <Button variant="outline">
                <MessageSquare className="h-4 w-4 mr-2" />
                Send SMS
              </Button>
              <Button variant="outline">
                <Phone className="h-4 w-4 mr-2" />
                Call Customer
              </Button>
              <Button variant="outline">
                <Receipt className="h-4 w-4 mr-2" />
                Create Invoice
              </Button>
              <Button>
                <ShoppingBag className="h-4 w-4 mr-2" />
                New Order
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default CustomerDetailsPage;