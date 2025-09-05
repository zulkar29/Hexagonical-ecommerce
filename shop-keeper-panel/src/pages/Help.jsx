import React, { useState } from 'react';
import {
  Search,
  Book,
  MessageCircle,
  Phone,
  Mail,
  Video,
  FileText,
  ExternalLink,
  ChevronRight,
  ChevronDown,
  Star,
  ThumbsUp,
  ThumbsDown,
  Download,
  Play,
  Clock,
  Users,
  HelpCircle,
  Lightbulb,
  AlertCircle,
  CheckCircle,
  ArrowRight,
  Bookmark,
  Share,
  Send,
  User,
  Building,
  CreditCard,
  Package,
  BarChart3,
  Settings,
  Shield,
  Database,
  Smartphone,
  Globe,
  Zap,
  Target,
  TrendingUp,
  Award,
  Coffee,
  Heart,
  MessageSquare,
  Calendar,
  MapPin,
  Headphones
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

const HelpPage = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [showContactDialog, setShowContactDialog] = useState(false);
  const [contactType, setContactType] = useState('general');

  // Help categories
  const categories = [
    {
      id: 'getting-started',
      name: 'Getting Started',
      icon: Zap,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
      articles: 12
    },
    {
      id: 'pos',
      name: 'Point of Sale',
      icon: CreditCard,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
      articles: 8
    },
    {
      id: 'inventory',
      name: 'Inventory Management',
      icon: Package,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
      articles: 15
    },
    {
      id: 'customers',
      name: 'Customer Management',
      icon: Users,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
      articles: 10
    },
    {
      id: 'reports',
      name: 'Reports & Analytics',
      icon: BarChart3,
      color: 'text-indigo-600',
      bgColor: 'bg-indigo-50',
      articles: 7
    },
    {
      id: 'settings',
      name: 'Settings & Configuration',
      icon: Settings,
      color: 'text-gray-600',
      bgColor: 'bg-gray-50',
      articles: 6
    }
  ];

  // FAQ data
  const faqs = [
    {
      id: 1,
      category: 'getting-started',
      question: 'How do I set up my business information?',
      answer: 'Go to Settings > Business tab and fill in your business name, address, phone number, and other details. This information will appear on receipts and invoices.',
      helpful: 45,
      notHelpful: 2
    },
    {
      id: 2,
      category: 'pos',
      question: 'How do I process a sale with bKash payment?',
      answer: 'In the POS screen, add items to cart, click "Process Payment", select "bKash" as payment method, and complete the transaction. The customer will receive a receipt.',
      helpful: 52,
      notHelpful: 1
    },
    {
      id: 3,
      category: 'inventory',
      question: 'How do I add new products to my inventory?',
      answer: 'Navigate to Inventory > Products, click "Add Product", fill in product details including name, price, stock quantity, and category. You can also upload a product image.',
      helpful: 38,
      notHelpful: 3
    },
    {
      id: 4,
      category: 'customers',
      question: 'How do I track customer purchase history?',
      answer: 'Go to Customers page, select a customer, and view their detailed profile. You\'ll see complete purchase history, favorite products, and loyalty points.',
      helpful: 29,
      notHelpful: 0
    },
    {
      id: 5,
      category: 'pos',
      question: 'Can I give discounts to customers?',
      answer: 'Yes! In the POS cart section, you can apply discounts either as a percentage or fixed amount. The discount will be reflected in the total and on the receipt.',
      helpful: 67,
      notHelpful: 2
    },
    {
      id: 6,
      category: 'inventory',
      question: 'How do I set up low stock alerts?',
      answer: 'When adding/editing products, set a minimum stock level. The system will automatically alert you when stock falls below this threshold.',
      helpful: 41,
      notHelpful: 1
    },
    {
      id: 7,
      category: 'reports',
      question: 'How do I generate sales reports?',
      answer: 'Visit the Reports section from the sidebar. You can generate daily, weekly, or monthly sales reports and export them as PDF or Excel files.',
      helpful: 33,
      notHelpful: 0
    },
    {
      id: 8,
      category: 'settings',
      question: 'How do I backup my data?',
      answer: 'Go to Settings > Backup tab. You can enable automatic backups or create manual backups. Data can be stored locally or in cloud services like Google Drive.',
      helpful: 28,
      notHelpful: 1
    }
  ];

  // Video tutorials
  const videoTutorials = [
    {
      id: 1,
      title: 'ShopOS Overview - Complete Tour',
      duration: '12:45',
      category: 'getting-started',
      thumbnail: '/api/placeholder/320/180',
      views: 1250,
      description: 'Complete walkthrough of ShopOS features and capabilities for Bangladesh businesses.'
    },
    {
      id: 2,
      title: 'Setting Up Your First Sale',
      duration: '8:30',
      category: 'pos',
      thumbnail: '/api/placeholder/320/180',
      views: 890,
      description: 'Step-by-step guide to process your first sale using the POS system.'
    },
    {
      id: 3,
      title: 'Managing Inventory Like a Pro',
      duration: '15:20',
      category: 'inventory',
      thumbnail: '/api/placeholder/320/180',
      views: 675,
      description: 'Advanced inventory management techniques and best practices.'
    },
    {
      id: 4,
      title: 'Customer Management & Loyalty',
      duration: '10:15',
      category: 'customers',
      thumbnail: '/api/placeholder/320/180',
      views: 542,
      description: 'Build customer relationships and manage loyalty programs effectively.'
    }
  ];

  // Guides and articles
  const guides = [
    {
      id: 1,
      title: 'Complete Setup Guide for New Users',
      category: 'getting-started',
      readTime: '15 min read',
      difficulty: 'Beginner',
      description: 'Everything you need to know to get started with ShopOS for your Bangladesh business.',
      lastUpdated: '2025-01-20'
    },
    {
      id: 2,
      title: 'Advanced POS Features and Tips',
      category: 'pos',
      readTime: '12 min read',
      difficulty: 'Intermediate',
      description: 'Master advanced POS features including bulk sales, custom receipts, and payment methods.',
      lastUpdated: '2025-01-18'
    },
    {
      id: 3,
      title: 'Inventory Optimization Strategies',
      category: 'inventory',
      readTime: '20 min read',
      difficulty: 'Advanced',
      description: 'Learn how to optimize your inventory levels and reduce waste while maximizing profits.',
      lastUpdated: '2025-01-15'
    },
    {
      id: 4,
      title: 'Building Customer Loyalty in Bangladesh',
      category: 'customers',
      readTime: '18 min read',
      difficulty: 'Intermediate',
      description: 'Strategies specific to Bangladesh market for building and maintaining customer loyalty.',
      lastUpdated: '2025-01-12'
    },
    {
      id: 5,
      title: 'Understanding Your Business Reports',
      category: 'reports',
      readTime: '14 min read',
      difficulty: 'Beginner',
      description: 'Make sense of your business data with comprehensive reporting and analytics.',
      lastUpdated: '2025-01-10'
    }
  ];

  // Filter FAQs based on search and category
  const filteredFaqs = faqs.filter(faq => {
    const matchesSearch = faq.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         faq.answer.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesCategory = selectedCategory === 'all' || faq.category === selectedCategory;
    return matchesSearch && matchesCategory;
  });

  const getDifficultyColor = (difficulty) => {
    switch (difficulty) {
      case 'Beginner': return 'bg-green-100 text-green-700';
      case 'Intermediate': return 'bg-yellow-100 text-yellow-700';
      case 'Advanced': return 'bg-red-100 text-red-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      {/* Header */}
      <div className="text-center space-y-4">
        <div className="flex items-center justify-center space-x-3">
          <div className="w-12 h-12 bg-primary rounded-full flex items-center justify-center">
            <HelpCircle className="h-6 w-6 text-primary-foreground" />
          </div>
          <div>
            <h1 className="text-3xl font-bold text-foreground">Help & Support Center</h1>
            <p className="text-muted-foreground">
              Get help, learn features, and grow your Bangladesh business
            </p>
          </div>
        </div>

        {/* Search Bar */}
        <div className="max-w-2xl mx-auto">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-5 w-5" />
            <Input
              placeholder="Search for help articles, guides, and tutorials..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 py-3 text-lg"
            />
          </div>
        </div>

        {/* Quick Access Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-4xl mx-auto">
          <Card className="hover:shadow-md transition-shadow cursor-pointer">
            <CardContent className="p-6 text-center">
              <Book className="h-8 w-8 text-blue-600 mx-auto mb-3" />
              <h3 className="font-semibold mb-2">User Guide</h3>
              <p className="text-sm text-muted-foreground">
                Complete documentation and tutorials
              </p>
            </CardContent>
          </Card>
          <Card className="hover:shadow-md transition-shadow cursor-pointer">
            <CardContent className="p-6 text-center">
              <Video className="h-8 w-8 text-green-600 mx-auto mb-3" />
              <h3 className="font-semibold mb-2">Video Tutorials</h3>
              <p className="text-sm text-muted-foreground">
                Step-by-step video guides
              </p>
            </CardContent>
          </Card>
          <Card className="hover:shadow-md transition-shadow cursor-pointer">
            <CardContent className="p-6 text-center">
              <MessageCircle className="h-8 w-8 text-purple-600 mx-auto mb-3" />
              <h3 className="font-semibold mb-2">Live Support</h3>
              <p className="text-sm text-muted-foreground">
                Chat with our support team
              </p>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Main Content Tabs */}
      <Card>
        <CardContent className="p-0">
          <Tabs defaultValue="faq" className="w-full">
            <div className="border-b">
              <TabsList className="grid w-full grid-cols-5 bg-transparent">
                <TabsTrigger value="faq" className="flex items-center space-x-2">
                  <HelpCircle className="h-4 w-4" />
                  <span>FAQ</span>
                </TabsTrigger>
                <TabsTrigger value="guides" className="flex items-center space-x-2">
                  <Book className="h-4 w-4" />
                  <span>Guides</span>
                </TabsTrigger>
                <TabsTrigger value="videos" className="flex items-center space-x-2">
                  <Video className="h-4 w-4" />
                  <span>Videos</span>
                </TabsTrigger>
                <TabsTrigger value="contact" className="flex items-center space-x-2">
                  <MessageCircle className="h-4 w-4" />
                  <span>Contact</span>
                </TabsTrigger>
                <TabsTrigger value="community" className="flex items-center space-x-2">
                  <Users className="h-4 w-4" />
                  <span>Community</span>
                </TabsTrigger>
              </TabsList>
            </div>

            {/* FAQ Tab */}
            <TabsContent value="faq" className="p-6 space-y-6">
              {/* Categories */}
              <div>
                <h3 className="text-lg font-semibold mb-4">Browse by Category</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {categories.map((category) => (
                    <Card 
                      key={category.id} 
                      className={cn(
                        "cursor-pointer hover:shadow-md transition-shadow",
                        selectedCategory === category.id && "ring-2 ring-primary"
                      )}
                      onClick={() => setSelectedCategory(category.id)}
                    >
                      <CardContent className="p-4">
                        <div className="flex items-center space-x-3">
                          <div className={cn("p-2 rounded-lg", category.bgColor)}>
                            <category.icon className={cn("h-5 w-5", category.color)} />
                          </div>
                          <div className="flex-1">
                            <h4 className="font-medium">{category.name}</h4>
                            <p className="text-sm text-muted-foreground">
                              {category.articles} articles
                            </p>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>

              {/* FAQ List */}
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold">
                    Frequently Asked Questions
                    {selectedCategory !== 'all' && (
                      <Badge variant="secondary" className="ml-2">
                        {categories.find(c => c.id === selectedCategory)?.name}
                      </Badge>
                    )}
                  </h3>
                  {selectedCategory !== 'all' && (
                    <Button 
                      variant="outline" 
                      size="sm"
                      onClick={() => setSelectedCategory('all')}
                    >
                      Show All
                    </Button>
                  )}
                </div>
                <Accordion type="single" collapsible className="space-y-2">
                  {filteredFaqs.map((faq) => (
                    <AccordionItem key={faq.id} value={faq.id.toString()}>
                      <AccordionTrigger className="text-left hover:no-underline">
                        <div className="flex items-start space-x-3 text-left">
                          <div className="w-6 h-6 bg-primary/10 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                            <span className="text-xs font-medium text-primary">Q</span>
                          </div>
                          <span className="font-medium">{faq.question}</span>
                        </div>
                      </AccordionTrigger>
                      <AccordionContent>
                        <div className="pl-9 space-y-4">
                          <p className="text-muted-foreground leading-relaxed">
                            {faq.answer}
                          </p>
                          <div className="flex items-center justify-between pt-2 border-t">
                            <div className="flex items-center space-x-4">
                              <span className="text-sm text-muted-foreground">
                                Was this helpful?
                              </span>
                              <div className="flex items-center space-x-2">
                                <Button variant="ghost" size="sm">
                                  <ThumbsUp className="h-4 w-4 mr-1" />
                                  {faq.helpful}
                                </Button>
                                <Button variant="ghost" size="sm">
                                  <ThumbsDown className="h-4 w-4 mr-1" />
                                  {faq.notHelpful}
                                </Button>
                              </div>
                            </div>
                            <div className="flex items-center space-x-2">
                              <Button variant="ghost" size="sm">
                                <Share className="h-4 w-4" />
                              </Button>
                              <Button variant="ghost" size="sm">
                                <Bookmark className="h-4 w-4" />
                              </Button>
                            </div>
                          </div>
                        </div>
                      </AccordionContent>
                    </AccordionItem>
                  ))}
                </Accordion>
              </div>
            </TabsContent>

            {/* Guides Tab */}
            <TabsContent value="guides" className="p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold mb-4">User Guides & Articles</h3>
                <div className="space-y-4">
                  {guides.map((guide) => (
                    <Card key={guide.id} className="hover:shadow-md transition-shadow">
                      <CardContent className="p-6">
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <div className="flex items-center space-x-3 mb-2">
                              <h4 className="text-lg font-semibold">{guide.title}</h4>
                              <Badge className={getDifficultyColor(guide.difficulty)}>
                                {guide.difficulty}
                              </Badge>
                            </div>
                            <p className="text-muted-foreground mb-3">
                              {guide.description}
                            </p>
                            <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                              <div className="flex items-center space-x-1">
                                <Clock className="h-4 w-4" />
                                <span>{guide.readTime}</span>
                              </div>
                              <div className="flex items-center space-x-1">
                                <Calendar className="h-4 w-4" />
                                <span>Updated {guide.lastUpdated}</span>
                              </div>
                            </div>
                          </div>
                          <div className="flex items-center space-x-2 ml-4">
                            <Button variant="outline" size="sm">
                              <Bookmark className="h-4 w-4" />
                            </Button>
                            <Button>
                              Read Guide
                              <ArrowRight className="h-4 w-4 ml-2" />
                            </Button>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>

              {/* Quick Start Guide */}
              <Card className="bg-gradient-to-r from-blue-50 to-purple-50 border-blue-200">
                <CardContent className="p-6">
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                      <Zap className="h-6 w-6 text-blue-600" />
                    </div>
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold mb-2">Quick Start Guide</h3>
                      <p className="text-muted-foreground mb-3">
                        New to ShopOS? Get up and running in just 10 minutes with our comprehensive quick start guide.
                      </p>
                      <Button className="bg-blue-600 hover:bg-blue-700">
                        <Play className="h-4 w-4 mr-2" />
                        Start Now
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            {/* Videos Tab */}
            <TabsContent value="videos" className="p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold mb-4">Video Tutorials</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {videoTutorials.map((video) => (
                    <Card key={video.id} className="overflow-hidden hover:shadow-md transition-shadow">
                      <div className="relative">
                        <img 
                          src={video.thumbnail} 
                          alt={video.title}
                          className="w-full h-48 object-cover"
                        />
                        <div className="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity cursor-pointer">
                          <div className="w-16 h-16 bg-white rounded-full flex items-center justify-center">
                            <Play className="h-8 w-8 text-gray-900 ml-1" />
                          </div>
                        </div>
                        <div className="absolute bottom-2 right-2 bg-black/80 text-white px-2 py-1 rounded text-sm">
                          {video.duration}
                        </div>
                      </div>
                      <CardContent className="p-4">
                        <h4 className="font-semibold mb-2">{video.title}</h4>
                        <p className="text-sm text-muted-foreground mb-3">
                          {video.description}
                        </p>
                        <div className="flex items-center justify-between text-sm text-muted-foreground">
                          <div className="flex items-center space-x-1">
                            <Users className="h-4 w-4" />
                            <span>{video.views.toLocaleString()} views</span>
                          </div>
                          <Badge variant="outline">
                            {categories.find(c => c.id === video.category)?.name}
                          </Badge>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>

              {/* Featured Playlist */}
              <Card className="bg-gradient-to-r from-green-50 to-blue-50 border-green-200">
                <CardContent className="p-6">
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                      <Video className="h-6 w-6 text-green-600" />
                    </div>
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold mb-2">Complete ShopOS Masterclass</h3>
                      <p className="text-muted-foreground mb-3">
                        A comprehensive video series covering everything from setup to advanced features. Perfect for Bangladesh business owners.
                      </p>
                      <Button className="bg-green-600 hover:bg-green-700">
                        <Play className="h-4 w-4 mr-2" />
                        Watch Playlist (8 videos)
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            {/* Contact Tab */}
            <TabsContent value="contact" className="p-6 space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Contact Methods */}
                <div className="space-y-4">
                  <h3 className="text-lg font-semibold">Get in Touch</h3>
                  <div className="space-y-4">
                    <Card>
                      <CardContent className="p-4">
                        <div className="flex items-center space-x-3">
                          <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                            <MessageCircle className="h-5 w-5 text-blue-600" />
                          </div>
                          <div className="flex-1">
                            <h4 className="font-medium">Live Chat</h4>
                            <p className="text-sm text-muted-foreground">
                              Chat with our support team
                            </p>
                            <p className="text-xs text-green-600 font-medium">• Online now</p>
                          </div>
                          <Button>Start Chat</Button>
                        </div>
                      </CardContent>
                    </Card>

                    <Card>
                      <CardContent className="p-4">
                        <div className="flex items-center space-x-3">
                          <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">
                            <Phone className="h-5 w-5 text-green-600" />
                          </div>
                          <div className="flex-1">
                            <h4 className="font-medium">Phone Support</h4>
                            <p className="text-sm text-muted-foreground">
                              +880 1700-000000
                            </p>
                            <p className="text-xs text-muted-foreground">
                              Sat-Thu, 9 AM - 6 PM (Dhaka time)
                            </p>
                          </div>
                          <Button variant="outline">Call Now</Button>
                        </div>
                      </CardContent>
                    </Card>

                    <Card>
                      <CardContent className="p-4">
                        <div className="flex items-center space-x-3">
                          <div className="w-10 h-10 bg-purple-100 rounded-full flex items-center justify-center">
                            <Mail className="h-5 w-5 text-purple-600" />
                          </div>
                          <div className="flex-1">
                            <h4 className="font-medium">Email Support</h4>
                            <p className="text-sm text-muted-foreground">
                              support@shopos.com.bd
                            </p>
                            <p className="text-xs text-muted-foreground">
                              Response within 24 hours
                            </p>
                          </div>
                          <Button variant="outline">Send Email</Button>
                        </div>
                      </CardContent>
                    </Card>

                    <Card>
                      <CardContent className="p-4">
                        <div className="flex items-center space-x-3">
                          <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                            <MapPin className="h-5 w-5 text-orange-600" />
                          </div>
                          <div className="flex-1">
                            <h4 className="font-medium">Office Visit</h4>
                            <p className="text-sm text-muted-foreground">
                              House 123, Road 45, Gulshan-2, Dhaka
                            </p>
                            <p className="text-xs text-muted-foreground">
                              By appointment only
                            </p>
                          </div>
                          <Button variant="outline">Book Visit</Button>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </div>

                {/* Contact Form */}
                <div>
                  <h3 className="text-lg font-semibold mb-4">Send us a Message</h3>
                  <Card>
                    <CardContent className="p-6">
                      <form onSubmit={(e) => {
                        e.preventDefault();
                        alert('Contact form submitted!');
                      }}>
                        <div className="space-y-4">
                          <div className="grid grid-cols-2 gap-4">
                            <div>
                              <Label htmlFor="name">Your Name</Label>
                              <Input id="name" placeholder="Enter your name" />
                            </div>
                            <div>
                              <Label htmlFor="email">Email</Label>
                              <Input id="email" type="email" placeholder="you@email.com" />
                            </div>
                          </div>
                          <div>
                            <Label htmlFor="subject">Subject</Label>
                            <Input id="subject" placeholder="How can we help you?" />
                          </div>
                          <div>
                            <Label htmlFor="message">Message</Label>
                            <Textarea id="message" placeholder="Type your message here..." rows={4} />
                          </div>
                          <Button type="submit" className="w-full mt-2">Send Message</Button>
                        </div>
                      </form>
                    </CardContent>
                  </Card>
                </div>
              </div>
            </TabsContent>

            {/* Community Tab */}
            <TabsContent value="community" className="p-6 space-y-6">
              <div className="text-center space-y-4">
                <Users className="h-10 w-10 mx-auto text-primary mb-2" />
                <h3 className="text-xl font-semibold">Join the ShopOS Community</h3>
                <p className="text-muted-foreground mb-4">
                  Connect with other business owners, share tips, and get peer support.
                </p>
                <Button asChild>
                  <a href="https://facebook.com/groups/shopos.bd" target="_blank" rel="noopener noreferrer">
                    Join Facebook Group
                  </a>
                </Button>
              </div>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  );
};

export default HelpPage;