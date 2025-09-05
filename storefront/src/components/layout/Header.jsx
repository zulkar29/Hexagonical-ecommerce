'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useAtom } from 'jotai';
import { 
  Search, 
  ShoppingCart, 
  Heart, 
  User, 
  Menu, 
  X, 
  GitCompare,
  ChevronDown,
  Truck
} from 'lucide-react';
import { 
  cartCountAtom, 
  wishlistCountAtom, 
  comparisonCountAtom,
  mobileMenuOpenAtom,
  searchOpenAtom
} from '@/store/atoms';
import SearchBar from '@/components/search/SearchBar';

export default function Header() {
  const [cartCount] = useAtom(cartCountAtom);
  const [wishlistCount] = useAtom(wishlistCountAtom);
  const [comparisonCount] = useAtom(comparisonCountAtom);
  const [mobileMenuOpen, setMobileMenuOpen] = useAtom(mobileMenuOpenAtom);
  const [searchOpen, setSearchOpen] = useAtom(searchOpenAtom);
  const [showCategories, setShowCategories] = useState(false);

  const categories = [
    { name: 'Electronics', href: '/category/electronics' },
    { name: 'Clothing', href: '/category/clothing' },
    { name: 'Accessories', href: '/category/accessories' },
    { name: 'Home & Kitchen', href: '/category/home-kitchen' }
  ];

  return (
    <header className="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-50">
      {/* Top bar */}
      <div className="bg-gradient-to-r from-gray-900 to-gray-800 text-white text-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-3">
            <div className="flex items-center space-x-6">
              <span className="flex items-center gap-2">
                <Truck className="w-4 h-4" />
                Free shipping on orders over $50
              </span>
            </div>
            <div className="flex items-center space-x-6">
              <Link href="/help" className="hover:text-blue-300 transition-colors duration-200 font-medium">
                Help
              </Link>
              <Link href="/contact" className="hover:text-blue-300 transition-colors duration-200 font-medium">
                Contact
              </Link>
              <div className="border-l border-gray-600 pl-4">
                <span className="text-xs opacity-90">24/7 Support Available</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main header */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-20">
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
                <span className="text-white font-bold text-xl">S</span>
              </div>
              <div className="flex flex-col">
                <span className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">Storefront</span>
                <span className="text-xs text-gray-500 -mt-1">Premium Shopping</span>
              </div>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-8">
            <Link href="/" className="text-gray-700 hover:text-blue-600 transition-all duration-200 font-medium relative group">
              Home
              <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-blue-600 to-purple-600 transition-all duration-200 group-hover:w-full"></span>
            </Link>
            
            <div 
              className="relative"
              onMouseEnter={() => setShowCategories(true)}
              onMouseLeave={() => setShowCategories(false)}
            >
              <button className="flex items-center text-gray-700 hover:text-blue-600 transition-all duration-200 font-medium relative group">
                Categories
                <ChevronDown className={`ml-1 w-4 h-4 transition-transform duration-200 ${showCategories ? 'rotate-180' : ''}`} />
                <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-blue-600 to-purple-600 transition-all duration-200 group-hover:w-full"></span>
              </button>
              
              {showCategories && (
                <div className="absolute top-full left-0 mt-2 w-56 bg-white rounded-2xl shadow-2xl border border-gray-100 py-3 z-50 animate-fade-in">
                  {categories.map((category) => (
                    <Link
                      key={category.name}
                      href={category.href}
                      className="block px-6 py-3 text-gray-700 hover:bg-gradient-to-r hover:from-blue-50 hover:to-purple-50 hover:text-blue-600 transition-all duration-200 rounded-lg mx-2"
                    >
                      {category.name}
                    </Link>
                  ))}
                </div>
              )}
            </div>
            
            <Link href="/products" className="text-gray-700 hover:text-blue-600 transition-all duration-200 font-medium relative group">
              Products
              <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-blue-600 to-purple-600 transition-all duration-200 group-hover:w-full"></span>
            </Link>
            <Link href="/deals" className="text-gray-700 hover:text-blue-600 transition-all duration-200 font-medium relative group">
              Deals
              <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-blue-600 to-purple-600 transition-all duration-200 group-hover:w-full"></span>
            </Link>
            <Link href="/about" className="text-gray-700 hover:text-blue-600 transition-all duration-200 font-medium relative group">
              About
              <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-gradient-to-r from-blue-600 to-purple-600 transition-all duration-200 group-hover:w-full"></span>
            </Link>
          </nav>

          {/* Search Bar - Desktop */}
          <div className="hidden lg:block flex-1 max-w-lg mx-8">
            <SearchBar />
          </div>

          {/* Action buttons */}
          <div className="flex items-center space-x-2">
            {/* Search button - Mobile */}
            <button
              onClick={() => setSearchOpen(!searchOpen)}
              className="lg:hidden p-3 text-gray-700 hover:text-blue-600 hover:bg-blue-50 transition-all duration-200 rounded-xl"
            >
              <Search className="w-5 h-5" />
            </button>

            {/* Comparison */}
            <Link
              href="/compare"
              className="relative p-3 text-gray-700 hover:text-blue-600 hover:bg-blue-50 transition-all duration-200 rounded-xl group"
            >
              <GitCompare className="w-5 h-5 group-hover:scale-110 transition-transform duration-200" />
              {comparisonCount > 0 && (
                <span className="absolute -top-1 -right-1 bg-gradient-to-r from-blue-500 to-blue-600 text-white text-xs rounded-full w-6 h-6 flex items-center justify-center font-semibold shadow-lg animate-pulse">
                  {comparisonCount}
                </span>
              )}
            </Link>

            {/* Wishlist */}
            <Link
              href="/wishlist"
              className="relative p-3 text-gray-700 hover:text-red-500 hover:bg-red-50 transition-all duration-200 rounded-xl group"
            >
              <Heart className="w-5 h-5 group-hover:scale-110 transition-transform duration-200" />
              {wishlistCount > 0 && (
                <span className="absolute -top-1 -right-1 bg-gradient-to-r from-red-500 to-pink-500 text-white text-xs rounded-full w-6 h-6 flex items-center justify-center font-semibold shadow-lg animate-pulse">
                  {wishlistCount}
                </span>
              )}
            </Link>

            {/* Cart */}
            <Link
              href="/cart"
              className="relative p-3 text-gray-700 hover:text-blue-600 hover:bg-blue-50 transition-all duration-200 rounded-xl group"
            >
              <ShoppingCart className="w-5 h-5 group-hover:scale-110 transition-transform duration-200" />
              {cartCount > 0 && (
                <span className="absolute -top-1 -right-1 bg-gradient-to-r from-blue-500 to-purple-500 text-white text-xs rounded-full w-6 h-6 flex items-center justify-center font-semibold shadow-lg animate-bounce">
                  {cartCount}
                </span>
              )}
            </Link>

            {/* User account */}
            <Link
              href="/account"
              className="p-3 text-gray-700 hover:text-blue-600 hover:bg-blue-50 transition-all duration-200 rounded-xl group"
            >
              <User className="w-5 h-5 group-hover:scale-110 transition-transform duration-200" />
            </Link>

            {/* Mobile menu button */}
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="md:hidden p-3 text-gray-700 hover:text-blue-600 hover:bg-blue-50 transition-all duration-200 rounded-xl"
            >
              {mobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Search */}
      {searchOpen && (
        <div className="lg:hidden border-t border-gray-200 p-4">
          <SearchBar />
        </div>
      )}

      {/* Mobile Menu */}
      {mobileMenuOpen && (
        <div className="md:hidden border-t border-gray-200">
          <div className="px-4 py-2 space-y-1">
            <Link
              href="/"
              className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors rounded-md"
              onClick={() => setMobileMenuOpen(false)}
            >
              Home
            </Link>
            
            <div className="px-3 py-2">
              <div className="text-gray-700 font-medium mb-2">Categories</div>
              <div className="pl-4 space-y-1">
                {categories.map((category) => (
                  <Link
                    key={category.name}
                    href={category.href}
                    className="block py-1 text-gray-600 hover:text-blue-600 transition-colors"
                    onClick={() => setMobileMenuOpen(false)}
                  >
                    {category.name}
                  </Link>
                ))}
              </div>
            </div>
            
            <Link
              href="/products"
              className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors rounded-md"
              onClick={() => setMobileMenuOpen(false)}
            >
              Products
            </Link>
            <Link
              href="/deals"
              className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors rounded-md"
              onClick={() => setMobileMenuOpen(false)}
            >
              Deals
            </Link>
            <Link
              href="/about"
              className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors rounded-md"
              onClick={() => setMobileMenuOpen(false)}
            >
              About
            </Link>
            
            <div className="border-t border-gray-200 pt-2 mt-2">
              <Link
                href="/account"
                className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors rounded-md"
                onClick={() => setMobileMenuOpen(false)}
              >
                My Account
              </Link>
              <Link
                href="/help"
                className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors rounded-md"
                onClick={() => setMobileMenuOpen(false)}
              >
                Help
              </Link>
              <Link
                href="/contact"
                className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors rounded-md"
                onClick={() => setMobileMenuOpen(false)}
              >
                Contact
              </Link>
            </div>
          </div>
        </div>
      )}
    </header>
  );
}