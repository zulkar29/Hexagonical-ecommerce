'use client';

import { Fragment } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAtom } from 'jotai';
import { X, Search, ShoppingCart, User, Home, Grid3X3, Tag } from 'lucide-react';
import { cartCountAtom } from '@/store/cartStore';

export default function MobileMenu({ isOpen, onClose }) {
  const [cartCount] = useAtom(cartCountAtom);

  const navigationItems = [
    { name: 'All Products', href: '/products' },
    { name: 'Clothing', href: '/categories/clothing' },
    { name: 'Electronics', href: '/categories/electronics' },
    { name: 'Home & Garden', href: '/categories/home-garden' },
  ];

  return (
    <Transition.Root show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={onClose}>
        <Transition.Child
          as={Fragment}
          enter="ease-in-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in-out duration-300"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-hidden">
          <div className="absolute inset-0 overflow-hidden">
            <div className="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10">
              <Transition.Child
                as={Fragment}
                enter="transform transition ease-in-out duration-300"
                enterFrom="translate-x-full"
                enterTo="translate-x-0"
                leave="transform transition ease-in-out duration-300"
                leaveFrom="translate-x-0"
                leaveTo="translate-x-full"
              >
                <Dialog.Panel className="pointer-events-auto w-screen max-w-sm">
                  <div className="flex h-full flex-col overflow-y-scroll bg-white shadow-xl">
                    {/* Header */}
                    <div className="flex items-center justify-between px-4 py-6 border-b border-gray-200">
                      <h2 className="text-lg font-medium text-gray-900">Menu</h2>
                      <button
                        type="button"
                        className="p-2 -mr-2 hover:bg-gray-100 rounded-lg transition-colors"
                        onClick={onClose}
                      >
                        <X className="h-6 w-6 text-gray-400" />
                      </button>
                    </div>

                    {/* Search */}
                    <div className="px-4 py-4 border-b border-gray-200">
                      <div className="flex items-center bg-gray-100 rounded-lg">
                        <Search className="h-5 w-5 text-gray-400 ml-3" />
                        <input
                          type="text"
                          placeholder="Search products..."
                          className="w-full px-3 py-2 bg-transparent border-none outline-none text-gray-700 placeholder-gray-400"
                        />
                      </div>
                    </div>

                    {/* Navigation */}
                    <div className="flex-1 px-4 py-6">
                      <nav className="space-y-1">
                        {navigationItems.map((item) => (
                          <Link
                            key={item.name}
                            href={item.href}
                            className="block px-3 py-2 text-base font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900 rounded-lg transition-colors"
                            onClick={onClose}
                          >
                            {item.name}
                          </Link>
                        ))}
                      </nav>
                    </div>

                    {/* Bottom Actions */}
                    <div className="border-t border-gray-200 px-4 py-6">
                      <div className="space-y-3">
                        <Link
                          href="/cart"
                          className="flex items-center justify-between w-full px-3 py-2 text-base font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900 rounded-lg transition-colors"
                          onClick={onClose}
                        >
                          <div className="flex items-center">
                            <ShoppingCart className="h-5 w-5 mr-3" />
                            Cart
                          </div>
                          {cartCount > 0 && (
                            <span className="bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center font-medium">
                              {cartCount}
                            </span>
                          )}
                        </Link>
                        
                        <Link
                          href="/account"
                          className="flex items-center w-full px-3 py-2 text-base font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900 rounded-lg transition-colors"
                          onClick={onClose}
                        >
                          <User className="h-5 w-5 mr-3" />
                          Account
                        </Link>
                      </div>
                    </div>
                  </div>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </div>
      </Dialog>
    </Transition.Root>
  );
}