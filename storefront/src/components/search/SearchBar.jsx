'use client';

import { useState, useEffect, useRef } from 'react';
import { useAtom } from 'jotai';
import { Search, X, Filter, TrendingUp } from 'lucide-react';
import { searchQueryAtom, filtersAtom, sortAtom } from '@/store/atoms';
import { getProducts } from '@/data/products';

export default function SearchBar({ showFilters = true, placeholder = "Search products..." }) {
  const [searchQuery, setSearchQuery] = useAtom(searchQueryAtom);
  const [filters, setFilters] = useAtom(filtersAtom);
  const [sort, setSort] = useAtom(sortAtom);
  const [isOpen, setIsOpen] = useState(false);
  const [suggestions, setSuggestions] = useState([]);
  const [recentSearches, setRecentSearches] = useState([]);
  const searchRef = useRef(null);
  const dropdownRef = useRef(null);

  // Load recent searches from localStorage
  useEffect(() => {
    const saved = localStorage.getItem('recentSearches');
    if (saved) {
      setRecentSearches(JSON.parse(saved));
    }
  }, []);

  // Generate suggestions based on search query
  useEffect(() => {
    if (searchQuery.length > 1) {
      const products = getProducts();
      const filtered = products
        .filter(product => 
          product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          product.category.toLowerCase().includes(searchQuery.toLowerCase()) ||
          product.brand?.toLowerCase().includes(searchQuery.toLowerCase())
        )
        .slice(0, 5);
      setSuggestions(filtered);
    } else {
      setSuggestions([]);
    }
  }, [searchQuery]);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleSearch = (query) => {
    setSearchQuery(query);
    setIsOpen(false);
    
    // Add to recent searches
    if (query.trim()) {
      const updated = [query, ...recentSearches.filter(s => s !== query)].slice(0, 5);
      setRecentSearches(updated);
      localStorage.setItem('recentSearches', JSON.stringify(updated));
    }
  };

  const clearSearch = () => {
    setSearchQuery('');
    setIsOpen(false);
    searchRef.current?.focus();
  };

  const clearRecentSearches = () => {
    setRecentSearches([]);
    localStorage.removeItem('recentSearches');
  };

  const popularSearches = [
    'wireless headphones',
    'summer dress',
    'laptop',
    'sneakers',
    'home decor'
  ];

  return (
    <div className="relative w-full max-w-2xl" ref={dropdownRef}>
      {/* Search Input */}
      <div className="relative">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search className="h-5 w-5 text-gray-400" />
        </div>
        <input
          ref={searchRef}
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          onFocus={() => setIsOpen(true)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              handleSearch(searchQuery);
            }
            if (e.key === 'Escape') {
              setIsOpen(false);
            }
          }}
          placeholder={placeholder}
          className="w-full pl-10 pr-12 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all duration-200"
        />
        {searchQuery && (
          <button
            onClick={clearSearch}
            className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
          >
            <X className="h-5 w-5" />
          </button>
        )}
      </div>

      {/* Search Dropdown */}
      {isOpen && (
        <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-lg shadow-lg z-50 max-h-96 overflow-y-auto">
          {/* Suggestions */}
          {suggestions.length > 0 && (
            <div className="p-4">
              <h3 className="text-sm font-medium text-gray-900 mb-3">Products</h3>
              <div className="space-y-2">
                {suggestions.map((product) => (
                  <button
                    key={product.id}
                    onClick={() => handleSearch(product.name)}
                    className="w-full flex items-center gap-3 p-2 hover:bg-gray-50 rounded-md text-left"
                  >
                    <img
                      src={product.images[0]}
                      alt={product.name}
                      className="w-10 h-10 object-cover rounded"
                    />
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {product.name}
                      </p>
                      <p className="text-sm text-gray-500">
                        {product.category} • ${product.price}
                      </p>
                    </div>
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Recent Searches */}
          {recentSearches.length > 0 && suggestions.length === 0 && (
            <div className="p-4 border-t border-gray-100">
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-sm font-medium text-gray-900">Recent Searches</h3>
                <button
                  onClick={clearRecentSearches}
                  className="text-xs text-gray-500 hover:text-gray-700"
                >
                  Clear all
                </button>
              </div>
              <div className="space-y-1">
                {recentSearches.map((search, index) => (
                  <button
                    key={index}
                    onClick={() => handleSearch(search)}
                    className="w-full flex items-center gap-2 p-2 hover:bg-gray-50 rounded-md text-left"
                  >
                    <Search className="h-4 w-4 text-gray-400" />
                    <span className="text-sm text-gray-700">{search}</span>
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Popular Searches */}
          {suggestions.length === 0 && (
            <div className="p-4 border-t border-gray-100">
              <h3 className="text-sm font-medium text-gray-900 mb-3 flex items-center gap-2">
                <TrendingUp className="h-4 w-4" />
                Popular Searches
              </h3>
              <div className="space-y-1">
                {popularSearches.map((search, index) => (
                  <button
                    key={index}
                    onClick={() => handleSearch(search)}
                    className="w-full flex items-center gap-2 p-2 hover:bg-gray-50 rounded-md text-left"
                  >
                    <TrendingUp className="h-4 w-4 text-gray-400" />
                    <span className="text-sm text-gray-700">{search}</span>
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Quick Filters */}
          {showFilters && (
            <div className="p-4 border-t border-gray-100">
              <h3 className="text-sm font-medium text-gray-900 mb-3 flex items-center gap-2">
                <Filter className="h-4 w-4" />
                Quick Filters
              </h3>
              <div className="flex flex-wrap gap-2">
                {[
                  { label: 'On Sale', key: 'onSale', value: true },
                  { label: 'Free Shipping', key: 'freeShipping', value: true },
                  { label: 'In Stock', key: 'inStock', value: true },
                  { label: 'New Arrivals', key: 'newArrivals', value: true }
                ].map((filter) => (
                  <button
                    key={filter.key}
                    onClick={() => {
                      setFilters(prev => ({
                        ...prev,
                        [filter.key]: prev[filter.key] ? null : filter.value
                      }));
                      setIsOpen(false);
                    }}
                    className={`px-3 py-1 text-xs rounded-full border transition-colors duration-200 ${
                      filters[filter.key]
                        ? 'bg-blue-100 border-blue-300 text-blue-700'
                        : 'bg-gray-100 border-gray-300 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    {filter.label}
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}