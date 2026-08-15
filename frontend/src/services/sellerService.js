import { api } from './api';

export const sellerService = {
  async registerStore(data) {
    const res = await api.post('/stores', data);
    return res.data;
  },

  async getMyStore() {
    const res = await api.get('/seller/store');
    return res.data;
  },

  async getStoreBalance() {
    const res = await api.get('/seller/balance');
    return res.data;
  },

  async createProduct(data) {
    const res = await api.post('/seller/products', data);
    return res.data;
  },

  async updateProduct(id, data) {
    const res = await api.put(`/seller/products/${id}`, data);
    return res.data;
  },

  async deleteProduct(id) {
    const res = await api.delete(`/seller/products/${id}`);
    return res.data;
  },

  async getStoreOrders() {
    const res = await api.get('/seller/orders');
    return res.data;
  },

  async shipOrderItem(orderItemId, data) {
    const res = await api.patch(`/seller/orders/${orderItemId}/ship`, data);
    return res.data;
  },

  async createStoreVoucher(data) {
    const res = await api.post('/seller/vouchers', data);
    return res.data;
  },

  async requestWithdrawal(data) {
    const res = await api.post('/seller/withdrawals', data);
    return res.data;
  },

  async getWithdrawals() {
    const res = await api.get('/seller/withdrawals');
    return res.data;
  },

  async getMutations() {
    const res = await api.get('/seller/mutations');
    return res.data;
  },
};
