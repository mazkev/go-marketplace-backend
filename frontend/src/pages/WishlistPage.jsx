import React from 'react';
import { Link } from 'react-router-dom';
import { Heart, Trash2 } from 'lucide-react';
import { useWishlist } from '../context/WishlistContext';
import { useAuth } from '../context/AuthContext';
import { ProductCard } from '../components/ProductCard';

export const WishlistPage = () => {
  const { wishlist, loading } = useWishlist();
  const { isAuthenticated, openLoginModal } = useAuth();

  if (!isAuthenticated) {
    return (
      <div className="container" style={{ padding: '60px 16px', textAlign: 'center' }}>
        <div className="card" style={{ maxWidth: '480px', margin: '0 auto', padding: '40px 24px' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>❤️</div>
          <h2 style={{ fontSize: '20px', fontWeight: 800, marginBottom: '8px' }}>Wishlist Saya</h2>
          <p style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '24px' }}>
            Masuk untuk melihat daftar produk favorit yang Anda simpan.
          </p>
          <button onClick={openLoginModal} className="btn btn-primary btn-block btn-lg">
            Masuk ke Akun
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '20px' }}>
        <Heart size={24} color="var(--color-accent-red)" fill="var(--color-accent-red)" />
        <h1 style={{ fontSize: '22px', fontWeight: 800 }}>Wishlist Saya ({wishlist.length})</h1>
      </div>

      {loading ? (
        <div style={{ textAlign: 'center', padding: '40px 0' }}>Memuat wishlist...</div>
      ) : wishlist.length === 0 ? (
        <div className="card" style={{ padding: '60px 20px', textAlign: 'center' }}>
          <div style={{ fontSize: '64px', marginBottom: '16px' }}>💔</div>
          <h3 style={{ fontSize: '18px', fontWeight: 800, marginBottom: '8px' }}>
            Belum ada barang di wishlist
          </h3>
          <p style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '24px' }}>
            Temukan barang yang kamu suka dan simpan ke wishlist untuk dibeli nanti!
          </p>
          <Link to="/" className="btn btn-primary btn-lg">
            Cari Produk Menarik
          </Link>
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))', gap: '16px' }}>
          {wishlist.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}
    </div>
  );
};
