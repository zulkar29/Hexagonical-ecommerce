import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";

const mockBanners = [
  { id: 1, title: "Summer Sale", image: "/vite.svg", status: "Active", start: "2025-05-01", end: "2025-06-01", description: "Up to 50% off!" },
  { id: 2, title: "New Arrivals", image: "/vite.svg", status: "Inactive", start: "2025-04-15", end: "2025-05-15", description: "Check out our latest products." },
  { id: 3, title: "Free Shipping Week", image: "/vite.svg", status: "Active", start: "2025-05-05", end: "2025-05-12", description: "Enjoy free shipping all week!" },
];

export default function EditBanner() {
  const { id } = useParams();
  const navigate = useNavigate();
  const banner = mockBanners.find((b) => String(b.id) === String(id));

  const [form, setForm] = useState(
    banner || { title: "", description: "", image: "", status: "Active", start: "", end: "" }
  );
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  if (!banner) {
    return (
      <div className="flex justify-center items-center min-h-[60vh] p-6">
        <Card className="w-full max-w-xl">
          <CardHeader>
            <CardTitle>Banner Not Found</CardTitle>
            <CardDescription>The banner you are trying to edit does not exist.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => navigate("/banners")}>Back to Banners</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleStatusChange = (value) => {
    setForm({ ...form, status: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);
    setError("");
    // Simulate API call
    setTimeout(() => {
      setSubmitting(false);
      navigate("/banners");
    }, 800);
  };

  return (
    <div className="flex justify-center items-center min-h-[60vh] p-6">
      <Card className="w-full max-w-xl">
        <CardHeader>
          <CardTitle>Edit Banner</CardTitle>
          <CardDescription>Update the details for this banner.</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-6" onSubmit={handleSubmit}>
            <div>
              <label className="block mb-1 font-medium" htmlFor="title">Title</label>
              <Input
                id="title"
                name="title"
                value={form.title}
                onChange={handleChange}
                required
                placeholder="Banner title"
                autoFocus
              />
            </div>
            <div>
              <label className="block mb-1 font-medium" htmlFor="description">Description</label>
              <Input
                id="description"
                name="description"
                value={form.description}
                onChange={handleChange}
                required
                placeholder="Short description"
              />
            </div>
            <div>
              <label className="block mb-1 font-medium" htmlFor="image">Image URL</label>
              <Input
                id="image"
                name="image"
                value={form.image}
                onChange={handleChange}
                required
                placeholder="Paste image URL or upload"
              />
            </div>
            <div className="flex gap-4">
              <div className="flex-1">
                <label className="block mb-1 font-medium" htmlFor="start">Start Date</label>
                <Input
                  id="start"
                  name="start"
                  type="date"
                  value={form.start}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className="flex-1">
                <label className="block mb-1 font-medium" htmlFor="end">End Date</label>
                <Input
                  id="end"
                  name="end"
                  type="date"
                  value={form.end}
                  onChange={handleChange}
                  required
                />
              </div>
            </div>
            <div>
              <label className="block mb-1 font-medium">Status</label>
              <Select value={form.status} onValueChange={handleStatusChange}>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Active">Active</SelectItem>
                  <SelectItem value="Inactive">Inactive</SelectItem>
                </SelectContent>
              </Select>
            </div>
            {error && <div className="text-destructive text-sm">{error}</div>}
            <div className="flex justify-end gap-2">
              <Button type="button" variant="ghost" onClick={() => navigate("/banners")}>Cancel</Button>
              <Button type="submit" disabled={submitting}>
                {submitting ? "Saving..." : "Save Changes"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
