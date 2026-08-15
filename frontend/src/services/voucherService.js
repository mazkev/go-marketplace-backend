import { api } from './api';

export const voucherService = {
  async getAvailableVouchers(storeId = null) {
    const params = {};
    if (storeId) params.store_id = storeId;
    const res = await api.get('/vouchers', { params });
    return res.data;
  },

  async applyVoucher(code, totalAmount, storeId = null) {
    const payload = {
      code,
      total_amount: totalAmount,
    };
    if (storeId) payload.store_id = storeId;
    const res = await api.post('/vouchers/apply', payload);
    return res.data;
  },
};
