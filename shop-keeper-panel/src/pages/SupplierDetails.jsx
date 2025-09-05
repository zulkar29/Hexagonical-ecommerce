import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  ArrowLeft,
  Edit3,
  Save,
  X,
  Building,
  Phone,
  Mail,
  MapPin,
  Calendar,
  DollarSign,
  Package,
  Star,
  StarOff,
  TrendingUp,
  TrendingDown,
  AlertTriangle,
  CheckCircle,
  Plus,
  Minus,
  History,
  ShoppingCart,
  Users,
  Eye,
  EyeOff,
  Copy,
  Trash2,
  MoreVertical,
  RefreshCw,
  Hash,
  CreditCard,
  Clock,
  Archive,
  FileText,
  Printer,
  Send,
  Download,
  Upload,
  Globe,
  Truck,
  Shield,
  Award,
  Briefcase,
  Factory,
  Store,
  User,
  UserCheck,
  MessageSquare,
  ExternalLink
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Textarea } from '@/components/ui/textarea';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';
import { Progress } from '@/components/ui/progress';
import { Switch } from '@/components/ui/switch';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';

const CompleteSupplierDetailsPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [isEditing, setIsEditing] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showContactDialog, setShowContactDialog] = useState(false);
  const [contactMessage, setContactMessage] = useState('');

  // Mock supplier data - in real app, fetch from API using the id
  const [supplier, setSupplier] = useState({
    id: parseInt(id) || 1,
    name: 'Rahman Traders',
    email: 'contact@rahmantraders.com',
    phone: '01712345678',
    alternatePhone: '02-9876543',
    website: 'www.rahmantraders.com',
    address: 'House 12, Road 5, Dhanmondi, Dhaka-1205',
    city: 'Dhaka',
    state: 'Dhaka Division',
    country: 'Bangladesh',
    postalCode: '1205',
    contactPerson: 'Abdul Rahman',
    designation: 'Sales Manager',
    businessType: 'Wholesale',
    registrationNumber: 'REG-123456789',
    taxId: 'BIN-123456789',
    bankName: 'Dutch Bangla Bank',
    accountNumber: '1234567890123',
    routingNumber: '090260001',
    swiftCode: 'DBBLBDDHXXX',
    paymentTerms: '30 days',
    creditLimit: 500000,
    leadTime: '5-7 days',
    minimumOrder: 50000,
    rating: 4.5,
    status: 'active',
    isPreferred: true,
    joinedDate: '2023-01-15',
    lastContact: '2024-07-18',
    lastOrderDate: '2024-07-20',
    productsCount: 12,
    activeProducts: 10,
    totalOrders: 45,
    totalValue: 250000,
    averageOrderValue: 5556,
    onTimeDelivery: 92,
    qualityRating: 4.3,
    communicationRating: 4.7,
    image: '/api/placeholder/80/80',
    tags: ['reliable', 'wholesale', 'food-items'],
    notes: 'Excellent supplier with consistent quality and timely deliveries. Primary contact for rice and grain products.',
    documents: [
      { id: 1, name: 'Trade License', type: 'PDF', uploadDate: '2023-01-15', status: 'verified' },
      { id: 2, name: 'Tax Certificate', type: 'PDF', uploadDate: '2023-01-16', status: 'verified' },
      { id: 3, name: 'Bank Statement', type: 'PDF', uploadDate: '2024-01-10', status: 'pending' }
    ]
  });

  // Mock supplier products
  const supplierProducts = [
    { id: 1, name: 'Basmati Rice 5kg', category: 'Rice & Grains', price: 500, stock: 50, lastOrder: '2024-07-20', orderCount: 15 },
    { id: 2, name: 'Jasmine Rice 2kg', category: 'Rice & Grains', price: 280, stock: 30, lastOrder: '2024-07-18', orderCount: 8 },
    { id: 3, name: 'Brown Rice 1kg', category: 'Rice & Grains', price: 150, stock: 25, lastOrder: '2024-07-15', orderCount: 12 },
    { id: 4, name: 'Rice Flour 500g', category: 'Flour', price: 80, stock: 40, lastOrder: '2024-07-10', orderCount: 6 }
  ];

  // Mock order history
  const orderHistory = [
    { id: 'ORD-001', date: '2024-07-20', products: 3, amount: 15000, status: 'delivered', paymentStatus: 'paid' },
    { id: 'ORD-002', date: '2024-07-15', products: 2, amount: 8500, status: 'delivered', paymentStatus: 'paid' },
    { id: 'ORD-003', date: '2024-07-10', products: 4, amount: 12000, status: 'delivered', paymentStatus: 'pending' },
    { id: 'ORD-004', date: '2024-07-05', products: 1, amount: 5000, status: 'delivered', paymentStatus: 'paid' },
    { id: 'ORD-005', date: '2024-06-30', products: 5, amount: 18000, status: 'delivered', paymentStatus: 'paid' }
  ];

  // Mock communication history
  const communicationHistory = [
    { id: 1, type: 'email', subject: 'New product inquiry', date: '2024-07-18', status: 'sent' },
    { id: 2, type: 'phone', subject: 'Order confirmation call', date: '2024-07-15', status: 'completed' },
    { id: 3, type: 'meeting', subject: 'Quarterly business review', date: '2024-07-01', status: 'completed' },
    { id: 4, type: 'email', subject: 'Payment reminder', date: '2024-06-25', status: 'sent' }
  ];

  const [editedSupplier, setEditedSupplier] = useState({ ...supplier });

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

  const handleSave = () => {
    // In real app, make API call to update supplier
    setSupplier({ ...editedSupplier });
    setIsEditing(false);
    console.log('Saving supplier:', editedSupplier);
  };

  const handleDelete = () => {
    // In real app, make API call to delete supplier
    console.log('Deleting supplier:', supplier.id);
    navigate('/inventory/suppliers');
  };

  const handleTogglePreferred = () => {
    setSupplier({ ...supplier, isPreferred: !supplier.isPreferred });
    console.log('Toggled preferred status:', !supplier.isPreferred);
  };

  const handleSendMessage = () => {
    console.log('Sending message:', contactMessage);
    setContactMessage('');
    setShowContactDialog(false);
    // Show success message
  };

  const renderStars = (rating) => {
    return Array.from({ length: 5 }, (_, i) => (
      <Star
        key={i}
        className={cn(
          "h-4 w-4",
          i < Math.floor(rating) ? "fill-yellow-400 text-yellow-400" : "text-gray-300"
        )}
      />
    ));
  };

  const getBusinessTypeIcon = (type) => {
    switch (type.toLowerCase()) {
      case 'wholesale': return Store;
      case 'manufacturer': return Factory;
      case 'distributor': return Truck;
      case 'retail': return Store;
      default: return Building;
    }
  };

  const BusinessTypeIcon = getBusinessTypeIcon(supplier.businessType);

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto p-6 space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Button 
              variant="ghost" 
              size="sm" 
              onClick={() => navigate('/inventory/suppliers')}
              className="p-2"
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div className="flex items-center space-x-4">
              <Avatar className="h-16 w-16">
                <AvatarImage src={supplier.image} alt={supplier.name} />
                <AvatarFallback>
                  <Building className="h-6 w-6" />
                </AvatarFallback>
              </Avatar>
              <div>
                <h1 className="text-2xl font-bold flex items-center">
                  {isEditing ? 'Edit Supplier' : supplier.name}
                  {supplier.isPreferred && !isEditing && (
                    <Star className="h-5 w-5 ml-2 fill-yellow-400 text-yellow-400" />
                  )}
                </h1>
                <div className="flex items-center space-x-4 text-sm text-muted-foreground mt-1">
                  <span className="flex items-center">
                    <BusinessTypeIcon className="h-3 w-3 mr-1" />
                    {supplier.businessType}
                  </span>
                  <span>•</span>
                  <span>Member since {formatDate(supplier.joinedDate)}</span>
                  <span>•</span>
                  <Badge variant={supplier.status === 'active' ? 'default' : 'secondary'}>
                    {supplier.status}
                  </Badge>
                </div>
              </div>
            </div>
          </div>

          <div className="flex items-center space-x-2">
            {!isEditing ? (
              <>
                <Button variant="outline" onClick={handleTogglePreferred}>
                  {supplier.isPreferred ? (
                    <>
                      <StarOff className="h-4 w-4 mr-2" />
                      Remove Preferred
                    </>
                  ) : (
                    <>
                      <Star className="h-4 w-4 mr-2" />
                      Mark Preferred
                    </>
                  )}
                </Button>
                <Dialog open={showContactDialog} onOpenChange={setShowContactDialog}>
                  <DialogTrigger asChild>
                    <Button variant="outline">
                      <MessageSquare className="h-4 w-4 mr-2" />
                      Contact
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Contact {supplier.name}</DialogTitle>
                      <DialogDescription>
                        Send a message to {supplier.contactPerson}
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                      <div>
                        <Label htmlFor="minimumOrder">Minimum Order</Label>
                        <Input
                          id="minimumOrder"
                          type="number"
                          value={editedSupplier.minimumOrder}
                          onChange={(e) => setEditedSupplier({ ...editedSupplier, minimumOrder: parseFloat(e.target.value) })}
                          placeholder="Minimum order amount"
                        />
                      </div>
                    </div>
                  </DialogContent>
                </Dialog>
              </>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Business Registration</Label>
                    <div className="space-y-2">
                      <div className="flex items-center space-x-2">
                        <Hash className="h-4 w-4 text-muted-foreground" />
                        <span className="text-sm font-mono">{supplier.registrationNumber}</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <FileText className="h-4 w-4 text-muted-foreground" />
                        <span className="text-sm font-mono">{supplier.taxId}</span>
                      </div>
                    </div>
                  </div>

                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Banking Details</Label>
                    <div className="space-y-2">
                      <div className="flex items-center space-x-2">
                        <Building className="h-4 w-4 text-muted-foreground" />
                        <span className="text-sm">{supplier.bankName}</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <CreditCard className="h-4 w-4 text-muted-foreground" />
                        <span className="text-sm font-mono">{supplier.accountNumber}</span>
                        <Button variant="ghost" size="sm">
                          <Copy className="h-3 w-3" />
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>
                
                <div className="space-y-4">
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Terms & Conditions</Label>
                    <div className="space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="text-sm">Payment Terms:</span>
                        <Badge variant="outline">{supplier.paymentTerms}</Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm">Credit Limit:</span>
                        <span className="text-sm font-medium">{formatCurrency(supplier.creditLimit)}</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm">Lead Time:</span>
                        <Badge variant="secondary">{supplier.leadTime}</Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm">Minimum Order:</span>
                        <span className="text-sm font-medium">{formatCurrency(supplier.minimumOrder)}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>

      {/* Notes */}
      <Card>
      <CardHeader>
        <CardTitle className="flex items-center">
          <FileText className="h-5 w-5 mr-2" />
          Notes & Comments
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isEditing ? (
          <Textarea
            value={editedSupplier.notes}
            onChange={(e) => setEditedSupplier({ ...editedSupplier, notes: e.target.value })}
            placeholder="Add notes about this supplier..."
            rows={4}
          />
        ) : (
          <p className="text-sm">{supplier.notes}</p>
        )}
      </CardContent>
    </Card>

    {/* Tabs for detailed information */}
    {!isEditing && (
      <Tabs defaultValue="products" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="products">Products</TabsTrigger>
          <TabsTrigger value="orders">Order History</TabsTrigger>
          <TabsTrigger value="communications">Communications</TabsTrigger>
          <TabsTrigger value="documents">Documents</TabsTrigger>
        </TabsList>
        
        <TabsContent value="products" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <span className="flex items-center">
                  <Package className="h-5 w-5 mr-2" />
                  Supplier Products ({supplierProducts.length})
                </span>
                <Button size="sm" variant="outline">
                  <Plus className="h-4 w-4 mr-2" />
                  Add Product
                </Button>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Product</TableHead>
                      <TableHead>Category</TableHead>
                      <TableHead>Price</TableHead>
                      <TableHead>Stock</TableHead>
                      <TableHead>Orders</TableHead>
                      <TableHead>Last Order</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {supplierProducts.map((product) => (
                      <TableRow key={product.id}>
                        <TableCell>
                          <div className="font-medium">{product.name}</div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="secondary">{product.category}</Badge>
                        </TableCell>
                        <TableCell>{formatCurrency(product.price)}</TableCell>
                        <TableCell>
                          <Badge variant={product.stock < 10 ? 'destructive' : 'default'}>
                            {product.stock}
                          </Badge>
                        </TableCell>
                        <TableCell>{product.orderCount}</TableCell>
                        <TableCell>{formatDate(product.lastOrder)}</TableCell>
                        <TableCell>
                          <div className="flex items-center space-x-1">
                            <Button size="sm" variant="ghost">
                              <Eye className="h-4 w-4" />
                            </Button>
                            <Button size="sm" variant="ghost">
                              <Edit3 className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="orders" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <ShoppingCart className="h-5 w-5 mr-2" />
                Order History ({orderHistory.length})
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Order ID</TableHead>
                      <TableHead>Date</TableHead>
                      <TableHead>Products</TableHead>
                      <TableHead>Amount</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Payment</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {orderHistory.map((order) => (
                      <TableRow key={order.id}>
                        <TableCell>
                          <span className="font-mono text-sm">{order.id}</span>
                        </TableCell>
                        <TableCell>{formatDate(order.date)}</TableCell>
                        <TableCell>{order.products} items</TableCell>
                        <TableCell>{formatCurrency(order.amount)}</TableCell>
                        <TableCell>
                          <Badge variant={
                            order.status === 'delivered' ? 'default' :
                            order.status === 'pending' ? 'secondary' : 'destructive'
                          }>
                            {order.status}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <Badge variant={
                            order.paymentStatus === 'paid' ? 'default' : 'destructive'
                          }>
                            {order.paymentStatus}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center space-x-1">
                            <Button size="sm" variant="ghost">
                              <Eye className="h-4 w-4" />
                            </Button>
                            <Button size="sm" variant="ghost">
                              <Download className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="communications" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <span className="flex items-center">
                  <MessageSquare className="h-5 w-5 mr-2" />
                  Communication History
                </span>
                <Button size="sm" onClick={() => setShowContactDialog(true)}>
                  <Plus className="h-4 w-4 mr-2" />
                  New Message
                </Button>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {communicationHistory.map((comm) => (
                  <div key={comm.id} className="flex items-start space-x-3 p-3 border rounded-lg">
                    <div className={cn(
                      "p-2 rounded-full",
                      comm.type === 'email' && "bg-blue-100 text-blue-600",
                      comm.type === 'phone' && "bg-green-100 text-green-600",
                      comm.type === 'meeting' && "bg-purple-100 text-purple-600"
                    )}>
                      {comm.type === 'email' && <Mail className="h-4 w-4" />}
                      {comm.type === 'phone' && <Phone className="h-4 w-4" />}
                      {comm.type === 'meeting' && <Users className="h-4 w-4" />}
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <p className="font-medium text-sm">{comm.subject}</p>
                        <span className="text-xs text-muted-foreground">{formatDate(comm.date)}</span>
                      </div>
                      <div className="flex items-center space-x-2 mt-1">
                        <Badge variant="outline" className="text-xs capitalize">
                          {comm.type}
                        </Badge>
                        <Badge 
                          variant={comm.status === 'completed' ? 'default' : 'secondary'}
                          className="text-xs"
                        >
                          {comm.status}
                        </Badge>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="documents" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <span className="flex items-center">
                  <FileText className="h-5 w-5 mr-2" />
                  Documents ({supplier.documents.length})
                </span>
                <Button size="sm" variant="outline">
                  <Upload className="h-4 w-4 mr-2" />
                  Upload Document
                </Button>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {supplier.documents.map((doc) => (
                  <div key={doc.id} className="flex items-center justify-between p-3 border rounded-lg">
                    <div className="flex items-center space-x-3">
                      <FileText className="h-8 w-8 text-muted-foreground" />
                      <div>
                        <p className="font-medium text-sm">{doc.name}</p>
                        <p className="text-xs text-muted-foreground">
                          {doc.type} • Uploaded {formatDate(doc.uploadDate)}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Badge 
                        variant={doc.status === 'verified' ? 'default' : 'secondary'}
                        className="text-xs"
                      >
                        {doc.status === 'verified' ? (
                          <>
                            <CheckCircle className="h-3 w-3 mr-1" />
                            Verified
                          </>
                        ) : (
                          <>
                            <Clock className="h-3 w-3 mr-1" />
                            Pending
                          </>
                        )}
                      </Badge>
                      <Button size="sm" variant="ghost">
                        <Download className="h-4 w-4" />
                      </Button>
                      <Button size="sm" variant="ghost">
                        <Eye className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    )}
  </div>

  {/* Sidebar */}
  <div className="space-y-6">
    {/* Performance Metrics */}
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Performance Metrics</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Overall Rating</span>
            <div className="flex items-center space-x-2">
              <div className="flex">
                {renderStars(supplier.rating)}
              </div>
              <span className="text-sm font-medium">{supplier.rating}</span>
            </div>
          </div>
          
          <div>
            <div className="flex items-center justify-between mb-1">
              <span className="text-sm text-muted-foreground">Quality Rating</span>
              <span className="text-sm font-medium">{supplier.qualityRating}/5</span>
            </div>
            <Progress value={(supplier.qualityRating / 5) * 100} className="h-2" />
          </div>

          <div>
            <div className="flex items-center justify-between mb-1">
              <span className="text-sm text-muted-foreground">Communication</span>
              <span className="text-sm font-medium">{supplier.communicationRating}/5</span>
            </div>
            <Progress value={(supplier.communicationRating / 5) * 100} className="h-2" />
          </div>

          <div>
            <div className="flex items-center justify-between mb-1">
              <span className="text-sm text-muted-foreground">On-time Delivery</span>
              <span className="text-sm font-medium">{supplier.onTimeDelivery}%</span>
            </div>
            <Progress value={supplier.onTimeDelivery} className="h-2" />
          </div>
        </div>
      </CardContent>
    </Card>

    {/* Quick Stats */}
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Quick Stats</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div className="text-center">
            <p className="text-2xl font-bold text-primary">{supplier.productsCount}</p>
            <p className="text-xs text-muted-foreground">Products</p>
          </div>
          <div className="text-center">
            <p className="text-2xl font-bold text-green-600">{supplier.totalOrders}</p>
            <p className="text-xs text-muted-foreground">Orders</p>
          </div>
        </div>

        <Separator />

        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Total Business</span>
            <span className="text-sm font-bold">{formatCurrency(supplier.totalValue)}</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Avg Order Value</span>
            <span className="text-sm font-medium">{formatCurrency(supplier.averageOrderValue)}</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Last Order</span>
            <span className="text-sm">{formatDate(supplier.lastOrderDate)}</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Last Contact</span>
            <span className="text-sm">{formatDate(supplier.lastContact)}</span>
          </div>
        </div>
      </CardContent>
    </Card>

    {/* Quick Actions */}
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Quick Actions</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <Button variant="outline" className="w-full justify-start">
          <Phone className="h-4 w-4 mr-2" />
          Call Supplier
        </Button>
        <Button variant="outline" className="w-full justify-start">
          <Mail className="h-4 w-4 mr-2" />
          Send Email
        </Button>
        <Button variant="outline" className="w-full justify-start">
          <Plus className="h-4 w-4 mr-2" />
          New Order
        </Button>
        <Button variant="outline" className="w-full justify-start">
          <FileText className="h-4 w-4 mr-2" />
          Generate Report
        </Button>
      </CardContent>
    </Card>

    {/* Activity Timeline */}
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Recent Activity</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div className="flex items-start space-x-3">
            <div className="w-2 h-2 bg-green-500 rounded-full mt-2"></div>
            <div>
              <p className="text-sm">Order delivered successfully</p>
              <p className="text-xs text-muted-foreground">2 days ago</p>
            </div>
          </div>
          <div className="flex items-start space-x-3">
            <div className="w-2 h-2 bg-blue-500 rounded-full mt-2"></div>
            <div>
              <p className="text-sm">Payment received</p>
              <p className="text-xs text-muted-foreground">4 days ago</p>
            </div>
          </div>
          <div className="flex items-start space-x-3">
            <div className="w-2 h-2 bg-yellow-500 rounded-full mt-2"></div>
            <div>
              <p className="text-sm">New order placed</p>
              <p className="text-xs text-muted-foreground">1 week ago</p>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>

  {/* Delete Confirmation Dialog */}
  <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
    <DialogContent className="max-w-md">
      <DialogHeader>
        <DialogTitle className="flex items-center text-destructive">
          <AlertTriangle className="h-5 w-5 mr-2" />
          Delete Supplier
        </DialogTitle>
        <DialogDescription>
          Are you sure you want to delete "{supplier.name}"? This action cannot be undone.
        </DialogDescription>
      </DialogHeader>
      <div className="bg-destructive/10 p-4 rounded-lg">
        <p className="text-sm text-destructive font-medium mb-2">This will permanently:</p>
        <ul className="text-sm text-destructive space-y-1">
          <li>• Remove the supplier from your database</li>
          <li>• Delete all contact information and history</li>
          <li>• Remove supplier associations from products</li>
          <li>• Archive all related orders and communications</li>
        </ul>
      </div>
      <DialogFooter>
        <Button variant="outline" onClick={() => setShowDeleteDialog(false)}>
          Cancel
        </Button>
        <Button variant="destructive" onClick={handleDelete}>
          <Trash2 className="h-4 w-4 mr-2" />
          Delete Supplier
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</div>
);
};

export default CompleteSupplierDetailsPage;