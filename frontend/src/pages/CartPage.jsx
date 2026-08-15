import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Trash2, Store, ShoppingBag, ArrowRight } from 'lucide-react';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';

export const CartPage = () => {
  const { cart, loading, updateQuantity, deleteItem, totalPrice, totalItems } = useCart();
  const { isAuthenticated, openLoginModal } = useAuth();
  const navigate = useNavigate();

  if (!isAuthenticated) {
    return (
      <div className="container" style={{ padding: '60px 16px', textAlign: 'center' }}>
        <div className="card" style={{ maxWidth: '480px', margin: '0 auto', padding: '40px 24px' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>🔒</div>
          <h2 style={{ fontSize: '20px', fontWeight: 800, marginBottom: '8px' }}>Silakan Masuk Terlebih Dahulu</h2>
          <p style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '24px' }}>
            Masuk untuk melihat dan mengelola barang-barang di keranjang belanja Anda.
          </p>
          <button onClick={openLoginModal} className="btn btn-primary btn-block btn-lg">
            Masuk ke Akun
          </button>
        </div>
      </div>
    );
  }

  const formatRupiah = (val) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0,
    }).format(val || 0);
  };

  const handleCheckout = () => {
    navigate('/checkout');
  };

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      <h1 style={{ fontSize: '22px', fontWeight: 800, marginBottom: '20px' }}>
        Keranjang Belanja
      </h1>

      {loading ? (
        <div style={{ textAlign: 'center', padding: '40px 0' }}>Memuat keranjang belanja...</div>
      ) : !cart.stores || cart.stores.length === 0 ? (
        <div className="card" style={{ padding: '60px 20px', textAlign: 'center' }}>
          <div style={{ fontSize: '64px', marginBottom: '16px' }}>🛒</div>
          <h3 style={{ fontSize: '18px', fontWeight: 800, marginBottom: '8px' }}>
            Wah, keranjang belanjamu kosong
          </h3>
          <p style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '24px' }}>
            Yuk, isi dengan barang-barang impianmu sekarang!
          </p>
          <Link to="/" className="btn btn-primary btn-lg">
            Mulai Belanja
          </Link>
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 340px', gap: '32px', alignItems: 'start' }}>
          {/* Left Column: Multi-Store Grouped Cart Items */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
            {cart.stores.map((storeGroup) => (
              <div key={storeGroup.store_id} className="card" style={{ overflow: 'hidden' }}>
                {/* Store Header */}
                <div style={{ background: '#F8F9FA', padding: '12px 20px', borderBottom: '1px solid var(--border-light)', display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <Store size={18} color="var(--color-primary)" />
                  <span style={{ fontWeight: 700, fontSize: '14px' }}>{storeGroup.store_name}</span>
                  <span className="badge badge-green" style={{ marginLeft: 'auto' }}>Kota Jakarta</span>
                </div>

                {/* Store Items List */}
                <div style={{ padding: '16px 20px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  {storeGroup.items.map((item) => (
                    <div
                      key={item.id}
                      style={{
                        display: 'grid',
                        gridTemplateColumns: '72px 1fr auto',
                        gap: '16px',
                        alignItems: 'center',
                        borderBottom: '1px solid #F1F3F5',
                        paddingBottom: '16px',
                      }}
                    >
                      {/* Product Thumbnail Placeholder */}
                      <div
                        style={{
                          width: '72px',
                          height: '72px',
                          borderRadius: '8px',
                          background: 'linear-gradient(135deg, #F0FDF4 0%, #DCFCE7 100%)',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          fontSize: '28px',
                        }}
                      >
                        🛍️
                      </div>

                      {/* Product Details */}
                      <div>
                        <Link
                          to={`/products/${item.product_id}`}
                          style={{ fontSize: '14px', fontWeight: 600, color: 'var(--text-primary)', marginBottom: '4px', display: 'block' }}
                        >
                          {item.product_name}
                        </Link>
                        {item.variant_name && (
                          <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginBottom: '6px' }}>
                            Varian: <span style={{ fontWeight: 600, color: 'var(--text-secondary)' }}>{item.variant_name}</span>
                          </div>
                        )}
                        <div style={{ fontSize: '15px', fontWeight: 800, color: 'var(--text-primary)' }}>
                          {formatRupiah(item.price)}
                        </div>
                      </div>

                      {/* Quantity & Delete Actions */}
                      <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                        <button
                          onClick={() => deleteItem(item.id)}
                          title="Hapus Item"
                          style={{ color: '#8D96AA', padding: '6px' }}
                        >
                          <Trash2 size={18} />
                        </button>

                        <div style={{ display: 'flex', alignItems: 'center', border: '1px solid var(--border-medium)', borderRadius: '6px' }}>
                          <button
                            disabled={item.quantity <= 1}
                            onClick={() => updateQuantity(item.id, item.quantity - 1)}
                            style={{ padding: '4px 10px', fontWeight: 700 }}
                          >
                            -
                          </button>
                          <span style={{ padding: '0 8px', fontSize: '13px', fontWeight: 600 }}>
                            {item.quantity}
                          </span>
                          <button
                            onClick={() => updateQuantity(item.id, item.quantity + 1)}
                            style={{ padding: '4px 10px', fontWeight: 700 }}
                          >
                            +
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}

                  {/* Store Subtotal */}
                  <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '13px', color: 'var(--text-secondary)', paddingTop: '4px' }}>
                    <span>Subtotal Pesanan Toko</span>
                    <span style={{ fontWeight: 700, color: 'var(--text-primary)' }}>
                      {formatRupiah(storeGroup.subtotal)}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Right Column: Checkout Summary Sidebar */}
          <div className="card" style={{ padding: '20px', position: 'sticky', top: '100px' }}>
            <div style={{ fontWeight: 800, fontSize: '16px', marginBottom: '16px' }}>
              Ringkasan Belanja
            </div>

            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '14px', marginBottom: '8px', color: 'var(--text-secondary)' }}>
              <span>Total Item ({totalItems} barang)</span>
              <span style={{ fontWeight: 600, color: 'var(--text-primary)' }}>{formatRupiah(totalPrice)}</span>
            </div>

            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '14px', marginBottom: '16px', color: 'var(--text-secondary)' }}>
              <span>Estimasi Ongkir</span>
              <span style={{ fontWeight: 600, color: 'var(--color-primary)' }}>Dihitung saat checkout</span>
            </div>

            <div style={{ borderTop: '1px solid var(--border-light)', paddingTop: '16px', marginBottom: '20px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span style={{ fontSize: '15px', fontWeight: 700 }}>Total Harga</span>
              <span style={{ fontSize: '20px', fontWeight: 900, color: 'var(--color-primary)' }}>
                {formatRupiah(totalPrice)}
              </span>
            </div>

            <button
              onClick={handleCheckout}
              className="btn btn-primary btn-block btn-lg"
            >
              <span>Beli ({totalItems})</span>
              <ArrowRight size={18} />
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
