import React, { useState } from 'react';
import {
  CreditCard,
  Users,
  TrendingUp,
  AlertCircle,
  Plus,
  Search,
  Filter,
  MoreHorizontal,
  Calendar,
  DollarSign,
  CheckCircle,
  XCircle,
  Clock,
  Package,
  Eye,
  FileText,
  Edit,
  RefreshCw
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Link, useNavigate } from 'react-router-dom';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'; 

// Helper function to count subscriptions by status and calculate stats
const calculateSubscriptionStats = (subs) => {
  const active = subs.filter(sub => sub.status === 'Active').length;
  const pending = subs.filter(sub => sub.status === 'Trial').length;
  const cancelled = subs.filter(sub => sub.status === 'Cancelled').length;
  const total = subs.length;
  
  // Calculate monthly revenue from active subscriptions
  const monthlyRevenue = subs
    .filter(sub => sub.status === 'Active')
    .reduce((total, sub) => {
      // For annual plans, divide by 12 to get monthly equivalent
      const monthlyAmount = sub.billingCycle === 'Annual' ? sub.amount / 12 : sub.amount;
      return total + monthlyAmount;
    }, 0);
    
  return [
    { label: 'Total Subscriptions', value: total, icon: CreditCard, color: 'text-blue-600', change: '+12%' },
    { label: 'Active', value: active, icon: CheckCircle, color: 'text-green-600', change: '+8%' },
    { label: 'Pending', value: pending, icon: Clock, color: 'text-yellow-600', change: '+15%' },
    { label: 'Cancelled', value: cancelled, icon: XCircle, color: 'text-red-600', change: '-5%' },
    { label: 'Monthly Revenue', value: `$${Math.round(monthlyRevenue).toLocaleString()}`, icon: DollarSign, color: 'text-emerald-600', change: '+10.5%' }
  ];
};

const subscriptions = [
  {
    id: 'sub_001',
    tenant: 'TechCorp Solutions',
    plan: 'Enterprise',
    status: 'Active',
    amount: 299,
    billingCycle: 'Monthly',
    nextBilling: '2024-02-15',
    startDate: '2023-08-15',
    paymentStatus: 'Current',
    features: ['Unlimited Users', 'Advanced Analytics', 'Priority Support']
  },
  {
    id: 'sub_002',
    tenant: 'StartupHub',
    plan: 'Professional',
    status: 'Active',
    amount: 99,
    billingCycle: 'Monthly',
    nextBilling: '2024-02-20',
    startDate: '2024-01-20',
    paymentStatus: 'Current',
    features: ['Up to 50 Users', 'Basic Analytics', 'Email Support']
  },
  {
    id: 'sub_003',
    tenant: 'RetailMax',
    plan: 'Basic',
    status: 'Trial',
    amount: 29,
    billingCycle: 'Monthly',
    nextBilling: '2024-02-10',
    startDate: '2024-02-01',
    paymentStatus: 'Current',
    features: ['Up to 10 Users', 'Basic Features']
  },
  {
    id: 'sub_004',
    tenant: 'GlobalTrade Inc',
    plan: 'Enterprise',
    status: 'Cancelled',
    amount: 299,
    billingCycle: 'Annual',
    nextBilling: null,
    startDate: '2023-05-10',
    paymentStatus: 'Cancelled',
    features: ['Unlimited Users', 'Advanced Analytics', 'Priority Support']
  },
  {
    id: 'sub_005',
    tenant: 'EcoFriendly Co',
    plan: 'Professional',
    status: 'Active',
    amount: 99,
    billingCycle: 'Monthly',
    nextBilling: '2024-02-25',
    startDate: '2023-11-25',
    paymentStatus: 'Overdue',
    features: ['Up to 50 Users', 'Basic Analytics', 'Email Support']
  }
];

export default function SubscriptionsHome() {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('All');
  const [planFilter, setPlanFilter] = useState('All');

  const filteredSubscriptions = subscriptions.filter(sub => {
    const matchesSearch = sub.tenant.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         sub.id.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'All' || sub.status === statusFilter;
    const matchesPlan = planFilter === 'All' || sub.plan === planFilter;
    return matchesSearch && matchesStatus && matchesPlan;
  });
  
  // Calculate subscription stats based on the current data
  const subscriptionStats = calculateSubscriptionStats(subscriptions);

  const getStatusBadge = (status) => {
    switch (status) {
      case 'Active':
        return <Badge variant="success">Active</Badge>;
      case 'Trial':
        return <Badge variant="secondary">Trial</Badge>;
      case 'Cancelled':
        return <Badge variant="destructive">Cancelled</Badge>;
      case 'Expired':
        return <Badge variant="outline">Expired</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };
  
  const getPaymentStatusBadge = (status) => {
    switch (status) {
      case 'Current':
        return <Badge variant="success">Current</Badge>;
      case 'Pending':
        return <Badge variant="secondary">Pending</Badge>;
      case 'Overdue':
        return <Badge variant="destructive">Overdue</Badge>;
      case 'Cancelled':
        return <Badge variant="outline">Cancelled</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center px-6">
          <div className="flex items-center gap-4">
            <CreditCard className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Subscriptions</h1>
            <Badge variant="outline" className="ml-2">Admin Panel</Badge>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Stats Row */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4 w-full">
          {subscriptionStats.map((stat) => (
            <div key={stat.label} className="flex items-center gap-3 bg-card p-3 rounded-lg border">
              <div className={`p-2 rounded-md ${stat.color} bg-opacity-10 shrink-0`}>
                <stat.icon className={`h-5 w-5 ${stat.color}`} />
              </div>
              <div className="flex-grow">
                <p className="text-sm text-muted-foreground">{stat.label}</p>
                <div className="flex items-center gap-1.5">
                  <p className="text-lg font-semibold">{stat.value}</p>
                  <span className={`text-xs ${stat.change.startsWith('+') ? 'text-green-600' : 'text-red-600'} font-medium`}>
                    {stat.change}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>

 

        {/* Subscriptions Table */}
        <Card>
          <CardHeader className="py-4">
                {/* Filters and Search */}
        <div className="flex flex-wrap items-center gap-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
            <Input
              type="text"
              placeholder="Search subscriptions by ID, tenant name or plan..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10 w-full"
            />
          </div>
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter by Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="All">All Statuses</SelectItem>
              <SelectItem value="Active">Active</SelectItem>
              <SelectItem value="Trial">Trial</SelectItem>
              <SelectItem value="Cancelled">Cancelled</SelectItem>
              <SelectItem value="Expired">Expired</SelectItem>
            </SelectContent>
          </Select>
          <Select value={planFilter} onValueChange={setPlanFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter by Plan" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="All">All Plans</SelectItem>
              <SelectItem value="Basic">Basic</SelectItem>
              <SelectItem value="Professional">Professional</SelectItem>
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
                    <TableHead>ID</TableHead>
                    <TableHead>Tenant</TableHead>
                    <TableHead>Plan</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Payment Status</TableHead>
                    <TableHead>Amount</TableHead>
                    <TableHead>Billing Cycle</TableHead>
                    <TableHead>Next Billing</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredSubscriptions.map((subscription) => (
                    <TableRow key={subscription.id}>
                      <TableCell>
                        <Link 
                          to={`/subscriptions/${subscription.id}`} 
                          className="text-primary hover:underline font-medium"
                        >
                          {subscription.id}
                        </Link>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <Users className="h-4 w-4 text-muted-foreground" />
                          {subscription.tenant}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">{subscription.plan}</Badge>
                      </TableCell>
                      <TableCell>
                        {getStatusBadge(subscription.status)}
                      </TableCell>
                      <TableCell>
                        {getPaymentStatusBadge(subscription.paymentStatus)}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1">
                          <DollarSign className="h-4 w-4 text-muted-foreground" />
                          {subscription.amount}
                        </div>
                      </TableCell>
                      <TableCell>{subscription.billingCycle}</TableCell>
                      <TableCell>
                        {subscription.nextBilling ? (
                          <div className="flex items-center gap-1">
                            <Calendar className="h-4 w-4 text-muted-foreground" />
                            {subscription.nextBilling}
                          </div>
                        ) : (
                          <span className="text-muted-foreground">-</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button size="sm" variant="ghost">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => navigate(`/subscriptions/${subscription.id}`)}>
                              <Eye className="h-4 w-4 mr-2" />
                              View Details
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => navigate('/billing')}>
                              <FileText className="h-4 w-4 mr-2" />
                              View Invoices
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => navigate(`/subscriptions/${subscription.id}/modify`)}>
                              <Edit className="h-4 w-4 mr-2" />
                              Edit Subscription
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem>
                              <RefreshCw className="h-4 w-4 mr-2" />
                              Renew Subscription
                            </DropdownMenuItem>
                            <DropdownMenuItem className="text-destructive">
                              <XCircle className="h-4 w-4 mr-2" />
                              Cancel Subscription
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
}
