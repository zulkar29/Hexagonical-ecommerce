'use client';

import { useAtom } from 'jotai';
import Link from 'next/link';
import { Heart, ShoppingCart, BarChart3, Star } from 'lucide-react';
import { addToCartAtom } from '@/store/cartStore';
import { wishlistAtom, addToWishlistAtom, removeFromWishlistAtom } from '@/store/userStore';
import { comparisonAtom, addToComparisonAtom, removeFromComparisonAtom } from '@/store/userStore';

export default function ProductListItem({ product }) {
  const [, addToCart] = useAtom(addToCartAtom);
  const [wishlist] = useAtom(wishlistAtom);
  const [, addToWishlist] = useAtom(addToWishlistAtom);
  const [, removeFromWishlist] = useAtom(removeFromWishlistAtom);
  const [comparison] = useAtom(comparisonAtom);
  const [, addToComparison] = useAtom(addToComparisonAtom);
  const [, removeFromComparison] = useAtom(removeFromComparisonAtom);

  const isInWishlist = wishlist.includes(product.id);
  const isInComparison = comparison.includes(product.id);

  const handleAddToCart = () => {
    addToCart({
      id: product.id,
      name: product.name,
      price: product.price,
      image: product.images[0],
      quantity: 1,
      variant: {
        size: product.variants?.sizes?.[0] || null,
        color: product.variants?.colors?.[0] || null
      }
    });
  };

  const handleWishlistToggle = () => {
    if (isInWishlist) {
      removeFromWishlist(product.id);
    } else {
      addToWishlist(product.id);
    }
  };

  const handleComparisonToggle = () => {
    if (isInComparison) {
      removeFromComparison(product.id);
    } else {
      if (comparison.length >= 4) {
        alert('You can only compare up to 4 products at a time.');
        return;
      }
      addToComparison(product.id);
    }
  };

  const discountPercentage = product.originalPrice && product.originalPrice > product.price
    ? Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100)
    : 0;

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow duration-200">
      <div className="flex flex-col sm:flex-row">
        {/* Product Image */}
        <div className="relative w-full sm:w-48 h-48 sm:h-auto flex-shrink-0">
          <Link href={`/product/${product.slug}`}>
            <img
              src={product.images[0]}
              alt={product.name}
              className="w-full h-full object-cover hover:scale-105 transition-transform duration-300"
            />
          </Link>
          
          {/* Discount Badge */}
          {discountPercentage > 0 && (
            <div className="absolute top-3 left-3 bg-red-500 text-white px-2 py-1 rounded-md text-xs font-medium">
              {discountPercentage}% OFF
            </div>
          )}
          
          {/* Stock Status */}
          {product.stock === 0 && (
            <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center">
              <span className="bg-white text-gray-900 px-3 py-1 rounded-md text-sm font-medium">
                Out of Stock
              </span>
            </div>
          )}
        </div>
        
        {/* Product Info */}
        <div className="flex-1 p-6">
          <div className="flex flex-col h-full">
            {/* Category */}
            <div className="mb-2">
              <span className="text-xs text-gray-500 uppercase tracking-wide">
                {product.category}
              </span>
            </div>
            
            {/* Product Name */}
            <Link href={`/product/${product.slug}`}>
              <h3 className="text-xl font-semibold text-gray-900 mb-2 hover:text-gray-700 transition-colors duration-200">
                {product.name}
              </h3>
            </Link>
            
            {/* Description */}
            <p className="text-gray-600 text-sm mb-4 line-clamp-2">
              {product.description}
            </p>
            
            {/* Rating */}
            <div className="flex items-center gap-1 mb-4">
              <div className="flex items-center">
                {[...Array(5)].map((_, i) => (
                  <Star
                    key={i}
                    className={`h-4 w-4 ${
                      i < Math.floor(product.rating) ? 'text-yellow-400 fill-current' : 'text-gray-300'
                    }`}
                  />
                ))}
              </div>
              <span className="text-sm text-gray-600 ml-1">
                {product.rating} ({product.reviews} reviews)
              </span>
            </div>
            
            {/* Features */}
            {product.features && product.features.length > 0 && (
              <div className="mb-4">
                <div className="flex flex-wrap gap-1">
                  {product.features.slice(0, 3).map((feature, index) => (
                    <span
                      key={index}
                      className="inline-flex items-center px-2 py-1 rounded-md bg-gray-100 text-gray-700 text-xs"
                    >
                      {feature}
                    </span>
                  ))}
                  {product.features.length > 3 && (
                    <span className="text-xs text-gray-500">+{product.features.length - 3} more</span>
                  )}
                </div>
              </div>
            )}
            
            {/* Available Colors */}
            {product.variants?.colors && (
              <div className="flex items-center gap-2 mb-4">
                <span className="text-sm text-gray-600">Colors:</span>
                <div className="flex items-center gap-1">
                  {product.variants.colors.slice(0, 5).map((color, index) => (
                    <div
                      key={index}
                      className="w-5 h-5 rounded-full border border-gray-300"
                      style={{ backgroundColor: color.toLowerCase() }}
                      title={color}
                    />
                  ))}
                  {product.variants.colors.length > 5 && (
                    <span className="text-xs text-gray-500 ml-1">+{product.variants.colors.length - 5}</span>
                  )}
                </div>
              </div>
            )}
            
            {/* Available Sizes */}
            {product.variants?.sizes && (
              <div className="flex items-center gap-2 mb-4">
                <span className="text-sm text-gray-600">Sizes:</span>
                <div className="flex items-center gap-1">
                  {product.variants.sizes.slice(0, 4).map((size, index) => (
                    <span
                      key={index}
                      className="inline-flex items-center px-2 py-1 rounded-md bg-gray-100 text-gray-700 text-xs"
                    >
                      {size}
                    </span>
                  ))}
                  {product.variants.sizes.length > 4 && (
                    <span className="text-xs text-gray-500">+{product.variants.sizes.length - 4}</span>
                  )}
                </div>
              </div>
            )}
            
            {/* Price and Actions */}
            <div className="mt-auto">
              <div className="flex items-center justify-between">
                {/* Price */}
                <div className="flex items-center gap-2">
                  <span className="text-2xl font-bold text-gray-900">
                    ${product.price}
                  </span>
                  {product.originalPrice && product.originalPrice > product.price && (
                    <span className="text-lg text-gray-500 line-through">
                      ${product.originalPrice}
                    </span>
                  )}
                </div>
                
                {/* Action Buttons */}
                <div className="flex items-center gap-2">
                  {/* Wishlist Button */}
                  <button
                    onClick={handleWishlistToggle}
                    className={`p-2 rounded-full border transition-colors duration-200 ${
                      isInWishlist
                        ? 'bg-red-50 border-red-200 text-red-600 hover:bg-red-100'
                        : 'bg-white border-gray-300 text-gray-600 hover:bg-gray-50'
                    }`}
                    title={isInWishlist ? 'Remove from wishlist' : 'Add to wishlist'}
                  >
                    <Heart className={`h-4 w-4 ${isInWishlist ? 'fill-current' : ''}`} />
                  </button>
                  
                  {/* Comparison Button */}
                  <button
                    onClick={handleComparisonToggle}
                    className={`p-2 rounded-full border transition-colors duration-200 ${
                      isInComparison
                        ? 'bg-blue-50 border-blue-200 text-blue-600 hover:bg-blue-100'
                        : 'bg-white border-gray-300 text-gray-600 hover:bg-gray-50'
                    }`}
                    title={isInComparison ? 'Remove from comparison' : 'Add to comparison'}
                  >
                    <BarChart3 className={`h-4 w-4 ${isInComparison ? 'fill-current' : ''}`} />
                  </button>
                  
                  {/* Add to Cart Button */}
                  <button
                    onClick={handleAddToCart}
                    disabled={product.stock === 0}
                    className="flex items-center gap-2 bg-black text-white px-6 py-2 rounded-md font-medium hover:bg-gray-800 transition-colors duration-200 disabled:bg-gray-300 disabled:cursor-not-allowed"
                  >
                    <ShoppingCart className="h-4 w-4" />
                    {product.stock === 0 ? 'Out of Stock' : 'Add to Cart'}
                  </button>
                </div>
              </div>
              
              {/* Stock Info */}
              <div className="flex items-center justify-between mt-3">
                <span className={`text-sm px-2 py-1 rounded-full ${
                  product.stock > 0
                    ? product.stock <= 5
                      ? 'bg-yellow-100 text-yellow-700'
                      : 'bg-green-100 text-green-700'
                    : 'bg-red-100 text-red-700'
                }`}>
                  {product.stock > 0
                    ? product.stock <= 5
                      ? `Only ${product.stock} left`
                      : 'In Stock'
                    : 'Out of Stock'
                  }
                </span>
                
                <Link
                  href={`/product/${product.slug}`}
                  className="text-sm text-blue-600 hover:text-blue-700 font-medium"
                >
                  View Details →
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}