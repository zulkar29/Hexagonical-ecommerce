import React from 'react';
import {
  TrendingUp,
  DollarSign,
  BarChart3,
  PieChart,
  Users,
  Package,
  Calendar,
  CreditCard,
  CheckCircle,
  AlertTriangle
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';

const summary = [
  { label: 'MRR', value: '৳245,000', icon: TrendingUp, color: 'text-green-600', change: '+12.5%' },
  { label: 'Total Revenue', value: '৳2,450,000', icon: DollarSign, color: 'text-blue-600', change: '+8.2%' },
  { label: 'Growth', value: '+8.2%', icon: BarChart3, color: 'text-green-600', change: '+1.1%' },
  { label: 'Churn', value: '2.8%', icon: AlertTriangle, color: 'text-red-600', change: '-0.5%' }
];

const planBreakdown = [
  { plan: 'Starter', revenue: 178000, color: 'text-green-600', count: 89 },
  { plan: 'Business', revenue: 225000, color: 'text-blue-600', count: 45 },
  { plan: 'Enterprise', revenue: 220000, color: 'text-yellow-600', count: 22 }
];

const recentPayments = [
  { id: 1, tenant: 'Rahman Electronics', amount: 5000, method: 'bKash', status: 'Completed', date: '2024-07-24' },
  { id: 2, tenant: 'Modern Pharmacy', amount: 10000, method: 'Bank Transfer', status: 'Completed', date: '2024-07-23' },
  { id: 3, tenant: 'Fatima Fashion', amount: 2000, method: 'Nagad', status: 'Pending', date: '2024-07-22' },
  { id: 4, tenant: 'Green Grocers', amount: 5000, method: 'bKash', status: 'Failed', date: '2024-07-21' },
  { id: 5, tenant: 'Tech Solutions', amount: 10000, method: 'Rocket', status: 'Completed', date: '2024-07-20' }
];

export default function RevenueAnalytics() {
  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center gap-4 px-6">
          <TrendingUp className="h-5 w-5 text-primary" />
          <h1 className="text-xl font-semibold">Revenue Analytics</h1>
          <Badge variant="outline" className="ml-2">This Month</Badge>
        </div>
      </div>
      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {summary.map((item) => (
            <Card key={item.label}>
              <CardContent className="p-6 flex items-center gap-4">
                <item.icon className={`h-8 w-8 ${item.color}`} />
                <div>
                  <p className="text-sm text-muted-foreground">{item.label}</p>
                  <p className="text-2xl font-bold">{item.value}</p>
                  <span className={`text-xs font-medium ${item.color}`}>{item.change}</span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Revenue Chart (Mock) */}
        <Card>
          <CardHeader>
            <CardTitle>Revenue Trend</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-64 flex items-center justify-center text-muted-foreground">
              {/* Replace with real chart in production */}
              <span>Revenue Chart (Mock)</span>
            </div>
          </CardContent>
        </Card>

        {/* Revenue by Plan Breakdown */}
        <Card>
          <CardHeader>
            <CardTitle>Revenue by Plan</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {planBreakdown.map((plan) => (
                <div key={plan.plan} className="flex flex-col items-center p-4 border rounded-lg">
                  <Package className={`h-6 w-6 mb-2 ${plan.color}`} />
                  <span className="font-semibold">{plan.plan}</span>
                  <span className="text-lg font-bold">৳{plan.revenue.toLocaleString()}</span>
                  <span className="text-xs text-muted-foreground">{plan.count} tenants</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Payments Table */}
        <Card>
          <CardHeader>
            <CardTitle>Recent Payments</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <table className="min-w-full text-sm border">
                <thead>
                  <tr className="bg-muted/50">
                    <th className="p-2 text-left font-medium">Tenant</th>
                    <th className="p-2 text-left font-medium">Amount</th>
                    <th className="p-2 text-left font-medium">Method</th>
                    <th className="p-2 text-left font-medium">Status</th>
                    <th className="p-2 text-left font-medium">Date</th>
                  </tr>
                </thead>
                <tbody>
                  {recentPayments.map((p) => (
                    <tr key={p.id} className="border-t">
                      <td className="p-2">{p.tenant}</td>
                      <td className="p-2">৳{p.amount.toLocaleString()}</td>
                      <td className="p-2">{p.method}</td>
                      <td className="p-2">
                        {p.status === 'Completed' ? (
                          <Badge variant="success">Completed</Badge>
                        ) : p.status === 'Pending' ? (
                          <Badge variant="secondary">Pending</Badge>
                        ) : (
                          <Badge variant="destructive">Failed</Badge>
                        )}
                      </td>
                      <td className="p-2">{p.date}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
} 