'use client';

import { useState } from 'react';
import Link from 'next/link';
import { 
  Menu, 
  X, 
  ShoppingCart,
  ArrowRight
} from 'lucide-react';
import LanguageSwitcher from '../ui/LanguageSwitcher';
import { useTranslations } from '@/hooks/useTranslations';

export default function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const { t } = useTranslations();

  return (
    <header className="bg-white/95 backdrop-blur-lg border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gray-900 rounded-lg flex items-center justify-center">
              <ShoppingCart className="w-5 h-5 text-white" />
            </div>
            <span className="text-xl font-bold text-gray-900">StoreBuilder</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden lg:flex items-center space-x-6">
            <Link href="/" className="text-gray-700 hover:text-gray-900 transition-colors font-medium">
              {t('navigation.home')}
            </Link>
            <Link href="/templates" className="text-gray-700 hover:text-gray-900 transition-colors font-medium">
              {t('themes.title')}
            </Link>
            <Link href="/pricing" className="text-gray-700 hover:text-gray-900 transition-colors font-medium">
              {t('navigation.pricing')}
            </Link>
            <Link href="/contact" className="text-gray-700 hover:text-gray-900 transition-colors font-medium">
              {t('navigation.contact')}
            </Link>
          </nav>

          {/* Action Buttons */}
          <div className="flex items-center space-x-3">
            <LanguageSwitcher />
            
            <Link
              href="/get-started"
              className="inline-flex items-center gap-2 bg-gray-900 text-white px-4 py-2 rounded-lg hover:bg-gray-800 transition-all duration-300 font-medium shadow-sm hover:shadow-md"
            >
              {t('navigation.getStarted')}
              <ArrowRight className="w-4 h-4" />
            </Link>

            {/* Mobile menu button */}
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="lg:hidden p-2 text-gray-700 hover:text-gray-900 hover:bg-gray-50 transition-all duration-200 rounded-lg"
            >
              {mobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Menu */}
      {mobileMenuOpen && (
        <div className="lg:hidden border-t border-gray-200 bg-white/95 backdrop-blur-lg">
          <div className="px-4 py-4 space-y-2">
            <Link
              href="/"
              className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-gray-900 transition-all duration-200 rounded-lg font-medium"
              onClick={() => setMobileMenuOpen(false)}
            >
              {t('navigation.home')}
            </Link>
            <Link
              href="/templates"
              className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-gray-900 transition-all duration-200 rounded-lg font-medium"
              onClick={() => setMobileMenuOpen(false)}
            >
              {t('themes.title')}
            </Link>
            <Link
              href="/pricing"
              className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-gray-900 transition-all duration-200 rounded-lg font-medium"
              onClick={() => setMobileMenuOpen(false)}
            >
              {t('navigation.pricing')}
            </Link>
            <Link
              href="/contact"
              className="block px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-gray-900 transition-all duration-200 rounded-lg font-medium"
              onClick={() => setMobileMenuOpen(false)}
            >
              {t('navigation.contact')}
            </Link>
            <div className="pt-2">
              <Link
                href="/get-started"
                className="block w-full text-center bg-gray-900 text-white px-4 py-3 rounded-lg font-medium hover:bg-gray-800 transition-all duration-300"
                onClick={() => setMobileMenuOpen(false)}
              >
                {t('navigation.getStarted')}
              </Link>
            </div>
          </div>
        </div>
      )}
    </header>
  );
}