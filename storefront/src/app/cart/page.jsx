'use client';

import { useAtom } from 'jotai';
import Image from 'next/image';
import Link from 'next/link';
import { Minus, Plus, Trash2, ShoppingBag, ArrowLeft } from 'lucide-react';
import { cartItemsAtom, updateQuantityAtom, removeFromCartAtom, clearCartAtom, cartTotalAtom, cartCountAtom } from '@/store/cartStore';

export default function CartPage() {
  const [cartItems] = useAtom(cartItemsAtom);
  const [, updateQuantity] = useAtom(updateQuantityAtom);
  const [, removeFromCart] = useAtom(removeFromCartAtom);
  const [, clearCart] = useAtom(clearCartAtom);
  const [cartTotal] = useAtom(cartTotalAtom);
  const [cartCount] = useAtom(cartCountAtom);

  const handleQuantityChange = (itemId, newQuantity) => {
    if (newQuantity <= 0) {
      removeFromCart(itemId);
    } else {
      updateQuantity(itemId, newQuantity);
    }
  };

  const handleRemoveItem = (itemId) => {
    removeFromCart(itemId);
  };

  const handleClearCart = () => {
    if (window.confirm('Are you sure you want to clear your cart?')) {
      clearCart();
    }
  };

  if (cartItems.length === 0) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-50 to-blue-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
          <div className="text-center bg-white rounded-3xl p-16 shadow-2xl max-w-2xl mx-auto">
            <div className="relative mb-8">
              <div className="w-32 h-32 mx-auto bg-gradient-to-br from-blue-100 to-purple-100 rounded-full flex items-center justify-center">
                <ShoppingBag className="h-16 w-16 text-blue-500" />
              </div>
              <div className="absolute -top-2 -right-2 w-8 h-8 bg-red-500 rounded-full flex items-center justify-center">
                <span className="text-white text-sm font-bold">0</span>
              </div>
            </div>
            <h1 className="text-4xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent mb-4">
              Your cart is empty
            </h1>
            <p className="text-xl text-gray-600 mb-10 leading-relaxed">
              Looks like you haven't added anything to your cart yet. Discover our amazing products!
            </p>
            <div className="space-y-4">
              <Link
                href="/products"
                className="inline-flex items-center gap-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white px-8 py-4 rounded-full hover:from-blue-700 hover:to-purple-700 transition-all duration-300 transform hover:scale-105 shadow-lg font-semibold text-lg"
              >
                <ArrowLeft className="h-5 w-5" />
                Start Shopping
              </Link>
              <div className="flex justify-center space-x-6 text-sm text-gray-500">
                <span className="flex items-center gap-1">
                  <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                  Free shipping over $50
                </span>
                <span className="flex items-center gap-1">
                  <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                  30-day returns
                </span>
                <span className="flex items-center gap-1">
                  <div className="w-2 h-2 bg-purple-500 rounded-full"></div>
                  Secure payments
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-blue-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-10">
          <div className="text-center mb-8">
            <h1 className="text-4xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent mb-4">
              Shopping Cart
            </h1>
            <p className="text-xl text-gray-600">
              {cartCount} {cartCount === 1 ? 'item' : 'items'} ready for checkout
            </p>
            <div className="w-24 h-1 bg-gradient-to-r from-blue-500 to-purple-500 mx-auto mt-4 rounded-full"></div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Cart Items */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-2xl shadow-xl border border-white/20 backdrop-blur-sm">
              {/* Cart Header */}
              <div className="px-8 py-6 border-b border-gray-100 flex justify-between items-center">
                <h2 className="text-2xl font-bold text-gray-900">Your Items</h2>
                <button
                  onClick={handleClearCart}
                  className="text-sm font-medium text-red-600 hover:text-red-700 hover:bg-red-50 px-4 py-2 rounded-full transition-all duration-200"
                >
                  Clear Cart
                </button>
              </div>

              {/* Cart Items List */}
              <div className="divide-y divide-gray-200">
                {cartItems.map((item) => {
                  const itemTotal = item.variant.price * item.quantity;
                  
                  return (
                    <div key={item.id} className="p-8 hover:bg-gray-50/50 transition-colors duration-200">
                      <div className="flex items-start gap-6">
                        {/* Product Image */}
                        <div className="relative w-28 h-28 bg-gray-100 rounded-2xl overflow-hidden flex-shrink-0 shadow-lg">
                          <Image
                            src={item.product.images[0]}
                            alt={item.product.name}
                            fill
                            className="object-cover"
                          />
                        </div>

                        {/* Product Details */}
                        <div className="flex-1 min-w-0">
                          <div className="flex justify-between items-start">
                            <div>
                              <h3 className="text-lg font-medium text-gray-900 mb-1">
                                <Link
                                  href={`/product/${item.product.slug}`}
                                  className="hover:text-blue-600 transition-colors duration-200"
                                >
                                  {item.product.name}
                                </Link>
                              </h3>
                              <p className="text-sm text-gray-500 mb-2">{item.product.category}</p>
                              
                              {/* Variant Details */}
                              <div className="flex flex-wrap gap-4 text-sm text-gray-600">
                                {item.variant.size && (
                                  <span>Size: <span className="font-medium">{item.variant.size}</span></span>
                                )}
                                {item.variant.color && (
                                  <span>Color: <span className="font-medium">{item.variant.color}</span></span>
                                )}
                                {item.variant.material && (
                                  <span>Material: <span className="font-medium">{item.variant.material}</span></span>
                                )}
                              </div>
                            </div>
                            
                            {/* Remove Button */}
                            <button
                              onClick={() => handleRemoveItem(item.id)}
                              className="p-3 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-full transition-all duration-200"
                              title="Remove item"
                            >
                              <Trash2 className="h-5 w-5" />
                            </button>
                          </div>

                          {/* Price and Quantity */}
                          <div className="flex items-center justify-between mt-4">
                            <div className="flex items-center gap-4">
                              {/* Quantity Controls */}
                              <div className="flex items-center border-2 border-gray-200 rounded-xl overflow-hidden shadow-sm">
                                <button
                                  onClick={() => handleQuantityChange(item.id, item.quantity - 1)}
                                  className="p-3 hover:bg-blue-50 hover:text-blue-600 transition-all duration-200 disabled:opacity-50"
                                  disabled={item.quantity <= 1}
                                >
                                  <Minus className="h-4 w-4" />
                                </button>
                                <span className="px-6 py-3 bg-gray-50 border-x-2 border-gray-200 min-w-[4rem] text-center font-semibold text-lg">
                                  {item.quantity}
                                </span>
                                <button
                                  onClick={() => handleQuantityChange(item.id, item.quantity + 1)}
                                  className="p-3 hover:bg-blue-50 hover:text-blue-600 transition-all duration-200 disabled:opacity-50"
                                  disabled={item.quantity >= item.variant.stock}
                                >
                                  <Plus className="h-4 w-4" />
                                </button>
                              </div>
                              
                              <span className="text-sm text-gray-500 font-medium">
                                Stock: {item.variant.stock}
                              </span>
                            </div>

                            {/* Price */}
                            <div className="text-right">
                              <div className="text-lg font-medium text-gray-900">
                                ${itemTotal.toFixed(2)}
                              </div>
                              <div className="text-sm text-gray-500">
                                ${item.variant.price.toFixed(2)} each
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>

          {/* Order Summary */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-2xl shadow-xl border border-white/20 backdrop-blur-sm p-8 sticky top-8">
              <h2 className="text-2xl font-bold text-gray-900 mb-8">Order Summary</h2>
              
              {/* Summary Details */}
              <div className="space-y-4 mb-6">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Subtotal ({cartCount} items)</span>
                  <span className="text-gray-900">${cartTotal.toFixed(2)}</span>
                </div>
                
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Shipping</span>
                  <span className="text-gray-900">
                    {cartTotal >= 50 ? 'Free' : '$9.99'}
                  </span>
                </div>
                
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Tax</span>
                  <span className="text-gray-900">${(cartTotal * 0.08).toFixed(2)}</span>
                </div>
                
                <div className="border-t border-gray-200 pt-4">
                  <div className="flex justify-between text-lg font-medium">
                    <span className="text-gray-900">Total</span>
                    <span className="text-gray-900">
                      ${(cartTotal + (cartTotal >= 50 ? 0 : 9.99) + (cartTotal * 0.08)).toFixed(2)}
                    </span>
                  </div>
                </div>
              </div>

              {/* Free Shipping Notice */}
              {cartTotal < 50 && (
                <div className="bg-blue-50 border border-blue-200 rounded-md p-4 mb-6">
                  <p className="text-sm text-blue-700">
                    Add <span className="font-medium">${(50 - cartTotal).toFixed(2)}</span> more to get free shipping!
                  </p>
                </div>
              )}

              {/* Action Buttons */}
              <div className="space-y-4">
                <Link
                  href="/checkout"
                  className="w-full bg-gradient-to-r from-blue-600 to-purple-600 text-white py-4 px-6 rounded-xl hover:from-blue-700 hover:to-purple-700 transition-all duration-300 text-center block font-semibold text-lg shadow-lg hover:shadow-xl transform hover:scale-[1.02]"
                >
                  Proceed to Checkout
                </Link>
                
                <Link
                  href="/products"
                  className="w-full border-2 border-gray-200 text-gray-700 py-4 px-6 rounded-xl hover:border-blue-300 hover:bg-blue-50 hover:text-blue-700 transition-all duration-300 text-center block font-semibold flex items-center justify-center gap-2"
                >
                  <ArrowLeft className="h-5 w-5" />
                  Continue Shopping
                </Link>
              </div>

              {/* Security Notice */}
              <div className="mt-6 pt-6 border-t border-gray-200">
                <div className="flex items-center gap-2 text-sm text-gray-500">
                  <svg className="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clipRule="evenodd" />
                  </svg>
                  <span>Secure checkout with SSL encryption</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}