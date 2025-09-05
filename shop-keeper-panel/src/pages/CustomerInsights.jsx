import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '@/components/ui/select';
import { Download, FileText, Users, Star, TrendingUp, Gift } from 'lucide-react';
import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';

const mockCustomers = [
  { id: 1, name: 'Fatima Khan', phone: '01712345678', email: 'fatima.khan@gmail.com', status: 'VIP', totalPurchases: 45780, totalOrders: 67, lastPurchase: '2025-01-20' },
  { id: 2, name: 'Ahmed Ali', phone: '01798765432', email: 'ahmed.ali@gmail.com', status: 'Regular', totalPurchases: 18500, totalOrders: 34, lastPurchase: '2025-01-19' },
  { id: 3, name: 'Rashida Begum', phone: '01656789012', email: 'rashida.begum@gmail.com', status: 'VIP', totalPurchases: 32000, totalOrders: 45, lastPurchase: '2025-01-18' },
  { id: 4, name: 'Mohammad Hasan', phone: '01534567890', email: 'hasan.mohammad@gmail.com', status: 'Silver', totalPurchases: 15600, totalOrders: 28, lastPurchase: '2025-01-17' },
  { id: 5, name: 'Nasir Ahmed', phone: '01623456789', email: 'nasir.ahmed@outlook.com', status: 'Gold', totalPurchases: 24500, totalOrders: 38, lastPurchase: '2025-01-16' },
];

const customerGrowth = [
  { month: 'Aug', customers: 120 },
  { month: 'Sep', customers: 140 },
  { month: 'Oct', customers: 170 },
  { month: 'Nov', customers: 200 },
  { month: 'Dec', customers: 230 },
  { month: 'Jan', customers: 260 },
];

const CustomerInsights = () => {
  const [status, setStatus] = useState('all');
  const [search, setSearch] = useState('');

  const filtered = mockCustomers.filter(c =>
    (status === 'all' || c.status === status) &&
    (c.name.toLowerCase().includes(search.toLowerCase()) || c.phone.includes(search) || c.email.toLowerCase().includes(search.toLowerCase()))
  );

  const totalCustomers = mockCustomers.length;
  const newCustomers = 3; // mock
  const vipCustomers = mockCustomers.filter(c => c.status === 'VIP').length;
  const totalRevenue = mockCustomers.reduce((sum, c) => sum + c.totalPurchases, 0);

  const getStatusBadge = (status) => {
    switch (status) {
      case 'VIP': return <Badge variant="secondary"><Star className="h-4 w-4 mr-1 inline" />VIP</Badge>;
      case 'Gold': return <Badge variant="outline" className="text-yellow-700 border-yellow-300"><Gift className="h-4 w-4 mr-1 inline" />Gold</Badge>;
      case 'Silver': return <Badge variant="outline" className="text-gray-700 border-gray-300">Silver</Badge>;
      default: return <Badge variant="outline">Regular</Badge>;
    }
  };

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Customer Insights</h1>
          <p className="text-muted-foreground">Analyze your customer base and discover growth opportunities.</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm"><FileText className="h-4 w-4 mr-2" />Export PDF</Button>
          <Button variant="outline" size="sm"><Download className="h-4 w-4 mr-2" />Export Excel</Button>
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <Users className="h-8 w-8 text-primary" />
            <div>
              <p className="text-sm text-muted-foreground">Total Customers</p>
              <p className="text-2xl font-bold">{totalCustomers}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <TrendingUp className="h-8 w-8 text-green-600" />
            <div>
              <p className="text-sm text-muted-foreground">New Customers</p>
              <p className="text-2xl font-bold">{newCustomers}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <Star className="h-8 w-8 text-yellow-500" />
            <div>
              <p className="text-sm text-muted-foreground">VIP Customers</p>
              <p className="text-2xl font-bold">{vipCustomers}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <Gift className="h-8 w-8 text-indigo-500" />
            <div>
              <p className="text-sm text-muted-foreground">Total Revenue</p>
              <p className="text-2xl font-bold">{formatCurrency(totalRevenue)}</p>
            </div>
          </CardContent>
        </Card>
      </div>
      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Customer Growth</CardTitle>
          <CardDescription>Monthly new customer trend</CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={220}>
            <AreaChart data={customerGrowth}>
              <defs>
                <linearGradient id="growthGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#22c55e" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
              <XAxis dataKey="month" />
              <YAxis />
              <Tooltip />
              <Area type="monotone" dataKey="customers" stroke="#22c55e" fillOpacity={1} fill="url(#growthGradient)" />
            </AreaChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-6">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-2 mb-4">
            <div className="flex gap-2 items-center">
              <Input
                placeholder="Search customers..."
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
                  <SelectItem value="VIP">VIP</SelectItem>
                  <SelectItem value="Gold">Gold</SelectItem>
                  <SelectItem value="Silver">Silver</SelectItem>
                  <SelectItem value="Regular">Regular</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Customer</TableHead>
                <TableHead>Contact</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Total Purchases</TableHead>
                <TableHead>Orders</TableHead>
                <TableHead>Last Purchase</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map(customer => (
                <TableRow key={customer.id}>
                  <TableCell>
                    <div className="flex items-center space-x-3">
                      <Avatar className="h-8 w-8">
                        <AvatarFallback>{customer.name.split(' ').map(n => n[0]).join('')}</AvatarFallback>
                      </Avatar>
                      <span className="font-medium">{customer.name}</span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="text-sm">{customer.phone}</p>
                      <p className="text-xs text-muted-foreground">{customer.email}</p>
                    </div>
                  </TableCell>
                  <TableCell>{getStatusBadge(customer.status)}</TableCell>
                  <TableCell>{formatCurrency(customer.totalPurchases)}</TableCell>
                  <TableCell>{customer.totalOrders}</TableCell>
                  <TableCell className="text-sm">{customer.lastPurchase}</TableCell>
                  <TableCell>
                    <Button size="sm" variant="outline">Details</Button>
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

export default CustomerInsights; 