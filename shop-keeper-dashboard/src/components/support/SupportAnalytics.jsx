import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import {
  HelpCircle,
  Clock,
  CheckCircle2,
  AlertTriangle,
  TrendingUp,
  TrendingDown,
  Users,
  MessageSquare,
} from "lucide-react";

const mockSupportData = {
  totalTickets: 156,
  openTickets: 23,
  pendingTickets: 12,
  closedTickets: 121,
  avgResponseTime: "2.4 hours",
  satisfactionScore: 4.6,
  ticketsToday: 8,
  ticketsThisWeek: 45,
  topCategories: [
    { name: "Order Issues", count: 45, percentage: 29 },
    { name: "Payment Problems", count: 32, percentage: 21 },
    { name: "Product Issues", count: 28, percentage: 18 },
    { name: "Shipping", count: 25, percentage: 16 },
    { name: "Technical Support", count: 15, percentage: 10 },
    { name: "General", count: 11, percentage: 7 }
  ],
  agentPerformance: [
    { name: "Sarah Wilson", tickets: 34, avgRating: 4.8, status: "online" },
    { name: "Mike Johnson", tickets: 28, avgRating: 4.6, status: "online" },
    { name: "David Lee", tickets: 22, avgRating: 4.5, status: "away" },
    { name: "Emma Davis", tickets: 19, avgRating: 4.7, status: "offline" }
  ],
  trends: {
    ticketsChange: "+12%",
    responseTimeChange: "-15%",
    satisfactionChange: "+8%"
  }
};

export default function SupportAnalytics() {
  const { 
    totalTickets, 
    openTickets, 
    pendingTickets, 
    closedTickets,
    avgResponseTime,
    satisfactionScore,
    ticketsToday,
    ticketsThisWeek,
    topCategories,
    agentPerformance,
    trends
  } = mockSupportData;

  const resolutionRate = ((closedTickets / totalTickets) * 100).toFixed(1);

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {/* Key Metrics */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-2">
            <HelpCircle className="h-5 w-5 text-blue-600" />
            <div className="space-y-1">
              <p className="text-2xl font-bold">{totalTickets}</p>
              <p className="text-xs text-muted-foreground">Total Tickets</p>
              <div className="flex items-center space-x-1">
                <TrendingUp className="h-3 w-3 text-green-600" />
                <span className="text-xs text-green-600">{trends.ticketsChange}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-2">
            <Clock className="h-5 w-5 text-yellow-600" />
            <div className="space-y-1">
              <p className="text-2xl font-bold">{avgResponseTime}</p>
              <p className="text-xs text-muted-foreground">Avg Response</p>
              <div className="flex items-center space-x-1">
                <TrendingDown className="h-3 w-3 text-green-600" />
                <span className="text-xs text-green-600">{trends.responseTimeChange}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-2">
            <CheckCircle2 className="h-5 w-5 text-green-600" />
            <div className="space-y-1">
              <p className="text-2xl font-bold">{resolutionRate}%</p>
              <p className="text-xs text-muted-foreground">Resolution Rate</p>
              <div className="flex items-center space-x-1">
                <TrendingUp className="h-3 w-3 text-green-600" />
                <span className="text-xs text-green-600">{trends.satisfactionChange}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-2">
            <AlertTriangle className="h-5 w-5 text-orange-600" />
            <div className="space-y-1">
              <p className="text-2xl font-bold">{openTickets + pendingTickets}</p>
              <p className="text-xs text-muted-foreground">Needs Attention</p>
              <div className="flex items-center space-x-2">
                <Badge variant="secondary" className="text-xs">
                  {openTickets} Open
                </Badge>
                <Badge variant="outline" className="text-xs">
                  {pendingTickets} Pending
                </Badge>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Ticket Categories */}
      <Card className="md:col-span-2">
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <MessageSquare className="h-5 w-5" />
            Top Categories
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {topCategories.map((category, index) => (
              <div key={index} className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="font-medium text-sm">{category.name}</div>
                  <Badge variant="outline" className="text-xs">
                    {category.count} tickets
                  </Badge>
                </div>
                <div className="flex items-center space-x-3">
                  <div className="w-32">
                    <Progress value={category.percentage} className="h-2" />
                  </div>
                  <div className="text-sm text-muted-foreground font-medium w-8">
                    {category.percentage}%
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Agent Performance */}
      <Card className="md:col-span-2">
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <Users className="h-5 w-5" />
            Agent Performance
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {agentPerformance.map((agent, index) => (
              <div key={index} className="flex items-center justify-between p-3 rounded-lg border">
                <div className="flex items-center space-x-3">
                  <div className="relative">
                    <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                      <Users className="h-5 w-5 text-primary" />
                    </div>
                    <div className={`absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-background ${
                      agent.status === 'online' ? 'bg-green-500' :
                      agent.status === 'away' ? 'bg-yellow-500' : 'bg-gray-400'
                    }`}></div>
                  </div>
                  <div>
                    <div className="font-medium text-sm">{agent.name}</div>
                    <div className="text-xs text-muted-foreground capitalize">
                      {agent.status} • {agent.tickets} tickets handled
                    </div>
                  </div>
                </div>
                <div className="text-right">
                  <div className="flex items-center space-x-1">
                    <span className="text-lg font-bold">{agent.avgRating}</span>
                    <span className="text-yellow-500">★</span>
                  </div>
                  <div className="text-xs text-muted-foreground">Avg Rating</div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Quick Stats */}
      <Card className="md:col-span-4">
        <CardHeader>
          <CardTitle className="text-lg">Quick Overview</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="text-center p-4 bg-muted/30 rounded-lg">
              <div className="text-2xl font-bold text-blue-600">{ticketsToday}</div>
              <div className="text-sm text-muted-foreground">Tickets Today</div>
            </div>
            <div className="text-center p-4 bg-muted/30 rounded-lg">
              <div className="text-2xl font-bold text-green-600">{ticketsThisWeek}</div>
              <div className="text-sm text-muted-foreground">This Week</div>
            </div>
            <div className="text-center p-4 bg-muted/30 rounded-lg">
              <div className="text-2xl font-bold text-purple-600">{satisfactionScore}</div>
              <div className="text-sm text-muted-foreground">Satisfaction Score</div>
            </div>
            <div className="text-center p-4 bg-muted/30 rounded-lg">
              <div className="text-2xl font-bold text-orange-600">{agentPerformance.filter(a => a.status === 'online').length}</div>
              <div className="text-sm text-muted-foreground">Agents Online</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}