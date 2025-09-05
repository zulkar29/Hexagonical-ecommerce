import { atom } from 'jotai';
import { atomWithStorage } from 'jotai/utils';

// User preferences with localStorage persistence
export const userPreferencesAtom = atomWithStorage('user-preferences', {
  theme: 'light',
  currency: 'USD',
  language: 'en'
});

// Theme atom (derived)
export const themeAtom = atom(
  (get) => get(userPreferencesAtom).theme,
  (get, set, newTheme) => {
    const preferences = get(userPreferencesAtom);
    set(userPreferencesAtom, { ...preferences, theme: newTheme });
  }
);

// Currency atom (derived)
export const currencyAtom = atom(
  (get) => get(userPreferencesAtom).currency,
  (get, set, newCurrency) => {
    const preferences = get(userPreferencesAtom);
    set(userPreferencesAtom, { ...preferences, currency: newCurrency });
  }
);

// Wishlist atom with localStorage persistence
export const wishlistAtom = atomWithStorage('wishlist', []);

// Add to wishlist atom (write-only)
export const addToWishlistAtom = atom(
  null,
  (get, set, product) => {
    const currentWishlist = get(wishlistAtom);
    const isAlreadyInWishlist = currentWishlist.some(item => item.id === product.id);
    
    if (!isAlreadyInWishlist) {
      set(wishlistAtom, [...currentWishlist, product]);
    }
  }
);

// Remove from wishlist atom (write-only)
export const removeFromWishlistAtom = atom(
  null,
  (get, set, productId) => {
    const currentWishlist = get(wishlistAtom);
    set(wishlistAtom, currentWishlist.filter(item => item.id !== productId));
  }
);

// Check if product is in wishlist atom (derived)
export const isInWishlistAtom = atom(
  null,
  (get, set, productId) => {
    const wishlist = get(wishlistAtom);
    return wishlist.some(item => item.id === productId);
  }
);

// Recently viewed products atom
export const recentlyViewedAtom = atomWithStorage('recently-viewed', []);

// Add to recently viewed atom (write-only)
export const addToRecentlyViewedAtom = atom(
  null,
  (get, set, product) => {
    const currentViewed = get(recentlyViewedAtom);
    const filteredViewed = currentViewed.filter(item => item.id !== product.id);
    const newViewed = [product, ...filteredViewed].slice(0, 10); // Keep only 10 items
    set(recentlyViewedAtom, newViewed);
  }
);