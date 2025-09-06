import { useForm } from "react-hook-form";
import { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowLeft, Upload, X, Plus, Loader2, ImageIcon, AlertCircle, Save, Clock } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Form, FormItem, FormLabel, FormControl, FormMessage, FormField } from "@/components/ui/form";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { useCategories, useCreateProduct, useUpdateProduct, useProduct } from "@/hooks/useApi";

const ProductForm = ({ productId, mode = "create" }) => {
  const navigate = useNavigate();
  const isEditMode = mode === "edit" && productId;

  // API hooks
  const { data: categoriesResponse, isLoading: categoriesLoading } = useCategories();
  const { data: productData, isLoading: productLoading } = useProduct(productId);
  const createProductMutation = useCreateProduct();
  const updateProductMutation = useUpdateProduct();

  // Extract categories data
  const categories = Array.isArray(categoriesResponse?.data) ? categoriesResponse.data : [];
  const product = productData?.data;

  // Form setup
  const methods = useForm({
    defaultValues: {
      name: "",
      price: "",
      stock: "",
      sku: "",
      description: "",
      category_id: "",
      status: "active",
      is_featured: false,
      meta_title: "",
      meta_description: "",
      tags: "",
      weight: "",
      length: "",
      width: "",
      height: "",
    }
  });

  const { handleSubmit, reset } = methods;
  
  // Local state
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [variants, setVariants] = useState([{ name: "", options: [""] }]);
  const [images, setImages] = useState([]);
  const [imageFiles, setImageFiles] = useState([]);
  const [autoSaveStatus, setAutoSaveStatus] = useState(null); // 'saving', 'saved', 'error'
  const [lastSaved, setLastSaved] = useState(null);

  // Load product data for edit mode
  useEffect(() => {
    if (isEditMode && product) {
      reset({
        name: product.name || "",
        price: product.price?.toString() || "",
        stock: product.stock?.toString() || "",
        sku: product.sku || "",
        description: product.description || "",
        category_id: product.category_id?.toString() || "",
        status: product.status || "active",
        is_featured: product.is_featured || false,
        meta_title: product.meta_title || "",
        meta_description: product.meta_description || "",
        tags: product.tags ? product.tags.join(", ") : "",
        weight: product.weight?.toString() || "",
        length: product.length?.toString() || "",
        width: product.width?.toString() || "",
        height: product.height?.toString() || "",
      });

      // Set categories
      if (product.categories && Array.isArray(product.categories)) {
        setSelectedCategories(product.categories.map(cat => cat.id));
      }

      // Set variants
      if (product.variants && Array.isArray(product.variants)) {
        setVariants(product.variants.length > 0 ? product.variants : [{ name: "", options: [""] }]);
      }

      // Set images
      if (product.images && Array.isArray(product.images)) {
        setImages(product.images);
      }
    }
  }, [isEditMode, product, reset]);

  const isSubmitting = createProductMutation.isPending || updateProductMutation.isPending;
  const isLoading = isEditMode ? productLoading : false;

  // Auto-save draft functionality
  const saveDraft = useCallback(async () => {
    if (isEditMode) return; // Only auto-save for new products
    
    const formData = methods.getValues();
    
    // Only save if there's meaningful content
    const hasContent = formData.name?.trim() || 
                      formData.description?.trim() || 
                      formData.price || 
                      selectedCategories.length > 0;
    
    if (!hasContent) return;
    
    try {
      setAutoSaveStatus('saving');
      
      const draftData = {
        ...formData,
        categories: selectedCategories,
        variants: variants.filter(v => v.name && v.options.some(o => o.trim())),
        isDraft: true
      };
      
      // Store in localStorage as fallback
      localStorage.setItem('product-draft', JSON.stringify({
        data: draftData,
        timestamp: Date.now()
      }));
      
      setAutoSaveStatus('saved');
      setLastSaved(new Date());
      
      // Clear status after 3 seconds
      setTimeout(() => setAutoSaveStatus(null), 3000);
    } catch (error) {
      setAutoSaveStatus('error');
      setTimeout(() => setAutoSaveStatus(null), 3000);
    }
  }, [isEditMode, methods, selectedCategories, variants]);

  // Auto-save every 30 seconds
  useEffect(() => {
    if (isEditMode) return;
    
    const interval = setInterval(saveDraft, 30000); // 30 seconds
    return () => clearInterval(interval);
  }, [saveDraft, isEditMode]);

  // Load draft on mount
  useEffect(() => {
    if (isEditMode) return;
    
    try {
      const savedDraft = localStorage.getItem('product-draft');
      if (savedDraft) {
        const { data, timestamp } = JSON.parse(savedDraft);
        
        // Only load if draft is less than 24 hours old
        const isRecent = Date.now() - timestamp < 24 * 60 * 60 * 1000;
        
        if (isRecent && data) {
          // Ask user if they want to restore
          const shouldRestore = window.confirm(
            'Found an unsaved draft from ' + new Date(timestamp).toLocaleString() + 
            '. Would you like to restore it?'
          );
          
          if (shouldRestore) {
            reset(data);
            if (data.categories) setSelectedCategories(data.categories);
            if (data.variants) setVariants(data.variants);
            setLastSaved(new Date(timestamp));
          } else {
            localStorage.removeItem('product-draft');
          }
        }
      }
    } catch (error) {
      console.error('Failed to load draft:', error);
    }
  }, [isEditMode, reset]);

  const onSubmit = async (data) => {
    try {
      // Validate required fields
      if (!data.name?.trim()) {
        toast.error("Product name is required");
        return;
      }
      
      if (!data.price || parseFloat(data.price) <= 0) {
        toast.error("Valid price is required");
        return;
      }
      
      if (data.stock === "" || parseInt(data.stock) < 0) {
        toast.error("Valid stock quantity is required");
        return;
      }

      // Validate images for new products
      if (!isEditMode && images.length === 0) {
        toast.error("At least one product image is required");
        return;
      }

      // Validate categories
      if (selectedCategories.length === 0) {
        toast.error("Please select at least one category");
        return;
      }

      // Prepare form data
      const formData = new FormData();
      
      // Basic fields with validation
      Object.keys(data).forEach(key => {
        const value = data[key];
        if (value !== "" && value !== null && value !== undefined) {
          // Convert boolean values
          if (typeof value === 'boolean') {
            formData.append(key, value ? '1' : '0');
          } else {
            formData.append(key, value.toString().trim());
          }
        }
      });

      // Categories
      selectedCategories.forEach(categoryId => {
        formData.append('category_ids[]', categoryId.toString());
      });

      // Variants - only include valid variants
      const validVariants = variants.filter(v => {
        const hasName = v.name && v.name.trim();
        const hasOptions = v.options && v.options.some(o => o && o.trim());
        return hasName && hasOptions;
      }).map(v => ({
        name: v.name.trim(),
        options: v.options.filter(o => o && o.trim()).map(o => o.trim())
      }));
      
      if (validVariants.length > 0) {
        formData.append('variants', JSON.stringify(validVariants));
      }

      // Tags processing
      if (data.tags && data.tags.trim()) {
        const tagsArray = data.tags
          .split(',')
          .map(tag => tag.trim())
          .filter(tag => tag && tag.length > 0)
          .slice(0, 10); // Limit to 10 tags
        
        if (tagsArray.length > 0) {
          formData.append('tags', JSON.stringify(tagsArray));
        }
      }

      // Images
      if (imageFiles.length > 0) {
        imageFiles.forEach(file => {
          formData.append('images[]', file);
        });
      }

      const loadingToast = toast.loading(`${isEditMode ? 'Updating' : 'Creating'} product...`);

      if (isEditMode) {
        await updateProductMutation.mutateAsync({ id: productId, data: formData });
        toast.success("Product updated successfully", { id: loadingToast });
      } else {
        await createProductMutation.mutateAsync(formData);
        toast.success("Product created successfully", { id: loadingToast });
        // Clear draft after successful creation
        localStorage.removeItem('product-draft');
      }
      
      navigate("/products");
    } catch (error) {
      console.error('Product submission error:', error);
      const errorMessage = error?.response?.data?.message || error?.message || `Failed to ${isEditMode ? 'update' : 'create'} product`;
      toast.error(errorMessage);
    }
  };

  const handleImageUpload = useCallback((event) => {
    const files = Array.from(event.target.files);
    const maxFiles = 10;
    const maxFileSize = 10 * 1024 * 1024; // 10MB
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif'];
    
    if (images.length + files.length > maxFiles) {
      toast.error(`Maximum ${maxFiles} images allowed`);
      return;
    }

    const validFiles = files.filter(file => {
      const isValidType = allowedTypes.includes(file.type.toLowerCase());
      const isValidSize = file.size <= maxFileSize;
      
      if (!isValidType) {
        toast.error(`${file.name} is not a valid image format. Use JPG, PNG, WebP, or GIF`);
        return false;
      }
      if (!isValidSize) {
        toast.error(`${file.name} is too large. Max size is 10MB`);
        return false;
      }
      
      return true;
    });

    if (validFiles.length > 0) {
      setImageFiles(prev => [...prev, ...validFiles]);
      
      // Create preview URLs with progress
      validFiles.forEach((file, index) => {
        const reader = new FileReader();
        reader.onloadstart = () => {
          setImages(prev => [...prev, { url: null, file: file, loading: true, id: Date.now() + index }]);
        };
        reader.onload = (e) => {
          setImages(prev => prev.map(img => 
            img.file === file 
              ? { ...img, url: e.target.result, loading: false }
              : img
          ));
        };
        reader.onerror = () => {
          toast.error(`Failed to load ${file.name}`);
          setImages(prev => prev.filter(img => img.file !== file));
          setImageFiles(prev => prev.filter(f => f !== file));
        };
        reader.readAsDataURL(file);
      });
    }
    
    // Reset file input
    event.target.value = '';
  }, [images.length]);

  const removeImage = useCallback((index) => {
    const imageToRemove = images[index];
    setImages(prev => prev.filter((_, i) => i !== index));
    
    // Remove from files array if it's a new upload
    if (imageToRemove?.file) {
      setImageFiles(prev => prev.filter(file => file !== imageToRemove.file));
    }
  }, [images]);


  // Variant management functions
  const addVariant = () => {
    setVariants([...variants, { name: "", options: [""] }]);
  };

  const removeVariant = (index) => {
    setVariants(variants.filter((_, i) => i !== index));
  };

  const updateVariant = (index, field, value) => {
    const newVariants = [...variants];
    newVariants[index][field] = value;
    setVariants(newVariants);
  };

  const addVariantOption = (variantIndex) => {
    const newVariants = [...variants];
    newVariants[variantIndex].options.push("");
    setVariants(newVariants);
  };

  const removeVariantOption = (variantIndex, optionIndex) => {
    const newVariants = [...variants];
    newVariants[variantIndex].options.splice(optionIndex, 1);
    setVariants(newVariants);
  };

  const updateVariantOption = (variantIndex, optionIndex, value) => {
    const newVariants = [...variants];
    newVariants[variantIndex].options[optionIndex] = value;
    setVariants(newVariants);
  };

  const toggleCategory = (categoryId) => {
    setSelectedCategories(prev =>
      prev.includes(categoryId)
        ? prev.filter(id => id !== categoryId)
        : [...prev, categoryId]
    );
  };

  if (isLoading) {
    return (
      <div className="space-y-4 sm:space-y-6 p-3 sm:p-4 lg:p-6">
        {/* Header Skeleton */}
        <div className="flex items-center gap-3">
          <div className="h-9 w-9 bg-muted rounded-md animate-pulse" />
          <div className="space-y-2">
            <div className="h-6 w-48 bg-muted rounded animate-pulse" />
            <div className="h-4 w-64 bg-muted rounded animate-pulse" />
          </div>
        </div>
        
        {/* Form Skeleton */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          <div className="lg:col-span-3 space-y-6">
            {/* Basic Information Card */}
            <Card>
              <CardHeader>
                <div className="h-5 w-32 bg-muted rounded animate-pulse" />
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <div className="h-4 w-24 bg-muted rounded animate-pulse" />
                    <div className="h-10 bg-muted rounded animate-pulse" />
                  </div>
                  <div className="space-y-2">
                    <div className="h-4 w-16 bg-muted rounded animate-pulse" />
                    <div className="h-10 bg-muted rounded animate-pulse" />
                  </div>
                </div>
                <div className="space-y-2">
                  <div className="h-4 w-20 bg-muted rounded animate-pulse" />
                  <div className="h-24 bg-muted rounded animate-pulse" />
                </div>
              </CardContent>
            </Card>
            
            {/* Additional Cards */}
            {[1, 2, 3].map(i => (
              <Card key={i}>
                <CardHeader>
                  <div className="h-5 w-40 bg-muted rounded animate-pulse" />
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="h-10 bg-muted rounded animate-pulse" />
                    <div className="h-10 bg-muted rounded animate-pulse" />
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
          
          {/* Sidebar Skeleton */}
          <div className="space-y-6">
            {[1, 2].map(i => (
              <Card key={i}>
                <CardHeader>
                  <div className="h-5 w-24 bg-muted rounded animate-pulse" />
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="h-10 bg-muted rounded animate-pulse" />
                    <div className="h-8 bg-muted rounded animate-pulse" />
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4 sm:space-y-6 p-3 sm:p-4 lg:p-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:justify-between">
        <div className="flex items-center gap-3">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/products")}
            className="h-9 w-9"
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <div>
              <h1 className="text-xl sm:text-2xl font-bold tracking-tight">
                {isEditMode ? "Edit Product" : "Add New Product"}
              </h1>
              <div className="flex items-center gap-2 mt-1">
                <p className="text-sm text-muted-foreground">
                  {isEditMode ? "Update product information" : "Create a new product for your store"}
                </p>
                {!isEditMode && (
                  <div className="flex items-center gap-1 text-xs text-muted-foreground">
                    {autoSaveStatus === 'saving' && (
                      <>
                        <Loader2 className="h-3 w-3 animate-spin" />
                        <span>Saving draft...</span>
                      </>
                    )}
                    {autoSaveStatus === 'saved' && (
                      <>
                        <Save className="h-3 w-3 text-green-600" />
                        <span className="text-green-600">Draft saved</span>
                      </>
                    )}
                    {autoSaveStatus === 'error' && (
                      <>
                        <AlertCircle className="h-3 w-3 text-destructive" />
                        <span className="text-destructive">Save failed</span>
                      </>
                    )}
                    {lastSaved && !autoSaveStatus && (
                      <>
                        <Clock className="h-3 w-3" />
                        <span>Last saved: {lastSaved.toLocaleTimeString()}</span>
                      </>
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      <Form {...methods}>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
            {/* Main Content - Left Side */}
            <div className="lg:col-span-3 space-y-6">
              {/* Basic Information */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Basic Information</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <FormField
                      name="name"
                      control={methods.control}
                      rules={{ 
                        required: "Product name is required",
                        minLength: { value: 2, message: "Product name must be at least 2 characters" },
                        maxLength: { value: 100, message: "Product name must be less than 100 characters" }
                      }}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Product Name *</FormLabel>
                          <FormControl>
                            <Input 
                              {...field} 
                              placeholder="Enter product name" 
                              className={methods.formState.errors.name ? "border-destructive" : ""}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      name="sku"
                      control={methods.control}
                      rules={{
                        pattern: {
                          value: /^[A-Za-z0-9-_]*$/,
                          message: "SKU can only contain letters, numbers, hyphens, and underscores"
                        },
                        maxLength: { value: 50, message: "SKU must be less than 50 characters" }
                      }}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>SKU</FormLabel>
                          <FormControl>
                            <Input 
                              {...field} 
                              placeholder="e.g., PROD-001" 
                              className={methods.formState.errors.sku ? "border-destructive" : ""}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  
                  <FormField
                    name="description"
                    control={methods.control}
                    rules={{
                      maxLength: { value: 2000, message: "Description must be less than 2000 characters" }
                    }}
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Description</FormLabel>
                        <FormControl>
                          <Textarea
                            {...field}
                            placeholder="Describe your product features, benefits, and specifications..."
                            className={`min-h-[120px] resize-none ${methods.formState.errors.description ? "border-destructive" : ""}`}
                          />
                        </FormControl>
                        <div className="flex justify-between text-xs text-muted-foreground">
                          <FormMessage />
                          <span>{field.value?.length || 0}/2000</span>
                        </div>
                      </FormItem>
                    )}
                  />
                </CardContent>
              </Card>

              {/* Pricing & Inventory */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Pricing & Inventory</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <FormField
                      name="price"
                      control={methods.control}
                      rules={{ 
                        required: "Price is required", 
                        min: { value: 0.01, message: "Price must be greater than $0.00" },
                        max: { value: 999999.99, message: "Price must be less than $999,999.99" },
                        pattern: {
                          value: /^\d+(\.\d{1,2})?$/,
                          message: "Please enter a valid price with up to 2 decimal places"
                        }
                      }}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Price ($) *</FormLabel>
                          <FormControl>
                            <div className="relative">
                              <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">$</span>
                              <Input
                                {...field}
                                type="number"
                                step="0.01"
                                min="0.01"
                                max="999999.99"
                                placeholder="0.00"
                                className={`pl-8 ${methods.formState.errors.price ? "border-destructive" : ""}`}
                              />
                            </div>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      name="stock"
                      control={methods.control}
                      rules={{ 
                        required: "Stock quantity is required", 
                        min: { value: 0, message: "Stock cannot be negative" },
                        max: { value: 999999, message: "Stock must be less than 999,999" },
                        pattern: {
                          value: /^\d+$/,
                          message: "Stock must be a whole number"
                        }
                      }}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Stock Quantity *</FormLabel>
                          <FormControl>
                            <Input
                              {...field}
                              type="number"
                              min="0"
                              max="999999"
                              step="1"
                              placeholder="0"
                              className={methods.formState.errors.stock ? "border-destructive" : ""}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </CardContent>
              </Card>

              {/* Product Images */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg flex items-center gap-2">
                    <ImageIcon className="h-5 w-5" />
                    Product Images
                    <span className="text-sm text-muted-foreground font-normal">({images.length}/10)</span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {/* Upload Area */}
                    <div className="relative border-2 border-dashed border-muted-foreground/25 rounded-lg p-6 hover:border-muted-foreground/40 transition-all duration-200 hover:bg-muted/10">
                      <Input
                        type="file"
                        multiple
                        accept="image/jpeg,image/jpg,image/png,image/webp,image/gif"
                        onChange={handleImageUpload}
                        className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
                        id="image-upload"
                        disabled={images.length >= 10}
                      />
                      <div className="text-center space-y-2">
                        <div className="mx-auto w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                          <Upload className="h-5 w-5 text-primary" />
                        </div>
                        <div className="space-y-1">
                          <div className="text-sm font-medium">
                            {images.length >= 10 ? "Maximum images reached" : "Click to upload or drag and drop"}
                          </div>
                          <p className="text-xs text-muted-foreground">
                            JPG, PNG, WebP, GIF up to 10MB each • Max 10 images
                          </p>
                        </div>
                      </div>
                    </div>
                    
                    {/* Image Previews */}
                    {images.length > 0 && (
                      <div className="space-y-3">
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <AlertCircle className="h-4 w-4" />
                          <span>First image will be used as the main product image</span>
                        </div>
                        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
                          {images.map((image, index) => (
                            <div 
                              key={image.id || index} 
                              className={`relative aspect-square border-2 rounded-lg overflow-hidden group ${
                                index === 0 ? 'border-primary ring-2 ring-primary/20' : 'border-muted-foreground/20'
                              }`}
                            >
                              {image.loading ? (
                                <div className="w-full h-full flex items-center justify-center bg-muted">
                                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                                </div>
                              ) : (
                                <>
                                  <img
                                    src={image.url || image}
                                    alt={`Product image ${index + 1}`}
                                    className="w-full h-full object-cover transition-transform group-hover:scale-105"
                                    loading="lazy"
                                  />
                                  <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors" />
                                  
                                  {/* Image Controls */}
                                  <div className="absolute top-2 right-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                    <Button
                                      type="button"
                                      variant="destructive"
                                      size="icon"
                                      className="h-6 w-6 bg-destructive/90 hover:bg-destructive"
                                      onClick={() => removeImage(index)}
                                    >
                                      <X className="h-3 w-3" />
                                    </Button>
                                  </div>
                                  
                                  {/* Main Image Badge */}
                                  {index === 0 && (
                                    <div className="absolute bottom-2 left-2">
                                      <Badge variant="default" className="text-xs px-2 py-0.5">
                                        Main
                                      </Badge>
                                    </div>
                                  )}
                                  
                                  {/* Image Number */}
                                  <div className="absolute bottom-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                    <Badge variant="secondary" className="text-xs px-1.5 py-0.5">
                                      {index + 1}
                                    </Badge>
                                  </div>
                                </>
                              )}
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>

              {/* Product Variants */}
              <Card>
                <CardHeader>
                  <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                    <div>
                      <CardTitle className="text-lg">Product Variants</CardTitle>
                      <p className="text-sm text-muted-foreground mt-1">
                        Add variants like size, color, or material to offer different options
                      </p>
                    </div>
                    <Button 
                      type="button" 
                      variant="outline" 
                      size="sm" 
                      onClick={addVariant}
                      className="shrink-0"
                      disabled={variants.length >= 5}
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      Add Variant
                    </Button>
                  </div>
                </CardHeader>
                <CardContent className="space-y-6">
                  {variants.length === 1 && !variants[0].name && !variants[0].options[0] ? (
                    <div className="text-center py-8 border-2 border-dashed border-muted-foreground/25 rounded-lg">
                      <div className="text-muted-foreground">
                        <p className="text-sm">No variants added yet</p>
                        <p className="text-xs mt-1">Add variants to offer different product options</p>
                      </div>
                    </div>
                  ) : (
                    variants.map((variant, variantIndex) => (
                      <div key={variantIndex} className="p-4 border rounded-lg space-y-4 bg-muted/20">
                        <div className="flex items-start gap-2">
                          <div className="flex-1">
                            <label className="text-sm font-medium mb-2 block">
                              Variant Name {variantIndex + 1}
                            </label>
                            <Input
                              placeholder="e.g., Size, Color, Material"
                              value={variant.name}
                              onChange={(e) => updateVariant(variantIndex, 'name', e.target.value)}
                              className="bg-background"
                            />
                          </div>
                          {variants.length > 1 && (
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              onClick={() => removeVariant(variantIndex)}
                              className="mt-6 text-destructive hover:text-destructive hover:bg-destructive/10"
                            >
                              <X className="h-4 w-4" />
                            </Button>
                          )}
                        </div>
                        
                        <div className="space-y-3">
                          <label className="text-sm font-medium block">
                            Options for {variant.name || 'this variant'}
                          </label>
                          <div className="space-y-2">
                            {variant.options.map((option, optionIndex) => (
                              <div key={optionIndex} className="flex items-center gap-2">
                                <Input
                                  placeholder={`Option ${optionIndex + 1} (e.g., Small, Red, Cotton)`}
                                  value={option}
                                  onChange={(e) => updateVariantOption(variantIndex, optionIndex, e.target.value)}
                                  className="flex-1 bg-background"
                                />
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => addVariantOption(variantIndex)}
                                  disabled={variant.options.length >= 10}
                                  title="Add option"
                                >
                                  <Plus className="h-4 w-4" />
                                </Button>
                                {variant.options.length > 1 && (
                                  <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => removeVariantOption(variantIndex, optionIndex)}
                                    className="text-destructive hover:text-destructive hover:bg-destructive/10"
                                    title="Remove option"
                                  >
                                    <X className="h-4 w-4" />
                                  </Button>
                                )}
                              </div>
                            ))}
                          </div>
                          
                          {/* Variant Preview */}
                          {variant.name && variant.options.some(opt => opt.trim()) && (
                            <div className="mt-3 p-3 bg-background rounded border">
                              <p className="text-xs font-medium text-muted-foreground mb-2">Preview:</p>
                              <div className="flex flex-wrap gap-1">
                                {variant.options
                                  .filter(opt => opt.trim())
                                  .map((option, idx) => (
                                    <Badge key={idx} variant="secondary" className="text-xs">
                                      {variant.name}: {option}
                                    </Badge>
                                  ))
                                }
                              </div>
                            </div>
                          )}
                        </div>
                      </div>
                    ))
                  )}
                  
                  {variants.length >= 5 && (
                    <p className="text-xs text-muted-foreground text-center">
                      Maximum 5 variants allowed
                    </p>
                  )}
                </CardContent>
              </Card>

              {/* SEO & Marketing */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">SEO & Marketing</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <FormField
                    name="meta_title"
                    control={methods.control}
                    rules={{
                      maxLength: { value: 60, message: "Meta title should be under 60 characters for better SEO" }
                    }}
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Meta Title</FormLabel>
                        <FormControl>
                          <Input 
                            {...field} 
                            placeholder="SEO-friendly title for search engines" 
                            className={methods.formState.errors.meta_title ? "border-destructive" : ""}
                          />
                        </FormControl>
                        <div className="flex justify-between text-xs text-muted-foreground">
                          <FormMessage />
                          <span className={field.value?.length > 60 ? "text-destructive" : ""}>
                            {field.value?.length || 0}/60
                          </span>
                        </div>
                      </FormItem>
                    )}
                  />
                  <FormField
                    name="meta_description"
                    control={methods.control}
                    rules={{
                      maxLength: { value: 160, message: "Meta description should be under 160 characters for better SEO" }
                    }}
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Meta Description</FormLabel>
                        <FormControl>
                          <Textarea
                            {...field}
                            placeholder="Compelling description for search engine results (150-160 characters)"
                            className={`min-h-[80px] resize-none ${methods.formState.errors.meta_description ? "border-destructive" : ""}`}
                          />
                        </FormControl>
                        <div className="flex justify-between text-xs text-muted-foreground">
                          <FormMessage />
                          <span className={field.value?.length > 160 ? "text-destructive" : ""}>
                            {field.value?.length || 0}/160
                          </span>
                        </div>
                      </FormItem>
                    )}
                  />
                  <FormField
                    name="tags"
                    control={methods.control}
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Product Tags</FormLabel>
                        <FormControl>
                          <Input {...field} placeholder="Enter tags separated by commas" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </CardContent>
              </Card>

              {/* Shipping & Dimensions */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Shipping & Dimensions</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <FormField
                      name="weight"
                      control={methods.control}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Weight (kg)</FormLabel>
                          <FormControl>
                            <Input {...field} type="number" step="0.01" placeholder="0.00" />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      name="length"
                      control={methods.control}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Length (cm)</FormLabel>
                          <FormControl>
                            <Input {...field} type="number" placeholder="0" />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      name="width"
                      control={methods.control}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Width (cm)</FormLabel>
                          <FormControl>
                            <Input {...field} type="number" placeholder="0" />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <FormField
                    name="height"
                    control={methods.control}
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Height (cm)</FormLabel>
                        <FormControl>
                          <Input {...field} type="number" placeholder="0" className="max-w-xs" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </CardContent>
              </Card>
            </div>

            {/* Sidebar - Right Side */}
            <div className="space-y-6 lg:sticky lg:top-6 lg:self-start">
              {/* Status */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Status</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <FormField
                    name="status"
                    control={methods.control}
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Product Status</FormLabel>
                        <FormControl>
                          <Select onValueChange={field.onChange} value={field.value}>
                            <SelectTrigger>
                              <SelectValue placeholder="Select status" />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="active">Active</SelectItem>
                              <SelectItem value="draft">Draft</SelectItem>
                              <SelectItem value="archived">Archived</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  
                  <FormField
                    name="is_featured"
                    control={methods.control}
                    render={({ field }) => (
                      <FormItem className="flex items-center justify-between">
                        <FormLabel>Featured Product</FormLabel>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                </CardContent>
              </Card>

              {/* Categories */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Categories</CardTitle>
                </CardHeader>
                <CardContent>
                  {categoriesLoading ? (
                    <div className="flex items-center justify-center py-4">
                      <Loader2 className="h-4 w-4 animate-spin" />
                    </div>
                  ) : (
                    <>
                      <div className="space-y-2">
                        {categories.map((category) => (
                          <div
                            key={category.id}
                            className={`p-2 rounded-md cursor-pointer transition-colors text-sm ${
                              selectedCategories.includes(category.id)
                                ? 'bg-primary/10 text-primary border border-primary/20'
                                : 'bg-muted hover:bg-muted/80'
                            }`}
                            onClick={() => toggleCategory(category.id)}
                          >
                            {category.name}
                          </div>
                        ))}
                      </div>
                      {selectedCategories.length > 0 && (
                        <div className="mt-3">
                          <p className="text-sm font-medium mb-2">Selected:</p>
                          <div className="flex flex-wrap gap-1">
                            {selectedCategories.map((categoryId) => {
                              const category = categories.find(cat => cat.id === categoryId);
                              return category ? (
                                <Badge
                                  key={categoryId}
                                  variant="secondary"
                                  className="text-xs"
                                >
                                  {category.name}
                                  <X
                                    className="h-3 w-3 ml-1 cursor-pointer"
                                    onClick={() => toggleCategory(categoryId)}
                                  />
                                </Badge>
                              ) : null;
                            })}
                          </div>
                        </div>
                      )}
                    </>
                  )}
                </CardContent>
              </Card>

              {/* Actions */}
              <Card>
                <CardContent className="pt-6">
                  <div className="flex flex-col gap-3">
                    <Button type="submit" className="w-full" disabled={isSubmitting}>
                      {isSubmitting && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                      {isEditMode ? "Update Product" : "Create Product"}
                    </Button>
                    
                    {!isEditMode && (
                      <Button
                        type="button"
                        variant="secondary"
                        onClick={saveDraft}
                        className="w-full"
                        disabled={isSubmitting}
                      >
                        <Save className="h-4 w-4 mr-2" />
                        Save as Draft
                      </Button>
                    )}
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => navigate("/products")}
                      className="w-full"
                      disabled={isSubmitting}
                    >
                      Cancel
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </form>
      </Form>
    </div>
  );
};

export default ProductForm;