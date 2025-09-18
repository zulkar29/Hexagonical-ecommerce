import React, { useState } from 'react';
import {
  Users,
  Shield,
  Plus,
  Search,
  Filter,
  MoreHorizontal,
  Edit,
  Trash2,
  Eye,
  UserPlus,
  Settings,
  Lock,
  Unlock,
  CheckCircle,
  XCircle,
  AlertTriangle
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
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
} from '@/components/ui/dialog';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';

// Mock data for roles and permissions
const mockRoles = [
  {
    id: 1,
    name: 'Super Admin',
    description: 'Full system access with all permissions',
    userCount: 2,
    permissions: ['all'],
    isSystem: true,
    createdAt: '2024-01-01',
    updatedAt: '2024-01-15'
  },
  {
    id: 2,
    name: 'Platform Admin',
    description: 'Manage tenants, subscriptions, and platform settings',
    userCount: 5,
    permissions: ['tenant_management', 'subscription_management', 'analytics', 'support'],
    isSystem: false,
    createdAt: '2024-01-05',
    updatedAt: '2024-01-18'
  },
  {
    id: 3,
    name: 'Support Manager',
    description: 'Handle customer support and basic tenant management',
    userCount: 8,
    permissions: ['support', 'tenant_view', 'user_management'],
    isSystem: false,
    createdAt: '2024-01-10',
    updatedAt: '2024-01-20'
  },
  {
    id: 4,
    name: 'Analytics Viewer',
    description: 'Read-only access to analytics and reports',
    userCount: 12,
    permissions: ['analytics_view', 'reports_view'],
    isSystem: false,
    createdAt: '2024-01-12',
    updatedAt: '2024-01-19'
  }
];

const mockUsers = [
  {
    id: 1,
    name: 'John Smith',
    email: 'john.smith@company.com',
    role: 'Super Admin',
    status: 'active',
    lastLogin: '2024-01-20T10:30:00Z',
    createdAt: '2024-01-01'
  },
  {
    id: 2,
    name: 'Sarah Johnson',
    email: 'sarah.johnson@company.com',
    role: 'Platform Admin',
    status: 'active',
    lastLogin: '2024-01-19T15:45:00Z',
    createdAt: '2024-01-05'
  },
  {
    id: 3,
    name: 'Mike Wilson',
    email: 'mike.wilson@company.com',
    role: 'Support Manager',
    status: 'active',
    lastLogin: '2024-01-20T09:15:00Z',
    createdAt: '2024-01-08'
  },
  {
    id: 4,
    name: 'Emily Davis',
    email: 'emily.davis@company.com',
    role: 'Analytics Viewer',
    status: 'inactive',
    lastLogin: '2024-01-15T14:20:00Z',
    createdAt: '2024-01-12'
  }
];

const availablePermissions = [
  {
    category: 'Tenant Management',
    permissions: [
      { key: 'tenant_management', name: 'Full Tenant Management', description: 'Create, edit, delete tenants' },
      { key: 'tenant_view', name: 'View Tenants', description: 'Read-only access to tenant information' },
      { key: 'tenant_suspend', name: 'Suspend Tenants', description: 'Suspend/reactivate tenant accounts' }
    ]
  },
  {
    category: 'Subscription Management',
    permissions: [
      { key: 'subscription_management', name: 'Manage Subscriptions', description: 'Handle billing and subscriptions' },
      { key: 'subscription_view', name: 'View Subscriptions', description: 'Read-only access to subscription data' }
    ]
  },
  {
    category: 'User Management',
    permissions: [
      { key: 'user_management', name: 'Manage Users', description: 'Create, edit, delete admin users' },
      { key: 'role_management', name: 'Manage Roles', description: 'Create and modify user roles' }
    ]
  },
  {
    category: 'Analytics & Reports',
    permissions: [
      { key: 'analytics', name: 'Full Analytics Access', description: 'Access all analytics and reports' },
      { key: 'analytics_view', name: 'View Analytics', description: 'Read-only access to analytics' },
      { key: 'reports_view', name: 'View Reports', description: 'Access to generated reports' }
    ]
  },
  {
    category: 'Support',
    permissions: [
      { key: 'support', name: 'Support Management', description: 'Handle customer support tickets' },
      { key: 'support_view', name: 'View Support', description: 'Read-only access to support tickets' }
    ]
  },
  {
    category: 'System',
    permissions: [
      { key: 'system_settings', name: 'System Settings', description: 'Modify platform settings' },
      { key: 'audit_logs', name: 'Audit Logs', description: 'Access system audit logs' }
    ]
  }
];

const getStatusColor = (status) => {
  switch (status) {
    case 'active': return 'bg-green-100 text-green-800';
    case 'inactive': return 'bg-gray-100 text-gray-800';
    case 'suspended': return 'bg-red-100 text-red-800';
    default: return 'bg-gray-100 text-gray-800';
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

export default function UserPermissions() {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedRole, setSelectedRole] = useState('all');
  const [showCreateRoleDialog, setShowCreateRoleDialog] = useState(false);
  const [showCreateUserDialog, setShowCreateUserDialog] = useState(false);
  const [showEditRoleDialog, setShowEditRoleDialog] = useState(false);
  const [editingRole, setEditingRole] = useState(null);
  const [newRole, setNewRole] = useState({
    name: '',
    description: '',
    permissions: []
  });
  const [newUser, setNewUser] = useState({
    name: '',
    email: '',
    role: ''
  });

  const filteredUsers = mockUsers.filter(user => {
    const matchesSearch = user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         user.email.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesRole = selectedRole === 'all' || user.role === selectedRole;
    return matchesSearch && matchesRole;
  });

  const handleCreateRole = () => {
    console.log('Creating role:', newRole);
    setShowCreateRoleDialog(false);
    setNewRole({ name: '', description: '', permissions: [] });
  };

  const handleCreateUser = () => {
    console.log('Creating user:', newUser);
    setShowCreateUserDialog(false);
    setNewUser({ name: '', email: '', role: '' });
  };

  const handleEditRole = (role) => {
    setEditingRole(role);
    setNewRole({
      name: role.name,
      description: role.description,
      permissions: role.permissions
    });
    setShowEditRoleDialog(true);
  };

  const handleUpdateRole = () => {
    console.log('Updating role:', editingRole.id, newRole);
    setShowEditRoleDialog(false);
    setEditingRole(null);
    setNewRole({ name: '', description: '', permissions: [] });
  };

  const togglePermission = (permissionKey) => {
    setNewRole(prev => ({
      ...prev,
      permissions: prev.permissions.includes(permissionKey)
        ? prev.permissions.filter(p => p !== permissionKey)
        : [...prev.permissions, permissionKey]
    }));
  };

  return (
    <div className="space-y-6 p-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-3xl font-bold tracking-tight">User Permissions</h2>
              <p className="text-muted-foreground">
                Manage admin users, roles, and permissions for the platform
              </p>
            </div>
            <div className="flex items-center space-x-2">
              <Button variant="outline" onClick={() => setShowCreateRoleDialog(true)}>
                <Shield className="h-4 w-4 mr-2" />
                Create Role
              </Button>
              <Button onClick={() => setShowCreateUserDialog(true)}>
                <UserPlus className="h-4 w-4 mr-2" />
                Add User
              </Button>
            </div>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Total Users</p>
                    <p className="text-2xl font-bold">{mockUsers.length}</p>
                  </div>
                  <Users className="h-8 w-8 text-blue-500" />
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Active Users</p>
                    <p className="text-2xl font-bold">{mockUsers.filter(u => u.status === 'active').length}</p>
                  </div>
                  <CheckCircle className="h-8 w-8 text-green-500" />
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Total Roles</p>
                    <p className="text-2xl font-bold">{mockRoles.length}</p>
                  </div>
                  <Shield className="h-8 w-8 text-purple-500" />
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Custom Roles</p>
                    <p className="text-2xl font-bold">{mockRoles.filter(r => !r.isSystem).length}</p>
                  </div>
                  <Settings className="h-8 w-8 text-orange-500" />
                </div>
              </CardContent>
            </Card>
          </div>

          <Tabs defaultValue="users" className="space-y-4">
            <TabsList>
              <TabsTrigger value="users">Users</TabsTrigger>
              <TabsTrigger value="roles">Roles</TabsTrigger>
              <TabsTrigger value="permissions">Permissions</TabsTrigger>
            </TabsList>

            {/* Users Tab */}
            <TabsContent value="users">
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle>Admin Users</CardTitle>
                      <CardDescription>Manage platform administrator accounts</CardDescription>
                    </div>
                    <div className="flex items-center space-x-2">
                      <div className="relative">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                        <Input
                          placeholder="Search users..."
                          value={searchQuery}
                          onChange={(e) => setSearchQuery(e.target.value)}
                          className="pl-10 w-64"
                        />
                      </div>
                      <Select value={selectedRole} onValueChange={setSelectedRole}>
                        <SelectTrigger className="w-40">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="all">All Roles</SelectItem>
                          {mockRoles.map(role => (
                            <SelectItem key={role.id} value={role.name}>{role.name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>User</TableHead>
                        <TableHead>Role</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Last Login</TableHead>
                        <TableHead>Created</TableHead>
                        <TableHead>Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {filteredUsers.map((user) => (
                        <TableRow key={user.id}>
                          <TableCell>
                            <div>
                              <p className="font-medium">{user.name}</p>
                              <p className="text-sm text-muted-foreground">{user.email}</p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <Badge variant="outline">{user.role}</Badge>
                          </TableCell>
                          <TableCell>
                            <Badge className={getStatusColor(user.status)}>
                              {user.status}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            <p className="text-sm">{formatDateTime(user.lastLogin)}</p>
                          </TableCell>
                          <TableCell>
                            <p className="text-sm">{formatDate(user.createdAt)}</p>
                          </TableCell>
                          <TableCell>
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="sm">
                                  <MoreHorizontal className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end">
                                <DropdownMenuItem>
                                  <Edit className="h-4 w-4 mr-2" />
                                  Edit User
                                </DropdownMenuItem>
                                <DropdownMenuItem>
                                  <Eye className="h-4 w-4 mr-2" />
                                  View Details
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                {user.status === 'active' ? (
                                  <DropdownMenuItem className="text-yellow-600">
                                    <Lock className="h-4 w-4 mr-2" />
                                    Suspend User
                                  </DropdownMenuItem>
                                ) : (
                                  <DropdownMenuItem className="text-green-600">
                                    <Unlock className="h-4 w-4 mr-2" />
                                    Activate User
                                  </DropdownMenuItem>
                                )}
                                <DropdownMenuItem className="text-red-600">
                                  <Trash2 className="h-4 w-4 mr-2" />
                                  Delete User
                                </DropdownMenuItem>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </CardContent>
              </Card>
            </TabsContent>

            {/* Roles Tab */}
            <TabsContent value="roles">
              <Card>
                <CardHeader>
                  <CardTitle>User Roles</CardTitle>
                  <CardDescription>Manage user roles and their permissions</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {mockRoles.map((role) => (
                      <Card key={role.id} className="relative">
                        <CardHeader>
                          <div className="flex items-center justify-between">
                            <div className="flex items-center space-x-2">
                              <Shield className="h-5 w-5 text-primary" />
                              <CardTitle className="text-lg">{role.name}</CardTitle>
                            </div>
                            {role.isSystem && (
                              <Badge variant="secondary" className="text-xs">
                                System
                              </Badge>
                            )}
                          </div>
                          <CardDescription>{role.description}</CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div className="space-y-3">
                            <div className="flex items-center justify-between text-sm">
                              <span className="text-muted-foreground">Users:</span>
                              <span className="font-medium">{role.userCount}</span>
                            </div>
                            <div className="flex items-center justify-between text-sm">
                              <span className="text-muted-foreground">Permissions:</span>
                              <span className="font-medium">{role.permissions.length}</span>
                            </div>
                            <div className="flex items-center justify-between text-sm">
                              <span className="text-muted-foreground">Updated:</span>
                              <span className="font-medium">{formatDate(role.updatedAt)}</span>
                            </div>
                            <Separator />
                            <div className="flex items-center justify-between">
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handleEditRole(role)}
                                disabled={role.isSystem}
                              >
                                <Edit className="h-4 w-4 mr-2" />
                                Edit
                              </Button>
                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button variant="ghost" size="sm">
                                    <MoreHorizontal className="h-4 w-4" />
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                  <DropdownMenuItem>
                                    <Eye className="h-4 w-4 mr-2" />
                                    View Details
                                  </DropdownMenuItem>
                                  <DropdownMenuItem>
                                    <Users className="h-4 w-4 mr-2" />
                                    View Users
                                  </DropdownMenuItem>
                                  {!role.isSystem && (
                                    <>
                                      <DropdownMenuSeparator />
                                      <DropdownMenuItem className="text-red-600">
                                        <Trash2 className="h-4 w-4 mr-2" />
                                        Delete Role
                                      </DropdownMenuItem>
                                    </>
                                  )}
                                </DropdownMenuContent>
                              </DropdownMenu>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            {/* Permissions Tab */}
            <TabsContent value="permissions">
              <Card>
                <CardHeader>
                  <CardTitle>Available Permissions</CardTitle>
                  <CardDescription>Overview of all available permissions in the system</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-6">
                    {availablePermissions.map((category) => (
                      <div key={category.category}>
                        <h3 className="text-lg font-semibold mb-3">{category.category}</h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                          {category.permissions.map((permission) => (
                            <Card key={permission.key} className="p-4">
                              <div className="space-y-2">
                                <div className="flex items-center space-x-2">
                                  <Shield className="h-4 w-4 text-primary" />
                                  <h4 className="font-medium">{permission.name}</h4>
                                </div>
                                <p className="text-sm text-muted-foreground">{permission.description}</p>
                                <Badge variant="outline" className="text-xs">
                                  {permission.key}
                                </Badge>
                              </div>
                            </Card>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>

          {/* Create Role Dialog */}
          <Dialog open={showCreateRoleDialog} onOpenChange={setShowCreateRoleDialog}>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>Create New Role</DialogTitle>
                <DialogDescription>
                  Define a new role with specific permissions
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="roleName">Role Name</Label>
                    <Input
                      id="roleName"
                      value={newRole.name}
                      onChange={(e) => setNewRole(prev => ({ ...prev, name: e.target.value }))}
                      placeholder="Enter role name"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="roleDescription">Description</Label>
                    <Input
                      id="roleDescription"
                      value={newRole.description}
                      onChange={(e) => setNewRole(prev => ({ ...prev, description: e.target.value }))}
                      placeholder="Enter role description"
                    />
                  </div>
                </div>
                <div className="space-y-3">
                  <Label>Permissions</Label>
                  <div className="space-y-4 max-h-60 overflow-y-auto">
                    {availablePermissions.map((category) => (
                      <div key={category.category}>
                        <h4 className="font-medium text-sm mb-2">{category.category}</h4>
                        <div className="space-y-2 pl-4">
                          {category.permissions.map((permission) => (
                            <div key={permission.key} className="flex items-center space-x-2">
                              <Checkbox
                                id={permission.key}
                                checked={newRole.permissions.includes(permission.key)}
                                onCheckedChange={() => togglePermission(permission.key)}
                              />
                              <Label htmlFor={permission.key} className="text-sm">
                                {permission.name}
                              </Label>
                            </div>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setShowCreateRoleDialog(false)}>
                  Cancel
                </Button>
                <Button onClick={handleCreateRole}>
                  Create Role
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>

          {/* Edit Role Dialog */}
          <Dialog open={showEditRoleDialog} onOpenChange={setShowEditRoleDialog}>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>Edit Role</DialogTitle>
                <DialogDescription>
                  Modify role permissions and details
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="editRoleName">Role Name</Label>
                    <Input
                      id="editRoleName"
                      value={newRole.name}
                      onChange={(e) => setNewRole(prev => ({ ...prev, name: e.target.value }))}
                      placeholder="Enter role name"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="editRoleDescription">Description</Label>
                    <Input
                      id="editRoleDescription"
                      value={newRole.description}
                      onChange={(e) => setNewRole(prev => ({ ...prev, description: e.target.value }))}
                      placeholder="Enter role description"
                    />
                  </div>
                </div>
                <div className="space-y-3">
                  <Label>Permissions</Label>
                  <div className="space-y-4 max-h-60 overflow-y-auto">
                    {availablePermissions.map((category) => (
                      <div key={category.category}>
                        <h4 className="font-medium text-sm mb-2">{category.category}</h4>
                        <div className="space-y-2 pl-4">
                          {category.permissions.map((permission) => (
                            <div key={permission.key} className="flex items-center space-x-2">
                              <Checkbox
                                id={`edit-${permission.key}`}
                                checked={newRole.permissions.includes(permission.key)}
                                onCheckedChange={() => togglePermission(permission.key)}
                              />
                              <Label htmlFor={`edit-${permission.key}`} className="text-sm">
                                {permission.name}
                              </Label>
                            </div>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setShowEditRoleDialog(false)}>
                  Cancel
                </Button>
                <Button onClick={handleUpdateRole}>
                  Update Role
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>

          {/* Create User Dialog */}
          <Dialog open={showCreateUserDialog} onOpenChange={setShowCreateUserDialog}>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New User</DialogTitle>
                <DialogDescription>
                  Create a new admin user account
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="userName">Full Name</Label>
                  <Input
                    id="userName"
                    value={newUser.name}
                    onChange={(e) => setNewUser(prev => ({ ...prev, name: e.target.value }))}
                    placeholder="Enter full name"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="userEmail">Email Address</Label>
                  <Input
                    id="userEmail"
                    type="email"
                    value={newUser.email}
                    onChange={(e) => setNewUser(prev => ({ ...prev, email: e.target.value }))}
                    placeholder="Enter email address"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="userRole">Role</Label>
                  <Select value={newUser.role} onValueChange={(value) => setNewUser(prev => ({ ...prev, role: value }))}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select a role" />
                    </SelectTrigger>
                    <SelectContent>
                      {mockRoles.map(role => (
                        <SelectItem key={role.id} value={role.name}>{role.name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setShowCreateUserDialog(false)}>
                  Cancel
                </Button>
                <Button onClick={handleCreateUser}>
                  Create User
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
    </div>
  );
}