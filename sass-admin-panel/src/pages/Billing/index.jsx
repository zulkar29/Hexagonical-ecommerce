import React, { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
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
  ArrowDownRight
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
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
    billingPeriod: 'Jan 2024'
  },
  {
    id: 'INV-2024-002',
    tenant: 'StartupHub',
    amount: 99.00,
    status: 'Pending',
    dueDate: '2024-01-20',
    paidDate: null,
    plan: 'Professional',
    billingPeriod: 'Jan 2024'
  },
  {
    id: 'INV-2024-003',
    tenant: 'RetailMax',
    amount: 29.00,
    status: 'Overdue',
    dueDate: '2024-01-10',
    paidDate: null,
    plan: 'Basic',
    billingPeriod: 'Jan 2024'
  },
  {
    id: 'INV-2024-004',
    tenant: 'GlobalTrade Inc',
    amount: 2988.00,
    status: 'Paid',
    dueDate: '2024-01-12',
    paidDate: '2024-01-11',
    plan: 'Enterprise Annual',
    billingPeriod: 'Jan 2024 - Dec 2024'
  },
  {
    id: 'INV-2024-005',
    tenant: 'EcoFriendly Co',
    amount: 99.00,
    status: 'Processing',
    dueDate: '2024-01-25',
    paidDate: null,
    plan: 'Professional',
    billingPeriod: 'Jan 2024'
  }
];

const recentTransactions = [
  { id: 'TXN-001', tenant: 'TechCorp Solutions', amount: 299.00, type: 'Payment', status: 'Completed', date: '2024-01-14' },
  { id: 'TXN-002', tenant: 'GlobalTrade Inc', amount: 2988.00, type: 'Payment', status: 'Completed', date: '2024-01-11' },
  { id: 'TXN-003', tenant: 'StartupHub', amount: 99.00, type: 'Refund', status: 'Processing', date: '2024-01-10' },
  { id: 'TXN-004', tenant: 'EcoFriendly Co', amount: 99.00, type: 'Payment', status: 'Failed', date: '2024-01-09' }
];

export default function BillingHome() {
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [activeTab, setActiveTab] = useState('invoices');

  const filteredInvoices = invoices.filter(invoice => {
    const matchesSearch = invoice.tenant.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         invoice.id.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || invoice.status.toLowerCase() === statusFilter;
    return matchesSearch && matchesStatus;
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

  const getTransactionStatusBadge = (status) => {
    switch (status) {
      case 'Completed':
        return <Badge variant="default" className="bg-green-600">Completed</Badge>;
      case 'Processing':
        return <Badge variant="outline" className="border-blue-600 text-blue-600">Processing</Badge>;
      case 'Failed':
        return <Badge variant="destructive">Failed</Badge>;
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

  return (
    <div className="space-y-6 p-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-3xl font-bold tracking-tight">Billing Management</h2>
              <p className="text-muted-foreground">
                Manage invoices, payments, and billing operations
              </p>
            </div>
            <div className="flex items-center space-x-3">
              <Button variant="outline">
                <Download className="h-4 w-4 mr-2" />
                Export
              </Button>
              <Link to="/billing/create">
                <Button>
                  <Plus className="h-4 w-4 mr-2" />
                  Create Invoice
                </Button>
              </Link>
            </div>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {billingStats.map((stat) => (
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

          {/* Main Content Tabs */}
          <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
            <TabsList className="grid w-full grid-cols-3">
              <TabsTrigger value="invoices">Invoices</TabsTrigger>
              <TabsTrigger value="transactions">Transactions</TabsTrigger>
              <TabsTrigger value="reports">Reports</TabsTrigger>
            </TabsList>

            <TabsContent value="invoices" className="space-y-6">
              {/* Filters and Search */}
              <Card>
                <CardContent className="p-6">
                  <div className="flex flex-col md:flex-row gap-4">
                    <div className="flex-1 relative">
                      <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                      <Input
                        type="text"
                        placeholder="Search invoices..."
                        className="pl-10"
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                      />
                    </div>
                    <div className="flex gap-2">
                      <Select value={statusFilter} onValueChange={setStatusFilter}>
                        <SelectTrigger className="w-40">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="all">All Status</SelectItem>
                          <SelectItem value="paid">Paid</SelectItem>
                          <SelectItem value="pending">Pending</SelectItem>
                          <SelectItem value="overdue">Overdue</SelectItem>
                          <SelectItem value="processing">Processing</SelectItem>
                        </SelectContent>
                      </Select>
                      <Button variant="outline">
                        <RefreshCw className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Invoices Table */}
              <Card>
                <CardHeader>
                  <CardTitle>Invoice Management</CardTitle>
                </CardHeader>
                <CardContent>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Invoice ID</TableHead>
                        <TableHead>Tenant</TableHead>
                        <TableHead>Amount</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Due Date</TableHead>
                        <TableHead>Plan</TableHead>
                        <TableHead>Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {filteredInvoices.map((invoice) => (
                        <TableRow key={invoice.id}>
                          <TableCell className="font-mono text-sm">{invoice.id}</TableCell>
                          <TableCell>{invoice.tenant}</TableCell>
                          <TableCell className="font-semibold">{formatCurrency(invoice.amount)}</TableCell>
                          <TableCell>{getStatusBadge(invoice.status)}</TableCell>
                          <TableCell>{formatDate(invoice.dueDate)}</TableCell>
                          <TableCell>
                            <Badge variant="outline">{invoice.plan}</Badge>
                          </TableCell>
                          <TableCell>
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="sm">
                                  <MoreHorizontal className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end">
                                <DropdownMenuItem>
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
                                <DropdownMenuItem asChild>
                                  <Link to={`/billing/${invoice.id}/edit`}>
                                    <Edit className="h-4 w-4 mr-2" />
                                    Edit Invoice
                                  </Link>
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
            </TabsContent>

            <TabsContent value="transactions" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Recent Transactions</CardTitle>
                </CardHeader>
                <CardContent>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Transaction ID</TableHead>
                        <TableHead>Tenant</TableHead>
                        <TableHead>Amount</TableHead>
                        <TableHead>Type</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Date</TableHead>
                        <TableHead>Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {recentTransactions.map((transaction) => (
                        <TableRow key={transaction.id}>
                          <TableCell className="font-mono text-sm">{transaction.id}</TableCell>
                          <TableCell>{transaction.tenant}</TableCell>
                          <TableCell className="font-semibold">{formatCurrency(transaction.amount)}</TableCell>
                          <TableCell>
                            <Badge variant={transaction.type === 'Payment' ? 'default' : 'secondary'}>
                              {transaction.type}
                            </Badge>
                          </TableCell>
                          <TableCell>{getTransactionStatusBadge(transaction.status)}</TableCell>
                          <TableCell>{formatDate(transaction.date)}</TableCell>
                          <TableCell>
                            <Button variant="ghost" size="sm">
                              <Eye className="h-4 w-4" />
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="reports" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Billing Reports</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <FileText className="h-6 w-6 mb-2" />
                      Revenue Report
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <Calendar className="h-6 w-6 mb-2" />
                      Monthly Summary
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <AlertCircle className="h-6 w-6 mb-2" />
                      Overdue Report
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <TrendingUp className="h-6 w-6 mb-2" />
                      Growth Analysis
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <DollarSign className="h-6 w-6 mb-2" />
                      Tax Report
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <Download className="h-6 w-6 mb-2" />
                      Export All
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
    </div>
  );
}
