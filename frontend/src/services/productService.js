import { api } from './api';

export const productService = {
  async getCategories() {
    const res = await api.get('/categories');
    return res.data;
  },

  async getProducts(params = {}) {
    const res = await api.get('/products', { params });
    return res.data;
  },

  async getProductById(id) {
    const res = await api.get(`/products/${id}`);
    return res.data;
  },

  async getProductReviews(id) {
    const res = await api.get(`/products/${id}/reviews`);
    return res.data;
  },

  async getStoreProfile(id) {
    const res = await api.get(`/stores/${id}`);
    return res.data;
  },
};
