'use client';

import { useState, useMemo } from 'react';
import { useParams } from 'next/navigation';
import Image from 'next/image';
import { useAtom } from 'jotai';
import { Star, Heart, ShoppingCart, Truck, Shield, RotateCcw, Share2, ChevronLeft, ChevronRight } from 'lucide-react';
import { getProductBySlug } from '@/data/products';
import { addToCartAtom } from '@/store/cartStore';
import { addToWishlistAtom, removeFromWishlistAtom, wishlistAtom, addToRecentlyViewedAtom } from '@/store/userStore';
import { addToComparisonAtom } from '@/store/searchStore';
import ProductCard from '@/components/product/ProductCard';

export default function ProductDetailPage() {
  const params = useParams();
  const product = getProductBySlug(params.slug);
  
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);
  const [selectedVariant, setSelectedVariant] = useState(product?.variants[0] || null);
  const [quantity, setQuantity] = useState(1);
  const [activeTab, setActiveTab] = useState('description');
  
  const [, addToCart] = useAtom(addToCartAtom);
  const [, addToWishlist] = useAtom(addToWishlistAtom);
  const [, removeFromWishlist] = useAtom(removeFromWishlistAtom);
  const [, addToComparison] = useAtom(addToComparisonAtom);
  const [, addToRecentlyViewed] = useAtom(addToRecentlyViewedAtom);
  const [wishlist] = useAtom(wishlistAtom);

  // Add to recently viewed when component mounts
  useState(() => {
    if (product) {
      addToRecentlyViewed(product);
    }
  }, [product, addToRecentlyViewed]);

  if (!product) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Product Not Found</h1>
          <p className="text-gray-600">The product you're looking for doesn't exist.</p>
        </div>
      </div>
    );
  }

  const isInWishlist = wishlist.some(item => item.id === product.id);
  const hasDiscount = selectedVariant.originalPrice > selectedVariant.price;
  const discountPercentage = hasDiscount 
    ? Math.round(((selectedVariant.originalPrice - selectedVariant.price) / selectedVariant.originalPrice) * 100)
    : 0;

  // Group variants by attribute for selection
  const variantOptions = useMemo(() => {
    const sizes = [...new Set(product.variants.map(v => v.size))];
    const colors = [...new Set(product.variants.map(v => v.color))];
    const materials = [...new Set(product.variants.map(v => v.material))];
    
    return { sizes, colors, materials };
  }, [product.variants]);

  const handleVariantChange = (attribute, value) => {
    const newVariant = product.variants.find(variant => {
      const currentSelection = {
        size: attribute === 'size' ? value : selectedVariant.size,
        color: attribute === 'color' ? value : selectedVariant.color,
        material: attribute === 'material' ? value : selectedVariant.material
      };
      
      return variant.size === currentSelection.size &&
             variant.color === currentSelection.color &&
             variant.material === currentSelection.material;
    });
    
    if (newVariant) {
      setSelectedVariant(newVariant);
    }
  };

  const handleAddToCart = () => {
    addToCart({
      product,
      variant: selectedVariant,
      quantity
    });
  };

  const handleWishlistToggle = () => {
    if (isInWishlist) {
      removeFromWishlist(product.id);
    } else {
      addToWishlist(product);
    }
  };

  const nextImage = () => {
    setSelectedImageIndex((prev) => 
      prev === product.images.length - 1 ? 0 : prev + 1
    );
  };

  const prevImage = () => {
    setSelectedImageIndex((prev) => 
      prev === 0 ? product.images.length - 1 : prev - 1
    );
  };

  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
          {/* Product Images */}
          <div className="space-y-4">
            {/* Main Image */}
            <div className="relative aspect-square bg-gray-100 rounded-lg overflow-hidden">
              <Image
                src={product.images[selectedImageIndex]}
                alt={product.name}
                fill
                className="object-cover"
                priority
              />
              
              {/* Image Navigation */}
              {product.images.length > 1 && (
                <>
                  <button
                    onClick={prevImage}
                    className="absolute left-4 top-1/2 transform -translate-y-1/2 bg-white bg-opacity-80 hover:bg-opacity-100 rounded-full p-2 transition-all duration-200"
                  >
                    <ChevronLeft className="h-5 w-5" />
                  </button>
                  <button
                    onClick={nextImage}
                    className="absolute right-4 top-1/2 transform -translate-y-1/2 bg-white bg-opacity-80 hover:bg-opacity-100 rounded-full p-2 transition-all duration-200"
                  >
                    <ChevronRight className="h-5 w-5" />
                  </button>
                </>
              )}
              
              {/* Discount Badge */}
              {hasDiscount && (
                <div className="absolute top-4 left-4 bg-red-500 text-white px-3 py-1 rounded-md font-medium">
                  -{discountPercentage}%
                </div>
              )}
            </div>
            
            {/* Thumbnail Images */}
            {product.images.length > 1 && (
              <div className="grid grid-cols-4 gap-4">
                {product.images.map((image, index) => (
                  <button
                    key={index}
                    onClick={() => setSelectedImageIndex(index)}
                    className={`relative aspect-square bg-gray-100 rounded-lg overflow-hidden border-2 transition-colors duration-200 ${
                      selectedImageIndex === index ? 'border-blue-500' : 'border-transparent hover:border-gray-300'
                    }`}
                  >
                    <Image
                      src={image}
                      alt={`${product.name} ${index + 1}`}
                      fill
                      className="object-cover"
                    />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Product Info */}
          <div className="space-y-6">
            {/* Basic Info */}
            <div>
              <p className="text-sm text-gray-500 mb-2">{product.category}</p>
              <h1 className="text-3xl font-bold text-gray-900 mb-4">{product.name}</h1>
              
              {/* Rating */}
              {product.rating && (
                <div className="flex items-center gap-2 mb-4">
                  <div className="flex items-center">
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        className={`h-5 w-5 ${
                          i < Math.floor(product.rating)
                            ? 'text-yellow-400 fill-current'
                            : 'text-gray-300'
                        }`}
                      />
                    ))}
                  </div>
                  <span className="text-sm text-gray-600">
                    {product.rating} ({product.reviewCount || 0} reviews)
                  </span>
                </div>
              )}
              
              {/* Price */}
              <div className="flex items-center gap-3 mb-6">
                <span className="text-3xl font-bold text-gray-900">
                  ${selectedVariant.price.toFixed(2)}
                </span>
                {hasDiscount && (
                  <span className="text-xl text-gray-500 line-through">
                    ${selectedVariant.originalPrice.toFixed(2)}
                  </span>
                )}
              </div>
            </div>

            {/* Variant Selection */}
            <div className="space-y-4">
              {/* Size Selection */}
              {variantOptions.sizes.length > 1 && (
                <div>
                  <h3 className="text-sm font-medium text-gray-900 mb-3">Size: {selectedVariant.size}</h3>
                  <div className="flex flex-wrap gap-2">
                    {variantOptions.sizes.map(size => {
                      const isAvailable = product.variants.some(v => 
                        v.size === size && 
                        v.color === selectedVariant.color && 
                        v.material === selectedVariant.material &&
                        v.stock > 0
                      );
                      
                      return (
                        <button
                          key={size}
                          onClick={() => handleVariantChange('size', size)}
                          disabled={!isAvailable}
                          className={`px-4 py-2 border rounded-md text-sm font-medium transition-colors duration-200 ${
                            selectedVariant.size === size
                              ? 'border-blue-500 bg-blue-50 text-blue-700'
                              : isAvailable
                              ? 'border-gray-300 hover:border-gray-400'
                              : 'border-gray-200 text-gray-400 cursor-not-allowed'
                          }`}
                        >
                          {size}
                        </button>
                      );
                    })}
                  </div>
                </div>
              )}

              {/* Color Selection */}
              {variantOptions.colors.length > 1 && (
                <div>
                  <h3 className="text-sm font-medium text-gray-900 mb-3">Color: {selectedVariant.color}</h3>
                  <div className="flex flex-wrap gap-2">
                    {variantOptions.colors.map(color => {
                      const isAvailable = product.variants.some(v => 
                        v.color === color && 
                        v.size === selectedVariant.size && 
                        v.material === selectedVariant.material &&
                        v.stock > 0
                      );
                      
                      return (
                        <button
                          key={color}
                          onClick={() => handleVariantChange('color', color)}
                          disabled={!isAvailable}
                          className={`w-8 h-8 rounded-full border-2 transition-all duration-200 ${
                            selectedVariant.color === color
                              ? 'border-blue-500 ring-2 ring-blue-200'
                              : isAvailable
                              ? 'border-gray-300 hover:border-gray-400'
                              : 'border-gray-200 opacity-50 cursor-not-allowed'
                          }`}
                          style={{ backgroundColor: color.toLowerCase() }}
                          title={color}
                        />
                      );
                    })}
                  </div>
                </div>
              )}

              {/* Material Selection */}
              {variantOptions.materials.length > 1 && (
                <div>
                  <h3 className="text-sm font-medium text-gray-900 mb-3">Material: {selectedVariant.material}</h3>
                  <div className="flex flex-wrap gap-2">
                    {variantOptions.materials.map(material => {
                      const isAvailable = product.variants.some(v => 
                        v.material === material && 
                        v.size === selectedVariant.size && 
                        v.color === selectedVariant.color &&
                        v.stock > 0
                      );
                      
                      return (
                        <button
                          key={material}
                          onClick={() => handleVariantChange('material', material)}
                          disabled={!isAvailable}
                          className={`px-4 py-2 border rounded-md text-sm font-medium transition-colors duration-200 ${
                            selectedVariant.material === material
                              ? 'border-blue-500 bg-blue-50 text-blue-700'
                              : isAvailable
                              ? 'border-gray-300 hover:border-gray-400'
                              : 'border-gray-200 text-gray-400 cursor-not-allowed'
                          }`}
                        >
                          {material}
                        </button>
                      );
                    })}
                  </div>
                </div>
              )}
            </div>

            {/* Stock Status */}
            <div className="flex items-center gap-2">
              <div className={`w-3 h-3 rounded-full ${
                selectedVariant.stock > 10 ? 'bg-green-500' :
                selectedVariant.stock > 0 ? 'bg-yellow-500' : 'bg-red-500'
              }`} />
              <span className="text-sm text-gray-600">
                {selectedVariant.stock > 10 ? 'In Stock' :
                 selectedVariant.stock > 0 ? `Only ${selectedVariant.stock} left` : 'Out of Stock'}
              </span>
            </div>

            {/* Quantity and Add to Cart */}
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <div className="flex items-center border border-gray-300 rounded-md">
                  <button
                    onClick={() => setQuantity(Math.max(1, quantity - 1))}
                    className="px-3 py-2 hover:bg-gray-50"
                  >
                    -
                  </button>
                  <span className="px-4 py-2 border-x border-gray-300">{quantity}</span>
                  <button
                    onClick={() => setQuantity(Math.min(selectedVariant.stock, quantity + 1))}
                    className="px-3 py-2 hover:bg-gray-50"
                  >
                    +
                  </button>
                </div>
                <span className="text-sm text-gray-600">Max: {selectedVariant.stock}</span>
              </div>

              <div className="flex gap-4">
                <button
                  onClick={handleAddToCart}
                  disabled={selectedVariant.stock === 0}
                  className="flex-1 bg-black text-white py-3 px-6 rounded-md hover:bg-gray-800 transition-colors duration-200 disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                >
                  <ShoppingCart className="h-5 w-5" />
                  {selectedVariant.stock === 0 ? 'Out of Stock' : 'Add to Cart'}
                </button>
                
                <button
                  onClick={handleWishlistToggle}
                  className={`p-3 border rounded-md transition-colors duration-200 ${
                    isInWishlist
                      ? 'border-red-500 bg-red-50 text-red-500'
                      : 'border-gray-300 hover:border-red-300 hover:bg-red-50 hover:text-red-500'
                  }`}
                >
                  <Heart className={`h-5 w-5 ${isInWishlist ? 'fill-current' : ''}`} />
                </button>
                
                <button
                  onClick={() => addToComparison(product)}
                  className="p-3 border border-gray-300 rounded-md hover:border-blue-300 hover:bg-blue-50 hover:text-blue-500 transition-colors duration-200"
                >
                  <Share2 className="h-5 w-5" />
                </button>
              </div>
            </div>

            {/* Features */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-6 border-t border-gray-200">
              <div className="flex items-center gap-3">
                <Truck className="h-5 w-5 text-green-600" />
                <div>
                  <p className="text-sm font-medium text-gray-900">Free Shipping</p>
                  <p className="text-xs text-gray-500">On orders over $50</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <RotateCcw className="h-5 w-5 text-blue-600" />
                <div>
                  <p className="text-sm font-medium text-gray-900">Easy Returns</p>
                  <p className="text-xs text-gray-500">30-day return policy</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <Shield className="h-5 w-5 text-purple-600" />
                <div>
                  <p className="text-sm font-medium text-gray-900">Secure Payment</p>
                  <p className="text-xs text-gray-500">SSL encrypted</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Product Details Tabs */}
        <div className="mt-16">
          <div className="border-b border-gray-200">
            <nav className="flex space-x-8">
              {['description', 'specifications', 'reviews'].map(tab => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={`py-4 px-1 border-b-2 font-medium text-sm capitalize transition-colors duration-200 ${
                    activeTab === tab
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  {tab}
                </button>
              ))}
            </nav>
          </div>

          <div className="py-8">
            {activeTab === 'description' && (
              <div className="prose max-w-none">
                <p className="text-gray-700 leading-relaxed">{product.description}</p>
              </div>
            )}
            
            {activeTab === 'specifications' && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 className="font-medium text-gray-900 mb-4">Product Details</h3>
                  <dl className="space-y-3">
                    <div className="flex justify-between">
                      <dt className="text-gray-500">SKU:</dt>
                      <dd className="text-gray-900">{selectedVariant.sku}</dd>
                    </div>
                    <div className="flex justify-between">
                      <dt className="text-gray-500">Category:</dt>
                      <dd className="text-gray-900">{product.category}</dd>
                    </div>
                    <div className="flex justify-between">
                      <dt className="text-gray-500">Material:</dt>
                      <dd className="text-gray-900">{selectedVariant.material}</dd>
                    </div>
                  </dl>
                </div>
              </div>
            )}
            
            {activeTab === 'reviews' && (
              <div className="text-center py-12">
                <p className="text-gray-500">No reviews yet. Be the first to review this product!</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}