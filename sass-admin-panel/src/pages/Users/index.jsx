import React from 'react';
import { Link } from 'react-router-dom';
import { Users, UserPlus, UserCheck, UserX, Clock, MoreHorizontal, Search, Filter, CheckCircle, AlertTriangle } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

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
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Users className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Users</h1>
            <Badge variant="outline" className="ml-2">Admin Panel</Badge>
          </div>
          <Link to="/users/create">
            <Button>
              <UserPlus className="h-4 w-4 mr-2" />
              Create User
            </Button>
          </Link>
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
            <CardTitle>Users</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Role</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Joined</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {users.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                          <span className="text-sm font-medium">{user.name.charAt(0)}</span>
                        </div>
                        <span className="font-medium">{user.name}</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground">{user.email}</TableCell>
                    <TableCell>
                      <Badge variant="outline">{user.role}</Badge>
                    </TableCell>
                    <TableCell>
                      <Badge 
                        variant={user.status === 'Active' ? 'default' : 
                               user.status === 'Pending' ? 'secondary' : 'destructive'}
                      >
                        {user.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground">{user.joined}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Link to={`/users/${user.id}/edit`}>
                          <Button size="sm" variant="outline">
                            Edit
                          </Button>
                        </Link>
                        <Button size="sm" variant="ghost">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}