import React, { useState } from 'react';
import { X, Lock, Mail, User, Phone, AlertCircle, ShoppingBag, Store } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

export const AuthModal = () => {
  const {
    isAuthModalOpen,
    setIsAuthModalOpen,
    authModalMode,
    setAuthModalMode,
    login,
    register,
  } = useAuth();

  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    phone: '',
  });
  const [role, setRole] = useState('buyer'); // 'buyer' | 'seller'
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  if (!isAuthModalOpen) return null;

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
    setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      if (authModalMode === 'login') {
        await login(formData.email, formData.password);
      } else {
        await register(formData.name, formData.email, formData.password, formData.phone, role);
      }
    } catch (err) {
      const errMsg = err.response?.data?.message || err.response?.data?.error || 'Terjadi kesalahan saat memproses';
      setError(errMsg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={() => setIsAuthModalOpen(false)}>
      <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        {/* Modal Header */}
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '16px 20px', borderBottom: '1px solid var(--border-light)' }}>
          <div style={{ fontWeight: 800, fontSize: '18px', color: 'var(--color-primary)' }}>
            {authModalMode === 'login' ? 'Masuk ke TokoMarket' : 'Daftar Akun Baru'}
          </div>
          <button onClick={() => setIsAuthModalOpen(false)} style={{ color: '#8D96AA' }}>
            <X size={20} />
          </button>
        </div>

        {/* Modal Body */}
        <div style={{ padding: '24px 20px' }}>
          {error && (
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '10px 14px', background: 'var(--color-accent-red-light)', color: 'var(--color-accent-red)', borderRadius: '8px', fontSize: '13px', marginBottom: '16px' }}>
              <AlertCircle size={16} />
              <span>{error}</span>
            </div>
          )}

          <form onSubmit={handleSubmit}>
            {authModalMode === 'register' && (
              <>
                {/* Role Selector Tabs */}
                <div style={{ marginBottom: '16px' }}>
                  <label className="form-label">Tipe Akun</label>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px' }}>
                    <button
                      type="button"
                      onClick={() => setRole('buyer')}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        gap: '6px',
                        padding: '10px 12px',
                        borderRadius: '8px',
                        border: role === 'buyer' ? '2px solid var(--color-primary)' : '1px solid var(--border-medium)',
                        background: role === 'buyer' ? 'var(--color-primary-light)' : 'white',
                        color: role === 'buyer' ? 'var(--color-primary-dark)' : 'var(--text-secondary)',
                        fontWeight: 700,
                        fontSize: '13px',
                      }}
                    >
                      <ShoppingBag size={16} /> Pembeli (Buyer)
                    </button>
                    <button
                      type="button"
                      onClick={() => setRole('seller')}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        gap: '6px',
                        padding: '10px 12px',
                        borderRadius: '8px',
                        border: role === 'seller' ? '2px solid var(--color-primary)' : '1px solid var(--border-medium)',
                        background: role === 'seller' ? 'var(--color-primary-light)' : 'white',
                        color: role === 'seller' ? 'var(--color-primary-dark)' : 'var(--text-secondary)',
                        fontWeight: 700,
                        fontSize: '13px',
                      }}
                    >
                      <Store size={16} /> Penjual (Seller)
                    </button>
                  </div>
                </div>

                <div className="form-group">
                  <label className="form-label">Nama Lengkap / Nama Toko</label>
                  <div style={{ position: 'relative' }}>
                    <input
                      type="text"
                      name="name"
                      required
                      placeholder={role === 'buyer' ? 'Contoh: Budi Santoso' : 'Contoh: Kevin Store'}
                      value={formData.name}
                      onChange={handleChange}
                      className="form-input"
                      style={{ paddingLeft: '38px' }}
                    />
                    <User size={16} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#8D96AA' }} />
                  </div>
                </div>
              </>
            )}

            <div className="form-group">
              <label className="form-label">Email</label>
              <div style={{ position: 'relative' }}>
                <input
                  type="email"
                  name="email"
                  required
                  placeholder="nama@email.com"
                  value={formData.email}
                  onChange={handleChange}
                  className="form-input"
                  style={{ paddingLeft: '38px' }}
                />
                <Mail size={16} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#8D96AA' }} />
              </div>
            </div>

            {authModalMode === 'register' && (
              <div className="form-group">
                <label className="form-label">Nomor Handphone (Opsional)</label>
                <div style={{ position: 'relative' }}>
                  <input
                    type="tel"
                    name="phone"
                    placeholder="081234567890"
                    value={formData.phone}
                    onChange={handleChange}
                    className="form-input"
                    style={{ paddingLeft: '38px' }}
                  />
                  <Phone size={16} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#8D96AA' }} />
                </div>
              </div>
            )}

            <div className="form-group">
              <label className="form-label">Kata Sandi</label>
              <div style={{ position: 'relative' }}>
                <input
                  type="password"
                  name="password"
                  required
                  placeholder="Minimal 6 karakter"
                  value={formData.password}
                  onChange={handleChange}
                  className="form-input"
                  style={{ paddingLeft: '38px' }}
                />
                <Lock size={16} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#8D96AA' }} />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="btn btn-primary btn-block btn-lg"
              style={{ marginTop: '8px' }}
            >
              {loading ? 'Memproses...' : authModalMode === 'login' ? 'Masuk' : `Daftar Sebagai ${role === 'seller' ? 'Penjual' : 'Pembeli'}`}
            </button>
          </form>

          {/* Switch Tab */}
          <div style={{ marginTop: '20px', textAlign: 'center', fontSize: '13px', color: 'var(--text-secondary)' }}>
            {authModalMode === 'login' ? (
              <>
                Belum punya akun TokoMarket?{' '}
                <button
                  type="button"
                  onClick={() => setAuthModalMode('register')}
                  style={{ color: 'var(--color-primary)', fontWeight: 700 }}
                >
                  Daftar
                </button>
              </>
            ) : (
              <>
                Sudah punya akun?{' '}
                <button
                  type="button"
                  onClick={() => setAuthModalMode('login')}
                  style={{ color: 'var(--color-primary)', fontWeight: 700 }}
                >
                  Masuk
                </button>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
