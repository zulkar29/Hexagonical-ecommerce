import React, { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  MessageCircle,
  ChevronLeft,
  Clock,
  Send,
  Paperclip,
  Tag,
  User,
  Building,
  CalendarClock,
  AlertCircle,
  ArrowUpRight,
  CheckCircle,
  MessageSquare,
  PencilLine,
  Info,
  MoreHorizontal,
  XCircle,
  ArchiveIcon,
  Share2,
  BellOff
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent, CardFooter, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

// Mock data for ticket details
const mockTicket = {
  id: 'TKT-1001',
  subject: 'Cannot access dashboard after upgrade',
  description: 'After upgrading to the latest version yesterday, I\'m unable to access the admin dashboard. I keep getting a "Loading..." message but it never completes loading. I\'ve tried clearing the browser cache and using different browsers, but the issue persists.',
  customer: 'John Smith',
  customerEmail: 'john.smith@acme.com',
  status: 'Open',
  priority: 'High',
  category: 'Technical',
  assignedTo: 'Sarah Johnson',
  createdAt: '2024-09-15T14:30:00Z',
  updatedAt: '2024-09-16T09:15:00Z',
  responseTime: '45m',
  tenant: 'Acme Corp',
  messages: [
    {
      id: 'm1',
      author: 'John Smith',
      authorEmail: 'john.smith@acme.com',
      isCustomer: true,
      content: 'After upgrading to the latest version yesterday, I\'m unable to access the admin dashboard. I keep getting a "Loading..." message but it never completes loading. I\'ve tried clearing the browser cache and using different browsers, but the issue persists.',
      timestamp: '2024-09-15T14:30:00Z',
      attachments: []
    },
    {
      id: 'm2',
      author: 'Sarah Johnson',
      authorEmail: 'sarah.johnson@support.com',
      isCustomer: false,
      content: 'Hi John, Im sorry to hear youre having trouble accessing the dashboard. Could you please provide the following information so we can better troubleshoot this issue:\n\n1. Your browser version\n2. The exact version you upgraded from and to\n3. Any error messages that appear in the browser console (you can access this with F12)\n\nIn the meantime, try accessing the dashboard using incognito/private browsing mode.',
      timestamp: '2024-09-15T15:15:00Z',
      attachments: []
    },
    {
      id: 'm3',
      author: 'John Smith',
      authorEmail: 'john.smith@acme.com',
      isCustomer: true,
      content: 'Thanks for the quick response, Sarah. Here are the details:\n\n1. Chrome v118.0.5993.88\n2. Upgraded from v2.4.1 to v3.0.0\n3. Console shows: "TypeError: Cannot read property \'data\' of undefined"\n\nI tried incognito mode but have the same issue.',
      timestamp: '2024-09-15T16:05:00Z',
      attachments: [
        { name: 'error_screenshot.png', size: '245KB' }
      ]
    },
    {
      id: 'm4',
      author: 'Sarah Johnson',
      authorEmail: 'sarah.johnson@support.com',
      isCustomer: false,
      content: 'Thank you for the information, John. The error suggests there might be an issue with how the new version is handling data. I\'ve opened an urgent ticket with our development team to investigate this further.\\n\\nAs a temporary workaround, you can try using our backup dashboard at https://backup-dashboard.acmecorp.com\\n\\nI\'ll update you as soon as I hear back from the development team.',
      timestamp: '2024-09-16T09:15:00Z',
      attachments: [
        { name: 'backup_access_guide.pdf', size: '1.2MB' }
      ]
    }
  ],
  notes: [
    {
      id: 'n1',
      author: 'Sarah Johnson',
      content: 'This appears to be related to the data parsing issue identified in the recent release. Engineering team has been notified.',
      timestamp: '2024-09-15T15:30:00Z'
    },
    {
      id: 'n2',
      author: 'Mike Wilson',
      content: 'Development confirmed this is a bug in v3.0.0. Fix will be included in v3.0.1 scheduled for tomorrow.',
      timestamp: '2024-09-16T10:20:00Z'
    }
  ],
  timeline: [
    { event: 'Ticket created', timestamp: '2024-09-15T14:30:00Z' },
    { event: 'Assigned to Sarah Johnson', timestamp: '2024-09-15T14:35:00Z' },
    { event: 'First response sent', timestamp: '2024-09-15T15:15:00Z' },
    { event: 'Internal note added', timestamp: '2024-09-15T15:30:00Z' },
    { event: 'Customer replied', timestamp: '2024-09-15T16:05:00Z' },
    { event: 'Escalated to Development', timestamp: '2024-09-15T16:45:00Z' },
    { event: 'Internal note added', timestamp: '2024-09-16T10:20:00Z' },
    { event: 'Agent replied', timestamp: '2024-09-16T09:15:00Z' }
  ],
  relatedTickets: [
    { id: 'TKT-985', subject: 'Dashboard loading issue on Firefox', status: 'Resolved' },
    { id: 'TKT-1003', subject: 'Cannot view analytics after v3.0.0 upgrade', status: 'In Progress' }
  ]
};

// Format date helper
const formatDateTime = (dateString) => {
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

// Time elapsed helper
const getTimeElapsed = (dateString) => {
  const now = new Date();
  const date = new Date(dateString);
  const diffMs = now - date;
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHrs = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays > 0) {
    return diffDays === 1 ? '1 day ago' : `${diffDays} days ago`;
  } else if (diffHrs > 0) {
    return diffHrs === 1 ? '1 hour ago' : `${diffHrs} hours ago`;
  } else {
    return diffMins <= 1 ? 'Just now' : `${diffMins} minutes ago`;
  }
};

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

export default function TicketDetail() {
  const { ticketId } = useParams();
  const [replyContent, setReplyContent] = useState('');
  const [ticket, setTicket] = useState(mockTicket);
  const [statusValue, setStatusValue] = useState(ticket.status);
  const [priorityValue, setPriorityValue] = useState(ticket.priority);
  const [assigneeValue, setAssigneeValue] = useState(ticket.assignedTo);
  const [noteContent, setNoteContent] = useState('');
  
  // In a real app, fetch the ticket details based on ticketId
  // useEffect(() => { fetchTicketDetails(ticketId) }, [ticketId]);
  
  const handleStatusChange = (value) => {
    setStatusValue(value);
    // In a real app, update the ticket status in the database
  };
  
  const handlePriorityChange = (value) => {
    setPriorityValue(value);
    // In a real app, update the ticket priority in the database
  };
  
  const handleAssigneeChange = (value) => {
    setAssigneeValue(value);
    // In a real app, update the ticket assignee in the database
  };
  
  const handleReply = () => {
    if (!replyContent.trim()) return;
    
    // In a real app, send the reply to the API
    console.log('Sending reply:', replyContent);
    
    // For demo: Add the reply to the messages
    const newMessage = {
      id: `m${ticket.messages.length + 1}`,
      author: 'Sarah Johnson',
      authorEmail: 'sarah.johnson@support.com',
      isCustomer: false,
      content: replyContent,
      timestamp: new Date().toISOString(),
      attachments: []
    };
    
    setTicket({
      ...ticket,
      messages: [...ticket.messages, newMessage],
      status: 'Waiting on Customer',
      updatedAt: new Date().toISOString()
    });
    
    setReplyContent('');
  };
  
  const handleAddNote = () => {
    if (!noteContent.trim()) return;
    
    // In a real app, send the note to the API
    console.log('Adding note:', noteContent);
    
    // For demo: Add the note
    const newNote = {
      id: `n${ticket.notes.length + 1}`,
      author: 'Sarah Johnson',
      content: noteContent,
      timestamp: new Date().toISOString()
    };
    
    setTicket({
      ...ticket,
      notes: [...ticket.notes, newNote]
    });
    
    setNoteContent('');
  };

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="sm" asChild>
              <Link to="/support">
                <ChevronLeft className="h-4 w-4 mr-1" />
                Back
              </Link>
            </Button>
            <div className="flex flex-col">
              <div className="flex items-center gap-2">
                <MessageCircle className="h-5 w-5 text-primary" />
                <h1 className="text-xl font-semibold">{ticket.id}</h1>
                {getStatusBadge(statusValue)}
                {getPriorityBadge(priorityValue)}
              </div>
              <p className="text-sm text-muted-foreground">{ticket.subject}</p>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  <MoreHorizontal className="h-4 w-4 mr-2" />
                  Actions
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem>
                  <ArrowUpRight className="h-4 w-4 mr-2" />
                  Escalate
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <XCircle className="h-4 w-4 mr-2" />
                  Close Ticket
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Share2 className="h-4 w-4 mr-2" />
                  Share
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
          </div>
        </div>
      </div>
      
      <div className="flex-1 overflow-auto p-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main ticket content */}
          <div className="lg:col-span-2 space-y-6">
            <Tabs defaultValue="conversation" className="w-full">
              <TabsList>
                <TabsTrigger value="conversation">Conversation</TabsTrigger>
                <TabsTrigger value="notes">Internal Notes</TabsTrigger>
                <TabsTrigger value="timeline">Timeline</TabsTrigger>
              </TabsList>
              
              {/* Conversation Tab */}
              <TabsContent value="conversation" className="space-y-4 pt-4">
                {ticket.messages.map((message) => (
                  <div key={message.id} className={`flex gap-4 ${message.isCustomer ? 'flex-row' : 'flex-row-reverse'}`}>
                    <Avatar className={`h-10 w-10 ${message.isCustomer ? 'bg-blue-100' : 'bg-emerald-100'}`}>
                      <AvatarFallback>{message.author.charAt(0)}</AvatarFallback>
                    </Avatar>
                    <div className={`flex-1 space-y-2 ${message.isCustomer ? '' : 'items-end'}`}>
                      <Card className={`${message.isCustomer ? 'bg-muted' : 'bg-primary/5'}`}>
                        <CardHeader className="py-3">
                          <div className="flex items-center justify-between">
                            <div>
                              <p className="font-medium">{message.author}</p>
                              <p className="text-xs text-muted-foreground">{message.authorEmail}</p>
                            </div>
                            <p className="text-xs text-muted-foreground">
                              {formatDateTime(message.timestamp)}
                            </p>
                          </div>
                        </CardHeader>
                        <CardContent className="pb-3">
                          <div className="whitespace-pre-line">{message.content}</div>
                          {message.attachments.length > 0 && (
                            <div className="mt-3 pt-3 border-t">
                              <p className="text-xs text-muted-foreground mb-2">Attachments:</p>
                              {message.attachments.map((attachment, idx) => (
                                <div key={idx} className="flex items-center gap-2 bg-background rounded-md p-2">
                                  <Paperclip className="h-3 w-3 text-muted-foreground" />
                                  <span className="text-sm">{attachment.name}</span>
                                  <span className="text-xs text-muted-foreground">({attachment.size})</span>
                                </div>
                              ))}
                            </div>
                          )}
                        </CardContent>
                      </Card>
                    </div>
                  </div>
                ))}
                
                {/* Reply box */}
                <Card>
                  <CardHeader className="py-3">
                    <CardTitle className="text-sm">Reply to {ticket.customer}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <Textarea 
                      placeholder="Type your reply here..." 
                      className="min-h-[150px]"
                      value={replyContent}
                      onChange={(e) => setReplyContent(e.target.value)}
                    />
                  </CardContent>
                  <CardFooter className="justify-between">
                    <div className="flex items-center gap-2">
                      <Button variant="outline" size="sm">
                        <Paperclip className="h-4 w-4 mr-2" />
                        Attach
                      </Button>
                    </div>
                    <Button onClick={handleReply} disabled={!replyContent.trim()}>
                      <Send className="h-4 w-4 mr-2" />
                      Send Reply
                    </Button>
                  </CardFooter>
                </Card>
              </TabsContent>
              
              {/* Internal Notes Tab */}
              <TabsContent value="notes" className="space-y-4 pt-4">
                {ticket.notes.map((note) => (
                  <Card key={note.id}>
                    <CardHeader className="py-3">
                      <div className="flex items-center justify-between">
                        <p className="font-medium">{note.author}</p>
                        <p className="text-xs text-muted-foreground">
                          {formatDateTime(note.timestamp)}
                        </p>
                      </div>
                    </CardHeader>
                    <CardContent className="pb-3">
                      {note.content}
                    </CardContent>
                  </Card>
                ))}
                
                {/* Add note */}
                <Card>
                  <CardHeader className="py-3">
                    <CardTitle className="text-sm">Add Internal Note</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <Textarea 
                      placeholder="Add a private note that is only visible to team members..." 
                      className="min-h-[100px]"
                      value={noteContent}
                      onChange={(e) => setNoteContent(e.target.value)}
                    />
                  </CardContent>
                  <CardFooter className="justify-end">
                    <Button onClick={handleAddNote} disabled={!noteContent.trim()}>
                      <PencilLine className="h-4 w-4 mr-2" />
                      Add Note
                    </Button>
                  </CardFooter>
                </Card>
              </TabsContent>
              
              {/* Timeline Tab */}
              <TabsContent value="timeline" className="pt-4">
                <Card>
                  <CardContent className="py-4">
                    <div className="space-y-4">
                      {ticket.timeline.map((item, idx) => (
                        <div key={idx} className="flex items-start gap-3">
                          <div className="mt-1 h-2 w-2 rounded-full bg-primary"></div>
                          <div>
                            <p className="font-medium">{item.event}</p>
                            <p className="text-sm text-muted-foreground">
                              {formatDateTime(item.timestamp)}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
          
          {/* Sidebar */}
          <div className="space-y-6">
            {/* Ticket Details */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Ticket Details</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <p className="text-sm text-muted-foreground">Status</p>
                    <Select value={statusValue} onValueChange={handleStatusChange}>
                      <SelectTrigger className="w-[140px]">
                        <SelectValue placeholder="Status" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="Open">Open</SelectItem>
                        <SelectItem value="In Progress">In Progress</SelectItem>
                        <SelectItem value="Waiting on Customer">Waiting on Customer</SelectItem>
                        <SelectItem value="Resolved">Resolved</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex justify-between">
                    <p className="text-sm text-muted-foreground">Priority</p>
                    <Select value={priorityValue} onValueChange={handlePriorityChange}>
                      <SelectTrigger className="w-[140px]">
                        <SelectValue placeholder="Priority" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="Critical">Critical</SelectItem>
                        <SelectItem value="High">High</SelectItem>
                        <SelectItem value="Medium">Medium</SelectItem>
                        <SelectItem value="Low">Low</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex justify-between">
                    <p className="text-sm text-muted-foreground">Assignee</p>
                    <Select value={assigneeValue} onValueChange={handleAssigneeChange}>
                      <SelectTrigger className="w-[140px]">
                        <SelectValue placeholder="Assignee" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="Unassigned">Unassigned</SelectItem>
                        <SelectItem value="Sarah Johnson">Sarah Johnson</SelectItem>
                        <SelectItem value="Mike Wilson">Mike Wilson</SelectItem>
                        <SelectItem value="Alex Peterson">Alex Peterson</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <Separator />
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <User className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm text-muted-foreground">Customer</p>
                      <p className="font-medium">{ticket.customer}</p>
                      <p className="text-xs text-muted-foreground">{ticket.customerEmail}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Building className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm text-muted-foreground">Company</p>
                      <p className="font-medium">{ticket.tenant}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Tag className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm text-muted-foreground">Category</p>
                      <p className="font-medium">{ticket.category}</p>
                    </div>
                  </div>
                </div>
                <Separator />
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <CalendarClock className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm text-muted-foreground">Created</p>
                      <p className="font-medium">{formatDate(ticket.createdAt)}</p>
                      <p className="text-xs text-muted-foreground">{getTimeElapsed(ticket.createdAt)}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Clock className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm text-muted-foreground">Last Updated</p>
                      <p className="font-medium">{formatDate(ticket.updatedAt)}</p>
                      <p className="text-xs text-muted-foreground">{getTimeElapsed(ticket.updatedAt)}</p>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
            
            {/* Related Tickets */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Related Tickets</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                {ticket.relatedTickets.map((related) => (
                  <div key={related.id} className="flex items-center justify-between">
                    <div>
                      <Link to={`/support/tickets/${related.id}`} className="font-mono text-xs hover:underline">
                        {related.id}
                      </Link>
                      <p className="text-sm truncate max-w-[200px]">{related.subject}</p>
                    </div>
                    {getStatusBadge(related.status)}
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
