import { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { 
  ArrowLeft, 
  Plus, 
  X, 
  Search, 
  User, 
  Package, 
  CreditCard, 
  Truck,
  Calculator,
  ShoppingCart,
  Loader2,
  CheckCircle,
  AlertCircle
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Form, FormItem, FormLabel, FormControl, FormMessage, FormField } from "@/components/ui/form";
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogTrigger,
  DialogFooter
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { useCreateOrder, useUsers, useProducts, useAvailableCoupons } from "@/hooks/useApi";
import { usersApi } from "@/lib/api";

const CreateOrder = () => {
  const navigate = useNavigate();
  const [currentStep, setCurrentStep] = useState(1);
  const [selectedCustomer, setSelectedCustomer] = useState(null);
  const [orderItems, setOrderItems] = useState([]);
  const [customerSearchQuery, setCustomerSearchQuery] = useState("");
  const [productSearchQuery, setProductSearchQuery] = useState("");
  const [showCustomerDialog, setShowCustomerDialog] = useState(false);
  const [showProductDialog, setShowProductDialog] = useState(false);
  const [appliedCoupon, setAppliedCoupon] = useState(null);
  const [orderSummary, setOrderSummary] = useState({
    subtotal: 0,
    discount: 0,
    tax: 0,
    shipping: 0,
    total: 0
  });
  const [isCreatingCustomer, setIsCreatingCustomer] = useState(false);
  const [newCustomerData, setNewCustomerData] = useState({
    name: "",
    email: "",
    phone: ""
  });

  // API hooks
  const { data: usersResponse } = useUsers({ 
    search: customerSearchQuery,
    per_page: 10 
  });
  const { data: productsResponse } = useProducts({ 
    search: productSearchQuery,
    per_page: 20 
  });
  const { data: couponsResponse } = useAvailableCoupons();
  const createOrderMutation = useCreateOrder();

  // Form setup
  const methods = useForm({
    defaultValues: {
      notes: "",
      shipping_method: "standard",
      payment_method: "cash",
      shipping_address: {
        street: "",
        city: "",
        state: "",
        postal_code: "",
        country: "US"
      },
      coupon_code: ""
    }
  });

  const users = usersResponse?.data || [];
  const products = productsResponse?.data || [];
  const coupons = couponsResponse?.data || [];

  // Calculate order totals
  const calculateTotals = useCallback(() => {
    const subtotal = orderItems.reduce((sum, item) => sum + (item.price * item.quantity), 0);
    const discount = appliedCoupon ? 
      (appliedCoupon.type === 'percentage' ? 
        (subtotal * appliedCoupon.value / 100) : 
        Math.min(appliedCoupon.value, subtotal)
      ) : 0;
    const tax = (subtotal - discount) * 0.08; // 8% tax rate
    const shipping = subtotal > 100 ? 0 : 15; // Free shipping over $100
    const total = subtotal - discount + tax + shipping;

    setOrderSummary({
      subtotal: subtotal,
      discount: discount,
      tax: tax,
      shipping: shipping,
      total: total
    });
  }, [orderItems, appliedCoupon]);

  useEffect(() => {
    calculateTotals();
  }, [calculateTotals]);

  // Add product to order
  const addProductToOrder = (product) => {
    // Check if product has sufficient stock
    if (product.stock <= 0) {
      toast.error(`${product.name} is out of stock`);
      return;
    }

    const existingItem = orderItems.find(item => item.product_id === product.id);
    
    if (existingItem) {
      // Check if we can add more (don't exceed stock)
      if (existingItem.quantity >= product.stock) {
        toast.error(`Cannot add more ${product.name} - insufficient stock`);
        return;
      }
      
      setOrderItems(items => 
        items.map(item => 
          item.product_id === product.id 
            ? { ...item, quantity: item.quantity + 1 }
            : item
        )
      );
    } else {
      setOrderItems(items => [...items, {
        product_id: product.id,
        name: product.name,
        price: parseFloat(product.price) || 0,
        quantity: 1,
        image: product.images?.[0] || null,
        sku: product.sku,
        max_stock: product.stock
      }]);
    }
    setShowProductDialog(false);
    toast.success(`${product.name} added to order`);
  };

  // Update item quantity
  const updateItemQuantity = (productId, newQuantity) => {
    if (newQuantity <= 0) {
      removeItemFromOrder(productId);
      return;
    }

    const item = orderItems.find(item => item.product_id === productId);
    if (item && item.max_stock && newQuantity > item.max_stock) {
      toast.error(`Cannot set quantity to ${newQuantity} - only ${item.max_stock} in stock`);
      return;
    }
    
    setOrderItems(items =>
      items.map(item =>
        item.product_id === productId
          ? { ...item, quantity: newQuantity }
          : item
      )
    );
  };

  // Remove item from order
  const removeItemFromOrder = (productId) => {
    setOrderItems(items => items.filter(item => item.product_id !== productId));
  };

  // Apply coupon
  const applyCoupon = (coupon) => {
    setAppliedCoupon(coupon);
    toast.success(`Coupon "${coupon.code}" applied`);
  };

  // Remove coupon
  const removeCoupon = () => {
    setAppliedCoupon(null);
    toast.success("Coupon removed");
  };

  // Create new customer
  const handleCreateCustomer = async () => {
    if (!newCustomerData.name.trim()) {
      toast.error("Customer name is required");
      return;
    }
    
    if (!newCustomerData.email.trim()) {
      toast.error("Customer email is required");
      return;
    }

    // Simple email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(newCustomerData.email)) {
      toast.error("Please enter a valid email address");
      return;
    }

    try {
      setIsCreatingCustomer(true);
      
      const customerData = {
        name: newCustomerData.name.trim(),
        email: newCustomerData.email.trim(),
        phone: newCustomerData.phone.trim() || null,
        role: 'customer'
      };

      const response = await usersApi.create(customerData);
      
      if (response.success && response.data) {
        setSelectedCustomer(response.data);
        setShowCustomerDialog(false);
        setNewCustomerData({ name: "", email: "", phone: "" });
        toast.success(`Customer "${response.data.name}" created successfully`);
      } else {
        throw new Error(response.message || "Failed to create customer");
      }
    } catch (error) {
      console.error("Customer creation error:", error);
      toast.error(error.message || "Failed to create customer");
    } finally {
      setIsCreatingCustomer(false);
    }
  };

  // Reset new customer form
  const resetCustomerForm = () => {
    setNewCustomerData({ name: "", email: "", phone: "" });
    setIsCreatingCustomer(false);
  };

  // Submit order
  const onSubmit = async (data) => {
    if (!selectedCustomer) {
      toast.error("Please select a customer");
      return;
    }
    
    if (orderItems.length === 0) {
      toast.error("Please add at least one product to the order");
      return;
    }

    try {
      const orderData = {
        customer_id: selectedCustomer.id,
        items: orderItems.map(item => ({
          product_id: item.product_id,
          quantity: item.quantity,
          price: item.price
        })),
        shipping_address: data.shipping_address,
        shipping_method: data.shipping_method,
        payment_method: data.payment_method,
        coupon_code: appliedCoupon?.code || null,
        notes: data.notes,
        subtotal: orderSummary.subtotal,
        discount_amount: orderSummary.discount,
        tax_amount: orderSummary.tax,
        shipping_amount: orderSummary.shipping,
        total_amount: orderSummary.total
      };

      await createOrderMutation.mutateAsync(orderData);
      toast.success("Order created successfully!");
      navigate("/orders");
    } catch (error) {
      toast.error(error.message || "Failed to create order");
    }
  };

  const steps = [
    { id: 1, title: "Customer", icon: User, description: "Select customer" },
    { id: 2, title: "Products", icon: Package, description: "Add products" },
    { id: 3, title: "Details", icon: CreditCard, description: "Shipping & payment" },
    { id: 4, title: "Review", icon: CheckCircle, description: "Confirm order" }
  ];

  return (
    <div className="space-y-6 p-3 sm:p-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/orders")}
          className="h-9 w-9"
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Create Custom Order</h1>
          <p className="text-muted-foreground">
            Manually create an order for a customer
          </p>
        </div>
      </div>

      {/* Progress Steps */}
      <div className="flex items-center justify-between mb-8">
        {steps.map((step, index) => (
          <div key={step.id} className="flex items-center">
            <div 
              className={`flex items-center justify-center w-10 h-10 rounded-full border-2 transition-colors ${
                currentStep >= step.id
                  ? 'bg-primary text-primary-foreground border-primary'
                  : 'border-muted-foreground text-muted-foreground'
              }`}
            >
              <step.icon className="h-5 w-5" />
            </div>
            <div className="ml-3">
              <p className={`text-sm font-medium ${
                currentStep >= step.id ? 'text-foreground' : 'text-muted-foreground'
              }`}>
                {step.title}
              </p>
              <p className="text-xs text-muted-foreground">{step.description}</p>
            </div>
            {index < steps.length - 1 && (
              <div 
                className={`h-0.5 w-16 mx-4 ${
                  currentStep > step.id ? 'bg-primary' : 'bg-muted'
                }`} 
              />
            )}
          </div>
        ))}
      </div>

      <Form {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)} className="space-y-6">
          {/* Step 1: Customer Selection */}
          {currentStep === 1 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <User className="h-5 w-5" />
                  Select Customer
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                {selectedCustomer ? (
                  <div className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center gap-3">
                      <Avatar>
                        <AvatarImage src={selectedCustomer.avatar} />
                        <AvatarFallback>{selectedCustomer.name[0]}</AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="font-medium">{selectedCustomer.name}</p>
                        <p className="text-sm text-muted-foreground">{selectedCustomer.email}</p>
                      </div>
                    </div>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setSelectedCustomer(null)}
                    >
                      Change
                    </Button>
                  </div>
                ) : (
                  <div className="space-y-4">
                    <div className="flex items-center gap-2">
                      <div className="relative flex-1">
                        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                        <Input
                          placeholder="Search customers..."
                          value={customerSearchQuery}
                          onChange={(e) => setCustomerSearchQuery(e.target.value)}
                          className="pl-10"
                        />
                      </div>
                      <Dialog open={showCustomerDialog} onOpenChange={setShowCustomerDialog}>
                        <DialogTrigger asChild>
                          <Button variant="outline">
                            <Plus className="h-4 w-4 mr-2" />
                            New Customer
                          </Button>
                        </DialogTrigger>
                        <DialogContent>
                          <DialogHeader>
                            <DialogTitle>Create New Customer</DialogTitle>
                            <DialogDescription>
                              Add a new customer to the system
                            </DialogDescription>
                          </DialogHeader>
                          
                          <div className="space-y-4">
                            <div>
                              <label className="text-sm font-medium">Name *</label>
                              <Input
                                placeholder="Customer name"
                                value={newCustomerData.name}
                                onChange={(e) => setNewCustomerData(prev => ({
                                  ...prev,
                                  name: e.target.value
                                }))}
                                className="mt-1"
                              />
                            </div>
                            
                            <div>
                              <label className="text-sm font-medium">Email *</label>
                              <Input
                                type="email"
                                placeholder="customer@example.com"
                                value={newCustomerData.email}
                                onChange={(e) => setNewCustomerData(prev => ({
                                  ...prev,
                                  email: e.target.value
                                }))}
                                className="mt-1"
                              />
                            </div>
                            
                            <div>
                              <label className="text-sm font-medium">Phone</label>
                              <Input
                                placeholder="Phone number (optional)"
                                value={newCustomerData.phone}
                                onChange={(e) => setNewCustomerData(prev => ({
                                  ...prev,
                                  phone: e.target.value
                                }))}
                                className="mt-1"
                              />
                            </div>
                          </div>

                          <DialogFooter className="gap-2">
                            <Button
                              variant="outline"
                              onClick={() => {
                                setShowCustomerDialog(false);
                                resetCustomerForm();
                              }}
                              disabled={isCreatingCustomer}
                            >
                              Cancel
                            </Button>
                            <Button
                              onClick={handleCreateCustomer}
                              disabled={isCreatingCustomer || !newCustomerData.name.trim() || !newCustomerData.email.trim()}
                            >
                              {isCreatingCustomer ? (
                                <>
                                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                  Creating...
                                </>
                              ) : (
                                <>
                                  <Plus className="mr-2 h-4 w-4" />
                                  Create Customer
                                </>
                              )}
                            </Button>
                          </DialogFooter>
                        </DialogContent>
                      </Dialog>
                    </div>

                    {users.length > 0 && (
                      <div className="space-y-2 max-h-60 overflow-y-auto">
                        {users.map((user) => (
                          <div
                            key={user.id}
                            className="flex items-center gap-3 p-3 border rounded-lg hover:bg-muted cursor-pointer transition-colors"
                            onClick={() => setSelectedCustomer(user)}
                          >
                            <Avatar className="h-8 w-8">
                              <AvatarImage src={user.avatar} />
                              <AvatarFallback>{user.name[0]}</AvatarFallback>
                            </Avatar>
                            <div className="flex-1 min-w-0">
                              <p className="font-medium truncate">{user.name}</p>
                              <p className="text-sm text-muted-foreground truncate">{user.email}</p>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                )}

                <div className="flex justify-end">
                  <Button
                    type="button"
                    onClick={() => setCurrentStep(2)}
                    disabled={!selectedCustomer}
                  >
                    Next: Add Products
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Step 2: Product Selection */}
          {currentStep === 2 && (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-2">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Package className="h-5 w-5" />
                      Add Products
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="relative">
                      <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                      <Input
                        placeholder="Search products..."
                        value={productSearchQuery}
                        onChange={(e) => setProductSearchQuery(e.target.value)}
                        className="pl-10"
                      />
                    </div>

                    {products.length > 0 && (
                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-h-96 overflow-y-auto">
                        {products.map((product) => (
                          <div
                            key={product.id}
                            className="flex items-center gap-3 p-3 border rounded-lg hover:bg-muted cursor-pointer transition-colors"
                            onClick={() => addProductToOrder(product)}
                          >
                            <div className="w-12 h-12 bg-muted rounded-md flex items-center justify-center flex-shrink-0">
                              {product.images?.[0] ? (
                                <img
                                  src={product.images[0]}
                                  alt={product.name}
                                  className="w-full h-full object-cover rounded-md"
                                />
                              ) : (
                                <Package className="h-6 w-6 text-muted-foreground" />
                              )}
                            </div>
                            <div className="flex-1 min-w-0">
                              <p className="font-medium truncate">{product.name}</p>
                              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <span>${product.price}</span>
                                <span>•</span>
                                <span className={`flex items-center gap-1 ${
                                  product.stock <= 0 
                                    ? 'text-destructive' 
                                    : product.stock <= 10 
                                    ? 'text-orange-600' 
                                    : 'text-green-600'
                                }`}>
                                  Stock: {product.stock}
                                  {product.stock <= 0 && (
                                    <AlertCircle className="h-3 w-3" />
                                  )}
                                </span>
                              </div>
                            </div>
                            <Plus className="h-4 w-4 text-primary" />
                          </div>
                        ))}
                      </div>
                    )}
                  </CardContent>
                </Card>
              </div>

              {/* Order Items */}
              <div>
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <ShoppingCart className="h-5 w-5" />
                      Order Items ({orderItems.length})
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {orderItems.length === 0 ? (
                      <div className="text-center text-muted-foreground py-8">
                        <Package className="h-12 w-12 mx-auto mb-4 opacity-50" />
                        <p>No products added yet</p>
                        <p className="text-sm">Search and click products to add them</p>
                      </div>
                    ) : (
                      <div className="space-y-3">
                        {orderItems.map((item) => (
                          <div key={item.product_id} className="flex items-center gap-3 p-3 border rounded-lg">
                            <div className="w-10 h-10 bg-muted rounded flex items-center justify-center flex-shrink-0">
                              {item.image ? (
                                <img
                                  src={item.image}
                                  alt={item.name}
                                  className="w-full h-full object-cover rounded"
                                />
                              ) : (
                                <Package className="h-4 w-4 text-muted-foreground" />
                              )}
                            </div>
                            <div className="flex-1 min-w-0">
                              <p className="text-sm font-medium truncate">{item.name}</p>
                              <p className="text-xs text-muted-foreground">${item.price}</p>
                            </div>
                            <div className="flex items-center gap-2">
                              <Button
                                variant="outline"
                                size="icon"
                                className="h-6 w-6"
                                onClick={() => updateItemQuantity(item.product_id, item.quantity - 1)}
                              >
                                -
                              </Button>
                              <span className="text-sm w-8 text-center">{item.quantity}</span>
                              <Button
                                variant="outline"
                                size="icon"
                                className="h-6 w-6"
                                onClick={() => updateItemQuantity(item.product_id, item.quantity + 1)}
                              >
                                +
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon"
                                className="h-6 w-6 text-destructive"
                                onClick={() => removeItemFromOrder(item.product_id)}
                              >
                                <X className="h-3 w-3" />
                              </Button>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                    
                    <Separator />
                    
                    {/* Order Summary Preview */}
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span>Subtotal:</span>
                        <span>${orderSummary.subtotal.toFixed(2)}</span>
                      </div>
                      {orderSummary.discount > 0 && (
                        <div className="flex justify-between text-green-600">
                          <span>Discount:</span>
                          <span>-${orderSummary.discount.toFixed(2)}</span>
                        </div>
                      )}
                      <div className="flex justify-between">
                        <span>Tax:</span>
                        <span>${orderSummary.tax.toFixed(2)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Shipping:</span>
                        <span>${orderSummary.shipping.toFixed(2)}</span>
                      </div>
                      <Separator />
                      <div className="flex justify-between font-medium">
                        <span>Total:</span>
                        <span>${orderSummary.total.toFixed(2)}</span>
                      </div>
                    </div>

                    <div className="flex gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setCurrentStep(1)}
                      >
                        Back
                      </Button>
                      <Button
                        size="sm"
                        onClick={() => setCurrentStep(3)}
                        disabled={orderItems.length === 0}
                        className="flex-1"
                      >
                        Next: Details
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          )}

          {/* Step 3: Shipping & Payment Details */}
          {currentStep === 3 && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="space-y-6">
                {/* Shipping Address */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Truck className="h-5 w-5" />
                      Shipping Address
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <FormField
                      name="shipping_address.street"
                      control={methods.control}
                      rules={{ required: "Street address is required" }}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Street Address *</FormLabel>
                          <FormControl>
                            <Input {...field} placeholder="123 Main Street" />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <div className="grid grid-cols-2 gap-4">
                      <FormField
                        name="shipping_address.city"
                        control={methods.control}
                        rules={{ required: "City is required" }}
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>City *</FormLabel>
                            <FormControl>
                              <Input {...field} placeholder="New York" />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        name="shipping_address.state"
                        control={methods.control}
                        rules={{ required: "State is required" }}
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>State *</FormLabel>
                            <FormControl>
                              <Input {...field} placeholder="NY" />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <FormField
                        name="shipping_address.postal_code"
                        control={methods.control}
                        rules={{ required: "Postal code is required" }}
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Postal Code *</FormLabel>
                            <FormControl>
                              <Input {...field} placeholder="10001" />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        name="shipping_address.country"
                        control={methods.control}
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Country</FormLabel>
                            <FormControl>
                              <Select onValueChange={field.onChange} value={field.value}>
                                <SelectTrigger>
                                  <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                  <SelectItem value="US">United States</SelectItem>
                                  <SelectItem value="CA">Canada</SelectItem>
                                  <SelectItem value="MX">Mexico</SelectItem>
                                </SelectContent>
                              </Select>
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>
                  </CardContent>
                </Card>

                {/* Shipping Method */}
                <Card>
                  <CardHeader>
                    <CardTitle>Shipping Method</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <FormField
                      name="shipping_method"
                      control={methods.control}
                      render={({ field }) => (
                        <FormItem>
                          <FormControl>
                            <div className="space-y-3">
                              {[
                                { value: 'standard', label: 'Standard Shipping', price: 15, time: '5-7 business days' },
                                { value: 'express', label: 'Express Shipping', price: 25, time: '2-3 business days' },
                                { value: 'overnight', label: 'Overnight Shipping', price: 45, time: '1 business day' }
                              ].map((option) => (
                                <div
                                  key={option.value}
                                  className={`flex items-center justify-between p-3 border rounded-lg cursor-pointer transition-colors ${
                                    field.value === option.value ? 'border-primary bg-primary/10' : 'hover:bg-muted'
                                  }`}
                                  onClick={() => field.onChange(option.value)}
                                >
                                  <div>
                                    <p className="font-medium">{option.label}</p>
                                    <p className="text-sm text-muted-foreground">{option.time}</p>
                                  </div>
                                  <div className="text-right">
                                    <p className="font-medium">${option.price}</p>
                                    <p className="text-xs text-muted-foreground">
                                      {orderSummary.subtotal > 100 && option.value === 'standard' ? 'FREE' : ''}
                                    </p>
                                  </div>
                                </div>
                              ))}
                            </div>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </CardContent>
                </Card>

                {/* Payment Method */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <CreditCard className="h-5 w-5" />
                      Payment Method
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <FormField
                      name="payment_method"
                      control={methods.control}
                      render={({ field }) => (
                        <FormItem>
                          <FormControl>
                            <div className="space-y-3">
                              {[
                                { value: 'cash', label: 'Cash on Delivery', description: 'Pay when order is delivered' },
                                { value: 'bank_transfer', label: 'Bank Transfer', description: 'Direct bank transfer' },
                                { value: 'check', label: 'Check', description: 'Payment by check' },
                                { value: 'credit_card', label: 'Credit Card', description: 'Process card payment separately' }
                              ].map((option) => (
                                <div
                                  key={option.value}
                                  className={`flex items-center gap-3 p-3 border rounded-lg cursor-pointer transition-colors ${
                                    field.value === option.value ? 'border-primary bg-primary/10' : 'hover:bg-muted'
                                  }`}
                                  onClick={() => field.onChange(option.value)}
                                >
                                  <div className={`w-4 h-4 rounded-full border-2 ${
                                    field.value === option.value ? 'border-primary bg-primary' : 'border-muted-foreground'
                                  }`} />
                                  <div>
                                    <p className="font-medium">{option.label}</p>
                                    <p className="text-sm text-muted-foreground">{option.description}</p>
                                  </div>
                                </div>
                              ))}
                            </div>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </CardContent>
                </Card>
              </div>

              {/* Order Summary & Coupons */}
              <div className="space-y-6">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Calculator className="h-5 w-5" />
                      Order Summary
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {/* Coupon Section */}
                    <div className="space-y-3">
                      <div className="flex items-center gap-2">
                        <FormField
                          name="coupon_code"
                          control={methods.control}
                          render={({ field }) => (
                            <FormItem className="flex-1">
                              <FormControl>
                                <Input
                                  {...field}
                                  placeholder="Enter coupon code"
                                  disabled={!!appliedCoupon}
                                />
                              </FormControl>
                            </FormItem>
                          )}
                        />
                        {!appliedCoupon ? (
                          <Button
                            type="button"
                            variant="outline"
                            onClick={() => {
                              const couponCode = methods.getValues('coupon_code');
                              const coupon = coupons.find(c => c.code === couponCode);
                              if (coupon) {
                                applyCoupon(coupon);
                              } else {
                                toast.error('Invalid coupon code');
                              }
                            }}
                          >
                            Apply
                          </Button>
                        ) : (
                          <Button
                            type="button"
                            variant="outline"
                            onClick={removeCoupon}
                          >
                            Remove
                          </Button>
                        )}
                      </div>
                      
                      {appliedCoupon && (
                        <div className="flex items-center gap-2 p-2 bg-green-50 border border-green-200 rounded">
                          <CheckCircle className="h-4 w-4 text-green-600" />
                          <span className="text-sm text-green-800">
                            Coupon "{appliedCoupon.code}" applied - {appliedCoupon.type === 'percentage' ? `${appliedCoupon.value}% off` : `$${appliedCoupon.value} off`}
                          </span>
                        </div>
                      )}
                      
                      {coupons.length > 0 && (
                        <div className="space-y-2">
                          <p className="text-sm font-medium">Available Coupons:</p>
                          <div className="space-y-1">
                            {coupons.slice(0, 3).map((coupon) => (
                              <div
                                key={coupon.id}
                                className="flex items-center justify-between p-2 border rounded cursor-pointer hover:bg-muted"
                                onClick={() => {
                                  methods.setValue('coupon_code', coupon.code);
                                  applyCoupon(coupon);
                                }}
                              >
                                <div>
                                  <p className="text-sm font-medium">{coupon.code}</p>
                                  <p className="text-xs text-muted-foreground">{coupon.description}</p>
                                </div>
                                <Badge variant="secondary">
                                  {coupon.type === 'percentage' ? `${coupon.value}%` : `$${coupon.value}`}
                                </Badge>
                              </div>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>

                    <Separator />

                    {/* Price Breakdown */}
                    <div className="space-y-3">
                      <div className="flex justify-between">
                        <span>Subtotal ({orderItems.length} items):</span>
                        <span>${orderSummary.subtotal.toFixed(2)}</span>
                      </div>
                      {orderSummary.discount > 0 && (
                        <div className="flex justify-between text-green-600">
                          <span>Discount:</span>
                          <span>-${orderSummary.discount.toFixed(2)}</span>
                        </div>
                      )}
                      <div className="flex justify-between">
                        <span>Tax (8%):</span>
                        <span>${orderSummary.tax.toFixed(2)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Shipping:</span>
                        <span>
                          {orderSummary.subtotal > 100 ? (
                            <span className="text-green-600">FREE</span>
                          ) : (
                            `$${orderSummary.shipping.toFixed(2)}`
                          )}
                        </span>
                      </div>
                      <Separator />
                      <div className="flex justify-between text-lg font-semibold">
                        <span>Total:</span>
                        <span>${orderSummary.total.toFixed(2)}</span>
                      </div>
                    </div>

                    {/* Order Items List */}
                    <Separator />
                    <div className="space-y-2">
                      <p className="font-medium">Items:</p>
                      {orderItems.map((item) => (
                        <div key={item.product_id} className="flex justify-between text-sm">
                          <span>{item.name} x{item.quantity}</span>
                          <span>${(item.price * item.quantity).toFixed(2)}</span>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>

                {/* Order Notes */}
                <Card>
                  <CardHeader>
                    <CardTitle>Order Notes</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <FormField
                      name="notes"
                      control={methods.control}
                      render={({ field }) => (
                        <FormItem>
                          <FormControl>
                            <Textarea
                              {...field}
                              placeholder="Add any special instructions or notes for this order..."
                              className="min-h-[100px] resize-none"
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </CardContent>
                </Card>

                {/* Action Buttons */}
                <div className="flex gap-3">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setCurrentStep(2)}
                    className="flex-1"
                  >
                    Back to Products
                  </Button>
                  <Button
                    type="button"
                    onClick={() => setCurrentStep(4)}
                    className="flex-1"
                  >
                    Review Order
                  </Button>
                </div>
              </div>
            </div>
          )}

          {/* Step 4: Order Review & Confirmation */}
          {currentStep === 4 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5" />
                  Review & Confirm Order
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* Customer Info */}
                  <div className="space-y-4">
                    <div>
                      <h3 className="font-semibold mb-3">Customer Information</h3>
                      <div className="flex items-center gap-3 p-3 border rounded-lg">
                        <Avatar>
                          <AvatarImage src={selectedCustomer?.avatar} />
                          <AvatarFallback>{selectedCustomer?.name[0]}</AvatarFallback>
                        </Avatar>
                        <div>
                          <p className="font-medium">{selectedCustomer?.name}</p>
                          <p className="text-sm text-muted-foreground">{selectedCustomer?.email}</p>
                        </div>
                      </div>
                    </div>

                    {/* Shipping Address */}
                    <div>
                      <h3 className="font-semibold mb-3">Shipping Address</h3>
                      <div className="p-3 border rounded-lg text-sm">
                        <p>{methods.getValues('shipping_address.street')}</p>
                        <p>
                          {methods.getValues('shipping_address.city')}, {methods.getValues('shipping_address.state')} {methods.getValues('shipping_address.postal_code')}
                        </p>
                        <p>{methods.getValues('shipping_address.country')}</p>
                      </div>
                    </div>

                    {/* Shipping & Payment Methods */}
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <h4 className="font-medium mb-2">Shipping</h4>
                        <p className="text-sm text-muted-foreground capitalize">
                          {methods.getValues('shipping_method').replace('_', ' ')}
                        </p>
                      </div>
                      <div>
                        <h4 className="font-medium mb-2">Payment</h4>
                        <p className="text-sm text-muted-foreground capitalize">
                          {methods.getValues('payment_method').replace('_', ' ')}
                        </p>
                      </div>
                    </div>
                  </div>

                  {/* Order Summary */}
                  <div className="space-y-4">
                    <div>
                      <h3 className="font-semibold mb-3">Order Summary</h3>
                      <div className="space-y-3">
                        {orderItems.map((item) => (
                          <div key={item.product_id} className="flex items-center gap-3">
                            <div className="w-12 h-12 bg-muted rounded flex items-center justify-center flex-shrink-0">
                              {item.image ? (
                                <img
                                  src={item.image}
                                  alt={item.name}
                                  className="w-full h-full object-cover rounded"
                                />
                              ) : (
                                <Package className="h-4 w-4 text-muted-foreground" />
                              )}
                            </div>
                            <div className="flex-1">
                              <p className="font-medium">{item.name}</p>
                              <p className="text-sm text-muted-foreground">
                                ${item.price} x {item.quantity}
                              </p>
                            </div>
                            <p className="font-medium">${(item.price * item.quantity).toFixed(2)}</p>
                          </div>
                        ))}
                      </div>
                    </div>

                    <Separator />

                    {/* Price Summary */}
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span>Subtotal:</span>
                        <span>${orderSummary.subtotal.toFixed(2)}</span>
                      </div>
                      {orderSummary.discount > 0 && (
                        <div className="flex justify-between text-green-600">
                          <span>Discount ({appliedCoupon?.code}):</span>
                          <span>-${orderSummary.discount.toFixed(2)}</span>
                        </div>
                      )}
                      <div className="flex justify-between">
                        <span>Tax:</span>
                        <span>${orderSummary.tax.toFixed(2)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Shipping:</span>
                        <span>${orderSummary.shipping.toFixed(2)}</span>
                      </div>
                      <Separator />
                      <div className="flex justify-between text-lg font-semibold">
                        <span>Total:</span>
                        <span>${orderSummary.total.toFixed(2)}</span>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Order Notes */}
                {methods.getValues('notes') && (
                  <div>
                    <h3 className="font-semibold mb-2">Order Notes</h3>
                    <p className="text-sm text-muted-foreground p-3 border rounded-lg bg-muted/50">
                      {methods.getValues('notes')}
                    </p>
                  </div>
                )}

                {/* Final Action Buttons */}
                <div className="flex gap-3 pt-4">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setCurrentStep(3)}
                    className="flex-1"
                  >
                    Back to Details
                  </Button>
                  <Button
                    type="submit"
                    disabled={createOrderMutation.isPending}
                    className="flex-1"
                  >
                    {createOrderMutation.isPending && (
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    )}
                    Create Order
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}
        </form>
      </Form>
    </div>
  );
};

export default CreateOrder;