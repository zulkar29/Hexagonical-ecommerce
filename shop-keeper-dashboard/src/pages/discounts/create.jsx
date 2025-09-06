import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";

const initialForm = {
  name: "",
  code: "",
  type: "Percentage",
  value: "",
  start: "",
  end: "",
  status: "Active",
};

export default function CreateDiscount() {
  const navigate = useNavigate();
  const [form, setForm] = useState(initialForm);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleTypeChange = (value) => {
    setForm({ ...form, type: value });
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
      navigate("/discounts");
    }, 800);
  };

  return (
    <div className="flex justify-center items-center min-h-[60vh] p-6">
      <Card className="w-full max-w-xl">
        <CardHeader>
          <CardTitle>Create Discount</CardTitle>
          <CardDescription>Fill in the details to add a new discount.</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-6" onSubmit={handleSubmit}>
            <div>
              <label className="block mb-1 font-medium" htmlFor="name">Name</label>
              <Input
                id="name"
                name="name"
                value={form.name}
                onChange={handleChange}
                required
                placeholder="Discount name"
                autoFocus
              />
            </div>
            <div>
              <label className="block mb-1 font-medium" htmlFor="code">Code</label>
              <Input
                id="code"
                name="code"
                value={form.code}
                onChange={handleChange}
                required
                placeholder="Discount code"
              />
            </div>
            <div className="flex gap-4">
              <div className="flex-1">
                <label className="block mb-1 font-medium">Type</label>
                <Select value={form.type} onValueChange={handleTypeChange}>
                  <SelectTrigger className="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="Percentage">Percentage</SelectItem>
                    <SelectItem value="Fixed">Fixed Amount</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="flex-1">
                <label className="block mb-1 font-medium" htmlFor="value">Value</label>
                <Input
                  id="value"
                  name="value"
                  type="number"
                  value={form.value}
                  onChange={handleChange}
                  required
                  placeholder={form.type === "Percentage" ? "%" : "$"}
                  min="0"
                />
              </div>
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
              <Button type="button" variant="ghost" onClick={() => navigate("/discounts")}>Cancel</Button>
              <Button type="submit" disabled={submitting}>
                {submitting ? "Creating..." : "Create Discount"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
