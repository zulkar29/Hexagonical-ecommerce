import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Calendar, Download, FileText } from 'lucide-react';
import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';

const mockSalesData = [
  { date: '2025-01-15', sales: 12000, orders: 18, avgOrder: 667 },
  { date: '2025-01-16', sales: 14500, orders: 22, avgOrder: 659 },
  { date: '2025-01-17', sales: 9800, orders: 14, avgOrder: 700 },
  { date: '2025-01-18', sales: 17200, orders: 25, avgOrder: 688 },
  { date: '2025-01-19', sales: 15800, orders: 21, avgOrder: 752 },
  { date: '2025-01-20', sales: 13400, orders: 19, avgOrder: 705 },
  { date: '2025-01-21', sales: 16200, orders: 23, avgOrder: 704 }
];

const SalesReportPage = () => {
  const [dateRange, setDateRange] = useState({ from: '2025-01-15', to: '2025-01-21' });
  const filteredData = mockSalesData; // In real app, filter by dateRange

  const totalSales = filteredData.reduce((sum, d) => sum + d.sales, 0);
  const totalOrders = filteredData.reduce((sum, d) => sum + d.orders, 0);
  const avgOrderValue = totalOrders ? Math.round(totalSales / totalOrders) : 0;

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  return (
    <div className="space-y-6 p-2 md:p-6">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-2xl font-bold">Sales Report</h1>
          <p className="text-muted-foreground">Overview of sales performance for your selected period.</p>
        </div>
        <div className="flex space-x-2">
          <Button variant="outline" size="sm"><FileText className="h-4 w-4 mr-2" /> Export PDF</Button>
          <Button variant="outline" size="sm"><Download className="h-4 w-4 mr-2" /> Export Excel</Button>
        </div>
      </div>

      {/* Date Range Filter */}
      <Card className="mb-4">
        <CardContent className="flex items-center space-x-4 p-4">
          <Calendar className="h-5 w-5 text-muted-foreground" />
          <span>Date Range:</span>
          <Input
            type="date"
            value={dateRange.from}
            onChange={e => setDateRange({ ...dateRange, from: e.target.value })}
            className="w-36"
          />
          <span>-</span>
          <Input
            type="date"
            value={dateRange.to}
            onChange={e => setDateRange({ ...dateRange, to: e.target.value })}
            className="w-36"
          />
        </CardContent>
      </Card>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <Card>
          <CardContent className="p-4">
            <CardTitle>Total Sales</CardTitle>
            <div className="text-2xl font-bold mt-2">{formatCurrency(totalSales)}</div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <CardTitle>Total Orders</CardTitle>
            <div className="text-2xl font-bold mt-2">{totalOrders}</div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <CardTitle>Avg. Order Value</CardTitle>
            <div className="text-2xl font-bold mt-2">{formatCurrency(avgOrderValue)}</div>
          </CardContent>
        </Card>
      </div>

      {/* Sales Trend Chart */}
      <Card className="mb-4">
        <CardHeader>
          <CardTitle>Sales Trend</CardTitle>
          <CardDescription>Daily sales for the selected period</CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <AreaChart data={filteredData}>
              <defs>
                <linearGradient id="salesGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#22c55e" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
              <XAxis dataKey="date" />
              <YAxis />
              <Tooltip formatter={(value) => [formatCurrency(value), 'Sales']} />
              <Area type="monotone" dataKey="sales" stroke="#22c55e" fillOpacity={1} fill="url(#salesGradient)" />
            </AreaChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Sales Table */}
      <Card>
        <CardHeader>
          <CardTitle>Sales by Day</CardTitle>
        </CardHeader>
        <CardContent className="overflow-x-auto p-0">
          <table className="min-w-full text-sm">
            <thead>
              <tr className="bg-muted">
                <th className="px-4 py-2 text-left">Date</th>
                <th className="px-4 py-2 text-right">Sales</th>
                <th className="px-4 py-2 text-right">Orders</th>
                <th className="px-4 py-2 text-right">Avg. Order Value</th>
              </tr>
            </thead>
            <tbody>
              {filteredData.map((row) => (
                <tr key={row.date} className="border-b last:border-0">
                  <td className="px-4 py-2">{row.date}</td>
                  <td className="px-4 py-2 text-right">{formatCurrency(row.sales)}</td>
                  <td className="px-4 py-2 text-right">{row.orders}</td>
                  <td className="px-4 py-2 text-right">{formatCurrency(row.avgOrder)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
};

export default SalesReportPage; 