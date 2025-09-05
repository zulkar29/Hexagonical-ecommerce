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
import {
  Search,
  Plus,
  Eye,
  Calendar,
  DollarSign,
  CheckCircle,
  Clock,
  XCircle,
  TrendingUp,
  TrendingDown
} from 'lucide-react';
import { cn } from '@/lib/utils';

// Mock data for payments
const mockPayments = [
  {
    id: 'PAY-2024-001',
    transactionId: 'TXN-001234567',
    date: '2024-02-01',
    time: '10:30 AM',
    type: 'received',
    method: 'bkash',
    amount: 16750,
    fee: 167,
    netAmount: 16583,
    status: 'completed',
    customerName: 'Fatima Khan',
    customerPhone: '01712345678',
    invoiceNumber: 'INV-2024-001',
    reference: 'bKash-TXN001234567',
    description: 'Payment for grocery order',
    reconciled: true,
    reconciledDate: '2024-02-01',
    reconciledBy: 'System Auto'
  },
  {
    id: 'PAY-2024-002',
    transactionId: 'TXN-001234568',
    date: '2024-02-01',
    time: '02:15 PM',
    type: 'received',
    method: 'cash',
    amount: 3000,
    fee: 0,
    netAmount: 3000,
    status: 'completed',
    customerName: 'Ahmed Ali',
    customerPhone: '01823456789',
    invoiceNumber: 'INV-2024-002',
    reference: 'CASH-001',
    description: 'Partial payment for household supplies',
    reconciled: true,
    reconciledDate: '2024-02-01',
    reconciledBy: 'Manager'
  },
  {
    id: 'PAY-2024-003',
    transactionId: 'TXN-001234569',
    date: '2024-02-02',
    time: '11:45 AM',
    type: 'received',
    method: 'nagad',
    amount: 13500,
    fee: 135,
    netAmount: 13365,
    status: 'completed',
    customerName: 'Rashida Begum',
    customerPhone: '01934567890',
    invoiceNumber: 'INV-2024-003',
    reference: 'Nagad-TXN987654321',
    description: 'Payment for premium grocery items',
    reconciled: true,
    reconciledDate: '2024-02-02',
    reconciledBy: 'System Auto'
  },
  {
    id: 'PAY-2024-004',
    transactionId: 'TXN-001234570',
    date: '2024-02-02',
    time: '04:20 PM',
    type: 'sent',
    method: 'bank_transfer',
    amount: 25000,
    fee: 50,
    netAmount: 25050,
    status: 'completed',
    supplierName: 'Tech Solutions BD',
    supplierPhone: '01934567890',
    invoiceNumber: 'AP-003',
    reference: 'BANK-TXN456789123',
    description: 'Payment for POS system maintenance',
    reconciled: true,
    reconciledDate: '2024-02-02',
    reconciledBy: 'Accountant'
  },
  {
    id: 'PAY-2024-005',
    transactionId: 'TXN-001234571',
    date: '2024-02-03',
    time: '09:30 AM',
    type: 'received',
    method: 'card',
    amount: 8750,
    fee: 175,
    netAmount: 8575,
    status: 'pending',
    customerName: 'Mohammad Hasan',
    customerPhone: '01645678901',
    invoiceNumber: 'INV-2024-004',
    reference: 'CARD-TXN789123456',
    description: 'Card payment for weekly groceries',
    reconciled: false,
    reconciledDate: null,
    reconciledBy: null
  },
  {
    id: 'PAY-2024-006',
    transactionId: 'TXN-001234572',
    date: '2024-02-03',
    time: '03:45 PM',
    type: 'sent',
    method: 'bkash',
    amount: 5000,
    fee: 50,
    netAmount: 5050,
    status: 'failed',
    supplierName: 'Fresh Vegetables Co.',
    supplierPhone: '01823456789',
    invoiceNumber: 'AP-002',
    reference: 'bKash-TXN123456789',
    description: 'Partial payment for vegetables supply',
    reconciled: false,
    reconciledDate: null,
    reconciledBy: null
  },
  {
    id: 'PAY-2024-007',
    transactionId: 'TXN-001234573',
    date: '2024-02-04',
    time: '01:15 PM',
    type: 'received',
    method: 'mobile_banking',
    amount: 12000,
    fee: 120,
    netAmount: 11880,
    status: 'processing',
    customerName: 'Nasir Ahmed',
    customerPhone: '01756789012',
    invoiceNumber: 'INV-2024-005',
    reference: 'MB-TXN654321987',
    description: 'Mobile banking payment for premium products',
    reconciled: false,
    reconciledDate: null,
    reconciledBy: null
  }
];

const Payments = () => {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [selectedPayment, setSelectedPayment] = useState(null);
  const [isViewDialogOpen, setIsViewDialogOpen] = useState(false);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-GB', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric'
    });
  };

  const getStatusBadge = (status) => {
    switch (status) {
      case 'completed':
        return (
          <Badge className="bg-green-100 text-green-800 hover:bg-green-100">
            <CheckCircle className="h-3 w-3 mr-1" />
            Completed
          </Badge>
        );
      case 'pending':
        return (
          <Badge className="bg-yellow-100 text-yellow-800 hover:bg-yellow-100">
            <Clock className="h-3 w-3 mr-1" />
            Pending
          </Badge>
        );
      case 'processing':
        return (
          <Badge className="bg-blue-100 text-blue-800 hover:bg-blue-100">
            <Clock className="h-3 w-3 mr-1" />
            Processing
          </Badge>
        );
      case 'failed':
        return (
          <Badge className="bg-red-100 text-red-800 hover:bg-red-100">
            <XCircle className="h-3 w-3 mr-1" />
            Failed
          </Badge>
        );
      default:
        return <Badge variant="outline">Unknown</Badge>;
    }
  };

  const getTypeBadge = (type) => {
    switch (type) {
      case 'received':
        return (
          <Badge className="bg-green-100 text-green-800 hover:bg-green-100">
            <TrendingUp className="h-3 w-3 mr-1" />
            Received
          </Badge>
        );
      case 'sent':
        return (
          <Badge className="bg-blue-100 text-blue-800 hover:bg-blue-100">
            <TrendingDown className="h-3 w-3 mr-1" />
            Sent
          </Badge>
        );
      default:
        return <Badge variant="outline">Unknown</Badge>;
    }
  };

  const filteredPayments = mockPayments.filter(payment => {
    const matchesSearch = payment.transactionId.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         payment.reference.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         (payment.customerName && payment.customerName.toLowerCase().includes(searchTerm.toLowerCase())) ||
                         (payment.supplierName && payment.supplierName.toLowerCase().includes(searchTerm.toLowerCase())) ||
                         (payment.invoiceNumber && payment.invoiceNumber.toLowerCase().includes(searchTerm.toLowerCase()));
    const matchesStatus = statusFilter === 'all' || payment.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  // Calculate summary statistics
  const totalPayments = mockPayments.length;
  const totalReceived = mockPayments.filter(p => p.type === 'received' && p.status === 'completed').reduce((sum, p) => sum + p.amount, 0);
  const totalSent = mockPayments.filter(p => p.type === 'sent' && p.status === 'completed').reduce((sum, p) => sum + p.amount, 0);
  const pendingPayments = mockPayments.filter(p => p.status === 'pending' || p.status === 'processing').length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Payment Management</h1>
          <p className="text-muted-foreground">
            Track and reconcile all payment transactions
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <Download className="h-4 w-4 mr-2" />
            Export Report
          </Button>
          <Button>
            <Plus className="h-4 w-4 mr-2" />
            Record Payment
          </Button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Payments</CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{totalPayments}</div>
            <p className="text-xs text-muted-foreground">
              All transactions
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Received</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{formatCurrency(totalReceived)}</div>
            <p className="text-xs text-muted-foreground">
              Total received
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Sent</CardTitle>
            <TrendingDown className="h-4 w-4 text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">{formatCurrency(totalSent)}</div>
            <p className="text-xs text-muted-foreground">
              Total sent
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pending</CardTitle>
            <Clock className="h-4 w-4 text-yellow-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">{pendingPayments}</div>
            <p className="text-xs text-muted-foreground">
              Awaiting completion
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Filters and Search */}
      <Card>
        <CardHeader>
          <CardTitle>Search & Filter</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col md:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search by transaction ID, reference, customer, or invoice..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[140px]">
                <Filter className="h-4 w-4 mr-2" />
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="processing">Processing</SelectItem>
                <SelectItem value="failed">Failed</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Payments Table */}
      <Card>
        <CardHeader>
          <CardTitle>Payment Transactions</CardTitle>
          <CardDescription>
            {filteredPayments.length} payments found
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Transaction</TableHead>
                <TableHead>Date & Time</TableHead>
                <TableHead>Party</TableHead>
                <TableHead>Amount</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredPayments.map((payment) => (
                <TableRow key={payment.id}>
                  <TableCell>
                    <div>
                      <p className="font-medium">{payment.transactionId}</p>
                      <p className="text-sm text-muted-foreground">{payment.reference}</p>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="font-medium">{formatDate(payment.date)}</p>
                      <p className="text-sm text-muted-foreground">{payment.time}</p>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="font-medium">
                        {payment.customerName || payment.supplierName}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {payment.customerPhone || payment.supplierPhone}
                      </p>
                      {payment.invoiceNumber && (
                        <p className="text-xs text-muted-foreground">
                          {payment.invoiceNumber}
                        </p>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <p className="font-medium">{formatCurrency(payment.amount)}</p>
                  </TableCell>
                  <TableCell>{getStatusBadge(payment.status)}</TableCell>
                  <TableCell>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => {
                        setSelectedPayment(payment);
                        setIsViewDialogOpen(true);
                      }}
                      title="View Details"
                    >
                      <Eye className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Payment Details Dialog */}
      <Dialog open={isViewDialogOpen} onOpenChange={setIsViewDialogOpen}>
        <DialogContent className="max-w-2xl">
          {selectedPayment && (
            <>
              <DialogHeader>
                <DialogTitle>Payment Details - {selectedPayment.transactionId}</DialogTitle>
                <DialogDescription>
                  Complete transaction information
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <h3 className="font-semibold mb-2">Transaction Information</h3>
                    <div className="space-y-1 text-sm">
                      <p><strong>Transaction ID:</strong> {selectedPayment.transactionId}</p>
                      <p><strong>Reference:</strong> {selectedPayment.reference}</p>
                      <p><strong>Date:</strong> {formatDate(selectedPayment.date)}</p>
                      <p><strong>Time:</strong> {selectedPayment.time}</p>
                      <p><strong>Type:</strong> {selectedPayment.type}</p>
                      <p><strong>Method:</strong> {selectedPayment.method}</p>
                      <p><strong>Status:</strong> {getStatusBadge(selectedPayment.status)}</p>
                    </div>
                  </div>
                  <div>
                    <h3 className="font-semibold mb-2">Party Information</h3>
                    <div className="space-y-1 text-sm">
                      <p><strong>Name:</strong> {selectedPayment.customerName || selectedPayment.supplierName}</p>
                      <p><strong>Phone:</strong> {selectedPayment.customerPhone || selectedPayment.supplierPhone}</p>
                      {selectedPayment.invoiceNumber && (
                        <p><strong>Invoice:</strong> {selectedPayment.invoiceNumber}</p>
                      )}
                      <p><strong>Description:</strong> {selectedPayment.description}</p>
                    </div>
                  </div>
                </div>
                
                <div>
                  <h3 className="font-semibold mb-2">Amount Details</h3>
                  <div className="text-sm">
                    <p className="text-muted-foreground">Amount</p>
                    <p className="font-medium text-lg">{formatCurrency(selectedPayment.amount)}</p>
                  </div>
                </div>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setIsViewDialogOpen(false)}>
                  Close
                </Button>
                <Button>
                  <Receipt className="h-4 w-4 mr-2" />
                  Download Receipt
                </Button>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>

      {/* Reconcile Dialog */}
      <Dialog open={isReconcileDialogOpen} onOpenChange={setIsReconcileDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Reconcile Payment</DialogTitle>
            <DialogDescription>
              Mark this payment as reconciled with your records
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="space-y-2">
              <Label>Payment Details</Label>
              <div className="p-3 bg-muted rounded-lg">
                <p className="text-sm font-medium">{selectedPayment?.transactionId}</p>
                <p className="text-sm text-muted-foreground">
                  {selectedPayment && formatCurrency(selectedPayment.netAmount)}
                </p>
                <p className="text-sm text-muted-foreground">
                  {selectedPayment?.customerName || selectedPayment?.supplierName}
                </p>
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="notes">Reconciliation Notes</Label>
              <Textarea
                id="notes"
                value={reconcileNotes}
                onChange={(e) => setReconcileNotes(e.target.value)}
                placeholder="Add notes about the reconciliation..."
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsReconcileDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleReconcile}>
              <CheckCircle className="h-4 w-4 mr-2" />
              Mark as Reconciled
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default Payments;