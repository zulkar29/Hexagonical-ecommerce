import React, { useState, useEffect } from 'react';
import {
  Search,
  Plus,
  Minus,
  X,
  ShoppingCart,
  CreditCard,
  User,
  Scan,
  Calculator,
  Receipt,
  Percent,
  DollarSign,
  Clock,
  Check,
  AlertCircle,
  Package,
  Hash,
  Users,
  Trash2,
  Edit3,
  Save,
  Phone,
  MapPin
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

const POSPage = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [cartItems, setCartItems] = useState([]);
  const [selectedCustomer, setSelectedCustomer] = useState(null);
  const [discount, setDiscount] = useState(0);
  const [discountType, setDiscountType] = useState('percentage'); // 'percentage' or 'fixed'
  const [paymentMethod, setPaymentMethod] = useState('cash');
  const [receivedAmount, setReceivedAmount] = useState('');
  const [showPaymentDialog, setShowPaymentDialog] = useState(false);
  const [showCustomerDialog, setShowCustomerDialog] = useState(false);
  const [currentTime, setCurrentTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  // Mock product data
  const products = [
    { id: 1, name: 'Basmati Rice 5kg', price: 500, stock: 50, category: 'Rice & Grains', barcode: '1234567890123' },
    { id: 2, name: 'Onion 1kg', price: 100, stock: 25, category: 'Vegetables', barcode: '1234567890124' },
    { id: 3, name: 'Milk 1L', price: 100, stock: 30, category: 'Dairy', barcode: '1234567890125' },
    { id: 4, name: 'Potato 1kg', price: 80, stock: 40, category: 'Vegetables', barcode: '1234567890126' },
    { id: 5, name: 'Cooking Oil 1L', price: 150, stock: 20, category: 'Cooking', barcode: '1234567890127' },
    { id: 6, name: 'Sugar 1kg', price: 120, stock: 5, category: 'Pantry', barcode: '1234567890128' },
    { id: 7, name: 'Tea Leaves 500g', price: 200, stock: 15, category: 'Beverages', barcode: '1234567890129' },
    { id: 8, name: 'Salt 1kg', price: 50, stock: 35, category: 'Pantry', barcode: '1234567890130' },
    { id: 9, name: 'Dal 1kg', price: 180, stock: 22, category: 'Pulses', barcode: '1234567890131' },
    { id: 10, name: 'Flour 2kg', price: 160, stock: 18, category: 'Flour', barcode: '1234567890132' }
  ];

  // Mock customers
  const customers = [
    { id: 1, name: 'Fatima Khan', phone: '01712345678', address: 'Dhanmondi, Dhaka', totalPurchases: 25000 },
    { id: 2, name: 'Ahmed Ali', phone: '01798765432', address: 'Gulshan, Dhaka', totalPurchases: 18500 },
    { id: 3, name: 'Rashida Begum', phone: '01656789012', address: 'Uttara, Dhaka', totalPurchases: 32000 },
    { id: 4, name: 'Mohammad Hasan', phone: '01534567890', address: 'Mirpur, Dhaka', totalPurchases: 15600 }
  ];

  // Filter products based on search
  const filteredProducts = products.filter(product =>
    product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    product.category.toLowerCase().includes(searchQuery.toLowerCase()) ||
    product.barcode.includes(searchQuery)
  );

  // Add product to cart
  const addToCart = (product) => {
    const existingItem = cartItems.find(item => item.id === product.id);
    if (existingItem) {
      setCartItems(cartItems.map(item =>
        item.id === product.id
          ? { ...item, quantity: Math.min(item.quantity + 1, product.stock) }
          : item
      ));
    } else {
      setCartItems([...cartItems, { ...product, quantity: 1 }]);
    }
  };

  // Update quantity
  const updateQuantity = (id, newQuantity) => {
    if (newQuantity <= 0) {
      removeFromCart(id);
    } else {
      const product = products.find(p => p.id === id);
      const validQuantity = Math.min(newQuantity, product.stock);
      setCartItems(cartItems.map(item =>
        item.id === id ? { ...item, quantity: validQuantity } : item
      ));
    }
  };

  // Remove from cart
  const removeFromCart = (id) => {
    setCartItems(cartItems.filter(item => item.id !== id));
  };

  // Clear cart
  const clearCart = () => {
    setCartItems([]);
    setSelectedCustomer(null);
    setDiscount(0);
    setReceivedAmount('');
  };

  // Calculate totals
  const subtotal = cartItems.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  const discountAmount = discountType === 'percentage' 
    ? (subtotal * discount) / 100 
    : Math.min(discount, subtotal);
  const total = subtotal - discountAmount;
  const change = receivedAmount ? Math.max(0, parseFloat(receivedAmount) - total) : 0;

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  // Payment methods
  const paymentMethods = [
    { id: 'cash', name: 'Cash', icon: DollarSign, color: 'bg-green-100 text-green-700' },
    { id: 'bkash', name: 'bKash', icon: CreditCard, color: 'bg-pink-100 text-pink-700' },
    { id: 'nagad', name: 'Nagad', icon: CreditCard, color: 'bg-orange-100 text-orange-700' },
    { id: 'rocket', name: 'Rocket', icon: CreditCard, color: 'bg-purple-100 text-purple-700' },
    { id: 'card', name: 'Card', icon: CreditCard, color: 'bg-blue-100 text-blue-700' }
  ];

  const handlePayment = () => {
    if (cartItems.length === 0) return;
    
    // Process payment logic here
    console.log('Processing payment:', {
      items: cartItems,
      customer: selectedCustomer,
      paymentMethod,
      total,
      receivedAmount,
      change
    });
    
    // Clear cart after successful payment
    clearCart();
    setShowPaymentDialog(false);
    
    // Show success message (you can implement toast here)
    alert('Payment processed successfully!');
  };

  return (
    <div className="flex h-screen bg-background">
      {/* Left Panel - Products */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <div className="p-6 border-b border-border">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-2xl font-bold text-foreground">Point of Sale</h1>
              <p className="text-sm text-muted-foreground">
                {currentTime.toLocaleDateString('en-US', { 
                  weekday: 'long', 
                  year: 'numeric', 
                  month: 'long', 
                  day: 'numeric' 
                })} • {currentTime.toLocaleTimeString('en-US', { 
                  hour12: true, 
                  hour: 'numeric', 
                  minute: '2-digit' 
                })}
              </p>
            </div>
            <div className="flex space-x-2">
              <Button variant="outline" size="sm">
                <Scan className="h-4 w-4 mr-2" />
                Scan Barcode
              </Button>
              <Button variant="outline" size="sm">
                <Calculator className="h-4 w-4 mr-2" />
                Calculator
              </Button>
            </div>
          </div>

          {/* Search */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
            <Input
              placeholder="Search products by name, category, or barcode..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
        </div>

        {/* Product Grid */}
        <div className="flex-1 p-6 overflow-y-auto">
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
            {filteredProducts.map((product) => (
              <Card 
                key={product.id} 
                className="cursor-pointer hover:shadow-md transition-shadow"
                onClick={() => addToCart(product)}
              >
                <CardContent className="p-4">
                  <div className="space-y-2">
                    <div className="flex items-start justify-between">
                      <Badge variant="secondary" className="text-xs">
                        {product.category}
                      </Badge>
                      <Badge 
                        variant={product.stock < 10 ? "destructive" : "default"}
                        className="text-xs"
                      >
                        {product.stock}
                      </Badge>
                    </div>
                    <h3 className="font-medium text-sm line-clamp-2">
                      {product.name}
                    </h3>
                    <p className="text-lg font-bold text-primary">
                      {formatCurrency(product.price)}
                    </p>
                    <Button size="sm" className="w-full">
                      <Plus className="h-3 w-3 mr-1" />
                      Add
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>

      {/* Right Panel - Cart */}
      <div className="w-96 border-l border-border flex flex-col bg-muted/30">
        {/* Cart Header */}
        <div className="p-6 border-b border-border">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">Order Summary</h2>
            <Badge variant="secondary">
              {cartItems.reduce((sum, item) => sum + item.quantity, 0)} items
            </Badge>
          </div>

          {/* Customer Selection */}
          <Dialog open={showCustomerDialog} onOpenChange={setShowCustomerDialog}>
            <DialogTrigger asChild>
              <Button variant="outline" className="w-full justify-start">
                {selectedCustomer ? (
                  <div className="flex items-center space-x-2">
                    <Avatar className="h-6 w-6">
                      <AvatarFallback className="text-xs">
                        {selectedCustomer.name.split(' ').map(n => n[0]).join('')}
                      </AvatarFallback>
                    </Avatar>
                    <span className="truncate">{selectedCustomer.name}</span>
                  </div>
                ) : (
                  <>
                    <User className="h-4 w-4 mr-2" />
                    Select Customer (Optional)
                  </>
                )}
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Select Customer</DialogTitle>
                <DialogDescription>
                  Choose a customer for this sale or proceed without one
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-3 max-h-64 overflow-y-auto">
                <Button
                  variant="outline"
                  className="w-full justify-start"
                  onClick={() => {
                    setSelectedCustomer(null);
                    setShowCustomerDialog(false);
                  }}
                >
                  <X className="h-4 w-4 mr-2" />
                  No Customer (Walk-in)
                </Button>
                {customers.map((customer) => (
                  <div
                    key={customer.id}
                    className="flex items-center justify-between p-3 border rounded-lg cursor-pointer hover:bg-muted"
                    onClick={() => {
                      setSelectedCustomer(customer);
                      setShowCustomerDialog(false);
                    }}
                  >
                    <div className="flex items-center space-x-3">
                      <Avatar className="h-8 w-8">
                        <AvatarFallback>
                          {customer.name.split(' ').map(n => n[0]).join('')}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="font-medium text-sm">{customer.name}</p>
                        <p className="text-xs text-muted-foreground">{customer.phone}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-medium">{formatCurrency(customer.totalPurchases)}</p>
                      <p className="text-xs text-muted-foreground">Total purchases</p>
                    </div>
                  </div>
                ))}
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setShowCustomerDialog(false)}>
                  Cancel
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>

        {/* Cart Items */}
        <div className="flex-1 overflow-y-auto">
          {cartItems.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-center p-6">
              <ShoppingCart className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="font-medium text-lg mb-2">Cart is empty</h3>
              <p className="text-sm text-muted-foreground">
                Add products from the left panel to start a sale
              </p>
            </div>
          ) : (
            <div className="p-4 space-y-3">
              {cartItems.map((item) => (
                <Card key={item.id}>
                  <CardContent className="p-4">
                    <div className="flex items-start justify-between mb-2">
                      <div className="flex-1 min-w-0">
                        <h4 className="font-medium text-sm truncate">{item.name}</h4>
                        <p className="text-xs text-muted-foreground">{formatCurrency(item.price)} each</p>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => removeFromCart(item.id)}
                        className="h-6 w-6 p-0 text-destructive"
                      >
                        <X className="h-3 w-3" />
                      </Button>
                    </div>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => updateQuantity(item.id, item.quantity - 1)}
                          className="h-8 w-8 p-0"
                        >
                          <Minus className="h-3 w-3" />
                        </Button>
                        <Input
                          type="number"
                          value={item.quantity}
                          onChange={(e) => updateQuantity(item.id, parseInt(e.target.value) || 0)}
                          className="w-16 h-8 text-center text-sm"
                          min="1"
                          max={item.stock}
                        />
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => updateQuantity(item.id, item.quantity + 1)}
                          className="h-8 w-8 p-0"
                          disabled={item.quantity >= item.stock}
                        >
                          <Plus className="h-3 w-3" />
                        </Button>
                      </div>
                      <p className="font-semibold text-sm">
                        {formatCurrency(item.price * item.quantity)}
                      </p>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>

        {/* Cart Footer */}
        {cartItems.length > 0 && (
          <div className="border-t border-border p-6 space-y-4">
            {/* Discount */}
            <div className="flex items-center space-x-2">
              <div className="flex-1">
                <Label htmlFor="discount" className="text-sm">Discount</Label>
                <div className="flex space-x-2">
                  <Input
                    id="discount"
                    type="number"
                    value={discount}
                    onChange={(e) => setDiscount(parseFloat(e.target.value) || 0)}
                    placeholder="0"
                    className="flex-1"
                  />
                  <Select value={discountType} onValueChange={setDiscountType}>
                    <SelectTrigger className="w-20">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="percentage">%</SelectItem>
                      <SelectItem value="fixed">৳</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </div>

            {/* Totals */}
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>Subtotal:</span>
                <span>{formatCurrency(subtotal)}</span>
              </div>
              {discountAmount > 0 && (
                <div className="flex justify-between text-sm text-green-600">
                  <span>Discount:</span>
                  <span>-{formatCurrency(discountAmount)}</span>
                </div>
              )}
              <Separator />
              <div className="flex justify-between font-semibold text-lg">
                <span>Total:</span>
                <span>{formatCurrency(total)}</span>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="space-y-2">
              <Dialog open={showPaymentDialog} onOpenChange={setShowPaymentDialog}>
                <DialogTrigger asChild>
                  <Button className="w-full" disabled={cartItems.length === 0}>
                    <CreditCard className="h-4 w-4 mr-2" />
                    Process Payment
                  </Button>
                </DialogTrigger>
                <DialogContent className="max-w-md">
                  <DialogHeader>
                    <DialogTitle>Process Payment</DialogTitle>
                    <DialogDescription>
                      Complete the sale for {formatCurrency(total)}
                    </DialogDescription>
                  </DialogHeader>
                  
                  <div className="space-y-4">
                    {/* Payment Method */}
                    <div>
                      <Label className="text-sm font-medium">Payment Method</Label>
                      <div className="grid grid-cols-2 gap-2 mt-2">
                        {paymentMethods.map((method) => (
                          <Button
                            key={method.id}
                            variant={paymentMethod === method.id ? "default" : "outline"}
                            size="sm"
                            onClick={() => setPaymentMethod(method.id)}
                            className="justify-start"
                          >
                            <method.icon className="h-4 w-4 mr-2" />
                            {method.name}
                          </Button>
                        ))}
                      </div>
                    </div>

                    {/* Amount Received (for cash) */}
                    {paymentMethod === 'cash' && (
                      <div>
                        <Label htmlFor="received" className="text-sm font-medium">
                          Amount Received
                        </Label>
                        <Input
                          id="received"
                          type="number"
                          value={receivedAmount}
                          onChange={(e) => setReceivedAmount(e.target.value)}
                          placeholder={total.toString()}
                          className="mt-1"
                        />
                        {change > 0 && (
                          <p className="text-sm text-green-600 mt-1">
                            Change: {formatCurrency(change)}
                          </p>
                        )}
                      </div>
                    )}

                    {/* Order Summary */}
                    <div className="bg-muted p-3 rounded-lg space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>Items ({cartItems.reduce((sum, item) => sum + item.quantity, 0)}):</span>
                        <span>{formatCurrency(subtotal)}</span>
                      </div>
                      {discountAmount > 0 && (
                        <div className="flex justify-between text-sm text-green-600">
                          <span>Discount:</span>
                          <span>-{formatCurrency(discountAmount)}</span>
                        </div>
                      )}
                      <Separator />
                      <div className="flex justify-between font-semibold">
                        <span>Total:</span>
                        <span>{formatCurrency(total)}</span>
                      </div>
                    </div>
                  </div>

                  <DialogFooter>
                    <Button variant="outline" onClick={() => setShowPaymentDialog(false)}>
                      Cancel
                    </Button>
                    <Button 
                      onClick={handlePayment}
                      disabled={paymentMethod === 'cash' && (!receivedAmount || parseFloat(receivedAmount) < total)}
                    >
                      <Check className="h-4 w-4 mr-2" />
                      Complete Sale
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>

              <div className="flex space-x-2">
                <Button variant="outline" onClick={clearCart} className="flex-1">
                  <Trash2 className="h-4 w-4 mr-2" />
                  Clear
                </Button>
                <Button variant="outline" className="flex-1">
                  <Receipt className="h-4 w-4 mr-2" />
                  Hold
                </Button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default POSPage;