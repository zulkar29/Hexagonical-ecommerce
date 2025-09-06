import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { ArrowLeft, ShieldCheck, Users, Settings, FileText, CreditCard, Package } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Form, FormItem, FormLabel, FormControl, FormMessage, FormField } from "@/components/ui/form";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

const permissions = [
  {
    category: "Products",
    icon: <Package className="h-4 w-4" />,
    items: [
      { id: "products_view", name: "View Products", description: "View product listings and details" },
      { id: "products_create", name: "Create Products", description: "Add new products to the store" },
      { id: "products_edit", name: "Edit Products", description: "Modify existing product information" },
      { id: "products_delete", name: "Delete Products", description: "Remove products from the store" },
    ]
  },
  {
    category: "Orders",
    icon: <FileText className="h-4 w-4" />,
    items: [
      { id: "orders_view", name: "View Orders", description: "View order listings and details" },
      { id: "orders_manage", name: "Manage Orders", description: "Process and update order status" },
      { id: "orders_cancel", name: "Cancel Orders", description: "Cancel customer orders" },
    ]
  },
  {
    category: "Customers",
    icon: <Users className="h-4 w-4" />,
    items: [
      { id: "customers_view", name: "View Customers", description: "View customer profiles and history" },
      { id: "customers_edit", name: "Edit Customers", description: "Modify customer information" },
      { id: "customers_delete", name: "Delete Customers", description: "Remove customer accounts" },
    ]
  },
  {
    category: "Payments",
    icon: <CreditCard className="h-4 w-4" />,
    items: [
      { id: "payments_view", name: "View Payments", description: "View payment transactions and history" },
      { id: "payments_refund", name: "Process Refunds", description: "Issue refunds to customers" },
    ]
  },
  {
    category: "Settings",
    icon: <Settings className="h-4 w-4" />,
    items: [
      { id: "settings_view", name: "View Settings", description: "Access store configuration" },
      { id: "settings_edit", name: "Edit Settings", description: "Modify store settings and configuration" },
      { id: "users_manage", name: "Manage Users", description: "Add, edit, and remove user accounts" },
    ]
  },
];

const CreateRole = () => {
  const navigate = useNavigate();
  const methods = useForm({
    defaultValues: {
      name: "",
      description: "",
      permissions: [],
    }
  });
  const { handleSubmit } = methods;
  const [selectedPermissions, setSelectedPermissions] = useState(new Set());

  const onSubmit = (data) => {
    const formData = {
      ...data,
      permissions: Array.from(selectedPermissions),
    };
    console.log("Role created:", formData);
    navigate("/user-roles");
  };

  const togglePermission = (permissionId) => {
    const newPermissions = new Set(selectedPermissions);
    if (newPermissions.has(permissionId)) {
      newPermissions.delete(permissionId);
    } else {
      newPermissions.add(permissionId);
    }
    setSelectedPermissions(newPermissions);
  };

  const toggleCategory = (category) => {
    const categoryPermissions = permissions.find(p => p.category === category)?.items.map(item => item.id) || [];
    const newPermissions = new Set(selectedPermissions);
    const allSelected = categoryPermissions.every(id => newPermissions.has(id));
    
    if (allSelected) {
      categoryPermissions.forEach(id => newPermissions.delete(id));
    } else {
      categoryPermissions.forEach(id => newPermissions.add(id));
    }
    setSelectedPermissions(newPermissions);
  };

  return (
    <div className="space-y-4 sm:space-y-6 p-3 sm:p-4 lg:p-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:justify-between">
        <div className="flex items-center gap-3">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/user-roles")}
            className="h-9 w-9"
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-xl sm:text-2xl font-bold tracking-tight">Add New Role</h1>
            <p className="text-sm text-muted-foreground">Create a new user role with specific permissions</p>
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
                  <CardTitle className="text-lg">Role Information</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 gap-4">
                    <FormField
                      name="name"
                      control={methods.control}
                      rules={{ required: "Role name is required" }}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Role Name *</FormLabel>
                          <FormControl>
                            <Input {...field} placeholder="Enter role name (e.g., Manager, Editor)" />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    
                    <FormField
                      name="description"
                      control={methods.control}
                      rules={{ required: "Description is required" }}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Description *</FormLabel>
                          <FormControl>
                            <Textarea
                              {...field}
                              placeholder="Describe what this role can do and its responsibilities..."
                              className="min-h-[100px] resize-none"
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </CardContent>
              </Card>

              {/* Permissions */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Permissions</CardTitle>
                  <p className="text-sm text-muted-foreground">
                    Select the permissions this role should have. You can select individual permissions or entire categories.
                  </p>
                </CardHeader>
                <CardContent className="space-y-6">
                  {permissions.map((category) => {
                    const categoryPermissions = category.items.map(item => item.id);
                    const selectedCount = categoryPermissions.filter(id => selectedPermissions.has(id)).length;
                    const allSelected = selectedCount === categoryPermissions.length;
                    const someSelected = selectedCount > 0 && !allSelected;

                    return (
                      <div key={category.category} className="space-y-4">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <div className="p-2 bg-primary/10 rounded-lg">
                              {category.icon}
                            </div>
                            <div>
                              <h3 className="font-medium">{category.category}</h3>
                              <p className="text-sm text-muted-foreground">
                                {selectedCount} of {categoryPermissions.length} selected
                              </p>
                            </div>
                          </div>
                          <Button
                            type="button"
                            variant={allSelected ? "default" : "outline"}
                            size="sm"
                            onClick={() => toggleCategory(category.category)}
                            className={someSelected && !allSelected ? "border-primary text-primary" : ""}
                          >
                            {allSelected ? "Deselect All" : someSelected ? "Select All" : "Select All"}
                          </Button>
                        </div>
                        
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pl-4">
                          {category.items.map((permission) => (
                            <div
                              key={permission.id}
                              className={`p-3 rounded-lg border cursor-pointer transition-colors ${
                                selectedPermissions.has(permission.id)
                                  ? 'bg-primary/5 border-primary/20'
                                  : 'hover:bg-muted/50'
                              }`}
                              onClick={() => togglePermission(permission.id)}
                            >
                              <div className="flex items-start justify-between">
                                <div className="flex-1">
                                  <div className="flex items-center gap-2">
                                    <span className="font-medium text-sm">{permission.name}</span>
                                    {selectedPermissions.has(permission.id) && (
                                      <ShieldCheck className="h-4 w-4 text-primary" />
                                    )}
                                  </div>
                                  <p className="text-xs text-muted-foreground mt-1">{permission.description}</p>
                                </div>
                              </div>
                            </div>
                          ))}
                        </div>
                        
                        {category !== permissions[permissions.length - 1] && (
                          <Separator />
                        )}
                      </div>
                    );
                  })}
                </CardContent>
              </Card>
            </div>

            {/* Sidebar - Right Side */}
            <div className="space-y-6 lg:sticky lg:top-6 lg:self-start">
              {/* Role Preview */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Role Preview</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-primary/10 rounded-lg">
                        <ShieldCheck className="h-5 w-5 text-primary" />
                      </div>
                      <div className="flex-1">
                        <h3 className="font-medium">New Role</h3>
                        <p className="text-xs text-muted-foreground">Draft</p>
                      </div>
                    </div>
                    
                    <Separator />
                    
                    <div className="space-y-2">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Total Permissions:</span>
                        <Badge variant="secondary">{selectedPermissions.size}</Badge>
                      </div>
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Categories:</span>
                        <Badge variant="outline">
                          {permissions.filter(cat => 
                            cat.items.some(item => selectedPermissions.has(item.id))
                          ).length}
                        </Badge>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Actions */}
              <Card>
                <CardContent className="pt-6">
                  <div className="flex flex-col gap-3">
                    <Button type="submit" className="w-full">
                      Create Role
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => navigate("/user-roles")}
                      className="w-full"
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

export default CreateRole;