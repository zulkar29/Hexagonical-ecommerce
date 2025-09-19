import React, { useState } from 'react';
import {
  Shield,
  CheckCircle,
  AlertTriangle,
  Plus,
  MoreHorizontal,
  Search,
  Edit,
  Trash2,
  Eye
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Link } from 'react-router-dom';

// Mock data for roles
const roles = [
  { 
    id: 1, 
    name: 'Admin', 
    description: 'Full access to all features', 
    status: 'Active',
    createdAt: '2024-01-01',
    updatedAt: '2024-01-15',
    usersCount: 2,
    isSystem: true
  },
  { 
    id: 2, 
    name: 'Manager', 
    description: 'Manage users and tenants', 
    status: 'Active',
    createdAt: '2024-01-05',
    updatedAt: '2024-01-18',
    usersCount: 5,
    isSystem: false
  },
  { 
    id: 3, 
    name: 'Support', 
    description: 'Handle support tickets', 
    status: 'Active',
    createdAt: '2024-01-10',
    updatedAt: '2024-01-19',
    usersCount: 8,
    isSystem: false
  },
  { 
    id: 4, 
    name: 'Viewer', 
    description: 'Read-only access', 
    status: 'Active',
    createdAt: '2024-01-12',
    updatedAt: '2024-01-19',
    usersCount: 10,
    isSystem: false
  },
  { 
    id: 5, 
    name: 'Suspended', 
    description: 'No access', 
    status: 'Inactive',
    createdAt: '2024-01-03',
    updatedAt: '2024-01-14',
    usersCount: 1,
    isSystem: false
  }
];

// Helper function to calculate role statistics
const calculateRoleStats = (roles) => {
  const totalRoles = roles.length;
  const activeRoles = roles.filter(role => role.status === 'Active').length;
  const inactiveRoles = roles.filter(role => role.status === 'Inactive').length;
  
  return [
    { 
      label: 'Total Roles', 
      value: totalRoles, 
      icon: Shield, 
      color: 'text-blue-600',
      change: '+1%'
    },
    { 
      label: 'Active Roles', 
      value: activeRoles, 
      icon: CheckCircle, 
      color: 'text-green-600',
      change: '0%'
    },
    { 
      label: 'Inactive Roles', 
      value: inactiveRoles, 
      icon: AlertTriangle, 
      color: 'text-amber-600',
      change: '+2%'
    }
  ];
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

export default function RoleManagement() {
  const [searchQuery, setSearchQuery] = useState('');
  
  // Calculate role statistics
  const roleStats = calculateRoleStats(roles);
  
  // Filter roles based on search query
  const filteredRoles = roles.filter(role => 
    role.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    role.description.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const getStatusBadge = (status) => {
    switch (status) {
      case 'Active':
        return <Badge variant="success">Active</Badge>;
      case 'Inactive':
        return <Badge variant="destructive">Inactive</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Shield className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Role Management</h1>
            <Badge variant="outline" className="ml-2">Permissions</Badge>
          </div>
          <div className="flex items-center space-x-2">
            <Button asChild>
              <Link to="/users/permissions/role/create">
                <Plus className="h-4 w-4 mr-2" />
                New Role
              </Link>
            </Button>
          </div>
        </div>
      </div>
      
      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Stats Row */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 w-full">
          {roleStats.map((stat) => (
            <div key={stat.label} className="flex items-center gap-3 bg-card p-3 rounded-lg border">
              <div className={`p-2 rounded-md ${stat.color} bg-opacity-10 shrink-0`}>
                <stat.icon className={`h-5 w-5 ${stat.color}`} />
              </div>
              <div className="flex-grow">
                <p className="text-sm text-muted-foreground">{stat.label}</p>
                <div className="flex items-center gap-1.5">
                  <p className="text-lg font-semibold">{stat.value}</p>
                  <span className={`text-xs ${stat.change.startsWith('+') ? 'text-green-600' : 'text-red-600'} font-medium`}>
                    {stat.change}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
        
        {/* Roles Table */}
        <Card>
          <CardHeader className="py-4">
            <div className="flex flex-wrap items-center gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  type="text"
                  placeholder="Search roles by name or description..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 w-full"
                />
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Role Name</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead>Users</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Last Modified</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredRoles.map((role) => (
                    <TableRow key={role.id}>
                      <TableCell>
                        <div className="flex items-center gap-3">
                          <Shield className="h-4 w-4 text-primary" />
                          <Link to={`/users/permissions/role/edit/${role.id}`} className="font-medium hover:underline">
                            {role.name}
                          </Link>
                        </div>
                      </TableCell>
                      <TableCell className="max-w-sm truncate">
                        {role.description}
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">{role.usersCount}</Badge>
                      </TableCell>
                      <TableCell>
                        {getStatusBadge(role.status)}
                      </TableCell>
                      <TableCell>{formatDate(role.updatedAt)}</TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button size="sm" variant="ghost">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem asChild>
                              <Link to={`/users/permissions/role/edit/${role.id}`}>
                                <Eye className="h-4 w-4 mr-2" />
                                View Role
                              </Link>
                            </DropdownMenuItem>
                            <DropdownMenuItem asChild disabled={role.isSystem}>
                              <Link to={`/users/permissions/role/edit/${role.id}`}>
                                <Edit className="h-4 w-4 mr-2" />
                                Edit Role
                              </Link>
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem disabled={role.isSystem} className="text-destructive">
                              <Trash2 className="h-4 w-4 mr-2" />
                              Delete Role
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
