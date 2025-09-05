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
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Search,
  Filter,
  Plus,
  Download,
  Eye,
  CreditCard,
  Calendar,
  DollarSign,
  AlertTriangle,
  CheckCircle,
  Clock,
  FileText,
  Banknote,
  Send,
  Phone,
  TrendingUp,
  TrendingDown,
  Target,
  Lightbulb,
  Users,
  BarChart3
} from 'lucide-react';
import { cn } from '@/lib/utils';

// Mock data for accounts receivable
const mockReceivables = [
  {
    id: 'AR-001',
    customer: 'Fatima Khan',
    invoiceNumber: 'INV-2024-101',
    invoiceDate: '2024-01-10',
    dueDate: '2024-02-10',
    amount: 15000,
    paidAmount: 0,
    remainingAmount: 15000,
    status: 'overdue',
    orderNumber: 'ORD-1234',
    description: 'Bulk grocery purchase for restaurant',
    paymentTerms: 30,
    priority: 'high',
    daysOverdue: 20,
    lateFee: 750, // 5% late fee
    customerRisk: 'medium',
    lastPaymentDate: '2023-12-15',
    phone: '01712345678',
    email: 'fatima.khan@email.com',
    address: 'House 45, Road 12, Dhanmondi, Dhaka'
  },
  {
    id: 'AR-002',
    customer: 'Ahmed Ali',
    invoiceNumber: 'INV-2024-102',
    invoiceDate: '2024-01-15',
    dueDate: '2024-02-15',
    amount: 8500,
    paidAmount: 3000,
    remainingAmount: 5500,
    status: 'partial',
    orderNumber: 'ORD-1235',
    description: 'Monthly household supplies',
    paymentTerms: 30,
    priority: 'medium',
    daysOverdue: 15,
    lateFee: 275, // 5% late fee on remaining
    customerRisk: 'low',
    lastPaymentDate: '2024-01-20',
    phone: '01823456789',
    email: 'ahmed.ali@email.com',
    address: 'Flat 8B, Gulshan Avenue, Gulshan-1, Dhaka'
  },
  {
    id: 'AR-003',
    customer: 'Rashida Begum',
    invoiceNumber: 'INV-2024-103',
    invoiceDate: '2024-01-20',
    dueDate: '2024-02-20',
    amount: 12000,
    paidAmount: 12000,
    remainingAmount: 0,
    status: 'paid',
    orderNumber: 'ORD-1236',
    description: 'Premium grocery items',
    paymentTerms: 30,
    priority: 'low',
    daysOverdue: 0,
    lateFee: 0,
    customerRisk: 'low',
    lastPaymentDate: '2024-02-18',
    phone: '01934567890',
    email: 'rashida.begum@email.com',
    address: 'Block C, House 23, Uttara Sector 7, Dhaka'
  },
  {
    id: 'AR-004',
    customer: 'Mohammad Hasan',
    invoiceNumber: 'INV-2024-104',
    invoiceDate: '2024-01-25',
    dueDate: '2024-02-25',
    amount: 6500,
    paidAmount: 0,
    remainingAmount: 6500,
    status: 'pending',
    orderNumber: 'ORD-1237',
    description: 'Weekly family groceries',
    paymentTerms: 30,
    priority: 'medium',
    daysOverdue: 5,
    lateFee: 325, // 5% late fee
    customerRisk: 'medium',
    lastPaymentDate: '2024-01-10',
    phone: '01645678901',
    email: 'hasan.mohammad@email.com',
    address: 'Road 5, Block B, Mirpur-10, Dhaka'
  },
  {
    id: 'AR-005',
    customer: 'Nasir Ahmed',
    invoiceNumber: 'INV-2024-105',
    invoiceDate: '2024-01-28',
    dueDate: '2024-02-12',
    amount: 9500,
    paidAmount: 0,
    remainingAmount: 9500,
    status: 'due_soon',
    orderNumber: 'ORD-1238',
    description: 'Premium brand products',
    paymentTerms: 15,
    priority: 'high',
    daysOverdue: 0,
    lateFee: 0,
    customerRisk: 'low',
    lastPaymentDate: '2024-01-15',
    phone: '01756789012',
    email: 'nasir.ahmed@email.com',
    address: 'House 78, New Eskaton Road, Ramna, Dhaka'
  },
  {
    id: 'AR-006',
    customer: 'Salma Khatun',
    invoiceNumber: 'INV-2024-106',
    invoiceDate: '2024-02-01',
    dueDate: '2024-03-01',
    amount: 4500,
    paidAmount: 0,
    remainingAmount: 4500,
    status: 'current',
    orderNumber: 'ORD-1239',
    description: 'Regular household items',
    paymentTerms: 30,
    priority: 'low',
    daysOverdue: 0,
    lateFee: 0,
    customerRisk: 'low',
    lastPaymentDate: '2024-01-25',
    phone: '01867890123',
    email: 'salma.khatun@email.com',
    address: 'Lane 3, Wari, Old Dhaka'
  }
];

// Calculate smart financial metrics
const calculateFinancialMetrics = (receivables) => {
  const totalAmount = receivables.reduce((sum, r) => sum + r.amount, 0);
  const totalReceived = receivables.reduce((sum, r) => sum + r.paidAmount, 0);
  const totalOutstanding = receivables.reduce((sum, r) => sum + r.remainingAmount, 0);
  const overdueAmount = receivables.filter(r => r.status === 'overdue').reduce((sum, r) => sum + r.remainingAmount, 0);
  
  // Enhanced smart insights and automated calculations
  const overdueCount = receivables.filter(r => r.status === 'overdue').length;
  const dueSoonAmount = receivables.filter(r => r.status === 'due_soon').reduce((sum, r) => sum + r.remainingAmount, 0);
  const dueSoonCount = receivables.filter(r => r.status === 'due_soon').length;
  const collectionRate = totalAmount > 0 ? ((totalReceived / totalAmount) * 100).toFixed(1) : 0;
  const overduePercentage = totalOutstanding > 0 ? ((overdueAmount / totalOutstanding) * 100).toFixed(1) : 0;
  
  // Customer insights
  const customerCount = new Set(receivables.map(r => r.customer)).size;
  const avgInvoiceAmount = receivables.length > 0 ? totalAmount / receivables.length : 0;
  const avgPaymentTerm = receivables.length > 0 ? 
    receivables.reduce((sum, r) => sum + r.paymentTerms, 0) / receivables.length : 30;
  
  // Risk assessment
  const highRiskCustomers = receivables.filter(r => r.customerRisk === 'high' || 
    (r.status === 'overdue' && r.remainingAmount > 10000)).length;
  const highRiskAmount = receivables.filter(r => r.status === 'overdue' && r.remainingAmount > 10000)
    .reduce((sum, r) => sum + r.remainingAmount, 0);
  
  // Payment behavior analysis
  const partialPayments = receivables.filter(r => r.status === 'partial').length;
  const onTimePayments = receivables.filter(r => r.status === 'paid' && r.daysOverdue === 0).length;
  
  return {
    totalAmount,
    totalReceived,
    totalOutstanding,
    overdueAmount,
    overdueCount,
    dueSoonAmount,
    dueSoonCount,
    collectionRate,
    overduePercentage,
    customerCount,
    avgInvoiceAmount,
    avgPaymentTerm,
    highRiskCustomers,
    highRiskAmount,
    partialPayments,
    onTimePayments
  };
};

// Utility functions for date calculations
const getDaysOverdue = (dueDate) => {
  const today = new Date();
  const due = new Date(dueDate);
  const diffTime = today - due;
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  return Math.max(0, diffDays);
};

const getDaysToDue = (dueDate) => {
  const today = new Date();
  const due = new Date(dueDate);
  const diffTime = due - today;
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  return Math.max(0, diffDays);
};

// Smart recommendations engine for receivables
const generateReceivableInsights = (metrics, receivables) => {
  const insights = [];
  
  // Collection performance insights
  if (metrics.overduePercentage > 25) {
    insights.push({
      type: 'warning',
      icon: AlertTriangle,
      title: 'High Overdue Rate',
      message: `${metrics.overduePercentage}% of receivables are overdue. Implement aggressive collection strategies.`,
      action: 'Review collection process'
    });
  }
  
  if (metrics.collectionRate > 85) {
    insights.push({
      type: 'success',
      icon: TrendingUp,
      title: 'Excellent Collection Rate',
      message: `${metrics.collectionRate}% collection rate indicates strong customer relationships.`,
      action: 'Maintain current practices'
    });
  }
  
  if (metrics.dueSoonAmount > 30000) {
    insights.push({
      type: 'info',
      icon: Calendar,
      title: 'Upcoming Collections',
      message: `৳${metrics.dueSoonAmount.toLocaleString()} due soon. Send proactive reminders to customers.`,
      action: 'Send payment reminders'
    });
  }
  
  // Customer relationship insights
  if (metrics.highRiskCustomers > 0) {
    insights.push({
      type: 'warning',
      icon: Users,
      title: 'High-Risk Customers',
      message: `${metrics.highRiskCustomers} customers require special attention for collections.`,
      action: 'Review customer credit terms'
    });
  }
  
  // Business growth insights
  if (metrics.avgInvoiceAmount > 8000) {
    insights.push({
      type: 'success',
      icon: BarChart3,
      title: 'Strong Average Sales',
      message: `High average invoice amount (৳${metrics.avgInvoiceAmount.toLocaleString()}) shows good customer value.`,
      action: 'Focus on customer retention'
    });
  }
  
  return insights.slice(0, 3); // Show top 3 insights
};

const AccountsReceivable = () => {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');

  const [selectedReceivable, setSelectedReceivable] = useState(null);
  const [isPaymentDialogOpen, setIsPaymentDialogOpen] = useState(false);
  const [isReminderDialogOpen, setIsReminderDialogOpen] = useState(false);
  const [paymentAmount, setPaymentAmount] = useState('');
  const [paymentMethod, setPaymentMethod] = useState('bank_transfer');
  const [paymentNotes, setPaymentNotes] = useState('');
  const [reminderMessage, setReminderMessage] = useState('');

  // Calculate smart financial metrics
  const metrics = calculateFinancialMetrics(mockReceivables);
  const smartInsights = generateReceivableInsights(metrics, mockReceivables);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-GB');
  };

  const getStatusBadge = (status) => {
    const statusConfig = {
      paid: { variant: 'default', label: 'Paid', icon: CheckCircle, color: 'text-green-600' },
      pending: { variant: 'secondary', label: 'Pending', icon: Clock, color: 'text-yellow-600' },
      overdue: { variant: 'destructive', label: 'Overdue', icon: AlertTriangle, color: 'text-red-600' },
      partial: { variant: 'outline', label: 'Partial', icon: CreditCard, color: 'text-blue-600' },
      due_soon: { variant: 'outline', label: 'Due Soon', icon: Calendar, color: 'text-orange-600' },
      current: { variant: 'secondary', label: 'Current', icon: CheckCircle, color: 'text-green-600' }
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





  const filteredReceivables = mockReceivables.filter(receivable => {
    const matchesSearch = receivable.customer.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         receivable.invoiceNumber.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || receivable.status === statusFilter;
    return matchesSearch && matchesStatus;
  });



  const handleRecordPayment = () => {
    // Handle payment recording logic here
    console.log('Recording payment:', {
      receivableId: selectedReceivable?.id,
      amount: paymentAmount,
      method: paymentMethod,
      notes: paymentNotes
    });
    setIsPaymentDialogOpen(false);
    setPaymentAmount('');
    setPaymentMethod('');
    setPaymentNotes('');
  };

  const handleSendReminder = () => {
    // Handle sending reminder logic here
    console.log('Sending reminder:', {
      receivableId: selectedReceivable?.id,
      message: reminderMessage
    });
    setIsReminderDialogOpen(false);
    setReminderMessage('');
  };

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Accounts Receivable</h1>
          <p className="text-muted-foreground">
            Track customer payments and outstanding invoices
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm">
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Button size="sm">
            <Plus className="h-4 w-4 mr-2" />
            Add Receivable
          </Button>
        </div>
      </div>

      {/* Smart Insights */}
      {smartInsights.length > 0 && (
        <div className="mb-6">
          <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <Lightbulb className="h-5 w-5 text-yellow-500" />
            Smart Insights
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {smartInsights.map((insight, index) => {
              const IconComponent = insight.icon;
              const colorClasses = {
                success: 'border-green-200 bg-green-50 text-green-800',
                warning: 'border-yellow-200 bg-yellow-50 text-yellow-800',
                info: 'border-blue-200 bg-blue-50 text-blue-800',
                error: 'border-red-200 bg-red-50 text-red-800'
              };
              
              return (
                <Card key={index} className={`border-l-4 ${colorClasses[insight.type]}`}>
                  <CardContent className="p-4">
                    <div className="flex items-start gap-3">
                      <IconComponent className="h-5 w-5 mt-0.5 flex-shrink-0" />
                      <div className="flex-1">
                        <h4 className="font-medium text-sm mb-1">{insight.title}</h4>
                        <p className="text-xs mb-2 opacity-90">{insight.message}</p>
                        <Button 
                          variant="outline" 
                          size="sm" 
                          className="text-xs h-7"
                        >
                          {insight.action}
                        </Button>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        </div>
      )}

      {/* Enhanced Financial Summary */}
      <div className="mb-6">
        <h3 className="text-lg font-semibold mb-4">Enhanced Financial Summary</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Outstanding</p>
                  <p className="text-2xl font-bold">{formatCurrency(metrics.totalOutstanding)}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    Collection Rate: {metrics.collectionRate}%
                  </p>
                </div>
                <DollarSign className="h-8 w-8 text-muted-foreground" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Overdue Amount</p>
                  <p className="text-2xl font-bold text-red-600">{formatCurrency(metrics.overdueAmount)}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {metrics.overdueCount} invoices ({metrics.overduePercentage}%)
                  </p>
                </div>
                <AlertTriangle className="h-8 w-8 text-red-600" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Due Soon</p>
                  <p className="text-2xl font-bold text-yellow-600">{formatCurrency(metrics.dueSoonAmount)}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {metrics.dueSoonCount} invoices
                  </p>
                </div>
                <Clock className="h-8 w-8 text-yellow-600" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Active Customers</p>
                  <p className="text-2xl font-bold">{metrics.customerCount}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    Avg: {formatCurrency(metrics.avgInvoiceAmount)}
                  </p>
                </div>
                <Users className="h-8 w-8 text-blue-600" />
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
      
      {/* Additional Smart Metrics */}
      <div className="mb-6">
        <h3 className="text-lg font-semibold mb-4">Collection Performance</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Collection Success</p>
                  <p className="text-2xl font-bold text-green-600">{metrics.collectionRate}%</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    Avg payment term: {Math.round(metrics.avgPaymentTerm)} days
                  </p>
                </div>
                <TrendingUp className="h-8 w-8 text-green-600" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Collection Health</p>
                  <p className={`text-2xl font-bold ${
                    metrics.overduePercentage < 15 ? 'text-green-600' : 
                    metrics.overduePercentage < 25 ? 'text-yellow-600' : 'text-red-600'
                  }`}>
                    {metrics.overduePercentage < 15 ? 'Excellent' : 
                     metrics.overduePercentage < 25 ? 'Good' : 'Needs Attention'}
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {metrics.overduePercentage}% overdue rate
                  </p>
                </div>
                <BarChart3 className={`h-8 w-8 ${
                  metrics.overduePercentage < 15 ? 'text-green-600' : 
                  metrics.overduePercentage < 25 ? 'text-yellow-600' : 'text-red-600'
                }`} />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">High Risk Exposure</p>
                  <p className="text-2xl font-bold text-red-600">{formatCurrency(metrics.highRiskAmount)}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    Large overdue amounts (&gt;৳10K)
                  </p>
                </div>
                <Target className="h-8 w-8 text-red-600" />
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  placeholder="Search by customer or invoice number..."
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
                <SelectItem value="current">Current</SelectItem>
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

      {/* Receivables Table */}
      <Card>
        <CardHeader>
          <CardTitle>Outstanding Receivables</CardTitle>
          <CardDescription>
            {filteredReceivables.length} receivables found
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Customer</TableHead>
                <TableHead>Invoice #</TableHead>
                <TableHead>Due Date</TableHead>
                <TableHead>Amount</TableHead>
                <TableHead>Remaining</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredReceivables.map((receivable) => (
                <TableRow key={receivable.id}>
                  <TableCell>
                    <div>
                      <p className="font-medium">{receivable.customer}</p>
                      <p className="text-sm text-muted-foreground">{receivable.phone}</p>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="font-medium">{receivable.invoiceNumber}</p>
                      <p className="text-sm text-muted-foreground">
                        {formatDate(receivable.invoiceDate)}
                      </p>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className={cn(
                        "font-medium",
                        receivable.status === 'overdue' && "text-red-600",
                        receivable.status === 'due_soon' && "text-orange-600"
                      )}>
                        {formatDate(receivable.dueDate)}
                      </p>
                      {receivable.status === 'overdue' && (
                        <p className="text-sm text-red-600">
                          {getDaysOverdue(receivable.dueDate)} days overdue
                        </p>
                      )}
                      {receivable.status === 'due_soon' && (
                        <p className="text-sm text-orange-600">
                          Due in {getDaysToDue(receivable.dueDate)} days
                        </p>
                      )}
                    </div>
                  </TableCell>
                  <TableCell className="font-medium">
                    {formatCurrency(receivable.amount)}
                  </TableCell>
                  <TableCell className="font-medium">
                    {formatCurrency(receivable.remainingAmount)}
                  </TableCell>
                  <TableCell>
                    {getStatusBadge(receivable.status)}
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Button variant="outline" size="sm">
                        <Eye className="h-4 w-4" />
                      </Button>
                      {receivable.remainingAmount > 0 && (
                        <>
                          <Button 
                            size="sm"
                            onClick={() => {
                              setSelectedReceivable(receivable);
                              setPaymentAmount(receivable.remainingAmount.toString());
                              setIsPaymentDialogOpen(true);
                            }}
                          >
                            <Banknote className="h-4 w-4 mr-1" />
                            Record
                          </Button>
                          <Button 
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              setSelectedReceivable(receivable);
                              setReminderMessage(`Dear ${receivable.customer}, this is a friendly reminder that your invoice ${receivable.invoiceNumber} for ${formatCurrency(receivable.remainingAmount)} is ${receivable.status === 'overdue' ? 'overdue' : 'due soon'}. Please arrange payment at your earliest convenience.`);
                              setIsReminderDialogOpen(true);
                            }}
                          >
                            <Send className="h-4 w-4 mr-1" />
                            Remind
                          </Button>
                        </>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Payment Recording Dialog */}
      <Dialog open={isPaymentDialogOpen} onOpenChange={setIsPaymentDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Record Payment</DialogTitle>
            <DialogDescription>
              Record payment received from {selectedReceivable?.customer}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="space-y-2">
              <Label>Invoice Details</Label>
              <div className="p-3 bg-muted rounded-lg">
                <p className="text-sm font-medium">{selectedReceivable?.invoiceNumber}</p>
                <p className="text-sm text-muted-foreground">
                  Due: {selectedReceivable && formatDate(selectedReceivable.dueDate)}
                </p>
                <p className="text-sm text-muted-foreground">
                  Outstanding: {selectedReceivable && formatCurrency(selectedReceivable.remainingAmount)}
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
                  <SelectItem value="card">Card</SelectItem>
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
            <Button onClick={handleRecordPayment}>
              Record Payment
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Payment Reminder Dialog */}
      <Dialog open={isReminderDialogOpen} onOpenChange={setIsReminderDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Send Payment Reminder</DialogTitle>
            <DialogDescription>
              Send payment reminder to {selectedReceivable?.customer}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="space-y-2">
              <Label>Customer Contact</Label>
              <div className="p-3 bg-muted rounded-lg">
                <p className="text-sm font-medium">{selectedReceivable?.customer}</p>
                <p className="text-sm text-muted-foreground flex items-center gap-2">
                  <Phone className="h-3 w-3" />
                  {selectedReceivable?.phone}
                </p>
                <p className="text-sm text-muted-foreground">
                  {selectedReceivable?.email}
                </p>
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="reminder-message">Reminder Message</Label>
              <Textarea
                id="reminder-message"
                value={reminderMessage}
                onChange={(e) => setReminderMessage(e.target.value)}
                className="min-h-[120px]"
                placeholder="Enter reminder message"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsReminderDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleSendReminder}>
              <Send className="h-4 w-4 mr-2" />
              Send Reminder
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default AccountsReceivable;