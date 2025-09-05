'use client';

import { createContext, useContext, useReducer, useEffect } from 'react';

const CartContext = createContext();

const cartReducer = (state, action) => {
  switch (action.type) {
    case 'ADD_TO_CART': {
      const { product, variant, quantity = 1 } = action.payload;
      const existingItemIndex = state.items.findIndex(
        item => item.product.id === product.id && item.variant.id === variant.id
      );

      if (existingItemIndex > -1) {
        const updatedItems = [...state.items];
        updatedItems[existingItemIndex].quantity += quantity;
        return {
          ...state,
          items: updatedItems,
          totalItems: state.totalItems + quantity,
          totalPrice: state.totalPrice + (variant.price * quantity)
        };
      } else {
        return {
          ...state,
          items: [...state.items, { product, variant, quantity }],
          totalItems: state.totalItems + quantity,
          totalPrice: state.totalPrice + (variant.price * quantity)
        };
      }
    }

    case 'REMOVE_FROM_CART': {
      const { productId, variantId } = action.payload;
      const itemToRemove = state.items.find(
        item => item.product.id === productId && item.variant.id === variantId
      );
      
      if (!itemToRemove) return state;

      const updatedItems = state.items.filter(
        item => !(item.product.id === productId && item.variant.id === variantId)
      );

      return {
        ...state,
        items: updatedItems,
        totalItems: state.totalItems - itemToRemove.quantity,
        totalPrice: state.totalPrice - (itemToRemove.variant.price * itemToRemove.quantity)
      };
    }

    case 'UPDATE_QUANTITY': {
      const { productId, variantId, quantity } = action.payload;
      const itemIndex = state.items.findIndex(
        item => item.product.id === productId && item.variant.id === variantId
      );

      if (itemIndex === -1) return state;

      const updatedItems = [...state.items];
      const oldQuantity = updatedItems[itemIndex].quantity;
      const priceDifference = (quantity - oldQuantity) * updatedItems[itemIndex].variant.price;

      if (quantity <= 0) {
        return cartReducer(state, {
          type: 'REMOVE_FROM_CART',
          payload: { productId, variantId }
        });
      }

      updatedItems[itemIndex].quantity = quantity;

      return {
        ...state,
        items: updatedItems,
        totalItems: state.totalItems + (quantity - oldQuantity),
        totalPrice: state.totalPrice + priceDifference
      };
    }

    case 'CLEAR_CART':
      return {
        items: [],
        totalItems: 0,
        totalPrice: 0
      };

    case 'LOAD_CART':
      return action.payload;

    default:
      return state;
  }
};

const initialState = {
  items: [],
  totalItems: 0,
  totalPrice: 0
};

export const CartProvider = ({ children }) => {
  const [state, dispatch] = useReducer(cartReducer, initialState);

  // Load cart from localStorage on mount
  useEffect(() => {
    const savedCart = localStorage.getItem('cart');
    if (savedCart) {
      try {
        const parsedCart = JSON.parse(savedCart);
        dispatch({ type: 'LOAD_CART', payload: parsedCart });
      } catch (error) {
        console.error('Error loading cart from localStorage:', error);
      }
    }
  }, []);

  // Save cart to localStorage whenever it changes
  useEffect(() => {
    localStorage.setItem('cart', JSON.stringify(state));
  }, [state]);

  const addToCart = (product, variant, quantity = 1) => {
    dispatch({
      type: 'ADD_TO_CART',
      payload: { product, variant, quantity }
    });
  };

  const removeFromCart = (productId, variantId) => {
    dispatch({
      type: 'REMOVE_FROM_CART',
      payload: { productId, variantId }
    });
  };

  const updateQuantity = (productId, variantId, quantity) => {
    dispatch({
      type: 'UPDATE_QUANTITY',
      payload: { productId, variantId, quantity }
    });
  };

  const clearCart = () => {
    dispatch({ type: 'CLEAR_CART' });
  };

  const getCartItemCount = () => {
    return state.totalItems;
  };

  const getCartTotal = () => {
    return state.totalPrice;
  };

  const isInCart = (productId, variantId) => {
    return state.items.some(
      item => item.product.id === productId && item.variant.id === variantId
    );
  };

  const getCartItem = (productId, variantId) => {
    return state.items.find(
      item => item.product.id === productId && item.variant.id === variantId
    );
  };

  const value = {
    cart: state,
    addToCart,
    removeFromCart,
    updateQuantity,
    clearCart,
    getCartItemCount,
    getCartTotal,
    isInCart,
    getCartItem
  };

  return (
    <CartContext.Provider value={value}>
      {children}
    </CartContext.Provider>
  );
};

export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};

export default CartContext;