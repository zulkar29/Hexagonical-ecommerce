import React from 'react';
import { useParams } from 'react-router-dom';
import {
  Store,
  Users,
  DollarSign,
  Package,
  Calendar,
  Mail,
  Phone,
  CheckCircle,
  AlertTriangle
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

const mockTenants = {
  1: {
    name: 'Rahman Electronics',
    plan: 'Business',
    status: 'Active',
    revenue: 15000,
    products: 245,
    owner: 'Abdul Rahman',
    joined: '2024-07-20',
    email: 'rahman@email.com',
    phone: '01712345678',
    recent: [
      { id: 1, action: 'Paid Invoice', date: '2024-07-24', amount: 5000 },
      { id: 2, action: 'Added Product', date: '2024-07-23', amount: null },
      { id: 3, action: 'Support Ticket', date: '2024-07-22', amount: null }
    ]
  },
  // Add more mock tenants as needed
};

export default function TenantDetail() {
  const { id } = useParams();
  const tenant = mockTenants[id] || mockTenants[1]; // fallback to 1 for demo

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center gap-4 px-6">
          <Store className="h-5 w-5 text-primary" />
          <h1 className="text-xl font-semibold">{tenant.name}</h1>
          <Badge variant="outline" className="ml-2">{tenant.plan}</Badge>
          <Badge variant={tenant.status === 'Active' ? 'success' : tenant.status === 'Trial' ? 'secondary' : 'destructive'} className="ml-2">
            {tenant.status}
          </Badge>
        </div>
      </div>
      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardContent className="p-6 flex items-center gap-4">
              <DollarSign className="h-8 w-8 text-green-600" />
              <div>
                <p className="text-sm text-muted-foreground">Revenue</p>
                <p className="text-2xl font-bold">৳{tenant.revenue.toLocaleString()}</p>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-6 flex items-center gap-4">
              <Package className="h-8 w-8 text-blue-600" />
              <div>
                <p className="text-sm text-muted-foreground">Products</p>
                <p className="text-2xl font-bold">{tenant.products}</p>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-6 flex items-center gap-4">
              <Users className="h-8 w-8 text-purple-600" />
              <div>
                <p className="text-sm text-muted-foreground">Owner</p>
                <p className="text-2xl font-bold">{tenant.owner}</p>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-6 flex items-center gap-4">
              <Calendar className="h-8 w-8 text-yellow-600" />
              <div>
                <p className="text-sm text-muted-foreground">Joined</p>
                <p className="text-2xl font-bold">{tenant.joined}</p>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Contact Info & Recent Activity */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle>Contact Info</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <div className="flex items-center gap-2">
                <Mail className="h-4 w-4 text-muted-foreground" />
                <span>{tenant.email}</span>
              </div>
              <div className="flex items-center gap-2">
                <Phone className="h-4 w-4 text-muted-foreground" />
                <span>{tenant.phone}</span>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Recent Activity</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {tenant.recent.map((a) => (
                <div key={a.id} className="flex items-center gap-2">
                  {a.action === 'Paid Invoice' ? (
                    <CheckCircle className="h-4 w-4 text-green-600" />
                  ) : a.action === 'Support Ticket' ? (
                    <AlertTriangle className="h-4 w-4 text-yellow-600" />
                  ) : (
                    <Users className="h-4 w-4 text-blue-600" />
                  )}
                  <span>{a.action}</span>
                  {a.amount && <span className="ml-2 text-xs text-muted-foreground">৳{a.amount}</span>}
                  <span className="ml-auto text-xs text-muted-foreground">{a.date}</span>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}