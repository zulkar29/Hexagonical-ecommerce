import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Shield,
  ArrowLeft,
  Save,
  AlertTriangle,
  CheckCircle,
  Info
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Separator } from '@/components/ui/separator';
import { Checkbox } from '@/components/ui/checkbox';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Switch } from '@/components/ui/switch';

// Available permissions grouped by category
const availablePermissions = [
  {
    category: 'Tenant Management',
    description: 'Permissions related to managing tenant accounts and organizations',
    permissions: [
      { 
        key: 'tenant_create', 
        name: 'Create Tenants', 
        description: 'Create new tenant accounts',
        risk: 'medium'
      },
      { 
        key: 'tenant_edit', 
        name: 'Edit Tenants', 
        description: 'Modify existing tenant information',
        risk: 'medium'
      },
      { 
        key: 'tenant_delete', 
        name: 'Delete Tenants', 
        description: 'Permanently delete tenant accounts',
        risk: 'high'
      },
      { 
        key: 'tenant_view', 
        name: 'View Tenants', 
        description: 'Read-only access to tenant information',
        risk: 'low'
      },
      { 
        key: 'tenant_suspend', 
        name: 'Suspend/Reactivate Tenants', 
        description: 'Suspend or reactivate tenant accounts',
        risk: 'high'
      }
    ]
  },
  {
    category: 'Subscription Management',
    description: 'Permissions for handling billing, subscriptions, and payments',
    permissions: [
      { 
        key: 'subscription_create', 
        name: 'Create Subscriptions', 
        description: 'Create new subscription plans',
        risk: 'medium'
      },
      { 
        key: 'subscription_edit', 
        name: 'Edit Subscriptions', 
        description: 'Modify subscription details and pricing',
        risk: 'high'
      },
      { 
        key: 'subscription_cancel', 
        name: 'Cancel Subscriptions', 
        description: 'Cancel active subscriptions',
        risk: 'high'
      },
      { 
        key: 'subscription_view', 
        name: 'View Subscriptions', 
        description: 'Read-only access to subscription data',
        risk: 'low'
      },
      { 
        key: 'billing_manage', 
        name: 'Manage Billing', 
        description: 'Handle billing issues and payments',
        risk: 'high'
      }
    ]
  },
  {
    category: 'User Management',
    description: 'Permissions for managing admin users and their access',
    permissions: [
      { 
        key: 'user_create', 
        name: 'Create Admin Users', 
        description: 'Create new admin user accounts',
        risk: 'high'
      },
      { 
        key: 'user_edit', 
        name: 'Edit Admin Users', 
        description: 'Modify admin user information',
        risk: 'medium'
      },
      { 
        key: 'user_delete', 
        name: 'Delete Admin Users', 
        description: 'Remove admin user accounts',
        risk: 'high'
      },
      { 
        key: 'user_view', 
        name: 'View Admin Users', 
        description: 'Read-only access to admin user data',
        risk: 'low'
      },
      { 
        key: 'role_manage', 
        name: 'Manage Roles', 
        description: 'Create and modify user roles and permissions',
        risk: 'high'
      }
    ]
  },
  {
    category: 'Analytics & Reports',
    description: 'Permissions for accessing platform analytics and generating reports',
    permissions: [
      { 
        key: 'analytics_view', 
        name: 'View Analytics', 
        description: 'Access platform analytics and metrics',
        risk: 'low'
      },
      { 
        key: 'analytics_export', 
        name: 'Export Analytics', 
        description: 'Export analytics data and reports',
        risk: 'medium'
      },
      { 
        key: 'reports_generate', 
        name: 'Generate Reports', 
        description: 'Create custom reports and dashboards',
        risk: 'low'
      },
      { 
        key: 'financial_reports', 
        name: 'Financial Reports', 
        description: 'Access financial and revenue reports',
        risk: 'medium'
      }
    ]
  },
  {
    category: 'Support Management',
    description: 'Permissions for handling customer support and tickets',
    permissions: [
      { 
        key: 'support_view', 
        name: 'View Support Tickets', 
        description: 'Read-only access to support tickets',
        risk: 'low'
      },
      { 
        key: 'support_respond', 
        name: 'Respond to Tickets', 
        description: 'Reply to and manage support tickets',
        risk: 'low'
      },
      { 
        key: 'support_escalate', 
        name: 'Escalate Tickets', 
        description: 'Escalate tickets to higher support tiers',
        risk: 'medium'
      },
      { 
        key: 'support_close', 
        name: 'Close Tickets', 
        description: 'Mark support tickets as resolved',
        risk: 'low'
      }
    ]
  },
  {
    category: 'System Administration',
    description: 'High-level system permissions and platform settings',
    permissions: [
      { 
        key: 'system_settings', 
        name: 'System Settings', 
        description: 'Modify platform-wide settings and configuration',
        risk: 'high'
      },
      { 
        key: 'audit_logs', 
        name: 'Audit Logs', 
        description: 'Access system audit logs and security events',
        risk: 'medium'
      },
      { 
        key: 'backup_restore', 
        name: 'Backup & Restore', 
        description: 'Manage system backups and data restoration',
        risk: 'high'
      },
      { 
        key: 'maintenance_mode', 
        name: 'Maintenance Mode', 
        description: 'Enable/disable platform maintenance mode',
        risk: 'high'
      }
    ]
  }
];

const getRiskColor = (risk) => {
  switch (risk) {
    case 'low': return 'bg-green-100 text-green-800 border-green-200';
    case 'medium': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    case 'high': return 'bg-red-100 text-red-800 border-red-200';
    default: return 'bg-gray-100 text-gray-800 border-gray-200';
  }
};

const getRiskIcon = (risk) => {
  switch (risk) {
    case 'low': return <CheckCircle className="h-3 w-3" />;
    case 'medium': return <Info className="h-3 w-3" />;
    case 'high': return <AlertTriangle className="h-3 w-3" />;
    default: return <Info className="h-3 w-3" />;
  }
};

export default function CreateRole() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    isActive: true,
    permissions: []
  });
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleInputChange = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    // Clear error when user starts typing
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: null }));
    }
  };

  const togglePermission = (permissionKey) => {
    setFormData(prev => ({
      ...prev,
      permissions: prev.permissions.includes(permissionKey)
        ? prev.permissions.filter(p => p !== permissionKey)
        : [...prev.permissions, permissionKey]
    }));
  };

  const toggleCategoryPermissions = (category) => {
    const categoryPermissions = category.permissions.map(p => p.key);
    const allSelected = categoryPermissions.every(p => formData.permissions.includes(p));
    
    setFormData(prev => ({
      ...prev,
      permissions: allSelected
        ? prev.permissions.filter(p => !categoryPermissions.includes(p))
        : [...new Set([...prev.permissions, ...categoryPermissions])]
    }));
  };

  const validateForm = () => {
    const newErrors = {};
    
    if (!formData.name.trim()) {
      newErrors.name = 'Role name is required';
    } else if (formData.name.length < 3) {
      newErrors.name = 'Role name must be at least 3 characters';
    }
    
    if (!formData.description.trim()) {
      newErrors.description = 'Role description is required';
    } else if (formData.description.length < 10) {
      newErrors.description = 'Description must be at least 10 characters';
    }
    
    if (formData.permissions.length === 0) {
      newErrors.permissions = 'At least one permission must be selected';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }
    
    setIsSubmitting(true);
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1500));
      
      console.log('Creating role:', formData);
      
      // Navigate back to permissions page
      navigate('/users/permissions');
    } catch (error) {
      console.error('Error creating role:', error);
      setErrors({ submit: 'Failed to create role. Please try again.' });
    } finally {
      setIsSubmitting(false);
    }
  };

  const selectedHighRiskPermissions = formData.permissions.filter(permKey => {
    const permission = availablePermissions
      .flatMap(cat => cat.permissions)
      .find(p => p.key === permKey);
    return permission?.risk === 'high';
  });

  return (
    <div className="space-y-6 p-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Button variant="ghost" onClick={() => navigate('/users/permissions')}>
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back
              </Button>
              <div>
                <h2 className="text-3xl font-bold tracking-tight">Create New Role</h2>
                <p className="text-muted-foreground">
                  Define a new role with specific permissions for admin users
                </p>
              </div>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Basic Information */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Shield className="h-5 w-5" />
                  <span>Role Information</span>
                </CardTitle>
                <CardDescription>
                  Basic details about the role
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="roleName">Role Name *</Label>
                    <Input
                      id="roleName"
                      value={formData.name}
                      onChange={(e) => handleInputChange('name', e.target.value)}
                      placeholder="e.g., Content Manager"
                      className={errors.name ? 'border-red-500' : ''}
                    />
                    {errors.name && (
                      <p className="text-sm text-red-600">{errors.name}</p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="roleStatus">Status</Label>
                    <div className="flex items-center space-x-2">
                      <Switch
                        id="roleStatus"
                        checked={formData.isActive}
                        onCheckedChange={(checked) => handleInputChange('isActive', checked)}
                      />
                      <Label htmlFor="roleStatus" className="text-sm">
                        {formData.isActive ? 'Active' : 'Inactive'}
                      </Label>
                    </div>
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="roleDescription">Description *</Label>
                  <Textarea
                    id="roleDescription"
                    value={formData.description}
                    onChange={(e) => handleInputChange('description', e.target.value)}
                    placeholder="Describe what this role is responsible for and its purpose..."
                    rows={3}
                    className={errors.description ? 'border-red-500' : ''}
                  />
                  {errors.description && (
                    <p className="text-sm text-red-600">{errors.description}</p>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Permissions */}
            <Card>
              <CardHeader>
                <CardTitle>Permissions</CardTitle>
                <CardDescription>
                  Select the permissions this role should have. Be careful with high-risk permissions.
                </CardDescription>
                {errors.permissions && (
                  <Alert className="border-red-200 bg-red-50">
                    <AlertTriangle className="h-4 w-4 text-red-600" />
                    <AlertDescription className="text-red-600">
                      {errors.permissions}
                    </AlertDescription>
                  </Alert>
                )}
              </CardHeader>
              <CardContent>
                <div className="space-y-6">
                  {availablePermissions.map((category) => {
                    const categoryPermissions = category.permissions.map(p => p.key);
                    const selectedInCategory = categoryPermissions.filter(p => formData.permissions.includes(p));
                    const allSelected = selectedInCategory.length === categoryPermissions.length;
                    const someSelected = selectedInCategory.length > 0;
                    
                    return (
                      <div key={category.category} className="space-y-3">
                        <div className="flex items-center justify-between">
                          <div>
                            <div className="flex items-center space-x-2">
                              <Checkbox
                                checked={allSelected}
                                ref={(el) => {
                                  if (el) el.indeterminate = someSelected && !allSelected;
                                }}
                                onCheckedChange={() => toggleCategoryPermissions(category)}
                              />
                              <h3 className="text-lg font-semibold">{category.category}</h3>
                              <Badge variant="outline" className="text-xs">
                                {selectedInCategory.length}/{categoryPermissions.length}
                              </Badge>
                            </div>
                            <p className="text-sm text-muted-foreground ml-6">{category.description}</p>
                          </div>
                        </div>
                        
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 ml-6">
                          {category.permissions.map((permission) => (
                            <div key={permission.key} className="flex items-start space-x-3 p-3 border rounded-lg hover:bg-muted/50">
                              <Checkbox
                                checked={formData.permissions.includes(permission.key)}
                                onCheckedChange={() => togglePermission(permission.key)}
                                className="mt-0.5"
                              />
                              <div className="flex-1 space-y-1">
                                <div className="flex items-center space-x-2">
                                  <Label className="text-sm font-medium cursor-pointer">
                                    {permission.name}
                                  </Label>
                                  <Badge 
                                    variant="outline" 
                                    className={`text-xs ${getRiskColor(permission.risk)}`}
                                  >
                                    {getRiskIcon(permission.risk)}
                                    <span className="ml-1">{permission.risk}</span>
                                  </Badge>
                                </div>
                                <p className="text-xs text-muted-foreground">
                                  {permission.description}
                                </p>
                              </div>
                            </div>
                          ))}
                        </div>
                        
                        {category !== availablePermissions[availablePermissions.length - 1] && (
                          <Separator className="mt-4" />
                        )}
                      </div>
                    );
                  })}
                </div>
              </CardContent>
            </Card>

            {/* High Risk Warning */}
            {selectedHighRiskPermissions.length > 0 && (
              <Alert className="border-orange-200 bg-orange-50">
                <AlertTriangle className="h-4 w-4 text-orange-600" />
                <AlertDescription className="text-orange-800">
                  <strong>Warning:</strong> This role includes {selectedHighRiskPermissions.length} high-risk permission(s). 
                  Users with this role will have significant system access. Please ensure this is intentional.
                </AlertDescription>
              </Alert>
            )}

            {/* Submit Error */}
            {errors.submit && (
              <Alert className="border-red-200 bg-red-50">
                <AlertTriangle className="h-4 w-4 text-red-600" />
                <AlertDescription className="text-red-600">
                  {errors.submit}
                </AlertDescription>
              </Alert>
            )}

            {/* Actions */}
            <div className="flex items-center justify-between pt-6 border-t">
              <div className="text-sm text-muted-foreground">
                {formData.permissions.length} permission(s) selected
              </div>
              <div className="flex items-center space-x-3">
                <Button 
                  type="button" 
                  variant="outline" 
                  onClick={() => navigate('/users/permissions')}
                  disabled={isSubmitting}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? (
                    <>
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                      Creating...
                    </>
                  ) : (
                    <>
                      <Save className="h-4 w-4 mr-2" />
                      Create Role
                    </>
                  )}
                </Button>
              </div>
            </div>
          </form>
    </div>
  );
}