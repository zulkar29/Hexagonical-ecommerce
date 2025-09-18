import React, { useState } from 'react';
import {
  ArrowLeft,
  CreditCard,
  Calendar,
  DollarSign,
  Users,
  Settings,
  AlertTriangle,
  CheckCircle,
  Edit,
  Trash2,
  Download,
  RefreshCw
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Link, useParams } from 'react-router-dom';

// Mock subscription data
const subscriptionData = {
  sub_001: {
    id: 'sub_001',
    tenant: 'TechCorp Solutions',
    tenantId: 'tenant_001',
    plan: 'Enterprise',
    status: 'Active',
    amount: 299,
    billingCycle: 'Monthly',
    nextBilling: '2024-02-15',
    startDate: '2023-08-15',
    features: ['Unlimited Users', 'Advanced Analytics', 'Priority Support', 'Custom Integrations'],
    paymentMethod: '**** **** **** 4242',
    billingHistory: [
      { date: '2024-01-15', amount: 299, status: 'Paid', invoice: 'INV-001' },
      { date: '2023-12-15', amount: 299, status: 'Paid', invoice: 'INV-002' },
      { date: '2023-11-15', amount: 299, status: 'Paid', invoice: 'INV-003' }
    ],
    usage: {
      users: { current: 45, limit: 'Unlimited' },
      storage: { current: '2.3 GB', limit: '100 GB' },
      apiCalls: { current: 15420, limit: 50000 }
    }
  }
};

export default function SubscriptionDetail() {
  const { id } = useParams();
  const [showCancelModal, setShowCancelModal] = useState(false);
  const subscription = subscriptionData[id] || subscriptionData['sub_001'];

  const getStatusBadge = (status) => {
    switch (status) {
      case 'Active':
        return <Badge variant="default" className="bg-green-600"><CheckCircle className="h-3 w-3 mr-1" />Active</Badge>;
      case 'Pending':
        return <Badge variant="secondary">Pending</Badge>;
      case 'Cancelled':
        return <Badge variant="destructive">Cancelled</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Link to="/subscriptions" className="text-muted-foreground hover:text-foreground">
              <ArrowLeft className="h-5 w-5" />
            </Link>
            <CreditCard className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Subscription Details</h1>
            <Badge variant="outline" className="ml-2">{subscription.id}</Badge>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm">
              <Edit className="h-4 w-4 mr-2" />
              Edit
            </Button>
            <Button variant="outline" size="sm" className="text-red-600 hover:text-red-700">
              <Trash2 className="h-4 w-4 mr-2" />
              Cancel
            </Button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Subscription Overview */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle>Subscription Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">Tenant</label>
                    <p className="text-lg font-semibold">{subscription.tenant}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">Status</label>
                    <div className="mt-1">{getStatusBadge(subscription.status)}</div>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">Plan</label>
                    <p className="text-lg font-semibold">{subscription.plan}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">Billing Cycle</label>
                    <p className="text-lg">{subscription.billingCycle}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">Amount</label>
                    <p className="text-lg font-semibold">${subscription.amount}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">Next Billing</label>
                    <p className="text-lg">{subscription.nextBilling}</p>
                  </div>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">Features</label>
                  <div className="flex flex-wrap gap-2 mt-2">
                    {subscription.features.map((feature, index) => (
                      <Badge key={index} variant="outline">{feature}</Badge>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          <div>
            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button className="w-full justify-start">
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Renew Subscription
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  <Settings className="h-4 w-4 mr-2" />
                  Modify Plan
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  <Download className="h-4 w-4 mr-2" />
                  Download Invoice
                </Button>
                <Button variant="outline" className="w-full justify-start text-red-600 hover:text-red-700">
                  <AlertTriangle className="h-4 w-4 mr-2" />
                  Suspend Account
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Usage Statistics */}
        <Card>
          <CardHeader>
            <CardTitle>Usage Statistics</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="text-center">
                <Users className="h-8 w-8 text-blue-600 mx-auto mb-2" />
                <p className="text-sm text-muted-foreground">Active Users</p>
                <p className="text-2xl font-bold">{subscription.usage.users.current}</p>
                <p className="text-xs text-muted-foreground">of {subscription.usage.users.limit}</p>
              </div>
              <div className="text-center">
                <DollarSign className="h-8 w-8 text-green-600 mx-auto mb-2" />
                <p className="text-sm text-muted-foreground">Storage Used</p>
                <p className="text-2xl font-bold">{subscription.usage.storage.current}</p>
                <p className="text-xs text-muted-foreground">of {subscription.usage.storage.limit}</p>
              </div>
              <div className="text-center">
                <Settings className="h-8 w-8 text-purple-600 mx-auto mb-2" />
                <p className="text-sm text-muted-foreground">API Calls</p>
                <p className="text-2xl font-bold">{subscription.usage.apiCalls.current.toLocaleString()}</p>
                <p className="text-xs text-muted-foreground">of {subscription.usage.apiCalls.limit.toLocaleString()}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Billing History */}
        <Card>
          <CardHeader>
            <CardTitle>Billing History</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <table className="min-w-full text-sm">
                <thead>
                  <tr className="border-b">
                    <th className="text-left p-4 font-medium">Date</th>
                    <th className="text-left p-4 font-medium">Amount</th>
                    <th className="text-left p-4 font-medium">Status</th>
                    <th className="text-left p-4 font-medium">Invoice</th>
                    <th className="text-left p-4 font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {subscription.billingHistory.map((bill, index) => (
                    <tr key={index} className="border-b hover:bg-muted/50">
                      <td className="p-4">
                        <div className="flex items-center gap-2">
                          <Calendar className="h-4 w-4 text-muted-foreground" />
                          {bill.date}
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="flex items-center gap-1">
                          <DollarSign className="h-4 w-4 text-muted-foreground" />
                          {bill.amount}
                        </div>
                      </td>
                      <td className="p-4">
                        <Badge variant="default" className="bg-green-600">{bill.status}</Badge>
                      </td>
                      <td className="p-4">
                        <span className="font-mono text-sm">{bill.invoice}</span>
                      </td>
                      <td className="p-4">
                        <Button size="sm" variant="outline">
                          <Download className="h-4 w-4 mr-2" />
                          Download
                        </Button>
                      </td>
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