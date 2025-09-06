import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Search, Plus, ArrowUpDown, Edit, Trash2, Loader2 } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { DeleteConfirmationDialog } from "@/components/ui/delete-confirmation-dialog";

const ITEMS_PER_PAGE = 5;

const mockDiscounts = [
  { id: 1, code: "WELCOME10", type: "Percentage", value: "10%", usage: 120, status: "Active", expires: "2025-06-01" },
  { id: 2, code: "FREESHIP", type: "Free Shipping", value: "-", usage: 80, status: "Active", expires: "2025-05-20" },
  { id: 3, code: "SUMMER25", type: "Amount", value: "$25", usage: 30, status: "Expired", expires: "2025-05-01" },
  { id: 4, code: "VIP50", type: "Percentage", value: "50%", usage: 5, status: "Inactive", expires: "2025-12-31" },
];

const Discounts = () => {
  const navigate = useNavigate();
  const [search, setSearch] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [sortConfig, setSortConfig] = useState({ key: null, direction: null });
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [discountToDelete, setDiscountToDelete] = useState(null);

  const { data: discounts = [], isLoading, isError } = useQuery({
    queryKey: ["discounts"],
    queryFn: async () => mockDiscounts,
  });

  const handleSort = (key) => {
    const direction = sortConfig.key === key && sortConfig.direction === "asc" ? "desc" : "asc";
    setSortConfig({ key, direction });
  };

  const handleDelete = (discount) => {
    setDiscountToDelete(discount);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (discountToDelete) {
      console.log("Deleting discount:", discountToDelete.code);
      // Add your delete logic here
      setDeleteDialogOpen(false);
    }
  };

  const filteredDiscounts = discounts
    .filter((discount) =>
      discount.code.toLowerCase().includes(search.toLowerCase()) ||
      discount.type.toLowerCase().includes(search.toLowerCase())
    )
    .sort((a, b) => {
      if (!sortConfig.key) return 0;
      const direction = sortConfig.direction === "asc" ? 1 : -1;
      return a[sortConfig.key] > b[sortConfig.key] ? direction : -direction;
    });

  const totalPages = Math.ceil(filteredDiscounts.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const paginatedDiscounts = filteredDiscounts.slice(startIndex, startIndex + ITEMS_PER_PAGE);

  if (isLoading) return (
    <div className="flex items-center justify-center h-[400px]">
      <Loader2 className="h-8 w-8 animate-spin" />
    </div>
  );

  if (isError) return (
    <Card>
      <CardContent className="flex items-center justify-center h-[400px] text-destructive">
        Error loading discounts. Please try again later.
      </CardContent>
    </Card>
  );

  return (
    <div className="space-y-6 p-6">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Discounts & Promo Codes</CardTitle>
            <CardDescription>
              Create and manage discount codes for your store
            </CardDescription>
          </div>
          <Button onClick={() => navigate("/discounts/create")}>
            <Plus className="mr-2 h-4 w-4" /> New Discount
          </Button>
        </CardHeader>
        <CardContent>
          {/* Search section */}
          <div className="flex flex-col gap-4 md:flex-row md:items-center mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search discounts..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>

          {/* Discounts Table */}
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>
                    <Button variant="ghost" onClick={() => handleSort("code")}>
                      Code
                      <ArrowUpDown className="ml-2 h-4 w-4" />
                    </Button>
                  </TableHead>
                  <TableHead>
                    <Button variant="ghost" onClick={() => handleSort("type")}>
                      Type
                      <ArrowUpDown className="ml-2 h-4 w-4" />
                    </Button>
                  </TableHead>
                  <TableHead>
                    <Button variant="ghost" onClick={() => handleSort("value")}>
                      Value
                      <ArrowUpDown className="ml-2 h-4 w-4" />
                    </Button>
                  </TableHead>
                  <TableHead>
                    <Button variant="ghost" onClick={() => handleSort("usage")}>
                      Usage
                      <ArrowUpDown className="ml-2 h-4 w-4" />
                    </Button>
                  </TableHead>
                  <TableHead>
                    <Button variant="ghost" onClick={() => handleSort("expires")}>
                      Expires
                      <ArrowUpDown className="ml-2 h-4 w-4" />
                    </Button>
                  </TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {paginatedDiscounts.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-muted-foreground">
                      No discounts found
                    </TableCell>
                  </TableRow>
                ) : (
                  paginatedDiscounts.map((discount) => (
                    <TableRow key={discount.id}>
                      <TableCell className="font-medium tracking-wider">{discount.code}</TableCell>
                      <TableCell>{discount.type}</TableCell>
                      <TableCell>{discount.value}</TableCell>
                      <TableCell>{discount.usage}</TableCell>
                      <TableCell>{discount.expires}</TableCell>
                      <TableCell>
                        <Badge variant={
                          discount.status === "Active" ? "default" :
                          discount.status === "Inactive" ? "secondary" :
                          "destructive"
                        }>
                          {discount.status}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => navigate(`/discounts/${discount.id}/edit`)}
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(discount)}
                            className="text-destructive hover:text-destructive/90"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex justify-center mt-6">
              <Pagination>
                <PaginationContent className="flex items-center gap-2">
                  <PaginationItem>
                    <PaginationPrevious
                      onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                      disabled={currentPage === 1}
                    />
                  </PaginationItem>

                  <PaginationItem>
                    <PaginationLink disabled className="pointer-events-none">
                      Page {currentPage} of {totalPages}
                    </PaginationLink>
                  </PaginationItem>

                  <PaginationItem>
                    <PaginationNext
                      onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                      disabled={currentPage === totalPages}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <DeleteConfirmationDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onConfirm={confirmDelete}
        itemName={discountToDelete?.code}
      />
    </div>
  );
};

export default Discounts;