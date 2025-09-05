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
  Eye,
  FileText,
  Calendar,
  DollarSign,
  CheckCircle,
  Clock,
  AlertTriangle,
  XCircle,
  Edit
} from 'lucide-react';
import { cn } from '@/lib/utils';

// Mock data for invoices
const mockInvoices = [
  {
    id: 'INV-2024-001',
    customerName: 'Fatima Khan',
    customerEmail: 'fatima.khan@email.com',
    customerPhone: '01712345678',
    customerAddress: 'House 23, Road 15, Dhanmondi, Dhaka',
    issueDate: '2024-01-10',
    dueDate: '2024-02-10',
    subtotal: 15000,
    tax: 2250,
    discount: 500,
    totalAmount: 16750,
    paidAmount: 16750,
    remainingAmount: 0,
    status: 'paid',
    paymentStatus: 'paid',
    paymentMethod: 'Bank Transfer',
    orderNumber: 'ORD-1234',
    notes: 'Monthly grocery order - premium items',
    items: [
      { name: 'Basmati Rice 5kg', quantity: 2, price: 3500, total: 7000 },
      { name: 'Premium Cooking Oil 2L', quantity: 3, price: 1200, total: 3600 },
      { name: 'Fresh Vegetables Mix', quantity: 1, price: 2500, total: 2500 },
      { name: 'Dairy Products Bundle', quantity: 1, price: 1900, total: 1900 }
    ]
  },
  {
    id: 'INV-2024-002',
    customerName: 'Ahmed Ali',
    customerEmail: 'ahmed.ali@email.com',
    customerPhone: '01823456789',
    customerAddress: 'Flat 8B, Gulshan Avenue, Gulshan-1, Dhaka',
    issueDate: '2024-01-15',
    dueDate: '2024-02-15',
    subtotal: 8500,
    tax: 1275,
    discount: 200,
    totalAmount: 9575,
    paidAmount: 3000,
    remainingAmount: 6575,
    status: 'partial',
    paymentStatus: 'partial',
    paymentMethod: 'Cash',
    orderNumber: 'ORD-1235',
    notes: 'Regular household supplies',
    items: [
      { name: 'Rice 2kg', quantity: 2, price: 1500, total: 3000 },
      { name: 'Lentils Mix 1kg', quantity: 3, price: 800, total: 2400 },
      { name: 'Spices Bundle', quantity: 1, price: 1200, total: 1200 },
      { name: 'Cleaning Supplies', quantity: 1, price: 1900, total: 1900 }
    ]
  },
  {
    id: 'INV-2024-003',
    customerName: 'Rashida Begum',
    customerEmail: 'rashida.begum@email.com',
    customerPhone: '01934567890',
    customerAddress: 'Block C, House 23, Uttara Sector 7, Dhaka',
    issueDate: '2024-01-20',
    dueDate: '2024-02-20',
    subtotal: 12000,
    tax: 1800,
    discount: 300,
    totalAmount: 13500,
    paidAmount: 13500,
    remainingAmount: 0,
    status: 'paid',
    paymentStatus: 'paid',
    paymentMethod: 'Mobile Banking',
    orderNumber: 'ORD-1236',
    notes: 'Premium grocery items for special occasion',
    items: [
      { name: 'Premium Basmati Rice 5kg', quantity: 1, price: 4000, total: 4000 },
      { name: 'Organic Vegetables', quantity: 2, price: 1500, total: 3000 },
      { name: 'Premium Meat 2kg', quantity: 1, price: 3500, total: 3500 },
      { name: 'Imported Fruits', quantity: 1, price: 1500, total: 1500 }
    ]
  },
  {
    id: 'INV-2024-004',
    customerName: 'Mohammad Hasan',
    customerEmail: 'hasan.mohammad@email.com',
    customerPhone: '01645678901',
    customerAddress: 'Road 5, Block B, Mirpur-10, Dhaka',
    issueDate: '2024-01-25',
    dueDate: '2024-02-25',
    subtotal: 6500,
    tax: 975,
    discount: 150,
    totalAmount: 7325,
    paidAmount: 0,
    remainingAmount: 7325,
    status: 'pending',
    paymentStatus: 'unpaid',
    paymentMethod: '',
    orderNumber: 'ORD-1237',
    notes: 'Weekly family groceries',
    items: [
      { name: 'Rice 3kg', quantity: 1, price: 2100, total: 2100 },
      { name: 'Vegetables 2kg', quantity: 1, price: 1200, total: 1200 },
      { name: 'Fish 1kg', quantity: 1, price: 1800, total: 1800 },
      { name: 'Household Items', quantity: 1, price: 1400, total: 1400 }
    ]
  },
  {
    id: 'INV-2024-005',
    customerName: 'Nasir Ahmed',
    customerEmail: 'nasir.ahmed@email.com',
    customerPhone: '01756789012',
    customerAddress: 'House 78, New Eskaton Road, Ramna, Dhaka',
    issueDate: '2024-01-28',
    dueDate: '2024-02-12',
    subtotal: 9500,
    tax: 1425,
    discount: 250,
    totalAmount: 10675,
    paidAmount: 0,
    remainingAmount: 10675,
    status: 'overdue',
    paymentStatus: 'unpaid',
    paymentMethod: '',
    orderNumber: 'ORD-1238',
    notes: 'Premium brand products',
    items: [
      { name: 'Premium Rice 5kg', quantity: 1, price: 3800, total: 3800 },
      { name: 'Imported Oil 2L', quantity: 2, price: 1400, total: 2800 },
      { name: 'Organic Vegetables', quantity: 1, price: 1600, total: 1600 },
      { name: 'Premium Dairy', quantity: 1, price: 1300, total: 1300 }
    ]
  }
];

const Invoices = () => {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');

  const [selectedInvoice, setSelectedInvoice] = useState(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
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
      case 'paid':
        return (
          <Badge className="bg-green-100 text-green-800 hover:bg-green-100">
            <CheckCircle className="h-3 w-3 mr-1" />
            Paid
          </Badge>
        );
      case 'partial':
        return (
          <Badge className="bg-yellow-100 text-yellow-800 hover:bg-yellow-100">
            <Clock className="h-3 w-3 mr-1" />
            Partial
          </Badge>
        );
      case 'pending':
        return (
          <Badge className="bg-blue-100 text-blue-800 hover:bg-blue-100">
            <Clock className="h-3 w-3 mr-1" />
            Pending
          </Badge>
        );
      case 'overdue':
        return (
          <Badge className="bg-red-100 text-red-800 hover:bg-red-100">
            <AlertTriangle className="h-3 w-3 mr-1" />
            Overdue
          </Badge>
        );
      case 'cancelled':
        return (
          <Badge className="bg-gray-100 text-gray-800 hover:bg-gray-100">
            <XCircle className="h-3 w-3 mr-1" />
            Cancelled
          </Badge>
        );
      default:
        return <Badge variant="outline">Unknown</Badge>;
    }
  };



  const getDaysOverdue = (dueDate) => {
    const today = new Date();
    const due = new Date(dueDate);
    const diffTime = today - due;
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays > 0 ? diffDays : 0;
  };

  const filteredInvoices = mockInvoices.filter(invoice => {
    const matchesSearch = invoice.customerName.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         invoice.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         invoice.orderNumber.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || invoice.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const totalInvoices = mockInvoices.length;
  const totalAmount = mockInvoices.reduce((sum, inv) => sum + inv.totalAmount, 0);
  const paidAmount = mockInvoices.reduce((sum, inv) => sum + inv.paidAmount, 0);
  const pendingAmount = mockInvoices.reduce((sum, inv) => sum + inv.remainingAmount, 0);
  const overdueInvoices = mockInvoices.filter(inv => inv.status === 'overdue').length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Invoice Management</h1>
          <p className="text-muted-foreground">
            Manage and track all your invoices
          </p>
        </div>
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              Create Invoice
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Create New Invoice</DialogTitle>
              <DialogDescription>
                Create a new invoice for your customer
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="customer">Customer Name</Label>
                  <Input id="customer" placeholder="Enter customer name" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="phone">Phone</Label>
                  <Input id="phone" placeholder="01XXXXXXXXX" />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="amount">Amount</Label>
                  <Input id="amount" type="number" placeholder="Enter amount" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="dueDate">Due Date</Label>
                  <Input id="dueDate" type="date" />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="notes">Notes</Label>
                <Textarea id="notes" placeholder="Additional notes" />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                Cancel
              </Button>
              <Button onClick={() => setIsCreateDialogOpen(false)}>
                Create Invoice
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Invoices</CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{totalInvoices}</div>
            <p className="text-xs text-muted-foreground">
              All time invoices
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Amount</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(totalAmount)}</div>
            <p className="text-xs text-muted-foreground">
              Total invoice value
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Paid Amount</CardTitle>
            <CheckCircle className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{formatCurrency(paidAmount)}</div>
            <p className="text-xs text-muted-foreground">
              Successfully collected
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pending Amount</CardTitle>
            <Clock className="h-4 w-4 text-yellow-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">{formatCurrency(pendingAmount)}</div>
            <p className="text-xs text-muted-foreground">
              Awaiting payment
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Overdue</CardTitle>
            <AlertTriangle className="h-4 w-4 text-red-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{overdueInvoices}</div>
            <p className="text-xs text-muted-foreground">
              Require attention
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
                  placeholder="Search by customer, invoice number, or order number..."
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
                <SelectItem value="paid">Paid</SelectItem>
                <SelectItem value="partial">Partial</SelectItem>
                <SelectItem value="overdue">Overdue</SelectItem>
                <SelectItem value="cancelled">Cancelled</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Invoices Table */}
      <Card>
        <CardHeader>
          <CardTitle>Invoices</CardTitle>
          <CardDescription>
            {filteredInvoices.length} invoices found
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Invoice #</TableHead>
                <TableHead>Customer</TableHead>
                <TableHead>Issue Date</TableHead>
                <TableHead>Due Date</TableHead>
                <TableHead>Amount</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredInvoices.map((invoice) => (
                <TableRow key={invoice.id}>
                  <TableCell>
                    <div>
                      <p className="font-medium">{invoice.id}</p>
                      <p className="text-sm text-muted-foreground">{invoice.orderNumber}</p>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="font-medium">{invoice.customerName}</p>
                      <p className="text-sm text-muted-foreground">{invoice.customerPhone}</p>
                    </div>
                  </TableCell>
                  <TableCell>{formatDate(invoice.issueDate)}</TableCell>
                  <TableCell>
                    <div>
                      <p className={cn(
                        "font-medium",
                        invoice.status === 'overdue' && "text-red-600"
                      )}>
                        {formatDate(invoice.dueDate)}
                      </p>
                      {invoice.status === 'overdue' && (
                        <p className="text-xs text-red-600">
                          {getDaysOverdue(invoice.dueDate)} days overdue
                        </p>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="font-medium">{formatCurrency(invoice.totalAmount)}</p>
                      {invoice.remainingAmount > 0 && (
                        <p className="text-sm text-muted-foreground">
                          Remaining: {formatCurrency(invoice.remainingAmount)}
                        </p>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>{getStatusBadge(invoice.status)}</TableCell>
                  <TableCell>
                    <div className="flex gap-1">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => {
                          setSelectedInvoice(invoice);
                          setIsViewDialogOpen(true);
                        }}
                        title="View Invoice"
                      >
                        <Eye className="h-4 w-4" />
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        title="Edit Invoice"
                      >
                        <Edit className="h-4 w-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Invoice Details Dialog */}
      <Dialog open={isViewDialogOpen} onOpenChange={setIsViewDialogOpen}>
        <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
          {selectedInvoice && (
            <>
              <DialogHeader>
                <DialogTitle>Invoice Details - {selectedInvoice.id}</DialogTitle>
                <DialogDescription>
                  Complete invoice information and items
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-6">
                {/* Invoice Header */}
                <div className="grid grid-cols-2 gap-6">
                  <div>
                    <h3 className="font-semibold mb-2">Customer Information</h3>
                    <div className="space-y-1 text-sm">
                      <p><strong>Name:</strong> {selectedInvoice.customerName}</p>
                      <p><strong>Email:</strong> {selectedInvoice.customerEmail}</p>
                      <p><strong>Phone:</strong> {selectedInvoice.customerPhone}</p>
                      <p><strong>Address:</strong> {selectedInvoice.customerAddress}</p>
                    </div>
                  </div>
                  <div>
                    <h3 className="font-semibold mb-2">Invoice Information</h3>
                    <div className="space-y-1 text-sm">
                      <p><strong>Invoice #:</strong> {selectedInvoice.id}</p>
                      <p><strong>Order #:</strong> {selectedInvoice.orderNumber}</p>
                      <p><strong>Issue Date:</strong> {formatDate(selectedInvoice.issueDate)}</p>
                      <p><strong>Due Date:</strong> {formatDate(selectedInvoice.dueDate)}</p>
                      <p><strong>Status:</strong> {getStatusBadge(selectedInvoice.status)}</p>
                      <p><strong>Payment:</strong> {getPaymentStatusBadge(selectedInvoice.paymentStatus)}</p>
                    </div>
                  </div>
                </div>

                {/* Invoice Items */}
                <div>
                  <h3 className="font-semibold mb-2">Invoice Items</h3>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Item</TableHead>
                        <TableHead>Quantity</TableHead>
                        <TableHead>Price</TableHead>
                        <TableHead>Total</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {selectedInvoice.items.map((item, index) => (
                        <TableRow key={index}>
                          <TableCell>{item.name}</TableCell>
                          <TableCell>{item.quantity}</TableCell>
                          <TableCell>{formatCurrency(item.price)}</TableCell>
                          <TableCell>{formatCurrency(item.total)}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>

                {/* Invoice Summary */}
                <div className="border-t pt-4">
                  <div className="flex justify-end">
                    <div className="w-64 space-y-2">
                      <div className="flex justify-between">
                        <span>Subtotal:</span>
                        <span>{formatCurrency(selectedInvoice.subtotal)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Tax (15%):</span>
                        <span>{formatCurrency(selectedInvoice.tax)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Discount:</span>
                        <span>-{formatCurrency(selectedInvoice.discount)}</span>
                      </div>
                      <div className="flex justify-between font-bold text-lg border-t pt-2">
                        <span>Total:</span>
                        <span>{formatCurrency(selectedInvoice.totalAmount)}</span>
                      </div>
                      <div className="flex justify-between text-green-600">
                        <span>Paid:</span>
                        <span>{formatCurrency(selectedInvoice.paidAmount)}</span>
                      </div>
                      {selectedInvoice.remainingAmount > 0 && (
                        <div className="flex justify-between text-red-600 font-medium">
                          <span>Remaining:</span>
                          <span>{formatCurrency(selectedInvoice.remainingAmount)}</span>
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                {/* Notes */}
                {selectedInvoice.notes && (
                  <div>
                    <h3 className="font-semibold mb-2">Notes</h3>
                    <p className="text-sm text-muted-foreground">{selectedInvoice.notes}</p>
                  </div>
                )}
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setIsViewDialogOpen(false)}>
                  Close
                </Button>
                <Button>
                  <Download className="h-4 w-4 mr-2" />
                  Download PDF
                </Button>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default Invoices;