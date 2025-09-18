import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
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
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
// DatePicker component not available, using Input with type="date"
import { toast } from 'sonner';

export default function EditInvoice() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState({
    invoiceNumber: '',
    tenant: '',
    billingPeriod: '',
    dueDate: null,
    plan: '',
    status: '',
    notes: '',
    items: [
      { description: '', quantity: 1, unitPrice: 0, amount: 0 }
    ]
  });

  const [errors, setErrors] = useState({});

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
          dueDate: new Date('2024-02-15'),
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
        
        setFormData(mockInvoice);
      } catch (error) {
        toast.error('Failed to load invoice data');
        navigate('/billing');
      } finally {
        setIsLoading(false);
      }
    };

    loadInvoice();
  }, [id, navigate]);

  const validateForm = () => {
    const newErrors = {};
    
    if (!formData.tenant) newErrors.tenant = 'Tenant is required';
    if (!formData.billingPeriod.trim()) newErrors.billingPeriod = 'Billing period is required';
    if (!formData.dueDate) newErrors.dueDate = 'Due date is required';
    if (!formData.plan) newErrors.plan = 'Plan is required';
    if (!formData.status) newErrors.status = 'Status is required';
    
    // Validate items
    const itemErrors = [];
    formData.items.forEach((item, index) => {
      const itemError = {};
      if (!item.description.trim()) itemError.description = 'Description is required';
      if (item.quantity <= 0) itemError.quantity = 'Quantity must be greater than 0';
      if (item.unitPrice <= 0) itemError.unitPrice = 'Unit price must be greater than 0';
      if (Object.keys(itemError).length > 0) {
        itemErrors[index] = itemError;
      }
    });
    
    if (itemErrors.length > 0) newErrors.items = itemErrors;
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleInputChange = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  const handleItemChange = (index, field, value) => {
    const updatedItems = [...formData.items];
    updatedItems[index] = { ...updatedItems[index], [field]: value };
    
    // Calculate amount for this item
    if (field === 'quantity' || field === 'unitPrice') {
      const quantity = field === 'quantity' ? value : updatedItems[index].quantity;
      const unitPrice = field === 'unitPrice' ? value : updatedItems[index].unitPrice;
      updatedItems[index].amount = quantity * unitPrice;
    }
    
    setFormData(prev => ({ ...prev, items: updatedItems }));
    
    // Clear item errors
    if (errors.items && errors.items[index] && errors.items[index][field]) {
      const newItemErrors = [...(errors.items || [])];
      delete newItemErrors[index][field];
      if (Object.keys(newItemErrors[index] || {}).length === 0) {
        newItemErrors[index] = undefined;
      }
      setErrors(prev => ({ ...prev, items: newItemErrors }));
    }
  };

  const addItem = () => {
    setFormData(prev => ({
      ...prev,
      items: [...prev.items, { description: '', quantity: 1, unitPrice: 0, amount: 0 }]
    }));
  };

  const removeItem = (index) => {
    if (formData.items.length > 1) {
      const updatedItems = formData.items.filter((_, i) => i !== index);
      setFormData(prev => ({ ...prev, items: updatedItems }));
    }
  };

  const calculateTotal = () => {
    return formData.items.reduce((total, item) => total + item.amount, 0);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      toast.error('Please fix the errors before submitting');
      return;
    }

    setIsSubmitting(true);
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1500));
      
      toast.success('Invoice updated successfully!');
      navigate('/billing');
    } catch (error) {
      toast.error('Failed to update invoice. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const tenants = [
    'TechCorp Solutions',
    'StartupHub',
    'RetailMax',
    'GlobalTrade Inc',
    'EcoFriendly Co'
  ];

  const plans = [
    { value: 'basic', label: 'Basic - $29/month', price: 29 },
    { value: 'professional', label: 'Professional - $99/month', price: 99 },
    { value: 'enterprise', label: 'Enterprise - $299/month', price: 299 },
    { value: 'enterprise-annual', label: 'Enterprise Annual - $2988/year', price: 2988 }
  ];

  const statuses = [
    { value: 'draft', label: 'Draft', color: 'bg-gray-100 text-gray-800' },
    { value: 'pending', label: 'Pending', color: 'bg-yellow-100 text-yellow-800' },
    { value: 'paid', label: 'Paid', color: 'bg-green-100 text-green-800' },
    { value: 'overdue', label: 'Overdue', color: 'bg-red-100 text-red-800' },
    { value: 'cancelled', label: 'Cancelled', color: 'bg-gray-100 text-gray-800' }
  ];

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  };

  const getStatusBadge = (status) => {
    const statusConfig = statuses.find(s => s.value === status);
    return statusConfig ? (
      <Badge className={statusConfig.color}>
        {statusConfig.label}
      </Badge>
    ) : null;
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading invoice...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
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
            <div className="h-6 w-px bg-border" />
            <FileText className="h-5 w-5 text-primary" />
            <div>
              <h1 className="text-xl font-semibold">Edit Invoice</h1>
              <p className="text-sm text-muted-foreground">{formData.invoiceNumber}</p>
            </div>
            {getStatusBadge(formData.status)}
          </div>
          <div className="flex items-center gap-3">
            <Button
              variant="outline"
              onClick={() => navigate(`/billing/invoice/${id}`)}
              className="flex items-center gap-2"
            >
              <Eye className="h-4 w-4" />
              View
            </Button>
            <Button
              variant="outline"
              onClick={() => navigate('/billing')}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              onClick={handleSubmit}
              disabled={isSubmitting}
              className="flex items-center gap-2"
            >
              <Save className="h-4 w-4" />
              {isSubmitting ? 'Updating...' : 'Update Invoice'}
            </Button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6">
        <div className="max-w-4xl mx-auto space-y-6">
          <form onSubmit={handleSubmit} className="space-y-6">
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
                    <Select value={formData.tenant} onValueChange={(value) => handleInputChange('tenant', value)}>
                      <SelectTrigger className={errors.tenant ? 'border-red-500' : ''}>
                        <SelectValue placeholder="Select tenant organization" />
                      </SelectTrigger>
                      <SelectContent>
                        {tenants.map((tenant) => (
                          <SelectItem key={tenant} value={tenant}>{tenant}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {errors.tenant && <p className="text-sm text-red-500">{errors.tenant}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="plan">Subscription Plan *</Label>
                    <Select value={formData.plan} onValueChange={(value) => handleInputChange('plan', value)}>
                      <SelectTrigger className={errors.plan ? 'border-red-500' : ''}>
                        <SelectValue placeholder="Select subscription plan" />
                      </SelectTrigger>
                      <SelectContent>
                        {plans.map((plan) => (
                          <SelectItem key={plan.value} value={plan.value}>{plan.label}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {errors.plan && <p className="text-sm text-red-500">{errors.plan}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="status">Status *</Label>
                    <Select value={formData.status} onValueChange={(value) => handleInputChange('status', value)}>
                      <SelectTrigger className={errors.status ? 'border-red-500' : ''}>
                        <SelectValue placeholder="Select invoice status" />
                      </SelectTrigger>
                      <SelectContent>
                        {statuses.map((status) => (
                          <SelectItem key={status.value} value={status.value}>{status.label}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {errors.status && <p className="text-sm text-red-500">{errors.status}</p>}
                  </div>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="billingPeriod">Billing Period *</Label>
                    <Input
                      id="billingPeriod"
                      value={formData.billingPeriod}
                      onChange={(e) => handleInputChange('billingPeriod', e.target.value)}
                      placeholder="e.g., January 2024 or Jan 2024 - Dec 2024"
                      className={errors.billingPeriod ? 'border-red-500' : ''}
                    />
                    {errors.billingPeriod && <p className="text-sm text-red-500">{errors.billingPeriod}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="dueDate">Due Date *</Label>
                    <Input
                      type="date"
                      value={formData.dueDate ? formData.dueDate.toISOString().split('T')[0] : ''}
                      onChange={(e) => handleInputChange('dueDate', e.target.value ? new Date(e.target.value) : null)}
                      className={errors.dueDate ? 'border-red-500' : ''}
                    />
                    {errors.dueDate && <p className="text-sm text-red-500">{errors.dueDate}</p>}
                  </div>
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="notes">Notes</Label>
                  <Textarea
                    id="notes"
                    value={formData.notes}
                    onChange={(e) => handleInputChange('notes', e.target.value)}
                    placeholder="Additional notes or terms for this invoice"
                    rows={3}
                  />
                </div>
              </CardContent>
            </Card>

            {/* Invoice Items */}
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle className="flex items-center gap-2">
                    <CreditCard className="h-5 w-5" />
                    Invoice Items
                  </CardTitle>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={addItem}
                    className="flex items-center gap-2"
                  >
                    <Plus className="h-4 w-4" />
                    Add Item
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                {formData.items.map((item, index) => (
                  <div key={index} className="border rounded-lg p-4 space-y-4">
                    <div className="flex items-center justify-between">
                      <h4 className="font-medium">Item {index + 1}</h4>
                      {formData.items.length > 1 && (
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => removeItem(index)}
                          className="text-red-600 hover:text-red-700"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      )}
                    </div>
                    
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                      <div className="md:col-span-2 space-y-2">
                        <Label>Description *</Label>
                        <Input
                          value={item.description}
                          onChange={(e) => handleItemChange(index, 'description', e.target.value)}
                          placeholder="Item description"
                          className={errors.items?.[index]?.description ? 'border-red-500' : ''}
                        />
                        {errors.items?.[index]?.description && (
                          <p className="text-sm text-red-500">{errors.items[index].description}</p>
                        )}
                      </div>
                      <div className="space-y-2">
                        <Label>Quantity *</Label>
                        <Input
                          type="number"
                          min="1"
                          value={item.quantity}
                          onChange={(e) => handleItemChange(index, 'quantity', parseInt(e.target.value) || 0)}
                          className={errors.items?.[index]?.quantity ? 'border-red-500' : ''}
                        />
                        {errors.items?.[index]?.quantity && (
                          <p className="text-sm text-red-500">{errors.items[index].quantity}</p>
                        )}
                      </div>
                      <div className="space-y-2">
                        <Label>Unit Price *</Label>
                        <Input
                          type="number"
                          min="0"
                          step="0.01"
                          value={item.unitPrice}
                          onChange={(e) => handleItemChange(index, 'unitPrice', parseFloat(e.target.value) || 0)}
                          className={errors.items?.[index]?.unitPrice ? 'border-red-500' : ''}
                        />
                        {errors.items?.[index]?.unitPrice && (
                          <p className="text-sm text-red-500">{errors.items[index].unitPrice}</p>
                        )}
                      </div>
                    </div>
                    
                    <div className="flex justify-end">
                      <div className="text-right">
                        <Label className="text-muted-foreground">Amount</Label>
                        <p className="text-lg font-semibold">{formatCurrency(item.amount)}</p>
                      </div>
                    </div>
                  </div>
                ))}
                
                {/* Total */}
                <div className="border-t pt-4">
                  <div className="flex justify-end">
                    <div className="text-right">
                      <Label className="text-muted-foreground">Total Amount</Label>
                      <p className="text-2xl font-bold text-primary">{formatCurrency(calculateTotal())}</p>
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