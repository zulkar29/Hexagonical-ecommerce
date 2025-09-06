import { useState } from "react";
import { useForm } from "react-hook-form";
import { 
  User, 
  Mail, 
  Phone, 
  MapPin, 
  Calendar, 
  Shield, 
  Activity, 
  Settings, 
  Camera,
  Edit3,
  Save,
  X,
  Clock,
  TrendingUp,
  ShoppingCart,
  Package,
  Users,
  Star,
  Bell,
  Key,
  Smartphone,
  Globe,
  Building
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Switch } from "@/components/ui/switch";
import { Separator } from "@/components/ui/separator";
import { Progress } from "@/components/ui/progress";
import { useNavigate } from "react-router-dom";

export default function Profile() {
  const navigate = useNavigate();
  const [isEditing, setIsEditing] = useState(false);
  const [activeTab, setActiveTab] = useState("overview");

  const { register, handleSubmit, formState: { errors } } = useForm({
    defaultValues: {
      name: "Admin User",
      email: "admin@shopvendor.com", 
      phone: "+1 555-123-4567",
      bio: "Store owner and administrator with 5+ years of e-commerce experience.",
      company: "ShopVendor Inc.",
      website: "https://shopvendor.com",
      location: "San Francisco, CA",
      timezone: "Pacific Time (PST)"
    },
  });

  const onSubmit = (data) => {
    console.log("Profile updated:", data);
    setIsEditing(false);
    // Save profile logic here
  };

  // Mock data for dashboard stats
  const stats = [
    { label: "Orders Managed", value: "2,847", icon: ShoppingCart, trend: "+12%", color: "text-blue-600" },
    { label: "Products Added", value: "156", icon: Package, trend: "+8%", color: "text-green-600" },
    { label: "Customers Served", value: "1,329", icon: Users, trend: "+15%", color: "text-purple-600" },
    { label: "Average Rating", value: "4.8", icon: Star, trend: "+0.2", color: "text-yellow-600" }
  ];

  const recentActivity = [
    { action: "Updated product inventory", time: "2 hours ago", icon: Package },
    { action: "Processed 5 orders", time: "4 hours ago", icon: ShoppingCart },
    { action: "Responded to customer inquiry", time: "6 hours ago", icon: Mail },
    { action: "Added new product category", time: "1 day ago", icon: Settings },
    { action: "Generated sales report", time: "2 days ago", icon: TrendingUp }
  ];

  const notifications = [
    { message: "New order #1234 requires attention", time: "5 min ago", type: "order" },
    { message: "Low stock alert for iPhone cases", time: "1 hour ago", type: "inventory" },
    { message: "Customer review pending approval", time: "3 hours ago", type: "review" }
  ];

  return (
    <div className="container mx-auto p-6 space-y-8">
      {/* Header Section */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
        <div className="flex items-center gap-6">
          <div className="relative">
            <Avatar className="w-20 h-20 ring-4 ring-background shadow-lg">
              <AvatarImage src="/vite.svg" alt="Admin User" />
              <AvatarFallback className="text-2xl font-bold">AU</AvatarFallback>
            </Avatar>
            <Button 
              size="icon" 
              variant="secondary" 
              className="absolute -bottom-2 -right-2 h-8 w-8 rounded-full cursor-pointer"
            >
              <Camera className="h-4 w-4" />
            </Button>
          </div>
          <div className="space-y-1">
            <h1 className="text-3xl font-bold">Admin User</h1>
            <p className="text-muted-foreground flex items-center gap-2">
              <Shield className="h-4 w-4" />
              Administrator • Store Owner
            </p>
            <div className="flex items-center gap-4 text-sm text-muted-foreground">
              <span className="flex items-center gap-1">
                <Mail className="h-4 w-4" />
                admin@shopvendor.com
              </span>
              <span className="flex items-center gap-1">
                <MapPin className="h-4 w-4" />
                San Francisco, CA
              </span>
            </div>
          </div>
        </div>
        <div className="flex gap-3">
          <Button 
            variant={isEditing ? "destructive" : "outline"} 
            onClick={() => setIsEditing(!isEditing)}
            className="cursor-pointer"
          >
            {isEditing ? (
              <>
                <X className="h-4 w-4 mr-2" />
                Cancel
              </>
            ) : (
              <>
                <Edit3 className="h-4 w-4 mr-2" />
                Edit Profile
              </>
            )}
          </Button>
          <Button 
            className="cursor-pointer"
            onClick={() => navigate('/settings')}
          >
            <Settings className="h-4 w-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid w-full grid-cols-2 lg:grid-cols-4">
          <TabsTrigger value="overview" className="cursor-pointer">Overview</TabsTrigger>
          <TabsTrigger value="profile" className="cursor-pointer">Profile</TabsTrigger>
          <TabsTrigger value="activity" className="cursor-pointer">Activity</TabsTrigger>
          <TabsTrigger value="settings" className="cursor-pointer">Settings</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          {/* Stats Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {stats.map((stat) => (
              <Card key={stat.label} className="hover:shadow-md transition-shadow cursor-pointer">
                <CardContent className="pt-6">
                  <div className="flex items-center justify-between">
                    <div className="space-y-2">
                      <p className="text-sm font-medium text-muted-foreground">{stat.label}</p>
                      <div className="flex items-baseline gap-2">
                        <p className="text-2xl font-bold">{stat.value}</p>
                        <span className="text-sm text-green-600 flex items-center gap-1">
                          <TrendingUp className="h-3 w-3" />
                          {stat.trend}
                        </span>
                      </div>
                    </div>
                    <div className={`p-3 rounded-full bg-muted ${stat.color}`}>
                      <stat.icon className="h-6 w-6" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Recent Activity */}
            <Card className="lg:col-span-2">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Activity className="h-5 w-5" />
                  Recent Activity
                </CardTitle>
                <CardDescription>Your latest actions in the dashboard</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {recentActivity.map((activity, index) => (
                  <div key={index} className="flex items-center gap-4 p-3 rounded-lg hover:bg-muted/50 transition-colors cursor-pointer">
                    <div className="p-2 rounded-full bg-primary/10">
                      <activity.icon className="h-4 w-4 text-primary" />
                    </div>
                    <div className="flex-1">
                      <p className="font-medium">{activity.action}</p>
                      <p className="text-sm text-muted-foreground flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        {activity.time}
                      </p>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>

            {/* Notifications & Quick Stats */}
            <div className="space-y-6">
              {/* Notifications */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Bell className="h-5 w-5" />
                    Notifications
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  {notifications.map((notification, index) => (
                    <div key={index} className="p-3 rounded-lg border border-border hover:bg-muted/50 transition-colors cursor-pointer">
                      <p className="text-sm font-medium">{notification.message}</p>
                      <p className="text-xs text-muted-foreground mt-1">{notification.time}</p>
                    </div>
                  ))}
                  <Button 
                    variant="outline" 
                    className="w-full cursor-pointer"
                    onClick={() => navigate('/notifications')}
                  >
                    View All Notifications
                  </Button>
                </CardContent>
              </Card>

              {/* Profile Completion */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Profile Completion</CardTitle>
                  <CardDescription>Complete your profile to unlock all features</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Progress</span>
                      <span>85%</span>
                    </div>
                    <Progress value={85} className="h-2" />
                  </div>
                  <div className="space-y-2 text-sm">
                    <div className="flex items-center justify-between">
                      <span>✓ Basic Info</span>
                      <Badge variant="secondary">Complete</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <span>✓ Contact Details</span>
                      <Badge variant="secondary">Complete</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <span>• Two-Factor Auth</span>
                      <Badge variant="outline">Pending</Badge>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>

        {/* Profile Tab */}
        <TabsContent value="profile" className="space-y-6">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Personal Information */}
              <Card>
                <CardHeader>
                  <CardTitle>Personal Information</CardTitle>
                  <CardDescription>Update your personal details and bio</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Full Name</label>
                    <Input 
                      {...register("name", { required: "Name is required" })} 
                      disabled={!isEditing}
                      className="cursor-pointer"
                    />
                    {errors.name && <p className="text-sm text-destructive">{errors.name.message}</p>}
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Email</label>
                    <Input 
                      type="email" 
                      {...register("email", { required: "Email is required" })} 
                      disabled={!isEditing}
                      className="cursor-pointer"
                    />
                    {errors.email && <p className="text-sm text-destructive">{errors.email.message}</p>}
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Phone</label>
                    <Input 
                      type="tel" 
                      {...register("phone")} 
                      disabled={!isEditing}
                      className="cursor-pointer"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Bio</label>
                    <Textarea 
                      {...register("bio")} 
                      disabled={!isEditing}
                      className="min-h-[100px] cursor-pointer"
                    />
                  </div>
                </CardContent>
              </Card>

              {/* Business Information */}
              <Card>
                <CardHeader>
                  <CardTitle>Business Information</CardTitle>
                  <CardDescription>Your company and business details</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Company</label>
                    <Input 
                      {...register("company")} 
                      disabled={!isEditing}
                      className="cursor-pointer"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Website</label>
                    <Input 
                      type="url" 
                      {...register("website")} 
                      disabled={!isEditing}
                      className="cursor-pointer"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Location</label>
                    <Input 
                      {...register("location")} 
                      disabled={!isEditing}
                      className="cursor-pointer"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Timezone</label>
                    <Input 
                      {...register("timezone")} 
                      disabled={!isEditing}
                      className="cursor-pointer"
                    />
                  </div>
                </CardContent>
              </Card>
            </div>

            {isEditing && (
              <div className="flex justify-end gap-3">
                <Button 
                  type="button" 
                  variant="outline" 
                  onClick={() => setIsEditing(false)}
                  className="cursor-pointer"
                >
                  Cancel
                </Button>
                <Button type="submit" className="cursor-pointer">
                  <Save className="h-4 w-4 mr-2" />
                  Save Changes
                </Button>
              </div>
            )}
          </form>
        </TabsContent>

        {/* Activity Tab */}
        <TabsContent value="activity" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <Card className="lg:col-span-2">
              <CardHeader>
                <CardTitle>Activity Timeline</CardTitle>
                <CardDescription>Detailed timeline of all your actions</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                {[...recentActivity, ...recentActivity].map((activity, index) => (
                  <div key={index} className="flex gap-4">
                    <div className="flex flex-col items-center">
                      <div className="p-2 rounded-full bg-primary/10">
                        <activity.icon className="h-4 w-4 text-primary" />
                      </div>
                      {index < recentActivity.length * 2 - 1 && (
                        <div className="w-px h-8 bg-border mt-2" />
                      )}
                    </div>
                    <div className="flex-1 pb-4">
                      <p className="font-medium">{activity.action}</p>
                      <p className="text-sm text-muted-foreground">{activity.time}</p>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Activity Summary</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-sm">This Week</span>
                    <Badge>47 actions</Badge>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm">This Month</span>
                    <Badge variant="secondary">203 actions</Badge>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm">Most Active</span>
                    <Badge variant="outline">Orders</Badge>
                  </div>
                </div>
                <Separator />
                <Button 
                  variant="outline" 
                  className="w-full cursor-pointer"
                  onClick={() => navigate('/logs')}
                >
                  View Full Activity Log
                </Button>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Settings Tab */}
        <TabsContent value="settings" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Notification Settings */}
            <Card>
              <CardHeader>
                <CardTitle>Notification Preferences</CardTitle>
                <CardDescription>Choose what notifications you want to receive</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="space-y-1">
                    <p className="font-medium">Email Notifications</p>
                    <p className="text-sm text-muted-foreground">Receive notifications via email</p>
                  </div>
                  <Switch defaultChecked className="cursor-pointer" />
                </div>
                <div className="flex items-center justify-between">
                  <div className="space-y-1">
                    <p className="font-medium">Order Alerts</p>
                    <p className="text-sm text-muted-foreground">Get notified about new orders</p>
                  </div>
                  <Switch defaultChecked className="cursor-pointer" />
                </div>
                <div className="flex items-center justify-between">
                  <div className="space-y-1">
                    <p className="font-medium">Low Stock Alerts</p>
                    <p className="text-sm text-muted-foreground">Alert when products are low in stock</p>
                  </div>
                  <Switch defaultChecked className="cursor-pointer" />
                </div>
                <div className="flex items-center justify-between">
                  <div className="space-y-1">
                    <p className="font-medium">Marketing Updates</p>
                    <p className="text-sm text-muted-foreground">Receive marketing and feature updates</p>
                  </div>
                  <Switch className="cursor-pointer" />
                </div>
              </CardContent>
            </Card>

            {/* Security Settings */}
            <Card>
              <CardHeader>
                <CardTitle>Security Settings</CardTitle>
                <CardDescription>Manage your account security</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-3">
                  <Button 
                    variant="outline" 
                    className="w-full justify-start cursor-pointer"
                    onClick={() => navigate('/settings/password')}
                  >
                    <Key className="h-4 w-4 mr-2" />
                    Change Password
                  </Button>
                  <Button 
                    variant="outline" 
                    className="w-full justify-start cursor-pointer"
                    onClick={() => navigate('/settings/2fa')}
                  >
                    <Smartphone className="h-4 w-4 mr-2" />
                    Two-Factor Authentication
                  </Button>
                  <Button 
                    variant="outline" 
                    className="w-full justify-start cursor-pointer"
                    onClick={() => navigate('/settings/sessions')}
                  >
                    <Globe className="h-4 w-4 mr-2" />
                    Active Sessions
                  </Button>
                </div>
                <Separator />
                <div className="pt-2">
                  <div className="flex items-center justify-between">
                    <div className="space-y-1">
                      <p className="font-medium">Two-Factor Authentication</p>
                      <p className="text-sm text-muted-foreground">Additional security for your account</p>
                    </div>
                    <Switch className="cursor-pointer" />
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}