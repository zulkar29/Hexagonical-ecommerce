import React, { useState, useMemo } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  CreditCard,
  DollarSign,
  TrendingUp,
  AlertCircle,
  Plus,
  Search,
  Filter,
  MoreHorizontal,
  Calendar,
  Download,
  RefreshCw,
  CheckCircle,
  XCircle,
  Clock,
  FileText,
  Send,
  Eye,
  Edit,
  Trash2,
  Users,
  ArrowUpRight,
  ArrowDownRight,
  Ban,
  RotateCcw
} from 'lucide-react';
import { Card, CardHeader, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
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

const billingStats = [
  { label: 'Total Revenue', value: '$2,847,650', icon: DollarSign, color: 'text-green-600', change: '+12.5%' },
  { label: 'Monthly Recurring Revenue', value: '$285,400', icon: TrendingUp, color: 'text-blue-600', change: '+8.3%' },
  { label: 'Outstanding Invoices', value: 23, icon: AlertCircle, color: 'text-yellow-600', change: '-5.2%' },
  { label: 'Processed This Month', value: 1247, icon: CheckCircle, color: 'text-purple-600', change: '+15.8%' }
];

const invoices = [
  {
    id: 'INV-2024-001',
    tenant: 'TechCorp Solutions',
    amount: 299.00,
    status: 'Paid',
    dueDate: '2024-01-15',
    paidDate: '2024-01-14',
    plan: 'Enterprise',
    billingPeriod: 'Jan 2024',
    subscriptionId: 'sub_001'
  },
  {
    id: 'INV-2024-002',
    tenant: 'StartupHub',
    amount: 99.00,
    status: 'Pending',
    dueDate: '2024-01-20',
    paidDate: null,
    plan: 'Professional',
    billingPeriod: 'Jan 2024',
    subscriptionId: 'sub_002'
  },
  {
    id: 'INV-2024-003',
    tenant: 'RetailMax',
    amount: 29.00,
    status: 'Overdue',
    dueDate: '2024-01-10',
    paidDate: null,
    plan: 'Basic',
    billingPeriod: 'Jan 2024',
    subscriptionId: 'sub_003'
  },
  {
    id: 'INV-2024-004',
    tenant: 'GlobalTrade Inc',
    amount: 2988.00,
    status: 'Paid',
    dueDate: '2024-01-12',
    paidDate: '2024-01-11',
    plan: 'Enterprise Annual',
    billingPeriod: 'Jan 2024 - Dec 2024',
    subscriptionId: 'sub_004'
  },
  {
    id: 'INV-2024-005',
    tenant: 'EcoFriendly Co',
    amount: 99.00,
    status: 'Processing',
    subscriptionId: 'sub_005',
    dueDate: '2024-01-25',
    paidDate: null,
    plan: 'Professional',
    billingPeriod: 'Jan 2024'
  }
];

const payments = [
  {
    id: 'PAY-2024-001',
    tenant: 'TechCorp Solutions',
    amount: 299.00,
    status: 'Completed',
    method: 'Credit Card',
    gateway: 'Stripe',
    date: '2024-01-15T10:30:00Z',
    transactionId: 'txn_1234567890',
    currency: 'USD',
    fee: 8.97
  },
  {
    id: 'PAY-2024-002',
    tenant: 'StartupHub',
    amount: 99.00,
    status: 'Failed',
    method: 'Credit Card',
    gateway: 'Stripe',
    date: '2024-01-15T09:15:00Z',
    transactionId: 'txn_0987654321',
    currency: 'USD',
    fee: 0,
    failureReason: 'Insufficient funds'
  },
  {
    id: 'PAY-2024-003',
    tenant: 'RetailMax',
    amount: 29.00,
    status: 'Pending',
    method: 'Bank Transfer',
    gateway: 'Plaid',
    date: '2024-01-15T08:45:00Z',
    transactionId: 'txn_1122334455',
    currency: 'USD',
    fee: 0.87
  },
  {
    id: 'PAY-2024-004',
    tenant: 'GlobalTrade Inc',
    amount: 2988.00,
    status: 'Completed',
    method: 'Wire Transfer',
    gateway: 'Manual',
    date: '2024-01-14T16:20:00Z',
    transactionId: 'wire_9988776655',
    currency: 'USD',
    fee: 25.00
  },
  {
    id: 'PAY-2024-005',
    tenant: 'EcoFriendly Co',
    amount: 99.00,
    status: 'Refunded',
    method: 'Credit Card',
    gateway: 'Stripe',
    date: '2024-01-14T14:10:00Z',
    transactionId: 'txn_5544332211',
    currency: 'USD',
    fee: -2.97,
    refundReason: 'Customer request'
  }
];


export default function BillingHome() {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [methodFilter, setMethodFilter] = useState('all');

  const filteredInvoices = invoices.filter(invoice => {
    const matchesSearch = invoice.tenant.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         invoice.id.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || invoice.status.toLowerCase() === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const filteredPayments = payments.filter(payment => {
    const matchesSearch = payment.tenant.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         payment.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         payment.transactionId.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || payment.status.toLowerCase() === statusFilter;
    const matchesMethod = methodFilter === 'all' || payment.method.toLowerCase().includes(methodFilter.toLowerCase());
    return matchesSearch && matchesStatus && matchesMethod;
  });

  const getStatusBadge = (status) => {
    switch (status) {
      case 'Paid':
        return <Badge variant="default" className="bg-green-600">Paid</Badge>;
      case 'Pending':
        return <Badge variant="secondary">Pending</Badge>;
      case 'Overdue':
        return <Badge variant="destructive">Overdue</Badge>;
      case 'Processing':
        return <Badge variant="outline" className="border-blue-600 text-blue-600">Processing</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  const getPaymentStatusBadge = (status) => {
    switch (status) {
      case 'Completed':
        return <Badge variant="default" className="bg-green-600">Completed</Badge>;
      case 'Pending':
        return <Badge variant="secondary">Pending</Badge>;
      case 'Failed':
        return <Badge variant="destructive">Failed</Badge>;
      case 'Refunded':
        return <Badge variant="outline" className="border-blue-600 text-blue-600">Refunded</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };


  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const formatDateTime = (dateString) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center px-6">
          <div className="flex items-center gap-4">
            <CreditCard className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Billing & Payments</h1>
            <Badge variant="outline" className="ml-2">Admin Panel</Badge>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Stats Row */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 w-full">
          {billingStats.map((stat) => (
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

        {/* Billing & Payments Table */}
        <Card>
          <CardHeader className="py-4">
            {/* Filters and Search */}
            <div className="flex flex-wrap items-center gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  type="text"
                  placeholder="Search by ID, tenant name, or transaction ID..."
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
                  <SelectItem value="all">All Statuses</SelectItem>
                  <SelectItem value="paid">Paid</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="overdue">Overdue</SelectItem>
                  <SelectItem value="failed">Failed</SelectItem>
                  <SelectItem value="processing">Processing</SelectItem>
                </SelectContent>
              </Select>
              <Select value={methodFilter} onValueChange={setMethodFilter}>
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Filter by Type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Types</SelectItem>
                  <SelectItem value="invoice">Invoices</SelectItem>
                  <SelectItem value="payment">Payments</SelectItem>
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
                    <TableHead>Amount</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Date</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {/* Invoices */}
                  {filteredInvoices.map((invoice) => (
                    <TableRow key={invoice.id}>
                      <TableCell>
                        <Link
                          to={`/billing/${invoice.id}`}
                          className="text-primary hover:underline font-medium"
                        >
                          {invoice.id}
                        </Link>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <FileText className="h-4 w-4 text-muted-foreground" />
                          {invoice.tenant}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1">
                          <DollarSign className="h-4 w-4 text-muted-foreground" />
                          {formatCurrency(invoice.amount)}
                        </div>
                      </TableCell>
                      <TableCell>
                        {getStatusBadge(invoice.status)}
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">Invoice</Badge>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1">
                          <Calendar className="h-4 w-4 text-muted-foreground" />
                          {formatDate(invoice.dueDate)}
                        </div>
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button size="sm" variant="ghost">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => navigate(`/billing/${invoice.id}`)}>
                              <Eye className="h-4 w-4 mr-2" />
                              View Details
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                              <Download className="h-4 w-4 mr-2" />
                              Download PDF
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                              <Send className="h-4 w-4 mr-2" />
                              Send Reminder
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => navigate(`/billing/${invoice.id}/edit`)}>
                              <Edit className="h-4 w-4 mr-2" />
                              Edit Invoice
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))}

                  {/* Payments */}
                  {filteredPayments.map((payment) => (
                    <TableRow key={payment.id}>
                      <TableCell>
                        <Link
                          to={`/billing/${payment.id}`}
                          className="text-primary hover:underline font-medium"
                        >
                          {payment.id}
                        </Link>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <CreditCard className="h-4 w-4 text-muted-foreground" />
                          {payment.tenant}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1">
                          <DollarSign className="h-4 w-4 text-muted-foreground" />
                          {formatCurrency(payment.amount)}
                        </div>
                      </TableCell>
                      <TableCell>
                        {getPaymentStatusBadge(payment.status)}
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">Payment</Badge>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1">
                          <Calendar className="h-4 w-4 text-muted-foreground" />
                          {formatDateTime(payment.date)}
                        </div>
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button size="sm" variant="ghost">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => navigate(`/billing/${payment.id}`)}>
                              <Eye className="h-4 w-4 mr-2" />
                              View Details
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                              <FileText className="h-4 w-4 mr-2" />
                              Download Receipt
                            </DropdownMenuItem>
                            {payment.status === 'Completed' && (
                              <DropdownMenuItem>
                                <RotateCcw className="h-4 w-4 mr-2" />
                                Process Refund
                              </DropdownMenuItem>
                            )}
                            {payment.status === 'Failed' && (
                              <DropdownMenuItem>
                                <RefreshCw className="h-4 w-4 mr-2" />
                                Retry Payment
                              </DropdownMenuItem>
                            )}
                            <DropdownMenuItem onClick={() => navigate(`/billing/${payment.id}/edit`)}>
                              <Edit className="h-4 w-4 mr-2" />
                              Edit Payment
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
