import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Search, Edit, Trash2, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";

const mockBanners = [
  { id: 1, title: "Summer Sale", image: "/vite.svg", status: "Active", start: "2025-05-01", end: "2025-06-01" },
  { id: 2, title: "New Arrivals", image: "/vite.svg", status: "Inactive", start: "2025-04-15", end: "2025-05-15" },
  { id: 3, title: "Free Shipping Week", image: "/vite.svg", status: "Active", start: "2025-05-05", end: "2025-05-12" },
];

export default function Banners() {
  const [search, setSearch] = useState("");
  const navigate = useNavigate();
  const { data: banners = [], isLoading, isError } = useQuery({
    queryKey: ["banners"],
    queryFn: async () => mockBanners,
  });

  const filtered = search
    ? banners.filter(b =>
        b.title.toLowerCase().includes(search.toLowerCase())
      )
    : banners;

  return (
    <div className="space-y-6 p-6">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Banners</CardTitle>
            <CardDescription>Manage your marketing banners</CardDescription>
          </div>
          <Button variant="default" size="sm" onClick={() => navigate("/banners/create")}> <Plus className="w-4 h-4 mr-1" /> New Banner</Button>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4 md:flex-row md:items-center mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input
                type="text"
                className="bg-muted border border-input text-foreground text-sm rounded-lg focus:ring-primary focus:border-primary block w-full pl-10 p-2.5"
                placeholder="Search banners..."
                value={search}
                onChange={e => setSearch(e.target.value)}
              />
            </div>
          </div>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Image</TableHead>
                  <TableHead>Title</TableHead>
                  <TableHead>Start</TableHead>
                  <TableHead>End</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filtered.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-muted-foreground">
                      No banners found
                    </TableCell>
                  </TableRow>
                ) : (
                  filtered.map((banner) => (
                    <TableRow key={banner.id}>
                      <TableCell><img src={banner.image} alt={banner.title} className="w-16 h-10 object-cover rounded-md border" /></TableCell>
                      <TableCell className="font-medium">{banner.title}</TableCell>
                      <TableCell className="text-muted-foreground">{banner.start}</TableCell>
                      <TableCell className="text-muted-foreground">{banner.end}</TableCell>
                      <TableCell>
                        <Badge variant={banner.status === "Active" ? "default" : "destructive"}>{banner.status}</Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button variant="ghost" size="icon" title="Edit" onClick={() => navigate(`/banners/${banner.id}/edit`)}>
                            <Edit className="h-4 w-4" />
                          </Button>
                          <Button variant="ghost" size="icon" title="Delete"><Trash2 className="h-4 w-4" /></Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}