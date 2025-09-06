import { useState, useRef, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Avatar } from "@/components/ui/avatar";
import { 
  MessageCircle, 
  Send, 
  X, 
  Minimize2, 
  Maximize2, 
  Phone, 
  Video,
  Paperclip,
  Smile,
  MoreVertical,
  User
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { toast } from "sonner";

const mockConversations = [
  {
    id: 1,
    customer: "Alice Johnson",
    email: "alice.j@email.com",
    status: "active",
    lastMessage: "I need help with my recent order",
    timestamp: "2 min ago",
    unreadCount: 3,
    messages: [
      { id: 1, sender: "customer", content: "Hello, I need help with my recent order", timestamp: "2:30 PM" },
      { id: 2, sender: "agent", content: "Hi Alice! I'd be happy to help you with your order. Can you please provide me with your order number?", timestamp: "2:31 PM" },
      { id: 3, sender: "customer", content: "Sure, it's #12345", timestamp: "2:32 PM" },
      { id: 4, sender: "customer", content: "I haven't received it yet and it's been 5 days", timestamp: "2:33 PM" },
    ]
  },
  {
    id: 2,
    customer: "Bob Smith",
    email: "bob.smith@email.com", 
    status: "waiting",
    lastMessage: "Thanks for your help!",
    timestamp: "15 min ago",
    unreadCount: 0,
    messages: [
      { id: 1, sender: "customer", content: "Hi, I have a question about returns", timestamp: "1:15 PM" },
      { id: 2, sender: "agent", content: "Of course! What would you like to know about our return policy?", timestamp: "1:16 PM" },
      { id: 3, sender: "customer", content: "Thanks for your help!", timestamp: "1:20 PM" },
    ]
  },
  {
    id: 3,
    customer: "Carol Davis",
    email: "carol.d@email.com",
    status: "new",
    lastMessage: "Hello, is anyone there?",
    timestamp: "30 min ago",
    unreadCount: 1,
    messages: [
      { id: 1, sender: "customer", content: "Hello, is anyone there?", timestamp: "2:00 PM" },
    ]
  }
];

export default function LiveChat({ isOpen, onToggle, isMinimized, onMinimize }) {
  const [conversations, setConversations] = useState(mockConversations);
  const [activeConversation, setActiveConversation] = useState(null);
  const [newMessage, setNewMessage] = useState("");
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [activeConversation?.messages]);

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (!newMessage.trim() || !activeConversation) return;

    const message = {
      id: Date.now(),
      sender: "agent",
      content: newMessage,
      timestamp: new Date().toLocaleTimeString('en-US', { 
        hour: '2-digit', 
        minute: '2-digit' 
      })
    };

    setConversations(prev => prev.map(conv => 
      conv.id === activeConversation.id 
        ? {
            ...conv,
            messages: [...conv.messages, message],
            lastMessage: newMessage,
            timestamp: "now"
          }
        : conv
    ));

    setActiveConversation(prev => ({
      ...prev,
      messages: [...prev.messages, message]
    }));

    setNewMessage("");
    toast.success("Message sent!");
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'active': return 'bg-green-500';
      case 'waiting': return 'bg-yellow-500';
      case 'new': return 'bg-blue-500';
      default: return 'bg-gray-500';
    }
  };

  const getStatusText = (status) => {
    switch (status) {
      case 'active': return 'Active';
      case 'waiting': return 'Waiting';
      case 'new': return 'New';
      default: return 'Offline';
    }
  };

  const totalUnreadCount = conversations.reduce((sum, conv) => sum + conv.unreadCount, 0);

  if (!isOpen) {
    return (
      <div className="fixed bottom-4 right-4 z-50">
        <Button
          onClick={onToggle}
          size="lg"
          className="rounded-full h-14 w-14 shadow-lg relative"
        >
          <MessageCircle className="h-6 w-6" />
          {totalUnreadCount > 0 && (
            <Badge 
              variant="destructive" 
              className="absolute -top-2 -right-2 h-6 w-6 rounded-full p-0 flex items-center justify-center text-xs"
            >
              {totalUnreadCount}
            </Badge>
          )}
        </Button>
      </div>
    );
  }

  return (
    <div className={`fixed bottom-4 right-4 z-50 bg-background border rounded-lg shadow-xl transition-all duration-300 ${
      isMinimized ? 'h-14 w-80' : 'h-96 w-80 md:h-[500px] md:w-96'
    }`}>
      {/* Chat Header */}
      <div className="flex items-center justify-between p-3 border-b bg-primary text-primary-foreground rounded-t-lg">
        <div className="flex items-center gap-2">
          <MessageCircle className="h-5 w-5" />
          <h3 className="font-semibold">
            {activeConversation ? activeConversation.customer : "Live Support"}
          </h3>
          {totalUnreadCount > 0 && !isMinimized && (
            <Badge variant="secondary" className="text-xs">
              {totalUnreadCount} new
            </Badge>
          )}
        </div>
        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="icon"
            onClick={onMinimize}
            className="h-8 w-8 text-primary-foreground hover:bg-primary-foreground/20"
          >
            {isMinimized ? <Maximize2 className="h-4 w-4" /> : <Minimize2 className="h-4 w-4" />}
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={onToggle}
            className="h-8 w-8 text-primary-foreground hover:bg-primary-foreground/20"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {!isMinimized && (
        <div className="flex h-full">
          {/* Conversations List */}
          <div className="w-full md:w-1/3 border-r">
            <div className="p-2">
              <div className="text-xs font-medium text-muted-foreground mb-2">
                Active Conversations ({conversations.length})
              </div>
              <div className="space-y-1 max-h-[200px] md:max-h-[350px] overflow-y-auto">
                {conversations.map((conversation) => (
                  <div
                    key={conversation.id}
                    onClick={() => setActiveConversation(conversation)}
                    className={`p-2 rounded cursor-pointer transition-colors ${
                      activeConversation?.id === conversation.id
                        ? 'bg-primary/10 border border-primary/20'
                        : 'hover:bg-muted/50'
                    }`}
                  >
                    <div className="flex items-center gap-2 mb-1">
                      <div className="relative">
                        <Avatar className="h-8 w-8">
                          <User className="h-4 w-4" />
                        </Avatar>
                        <div className={`absolute -bottom-0 -right-0 h-3 w-3 rounded-full border-2 border-background ${getStatusColor(conversation.status)}`}></div>
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between">
                          <div className="font-medium text-sm truncate">{conversation.customer}</div>
                          {conversation.unreadCount > 0 && (
                            <Badge variant="secondary" className="h-4 text-xs ml-1">
                              {conversation.unreadCount}
                            </Badge>
                          )}
                        </div>
                        <div className="text-xs text-muted-foreground truncate">
                          {conversation.lastMessage}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {conversation.timestamp}
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Chat Window */}
          <div className="flex-1 flex flex-col">
            {activeConversation ? (
              <>
                {/* Chat Header */}
                <div className="p-3 border-b bg-muted/30">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Avatar className="h-8 w-8">
                        <User className="h-4 w-4" />
                      </Avatar>
                      <div>
                        <div className="font-medium text-sm">{activeConversation.customer}</div>
                        <div className="flex items-center gap-1">
                          <div className={`h-2 w-2 rounded-full ${getStatusColor(activeConversation.status)}`}></div>
                          <div className="text-xs text-muted-foreground">{getStatusText(activeConversation.status)}</div>
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-1">
                      <Button variant="ghost" size="icon" className="h-8 w-8">
                        <Phone className="h-4 w-4" />
                      </Button>
                      <Button variant="ghost" size="icon" className="h-8 w-8">
                        <Video className="h-4 w-4" />
                      </Button>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8">
                            <MoreVertical className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem>View Profile</DropdownMenuItem>
                          <DropdownMenuItem>Create Ticket</DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem>End Conversation</DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                </div>

                {/* Messages */}
                <div className="flex-1 p-3 overflow-y-auto space-y-3">
                  {activeConversation.messages.map((message) => (
                    <div
                      key={message.id}
                      className={`flex ${message.sender === 'agent' ? 'justify-end' : 'justify-start'}`}
                    >
                      <div className={`max-w-[80%] px-3 py-2 rounded-lg text-sm ${
                        message.sender === 'agent'
                          ? 'bg-primary text-primary-foreground rounded-br-sm'
                          : 'bg-muted rounded-bl-sm'
                      }`}>
                        <div>{message.content}</div>
                        <div className={`text-xs mt-1 ${
                          message.sender === 'agent' ? 'text-primary-foreground/70' : 'text-muted-foreground'
                        }`}>
                          {message.timestamp}
                        </div>
                      </div>
                    </div>
                  ))}
                  <div ref={messagesEndRef} />
                </div>

                {/* Message Input */}
                <form onSubmit={handleSendMessage} className="p-3 border-t">
                  <div className="flex items-center gap-2">
                    <Button type="button" variant="ghost" size="icon" className="h-8 w-8">
                      <Paperclip className="h-4 w-4" />
                    </Button>
                    <Input
                      value={newMessage}
                      onChange={(e) => setNewMessage(e.target.value)}
                      placeholder="Type a message..."
                      className="flex-1"
                    />
                    <Button type="button" variant="ghost" size="icon" className="h-8 w-8">
                      <Smile className="h-4 w-4" />
                    </Button>
                    <Button type="submit" size="icon" className="h-8 w-8">
                      <Send className="h-4 w-4" />
                    </Button>
                  </div>
                </form>
              </>
            ) : (
              <div className="flex-1 flex items-center justify-center p-8">
                <div className="text-center">
                  <MessageCircle className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
                  <p className="text-muted-foreground">Select a conversation to start chatting</p>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}