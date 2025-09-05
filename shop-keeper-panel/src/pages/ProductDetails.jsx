import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  ArrowLeft,
  Edit3,
  Save,
  X,
  Package,
  Star,
  DollarSign,
  AlertTriangle,
  Plus,
  Minus,
  ShoppingCart,
  Eye,
  Trash2,
  Upload,
  Download,
  FileText,
  Camera,
  BarChart3,
  TrendingUp,
  TrendingDown,
  Archive,
  CheckCircle,
  Clock,
  XCircle,
  Tag,
  Warehouse,
  Building
} from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Textarea } from '@/components/ui/textarea';
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
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';

const ProductDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [isEditing, setIsEditing] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showImageUpload, setShowImageUpload] = useState(false);

  // Mock product data - in real app, fetch from API using the id
  const [product, setProduct] = useState({
    id: parseInt(id) || 1,
    name: 'Basmati Rice 5kg',
    sku: 'BR-5KG-001',
    category: 'Rice & Grains',
    subcategory: 'Basmati Rice',
    brand: 'Kalijeera',
    description: 'Premium quality aged basmati rice, perfect for biryanis and special dishes. Long grain with excellent aroma and texture.',
    price: 500,
    costPrice: 420,
    discountPrice: 475,
    currency: '৳',
    stock: 150,
    minStock: 20,
    maxStock: 500,
    unit: 'kg',
    weight: 5,
    dimensions: { length: 30, width: 20, height: 8 },
    barcode: '8901030875029',
    supplier: 'Rahman Traders',
    supplierId: 1,
    status: 'active',
    featured: true,
    tags: ['premium', 'organic', 'gluten-free'],
    images: [
      '/api/placeholder/400/400',
      '/api/placeholder/400/400',
      '/api/placeholder/400/400'
    ],
    createdDate: '2024-01-15',
    updatedDate: '2024-07-20',
    totalSales: 450,
    revenue: 225000,
    rating: 4.5,
    reviewCount: 89,
    location: 'Warehouse A - Section 2',
    expiryDate: '2025-01-15',
    batchNumber: 'BR240715'
  });

  // Mock sales history
  const salesHistory = [
    { date: '2024-07-20', quantity: 15, amount: 7500, customer: 'Fatima Khan' },
    { date: '2024-07-18', quantity: 8, amount: 4000, customer: 'Ahmed Ali' },
    { date: '2024-07-15', quantity: 12, amount: 6000, customer: 'Rashida Begum' },
    { date: '2024-07-10', quantity: 6, amount: 3000, customer: 'Mohammad Hasan' },
    { date: '2024-07-05', quantity: 20, amount: 10000, customer: 'Nasir Ahmed' }
  ];

  // Mock inventory movements
  const inventoryHistory = [
    { date: '2024-07-20', type: 'sale', quantity: -15, balance: 150, reference: 'ORD-1234' },
    { date: '2024-07-15', type: 'purchase', quantity: +50, balance: 165, reference: 'PO-5678' },
    { date: '2024-07-10', type: 'sale', quantity: -8, balance: 115, reference: 'ORD-1235' },
    { date: '2024-07-05', type: 'adjustment', quantity: +5, balance: 123, reference: 'ADJ-001' }
  ];

  const [editedProduct, setEditedProduct] = useState({ ...product });

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

  const [errors, setErrors] = useState({});

  const validateForm = () => {
    const newErrors = {};
    
    if (!editedProduct.name?.trim()) {
      newErrors.name = 'Product name is required';
    }
    
    if (!editedProduct.category?.trim()) {
      newErrors.category = 'Category is required';
    }
    
    if (!editedProduct.price || editedProduct.price <= 0) {
      newErrors.price = 'Valid selling price is required';
    }
    
    if (!editedProduct.costPrice || editedProduct.costPrice <= 0) {
      newErrors.costPrice = 'Valid cost price is required';
    }
    
    if (editedProduct.costPrice >= editedProduct.price) {
      newErrors.costPrice = 'Cost price must be less than selling price';
    }
    
    if (editedProduct.stock < 0) {
      newErrors.stock = 'Stock cannot be negative';
    }
    
    if (editedProduct.minStock < 0) {
      newErrors.minStock = 'Min stock cannot be negative';
    }
    
    if (editedProduct.maxStock < editedProduct.minStock) {
      newErrors.maxStock = 'Max stock must be greater than min stock';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSave = () => {
    if (!validateForm()) {
      return;
    }
    
    // Update timestamps
    const updatedProduct = {
      ...editedProduct,
      updatedDate: new Date().toISOString().split('T')[0]
    };
    
    setProduct(updatedProduct);
    setIsEditing(false);
    setErrors({});
    console.log('Saving product:', updatedProduct);
    
    // In real app, make API call here
    // toast.success('Product updated successfully');
  };

  const handleDelete = () => {
    console.log('Deleting product:', product.id);
    // In real app, make API call here
    navigate('/inventory/products');
    // toast.success('Product deleted successfully');
  };

  const handleCancel = () => {
    setEditedProduct({ ...product });
    setIsEditing(false);
    setErrors({});
  };

  const handleToggleFeatured = () => {
    setProduct({ ...product, featured: !product.featured });
  };

  const handleImageUpload = (event) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      // In real app, upload to server and get URL
      console.log('Uploading images:', files);
      setShowImageUpload(false);
    }
  };

  const getStockStatus = (stock, minStock) => {
    if (stock === 0) return { label: 'Out of Stock', variant: 'destructive', icon: XCircle };
    if (stock <= minStock) return { label: 'Low Stock', variant: 'default', icon: AlertTriangle };
    return { label: 'In Stock', variant: 'secondary', icon: CheckCircle };
  };

  const stockStatus = getStockStatus(product.stock, product.minStock);
  const StockIcon = stockStatus.icon;

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      <Button variant="outline" onClick={() => navigate('/inventory/products')} className="mb-4">
        <ArrowLeft className="h-4 w-4 mr-2" /> Back to Products
      </Button>
      
      {/* Header Section */}
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-3xl font-bold text-foreground flex items-center">
            {isEditing ? 'Edit Product' : product.name}
            {product.featured && !isEditing && (
              <Star className="h-6 w-6 ml-3 fill-yellow-400 text-yellow-400" />
            )}
          </h1>
          <p className="text-muted-foreground">
            SKU: {product.sku} • {product.category} • 
            <Badge variant={product.status === 'active' ? 'default' : 'secondary'} className="ml-2">
              {product.status}
            </Badge>
          </p>
        </div>
        <div className="flex gap-2">
          {!isEditing ? (
            <>
              <Button variant="outline" size="sm" onClick={handleToggleFeatured}>
                {product.featured ? (
                  <>
                    <Star className="h-4 w-4 mr-2" />
                    Featured
                  </>
                ) : (
                  <>
                    <Star className="h-4 w-4 mr-2" />
                    Mark Featured
                  </>
                )}
              </Button>
              <Dialog open={showImageUpload} onOpenChange={setShowImageUpload}>
                <DialogTrigger asChild>
                  <Button variant="outline" size="sm">
                    <Upload className="h-4 w-4 mr-2" />
                    Upload Images
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Upload Product Images</DialogTitle>
                    <DialogDescription>
                      Upload high-quality images of your product. You can upload multiple images.
                    </DialogDescription>
                  </DialogHeader>
                  <div className="space-y-4">
                    <div className="border-2 border-dashed border-muted-foreground/25 rounded-lg p-8 text-center">
                      <Camera className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
                      <div className="space-y-2">
                        <Label htmlFor="image-upload" className="cursor-pointer">
                          <span className="text-sm font-medium">Click to upload images</span>
                          <Input
                            id="image-upload"
                            type="file"
                            multiple
                            accept="image/*"
                            className="hidden"
                            onChange={handleImageUpload}
                          />
                        </Label>
                        <p className="text-xs text-muted-foreground">
                          PNG, JPG, GIF up to 10MB each
                        </p>
                      </div>
                    </div>
                  </div>
                </DialogContent>
              </Dialog>
              <Button variant="outline" size="sm" onClick={() => setIsEditing(true)}>
                <Edit3 className="h-4 w-4 mr-2" />
                Edit
              </Button>
            </>
          ) : (
            <>
              <Button variant="outline" size="sm" onClick={handleCancel}>
                <X className="h-4 w-4 mr-2" />
                Cancel
              </Button>
              <Button size="sm" onClick={handleSave}>
                <Save className="h-4 w-4 mr-2" />
                Save Changes
              </Button>
              <Button variant="destructive" size="sm" onClick={() => setShowDeleteDialog(true)}>
                <Trash2 className="h-4 w-4 mr-2" />
                Delete
              </Button>
            </>
          )}
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <StockIcon className="h-8 w-8 text-primary" />
            <div>
              <p className="text-sm text-muted-foreground">Stock</p>
              <p className="text-2xl font-bold">{product.stock}</p>
              <Badge variant={stockStatus.variant} className="mt-1 text-xs">
                {stockStatus.label}
              </Badge>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <DollarSign className="h-8 w-8 text-green-600" />
            <div>
              <p className="text-sm text-muted-foreground">Price</p>
              <p className="text-2xl font-bold">{formatCurrency(product.price)}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <ShoppingCart className="h-8 w-8 text-blue-600" />
            <div>
              <p className="text-sm text-muted-foreground">Total Sales</p>
              <p className="text-2xl font-bold">{product.totalSales}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <BarChart3 className="h-8 w-8 text-purple-600" />
            <div>
              <p className="text-sm text-muted-foreground">Revenue</p>
              <p className="text-2xl font-bold">{formatCurrency(product.revenue)}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Product Images & Information */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center gap-2 mb-4">
              <Camera className="h-5 w-5 text-primary" />
              <h3 className="text-lg font-semibold">Product Images</h3>
            </div>
            <div className="grid grid-cols-1 gap-4">
              {product.images.map((image, index) => (
                <div key={index} className="aspect-square bg-muted rounded-lg overflow-hidden">
                  <img 
                    src={image} 
                    alt={`${product.name} ${index + 1}`}
                    className="w-full h-full object-cover"
                  />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center gap-2 mb-4">
                <Package className="h-5 w-5 text-primary" />
                <h3 className="text-lg font-semibold">Product Information</h3>
              </div>
              
              {!isEditing ? (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-3">
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Name</Label>
                      <p className="font-medium">{product.name}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Brand</Label>
                      <p>{product.brand}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Category</Label>
                      <p>{product.category} → {product.subcategory}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Barcode</Label>
                      <p className="font-mono text-sm">{product.barcode}</p>
                    </div>
                  </div>
                  <div className="space-y-3">
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Supplier</Label>
                      <p>{product.supplier}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Location</Label>
                      <p>{product.location}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Weight</Label>
                      <p>{product.weight} {product.unit}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Batch Number</Label>
                      <p className="font-mono text-sm">{product.batchNumber}</p>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="name">Product Name *</Label>
                      <Input
                        id="name"
                        value={editedProduct.name}
                        onChange={(e) => setEditedProduct({ ...editedProduct, name: e.target.value })}
                        placeholder="Enter product name"
                        className={errors.name ? 'border-destructive' : ''}
                      />
                      {errors.name && <p className="text-sm text-destructive mt-1">{errors.name}</p>}
                    </div>
                    <div>
                      <Label htmlFor="brand">Brand</Label>
                      <Input
                        id="brand"
                        value={editedProduct.brand}
                        onChange={(e) => setEditedProduct({ ...editedProduct, brand: e.target.value })}
                        placeholder="Enter brand name"
                      />
                    </div>
                    <div>
                      <Label htmlFor="category">Category *</Label>
                      <Select 
                        value={editedProduct.category} 
                        onValueChange={(value) => setEditedProduct({ ...editedProduct, category: value })}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Select category" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="Rice & Grains">Rice & Grains</SelectItem>
                          <SelectItem value="Vegetables">Vegetables</SelectItem>
                          <SelectItem value="Dairy">Dairy</SelectItem>
                          <SelectItem value="Cooking">Cooking</SelectItem>
                          <SelectItem value="Pantry">Pantry</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div>
                      <Label htmlFor="barcode">Barcode</Label>
                      <Input
                        id="barcode"
                        value={editedProduct.barcode}
                        onChange={(e) => setEditedProduct({ ...editedProduct, barcode: e.target.value })}
                        placeholder="Enter barcode"
                      />
                    </div>
                  </div>
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="supplier">Supplier</Label>
                      <Input
                        id="supplier"
                        value={editedProduct.supplier}
                        onChange={(e) => setEditedProduct({ ...editedProduct, supplier: e.target.value })}
                        placeholder="Enter supplier name"
                      />
                    </div>
                    <div>
                      <Label htmlFor="location">Storage Location</Label>
                      <Input
                        id="location"
                        value={editedProduct.location}
                        onChange={(e) => setEditedProduct({ ...editedProduct, location: e.target.value })}
                        placeholder="e.g., Warehouse A - Section 2"
                      />
                    </div>
                    <div className="grid grid-cols-2 gap-2">
                      <div>
                        <Label htmlFor="weight">Weight</Label>
                        <Input
                          id="weight"
                          type="number"
                          value={editedProduct.weight}
                          onChange={(e) => setEditedProduct({ ...editedProduct, weight: parseFloat(e.target.value) || 0 })}
                          placeholder="0"
                        />
                      </div>
                      <div>
                        <Label htmlFor="unit">Unit</Label>
                        <Select 
                          value={editedProduct.unit} 
                          onValueChange={(value) => setEditedProduct({ ...editedProduct, unit: value })}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Unit" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="kg">kg</SelectItem>
                            <SelectItem value="g">g</SelectItem>
                            <SelectItem value="L">L</SelectItem>
                            <SelectItem value="ml">ml</SelectItem>
                            <SelectItem value="pcs">pcs</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>
                    <div>
                      <Label htmlFor="batch">Batch Number</Label>
                      <Input
                        id="batch"
                        value={editedProduct.batchNumber}
                        onChange={(e) => setEditedProduct({ ...editedProduct, batchNumber: e.target.value })}
                        placeholder="Enter batch number"
                      />
                    </div>
                  </div>
                </div>
              )}
              
              <div className="mt-6">
                <Label className="text-sm font-medium text-muted-foreground">Description</Label>
                {!isEditing ? (
                  <p className="mt-1 text-sm">{product.description}</p>
                ) : (
                  <Textarea
                    className="mt-1"
                    value={editedProduct.description}
                    onChange={(e) => setEditedProduct({ ...editedProduct, description: e.target.value })}
                    placeholder="Enter product description"
                    rows={3}
                  />
                )}
              </div>
              
              <div className="mt-4">
                <Label className="text-sm font-medium text-muted-foreground">Tags</Label>
                {!isEditing ? (
                  <div className="flex gap-2 mt-1">
                    {product.tags.map((tag, index) => (
                      <Badge key={index} variant="outline">{tag}</Badge>
                    ))}
                  </div>
                ) : (
                  <Input
                    className="mt-1"
                    value={editedProduct.tags.join(', ')}
                    onChange={(e) => setEditedProduct({ 
                      ...editedProduct, 
                      tags: e.target.value.split(',').map(tag => tag.trim()).filter(tag => tag) 
                    })}
                    placeholder="Enter tags separated by commas"
                  />
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Pricing & Inventory */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center gap-2 mb-4">
              <DollarSign className="h-5 w-5 text-green-600" />
              <h3 className="text-lg font-semibold">Pricing</h3>
            </div>
            
            {!isEditing ? (
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Selling Price</span>
                  <span className="font-bold text-lg">{formatCurrency(product.price)}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Cost Price</span>
                  <span className="font-medium">{formatCurrency(product.costPrice)}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Discount Price</span>
                  <span className="font-medium">{formatCurrency(product.discountPrice)}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Profit Margin</span>
                  <Badge variant="secondary">
                    {Math.round(((product.price - product.costPrice) / product.price) * 100)}%
                  </Badge>
                </div>
              </div>
            ) : (
              <div className="space-y-4">
                <div>
                  <Label htmlFor="price">Selling Price *</Label>
                  <Input
                    id="price"
                    type="number"
                    value={editedProduct.price}
                    onChange={(e) => setEditedProduct({ ...editedProduct, price: parseFloat(e.target.value) || 0 })}
                    placeholder="0.00"
                    step="0.01"
                  />
                </div>
                <div>
                  <Label htmlFor="costPrice">Cost Price *</Label>
                  <Input
                    id="costPrice"
                    type="number"
                    value={editedProduct.costPrice}
                    onChange={(e) => setEditedProduct({ ...editedProduct, costPrice: parseFloat(e.target.value) || 0 })}
                    placeholder="0.00"
                    step="0.01"
                  />
                </div>
                <div>
                  <Label htmlFor="discountPrice">Discount Price</Label>
                  <Input
                    id="discountPrice"
                    type="number"
                    value={editedProduct.discountPrice}
                    onChange={(e) => setEditedProduct({ ...editedProduct, discountPrice: parseFloat(e.target.value) || 0 })}
                    placeholder="0.00"
                    step="0.01"
                  />
                </div>
                <div className="bg-muted p-3 rounded-lg">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Profit Margin</span>
                    <Badge variant="secondary">
                      {Math.round(((editedProduct.price - editedProduct.costPrice) / editedProduct.price) * 100) || 0}%
                    </Badge>
                  </div>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center gap-2 mb-4">
              <Warehouse className="h-5 w-5 text-blue-600" />
              <h3 className="text-lg font-semibold">Inventory</h3>
            </div>
            
            {!isEditing ? (
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Current Stock</span>
                  <Badge variant={stockStatus.variant}>{product.stock} units</Badge>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Min Stock Level</span>
                  <span>{product.minStock} units</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Max Stock Level</span>
                  <span>{product.maxStock} units</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Expiry Date</span>
                  <span>{formatDate(product.expiryDate)}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Last Updated</span>
                  <span>{formatDate(product.updatedDate)}</span>
                </div>
              </div>
            ) : (
              <div className="space-y-4">
                <div>
                  <Label htmlFor="stock">Current Stock *</Label>
                  <Input
                    id="stock"
                    type="number"
                    value={editedProduct.stock}
                    onChange={(e) => setEditedProduct({ ...editedProduct, stock: parseInt(e.target.value) || 0 })}
                    placeholder="0"
                  />
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <Label htmlFor="minStock">Min Stock</Label>
                    <Input
                      id="minStock"
                      type="number"
                      value={editedProduct.minStock}
                      onChange={(e) => setEditedProduct({ ...editedProduct, minStock: parseInt(e.target.value) || 0 })}
                      placeholder="0"
                    />
                  </div>
                  <div>
                    <Label htmlFor="maxStock">Max Stock</Label>
                    <Input
                      id="maxStock"
                      type="number"
                      value={editedProduct.maxStock}
                      onChange={(e) => setEditedProduct({ ...editedProduct, maxStock: parseInt(e.target.value) || 0 })}
                      placeholder="0"
                    />
                  </div>
                </div>
                <div>
                  <Label htmlFor="expiryDate">Expiry Date</Label>
                  <Input
                    id="expiryDate"
                    type="date"
                    value={editedProduct.expiryDate}
                    onChange={(e) => setEditedProduct({ ...editedProduct, expiryDate: e.target.value })}
                  />
                </div>
                <div className="bg-muted p-3 rounded-lg">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Stock Status</span>
                    <Badge variant={getStockStatus(editedProduct.stock, editedProduct.minStock).variant}>
                      {getStockStatus(editedProduct.stock, editedProduct.minStock).label}
                    </Badge>
                  </div>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Sales History */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center gap-2 mb-4">
            <BarChart3 className="h-5 w-5 text-primary" />
            <h3 className="text-lg font-semibold">Recent Sales</h3>
          </div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Date</TableHead>
                <TableHead>Customer</TableHead>
                <TableHead className="text-center">Quantity</TableHead>
                <TableHead className="text-right">Amount</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {salesHistory.map((sale, index) => (
                <TableRow key={index}>
                  <TableCell>{formatDate(sale.date)}</TableCell>
                  <TableCell>{sale.customer}</TableCell>
                  <TableCell className="text-center">{sale.quantity} units</TableCell>
                  <TableCell className="text-right font-bold">
                    {formatCurrency(sale.amount)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center text-destructive">
              <AlertTriangle className="h-5 w-5 mr-2" />
              Delete Product
            </DialogTitle>
            <DialogDescription>
              Are you sure you want to delete "{product.name}"? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          
          <div className="bg-destructive/10 p-4 rounded-lg">
            <p className="text-sm text-destructive font-medium mb-2">This will permanently:</p>
            <ul className="text-sm text-destructive space-y-1">
              <li>• Remove the product from your inventory</li>
              <li>• Delete all product information and history</li>
              <li>• Remove from all sales records</li>
              <li>• Archive all related transactions</li>
            </ul>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowDeleteDialog(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              <Trash2 className="h-4 w-4 mr-2" />
              Delete Product
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default ProductDetails;