import { useState, useEffect, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { 
  ArrowLeft, 
  Printer, 
  Download, 
  Edit, 
  Truck,
  Package,
  CreditCard,
  User,
  Clock,
  Phone,
  Mail,
  CheckCircle,
  RefreshCw,
  MessageSquare,
  ExternalLink,
  Copy
} from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Separator } from "@/components/ui/separator";
import { Progress } from "@/components/ui/progress";
import { toast } from "sonner";
import { useOrder, useUpdateOrderStatus } from "@/hooks/useApi";
import { generateInvoicePDF } from "@/utils/pdfGenerator";
import { useReactToPrint } from "react-to-print";
import InvoicePDF from "@/components/invoices/InvoicePDF";

// API calls - get order data from backend
const getOrderData = (order) => {
  if (!order) return null;
  
  // Transform API response to match expected structure
  return {
    ...order,
    date_created: order.created_at,
    date_updated: order.updated_at,
    // Ensure we have the right field names
    id: order.id || order.order_number,
    total: order.final_price || order.total,
    subtotal: order.total_price || order.subtotal,
    discount: order.discount_amount || order.discount || 0,
    shipping_cost: order.shipping_cost || 0,
    payment_status: order.payment_status || 'pending',
    // Mock shipping data structure until we have it in backend
    shipping: {
      method: 'Standard Shipping',
      carrier: 'Local Delivery',
      tracking_number: null,
      estimated_delivery: new Date(Date.now() + 3 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      address: typeof order.shipping_address === 'string' ? 
        { street: order.shipping_address } : 
        (order.shipping_address || order.billing?.address || {
          street: "Shipping address not available",
          city: "",
          state: "",
          postal_code: "",
          country: "Bangladesh"
        })
    },
    // Mock timeline until we implement it
    timeline: [
      {
        status: "placed",
        timestamp: order.created_at,
        description: "Order placed by customer",
        user: "System"
      },
      {
        status: order.status,
        timestamp: order.updated_at,
        description: `Order status: ${order.status}`,
        user: "Admin"
      }
    ]
  };
};

export default function OrderDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [status, setStatus] = useState("");
  const [isEditingNotes, setIsEditingNotes] = useState(false);
  const [internalNotes, setInternalNotes] = useState("");
  const [trackingDialog, setTrackingDialog] = useState(false);
  const [newTracking, setNewTracking] = useState("");
  const [isGeneratingPDF, setIsGeneratingPDF] = useState(false);
  const printRef = useRef();
  
  // API calls
  const { 
    data: orderResponse, 
    isLoading, 
    isError, 
    error 
  } = useOrder(id);
  
  // Transform the order data
  const order = orderResponse?.data ? getOrderData(orderResponse.data) : null;
  
  const updateOrderStatusMutation = useUpdateOrderStatus();

  useEffect(() => {
    if (order) {
      setStatus(order.status);
      setInternalNotes(order.internal_notes || "");
    }
  }, [order]);

  const handleStatusUpdate = async (newStatus) => {
    try {
      await updateOrderStatusMutation.mutateAsync({ id, status: newStatus });
      setStatus(newStatus);
      toast.success(`Order status updated to ${newStatus}`);
    } catch (error) {
      toast.error(error.message || "Failed to update order status");
    }
  };

  const handleNotesUpdate = async () => {
    try {
      // TODO: Add API call for updating notes
      // await ordersApi.updateNotes(id, internalNotes);
      setIsEditingNotes(false);
      toast.success("Internal notes updated");
    } catch (error) {
      toast.error("Failed to update notes");
    }
  };

  const handleTrackingUpdate = async () => {
    try {
      // TODO: Add API call for updating tracking
      // await ordersApi.updateTracking(id, newTracking);
      setTrackingDialog(false);
      setNewTracking("");
      toast.success("Tracking number updated");
    } catch (error) {
      toast.error("Failed to update tracking number");
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'pending':
        return 'secondary';
      case 'processing':
        return 'default';
      case 'shipped':
        return 'outline';
      case 'delivered':
        return 'default';
      case 'cancelled':
        return 'destructive';
      default:
        return 'secondary';
    }
  };

  const getStatusProgress = (status) => {
    switch (status) {
      case 'pending':
        return 20;
      case 'processing':
        return 40;
      case 'shipped':
        return 70;
      case 'delivered':
        return 100;
      case 'cancelled':
        return 0;
      default:
        return 0;
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    toast.success("Copied to clipboard");
  };

  const handlePrintInvoice = useReactToPrint({
    content: () => printRef.current,
    documentTitle: `Invoice-${order?.id || 'unknown'}`,
    onAfterPrint: () => {
      toast.success("Invoice printed successfully!");
    },
    onPrintError: () => {
      toast.error("Failed to print invoice. Please try again.");
    }
  });
  
  const handleDownloadPDF = async () => {
    if (!order) {
      toast.error("No order data available");
      return;
    }
    
    setIsGeneratingPDF(true);
    try {
      await generateInvoicePDF(order);
      toast.success("PDF invoice downloaded successfully!");
    } catch (error) {
      console.error("PDF generation error:", error);
      toast.error("Failed to generate PDF. Please try again.");
    } finally {
      setIsGeneratingPDF(false);
    }
  };

  const generateInvoiceHTML = (order, currentStatus) => {
    return `
      <!DOCTYPE html>
      <html>
      <head>
        <title>Invoice ${order.id}</title>
        <style>
          body { font-family: Arial, sans-serif; margin: 20px; color: #333; }
          .header { display: flex; justify-content: space-between; border-bottom: 2px solid #333; padding-bottom: 20px; margin-bottom: 30px; }
          .company-info h1 { margin: 0; font-size: 28px; }
          .invoice-info { text-align: right; }
          .customer-info { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-bottom: 30px; }
          .section-title { font-weight: bold; margin-bottom: 10px; font-size: 16px; color: #333; }
          table { width: 100%; border-collapse: collapse; margin: 20px 0; }
          th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
          th { background-color: #f8f9fa; font-weight: bold; }
          .totals { margin-top: 30px; }
          .total-row { display: flex; justify-content: space-between; margin: 5px 0; }
          .grand-total { font-size: 18px; font-weight: bold; border-top: 2px solid #333; padding-top: 10px; margin-top: 15px; }
          .footer { margin-top: 50px; text-align: center; font-size: 12px; color: #666; }
          @media print { body { -webkit-print-color-adjust: exact; } }
        </style>
      </head>
      <body>
        <div class="header">
          <div class="company-info">
            <h1>ShopVendor</h1>
            <p>Single Vendor E-commerce</p>
            <p>📧 admin@shopvendor.com | 📞 (555) 123-4567</p>
          </div>
          <div class="invoice-info">
            <h2>INVOICE</h2>
            <p><strong>Order ID:</strong> ${order.id}</p>
            <p><strong>Date:</strong> ${formatDate(order.date_created)}</p>
            <p><strong>Status:</strong> ${currentStatus.toUpperCase()}</p>
          </div>
        </div>
        
        <div class="customer-info">
          <div>
            <div class="section-title">Bill To:</div>
            <p><strong>${order.customer.name}</strong></p>
            <p>${order.customer.email}</p>
            <p>${order.customer.phone}</p>
            <p>${order.billing.address.street}</p>
            <p>${order.billing.address.city}, ${order.billing.address.state} ${order.billing.address.postal_code}</p>
            <p>${order.billing.address.country}</p>
          </div>
          <div>
            <div class="section-title">Ship To:</div>
            <p>${order.shipping.address.street}</p>
            ${order.shipping.address.apartment ? `<p>${order.shipping.address.apartment}</p>` : ''}
            <p>${order.shipping.address.city}, ${order.shipping.address.state} ${order.shipping.address.postal_code}</p>
            <p>${order.shipping.address.country}</p>
            <br>
            <div class="section-title">Payment Method:</div>
            <p>${order.billing.payment_method}</p>
            <p>Transaction: ${order.billing.transaction_id}</p>
          </div>
        </div>
        
        <table>
          <thead>
            <tr>
              <th>Item</th>
              <th>SKU</th>
              <th>Variant</th>
              <th>Qty</th>
              <th>Unit Price</th>
              <th>Total</th>
            </tr>
          </thead>
          <tbody>
            ${order.items.map(item => `
              <tr>
                <td>${item.name}</td>
                <td>${item.sku}</td>
                <td>${item.variant || '-'}</td>
                <td>${item.qty}</td>
                <td>${formatCurrency(item.price)}</td>
                <td>${formatCurrency(item.price * item.qty)}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>
        
        <div class="totals">
          <div class="total-row">
            <span>Subtotal:</span>
            <span>${formatCurrency(order.subtotal)}</span>
          </div>
          ${order.discount > 0 ? `
            <div class="total-row">
              <span>Discount:</span>
              <span>-${formatCurrency(order.discount)}</span>
            </div>
          ` : ''}
          <div class="total-row">
            <span>Shipping:</span>
            <span>${formatCurrency(order.shipping_cost)}</span>
          </div>
          <div class="total-row grand-total">
            <span>Grand Total:</span>
            <span>${formatCurrency(order.total)}</span>
          </div>
        </div>
        
        ${order.notes ? `
          <div style="margin-top: 30px;">
            <div class="section-title">Notes:</div>
            <p>${order.notes}</p>
          </div>
        ` : ''}
        
        <div class="footer">
          <p>Thank you for your business!</p>
          <p>For support, contact us at support@shopvendor.com</p>
        </div>
        
        <script>
          window.onload = function() {
            window.print();
          };
        </script>
      </body>
      </html>
    `;
  };

  if (isError) {
    return (
      <div className="space-y-6 p-6">
        <Card>
          <CardHeader>
            <CardTitle>Error Loading Order</CardTitle>
            <CardDescription>{error?.message || 'Failed to load order details'}</CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => navigate("/orders")}>
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Orders
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="space-y-6 p-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-center py-8">
              <RefreshCw className="h-8 w-8 animate-spin" />
              <span className="ml-2">Loading order details...</span>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!order) {
    return (
      <div className="space-y-6 p-6">
        <Card>
          <CardHeader>
            <CardTitle>Order Not Found</CardTitle>
            <CardDescription>The order you are looking for does not exist.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => navigate("/orders")}>
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Orders
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/orders")}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold tracking-tight">Order {order.id}</h1>
          <p className="text-muted-foreground">
            Created {formatDate(order.date_created)} • Last updated {formatDate(order.date_updated)}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button 
            variant="outline" 
            size="sm" 
            onClick={handlePrintInvoice}
            disabled={!order}
          >
            <Printer className="mr-2 h-4 w-4" />
            Print Invoice
          </Button>
          <Button 
            variant="outline" 
            size="sm" 
            onClick={handleDownloadPDF}
            disabled={isGeneratingPDF || !order}
          >
            {isGeneratingPDF ? (
              <RefreshCw className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Download className="mr-2 h-4 w-4" />
            )}
            {isGeneratingPDF ? 'Generating...' : 'Download PDF'}
          </Button>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Order Status & Progress */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <Package className="h-5 w-5" />
                  Order Status
                </CardTitle>
                <div className="flex items-center gap-3">
                  <Select value={status} onValueChange={handleStatusUpdate}>
                    <SelectTrigger className="w-40">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="pending">Pending</SelectItem>
                      <SelectItem value="processing">Processing</SelectItem>
                      <SelectItem value="shipped">Shipped</SelectItem>
                      <SelectItem value="delivered">Delivered</SelectItem>
                      <SelectItem value="cancelled">Cancelled</SelectItem>
                    </SelectContent>
                  </Select>
                  <Badge variant={getStatusColor(status)} className="capitalize">
                    {status}
                  </Badge>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <div className="flex justify-between text-sm mb-2">
                    <span>Order Progress</span>
                    <span>{getStatusProgress(status)}%</span>
                  </div>
                  <Progress value={getStatusProgress(status)} className="h-2" />
                </div>
                
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                  <div className={`p-3 rounded-lg border ${status === 'pending' ? 'bg-blue-50 border-blue-200' : 'bg-gray-50'}`}>
                    <div className="font-medium text-sm">Pending</div>
                    <div className="text-xs text-muted-foreground mt-1">Order placed</div>
                  </div>
                  <div className={`p-3 rounded-lg border ${status === 'processing' ? 'bg-yellow-50 border-yellow-200' : 'bg-gray-50'}`}>
                    <div className="font-medium text-sm">Processing</div>
                    <div className="text-xs text-muted-foreground mt-1">Being prepared</div>
                  </div>
                  <div className={`p-3 rounded-lg border ${status === 'shipped' ? 'bg-purple-50 border-purple-200' : 'bg-gray-50'}`}>
                    <div className="font-medium text-sm">Shipped</div>
                    <div className="text-xs text-muted-foreground mt-1">In transit</div>
                  </div>
                  <div className={`p-3 rounded-lg border ${status === 'delivered' ? 'bg-green-50 border-green-200' : 'bg-gray-50'}`}>
                    <div className="font-medium text-sm">Delivered</div>
                    <div className="text-xs text-muted-foreground mt-1">Completed</div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Order Items */}
          <Card>
            <CardHeader>
              <CardTitle>Order Items</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {order.items.map((item) => (
                  <div key={item.id} className="flex items-center gap-4 p-4 border rounded-lg">
                    <div className="w-16 h-16 bg-gray-100 rounded-lg flex items-center justify-center">
                      <Package className="h-8 w-8 text-gray-400" />
                    </div>
                    <div className="flex-1">
                      <h4 className="font-medium">{item.name}</h4>
                      <p className="text-sm text-muted-foreground">SKU: {item.sku}</p>
                      {item.variant && <p className="text-sm text-muted-foreground">{item.variant}</p>}
                    </div>
                    <div className="text-right">
                      <p className="font-medium">Qty: {item.qty}</p>
                      <p className="text-sm text-muted-foreground">{formatCurrency(item.price)} each</p>
                      <p className="font-medium">{formatCurrency(item.price * item.qty)}</p>
                    </div>
                  </div>
                ))}
              </div>

              <Separator className="my-6" />

              {/* Order Totals */}
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span>Subtotal:</span>
                  <span>{formatCurrency(order.subtotal)}</span>
                </div>
                {order.discount > 0 && (
                  <div className="flex justify-between text-green-600">
                    <span>Discount:</span>
                    <span>-{formatCurrency(order.discount)}</span>
                  </div>
                )}
                <div className="flex justify-between">
                  <span>Shipping:</span>
                  <span>{formatCurrency(order.shipping_cost)}</span>
                </div>
                <Separator />
                <div className="flex justify-between text-lg font-bold">
                  <span>Total:</span>
                  <span>{formatCurrency(order.total)}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Order Timeline */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="h-5 w-5" />
                Order Timeline
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {order.timeline.map((event, index) => (
                  <div key={index} className="flex items-start gap-3">
                    <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0">
                      <CheckCircle className="h-4 w-4 text-blue-600" />
                    </div>
                    <div className="flex-1">
                      <p className="font-medium capitalize">{event.status}</p>
                      <p className="text-sm text-muted-foreground">{event.description}</p>
                      <p className="text-xs text-muted-foreground">
                        {formatDate(event.timestamp)} • {event.user}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Customer Information */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-5 w-5" />
                Customer
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center">
                  <User className="h-5 w-5 text-gray-600" />
                </div>
                <div>
                  <p className="font-medium">{order.customer.name}</p>
                  <p className="text-sm text-muted-foreground">ID: {order.customer.id}</p>
                </div>
              </div>
              
              <Separator />
              
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">{order.customer.email}</span>
                  <Button 
                    variant="ghost" 
                    size="sm" 
                    onClick={() => copyToClipboard(order.customer.email)}
                  >
                    <Copy className="h-3 w-3" />
                  </Button>
                </div>
                <div className="flex items-center gap-2">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">{order.customer.phone}</span>
                  <Button 
                    variant="ghost" 
                    size="sm" 
                    onClick={() => copyToClipboard(order.customer.phone)}
                  >
                    <Copy className="h-3 w-3" />
                  </Button>
                </div>
              </div>
              
              <Separator />
              
              <div className="space-y-2">
                <Button variant="outline" size="sm" className="w-full">
                  <MessageSquare className="mr-2 h-4 w-4" />
                  Contact Customer
                </Button>
                <Button 
                  variant="outline" 
                  size="sm" 
                  className="w-full"
                  onClick={() => navigate(`/customers/${order.customer.id}`)}
                >
                  <ExternalLink className="mr-2 h-4 w-4" />
                  View Profile
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Shipping Information */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Truck className="h-5 w-5" />
                Shipping
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-3">
                <div>
                  <Label className="text-sm font-medium">Method</Label>
                  <p className="text-sm">{order.shipping.method} - {order.shipping.carrier}</p>
                </div>
                
                <div>
                  <Label className="text-sm font-medium">Estimated Delivery</Label>
                  <p className="text-sm">{new Date(order.shipping.estimated_delivery).toLocaleDateString()}</p>
                </div>
                
                <div>
                  <Label className="text-sm font-medium">Tracking Number</Label>
                  <div className="flex items-center gap-2">
                    {order.shipping.tracking_number ? (
                      <>
                        <code className="text-sm bg-muted px-2 py-1 rounded">
                          {order.shipping.tracking_number}
                        </code>
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => copyToClipboard(order.shipping.tracking_number)}
                        >
                          <Copy className="h-3 w-3" />
                        </Button>
                      </>
                    ) : (
                      <span className="text-sm text-muted-foreground">Not available</span>
                    )}
                  </div>
                </div>
              </div>
              
              <Separator />
              
              <div>
                <Label className="text-sm font-medium mb-2 block">Shipping Address</Label>
                <div className="text-sm text-muted-foreground space-y-1">
                  <p>{order.shipping.address.street}</p>
                  {order.shipping.address.apartment && <p>{order.shipping.address.apartment}</p>}
                  <p>{order.shipping.address.city}, {order.shipping.address.state} {order.shipping.address.postal_code}</p>
                  <p>{order.shipping.address.country}</p>
                </div>
              </div>
              
              <Separator />
              
              <Dialog open={trackingDialog} onOpenChange={setTrackingDialog}>
                <DialogTrigger asChild>
                  <Button variant="outline" size="sm" className="w-full">
                    <Edit className="mr-2 h-4 w-4" />
                    Update Tracking
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Update Tracking Number</DialogTitle>
                    <DialogDescription>
                      Enter the tracking number for this shipment.
                    </DialogDescription>
                  </DialogHeader>
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="tracking">Tracking Number</Label>
                      <Input
                        id="tracking"
                        value={newTracking}
                        onChange={(e) => setNewTracking(e.target.value)}
                        placeholder="Enter tracking number"
                      />
                    </div>
                  </div>
                  <DialogFooter>
                    <Button variant="outline" onClick={() => setTrackingDialog(false)}>
                      Cancel
                    </Button>
                    <Button onClick={handleTrackingUpdate} disabled={!newTracking.trim()}>
                      Update
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </CardContent>
          </Card>

          {/* Payment Information */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <CreditCard className="h-5 w-5" />
                Payment
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-3">
                <div>
                  <Label className="text-sm font-medium">Method</Label>
                  <p className="text-sm">{order.billing.payment_method}</p>
                </div>
                
                <div>
                  <Label className="text-sm font-medium">Status</Label>
                  <Badge variant={order.payment_status === 'paid' ? 'default' : 'secondary'}>
                    {order.payment_status}
                  </Badge>
                </div>
                
                <div>
                  <Label className="text-sm font-medium">Transaction ID</Label>
                  <div className="flex items-center gap-2">
                    <code className="text-sm bg-muted px-2 py-1 rounded">
                      {order.billing.transaction_id}
                    </code>
                    <Button 
                      variant="ghost" 
                      size="sm"
                      onClick={() => copyToClipboard(order.billing.transaction_id)}
                    >
                      <Copy className="h-3 w-3" />
                    </Button>
                  </div>
                </div>
              </div>
              
              <Separator />
              
              <div>
                <Label className="text-sm font-medium mb-2 block">Billing Address</Label>
                <div className="text-sm text-muted-foreground space-y-1">
                  <p>{order.billing.address.street}</p>
                  {order.billing.address.apartment && <p>{order.billing.address.apartment}</p>}
                  <p>{order.billing.address.city}, {order.billing.address.state} {order.billing.address.postal_code}</p>
                  <p>{order.billing.address.country}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Notes */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <MessageSquare className="h-5 w-5" />
                Notes
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {order.notes && (
                <div>
                  <Label className="text-sm font-medium">Customer Notes</Label>
                  <p className="text-sm text-muted-foreground bg-muted p-3 rounded-lg mt-1">
                    {order.notes}
                  </p>
                </div>
              )}
              
              <div>
                <Label className="text-sm font-medium">Internal Notes</Label>
                {isEditingNotes ? (
                  <div className="space-y-2 mt-1">
                    <Textarea
                      value={internalNotes}
                      onChange={(e) => setInternalNotes(e.target.value)}
                      placeholder="Add internal notes..."
                      rows={3}
                    />
                    <div className="flex gap-2">
                      <Button size="sm" onClick={handleNotesUpdate}>
                        Save
                      </Button>
                      <Button 
                        variant="outline" 
                        size="sm" 
                        onClick={() => setIsEditingNotes(false)}
                      >
                        Cancel
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="mt-1">
                    <p className="text-sm text-muted-foreground bg-muted p-3 rounded-lg min-h-[60px]">
                      {internalNotes || "No internal notes"}
                    </p>
                    <Button 
                      variant="outline" 
                      size="sm" 
                      className="mt-2"
                      onClick={() => setIsEditingNotes(true)}
                    >
                      <Edit className="mr-2 h-4 w-4" />
                      Edit Notes
                    </Button>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
      
      {/* Hidden Invoice Component for Printing */}
      <div style={{ display: 'none' }}>
        <InvoicePDF 
          ref={printRef} 
          order={order} 
          companyInfo={{
            name: "ShopVendor",
            subtitle: "Single Vendor E-commerce",
            email: "admin@shopvendor.com",
            phone: "(555) 123-4567",
            website: "www.shopvendor.com"
          }} 
        />
      </div>
    </div>
  );
}