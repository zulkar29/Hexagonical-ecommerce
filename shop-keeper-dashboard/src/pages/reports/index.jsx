import { useState, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { 
  Plus, 
  FileText, 
  Download, 
  Search, 
  MoreHorizontal,
  Calendar,
  BarChart3,
  TrendingUp,
  Users,
  Package,
  ShoppingCart,
  Eye,
  Trash2
} from "lucide-react";
import { toast } from "sonner";

// Enhanced mock data
const generateReportsData = () => [
  { 
    id: 1, 
    name: "Sales Performance Report", 
    description: "Monthly sales analysis and performance metrics",
    category: "Sales",
    period: "May 2025", 
    type: "PDF", 
    size: "1.2MB", 
    date: "2025-06-01",
    status: "completed",
    generated_by: "Admin User",
    downloads: 45
  },
  { 
    id: 2, 
    name: "Order Analytics Report", 
    description: "Detailed order statistics and trends",
    category: "Orders",
    period: "May 2025", 
    type: "CSV", 
    size: "800KB", 
    date: "2025-06-01",
    status: "completed",
    generated_by: "Sales Manager",
    downloads: 23
  },
  { 
    id: 3, 
    name: "Customer Insights Report", 
    description: "Customer behavior and demographics analysis",
    category: "Customers",
    period: "May 2025", 
    type: "PDF", 
    size: "950KB", 
    date: "2025-06-01",
    status: "completed",
    generated_by: "Marketing Team",
    downloads: 67
  },
  { 
    id: 4, 
    name: "Inventory Status Report", 
    description: "Stock levels and inventory management overview",
    category: "Inventory",
    period: "May 2025", 
    type: "CSV", 
    size: "600KB", 
    date: "2025-06-01",
    status: "generating",
    generated_by: "Inventory Manager",
    downloads: 0
  },
  { 
    id: 5, 
    name: "Financial Summary", 
    description: "Revenue, expenses and profit analysis",
    category: "Finance",
    period: "Q2 2025", 
    type: "PDF", 
    size: "2.1MB", 
    date: "2025-06-15",
    status: "completed",
    generated_by: "Finance Team",
    downloads: 89
  },
  { 
    id: 6, 
    name: "Marketing Campaign Analysis", 
    description: "Campaign performance and ROI metrics",
    category: "Marketing",
    period: "May 2025", 
    type: "XLSX", 
    size: "1.5MB", 
    date: "2025-06-03",
    status: "failed",
    generated_by: "Marketing Manager",
    downloads: 0
  }
];

const ITEMS_PER_PAGE = 10;

export default function Reports() {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [currentPage, setCurrentPage] = useState(1);

  const allReports = generateReportsData();

  // Filter reports
  const filteredReports = useMemo(() => {
    let filtered = allReports;

    // Apply search filter
    if (searchQuery?.trim()) {
      const query = searchQuery.trim().toLowerCase();
      filtered = filtered.filter(report => 
        report.name.toLowerCase().includes(query) ||
        report.description.toLowerCase().includes(query) ||
        report.category.toLowerCase().includes(query)
      );
    }

    // Apply category filter
    if (categoryFilter !== 'all') {
      filtered = filtered.filter(report => 
        report.category.toLowerCase() === categoryFilter.toLowerCase()
      );
    }

    // Apply status filter
    if (statusFilter !== 'all') {
      filtered = filtered.filter(report => report.status === statusFilter);
    }

    return filtered;
  }, [searchQuery, categoryFilter, statusFilter, allReports]);

  // Pagination
  const totalPages = Math.ceil(filteredReports.length / ITEMS_PER_PAGE);
  const paginatedReports = filteredReports.slice(
    (currentPage - 1) * ITEMS_PER_PAGE,
    currentPage * ITEMS_PER_PAGE
  );

  const handleDownload = (report) => {
    if (report.status === 'completed') {
      toast.success(`Downloading ${report.name}...`);
      // Simulate download
    } else {
      toast.error("Report is not ready for download");
    }
  };

  const handleDelete = (report) => {
    toast.success(`Report "${report.name}" deleted successfully`);
    // Handle delete logic here
  };

  const getStatusBadge = (status) => {
    const variants = {
      completed: "default",
      generating: "secondary", 
      failed: "destructive"
    };
    
    const labels = {
      completed: "Completed",
      generating: "Generating",
      failed: "Failed"
    };

    return (
      <Badge variant={variants[status] || "secondary"}>
        {labels[status] || status}
      </Badge>
    );
  };

  const getCategoryIcon = (category) => {
    const icons = {
      Sales: BarChart3,
      Orders: ShoppingCart,
      Customers: Users,
      Inventory: Package,
      Finance: TrendingUp,
      Marketing: BarChart3
    };
    
    const IconComponent = icons[category] || FileText;
    return <IconComponent className="h-4 w-4" />;
  };

  const getTypeColor = (type) => {
    const colors = {
      PDF: "text-red-600 bg-red-50",
      CSV: "text-green-600 bg-green-50", 
      XLSX: "text-blue-600 bg-blue-50"
    };
    return colors[type] || "text-gray-600 bg-gray-50";
  };

  return (
    <div className="space-y-6 p-6">
      <Card className="w-full">
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <FileText className="h-5 w-5" />
              Reports ({allReports.length})
            </CardTitle>
            <CardDescription>
              Generate, manage and download business reports
            </CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm">
              <Calendar className="mr-2 h-4 w-4" />
              Schedule Report
            </Button>
            <Button onClick={() => navigate("/reports/create")}>
              <Plus className="mr-2 h-4 w-4" />
              Create Report
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {/* Search and Filter section */}
          <div className="flex flex-col gap-4 md:flex-row md:items-center mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search reports by name, description, or category..."
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value);
                  setCurrentPage(1);
                }}
                className="pl-10"
              />
            </div>
            <div className="flex gap-2">
              <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                <SelectTrigger className="w-40">
                  <SelectValue placeholder="Category" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Categories</SelectItem>
                  <SelectItem value="sales">Sales</SelectItem>
                  <SelectItem value="orders">Orders</SelectItem>
                  <SelectItem value="customers">Customers</SelectItem>
                  <SelectItem value="inventory">Inventory</SelectItem>
                  <SelectItem value="finance">Finance</SelectItem>
                  <SelectItem value="marketing">Marketing</SelectItem>
                </SelectContent>
              </Select>
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                  <SelectItem value="generating">Generating</SelectItem>
                  <SelectItem value="failed">Failed</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Reports Table */}
          <div className="w-full overflow-x-auto">
            <div className="rounded-md border">
              <Table className="w-full min-w-full">
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[30%]">Report</TableHead>
                    <TableHead className="w-[15%]">Category</TableHead>
                    <TableHead className="w-[12%]">Period</TableHead>
                    <TableHead className="w-[10%]">Type</TableHead>
                    <TableHead className="w-[10%]">Status</TableHead>
                    <TableHead className="w-[13%]">Generated</TableHead>
                    <TableHead className="w-[10%] text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {paginatedReports.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                        {searchQuery?.trim() ? `No reports found matching "${searchQuery.trim()}"` : "No reports found"}
                      </TableCell>
                    </TableRow>
                  ) : (
                    paginatedReports.map((report) => (
                      <TableRow key={report.id}>
                        <TableCell>
                          <div className="flex items-start gap-3">
                            <div className="flex-shrink-0 mt-1">
                              {getCategoryIcon(report.category)}
                            </div>
                            <div className="min-w-0">
                              <div className="font-medium text-sm">{report.name}</div>
                              <div className="text-xs text-muted-foreground line-clamp-1" title={report.description}>
                                {report.description}
                              </div>
                              <div className="text-xs text-muted-foreground mt-1">
                                By {report.generated_by}
                              </div>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline" className="text-xs">
                            {report.category}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground">
                          {report.period}
                        </TableCell>
                        <TableCell>
                          <Badge 
                            variant="secondary" 
                            className={`text-xs ${getTypeColor(report.type)}`}
                          >
                            {report.type}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          {getStatusBadge(report.status)}
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground">
                          <div>{new Date(report.date).toLocaleDateString()}</div>
                          {report.downloads > 0 && (
                            <div className="text-xs text-muted-foreground">
                              {report.downloads} downloads
                            </div>
                          )}
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex items-center justify-end gap-2">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleDownload(report)}
                              disabled={report.status !== 'completed'}
                              title="Download Report"
                            >
                              <Download className="h-4 w-4" />
                            </Button>
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="icon">
                                  <MoreHorizontal className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end" className="w-48">
                                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem>
                                  <Eye className="mr-2 h-4 w-4" />
                                  View Details
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                  onClick={() => handleDownload(report)}
                                  disabled={report.status !== 'completed'}
                                >
                                  <Download className="mr-2 h-4 w-4" />
                                  Download
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem
                                  onClick={() => handleDelete(report)}
                                  className="text-destructive"
                                >
                                  <Trash2 className="mr-2 h-4 w-4" />
                                  Delete
                                </DropdownMenuItem>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between mt-6">
              <div className="text-sm text-muted-foreground">
                Showing {((currentPage - 1) * ITEMS_PER_PAGE) + 1} to {Math.min(currentPage * ITEMS_PER_PAGE, filteredReports.length)} of {filteredReports.length} reports
              </div>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                >
                  Previous
                </Button>
                <span className="text-sm text-muted-foreground">
                  Page {currentPage} of {totalPages}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                  disabled={currentPage === totalPages}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}