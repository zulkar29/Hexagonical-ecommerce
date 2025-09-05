import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  ArrowLeft,
  Save,
  Upload,
  Camera,
  Package,
  DollarSign,
  Warehouse,
  Tag,
  AlertCircle,
  Plus,
  X
} from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';

const ProductCreate = () => {
  const navigate = useNavigate();
  const [errors, setErrors] = useState({});
  const [uploading, setUploading] = useState(false);

  const [product, setProduct] = useState({
    name: '',
    sku: '',
    category: '',
    subcategory: '',
    brand: '',
    description: '',
    price: '',
    costPrice: '',
    discountPrice: '',
    stock: '',
    minStock: '',
    maxStock: '',
    unit: 'kg',
    weight: '',
    barcode: '',
    supplier: '',
    location: '',
    expiryDate: '',
    batchNumber: '',
    status: 'active',
    featured: false,
    tags: [],
    images: []
  });

  const [newTag, setNewTag] = useState('');

  const validateForm = () => {
    const newErrors = {};

    if (!product.name?.trim()) {
      newErrors.name = 'Product name is required';
    }

    if (!product.category?.trim()) {
      newErrors.category = 'Category is required';
    }

    if (!product.price || parseFloat(product.price) <= 0) {
      newErrors.price = 'Valid selling price is required';
    }

    if (!product.costPrice || parseFloat(product.costPrice) <= 0) {
      newErrors.costPrice = 'Valid cost price is required';
    }

    if (parseFloat(product.costPrice) >= parseFloat(product.price)) {
      newErrors.costPrice = 'Cost price must be less than selling price';
    }

    if (product.stock && parseInt(product.stock) < 0) {
      newErrors.stock = 'Stock cannot be negative';
    }

    if (product.minStock && parseInt(product.minStock) < 0) {
      newErrors.minStock = 'Min stock cannot be negative';
    }

    if (product.maxStock && product.minStock && parseInt(product.maxStock) < parseInt(product.minStock)) {
      newErrors.maxStock = 'Max stock must be greater than min stock';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSave = async () => {
    if (!validateForm()) {
      return;
    }

    setUploading(true);
    
    try {
      // Generate SKU if not provided
      const finalProduct = {
        ...product,
        sku: product.sku || generateSKU(product.name, product.category),
        createdDate: new Date().toISOString().split('T')[0],
        updatedDate: new Date().toISOString().split('T')[0],
        id: Date.now() // Mock ID generation
      };

      console.log('Creating product:', finalProduct);
      
      // In real app, make API call here
      // await createProduct(finalProduct);
      
      // Navigate to the new product details page
      navigate(`/inventory/products/${finalProduct.id}`, {
        state: { message: 'Product created successfully!' }
      });
    } catch (error) {
      console.error('Error creating product:', error);
    } finally {
      setUploading(false);
    }
  };

  const generateSKU = (name, category) => {
    const namePrefix = name.substring(0, 3).toUpperCase();
    const categoryPrefix = category.substring(0, 3).toUpperCase();
    const randomNum = Math.floor(Math.random() * 1000).toString().padStart(3, '0');
    return `${categoryPrefix}-${namePrefix}-${randomNum}`;
  };

  const handleImageUpload = (event) => {
    const files = Array.from(event.target.files);
    if (files.length > 0) {
      // In real app, upload to server
      const newImages = files.map(file => URL.createObjectURL(file));
      setProduct({ ...product, images: [...product.images, ...newImages] });
    }
  };

  const removeImage = (index) => {
    const newImages = product.images.filter((_, i) => i !== index);
    setProduct({ ...product, images: newImages });
  };

  const addTag = () => {
    if (newTag.trim() && !product.tags.includes(newTag.trim())) {
      setProduct({ ...product, tags: [...product.tags, newTag.trim()] });
      setNewTag('');
    }
  };

  const removeTag = (tagToRemove) => {
    setProduct({ ...product, tags: product.tags.filter(tag => tag !== tagToRemove) });
  };

  const formatCurrency = (amount) => `৳${parseFloat(amount || 0).toLocaleString()}`;

  const profitMargin = product.price && product.costPrice 
    ? Math.round(((parseFloat(product.price) - parseFloat(product.costPrice)) / parseFloat(product.price)) * 100)
    : 0;

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      <Button variant="outline" onClick={() => navigate('/inventory/products')} className="mb-4">
        <ArrowLeft className="h-4 w-4 mr-2" /> Back to Products
      </Button>
      
      {/* Header Section */}
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Create New Product</h1>
          <p className="text-muted-foreground">Add a new product to your inventory</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate('/inventory/products')}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={uploading}>
            <Save className="h-4 w-4 mr-2" />
            {uploading ? 'Creating...' : 'Create Product'}
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Product Images */}
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center gap-2 mb-4">
              <Camera className="h-5 w-5 text-primary" />
              <h3 className="text-lg font-semibold">Product Images</h3>
            </div>
            
            <div className="space-y-4">
              {/* Upload Area */}
              <div className="border-2 border-dashed border-muted-foreground/25 rounded-lg p-6 text-center">
                <Camera className="mx-auto h-8 w-8 text-muted-foreground mb-2" />
                <Label htmlFor="image-upload" className="cursor-pointer">
                  <span className="text-sm font-medium text-primary hover:underline">
                    Click to upload images
                  </span>
                  <Input
                    id="image-upload"
                    type="file"
                    multiple
                    accept="image/*"
                    className="hidden"
                    onChange={handleImageUpload}
                  />
                </Label>
                <p className="text-xs text-muted-foreground mt-1">
                  PNG, JPG, GIF up to 10MB each
                </p>
              </div>

              {/* Image Preview */}
              {product.images.length > 0 && (
                <div className="grid grid-cols-2 gap-2">
                  {product.images.map((image, index) => (
                    <div key={index} className="relative group">
                      <img 
                        src={image} 
                        alt={`Product ${index + 1}`}
                        className="w-full h-24 object-cover rounded border"
                      />
                      <Button
                        size="sm"
                        variant="destructive"
                        className="absolute top-1 right-1 h-6 w-6 p-0 opacity-0 group-hover:opacity-100"
                        onClick={() => removeImage(index)}
                      >
                        <X className="h-3 w-3" />
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Basic Information */}
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center gap-2 mb-4">
                <Package className="h-5 w-5 text-primary" />
                <h3 className="text-lg font-semibold">Basic Information</h3>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="md:col-span-2">
                  <Label htmlFor="name">Product Name *</Label>
                  <Input
                    id="name"
                    value={product.name}
                    onChange={(e) => setProduct({ ...product, name: e.target.value })}
                    placeholder="Enter product name"
                    className={errors.name ? 'border-destructive' : ''}
                  />
                  {errors.name && <p className="text-sm text-destructive mt-1">{errors.name}</p>}
                </div>

                <div>
                  <Label htmlFor="brand">Brand</Label>
                  <Input
                    id="brand"
                    value={product.brand}
                    onChange={(e) => setProduct({ ...product, brand: e.target.value })}
                    placeholder="Enter brand name"
                  />
                </div>

                <div>
                  <Label htmlFor="sku">SKU</Label>
                  <Input
                    id="sku"
                    value={product.sku}
                    onChange={(e) => setProduct({ ...product, sku: e.target.value })}
                    placeholder="Auto-generated if empty"
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    Leave empty to auto-generate
                  </p>
                </div>

                <div>
                  <Label htmlFor="category">Category *</Label>
                  <Select 
                    value={product.category} 
                    onValueChange={(value) => setProduct({ ...product, category: value })}
                  >
                    <SelectTrigger className={errors.category ? 'border-destructive' : ''}>
                      <SelectValue placeholder="Select category" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="Rice & Grains">Rice & Grains</SelectItem>
                      <SelectItem value="Vegetables">Vegetables</SelectItem>
                      <SelectItem value="Dairy">Dairy</SelectItem>
                      <SelectItem value="Cooking">Cooking</SelectItem>
                      <SelectItem value="Pantry">Pantry</SelectItem>
                      <SelectItem value="Beverages">Beverages</SelectItem>
                      <SelectItem value="Snacks">Snacks</SelectItem>
                    </SelectContent>
                  </Select>
                  {errors.category && <p className="text-sm text-destructive mt-1">{errors.category}</p>}
                </div>

                <div>
                  <Label htmlFor="supplier">Supplier</Label>
                  <Input
                    id="supplier"
                    value={product.supplier}
                    onChange={(e) => setProduct({ ...product, supplier: e.target.value })}
                    placeholder="Enter supplier name"
                  />
                </div>

                <div className="md:col-span-2">
                  <Label htmlFor="description">Description</Label>
                  <Textarea
                    id="description"
                    value={product.description}
                    onChange={(e) => setProduct({ ...product, description: e.target.value })}
                    placeholder="Enter product description"
                    rows={3}
                  />
                </div>

                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <Label htmlFor="weight">Weight</Label>
                    <Input
                      id="weight"
                      type="number"
                      value={product.weight}
                      onChange={(e) => setProduct({ ...product, weight: e.target.value })}
                      placeholder="0"
                    />
                  </div>
                  <div>
                    <Label htmlFor="unit">Unit</Label>
                    <Select 
                      value={product.unit} 
                      onValueChange={(value) => setProduct({ ...product, unit: value })}
                    >
                      <SelectTrigger>
                        <SelectValue />
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
                  <Label htmlFor="barcode">Barcode</Label>
                  <Input
                    id="barcode"
                    value={product.barcode}
                    onChange={(e) => setProduct({ ...product, barcode: e.target.value })}
                    placeholder="Enter barcode"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Pricing */}
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center gap-2 mb-4">
                <DollarSign className="h-5 w-5 text-green-600" />
                <h3 className="text-lg font-semibold">Pricing</h3>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="price">Selling Price *</Label>
                  <Input
                    id="price"
                    type="number"
                    value={product.price}
                    onChange={(e) => setProduct({ ...product, price: e.target.value })}
                    placeholder="0.00"
                    step="0.01"
                    className={errors.price ? 'border-destructive' : ''}
                  />
                  {errors.price && <p className="text-sm text-destructive mt-1">{errors.price}</p>}
                </div>

                <div>
                  <Label htmlFor="costPrice">Cost Price *</Label>
                  <Input
                    id="costPrice"
                    type="number"
                    value={product.costPrice}
                    onChange={(e) => setProduct({ ...product, costPrice: e.target.value })}
                    placeholder="0.00"
                    step="0.01"
                    className={errors.costPrice ? 'border-destructive' : ''}
                  />
                  {errors.costPrice && <p className="text-sm text-destructive mt-1">{errors.costPrice}</p>}
                </div>

                <div>
                  <Label htmlFor="discountPrice">Discount Price</Label>
                  <Input
                    id="discountPrice"
                    type="number"
                    value={product.discountPrice}
                    onChange={(e) => setProduct({ ...product, discountPrice: e.target.value })}
                    placeholder="0.00"
                    step="0.01"
                  />
                </div>

                <div className="flex items-end">
                  <div className="bg-muted p-3 rounded-lg w-full">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-muted-foreground">Profit Margin</span>
                      <Badge variant="secondary">{profitMargin}%</Badge>
                    </div>
                    <div className="flex justify-between items-center mt-1">
                      <span className="text-sm text-muted-foreground">Profit</span>
                      <span className="text-sm font-medium">
                        {formatCurrency((parseFloat(product.price || 0) - parseFloat(product.costPrice || 0)))}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Inventory */}
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center gap-2 mb-4">
                <Warehouse className="h-5 w-5 text-blue-600" />
                <h3 className="text-lg font-semibold">Inventory</h3>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <Label htmlFor="stock">Initial Stock</Label>
                  <Input
                    id="stock"
                    type="number"
                    value={product.stock}
                    onChange={(e) => setProduct({ ...product, stock: e.target.value })}
                    placeholder="0"
                    className={errors.stock ? 'border-destructive' : ''}
                  />
                  {errors.stock && <p className="text-sm text-destructive mt-1">{errors.stock}</p>}
                </div>

                <div>
                  <Label htmlFor="minStock">Min Stock Level</Label>
                  <Input
                    id="minStock"
                    type="number"
                    value={product.minStock}
                    onChange={(e) => setProduct({ ...product, minStock: e.target.value })}
                    placeholder="0"
                    className={errors.minStock ? 'border-destructive' : ''}
                  />
                  {errors.minStock && <p className="text-sm text-destructive mt-1">{errors.minStock}</p>}
                </div>

                <div>
                  <Label htmlFor="maxStock">Max Stock Level</Label>
                  <Input
                    id="maxStock"
                    type="number"
                    value={product.maxStock}
                    onChange={(e) => setProduct({ ...product, maxStock: e.target.value })}
                    placeholder="0"
                    className={errors.maxStock ? 'border-destructive' : ''}
                  />
                  {errors.maxStock && <p className="text-sm text-destructive mt-1">{errors.maxStock}</p>}
                </div>

                <div>
                  <Label htmlFor="location">Storage Location</Label>
                  <Input
                    id="location"
                    value={product.location}
                    onChange={(e) => setProduct({ ...product, location: e.target.value })}
                    placeholder="e.g., Warehouse A - Section 2"
                  />
                </div>

                <div>
                  <Label htmlFor="expiryDate">Expiry Date</Label>
                  <Input
                    id="expiryDate"
                    type="date"
                    value={product.expiryDate}
                    onChange={(e) => setProduct({ ...product, expiryDate: e.target.value })}
                  />
                </div>

                <div>
                  <Label htmlFor="batch">Batch Number</Label>
                  <Input
                    id="batch"
                    value={product.batchNumber}
                    onChange={(e) => setProduct({ ...product, batchNumber: e.target.value })}
                    placeholder="Enter batch number"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Additional Settings */}
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center gap-2 mb-4">
                <Tag className="h-5 w-5 text-purple-600" />
                <h3 className="text-lg font-semibold">Additional Settings</h3>
              </div>
              
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <Label>Featured Product</Label>
                    <p className="text-sm text-muted-foreground">Show this product prominently</p>
                  </div>
                  <Switch
                    checked={product.featured}
                    onCheckedChange={(checked) => setProduct({ ...product, featured: checked })}
                  />
                </div>

                <div>
                  <Label>Product Status</Label>
                  <Select 
                    value={product.status} 
                    onValueChange={(value) => setProduct({ ...product, status: value })}
                  >
                    <SelectTrigger className="mt-1">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="active">Active</SelectItem>
                      <SelectItem value="inactive">Inactive</SelectItem>
                      <SelectItem value="draft">Draft</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label>Tags</Label>
                  <div className="flex gap-2 mt-1">
                    <Input
                      value={newTag}
                      onChange={(e) => setNewTag(e.target.value)}
                      placeholder="Add a tag"
                      onKeyPress={(e) => e.key === 'Enter' && addTag()}
                    />
                    <Button type="button" variant="outline" onClick={addTag}>
                      <Plus className="h-4 w-4" />
                    </Button>
                  </div>
                  {product.tags.length > 0 && (
                    <div className="flex gap-2 mt-2 flex-wrap">
                      {product.tags.map((tag, index) => (
                        <Badge key={index} variant="outline" className="gap-1">
                          {tag}
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-auto p-0 ml-1"
                            onClick={() => removeTag(tag)}
                          >
                            <X className="h-3 w-3" />
                          </Button>
                        </Badge>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default ProductCreate;