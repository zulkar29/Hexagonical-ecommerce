import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '@/components/ui/select';
import { Download, FileText, CheckCircle, Clock, XCircle, Eye, CreditCard, Receipt, DollarSign } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const mockOrders = [
  { id: 1234, customer: 'Fatima Khan', date: '2025-01-20', amount: 2450, status: 'completed', paymentStatus: 'paid', paymentMethod: 'bKash', invoiceStatus: 'generated' },
  { id: 1235, customer: 'Ahmed Ali', date: '2025-01-19', amount: 1890, status: 'pending', paymentStatus: 'unpaid', paymentMethod: 'Cash', invoiceStatus: 'pending' },
  { id: 1236, customer: 'Rashida Begum', date: '2025-01-18', amount: 3200, status: 'completed', paymentStatus: 'paid', paymentMethod: 'Nagad', invoiceStatus: 'generated' },
  { id: 1237, customer: 'Mohammad Hasan', date: '2025-01-17', amount: 1650, status: 'cancelled', paymentStatus: 'refunded', paymentMethod: 'Card', invoiceStatus: 'cancelled' },
  { id: 1238, customer: 'Nasir Ahmed', date: '2025-01-16', amount: 2780, status: 'completed', paymentStatus: 'paid', paymentMethod: 'Cash', invoiceStatus: 'generated' },
];

const Orders = () => {
  const [status, setStatus] = useState('all');
  const [search, setSearch] = useState('');
  const navigate = useNavigate();

  const filtered = mockOrders.filter(o =>
    (status === 'all' || o.status === status) &&
    (o.customer.toLowerCase().includes(search.toLowerCase()) || o.id.toString().includes(search))
  );

  const totalOrders = mockOrders.length;
  const completed = mockOrders.filter(o => o.status === 'completed').length;
  const pending = mockOrders.filter(o => o.status === 'pending').length;
  const cancelled = mockOrders.filter(o => o.status === 'cancelled').length;

  const getStatusBadge = (status) => {
    switch (status) {
      case 'completed': return <Badge variant="secondary"><CheckCircle className="h-4 w-4 mr-1 inline" />Completed</Badge>;
      case 'pending': return <Badge variant="outline"><Clock className="h-4 w-4 mr-1 inline" />Pending</Badge>;
      case 'cancelled': return <Badge variant="destructive"><XCircle className="h-4 w-4 mr-1 inline" />Cancelled</Badge>;
      default: return <Badge variant="outline">Unknown</Badge>;
    }
  };

  const getPaymentStatusBadge = (status) => {
    switch (status) {
      case 'paid': return <Badge variant="secondary" className="bg-green-100 text-green-800"><CheckCircle className="h-3 w-3 mr-1 inline" />Paid</Badge>;
      case 'unpaid': return <Badge variant="destructive"><Clock className="h-3 w-3 mr-1 inline" />Unpaid</Badge>;
      case 'refunded': return <Badge variant="outline"><XCircle className="h-3 w-3 mr-1 inline" />Refunded</Badge>;
      default: return <Badge variant="outline">Unknown</Badge>;
    }
  };

  const getInvoiceStatusBadge = (status) => {
    switch (status) {
      case 'generated': return <Badge variant="secondary" className="bg-blue-100 text-blue-800"><Receipt className="h-3 w-3 mr-1 inline" />Generated</Badge>;
      case 'pending': return <Badge variant="outline"><Clock className="h-3 w-3 mr-1 inline" />Pending</Badge>;
      case 'cancelled': return <Badge variant="destructive"><XCircle className="h-3 w-3 mr-1 inline" />Cancelled</Badge>;
      default: return <Badge variant="outline">Unknown</Badge>;
    }
  };

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Orders</h1>
          <p className="text-muted-foreground">Manage and track all customer orders.</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm"><FileText className="h-4 w-4 mr-2" />Export PDF</Button>
          <Button variant="outline" size="sm"><Download className="h-4 w-4 mr-2" />Export Excel</Button>
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <CheckCircle className="h-8 w-8 text-primary" />
            <div>
              <p className="text-sm text-muted-foreground">Total Orders</p>
              <p className="text-2xl font-bold">{totalOrders}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <Clock className="h-8 w-8 text-yellow-500" />
            <div>
              <p className="text-sm text-muted-foreground">Pending</p>
              <p className="text-2xl font-bold">{pending}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <XCircle className="h-8 w-8 text-destructive" />
            <div>
              <p className="text-sm text-muted-foreground">Cancelled</p>
              <p className="text-2xl font-bold">{cancelled}</p>
            </div>
          </CardContent>
        </Card>
      </div>
      <Card>
        <CardContent className="p-6">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-2 mb-4">
            <div className="flex gap-2 items-center">
              <Input
                placeholder="Search orders..."
                value={search}
                onChange={e => setSearch(e.target.value)}
                className="max-w-xs"
              />
              <Select value={status} onValueChange={setStatus}>
                <SelectTrigger className="w-40">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="cancelled">Cancelled</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Order ID</TableHead>
                <TableHead>Customer</TableHead>
                <TableHead>Date</TableHead>
                <TableHead>Amount</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Payment</TableHead>
                <TableHead>Invoice</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map(order => (
                <TableRow key={order.id}>
                  <TableCell>#{order.id}</TableCell>
                  <TableCell>{order.customer}</TableCell>
                  <TableCell>{order.date}</TableCell>
                  <TableCell>{formatCurrency(order.amount)}</TableCell>
                  <TableCell>{getStatusBadge(order.status)}</TableCell>
                  <TableCell>
                    <div className="space-y-1">
                      {getPaymentStatusBadge(order.paymentStatus)}
                      <div className="text-xs text-muted-foreground">{order.paymentMethod}</div>
                    </div>
                  </TableCell>
                  <TableCell>{getInvoiceStatusBadge(order.invoiceStatus)}</TableCell>
                  <TableCell>
                    <div className="flex gap-1">
                      <Button size="sm" variant="outline" onClick={() => navigate(`/orders/${order.id}`)}><Eye className="h-4 w-4" /></Button>
                      {order.invoiceStatus === 'generated' && (
                        <Button size="sm" variant="outline" title="Download Invoice"><FileText className="h-4 w-4" /></Button>
                      )}
                      {order.paymentStatus === 'unpaid' && (
                        <Button size="sm" variant="outline" title="Process Payment"><CreditCard className="h-4 w-4" /></Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default Orders; 