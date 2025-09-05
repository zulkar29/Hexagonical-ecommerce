import { atom } from 'jotai';

// Search query atom
export const searchQueryAtom = atom('');

// Search filters atom
export const searchFiltersAtom = atom({
  category: '',
  priceRange: { min: 0, max: 1000 },
  sizes: [],
  colors: [],
  materials: [],
  inStock: false,
  onSale: false
});

// Sort options atom
export const sortOptionAtom = atom('featured'); // featured, price-low, price-high, newest, rating

// Search results atom (derived)
export const searchResultsAtom = atom((get) => {
  // This would typically fetch from an API
  // For now, we'll return empty array as placeholder
  return [];
});

// Active filters count atom (derived)
export const activeFiltersCountAtom = atom((get) => {
  const filters = get(searchFiltersAtom);
  let count = 0;
  
  if (filters.category) count++;
  if (filters.sizes.length > 0) count++;
  if (filters.colors.length > 0) count++;
  if (filters.materials.length > 0) count++;
  if (filters.inStock) count++;
  if (filters.onSale) count++;
  if (filters.priceRange.min > 0 || filters.priceRange.max < 1000) count++;
  
  return count;
});

// Update filter atom (write-only)
export const updateFilterAtom = atom(
  null,
  (get, set, { key, value }) => {
    const currentFilters = get(searchFiltersAtom);
    set(searchFiltersAtom, { ...currentFilters, [key]: value });
  }
);

// Clear all filters atom (write-only)
export const clearFiltersAtom = atom(
  null,
  (get, set) => {
    set(searchFiltersAtom, {
      category: '',
      priceRange: { min: 0, max: 1000 },
      sizes: [],
      colors: [],
      materials: [],
      inStock: false,
      onSale: false
    });
  }
);

// Product comparison atom
export const comparisonAtom = atom([]);

// Add to comparison atom (write-only)
export const addToComparisonAtom = atom(
  null,
  (get, set, product) => {
    const currentComparison = get(comparisonAtom);
    if (currentComparison.length < 4 && !currentComparison.some(p => p.id === product.id)) {
      set(comparisonAtom, [...currentComparison, product]);
    }
  }
);

// Remove from comparison atom (write-only)
export const removeFromComparisonAtom = atom(
  null,
  (get, set, productId) => {
    const currentComparison = get(comparisonAtom);
    set(comparisonAtom, currentComparison.filter(p => p.id !== productId));
  }
);

// Clear comparison atom (write-only)
export const clearComparisonAtom = atom(
  null,
  (get, set) => {
    set(comparisonAtom, []);
  }
);