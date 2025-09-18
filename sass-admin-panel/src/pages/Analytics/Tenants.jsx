import React from 'react';
import {
  Users,
  UserPlus,
  TrendingUp,
  TrendingDown,
  Package,
  PieChart,
  CheckCircle,
  AlertTriangle,
  Clock
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

const summary = [
  { label: 'Active Tenants', value: 156, icon: Users, color: 'text-green-600' },
  { label: 'New Signups', value: 24, icon: UserPlus, color: 'text-blue-600' },
  { label: 'Churned', value: 5, icon: TrendingDown, color: 'text-red-600' },
  { label: 'Trial Tenants', value: 12, icon: Clock, color: 'text-yellow-600' }
];

const planDistribution = [
  { plan: 'Starter', count: 89, color: 'text-green-600' },
  { plan: 'Business', count: 45, color: 'text-blue-600' },
  { plan: 'Enterprise', count: 22, color: 'text-yellow-600' }
];

const recentTenants = [
  { id: 1, name: 'Rahman Electronics', owner: 'Abdul Rahman', plan: 'Business', status: 'Active', joined: '2024-07-20' },
  { id: 2, name: 'Fatima Fashion', owner: 'Fatima Khatun', plan: 'Starter', status: 'Trial', joined: '2024-07-22' },
  { id: 3, name: 'Modern Pharmacy', owner: 'Dr. Ahmed Ali', plan: 'Enterprise', status: 'Active', joined: '2024-07-18' },
  { id: 4, name: 'Green Grocers', owner: 'Mohammad Hasan', plan: 'Business', status: 'Suspended', joined: '2024-07-15' },
  { id: 5, name: 'Tech Solutions', owner: 'Rashida Begum', plan: 'Enterprise', status: 'Active', joined: '2024-07-21' }
];

export default function TenantAnalytics() {
  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center gap-4 px-6">
          <PieChart className="h-5 w-5 text-primary" />
          <h1 className="text-xl font-semibold">Tenant Analytics</h1>
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
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Tenants Growth Chart (Mock) */}
        <Card>
          <CardHeader>
            <CardTitle>Tenants Growth Trend</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-64 flex items-center justify-center text-muted-foreground">
              {/* Replace with real chart in production */}
              <span>Tenants Growth Chart (Mock)</span>
            </div>
          </CardContent>
        </Card>

        {/* Plan Distribution Breakdown */}
        <Card>
          <CardHeader>
            <CardTitle>Plan Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {planDistribution.map((plan) => (
                <div key={plan.plan} className="flex flex-col items-center p-4 border rounded-lg">
                  <Package className={`h-6 w-6 mb-2 ${plan.color}`} />
                  <span className="font-semibold">{plan.plan}</span>
                  <span className="text-lg font-bold">{plan.count} tenants</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Tenants Table */}
        <Card>
          <CardHeader>
            <CardTitle>Recent Tenants</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow className="bg-muted/50">
                    <TableHead className="p-2">Name</TableHead>
                    <TableHead className="p-2">Owner</TableHead>
                    <TableHead className="p-2">Plan</TableHead>
                    <TableHead className="p-2">Status</TableHead>
                    <TableHead className="p-2">Joined</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {recentTenants.map((t) => (
                    <TableRow key={t.id} className="border-t">
                      <TableCell className="p-2">{t.name}</TableCell>
                      <TableCell className="p-2">{t.owner}</TableCell>
                      <TableCell className="p-2">{t.plan}</TableCell>
                      <TableCell className="p-2">
                        {t.status === 'Active' ? (
                          <Badge variant="success">Active</Badge>
                        ) : t.status === 'Trial' ? (
                          <Badge variant="secondary">Trial</Badge>
                        ) : t.status === 'Suspended' ? (
                          <Badge variant="destructive">Suspended</Badge>
                        ) : (
                          <Badge>{t.status}</Badge>
                        )}
                      </TableCell>
                      <TableCell className="p-2">{t.joined}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
} 