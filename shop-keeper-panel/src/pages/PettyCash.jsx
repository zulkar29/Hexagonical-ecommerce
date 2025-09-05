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
  Plus,
  Download,
  Eye,
  Wallet,
  FileText,
  TrendingDown,
  TrendingUp,
  Receipt,
  PlusCircle,
  Car,
  Coffee,
  Wrench,
  Zap,
  ShoppingBag,
  AlertTriangle,
  Target,
  Lightbulb,
  BarChart3,
  Calendar,
  DollarSign,
  Filter
} from 'lucide-react';
import { cn } from '@/lib/utils';

// Mock data for petty cash transactions
const mockTransactions = [
  {
    id: 'PC-001',
    date: '2024-02-01',
    type: 'expense',
    category: 'Office Supplies',
    description: 'Printer paper and pens',
    amount: 850,
    recipient: 'Stationery Shop',
    paymentMethod: 'cash',
    receipt: 'RCP-001',
    approvedBy: 'Manager',
    status: 'approved'
  },
  {
    id: 'PC-002',
    date: '2024-02-01',
    type: 'income',
    category: 'Cash Deposit',
    description: 'Initial petty cash fund',
    amount: 10000,
    recipient: 'Petty Cash Fund',
    paymentMethod: 'bank_transfer',
    receipt: 'DEP-001',
    approvedBy: 'Owner',
    status: 'approved'
  },
  {
    id: 'PC-003',
    date: '2024-02-02',
    type: 'expense',
    category: 'Transportation',
    description: 'Delivery fuel cost',
    amount: 500,
    recipient: 'Fuel Station',
    paymentMethod: 'cash',
    receipt: 'RCP-002',
    approvedBy: 'Manager',
    status: 'approved'
  },
  {
    id: 'PC-004',
    date: '2024-02-03',
    type: 'expense',
    category: 'Refreshments',
    description: 'Tea and snacks for staff',
    amount: 300,
    recipient: 'Local Tea Shop',
    paymentMethod: 'cash',
    receipt: 'RCP-003',
    approvedBy: 'Manager',
    status: 'approved'
  },
  {
    id: 'PC-005',
    date: '2024-02-04',
    type: 'expense',
    category: 'Maintenance',
    description: 'AC repair service',
    amount: 1200,
    recipient: 'AC Service Center',
    paymentMethod: 'cash',
    receipt: 'RCP-004',
    approvedBy: 'Manager',
    status: 'pending'
  },
  {
    id: 'PC-006',
    date: '2024-02-05',
    type: 'expense',
    category: 'Utilities',
    description: 'Mobile bill payment',
    amount: 800,
    recipient: 'Telecom Company',
    paymentMethod: 'cash',
    receipt: 'RCP-005',
    approvedBy: 'Manager',
    status: 'approved'
  }
];

const expenseCategories = [
  { value: 'office_supplies', label: 'Office Supplies', icon: FileText },
  { value: 'transportation', label: 'Transportation', icon: Car },
  { value: 'refreshments', label: 'Refreshments', icon: Coffee },
  { value: 'maintenance', label: 'Maintenance', icon: Wrench },
  { value: 'utilities', label: 'Utilities', icon: Zap },
  { value: 'miscellaneous', label: 'Miscellaneous', icon: ShoppingBag }
];

// Enhanced petty cash calculations with smart insights
const calculatePettyCashMetrics = (transactions) => {
  const approvedTransactions = transactions.filter(t => t.status === 'approved');
  const totalIncome = approvedTransactions.filter(t => t.type === 'income').reduce((sum, t) => sum + t.amount, 0);
  const totalExpenses = approvedTransactions.filter(t => t.type === 'expense').reduce((sum, t) => sum + t.amount, 0);
  const currentBalance = totalIncome - totalExpenses;
  
  // Pending transactions analysis
  const pendingExpenses = transactions.filter(t => t.type === 'expense' && t.status === 'pending').reduce((sum, t) => sum + t.amount, 0);
  const pendingCount = transactions.filter(t => t.status === 'pending').length;
  
  // Category analysis
  const expensesByCategory = {};
  approvedTransactions.filter(t => t.type === 'expense').forEach(t => {
    expensesByCategory[t.category] = (expensesByCategory[t.category] || 0) + t.amount;
  });
  const topCategory = Object.entries(expensesByCategory).sort((a, b) => b[1] - a[1])[0];
  
  // Time-based analysis
  const thisMonth = new Date().getMonth();
  const thisYear = new Date().getFullYear();
  const monthlyExpenses = approvedTransactions.filter(t => {
    const transactionDate = new Date(t.date);
    return t.type === 'expense' && 
           transactionDate.getMonth() === thisMonth && 
           transactionDate.getFullYear() === thisYear;
  }).reduce((sum, t) => sum + t.amount, 0);
  
  // Risk assessment
  const balanceRatio = totalIncome > 0 ? (currentBalance / totalIncome) * 100 : 0;
  const avgTransactionAmount = approvedTransactions.length > 0 ? 
    totalExpenses / approvedTransactions.filter(t => t.type === 'expense').length : 0;
  
  // Cash flow health
  const largeExpenses = approvedTransactions.filter(t => t.type === 'expense' && t.amount > 1000).length;
  const smallExpenses = approvedTransactions.filter(t => t.type === 'expense' && t.amount <= 500).length;
  
  return {
    totalIncome,
    totalExpenses,
    currentBalance,
    pendingExpenses,
    pendingCount,
    expensesByCategory,
    topCategory,
    monthlyExpenses,
    balanceRatio,
    avgTransactionAmount,
    largeExpenses,
    smallExpenses
  };
};

// Smart recommendations for petty cash management
const generatePettyCashInsights = (metrics) => {
  const insights = [];
  
  // Balance management insights
  if (metrics.balanceRatio < 20) {
    insights.push({
      type: 'warning',
      icon: AlertTriangle,
      title: 'Low Cash Balance',
      message: `Only ${metrics.balanceRatio.toFixed(1)}% of funds remaining. Consider replenishing soon.`,
      action: 'Replenish petty cash'
    });
  }
  
  if (metrics.pendingExpenses > metrics.currentBalance) {
    insights.push({
      type: 'error',
      icon: Target,
      title: 'Insufficient Funds',
      message: `Pending expenses (৳${metrics.pendingExpenses.toLocaleString()}) exceed current balance.`,
      action: 'Review pending transactions'
    });
  }
  
  // Spending pattern insights
  if (metrics.topCategory && metrics.topCategory[1] > metrics.totalExpenses * 0.4) {
    insights.push({
      type: 'info',
      icon: BarChart3,
      title: 'High Category Spending',
      message: `${metrics.topCategory[0]} accounts for ${((metrics.topCategory[1]/metrics.totalExpenses)*100).toFixed(1)}% of expenses.`,
      action: 'Review spending patterns'
    });
  }
  
  // Transaction management insights
  if (metrics.pendingCount > 3) {
    insights.push({
      type: 'warning',
      icon: Calendar,
      title: 'Pending Approvals',
      message: `${metrics.pendingCount} transactions awaiting approval. Process them promptly.`,
      action: 'Review pending items'
    });
  }
  
  // Efficiency insights
  if (metrics.avgTransactionAmount < 300) {
    insights.push({
      type: 'success',
      icon: DollarSign,
      title: 'Efficient Spending',
      message: `Low average transaction amount (৳${metrics.avgTransactionAmount.toFixed(0)}) shows good cost control.`,
      action: 'Maintain current practices'
    });
  }
  
  return insights.slice(0, 3); // Show top 3 insights
};

const PettyCash = () => {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [typeFilter, setTypeFilter] = useState('all');
  const [categoryFilter, setCategoryFilter] = useState('all');

  const [isTransactionDialogOpen, setIsTransactionDialogOpen] = useState(false);
  const [isReplenishDialogOpen, setIsReplenishDialogOpen] = useState(false);
  const [transactionType, setTransactionType] = useState('expense');
  const [transactionData, setTransactionData] = useState({
    date: new Date().toISOString().split('T')[0],
    category: '',
    description: '',
    amount: '',
    recipient: '',
    paymentMethod: 'cash',
    receipt: ''
  });
  const [replenishAmount, setReplenishAmount] = useState('');
  const [replenishMethod, setReplenishMethod] = useState('');
  const [replenishNotes, setReplenishNotes] = useState('');

  // Calculate metrics and insights
  const metrics = calculatePettyCashMetrics(mockTransactions);
  const smartInsights = generatePettyCashInsights(metrics);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-GB');
  };

  const getTypeIcon = (type) => {
    return type === 'income' ? TrendingUp : TrendingDown;
  };

  const getTypeBadge = (type) => {
    const Icon = getTypeIcon(type);
    return (
      <Badge variant={type === 'income' ? 'default' : 'secondary'} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {type === 'income' ? 'Income' : 'Expense'}
      </Badge>
    );
  };

  const getStatusBadge = (status) => {
    const statusConfig = {
      approved: { variant: 'default', label: 'Approved' },
      pending: { variant: 'secondary', label: 'Pending' },
      rejected: { variant: 'destructive', label: 'Rejected' }
    };
    
    const config = statusConfig[status] || statusConfig.pending;
    
    return (
      <Badge variant={config.variant}>
        {config.label}
      </Badge>
    );
  };

  const getCategoryIcon = (category) => {
    const categoryIconMap = {
      'Office Supplies': FileText,
      'Transportation': Car,
      'Refreshments': Coffee,
      'Maintenance': Wrench,
      'Utilities': Zap,
      'Cash Deposit': Wallet,
      'Miscellaneous': ShoppingBag
    };
    
    return categoryIconMap[category] || ShoppingBag;
  };



  const filteredTransactions = mockTransactions.filter(transaction => {
    const matchesSearch = transaction.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         transaction.recipient.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesType = typeFilter === 'all' || transaction.type === typeFilter;
    const matchesCategory = categoryFilter === 'all' || 
      (categoryFilter === 'office' && transaction.category === 'Office Supplies') ||
      (categoryFilter === 'transport' && transaction.category === 'Transportation') ||
      (categoryFilter === 'refresh' && transaction.category === 'Refreshments') ||
      (categoryFilter === 'maintenance' && transaction.category === 'Maintenance') ||
      (categoryFilter === 'utilities' && transaction.category === 'Utilities');
    return matchesSearch && matchesType && matchesCategory;
  });



  const handleAddTransaction = () => {
    // Handle adding transaction logic here
    console.log('Adding transaction:', {
      type: transactionType,
      ...transactionData
    });
    setIsTransactionDialogOpen(false);
    setTransactionData({
      date: new Date().toISOString().split('T')[0],
      category: '',
      description: '',
      amount: '',
      recipient: '',
      paymentMethod: 'cash',
      receipt: ''
    });
  };

  const handleReplenish = () => {
    // Handle replenishing petty cash logic here
    console.log('Replenishing petty cash:', {
      amount: replenishAmount,
      method: replenishMethod,
      notes: replenishNotes
    });
    setIsReplenishDialogOpen(false);
    setReplenishAmount('');
    setReplenishMethod('');
    setReplenishNotes('');
  };

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Petty Cash Management</h1>
          <p className="text-muted-foreground">
            Track small expenses and manage petty cash fund
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm">
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Button 
            variant="outline" 
            size="sm"
            onClick={() => setIsReplenishDialogOpen(true)}
          >
            <PlusCircle className="h-4 w-4 mr-2" />
            Replenish
          </Button>
          <Button 
            size="sm"
            onClick={() => setIsTransactionDialogOpen(true)}
          >
            <Plus className="h-4 w-4 mr-2" />
            Add Transaction
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

      {/* Enhanced Summary Cards */}
      <div className="mb-6">
        <h3 className="text-lg font-semibold mb-4">Financial Overview</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Current Balance</p>
                  <p className={cn(
                    "text-2xl font-bold",
                    metrics.currentBalance < 1000 ? "text-red-600" : "text-green-600"
                  )}>
                    {formatCurrency(metrics.currentBalance)}
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {metrics.balanceRatio.toFixed(1)}% of initial fund
                  </p>
                </div>
                <Wallet className="h-8 w-8 text-blue-600" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Expenses</p>
                  <p className="text-2xl font-bold text-red-600">{formatCurrency(metrics.totalExpenses)}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    Avg: ৳{metrics.avgTransactionAmount.toFixed(0)} per transaction
                  </p>
                </div>
                <TrendingDown className="h-8 w-8 text-red-600" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Pending Expenses</p>
                  <p className="text-2xl font-bold text-orange-600">{formatCurrency(metrics.pendingExpenses)}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {metrics.pendingCount} transactions
                  </p>
                </div>
                <Calendar className="h-8 w-8 text-orange-600" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Monthly Expenses</p>
                  <p className="text-2xl font-bold">{formatCurrency(metrics.monthlyExpenses)}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    This month's spending
                  </p>
                </div>
                <BarChart3 className="h-8 w-8 text-purple-600" />
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
      
      {/* Category Analysis */}
      {metrics.topCategory && (
        <div className="mb-6">
          <h3 className="text-lg font-semibold mb-4">Spending Analysis</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Top Category</p>
                    <p className="text-2xl font-bold">{metrics.topCategory[0]}</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {formatCurrency(metrics.topCategory[1])} ({((metrics.topCategory[1]/metrics.totalExpenses)*100).toFixed(1)}%)
                    </p>
                  </div>
                  <Target className="h-8 w-8 text-blue-600" />
                </div>
              </CardContent>
            </Card>
            
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Large Expenses</p>
                    <p className="text-2xl font-bold text-red-600">{metrics.largeExpenses}</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      Transactions &gt; ৳1,000
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
                    <p className="text-sm font-medium text-muted-foreground">Small Expenses</p>
                    <p className="text-2xl font-bold text-green-600">{metrics.smallExpenses}</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      Transactions ≤ ৳500
                    </p>
                  </div>
                  <DollarSign className="h-8 w-8 text-green-600" />
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      )}

      {/* Low Balance Alert */}
      {metrics.currentBalance < 1000 && (
        <Card className="border-red-200 bg-red-50">
          <CardContent className="p-4 flex items-center gap-4">
            <AlertTriangle className="h-6 w-6 text-red-600" />
            <div className="flex-1">
              <p className="font-medium text-red-800">Low Petty Cash Balance</p>
              <p className="text-sm text-red-600">
                Current balance is {formatCurrency(metrics.currentBalance)}. Consider replenishing the fund.
              </p>
            </div>
            <Button 
              variant="outline" 
              size="sm" 
              className="border-red-300 text-red-700 hover:bg-red-100"
              onClick={() => setIsReplenishDialogOpen(true)}
            >
              Replenish Now
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  placeholder="Search transactions..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <Select value={typeFilter} onValueChange={setTypeFilter}>
              <SelectTrigger className="w-[140px]">
                <SelectValue placeholder="Type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Types</SelectItem>
                <SelectItem value="income">Income</SelectItem>
                <SelectItem value="expense">Expense</SelectItem>
              </SelectContent>
            </Select>
            <Select value={categoryFilter} onValueChange={setCategoryFilter}>
              <SelectTrigger className="w-[180px]">
                <Filter className="h-4 w-4 mr-2" />
                <SelectValue placeholder="Category" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Categories</SelectItem>
                <SelectItem value="office">Office Supplies</SelectItem>
                <SelectItem value="transport">Transportation</SelectItem>
                <SelectItem value="refresh">Refreshments</SelectItem>
                <SelectItem value="maintenance">Maintenance</SelectItem>
                <SelectItem value="utilities">Utilities</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Transactions Table */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Transactions</CardTitle>
          <CardDescription>
            {filteredTransactions.length} transactions found
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Date</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Category</TableHead>
                <TableHead>Description</TableHead>
                <TableHead>Amount</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredTransactions.map((transaction) => {
                const CategoryIcon = getCategoryIcon(transaction.category);
                return (
                  <TableRow key={transaction.id}>
                    <TableCell>
                      <p className="font-medium">{formatDate(transaction.date)}</p>
                      <p className="text-sm text-muted-foreground">{transaction.id}</p>
                    </TableCell>
                    <TableCell>
                      {getTypeBadge(transaction.type)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <CategoryIcon className="h-4 w-4 text-muted-foreground" />
                        <span className="font-medium">{transaction.category}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div>
                        <p className="font-medium">{transaction.description}</p>
                        <p className="text-sm text-muted-foreground">{transaction.recipient}</p>
                      </div>
                    </TableCell>
                    <TableCell>
                      <p className={cn(
                        "font-bold",
                        transaction.type === 'income' ? "text-green-600" : "text-red-600"
                      )}>
                        {transaction.type === 'income' ? '+' : '-'}{formatCurrency(transaction.amount)}
                      </p>
                    </TableCell>
                    <TableCell>
                      {getStatusBadge(transaction.status)}
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-2">
                        <Button variant="outline" size="sm">
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button variant="outline" size="sm">
                          <Receipt className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Add Transaction Dialog */}
      <Dialog open={isTransactionDialogOpen} onOpenChange={setIsTransactionDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Add Transaction</DialogTitle>
            <DialogDescription>
              Record a new petty cash transaction
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="type" className="text-right">
                Type
              </Label>
              <Select value={transactionType} onValueChange={setTransactionType}>
                <SelectTrigger className="col-span-3">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="expense">Expense</SelectItem>
                  <SelectItem value="income">Income</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="date" className="text-right">
                Date
              </Label>
              <Input
                id="date"
                type="date"
                value={transactionData.date}
                onChange={(e) => setTransactionData({...transactionData, date: e.target.value})}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="category" className="text-right">
                Category
              </Label>
              <Select 
                value={transactionData.category} 
                onValueChange={(value) => setTransactionData({...transactionData, category: value})}
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select category" />
                </SelectTrigger>
                <SelectContent>
                  {expenseCategories.map((cat) => (
                    <SelectItem key={cat.value} value={cat.label}>
                      {cat.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="description" className="text-right">
                Description
              </Label>
              <Input
                id="description"
                value={transactionData.description}
                onChange={(e) => setTransactionData({...transactionData, description: e.target.value})}
                className="col-span-3"
                placeholder="Enter description"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="amount" className="text-right">
                Amount
              </Label>
              <Input
                id="amount"
                type="number"
                value={transactionData.amount}
                onChange={(e) => setTransactionData({...transactionData, amount: e.target.value})}
                className="col-span-3"
                placeholder="Enter amount"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="recipient" className="text-right">
                Recipient
              </Label>
              <Input
                id="recipient"
                value={transactionData.recipient}
                onChange={(e) => setTransactionData({...transactionData, recipient: e.target.value})}
                className="col-span-3"
                placeholder="Enter recipient/vendor"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="receipt" className="text-right">
                Receipt #
              </Label>
              <Input
                id="receipt"
                value={transactionData.receipt}
                onChange={(e) => setTransactionData({...transactionData, receipt: e.target.value})}
                className="col-span-3"
                placeholder="Enter receipt number"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsTransactionDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleAddTransaction}>
              Add Transaction
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Replenish Dialog */}
      <Dialog open={isReplenishDialogOpen} onOpenChange={setIsReplenishDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Replenish Petty Cash</DialogTitle>
            <DialogDescription>
              Add funds to the petty cash account
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="space-y-2">
              <Label>Current Balance</Label>
              <div className="p-3 bg-muted rounded-lg">
                <p className="text-lg font-bold">{formatCurrency(metrics.currentBalance)}</p>
              </div>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="replenish-amount" className="text-right">
                Amount
              </Label>
              <Input
                id="replenish-amount"
                type="number"
                value={replenishAmount}
                onChange={(e) => setReplenishAmount(e.target.value)}
                className="col-span-3"
                placeholder="Enter replenish amount"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="replenish-method" className="text-right">
                Method
              </Label>
              <Select value={replenishMethod} onValueChange={setReplenishMethod}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select method" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bank_transfer">Bank Transfer</SelectItem>
                  <SelectItem value="cash_deposit">Cash Deposit</SelectItem>
                  <SelectItem value="check">Check</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="replenish-notes" className="text-right">
                Notes
              </Label>
              <Textarea
                id="replenish-notes"
                value={replenishNotes}
                onChange={(e) => setReplenishNotes(e.target.value)}
                className="col-span-3"
                placeholder="Replenishment notes (optional)"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsReplenishDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleReplenish}>
              Replenish Fund
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default PettyCash;