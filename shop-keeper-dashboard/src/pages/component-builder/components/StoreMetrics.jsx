import React from 'react';
import {
  Eye,
  ShoppingCart,
  DollarSign,
  Users,
  TrendingUp,
  Activity
} from 'lucide-react';

const StoreMetrics = ({ designData }) => {
  // Mock metrics for demo - in real app would come from analytics
  const metrics = {
    views: 1247,
    conversions: 23,
    revenue: 4567,
    customers: 89,
    components: designData?.components?.length || 0,
    lastUpdated: new Date().toLocaleDateString()
  };

  const conversionRate = ((metrics.conversions / metrics.views) * 100).toFixed(1);

  const MetricCard = ({ icon: Icon, title, value, subtitle, trend }) => (
    <div className="bg-white p-4 rounded-lg border border-gray-200 hover:shadow-md transition-shadow">
      <div className="flex items-center justify-between mb-2">
        <div className="p-2 bg-blue-50 rounded-lg">
          <Icon className="w-5 h-5 text-blue-600" />
        </div>
        {trend && (
          <div className="flex items-center text-green-600 text-sm">
            <TrendingUp className="w-4 h-4 mr-1" />
            {trend}
          </div>
        )}
      </div>
      <div className="mb-1">
        <div className="text-2xl font-bold text-gray-900">{value}</div>
        <div className="text-sm text-gray-600">{title}</div>
      </div>
      {subtitle && (
        <div className="text-xs text-gray-500">{subtitle}</div>
      )}
    </div>
  );

  return (
    <div className="p-6 bg-gray-50">
      <div className="mb-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-2">Store Performance</h3>
        <p className="text-sm text-gray-600">
          Analytics for your current store design • Last updated {metrics.lastUpdated}
        </p>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mb-6">
        <MetricCard
          icon={Eye}
          title="Page Views"
          value={metrics.views.toLocaleString()}
          subtitle="Last 30 days"
          trend="+12%"
        />

        <MetricCard
          icon={ShoppingCart}
          title="Conversions"
          value={metrics.conversions}
          subtitle={`${conversionRate}% conversion rate`}
          trend="+8%"
        />

        <MetricCard
          icon={DollarSign}
          title="Revenue"
          value={`$${metrics.revenue.toLocaleString()}`}
          subtitle="This month"
          trend="+15%"
        />

        <MetricCard
          icon={Users}
          title="Customers"
          value={metrics.customers}
          subtitle="Active customers"
        />

        <MetricCard
          icon={Activity}
          title="Components"
          value={metrics.components}
          subtitle="In current design"
        />

        <MetricCard
          icon={TrendingUp}
          title="Performance"
          value="92%"
          subtitle="Design score"
        />
      </div>

      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <h4 className="font-medium text-gray-900 mb-3">Quick Insights</h4>
        <div className="space-y-2 text-sm">
          <div className="flex items-center justify-between">
            <span className="text-gray-600">Mobile Traffic</span>
            <span className="font-medium">68%</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-gray-600">Avg. Session Duration</span>
            <span className="font-medium">2m 34s</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-gray-600">Top Product Category</span>
            <span className="font-medium">Electronics</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-gray-600">Design Engagement</span>
            <span className="font-medium text-green-600">High</span>
          </div>
        </div>
      </div>

      <div className="mt-4 text-xs text-gray-500 text-center">
        💡 Tip: Use lighter themes for better mobile performance
      </div>
    </div>
  );
};

export default StoreMetrics;