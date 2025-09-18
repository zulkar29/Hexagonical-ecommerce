import React, { useState } from 'react';
import {
  ArrowLeft,
  MessageSquare,
  User,
  Calendar,
  Clock,
  Tag,
  AlertTriangle,
  CheckCircle,
  Send,
  Paperclip,
  Phone,
  Mail,
  Edit,
  MoreHorizontal
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { Link, useParams } from 'react-router-dom';

const ticketData = {
  id: 'TKT-001',
  title: 'Unable to access billing dashboard',
  description: 'Customer is experiencing issues accessing the billing dashboard. Error message appears when trying to navigate to billing section.',
  customer: {
    name: 'John Smith',
    email: 'john.smith@techcorp.com',
    phone: '+1 (555) 123-4567',
    tenant: 'TechCorp Solutions',
    tenantId: 'tenant_123'
  },
  priority: 'High',
  status: 'Open',
  category: 'Billing',
  created: '2024-01-15T10:30:00Z',
  lastUpdate: '2024-01-15T14:45:00Z',
  assignee: 'Sarah Johnson',
  tags: ['billing', 'dashboard', 'access-issue']
};

const ticketHistory = [
  {
    id: 1,
    type: 'comment',
    author: 'John Smith',
    role: 'Customer',
    content: 'I\'m getting a 403 error when trying to access the billing dashboard. This started happening yesterday.',
    timestamp: '2024-01-15T10:30:00Z',
    attachments: []
  },
  {
    id: 2,
    type: 'status_change',
    author: 'System',
    role: 'System',
    content: 'Ticket assigned to Sarah Johnson',
    timestamp: '2024-01-15T11:00:00Z'
  },
  {
    id: 3,
    type: 'comment',
    author: 'Sarah Johnson',
    role: 'Support Agent',
    content: 'Hi John, I\'m looking into this issue. Can you please try clearing your browser cache and cookies, then attempt to access the billing dashboard again?',
    timestamp: '2024-01-15T11:15:00Z',
    attachments: []
  },
  {
    id: 4,
    type: 'comment',
    author: 'John Smith',
    role: 'Customer',
    content: 'I tried clearing the cache but the issue persists. I\'ve attached a screenshot of the error.',
    timestamp: '2024-01-15T14:45:00Z',
    attachments: ['error_screenshot.png']
  }
];

export default function SupportDetail() {
  const { ticketId } = useParams();
  const [newMessage, setNewMessage] = useState('');
  const [ticketStatus, setTicketStatus] = useState(ticketData.status);
  const [ticketPriority, setTicketPriority] = useState(ticketData.priority);
  const [assignee, setAssignee] = useState(ticketData.assignee);

  const handleSendMessage = () => {
    if (newMessage.trim()) {
      // Handle sending message
      console.log('Sending message:', newMessage);
      setNewMessage('');
    }
  };

  const getStatusBadge = (status) => {
    switch (status) {
      case 'Open':
        return <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">Open</Badge>;
      case 'In Progress':
        return <Badge variant="outline" className="bg-blue-50 text-blue-700 border-blue-200">In Progress</Badge>;
      case 'Resolved':
        return <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">Resolved</Badge>;
      case 'Closed':
        return <Badge variant="outline" className="bg-gray-50 text-gray-700 border-gray-200">Closed</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  const getPriorityBadge = (priority) => {
    switch (priority) {
      case 'High':
        return <Badge variant="destructive">High</Badge>;
      case 'Medium':
        return <Badge variant="outline" className="bg-orange-50 text-orange-700 border-orange-200">Medium</Badge>;
      case 'Low':
        return <Badge variant="secondary">Low</Badge>;
      default:
        return <Badge variant="outline">{priority}</Badge>;
    }
  };

  const formatTimestamp = (timestamp) => {
    return new Date(timestamp).toLocaleString();
  };

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Link to="/support" className="flex items-center gap-2 text-muted-foreground hover:text-foreground">
              <ArrowLeft className="h-4 w-4" />
              Back to Support
            </Link>
            <Separator orientation="vertical" className="h-6" />
            <MessageSquare className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">{ticketData.id}</h1>
            {getStatusBadge(ticketStatus)}
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm">
              <Edit className="h-4 w-4 mr-2" />
              Edit
            </Button>
            <Button variant="outline" size="sm">
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 p-6">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Ticket Details */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span>{ticketData.title}</span>
                  <div className="flex items-center gap-2">
                    {getPriorityBadge(ticketPriority)}
                    <Badge variant="outline">{ticketData.category}</Badge>
                  </div>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground mb-4">{ticketData.description}</p>
                <div className="flex flex-wrap gap-2">
                  {ticketData.tags.map((tag) => (
                    <Badge key={tag} variant="secondary" className="text-xs">
                      <Tag className="h-3 w-3 mr-1" />
                      {tag}
                    </Badge>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Conversation History */}
            <Card>
              <CardHeader>
                <CardTitle>Conversation</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-6">
                  {ticketHistory.map((item) => (
                    <div key={item.id} className="flex gap-4">
                      <div className="flex-shrink-0">
                        {item.type === 'comment' ? (
                          <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                            <User className="h-4 w-4 text-primary" />
                          </div>
                        ) : (
                          <div className="w-8 h-8 bg-muted rounded-full flex items-center justify-center">
                            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
                          </div>
                        )}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="font-medium text-sm">{item.author}</span>
                          <Badge variant="outline" className="text-xs">{item.role}</Badge>
                          <span className="text-xs text-muted-foreground">
                            {formatTimestamp(item.timestamp)}
                          </span>
                        </div>
                        <div className="bg-muted/50 rounded-lg p-3">
                          <p className="text-sm">{item.content}</p>
                          {item.attachments && item.attachments.length > 0 && (
                            <div className="mt-2 flex gap-2">
                              {item.attachments.map((attachment, index) => (
                                <div key={index} className="flex items-center gap-1 text-xs text-primary">
                                  <Paperclip className="h-3 w-3" />
                                  {attachment}
                                </div>
                              ))}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Reply Section */}
            <Card>
              <CardHeader>
                <CardTitle>Reply</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <Textarea
                    placeholder="Type your response..."
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    rows={4}
                  />
                  <div className="flex items-center justify-between">
                    <Button variant="outline" size="sm">
                      <Paperclip className="h-4 w-4 mr-2" />
                      Attach File
                    </Button>
                    <Button onClick={handleSendMessage} disabled={!newMessage.trim()}>
                      <Send className="h-4 w-4 mr-2" />
                      Send Reply
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Ticket Properties */}
            <Card>
              <CardHeader>
                <CardTitle>Ticket Properties</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <label className="text-sm font-medium mb-2 block">Status</label>
                  <Select value={ticketStatus} onValueChange={setTicketStatus}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="Open">Open</SelectItem>
                      <SelectItem value="In Progress">In Progress</SelectItem>
                      <SelectItem value="Resolved">Resolved</SelectItem>
                      <SelectItem value="Closed">Closed</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <label className="text-sm font-medium mb-2 block">Priority</label>
                  <Select value={ticketPriority} onValueChange={setTicketPriority}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="Low">Low</SelectItem>
                      <SelectItem value="Medium">Medium</SelectItem>
                      <SelectItem value="High">High</SelectItem>
                      <SelectItem value="Critical">Critical</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <label className="text-sm font-medium mb-2 block">Assignee</label>
                  <Select value={assignee} onValueChange={setAssignee}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="Sarah Johnson">Sarah Johnson</SelectItem>
                      <SelectItem value="David Wilson">David Wilson</SelectItem>
                      <SelectItem value="Mike Chen">Mike Chen</SelectItem>
                      <SelectItem value="Unassigned">Unassigned</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardContent>
            </Card>

            {/* Customer Information */}
            <Card>
              <CardHeader>
                <CardTitle>Customer Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center gap-3">
                  <User className="h-4 w-4 text-muted-foreground" />
                  <div>
                    <p className="font-medium">{ticketData.customer.name}</p>
                    <p className="text-sm text-muted-foreground">{ticketData.customer.tenant}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <a href={`mailto:${ticketData.customer.email}`} className="text-sm text-primary hover:underline">
                    {ticketData.customer.email}
                  </a>
                </div>
                <div className="flex items-center gap-3">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <a href={`tel:${ticketData.customer.phone}`} className="text-sm text-primary hover:underline">
                    {ticketData.customer.phone}
                  </a>
                </div>
                <Separator />
                <div className="space-y-2">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">Created:</span>
                    <span>{formatTimestamp(ticketData.created)}</span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">Last Update:</span>
                    <span>{formatTimestamp(ticketData.lastUpdate)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Quick Actions */}
            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <Button variant="outline" className="w-full justify-start">
                  <CheckCircle className="h-4 w-4 mr-2" />
                  Mark as Resolved
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  <User className="h-4 w-4 mr-2" />
                  Escalate to Manager
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  <MessageSquare className="h-4 w-4 mr-2" />
                  Create Internal Note
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}