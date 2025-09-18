import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  ArrowLeft,
  Save,
  User,
  Building,
  Tag,
  AlertTriangle,
  MessageSquare,
  Calendar,
  Users,
  Clock,
  CheckCircle
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { toast } from 'sonner';

export default function EditTicket() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    customer: '',
    customerEmail: '',
    tenant: '',
    priority: '',
    category: '',
    status: '',
    assignee: '',
    created: '',
    lastUpdate: ''
  });

  const [errors, setErrors] = useState({});

  // Mock data for the ticket
  const mockTicket = {
    id: 'TKT-001',
    title: 'Unable to access billing dashboard',
    description: 'Customer is experiencing issues accessing the billing dashboard. The page loads but shows a blank screen after login. This started happening yesterday after the latest update.',
    customer: 'John Smith',
    customerEmail: 'john.smith@techcorp.com',
    tenant: 'TechCorp Solutions',
    priority: 'high',
    status: 'open',
    category: 'billing',
    assignee: 'Sarah Johnson',
    created: '2024-01-15',
    lastUpdate: '2024-01-15'
  };

  useEffect(() => {
    // Simulate loading ticket data
    const loadTicket = async () => {
      try {
        await new Promise(resolve => setTimeout(resolve, 1000));
        setFormData(mockTicket);
      } catch (error) {
        toast.error('Failed to load ticket data');
      } finally {
        setIsLoading(false);
      }
    };

    loadTicket();
  }, [id]);

  const validateForm = () => {
    const newErrors = {};
    
    if (!formData.title.trim()) newErrors.title = 'Title is required';
    if (!formData.description.trim()) newErrors.description = 'Description is required';
    if (!formData.customer.trim()) newErrors.customer = 'Customer name is required';
    if (!formData.customerEmail.trim()) {
      newErrors.customerEmail = 'Customer email is required';
    } else if (!/\S+@\S+\.\S+/.test(formData.customerEmail)) {
      newErrors.customerEmail = 'Please enter a valid email address';
    }
    if (!formData.tenant) newErrors.tenant = 'Tenant is required';
    if (!formData.priority) newErrors.priority = 'Priority is required';
    if (!formData.category) newErrors.category = 'Category is required';
    if (!formData.status) newErrors.status = 'Status is required';
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleInputChange = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
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
      
      toast.success('Support ticket updated successfully!');
      navigate('/support');
    } catch (error) {
      toast.error('Failed to update ticket. Please try again.');
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

  const assignees = [
    'Sarah Johnson',
    'David Wilson',
    'Mike Chen',
    'Emma Davis'
  ];

  const getStatusBadge = (status) => {
    switch (status) {
      case 'open':
        return <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">Open</Badge>;
      case 'in-progress':
        return <Badge variant="outline" className="bg-blue-50 text-blue-700 border-blue-200">In Progress</Badge>;
      case 'resolved':
        return <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">Resolved</Badge>;
      case 'closed':
        return <Badge variant="outline" className="bg-gray-50 text-gray-700 border-gray-200">Closed</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  const getPriorityBadge = (priority) => {
    switch (priority) {
      case 'high':
        return <Badge variant="destructive">High</Badge>;
      case 'medium':
        return <Badge variant="outline" className="bg-orange-50 text-orange-700 border-orange-200">Medium</Badge>;
      case 'low':
        return <Badge variant="secondary">Low</Badge>;
      default:
        return <Badge variant="outline">{priority}</Badge>;
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading ticket...</p>
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
              onClick={() => navigate('/support')}
              className="flex items-center gap-2"
            >
              <ArrowLeft className="h-4 w-4" />
              Back to Support
            </Button>
            <div className="h-6 w-px bg-border" />
            <MessageSquare className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Edit Ticket #{formData.id}</h1>
            <div className="flex items-center gap-2">
              {getStatusBadge(formData.status)}
              {getPriorityBadge(formData.priority)}
            </div>
          </div>
          <div className="flex items-center gap-3">
            <Button
              variant="outline"
              onClick={() => navigate('/support')}
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
              {isSubmitting ? 'Updating...' : 'Update Ticket'}
            </Button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-4 md:p-6">
        <div className="max-w-7xl mx-auto space-y-6">
          {/* Ticket Metadata */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="h-5 w-5" />
                Ticket Information
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                <div>
                  <Label className="text-muted-foreground">Ticket ID</Label>
                  <p className="font-medium">{formData.id}</p>
                </div>
                <div>
                  <Label className="text-muted-foreground">Created</Label>
                  <p className="font-medium">{new Date(formData.created).toLocaleDateString()}</p>
                </div>
                <div>
                  <Label className="text-muted-foreground">Last Updated</Label>
                  <p className="font-medium">{new Date(formData.lastUpdate).toLocaleDateString()}</p>
                </div>
                <div>
                  <Label className="text-muted-foreground">Current Status</Label>
                  <div className="mt-1">{getStatusBadge(formData.status)}</div>
                </div>
              </div>
            </CardContent>
          </Card>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Ticket Details */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MessageSquare className="h-5 w-5" />
                  Ticket Details
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="title">Title *</Label>
                    <Input
                      id="title"
                      value={formData.title}
                      onChange={(e) => handleInputChange('title', e.target.value)}
                      placeholder="Brief description of the issue"
                      className={errors.title ? 'border-red-500' : ''}
                    />
                    {errors.title && <p className="text-sm text-red-500">{errors.title}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="category">Category *</Label>
                    <Select value={formData.category} onValueChange={(value) => handleInputChange('category', value)}>
                      <SelectTrigger className={errors.category ? 'border-red-500' : ''}>
                        <SelectValue placeholder="Select category" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="billing">Billing</SelectItem>
                        <SelectItem value="technical">Technical</SelectItem>
                        <SelectItem value="account">Account</SelectItem>
                        <SelectItem value="feature-request">Feature Request</SelectItem>
                        <SelectItem value="bug-report">Bug Report</SelectItem>
                        <SelectItem value="general">General</SelectItem>
                      </SelectContent>
                    </Select>
                    {errors.category && <p className="text-sm text-red-500">{errors.category}</p>}
                  </div>
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="description">Description *</Label>
                  <Textarea
                    id="description"
                    value={formData.description}
                    onChange={(e) => handleInputChange('description', e.target.value)}
                    placeholder="Detailed description of the issue or request"
                    rows={4}
                    className={errors.description ? 'border-red-500' : ''}
                  />
                  {errors.description && <p className="text-sm text-red-500">{errors.description}</p>}
                </div>
              </CardContent>
            </Card>

            {/* Customer Information */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <User className="h-5 w-5" />
                  Customer Information
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="customer">Customer Name *</Label>
                    <Input
                      id="customer"
                      value={formData.customer}
                      onChange={(e) => handleInputChange('customer', e.target.value)}
                      placeholder="Customer full name"
                      className={errors.customer ? 'border-red-500' : ''}
                    />
                    {errors.customer && <p className="text-sm text-red-500">{errors.customer}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="customerEmail">Customer Email *</Label>
                    <Input
                      id="customerEmail"
                      type="email"
                      value={formData.customerEmail}
                      onChange={(e) => handleInputChange('customerEmail', e.target.value)}
                      placeholder="customer@example.com"
                      className={errors.customerEmail ? 'border-red-500' : ''}
                    />
                    {errors.customerEmail && <p className="text-sm text-red-500">{errors.customerEmail}</p>}
                  </div>
                </div>
                
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
              </CardContent>
            </Card>

            {/* Status & Assignment */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Tag className="h-5 w-5" />
                  Status & Assignment
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="status">Status *</Label>
                    <Select value={formData.status} onValueChange={(value) => handleInputChange('status', value)}>
                      <SelectTrigger className={errors.status ? 'border-red-500' : ''}>
                        <SelectValue placeholder="Select status" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="open">Open</SelectItem>
                        <SelectItem value="in-progress">In Progress</SelectItem>
                        <SelectItem value="resolved">Resolved</SelectItem>
                        <SelectItem value="closed">Closed</SelectItem>
                      </SelectContent>
                    </Select>
                    {errors.status && <p className="text-sm text-red-500">{errors.status}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="priority">Priority *</Label>
                    <Select value={formData.priority} onValueChange={(value) => handleInputChange('priority', value)}>
                      <SelectTrigger className={errors.priority ? 'border-red-500' : ''}>
                        <SelectValue placeholder="Select priority level" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="low">Low</SelectItem>
                        <SelectItem value="medium">Medium</SelectItem>
                        <SelectItem value="high">High</SelectItem>
                      </SelectContent>
                    </Select>
                    {errors.priority && <p className="text-sm text-red-500">{errors.priority}</p>}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="assignee">Assignee</Label>
                    <Select value={formData.assignee} onValueChange={(value) => handleInputChange('assignee', value)}>
                      <SelectTrigger>
                        <SelectValue placeholder="Assign to team member" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="unassigned">Unassigned</SelectItem>
                        {assignees.map((assignee) => (
                          <SelectItem key={assignee} value={assignee}>{assignee}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
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