import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";

const initialForm = {
  title: "",
  description: "",
  image: "",
  status: "Active",
  start: "",
  end: "",
};

export default function CreateBanner() {
  const navigate = useNavigate();
  const [form, setForm] = useState(initialForm);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

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
          <CardTitle>Create Banner</CardTitle>
          <CardDescription>Fill in the details to add a new banner.</CardDescription>
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
                {submitting ? "Creating..." : "Create Banner"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
