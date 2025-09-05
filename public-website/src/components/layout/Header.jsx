'use client';

import { useState } from 'react';
import Link from 'next/link';
import { 
  Menu, 
  X, 
  ChevronDown, 
  Layers,
  ArrowRight,
  Play
} from 'lucide-react';

export default function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [showSolutions, setShowSolutions] = useState(false);
  const [showResources, setShowResources] = useState(false);

  const solutions = [
    {
      name: 'Fashion & Apparel',
      description: 'Size charts, color swatches, lookbooks',
      href: '/themes/fashion',
      icon: '👗'
    },
    {
      name: 'Electronics & Tech',
      description: 'Spec comparison, tech reviews',
      href: '/themes/electronics',
      icon: '📱'
    },
    {
      name: 'Marketplace',
      description: 'Multi-vendor, commission tracking',
      href: '/themes/marketplace',
      icon: '🏪'
    },
    {
      name: 'B2B Wholesale',
      description: 'Bulk pricing, account management',
      href: '/themes/b2b',
      icon: '🏢'
    }
  ];

  const resources = [
    {
      name: 'Documentation',
      description: 'Get started guides and API docs',
      href: '/docs',
      icon: '📚'
    },
    {
      name: 'Theme Gallery',
      description: 'Browse all available themes',
      href: '/themes',
      icon: '🎨'
    },
    {
      name: 'Success Stories',
      description: 'Customer case studies',
      href: '/customers',
      icon: '⭐'
    },
    {
      name: 'Blog',
      description: 'Ecommerce tips and updates',
      href: '/blog',
      icon: '📝'
    }
  ];

  return (
    <header className="bg-white/95 backdrop-blur-lg border-b border-gray-200 sticky top-0 z-50">
      {/* Main Navigation */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-20">
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
                <Layers className="w-6 h-6 text-white" />
              </div>
              <div className="flex flex-col">
                <span className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                  Hexagonal
                </span>
                <span className="text-xs text-gray-500 -mt-1">Ecommerce SaaS</span>
              </div>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden lg:flex items-center space-x-8">
            {/* Solutions Dropdown */}
            <div 
              className="relative"
              onMouseEnter={() => setShowSolutions(true)}
              onMouseLeave={() => setShowSolutions(false)}
            >
              <button className="flex items-center text-gray-700 hover:text-blue-600 transition-all duration-200 font-medium group">
                Solutions
                <ChevronDown className={`ml-1 w-4 h-4 transition-transform duration-200 ${showSolutions ? 'rotate-180' : ''}`} />
                <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-blue-600 to-purple-600 transition-all duration-200 group-hover:w-full"></span>
              </button>
              
              {showSolutions && (
                <div className="absolute top-full left-0 mt-2 w-80 bg-white rounded-2xl shadow-2xl border border-gray-100 py-4 animate-fade-in">
                  <div className="px-4 pb-2 mb-2 border-b border-gray-100">
                    <h3 className="text-sm font-semibold text-gray-900">Industry Solutions</h3>
                  </div>
                  {solutions.map((solution) => (
                    <Link
                      key={solution.name}
                      href={solution.href}
                      className="flex items-start gap-3 px-4 py-3 hover:bg-gradient-to-r hover:from-blue-50 hover:to-purple-50 transition-all duration-200 rounded-lg mx-2"
                    >
                      <span className="text-2xl">{solution.icon}</span>
                      <div>
                        <div className="font-medium text-gray-900">{solution.name}</div>
                        <div className="text-sm text-gray-600">{solution.description}</div>
                      </div>
                    </Link>
                  ))}
                  <div className="px-4 pt-2 mt-2 border-t border-gray-100">
                    <Link
                      href="/themes"
                      className="inline-flex items-center gap-2 text-blue-600 hover:text-purple-600 font-medium text-sm"
                    >
                      View all themes
                      <ArrowRight className="w-4 h-4" />
                    </Link>
                  </div>
                </div>
              )}
            </div>

            <Link href="/pricing" className="text-gray-700 hover:text-blue-600 transition-all duration-200 font-medium relative group">
              Pricing
              <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-blue-600 to-purple-600 transition-all duration-200 group-hover:w-full"></span>
            </Link>

            {/* Resources Dropdown */}
            <div 
              className="relative"
              onMouseEnter={() => setShowResources(true)}
              onMouseLeave={() => setShowResources(false)}
            >
              <button className="flex items-center text-gray-700 hover:text-blue-600 transition-all duration-200 font-medium group">
                Resources
                <ChevronDown className={`ml-1 w-4 h-4 transition-transform duration-200 ${showResources ? 'rotate-180' : ''}`} />
                <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-blue-600 to-purple-600 transition-all duration-200 group-hover:w-full"></span>
              </button>
              
              {showResources && (
                <div className="absolute top-full right-0 mt-2 w-72 bg-white rounded-2xl shadow-2xl border border-gray-100 py-4 animate-fade-in">
                  {resources.map((resource) => (
                    <Link
                      key={resource.name}
                      href={resource.href}
                      className="flex items-start gap-3 px-4 py-3 hover:bg-gradient-to-r hover:from-blue-50 hover:to-purple-50 transition-all duration-200 rounded-lg mx-2"
                    >
                      <span className="text-xl">{resource.icon}</span>
                      <div>
                        <div className="font-medium text-gray-900">{resource.name}</div>
                        <div className="text-sm text-gray-600">{resource.description}</div>
                      </div>
                    </Link>
                  ))}
                </div>
              )}
            </div>

            <Link href="/about" className="text-gray-700 hover:text-blue-600 transition-all duration-200 font-medium relative group">
              About
              <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-blue-600 to-purple-600 transition-all duration-200 group-hover:w-full"></span>
            </Link>
          </nav>

          {/* Action Buttons */}
          <div className="flex items-center space-x-4">
            <Link
              href="/login"
              className="hidden sm:block text-gray-700 hover:text-blue-600 transition-colors font-medium"
            >
              Login
            </Link>
            
            <button className="group hidden sm:flex items-center gap-2 border border-gray-300 text-gray-700 px-4 py-2 rounded-full hover:border-blue-300 hover:text-blue-600 hover:bg-blue-50 transition-all duration-300">
              <Play className="w-4 h-4" />
              <span className="font-medium">Demo</span>
            </button>
            
            <Link
              href="/get-started"
              className="inline-flex items-center gap-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white px-6 py-3 rounded-full hover:from-blue-700 hover:to-purple-700 transition-all duration-300 font-semibold shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              Start Free Trial
              <ArrowRight className="w-4 h-4" />
            </Link>

            {/* Mobile menu button */}
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="lg:hidden p-2 text-gray-700 hover:text-blue-600 hover:bg-blue-50 transition-all duration-200 rounded-xl"
            >
              {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Menu */}
      {mobileMenuOpen && (
        <div className="lg:hidden border-t border-gray-200 bg-white/95 backdrop-blur-lg">
          <div className="px-4 py-6 space-y-1">
            <div className="space-y-1">
              <div className="px-3 py-2 text-sm font-semibold text-gray-900">Solutions</div>
              <div className="pl-6 space-y-1">
                {solutions.map((solution) => (
                  <Link
                    key={solution.name}
                    href={solution.href}
                    className="flex items-center gap-3 px-3 py-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 transition-all duration-200 rounded-lg"
                    onClick={() => setMobileMenuOpen(false)}
                  >
                    <span>{solution.icon}</span>
                    <span className="text-sm">{solution.name}</span>
                  </Link>
                ))}
              </div>
            </div>
            
            <Link
              href="/pricing"
              className="block px-3 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-all duration-200 rounded-lg font-medium"
              onClick={() => setMobileMenuOpen(false)}
            >
              Pricing
            </Link>

            <div className="space-y-1">
              <div className="px-3 py-2 text-sm font-semibold text-gray-900">Resources</div>
              <div className="pl-6 space-y-1">
                {resources.map((resource) => (
                  <Link
                    key={resource.name}
                    href={resource.href}
                    className="flex items-center gap-3 px-3 py-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 transition-all duration-200 rounded-lg"
                    onClick={() => setMobileMenuOpen(false)}
                  >
                    <span>{resource.icon}</span>
                    <span className="text-sm">{resource.name}</span>
                  </Link>
                ))}
              </div>
            </div>
            
            <Link
              href="/about"
              className="block px-3 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-all duration-200 rounded-lg font-medium"
              onClick={() => setMobileMenuOpen(false)}
            >
              About
            </Link>
            
            <div className="border-t border-gray-200 pt-4 mt-4 space-y-3">
              <Link
                href="/login"
                className="block px-3 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-all duration-200 rounded-lg font-medium"
                onClick={() => setMobileMenuOpen(false)}
              >
                Login
              </Link>
              
              <button 
                className="flex items-center gap-2 w-full px-3 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-all duration-200 rounded-lg font-medium"
                onClick={() => setMobileMenuOpen(false)}
              >
                <Play className="w-4 h-4" />
                Watch Demo
              </button>
              
              <Link
                href="/get-started"
                className="block w-full text-center bg-gradient-to-r from-blue-600 to-purple-600 text-white px-4 py-3 rounded-lg font-semibold hover:from-blue-700 hover:to-purple-700 transition-all duration-300"
                onClick={() => setMobileMenuOpen(false)}
              >
                Start Free Trial
              </Link>
            </div>
          </div>
        </div>
      )}
    </header>
  );
}