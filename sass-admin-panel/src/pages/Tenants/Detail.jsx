import React from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  Store,
  Users,
  DollarSign,
  Package,
  Mail,
  Phone,
  MapPin,
  Globe,
  Star,
  Activity,
  Clock,
  TrendingUp,
  Edit,
  Settings,
  BarChart3,
  CreditCard,
  ArrowLeft,
  MoreHorizontal,
  ShoppingCart,
  UserCheck,
  Ticket,
  MessageSquare,
  HelpCircle,
  ExternalLink,
  CheckCircle2,
  Zap,
  History,
  Receipt
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
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

const mockTenants = {
  1: {
    id: 1,
    name: 'Rahman Electronics',
    businessType: 'Electronics & Technology',
    plan: 'Business',
    status: 'active',
    monthlyRevenue: 5000,
    totalRevenue: 35000,
    productsCount: 245,
    ordersCount: 1250,
    customersCount: 89,
    owner: 'Abdul Rahman',
    email: 'rahman@email.com',
    phone: '01712345678',
    address: 'House 15, Road 7, Dhanmondi, Dhaka 1205',
    location: 'Dhanmondi, Dhaka',
    website: 'https://rahmanelectronics.com',
    joinDate: '2024-01-15',
    subscriptionEnd: '2024-08-15',
    lastLogin: '2024-07-24T10:30:00',
    rating: 4.8,
    notes: 'Premium customer with excellent payment history',
    recentOrders: [
      { id: 'ORD-001', customer: 'Ahmed Khan', amount: 1250, status: 'Completed', date: '2024-07-24' },
      { id: 'ORD-002', customer: 'Fatima Ahmed', amount: 850, status: 'Processing', date: '2024-07-23' },
      { id: 'ORD-003', customer: 'Mohammad Ali', amount: 2100, status: 'Completed', date: '2024-07-22' },
      { id: 'ORD-004', customer: 'Rashida Begum', amount: 450, status: 'Cancelled', date: '2024-07-21' }
    ],
    recentActivity: [
      { id: 1, action: 'Paid Monthly Invoice', date: '2024-07-24T09:00:00', amount: 5000, type: 'payment' },
      { id: 2, action: 'Added 15 new products', date: '2024-07-23T14:30:00', amount: null, type: 'product' },
      { id: 3, action: 'Customer support ticket resolved', date: '2024-07-22T11:15:00', amount: null, type: 'support' },
      { id: 4, action: 'Updated business profile', date: '2024-07-21T16:45:00', amount: null, type: 'profile' },
      { id: 5, action: 'Subscription renewed', date: '2024-07-15T10:00:00', amount: 5000, type: 'subscription' }
    ],
    stats: {
      totalOrders: 1250,
      avgOrderValue: 1580,
      conversionRate: 3.2,
      customerRetention: 78,
      monthlyGrowth: 12.5,
      lastMonthRevenue: 4500
    },
    subscriptionHistory: [
      {
        id: 'SUB-001',
        plan: 'Business',
        startDate: '2024-01-15',
        endDate: '2024-08-15',
        amount: 5000,
        status: 'active',
        paymentMethod: 'Credit Card',
        autoRenew: true
      },
      {
        id: 'SUB-002',
        plan: 'Starter',
        startDate: '2023-07-15',
        endDate: '2024-01-15',
        amount: 2000,
        status: 'completed',
        paymentMethod: 'Bank Transfer',
        autoRenew: false
      }
    ],
    supportTickets: [
      {
        id: 'TKT-001',
        title: 'Payment gateway integration issue',
        description: 'Having trouble with Stripe payment integration for international customers',
        status: 'open',
        priority: 'high',
        category: 'Technical',
        createdAt: '2024-07-22T09:30:00',
        updatedAt: '2024-07-23T14:20:00',
        assignedTo: 'Support Team',
        messages: 3
      },
      {
        id: 'TKT-002',
        title: 'Feature request: Bulk product import',
        description: 'Need ability to import products in bulk via CSV file',
        status: 'in_progress',
        priority: 'medium',
        category: 'Feature Request',
        createdAt: '2024-07-20T11:15:00',
        updatedAt: '2024-07-21T16:45:00',
        assignedTo: 'Development Team',
        messages: 5
      },
      {
        id: 'TKT-003',
        title: 'Account verification documents',
        description: 'Submitted business license for account verification',
        status: 'resolved',
        priority: 'low',
        category: 'Account',
        createdAt: '2024-07-18T08:00:00',
        updatedAt: '2024-07-19T10:30:00',
        assignedTo: 'Compliance Team',
        messages: 2
      }
    ],
    billingHistory: [
      { id: 'INV-001', date: '2024-07-15', amount: 5000, status: 'paid', description: 'Business Plan - Monthly Subscription', dueDate: '2024-07-30' },
      { id: 'INV-002', date: '2024-06-15', amount: 5000, status: 'paid', description: 'Business Plan - Monthly Subscription', dueDate: '2024-06-30' },
      { id: 'INV-003', date: '2024-05-15', amount: 5000, status: 'paid', description: 'Business Plan - Monthly Subscription', dueDate: '2024-05-30' },
      { id: 'INV-004', date: '2024-04-15', amount: 5000, status: 'overdue', description: 'Business Plan - Monthly Subscription', dueDate: '2024-04-30' }
    ],
    integrations: [
      { name: 'Stripe', status: 'connected', type: 'Payment Gateway', lastSync: '2024-07-24T10:30:00' },
      { name: 'Google Analytics', status: 'connected', type: 'Analytics', lastSync: '2024-07-23T18:45:00' },
      { name: 'Facebook Pixel', status: 'disconnected', type: 'Marketing', lastSync: null },
      { name: 'Mailchimp', status: 'connected', type: 'Email Marketing', lastSync: '2024-07-24T06:15:00' }
    ]
  },
  // Add more mock tenants as needed
};

export default function TenantDetail() {
  const { id } = useParams();
  const tenant = mockTenants[id] || mockTenants[1]; // fallback to 1 for demo

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  const formatDate = (dateString) => new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

  const formatDateTime = (dateString) => new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });

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

  const getOrderStatusColor = (status) => {
    switch (status) {
      case 'Completed': return 'bg-green-100 text-green-700';
      case 'Processing': return 'bg-blue-100 text-blue-700';
      case 'Cancelled': return 'bg-red-100 text-red-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getTicketStatusColor = (status) => {
    switch (status) {
      case 'open': return 'bg-red-100 text-red-700';
      case 'in_progress': return 'bg-blue-100 text-blue-700';
      case 'resolved': return 'bg-green-100 text-green-700';
      case 'closed': return 'bg-gray-100 text-gray-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getPriorityColor = (priority) => {
    switch (priority) {
      case 'high': return 'bg-red-100 text-red-700';
      case 'medium': return 'bg-yellow-100 text-yellow-700';
      case 'low': return 'bg-green-100 text-green-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getBillingStatusColor = (status) => {
    switch (status) {
      case 'paid': return 'bg-green-100 text-green-700';
      case 'pending': return 'bg-yellow-100 text-yellow-700';
      case 'overdue': return 'bg-red-100 text-red-700';
      case 'failed': return 'bg-red-100 text-red-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getIntegrationStatusColor = (status) => {
    switch (status) {
      case 'connected': return 'bg-green-100 text-green-700';
      case 'disconnected': return 'bg-red-100 text-red-700';
      case 'error': return 'bg-yellow-100 text-yellow-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getActivityIcon = (type) => {
    switch (type) {
      case 'payment': return <DollarSign className="h-4 w-4 text-green-600" />;
      case 'product': return <Package className="h-4 w-4 text-blue-600" />;
      case 'support': return <Ticket className="h-4 w-4 text-yellow-600" />;
      case 'profile': return <Settings className="h-4 w-4 text-purple-600" />;
      case 'subscription': return <CreditCard className="h-4 w-4 text-orange-600" />;
      default: return <Activity className="h-4 w-4 text-gray-600" />;
    }
  };

  const getRatingStars = (rating) => {
    return Array.from({ length: 5 }, (_, i) => (
      <Star
        key={i}
        className={`h-4 w-4 ${
          i < Math.floor(rating)
            ? 'text-yellow-400 fill-current'
            : 'text-gray-300'
        }`}
      />
    ));
  };

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Link to="/tenants">
              <Button variant="ghost" size="sm">
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Tenants
              </Button>
            </Link>
            <div className="flex items-center gap-3">
              <Avatar className="h-10 w-10">
                <AvatarFallback className="text-sm font-medium">
                  {tenant.name.split(' ').map(n => n[0]).join('')}
                </AvatarFallback>
              </Avatar>
              <div>
                <h1 className="text-xl font-semibold">{tenant.name}</h1>
                <p className="text-sm text-muted-foreground">{tenant.businessType}</p>
              </div>
            </div>
            <Badge className={getPlanColor(tenant.plan)}>{tenant.plan}</Badge>
            <Badge className={getStatusColor(tenant.status)}>{tenant.status}</Badge>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm">
              <Edit className="h-4 w-4 mr-2" />
              Edit
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem>
                  <BarChart3 className="h-4 w-4 mr-2" />
                  View Analytics
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <CreditCard className="h-4 w-4 mr-2" />
                  Billing History
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Mail className="h-4 w-4 mr-2" />
                  Send Message
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem className="text-orange-600">
                  <Clock className="h-4 w-4 mr-2" />
                  Suspend Account
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Key Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
          <div className="flex items-center gap-3 bg-card p-4 rounded-lg border">
            <div className="p-2 rounded-md text-green-600 bg-green-100 shrink-0">
              <DollarSign className="h-5 w-5" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Monthly Revenue</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{formatCurrency(tenant.monthlyRevenue)}</p>
                <span className="text-xs text-green-600 font-medium">
                  +{((tenant.monthlyRevenue - tenant.stats.lastMonthRevenue) / tenant.stats.lastMonthRevenue * 100).toFixed(1)}%
                </span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3 bg-card p-4 rounded-lg border">
            <div className="p-2 rounded-md text-blue-600 bg-blue-100 shrink-0">
              <ShoppingCart className="h-5 w-5" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Total Orders</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{tenant.ordersCount.toLocaleString()}</p>
                <span className="text-xs text-blue-600 font-medium">+12%</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3 bg-card p-4 rounded-lg border">
            <div className="p-2 rounded-md text-purple-600 bg-purple-100 shrink-0">
              <Package className="h-5 w-5" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Products</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{tenant.productsCount}</p>
                <span className="text-xs text-purple-600 font-medium">+8%</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3 bg-card p-4 rounded-lg border">
            <div className="p-2 rounded-md text-orange-600 bg-orange-100 shrink-0">
              <UserCheck className="h-5 w-5" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Customers</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{tenant.customersCount}</p>
                <span className="text-xs text-orange-600 font-medium">+15%</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3 bg-card p-4 rounded-lg border">
            <div className="p-2 rounded-md text-yellow-600 bg-yellow-100 shrink-0">
              <TrendingUp className="h-5 w-5" />
            </div>
            <div className="flex-grow">
              <p className="text-sm text-muted-foreground">Avg Order Value</p>
              <div className="flex items-center gap-1.5">
                <p className="text-lg font-semibold">{formatCurrency(tenant.stats.avgOrderValue)}</p>
                <span className="text-xs text-yellow-600 font-medium">+5%</span>
              </div>
            </div>
          </div>
        </div>

        {/* Main Content - Comprehensive Tabs Layout */}
        <div className="grid grid-cols-1 xl:grid-cols-4 gap-6">
          {/* Left Sidebar - Business Info & Subscription */}
          <div className="space-y-6">
            {/* Business Information */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Store className="h-5 w-5" />
                  Business Info
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <Users className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium text-sm">Owner:</span>
                    <span className="text-sm">{tenant.owner}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Mail className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium text-sm">Email:</span>
                    <span className="text-sm truncate">{tenant.email}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Phone className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium text-sm">Phone:</span>
                    <span className="text-sm">{tenant.phone}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium text-sm">Location:</span>
                    <span className="text-sm">{tenant.location}</span>
                  </div>
                  {tenant.website && (
                    <div className="flex items-center gap-2">
                      <Globe className="h-4 w-4 text-muted-foreground" />
                      <span className="font-medium text-sm">Website:</span>
                      <a href={tenant.website} target="_blank" rel="noopener noreferrer" className="text-primary hover:underline text-sm flex items-center gap-1">
                        View Site <ExternalLink className="h-3 w-3" />
                      </a>
                    </div>
                  )}
                  <div className="flex items-center gap-2">
                    <Star className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium text-sm">Rating:</span>
                    <div className="flex items-center gap-1">
                      {getRatingStars(tenant.rating)}
                      <span className="text-sm ml-1">{tenant.rating}</span>
                    </div>
                  </div>
                </div>
                {tenant.notes && (
                  <div className="border-t pt-3">
                    <p className="font-medium text-sm mb-1">Notes:</p>
                    <p className="text-xs text-muted-foreground">{tenant.notes}</p>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Current Subscription */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Zap className="h-5 w-5" />
                  Current Subscription
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Plan:</span>
                  <Badge className={getPlanColor(tenant.plan)}>{tenant.plan}</Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Status:</span>
                  <Badge className={getStatusColor(tenant.status)}>{tenant.status}</Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Monthly Fee:</span>
                  <span className="text-sm font-semibold">{formatCurrency(tenant.monthlyRevenue)}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Joined:</span>
                  <span className="text-sm">{formatDate(tenant.joinDate)}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Expires:</span>
                  <span className="text-sm">{formatDate(tenant.subscriptionEnd)}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Auto Renew:</span>
                  <span className="text-sm flex items-center gap-1">
                    <CheckCircle2 className="h-3 w-3 text-green-600" />
                    Enabled
                  </span>
                </div>
              </CardContent>
            </Card>

            {/* Quick Stats */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="h-5 w-5" />
                  Performance
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium">Conversion Rate</span>
                  <span className="text-sm">{tenant.stats.conversionRate}%</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium">Retention</span>
                  <span className="text-sm">{tenant.stats.customerRetention}%</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium">Growth</span>
                  <span className="text-sm text-green-600">+{tenant.stats.monthlyGrowth}%</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium">Total Revenue</span>
                  <span className="text-sm font-semibold">{formatCurrency(tenant.totalRevenue)}</span>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Main Content Area */}
          <div className="xl:col-span-3 space-y-6">
            {/* Support Tickets */}
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <Ticket className="h-5 w-5" />
                  Support Tickets
                </CardTitle>
                <Button variant="outline" size="sm">
                  <HelpCircle className="h-4 w-4 mr-2" />
                  New Ticket
                </Button>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {tenant.supportTickets.map((ticket) => (
                    <div key={ticket.id} className="border rounded-lg p-4 hover:bg-muted/50 transition-colors">
                      <div className="flex items-start justify-between mb-3">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <span className="font-medium text-sm">{ticket.title}</span>
                            <Badge className={getTicketStatusColor(ticket.status)} variant="outline">
                              {ticket.status.replace('_', ' ')}
                            </Badge>
                            <Badge className={getPriorityColor(ticket.priority)} variant="outline">
                              {ticket.priority}
                            </Badge>
                          </div>
                          <p className="text-sm text-muted-foreground mb-2">{ticket.description}</p>
                          <div className="flex items-center gap-4 text-xs text-muted-foreground">
                            <span className="flex items-center gap-1">
                              <Ticket className="h-3 w-3" />
                              {ticket.id}
                            </span>
                            <span className="flex items-center gap-1">
                              <Users className="h-3 w-3" />
                              {ticket.assignedTo}
                            </span>
                            <span className="flex items-center gap-1">
                              <MessageSquare className="h-3 w-3" />
                              {ticket.messages} messages
                            </span>
                            <span className="flex items-center gap-1">
                              <Clock className="h-3 w-3" />
                              Updated {formatDateTime(ticket.updatedAt)}
                            </span>
                          </div>
                        </div>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="sm">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem>
                              <MessageSquare className="h-4 w-4 mr-2" />
                              View Ticket
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                              <Edit className="h-4 w-4 mr-2" />
                              Update Status
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                              <Users className="h-4 w-4 mr-2" />
                              Reassign
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Subscription History */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <History className="h-5 w-5" />
                  Subscription History
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Plan</TableHead>
                        <TableHead>Period</TableHead>
                        <TableHead>Amount</TableHead>
                        <TableHead>Payment Method</TableHead>
                        <TableHead>Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {tenant.subscriptionHistory.map((sub) => (
                        <TableRow key={sub.id}>
                          <TableCell>
                            <Badge className={getPlanColor(sub.plan)} variant="outline">
                              {sub.plan}
                            </Badge>
                          </TableCell>
                          <TableCell className="text-sm">
                            {formatDate(sub.startDate)} - {formatDate(sub.endDate)}
                          </TableCell>
                          <TableCell className="font-semibold">{formatCurrency(sub.amount)}</TableCell>
                          <TableCell className="text-sm">{sub.paymentMethod}</TableCell>
                          <TableCell>
                            <Badge className={getStatusColor(sub.status)} variant="outline">
                              {sub.status}
                            </Badge>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>

            {/* Billing History */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Receipt className="h-5 w-5" />
                  Billing History
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Invoice</TableHead>
                        <TableHead>Description</TableHead>
                        <TableHead>Amount</TableHead>
                        <TableHead>Due Date</TableHead>
                        <TableHead>Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {tenant.billingHistory.map((bill) => (
                        <TableRow key={bill.id}>
                          <TableCell className="font-mono text-sm">{bill.id}</TableCell>
                          <TableCell className="text-sm">{bill.description}</TableCell>
                          <TableCell className="font-semibold">{formatCurrency(bill.amount)}</TableCell>
                          <TableCell className="text-sm">{formatDate(bill.dueDate)}</TableCell>
                          <TableCell>
                            <Badge className={getBillingStatusColor(bill.status)}>
                              {bill.status}
                            </Badge>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>

            {/* Integrations */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Zap className="h-5 w-5" />
                  Integrations
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {tenant.integrations.map((integration, index) => (
                    <div key={index} className="border rounded-lg p-3 flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <div className={`w-2 h-2 rounded-full ${
                          integration.status === 'connected' ? 'bg-green-500' : 'bg-red-500'
                        }`} />
                        <div>
                          <p className="font-medium text-sm">{integration.name}</p>
                          <p className="text-xs text-muted-foreground">{integration.type}</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <Badge className={getIntegrationStatusColor(integration.status)} variant="outline">
                          {integration.status}
                        </Badge>
                        {integration.lastSync && (
                          <p className="text-xs text-muted-foreground mt-1">
                            Last sync: {formatDateTime(integration.lastSync)}
                          </p>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Recent Orders */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <ShoppingCart className="h-5 w-5" />
                  Recent Orders
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Order ID</TableHead>
                        <TableHead>Customer</TableHead>
                        <TableHead>Amount</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Date</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {tenant.recentOrders.map((order) => (
                        <TableRow key={order.id}>
                          <TableCell className="font-mono text-sm">{order.id}</TableCell>
                          <TableCell>{order.customer}</TableCell>
                          <TableCell className="font-semibold">{formatCurrency(order.amount)}</TableCell>
                          <TableCell>
                            <Badge className={getOrderStatusColor(order.status)}>
                              {order.status}
                            </Badge>
                          </TableCell>
                          <TableCell>{formatDate(order.date)}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>

            {/* Recent Activity */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Activity className="h-5 w-5" />
                  Recent Activity
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {tenant.recentActivity.map((activity) => (
                    <div key={activity.id} className="flex items-start gap-3 p-3 rounded-lg border hover:bg-muted/50">
                      <div className="mt-0.5">
                        {getActivityIcon(activity.type)}
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm font-medium">{activity.action}</p>
                        <div className="flex items-center gap-2 text-xs text-muted-foreground">
                          <span>{formatDateTime(activity.date)}</span>
                          {activity.amount && (
                            <>
                              <span>•</span>
                              <span className="font-medium">{formatCurrency(activity.amount)}</span>
                            </>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}