'use client';

import { useState, useMemo } from 'react';
import { useAtom } from 'jotai';
import { Filter, Grid, List, ChevronDown } from 'lucide-react';
import ProductCard from '@/components/product/ProductCard';
import { products, categories } from '@/data/products';
import { searchFiltersAtom, sortOptionAtom, updateFilterAtom, clearFiltersAtom, activeFiltersCountAtom } from '@/store/searchStore';

export default function ProductsPage() {
  const [viewMode, setViewMode] = useState('grid');
  const [showFilters, setShowFilters] = useState(false);
  const [filters] = useAtom(searchFiltersAtom);
  const [sortOption, setSortOption] = useAtom(sortOptionAtom);
  const [, updateFilter] = useAtom(updateFilterAtom);
  const [, clearFilters] = useAtom(clearFiltersAtom);
  const [activeFiltersCount] = useAtom(activeFiltersCountAtom);

  // Filter and sort products
  const filteredAndSortedProducts = useMemo(() => {
    let filtered = products;

    // Apply category filter
    if (filters.category) {
      filtered = filtered.filter(product => product.category === filters.category);
    }

    // Apply price range filter
    filtered = filtered.filter(product => {
      const price = product.variants[0].price;
      return price >= filters.priceRange.min && price <= filters.priceRange.max;
    });

    // Apply size filter
    if (filters.sizes.length > 0) {
      filtered = filtered.filter(product =>
        product.variants.some(variant => filters.sizes.includes(variant.size))
      );
    }

    // Apply color filter
    if (filters.colors.length > 0) {
      filtered = filtered.filter(product =>
        product.variants.some(variant => filters.colors.includes(variant.color))
      );
    }

    // Apply material filter
    if (filters.materials.length > 0) {
      filtered = filtered.filter(product =>
        product.variants.some(variant => filters.materials.includes(variant.material))
      );
    }

    // Apply stock filter
    if (filters.inStock) {
      filtered = filtered.filter(product =>
        product.variants.some(variant => variant.stock > 0)
      );
    }

    // Apply sale filter
    if (filters.onSale) {
      filtered = filtered.filter(product =>
        product.variants.some(variant => variant.originalPrice > variant.price)
      );
    }

    // Sort products
    switch (sortOption) {
      case 'price-low':
        filtered.sort((a, b) => a.variants[0].price - b.variants[0].price);
        break;
      case 'price-high':
        filtered.sort((a, b) => b.variants[0].price - a.variants[0].price);
        break;
      case 'newest':
        filtered.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
        break;
      case 'rating':
        filtered.sort((a, b) => (b.rating || 0) - (a.rating || 0));
        break;
      default: // featured
        break;
    }

    return filtered;
  }, [products, filters, sortOption]);

  const sortOptions = [
    { value: 'featured', label: 'Featured' },
    { value: 'price-low', label: 'Price: Low to High' },
    { value: 'price-high', label: 'Price: High to Low' },
    { value: 'newest', label: 'Newest' },
    { value: 'rating', label: 'Customer Rating' }
  ];

  const availableSizes = [...new Set(products.flatMap(p => p.variants.map(v => v.size)))];
  const availableColors = [...new Set(products.flatMap(p => p.variants.map(v => v.color)))];
  const availableMaterials = [...new Set(products.flatMap(p => p.variants.map(v => v.material)))];

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">All Products</h1>
          <p className="text-gray-600">Discover our complete collection of premium products</p>
        </div>

        <div className="flex flex-col lg:flex-row gap-8">
          {/* Filters Sidebar */}
          <div className={`lg:w-64 ${showFilters ? 'block' : 'hidden lg:block'}`}>
            <div className="bg-white rounded-lg shadow-sm p-6 sticky top-4">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold text-gray-900">Filters</h3>
                {activeFiltersCount > 0 && (
                  <button
                    onClick={clearFilters}
                    className="text-sm text-blue-600 hover:text-blue-700"
                  >
                    Clear All ({activeFiltersCount})
                  </button>
                )}
              </div>

              {/* Category Filter */}
              <div className="mb-6">
                <h4 className="font-medium text-gray-900 mb-3">Category</h4>
                <div className="space-y-2">
                  {categories.map(category => (
                    <label key={category.id} className="flex items-center">
                      <input
                        type="radio"
                        name="category"
                        value={category.name}
                        checked={filters.category === category.name}
                        onChange={(e) => updateFilter({ key: 'category', value: e.target.value })}
                        className="mr-2"
                      />
                      <span className="text-sm text-gray-700">{category.name}</span>
                    </label>
                  ))}
                </div>
              </div>

              {/* Price Range Filter */}
              <div className="mb-6">
                <h4 className="font-medium text-gray-900 mb-3">Price Range</h4>
                <div className="space-y-3">
                  <div>
                    <label className="block text-sm text-gray-700 mb-1">Min: ${filters.priceRange.min}</label>
                    <input
                      type="range"
                      min="0"
                      max="1000"
                      value={filters.priceRange.min}
                      onChange={(e) => updateFilter({
                        key: 'priceRange',
                        value: { ...filters.priceRange, min: parseInt(e.target.value) }
                      })}
                      className="w-full"
                    />
                  </div>
                  <div>
                    <label className="block text-sm text-gray-700 mb-1">Max: ${filters.priceRange.max}</label>
                    <input
                      type="range"
                      min="0"
                      max="1000"
                      value={filters.priceRange.max}
                      onChange={(e) => updateFilter({
                        key: 'priceRange',
                        value: { ...filters.priceRange, max: parseInt(e.target.value) }
                      })}
                      className="w-full"
                    />
                  </div>
                </div>
              </div>

              {/* Size Filter */}
              <div className="mb-6">
                <h4 className="font-medium text-gray-900 mb-3">Size</h4>
                <div className="grid grid-cols-3 gap-2">
                  {availableSizes.map(size => (
                    <label key={size} className="flex items-center">
                      <input
                        type="checkbox"
                        checked={filters.sizes.includes(size)}
                        onChange={(e) => {
                          const newSizes = e.target.checked
                            ? [...filters.sizes, size]
                            : filters.sizes.filter(s => s !== size);
                          updateFilter({ key: 'sizes', value: newSizes });
                        }}
                        className="mr-1"
                      />
                      <span className="text-sm text-gray-700">{size}</span>
                    </label>
                  ))}
                </div>
              </div>

              {/* Color Filter */}
              <div className="mb-6">
                <h4 className="font-medium text-gray-900 mb-3">Color</h4>
                <div className="grid grid-cols-2 gap-2">
                  {availableColors.map(color => (
                    <label key={color} className="flex items-center">
                      <input
                        type="checkbox"
                        checked={filters.colors.includes(color)}
                        onChange={(e) => {
                          const newColors = e.target.checked
                            ? [...filters.colors, color]
                            : filters.colors.filter(c => c !== color);
                          updateFilter({ key: 'colors', value: newColors });
                        }}
                        className="mr-2"
                      />
                      <span className="text-sm text-gray-700">{color}</span>
                    </label>
                  ))}
                </div>
              </div>

              {/* Additional Filters */}
              <div className="space-y-3">
                <label className="flex items-center">
                  <input
                    type="checkbox"
                    checked={filters.inStock}
                    onChange={(e) => updateFilter({ key: 'inStock', value: e.target.checked })}
                    className="mr-2"
                  />
                  <span className="text-sm text-gray-700">In Stock Only</span>
                </label>
                <label className="flex items-center">
                  <input
                    type="checkbox"
                    checked={filters.onSale}
                    onChange={(e) => updateFilter({ key: 'onSale', value: e.target.checked })}
                    className="mr-2"
                  />
                  <span className="text-sm text-gray-700">On Sale</span>
                </label>
              </div>
            </div>
          </div>

          {/* Products Grid */}
          <div className="flex-1">
            {/* Toolbar */}
            <div className="bg-white rounded-lg shadow-sm p-4 mb-6">
              <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div className="flex items-center gap-4">
                  <button
                    onClick={() => setShowFilters(!showFilters)}
                    className="lg:hidden flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
                  >
                    <Filter className="h-4 w-4" />
                    Filters
                    {activeFiltersCount > 0 && (
                      <span className="bg-blue-600 text-white text-xs rounded-full px-2 py-1">
                        {activeFiltersCount}
                      </span>
                    )}
                  </button>
                  <span className="text-sm text-gray-600">
                    {filteredAndSortedProducts.length} products
                  </span>
                </div>

                <div className="flex items-center gap-4">
                  {/* Sort Dropdown */}
                  <div className="relative">
                    <select
                      value={sortOption}
                      onChange={(e) => setSortOption(e.target.value)}
                      className="appearance-none bg-white border border-gray-300 rounded-md px-4 py-2 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      {sortOptions.map(option => (
                        <option key={option.value} value={option.value}>
                          {option.label}
                        </option>
                      ))}
                    </select>
                    <ChevronDown className="absolute right-2 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400 pointer-events-none" />
                  </div>

                  {/* View Mode Toggle */}
                  <div className="flex border border-gray-300 rounded-md">
                    <button
                      onClick={() => setViewMode('grid')}
                      className={`p-2 ${viewMode === 'grid' ? 'bg-gray-100' : 'hover:bg-gray-50'}`}
                    >
                      <Grid className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => setViewMode('list')}
                      className={`p-2 ${viewMode === 'list' ? 'bg-gray-100' : 'hover:bg-gray-50'}`}
                    >
                      <List className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              </div>
            </div>

            {/* Products */}
            {filteredAndSortedProducts.length > 0 ? (
              <div className={`grid gap-6 ${
                viewMode === 'grid'
                  ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4'
                  : 'grid-cols-1'
              }`}>
                {filteredAndSortedProducts.map(product => (
                  <ProductCard
                    key={product.id}
                    product={product}
                    className={viewMode === 'list' ? 'flex-row' : ''}
                  />
                ))}
              </div>
            ) : (
              <div className="text-center py-12">
                <p className="text-gray-500 text-lg mb-4">No products found matching your criteria</p>
                <button
                  onClick={clearFilters}
                  className="text-blue-600 hover:text-blue-700 font-medium"
                >
                  Clear all filters
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}