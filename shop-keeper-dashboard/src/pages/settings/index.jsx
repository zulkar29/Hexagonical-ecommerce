import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { Separator } from "@/components/ui/separator";
import { toast } from "sonner";

// Form schemas
const storeFormSchema = z.object({
  storeName: z.string().min(2, { message: "Store name must be at least 2 characters." }),
  email: z.string().email({ message: "Please enter a valid email address." }),
  phone: z.string().optional(),
  address: z.string().optional(),
  description: z.string().optional(),
  logo: z.string().optional(),
});

const shippingFormSchema = z.object({
  enableFreeShipping: z.boolean().default(false),
  freeShippingThreshold: z.string().optional(),
  flatRateShipping: z.string().optional(),
  localPickup: z.boolean().default(false),
});

const taxFormSchema = z.object({
  enableTax: z.boolean().default(true),
  taxRate: z.string().optional(),
  includeTaxInPrice: z.boolean().default(false),
});

const notificationFormSchema = z.object({
  orderConfirmation: z.boolean().default(true),
  orderStatusUpdate: z.boolean().default(true),
  lowStockAlert: z.boolean().default(true),
  marketingEmails: z.boolean().default(false),
});

const appearanceFormSchema = z.object({
  theme: z.enum(["light", "dark", "system"]).default("system"),
  primaryColor: z.string().default("#0f172a"),
  accentColor: z.string().default("#3b82f6"),
});

const Settings = () => {
  const [activeTab, setActiveTab] = useState("store");

  // Initialize forms
  const storeForm = useForm({
    resolver: zodResolver(storeFormSchema),
    defaultValues: {
      storeName: "ShopVendor",
      email: "admin@shopvendor.com",
      phone: "+1 555-123-4567",
      address: "123 Main St, City, Country",
      description: "Your one-stop shop for quality products.",
      logo: "",
    },
  });

  const shippingForm = useForm({
    resolver: zodResolver(shippingFormSchema),
    defaultValues: {
      enableFreeShipping: true,
      freeShippingThreshold: "50",
      flatRateShipping: "5",
      localPickup: true,
    },
  });

  const taxForm = useForm({
    resolver: zodResolver(taxFormSchema),
    defaultValues: {
      enableTax: true,
      taxRate: "7.5",
      includeTaxInPrice: false,
    },
  });

  const notificationForm = useForm({
    resolver: zodResolver(notificationFormSchema),
    defaultValues: {
      orderConfirmation: true,
      orderStatusUpdate: true,
      lowStockAlert: true,
      marketingEmails: false,
    },
  });

  const appearanceForm = useForm({
    resolver: zodResolver(appearanceFormSchema),
    defaultValues: {
      theme: "system",
      primaryColor: "#0f172a",
      accentColor: "#3b82f6",
    },
  });

  // Form submission handlers
  const onStoreSubmit = (data) => {
    console.log("Store settings:", data);
    toast.success("Store settings saved successfully!");
  };

  const onShippingSubmit = (data) => {
    console.log("Shipping settings:", data);
    toast.success("Shipping settings saved successfully!");
  };

  const onTaxSubmit = (data) => {
    console.log("Tax settings:", data);
    toast.success("Tax settings saved successfully!");
  };

  const onNotificationSubmit = (data) => {
    console.log("Notification settings:", data);
    toast.success("Notification settings saved successfully!");
  };

  const onAppearanceSubmit = (data) => {
    console.log("Appearance settings:", data);
    toast.success("Appearance settings saved successfully!");
  };

  return (
    <div className="container space-y-6 p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground">Manage your store settings and preferences</p>
      </div>

      <Tabs defaultValue="store" value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid grid-cols-5 mb-8">
          <TabsTrigger value="store">Store</TabsTrigger>
          <TabsTrigger value="shipping">Shipping</TabsTrigger>
          <TabsTrigger value="tax">Tax</TabsTrigger>
          <TabsTrigger value="notifications">Notifications</TabsTrigger>
          <TabsTrigger value="appearance">Appearance</TabsTrigger>
        </TabsList>

        {/* Store Settings */}
        <TabsContent value="store">
          <Card>
            <CardHeader>
              <CardTitle>Store Information</CardTitle>
              <CardDescription>
                Manage your store details and contact information
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...storeForm}>
                <form id="store-form" onSubmit={storeForm.handleSubmit(onStoreSubmit)} className="space-y-6">
                  <FormField
                    control={storeForm.control}
                    name="storeName"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Store Name</FormLabel>
                        <FormControl>
                          <Input placeholder="Your store name" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  
                  <div className="grid grid-cols-2 gap-4">
                    <FormField
                      control={storeForm.control}
                      name="email"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Email</FormLabel>
                          <FormControl>
                            <Input type="email" placeholder="contact@example.com" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    
                    <FormField
                      control={storeForm.control}
                      name="phone"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Phone</FormLabel>
                          <FormControl>
                            <Input placeholder="+1 (555) 000-0000" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  
                  <FormField
                    control={storeForm.control}
                    name="address"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Address</FormLabel>
                        <FormControl>
                          <Textarea placeholder="Store address" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  
                  <FormField
                    control={storeForm.control}
                    name="description"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Store Description</FormLabel>
                        <FormControl>
                          <Textarea placeholder="Describe your store" {...field} />
                        </FormControl>
                        <FormDescription>
                          This will be displayed on your store's homepage.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </form>
              </Form>
            </CardContent>
            <CardFooter className="flex justify-end">
              <Button type="submit" form="store-form">Save Changes</Button>
            </CardFooter>
          </Card>
        </TabsContent>

        {/* Shipping Settings */}
        <TabsContent value="shipping">
          <Card>
            <CardHeader>
              <CardTitle>Shipping Options</CardTitle>
              <CardDescription>
                Configure how you'll ship products to your customers
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...shippingForm}>
                <form id="shipping-form" onSubmit={shippingForm.handleSubmit(onShippingSubmit)} className="space-y-6">
                  <FormField
                    control={shippingForm.control}
                    name="enableFreeShipping"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                          <FormLabel className="text-base">Free Shipping</FormLabel>
                          <FormDescription>
                            Offer free shipping to your customers
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                  
                  <FormField
                    control={shippingForm.control}
                    name="freeShippingThreshold"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Free Shipping Threshold ($)</FormLabel>
                        <FormControl>
                          <Input type="number" placeholder="50" {...field} />
                        </FormControl>
                        <FormDescription>
                          Minimum order amount to qualify for free shipping
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  
                  <FormField
                    control={shippingForm.control}
                    name="flatRateShipping"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Flat Rate Shipping ($)</FormLabel>
                        <FormControl>
                          <Input type="number" placeholder="5" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  
                  <FormField
                    control={shippingForm.control}
                    name="localPickup"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                          <FormLabel className="text-base">Local Pickup</FormLabel>
                          <FormDescription>
                            Allow customers to pick up orders from your location
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                </form>
              </Form>
            </CardContent>
            <CardFooter className="flex justify-end">
              <Button type="submit" form="shipping-form">Save Changes</Button>
            </CardFooter>
          </Card>
        </TabsContent>

        {/* Tax Settings */}
        <TabsContent value="tax">
          <Card>
            <CardHeader>
              <CardTitle>Tax Configuration</CardTitle>
              <CardDescription>
                Set up tax rates and rules for your store
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...taxForm}>
                <form id="tax-form" onSubmit={taxForm.handleSubmit(onTaxSubmit)} className="space-y-6">
                  <FormField
                    control={taxForm.control}
                    name="enableTax"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                          <FormLabel className="text-base">Enable Tax Calculation</FormLabel>
                          <FormDescription>
                            Apply tax to orders based on your settings
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                  
                  <FormField
                    control={taxForm.control}
                    name="taxRate"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Default Tax Rate (%)</FormLabel>
                        <FormControl>
                          <Input type="number" placeholder="7.5" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  
                  <FormField
                    control={taxForm.control}
                    name="includeTaxInPrice"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                          <FormLabel className="text-base">Include Tax in Prices</FormLabel>
                          <FormDescription>
                            Display product prices with tax included
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                </form>
              </Form>
            </CardContent>
            <CardFooter className="flex justify-end">
              <Button type="submit" form="tax-form">Save Changes</Button>
            </CardFooter>
          </Card>
        </TabsContent>

        {/* Notification Settings */}
        <TabsContent value="notifications">
          <Card>
            <CardHeader>
              <CardTitle>Notification Preferences</CardTitle>
              <CardDescription>
                Manage email notifications for you and your customers
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...notificationForm}>
                <form id="notification-form" onSubmit={notificationForm.handleSubmit(onNotificationSubmit)} className="space-y-6">
                  <div className="space-y-4">
                    <h3 className="text-lg font-medium">Customer Notifications</h3>
                    <Separator />
                    
                    <FormField
                      control={notificationForm.control}
                      name="orderConfirmation"
                      render={({ field }) => (
                        <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                          <div className="space-y-0.5">
                            <FormLabel className="text-base">Order Confirmation</FormLabel>
                            <FormDescription>
                              Send confirmation emails when orders are placed
                            </FormDescription>
                          </div>
                          <FormControl>
                            <Switch
                              checked={field.value}
                              onCheckedChange={field.onChange}
                            />
                          </FormControl>
                        </FormItem>
                      )}
                    />
                    
                    <FormField
                      control={notificationForm.control}
                      name="orderStatusUpdate"
                      render={({ field }) => (
                        <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                          <div className="space-y-0.5">
                            <FormLabel className="text-base">Order Status Updates</FormLabel>
                            <FormDescription>
                              Notify customers when their order status changes
                            </FormDescription>
                          </div>
                          <FormControl>
                            <Switch
                              checked={field.value}
                              onCheckedChange={field.onChange}
                            />
                          </FormControl>
                        </FormItem>
                      )}
                    />
                  </div>
                  
                  <div className="space-y-4">
                    <h3 className="text-lg font-medium">Admin Notifications</h3>
                    <Separator />
                    
                    <FormField
                      control={notificationForm.control}
                      name="lowStockAlert"
                      render={({ field }) => (
                        <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                          <div className="space-y-0.5">
                            <FormLabel className="text-base">Low Stock Alerts</FormLabel>
                            <FormDescription>
                              Get notified when product inventory is running low
                            </FormDescription>
                          </div>
                          <FormControl>
                            <Switch
                              checked={field.value}
                              onCheckedChange={field.onChange}
                            />
                          </FormControl>
                        </FormItem>
                      )}
                    />
                    
                    <FormField
                      control={notificationForm.control}
                      name="marketingEmails"
                      render={({ field }) => (
                        <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                          <div className="space-y-0.5">
                            <FormLabel className="text-base">Marketing Emails</FormLabel>
                            <FormDescription>
                              Receive updates about new features and promotions
                            </FormDescription>
                          </div>
                          <FormControl>
                            <Switch
                              checked={field.value}
                              onCheckedChange={field.onChange}
                            />
                          </FormControl>
                        </FormItem>
                      )}
                    />
                  </div>
                </form>
              </Form>
            </CardContent>
            <CardFooter className="flex justify-end">
              <Button type="submit" form="notification-form">Save Changes</Button>
            </CardFooter>
          </Card>
        </TabsContent>

        {/* Appearance Settings */}
        <TabsContent value="appearance">
          <Card>
            <CardHeader>
              <CardTitle>Appearance</CardTitle>
              <CardDescription>
                Customize how your dashboard looks
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...appearanceForm}>
                <form id="appearance-form" onSubmit={appearanceForm.handleSubmit(onAppearanceSubmit)} className="space-y-6">
                  <FormField
                    control={appearanceForm.control}
                    name="theme"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Theme</FormLabel>
                        <Select onValueChange={field.onChange} defaultValue={field.value}>
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Select a theme" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            <SelectItem value="light">Light</SelectItem>
                            <SelectItem value="dark">Dark</SelectItem>
                            <SelectItem value="system">System</SelectItem>
                          </SelectContent>
                        </Select>
                        <FormDescription>
                          Select your preferred theme for the dashboard
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  
                  <div className="grid grid-cols-2 gap-4">
                    <FormField
                      control={appearanceForm.control}
                      name="primaryColor"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Primary Color</FormLabel>
                          <FormControl>
                            <div className="flex items-center gap-2">
                              <Input type="color" {...field} className="w-12 h-10 p-1" />
                              <Input type="text" value={field.value} onChange={field.onChange} />
                            </div>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    
                    <FormField
                      control={appearanceForm.control}
                      name="accentColor"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Accent Color</FormLabel>
                          <FormControl>
                            <div className="flex items-center gap-2">
                              <Input type="color" {...field} className="w-12 h-10 p-1" />
                              <Input type="text" value={field.value} onChange={field.onChange} />
                            </div>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </form>
              </Form>
            </CardContent>
            <CardFooter className="flex justify-end">
              <Button type="submit" form="appearance-form">Save Changes</Button>
            </CardFooter>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Settings;
