import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
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
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';
import {
  History,
  Search,
  Filter,
  Download,
  RefreshCw,
  Calendar,
  User,
  Share2,
  Eye,
  MoreHorizontal,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  TrendingUp,
  Users,
  DollarSign,
  MousePointer
} from 'lucide-react';
import { toast } from 'sonner';
import referralService from '@/services/referralService';

const ReferralHistory = () => {
  const [referrals, setReferrals] = useState([]);
  const [filteredReferrals, setFilteredReferrals] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [dateRange, setDateRange] = useState('30d');
  const [selectedReferral, setSelectedReferral] = useState(null);
  const [showDetailsDialog, setShowDetailsDialog] = useState(false);
  const [activeTab, setActiveTab] = useState('all');

  const statusOptions = [
    { value: 'all', label: 'All Statuses' },
    { value: 'active', label: 'Active' },
    { value: 'expired', label: 'Expired' },
    { value: 'disabled', label: 'Disabled' },
    { value: 'completed', label: 'Completed' }
  ];

  const dateRangeOptions = [
    { value: '7d', label: 'Last 7 days' },
    { value: '30d', label: 'Last 30 days' },
    { value: '90d', label: 'Last 3 months' },
    { value: '1y', label: 'Last year' },
    { value: 'all', label: 'All time' }
  ];

  useEffect(() => {
    fetchReferrals();
  }, [statusFilter, dateRange, activeTab]);

  useEffect(() => {
    filterReferrals();
  }, [referrals, searchTerm]);

  const fetchReferrals = async () => {
    try {
      setLoading(true);
      const response = await referralService.getReferralHistory({
        status: statusFilter === 'all' ? undefined : statusFilter,
        dateRange: dateRange === 'all' ? undefined : dateRange,
        type: activeTab === 'all' ? undefined : activeTab,
        limit: 100
      });
      setReferrals(response.data || []);
    } catch (error) {
      toast.error('Failed to fetch referral history');
      console.error('Referral history fetch error:', error);
    } finally {
      setLoading(false);
    }
  };

  const filterReferrals = () => {
    let filtered = referrals;
    
    if (searchTerm) {
      filtered = filtered.filter(referral => 
        referral.code?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        referral.referrerEmail?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        referral.refereeEmail?.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }
    
    setFilteredReferrals(filtered);
  };

  const getStatusBadge = (status) => {
    const statusConfig = {
      active: { variant: 'default', icon: CheckCircle, color: 'text-green-600' },
      expired: { variant: 'secondary', icon: Clock, color: 'text-gray-600' },
      disabled: { variant: 'destructive', icon: XCircle, color: 'text-red-600' },
      completed: { variant: 'default', icon: CheckCircle, color: 'text-blue-600' }
    };
    
    const config = statusConfig[status] || statusConfig.active;
    const Icon = config.icon;
    
    return (
      <Badge variant={config.variant} className="flex items-center gap-1">
        <Icon className="h-3 w-3" />
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </Badge>
    );
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-BD', {
      style: 'currency',
      currency: 'BDT',
      minimumFractionDigits: 0
    }).format(amount);
  };

  const handleViewDetails = (referral) => {
    setSelectedReferral(referral);
    setShowDetailsDialog(true);
  };

  const exportHistory = async () => {
    try {
      const response = await referralService.exportReferralHistory({
        status: statusFilter === 'all' ? undefined : statusFilter,
        dateRange: dateRange === 'all' ? undefined : dateRange,
        type: activeTab === 'all' ? undefined : activeTab,
        format: 'csv'
      });
      
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `referral-history-${dateRange}.csv`);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      toast.success('Referral history exported successfully');
    } catch (error) {
      toast.error('Failed to export referral history');
    }
  };

  const getHistoryStats = () => {
    const stats = {
      total: filteredReferrals.length,
      active: filteredReferrals.filter(r => r.status === 'active').length,
      completed: filteredReferrals.filter(r => r.status === 'completed').length,
      totalClicks: filteredReferrals.reduce((sum, r) => sum + (r.clicks || 0), 0),
      totalConversions: filteredReferrals.reduce((sum, r) => sum + (r.conversions || 0), 0),
      totalCommissions: filteredReferrals.reduce((sum, r) => sum + (r.totalCommissions || 0), 0)
    };
    
    stats.conversionRate = stats.totalClicks > 0 ? stats.totalConversions / stats.totalClicks : 0;
    
    return stats;
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-bold">Referral History</h2>
          <div className="flex gap-2">
            <div className="w-32 h-10 bg-muted animate-pulse rounded" />
            <div className="w-24 h-10 bg-muted animate-pulse rounded" />
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="h-24 bg-muted animate-pulse rounded-lg" />
          ))}
        </div>
        <div className="h-96 bg-muted animate-pulse rounded-lg" />
      </div>
    );
  }

  const stats = getHistoryStats();

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl font-bold">Referral History</h2>
          <p className="text-muted-foreground">
            Complete history of all referral activities
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={fetchReferrals}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
          <Button variant="outline" onClick={exportHistory}>
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Referrals</p>
                <p className="text-2xl font-bold">{stats.total}</p>
              </div>
              <Share2 className="h-8 w-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Clicks</p>
                <p className="text-2xl font-bold">{stats.totalClicks.toLocaleString()}</p>
              </div>
              <MousePointer className="h-8 w-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Conversions</p>
                <p className="text-2xl font-bold">{stats.totalConversions}</p>
              </div>
              <TrendingUp className="h-8 w-8 text-green-500" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Commissions</p>
                <p className="text-2xl font-bold">{formatCurrency(stats.totalCommissions)}</p>
              </div>
              <DollarSign className="h-8 w-8 text-orange-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="all">All Referrals</TabsTrigger>
          <TabsTrigger value="codes">Referral Codes</TabsTrigger>
          <TabsTrigger value="links">Referral Links</TabsTrigger>
          <TabsTrigger value="campaigns">Campaigns</TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="space-y-4">
          {/* Filters */}
          <Card>
            <CardContent className="p-4">
              <div className="flex flex-col sm:flex-row gap-4">
                <div className="flex-1">
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                      placeholder="Search by code, referrer, or referee..."
                      value={searchTerm}
                      onChange={(e) => setSearchTerm(e.target.value)}
                      className="pl-10"
                    />
                  </div>
                </div>
                <Select value={statusFilter} onValueChange={setStatusFilter}>
                  <SelectTrigger className="w-40">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {statusOptions.map(option => (
                      <SelectItem key={option.value} value={option.value}>
                        {option.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <Select value={dateRange} onValueChange={setDateRange}>
                  <SelectTrigger className="w-40">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {dateRangeOptions.map(option => (
                      <SelectItem key={option.value} value={option.value}>
                        {option.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </CardContent>
          </Card>

          {/* History Table */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <History className="h-5 w-5" />
                Referral History
              </CardTitle>
              <CardDescription>
                {filteredReferrals.length} referral(s) found
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Code/Link</TableHead>
                      <TableHead>Type</TableHead>
                      <TableHead>Referrer</TableHead>
                      <TableHead>Clicks</TableHead>
                      <TableHead>Conversions</TableHead>
                      <TableHead>Commission</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredReferrals.map((referral) => (
                      <TableRow key={referral.id}>
                        <TableCell>
                          <div className="font-mono text-sm">
                            {referral.code || referral.linkId}
                          </div>
                          {referral.description && (
                            <div className="text-xs text-muted-foreground mt-1">
                              {referral.description}
                            </div>
                          )}
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline">
                            {referral.type || 'Code'}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <User className="h-4 w-4 text-muted-foreground" />
                            <span className="text-sm">{referral.referrerEmail}</span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <MousePointer className="h-4 w-4 text-muted-foreground" />
                            <span>{referral.clicks || 0}</span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <TrendingUp className="h-4 w-4 text-muted-foreground" />
                            <span>{referral.conversions || 0}</span>
                            {referral.clicks > 0 && (
                              <span className="text-xs text-muted-foreground">
                                ({((referral.conversions || 0) / referral.clicks * 100).toFixed(1)}%)
                              </span>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <DollarSign className="h-4 w-4 text-muted-foreground" />
                            <span className="font-medium">
                              {formatCurrency(referral.totalCommissions || 0)}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          {getStatusBadge(referral.status)}
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <Calendar className="h-4 w-4 text-muted-foreground" />
                            <span className="text-sm">
                              {new Date(referral.createdAt).toLocaleDateString()}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" className="h-8 w-8 p-0">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuLabel>Actions</DropdownMenuLabel>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem
                                onClick={() => handleViewDetails(referral)}
                              >
                                <Eye className="h-4 w-4 mr-2" />
                                View Details
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
              
              {filteredReferrals.length === 0 && (
                <div className="text-center py-8">
                  <AlertCircle className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-medium mb-2">No referrals found</h3>
                  <p className="text-muted-foreground">
                    {searchTerm || statusFilter !== 'all' 
                      ? 'Try adjusting your filters'
                      : 'No referral activity recorded yet'
                    }
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Details Dialog */}
      <Dialog open={showDetailsDialog} onOpenChange={setShowDetailsDialog}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Referral Details</DialogTitle>
            <DialogDescription>
              Detailed information about this referral
            </DialogDescription>
          </DialogHeader>
          
          {selectedReferral && (
            <div className="space-y-6">
              {/* Basic Info */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-muted-foreground">Code/ID</label>
                  <p className="font-mono">{selectedReferral.code || selectedReferral.linkId}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">Type</label>
                  <p>{selectedReferral.type || 'Code'}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">Status</label>
                  <div className="mt-1">
                    {getStatusBadge(selectedReferral.status)}
                  </div>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">Created</label>
                  <p>{new Date(selectedReferral.createdAt).toLocaleString()}</p>
                </div>
              </div>

              {/* Performance Metrics */}
              <div>
                <h4 className="font-medium mb-3">Performance Metrics</h4>
                <div className="grid grid-cols-3 gap-4">
                  <div className="text-center p-3 bg-muted rounded-lg">
                    <p className="text-2xl font-bold">{selectedReferral.clicks || 0}</p>
                    <p className="text-sm text-muted-foreground">Clicks</p>
                  </div>
                  <div className="text-center p-3 bg-muted rounded-lg">
                    <p className="text-2xl font-bold">{selectedReferral.conversions || 0}</p>
                    <p className="text-sm text-muted-foreground">Conversions</p>
                  </div>
                  <div className="text-center p-3 bg-muted rounded-lg">
                    <p className="text-2xl font-bold">
                      {selectedReferral.clicks > 0 
                        ? `${((selectedReferral.conversions || 0) / selectedReferral.clicks * 100).toFixed(1)}%`
                        : '0%'
                      }
                    </p>
                    <p className="text-sm text-muted-foreground">Conversion Rate</p>
                  </div>
                </div>
              </div>

              {/* Commission Info */}
              <div>
                <h4 className="font-medium mb-3">Commission Information</h4>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">Commission Rate</label>
                    <p>
                      {selectedReferral.commissionRate}{
                        selectedReferral.commissionType === 'percentage' ? '%' : '৳'
                      }
                    </p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">Total Earned</label>
                    <p className="font-medium">
                      {formatCurrency(selectedReferral.totalCommissions || 0)}
                    </p>
                  </div>
                </div>
              </div>

              {/* Additional Details */}
              {selectedReferral.description && (
                <div>
                  <label className="text-sm font-medium text-muted-foreground">Description</label>
                  <p className="mt-1">{selectedReferral.description}</p>
                </div>
              )}

              {selectedReferral.expiresAt && (
                <div>
                  <label className="text-sm font-medium text-muted-foreground">Expires</label>
                  <p className="mt-1">{new Date(selectedReferral.expiresAt).toLocaleString()}</p>
                </div>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default ReferralHistory;