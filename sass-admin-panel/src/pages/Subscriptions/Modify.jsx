import React, { useState, useEffect } from 'react';
import {
  ArrowLeft,
  Package,
  Calendar,
  DollarSign,
  Users,
  Settings,
  Save,
  X,
  AlertTriangle,
  CheckCircle,
  Clock,
  CreditCard,
  Pause,
  Play,
  Trash2
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
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Link, useParams } from 'react-router-dom';

const subscriptionPlans = [
  {
    id: 'starter',
    name: 'Starter Plan',
    description: 'Perfect for small businesses getting started',
    basePrice: 29,
    features: ['Up to 100 products', 'Basic analytics', 'Email support', '5GB storage']
  },
  {
    id: 'professional',
    name: 'Professional Plan',
    description: 'Advanced features for growing businesses',
    basePrice: 79,
    features: ['Up to 1000 products', 'Advanced analytics', 'Priority support', '50GB storage', 'API access']
  },
  {
    id: 'enterprise',
    name: 'Enterprise Plan',
    description: 'Full-featured solution for large organizations',
    basePrice: 199,
    features: ['Unlimited products', 'Custom analytics', '24/7 support', 'Unlimited storage', 'Full API access', 'Custom integrations']
  }
];

const addOns = [
  { id: 'extra_storage', name: 'Extra Storage (100GB)', price: 10 },
  { id: 'premium_support', name: 'Premium Support', price: 25 },
  { id: 'custom_domain', name: 'Custom Domain', price: 15 },
  { id: 'advanced_security', name: 'Advanced Security', price: 30 }
];

// Mock subscription data
const mockSubscription = {
  id: 'sub_123',
  tenantId: 'tenant_1',
  tenantName: 'TechCorp Solutions',
  tenantEmail: 'admin@techcorp.com',
  planId: 'professional',
  status: 'active',
  billingCycle: 'monthly',
  currentPrice: 79,
  nextBillingDate: '2024-02-15',
  startDate: '2024-01-15',
  autoRenew: true,
  trialDays: 0,
  selectedAddOns: ['extra_storage', 'premium_support'],
  notes: 'Customer requested priority support due to high volume usage.',
  discountPercent: 10
};

export default function ModifySubscription() {
  const { id } = useParams();
  const [subscription, setSubscription] = useState(mockSubscription);
  const [formData, setFormData] = useState({
    planId: '',
    billingCycle: '',
    customPrice: '',
    discountPercent: '',
    notes: '',
    autoRenew: true,
    selectedAddOns: []
  });
  const [selectedPlan, setSelectedPlan] = useState(null);
  const [customPricing, setCustomPricing] = useState(false);
  const [showCancelDialog, setShowCancelDialog] = useState(false);
  const [actionType, setActionType] = useState('modify'); // modify, pause, cancel

  useEffect(() => {
    // Initialize form with current subscription data
    setFormData({
      planId: subscription.planId,
      billingCycle: subscription.billingCycle,
      customPrice: subscription.currentPrice.toString(),
      discountPercent: subscription.discountPercent?.toString() || '',
      notes: subscription.notes,
      autoRenew: subscription.autoRenew,
      selectedAddOns: subscription.selectedAddOns
    });
    
    const plan = subscriptionPlans.find(p => p.id === subscription.planId);
    setSelectedPlan(plan);
    setCustomPricing(subscription.currentPrice !== plan?.basePrice);
  }, [subscription]);

  const handleInputChange = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handlePlanSelect = (planId) => {
    const plan = subscriptionPlans.find(p => p.id === planId);
    setSelectedPlan(plan);
    handleInputChange('planId', planId);
  };

  const handleAddOnToggle = (addOnId) => {
    setFormData(prev => ({
      ...prev,
      selectedAddOns: prev.selectedAddOns.includes(addOnId)
        ? prev.selectedAddOns.filter(id => id !== addOnId)
        : [...prev.selectedAddOns, addOnId]
    }));
  };

  const calculateNewPrice = () => {
    if (!selectedPlan) return 0;
    
    let basePrice = customPricing && formData.customPrice 
      ? parseFloat(formData.customPrice) 
      : selectedPlan.basePrice;
    
    if (formData.billingCycle === 'yearly') {
      basePrice = basePrice * 12 * 0.9; // 10% discount for yearly
    }
    
    const addOnsPrice = formData.selectedAddOns.reduce((total, addOnId) => {
      const addOn = addOns.find(a => a.id === addOnId);
      return total + (addOn ? addOn.price : 0);
    }, 0);
    
    let total = basePrice + addOnsPrice;
    
    if (formData.discountPercent) {
      total = total * (1 - parseFloat(formData.discountPercent) / 100);
    }
    
    return total;
  };

  const getStatusBadge = (status) => {
    const statusConfig = {
      active: { color: 'bg-green-100 text-green-800', icon: CheckCircle },
      paused: { color: 'bg-yellow-100 text-yellow-800', icon: Pause },
      cancelled: { color: 'bg-red-100 text-red-800', icon: X },
      trial: { color: 'bg-blue-100 text-blue-800', icon: Clock }
    };
    
    const config = statusConfig[status] || statusConfig.active;
    const Icon = config.icon;
    
    return (
      <Badge className={config.color}>
        <Icon className="h-3 w-3 mr-1" />
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </Badge>
    );
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log('Modifying subscription:', { id, formData, actionType });
    // Handle form submission
  };

  const handlePauseSubscription = () => {
    setActionType('pause');
    console.log('Pausing subscription:', id);
  };

  const handleCancelSubscription = () => {
    setActionType('cancel');
    setShowCancelDialog(true);
  };

  const confirmCancellation = () => {
    console.log('Cancelling subscription:', id);
    setShowCancelDialog(false);
  };

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Link to="/subscriptions" className="flex items-center gap-2 text-muted-foreground hover:text-foreground">
              <ArrowLeft className="h-4 w-4" />
              Back to Subscriptions
            </Link>
            <Separator orientation="vertical" className="h-6" />
            <Package className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Modify Subscription</h1>
            {getStatusBadge(subscription.status)}
          </div>
          <div className="flex items-center gap-2">
            {subscription.status === 'active' && (
              <>
                <Button variant="outline" onClick={handlePauseSubscription}>
                  <Pause className="h-4 w-4 mr-2" />
                  Pause
                </Button>
                <Button variant="outline" onClick={handleCancelSubscription}>
                  <Trash2 className="h-4 w-4 mr-2" />
                  Cancel
                </Button>
              </>
            )}
            {subscription.status === 'paused' && (
              <Button variant="outline">
                <Play className="h-4 w-4 mr-2" />
                Resume
              </Button>
            )}
            <Button variant="outline">
              <X className="h-4 w-4 mr-2" />
              Discard Changes
            </Button>
            <Button type="submit" form="subscription-form">
              <Save className="h-4 w-4 mr-2" />
              Save Changes
            </Button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6 space-y-6">
          {/* Current Subscription Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="h-5 w-5" />
                Current Subscription
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Tenant</Label>
                  <div className="mt-1">
                    <div className="font-medium">{subscription.tenantName}</div>
                    <div className="text-sm text-muted-foreground">{subscription.tenantEmail}</div>
                  </div>
                </div>
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Current Plan</Label>
                  <div className="mt-1">
                    <div className="font-medium">{selectedPlan?.name}</div>
                    <div className="text-sm text-muted-foreground">${subscription.currentPrice}/{subscription.billingCycle}</div>
                  </div>
                </div>
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Next Billing</Label>
                  <div className="mt-1">
                    <div className="font-medium">{subscription.nextBillingDate}</div>
                    <div className="text-sm text-muted-foreground">
                      {subscription.autoRenew ? 'Auto-renewal enabled' : 'Auto-renewal disabled'}
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <form id="subscription-form" onSubmit={handleSubmit} className="space-y-6">
            {/* Plan Modification */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Package className="h-5 w-5" />
                  Change Plan
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {subscriptionPlans.map((plan) => (
                    <div
                      key={plan.id}
                      className={`border rounded-lg p-4 cursor-pointer transition-colors ${
                        formData.planId === plan.id
                          ? 'border-primary bg-primary/5'
                          : 'border-border hover:border-primary/50'
                      }`}
                      onClick={() => handlePlanSelect(plan.id)}
                    >
                      <div className="flex items-center justify-between mb-2">
                        <h3 className="font-semibold">{plan.name}</h3>
                        <div className="text-right">
                          <div className="text-2xl font-bold">${plan.basePrice}</div>
                          <div className="text-sm text-muted-foreground">/month</div>
                        </div>
                      </div>
                      {plan.id === subscription.planId && (
                        <Badge className="mb-2">Current Plan</Badge>
                      )}
                      <p className="text-sm text-muted-foreground mb-3">{plan.description}</p>
                      <ul className="space-y-1">
                        {plan.features.slice(0, 3).map((feature, index) => (
                          <li key={index} className="text-sm flex items-center gap-2">
                            <div className="w-1 h-1 bg-primary rounded-full"></div>
                            {feature}
                          </li>
                        ))}
                      </ul>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Billing Modifications */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <DollarSign className="h-5 w-5" />
                  Billing Settings
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="billingCycle">Billing Cycle</Label>
                    <Select value={formData.billingCycle} onValueChange={(value) => handleInputChange('billingCycle', value)}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="monthly">Monthly</SelectItem>
                        <SelectItem value="yearly">Yearly (10% discount)</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex items-center space-x-2 pt-6">
                    <Switch
                      checked={formData.autoRenew}
                      onCheckedChange={(checked) => handleInputChange('autoRenew', checked)}
                    />
                    <Label>Auto-renewal</Label>
                  </div>
                </div>

                <div className="flex items-center space-x-2">
                  <Switch
                    checked={customPricing}
                    onCheckedChange={setCustomPricing}
                  />
                  <Label>Custom Pricing</Label>
                </div>

                {customPricing && (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <Label htmlFor="customPrice">Custom Price ($)</Label>
                      <Input
                        type="number"
                        step="0.01"
                        value={formData.customPrice}
                        onChange={(e) => handleInputChange('customPrice', e.target.value)}
                        placeholder="Enter custom price"
                      />
                    </div>
                    <div>
                      <Label htmlFor="discountPercent">Discount (%)</Label>
                      <Input
                        type="number"
                        min="0"
                        max="100"
                        value={formData.discountPercent}
                        onChange={(e) => handleInputChange('discountPercent', e.target.value)}
                        placeholder="Enter discount percentage"
                      />
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Add-ons Modification */}
            <Card>
              <CardHeader>
                <CardTitle>Manage Add-ons</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {addOns.map((addOn) => (
                    <div key={addOn.id} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center space-x-3">
                        <Switch
                          checked={formData.selectedAddOns.includes(addOn.id)}
                          onCheckedChange={() => handleAddOnToggle(addOn.id)}
                        />
                        <div>
                          <div className="font-medium">{addOn.name}</div>
                          <div className="text-sm text-muted-foreground">${addOn.price}/month</div>
                        </div>
                      </div>
                      {subscription.selectedAddOns.includes(addOn.id) && (
                        <Badge variant="outline">Current</Badge>
                      )}
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Notes */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Settings className="h-5 w-5" />
                  Additional Notes
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div>
                  <Label htmlFor="notes">Notes</Label>
                  <Textarea
                    value={formData.notes}
                    onChange={(e) => handleInputChange('notes', e.target.value)}
                    placeholder="Add any additional notes about this modification..."
                    rows={3}
                  />
                </div>
              </CardContent>
            </Card>

            {/* Price Comparison */}
            {selectedPlan && (
              <Card>
                <CardHeader>
                  <CardTitle>Price Comparison</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                      <h4 className="font-medium mb-3">Current</h4>
                      <div className="space-y-2">
                        <div className="flex justify-between">
                          <span>Plan:</span>
                          <span>{subscriptionPlans.find(p => p.id === subscription.planId)?.name}</span>
                        </div>
                        <div className="flex justify-between">
                          <span>Price:</span>
                          <span>${subscription.currentPrice}/{subscription.billingCycle}</span>
                        </div>
                      </div>
                    </div>
                    <div>
                      <h4 className="font-medium mb-3">New</h4>
                      <div className="space-y-2">
                        <div className="flex justify-between">
                          <span>Plan:</span>
                          <span>{selectedPlan.name}</span>
                        </div>
                        <div className="flex justify-between">
                          <span>Price:</span>
                          <span className="font-semibold">${calculateNewPrice().toFixed(2)}/{formData.billingCycle}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                  
                  {calculateNewPrice() !== subscription.currentPrice && (
                    <Alert className="mt-4">
                      <AlertTriangle className="h-4 w-4" />
                      <AlertDescription>
                        {calculateNewPrice() > subscription.currentPrice 
                          ? `Price will increase by $${(calculateNewPrice() - subscription.currentPrice).toFixed(2)}/${formData.billingCycle}. Changes will take effect on the next billing cycle.`
                          : `Price will decrease by $${(subscription.currentPrice - calculateNewPrice()).toFixed(2)}/${formData.billingCycle}. Customer will receive a prorated credit.`
                        }
                      </AlertDescription>
                    </Alert>
                  )}
                </CardContent>
              </Card>
            )}
          </form>
      </div>

      {/* Cancel Confirmation Dialog */}
      {showCancelDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <Card className="w-full max-w-md">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-red-600">
                <AlertTriangle className="h-5 w-5" />
                Cancel Subscription
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <p>Are you sure you want to cancel this subscription? This action cannot be undone.</p>
              <div className="flex justify-end gap-2">
                <Button variant="outline" onClick={() => setShowCancelDialog(false)}>
                  Keep Subscription
                </Button>
                <Button variant="destructive" onClick={confirmCancellation}>
                  Cancel Subscription
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}