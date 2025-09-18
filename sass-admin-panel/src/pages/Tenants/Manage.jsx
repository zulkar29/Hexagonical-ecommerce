import React, { useState } from 'react';
import {
  ArrowLeft,
  Building2,
  Users,
  CreditCard,
  Settings,
  Activity,
  MoreHorizontal,
  Edit,
  Trash2,
  Pause,
  Play,
  Mail,
  Phone,
  Globe,
  MapPin,
  Calendar,
  DollarSign,
  TrendingUp,
  AlertTriangle,
  CheckCircle,
  Clock,
  Package,
  Database,
  Shield,
  Eye,
  Download,
  RefreshCw
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Link, useParams } from 'react-router-dom';

// Mock tenant data
const mockTenant = {
  id: 'tenant-001',
  companyName: 'Acme Corporation',
  companyEmail: 'admin@acme.com',
  status: 'active',
  plan: 'Professional',
  planPrice: 79,
  billingCycle: 'monthly',
  trialEndsAt: null,
  createdAt: '2024-01-15',
  lastLoginAt: '2024-01-20T10:30:00Z',
  industry: 'Technology',
  companySize: '51-200',
  website: 'https://acme.com',
  phone: '+1 (555) 123-4567',
  address: {
    street: '123 Business Street',
    city: 'San Francisco',
    state: 'CA',
    zipCode: '94105',
    country: 'United States'
  },
  admin: {
    firstName: 'John',
    lastName: 'Doe',
    email: 'john.doe@acme.com',
    phone: '+1 (555) 987-6543',
    lastLogin: '2024-01-20T10:30:00Z'
  },
  settings: {
    timezone: 'America/Los_Angeles',
    currency: 'USD',
    language: 'en',
    customDomain: 'shop.acme.com',
    enableNotifications: true,
    enableAnalytics: true
  },
  usage: {
    products: 245,
    orders: 1250,
    customers: 890,
    storage: 15.2, // GB
    apiCalls: 45000
  },
  limits: {
    products: 1000,
    storage: 50, // GB
    apiCalls: 100000
  },
  billing: {
    nextBillingDate: '2024-02-15',
    lastPayment: {
      amount: 79,
      date: '2024-01-15',
      status: 'paid'
    },
    paymentMethod: {
      type: 'card',
      last4: '4242',
      brand: 'visa'
    }
  }
};

const mockActivity = [
  {
    id: 1,
    type: 'login',
    description: 'Admin user logged in',
    timestamp: '2024-01-20T10:30:00Z',
    user: 'john.doe@acme.com'
  },
  {
    id: 2,
    type: 'subscription',
    description: 'Subscription renewed',
    timestamp: '2024-01-15T09:00:00Z',
    user: 'system'
  },
  {
    id: 3,
    type: 'settings',
    description: 'Custom domain configured',
    timestamp: '2024-01-14T14:22:00Z',
    user: 'john.doe@acme.com'
  },
  {
    id: 4,
    type: 'user',
    description: 'New user added to account',
    timestamp: '2024-01-12T11:15:00Z',
    user: 'john.doe@acme.com'
  },
  {
    id: 5,
    type: 'billing',
    description: 'Payment method updated',
    timestamp: '2024-01-10T16:45:00Z',
    user: 'john.doe@acme.com'
  }
];

const getStatusColor = (status) => {
  switch (status) {
    case 'active': return 'bg-green-100 text-green-800';
    case 'trial': return 'bg-blue-100 text-blue-800';
    case 'suspended': return 'bg-yellow-100 text-yellow-800';
    case 'cancelled': return 'bg-red-100 text-red-800';
    default: return 'bg-gray-100 text-gray-800';
  }
};

const getActivityIcon = (type) => {
  switch (type) {
    case 'login': return Users;
    case 'subscription': return Package;
    case 'settings': return Settings;
    case 'user': return Users;
    case 'billing': return CreditCard;
    default: return Activity;
  }
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

const formatDateTime = (dateString) => {
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

export default function TenantManagement() {
  const { id } = useParams();
  const [tenant, setTenant] = useState(mockTenant);
  const [isEditing, setIsEditing] = useState(false);
  const [showSuspendDialog, setShowSuspendDialog] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [editForm, setEditForm] = useState({
    companyName: tenant.companyName,
    companyEmail: tenant.companyEmail,
    website: tenant.website,
    phone: tenant.phone,
    customDomain: tenant.settings.customDomain
  });

  const handleSave = () => {
    setTenant(prev => ({
      ...prev,
      companyName: editForm.companyName,
      companyEmail: editForm.companyEmail,
      website: editForm.website,
      phone: editForm.phone,
      settings: {
        ...prev.settings,
        customDomain: editForm.customDomain
      }
    }));
    setIsEditing(false);
  };

  const handleSuspend = () => {
    setTenant(prev => ({ ...prev, status: 'suspended' }));
    setShowSuspendDialog(false);
  };

  const handleReactivate = () => {
    setTenant(prev => ({ ...prev, status: 'active' }));
  };

  const usagePercentage = (used, limit) => {
    return Math.round((used / limit) * 100);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Link to="/tenants" className="flex items-center gap-2 text-muted-foreground hover:text-foreground">
              <ArrowLeft className="h-4 w-4" />
              Back to Tenants
            </Link>
            <Separator orientation="vertical" className="h-6" />
            <Building2 className="h-5 w-5 text-primary" />
            <div>
              <h1 className="text-xl font-semibold">{tenant.companyName}</h1>
              <p className="text-sm text-muted-foreground">Tenant ID: {tenant.id}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Badge className={getStatusColor(tenant.status)}>
              {tenant.status.charAt(0).toUpperCase() + tenant.status.slice(1)}
            </Badge>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                <DropdownMenuItem onClick={() => setIsEditing(true)}>
                  <Edit className="h-4 w-4 mr-2" />
                  Edit Details
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Eye className="h-4 w-4 mr-2" />
                  View as Tenant
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Download className="h-4 w-4 mr-2" />
                  Export Data
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                {tenant.status === 'active' ? (
                  <DropdownMenuItem onClick={() => setShowSuspendDialog(true)} className="text-yellow-600">
                    <Pause className="h-4 w-4 mr-2" />
                    Suspend Account
                  </DropdownMenuItem>
                ) : (
                  <DropdownMenuItem onClick={handleReactivate} className="text-green-600">
                    <Play className="h-4 w-4 mr-2" />
                    Reactivate Account
                  </DropdownMenuItem>
                )}
                <DropdownMenuItem onClick={() => setShowDeleteDialog(true)} className="text-red-600">
                  <Trash2 className="h-4 w-4 mr-2" />
                  Delete Tenant
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      <div className="w-full p-6 space-y-6">
          {/* Overview Cards */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Monthly Revenue</p>
                    <p className="text-2xl font-bold">${tenant.planPrice}</p>
                  </div>
                  <DollarSign className="h-8 w-8 text-green-500" />
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Total Users</p>
                    <p className="text-2xl font-bold">{tenant.usage.customers}</p>
                  </div>
                  <Users className="h-8 w-8 text-blue-500" />
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Total Orders</p>
                    <p className="text-2xl font-bold">{tenant.usage.orders}</p>
                  </div>
                  <Package className="h-8 w-8 text-purple-500" />
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Storage Used</p>
                    <p className="text-2xl font-bold">{tenant.usage.storage}GB</p>
                  </div>
                  <Database className="h-8 w-8 text-orange-500" />
                </div>
              </CardContent>
            </Card>
          </div>

          <Tabs defaultValue="overview" className="space-y-4">
            <TabsList>
              <TabsTrigger value="overview">Overview</TabsTrigger>
              <TabsTrigger value="subscription">Subscription</TabsTrigger>
              <TabsTrigger value="usage">Usage & Limits</TabsTrigger>
              <TabsTrigger value="settings">Settings</TabsTrigger>
              <TabsTrigger value="activity">Activity</TabsTrigger>
            </TabsList>

            <TabsContent value="overview" className="space-y-4">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Company Information */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center justify-between">
                      <span className="flex items-center gap-2">
                        <Building2 className="h-5 w-5" />
                        Company Information
                      </span>
                      {!isEditing && (
                        <Button variant="outline" size="sm" onClick={() => setIsEditing(true)}>
                          <Edit className="h-4 w-4 mr-2" />
                          Edit
                        </Button>
                      )}
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {isEditing ? (
                      <div className="space-y-4">
                        <div>
                          <Label htmlFor="companyName">Company Name</Label>
                          <Input
                            id="companyName"
                            value={editForm.companyName}
                            onChange={(e) => setEditForm(prev => ({ ...prev, companyName: e.target.value }))}
                          />
                        </div>
                        <div>
                          <Label htmlFor="companyEmail">Company Email</Label>
                          <Input
                            id="companyEmail"
                            value={editForm.companyEmail}
                            onChange={(e) => setEditForm(prev => ({ ...prev, companyEmail: e.target.value }))}
                          />
                        </div>
                        <div>
                          <Label htmlFor="website">Website</Label>
                          <Input
                            id="website"
                            value={editForm.website}
                            onChange={(e) => setEditForm(prev => ({ ...prev, website: e.target.value }))}
                          />
                        </div>
                        <div>
                          <Label htmlFor="phone">Phone</Label>
                          <Input
                            id="phone"
                            value={editForm.phone}
                            onChange={(e) => setEditForm(prev => ({ ...prev, phone: e.target.value }))}
                          />
                        </div>
                        <div>
                          <Label htmlFor="customDomain">Custom Domain</Label>
                          <Input
                            id="customDomain"
                            value={editForm.customDomain}
                            onChange={(e) => setEditForm(prev => ({ ...prev, customDomain: e.target.value }))}
                          />
                        </div>
                        <div className="flex gap-2">
                          <Button onClick={handleSave}>
                            <CheckCircle className="h-4 w-4 mr-2" />
                            Save Changes
                          </Button>
                          <Button variant="outline" onClick={() => setIsEditing(false)}>
                            Cancel
                          </Button>
                        </div>
                      </div>
                    ) : (
                      <div className="space-y-3">
                        <div className="flex items-center gap-2">
                          <Mail className="h-4 w-4 text-muted-foreground" />
                          <span>{tenant.companyEmail}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <Phone className="h-4 w-4 text-muted-foreground" />
                          <span>{tenant.phone}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <Globe className="h-4 w-4 text-muted-foreground" />
                          <a href={tenant.website} target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">
                            {tenant.website}
                          </a>
                        </div>
                        <div className="flex items-center gap-2">
                          <MapPin className="h-4 w-4 text-muted-foreground" />
                          <span>{tenant.address.city}, {tenant.address.state}, {tenant.address.country}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <Building2 className="h-4 w-4 text-muted-foreground" />
                          <span>{tenant.industry} • {tenant.companySize} employees</span>
                        </div>
                        {tenant.settings.customDomain && (
                          <div className="flex items-center gap-2">
                            <Globe className="h-4 w-4 text-muted-foreground" />
                            <span>Custom Domain: {tenant.settings.customDomain}</span>
                          </div>
                        )}
                      </div>
                    )}
                  </CardContent>
                </Card>

                {/* Admin User */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Users className="h-5 w-5" />
                      Admin User
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    <div>
                      <p className="font-medium">{tenant.admin.firstName} {tenant.admin.lastName}</p>
                      <p className="text-sm text-muted-foreground">Primary Administrator</p>
                    </div>
                    <div className="flex items-center gap-2">
                      <Mail className="h-4 w-4 text-muted-foreground" />
                      <span>{tenant.admin.email}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <Phone className="h-4 w-4 text-muted-foreground" />
                      <span>{tenant.admin.phone}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <Clock className="h-4 w-4 text-muted-foreground" />
                      <span>Last login: {formatDateTime(tenant.admin.lastLogin)}</span>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            <TabsContent value="subscription" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Package className="h-5 w-5" />
                    Current Subscription
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                      <p className="text-sm font-medium text-muted-foreground">Plan</p>
                      <p className="text-lg font-semibold">{tenant.plan}</p>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-muted-foreground">Price</p>
                      <p className="text-lg font-semibold">${tenant.planPrice}/{tenant.billingCycle}</p>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-muted-foreground">Next Billing</p>
                      <p className="text-lg font-semibold">{formatDate(tenant.billing.nextBillingDate)}</p>
                    </div>
                  </div>
                  
                  <Separator />
                  
                  <div>
                    <h4 className="font-medium mb-3">Billing Information</h4>
                    <div className="space-y-2">
                      <div className="flex justify-between">
                        <span>Payment Method:</span>
                        <span className="capitalize">{tenant.billing.paymentMethod.brand} •••• {tenant.billing.paymentMethod.last4}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Last Payment:</span>
                        <span>${tenant.billing.lastPayment.amount} on {formatDate(tenant.billing.lastPayment.date)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Status:</span>
                        <Badge className="bg-green-100 text-green-800">{tenant.billing.lastPayment.status}</Badge>
                      </div>
                    </div>
                  </div>
                  
                  <div className="flex gap-2">
                    <Link to={`/subscriptions/${tenant.id}/modify`}>
                      <Button variant="outline">
                        <Edit className="h-4 w-4 mr-2" />
                        Modify Subscription
                      </Button>
                    </Link>
                    <Button variant="outline">
                      <CreditCard className="h-4 w-4 mr-2" />
                      Update Billing
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="usage" className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Card>
                  <CardHeader>
                    <CardTitle>Products</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>{tenant.usage.products} / {tenant.limits.products}</span>
                        <span>{usagePercentage(tenant.usage.products, tenant.limits.products)}%</span>
                      </div>
                      <Progress value={usagePercentage(tenant.usage.products, tenant.limits.products)} />
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardHeader>
                    <CardTitle>Storage</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>{tenant.usage.storage}GB / {tenant.limits.storage}GB</span>
                        <span>{usagePercentage(tenant.usage.storage, tenant.limits.storage)}%</span>
                      </div>
                      <Progress value={usagePercentage(tenant.usage.storage, tenant.limits.storage)} />
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardHeader>
                    <CardTitle>API Calls</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>{tenant.usage.apiCalls.toLocaleString()} / {tenant.limits.apiCalls.toLocaleString()}</span>
                        <span>{usagePercentage(tenant.usage.apiCalls, tenant.limits.apiCalls)}%</span>
                      </div>
                      <Progress value={usagePercentage(tenant.usage.apiCalls, tenant.limits.apiCalls)} />
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardHeader>
                    <CardTitle>Customers</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>{tenant.usage.customers} active</span>
                        <TrendingUp className="h-4 w-4 text-green-500" />
                      </div>
                      <p className="text-sm text-muted-foreground">No limit on current plan</p>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            <TabsContent value="settings" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Settings className="h-5 w-5" />
                    Tenant Settings
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                      <Label>Timezone</Label>
                      <p className="text-sm">{tenant.settings.timezone}</p>
                    </div>
                    <div>
                      <Label>Currency</Label>
                      <p className="text-sm">{tenant.settings.currency}</p>
                    </div>
                    <div>
                      <Label>Language</Label>
                      <p className="text-sm">{tenant.settings.language}</p>
                    </div>
                  </div>
                  
                  <Separator />
                  
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <Label>Email Notifications</Label>
                        <p className="text-sm text-muted-foreground">Send system notifications via email</p>
                      </div>
                      <Switch checked={tenant.settings.enableNotifications} />
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <Label>Analytics Tracking</Label>
                        <p className="text-sm text-muted-foreground">Enable usage analytics and reporting</p>
                      </div>
                      <Switch checked={tenant.settings.enableAnalytics} />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="activity" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <Activity className="h-5 w-5" />
                      Recent Activity
                    </span>
                    <Button variant="outline" size="sm">
                      <RefreshCw className="h-4 w-4 mr-2" />
                      Refresh
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {mockActivity.map((activity) => {
                      const Icon = getActivityIcon(activity.type);
                      return (
                        <div key={activity.id} className="flex items-start gap-3 p-3 border rounded-lg">
                          <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                            <Icon className="h-4 w-4" />
                          </div>
                          <div className="flex-1">
                            <p className="text-sm font-medium">{activity.description}</p>
                            <div className="flex items-center gap-2 text-xs text-muted-foreground mt-1">
                              <span>{formatDateTime(activity.timestamp)}</span>
                              <span>•</span>
                              <span>{activity.user}</span>
                            </div>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
      </div>

      {/* Suspend Dialog */}
      <Dialog open={showSuspendDialog} onOpenChange={setShowSuspendDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Suspend Tenant Account</DialogTitle>
            <DialogDescription>
              Are you sure you want to suspend this tenant account? The tenant will lose access to their account until reactivated.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowSuspendDialog(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleSuspend}>
              Suspend Account
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Tenant Account</DialogTitle>
            <DialogDescription>
              This action cannot be undone. This will permanently delete the tenant account and all associated data.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowDeleteDialog(false)}>
              Cancel
            </Button>
            <Button variant="destructive">
              Delete Account
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}