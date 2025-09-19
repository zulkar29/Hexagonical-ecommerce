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
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
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
  DollarSign,
  Clock,
  CheckCircle,
  XCircle,
  MoreHorizontal,
  Search,
  Filter,
  Download,
  RefreshCw,
  Calendar,
  User,
  CreditCard,
  AlertCircle,
  TrendingUp
} from 'lucide-react';
import { toast } from 'sonner';
import referralService from '@/services/referralService';

const CommissionTracking = () => {
  const [commissions, setCommissions] = useState([]);
  const [filteredCommissions, setFilteredCommissions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [dateRange, setDateRange] = useState('30d');
  const [selectedCommission, setSelectedCommission] = useState(null);
  const [showPaymentDialog, setShowPaymentDialog] = useState(false);
  const [processingPayment, setProcessingPayment] = useState(false);
  const [bulkSelection, setBulkSelection] = useState([]);
  const [showBulkActions, setShowBulkActions] = useState(false);

  const statusOptions = [
    { value: 'all', label: 'All Statuses' },
    { value: 'pending', label: 'Pending' },
    { value: 'approved', label: 'Approved' },
    { value: 'paid', label: 'Paid' },
    { value: 'rejected', label: 'Rejected' }
  ];

  const dateRangeOptions = [
    { value: '7d', label: 'Last 7 days' },
    { value: '30d', label: 'Last 30 days' },
    { value: '90d', label: 'Last 3 months' },
    { value: '1y', label: 'Last year' }
  ];

  useEffect(() => {
    fetchCommissions();
  }, [statusFilter, dateRange]);

  useEffect(() => {
    filterCommissions();
  }, [commissions, searchTerm]);

  const fetchCommissions = async () => {
    try {
      setLoading(true);
      const response = await referralService.getCommissions({
        status: statusFilter === 'all' ? undefined : statusFilter,
        dateRange,
        limit: 100
      });
      setCommissions(response.data || []);
    } catch (error) {
      toast.error('Failed to fetch commissions');
      console.error('Commission fetch error:', error);
    } finally {
      setLoading(false);
    }
  };

  const filterCommissions = () => {
    let filtered = commissions;
    
    if (searchTerm) {
      filtered = filtered.filter(commission => 
        commission.referralCode?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        commission.referrerEmail?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        commission.refereeEmail?.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }
    
    setFilteredCommissions(filtered);
  };

  const getStatusBadge = (status) => {
    const statusConfig = {
      pending: { variant: 'secondary', icon: Clock, color: 'text-yellow-600' },
      approved: { variant: 'default', icon: CheckCircle, color: 'text-blue-600' },
      paid: { variant: 'default', icon: CheckCircle, color: 'text-green-600' },
      rejected: { variant: 'destructive', icon: XCircle, color: 'text-red-600' }
    };
    
    const config = statusConfig[status] || statusConfig.pending;
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

  const handleStatusUpdate = async (commissionId, newStatus) => {
    try {
      await referralService.updateCommissionStatus(commissionId, newStatus);
      toast.success(`Commission ${newStatus} successfully`);
      fetchCommissions();
    } catch (error) {
      toast.error(`Failed to ${newStatus} commission`);
    }
  };

  const handlePayment = async (commission) => {
    setSelectedCommission(commission);
    setShowPaymentDialog(true);
  };

  const processPayment = async () => {
    if (!selectedCommission) return;
    
    try {
      setProcessingPayment(true);
      await referralService.processCommissionPayment(selectedCommission.id, {
        method: 'bank_transfer',
        notes: 'Processed via admin panel'
      });
      toast.success('Payment processed successfully');
      setShowPaymentDialog(false);
      setSelectedCommission(null);
      fetchCommissions();
    } catch (error) {
      toast.error('Failed to process payment');
    } finally {
      setProcessingPayment(false);
    }
  };

  const handleBulkSelection = (commissionId, checked) => {
    if (checked) {
      setBulkSelection(prev => [...prev, commissionId]);
    } else {
      setBulkSelection(prev => prev.filter(id => id !== commissionId));
    }
  };

  const handleBulkAction = async (action) => {
    if (bulkSelection.length === 0) {
      toast.error('Please select commissions first');
      return;
    }

    try {
      await referralService.bulkUpdateCommissions(bulkSelection, { status: action });
      toast.success(`${bulkSelection.length} commissions ${action} successfully`);
      setBulkSelection([]);
      setShowBulkActions(false);
      fetchCommissions();
    } catch (error) {
      toast.error(`Failed to ${action} selected commissions`);
    }
  };

  const exportCommissions = async () => {
    try {
      const response = await referralService.exportCommissions({
        status: statusFilter === 'all' ? undefined : statusFilter,
        dateRange,
        format: 'csv'
      });
      
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `commissions-${dateRange}.csv`);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      toast.success('Commissions exported successfully');
    } catch (error) {
      toast.error('Failed to export commissions');
    }
  };

  const getTotalCommissions = () => {
    return filteredCommissions.reduce((total, commission) => total + commission.amount, 0);
  };

  const getCommissionsByStatus = () => {
    const statusCounts = filteredCommissions.reduce((acc, commission) => {
      acc[commission.status] = (acc[commission.status] || 0) + 1;
      return acc;
    }, {});
    return statusCounts;
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-bold">Commission Tracking</h2>
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

  const statusCounts = getCommissionsByStatus();

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl font-bold">Commission Tracking</h2>
          <p className="text-muted-foreground">
            Manage and track referral commissions
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={fetchCommissions}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
          <Button variant="outline" onClick={exportCommissions}>
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
                <p className="text-sm text-muted-foreground">Total Amount</p>
                <p className="text-2xl font-bold">{formatCurrency(getTotalCommissions())}</p>
              </div>
              <DollarSign className="h-8 w-8 text-green-500" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Pending</p>
                <p className="text-2xl font-bold">{statusCounts.pending || 0}</p>
              </div>
              <Clock className="h-8 w-8 text-yellow-500" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Approved</p>
                <p className="text-2xl font-bold">{statusCounts.approved || 0}</p>
              </div>
              <CheckCircle className="h-8 w-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Paid</p>
                <p className="text-2xl font-bold">{statusCounts.paid || 0}</p>
              </div>
              <CheckCircle className="h-8 w-8 text-green-500" />
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
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search by code, referrer, or referee email..."
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

      {/* Bulk Actions */}
      {bulkSelection.length > 0 && (
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">
                {bulkSelection.length} commission(s) selected
              </span>
              <div className="flex gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => handleBulkAction('approved')}
                >
                  Approve Selected
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => handleBulkAction('rejected')}
                >
                  Reject Selected
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setBulkSelection([])}
                >
                  Clear Selection
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Commissions Table */}
      <Card>
        <CardHeader>
          <CardTitle>Commission History</CardTitle>
          <CardDescription>
            {filteredCommissions.length} commission(s) found
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-12">
                    <input
                      type="checkbox"
                      onChange={(e) => {
                        if (e.target.checked) {
                          setBulkSelection(filteredCommissions.map(c => c.id));
                        } else {
                          setBulkSelection([]);
                        }
                      }}
                      checked={bulkSelection.length === filteredCommissions.length && filteredCommissions.length > 0}
                    />
                  </TableHead>
                  <TableHead>Referral Code</TableHead>
                  <TableHead>Referrer</TableHead>
                  <TableHead>Referee</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredCommissions.map((commission) => (
                  <TableRow key={commission.id}>
                    <TableCell>
                      <input
                        type="checkbox"
                        checked={bulkSelection.includes(commission.id)}
                        onChange={(e) => handleBulkSelection(commission.id, e.target.checked)}
                      />
                    </TableCell>
                    <TableCell className="font-mono">
                      {commission.referralCode}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <User className="h-4 w-4 text-muted-foreground" />
                        {commission.referrerEmail}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <User className="h-4 w-4 text-muted-foreground" />
                        {commission.refereeEmail}
                      </div>
                    </TableCell>
                    <TableCell className="font-medium">
                      {formatCurrency(commission.amount)}
                    </TableCell>
                    <TableCell>
                      {getStatusBadge(commission.status)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Calendar className="h-4 w-4 text-muted-foreground" />
                        {new Date(commission.createdAt).toLocaleDateString()}
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
                          {commission.status === 'pending' && (
                            <>
                              <DropdownMenuItem
                                onClick={() => handleStatusUpdate(commission.id, 'approved')}
                              >
                                <CheckCircle className="h-4 w-4 mr-2" />
                                Approve
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                onClick={() => handleStatusUpdate(commission.id, 'rejected')}
                              >
                                <XCircle className="h-4 w-4 mr-2" />
                                Reject
                              </DropdownMenuItem>
                            </>
                          )}
                          {commission.status === 'approved' && (
                            <DropdownMenuItem
                              onClick={() => handlePayment(commission)}
                            >
                              <CreditCard className="h-4 w-4 mr-2" />
                              Process Payment
                            </DropdownMenuItem>
                          )}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
          
          {filteredCommissions.length === 0 && (
            <div className="text-center py-8">
              <AlertCircle className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">No commissions found</h3>
              <p className="text-muted-foreground">
                {searchTerm || statusFilter !== 'all' 
                  ? 'Try adjusting your filters'
                  : 'No commissions have been generated yet'
                }
              </p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Payment Dialog */}
      <Dialog open={showPaymentDialog} onOpenChange={setShowPaymentDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Process Commission Payment</DialogTitle>
            <DialogDescription>
              Confirm payment for commission #{selectedCommission?.id}
            </DialogDescription>
          </DialogHeader>
          
          {selectedCommission && (
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium">Referral Code</label>
                  <p className="font-mono">{selectedCommission.referralCode}</p>
                </div>
                <div>
                  <label className="text-sm font-medium">Amount</label>
                  <p className="font-medium">{formatCurrency(selectedCommission.amount)}</p>
                </div>
                <div>
                  <label className="text-sm font-medium">Referrer</label>
                  <p>{selectedCommission.referrerEmail}</p>
                </div>
                <div>
                  <label className="text-sm font-medium">Date</label>
                  <p>{new Date(selectedCommission.createdAt).toLocaleDateString()}</p>
                </div>
              </div>
            </div>
          )}
          
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setShowPaymentDialog(false)}
            >
              Cancel
            </Button>
            <Button
              onClick={processPayment}
              disabled={processingPayment}
            >
              {processingPayment ? (
                <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
              ) : (
                <CreditCard className="h-4 w-4 mr-2" />
              )}
              Process Payment
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default CommissionTracking;