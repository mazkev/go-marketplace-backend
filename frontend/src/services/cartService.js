import { api } from './api';

export const cartService = {
  async getCart() {
    const res = await api.get('/cart');
    return res.data;
  },

  async addToCart(productId, quantity = 1, variantId = null) {
    const payload = {
      product_id: productId,
      quantity,
    };
    if (variantId) {
      payload.variant_id = variantId;
    }
    const res = await api.post('/cart/items', payload);
    return res.data;
  },

  async updateQuantity(cartId, quantity) {
    const res = await api.put(`/cart/items/${cartId}`, { quantity });
    return res.data;
  },

  async deleteItem(cartId) {
    const res = await api.delete(`/cart/items/${cartId}`);
    return res.data;
  },
};
