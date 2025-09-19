import React, { useState } from 'react';
import {
  MessageCircle,
  Users,
  AlertCircle,
  CheckCircle2,
  Clock,
  Search,
  Filter,
  Plus,
  MoreHorizontal,
  UserRound,
  Eye,
  ArchiveIcon,
  BellOff,
  Tag,
  ArrowUpRight
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Link } from 'react-router-dom';

// Mock data for support tickets
const mockTickets = [
  {
    id: 'TKT-1001',
    subject: 'Cannot access dashboard after upgrade',
    customer: 'John Smith',
    customerEmail: 'john.smith@acme.com',
    status: 'Open',
    priority: 'High',
    category: 'Technical',
    assignedTo: 'Sarah Johnson',
    createdAt: '2024-09-15T14:30:00Z',
    updatedAt: '2024-09-16T09:15:00Z',
    responseTime: '45m',
    tenant: 'Acme Corp'
  },
  {
    id: 'TKT-1002',
    subject: 'Billing issue with last month\'s invoice',
    customer: 'Emily Rogers',
    customerEmail: 'emily@globex.net',
    status: 'In Progress',
    priority: 'Medium',
    category: 'Billing',
    assignedTo: 'Mike Wilson',
    createdAt: '2024-09-14T11:20:00Z',
    updatedAt: '2024-09-16T10:45:00Z',
    responseTime: '2h 15m',
    tenant: 'Globex Industries'
  },
  {
    id: 'TKT-1003',
    subject: 'How to configure custom notifications',
    customer: 'David Lee',
    customerEmail: 'david.lee@startechco.com',
    status: 'Resolved',
    priority: 'Low',
    category: 'Help',
    assignedTo: 'Sarah Johnson',
    createdAt: '2024-09-10T09:00:00Z',
    updatedAt: '2024-09-12T16:30:00Z',
    responseTime: '1h 45m',
    tenant: 'StarTech Co'
  },
  {
    id: 'TKT-1004',
    subject: 'Need to add additional users to our account',
    customer: 'Lisa Wang',
    customerEmail: 'lwang@abctech.com',
    status: 'Open',
    priority: 'Medium',
    category: 'Account',
    assignedTo: 'Unassigned',
    createdAt: '2024-09-17T08:45:00Z',
    updatedAt: '2024-09-17T08:45:00Z',
    responseTime: '-',
    tenant: 'ABC Tech'
  },
  {
    id: 'TKT-1005',
    subject: 'API integration not working after latest update',
    customer: 'Robert Chen',
    customerEmail: 'robert@megacorp.com',
    status: 'In Progress',
    priority: 'Critical',
    category: 'Technical',
    assignedTo: 'Mike Wilson',
    createdAt: '2024-09-16T16:10:00Z',
    updatedAt: '2024-09-17T09:30:00Z',
    responseTime: '30m',
    tenant: 'MegaCorp'
  },
  {
    id: 'TKT-1006',
    subject: 'Request for extended trial period',
    customer: 'Sophie Martin',
    customerEmail: 'sophie@newstartup.co',
    status: 'Waiting on Customer',
    priority: 'Low',
    category: 'Sales',
    assignedTo: 'Alex Peterson',
    createdAt: '2024-09-15T11:45:00Z',
    updatedAt: '2024-09-16T14:20:00Z',
    responseTime: '1h 10m',
    tenant: 'NewStartup Inc.'
  },
  {
    id: 'TKT-1007',
    subject: 'Security concern regarding user permissions',
    customer: 'Thomas Brown',
    customerEmail: 'tbrown@securedata.org',
    status: 'Open',
    priority: 'High',
    category: 'Security',
    assignedTo: 'Sarah Johnson',
    createdAt: '2024-09-17T07:20:00Z',
    updatedAt: '2024-09-17T09:15:00Z',
    responseTime: '1h 55m',
    tenant: 'SecureData Solutions'
  }
];

// Helper function to calculate support statistics
const calculateSupportStats = (tickets) => {
  const totalTickets = tickets.length;
  const openTickets = tickets.filter(ticket => ticket.status === 'Open').length;
  const inProgressTickets = tickets.filter(ticket => ticket.status === 'In Progress').length;
  const resolvedTickets = tickets.filter(ticket => ticket.status === 'Resolved').length;
  
  return [
    { 
      label: 'Total Tickets', 
      value: totalTickets, 
      icon: MessageCircle, 
      color: 'text-blue-600',
      change: '+12%' 
    },
    { 
      label: 'Open', 
      value: openTickets, 
      icon: AlertCircle, 
      color: 'text-amber-600',
      change: '+5%' 
    },
    { 
      label: 'In Progress', 
      value: inProgressTickets, 
      icon: Clock, 
      color: 'text-indigo-600',
      change: '-2%' 
    },
    { 
      label: 'Resolved', 
      value: resolvedTickets, 
      icon: CheckCircle2, 
      color: 'text-green-600',
      change: '+8%' 
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

const formatDateTime = (dateString) => {
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

export default function CustomerSupport() {
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [priorityFilter, setPriorityFilter] = useState('all');
  const [categoryFilter, setCategoryFilter] = useState('all');
  
  // Calculate support statistics
  const supportStats = calculateSupportStats(mockTickets);
  
  // Filter tickets based on search query and filters
  const filteredTickets = mockTickets.filter(ticket => {
    const matchesSearch = 
      ticket.subject.toLowerCase().includes(searchQuery.toLowerCase()) ||
      ticket.customer.toLowerCase().includes(searchQuery.toLowerCase()) ||
      ticket.customerEmail.toLowerCase().includes(searchQuery.toLowerCase()) ||
      ticket.id.toLowerCase().includes(searchQuery.toLowerCase());
      
    const matchesStatus = statusFilter === 'all' || ticket.status === statusFilter;
    const matchesPriority = priorityFilter === 'all' || ticket.priority === priorityFilter;
    const matchesCategory = categoryFilter === 'all' || ticket.category === categoryFilter;
    
    return matchesSearch && matchesStatus && matchesPriority && matchesCategory;
  });

  const getStatusBadge = (status) => {
    switch (status) {
      case 'Open':
        return <Badge variant="outline" className="bg-amber-50 text-amber-700 border-amber-200">Open</Badge>;
      case 'In Progress':
        return <Badge variant="outline" className="bg-blue-50 text-blue-700 border-blue-200">In Progress</Badge>;
      case 'Resolved':
        return <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">Resolved</Badge>;
      case 'Waiting on Customer':
        return <Badge variant="outline" className="bg-purple-50 text-purple-700 border-purple-200">Waiting on Customer</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  const getPriorityBadge = (priority) => {
    switch (priority) {
      case 'Critical':
        return <Badge variant="destructive">Critical</Badge>;
      case 'High':
        return <Badge className="bg-orange-500">High</Badge>;
      case 'Medium':
        return <Badge variant="secondary">Medium</Badge>;
      case 'Low':
        return <Badge variant="outline">Low</Badge>;
      default:
        return <Badge variant="outline">{priority}</Badge>;
    }
  };

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <MessageCircle className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Customer Support</h1>
            <Badge variant="outline" className="ml-2">Tickets</Badge>
          </div>
          <div className="flex items-center space-x-2">
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              New Ticket
            </Button>
          </div>
        </div>
      </div>
      
      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Stats Row */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 w-full">
          {supportStats.map((stat) => (
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
        
        {/* Tickets Table */}
        <Card>
          <CardHeader className="py-4">
            <div className="flex flex-wrap items-center gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  type="text"
                  placeholder="Search tickets by ID, subject, customer..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 w-full"
                />
              </div>
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-[130px]">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="Open">Open</SelectItem>
                  <SelectItem value="In Progress">In Progress</SelectItem>
                  <SelectItem value="Resolved">Resolved</SelectItem>
                  <SelectItem value="Waiting on Customer">Waiting</SelectItem>
                </SelectContent>
              </Select>
              <Select value={priorityFilter} onValueChange={setPriorityFilter}>
                <SelectTrigger className="w-[130px]">
                  <SelectValue placeholder="Priority" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Priorities</SelectItem>
                  <SelectItem value="Critical">Critical</SelectItem>
                  <SelectItem value="High">High</SelectItem>
                  <SelectItem value="Medium">Medium</SelectItem>
                  <SelectItem value="Low">Low</SelectItem>
                </SelectContent>
              </Select>
              <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                <SelectTrigger className="w-[130px]">
                  <SelectValue placeholder="Category" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Categories</SelectItem>
                  <SelectItem value="Technical">Technical</SelectItem>
                  <SelectItem value="Billing">Billing</SelectItem>
                  <SelectItem value="Account">Account</SelectItem>
                  <SelectItem value="Security">Security</SelectItem>
                  <SelectItem value="Sales">Sales</SelectItem>
                  <SelectItem value="Help">Help</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardHeader>
          <CardContent>
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Ticket ID</TableHead>
                    <TableHead>Subject</TableHead>
                    <TableHead>Customer</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Priority</TableHead>
                    <TableHead>Category</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredTickets.map((ticket) => (
                    <TableRow key={ticket.id}>
                      <TableCell className="font-mono text-xs">
                        {ticket.id}
                      </TableCell>
                      <TableCell>
                        <div className="font-medium max-w-[200px] truncate">
                          {ticket.subject}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {ticket.tenant}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <div className="w-7 h-7 bg-primary/10 rounded-full flex items-center justify-center">
                            <span className="text-xs font-medium">{ticket.customer.charAt(0)}</span>
                          </div>
                          <div>
                            <div className="text-sm">{ticket.customer}</div>
                            <div className="text-xs text-muted-foreground">{ticket.customerEmail}</div>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        {getStatusBadge(ticket.status)}
                      </TableCell>
                      <TableCell>
                        {getPriorityBadge(ticket.priority)}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1">
                          <Tag className="h-3 w-3" />
                          <span className="text-sm">{ticket.category}</span>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="text-sm">{formatDate(ticket.createdAt)}</div>
                        <div className="text-xs text-muted-foreground">
                          {ticket.assignedTo !== 'Unassigned' ? `Assigned to ${ticket.assignedTo}` : 'Unassigned'}
                        </div>
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button size="sm" variant="ghost">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem asChild>
                              <Link to={`/support/tickets/${ticket.id}`}>
                                <Eye className="h-4 w-4 mr-2" />
                                View Ticket
                              </Link>
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                              <ArrowUpRight className="h-4 w-4 mr-2" />
                              Escalate
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem>
                              <ArchiveIcon className="h-4 w-4 mr-2" />
                              Archive
                            </DropdownMenuItem>
                            <DropdownMenuItem className="text-amber-600">
                              <BellOff className="h-4 w-4 mr-2" />
                              Mute Notifications
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
