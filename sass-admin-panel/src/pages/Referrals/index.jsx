import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Users,
  DollarSign,
  TrendingUp,
  Share2,
  QrCode,
  Copy,
  Eye,
  Calendar,
  Target,
  Code,
  BarChart3,
  History,
  Wallet
} from 'lucide-react';
import { toast } from 'sonner';
import ReferralCodeGenerator from './ReferralCodeGenerator';
import ReferralStatistics from './ReferralStatistics';
import CommissionTracking from './CommissionTracking';
import ReferralHistory from './ReferralHistory';

const ReferralsHome = () => {
  const [referralStats, setReferralStats] = useState({
    totalReferrals: 0,
    activeReferrals: 0,
    totalCommissions: 0,
    pendingCommissions: 0,
    conversionRate: 0,
    thisMonthReferrals: 0
  });

  const [recentReferrals, setRecentReferrals] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('overview');

  useEffect(() => {
    fetchReferralData();
  }, []);

  const fetchReferralData = async () => {
    try {
      setLoading(true);
      // Mock data - replace with actual API calls
      setReferralStats({
        totalReferrals: 156,
        activeReferrals: 89,
        totalCommissions: 12450,
        pendingCommissions: 2340,
        conversionRate: 24.5,
        thisMonthReferrals: 23
      });

      setRecentReferrals([
        {
          id: 1,
          referralCode: 'REF001',
          referredUser: 'john@example.com',
          status: 'active',
          commission: 150,
          createdAt: '2024-01-15',
          convertedAt: '2024-01-16'
        },
        {
          id: 2,
          referralCode: 'REF002',
          referredUser: 'jane@example.com',
          status: 'pending',
          commission: 0,
          createdAt: '2024-01-14',
          convertedAt: null
        }
      ]);
    } catch (error) {
      toast.error('Failed to fetch referral data');
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  const getStatusBadge = (status) => {
    const variants = {
      active: 'default',
      pending: 'secondary',
      expired: 'destructive'
    };
    return <Badge variant={variants[status]}>{status}</Badge>;
  };

  const statsCards = [
    {
      title: 'Total Referrals',
      value: referralStats.totalReferrals,
      icon: Users,
      description: '+12% from last month',
      color: 'text-blue-500'
    },
    {
      title: 'Active Referrals',
      value: referralStats.activeReferrals,
      icon: Target,
      description: `${referralStats.conversionRate}% conversion rate`,
      color: 'text-green-500'
    },
    {
      title: 'Total Commissions',
      value: formatCurrency(referralStats.totalCommissions),
      icon: DollarSign,
      description: '+8% from last month',
      color: 'text-emerald-500'
    },
    {
      title: 'Pending Commissions',
      value: formatCurrency(referralStats.pendingCommissions),
      icon: TrendingUp,
      description: 'Awaiting payment',
      color: 'text-orange-500'
    }
  ];

  if (loading) {
    return (
      <div className="p-6 space-y-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-32 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Referral Management</h1>
          <p className="text-muted-foreground">
            Manage referral codes, track commissions, and monitor performance
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <QrCode className="h-4 w-4 mr-2" />
            Generate QR Code
          </Button>
          <Button>
            <Share2 className="h-4 w-4 mr-2" />
            Create Referral
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {statsCards.map((stat, index) => {
          const Icon = stat.icon;
          return (
            <Card key={index}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
                <Icon className={`h-4 w-4 ${stat.color}`} />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                <p className="text-xs text-muted-foreground">{stat.description}</p>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="codes">Referral Codes</TabsTrigger>
          <TabsTrigger value="commissions">Commissions</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="history">History</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Recent Referrals */}
            <Card>
              <CardHeader>
                <CardTitle>Recent Referrals</CardTitle>
                <CardDescription>
                  Latest referral activities and conversions
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {recentReferrals.map((referral) => (
                    <div key={referral.id} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="font-medium">{referral.referralCode}</span>
                          {getStatusBadge(referral.status)}
                        </div>
                        <p className="text-sm text-muted-foreground">{referral.referredUser}</p>
                        <p className="text-xs text-muted-foreground">
                          Created: {referral.createdAt}
                        </p>
                      </div>
                      <div className="text-right">
                        <p className="font-medium">{formatCurrency(referral.commission)}</p>
                        <Button variant="ghost" size="sm">
                          <Eye className="h-3 w-3" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Quick Actions */}
            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
                <CardDescription>
                  Navigate to different sections
                </CardDescription>
              </CardHeader>
              <CardContent className="grid grid-cols-2 gap-3">
                <Button className="justify-start" variant="outline" onClick={() => setActiveTab('codes')}>
                  <Code className="h-4 w-4 mr-2" />
                  Manage Codes
                </Button>
                <Button className="justify-start" variant="outline" onClick={() => setActiveTab('analytics')}>
                  <BarChart3 className="h-4 w-4 mr-2" />
                  Analytics
                </Button>
                <Button className="justify-start" variant="outline" onClick={() => setActiveTab('commissions')}>
                  <Wallet className="h-4 w-4 mr-2" />
                  Commissions
                </Button>
                <Button className="justify-start" variant="outline" onClick={() => setActiveTab('history')}>
                  <History className="h-4 w-4 mr-2" />
                  History
                </Button>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="codes">
          <ReferralCodeGenerator />
        </TabsContent>

        <TabsContent value="commissions">
          <CommissionTracking />
        </TabsContent>

        <TabsContent value="analytics">
          <ReferralStatistics />
        </TabsContent>
        
        <TabsContent value="history">
          <ReferralHistory />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default ReferralsHome;