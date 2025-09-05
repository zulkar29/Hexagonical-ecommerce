import React from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { useNavigate, useParams } from 'react-router-dom';
import { ArrowLeft, CheckCircle, Clock, XCircle, CreditCard, User, FileText, Download, DollarSign, Receipt } from 'lucide-react';

const mockOrder = {
  id: 1234,
  status: 'completed',
  date: '2025-01-20',
  customer: {
    name: 'Fatima Khan',
    phone: '01712345678',
    email: 'fatima.khan@gmail.com',
  },
  amount: 2450,
  payment: {
    method: 'bKash',
    status: 'paid',
    paidAmount: 2450,
    transactionId: 'TXN123456789',
    paidDate: '2025-01-20',
  },
  invoice: {
    status: 'generated',
    number: 'INV-2025-001234',
    generatedDate: '2025-01-20',
    dueDate: '2025-02-20',
  },
  products: [
    { id: 1, name: 'Basmati Rice 5kg', quantity: 2, price: 500 },
    { id: 2, name: 'Cooking Oil 1L', quantity: 3, price: 150 },
    { id: 3, name: 'Onion 1kg', quantity: 5, price: 100 },
  ],
};

const getStatusBadge = (status) => {
  switch (status) {
    case 'completed': return <Badge variant="secondary"><CheckCircle className="h-4 w-4 mr-1 inline" />Completed</Badge>;
    case 'pending': return <Badge variant="outline"><Clock className="h-4 w-4 mr-1 inline" />Pending</Badge>;
    case 'cancelled': return <Badge variant="destructive"><XCircle className="h-4 w-4 mr-1 inline" />Cancelled</Badge>;
    default: return <Badge variant="outline">Unknown</Badge>;
  }
};

const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

const OrderDetails = () => {
  const navigate = useNavigate();
  const order = mockOrder; // Replace with real lookup
  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      <Button variant="outline" onClick={() => navigate('/orders')} className="mb-4">
        <ArrowLeft className="h-4 w-4 mr-2" /> Back to Orders
      </Button>
      
      {/* Header Section */}
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Order #{order.id}</h1>
          <p className="text-muted-foreground">Placed on {order.date} • {getStatusBadge(order.status)}</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm">
            <Download className="h-4 w-4 mr-2" />Download Invoice
          </Button>
          <Button variant="outline" size="sm">
            <Receipt className="h-4 w-4 mr-2" />Print Invoice
          </Button>
          {order.payment.status !== 'paid' && (
            <Button variant="default" size="sm">
              <CreditCard className="h-4 w-4 mr-2" />Process Payment
            </Button>
          )}
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <DollarSign className="h-8 w-8 text-primary" />
            <div>
              <p className="text-sm text-muted-foreground">Total Amount</p>
              <p className="text-2xl font-bold">{formatCurrency(order.amount)}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <Receipt className="h-8 w-8 text-blue-600" />
            <div>
              <p className="text-sm text-muted-foreground">Items</p>
              <p className="text-2xl font-bold">{order.products.length}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <CreditCard className="h-8 w-8 text-green-600" />
            <div>
              <p className="text-sm text-muted-foreground">Payment Status</p>
              <p className="text-2xl font-bold">{order.payment.status === 'paid' ? 'Paid' : 'Unpaid'}</p>
            </div>
          </CardContent>
        </Card>
      </div>
      
      {/* Customer Information Card */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center gap-2 mb-4">
            <User className="h-5 w-5 text-primary" />
            <h3 className="text-lg font-semibold">Customer Information</h3>
          </div>
          <div className="flex items-center gap-4">
            <Avatar className="h-12 w-12">
              <AvatarFallback className="bg-primary/10 text-primary font-semibold">
                {order.customer.name.split(' ').map(n => n[0]).join('')}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1">
              <p className="font-semibold text-lg">{order.customer.name}</p>
              <div className="flex flex-col sm:flex-row sm:gap-6 gap-1 mt-1">
                <p className="text-sm text-muted-foreground">{order.customer.phone}</p>
                <p className="text-sm text-muted-foreground">{order.customer.email}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Payment & Invoice Cards */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center gap-2 mb-4">
              <CreditCard className="h-5 w-5 text-green-600" />
              <h3 className="text-lg font-semibold">Payment Details</h3>
            </div>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-muted-foreground font-medium">Status</span>
                <Badge variant={order.payment.status === 'paid' ? 'secondary' : 'destructive'} 
                       className={order.payment.status === 'paid' ? 'bg-green-100 text-green-800' : ''}>
                  {order.payment.status === 'paid' ? 'Paid' : 'Unpaid'}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-muted-foreground font-medium">Method</span>
                <span className="font-semibold">{order.payment.method}</span>
              </div>
              {order.payment.status === 'paid' && (
                <>
                  <div className="flex justify-between items-center">
                    <span className="text-muted-foreground font-medium">Amount</span>
                    <span className="font-bold text-green-700">{formatCurrency(order.payment.paidAmount)}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-muted-foreground font-medium">Transaction ID</span>
                    <span className="font-mono text-sm">{order.payment.transactionId}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-muted-foreground font-medium">Date</span>
                    <span>{order.payment.paidDate}</span>
                  </div>
                </>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center gap-2 mb-4">
              <FileText className="h-5 w-5 text-blue-600" />
              <h3 className="text-lg font-semibold">Invoice Details</h3>
            </div>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-muted-foreground font-medium">Status</span>
                <Badge variant={order.invoice.status === 'generated' ? 'secondary' : 'outline'} 
                       className={order.invoice.status === 'generated' ? 'bg-blue-100 text-blue-800' : ''}>
                  {order.invoice.status === 'generated' ? 'Generated' : 'Pending'}
                </Badge>
              </div>
              {order.invoice.status === 'generated' && (
                <>
                  <div className="flex justify-between items-center">
                    <span className="text-muted-foreground font-medium">Invoice #</span>
                    <span className="font-mono text-sm font-semibold">{order.invoice.number}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-muted-foreground font-medium">Generated</span>
                    <span>{order.invoice.generatedDate}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-muted-foreground font-medium">Due Date</span>
                    <span className="font-medium">{order.invoice.dueDate}</span>
                  </div>
                </>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Products Card */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center gap-2 mb-4">
            <Receipt className="h-5 w-5 text-primary" />
            <h3 className="text-lg font-semibold">Order Items</h3>
          </div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Product</TableHead>
                <TableHead className="text-center">Quantity</TableHead>
                <TableHead className="text-right">Price</TableHead>
                <TableHead className="text-right">Subtotal</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {order.products.map(product => (
                <TableRow key={product.id}>
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-primary/10 rounded-lg flex items-center justify-center">
                        <Receipt className="h-4 w-4 text-primary" />
                      </div>
                      <div>
                        <div className="font-medium">{product.name}</div>
                        <div className="text-sm text-muted-foreground">SKU: #{product.id}</div>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell className="text-center">
                    <Badge variant="outline">{product.quantity}</Badge>
                  </TableCell>
                  <TableCell className="text-right font-medium">
                    {formatCurrency(product.price)}
                  </TableCell>
                  <TableCell className="text-right font-bold">
                    {formatCurrency(product.price * product.quantity)}
                  </TableCell>
                </TableRow>
              ))}
              <TableRow className="border-t-2">
                <TableCell colSpan={3} className="font-semibold">Total</TableCell>
                <TableCell className="text-right font-bold text-lg text-primary">
                  {formatCurrency(order.amount)}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default OrderDetails; 