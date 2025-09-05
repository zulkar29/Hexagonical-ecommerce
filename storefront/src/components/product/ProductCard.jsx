'use client';

import { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useAtom } from 'jotai';
import { Heart, ShoppingCart, Star, Eye, GitCompare } from 'lucide-react';
import { cartAtom, wishlistAtom, comparisonAtom } from '@/store/atoms';

export default function ProductCard({ product, viewMode = 'grid' }) {
  const [cart, setCart] = useAtom(cartAtom);
  const [wishlist, setWishlist] = useAtom(wishlistAtom);
  const [comparison, setComparison] = useAtom(comparisonAtom);
  const [imageLoading, setImageLoading] = useState(true);
  const [selectedVariant, setSelectedVariant] = useState({
    color: product.variants?.colors?.[0] || null,
    size: product.variants?.sizes?.[0] || null
  });

  const isInWishlist = wishlist.some(item => item.id === product.id);
  const isInComparison = comparison.some(item => item.id === product.id);
  const isInCart = cart.some(item => item.id === product.id);
  
  const discountPercentage = product.originalPrice 
    ? Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100)
    : 0;

  const handleAddToCart = (e) => {
    e.preventDefault();
    e.stopPropagation();
    
    const cartItem = {
      id: product.id,
      name: product.name,
      price: product.price,
      image: product.images[0],
      variant: selectedVariant,
      quantity: 1
    };
    
    const existingItem = cart.find(item => 
      item.id === product.id && 
      JSON.stringify(item.variant) === JSON.stringify(selectedVariant)
    );
    
    if (existingItem) {
      setCart(cart.map(item => 
        item.id === product.id && JSON.stringify(item.variant) === JSON.stringify(selectedVariant)
          ? { ...item, quantity: item.quantity + 1 }
          : item
      ));
    } else {
      setCart([...cart, cartItem]);
    }
  };

  const handleToggleWishlist = (e) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (isInWishlist) {
      setWishlist(wishlist.filter(item => item.id !== product.id));
    } else {
      setWishlist([...wishlist, product]);
    }
  };

  const handleToggleComparison = (e) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (isInComparison) {
      setComparison(comparison.filter(item => item.id !== product.id));
    } else if (comparison.length < 4) {
      setComparison([...comparison, product]);
    }
  };

  const renderRating = () => {
    const stars = [];
    const fullStars = Math.floor(product.rating);
    const hasHalfStar = product.rating % 1 !== 0;
    
    for (let i = 0; i < fullStars; i++) {
      stars.push(
        <Star key={i} className="w-4 h-4 fill-yellow-400 text-yellow-400" />
      );
    }
    
    if (hasHalfStar) {
      stars.push(
        <Star key="half" className="w-4 h-4 fill-yellow-400/50 text-yellow-400" />
      );
    }
    
    const emptyStars = 5 - Math.ceil(product.rating);
    for (let i = 0; i < emptyStars; i++) {
      stars.push(
        <Star key={`empty-${i}`} className="w-4 h-4 text-gray-300" />
      );
    }
    
    return stars;
  };

  if (viewMode === 'list') {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 hover:shadow-md transition-shadow">
        <div className="flex gap-4">
          <div className="relative w-32 h-32 flex-shrink-0">
            <Link href={`/products/${product.slug}`}>
              <Image
                src={product.images[0]}
                alt={product.name}
                fill
                className="object-cover rounded-lg"
                onLoad={() => setImageLoading(false)}
              />
            </Link>
            {discountPercentage > 0 && (
              <span className="absolute top-2 left-2 bg-red-500 text-white text-xs px-2 py-1 rounded">
                -{discountPercentage}%
              </span>
            )}
            {product.stock < 10 && product.stock > 0 && (
              <span className="absolute top-2 right-2 bg-orange-500 text-white text-xs px-2 py-1 rounded">
                Low Stock
              </span>
            )}
            {product.stock === 0 && (
              <span className="absolute top-2 right-2 bg-gray-500 text-white text-xs px-2 py-1 rounded">
                Out of Stock
              </span>
            )}
          </div>
          
          <div className="flex-1">
            <div className="flex justify-between items-start mb-2">
              <div>
                <p className="text-sm text-gray-500 mb-1">{product.category}</p>
                <Link href={`/products/${product.slug}`}>
                  <h3 className="font-semibold text-lg hover:text-blue-600 transition-colors">
                    {product.name}
                  </h3>
                </Link>
              </div>
              <button
                onClick={handleToggleWishlist}
                className={`p-2 rounded-full transition-colors ${
                  isInWishlist 
                    ? 'text-red-500 bg-red-50' 
                    : 'text-gray-400 hover:text-red-500 hover:bg-red-50'
                }`}
              >
                <Heart className={`w-5 h-5 ${isInWishlist ? 'fill-current' : ''}`} />
              </button>
            </div>
            
            <p className="text-gray-600 text-sm mb-3 line-clamp-2">
              {product.description}
            </p>
            
            <div className="flex items-center gap-2 mb-3">
              <div className="flex items-center">
                {renderRating()}
              </div>
              <span className="text-sm text-gray-500">({product.reviews})</span>
            </div>
            
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-xl font-bold text-gray-900">
                  ${product.price}
                </span>
                {product.originalPrice && (
                  <span className="text-sm text-gray-500 line-through">
                    ${product.originalPrice}
                  </span>
                )}
              </div>
              
              <div className="flex items-center gap-2">
                <button
                  onClick={handleToggleComparison}
                  disabled={!isInComparison && comparison.length >= 4}
                  className={`p-2 rounded-full transition-colors ${
                    isInComparison 
                      ? 'text-blue-500 bg-blue-50' 
                      : 'text-gray-400 hover:text-blue-500 hover:bg-blue-50 disabled:opacity-50 disabled:cursor-not-allowed'
                  }`}
                >
                  <GitCompare className="w-4 h-4" />
                </button>
                
                <button
                  onClick={handleAddToCart}
                  disabled={product.stock === 0}
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <ShoppingCart className="w-4 h-4" />
                  {isInCart ? 'Added' : 'Add to Cart'}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-2xl transition-all duration-500 group transform hover:-translate-y-2">
      <div className="relative aspect-square overflow-hidden bg-gray-50">
        <Link href={`/products/${product.slug}`}>
          <Image
            src={product.images[0]}
            alt={product.name}
            fill
            className="object-cover group-hover:scale-110 transition-transform duration-700"
            onLoad={() => setImageLoading(false)}
          />
        </Link>
        
        {/* Badges */}
        <div className="absolute top-4 left-4 flex flex-col gap-2">
          {discountPercentage > 0 && (
            <span className="bg-gradient-to-r from-red-500 to-pink-500 text-white text-xs px-3 py-1.5 rounded-full font-semibold shadow-lg animate-pulse">
              -{discountPercentage}%
            </span>
          )}
          {product.stock < 10 && product.stock > 0 && (
            <span className="bg-gradient-to-r from-orange-500 to-amber-500 text-white text-xs px-3 py-1.5 rounded-full font-semibold shadow-lg">
              Low Stock
            </span>
          )}
          {product.stock === 0 && (
            <span className="bg-gradient-to-r from-gray-500 to-gray-600 text-white text-xs px-3 py-1.5 rounded-full font-semibold shadow-lg">
              Out of Stock
            </span>
          )}
        </div>
        
        {/* Action buttons */}
        <div className="absolute top-4 right-4 flex flex-col gap-2 opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-x-2 group-hover:translate-x-0">
          <button
            onClick={handleToggleWishlist}
            className={`p-3 rounded-full backdrop-blur-md shadow-lg transition-all duration-300 transform hover:scale-110 ${
              isInWishlist 
                ? 'text-red-500 bg-white/95 shadow-red-200' 
                : 'text-gray-600 bg-white/95 hover:text-red-500 hover:shadow-red-100'
            }`}
          >
            <Heart className={`w-5 h-5 ${isInWishlist ? 'fill-current' : ''}`} />
          </button>
          
          <button
            onClick={handleToggleComparison}
            disabled={!isInComparison && comparison.length >= 4}
            className={`p-3 rounded-full backdrop-blur-md shadow-lg transition-all duration-300 transform hover:scale-110 ${
              isInComparison 
                ? 'text-blue-500 bg-white/95 shadow-blue-200' 
                : 'text-gray-600 bg-white/95 hover:text-blue-500 hover:shadow-blue-100 disabled:opacity-50 disabled:cursor-not-allowed'
            }`}
          >
            <GitCompare className="w-5 h-5" />
          </button>
          
          <Link
            href={`/products/${product.slug}`}
            className="p-3 rounded-full bg-white/95 text-gray-600 hover:text-blue-500 transition-all duration-300 transform hover:scale-110 shadow-lg hover:shadow-blue-100"
          >
            <Eye className="w-5 h-5" />
          </Link>
        </div>
        
        {imageLoading && (
          <div className="absolute inset-0 bg-gray-200 animate-pulse" />
        )}
      </div>
      
      <div className="p-6">
        <div className="mb-3">
          <p className="text-xs font-medium text-blue-600 uppercase tracking-wider mb-2">{product.category}</p>
          <Link href={`/products/${product.slug}`}>
            <h3 className="font-bold text-gray-900 hover:text-blue-600 transition-colors line-clamp-2 text-lg leading-tight">
              {product.name}
            </h3>
          </Link>
        </div>
        
        <div className="flex items-center gap-2 mb-4">
          <div className="flex items-center">
            {renderRating()}
          </div>
          <span className="text-sm text-gray-500 font-medium">({product.reviews})</span>
        </div>
        
        {/* Variant selection */}
        {product.variants?.colors && product.variants.colors.length > 1 && (
          <div className="mb-3">
            <div className="flex gap-1">
              {product.variants.colors.slice(0, 4).map((color) => (
                <button
                  key={color}
                  onClick={(e) => {
                    e.preventDefault();
                    setSelectedVariant(prev => ({ ...prev, color }));
                  }}
                  className={`w-6 h-6 rounded-full border-2 transition-all ${
                    selectedVariant.color === color 
                      ? 'border-gray-900 scale-110' 
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                  style={{ 
                    backgroundColor: color.toLowerCase() === 'white' ? '#ffffff' : 
                                   color.toLowerCase() === 'black' ? '#000000' :
                                   color.toLowerCase() === 'gray' ? '#6b7280' :
                                   color.toLowerCase() === 'navy' ? '#1e3a8a' :
                                   color.toLowerCase() === 'brown' ? '#92400e' :
                                   color.toLowerCase() === 'tan' ? '#d2b48c' :
                                   color.toLowerCase() === 'silver' ? '#c0c0c0' :
                                   color.toLowerCase() === 'blue' ? '#3b82f6' :
                                   color.toLowerCase() === 'green' ? '#10b981' :
                                   color.toLowerCase() === 'rose gold' ? '#e8b4b8' :
                                   color.toLowerCase()
                  }}
                  title={color}
                />
              ))}
              {product.variants.colors.length > 4 && (
                <span className="text-xs text-gray-500 self-center ml-1">
                  +{product.variants.colors.length - 4}
                </span>
              )}
            </div>
          </div>
        )}
        
        <div className="flex items-center justify-between">
          <div className="flex flex-col">
            <div className="flex items-center gap-2">
              <span className="text-xl font-bold text-gray-900">
                ${product.price}
              </span>
              {product.originalPrice && (
                <span className="text-sm text-gray-500 line-through">
                  ${product.originalPrice}
                </span>
              )}
            </div>
            {discountPercentage > 0 && (
              <span className="text-xs text-green-600 font-medium">
                Save ${(product.originalPrice - product.price).toFixed(2)}
              </span>
            )}
          </div>
          
          <button
            onClick={handleAddToCart}
            disabled={product.stock === 0}
            className="bg-gradient-to-r from-blue-600 to-blue-700 text-white px-4 py-3 rounded-xl hover:from-blue-700 hover:to-blue-800 transition-all duration-300 disabled:from-gray-400 disabled:to-gray-500 disabled:cursor-not-allowed flex items-center gap-2 text-sm font-semibold shadow-lg hover:shadow-xl transform hover:scale-105"
          >
            <ShoppingCart className="w-4 h-4" />
            {isInCart ? 'Added' : 'Add to Cart'}
          </button>
        </div>
      </div>
    </div>
  );
}