'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { 
  ArrowRight, 
  Check, 
  Star, 
  Shield, 
  Globe, 
  Palette, 
  BarChart3,
  Smartphone,
  ShoppingCart,
  Settings,
  Play,
  ChevronRight,
  Code,
  Building,
  Layers
} from 'lucide-react';

export default function HomePage() {
  const [activeDemo, setActiveDemo] = useState(0);
  const [stats, setStats] = useState({ customers: 0, stores: 0, revenue: 0 });

  // Animate stats on mount
  useEffect(() => {
    const animateStats = () => {
      const targets = { customers: 2500, stores: 450, revenue: 120 };
      const duration = 2000;
      const steps = 60;
      const stepDuration = duration / steps;
      
      let step = 0;
      const timer = setInterval(() => {
        step++;
        const progress = step / steps;
        
        setStats({
          customers: Math.floor(targets.customers * progress),
          stores: Math.floor(targets.stores * progress),
          revenue: Math.floor(targets.revenue * progress)
        });
        
        if (step >= steps) clearInterval(timer);
      }, stepDuration);
    };
    
    animateStats();
  }, []);

  const demos = [
    {
      name: 'Fashion Store',
      image: '/demos/fashion-demo.jpg',
      features: ['Size Charts', 'Color Swatches', 'Lookbooks'],
      color: 'from-pink-500 to-rose-600'
    },
    {
      name: 'Electronics Shop',
      image: '/demos/electronics-demo.jpg', 
      features: ['Spec Comparison', 'Tech Reviews', 'Warranties'],
      color: 'from-blue-500 to-cyan-600'
    },
    {
      name: 'Marketplace',
      image: '/demos/marketplace-demo.jpg',
      features: ['Multi-Vendor', 'Commission Tracking', 'Seller Tools'],
      color: 'from-green-500 to-emerald-600'
    }
  ];

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative min-h-screen flex items-center bg-gradient-to-br from-blue-50 via-white to-purple-50 overflow-hidden">
        {/* Background Elements */}
        <div className="absolute inset-0">
          <div className="absolute top-20 right-20 w-72 h-72 bg-blue-200/30 rounded-full blur-3xl"></div>
          <div className="absolute bottom-20 left-20 w-96 h-96 bg-purple-200/30 rounded-full blur-3xl"></div>
          <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-gradient-to-r from-blue-100/20 to-purple-100/20 rounded-full blur-3xl"></div>
        </div>
        
        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            {/* Left Content */}
            <div className="text-center lg:text-left">
              <div className="inline-flex items-center gap-2 bg-blue-100 text-blue-800 px-4 py-2 rounded-full text-sm font-semibold mb-6 animate-fade-in">
                <Layers className="w-4 h-4" />
                Hexagonal Architecture • Enterprise SaaS
              </div>
              
              <h1 className="text-5xl lg:text-6xl font-bold text-gray-900 mb-6 leading-tight animate-slide-up">
                Launch Your 
                <span className="bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent"> Ecommerce Empire</span> 
                in Minutes
              </h1>
              
              <p className="text-xl text-gray-600 mb-8 leading-relaxed animate-slide-up" style={{animationDelay: '0.1s'}}>
                The most powerful SaaS ecommerce platform with beautiful themes, 
                advanced features, and enterprise-grade hexagonal architecture. 
                Choose your theme, customize your brand, and start selling instantly.
              </p>
              
              <div className="flex flex-col sm:flex-row gap-4 mb-8 animate-slide-up" style={{animationDelay: '0.2s'}}>
                <Link
                  href="/get-started"
                  className="group inline-flex items-center justify-center gap-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white px-8 py-4 rounded-full hover:from-blue-700 hover:to-purple-700 transition-all duration-300 transform hover:scale-105 shadow-2xl font-semibold text-lg"
                >
                  Start Free Trial
                  <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                </Link>
                
                <button className="group inline-flex items-center justify-center gap-3 border-2 border-gray-300 text-gray-700 px-8 py-4 rounded-full hover:border-blue-300 hover:text-blue-600 hover:bg-blue-50 transition-all duration-300 backdrop-blur-sm">
                  <Play className="w-5 h-5" />
                  Watch Demo
                </button>
              </div>
              
              {/* Stats */}
              <div className="flex flex-col sm:flex-row gap-8 justify-center lg:justify-start animate-fade-in" style={{animationDelay: '0.3s'}}>
                <div className="text-center">
                  <div className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                    {stats.customers.toLocaleString()}+
                  </div>
                  <div className="text-gray-600 font-medium">Happy Customers</div>
                </div>
                <div className="text-center">
                  <div className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                    {stats.stores.toLocaleString()}+
                  </div>
                  <div className="text-gray-600 font-medium">Active Stores</div>
                </div>
                <div className="text-center">
                  <div className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                    ${stats.revenue.toLocaleString()}M+
                  </div>
                  <div className="text-gray-600 font-medium">Revenue Generated</div>
                </div>
              </div>
            </div>
            
            {/* Right Content - Demo Preview */}
            <div className="relative animate-scale-in" style={{animationDelay: '0.4s'}}>
              <div className="relative bg-white rounded-3xl shadow-2xl p-6 transform rotate-3 hover:rotate-0 transition-transform duration-500">
                <div className="aspect-video bg-gradient-to-br from-blue-100 to-purple-100 rounded-2xl mb-4 overflow-hidden relative">
                  <div className="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-purple-500/10"></div>
                  <div className="h-full flex items-center justify-center relative z-10">
                    <div className="text-center p-8">
                      <div className="w-20 h-20 bg-gradient-to-r from-blue-600 to-purple-600 rounded-3xl mx-auto mb-4 flex items-center justify-center shadow-lg">
                        <ShoppingCart className="w-10 h-10 text-white" />
                      </div>
                      <h3 className="text-2xl font-bold text-gray-900 mb-2">Your Store Preview</h3>
                      <p className="text-gray-600">Professional ecommerce in minutes</p>
                      <div className="mt-4 flex justify-center gap-2">
                        <div className="w-12 h-2 bg-gray-300 rounded-full"></div>
                        <div className="w-8 h-2 bg-gray-200 rounded-full"></div>
                        <div className="w-16 h-2 bg-gray-200 rounded-full"></div>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="flex gap-2 justify-center">
                  <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                  <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                  <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                </div>
              </div>
              
              {/* Floating Elements */}
              <div className="absolute -top-6 -right-6 bg-white rounded-2xl shadow-lg p-4 animate-bounce">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
                  <span className="text-sm font-medium text-gray-900">Live Store</span>
                </div>
              </div>
              
              <div className="absolute -bottom-6 -left-6 bg-white rounded-2xl shadow-lg p-4">
                <div className="flex items-center gap-3">
                  <BarChart3 className="w-6 h-6 text-blue-600" />
                  <div>
                    <div className="text-lg font-bold text-gray-900">↗ 347%</div>
                    <div className="text-xs text-gray-600">Sales Growth</div>
                  </div>
                </div>
              </div>

              <div className="absolute top-1/2 -left-8 bg-white rounded-2xl shadow-lg p-3">
                <div className="flex items-center gap-2">
                  <Building className="w-5 h-5 text-purple-600" />
                  <span className="text-sm font-medium text-gray-900">SaaS</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Overview */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-6">
              Everything You Need to 
              <span className="bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent"> Succeed Online</span>
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              Built with hexagonal architecture for maximum flexibility, scalability, and performance. 
              Launch faster, customize deeper, scale infinitely.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {[
              {
                icon: Palette,
                title: 'Industry-Specific Themes',
                description: 'Choose from 20+ professional themes designed for fashion, electronics, marketplace, B2B, and more.',
                color: 'from-purple-500 to-pink-500'
              },
              {
                icon: Settings,
                title: 'No-Code Customization',
                description: 'Customize colors, layouts, and features with our intuitive visual editor. No coding required.',
                color: 'from-blue-500 to-cyan-500'
              },
              {
                icon: Code,
                title: 'Hexagonal Architecture',
                description: 'Enterprise-grade architecture that adapts to your business needs and scales infinitely.',
                color: 'from-green-500 to-emerald-500'
              },
              {
                icon: Smartphone,
                title: 'Mobile-First Design',
                description: 'All themes are fully responsive and optimized for mobile commerce with PWA support.',
                color: 'from-orange-500 to-red-500'
              },
              {
                icon: Shield,
                title: 'Enterprise Security',
                description: 'PCI DSS compliant with advanced security, fraud protection, and SSL certificates.',
                color: 'from-indigo-500 to-purple-500'
              },
              {
                icon: Globe,
                title: 'Global Commerce',
                description: 'Multi-currency, multi-language, international shipping, and local payment methods.',
                color: 'from-teal-500 to-blue-500'
              }
            ].map((feature, index) => (
              <div 
                key={index}
                className="group bg-white rounded-2xl p-8 shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-2 border border-gray-100"
              >
                <div className={`w-14 h-14 bg-gradient-to-r ${feature.color} rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg`}>
                  <feature.icon className="w-7 h-7 text-white" />
                </div>
                <h3 className="text-xl font-bold text-gray-900 mb-3">{feature.title}</h3>
                <p className="text-gray-600 leading-relaxed mb-4">{feature.description}</p>
                <div className="flex items-center text-blue-600 group-hover:text-purple-600 transition-colors cursor-pointer">
                  <span className="text-sm font-semibold">Learn more</span>
                  <ChevronRight className="w-4 h-4 ml-1 group-hover:translate-x-1 transition-transform" />
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Theme Showcase */}
      <section className="py-20 bg-gradient-to-br from-gray-50 to-blue-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-6">
              Themes Built for Every Industry
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              From fashion boutiques to B2B marketplaces, we have specialized themes 
              that understand your business and convert visitors into customers.
            </p>
          </div>
          
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {demos.map((demo, index) => (
              <div 
                key={index}
                className={`group relative bg-white rounded-3xl overflow-hidden shadow-xl hover:shadow-2xl transition-all duration-500 transform hover:-translate-y-2 cursor-pointer ${
                  activeDemo === index ? 'ring-4 ring-blue-500 ring-opacity-50' : ''
                }`}
                onClick={() => setActiveDemo(index)}
              >
                <div className="aspect-video bg-gray-200 relative overflow-hidden">
                  <div className={`absolute inset-0 bg-gradient-to-br ${demo.color} opacity-80 flex items-center justify-center`}>
                    <div className="text-center p-8">
                      <div className="w-16 h-16 bg-white/20 backdrop-blur-sm rounded-2xl mx-auto mb-4 flex items-center justify-center border border-white/30">
                        <span className="text-2xl font-bold text-white">{demo.name.charAt(0)}</span>
                      </div>
                      <h3 className="text-xl font-bold text-white mb-2">{demo.name}</h3>
                      <div className="space-y-2">
                        <div className="h-2 bg-white/30 rounded-full w-full"></div>
                        <div className="h-2 bg-white/20 rounded-full w-3/4 mx-auto"></div>
                        <div className="h-2 bg-white/20 rounded-full w-1/2 mx-auto"></div>
                      </div>
                    </div>
                  </div>
                  <div className="absolute top-4 right-4 bg-white/90 backdrop-blur-sm rounded-full p-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <Play className="w-4 h-4 text-gray-700" />
                  </div>
                </div>
                
                <div className="p-6">
                  <h3 className="text-xl font-bold text-gray-900 mb-3">{demo.name}</h3>
                  <p className="text-gray-600 mb-4">Perfect for {demo.name.toLowerCase()} businesses with specialized features</p>
                  <div className="flex flex-wrap gap-2">
                    {demo.features.map((feature, fIndex) => (
                      <span 
                        key={fIndex}
                        className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium"
                      >
                        {feature}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
          
          <div className="text-center mt-12">
            <Link
              href="/themes"
              className="inline-flex items-center gap-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white px-8 py-4 rounded-full hover:from-blue-700 hover:to-purple-700 transition-all duration-300 font-semibold shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              Explore All 20+ Themes
              <ArrowRight className="w-5 h-5" />
            </Link>
          </div>
        </div>
      </section>

      {/* Social Proof */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-6">
              Loved by Entrepreneurs Worldwide
            </h2>
            <p className="text-xl text-gray-600">
              Join thousands of successful businesses already using our platform
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {[
              {
                name: "Sarah Johnson",
                role: "Fashion Boutique Owner",
                content: "I launched my fashion store in 30 minutes using the Fashion theme. The size charts and color swatches are perfect for my customers. Sales increased 250% in the first month!",
                rating: 5,
                avatar: "SJ",
                company: "Bella Fashion"
              },
              {
                name: "Marcus Rodriguez", 
                role: "Electronics Retailer",
                content: "The spec comparison features and technical review system saved me months of development. My customers love comparing products side-by-side. Best decision I made for my business.",
                rating: 5,
                avatar: "MR",
                company: "TechHub Pro"
              },
              {
                name: "Lisa Chen",
                role: "Marketplace Founder", 
                content: "Managing 200+ vendors was impossible before. The multi-vendor dashboard and automated commission system changed everything. My marketplace now runs itself!",
                rating: 5,
                avatar: "LC",
                company: "Global Marketplace"
              }
            ].map((testimonial, index) => (
              <div key={index} className="bg-gradient-to-br from-gray-50 to-white rounded-2xl p-8 shadow-lg hover:shadow-xl transition-shadow duration-300 border border-gray-100">
                <div className="flex items-center mb-4">
                  {[...Array(testimonial.rating)].map((_, i) => (
                    <Star key={i} className="w-5 h-5 text-yellow-400 fill-current" />
                  ))}
                </div>
                <p className="text-gray-700 mb-6 italic leading-relaxed">"{testimonial.content}"</p>
                <div className="flex items-center">
                  <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white font-bold text-lg">
                    {testimonial.avatar}
                  </div>
                  <div className="ml-4">
                    <div className="font-bold text-gray-900">{testimonial.name}</div>
                    <div className="text-gray-600 text-sm">{testimonial.role}</div>
                    <div className="text-blue-600 text-sm font-medium">{testimonial.company}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Architecture Highlight */}
      <section className="py-20 bg-gradient-to-r from-gray-900 to-black text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <div>
              <div className="inline-flex items-center gap-2 bg-blue-600/20 text-blue-300 px-4 py-2 rounded-full text-sm font-semibold mb-6">
                <Layers className="w-4 h-4" />
                Enterprise Architecture
              </div>
              
              <h2 className="text-4xl font-bold mb-6">
                Built with 
                <span className="bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent"> Hexagonal Architecture</span>
              </h2>
              
              <p className="text-xl text-gray-300 mb-8 leading-relaxed">
                Unlike traditional ecommerce platforms, we use hexagonal architecture 
                to ensure your store is flexible, maintainable, and scales with your business.
              </p>
              
              <div className="space-y-4">
                {[
                  'Clean separation of business logic and external dependencies',
                  'Easy integration with any payment gateway, shipping provider, or ERP',
                  'Modular design that adapts to your specific business needs',
                  'Enterprise-grade scalability and performance'
                ].map((benefit, index) => (
                  <div key={index} className="flex items-start gap-3">
                    <Check className="w-6 h-6 text-green-400 flex-shrink-0 mt-0.5" />
                    <span className="text-gray-300">{benefit}</span>
                  </div>
                ))}
              </div>
            </div>
            
            <div className="relative">
              <div className="bg-gradient-to-br from-blue-900/50 to-purple-900/50 rounded-3xl p-8 border border-blue-500/20">
                <div className="text-center">
                  <div className="w-20 h-20 bg-gradient-to-r from-blue-500 to-purple-500 rounded-2xl mx-auto mb-6 flex items-center justify-center">
                    <Code className="w-10 h-10 text-white" />
                  </div>
                  <h3 className="text-2xl font-bold text-white mb-4">Hexagonal Core</h3>
                  <div className="space-y-3">
                    <div className="bg-blue-600/30 rounded-lg p-3 border border-blue-500/30">
                      <span className="text-blue-200 font-medium">Domain Layer</span>
                    </div>
                    <div className="bg-purple-600/30 rounded-lg p-3 border border-purple-500/30">
                      <span className="text-purple-200 font-medium">Application Layer</span>
                    </div>
                    <div className="bg-green-600/30 rounded-lg p-3 border border-green-500/30">
                      <span className="text-green-200 font-medium">Infrastructure Layer</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-r from-blue-600 to-purple-600">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-4xl lg:text-5xl font-bold text-white mb-6">
            Ready to Launch Your Ecommerce Empire?
          </h2>
          <p className="text-xl text-blue-100 mb-8 leading-relaxed">
            Join thousands of successful entrepreneurs who chose our platform to build their online business. 
            Start with any theme, customize to match your brand, and launch in minutes.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mb-8">
            <Link
              href="/get-started"
              className="inline-flex items-center gap-3 bg-white text-blue-600 px-10 py-5 rounded-full hover:bg-gray-100 transition-all duration-300 transform hover:scale-105 shadow-2xl font-bold text-lg"
            >
              Start Your Free Trial
              <ArrowRight className="w-6 h-6" />
            </Link>
            
            <div className="text-white/80 text-sm">
              No credit card required • 14-day free trial • Cancel anytime
            </div>
          </div>
          
          <div className="flex flex-wrap justify-center items-center gap-6 text-white/80">
            <div className="flex items-center gap-2">
              <Check className="w-5 h-5" />
              <span>Free SSL Certificate</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="w-5 h-5" />
              <span>24/7 Premium Support</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="w-5 h-5" />
              <span>No Setup Fees</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="w-5 h-5" />
              <span>Money-Back Guarantee</span>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}