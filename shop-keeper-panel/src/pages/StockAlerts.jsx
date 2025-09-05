import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  AlertTriangle,
  AlertCircle,
  XCircle,
  CheckCircle,
  Package,
  TrendingDown,
  TrendingUp,
  Clock,
  Bell,
  BellOff,
  Save,
  Search,
  Download,
  Upload,
  RefreshCw,
  Settings,
  Plus,
  Minus,
  Eye,
  Edit,
  MoreVertical,
  Trash2,
  Archive,
  Star,
  Package2,
  ShoppingCart,
  Calendar,
  Users,
  Building,
  Phone,
  Mail,
  ExternalLink,
  FileText,
  BarChart3,
  Target,
  Zap,
  Activity
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Checkbox } from '@/components/ui/checkbox';
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from '@/components/ui/dialog';
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';

const StockAlertsPage = () => {
  const navigate = useNavigate();
  const [search, setSearch] = useState('');
  const [selectedAlerts, setSelectedAlerts] = useState([]);
  const [alertTypeFilter, setAlertTypeFilter] = useState('all');
  const [priorityFilter, setPriorityFilter] = useState('all');
  const [statusFilter, setStatusFilter] = useState('all');
  const [showBulkActions, setShowBulkActions] = useState(false);
  const [showSettingsDialog, setShowSettingsDialog] = useState(false);
  const [showRestockDialog, setShowRestockDialog] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [restockQuantity, setRestockQuantity] = useState('');
  const [restockNotes, setRestockNotes] = useState('');
  const [itemsPerPage, setItemsPerPage] = useState(15);
  const [currentPage, setCurrentPage] = useState(1);
  const [alertSettings, setAlertSettings] = useState({
    lowStockEnabled: true,
    outOfStockEnabled: true,
    overStockEnabled: false,
    expiryEnabled: true,
    emailNotifications: true,
    smsNotifications: false,
    defaultLowStockThreshold: 10
  });

  // Mock stock alerts data
  const mockAlerts = [
    {
      id: 1,
      productId: 1,
      productName: 'Basmati Rice 5kg',
      category: 'Rice & Grains',
      currentStock: 3,
      minStock: 10,
      maxStock: 200,
      alertType: 'low_stock',
      priority: 'high',
      status: 'active',
      createdAt: '2024-07-24T08:00:00Z',
      lastRestocked: '2024-07-10',
      supplier: 'Rahman Traders',
      supplierPhone: '01712345678',
      supplierEmail: 'contact@rahmantraders.com',
      avgDailySales: 2.5,
      daysUntilEmpty: 1,
      reorderPoint: 15,
      suggestedOrderQty: 50,
      unitPrice: 500,
      image: '/api/placeholder/40/40',
      locationInStore: 'A-1-001'
    },
    {
      id: 2,
      productId: 2,
      productName: 'Onion 1kg',
      category: 'Vegetables',
      currentStock: 0,
      minStock: 10,
      maxStock: 100,
      alertType: 'out_of_stock',
      priority: 'critical',
      status: 'active',
      createdAt: '2024-07-23T14:30:00Z',
      lastRestocked: '2024-07-15',
      supplier: 'Green Agro',
      supplierPhone: '01798765432',
      supplierEmail: 'info@greenagro.com',
      avgDailySales: 8.2,
      daysUntilEmpty: 0,
      reorderPoint: 20,
      suggestedOrderQty: 80,
      unitPrice: 100,
      image: '/api/placeholder/40/40',
      locationInStore: 'B-2-005'
    },
    {
      id: 3,
      productId: 3,
      productName: 'Sugar 1kg',
      category: 'Pantry',
      currentStock: 8,
      minStock: 10,
      maxStock: 100,
      alertType: 'low_stock',
      priority: 'medium',
      status: 'active',
      createdAt: '2024-07-22T10:15:00Z',
      lastRestocked: '2024-07-08',
      supplier: 'Sweet Supply Co.',
      supplierPhone: '01912345678',
      supplierEmail: 'orders@sweetco.com',
      avgDailySales: 1.8,
      daysUntilEmpty: 4,
      reorderPoint: 15,
      suggestedOrderQty: 40,
      unitPrice: 120,
      image: '/api/placeholder/40/40',
      locationInStore: 'C-1-010'
    },
    {
      id: 4,
      productId: 4,
      productName: 'Cooking Oil 1L',
      category: 'Cooking',
      currentStock: 5,
      minStock: 8,
      maxStock: 60,
      alertType: 'low_stock',
      priority: 'high',
      status: 'acknowledged',
      createdAt: '2024-07-21T16:45:00Z',
      lastRestocked: '2024-07-12',
      supplier: 'Oil Mart',
      supplierPhone: '01812345678',
      supplierEmail: 'procurement@oilmart.bd',
      avgDailySales: 1.2,
      daysUntilEmpty: 4,
      reorderPoint: 12,
      suggestedOrderQty: 30,
      unitPrice: 150,
      image: '/api/placeholder/40/40',
      locationInStore: 'D-3-002'
    },
    {
      id: 5,
      productId: 5,
      productName: 'Milk 1L',
      category: 'Dairy',
      currentStock: 220,
      minStock: 15,
      maxStock: 80,
      alertType: 'overstock',
      priority: 'low',
      status: 'active',
      createdAt: '2024-07-20T09:20:00Z',
      lastRestocked: '2024-07-19',
      supplier: 'Milk Fresh Ltd.',
      supplierPhone: '01656789012',
      supplierEmail: 'orders@milkfresh.bd',
      avgDailySales: 12.5,
      daysUntilEmpty: 18,
      reorderPoint: 25,
      suggestedOrderQty: 0,
      unitPrice: 100,
      image: '/api/placeholder/40/40',
      locationInStore: 'E-1-003'
    },
    {
      id: 6,
      productId: 6,
      productName: 'Potato 1kg',
      category: 'Vegetables',
      currentStock: 2,
      minStock: 10,
      maxStock: 150,
      alertType: 'low_stock',
      priority: 'high',
      status: 'resolved',
      createdAt: '2024-07-19T11:30:00Z',
      lastRestocked: '2024-07-18',
      supplier: 'Potato House',
      supplierPhone: '01534567890',
      supplierEmail: 'sales@potatohouse.com',
      avgDailySales: 6.8,
      daysUntilEmpty: 0,
      reorderPoint: 20,
      suggestedOrderQty: 60,
      unitPrice: 80,
      image: '/api/placeholder/40/40',
      locationInStore: 'B-1-008'
    }
  ];

  // Filter and sort alerts
  const filteredAlerts = useMemo(() => {
    return mockAlerts.filter(alert => {
      const matchesSearch = alert.productName.toLowerCase().includes(search.toLowerCase()) ||
                           alert.category.toLowerCase().includes(search.toLowerCase()) ||
                           alert.supplier.toLowerCase().includes(search.toLowerCase());

      const matchesType = alertTypeFilter === 'all' || alert.alertType === alertTypeFilter;
      const matchesPriority = priorityFilter === 'all' || alert.priority === priorityFilter;
      const matchesStatus = statusFilter === 'all' || alert.status === statusFilter;

      return matchesSearch && matchesType && matchesPriority && matchesStatus;
    });
  }, [mockAlerts, search, alertTypeFilter, priorityFilter, statusFilter]);

  // Pagination
  const totalPages = Math.ceil(filteredAlerts.length / itemsPerPage);
  const paginatedAlerts = filteredAlerts.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  // Stats calculations
  const stats = useMemo(() => {
    const totalAlerts = mockAlerts.length;
    const activeAlerts = mockAlerts.filter(a => a.status === 'active').length;
    const criticalAlerts = mockAlerts.filter(a => a.priority === 'critical').length;
    const outOfStockAlerts = mockAlerts.filter(a => a.alertType === 'out_of_stock').length;
    const lowStockAlerts = mockAlerts.filter(a => a.alertType === 'low_stock').length;
    const overStockAlerts = mockAlerts.filter(a => a.alertType === 'overstock').length;
    const totalValue = mockAlerts.reduce((sum, a) => sum + (a.currentStock * a.unitPrice), 0);
    const requiredValue = mockAlerts.reduce((sum, a) => {
      if (a.alertType === 'low_stock' || a.alertType === 'out_of_stock') {
        return sum + (a.suggestedOrderQty * a.unitPrice);
      }
      return sum;
    }, 0);

    return {
      totalAlerts,
      activeAlerts,
      criticalAlerts,
      outOfStockAlerts,
      lowStockAlerts,
      overStockAlerts,
      totalValue,
      requiredValue
    };
  }, [mockAlerts]);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  const formatDate = (dateString) => new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

  const getAlertTypeInfo = (type) => {
    switch (type) {
      case 'out_of_stock':
        return { 
          icon: XCircle, 
          color: 'destructive', 
          label: 'Out of Stock',
          bgColor: 'bg-red-50 border-red-200',
          iconColor: 'text-red-600'
        };
      case 'low_stock':
        return { 
          icon: AlertTriangle, 
          color: 'warning', 
          label: 'Low Stock',
          bgColor: 'bg-yellow-50 border-yellow-200',
          iconColor: 'text-yellow-600'
        };
      case 'overstock':
        return { 
          icon: TrendingUp, 
          color: 'secondary', 
          label: 'Overstock',
          bgColor: 'bg-blue-50 border-blue-200',  
          iconColor: 'text-blue-600'
        };
      default:
        return { 
          icon: AlertCircle, 
          color: 'default', 
          label: 'Alert',
          bgColor: 'bg-gray-50 border-gray-200',
          iconColor: 'text-gray-600'
        };
    }
  };

  const getPriorityInfo = (priority) => {
    switch (priority) {
      case 'critical':
        return { color: 'destructive', label: 'Critical' };
      case 'high':
        return { color: 'warning', label: 'High' };
      case 'medium':
        return { color: 'default', label: 'Medium' };
      case 'low':
        return { color: 'secondary', label: 'Low' };
      default:
        return { color: 'default', label: 'Normal' };
    }
  };

  const handleSelectAlert = (alertId) => {
    setSelectedAlerts(prev =>
      prev.includes(alertId)
        ? prev.filter(id => id !== alertId)
        : [...prev, alertId]
    );
  };

  const handleSelectAll = () => {
    if (selectedAlerts.length === paginatedAlerts.length) {
      setSelectedAlerts([]);
    } else {
      setSelectedAlerts(paginatedAlerts.map(a => a.id));
    }
  };

  const handleBulkAction = (action) => {
    console.log(`Bulk action: ${action} on alerts:`, selectedAlerts);
    setSelectedAlerts([]);
  };

  const handleRestock = () => {
    console.log('Restocking product:', {
      product: selectedProduct,
      quantity: restockQuantity,
      notes: restockNotes
    });
    setRestockQuantity('');
    setRestockNotes('');
    setSelectedProduct(null);
    setShowRestockDialog(false);
  };

  const openRestockDialog = (alert) => {
    setSelectedProduct(alert);
    setRestockQuantity(alert.suggestedOrderQty.toString());
    setShowRestockDialog(true);
  };

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto p-6 space-y-6">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold flex items-center">
              <AlertTriangle className="h-8 w-8 mr-3 text-orange-600" />
              Stock Alerts
            </h1>
            <p className="text-muted-foreground mt-1">
              Monitor and manage your inventory alerts
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm" onClick={() => setShowSettingsDialog(true)}>
              <Settings className="h-4 w-4 mr-2" />
              Settings
            </Button>
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
            <Button>
              <RefreshCw className="h-4 w-4 mr-2" />
              Refresh Alerts
            </Button>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-8 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Total Alerts</p>
                  <p className="text-2xl font-bold">{stats.totalAlerts}</p>
                </div>
                <Bell className="h-8 w-8 text-primary" />
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Active</p>
                  <p className="text-2xl font-bold text-orange-600">{stats.activeAlerts}</p>
                </div>
                <Activity className="h-8 w-8 text-orange-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Critical</p>
                  <p className="text-2xl font-bold text-red-600">{stats.criticalAlerts}</p>
                </div>
                <Zap className="h-8 w-8 text-red-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Out of Stock</p>
                  <p className="text-2xl font-bold text-red-600">{stats.outOfStockAlerts}</p>
                </div>
                <XCircle className="h-8 w-8 text-red-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Low Stock</p>
                  <p className="text-2xl font-bold text-yellow-600">{stats.lowStockAlerts}</p>
                </div>
                <AlertTriangle className="h-8 w-8 text-yellow-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Overstock</p>
                  <p className="text-2xl font-bold text-blue-600">{stats.overStockAlerts}</p>
                </div>
                <TrendingUp className="h-8 w-8 text-blue-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Current Value</p>
                  <p className="text-lg font-bold">{formatCurrency(stats.totalValue)}</p>
                </div>
                <BarChart3 className="h-8 w-8 text-green-600" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Restock Value</p>
                  <p className="text-lg font-bold">{formatCurrency(stats.requiredValue)}</p>
                </div>
                <Target className="h-8 w-8 text-purple-600" />
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Filters and Controls */}
        <Card>
          <CardContent className="p-6">
            <div className="flex flex-col lg:flex-row gap-4 mb-6">
              {/* Search */}
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  placeholder="Search products, categories, or suppliers..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-10"
                />
              </div>

              {/* Filters */}
              <div className="flex flex-wrap gap-2">
                <Select value={alertTypeFilter} onValueChange={setAlertTypeFilter}>
                  <SelectTrigger className="w-40">
                    <SelectValue placeholder="Alert Type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Types</SelectItem>
                    <SelectItem value="out_of_stock">Out of Stock</SelectItem>
                    <SelectItem value="low_stock">Low Stock</SelectItem>
                    <SelectItem value="overstock">Overstock</SelectItem>
                  </SelectContent>
                </Select>

                <Select value={priorityFilter} onValueChange={setPriorityFilter}>
                  <SelectTrigger className="w-32">
                    <SelectValue placeholder="Priority" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Priority</SelectItem>
                    <SelectItem value="critical">Critical</SelectItem>
                    <SelectItem value="high">High</SelectItem>
                    <SelectItem value="medium">Medium</SelectItem>
                    <SelectItem value="low">Low</SelectItem>
                  </SelectContent>
                </Select>

                <Select value={statusFilter} onValueChange={setStatusFilter}>
                  <SelectTrigger className="w-32">
                    <SelectValue placeholder="Status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Status</SelectItem>
                    <SelectItem value="active">Active</SelectItem>
                    <SelectItem value="acknowledged">Acknowledged</SelectItem>
                    <SelectItem value="resolved">Resolved</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Bulk Actions */}
            {selectedAlerts.length > 0 && (
              <div className="flex items-center justify-between mb-4 p-3 bg-primary/10 rounded-lg">
                <span className="text-sm font-medium">
                  {selectedAlerts.length} alert{selectedAlerts.length > 1 ? 's' : ''} selected
                </span>
                <div className="flex items-center gap-2">
                  <Button size="sm" variant="outline" onClick={() => handleBulkAction('acknowledge')}>
                    <CheckCircle className="h-4 w-4 mr-1" />
                    Acknowledge
                  </Button>
                  <Button size="sm" variant="outline" onClick={() => handleBulkAction('resolve')}>
                    <CheckCircle className="h-4 w-4 mr-1" />
                    Resolve
                  </Button>
                  <Button size="sm" variant="outline" onClick={() => handleBulkAction('export')}>
                    <Download className="h-4 w-4 mr-1" />
                    Export
                  </Button>
                  <Button size="sm" variant="ghost" onClick={() => setSelectedAlerts([])}>
                    <XCircle className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            )}

            {/* Alerts Table */}
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-12">
                      <Checkbox
                        checked={selectedAlerts.length === paginatedAlerts.length}
                        onCheckedChange={handleSelectAll}
                      />
                    </TableHead>
                    <TableHead>Product</TableHead>
                    <TableHead>Alert Type</TableHead>
                    <TableHead>Priority</TableHead>
                    <TableHead>Current Stock</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Days Until Empty</TableHead>
                    <TableHead>Supplier</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {paginatedAlerts.map((alert) => {
                    const alertType = getAlertTypeInfo(alert.alertType);
                    const priority = getPriorityInfo(alert.priority);
                    const AlertIcon = alertType.icon;
                    
                    return (
                      <TableRow key={alert.id} className={cn("hover:bg-muted/50", alertType.bgColor)}>
                        <TableCell>
                          <Checkbox
                            checked={selectedAlerts.includes(alert.id)}
                            onCheckedChange={() => handleSelectAlert(alert.id)}
                          />
                        </TableCell>
                        
                        <TableCell>
                          <div className="flex items-center space-x-3">
                            <Avatar className="h-10 w-10">
                              <AvatarImage src={alert.image} alt={alert.productName} />
                              <AvatarFallback>
                                <Package className="h-4 w-4" />
                              </AvatarFallback>
                            </Avatar>
                            <div>
                              <p 
                                className="font-medium cursor-pointer hover:text-primary"
                                onClick={() => navigate(`/inventory/products/${alert.productId}`)}
                              >
                                {alert.productName}
                              </p>
                              <div className="flex items-center space-x-2 text-xs text-muted-foreground">
                                <span>{alert.category}</span>
                                <span>•</span>
                                <span>{alert.locationInStore}</span>
                              </div>
                            </div>
                          </div>
                        </TableCell>

                        <TableCell>
                          <div className="flex items-center space-x-2">
                            <AlertIcon className={cn("h-4 w-4", alertType.iconColor)} />
                            <Badge variant={alertType.color} className="text-xs">
                              {alertType.label}
                            </Badge>
                          </div>
                        </TableCell>

                        <TableCell>
                          <Badge variant={priority.color} className="text-xs">
                            {priority.label}
                          </Badge>
                        </TableCell>

                        <TableCell>
                          <div className="space-y-1">
                            <div className="flex items-center space-x-2">
                              <span className="font-medium">{alert.currentStock}</span>
                              <span className="text-xs text-muted-foreground">
                                / {alert.maxStock}
                              </span>
                            </div>
                            <Progress 
                              value={Math.min((alert.currentStock / alert.maxStock) * 100, 100)} 
                              className="h-1 w-16"
                            />
                            <div className="text-xs text-muted-foreground">
                              Min: {alert.minStock}
                            </div>
                          </div>
                        </TableCell>

                        <TableCell>
                          <Badge 
                            variant={
                              alert.status === 'active' ? 'default' :
                              alert.status === 'acknowledged' ? 'secondary' : 'outline'
                            }
                            className="text-xs"
                          >
                            {alert.status}
                          </Badge>
                        </TableCell>

                        <TableCell>
                          <div className="text-center">
                            <p className={cn(
                              "font-medium",
                              alert.daysUntilEmpty <= 1 ? "text-red-600" :
                              alert.daysUntilEmpty <= 3 ? "text-yellow-600" : "text-green-600"
                            )}>
                              {alert.daysUntilEmpty}
                            </p>
                            <p className="text-xs text-muted-foreground">days</p>
                          </div>
                        </TableCell>

                        <TableCell>
                          <div className="space-y-1">
                            <p className="text-sm font-medium">{alert.supplier}</p>
                            <div className="flex items-center space-x-2">
                              <Button 
                                variant="ghost" 
                                size="sm" 
                                className="h-6 w-6 p-0"
                                onClick={() => window.open(`tel:${alert.supplierPhone}`)}
                              >
                                <Phone className="h-3 w-3" />
                              </Button>
                              <Button 
                                variant="ghost" 
                                size="sm" 
                                className="h-6 w-6 p-0"
                                onClick={() => window.open(`mailto:${alert.supplierEmail}`)}
                              >
                                <Mail className="h-3 w-3" />
                              </Button>
                            </div>
                          </div>
                        </TableCell>

                        <TableCell>
                          <div className="flex items-center space-x-1">
                            {(alert.alertType === 'low_stock' || alert.alertType === 'out_of_stock') && (
                              <Button 
                                size="sm" 
                                onClick={() => openRestockDialog(alert)}
                                className="h-8 px-2"
                              >
                                <Plus className="h-3 w-3 mr-1" />
                                Restock
                              </Button>
                            )}
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button size="sm" variant="ghost" className="h-8 w-8 p-0">
                                  <MoreVertical className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end">
                                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                <DropdownMenuItem onClick={() => navigate(`/inventory/products/${alert.productId}`)}>
                                  <Eye className="h-4 w-4 mr-2" />
                                  View Product
                                </DropdownMenuItem>
                                <DropdownMenuItem>
                                  <CheckCircle className="h-4 w-4 mr-2" />
                                  Acknowledge Alert
                                </DropdownMenuItem>
                                <DropdownMenuItem>
                                  <Building className="h-4 w-4 mr-2" />
                                  Contact Supplier
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem>
                                  <Archive className="h-4 w-4 mr-2" />
                                  Archive Alert
                                </DropdownMenuItem>
                                <DropdownMenuItem className="text-destructive">
                                  <Trash2 className="h-4 w-4 mr-2" />
                                  Delete Alert
                                </DropdownMenuItem>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </div>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </div>

            {/* Pagination */}
            <div className="flex items-center justify-between mt-4">
              <div className="flex items-center space-x-2">
                <span className="text-sm text-muted-foreground">Show</span>
                <Select value={itemsPerPage.toString()} onValueChange={(value) => setItemsPerPage(parseInt(value))}>
                  <SelectTrigger className="w-20">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="10">10</SelectItem>
                    <SelectItem value="15">15</SelectItem>
                    <SelectItem value="25">25</SelectItem>
                    <SelectItem value="50">50</SelectItem>
                  </SelectContent>
                </Select>
                <span className="text-sm text-muted-foreground">
                  of {filteredAlerts.length} alerts
                </span>
              </div>

              <div className="flex items-center space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                  disabled={currentPage === 1}
                >
                  Previous
                </Button>
                <span className="text-sm">
                  Page {currentPage} of {totalPages}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                  disabled={currentPage === totalPages}
                >
                  Next
                </Button>
              </div>
            </div>

            {/* No Results */}
            {filteredAlerts.length === 0 && (
              <div className="text-center py-12">
                <CheckCircle className="h-12 w-12 mx-auto text-green-600 mb-4" />
                <h3 className="text-lg font-semibold mb-2">No alerts found</h3>
                <p className="text-muted-foreground mb-4">
                  {search || alertTypeFilter !== 'all' || priorityFilter !== 'all' || statusFilter !== 'all'
                    ? 'Try adjusting your search or filters'
                    : 'All your inventory levels are looking good!'
                  }
                </p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Restock Dialog */}
        <Dialog open={showRestockDialog} onOpenChange={setShowRestockDialog}>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>Restock Product</DialogTitle>
              <DialogDescription>
                {selectedProduct && `Add stock for ${selectedProduct.productName}`}
              </DialogDescription>
            </DialogHeader>
            
            {selectedProduct && (
              <div className="space-y-4">
                {/* Product Info */}
                <div className="flex items-center space-x-3 p-3 bg-muted/50 rounded-lg">
                  <Avatar className="h-12 w-12">
                    <AvatarImage src={selectedProduct.image} alt={selectedProduct.productName} />
                    <AvatarFallback>
                      <Package className="h-4 w-4" />
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-medium">{selectedProduct.productName}</p>
                    <p className="text-sm text-muted-foreground">
                      Current: {selectedProduct.currentStock} | Min: {selectedProduct.minStock}
                    </p>
                  </div>
                </div>

                {/* Supplier Info */}
                <div className="p-3 border rounded-lg">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium">Supplier</span>
                    <div className="flex items-center space-x-2">
                      <Button 
                        variant="ghost" 
                        size="sm" 
                        className="h-6 w-6 p-0"
                        onClick={() => window.open(`tel:${selectedProduct.supplierPhone}`)}
                      >
                        <Phone className="h-3 w-3" />
                      </Button>
                      <Button 
                        variant="ghost" 
                        size="sm" 
                        className="h-6 w-6 p-0"
                        onClick={() => window.open(`mailto:${selectedProduct.supplierEmail}`)}
                      >
                        <Mail className="h-3 w-3" />
                      </Button>
                    </div>
                  </div>
                  <p className="text-sm">{selectedProduct.supplier}</p>
                  <p className="text-xs text-muted-foreground">{selectedProduct.supplierPhone}</p>
                </div>

                {/* Restock Quantity */}
                <div>
                  <Label htmlFor="quantity" className="text-sm font-medium">
                    Restock Quantity
                  </Label>
                  <Input
                    id="quantity"
                    type="number"
                    value={restockQuantity}
                    onChange={(e) => setRestockQuantity(e.target.value)}
                    placeholder="Enter quantity to add"
                    className="mt-1"
                    min="1"
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    Suggested: {selectedProduct.suggestedOrderQty} units
                  </p>
                </div>

                {/* Cost Calculation */}
                {restockQuantity && (
                  <div className="p-3 bg-green-50 border border-green-200 rounded-lg">
                    <div className="flex justify-between items-center">
                      <span className="text-sm">Total Cost:</span>
                      <span className="font-semibold text-green-700">
                        {formatCurrency(parseInt(restockQuantity) * selectedProduct.unitPrice)}
                      </span>
                    </div>
                    <div className="flex justify-between items-center text-sm text-muted-foreground">
                      <span>New Stock Level:</span>
                      <span>{selectedProduct.currentStock + parseInt(restockQuantity || 0)} units</span>
                    </div>
                  </div>
                )}

                {/* Notes */}
                <div>
                  <Label htmlFor="notes" className="text-sm font-medium">
                    Notes (Optional)
                  </Label>
                  <Textarea
                    id="notes"
                    value={restockNotes}
                    onChange={(e) => setRestockNotes(e.target.value)}
                    placeholder="Add any notes about this restock..."
                    className="mt-1"
                    rows={3}
                  />
                </div>
              </div>
            )}

            <DialogFooter>
              <Button variant="outline" onClick={() => setShowRestockDialog(false)}>
                Cancel
              </Button>
              <Button 
                onClick={handleRestock}
                disabled={!restockQuantity || parseInt(restockQuantity) <= 0}
              >
                <Plus className="h-4 w-4 mr-2" />
                Add Stock
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Settings Dialog */}
        <Dialog open={showSettingsDialog} onOpenChange={setShowSettingsDialog}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Alert Settings</DialogTitle>
              <DialogDescription>
                Configure your stock alert preferences and thresholds
              </DialogDescription>
            </DialogHeader>
            
            <Tabs defaultValue="alerts" className="w-full">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="alerts">Alert Types</TabsTrigger>
                <TabsTrigger value="notifications">Notifications</TabsTrigger>
              </TabsList>
              
              <TabsContent value="alerts" className="space-y-4">
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <Label className="text-sm font-medium">Low Stock Alerts</Label>
                      <p className="text-xs text-muted-foreground">
                        Alert when products fall below minimum stock level
                      </p>
                    </div>
                    <Switch 
                      checked={alertSettings.lowStockEnabled}
                      onCheckedChange={(checked) => 
                        setAlertSettings({...alertSettings, lowStockEnabled: checked})
                      }
                    />
                  </div>

                  <div className="flex items-center justify-between">
                    <div>
                      <Label className="text-sm font-medium">Out of Stock Alerts</Label>
                      <p className="text-xs text-muted-foreground">
                        Alert when products are completely out of stock
                      </p>
                    </div>
                    <Switch 
                      checked={alertSettings.outOfStockEnabled}
                      onCheckedChange={(checked) => 
                        setAlertSettings({...alertSettings, outOfStockEnabled: checked})
                      }
                    />
                  </div>

                  <div className="flex items-center justify-between">
                    <div>
                      <Label className="text-sm font-medium">Overstock Alerts</Label>
                      <p className="text-xs text-muted-foreground">
                        Alert when products exceed maximum stock level
                      </p>
                    </div>
                    <Switch 
                      checked={alertSettings.overStockEnabled}
                      onCheckedChange={(checked) => 
                        setAlertSettings({...alertSettings, overStockEnabled: checked})
                      }
                    />
                  </div>

                  <div className="flex items-center justify-between">
                    <div>
                      <Label className="text-sm font-medium">Expiry Alerts</Label>
                      <p className="text-xs text-muted-foreground">
                        Alert for products nearing expiration date
                      </p>
                    </div>
                    <Switch 
                      checked={alertSettings.expiryEnabled}
                      onCheckedChange={(checked) => 
                        setAlertSettings({...alertSettings, expiryEnabled: checked})
                      }
                    />
                  </div>

                  <Separator />

                  <div>
                    <Label htmlFor="threshold" className="text-sm font-medium">
                      Default Low Stock Threshold
                    </Label>
                    <Input
                      id="threshold"
                      type="number"
                      value={alertSettings.defaultLowStockThreshold}
                      onChange={(e) => 
                        setAlertSettings({
                          ...alertSettings, 
                          defaultLowStockThreshold: parseInt(e.target.value)
                        })
                      }
                      className="mt-2 w-32"
                      min="1"
                    />
                    <p className="text-xs text-muted-foreground mt-1">
                      Applied to new products without specific thresholds
                    </p>
                  </div>
                </div>
              </TabsContent>
              
              <TabsContent value="notifications" className="space-y-4">
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <Label className="text-sm font-medium">Email Notifications</Label>
                      <p className="text-xs text-muted-foreground">
                        Receive alerts via email
                      </p>
                    </div>
                    <Switch 
                      checked={alertSettings.emailNotifications}
                      onCheckedChange={(checked) => 
                        setAlertSettings({...alertSettings, emailNotifications: checked})
                      }
                    />
                  </div>

                  <div className="flex items-center justify-between">
                    <div>
                      <Label className="text-sm font-medium">SMS Notifications</Label>
                      <p className="text-xs text-muted-foreground">
                        Receive alerts via SMS
                      </p>
                    </div>
                    <Switch 
                      checked={alertSettings.smsNotifications}
                      onCheckedChange={(checked) => 
                        setAlertSettings({...alertSettings, smsNotifications: checked})
                      }
                    />
                  </div>

                  {alertSettings.emailNotifications && (
                    <div>
                      <Label htmlFor="email" className="text-sm font-medium">
                        Notification Email
                      </Label>
                      <Input
                        id="email"
                        type="email"
                        placeholder="Enter email address"
                        className="mt-2"
                      />
                    </div>
                  )}

                  {alertSettings.smsNotifications && (
                    <div>
                      <Label htmlFor="phone" className="text-sm font-medium">
                        Notification Phone
                      </Label>
                      <Input
                        id="phone"
                        type="tel"
                        placeholder="Enter phone number"
                        className="mt-2"
                      />
                    </div>
                  )}
                </div>
              </TabsContent>
            </Tabs>

            <DialogFooter>
              <Button variant="outline" onClick={() => setShowSettingsDialog(false)}>
                Cancel
              </Button>
              <Button onClick={() => {
                console.log('Saving alert settings:', alertSettings);
                setShowSettingsDialog(false);
              }}>
                <Save className="h-4 w-4 mr-2" />
                Save Settings
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
};

export default StockAlertsPage;