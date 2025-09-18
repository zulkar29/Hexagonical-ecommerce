import React, { useState } from 'react';
import {
  CreditCard,
  DollarSign,
  TrendingUp,
  AlertTriangle,
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
  ArrowUpRight,
  ArrowDownLeft,
  Shield,
  Eye,
  Ban,
  RotateCcw,
  FileText
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
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
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

const paymentStats = [
  { label: 'Total Processed', value: '$3,247,890', icon: DollarSign, color: 'text-green-600', change: '+18.2%' },
  { label: 'Success Rate', value: '98.7%', icon: CheckCircle, color: 'text-blue-600', change: '+0.3%' },
  { label: 'Failed Payments', value: 47, icon: XCircle, color: 'text-red-600', change: '-12.5%' },
  { label: 'Pending Review', value: 23, icon: Clock, color: 'text-yellow-600', change: '+5.1%' }
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

const paymentMethods = [
  { name: 'Credit/Debit Cards', count: 1247, percentage: 78.5, icon: CreditCard },
  { name: 'Bank Transfers', count: 234, percentage: 14.7, icon: ArrowDownLeft },
  { name: 'Wire Transfers', count: 89, percentage: 5.6, icon: ArrowUpRight },
  { name: 'Digital Wallets', count: 18, percentage: 1.2, icon: Shield }
];

const fraudAlerts = [
  { id: 'FRD-001', tenant: 'SuspiciousUser', amount: 5000, reason: 'Unusual transaction amount', severity: 'High', date: '2024-01-15' },
  { id: 'FRD-002', tenant: 'TestAccount', amount: 1, reason: 'Multiple failed attempts', severity: 'Medium', date: '2024-01-15' },
  { id: 'FRD-003', tenant: 'NewTenant', amount: 999, reason: 'First-time large payment', severity: 'Low', date: '2024-01-14' }
];

export default function PaymentsHome() {
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [methodFilter, setMethodFilter] = useState('all');
  const [activeTab, setActiveTab] = useState('payments');

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

  const getRiskBadge = (risk) => {
    switch (risk) {
      case 'High':
        return <Badge variant="destructive">High</Badge>;
      case 'Medium':
        return <Badge variant="secondary">Medium</Badge>;
      case 'Low':
        return <Badge variant="default" className="bg-green-600">Low</Badge>;
      default:
        return <Badge variant="outline">{risk}</Badge>;
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
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
              <h2 className="text-3xl font-bold tracking-tight">Payment Management</h2>
              <p className="text-muted-foreground">
                Monitor transactions, process refunds, and manage payment security
              </p>
            </div>
            <div className="flex items-center space-x-3">
              <Button variant="outline">
                <Download className="h-4 w-4 mr-2" />
                Export Report
              </Button>
              <Button>
                <Plus className="h-4 w-4 mr-2" />
                Manual Payment
              </Button>
            </div>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {paymentStats.map((stat) => (
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
            <TabsList className="grid w-full grid-cols-4">
              <TabsTrigger value="payments">Payments</TabsTrigger>
              <TabsTrigger value="methods">Payment Methods</TabsTrigger>
              <TabsTrigger value="fraud">Fraud Detection</TabsTrigger>
              <TabsTrigger value="analytics">Analytics</TabsTrigger>
            </TabsList>

            <TabsContent value="payments" className="space-y-6">
              {/* Filters and Search */}
              <Card>
                <CardContent className="p-6">
                  <div className="flex flex-col md:flex-row gap-4">
                    <div className="flex-1 relative">
                      <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                      <Input
                        type="text"
                        placeholder="Search payments..."
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
                          <SelectItem value="completed">Completed</SelectItem>
                          <SelectItem value="pending">Pending</SelectItem>
                          <SelectItem value="failed">Failed</SelectItem>
                          <SelectItem value="refunded">Refunded</SelectItem>
                        </SelectContent>
                      </Select>
                      <Select value={methodFilter} onValueChange={setMethodFilter}>
                        <SelectTrigger className="w-40">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="all">All Methods</SelectItem>
                          <SelectItem value="credit">Credit Card</SelectItem>
                          <SelectItem value="bank">Bank Transfer</SelectItem>
                          <SelectItem value="wire">Wire Transfer</SelectItem>
                        </SelectContent>
                      </Select>
                      <Button variant="outline">
                        <RefreshCw className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Payments Table */}
              <Card>
                <CardHeader>
                  <CardTitle>Payment Transactions</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="overflow-x-auto">
                    <Table>
                      <TableHeader>
                        <TableRow className="border-b">
                          <TableHead className="text-left p-4">Payment ID</TableHead>
                          <TableHead className="text-left p-4">Tenant</TableHead>
                          <TableHead className="text-left p-4">Amount</TableHead>
                          <TableHead className="text-left p-4">Status</TableHead>
                          <TableHead className="text-left p-4">Method</TableHead>
                          <TableHead className="text-left p-4">Gateway</TableHead>
                          <TableHead className="text-left p-4">Date</TableHead>
                          <TableHead className="text-left p-4">Actions</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {filteredPayments.map((payment) => (
                          <TableRow key={payment.id} className="border-b hover:bg-muted/50">
                            <TableCell className="p-4 font-mono text-sm">{payment.id}</TableCell>
                            <TableCell className="p-4">{payment.tenant}</TableCell>
                            <TableCell className="p-4 font-semibold">{formatCurrency(payment.amount)}</TableCell>
                            <TableCell className="p-4">{getStatusBadge(payment.status)}</TableCell>
                            <TableCell className="p-4">
                              <Badge variant="outline">{payment.method}</Badge>
                            </TableCell>
                            <TableCell className="p-4 text-sm text-muted-foreground">{payment.gateway}</TableCell>
                            <TableCell className="p-4 text-sm">{formatDateTime(payment.date)}</TableCell>
                            <TableCell className="p-4">
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
                                  <DropdownMenuItem>
                                    <FileText className="h-4 w-4 mr-2" />
                                    Download Receipt
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
            </TabsContent>

            <TabsContent value="methods" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Payment Method Distribution</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {paymentMethods.map((method) => (
                      <div key={method.name} className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center gap-4">
                          <method.icon className="h-6 w-6 text-muted-foreground" />
                          <div>
                            <p className="font-medium">{method.name}</p>
                            <p className="text-sm text-muted-foreground">{method.count} transactions</p>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="font-semibold">{method.percentage}%</p>
                          <div className="w-24 h-2 bg-muted rounded-full mt-1">
                            <div 
                              className="h-full bg-primary rounded-full" 
                              style={{ width: `${method.percentage}%` }}
                            />
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="fraud" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Fraud Detection Alerts</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="overflow-x-auto">
                    <Table>
                      <TableHeader>
                        <TableRow className="border-b">
                          <TableHead className="text-left p-4">Alert ID</TableHead>
                          <TableHead className="text-left p-4">Tenant</TableHead>
                          <TableHead className="text-left p-4">Amount</TableHead>
                          <TableHead className="text-left p-4">Reason</TableHead>
                          <TableHead className="text-left p-4">Severity</TableHead>
                          <TableHead className="text-left p-4">Date</TableHead>
                          <TableHead className="text-left p-4">Actions</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {fraudAlerts.map((alert) => (
                          <TableRow key={alert.id} className="border-b hover:bg-muted/50">
                            <TableCell className="p-4 font-mono text-sm">{alert.id}</TableCell>
                            <TableCell className="p-4">{alert.tenant}</TableCell>
                            <TableCell className="p-4 font-semibold">{formatCurrency(alert.amount)}</TableCell>
                            <TableCell className="p-4">{alert.reason}</TableCell>
                            <TableCell className="p-4">test</TableCell>
                            <TableCell className="p-4">{formatDate(alert.date)}</TableCell>
                            <TableCell className="p-4">
                              <div className="flex gap-2">
                                <Button variant="outline" size="sm">
                                  <Eye className="h-4 w-4" />
                                </Button>
                                <Button variant="outline" size="sm">
                                  <Ban className="h-4 w-4" />
                                </Button>
                              </div>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="analytics" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Payment Analytics</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <TrendingUp className="h-6 w-6 mb-2" />
                      Revenue Trends
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <Calendar className="h-6 w-6 mb-2" />
                      Monthly Reports
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <AlertTriangle className="h-6 w-6 mb-2" />
                      Failure Analysis
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <Shield className="h-6 w-6 mb-2" />
                      Security Report
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <DollarSign className="h-6 w-6 mb-2" />
                      Fee Analysis
                    </Button>
                    <Button variant="outline" className="h-24 flex flex-col items-center justify-center">
                      <Download className="h-6 w-6 mb-2" />
                      Export Data
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
    </div>
  );
}
