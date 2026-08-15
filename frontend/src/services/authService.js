import { api } from './api';

export const authService = {
  async register(data) {
    const res = await api.post('/auth/register', data);
    return res.data;
  },

  async login(data) {
    const res = await api.post('/auth/login', data);
    return res.data;
  },

  async getProfile() {
    const res = await api.get('/auth/profile');
    return res.data;
  },

  async getAddresses() {
    const res = await api.get('/user/addresses');
    return res.data;
  },

  async addAddress(data) {
    const res = await api.post('/user/addresses', data);
    return res.data;
  },
};
