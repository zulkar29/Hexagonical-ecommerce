import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Mail, CheckCircle, XCircle, Loader2, ArrowLeft, Paperclip, Phone, Clock, AlertCircle, Zap, Send, User, Users } from "lucide-react";
import { toast } from "sonner";

const mockTickets = [
  { 
    id: 1, 
    subject: "Order not received", 
    customer: "John Doe", 
    email: "john.doe@example.com", 
    phone: "+1 (555) 123-4567",
    date: "2025-05-05", 
    status: "Open", 
    priority: "high",
    category: "order",
    assignedTo: "Sarah Wilson",
    message: "I placed an order a week ago and haven't received it yet. Order number #12345. Can you please help me track it?" 
  },
  { 
    id: 2, 
    subject: "Refund request", 
    customer: "Jane Smith", 
    email: "jane.smith@example.com", 
    phone: "+1 (555) 987-6543",
    date: "2025-05-04", 
    status: "Pending", 
    priority: "medium",
    category: "payment",
    assignedTo: "Mike Johnson",
    message: "I would like a refund for my last purchase. The product doesn't meet my expectations." 
  },
  { 
    id: 3, 
    subject: "Product damaged", 
    customer: "Emily Davis", 
    email: "emily.davis@example.com", 
    phone: null,
    date: "2025-05-03", 
    status: "Closed", 
    priority: "urgent",
    category: "product",
    assignedTo: "David Lee",
    message: "The product arrived damaged. Please advise on next steps for replacement or refund." 
  },
  { 
    id: 4, 
    subject: "Payment issue", 
    customer: "Michael Brown", 
    email: "michael.b@example.com", 
    phone: "+1 (555) 456-7890",
    date: "2025-05-02", 
    status: "Open", 
    priority: "low",
    category: "payment",
    assignedTo: null,
    message: "My payment did not go through but money was deducted from my account." 
  },
];

export default function SupportDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [reply, setReply] = useState("");
  const [attachments, setAttachments] = useState([]);
  const [currentStatus, setCurrentStatus] = useState("");
  const [currentPriority, setCurrentPriority] = useState("");
  const [assignedAgent, setAssignedAgent] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  const ticket = mockTickets.find(t => t.id === Number(id));

  React.useEffect(() => {
    if (ticket) {
      setCurrentStatus(ticket.status);
      setCurrentPriority(ticket.priority);
      setAssignedAgent(ticket.assignedTo || "unassigned");
    }
  }, [ticket]);

  const mockMessages = [
    { 
      sender: "customer", 
      message: ticket?.message, 
      date: ticket?.date, 
      time: "2:30 PM",
      attachments: []
    },
    { 
      sender: "agent", 
      message: "Thank you for reaching out. We're looking into your issue and will get back to you shortly.", 
      date: "2025-05-06", 
      time: "3:45 PM",
      agentName: "Sarah Wilson",
      attachments: []
    },
    {
      sender: "customer",
      message: "Thank you for the quick response. Do you have an estimated timeline for resolution?",
      date: "2025-05-06",
      time: "4:15 PM",
      attachments: []
    }
  ];

  const [messages, setMessages] = useState(mockMessages);

  const getPriorityIcon = (priority) => {
    switch (priority) {
      case 'low': return <Clock className="w-4 h-4 text-blue-500" />;
      case 'medium': return <AlertCircle className="w-4 h-4 text-yellow-500" />;
      case 'high': return <AlertCircle className="w-4 h-4 text-orange-500" />;
      case 'urgent': return <Zap className="w-4 h-4 text-red-500" />;
      default: return null;
    }
  };

  const getPriorityColor = (priority) => {
    switch (priority) {
      case 'low': return 'text-blue-600 bg-blue-50 border-blue-200';
      case 'medium': return 'text-yellow-600 bg-yellow-50 border-yellow-200';
      case 'high': return 'text-orange-600 bg-orange-50 border-orange-200';
      case 'urgent': return 'text-red-600 bg-red-50 border-red-200';
      default: return 'text-gray-600 bg-gray-50 border-gray-200';
    }
  };

  const handleFileAttach = (event) => {
    const files = Array.from(event.target.files);
    const validFiles = files.filter(file => file.size <= 5 * 1024 * 1024);
    setAttachments(prev => [...prev, ...validFiles]);
  };

  const removeAttachment = (index) => {
    setAttachments(prev => prev.filter((_, i) => i !== index));
  };

  const handleReply = async () => {
    if (!reply.trim()) return;
    
    setIsSubmitting(true);
    try {
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call
      
      setMessages([
        ...messages,
        { 
          sender: "agent", 
          message: reply, 
          date: new Date().toISOString().slice(0, 10),
          time: new Date().toLocaleTimeString('en-US', { 
            hour: '2-digit', 
            minute: '2-digit' 
          }),
          agentName: "You",
          attachments: [...attachments]
        }
      ]);
      
      toast.success("Reply sent successfully!");
      setReply("");
      setAttachments([]);
    } catch {
      toast.error("Failed to send reply. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleStatusUpdate = async (newStatus) => {
    try {
      setCurrentStatus(newStatus);
      toast.success(`Ticket status updated to ${newStatus}`);
    } catch {
      toast.error("Failed to update status");
    }
  };

  const handlePriorityUpdate = async (newPriority) => {
    try {
      setCurrentPriority(newPriority);
      toast.success(`Priority updated to ${newPriority}`);
    } catch {
      toast.error("Failed to update priority");
    }
  };

  const handleAgentAssignment = async (agent) => {
    try {
      setAssignedAgent(agent);
      toast.success(`Ticket ${agent === "unassigned" ? "unassigned" : `assigned to ${agent}`}`);
    } catch {
      toast.error("Failed to assign agent");
    }
  };

  if (!ticket) {
    return (
      <div className="space-y-6 p-6">
        <Card>
          <CardHeader className="px-2">
            <CardTitle>Ticket Not Found</CardTitle>
            <CardDescription>The support ticket you are looking for does not exist.</CardDescription>
          </CardHeader>
          <CardContent className="px-8">
            <Button onClick={() => navigate("/support")}> <ArrowLeft className="mr-2 h-4 w-4" />Back to Tickets</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6 p-6">
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main Ticket Content */}
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardHeader className="px-8 pb-4">
              <div className="flex items-center gap-2 mb-4">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => navigate("/support")}
                  className="mr-2"
                >
                  <ArrowLeft className="h-4 w-4" />
                </Button>
                <div className="flex-1">
                  <CardTitle className="text-xl mb-2">{ticket.subject}</CardTitle>
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    <span>Created {ticket.date}</span>
                    <span>•</span>
                    <span>Ticket #{ticket.id}</span>
                  </div>
                </div>
              </div>
            </CardHeader>
            
            <CardContent className="px-8 space-y-6">
              {/* Customer Info */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 p-4 bg-muted/30 rounded-lg">
                <div className="flex items-center gap-2">
                  <User className="w-4 h-4 text-muted-foreground" />
                  <div>
                    <div className="font-medium">{ticket.customer}</div>
                    <div className="text-sm text-muted-foreground">Customer</div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Mail className="w-4 h-4 text-muted-foreground" />
                  <div>
                    <div className="font-medium">{ticket.email}</div>
                    <div className="text-sm text-muted-foreground">Email</div>
                  </div>
                </div>
                {ticket.phone && (
                  <div className="flex items-center gap-2">
                    <Phone className="w-4 h-4 text-muted-foreground" />
                    <div>
                      <div className="font-medium">{ticket.phone}</div>
                      <div className="text-sm text-muted-foreground">Phone</div>
                    </div>
                  </div>
                )}
              </div>

              {/* Conversation */}
              <div>
                <h3 className="font-semibold mb-4">Conversation</h3>
                <div className="space-y-4 max-h-[500px] overflow-y-auto">
                  {messages.map((msg, idx) => (
                    <div key={idx} className={`flex ${msg.sender === "agent" ? "justify-end" : "justify-start"}`}>
                      <div className={`max-w-[80%] ${msg.sender === "agent" ? "order-2" : "order-1"}`}>
                        <div className={`px-4 py-3 rounded-lg text-sm shadow-sm
                          ${msg.sender === "agent" 
                            ? "bg-primary text-primary-foreground rounded-br-sm" 
                            : "bg-white border border-border rounded-bl-sm"}
                        `}>
                          {msg.sender === "agent" && msg.agentName && (
                            <div className="text-xs opacity-75 mb-1">{msg.agentName}</div>
                          )}
                          <div className="whitespace-pre-line">{msg.message}</div>
                          {msg.attachments && msg.attachments.length > 0 && (
                            <div className="mt-2 space-y-1">
                              {msg.attachments.map((file, fileIdx) => (
                                <div key={fileIdx} className="text-xs flex items-center gap-1">
                                  <Paperclip className="w-3 h-3" />
                                  {file.name}
                                </div>
                              ))}
                            </div>
                          )}
                          <div className="text-xs opacity-70 mt-2">
                            {msg.date} at {msg.time}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Reply Section */}
              <div className="border-t pt-6">
                <div className="space-y-4">
                  <h4 className="font-semibold">Add Reply</h4>
                  <Textarea
                    placeholder="Type your response here..."
                    value={reply}
                    onChange={(e) => setReply(e.target.value)}
                    className="min-h-[100px] resize-none"
                    disabled={currentStatus === 'Closed'}
                  />
                  
                  {/* Attachments */}
                  <div className="space-y-2">
                    <div className="flex items-center gap-2">
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => document.getElementById('reply-file-upload').click()}
                        disabled={currentStatus === 'Closed'}
                      >
                        <Paperclip className="w-4 h-4 mr-2" />
                        Attach Files
                      </Button>
                      <span className="text-xs text-muted-foreground">Max 5MB per file</span>
                    </div>
                    <input
                      id="reply-file-upload"
                      type="file"
                      multiple
                      onChange={handleFileAttach}
                      className="hidden"
                      accept="image/*,.pdf,.doc,.docx,.txt"
                    />
                    
                    {attachments.length > 0 && (
                      <div className="space-y-2">
                        {attachments.map((file, index) => (
                          <div key={index} className="flex items-center justify-between p-2 bg-muted/50 rounded-md text-sm">
                            <span className="truncate">{file.name}</span>
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              onClick={() => removeAttachment(index)}
                              className="h-6 w-6 p-0"
                            >
                              ×
                            </Button>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>

                  <div className="flex justify-end">
                    <Button 
                      onClick={handleReply} 
                      disabled={!reply.trim() || currentStatus === 'Closed' || isSubmitting}
                    >
                      {isSubmitting ? (
                        <>
                          <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                          Sending...
                        </>
                      ) : (
                        <>
                          <Send className="w-4 h-4 mr-2" />
                          Send Reply
                        </>
                      )}
                    </Button>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Ticket Sidebar */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Ticket Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Status */}
              <div>
                <label className="text-sm font-medium mb-2 block">Status</label>
                <Select value={currentStatus} onValueChange={handleStatusUpdate}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="Open">
                      <div className="flex items-center gap-2">
                        <CheckCircle className="w-4 h-4 text-green-600" />
                        Open
                      </div>
                    </SelectItem>
                    <SelectItem value="Pending">
                      <div className="flex items-center gap-2">
                        <Loader2 className="w-4 h-4 text-yellow-600" />
                        Pending
                      </div>
                    </SelectItem>
                    <SelectItem value="Closed">
                      <div className="flex items-center gap-2">
                        <XCircle className="w-4 h-4 text-gray-600" />
                        Closed
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Priority */}
              <div>
                <label className="text-sm font-medium mb-2 block">Priority</label>
                <Select value={currentPriority} onValueChange={handlePriorityUpdate}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="low">
                      <div className="flex items-center gap-2">
                        <Clock className="w-4 h-4 text-blue-500" />
                        Low
                      </div>
                    </SelectItem>
                    <SelectItem value="medium">
                      <div className="flex items-center gap-2">
                        <AlertCircle className="w-4 h-4 text-yellow-500" />
                        Medium
                      </div>
                    </SelectItem>
                    <SelectItem value="high">
                      <div className="flex items-center gap-2">
                        <AlertCircle className="w-4 h-4 text-orange-500" />
                        High
                      </div>
                    </SelectItem>
                    <SelectItem value="urgent">
                      <div className="flex items-center gap-2">
                        <Zap className="w-4 h-4 text-red-500" />
                        Urgent
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Category */}
              <div>
                <label className="text-sm font-medium mb-2 block">Category</label>
                <Badge variant="outline" className="capitalize w-full justify-center py-2">
                  {ticket.category}
                </Badge>
              </div>

              {/* Assigned Agent */}
              <div>
                <label className="text-sm font-medium mb-2 block">Assigned Agent</label>
                <Select value={assignedAgent} onValueChange={handleAgentAssignment}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select agent..." />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="unassigned">
                      <div className="flex items-center gap-2">
                        <Users className="w-4 h-4 text-muted-foreground" />
                        Unassigned
                      </div>
                    </SelectItem>
                    <SelectItem value="Sarah Wilson">
                      <div className="flex items-center gap-2">
                        <User className="w-4 h-4" />
                        Sarah Wilson
                      </div>
                    </SelectItem>
                    <SelectItem value="Mike Johnson">
                      <div className="flex items-center gap-2">
                        <User className="w-4 h-4" />
                        Mike Johnson
                      </div>
                    </SelectItem>
                    <SelectItem value="David Lee">
                      <div className="flex items-center gap-2">
                        <User className="w-4 h-4" />
                        David Lee
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Priority Badge */}
              <div>
                <label className="text-sm font-medium mb-2 block">Current Priority</label>
                <div className={`inline-flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium border w-full justify-center ${getPriorityColor(currentPriority)}`}>
                  {getPriorityIcon(currentPriority)}
                  {currentPriority} priority
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
