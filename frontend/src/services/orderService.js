import { api } from './api';

export const orderService = {
  async checkout(payload) {
    const res = await api.post('/orders/checkout', payload);
    return res.data;
  },

  async getUserOrders() {
    const res = await api.get('/orders');
    return res.data;
  },

  async getOrderById(orderId) {
    const res = await api.get(`/orders/${orderId}`);
    return res.data;
  },

  async completeOrderItem(orderItemId) {
    const res = await api.post(`/orders/items/${orderItemId}/complete`);
    return res.data;
  },

  async createReview(data) {
    const res = await api.post('/reviews', data);
    return res.data;
  },

  async simulatePaymentWebhook(invoiceNumber, amount) {
    const res = await api.post('/payments/webhook', {
      invoice_number: invoiceNumber,
      amount: amount,
      payment_status: 'SETTLEMENT',
    });
    return res.data;
  },
};
