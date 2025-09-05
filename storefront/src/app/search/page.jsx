'use client';

import { useEffect } from 'react';
import { useAtom } from 'jotai';
import { useSearchParams } from 'next/navigation';
import { searchQueryAtom } from '@/store/searchStore';
import SearchBar from '@/components/search/SearchBar';
import SearchResults from '@/components/search/SearchResults';

export default function SearchPage() {
  const searchParams = useSearchParams();
  const [, setSearchQuery] = useAtom(searchQueryAtom);

  // Set search query from URL params on page load
  useEffect(() => {
    const query = searchParams.get('q');
    if (query) {
      setSearchQuery(query);
    }
  }, [searchParams, setSearchQuery]);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Search Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-center">
            <SearchBar placeholder="Search for products, brands, categories..." />
          </div>
        </div>
      </div>

      {/* Search Results */}
      <SearchResults />
    </div>
  );
}