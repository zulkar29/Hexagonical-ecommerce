import React from 'react';
import {
  Shield,
  CheckCircle,
  AlertTriangle,
  Plus,
  MoreHorizontal
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';

const summary = [
  { label: 'Total Roles', value: 5, icon: Shield, color: 'text-blue-600' },
  { label: 'Active', value: 4, icon: CheckCircle, color: 'text-green-600' },
  { label: 'Inactive', value: 1, icon: AlertTriangle, color: 'text-red-600' }
];

const roles = [
  { id: 1, name: 'Admin', description: 'Full access to all features', status: 'Active' },
  { id: 2, name: 'Manager', description: 'Manage users and tenants', status: 'Active' },
  { id: 3, name: 'Support', description: 'Handle support tickets', status: 'Active' },
  { id: 4, name: 'Viewer', description: 'Read-only access', status: 'Active' },
  { id: 5, name: 'Suspended', description: 'No access', status: 'Inactive' }
];

export default function RoleManagement() {
  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center gap-4 px-6">
          <Shield className="h-5 w-5 text-primary" />
          <h1 className="text-xl font-semibold">Role Management</h1>
          <Badge variant="outline" className="ml-2">Permissions</Badge>
          <div className="ml-auto flex gap-2">
            <Button asChild size="sm">
              <Link to="/users/permissions/role/create">
                <Plus className="h-4 w-4 mr-2" />
                New Role
              </Link>
            </Button>
          </div>
        </div>
      </div>
      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {summary.map((item) => (
            <Card key={item.label}>
              <CardContent className="p-6 flex items-center gap-4">
                <item.icon className={`h-8 w-8 ${item.color}`} />
                <div>
                  <p className="text-sm text-muted-foreground">{item.label}</p>
                  <p className="text-2xl font-bold">{item.value}</p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Roles Table */}
        <Card>
          <CardHeader>
            <CardTitle>Roles List</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <table className="min-w-full text-sm border">
                <thead>
                  <tr className="bg-muted/50">
                    <th className="p-2 text-left font-medium">Role Name</th>
                    <th className="p-2 text-left font-medium">Description</th>
                    <th className="p-2 text-left font-medium">Status</th>
                    <th className="p-2 text-left font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {roles.map((r) => (
                    <tr key={r.id} className="border-t">
                      <td className="p-2">
                        <Link to={`/users/permissions/role/edit/${r.id}`} className="text-primary underline hover:opacity-80">{r.name}</Link>
                      </td>
                      <td className="p-2">{r.description}</td>
                      <td className="p-2">
                        {r.status === 'Active' ? (
                          <Badge variant="success">Active</Badge>
                        ) : (
                          <Badge variant="destructive">Inactive</Badge>
                        )}
                      </td>
                      <td className="p-2">
                        <Button size="icon" variant="outline">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
