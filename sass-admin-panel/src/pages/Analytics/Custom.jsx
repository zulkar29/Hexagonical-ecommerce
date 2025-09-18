import React, { useState } from 'react';
import {
  BarChart3,
  LineChart,
  PieChart,
  TrendingUp,
  Calendar,
  Download,
  Filter,
  Plus,
  Settings,
  Eye,
  Trash2,
  Edit,
  Share,
  Clock,
  Users,
  DollarSign,
  Activity,
  Target,
  Zap
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

const reportTemplates = [
  {
    id: 'revenue-analysis',
    name: 'Revenue Analysis',
    description: 'Comprehensive revenue breakdown by tenant, plan, and time period',
    icon: DollarSign,
    category: 'Financial',
    metrics: ['Total Revenue', 'MRR', 'ARR', 'Churn Rate'],
    lastUsed: '2024-01-15'
  },
  {
    id: 'user-engagement',
    name: 'User Engagement',
    description: 'Track user activity, feature adoption, and engagement metrics',
    icon: Users,
    category: 'User Analytics',
    metrics: ['Active Users', 'Session Duration', 'Feature Usage', 'Retention'],
    lastUsed: '2024-01-14'
  },
  {
    id: 'performance-metrics',
    name: 'Performance Metrics',
    description: 'System performance, uptime, and technical KPIs',
    icon: Activity,
    category: 'Technical',
    metrics: ['Uptime', 'Response Time', 'Error Rate', 'API Usage'],
    lastUsed: '2024-01-13'
  },
  {
    id: 'conversion-funnel',
    name: 'Conversion Funnel',
    description: 'Track user journey from signup to paid subscription',
    icon: Target,
    category: 'Marketing',
    metrics: ['Signup Rate', 'Trial Conversion', 'Upgrade Rate', 'Churn'],
    lastUsed: '2024-01-12'
  }
];

const savedReports = [
  {
    id: 'RPT-001',
    name: 'Q1 2024 Revenue Report',
    template: 'Revenue Analysis',
    createdBy: 'Admin User',
    createdAt: '2024-01-15T10:30:00Z',
    lastRun: '2024-01-15T14:20:00Z',
    status: 'Ready',
    schedule: 'Weekly',
    recipients: ['admin@company.com', 'finance@company.com']
  },
  {
    id: 'RPT-002',
    name: 'User Engagement Dashboard',
    template: 'User Engagement',
    createdBy: 'Product Manager',
    createdAt: '2024-01-14T09:15:00Z',
    lastRun: '2024-01-15T08:00:00Z',
    status: 'Running',
    schedule: 'Daily',
    recipients: ['product@company.com']
  },
  {
    id: 'RPT-003',
    name: 'System Health Report',
    template: 'Performance Metrics',
    createdBy: 'DevOps Team',
    createdAt: '2024-01-13T16:45:00Z',
    lastRun: '2024-01-15T12:00:00Z',
    status: 'Completed',
    schedule: 'Hourly',
    recipients: ['devops@company.com', 'support@company.com']
  }
];

const chartTypes = [
  { id: 'line', name: 'Line Chart', icon: LineChart, description: 'Show trends over time' },
  { id: 'bar', name: 'Bar Chart', icon: BarChart3, description: 'Compare values across categories' },
  { id: 'pie', name: 'Pie Chart', icon: PieChart, description: 'Show proportions and percentages' },
  { id: 'area', name: 'Area Chart', icon: TrendingUp, description: 'Visualize cumulative data' }
];

export default function CustomReports() {
  const [activeTab, setActiveTab] = useState('templates');
  const [searchTerm, setSearchTerm] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('all');
  const [statusFilter, setStatusFilter] = useState('all');

  const filteredTemplates = reportTemplates.filter(template => {
    const matchesSearch = template.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         template.description.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesCategory = categoryFilter === 'all' || template.category === categoryFilter;
    return matchesSearch && matchesCategory;
  });

  const filteredReports = savedReports.filter(report => {
    const matchesSearch = report.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         report.template.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || report.status.toLowerCase() === statusFilter.toLowerCase();
    return matchesSearch && matchesStatus;
  });

  const getStatusBadge = (status) => {
    switch (status) {
      case 'Ready':
        return <Badge variant="default" className="bg-green-600">Ready</Badge>;
      case 'Running':
        return <Badge variant="outline" className="border-blue-600 text-blue-600">Running</Badge>;
      case 'Completed':
        return <Badge variant="secondary">Completed</Badge>;
      case 'Failed':
        return <Badge variant="destructive">Failed</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  const formatDateTime = (dateString) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="space-y-6 p-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-3xl font-bold tracking-tight">Custom Reports</h2>
              <p className="text-muted-foreground">
                Create, schedule, and manage custom analytics reports
              </p>
            </div>
            <div className="flex items-center space-x-3">
              <Button variant="outline">
                <Download className="h-4 w-4 mr-2" />
                Export All
              </Button>
              <Button>
                <Plus className="h-4 w-4 mr-2" />
                Create Report
              </Button>
            </div>
          </div>

          {/* Main Content Tabs */}
          <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
            <TabsList className="grid w-full grid-cols-4">
              <TabsTrigger value="templates">Templates</TabsTrigger>
              <TabsTrigger value="saved">Saved Reports</TabsTrigger>
              <TabsTrigger value="builder">Report Builder</TabsTrigger>
              <TabsTrigger value="scheduled">Scheduled</TabsTrigger>
            </TabsList>

            <TabsContent value="templates" className="space-y-6">
              {/* Search and Filters */}
              <Card>
                <CardContent className="p-6">
                  <div className="flex flex-col md:flex-row gap-4">
                    <div className="flex-1">
                      <Input
                        type="text"
                        placeholder="Search templates..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                      />
                    </div>
                    <div className="flex gap-2">
                      <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                        <SelectTrigger className="w-40">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="all">All Categories</SelectItem>
                          <SelectItem value="Financial">Financial</SelectItem>
                          <SelectItem value="User Analytics">User Analytics</SelectItem>
                          <SelectItem value="Technical">Technical</SelectItem>
                          <SelectItem value="Marketing">Marketing</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Templates Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredTemplates.map((template) => (
                  <Card key={template.id} className="hover:shadow-lg transition-shadow">
                    <CardContent className="p-6">
                      <div className="flex items-start justify-between mb-4">
                        <div className="flex items-center gap-3">
                          <template.icon className="h-8 w-8 text-primary" />
                          <div>
                            <h3 className="font-semibold">{template.name}</h3>
                            <Badge variant="outline" className="text-xs">
                              {template.category}
                            </Badge>
                          </div>
                        </div>
                      </div>
                      <p className="text-sm text-muted-foreground mb-4">
                        {template.description}
                      </p>
                      <div className="space-y-3">
                        <div>
                          <p className="text-xs font-medium text-muted-foreground mb-2">Key Metrics:</p>
                          <div className="flex flex-wrap gap-1">
                            {template.metrics.map((metric) => (
                              <Badge key={metric} variant="secondary" className="text-xs">
                                {metric}
                              </Badge>
                            ))}
                          </div>
                        </div>
                        <div className="flex items-center justify-between pt-2">
                          <p className="text-xs text-muted-foreground">
                            Last used: {new Date(template.lastUsed).toLocaleDateString()}
                          </p>
                          <Button size="sm">
                            Use Template
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </TabsContent>

            <TabsContent value="saved" className="space-y-6">
              {/* Search and Filters */}
              <Card>
                <CardContent className="p-6">
                  <div className="flex flex-col md:flex-row gap-4">
                    <div className="flex-1">
                      <Input
                        type="text"
                        placeholder="Search saved reports..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                      />
                    </div>
                    <div className="flex gap-2">
                      <Select value={statusFilter} onValueChange={setStatusFilter}>
                        <SelectTrigger className="w-40">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="all">All Status</SelectItem>
                          <SelectItem value="ready">Ready</SelectItem>
                          <SelectItem value="running">Running</SelectItem>
                          <SelectItem value="completed">Completed</SelectItem>
                          <SelectItem value="failed">Failed</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Saved Reports Table */}
              <Card>
                <CardHeader>
                  <CardTitle>Saved Reports</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead>
                        <tr className="border-b">
                          <th className="text-left p-4 font-medium">Report Name</th>
                          <th className="text-left p-4 font-medium">Template</th>
                          <th className="text-left p-4 font-medium">Status</th>
                          <th className="text-left p-4 font-medium">Schedule</th>
                          <th className="text-left p-4 font-medium">Last Run</th>
                          <th className="text-left p-4 font-medium">Created By</th>
                          <th className="text-left p-4 font-medium">Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        {filteredReports.map((report) => (
                          <tr key={report.id} className="border-b hover:bg-muted/50">
                            <td className="p-4">
                              <div>
                                <p className="font-medium">{report.name}</p>
                                <p className="text-sm text-muted-foreground">{report.id}</p>
                              </div>
                            </td>
                            <td className="p-4">
                              <Badge variant="outline">{report.template}</Badge>
                            </td>
                            <td className="p-4">{getStatusBadge(report.status)}</td>
                            <td className="p-4">
                              <Badge variant="secondary">{report.schedule}</Badge>
                            </td>
                            <td className="p-4 text-sm">{formatDateTime(report.lastRun)}</td>
                            <td className="p-4 text-sm">{report.createdBy}</td>
                            <td className="p-4">
                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button variant="ghost" size="sm">
                                    <Settings className="h-4 w-4" />
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                  <DropdownMenuItem>
                                    <Eye className="h-4 w-4 mr-2" />
                                    View Report
                                  </DropdownMenuItem>
                                  <DropdownMenuItem>
                                    <Edit className="h-4 w-4 mr-2" />
                                    Edit
                                  </DropdownMenuItem>
                                  <DropdownMenuItem>
                                    <Share className="h-4 w-4 mr-2" />
                                    Share
                                  </DropdownMenuItem>
                                  <DropdownMenuItem>
                                    <Download className="h-4 w-4 mr-2" />
                                    Export
                                  </DropdownMenuItem>
                                  <DropdownMenuItem className="text-red-600">
                                    <Trash2 className="h-4 w-4 mr-2" />
                                    Delete
                                  </DropdownMenuItem>
                                </DropdownMenuContent>
                              </DropdownMenu>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="builder" className="space-y-6">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Chart Types */}
                <Card>
                  <CardHeader>
                    <CardTitle>Chart Types</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-2 gap-4">
                      {chartTypes.map((chart) => (
                        <div key={chart.id} className="p-4 border rounded-lg hover:bg-muted/50 cursor-pointer">
                          <div className="flex items-center gap-3 mb-2">
                            <chart.icon className="h-6 w-6 text-primary" />
                            <p className="font-medium">{chart.name}</p>
                          </div>
                          <p className="text-sm text-muted-foreground">{chart.description}</p>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>

                {/* Data Sources */}
                <Card>
                  <CardHeader>
                    <CardTitle>Data Sources</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      <div className="p-3 border rounded-lg">
                        <div className="flex items-center gap-2 mb-1">
                          <Users className="h-4 w-4 text-blue-600" />
                          <p className="font-medium">User Data</p>
                        </div>
                        <p className="text-sm text-muted-foreground">User accounts, activity, and engagement metrics</p>
                      </div>
                      <div className="p-3 border rounded-lg">
                        <div className="flex items-center gap-2 mb-1">
                          <DollarSign className="h-4 w-4 text-green-600" />
                          <p className="font-medium">Revenue Data</p>
                        </div>
                        <p className="text-sm text-muted-foreground">Subscriptions, payments, and financial metrics</p>
                      </div>
                      <div className="p-3 border rounded-lg">
                        <div className="flex items-center gap-2 mb-1">
                          <Activity className="h-4 w-4 text-orange-600" />
                          <p className="font-medium">System Data</p>
                        </div>
                        <p className="text-sm text-muted-foreground">Performance, uptime, and technical metrics</p>
                      </div>
                      <div className="p-3 border rounded-lg">
                        <div className="flex items-center gap-2 mb-1">
                          <Target className="h-4 w-4 text-purple-600" />
                          <p className="font-medium">Marketing Data</p>
                        </div>
                        <p className="text-sm text-muted-foreground">Campaigns, conversions, and acquisition metrics</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Report Builder Actions */}
              <Card>
                <CardContent className="p-6">
                  <div className="text-center space-y-4">
                    <Zap className="h-12 w-12 text-primary mx-auto" />
                    <div>
                      <h3 className="text-lg font-semibold">Build Custom Report</h3>
                      <p className="text-muted-foreground">Drag and drop to create your custom analytics report</p>
                    </div>
                    <Button size="lg">
                      <Plus className="h-4 w-4 mr-2" />
                      Start Building
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="scheduled" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Scheduled Reports</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {savedReports.filter(report => report.schedule !== 'Manual').map((report) => (
                      <div key={report.id} className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center gap-4">
                          <Clock className="h-6 w-6 text-muted-foreground" />
                          <div>
                            <p className="font-medium">{report.name}</p>
                            <p className="text-sm text-muted-foreground">
                              Runs {report.schedule.toLowerCase()} • Next: Tomorrow 9:00 AM
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          {getStatusBadge(report.status)}
                          <Button variant="outline" size="sm">
                            <Settings className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
    </div>
  );
}