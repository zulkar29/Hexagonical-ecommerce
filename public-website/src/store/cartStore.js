import { atom } from 'jotai';
import { atomWithStorage } from 'jotai/utils';

// Cart items atom with localStorage persistence
export const cartItemsAtom = atomWithStorage('cart-items', []);

// Cart total count atom (derived)
export const cartCountAtom = atom((get) => {
  const items = get(cartItemsAtom);
  return items.reduce((total, item) => total + item.quantity, 0);
});

// Cart total price atom (derived)
export const cartTotalAtom = atom((get) => {
  const items = get(cartItemsAtom);
  return items.reduce((total, item) => total + (item.price * item.quantity), 0);
});

// Add item to cart atom (write-only)
export const addToCartAtom = atom(
  null,
  (get, set, { product, variant, quantity = 1 }) => {
    const currentItems = get(cartItemsAtom);
    const existingItemIndex = currentItems.findIndex(
      item => item.productId === product.id && 
              item.variantId === variant.id
    );

    if (existingItemIndex >= 0) {
      // Update existing item quantity
      const updatedItems = [...currentItems];
      updatedItems[existingItemIndex].quantity += quantity;
      set(cartItemsAtom, updatedItems);
    } else {
      // Add new item
      const newItem = {
        id: `${product.id}-${variant.id}`,
        productId: product.id,
        variantId: variant.id,
        name: product.name,
        price: variant.price,
        image: product.images[0],
        variant: {
          size: variant.size,
          color: variant.color,
          material: variant.material
        },
        quantity,
        sku: variant.sku
      };
      set(cartItemsAtom, [...currentItems, newItem]);
    }
  }
);

// Remove item from cart atom (write-only)
export const removeFromCartAtom = atom(
  null,
  (get, set, itemId) => {
    const currentItems = get(cartItemsAtom);
    set(cartItemsAtom, currentItems.filter(item => item.id !== itemId));
  }
);

// Update item quantity atom (write-only)
export const updateQuantityAtom = atom(
  null,
  (get, set, { itemId, quantity }) => {
    const currentItems = get(cartItemsAtom);
    if (quantity <= 0) {
      set(cartItemsAtom, currentItems.filter(item => item.id !== itemId));
    } else {
      const updatedItems = currentItems.map(item =>
        item.id === itemId ? { ...item, quantity } : item
      );
      set(cartItemsAtom, updatedItems);
    }
  }
);

// Clear cart atom (write-only)
export const clearCartAtom = atom(
  null,
  (get, set) => {
    set(cartItemsAtom, []);
  }
);