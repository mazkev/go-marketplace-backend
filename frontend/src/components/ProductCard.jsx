import React from 'react';
import { Link } from 'react-router-dom';
import { Star, MapPin, Heart, ShieldCheck } from 'lucide-react';
import { useWishlist } from '../context/WishlistContext';
import { useAuth } from '../context/AuthContext';

export const ProductCard = ({ product }) => {
  const { isWishlisted, toggleWishlist } = useWishlist();
  const { isAuthenticated, openLoginModal } = useAuth();

  const isSaved = isWishlisted(product.id);

  const handleWishlistClick = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (!isAuthenticated) {
      openLoginModal();
      return;
    }
    toggleWishlist(product.id);
  };

  const formatRupiah = (val) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0,
    }).format(val || 0);
  };

  return (
    <Link
      to={`/products/${product.id}`}
      className="card"
      style={{
        display: 'flex',
        flexDirection: 'column',
        textDecoration: 'none',
        overflow: 'hidden',
        transition: 'transform 0.2s, box-shadow 0.2s',
        position: 'relative',
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = 'translateY(-4px)';
        e.currentTarget.style.boxShadow = 'var(--shadow-md)';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = 'translateY(0)';
        e.currentTarget.style.boxShadow = 'var(--shadow-sm)';
      }}
    >
      {/* Product Image Placeholder with Badge */}
      <div
        style={{
          width: '100%',
          paddingTop: '100%',
          position: 'relative',
          background: 'linear-gradient(135deg, #F0FDF4 0%, #DCFCE7 100%)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div
          style={{
            position: 'absolute',
            inset: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: '36px',
          }}
        >
          🛍️
        </div>

        {/* Wishlist Button */}
        <button
          onClick={handleWishlistClick}
          style={{
            position: 'absolute',
            top: '8px',
            right: '8px',
            background: 'rgba(255, 255, 255, 0.9)',
            border: 'none',
            borderRadius: '50%',
            width: '32px',
            height: '32px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
            color: isSaved ? 'var(--color-accent-red)' : '#8D96AA',
          }}
        >
          <Heart size={16} fill={isSaved ? 'var(--color-accent-red)' : 'none'} />
        </button>

        {/* Free Shipping / Official Badge */}
        <div style={{ position: 'absolute', bottom: '8px', left: '8px', display: 'flex', gap: '4px' }}>
          <span className="badge badge-green">Bebas Ongkir</span>
        </div>
      </div>

      {/* Product Info */}
      <div style={{ padding: '12px', display: 'flex', flexDirection: 'column', flex: 1, justifyContent: 'space-between' }}>
        <div>
          {/* Title */}
          <div
            style={{
              fontSize: '13px',
              fontWeight: 500,
              color: 'var(--text-primary)',
              lineHeight: '1.4',
              height: '36px',
              overflow: 'hidden',
              display: '-webkit-box',
              WebkitLineClamp: 2,
              WebkitBoxOrient: 'vertical',
              marginBottom: '6px',
            }}
          >
            {product.name}
          </div>

          {/* Price */}
          <div style={{ fontSize: '15px', fontWeight: 800, color: 'var(--text-primary)', marginBottom: '4px' }}>
            {formatRupiah(product.price)}
          </div>

          {/* Discount tag mockup */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '4px', marginBottom: '8px' }}>
            <span style={{ fontSize: '10px', fontWeight: 700, color: 'var(--color-accent-red)', background: 'var(--color-accent-red-light)', padding: '1px 4px', borderRadius: '4px' }}>
              10%
            </span>
            <span style={{ fontSize: '11px', color: 'var(--text-muted)', textDecoration: 'line-through' }}>
              {formatRupiah(product.price * 1.1)}
            </span>
          </div>
        </div>

        <div>
          {/* Store Info */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '4px', fontSize: '11px', color: 'var(--text-secondary)', marginBottom: '4px' }}>
            <ShieldCheck size={13} color="var(--color-primary)" />
            <span style={{ fontWeight: 600, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
              {product.store?.store_name || 'Official Store'}
            </span>
          </div>

          {/* Location & Rating */}
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', fontSize: '11px', color: 'var(--text-muted)' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '3px' }}>
              <MapPin size={11} />
              <span>Kota Jakarta</span>
            </div>

            <div style={{ display: 'flex', alignItems: 'center', gap: '3px' }}>
              <Star size={12} fill="#FFC400" color="#FFC400" />
              <span style={{ fontWeight: 700, color: 'var(--text-primary)' }}>
                {product.rating_avg ? product.rating_avg.toFixed(1) : '5.0'}
              </span>
              <span>({product.rating_count || 1})</span>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
};
