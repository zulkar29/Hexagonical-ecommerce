import React from 'react';
import {
  Users,
  UserPlus,
  CheckCircle,
  AlertTriangle,
  Clock,
  MoreHorizontal,
  Mail,
  Calendar
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';

const summary = [
  { label: 'Total Users', value: 120, icon: Users, color: 'text-blue-600' },
  { label: 'Active', value: 98, icon: CheckCircle, color: 'text-green-600' },
  { label: 'Pending', value: 12, icon: Clock, color: 'text-yellow-600' },
  { label: 'Suspended', value: 10, icon: AlertTriangle, color: 'text-red-600' }
];

const users = [
  { id: 1, name: 'Abdul Rahman', email: 'abdul@shopowner.com', role: 'Admin', status: 'Active', joined: '2024-01-10' },
  { id: 2, name: 'Fatima Khatun', email: 'fatima@shopowner.com', role: 'Manager', status: 'Pending', joined: '2024-03-15' },
  { id: 3, name: 'Rashida Begum', email: 'rashida@shopowner.com', role: 'Support', status: 'Active', joined: '2024-02-20' },
  { id: 4, name: 'Nasir Ahmed', email: 'nasir@shopowner.com', role: 'Admin', status: 'Suspended', joined: '2023-12-05' },
  { id: 5, name: 'Dr. Ahmed Ali', email: 'ahmed@shopowner.com', role: 'Manager', status: 'Active', joined: '2024-04-01' }
];

export default function UsersHome() {
  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center gap-4 px-6">
          <Users className="h-5 w-5 text-primary" />
          <h1 className="text-xl font-semibold">Users</h1>
          <Badge variant="outline" className="ml-2">Admin Panel</Badge>
        </div>
      </div>
      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
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

        {/* Users Table */}
        <Card>
          <CardHeader>
            <CardTitle>Users List</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <table className="min-w-full text-sm border">
                <thead>
                  <tr className="bg-muted/50">
                    <th className="p-2 text-left font-medium">Name</th>
                    <th className="p-2 text-left font-medium">Email</th>
                    <th className="p-2 text-left font-medium">Role</th>
                    <th className="p-2 text-left font-medium">Status</th>
                    <th className="p-2 text-left font-medium">Joined</th>
                    <th className="p-2 text-left font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((u) => (
                    <tr key={u.id} className="border-t">
                      <td className="p-2">
                        <Link to={`/users/${u.id}`} className="text-primary underline hover:opacity-80">{u.name}</Link>
                      </td>
                      <td className="p-2 flex items-center gap-2">
                        <Mail className="h-4 w-4 text-muted-foreground" />
                        {u.email}
                      </td>
                      <td className="p-2">{u.role}</td>
                      <td className="p-2">
                        {u.status === 'Active' ? (
                          <Badge variant="success">Active</Badge>
                        ) : u.status === 'Pending' ? (
                          <Badge variant="secondary">Pending</Badge>
                        ) : (
                          <Badge variant="destructive">Suspended</Badge>
                        )}
                      </td>
                      <td className="p-2 flex items-center gap-2">
                        <Calendar className="h-4 w-4 text-muted-foreground" />
                        {u.joined}
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
