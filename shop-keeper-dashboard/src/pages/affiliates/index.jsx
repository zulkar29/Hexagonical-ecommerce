import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { 
  Users, 
  DollarSign, 
  TrendingUp, 
  Eye, 
  Plus,
  Search,
  Filter,
  Download,
  MoreHorizontal,
  ArrowUpDown,
  CheckCircle,
  XCircle,
  Clock,
  ExternalLink,
  Calendar
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
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
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

const Affiliates = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState("overview");
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");

  // Mock data - replace with actual API calls
  const affiliateStats = {
    totalAffiliates: 142,
    activeAffiliates: 89,
    pendingAffiliates: 12,
    totalCommissions: 45678.50,
    monthlyCommissions: 8234.25,
    conversionRate: 3.2
  };

  const affiliates = [
    {
      id: 1,
      name: "John Smith",
      email: "john@example.com",
      avatar: null,
      status: "active",
      joinDate: "2024-01-15",
      totalSales: 12450.00,
      commission: 1245.00,
      referrals: 24,
      conversionRate: 4.2
    },
    {
      id: 2,
      name: "Sarah Johnson",
      email: "sarah@example.com",
      avatar: null,
      status: "pending",
      joinDate: "2024-02-20",
      totalSales: 8750.00,
      commission: 875.00,
      referrals: 15,
      conversionRate: 3.8
    },
    {
      id: 3,
      name: "Mike Davis",
      email: "mike@example.com",
      avatar: null,
      status: "inactive",
      joinDate: "2024-01-05",
      totalSales: 3200.00,
      commission: 320.00,
      referrals: 8,
      conversionRate: 2.1
    }
  ];

  const getStatusVariant = (status) => {
    switch (status) {
      case 'active':
        return 'default';
      case 'pending':
        return 'secondary';
      case 'inactive':
        return 'destructive';
      default:
        return 'secondary';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'active':
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      case 'pending':
        return <Clock className="h-4 w-4 text-yellow-500" />;
      case 'inactive':
        return <XCircle className="h-4 w-4 text-red-500" />;
      default:
        return null;
    }
  };

  const handleStatusUpdate = (affiliateId, newStatus) => {
    // TODO: Implement status update API call
    console.log(`Updating affiliate ${affiliateId} to ${newStatus}`);
    toast.success(`Affiliate status updated to ${newStatus}`);
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
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Affiliate Program</h1>
          <p className="text-muted-foreground">
            Manage your affiliate partners and track their performance
          </p>
        </div>
        <Button onClick={() => navigate("/affiliates/register")}>
          <Plus className="mr-2 h-4 w-4" />
          Add Affiliate
        </Button>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="affiliates">Affiliates</TabsTrigger>
          <TabsTrigger value="commissions">Commissions</TabsTrigger>
          <TabsTrigger value="payouts">Payouts</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          {/* Stats Cards */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Affiliates</CardTitle>
                <Users className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{affiliateStats.totalAffiliates}</div>
                <p className="text-xs text-muted-foreground">
                  +12 from last month
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Active Affiliates</CardTitle>
                <CheckCircle className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{affiliateStats.activeAffiliates}</div>
                <p className="text-xs text-muted-foreground">
                  {((affiliateStats.activeAffiliates / affiliateStats.totalAffiliates) * 100).toFixed(1)}% of total
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Commissions</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">${affiliateStats.totalCommissions.toLocaleString()}</div>
                <p className="text-xs text-muted-foreground">
                  ${affiliateStats.monthlyCommissions.toLocaleString()} this month
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Conversion Rate</CardTitle>
                <TrendingUp className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{affiliateStats.conversionRate}%</div>
                <p className="text-xs text-muted-foreground">
                  +0.5% from last month
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Recent Activity */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Affiliate Activity</CardTitle>
              <CardDescription>
                Latest sign-ups and performance updates
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {affiliates.slice(0, 3).map((affiliate) => (
                  <div key={affiliate.id} className="flex items-center justify-between border-b pb-3 last:border-b-0 last:pb-0">
                    <div className="flex items-center space-x-3">
                      <Avatar className="h-9 w-9">
                        <AvatarImage src={affiliate.avatar} />
                        <AvatarFallback>{affiliate.name.split(' ').map(n => n[0]).join('')}</AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="text-sm font-medium">{affiliate.name}</p>
                        <p className="text-sm text-muted-foreground">{affiliate.email}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-medium">${affiliate.commission.toFixed(2)}</p>
                      <div className="flex items-center">
                        {getStatusIcon(affiliate.status)}
                        <span className="ml-1 text-xs text-muted-foreground capitalize">
                          {affiliate.status}
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Affiliates Tab */}
        <TabsContent value="affiliates" className="space-y-4">
          {/* Search and Filters */}
          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col gap-4 md:flex-row md:items-center">
                <div className="relative flex-1">
                  <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    placeholder="Search affiliates..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10"
                  />
                </div>
                <Select value={statusFilter} onValueChange={setStatusFilter}>
                  <SelectTrigger className="w-full md:w-[180px]">
                    <SelectValue placeholder="Filter by status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Status</SelectItem>
                    <SelectItem value="active">Active</SelectItem>
                    <SelectItem value="pending">Pending</SelectItem>
                    <SelectItem value="inactive">Inactive</SelectItem>
                  </SelectContent>
                </Select>
                <Button variant="outline">
                  <Download className="mr-2 h-4 w-4" />
                  Export
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Affiliates Table */}
          <Card>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Affiliate</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Join Date</TableHead>
                    <TableHead>Referrals</TableHead>
                    <TableHead>Total Sales</TableHead>
                    <TableHead>Commission</TableHead>
                    <TableHead>Conversion Rate</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {affiliates.map((affiliate) => (
                    <TableRow key={affiliate.id}>
                      <TableCell>
                        <div className="flex items-center space-x-3">
                          <Avatar className="h-8 w-8">
                            <AvatarImage src={affiliate.avatar} />
                            <AvatarFallback>{affiliate.name.split(' ').map(n => n[0]).join('')}</AvatarFallback>
                          </Avatar>
                          <div>
                            <p className="font-medium">{affiliate.name}</p>
                            <p className="text-sm text-muted-foreground">{affiliate.email}</p>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant={getStatusVariant(affiliate.status)} className="capitalize">
                          {affiliate.status}
                        </Badge>
                      </TableCell>
                      <TableCell>{formatDate(affiliate.joinDate)}</TableCell>
                      <TableCell>{affiliate.referrals}</TableCell>
                      <TableCell>${affiliate.totalSales.toFixed(2)}</TableCell>
                      <TableCell>${affiliate.commission.toFixed(2)}</TableCell>
                      <TableCell>{affiliate.conversionRate}%</TableCell>
                      <TableCell className="text-right">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" className="h-8 w-8 p-0">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuLabel>Actions</DropdownMenuLabel>
                            <DropdownMenuItem onClick={() => navigate(`/affiliates/${affiliate.id}`)}>
                              <Eye className="mr-2 h-4 w-4" />
                              View Details
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                              <ExternalLink className="mr-2 h-4 w-4" />
                              View Links
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            {affiliate.status === 'pending' && (
                              <DropdownMenuItem 
                                onClick={() => handleStatusUpdate(affiliate.id, 'active')}
                              >
                                <CheckCircle className="mr-2 h-4 w-4" />
                                Approve
                              </DropdownMenuItem>
                            )}
                            {affiliate.status === 'active' && (
                              <DropdownMenuItem 
                                onClick={() => handleStatusUpdate(affiliate.id, 'inactive')}
                              >
                                <XCircle className="mr-2 h-4 w-4" />
                                Deactivate
                              </DropdownMenuItem>
                            )}
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

        {/* Commissions Tab */}
        <TabsContent value="commissions" className="space-y-4">
          {/* Commission Stats */}
          <div className="grid gap-4 md:grid-cols-3">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Commissions</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">${affiliateStats.totalCommissions.toLocaleString()}</div>
                <p className="text-xs text-muted-foreground">
                  All time earnings
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">This Month</CardTitle>
                <TrendingUp className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">${affiliateStats.monthlyCommissions.toLocaleString()}</div>
                <p className="text-xs text-muted-foreground">
                  Current month earnings
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Pending Payouts</CardTitle>
                <Clock className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">$12,345</div>
                <p className="text-xs text-muted-foreground">
                  Ready for payout
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Commission Filters */}
          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col gap-4 md:flex-row md:items-center">
                <div className="relative flex-1">
                  <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    placeholder="Search by affiliate or order..."
                    className="pl-10"
                  />
                </div>
                <Select defaultValue="all">
                  <SelectTrigger className="w-full md:w-[180px]">
                    <SelectValue placeholder="Filter by status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Status</SelectItem>
                    <SelectItem value="pending">Pending</SelectItem>
                    <SelectItem value="paid">Paid</SelectItem>
                    <SelectItem value="cancelled">Cancelled</SelectItem>
                  </SelectContent>
                </Select>
                <Select defaultValue="thisMonth">
                  <SelectTrigger className="w-full md:w-[180px]">
                    <SelectValue placeholder="Time period" />
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

          {/* Commission Transactions Table */}
          <Card>
            <CardHeader>
              <CardTitle>Commission Transactions</CardTitle>
              <CardDescription>
                All commission transactions across your affiliate program
              </CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Affiliate</TableHead>
                    <TableHead>Order ID</TableHead>
                    <TableHead>Date</TableHead>
                    <TableHead>Sale Amount</TableHead>
                    <TableHead>Commission Rate</TableHead>
                    <TableHead>Commission Amount</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {/* Mock commission data */}
                  {[
                    {
                      id: 1,
                      affiliate: { name: "John Smith", email: "john@example.com" },
                      orderId: "ORD-001234",
                      date: "2024-03-15",
                      saleAmount: 299.99,
                      rate: 10,
                      commission: 29.99,
                      status: "paid"
                    },
                    {
                      id: 2,
                      affiliate: { name: "Sarah Johnson", email: "sarah@example.com" },
                      orderId: "ORD-001235",
                      date: "2024-03-14",
                      saleAmount: 149.99,
                      rate: 8,
                      commission: 11.99,
                      status: "pending"
                    },
                    {
                      id: 3,
                      affiliate: { name: "Mike Davis", email: "mike@example.com" },
                      orderId: "ORD-001236",
                      date: "2024-03-13",
                      saleAmount: 449.99,
                      rate: 12,
                      commission: 53.99,
                      status: "paid"
                    }
                  ].map((commission) => (
                    <TableRow key={commission.id}>
                      <TableCell>
                        <div>
                          <p className="font-medium">{commission.affiliate.name}</p>
                          <p className="text-sm text-muted-foreground">{commission.affiliate.email}</p>
                        </div>
                      </TableCell>
                      <TableCell className="font-mono text-sm">{commission.orderId}</TableCell>
                      <TableCell>{formatDate(commission.date)}</TableCell>
                      <TableCell>${commission.saleAmount.toFixed(2)}</TableCell>
                      <TableCell>{commission.rate}%</TableCell>
                      <TableCell className="font-medium">${commission.commission.toFixed(2)}</TableCell>
                      <TableCell>
                        <Badge variant={commission.status === 'paid' ? 'default' : commission.status === 'pending' ? 'secondary' : 'destructive'}>
                          {commission.status}
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
                            {commission.status === 'pending' && (
                              <>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem>
                                  <CheckCircle className="mr-2 h-4 w-4" />
                                  Mark as Paid
                                </DropdownMenuItem>
                              </>
                            )}
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

        {/* Payouts Tab */}
        <TabsContent value="payouts" className="space-y-4">
          {/* Payout Stats */}
          <div className="grid gap-4 md:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Pending Payouts</CardTitle>
                <Clock className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">$12,345</div>
                <p className="text-xs text-muted-foreground">
                  15 affiliates
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">This Month Paid</CardTitle>
                <CheckCircle className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">$8,234</div>
                <p className="text-xs text-muted-foreground">
                  12 transactions
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Paid</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">$156,789</div>
                <p className="text-xs text-muted-foreground">
                  All time payouts
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Next Payout</CardTitle>
                <Calendar className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">Mar 31</div>
                <p className="text-xs text-muted-foreground">
                  Monthly cycle
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Bulk Actions */}
          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div className="flex flex-col gap-4 md:flex-row md:items-center">
                  <div className="relative flex-1 md:w-64">
                    <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      placeholder="Search payouts..."
                      className="pl-10"
                    />
                  </div>
                  <Select defaultValue="all">
                    <SelectTrigger className="w-full md:w-[180px]">
                      <SelectValue placeholder="Filter by status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Status</SelectItem>
                      <SelectItem value="pending">Pending</SelectItem>
                      <SelectItem value="processing">Processing</SelectItem>
                      <SelectItem value="completed">Completed</SelectItem>
                      <SelectItem value="failed">Failed</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex gap-2">
                  <Button variant="outline">
                    <Download className="mr-2 h-4 w-4" />
                    Export
                  </Button>
                  <Button>
                    <Plus className="mr-2 h-4 w-4" />
                    Process Payouts
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Pending Payouts */}
          <Card>
            <CardHeader>
              <CardTitle>Pending Payouts</CardTitle>
              <CardDescription>
                Affiliates ready for payment (minimum $50 threshold met)
              </CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Affiliate</TableHead>
                    <TableHead>Earnings</TableHead>
                    <TableHead>Payment Method</TableHead>
                    <TableHead>Days Pending</TableHead>
                    <TableHead>Last Payout</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {[
                    {
                      id: 1,
                      affiliate: { name: "John Smith", email: "john@example.com" },
                      earnings: 567.89,
                      paymentMethod: "PayPal",
                      daysPending: 5,
                      lastPayout: "2024-02-15",
                      paymentEmail: "john@paypal.com"
                    },
                    {
                      id: 2,
                      affiliate: { name: "Sarah Johnson", email: "sarah@example.com" },
                      earnings: 234.56,
                      paymentMethod: "Bank Transfer",
                      daysPending: 12,
                      lastPayout: "2024-01-31",
                      paymentEmail: null
                    },
                    {
                      id: 3,
                      affiliate: { name: "Mike Davis", email: "mike@example.com" },
                      earnings: 89.45,
                      paymentMethod: "Stripe",
                      daysPending: 3,
                      lastPayout: "2024-02-28",
                      paymentEmail: "mike@stripe.com"
                    }
                  ].map((payout) => (
                    <TableRow key={payout.id}>
                      <TableCell>
                        <div className="flex items-center space-x-3">
                          <Avatar className="h-8 w-8">
                            <AvatarFallback>{payout.affiliate.name.split(' ').map(n => n[0]).join('')}</AvatarFallback>
                          </Avatar>
                          <div>
                            <p className="font-medium">{payout.affiliate.name}</p>
                            <p className="text-sm text-muted-foreground">{payout.affiliate.email}</p>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="font-medium">${payout.earnings.toFixed(2)}</TableCell>
                      <TableCell>
                        <Badge variant="outline">{payout.paymentMethod}</Badge>
                      </TableCell>
                      <TableCell>{payout.daysPending} days</TableCell>
                      <TableCell>{formatDate(payout.lastPayout)}</TableCell>
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
                              <DollarSign className="mr-2 h-4 w-4" />
                              Process Payout
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                              <Eye className="mr-2 h-4 w-4" />
                              View Details
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem>
                              <ExternalLink className="mr-2 h-4 w-4" />
                              Contact Affiliate
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

          {/* Payout History */}
          <Card>
            <CardHeader>
              <CardTitle>Payout History</CardTitle>
              <CardDescription>
                Complete history of all affiliate payouts
              </CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Payout ID</TableHead>
                    <TableHead>Affiliate</TableHead>
                    <TableHead>Amount</TableHead>
                    <TableHead>Method</TableHead>
                    <TableHead>Date</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {[
                    {
                      id: "PAY-001234",
                      affiliate: { name: "John Smith", email: "john@example.com" },
                      amount: 456.78,
                      method: "PayPal",
                      date: "2024-03-01",
                      status: "completed",
                      transactionId: "TXN-789012"
                    },
                    {
                      id: "PAY-001233",
                      affiliate: { name: "Sarah Johnson", email: "sarah@example.com" },
                      amount: 234.56,
                      method: "Bank Transfer",
                      date: "2024-03-01",
                      status: "processing",
                      transactionId: null
                    },
                    {
                      id: "PAY-001232",
                      affiliate: { name: "Mike Davis", email: "mike@example.com" },
                      amount: 123.45,
                      method: "Stripe",
                      date: "2024-02-28",
                      status: "failed",
                      transactionId: null
                    }
                  ].map((payout) => (
                    <TableRow key={payout.id}>
                      <TableCell className="font-mono text-sm">{payout.id}</TableCell>
                      <TableCell>
                        <div>
                          <p className="font-medium">{payout.affiliate.name}</p>
                          <p className="text-sm text-muted-foreground">{payout.affiliate.email}</p>
                        </div>
                      </TableCell>
                      <TableCell className="font-medium">${payout.amount.toFixed(2)}</TableCell>
                      <TableCell>
                        <Badge variant="outline">{payout.method}</Badge>
                      </TableCell>
                      <TableCell>{formatDate(payout.date)}</TableCell>
                      <TableCell>
                        <Badge variant={
                          payout.status === 'completed' ? 'default' : 
                          payout.status === 'processing' ? 'secondary' : 
                          'destructive'
                        }>
                          {payout.status}
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
                            {payout.status === 'failed' && (
                              <DropdownMenuItem>
                                <ArrowUpDown className="mr-2 h-4 w-4" />
                                Retry Payout
                              </DropdownMenuItem>
                            )}
                            <DropdownMenuItem>
                              <Download className="mr-2 h-4 w-4" />
                              Download Receipt
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
      </Tabs>
    </div>
  );
};

export default Affiliates;