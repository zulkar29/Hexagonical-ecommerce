import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
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
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

import {
  Search,
  Plus,
  Eye,
  CreditCard,
  Calendar,
  DollarSign,
  CheckCircle,
  Clock,
  Building2,
  Download,
  Banknote,
  Filter,
  AlertTriangle,
  TrendingUp,
  TrendingDown,
  Target,
  Lightbulb
} from 'lucide-react';
import { cn } from '@/lib/utils';

// Realistic mock data for accounts payable with smart calculations
const mockPayables = [
  {
    id: 'AP-001',
    supplier: 'ABC Wholesale Ltd.',
    invoiceNumber: 'INV-2024-001',
    invoiceDate: '2024-01-15',
    dueDate: '2024-02-15',
    amount: 45000,
    paidAmount: 0,
    remainingAmount: 45000,
    status: 'overdue',
    category: 'Inventory',
    description: 'Monthly inventory stock purchase',
    paymentTerms: 30,
    priority: 'high',
    daysOverdue: 15,
    lateFee: 2250, // 5% late fee
    contact: 'Mr. Karim Ahmed',
    phone: '01712345678'
  },
  {
    id: 'AP-002',
    supplier: 'Fresh Vegetables Co.',
    invoiceNumber: 'INV-2024-002',
    invoiceDate: '2024-01-20',
    dueDate: '2024-02-20',
    amount: 12500,
    paidAmount: 5000,
    remainingAmount: 7500,
    status: 'partial',
    category: 'Supplies',
    description: 'Weekly fresh vegetables supply',
    paymentTerms: 30,
    priority: 'medium',
    daysOverdue: 10,
    lateFee: 375, // 5% late fee on remaining
    contact: 'Ms. Fatima Begum',
    phone: '01823456789'
  },
  {
    id: 'AP-003',
    supplier: 'Tech Solutions BD',
    invoiceNumber: 'INV-2024-003',
    invoiceDate: '2024-01-25',
    dueDate: '2024-02-25',
    amount: 25000,
    paidAmount: 25000,
    remainingAmount: 0,
    status: 'paid',
    category: 'Technology',
    description: 'POS system maintenance and upgrade',
    paymentTerms: 30,
    priority: 'low',
    daysOverdue: 0,
    lateFee: 0,
    contact: 'Mr. Rahman Sheikh',
    phone: '01934567890'
  },
  {
    id: 'AP-004',
    supplier: 'Office Supplies Plus',
    invoiceNumber: 'INV-2024-004',
    invoiceDate: '2024-01-28',
    dueDate: '2024-02-28',
    amount: 8500,
    paidAmount: 0,
    remainingAmount: 8500,
    status: 'current',
    category: 'Office',
    description: 'Monthly office supplies and stationery',
    paymentTerms: 30,
    priority: 'low',
    daysOverdue: 0,
    lateFee: 0,
    contact: 'Mr. Nasir Uddin',
    phone: '01645678901'
  },
  {
    id: 'AP-005',
    supplier: 'Utility Services Ltd.',
    invoiceNumber: 'UTIL-2024-001',
    invoiceDate: '2024-01-30',
    dueDate: '2024-02-15',
    amount: 15000,
    paidAmount: 0,
    remainingAmount: 15000,
    status: 'due_soon',
    category: 'Utilities',
    description: 'Monthly electricity and water bills',
    paymentTerms: 15,
    priority: 'high',
    daysOverdue: 0,
    lateFee: 0,
    contact: 'Customer Service',
    phone: '01756789012'
  }
];

// Enhanced financial calculations with smart insights
const calculateFinancialMetrics = (payables) => {
  const totalAmount = payables.reduce((sum, p) => sum + p.amount, 0);
  const totalPaid = payables.reduce((sum, p) => sum + p.paidAmount, 0);
  const totalOutstanding = payables.reduce((sum, p) => sum + p.remainingAmount, 0);
  const overdueAmount = payables.filter(p => p.status === 'overdue').reduce((sum, p) => sum + p.remainingAmount, 0);
  
  // Smart insights and automated calculations
  const overdueCount = payables.filter(p => p.status === 'overdue').length;
  const dueSoonAmount = payables.filter(p => p.status === 'due_soon').reduce((sum, p) => sum + p.remainingAmount, 0);
  const dueSoonCount = payables.filter(p => p.status === 'due_soon').length;
  const paymentRate = totalAmount > 0 ? ((totalPaid / totalAmount) * 100).toFixed(1) : 0;
  const overduePercentage = totalOutstanding > 0 ? ((overdueAmount / totalOutstanding) * 100).toFixed(1) : 0;
  
  // Cash flow predictions
  const avgPaymentTerm = payables.length > 0 ? 
    payables.reduce((sum, p) => sum + p.paymentTerms, 0) / payables.length : 30;
  
  // Supplier performance insights
  const supplierCount = new Set(payables.map(p => p.supplier)).size;
  const avgInvoiceAmount = payables.length > 0 ? totalAmount / payables.length : 0;
  
  // Risk assessment
  const highRiskAmount = payables.filter(p => p.status === 'overdue' && p.remainingAmount > 20000)
    .reduce((sum, p) => sum + p.remainingAmount, 0);
  
  return {
    totalAmount,
    totalPaid,
    totalOutstanding,
    overdueAmount,
    overdueCount,
    dueSoonAmount,
    dueSoonCount,
    paymentRate,
    overduePercentage,
    avgPaymentTerm,
    supplierCount,
    avgInvoiceAmount,
    highRiskAmount
  };
};

// Smart recommendations engine
const generateSmartInsights = (metrics, payables) => {
  const insights = [];
  
  // Cash flow insights
  if (metrics.overduePercentage > 20) {
    insights.push({
      type: 'warning',
      icon: AlertTriangle,
      title: 'High Overdue Rate',
      message: `${metrics.overduePercentage}% of outstanding payments are overdue. Consider prioritizing collections.`,
      action: 'Review overdue payments'
    });
  }
  
  if (metrics.dueSoonAmount > 50000) {
    insights.push({
      type: 'info',
      icon: Calendar,
      title: 'Upcoming Payments',
      message: `৳${metrics.dueSoonAmount.toLocaleString()} in payments due soon. Plan cash flow accordingly.`,
      action: 'Schedule payments'
    });
  }
  
  if (metrics.paymentRate > 80) {
    insights.push({
      type: 'success',
      icon: TrendingUp,
      title: 'Good Payment Performance',
      message: `${metrics.paymentRate}% payment rate shows healthy cash management.`,
      action: 'Maintain current practices'
    });
  }
  
  // Supplier relationship insights
  if (metrics.supplierCount < 5) {
    insights.push({
      type: 'info',
      icon: Building2,
      title: 'Supplier Diversification',
      message: 'Consider diversifying suppliers to reduce dependency risk.',
      action: 'Explore new suppliers'
    });
  }
  
  // Cost optimization insights
  if (metrics.avgInvoiceAmount > 30000) {
    insights.push({
      type: 'info',
      icon: Target,
      title: 'Bulk Purchase Opportunity',
      message: 'High average invoice amounts suggest good bulk purchasing power.',
      action: 'Negotiate better terms'
    });
  }
  
  return insights.slice(0, 3); // Show top 3 insights
};

const AccountsPayable = () => {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [selectedPayable, setSelectedPayable] = useState(null);
  const [isPaymentDialogOpen, setIsPaymentDialogOpen] = useState(false);
  const [paymentAmount, setPaymentAmount] = useState('');
  const [paymentMethod, setPaymentMethod] = useState('bank_transfer');
  const [paymentNotes, setPaymentNotes] = useState('');

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-GB');
  };

  const metrics = calculateFinancialMetrics(mockPayables);
  const smartInsights = generateSmartInsights(metrics, mockPayables);

  const getStatusBadge = (status) => {
    const statusConfig = {
      paid: { variant: 'default', label: 'Paid', icon: CheckCircle, color: 'text-green-600' },
      pending: { variant: 'secondary', label: 'Pending', icon: Clock, color: 'text-yellow-600' },
      overdue: { variant: 'destructive', label: 'Overdue', icon: AlertTriangle, color: 'text-red-600' },
      partial: { variant: 'outline', label: 'Partial', icon: CreditCard, color: 'text-blue-600' },
      due_soon: { variant: 'outline', label: 'Due Soon', icon: Calendar, color: 'text-orange-600' }
    };
    
    const config = statusConfig[status] || statusConfig.pending;
    const Icon = config.icon;
    
    return (
      <Badge variant={config.variant} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    );
  };



  const filteredPayables = mockPayables.filter(payable => {
    const matchesSearch = payable.supplier.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         payable.invoiceNumber.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         payable.description.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || payable.status === statusFilter;
    return matchesSearch && matchesStatus;
  }).sort((a, b) => {
    // Simple sorting: overdue first, then by due date
    if (a.status === 'overdue' && b.status !== 'overdue') return -1;
    if (b.status === 'overdue' && a.status !== 'overdue') return 1;
    return new Date(a.dueDate) - new Date(b.dueDate);
  });

  const handleMakePayment = () => {
    // Handle payment logic here
    console.log('Making payment:', {
      payableId: selectedPayable?.id,
      amount: paymentAmount,
      method: paymentMethod,
      notes: paymentNotes
    });
    setIsPaymentDialogOpen(false);
    setPaymentAmount('');
    setPaymentMethod('');
    setPaymentNotes('');
  };

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Accounts Payable</h1>
          <p className="text-muted-foreground">
            Manage supplier payments and outstanding invoices
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm">
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Button size="sm">
            <Plus className="h-4 w-4 mr-2" />
            Add Payable
          </Button>
        </div>
      </div>

      {/* Smart Insights */}
      {smartInsights.length > 0 && (
        <Card className="border-l-4 border-l-blue-500">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Lightbulb className="h-5 w-5 text-blue-500" />
              Smart Insights
            </CardTitle>
            <CardDescription>
              AI-powered recommendations for better financial management
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {smartInsights.map((insight, index) => {
                const Icon = insight.icon;
                const bgColor = {
                  warning: 'bg-red-50 border-red-200',
                  info: 'bg-blue-50 border-blue-200',
                  success: 'bg-green-50 border-green-200'
                }[insight.type];
                const iconColor = {
                  warning: 'text-red-500',
                  info: 'text-blue-500',
                  success: 'text-green-500'
                }[insight.type];
                
                return (
                  <div key={index} className={`p-3 rounded-lg border ${bgColor}`}>
                    <div className="flex items-start gap-3">
                      <Icon className={`h-5 w-5 mt-0.5 ${iconColor}`} />
                      <div className="flex-1">
                        <h4 className="font-medium text-sm">{insight.title}</h4>
                        <p className="text-sm text-muted-foreground mt-1">{insight.message}</p>
                        <Button variant="link" size="sm" className="p-0 h-auto mt-1 text-xs">
                          {insight.action} →
                        </Button>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Enhanced Financial Summary */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="hover:shadow-md transition-shadow">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <DollarSign className="h-8 w-8 text-blue-600" />
                <div>
                  <p className="text-sm text-muted-foreground">Total Outstanding</p>
                  <p className="text-2xl font-bold">{formatCurrency(metrics.totalOutstanding)}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-xs text-muted-foreground">Payment Rate</p>
                <p className="text-sm font-semibold text-green-600">{metrics.paymentRate}%</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card className="hover:shadow-md transition-shadow">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <AlertTriangle className="h-8 w-8 text-red-600" />
                <div>
                  <p className="text-sm text-muted-foreground">Overdue Amount</p>
                  <p className="text-2xl font-bold text-red-600">{formatCurrency(metrics.overdueAmount)}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-xs text-muted-foreground">{metrics.overdueCount} invoices</p>
                <p className="text-sm font-semibold text-red-600">{metrics.overduePercentage}% of total</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card className="hover:shadow-md transition-shadow">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Calendar className="h-8 w-8 text-orange-600" />
                <div>
                  <p className="text-sm text-muted-foreground">Due Soon</p>
                  <p className="text-2xl font-bold text-orange-600">{formatCurrency(metrics.dueSoonAmount)}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-xs text-muted-foreground">{metrics.dueSoonCount} invoices</p>
                <p className="text-sm font-semibold text-orange-600">Next 7 days</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card className="hover:shadow-md transition-shadow">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Building2 className="h-8 w-8 text-purple-600" />
                <div>
                  <p className="text-sm text-muted-foreground">Active Suppliers</p>
                  <p className="text-2xl font-bold">{metrics.supplierCount}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-xs text-muted-foreground">Avg Invoice</p>
                <p className="text-sm font-semibold text-purple-600">{formatCurrency(metrics.avgInvoiceAmount)}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
      
      {/* Additional Smart Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="border-l-4 border-l-green-500">
          <CardContent className="p-4">
            <div className="flex items-center gap-3">
              <TrendingUp className="h-6 w-6 text-green-600" />
              <div>
                <p className="text-sm text-muted-foreground">Payment Performance</p>
                <p className="text-lg font-bold text-green-600">{metrics.paymentRate}% Success Rate</p>
                <p className="text-xs text-muted-foreground">Avg payment term: {Math.round(metrics.avgPaymentTerm)} days</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card className="border-l-4 border-l-blue-500">
          <CardContent className="p-4">
            <div className="flex items-center gap-3">
              <Target className="h-6 w-6 text-blue-600" />
              <div>
                <p className="text-sm text-muted-foreground">Cash Flow Health</p>
                <p className="text-lg font-bold">
                  {metrics.overduePercentage < 10 ? 'Excellent' : 
                   metrics.overduePercentage < 20 ? 'Good' : 
                   metrics.overduePercentage < 30 ? 'Fair' : 'Poor'}
                </p>
                <p className="text-xs text-muted-foreground">{metrics.overduePercentage}% overdue ratio</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card className="border-l-4 border-l-yellow-500">
          <CardContent className="p-4">
            <div className="flex items-center gap-3">
              <AlertTriangle className="h-6 w-6 text-yellow-600" />
              <div>
                <p className="text-sm text-muted-foreground">High Risk Exposure</p>
                <p className="text-lg font-bold text-yellow-600">{formatCurrency(metrics.highRiskAmount)}</p>
                <p className="text-xs text-muted-foreground">Large overdue amounts (&gt;৳20K)</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  placeholder="Search by supplier or invoice number..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[180px]">
                <Filter className="h-4 w-4 mr-2" />
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="overdue">Overdue</SelectItem>
                <SelectItem value="due_soon">Due Soon</SelectItem>
                <SelectItem value="partial">Partial</SelectItem>
                <SelectItem value="paid">Paid</SelectItem>
              </SelectContent>
            </Select>

          </div>
        </CardContent>
      </Card>

      {/* Payables Table */}
      <Card>
        <CardHeader>
          <CardTitle>Outstanding Payables</CardTitle>
          <CardDescription>
            {filteredPayables.length} payables found
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Supplier</TableHead>
                <TableHead>Invoice #</TableHead>
                <TableHead>Due Date</TableHead>
                <TableHead>Amount</TableHead>
                <TableHead>Remaining</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredPayables.map((payable) => (
                <TableRow key={payable.id}>
                  <TableCell>
                    <div>
                      <p className="font-medium">{payable.supplier}</p>
                      <p className="text-sm text-muted-foreground">{payable.category}</p>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="font-medium">{payable.invoiceNumber}</p>
                      <p className="text-sm text-muted-foreground">
                        {formatDate(payable.invoiceDate)}
                      </p>
                    </div>
                  </TableCell>
                  <TableCell className="font-medium">
                    {formatDate(payable.dueDate)}
                  </TableCell>
                  <TableCell className="font-medium">
                    {formatCurrency(payable.amount)}
                  </TableCell>
                  <TableCell className="font-medium">
                    {formatCurrency(payable.remainingAmount)}
                  </TableCell>
                  <TableCell>
                    {getStatusBadge(payable.status)}
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Button variant="outline" size="sm">
                        <Eye className="h-4 w-4" />
                      </Button>
                      {payable.remainingAmount > 0 && (
                        <Button 
                          size="sm"
                          onClick={() => {
                            setSelectedPayable(payable);
                            setPaymentAmount(payable.remainingAmount.toString());
                            setIsPaymentDialogOpen(true);
                          }}
                        >
                          <Banknote className="h-4 w-4 mr-1" />
                          Pay
                        </Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Payment Dialog */}
      <Dialog open={isPaymentDialogOpen} onOpenChange={setIsPaymentDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Make Payment</DialogTitle>
            <DialogDescription>
              Record payment for {selectedPayable?.supplier}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="space-y-2">
              <Label>Invoice Details</Label>
              <div className="p-3 bg-muted rounded-lg">
                <p className="text-sm font-medium">{selectedPayable?.invoiceNumber}</p>
                <p className="text-sm text-muted-foreground">
                  Due: {selectedPayable && formatDate(selectedPayable.dueDate)}
                </p>
                <p className="text-sm text-muted-foreground">
                  Outstanding: {selectedPayable && formatCurrency(selectedPayable.remainingAmount)}
                </p>
              </div>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="amount" className="text-right">
                Amount
              </Label>
              <Input
                id="amount"
                type="number"
                value={paymentAmount}
                onChange={(e) => setPaymentAmount(e.target.value)}
                className="col-span-3"
                placeholder="Enter payment amount"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="method" className="text-right">
                Method
              </Label>
              <Select value={paymentMethod} onValueChange={setPaymentMethod}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select payment method" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="cash">Cash</SelectItem>
                  <SelectItem value="bank_transfer">Bank Transfer</SelectItem>
                  <SelectItem value="bkash">bKash</SelectItem>
                  <SelectItem value="nagad">Nagad</SelectItem>
                  <SelectItem value="check">Check</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="notes" className="text-right">
                Notes
              </Label>
              <Textarea
                id="notes"
                value={paymentNotes}
                onChange={(e) => setPaymentNotes(e.target.value)}
                className="col-span-3"
                placeholder="Payment notes (optional)"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsPaymentDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleMakePayment}>
              Record Payment
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default AccountsPayable;