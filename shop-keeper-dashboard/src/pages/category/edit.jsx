import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { 
  Select, 
  SelectTrigger, 
  SelectValue, 
  SelectContent, 
  SelectItem 
} from "@/components/ui/select";
import { 
  ArrowLeft,
  Save,
  FolderTree,
  Image as ImageIcon,
  Upload,
  X,
  Tag,
  Settings,
  ChevronRight
} from "lucide-react";
import { Switch } from "@/components/ui/switch";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";

// Mock categories for parent selection (same as create.jsx)
const mockCategories = [
  { id: 1, name: "Electronics", slug: "electronics", level: 0, children: [
    { id: 2, name: "Smartphones", slug: "smartphones", level: 1, children: [] },
    { id: 6, name: "Laptops", slug: "laptops", level: 1, children: [] },
    { id: 7, name: "Headphones", slug: "headphones", level: 1, children: [] }
  ]},
  { id: 9, name: "Fashion", slug: "fashion", level: 0, children: [
    { id: 10, name: "Men's Clothing", slug: "mens-clothing", level: 1, children: [] },
    { id: 11, name: "Women's Clothing", slug: "womens-clothing", level: 1, children: [] }
  ]},
  { id: 13, name: "Home & Garden", slug: "home-garden", level: 0, children: [
    { id: 14, name: "Kitchen", slug: "kitchen", level: 1, children: [] }
  ]}
];

// Mock detailed category data for editing
const generateCategoryDetail = (id) => {
  const categoriesData = {
    1: {
      id: 1,
      name: "Electronics",
      slug: "electronics",
      description: "Electronic devices, gadgets, and accessories for modern life",
      parent_id: null,
      is_active: true,
      sort_order: 1,
      seo_title: "Electronics - Latest Gadgets & Devices",
      seo_description: "Shop the latest electronics including smartphones, laptops, headphones and more.",
      meta_keywords: "electronics, gadgets, smartphones, laptops, headphones",
      image: null
    },
    2: {
      id: 2,
      name: "Smartphones",
      slug: "smartphones",
      description: "Mobile phones and smartphone accessories",
      parent_id: 1,
      is_active: true,
      sort_order: 1,
      seo_title: "Smartphones - Latest Mobile Phones",
      seo_description: "Browse our collection of latest smartphones from top brands.",
      meta_keywords: "smartphones, mobile phones, iPhone, Android",
      image: null
    }
  };
  return categoriesData[id] || null;
};

const flattenCategories = (categories, level = 0) => {
  let result = [];
  categories.forEach(category => {
    result.push({ ...category, level });
    if (category.children && category.children.length > 0) {
      result = result.concat(flattenCategories(category.children, level + 1));
    }
  });
  return result;
};

export default function EditCategory() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [category, setCategory] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const [form, setForm] = useState({
    name: "",
    slug: "",
    description: "",
    parent_id: "",
    is_active: true,
    sort_order: 0,
    seo_title: "",
    seo_description: "",
    meta_keywords: "",
    image: null
  });
  const [submitting, setSubmitting] = useState(false);
  const [slugManuallyEdited, setSlugManuallyEdited] = useState(false);
  const [selectedParent, setSelectedParent] = useState(null);

  const flatCategories = flattenCategories(mockCategories);

  useEffect(() => {
    const fetchCategory = async () => {
      setIsLoading(true);
      try {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 1000));
        const categoryData = generateCategoryDetail(parseInt(id));
        
        if (categoryData) {
          setCategory(categoryData);
          setForm({
            name: categoryData.name,
            slug: categoryData.slug,
            description: categoryData.description,
            parent_id: categoryData.parent_id ? categoryData.parent_id.toString() : "",
            is_active: categoryData.is_active,
            sort_order: categoryData.sort_order,
            seo_title: categoryData.seo_title,
            seo_description: categoryData.seo_description,
            meta_keywords: categoryData.meta_keywords,
            image: categoryData.image
          });
          
          if (categoryData.parent_id) {
            const parent = flatCategories.find(cat => cat.id === categoryData.parent_id);
            setSelectedParent(parent);
          }
        }
      } catch {
        toast.error("Failed to load category details");
        navigate("/category");
      } finally {
        setIsLoading(false);
      }
    };

    if (id) {
      fetchCategory();
    }
  }, [id, navigate, flatCategories]);

  const generateSlug = (name) => {
    return name
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .trim();
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));

    if (name === 'name' && !slugManuallyEdited) {
      setForm(prev => ({ ...prev, slug: generateSlug(value) }));
    }

    if (name === 'slug') {
      setSlugManuallyEdited(true);
    }
  };

  const handleParentChange = (value) => {
    const actualValue = value === "none" ? "" : value;
    setForm(prev => ({ ...prev, parent_id: actualValue }));
    const parent = value === "none" ? null : flatCategories.find(cat => cat.id === parseInt(value));
    setSelectedParent(parent);
  };

  const handleSwitchChange = (checked) => {
    setForm(prev => ({ ...prev, is_active: checked }));
  };

  const handleImageUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      if (file.size > 5 * 1024 * 1024) { // 5MB limit
        toast.error("File size must be less than 5MB");
        return;
      }
      setForm(prev => ({ ...prev, image: file }));
      toast.success("Image uploaded successfully");
    }
  };

  const removeImage = () => {
    setForm(prev => ({ ...prev, image: null }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1500));
      
      console.log("Updating category:", form);
      toast.success("Category updated successfully!");
      navigate("/category");
    } catch {
      toast.error("Failed to update category. Please try again.");
    } finally {
      setSubmitting(false);
    }
  };

  const renderCategoryOption = (category) => {
    const indent = "  ".repeat(category.level);
    return (
      <SelectItem key={category.id} value={category.id.toString()}>
        {indent}{category.name}
      </SelectItem>
    );
  };

  if (isLoading) {
    return (
      <div className="space-y-6 p-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mr-2"></div>
              <span>Loading category details...</span>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!category) {
    return (
      <div className="space-y-6 p-6">
        <Card>
          <CardHeader>
            <CardTitle>Category Not Found</CardTitle>
            <CardDescription>The category you are trying to edit does not exist.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => navigate("/category")}>
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Categories
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
          onClick={() => navigate("/category")}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold tracking-tight">Edit Category</h1>
          <p className="text-muted-foreground">
            Update the details for "{category.name}"
          </p>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main Form */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle>Category Information</CardTitle>
              <CardDescription>
                Update the details for this category
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-6">
                {/* Basic Information */}
                <div className="space-y-4">
                  <h3 className="text-lg font-medium flex items-center gap-2">
                    <Tag className="h-5 w-5" />
                    Basic Information
                  </h3>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="name">Category Name *</Label>
                      <Input
                        id="name"
                        name="name"
                        value={form.name}
                        onChange={handleChange}
                        placeholder="Enter category name"
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="slug">URL Slug *</Label>
                      <Input
                        id="slug"
                        name="slug"
                        value={form.slug}
                        onChange={handleChange}
                        placeholder="category-url-slug"
                        required
                      />
                      <p className="text-xs text-muted-foreground">
                        URL: /category/{form.slug || 'category-slug'}
                      </p>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="description">Description</Label>
                    <Textarea
                      id="description"
                      name="description"
                      value={form.description}
                      onChange={handleChange}
                      placeholder="Brief description of the category"
                      rows={3}
                    />
                  </div>
                </div>

                <Separator />

                {/* Category Hierarchy */}
                <div className="space-y-4">
                  <h3 className="text-lg font-medium flex items-center gap-2">
                    <FolderTree className="h-5 w-5" />
                    Category Hierarchy
                  </h3>
                  
                  <div className="space-y-2">
                    <Label htmlFor="parent_id">Parent Category</Label>
                    <Select value={form.parent_id || "none"} onValueChange={handleParentChange}>
                      <SelectTrigger>
                        <SelectValue placeholder="Select parent category (optional)" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="none">No Parent (Root Category)</SelectItem>
                        {flatCategories.filter(cat => cat.id !== parseInt(id)).map(renderCategoryOption)}
                      </SelectContent>
                    </Select>
                    {selectedParent && (
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <span>Parent category:</span>
                        <Badge variant="outline">{selectedParent.name}</Badge>
                      </div>
                    )}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="sort_order">Sort Order</Label>
                    <Input
                      id="sort_order"
                      name="sort_order"
                      type="number"
                      value={form.sort_order}
                      onChange={handleChange}
                      placeholder="0"
                      min="0"
                    />
                    <p className="text-xs text-muted-foreground">
                      Lower numbers appear first in category lists
                    </p>
                  </div>
                </div>

                <Separator />

                {/* Category Image */}
                <div className="space-y-4">
                  <h3 className="text-lg font-medium flex items-center gap-2">
                    <ImageIcon className="h-5 w-5" />
                    Category Image
                  </h3>
                  
                  <div className="space-y-4">
                    {form.image ? (
                      <div className="relative">
                        <div className="w-full h-48 bg-gray-100 rounded-lg flex items-center justify-center border">
                          <div className="text-center">
                            <ImageIcon className="h-12 w-12 text-gray-400 mx-auto mb-2" />
                            <p className="text-sm text-gray-600">{form.image.name || 'Current image'}</p>
                            {form.image.size && (
                              <p className="text-xs text-gray-500">
                                {(form.image.size / 1024 / 1024).toFixed(2)} MB
                              </p>
                            )}
                          </div>
                        </div>
                        <Button
                          type="button"
                          variant="destructive"
                          size="icon"
                          className="absolute top-2 right-2"
                          onClick={removeImage}
                        >
                          <X className="h-4 w-4" />
                        </Button>
                      </div>
                    ) : (
                      <div className="border-2 border-dashed border-muted-foreground/25 rounded-lg p-8 text-center">
                        <ImageIcon className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
                        <div className="space-y-2">
                          <p className="text-sm text-muted-foreground">
                            Upload a category image
                          </p>
                          <input
                            type="file"
                            id="image"
                            accept="image/*"
                            onChange={handleImageUpload}
                            className="hidden"
                          />
                          <Button
                            type="button"
                            variant="outline"
                            onClick={() => document.getElementById("image").click()}
                          >
                            <Upload className="mr-2 h-4 w-4" />
                            Choose Image
                          </Button>
                        </div>
                      </div>
                    )}
                    <p className="text-xs text-muted-foreground">
                      Recommended: 800x600px, max 5MB, JPG or PNG
                    </p>
                  </div>
                </div>

                <Separator />

                {/* SEO Settings */}
                <div className="space-y-4">
                  <h3 className="text-lg font-medium flex items-center gap-2">
                    <Settings className="h-5 w-5" />
                    SEO Settings
                  </h3>
                  
                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="seo_title">SEO Title</Label>
                      <Input
                        id="seo_title"
                        name="seo_title"
                        value={form.seo_title}
                        onChange={handleChange}
                        placeholder="SEO optimized title"
                        maxLength={60}
                      />
                      <p className="text-xs text-muted-foreground">
                        {form.seo_title.length}/60 characters
                      </p>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="seo_description">Meta Description</Label>
                      <Textarea
                        id="seo_description"
                        name="seo_description"
                        value={form.seo_description}
                        onChange={handleChange}
                        placeholder="Brief description for search engines"
                        maxLength={160}
                        rows={3}
                      />
                      <p className="text-xs text-muted-foreground">
                        {form.seo_description.length}/160 characters
                      </p>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="meta_keywords">Meta Keywords</Label>
                      <Input
                        id="meta_keywords"
                        name="meta_keywords"
                        value={form.meta_keywords}
                        onChange={handleChange}
                        placeholder="keyword1, keyword2, keyword3"
                      />
                      <p className="text-xs text-muted-foreground">
                        Separate keywords with commas
                      </p>
                    </div>
                  </div>
                </div>

                {/* Submit Buttons */}
                <div className="flex justify-end gap-3 pt-6">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => navigate("/category")}
                  >
                    Cancel
                  </Button>
                  <Button type="submit" disabled={submitting}>
                    {submitting ? (
                      <>
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                        Updating...
                      </>
                    ) : (
                      <>
                        <Save className="mr-2 h-4 w-4" />
                        Update Category
                      </>
                    )}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Category Status */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Category Status</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label htmlFor="is_active">Active Status</Label>
                  <p className="text-sm text-muted-foreground">
                    Make category visible to customers
                  </p>
                </div>
                <Switch
                  id="is_active"
                  checked={form.is_active}
                  onCheckedChange={handleSwitchChange}
                />
              </div>
              
              <div className="p-3 bg-muted rounded-lg">
                <p className="text-sm font-medium">
                  Status: {form.is_active ? "Active" : "Inactive"}
                </p>
                <p className="text-xs text-muted-foreground">
                  {form.is_active 
                    ? "Category will be visible to customers and in navigation"
                    : "Category will be hidden from customers"
                  }
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Category Preview */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Category Preview</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="p-4 border rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <FolderTree className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">{form.name || "Category Name"}</span>
                </div>
                
                {selectedParent && (
                  <div className="flex items-center gap-1 text-xs text-muted-foreground mb-2">
                    <span>{selectedParent.name}</span>
                    <ChevronRight className="h-3 w-3" />
                    <span>{form.name || "Category"}</span>
                  </div>
                )}
                
                <p className="text-sm text-muted-foreground">
                  {form.description || "Category description will appear here"}
                </p>
                
                <div className="flex items-center gap-2 mt-3">
                  <Badge variant={form.is_active ? "default" : "secondary"}>
                    {form.is_active ? "Active" : "Inactive"}
                  </Badge>
                  {form.slug && (
                    <Badge variant="outline" className="font-mono text-xs">
                      /{form.slug}
                    </Badge>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Tips */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Tips</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="text-sm text-muted-foreground space-y-2">
                <p><strong>Hierarchy:</strong> Changing parent category affects URL structure</p>
                <p><strong>SEO:</strong> Update meta information for better search rankings</p>
                <p><strong>Status:</strong> Inactive categories hide all products</p>
                <p><strong>Slug:</strong> URL changes may affect existing links</p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
