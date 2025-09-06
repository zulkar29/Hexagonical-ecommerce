import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { 
  ArrowLeft, 
  Phone, 
  Globe, 
  MapPin, 
  Calendar,
  DollarSign,
  TrendingUp,
  Users,
  ExternalLink,
  Edit,
  CheckCircle,
  XCircle,
  Clock,
  BarChart3
} from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Separator } from "@/components/ui/separator";

const AffiliateDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [affiliate, setAffiliate] = useState(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("overview");

  // Mock data - replace with actual API call
  useEffect(() => {
    const fetchAffiliate = async () => {
      try {
        // TODO: Replace with actual API call
        // const response = await affiliatesApi.getById(id);
        
        // Mock data
        const mockAffiliate = {
          id: parseInt(id),
          name: "John Smith",
          email: "john@example.com",
          avatar: null,
          status: "active",
          joinDate: "2024-01-15",
          phone: "+1 (555) 123-4567",
          country: "United States",
          state: "California",
          city: "San Francisco",
          website: "https://johnsmith.com",
          businessName: "Smith Marketing",
          niche: "Technology",
          audienceSize: "10k-50k",
          socialMedia: "@johnsmith_tech",
          totalSales: 45678.90,
          commission: 4567.89,
          referrals: 156,
          conversionRate: 4.2,
          monthlyEarnings: 1234.56,
          links: [
            {
              id: 1,
              name: "Homepage Banner",
              url: "https://store.com/ref/john123",
              clicks: 1234,
              conversions: 45,
              revenue: 2345.67,
              createdAt: "2024-01-20"
            },
            {
              id: 2,
              name: "Product Review Link",
              url: "https://store.com/product/abc?ref=john123",
              clicks: 856,
              conversions: 23,
              revenue: 1567.89,
              createdAt: "2024-02-10"
            }
          ],
          recentCommissions: [
            {
              id: 1,
              orderId: "ORD-001",
              amount: 156.78,
              commission: 15.68,
              date: "2024-03-15",
              status: "paid"
            },
            {
              id: 2,
              orderId: "ORD-002",
              amount: 234.56,
              commission: 23.46,
              date: "2024-03-14",
              status: "pending"
            }
          ]
        };
        
        setAffiliate(mockAffiliate);
      } catch (error) {
        toast.error("Failed to load affiliate details");
        console.error(error);
      } finally {
        setLoading(false);
      }
    };

    fetchAffiliate();
  }, [id]);

  const getStatusIcon = (status) => {
    switch (status) {
      case 'active':
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      case 'pending':
        return <Clock className="h-4 w-4 text-yellow-500" />;
      case 'inactive':
        return <XCircle className="h-4 w-4 text-red-500" />;
      default:
        return null;
    }
  };

  const getStatusVariant = (status) => {
    switch (status) {
      case 'active':
        return 'default';
      case 'pending':
        return 'secondary';
      case 'inactive':
        return 'destructive';
      default:
        return 'secondary';
    }
  };

  const handleStatusUpdate = async (newStatus) => {
    try {
      // TODO: Implement API call
      // await affiliatesApi.updateStatus(id, newStatus);
      setAffiliate(prev => ({ ...prev, status: newStatus }));
      toast.success(`Affiliate status updated to ${newStatus}`);
    } catch (error) {
      toast.error("Failed to update status");
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
          <p className="mt-2 text-muted-foreground">Loading affiliate details...</p>
        </div>
      </div>
    );
  }

  if (!affiliate) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <p className="text-muted-foreground">Affiliate not found</p>
          <Button onClick={() => navigate("/affiliates")} className="mt-4">
            Back to Affiliates
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6 p-6 max-w-7xl mx-auto">
      {/* Header Card */}
      <Card>
        <CardHeader className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 pb-4">
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="icon"
              className="rounded-full"
              onClick={() => navigate("/affiliates")}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div className="flex items-center gap-4">
              <Avatar className="h-16 w-16">
                <AvatarImage src={affiliate.avatar} />
                <AvatarFallback className="text-lg">
                  {affiliate.name.split(' ').map(n => n[0]).join('')}
                </AvatarFallback>
              </Avatar>
              <div>
                <h1 className="text-2xl font-bold">{affiliate.name}</h1>
                <p className="text-muted-foreground">{affiliate.email}</p>
                <div className="flex items-center gap-2 mt-1">
                  {getStatusIcon(affiliate.status)}
                  <Badge variant={getStatusVariant(affiliate.status)} className="capitalize">
                    {affiliate.status}
                  </Badge>
                </div>
              </div>
            </div>
          </div>
          <div className="flex gap-2">
            {affiliate.status === 'pending' && (
              <Button onClick={() => handleStatusUpdate('active')}>
                <CheckCircle className="mr-2 h-4 w-4" />
                Approve
              </Button>
            )}
            {affiliate.status === 'active' && (
              <Button 
                variant="outline" 
                onClick={() => handleStatusUpdate('inactive')}
              >
                <XCircle className="mr-2 h-4 w-4" />
                Deactivate
              </Button>
            )}
            <Button variant="outline">
              <Edit className="mr-2 h-4 w-4" />
              Edit
            </Button>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Contact & Business Info */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="space-y-3">
              <h3 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">Contact Info</h3>
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">{affiliate.phone || 'Not provided'}</span>
                </div>
                <div className="flex items-center gap-2">
                  <MapPin className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">
                    {[affiliate.city, affiliate.state, affiliate.country].filter(Boolean).join(', ')}
                  </span>
                </div>
              </div>
            </div>
            
            <div className="space-y-3">
              <h3 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">Business Info</h3>
              <div className="space-y-2">
                <div>
                  <p className="text-xs text-muted-foreground">Business Name</p>
                  <p className="text-sm font-medium">{affiliate.businessName || 'Not provided'}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Niche</p>
                  <p className="text-sm font-medium">{affiliate.niche}</p>
                </div>
              </div>
            </div>
            
            <div className="space-y-3">
              <h3 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">Online Presence</h3>
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Globe className="h-4 w-4 text-muted-foreground" />
                  <a 
                    href={affiliate.website} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className="text-sm text-primary hover:underline flex items-center gap-1"
                  >
                    Website <ExternalLink className="h-3 w-3" />
                  </a>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Social Media</p>
                  <p className="text-sm font-medium">{affiliate.socialMedia || 'Not provided'}</p>
                </div>
              </div>
            </div>
            
            <div className="space-y-3">
              <h3 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">Partnership</h3>
              <div className="space-y-2">
                <div>
                  <p className="text-xs text-muted-foreground">Join Date</p>
                  <p className="text-sm font-medium">{formatDate(affiliate.joinDate)}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Audience Size</p>
                  <p className="text-sm font-medium">{affiliate.audienceSize}</p>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Performance Stats */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Sales</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(affiliate.totalSales)}</div>
            <p className="text-xs text-muted-foreground">
              Lifetime sales volume
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Commission</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(affiliate.commission)}</div>
            <p className="text-xs text-muted-foreground">
              Lifetime earnings
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Referrals</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{affiliate.referrals}</div>
            <p className="text-xs text-muted-foreground">
              Total customers referred
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Conversion Rate</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{affiliate.conversionRate}%</div>
            <p className="text-xs text-muted-foreground">
              Click to sale conversion
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Detailed Information Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid w-full grid-cols-3 lg:grid-cols-3">
          <TabsTrigger value="activity">Recent Activity</TabsTrigger>
          <TabsTrigger value="links">Affiliate Links</TabsTrigger>
          <TabsTrigger value="commissions">Commission History</TabsTrigger>
        </TabsList>

        <TabsContent value="activity" className="space-y-4">
          <div className="grid gap-4 lg:grid-cols-2">
            {/* Recent Transactions */}
            <Card>
              <CardHeader>
                <CardTitle>Recent Transactions</CardTitle>
                <CardDescription>Latest commission transactions</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {affiliate.recentCommissions.map((commission) => (
                    <div key={commission.id} className="flex items-center justify-between p-3 bg-muted/30 rounded-lg">
                      <div>
                        <p className="font-medium text-sm">{commission.orderId}</p>
                        <p className="text-xs text-muted-foreground">{formatDate(commission.date)}</p>
                      </div>
                      <div className="text-right">
                        <p className="font-medium text-sm">{formatCurrency(commission.commission)}</p>
                        <Badge variant={commission.status === 'paid' ? 'default' : 'secondary'} className="text-xs">
                          {commission.status}
                        </Badge>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Monthly Performance */}
            <Card>
              <CardHeader>
                <CardTitle>Monthly Performance</CardTitle>
                <CardDescription>Current month statistics</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Monthly Earnings</span>
                    <span className="font-semibold">{formatCurrency(affiliate.monthlyEarnings)}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">New Referrals</span>
                    <span className="font-semibold">12</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Conversion Rate</span>
                    <span className="font-semibold">{affiliate.conversionRate}%</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Click-through Rate</span>
                    <span className="font-semibold">6.8%</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="links" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Affiliate Links Performance</CardTitle>
              <CardDescription>Track performance of affiliate tracking links</CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Link Name</TableHead>
                    <TableHead>URL</TableHead>
                    <TableHead>Clicks</TableHead>
                    <TableHead>Conversions</TableHead>
                    <TableHead>Revenue</TableHead>
                    <TableHead>Created</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {affiliate.links.map((link) => (
                    <TableRow key={link.id}>
                      <TableCell className="font-medium">{link.name}</TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <code className="text-xs bg-muted px-2 py-1 rounded max-w-[200px] truncate">
                            {link.url}
                          </code>
                          <Button size="icon" variant="ghost" className="h-6 w-6">
                            <ExternalLink className="h-3 w-3" />
                          </Button>
                        </div>
                      </TableCell>
                      <TableCell>{link.clicks.toLocaleString()}</TableCell>
                      <TableCell>{link.conversions}</TableCell>
                      <TableCell>{formatCurrency(link.revenue)}</TableCell>
                      <TableCell>{formatDate(link.createdAt)}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="commissions" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Complete Commission History</CardTitle>
              <CardDescription>All commission transactions and payments</CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Order ID</TableHead>
                    <TableHead>Date</TableHead>
                    <TableHead>Sale Amount</TableHead>
                    <TableHead>Commission</TableHead>
                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {affiliate.recentCommissions.map((commission) => (
                    <TableRow key={commission.id}>
                      <TableCell className="font-medium font-mono text-sm">{commission.orderId}</TableCell>
                      <TableCell>{formatDate(commission.date)}</TableCell>
                      <TableCell>{formatCurrency(commission.amount)}</TableCell>
                      <TableCell>{formatCurrency(commission.commission)}</TableCell>
                      <TableCell>
                        <Badge variant={commission.status === 'paid' ? 'default' : 'secondary'}>
                          {commission.status}
                        </Badge>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default AffiliateDetails;