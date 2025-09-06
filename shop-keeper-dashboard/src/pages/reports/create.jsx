import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { 
  ArrowLeft,
  Save,
  FileText,
  Calendar,
  Settings,
  BarChart3,
  TrendingUp,
  Users,
  Package,
  ShoppingCart,
  DollarSign,
  ChevronRight
} from "lucide-react";
import { toast } from "sonner";

const reportTemplates = [
  {
    id: "sales_performance",
    name: "Sales Performance",
    description: "Revenue, conversion rates, and sales trends analysis",
    category: "Sales",
    icon: BarChart3,
    metrics: ["Revenue", "Orders", "Conversion Rate", "Average Order Value"],
    estimatedTime: "2-3 minutes"
  },
  {
    id: "order_analytics",
    name: "Order Analytics", 
    description: "Order statistics, fulfillment rates, and delivery metrics",
    category: "Orders",
    icon: ShoppingCart,
    metrics: ["Total Orders", "Order Status", "Fulfillment Rate", "Delivery Time"],
    estimatedTime: "1-2 minutes"
  },
  {
    id: "customer_insights",
    name: "Customer Insights",
    description: "Customer behavior, demographics, and retention analysis",
    category: "Customers", 
    icon: Users,
    metrics: ["Customer Count", "Retention Rate", "Demographics", "Lifetime Value"],
    estimatedTime: "3-4 minutes"
  },
  {
    id: "inventory_status",
    name: "Inventory Status",
    description: "Stock levels, inventory turnover, and supply chain metrics",
    category: "Inventory",
    icon: Package,
    metrics: ["Stock Levels", "Turnover Rate", "Reorder Points", "Supplier Performance"],
    estimatedTime: "2-3 minutes"
  },
  {
    id: "financial_summary",
    name: "Financial Summary",
    description: "Revenue, expenses, profit margins, and financial health",
    category: "Finance",
    icon: DollarSign,
    metrics: ["Revenue", "Expenses", "Profit Margin", "Cash Flow"],
    estimatedTime: "4-5 minutes"
  },
  {
    id: "marketing_performance",
    name: "Marketing Performance",
    description: "Campaign metrics, ROI analysis, and marketing effectiveness",
    category: "Marketing",
    icon: TrendingUp,
    metrics: ["Campaign ROI", "Lead Generation", "Conversion Funnel", "Channel Performance"],
    estimatedTime: "3-4 minutes"
  }
];

export default function CreateReport() {
  const navigate = useNavigate();
  const [selectedTemplate, setSelectedTemplate] = useState("");
  const [form, setForm] = useState({
    name: "",
    description: "",
    category: "",
    format: "PDF",
    period: "last_30_days",
    scheduled: false,
    schedule_frequency: "weekly",
    include_charts: true,
    include_summary: true,
    email_recipients: ""
  });
  const [submitting, setSubmitting] = useState(false);

  const handleTemplateSelect = (template) => {
    setSelectedTemplate(template.id);
    setForm(prev => ({
      ...prev,
      name: template.name + " Report",
      description: template.description,
      category: template.category
    }));
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));
  };

  const handleSwitchChange = (name, checked) => {
    setForm(prev => ({ ...prev, [name]: checked }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      console.log("Creating report:", { ...form, template: selectedTemplate });
      toast.success("Report generation started! You'll be notified when it's ready.");
      navigate("/reports");
    } catch (error) {
      toast.error("Failed to create report. Please try again.");
    } finally {
      setSubmitting(false);
    }
  };

  const selectedTemplateData = reportTemplates.find(t => t.id === selectedTemplate);

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/reports")}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold tracking-tight">Create New Report</h1>
          <p className="text-muted-foreground">
            Generate custom business reports and analytics
          </p>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Template Selection */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle>Choose Report Template</CardTitle>
              <CardDescription>
                Select a pre-built template or create a custom report
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-2">
                {reportTemplates.map((template) => {
                  const IconComponent = template.icon;
                  const isSelected = selectedTemplate === template.id;
                  
                  return (
                    <div
                      key={template.id}
                      className={`relative p-4 border rounded-lg cursor-pointer transition-all hover:shadow-md ${
                        isSelected ? "border-primary bg-primary/5" : "border-border hover:border-primary/50"
                      }`}
                      onClick={() => handleTemplateSelect(template)}
                    >
                      <div className="flex items-start gap-3">
                        <div className={`p-2 rounded-lg ${isSelected ? "bg-primary text-primary-foreground" : "bg-muted"}`}>
                          <IconComponent className="h-5 w-5" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-1">
                            <h3 className="font-medium text-sm">{template.name}</h3>
                            <Badge variant="outline" className="text-xs">
                              {template.category}
                            </Badge>
                          </div>
                          <p className="text-xs text-muted-foreground mb-3 line-clamp-2">
                            {template.description}
                          </p>
                          <div className="space-y-2">
                            <div className="text-xs text-muted-foreground">
                              <strong>Key Metrics:</strong> {template.metrics.slice(0, 2).join(", ")}
                              {template.metrics.length > 2 && ` +${template.metrics.length - 2} more`}
                            </div>
                            <div className="text-xs text-muted-foreground">
                              <strong>Est. Time:</strong> {template.estimatedTime}
                            </div>
                          </div>
                        </div>
                      </div>
                      {isSelected && (
                        <div className="absolute top-2 right-2">
                          <div className="w-2 h-2 bg-primary rounded-full"></div>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          {/* Report Configuration */}
          {selectedTemplate && (
            <Card className="mt-6">
              <CardHeader>
                <CardTitle>Report Configuration</CardTitle>
                <CardDescription>
                  Customize your report settings and parameters
                </CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-6">
                  {/* Basic Information */}
                  <div className="space-y-4">
                    <h3 className="text-lg font-medium flex items-center gap-2">
                      <FileText className="h-5 w-5" />
                      Basic Information
                    </h3>
                    
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="name">Report Name *</Label>
                        <Input
                          id="name"
                          name="name"
                          value={form.name}
                          onChange={handleChange}
                          placeholder="Enter report name"
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="format">Output Format</Label>
                        <Select value={form.format} onValueChange={(value) => setForm(prev => ({ ...prev, format: value }))}>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="PDF">PDF Document</SelectItem>
                            <SelectItem value="CSV">CSV Spreadsheet</SelectItem>
                            <SelectItem value="XLSX">Excel Workbook</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="description">Description</Label>
                      <Textarea
                        id="description"
                        name="description"
                        value={form.description}
                        onChange={handleChange}
                        placeholder="Brief description of the report"
                        rows={3}
                      />
                    </div>
                  </div>

                  <Separator />

                  {/* Time Period */}
                  <div className="space-y-4">
                    <h3 className="text-lg font-medium flex items-center gap-2">
                      <Calendar className="h-5 w-5" />
                      Time Period
                    </h3>
                    
                    <div className="space-y-2">
                      <Label htmlFor="period">Data Period</Label>
                      <Select value={form.period} onValueChange={(value) => setForm(prev => ({ ...prev, period: value }))}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="last_7_days">Last 7 Days</SelectItem>
                          <SelectItem value="last_30_days">Last 30 Days</SelectItem>
                          <SelectItem value="last_90_days">Last 90 Days</SelectItem>
                          <SelectItem value="last_6_months">Last 6 Months</SelectItem>
                          <SelectItem value="last_year">Last Year</SelectItem>
                          <SelectItem value="current_month">Current Month</SelectItem>
                          <SelectItem value="current_quarter">Current Quarter</SelectItem>
                          <SelectItem value="current_year">Current Year</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>

                  <Separator />

                  {/* Report Options */}
                  <div className="space-y-4">
                    <h3 className="text-lg font-medium flex items-center gap-2">
                      <Settings className="h-5 w-5" />
                      Report Options
                    </h3>
                    
                    <div className="space-y-4">
                      <div className="flex items-center justify-between">
                        <div className="space-y-1">
                          <Label htmlFor="include_charts">Include Charts & Graphs</Label>
                          <p className="text-sm text-muted-foreground">
                            Add visual charts and graphs to the report
                          </p>
                        </div>
                        <Switch
                          id="include_charts"
                          checked={form.include_charts}
                          onCheckedChange={(checked) => handleSwitchChange("include_charts", checked)}
                        />
                      </div>

                      <div className="flex items-center justify-between">
                        <div className="space-y-1">
                          <Label htmlFor="include_summary">Executive Summary</Label>
                          <p className="text-sm text-muted-foreground">
                            Include a summary section with key insights
                          </p>
                        </div>
                        <Switch
                          id="include_summary"
                          checked={form.include_summary}
                          onCheckedChange={(checked) => handleSwitchChange("include_summary", checked)}
                        />
                      </div>

                      <div className="flex items-center justify-between">
                        <div className="space-y-1">
                          <Label htmlFor="scheduled">Schedule Report</Label>
                          <p className="text-sm text-muted-foreground">
                            Automatically generate this report on a schedule
                          </p>
                        </div>
                        <Switch
                          id="scheduled"
                          checked={form.scheduled}
                          onCheckedChange={(checked) => handleSwitchChange("scheduled", checked)}
                        />
                      </div>

                      {form.scheduled && (
                        <div className="ml-6 space-y-4">
                          <div className="space-y-2">
                            <Label htmlFor="schedule_frequency">Frequency</Label>
                            <Select value={form.schedule_frequency} onValueChange={(value) => setForm(prev => ({ ...prev, schedule_frequency: value }))}>
                              <SelectTrigger>
                                <SelectValue />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="daily">Daily</SelectItem>
                                <SelectItem value="weekly">Weekly</SelectItem>
                                <SelectItem value="monthly">Monthly</SelectItem>
                                <SelectItem value="quarterly">Quarterly</SelectItem>
                              </SelectContent>
                            </Select>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="email_recipients">Email Recipients</Label>
                            <Input
                              id="email_recipients"
                              name="email_recipients"
                              value={form.email_recipients}
                              onChange={handleChange}
                              placeholder="admin@company.com, manager@company.com"
                            />
                            <p className="text-xs text-muted-foreground">
                              Separate multiple email addresses with commas
                            </p>
                          </div>
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Submit Buttons */}
                  <div className="flex justify-end gap-3 pt-6">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => navigate("/reports")}
                    >
                      Cancel
                    </Button>
                    <Button type="submit" disabled={submitting}>
                      {submitting ? (
                        <>
                          <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                          Generating...
                        </>
                      ) : (
                        <>
                          <Save className="mr-2 h-4 w-4" />
                          Generate Report
                        </>
                      )}
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Selected Template Preview */}
          {selectedTemplateData && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Selected Template</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-primary text-primary-foreground rounded-lg">
                    <selectedTemplateData.icon className="h-5 w-5" />
                  </div>
                  <div>
                    <div className="font-medium">{selectedTemplateData.name}</div>
                    <div className="text-sm text-muted-foreground">{selectedTemplateData.category}</div>
                  </div>
                </div>
                
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Description</Label>
                  <p className="text-sm mt-1">{selectedTemplateData.description}</p>
                </div>
                
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Key Metrics</Label>
                  <div className="flex flex-wrap gap-1 mt-1">
                    {selectedTemplateData.metrics.map((metric, index) => (
                      <Badge key={index} variant="outline" className="text-xs">
                        {metric}
                      </Badge>
                    ))}
                  </div>
                </div>
                
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Estimated Generation Time</Label>
                  <p className="text-sm mt-1">{selectedTemplateData.estimatedTime}</p>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Help & Tips */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Tips</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="text-sm text-muted-foreground space-y-2">
                <p><strong>Report Name:</strong> Use descriptive names for easy identification</p>
                <p><strong>Scheduling:</strong> Scheduled reports are sent automatically via email</p>
                <p><strong>Format:</strong> PDF for presentations, CSV/Excel for data analysis</p>
                <p><strong>Performance:</strong> Larger date ranges may take longer to generate</p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}