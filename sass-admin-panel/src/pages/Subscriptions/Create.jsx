import React, { useState } from 'react';
import {
  ArrowLeft,
  Plus,
  Minus,
  Calendar,
  DollarSign,
  Users,
  Package,
  Settings,
  Save,
  X
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
import { Link } from 'react-router-dom';

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

const tenantsList = [
  { id: 'tenant_1', name: 'TechCorp Solutions', email: 'admin@techcorp.com' },
  { id: 'tenant_2', name: 'StartupHub', email: 'contact@startuphub.com' },
  { id: 'tenant_3', name: 'RetailMax', email: 'info@retailmax.com' },
  { id: 'tenant_4', name: 'GlobalTrade Inc', email: 'admin@globaltrade.com' }
];

const addOns = [
  { id: 'extra_storage', name: 'Extra Storage (100GB)', price: 10 },
  { id: 'premium_support', name: 'Premium Support', price: 25 },
  { id: 'custom_domain', name: 'Custom Domain', price: 15 },
  { id: 'advanced_security', name: 'Advanced Security', price: 30 }
];

export default function CreateSubscription() {
  const [formData, setFormData] = useState({
    tenantId: '',
    planId: '',
    billingCycle: 'monthly',
    startDate: new Date().toISOString().split('T')[0],
    customPrice: '',
    discountPercent: '',
    notes: '',
    autoRenew: true,
    trialDays: '14',
    selectedAddOns: []
  });

  const [selectedPlan, setSelectedPlan] = useState(null);
  const [customPricing, setCustomPricing] = useState(false);

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

  const calculateTotalPrice = () => {
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

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log('Creating subscription:', formData);
    // Handle form submission
  };

  return (
    <div className="p-6">
      {/* Page Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Link to="/subscriptions" className="flex items-center gap-2 text-muted-foreground hover:text-foreground">
              <ArrowLeft className="h-4 w-4" />
              Back to Subscriptions
            </Link>
            <Separator orientation="vertical" className="h-6" />
            <Package className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Create New Subscription</h1>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline">
              <X className="h-4 w-4 mr-2" />
              Cancel
            </Button>
            <Button type="submit" form="subscription-form">
              <Save className="h-4 w-4 mr-2" />
              Create Subscription
            </Button>
          </div>
        </div>
      </div>

      <div className="space-y-6">
        <form id="subscription-form" onSubmit={handleSubmit} className="space-y-6">
          {/* Tenant Selection */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="h-5 w-5" />
                Tenant Information
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label htmlFor="tenant">Select Tenant</Label>
                <Select value={formData.tenantId} onValueChange={(value) => handleInputChange('tenantId', value)}>
                  <SelectTrigger>
                    <SelectValue placeholder="Choose a tenant" />
                  </SelectTrigger>
                  <SelectContent>
                    {tenantsList.map((tenant) => (
                      <SelectItem key={tenant.id} value={tenant.id}>
                        <div>
                          <div className="font-medium">{tenant.name}</div>
                          <div className="text-sm text-muted-foreground">{tenant.email}</div>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </CardContent>
          </Card>

          {/* Plan Selection */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Package className="h-5 w-5" />
                Subscription Plan
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
                    <p className="text-sm text-muted-foreground mb-3">{plan.description}</p>
                    <ul className="space-y-1">
                      {plan.features.map((feature, index) => (
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

          {/* Billing Configuration */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <DollarSign className="h-5 w-5" />
                Billing Configuration
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
                <div>
                  <Label htmlFor="startDate">Start Date</Label>
                  <Input
                    type="date"
                    value={formData.startDate}
                    onChange={(e) => handleInputChange('startDate', e.target.value)}
                  />
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

          {/* Add-ons */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Plus className="h-5 w-5" />
                Add-ons
              </CardTitle>
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
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Trial and Auto-renewal */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Settings className="h-5 w-5" />
                Additional Settings
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="trialDays">Trial Period (days)</Label>
                  <Input
                    type="number"
                    min="0"
                    value={formData.trialDays}
                    onChange={(e) => handleInputChange('trialDays', e.target.value)}
                  />
                </div>
                <div className="flex items-center space-x-2 pt-6">
                  <Switch
                    checked={formData.autoRenew}
                    onCheckedChange={(checked) => handleInputChange('autoRenew', checked)}
                  />
                  <Label>Auto-renewal</Label>
                </div>
              </div>
              <div>
                <Label htmlFor="notes">Notes</Label>
                <Textarea
                  value={formData.notes}
                  onChange={(e) => handleInputChange('notes', e.target.value)}
                  placeholder="Add any additional notes or special instructions..."
                  rows={3}
                />
              </div>
            </CardContent>
          </Card>

          {/* Summary */}
          {selectedPlan && (
            <Card>
              <CardHeader>
                <CardTitle>Subscription Summary</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span>Plan:</span>
                    <span className="font-medium">{selectedPlan.name}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Billing Cycle:</span>
                    <span className="font-medium capitalize">{formData.billingCycle}</span>
                  </div>
                  {formData.selectedAddOns.length > 0 && (
                    <div>
                      <div className="font-medium mb-2">Add-ons:</div>
                      {formData.selectedAddOns.map(addOnId => {
                        const addOn = addOns.find(a => a.id === addOnId);
                        return (
                          <div key={addOnId} className="flex justify-between text-sm ml-4">
                            <span>{addOn?.name}</span>
                            <span>${addOn?.price}/month</span>
                          </div>
                        );
                      })}
                    </div>
                  )}
                  <Separator />
                  <div className="flex justify-between text-lg font-semibold">
                    <span>Total:</span>
                    <span>${calculateTotalPrice().toFixed(2)}/{formData.billingCycle === 'yearly' ? 'year' : 'month'}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </form>
      </div>
    </div>
  );
}