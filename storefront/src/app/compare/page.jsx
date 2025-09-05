'use client';

import { useAtom } from 'jotai';
import Link from 'next/link';
import { X, ShoppingCart, Heart, Star, Check, Minus } from 'lucide-react';
import { comparisonAtom, removeFromComparisonAtom, clearComparisonAtom } from '@/store/userStore';
import { addToCartAtom } from '@/store/cartStore';
import { addToWishlistAtom } from '@/store/userStore';
import { getProducts } from '@/data/products';

export default function ComparePage() {
  const [comparison] = useAtom(comparisonAtom);
  const [, removeFromComparison] = useAtom(removeFromComparisonAtom);
  const [, clearComparison] = useAtom(clearComparisonAtom);
  const [, addToCart] = useAtom(addToCartAtom);
  const [, addToWishlist] = useAtom(addToWishlistAtom);

  // Get full product details for comparison items
  const products = getProducts();
  const comparisonProducts = products.filter(product => comparison.includes(product.id));

  const handleAddToCart = (product) => {
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

  const handleAddToWishlist = (productId) => {
    addToWishlist(productId);
  };

  // Comparison features to display
  const comparisonFeatures = [
    { key: 'price', label: 'Price', type: 'price' },
    { key: 'rating', label: 'Rating', type: 'rating' },
    { key: 'reviews', label: 'Reviews', type: 'number' },
    { key: 'category', label: 'Category', type: 'text' },
    { key: 'stock', label: 'Stock', type: 'stock' },
    { key: 'material', label: 'Material', type: 'array' },
    { key: 'sizes', label: 'Available Sizes', type: 'variants' },
    { key: 'colors', label: 'Available Colors', type: 'variants' },
    { key: 'features', label: 'Features', type: 'array' }
  ];

  const renderFeatureValue = (product, feature) => {
    const value = feature.key === 'sizes' || feature.key === 'colors' 
      ? product.variants?.[feature.key]
      : product[feature.key];

    switch (feature.type) {
      case 'price':
        return (
          <div className="flex items-center gap-2">
            <span className="text-lg font-bold text-gray-900">${value}</span>
            {product.originalPrice && product.originalPrice > value && (
              <span className="text-sm text-gray-500 line-through">${product.originalPrice}</span>
            )}
          </div>
        );
      
      case 'rating':
        return (
          <div className="flex items-center gap-1">
            {[...Array(5)].map((_, i) => (
              <Star
                key={i}
                className={`h-4 w-4 ${
                  i < Math.floor(value) ? 'text-yellow-400 fill-current' : 'text-gray-300'
                }`}
              />
            ))}
            <span className="text-sm text-gray-600 ml-1">({value})</span>
          </div>
        );
      
      case 'number':
        return <span className="text-gray-900">{value?.toLocaleString()}</span>;
      
      case 'stock':
        return (
          <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
            value > 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
          }`}>
            {value > 0 ? `${value} in stock` : 'Out of stock'}
          </span>
        );
      
      case 'array':
        return value && value.length > 0 ? (
          <div className="flex flex-wrap gap-1">
            {value.map((item, index) => (
              <span key={index} className="inline-flex items-center px-2 py-1 rounded-md bg-gray-100 text-gray-700 text-xs">
                {item}
              </span>
            ))}
          </div>
        ) : (
          <span className="text-gray-400">-</span>
        );
      
      case 'variants':
        return value && value.length > 0 ? (
          <div className="flex flex-wrap gap-1">
            {feature.key === 'colors' ? (
              value.map((color, index) => (
                <div
                  key={index}
                  className="w-6 h-6 rounded-full border border-gray-300"
                  style={{ backgroundColor: color.toLowerCase() }}
                  title={color}
                />
              ))
            ) : (
              value.map((size, index) => (
                <span key={index} className="inline-flex items-center px-2 py-1 rounded-md bg-gray-100 text-gray-700 text-xs">
                  {size}
                </span>
              ))
            )}
          </div>
        ) : (
          <span className="text-gray-400">-</span>
        );
      
      default:
        return <span className="text-gray-900">{value || '-'}</span>;
    }
  };

  if (comparisonProducts.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="text-center">
            <div className="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-6">
              <div className="flex items-center gap-1">
                <div className="w-6 h-6 bg-gray-400 rounded"></div>
                <div className="w-6 h-6 bg-gray-400 rounded"></div>
              </div>
            </div>
            <h1 className="text-3xl font-bold text-gray-900 mb-4">No Products to Compare</h1>
            <p className="text-lg text-gray-600 mb-8 max-w-md mx-auto">
              Add products to compare their features, prices, and specifications side by side.
            </p>
            <Link
              href="/products"
              className="inline-flex items-center gap-2 bg-black text-white px-6 py-3 rounded-md font-medium hover:bg-gray-800 transition-colors duration-200"
            >
              Browse Products
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Product Comparison</h1>
            <p className="text-gray-600 mt-1">
              Comparing {comparisonProducts.length} {comparisonProducts.length === 1 ? 'product' : 'products'}
            </p>
          </div>
          {comparisonProducts.length > 0 && (
            <button
              onClick={() => clearComparison()}
              className="px-4 py-2 text-red-600 hover:text-red-700 font-medium border border-red-300 rounded-md hover:bg-red-50 transition-colors duration-200"
            >
              Clear All
            </button>
          )}
        </div>

        {/* Comparison Table */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              {/* Product Headers */}
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left p-6 bg-gray-50 font-medium text-gray-900 w-48">
                    Products
                  </th>
                  {comparisonProducts.map((product) => (
                    <th key={product.id} className="p-6 bg-gray-50 min-w-80">
                      <div className="text-center">
                        {/* Product Image */}
                        <div className="relative mb-4">
                          <img
                            src={product.images[0]}
                            alt={product.name}
                            className="w-32 h-32 object-cover rounded-lg mx-auto"
                          />
                          <button
                            onClick={() => removeFromComparison(product.id)}
                            className="absolute -top-2 -right-2 p-1 bg-red-500 text-white rounded-full hover:bg-red-600 transition-colors duration-200"
                          >
                            <X className="h-4 w-4" />
                          </button>
                        </div>
                        
                        {/* Product Name */}
                        <Link href={`/product/${product.slug}`}>
                          <h3 className="font-semibold text-gray-900 hover:text-gray-700 transition-colors duration-200 mb-2">
                            {product.name}
                          </h3>
                        </Link>
                        
                        {/* Action Buttons */}
                        <div className="flex flex-col gap-2">
                          <button
                            onClick={() => handleAddToCart(product)}
                            disabled={product.stock === 0}
                            className="flex items-center justify-center gap-2 bg-black text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-800 transition-colors duration-200 disabled:bg-gray-300 disabled:cursor-not-allowed"
                          >
                            <ShoppingCart className="h-4 w-4" />
                            Add to Cart
                          </button>
                          <button
                            onClick={() => handleAddToWishlist(product.id)}
                            className="flex items-center justify-center gap-2 border border-gray-300 text-gray-700 px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-50 transition-colors duration-200"
                          >
                            <Heart className="h-4 w-4" />
                            Add to Wishlist
                          </button>
                        </div>
                      </div>
                    </th>
                  ))}
                </tr>
              </thead>
              
              {/* Feature Comparison */}
              <tbody>
                {comparisonFeatures.map((feature, index) => (
                  <tr key={feature.key} className={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>
                    <td className="p-6 font-medium text-gray-900 border-r border-gray-200">
                      {feature.label}
                    </td>
                    {comparisonProducts.map((product) => (
                      <td key={product.id} className="p-6 text-center">
                        {renderFeatureValue(product, feature)}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Add More Products */}
        <div className="text-center mt-8">
          <Link
            href="/products"
            className="inline-flex items-center gap-2 text-blue-600 hover:text-blue-700 font-medium"
          >
            Add More Products to Compare
          </Link>
        </div>

        {/* Comparison Tips */}
        <div className="mt-16 bg-white rounded-lg p-6 border border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Comparison Tips</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="text-center">
              <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-3">
                <Check className="h-6 w-6 text-blue-600" />
              </div>
              <h3 className="font-medium text-gray-900 mb-2">Compare Features</h3>
              <p className="text-sm text-gray-600">
                See detailed specifications and features side by side.
              </p>
            </div>
            
            <div className="text-center">
              <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-3">
                <Star className="h-6 w-6 text-green-600" />
              </div>
              <h3 className="font-medium text-gray-900 mb-2">Check Ratings</h3>
              <p className="text-sm text-gray-600">
                Compare customer ratings and reviews to make informed decisions.
              </p>
            </div>
            
            <div className="text-center">
              <div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-3">
                <ShoppingCart className="h-6 w-6 text-purple-600" />
              </div>
              <h3 className="font-medium text-gray-900 mb-2">Quick Purchase</h3>
              <p className="text-sm text-gray-600">
                Add your preferred products directly to cart from comparison.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}