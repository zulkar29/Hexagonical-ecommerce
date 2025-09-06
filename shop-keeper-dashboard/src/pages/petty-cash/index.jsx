import { useState } from "react";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { 
  Wallet,
  Plus,
  Minus,
  TrendingUp,
  TrendingDown,
  Search,
  Download,
  MoreHorizontal,
  Eye,
  Edit,
  Trash2,
  Receipt,
  AlertCircle
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
// import { pettyCashApi } from "@/lib/api";

const PettyCash = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState("overview");
  const [searchQuery, setSearchQuery] = useState("");
  const [typeFilter, setTypeFilter] = useState("all");
  const [dateFilter, setDateFilter] = useState("thisMonth");

  // Mock data - replace with actual API calls
  const pettyCashStats = {
    currentBalance: 2450.75,
    totalIn: 5000.00,
    totalOut: 2549.25,
    monthlyIn: 1200.00,
    monthlyOut: 875.50,
    transactionCount: 127,
    lastReplenishment: "2024-03-15"
  };

  const transactions = [
    {
      id: 1,
      type: "out",
      category: "Office Supplies",
      description: "Printer paper and ink cartridges",
      amount: 89.50,
      date: "2024-03-18",
      receipt: "REC-001234",
      handledBy: "John Admin",
      approved: true
    },
    {
      id: 2,
      type: "in",
      category: "Replenishment",
      description: "Monthly petty cash replenishment",
      amount: 500.00,
      date: "2024-03-15",
      receipt: "REP-001235",
      handledBy: "Finance Team",
      approved: true
    },
    {
      id: 3,
      type: "out",
      category: "Travel",
      description: "Client meeting taxi fare",
      amount: 25.75,
      date: "2024-03-14",
      receipt: "REC-001233",
      handledBy: "Sales Team",
      approved: true
    },
    {
      id: 4,
      type: "out",
      category: "Refreshments",
      description: "Office coffee and snacks",
      amount: 45.20,
      date: "2024-03-13",
      receipt: "REC-001232",
      handledBy: "HR Team",
      approved: false
    },
    {
      id: 5,
      type: "out",
      category: "Maintenance",
      description: "Emergency plumbing repair",
      amount: 150.00,
      date: "2024-03-12",
      receipt: "REC-001231",
      handledBy: "Facility Manager",
      approved: true
    }
  ];


  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  };

  const getTransactionIcon = (type) => {
    return type === 'in' ? (
      <Plus className="h-4 w-4 text-green-500" />
    ) : (
      <Minus className="h-4 w-4 text-red-500" />
    );
  };

  const getTransactionVariant = (type) => {
    return type === 'in' ? 'default' : 'secondary';
  };

  const handleDeleteTransaction = (id) => {
    // TODO: Implement delete API call
    console.log(`Deleting transaction ${id}`);
    toast.success("Transaction deleted successfully");
  };

  const handleApproveTransaction = (id) => {
    // TODO: Implement approval API call
    console.log(`Approving transaction ${id}`);
    toast.success("Transaction approved successfully");
  };

  const filteredTransactions = transactions.filter(transaction => {
    const matchesSearch = transaction.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         transaction.category.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesType = typeFilter === 'all' || transaction.type === typeFilter;
    return matchesSearch && matchesType;
  });

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Petty Cash Management</h1>
          <p className="text-muted-foreground">
            Track and manage small cash transactions and expenses
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <TrendingUp className="mr-2 h-4 w-4" />
            Replenish Cash
          </Button>
          <Button onClick={() => navigate("/petty-cash/new")}>
            <Plus className="mr-2 h-4 w-4" />
            New Transaction
          </Button>
        </div>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="transactions">Transactions</TabsTrigger>
          <TabsTrigger value="reports">Reports</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          {/* Balance Summary */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Current Balance</CardTitle>
                <Wallet className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatCurrency(pettyCashStats.currentBalance)}</div>
                <p className="text-xs text-muted-foreground">
                  {pettyCashStats.currentBalance > 500 ? "Healthy balance" : "Consider replenishment"}
                </p>
              </CardContent>
            </Card>
            
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Monthly In</CardTitle>
                <TrendingUp className="h-4 w-4 text-green-500" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">{formatCurrency(pettyCashStats.monthlyIn)}</div>
                <p className="text-xs text-muted-foreground">
                  Replenishments this month
                </p>
              </CardContent>
            </Card>
            
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Monthly Out</CardTitle>
                <TrendingDown className="h-4 w-4 text-red-500" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-red-600">{formatCurrency(pettyCashStats.monthlyOut)}</div>
                <p className="text-xs text-muted-foreground">
                  Expenses this month
                </p>
              </CardContent>
            </Card>
            
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Transactions</CardTitle>
                <Receipt className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{pettyCashStats.transactionCount}</div>
                <p className="text-xs text-muted-foreground">
                  All time transactions
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Recent Transactions */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Transactions</CardTitle>
              <CardDescription>
                Latest petty cash activities
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {transactions.slice(0, 5).map((transaction) => (
                  <div key={transaction.id} className="flex items-center justify-between border-b pb-3 last:border-b-0 last:pb-0">
                    <div className="flex items-center space-x-3">
                      {getTransactionIcon(transaction.type)}
                      <div>
                        <p className="text-sm font-medium">{transaction.description}</p>
                        <p className="text-sm text-muted-foreground">{transaction.category} • {formatDate(transaction.date)}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className={`text-sm font-medium ${transaction.type === 'in' ? 'text-green-600' : 'text-red-600'}`}>
                        {transaction.type === 'in' ? '+' : '-'}{formatCurrency(transaction.amount)}
                      </p>
                      <div className="flex items-center gap-2">
                        {!transaction.approved && (
                          <AlertCircle className="h-3 w-3 text-yellow-500" />
                        )}
                        <span className="text-xs text-muted-foreground">
                          {transaction.handledBy}
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Low Balance Alert */}
          {pettyCashStats.currentBalance < 500 && (
            <Card className="border-yellow-200 bg-yellow-50">
              <CardContent className="pt-6">
                <div className="flex items-center gap-3">
                  <AlertCircle className="h-5 w-5 text-yellow-600" />
                  <div>
                    <h3 className="font-semibold text-yellow-800">Low Balance Alert</h3>
                    <p className="text-sm text-yellow-700">
                      Current balance is below recommended threshold. Consider replenishing petty cash.
                    </p>
                  </div>
                  <Button variant="outline" className="ml-auto">
                    Replenish Now
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        {/* Transactions Tab */}
        <TabsContent value="transactions" className="space-y-4">
          {/* Search and Filters */}
          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col gap-4 md:flex-row md:items-center">
                <div className="relative flex-1">
                  <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    placeholder="Search transactions..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10"
                  />
                </div>
                <Select value={typeFilter} onValueChange={setTypeFilter}>
                  <SelectTrigger className="w-full md:w-[180px]">
                    <SelectValue placeholder="Filter by type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Types</SelectItem>
                    <SelectItem value="in">Cash In</SelectItem>
                    <SelectItem value="out">Cash Out</SelectItem>
                  </SelectContent>
                </Select>
                <Select value={dateFilter} onValueChange={setDateFilter}>
                  <SelectTrigger className="w-full md:w-[180px]">
                    <SelectValue placeholder="Date range" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="thisMonth">This Month</SelectItem>
                    <SelectItem value="lastMonth">Last Month</SelectItem>
                    <SelectItem value="last3Months">Last 3 Months</SelectItem>
                    <SelectItem value="thisYear">This Year</SelectItem>
                  </SelectContent>
                </Select>
                <Button variant="outline">
                  <Download className="mr-2 h-4 w-4" />
                  Export
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Transactions Table */}
          <Card>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Type</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead>Category</TableHead>
                    <TableHead>Amount</TableHead>
                    <TableHead>Date</TableHead>
                    <TableHead>Handled By</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredTransactions.map((transaction) => (
                    <TableRow key={transaction.id}>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          {getTransactionIcon(transaction.type)}
                          <Badge variant={getTransactionVariant(transaction.type)}>
                            {transaction.type === 'in' ? 'Cash In' : 'Cash Out'}
                          </Badge>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div>
                          <p className="font-medium">{transaction.description}</p>
                          <p className="text-sm text-muted-foreground">{transaction.receipt}</p>
                        </div>
                      </TableCell>
                      <TableCell>{transaction.category}</TableCell>
                      <TableCell>
                        <span className={transaction.type === 'in' ? 'text-green-600' : 'text-red-600'}>
                          {transaction.type === 'in' ? '+' : '-'}{formatCurrency(transaction.amount)}
                        </span>
                      </TableCell>
                      <TableCell>{formatDate(transaction.date)}</TableCell>
                      <TableCell>{transaction.handledBy}</TableCell>
                      <TableCell>
                        <Badge variant={transaction.approved ? 'default' : 'secondary'}>
                          {transaction.approved ? 'Approved' : 'Pending'}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" className="h-8 w-8 p-0">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuLabel>Actions</DropdownMenuLabel>
                            <DropdownMenuItem>
                              <Eye className="mr-2 h-4 w-4" />
                              View Details
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => navigate(`/petty-cash/${transaction.id}/edit`)}>
                              <Edit className="mr-2 h-4 w-4" />
                              Edit
                            </DropdownMenuItem>
                            {!transaction.approved && (
                              <DropdownMenuItem onClick={() => handleApproveTransaction(transaction.id)}>
                                <Receipt className="mr-2 h-4 w-4" />
                                Approve
                              </DropdownMenuItem>
                            )}
                            <DropdownMenuSeparator />
                            <DropdownMenuItem 
                              onClick={() => handleDeleteTransaction(transaction.id)}
                              className="text-red-600"
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Delete
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

        {/* Reports Tab */}
        <TabsContent value="reports">
          <Card>
            <CardHeader>
              <CardTitle>Petty Cash Reports</CardTitle>
              <CardDescription>
                Generate detailed reports on petty cash usage and trends
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-center py-8 text-muted-foreground">
                Reports and analytics coming soon...
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default PettyCash;