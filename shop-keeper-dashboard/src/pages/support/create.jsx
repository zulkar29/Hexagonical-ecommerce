import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Form, FormItem, FormLabel, FormControl, FormMessage } from "@/components/ui/form";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { useNavigate } from "react-router-dom";
import { ArrowLeft, Paperclip, AlertCircle, Clock, Zap } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

const formSchema = z.object({
  subject: z.string().min(5, "Subject must be at least 5 characters"),
  customer: z.string().min(2, "Customer name is required"),
  email: z.string().email("Please enter a valid email address"),
  phone: z.string().optional(),
  priority: z.enum(["low", "medium", "high", "urgent"], {
    required_error: "Please select a priority level",
  }),
  category: z.enum(["general", "order", "payment", "product", "shipping", "technical"], {
    required_error: "Please select a category",
  }),
  message: z.string().min(10, "Message must be at least 10 characters"),
});

export default function CreateSupportTicket() {
  const [attachments, setAttachments] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  const methods = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      subject: "",
      customer: "",
      email: "",
      phone: "",
      priority: "",
      category: "",
      message: "",
    },
  });
  const { handleSubmit, register, setValue, watch, formState: { errors } } = methods;
  const navigate = useNavigate();

  const handleFileAttach = (event) => {
    const files = Array.from(event.target.files);
    const validFiles = files.filter(file => file.size <= 5 * 1024 * 1024); // 5MB limit
    setAttachments(prev => [...prev, ...validFiles]);
  };

  const removeAttachment = (index) => {
    setAttachments(prev => prev.filter((_, i) => i !== index));
  };

  const onSubmit = async () => {
    setIsSubmitting(true);
    try {
      // In a real app, send data to API with attachments
      await new Promise(resolve => setTimeout(resolve, 1500)); // Simulate API call
      toast.success("Support ticket created successfully!");
      navigate("/support");
    } catch {
      toast.error("Failed to create support ticket. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const getPriorityIcon = (priority) => {
    switch (priority) {
      case 'low': return <Clock className="w-4 h-4 text-blue-500" />;
      case 'medium': return <AlertCircle className="w-4 h-4 text-yellow-500" />;
      case 'high': return <AlertCircle className="w-4 h-4 text-orange-500" />;
      case 'urgent': return <Zap className="w-4 h-4 text-red-500" />;
      default: return null;
    }
  };

  const priority = watch('priority');

  return (
    <div className="space-y-6 p-6">
      <Card>
        <CardHeader className="px-8">
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => navigate("/support")}
              className="mr-2"
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <CardTitle>Create Support Ticket</CardTitle>
              <CardDescription>Fill out the form to submit a new support request.</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="px-8">
          <Form {...methods}>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <FormItem>
                  <FormLabel>Subject *</FormLabel>
                  <FormControl>
                    <Input {...register("subject")} placeholder="Brief description of your issue" />
                  </FormControl>
                  <FormMessage>{errors.subject?.message}</FormMessage>
                </FormItem>
                
                <FormItem>
                  <FormLabel>Category *</FormLabel>
                  <Select onValueChange={(value) => setValue('category', value)}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select category" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="general">General Inquiry</SelectItem>
                      <SelectItem value="order">Order Issues</SelectItem>
                      <SelectItem value="payment">Payment Problems</SelectItem>
                      <SelectItem value="product">Product Issues</SelectItem>
                      <SelectItem value="shipping">Shipping & Delivery</SelectItem>
                      <SelectItem value="technical">Technical Support</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage>{errors.category?.message}</FormMessage>
                </FormItem>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <FormItem>
                  <FormLabel>Customer Name *</FormLabel>
                  <FormControl>
                    <Input {...register("customer")} placeholder="Your full name" />
                  </FormControl>
                  <FormMessage>{errors.customer?.message}</FormMessage>
                </FormItem>
                
                <FormItem>
                  <FormLabel>Email *</FormLabel>
                  <FormControl>
                    <Input type="email" {...register("email")} placeholder="you@email.com" />
                  </FormControl>
                  <FormMessage>{errors.email?.message}</FormMessage>
                </FormItem>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <FormItem>
                  <FormLabel>Phone Number</FormLabel>
                  <FormControl>
                    <Input {...register("phone")} placeholder="+1 (555) 123-4567" />
                  </FormControl>
                  <FormMessage>{errors.phone?.message}</FormMessage>
                </FormItem>
                
                <FormItem>
                  <FormLabel>Priority *</FormLabel>
                  <Select onValueChange={(value) => setValue('priority', value)}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select priority level" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="low">
                        <div className="flex items-center gap-2">
                          <Clock className="w-4 h-4 text-blue-500" />
                          Low - General inquiry
                        </div>
                      </SelectItem>
                      <SelectItem value="medium">
                        <div className="flex items-center gap-2">
                          <AlertCircle className="w-4 h-4 text-yellow-500" />
                          Medium - Standard issue
                        </div>
                      </SelectItem>
                      <SelectItem value="high">
                        <div className="flex items-center gap-2">
                          <AlertCircle className="w-4 h-4 text-orange-500" />
                          High - Important issue
                        </div>
                      </SelectItem>
                      <SelectItem value="urgent">
                        <div className="flex items-center gap-2">
                          <Zap className="w-4 h-4 text-red-500" />
                          Urgent - Critical issue
                        </div>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage>{errors.priority?.message}</FormMessage>
                </FormItem>
              </div>

              <FormItem>
                <FormLabel>Message *</FormLabel>
                <FormControl>
                  <Textarea 
                    {...register("message")} 
                    placeholder="Provide detailed information about your issue..."
                    className="min-h-[120px] resize-none"
                  />
                </FormControl>
                <FormMessage>{errors.message?.message}</FormMessage>
              </FormItem>

              <FormItem>
                <FormLabel>Attachments</FormLabel>
                <div className="space-y-4">
                  <div className="flex items-center gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => document.getElementById('file-upload').click()}
                    >
                      <Paperclip className="w-4 h-4 mr-2" />
                      Attach Files
                    </Button>
                    <span className="text-sm text-muted-foreground">
                      Max 5MB per file
                    </span>
                  </div>
                  <input
                    id="file-upload"
                    type="file"
                    multiple
                    onChange={handleFileAttach}
                    className="hidden"
                    accept="image/*,.pdf,.doc,.docx,.txt"
                  />
                  
                  {attachments.length > 0 && (
                    <div className="space-y-2">
                      {attachments.map((file, index) => (
                        <div key={index} className="flex items-center justify-between p-2 bg-muted rounded-md">
                          <span className="text-sm truncate">{file.name}</span>
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => removeAttachment(index)}
                          >
                            ×
                          </Button>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </FormItem>

              <div className="flex justify-end gap-2 pt-6">
                <Button 
                  type="button" 
                  variant="secondary" 
                  onClick={() => navigate("/support")}
                  disabled={isSubmitting}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? "Creating..." : "Submit Ticket"}
                  {priority && <span className="ml-2">{getPriorityIcon(priority)}</span>}
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
