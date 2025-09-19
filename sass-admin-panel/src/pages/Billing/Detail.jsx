import React from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  ArrowLeft,
  Download,
  Calendar,
  Building,
  CheckCircle,
  XCircle,
  Clock,
  AlertTriangle,
  Receipt,
  FileText,
  Eye,
  Send
} from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';

// Mock subscription billing data
const mockBillings = {
  1: {
    id: 'INV-2024-001',
    tenantName: 'Rahman Electronics',
    tenantEmail: 'rahman@email.com',
    tenantId: 'TEN-001',
    amount: 5000,
    status: 'paid',
    method: 'bKash',
    transactionId: 'BKS789123456',
    issueDate: '2024-07-24',
    dueDate: '2024-08-24',
    paidDate: '2024-07-24',
    billingPeriod: '2024-07-01 to 2024-07-31',
    plan: {
      name: 'Business Plan',
      features: ['Up to 1000 products', 'Advanced analytics', 'Priority support', '10GB storage'],
      users: 5,
      maxUsers: 10,
      storage: '7.2GB used of 10GB'
    },
    usageMetrics: {
      apiCalls: { used: 45000, limit: 50000 },
      storage: { used: 7.2, limit: 10, unit: 'GB' },
      users: { active: 5, limit: 10 },
      bandwidth: { used: 125, limit: 200, unit: 'GB' }
    },
    subscriptionDetails: {
      startDate: '2024-01-15',
      renewalDate: '2024-08-15',
      billingCycle: 'Monthly',
      autoRenew: true,
      proration: 0
    },
    charges: {
      baseSubscription: 4500,
      additionalUsers: 0,
      overage: 500,
      discounts: 0,
      tax: 0,
      total: 5000
    },
    paymentHistory: [
      { date: '2024-07-24', amount: 5000, method: 'bKash', status: 'Success' },
      { date: '2024-06-24', amount: 5000, method: 'bKash', status: 'Success' },
      { date: '2024-05-24', amount: 4500, method: 'Bank Transfer', status: 'Success' }
    ]
  },
  2: {
    id: 'INV-2024-002',
    tenantName: 'Modern Pharmacy',
    tenantEmail: 'ahmed@email.com',
    tenantId: 'TEN-002',
    amount: 10000,
    status: 'overdue',
    method: 'Bank Transfer',
    transactionId: 'BT456789123',
    issueDate: '2024-07-15',
    dueDate: '2024-08-15',
    paidDate: null,
    billingPeriod: '2024-07-01 to 2024-07-31',
    plan: {
      name: 'Enterprise Plan',
      features: ['Unlimited products', 'Advanced analytics', 'Priority support', '50GB storage', 'API access'],
      users: 15,
      maxUsers: 25,
      storage: '42GB used of 50GB'
    },
    usageMetrics: {
      apiCalls: { used: 180000, limit: 200000 },
      storage: { used: 42, limit: 50, unit: 'GB' },
      users: { active: 15, limit: 25 },
      bandwidth: { used: 380, limit: 500, unit: 'GB' }
    },
    subscriptionDetails: {
      startDate: '2024-01-01',
      renewalDate: '2024-08-15',
      billingCycle: 'Monthly',
      autoRenew: true,
      proration: 0
    },
    charges: {
      baseSubscription: 9000,
      additionalUsers: 1000,
      overage: 0,
      discounts: 0,
      tax: 0,
      total: 10000
    },
    paymentHistory: [
      { date: '2024-06-15', amount: 10000, method: 'Bank Transfer', status: 'Success' },
      { date: '2024-05-15', amount: 9500, method: 'Bank Transfer', status: 'Success' },
      { date: '2024-04-15', amount: 9000, method: 'Credit Card', status: 'Success' }
    ]
  }
};

export default function BillingDetail() {
  const { id } = useParams();
  const billing = mockBillings[id] || mockBillings[1];

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => {
    if (!dateString) return 'Not paid';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'paid': return 'bg-green-100 text-green-700 border-green-200';
      case 'pending': return 'bg-yellow-100 text-yellow-700 border-yellow-200';
      case 'overdue': return 'bg-red-100 text-red-700 border-red-200';
      case 'failed': return 'bg-red-100 text-red-700 border-red-200';
      default: return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'paid': return <CheckCircle className="h-4 w-4" />;
      case 'pending': return <Clock className="h-4 w-4" />;
      case 'overdue': return <AlertTriangle className="h-4 w-4" />;
      case 'failed': return <XCircle className="h-4 w-4" />;
      default: return <Clock className="h-4 w-4" />;
    }
  };

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Link to="/billing">
              <Button variant="ghost" size="sm">
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Billing
              </Button>
            </Link>
            <div className="flex items-center gap-3">
              <Receipt className="h-8 w-8 text-muted-foreground" />
              <div>
                <h1 className="text-xl font-semibold">Invoice {billing.id}</h1>
                <p className="text-sm text-muted-foreground">{billing.tenantName}</p>
              </div>
            </div>
            <Badge className={getStatusColor(billing.status)}>
              {getStatusIcon(billing.status)}
              <span className="ml-1 capitalize">{billing.status}</span>
            </Badge>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Download PDF
            </Button>
            <Button variant="outline" size="sm">
              <Send className="h-4 w-4 mr-2" />
              Send Invoice
            </Button>
            <Button size="sm">
              <Eye className="h-4 w-4 mr-2" />
              Print
            </Button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Invoice Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Invoice Header */}
            <Card>
              <CardContent className="p-6">
                <div className="flex justify-between items-start mb-6">
                  <div>
                    <h2 className="text-xl font-bold text-primary">Subscription Invoice</h2>
                    <p className="text-sm text-muted-foreground mt-1">#{billing.id}</p>
                    <p className="text-sm text-muted-foreground">Billing Period: {billing.billingPeriod}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-2xl font-bold">{formatCurrency(billing.amount)}</p>
                    <p className="text-sm text-muted-foreground">Total Amount</p>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-6">
                  <div>
                    <h4 className="font-semibold mb-2 flex items-center gap-2">
                      <Building className="h-4 w-4" />
                      Tenant Details
                    </h4>
                    <div className="space-y-1 text-sm">
                      <p className="font-medium">{billing.tenantName}</p>
                      <p className="text-muted-foreground">{billing.tenantEmail}</p>
                      <p className="text-muted-foreground">ID: {billing.tenantId}</p>
                    </div>
                  </div>

                  <div>
                    <h4 className="font-semibold mb-2 flex items-center gap-2">
                      <Calendar className="h-4 w-4" />
                      Invoice Dates
                    </h4>
                    <div className="space-y-1 text-sm">
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Issued:</span>
                        <span>{formatDate(billing.issueDate)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Due:</span>
                        <span>{formatDate(billing.dueDate)}</span>
                      </div>
                      {billing.paidDate && (
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Paid:</span>
                          <span className="text-green-600">{formatDate(billing.paidDate)}</span>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Subscription Plan Details */}
            <Card>
              <CardContent className="p-6">
                <h4 className="font-semibold mb-4">Subscription Plan: {billing.plan.name}</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <h5 className="font-medium mb-2">Plan Features</h5>
                    <ul className="space-y-1 text-sm text-muted-foreground">
                      {billing.plan.features.map((feature, index) => (
                        <li key={index} className="flex items-center gap-2">
                          <CheckCircle className="h-3 w-3 text-green-600" />
                          {feature}
                        </li>
                      ))}
                    </ul>
                  </div>
                  <div>
                    <h5 className="font-medium mb-2">Usage Summary</h5>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Active Users:</span>
                        <span>{billing.plan.users} of {billing.plan.maxUsers}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Storage:</span>
                        <span>{billing.plan.storage}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Billing Breakdown */}
            <Card>
              <CardContent className="p-6">
                <h4 className="font-semibold mb-4">Billing Breakdown</h4>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span>{billing.plan.name} (Monthly)</span>
                    <span>{formatCurrency(billing.charges.baseSubscription)}</span>
                  </div>
                  {billing.charges.additionalUsers > 0 && (
                    <div className="flex justify-between">
                      <span>Additional Users</span>
                      <span>{formatCurrency(billing.charges.additionalUsers)}</span>
                    </div>
                  )}
                  {billing.charges.overage > 0 && (
                    <div className="flex justify-between">
                      <span>Usage Overage</span>
                      <span>{formatCurrency(billing.charges.overage)}</span>
                    </div>
                  )}
                  {billing.charges.discounts > 0 && (
                    <div className="flex justify-between text-green-600">
                      <span>Discounts</span>
                      <span>-{formatCurrency(billing.charges.discounts)}</span>
                    </div>
                  )}
                  <Separator />
                  <div className="flex justify-between font-semibold text-lg">
                    <span>Total</span>
                    <span>{formatCurrency(billing.charges.total)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Payment Status */}
            {billing.status === 'paid' && (
              <Card className="border-green-200">
                <CardContent className="p-6">
                  <div className="flex items-center gap-3 mb-4">
                    <CheckCircle className="h-5 w-5 text-green-600" />
                    <h4 className="font-semibold text-green-700">Payment Received</h4>
                  </div>
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Payment Method:</span>
                      <span>{billing.method}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Transaction ID:</span>
                      <span className="font-mono">{billing.transactionId}</span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            )}

            {billing.status === 'overdue' && (
              <Card className="border-red-200">
                <CardContent className="p-6">
                  <div className="flex items-center gap-3 mb-4">
                    <AlertTriangle className="h-5 w-5 text-red-600" />
                    <h4 className="font-semibold text-red-700">Payment Overdue</h4>
                  </div>
                  <p className="text-sm text-red-600 mb-4">
                    This invoice is past due. Please follow up with the tenant for payment.
                  </p>
                  <div className="flex gap-2">
                    <Button size="sm" variant="outline" className="text-red-600 border-red-200">
                      Send Reminder
                    </Button>
                    <Button size="sm" variant="outline" className="text-red-600 border-red-200">
                      Mark as Paid
                    </Button>
                  </div>
                </CardContent>
              </Card>
            )}
          </div>

          {/* Right Sidebar */}
          <div className="space-y-6">
            {/* Subscription Info */}
            <Card>
              <CardContent className="p-6">
                <h4 className="font-semibold mb-4">Subscription Details</h4>
                <div className="space-y-3 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Start Date:</span>
                    <span>{formatDate(billing.subscriptionDetails.startDate)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Next Renewal:</span>
                    <span>{formatDate(billing.subscriptionDetails.renewalDate)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Billing Cycle:</span>
                    <span>{billing.subscriptionDetails.billingCycle}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Auto Renew:</span>
                    <span className={billing.subscriptionDetails.autoRenew ? 'text-green-600' : 'text-red-600'}>
                      {billing.subscriptionDetails.autoRenew ? 'Enabled' : 'Disabled'}
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Usage Metrics */}
            <Card>
              <CardContent className="p-6">
                <h4 className="font-semibold mb-4">Current Usage</h4>
                <div className="space-y-4">
                  <div>
                    <div className="flex justify-between text-sm mb-1">
                      <span>API Calls</span>
                      <span>{billing.usageMetrics.apiCalls.used.toLocaleString()} / {billing.usageMetrics.apiCalls.limit.toLocaleString()}</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div className="bg-blue-600 h-2 rounded-full" style={{width: `${(billing.usageMetrics.apiCalls.used / billing.usageMetrics.apiCalls.limit) * 100}%`}}></div>
                    </div>
                  </div>

                  <div>
                    <div className="flex justify-between text-sm mb-1">
                      <span>Storage</span>
                      <span>{billing.usageMetrics.storage.used}{billing.usageMetrics.storage.unit} / {billing.usageMetrics.storage.limit}{billing.usageMetrics.storage.unit}</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div className="bg-green-600 h-2 rounded-full" style={{width: `${(billing.usageMetrics.storage.used / billing.usageMetrics.storage.limit) * 100}%`}}></div>
                    </div>
                  </div>

                  <div>
                    <div className="flex justify-between text-sm mb-1">
                      <span>Active Users</span>
                      <span>{billing.usageMetrics.users.active} / {billing.usageMetrics.users.limit}</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div className="bg-purple-600 h-2 rounded-full" style={{width: `${(billing.usageMetrics.users.active / billing.usageMetrics.users.limit) * 100}%`}}></div>
                    </div>
                  </div>

                  <div>
                    <div className="flex justify-between text-sm mb-1">
                      <span>Bandwidth</span>
                      <span>{billing.usageMetrics.bandwidth.used}{billing.usageMetrics.bandwidth.unit} / {billing.usageMetrics.bandwidth.limit}{billing.usageMetrics.bandwidth.unit}</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div className="bg-orange-600 h-2 rounded-full" style={{width: `${(billing.usageMetrics.bandwidth.used / billing.usageMetrics.bandwidth.limit) * 100}%`}}></div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Payment History */}
            <Card>
              <CardContent className="p-6">
                <h4 className="font-semibold mb-4">Recent Payments</h4>
                <div className="space-y-3">
                  {billing.paymentHistory.map((payment, index) => (
                    <div key={index} className="flex justify-between items-center text-sm">
                      <div>
                        <p className="font-medium">{formatCurrency(payment.amount)}</p>
                        <p className="text-muted-foreground">{formatDate(payment.date)}</p>
                      </div>
                      <div className="text-right">
                        <Badge className={payment.status === 'Success' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}>
                          {payment.status}
                        </Badge>
                        <p className="text-muted-foreground text-xs mt-1">{payment.method}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}