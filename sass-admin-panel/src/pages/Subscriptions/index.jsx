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
import { Link } from 'react-router-dom';
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

const subscriptionStats = [
  { label: 'Total Subscriptions', value: 1247, icon: CreditCard, color: 'text-blue-600', change: '+12%' },
  { label: 'Active', value: 1089, icon: CheckCircle, color: 'text-green-600', change: '+8%' },
  { label: 'Pending', value: 98, icon: Clock, color: 'text-yellow-600', change: '+15%' },
  { label: 'Cancelled', value: 60, icon: XCircle, color: 'text-red-600', change: '-5%' },
  { label: 'Monthly Revenue', value: '$124,350', icon: DollarSign, color: 'text-emerald-600', change: '+10.5%' }
];

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
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <CreditCard className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Subscriptions</h1>
            <Badge variant="outline" className="ml-2">Admin Panel</Badge>
          </div>
          <Button className="flex items-center gap-2">
            <Plus className="h-4 w-4" />
            New Subscription
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {subscriptionStats.map((stat) => (
            <Card key={stat.label}>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <stat.icon className={`h-8 w-8 ${stat.color}`} />
                    <div>
                      <p className="text-sm text-muted-foreground">{stat.label}</p>
                      <p className="text-2xl font-bold">{stat.value}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className={`text-sm font-medium ${
                      stat.change.startsWith('+') ? 'text-green-600' : 'text-red-600'
                    }`}>
                      {stat.change}
                    </p>
                    <p className="text-xs text-muted-foreground">vs last month</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Filters and Search */}
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  type="text"
                  placeholder="Search subscriptions..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-[140px]">
                  <SelectValue placeholder="All Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="cancelled">Cancelled</SelectItem>
                  <SelectItem value="expired">Expired</SelectItem>
                </SelectContent>
              </Select>
              <Select value={planFilter} onValueChange={setPlanFilter}>
                <SelectTrigger className="w-[140px]">
                  <SelectValue placeholder="All Plans" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Plans</SelectItem>
                  <SelectItem value="basic">Basic</SelectItem>
                  <SelectItem value="pro">Pro</SelectItem>
                  <SelectItem value="enterprise">Enterprise</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        {/* Subscriptions Table */}
        <Card>
          <CardHeader>
            <CardTitle>Subscription Management</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Subscription ID</TableHead>
                  <TableHead>Tenant</TableHead>
                  <TableHead>Plan</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Payment Status</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>Billing Cycle</TableHead>
                  <TableHead>Next Billing</TableHead>
                  <TableHead>Payment Status</TableHead>
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
                      {getPaymentStatusBadge(subscription.paymentStatus)}
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button size="sm" variant="ghost">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem>
                            <Eye className="h-4 w-4 mr-2" />
                            View Details
                          </DropdownMenuItem>
                          <DropdownMenuItem>
                            <FileText className="h-4 w-4 mr-2" />
                            View Invoices
                          </DropdownMenuItem>
                          <DropdownMenuItem>
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
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
