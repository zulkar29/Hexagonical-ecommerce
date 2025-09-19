import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useForm, useFieldArray } from 'react-hook-form';
import {
  ArrowLeft,
  Save,
  Building,
  CreditCard,
  Calendar,
  DollarSign,
  FileText,
  Plus,
  Trash2,
  Eye
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { toast } from 'sonner';

export default function EditInvoice() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { 
    register, 
    handleSubmit, 
    control, 
    watch, 
    setValue, 
    formState: { errors } 
  } = useForm({
    defaultValues: {
      invoiceNumber: '',
      tenant: '',
      billingPeriod: '',
      dueDate: '',
      plan: '',
      status: '',
      notes: '',
      items: [
        { description: '', quantity: 1, unitPrice: 0, amount: 0 }
      ]
    }
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'items'
  });

  const watchedItems = watch('items');

  // Mock data for existing invoice
  useEffect(() => {
    const loadInvoice = async () => {
      setIsLoading(true);
      try {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // Mock invoice data
        const mockInvoice = {
          invoiceNumber: `INV-${id}`,
          tenant: 'TechCorp Solutions',
          billingPeriod: 'January 2024',
          dueDate: '2024-02-15',
          plan: 'professional',
          status: 'pending',
          notes: 'Professional plan subscription for January 2024',
          items: [
            {
              description: 'Professional Plan Subscription',
              quantity: 1,
              unitPrice: 99,
              amount: 99
            },
            {
              description: 'Additional User Licenses (5 users)',
              quantity: 5,
              unitPrice: 15,
              amount: 75
            }
          ]
        };
        
        // Set form values using React Hook Form
        Object.keys(mockInvoice).forEach(key => {
          setValue(key, mockInvoice[key]);
        });
      } catch (error) {
        toast.error('Failed to load invoice data');
        navigate('/billing');
      } finally {
        setIsLoading(false);
      }
    };

    loadInvoice();
  }, [id, navigate, setValue]);

  // Calculate amount when quantity or unitPrice changes
  useEffect(() => {
    watchedItems?.forEach((item, index) => {
      const amount = (item.quantity || 0) * (item.unitPrice || 0);
      if (item.amount !== amount) {
        setValue(`items.${index}.amount`, amount);
      }
    });
  }, [watchedItems, setValue]);

  const addItem = () => {
    append({ description: '', quantity: 1, unitPrice: 0, amount: 0 });
  };

  const removeItem = (index) => {
    if (fields.length > 1) {
      remove(index);
    }
  };

  const onSubmit = async (data) => {
    setIsSubmitting(true);
    try {
      console.log('Submitting invoice data:', data);
      // Handle form submission with API call
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call
      
      toast.success('Invoice updated successfully!');
      navigate('/billing');
    } catch (error) {
      console.error('Error updating invoice:', error);
      toast.error('Failed to update invoice');
    } finally {
      setIsSubmitting(false);
    }
  };

  const calculateSubtotal = () => {
    return watchedItems?.reduce((sum, item) => sum + (item.amount || 0), 0) || 0;
  };

  const calculateTax = () => {
    return calculateSubtotal() * 0.15; // 15% tax
  };

  const calculateTotal = () => {
    return calculateSubtotal() + calculateTax();
  };

  const tenants = [
    'TechCorp Solutions',
    'StartupHub',
    'RetailMax', 
    'GlobalTrade Inc',
    'EcoFriendly Co'
  ];

  const plans = [
    { value: 'starter', label: 'Starter Plan ($29/month)' },
    { value: 'professional', label: 'Professional Plan ($79/month)' },
    { value: 'enterprise', label: 'Enterprise Plan ($199/month)' },
    { value: 'custom', label: 'Custom Plan' }
  ];

  const statuses = [
    { value: 'draft', label: 'Draft' },
    { value: 'pending', label: 'Pending' },
    { value: 'paid', label: 'Paid' },
    { value: 'overdue', label: 'Overdue' },
    { value: 'cancelled', label: 'Cancelled' }
  ];

  const getStatusBadge = (status) => {
    const statusConfig = {
      draft: { variant: 'secondary', label: 'Draft' },
      pending: { variant: 'default', label: 'Pending' },
      paid: { variant: 'default', label: 'Paid' },
      overdue: { variant: 'destructive', label: 'Overdue' },
      cancelled: { variant: 'outline', label: 'Cancelled' }
    };
    
    const config = statusConfig[status] || { variant: 'secondary', label: status };
    return <Badge variant={config.variant}>{config.label}</Badge>;
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-2 text-muted-foreground">Loading invoice...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-4 py-3">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate('/billing')}
              className="flex items-center gap-2"
            >
              <ArrowLeft className="h-4 w-4" />
              Back to Billing
            </Button>
            <div className="flex items-center gap-3">
              <FileText className="h-5 w-5 text-blue-600" />
              <div>
                <h1 className="text-xl font-semibold text-gray-900">Edit Invoice</h1>
                <p className="text-sm text-muted-foreground">{watch('invoiceNumber')}</p>
              </div>
            </div>
            {getStatusBadge(watch('status'))}
          </div>
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm">
              <Eye className="h-4 w-4 mr-2" />
              Preview
            </Button>
            <Button 
              type="submit" 
              form="invoice-form"
              disabled={isSubmitting}
              size="sm"
            >
              <Save className="h-4 w-4" />
              {isSubmitting ? 'Updating...' : 'Update Invoice'}
            </Button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-4 md:p-6">
        <div className="max-w-7xl mx-auto space-y-6">
          <form id="invoice-form" onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Invoice Information */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  Invoice Information
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="tenant">Tenant *</Label>
                    <select
                      {...register("tenant", { required: "Tenant is required" })}
                      className={`w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${errors.tenant ? 'border-red-500' : ''}`}
                    >
                      <option value="">Select tenant</option>
                      {tenants.map((tenant) => (
                        <option key={tenant} value={tenant}>{tenant}</option>
                      ))}
                    </select>
                    {errors.tenant && <p className="text-sm text-red-500">{errors.tenant.message}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="plan">Plan *</Label>
                    <select
                      {...register("plan", { required: "Plan is required" })}
                      className={`w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${errors.plan ? 'border-red-500' : ''}`}
                    >
                      <option value="">Select plan</option>
                      {plans.map((plan) => (
                        <option key={plan.value} value={plan.value}>{plan.label}</option>
                      ))}
                    </select>
                    {errors.plan && <p className="text-sm text-red-500">{errors.plan.message}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="status">Status *</Label>
                    <select
                      {...register("status", { required: "Status is required" })}
                      className={`w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${errors.status ? 'border-red-500' : ''}`}
                    >
                      <option value="">Select status</option>
                      {statuses.map((status) => (
                        <option key={status.value} value={status.value}>{status.label}</option>
                      ))}
                    </select>
                    {errors.status && <p className="text-sm text-red-500">{errors.status.message}</p>}
                  </div>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="billingPeriod">Billing Period *</Label>
                    <Input
                      id="billingPeriod"
                      {...register("billingPeriod", { required: "Billing period is required" })}
                      placeholder="e.g., January 2024 or Jan 2024 - Dec 2024"
                      className={errors.billingPeriod ? 'border-red-500' : ''}
                    />
                    {errors.billingPeriod && <p className="text-sm text-red-500">{errors.billingPeriod.message}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="dueDate">Due Date *</Label>
                    <Input
                      type="date"
                      {...register("dueDate", { required: "Due date is required" })}
                      className={errors.dueDate ? 'border-red-500' : ''}
                    />
                    {errors.dueDate && <p className="text-sm text-red-500">{errors.dueDate.message}</p>}
                  </div>
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="notes">Notes</Label>
                  <Textarea
                    id="notes"
                    {...register("notes")}
                    placeholder="Additional notes or terms for this invoice"
                    rows={3}
                  />
                </div>
              </CardContent>
            </Card>

            {/* Line Items */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <DollarSign className="h-5 w-5" />
                    Line Items
                  </div>
                  <Button type="button" variant="outline" size="sm" onClick={addItem}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Item
                  </Button>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {fields.map((field, index) => (
                    <div key={field.id} className="grid grid-cols-12 gap-4 items-start p-4 border border-gray-200 rounded-lg">
                      <div className="col-span-5">
                        <Label htmlFor={`items.${index}.description`}>Description *</Label>
                        <Input
                          {...register(`items.${index}.description`, { required: "Description is required" })}
                          placeholder="Item description"
                          className={errors.items?.[index]?.description ? 'border-red-500' : ''}
                        />
                        {errors.items?.[index]?.description && (
                          <p className="text-sm text-red-500 mt-1">{errors.items[index].description.message}</p>
                        )}
                      </div>
                      <div className="col-span-2">
                        <Label htmlFor={`items.${index}.quantity`}>Quantity *</Label>
                        <Input
                          type="number"
                          min="1"
                          step="1"
                          {...register(`items.${index}.quantity`, { 
                            required: "Quantity is required",
                            min: { value: 1, message: "Quantity must be at least 1" }
                          })}
                          className={errors.items?.[index]?.quantity ? 'border-red-500' : ''}
                        />
                        {errors.items?.[index]?.quantity && (
                          <p className="text-sm text-red-500 mt-1">{errors.items[index].quantity.message}</p>
                        )}
                      </div>
                      <div className="col-span-2">
                        <Label htmlFor={`items.${index}.unitPrice`}>Unit Price *</Label>
                        <Input
                          type="number"
                          min="0"
                          step="0.01"
                          {...register(`items.${index}.unitPrice`, { 
                            required: "Unit price is required",
                            min: { value: 0, message: "Unit price must be positive" }
                          })}
                          className={errors.items?.[index]?.unitPrice ? 'border-red-500' : ''}
                        />
                        {errors.items?.[index]?.unitPrice && (
                          <p className="text-sm text-red-500 mt-1">{errors.items[index].unitPrice.message}</p>
                        )}
                      </div>
                      <div className="col-span-2">
                        <Label>Amount</Label>
                        <Input
                          {...register(`items.${index}.amount`)}
                          readOnly
                          className="bg-gray-50"
                        />
                      </div>
                      <div className="col-span-1 flex justify-end pt-6">
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          onClick={() => removeItem(index)}
                          disabled={fields.length === 1}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  ))}
                  
                  {/* Summary */}
                  <div className="border-t pt-4">
                    <div className="flex justify-end">
                      <div className="w-64 space-y-2">
                        <div className="flex justify-between">
                          <span>Subtotal:</span>
                          <span>${calculateSubtotal().toFixed(2)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span>Tax (15%):</span>
                          <span>${calculateTax().toFixed(2)}</span>
                        </div>
                        <div className="flex justify-between font-semibold text-lg border-t pt-2">
                          <span>Total:</span>
                          <span>${calculateTotal().toFixed(2)}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </form>
        </div>
      </div>
    </div>
  );
}