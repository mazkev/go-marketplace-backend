import React, { createContext, useContext, useState, useEffect } from 'react';
import { cartService } from '../services/cartService';
import { useAuth } from './AuthContext';

const CartContext = createContext(null);

export const CartProvider = ({ children }) => {
  const { isAuthenticated } = useAuth();
  const [cart, setCart] = useState({ stores: [], total_items: 0, total_price: 0 });
  const [loading, setLoading] = useState(false);

  const fetchCart = async () => {
    if (!isAuthenticated) {
      setCart({ stores: [], total_items: 0, total_price: 0 });
      return;
    }
    try {
      setLoading(true);
      const res = await cartService.getCart();
      if (res.data) {
        setCart(res.data);
      }
    } catch (err) {
      console.error('Error fetching cart:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCart();
  }, [isAuthenticated]);

  const addToCart = async (productId, quantity = 1, variantId = null) => {
    const res = await cartService.addToCart(productId, quantity, variantId);
    await fetchCart();
    return res;
  };

  const updateQuantity = async (cartId, quantity) => {
    const res = await cartService.updateQuantity(cartId, quantity);
    await fetchCart();
    return res;
  };

  const deleteItem = async (cartId) => {
    const res = await cartService.deleteItem(cartId);
    await fetchCart();
    return res;
  };

  return (
    <CartContext.Provider
      value={{
        cart,
        loading,
        fetchCart,
        addToCart,
        updateQuantity,
        deleteItem,
        totalItems: cart.total_items || 0,
        totalPrice: cart.total_price || 0,
      }}
    >
      {children}
    </CartContext.Provider>
  );
};

export const useCart = () => useContext(CartContext);
