import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Star, CheckCircle, XCircle, Trash2, Search, ArrowUpDown, Loader2 } from "lucide-react";
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
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

const ITEMS_PER_PAGE = 5;

const mockReviews = [
  { id: 1, product: "Smartphone X", customer: "John Doe", rating: 5, comment: "Excellent product!", date: "2025-05-01", status: "Approved" },
  { id: 2, product: "Running Shoes", customer: "Jane Smith", rating: 4, comment: "Very comfortable.", date: "2025-04-28", status: "Pending" },
  { id: 3, product: "Coffee Maker", customer: "Emily Davis", rating: 2, comment: "Stopped working after a week.", date: "2025-04-25", status: "Rejected" },
  { id: 4, product: "Laptop Pro", customer: "Michael Brown", rating: 5, comment: "Super fast and reliable.", date: "2025-04-20", status: "Approved" },
  { id: 5, product: "Wireless Headphones", customer: "Sarah Wilson", rating: 3, comment: "Sound is okay.", date: "2025-04-18", status: "Pending" },
];

export default function Reviews() {
  const [search, setSearch] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const { data: reviews = [], isLoading, isError } = useQuery({
    queryKey: ["reviews"],
    queryFn: async () => mockReviews,
  });

  const handleStatusUpdate = (id, newStatus) => {
    // In a real app, update status via API
    // For demo, update local state
    setReviews(reviews =>
      reviews.map(r =>
        r.id === id ? { ...r, status: newStatus } : r
      )
    );
  };

  const filtered = search
    ? reviews.filter(r =>
        r.product.toLowerCase().includes(search.toLowerCase()) ||
        r.customer.toLowerCase().includes(search.toLowerCase()) ||
        r.comment.toLowerCase().includes(search.toLowerCase())
      )
    : reviews;

  const totalPages = Math.ceil(filtered.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const paginatedReviews = filtered.slice(startIndex, startIndex + ITEMS_PER_PAGE);

  const handleApprove = (id) => {
    handleStatusUpdate(id, "Approved");
  };
  const handleReject = (id) => {
    handleStatusUpdate(id, "Rejected");
  };
  const handleDelete = (id) => {
    // Delete logic here
  };

  if (isLoading) return (
    <div className="flex items-center justify-center h-[400px]">
      <Loader2 className="h-8 w-8 animate-spin" />
    </div>
  );

  if (isError) return (
    <Card>
      <CardContent className="flex items-center justify-center h-[400px] text-destructive">
        Error loading reviews. Please try again later.
      </CardContent>
    </Card>
  );

  return (
    <div className="space-y-6 p-6">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Product Reviews</CardTitle>
            <CardDescription>View and manage all product reviews</CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          {/* Search section */}
          <div className="flex flex-col gap-4 md:flex-row md:items-center mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search reviews..."
                value={search}
                onChange={e => setSearch(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>

          {/* Reviews Table */}
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product</TableHead>
                  <TableHead>Customer</TableHead>
                  <TableHead>Rating</TableHead>
                  <TableHead>Comment</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {paginatedReviews.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-muted-foreground h-32">
                      No reviews found
                    </TableCell>
                  </TableRow>
                ) : (
                  paginatedReviews.map((review) => (
                    <TableRow key={review.id}>
                      <TableCell className="font-medium">{review.product}</TableCell>
                      <TableCell className="text-muted-foreground">{review.customer}</TableCell>
                      <TableCell>
                        <span className="flex items-center gap-1">
                          {[...Array(5)].map((_, i) => (
                            <Star key={i} className={`w-4 h-4 ${i < review.rating ? 'text-yellow-400' : 'text-muted-foreground'}`} fill={i < review.rating ? '#facc15' : 'none'} />
                          ))}
                        </span>
                      </TableCell>
                      <TableCell className="text-foreground max-w-xs truncate">{review.comment}</TableCell>
                      <TableCell className="text-muted-foreground">{review.date}</TableCell>
                      <TableCell>
                        <Select value={review.status} onValueChange={value => handleStatusUpdate(review.id, value)}>
                          <SelectTrigger className="w-[120px] h-8 text-xs">
                            <SelectValue placeholder="Status" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="Approved">Approved</SelectItem>
                            <SelectItem value="Pending">Pending</SelectItem>
                            <SelectItem value="Rejected">Rejected</SelectItem>
                          </SelectContent>
                        </Select>
                      </TableCell>
                      <TableCell className="text-right">
                          <Button variant="ghost" size="icon" onClick={() => handleDelete(review.id)} title="Delete"><Trash2 className="w-4 h-4" /></Button>
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
                    <PaginationLink isActive>{currentPage}</PaginationLink>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLink disabled className="pointer-events-none">
                      of {totalPages}
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
    </div>
  );
}