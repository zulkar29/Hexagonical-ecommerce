import { atom } from 'jotai';
import { atomWithStorage } from 'jotai/utils';

// Cart state
export const cartAtom = atomWithStorage('cart', []);

// Wishlist state
export const wishlistAtom = atomWithStorage('wishlist', []);

// Comparison state (max 4 items)
export const comparisonAtom = atomWithStorage('comparison', []);

// Search state
export const searchQueryAtom = atom('');
export const filtersAtom = atom({
  categories: [],
  priceRange: { min: null, max: null },
  minRating: null,
  brands: [],
  features: [],
  sizes: [],
  colors: [],
  materials: []
});
export const sortAtom = atom('relevance');
export const searchResultsAtom = atom([]);

// Legacy aliases for backward compatibility
export const searchFiltersAtom = filtersAtom;
export const searchSortAtom = sortAtom;

// User preferences
export const themeAtom = atomWithStorage('theme', 'light');
export const currencyAtom = atomWithStorage('currency', 'USD');
export const languageAtom = atomWithStorage('language', 'en');

// UI state
export const mobileMenuOpenAtom = atom(false);
export const cartOpenAtom = atom(false);
export const searchOpenAtom = atom(false);

// Derived atoms
export const cartCountAtom = atom((get) => {
  const cart = get(cartAtom);
  return cart.reduce((total, item) => total + item.quantity, 0);
});

export const cartTotalAtom = atom((get) => {
  const cart = get(cartAtom);
  return cart.reduce((total, item) => total + (item.price * item.quantity), 0);
});

export const wishlistCountAtom = atom((get) => {
  const wishlist = get(wishlistAtom);
  return wishlist.length;
});

export const comparisonCountAtom = atom((get) => {
  const comparison = get(comparisonAtom);
  return comparison.length;
});

// Search derived atoms
export const filteredProductsAtom = atom((get) => {
  const query = get(searchQueryAtom);
  const filters = get(filtersAtom);
  const sort = get(sortAtom);
  const results = get(searchResultsAtom);
  
  // This would typically be handled by the search component
  // For now, return the search results as-is
  return results;
});

// Actions
export const addToCartAtom = atom(
  null,
  (get, set, product) => {
    const cart = get(cartAtom);
    const existingItem = cart.find(item => 
      item.id === product.id && 
      JSON.stringify(item.variant) === JSON.stringify(product.variant)
    );
    
    if (existingItem) {
      set(cartAtom, cart.map(item => 
        item.id === product.id && JSON.stringify(item.variant) === JSON.stringify(product.variant)
          ? { ...item, quantity: item.quantity + (product.quantity || 1) }
          : item
      ));
    } else {
      set(cartAtom, [...cart, { ...product, quantity: product.quantity || 1 }]);
    }
  }
);

export const removeFromCartAtom = atom(
  null,
  (get, set, { id, variant }) => {
    const cart = get(cartAtom);
    set(cartAtom, cart.filter(item => 
      !(item.id === id && JSON.stringify(item.variant) === JSON.stringify(variant))
    ));
  }
);

export const updateCartQuantityAtom = atom(
  null,
  (get, set, { id, variant, quantity }) => {
    const cart = get(cartAtom);
    if (quantity <= 0) {
      set(removeFromCartAtom, { id, variant });
      return;
    }
    
    set(cartAtom, cart.map(item => 
      item.id === id && JSON.stringify(item.variant) === JSON.stringify(variant)
        ? { ...item, quantity }
        : item
    ));
  }
);

export const clearCartAtom = atom(
  null,
  (get, set) => {
    set(cartAtom, []);
  }
);

export const toggleWishlistAtom = atom(
  null,
  (get, set, product) => {
    const wishlist = get(wishlistAtom);
    const isInWishlist = wishlist.some(item => item.id === product.id);
    
    if (isInWishlist) {
      set(wishlistAtom, wishlist.filter(item => item.id !== product.id));
    } else {
      set(wishlistAtom, [...wishlist, product]);
    }
  }
);

export const toggleComparisonAtom = atom(
  null,
  (get, set, product) => {
    const comparison = get(comparisonAtom);
    const isInComparison = comparison.some(item => item.id === product.id);
    
    if (isInComparison) {
      set(comparisonAtom, comparison.filter(item => item.id !== product.id));
    } else if (comparison.length < 4) {
      set(comparisonAtom, [...comparison, product]);
    }
  }
);

export const clearComparisonAtom = atom(
  null,
  (get, set) => {
    set(comparisonAtom, []);
  }
);

export const clearFiltersAtom = atom(
  null,
  (get, set) => {
    set(filtersAtom, {
      categories: [],
      priceRange: { min: null, max: null },
      minRating: null,
      brands: [],
      features: [],
      sizes: [],
      colors: [],
      materials: []
    });
  }
);

// Action alias for SearchFilters component
export const clearFiltersAction = clearFiltersAtom;

// Checkout state
export const checkoutStepAtom = atom(1);
export const shippingInfoAtom = atom({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  address: '',
  city: '',
  state: '',
  zipCode: '',
  country: ''
});
export const paymentInfoAtom = atom({
  cardNumber: '',
  expiryDate: '',
  cvv: '',
  cardholderName: ''
});
export const orderNotesAtom = atom('');

// Notification state
export const notificationsAtom = atom([]);

export const addNotificationAtom = atom(
  null,
  (get, set, notification) => {
    const notifications = get(notificationsAtom);
    const newNotification = {
      id: Date.now(),
      timestamp: new Date(),
      ...notification
    };
    set(notificationsAtom, [...notifications, newNotification]);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
      set(removeNotificationAtom, newNotification.id);
    }, 5000);
  }
);

export const removeNotificationAtom = atom(
  null,
  (get, set, id) => {
    const notifications = get(notificationsAtom);
    set(notificationsAtom, notifications.filter(n => n.id !== id));
  }
);

// Recent searches
export const recentSearchesAtom = atomWithStorage('recentSearches', []);

export const addRecentSearchAtom = atom(
  null,
  (get, set, query) => {
    if (!query.trim()) return;
    
    const recent = get(recentSearchesAtom);
    const filtered = recent.filter(item => item !== query);
    const updated = [query, ...filtered].slice(0, 10); // Keep only 10 recent searches
    
    set(recentSearchesAtom, updated);
  }
);

export const clearRecentSearchesAtom = atom(
  null,
  (get, set) => {
    set(recentSearchesAtom, []);
  }
);