import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '@/components/ui/select';
import { Download, FileText, AlertTriangle, Package, XCircle, CheckCircle } from 'lucide-react';

const mockInventory = [
  { id: 1, name: 'Basmati Rice 5kg', category: 'Rice & Grains', stock: 5, minStock: 20, status: 'critical' },
  { id: 2, name: 'Onion 1kg', category: 'Vegetables', stock: 25, minStock: 10, status: 'ok' },
  { id: 3, name: 'Milk 1L', category: 'Dairy', stock: 0, minStock: 10, status: 'out' },
  { id: 4, name: 'Potato 1kg', category: 'Vegetables', stock: 8, minStock: 15, status: 'low' },
  { id: 5, name: 'Cooking Oil 1L', category: 'Cooking', stock: 20, minStock: 10, status: 'ok' },
  { id: 6, name: 'Sugar 1kg', category: 'Pantry', stock: 3, minStock: 20, status: 'critical' },
  { id: 7, name: 'Tea Leaves 500g', category: 'Beverages', stock: 15, minStock: 10, status: 'ok' },
  { id: 8, name: 'Salt 1kg', category: 'Pantry', stock: 12, minStock: 25, status: 'low' },
];

const categories = [
  'All',
  'Rice & Grains',
  'Vegetables',
  'Dairy',
  'Cooking',
  'Pantry',
  'Beverages',
];

const InventoryReport = () => {
  const [category, setCategory] = useState('All');
  const [showLowStock, setShowLowStock] = useState(false);
  const [search, setSearch] = useState('');

  const filtered = mockInventory.filter(item =>
    (category === 'All' || item.category === category) &&
    (!showLowStock || item.status === 'critical' || item.status === 'low' || item.status === 'out') &&
    (item.name.toLowerCase().includes(search.toLowerCase()) || item.category.toLowerCase().includes(search.toLowerCase()))
  );

  const totalProducts = mockInventory.length;
  const lowStock = mockInventory.filter(i => i.status === 'critical' || i.status === 'low').length;
  const outOfStock = mockInventory.filter(i => i.status === 'out').length;

  const getStatus = (status) => {
    switch (status) {
      case 'critical':
        return <Badge variant="destructive"><AlertTriangle className="h-4 w-4 mr-1 inline" />Critical</Badge>;
      case 'low':
        return <Badge variant="secondary"><AlertTriangle className="h-4 w-4 mr-1 inline" />Low</Badge>;
      case 'out':
        return <Badge variant="outline" className="text-destructive"><XCircle className="h-4 w-4 mr-1 inline" />Out</Badge>;
      default:
        return <Badge variant="outline" className="text-green-700 border-green-300"><CheckCircle className="h-4 w-4 mr-1 inline" />OK</Badge>;
    }
  };

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Inventory Report</h1>
          <p className="text-muted-foreground">Monitor your stock levels and identify low or out-of-stock items.</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm"><FileText className="h-4 w-4 mr-2" />Export PDF</Button>
          <Button variant="outline" size="sm"><Download className="h-4 w-4 mr-2" />Export Excel</Button>
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <Package className="h-8 w-8 text-primary" />
            <div>
              <p className="text-sm text-muted-foreground">Total Products</p>
              <p className="text-2xl font-bold">{totalProducts}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <AlertTriangle className="h-8 w-8 text-orange-500" />
            <div>
              <p className="text-sm text-muted-foreground">Low Stock</p>
              <p className="text-2xl font-bold">{lowStock}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <XCircle className="h-8 w-8 text-destructive" />
            <div>
              <p className="text-sm text-muted-foreground">Out of Stock</p>
              <p className="text-2xl font-bold">{outOfStock}</p>
            </div>
          </CardContent>
        </Card>
      </div>
      <Card>
        <CardContent className="p-6">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-2 mb-4">
            <div className="flex gap-2 items-center">
              <Input
                placeholder="Search inventory..."
                value={search}
                onChange={e => setSearch(e.target.value)}
                className="max-w-xs"
              />
              <Select value={category} onValueChange={setCategory}>
                <SelectTrigger className="w-40">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {categories.map(cat => (
                    <SelectItem key={cat} value={cat}>{cat}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Button
                size="sm"
                variant={showLowStock ? 'default' : 'outline'}
                onClick={() => setShowLowStock(v => !v)}
              >
                Show Low/Out Stock
              </Button>
            </div>
          </div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Product</TableHead>
                <TableHead>Category</TableHead>
                <TableHead>Stock</TableHead>
                <TableHead>Min Stock</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map(item => (
                <TableRow key={item.id}>
                  <TableCell>{item.name}</TableCell>
                  <TableCell>{item.category}</TableCell>
                  <TableCell>
                    <Badge variant={item.stock < 10 ? 'destructive' : 'default'}>{item.stock}</Badge>
                  </TableCell>
                  <TableCell>{item.minStock}</TableCell>
                  <TableCell>{getStatus(item.status)}</TableCell>
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

export default InventoryReport; 