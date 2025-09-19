import React, { useState } from 'react';
import { useForm, useFieldArray } from 'react-hook-form';
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
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import { Link } from 'react-router-dom';
import { toast } from 'sonner';

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
  const [selectedPlan, setSelectedPlan] = useState(null);
  const [customPricing, setCustomPricing] = useState(false);

  const { 
    register, 
    handleSubmit, 
    watch, 
    setValue, 
    formState: { errors, isSubmitting } 
  } = useForm({
    defaultValues: {
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
    }
  });

  const watchedValues = watch();

  const handlePlanSelect = (planId) => {
    const plan = subscriptionPlans.find(p => p.id === planId);
    setSelectedPlan(plan);
    setValue('planId', planId);
  };

  const handleAddOnToggle = (addOnId) => {
    const currentAddOns = watchedValues.selectedAddOns || [];
    const newAddOns = currentAddOns.includes(addOnId)
      ? currentAddOns.filter(id => id !== addOnId)
      : [...currentAddOns, addOnId];
    setValue('selectedAddOns', newAddOns);
  };

  const calculateTotalPrice = () => {
    if (!selectedPlan) return 0;
    
    let basePrice = customPricing && watchedValues.customPrice 
      ? parseFloat(watchedValues.customPrice) 
      : selectedPlan.basePrice;
    
    if (watchedValues.billingCycle === 'yearly') {
      basePrice = basePrice * 12 * 0.9; // 10% discount for yearly
    }
    
    const addOnsPrice = (watchedValues.selectedAddOns || []).reduce((total, addOnId) => {
      const addOn = addOns.find(a => a.id === addOnId);
      return total + (addOn ? addOn.price : 0);
    }, 0);
    
    let total = basePrice + addOnsPrice;
    
    if (watchedValues.discountPercent) {
      total = total * (1 - parseFloat(watchedValues.discountPercent) / 100);
    }
    
    return total;
  };

  const onSubmit = async (data) => {
    try {
      console.log('Creating subscription:', data);
      // Handle form submission with API call
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call
      
      toast.success('Subscription created successfully!');
      // navigate('/subscriptions');
    } catch (error) {
      console.error('Error creating subscription:', error);
      toast.error('Failed to create subscription');
    }
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
            <Button type="submit" form="subscription-form" disabled={isSubmitting}>
              <Save className="h-4 w-4 mr-2" />
              {isSubmitting ? 'Creating...' : 'Create Subscription'}
            </Button>
          </div>
        </div>
      </div>

      <form id="subscription-form" onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Tenant Information */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              Tenant Information
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="tenantId">Select Tenant *</Label>
              <select
                {...register("tenantId", { required: "Tenant is required" })}
                className={`w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${errors.tenantId ? 'border-red-500' : ''}`}
              >
                <option value="">Choose a tenant organization</option>
                {tenantsList.map((tenant) => (
                  <option key={tenant.id} value={tenant.id}>
                    {tenant.name} ({tenant.email})
                  </option>
                ))}
              </select>
              {errors.tenantId && <p className="text-sm text-red-500">{errors.tenantId.message}</p>}
            </div>
          </CardContent>
        </Card>

        {/* Plan Selection */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5" />
              Select Subscription Plan
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {subscriptionPlans.map((plan) => (
                <div
                  key={plan.id}
                  className={`relative p-4 border-2 rounded-lg cursor-pointer transition-all ${
                    watchedValues.planId === plan.id
                      ? 'border-blue-500 bg-blue-50'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                  onClick={() => handlePlanSelect(plan.id)}
                >
                  {watchedValues.planId === plan.id && (
                    <div className="absolute -top-2 -right-2">
                      <Badge className="bg-blue-500">Selected</Badge>
                    </div>
                  )}
                  <h3 className="font-semibold text-lg">{plan.name}</h3>
                  <p className="text-gray-600 text-sm mb-2">{plan.description}</p>
                  <div className="text-2xl font-bold text-blue-600 mb-3">${plan.basePrice}/month</div>
                  <ul className="space-y-1 text-sm">
                    {plan.features.map((feature, index) => (
                      <li key={index} className="flex items-center gap-2">
                        <div className="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
                        {feature}
                      </li>
                    ))}
                  </ul>
                </div>
              ))}
            </div>
            {errors.planId && <p className="text-sm text-red-500 mt-2">Please select a plan</p>}
          </CardContent>
        </Card>

        {selectedPlan && (
          <>
            {/* Billing Configuration */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Calendar className="h-5 w-5" />
                  Billing Configuration
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="billingCycle">Billing Cycle</Label>
                    <select
                      {...register("billingCycle")}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="monthly">Monthly</option>
                      <option value="yearly">Yearly (10% discount)</option>
                    </select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="startDate">Start Date</Label>
                    <Input
                      type="date"
                      {...register("startDate", { required: "Start date is required" })}
                      className={errors.startDate ? 'border-red-500' : ''}
                    />
                    {errors.startDate && <p className="text-sm text-red-500">{errors.startDate.message}</p>}
                  </div>
                </div>

                <div className="flex items-center gap-4">
                  <div className="flex items-center space-x-2">
                    <Switch
                      checked={customPricing}
                      onCheckedChange={setCustomPricing}
                    />
                    <Label>Custom Pricing</Label>
                  </div>
                  {customPricing && (
                    <div className="flex-1 max-w-xs">
                      <Input
                        type="number"
                        min="0"
                        step="0.01"
                        {...register("customPrice")}
                        placeholder="Enter custom price"
                      />
                    </div>
                  )}
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="trialDays">Trial Period (days)</Label>
                    <Input
                      type="number"
                      min="0"
                      max="90"
                      {...register("trialDays")}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="discountPercent">Discount (%)</Label>
                    <Input
                      type="number"
                      min="0"
                      max="100"
                      step="0.01"
                      {...register("discountPercent")}
                      placeholder="0.00"
                    />
                  </div>
                </div>

                <div className="flex items-center space-x-2">
                  <Switch
                    {...register("autoRenew")}
                    defaultChecked={true}
                  />
                  <Label>Auto-renewal enabled</Label>
                </div>
              </CardContent>
            </Card>

            {/* Add-ons */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Settings className="h-5 w-5" />
                  Add-on Services
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {addOns.map((addOn) => (
                    <div
                      key={addOn.id}
                      className={`p-4 border rounded-lg cursor-pointer transition-all ${
                        (watchedValues.selectedAddOns || []).includes(addOn.id)
                          ? 'border-blue-500 bg-blue-50'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                      onClick={() => handleAddOnToggle(addOn.id)}
                    >
                      <div className="flex items-center justify-between">
                        <div>
                          <h4 className="font-medium">{addOn.name}</h4>
                          <p className="text-sm text-gray-600">+${addOn.price}/month</p>
                        </div>
                        <div className={`w-4 h-4 rounded border-2 ${
                          (watchedValues.selectedAddOns || []).includes(addOn.id)
                            ? 'bg-blue-500 border-blue-500'
                            : 'border-gray-300'
                        }`}>
                          {(watchedValues.selectedAddOns || []).includes(addOn.id) && (
                            <div className="w-full h-full flex items-center justify-center">
                              <div className="w-2 h-2 bg-white rounded-sm"></div>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Notes */}
            <Card>
              <CardHeader>
                <CardTitle>Additional Notes</CardTitle>
              </CardHeader>
              <CardContent>
                <Textarea
                  {...register("notes")}
                  placeholder="Add any special terms, conditions, or notes for this subscription..."
                  rows={4}
                />
              </CardContent>
            </Card>

            {/* Pricing Summary */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <DollarSign className="h-5 w-5" />
                  Pricing Summary
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span>{selectedPlan.name} ({watchedValues.billingCycle})</span>
                    <span>
                      ${customPricing && watchedValues.customPrice 
                        ? parseFloat(watchedValues.customPrice || 0).toFixed(2) 
                        : selectedPlan.basePrice.toFixed(2)}
                      {watchedValues.billingCycle === 'yearly' && !customPricing && (
                        <span className="text-green-600 ml-2">(10% yearly discount applied)</span>
                      )}
                    </span>
                  </div>
                  
                  {(watchedValues.selectedAddOns || []).map((addOnId) => {
                    const addOn = addOns.find(a => a.id === addOnId);
                    return addOn ? (
                      <div key={addOnId} className="flex justify-between text-sm">
                        <span>{addOn.name}</span>
                        <span>+${addOn.price.toFixed(2)}</span>
                      </div>
                    ) : null;
                  })}
                  
                  {watchedValues.discountPercent && (
                    <div className="flex justify-between text-green-600">
                      <span>Discount ({watchedValues.discountPercent}%)</span>
                      <span>-${(calculateTotalPrice() * parseFloat(watchedValues.discountPercent) / 100).toFixed(2)}</span>
                    </div>
                  )}
                  
                  <Separator />
                  <div className="flex justify-between text-lg font-semibold">
                    <span>Total {watchedValues.billingCycle === 'yearly' ? 'Annual' : 'Monthly'} Cost:</span>
                    <span>${calculateTotalPrice().toFixed(2)}</span>
                  </div>
                  
                  {watchedValues.trialDays && parseInt(watchedValues.trialDays) > 0 && (
                    <p className="text-sm text-blue-600">
                      *Includes {watchedValues.trialDays}-day free trial
                    </p>
                  )}
                </div>
              </CardContent>
            </Card>
          </>
        )}
      </form>
    </div>
  );
}