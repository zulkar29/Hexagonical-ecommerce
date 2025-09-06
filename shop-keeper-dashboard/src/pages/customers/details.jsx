import { useParams, useNavigate } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";

const mockCustomers = [
  { id: 1, name: "John Doe", email: "john@example.com", phone: "123-456-7890", status: "Active", joined: "2024-01-10", orders: 5, address: "123 Main St, Springfield, USA", notes: "VIP customer. Loves discounts." },
  { id: 2, name: "Jane Smith", email: "jane@example.com", phone: "987-654-3210", status: "Inactive", joined: "2023-11-22", orders: 2, address: "456 Oak Ave, Metropolis, USA", notes: "Requested newsletter subscription." },
  { id: 3, name: "Alice Brown", email: "alice@example.com", phone: "555-123-4567", status: "Active", joined: "2024-03-05", orders: 8, address: "789 Pine Rd, Gotham, USA", notes: "Prefers express shipping." },
];

export default function CustomerDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const customer = mockCustomers.find(c => String(c.id) === String(id));

  if (!customer) {
    return (
      <div className="space-y-6 p-6">
        <Card>
          <CardHeader>
            <CardTitle>Customer Not Found</CardTitle>
            <CardDescription>The customer you are looking for does not exist.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => navigate("/customers")}> <ArrowLeft className="mr-2 h-4 w-4" />Back to Customers</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Mock order and support history for demo
  const orderHistory = [
    { id: "ORD-001", date: "2024-04-10", total: "$120.00", status: "Completed" },
    { id: "ORD-002", date: "2024-03-22", total: "$75.50", status: "Pending" },
  ];
  const supportTickets = [
    { id: "SUP-101", subject: "Order not received", status: "Open", date: "2024-04-12" },
    { id: "SUP-099", subject: "Refund request", status: "Closed", date: "2024-03-25" },
  ];

  return (
    <div className="space-y-6 p-6">
      <Card>
        <CardHeader className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 pb-0">
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="rounded-full p-2 mr-2"
              onClick={() => navigate("/customers")}
              aria-label="Back to Customers"
            >
              <ArrowLeft className="h-5 w-5" />
            </Button>
            <CardTitle className="truncate max-w-[60vw]">{customer.name}</CardTitle>
          </div>
          <Badge variant={customer.status === "Active" ? "default" : "destructive"} className="text-base px-4 py-1">
            {customer.status}
          </Badge>
        </CardHeader>
        <CardContent className="space-y-8 px-8 pt-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <div className="font-semibold text-lg mb-2">Contact Info</div>
              <div className="grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
                <div className="font-medium">Email:</div>
                <div className="text-muted-foreground">{customer.email}</div>
                <div className="font-medium">Phone:</div>
                <div className="text-muted-foreground">{customer.phone}</div>
                <div className="font-medium">Joined:</div>
                <div className="text-muted-foreground">{customer.joined}</div>
                <div className="font-medium">Orders:</div>
                <div className="text-muted-foreground">{customer.orders}</div>
              </div>
            </div>
            <div className="space-y-2">
              <div className="font-semibold text-lg mb-2">Address</div>
              <div className="text-muted-foreground text-sm">{customer.address}</div>
              <div className="font-semibold text-lg mt-4 mb-2">Notes</div>
              <div className="text-muted-foreground text-sm">{customer.notes || 'No additional notes.'}</div>
            </div>
          </div>

          {/* Order History */}
          <div>
            <div className="font-semibold text-lg mb-2">Order History</div>
            <div className="rounded-md border overflow-x-auto">
              <table className="min-w-full text-sm">
                <thead className="bg-muted">
                  <tr>
                    <th className="px-4 py-2 text-left">Order ID</th>
                    <th className="px-4 py-2 text-left">Date</th>
                    <th className="px-4 py-2 text-left">Total</th>
                    <th className="px-4 py-2 text-left">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {orderHistory.map((order) => (
                    <tr key={order.id}>
                      <td className="px-4 py-2">{order.id}</td>
                      <td className="px-4 py-2">{order.date}</td>
                      <td className="px-4 py-2">{order.total}</td>
                      <td className="px-4 py-2">
                        <Badge variant={order.status === "Completed" ? "default" : order.status === "Pending" ? "secondary" : "destructive"}>{order.status}</Badge>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Support Ticket History */}
          <div>
            <div className="font-semibold text-lg mb-2">Support Tickets</div>
            <div className="rounded-md border overflow-x-auto">
              <table className="min-w-full text-sm">
                <thead className="bg-muted">
                  <tr>
                    <th className="px-4 py-2 text-left">Ticket ID</th>
                    <th className="px-4 py-2 text-left">Subject</th>
                    <th className="px-4 py-2 text-left">Status</th>
                    <th className="px-4 py-2 text-left">Date</th>
                  </tr>
                </thead>
                <tbody>
                  {supportTickets.map((ticket) => (
                    <tr key={ticket.id}>
                      <td className="px-4 py-2">{ticket.id}</td>
                      <td className="px-4 py-2">{ticket.subject}</td>
                      <td className="px-4 py-2">
                        <Badge variant={ticket.status === "Open" ? "secondary" : "default"}>{ticket.status}</Badge>
                      </td>
                      <td className="px-4 py-2">{ticket.date}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
