import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Star,
  MapPin,
  ShieldCheck,
  Heart,
  ShoppingCart,
  CheckCircle,
  AlertCircle,
  Truck,
  MessageCircle,
} from 'lucide-react';
import { productService } from '../services/productService';
import { useCart } from '../context/CartContext';
import { useWishlist } from '../context/WishlistContext';
import { useAuth } from '../context/AuthContext';

export const ProductDetailPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { addToCart } = useCart();
  const { isWishlisted, toggleWishlist } = useWishlist();
  const { isAuthenticated, openLoginModal } = useAuth();

  const [product, setProduct] = useState(null);
  const [reviews, setReviews] = useState([]);
  const [selectedVariant, setSelectedVariant] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [loading, setLoading] = useState(true);
  const [actionSuccess, setActionSuccess] = useState('');
  const [actionError, setActionError] = useState('');

  useEffect(() => {
    fetchProductDetails();
  }, [id]);

  const fetchProductDetails = async () => {
    setLoading(true);
    try {
      const [prodRes, revRes] = await Promise.all([
        productService.getProductById(id),
        productService.getProductReviews(id),
      ]);
      if (prodRes.data) {
        setProduct(prodRes.data);
        if (prodRes.data.variants && prodRes.data.variants.length > 0) {
          setSelectedVariant(prodRes.data.variants[0]);
        }
      }
      if (revRes.data) setReviews(revRes.data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="container" style={{ padding: '60px 16px', textAlign: 'center' }}>
        Memuat detail produk...
      </div>
    );
  }

  if (!product) {
    return (
      <div className="container" style={{ padding: '60px 16px', textAlign: 'center' }}>
        Produk tidak ditemukan.
      </div>
    );
  }

  const currentPrice = selectedVariant?.price_override != null ? selectedVariant.price_override : product.price;
  const currentStock = selectedVariant ? selectedVariant.stock : product.stock;
  const isSaved = isWishlisted(product.id);

  const formatRupiah = (val) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0,
    }).format(val || 0);
  };

  const handleAddToCart = async () => {
    if (!isAuthenticated) {
      openLoginModal();
      return;
    }
    setActionError('');
    setActionSuccess('');
    try {
      await addToCart(product.id, quantity, selectedVariant?.id);
      setActionSuccess('Produk berhasil ditambahkan ke keranjang!');
      setTimeout(() => setActionSuccess(''), 3000);
    } catch (err) {
      setActionError(err.response?.data?.message || 'Gagal menambahkan ke keranjang');
    }
  };

  const handleBuyNow = async () => {
    if (!isAuthenticated) {
      openLoginModal();
      return;
    }
    try {
      await addToCart(product.id, quantity, selectedVariant?.id);
      navigate('/cart');
    } catch (err) {
      setActionError(err.response?.data?.message || 'Gagal memproses pesanan');
    }
  };

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      {/* Alerts */}
      {actionSuccess && (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '12px 16px', background: 'var(--color-primary-light)', color: 'var(--color-primary-dark)', borderRadius: '8px', marginBottom: '16px', fontWeight: 600, fontSize: '14px' }}>
          <CheckCircle size={18} />
          <span>{actionSuccess}</span>
        </div>
      )}
      {actionError && (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '12px 16px', background: 'var(--color-accent-red-light)', color: 'var(--color-accent-red)', borderRadius: '8px', marginBottom: '16px', fontWeight: 600, fontSize: '14px' }}>
          <AlertCircle size={18} />
          <span>{actionError}</span>
        </div>
      )}

      {/* Main Grid: Left Gallery | Center Info | Right Purchase Box */}
      <div style={{ display: 'grid', gridTemplateColumns: '380px 1fr 300px', gap: '32px', alignItems: 'start' }}>
        {/* Left Column: Image Gallery */}
        <div>
          <div
            className="card"
            style={{
              width: '100%',
              height: '380px',
              borderRadius: '16px',
              background: 'linear-gradient(135deg, #F0FDF4 0%, #DCFCE7 100%)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: '96px',
              overflow: 'hidden',
              marginBottom: '16px',
            }}
          >
            🛍️
          </div>
          <div style={{ display: 'flex', gap: '12px' }}>
            {[1, 2, 3].map((_, i) => (
              <div
                key={i}
                style={{
                  width: '64px',
                  height: '64px',
                  borderRadius: '8px',
                  border: i === 0 ? '2px solid var(--color-primary)' : '1px solid var(--border-light)',
                  background: '#F9FAFB',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '24px',
                  cursor: 'pointer',
                }}
              >
                🛍️
              </div>
            ))}
          </div>
        </div>

        {/* Center Column: Product Details & Store */}
        <div>
          <h1 style={{ fontSize: '20px', fontWeight: 800, lineHeight: '1.4', marginBottom: '12px' }}>
            {product.name}
          </h1>

          {/* Rating & Sold Stats */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '16px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
              <Star size={16} fill="#FFC400" color="#FFC400" />
              <span style={{ fontWeight: 700, color: 'var(--text-primary)' }}>
                {product.rating_avg ? product.rating_avg.toFixed(1) : '5.0'}
              </span>
              <span>({product.rating_count || 1} ulasan)</span>
            </div>
            <span>•</span>
            <span>Terjual 50+</span>
            <span>•</span>
            <span style={{ color: 'var(--color-primary)', fontWeight: 600 }}>Kategori: {product.category?.name || 'Produk'}</span>
          </div>

          {/* Price */}
          <div style={{ fontSize: '28px', fontWeight: 900, color: 'var(--text-primary)', marginBottom: '16px' }}>
            {formatRupiah(currentPrice)}
          </div>

          {/* Variant Selector */}
          {product.variants && product.variants.length > 0 && (
            <div style={{ marginBottom: '24px', borderTop: '1px solid var(--border-light)', paddingTop: '16px' }}>
              <div style={{ fontSize: '13px', fontWeight: 700, marginBottom: '8px' }}>
                Pilih Varian: <span style={{ color: 'var(--color-primary)' }}>{selectedVariant?.variant_name}</span>
              </div>
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                {product.variants.map((v) => {
                  const isSelected = selectedVariant?.id === v.id;
                  return (
                    <button
                      key={v.id}
                      onClick={() => setSelectedVariant(v)}
                      style={{
                        padding: '8px 16px',
                        borderRadius: '8px',
                        border: isSelected ? '2px solid var(--color-primary)' : '1px solid var(--border-medium)',
                        background: isSelected ? 'var(--color-primary-light)' : 'white',
                        color: isSelected ? 'var(--color-primary-dark)' : 'var(--text-primary)',
                        fontWeight: 600,
                        fontSize: '13px',
                      }}
                    >
                      {v.variant_name}
                    </button>
                  );
                })}
              </div>
            </div>
          )}

          {/* Description */}
          <div style={{ borderTop: '1px solid var(--border-light)', paddingTop: '16px', marginBottom: '24px' }}>
            <div style={{ fontSize: '14px', fontWeight: 700, marginBottom: '8px' }}>Deskripsi Produk</div>
            <p style={{ fontSize: '14px', lineHeight: '1.6', color: 'var(--text-secondary)', whiteSpace: 'pre-line' }}>
              {product.description || 'Tidak ada deskripsi untuk produk ini.'}
            </p>
          </div>

          {/* Store Info Card */}
          <div className="card" style={{ padding: '16px', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <div style={{ width: '48px', height: '48px', borderRadius: '50%', background: 'var(--color-primary-light)', color: 'var(--color-primary-dark)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 800, fontSize: '20px' }}>
                🏬
              </div>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '6px', fontWeight: 700, fontSize: '14px' }}>
                  <ShieldCheck size={16} color="var(--color-primary)" />
                  <span>{product.store?.store_name || 'Official Seller'}</span>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '4px', fontSize: '12px', color: 'var(--text-secondary)' }}>
                  <MapPin size={12} />
                  <span>Kota Jakarta • Online</span>
                </div>
              </div>
            </div>
            <button className="btn btn-outline btn-sm">
              <MessageCircle size={14} /> Chat Toko
            </button>
          </div>
        </div>

        {/* Right Column: Sticky Purchase Box */}
        <div className="card" style={{ padding: '20px', position: 'sticky', top: '100px' }}>
          <div style={{ fontWeight: 800, fontSize: '15px', marginBottom: '16px' }}>Atur Jumlah & Catatan</div>

          {/* Quantity Controls */}
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '16px' }}>
            <div style={{ display: 'flex', alignItems: 'center', border: '1px solid var(--border-medium)', borderRadius: '8px' }}>
              <button
                disabled={quantity <= 1}
                onClick={() => setQuantity(Math.max(1, quantity - 1))}
                style={{ padding: '6px 12px', fontWeight: 700, fontSize: '16px' }}
              >
                -
              </button>
              <span style={{ padding: '0 8px', fontWeight: 600, fontSize: '14px' }}>{quantity}</span>
              <button
                disabled={quantity >= currentStock}
                onClick={() => setQuantity(Math.min(currentStock, quantity + 1))}
                style={{ padding: '6px 12px', fontWeight: 700, fontSize: '16px' }}
              >
                +
              </button>
            </div>
            <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>
              Sisa Stok: <span style={{ fontWeight: 700, color: 'var(--text-primary)' }}>{currentStock}</span>
            </div>
          </div>

          {/* Subtotal */}
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', borderTop: '1px solid var(--border-light)', paddingTop: '12px' }}>
            <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>Subtotal</span>
            <span style={{ fontSize: '18px', fontWeight: 900, color: 'var(--text-primary)' }}>
              {formatRupiah(currentPrice * quantity)}
            </span>
          </div>

          {/* Buttons */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <button
              onClick={handleAddToCart}
              className="btn btn-primary btn-block btn-lg"
            >
              <ShoppingCart size={18} /> + Keranjang
            </button>
            <button
              onClick={handleBuyNow}
              className="btn btn-outline btn-block"
            >
              Beli Langsung
            </button>
            <button
              onClick={() => {
                if (!isAuthenticated) openLoginModal();
                else toggleWishlist(product.id);
              }}
              className="btn btn-secondary btn-block btn-sm"
              style={{ color: isSaved ? 'var(--color-accent-red)' : 'var(--text-secondary)' }}
            >
              <Heart size={14} fill={isSaved ? 'var(--color-accent-red)' : 'none'} />
              {isSaved ? 'Hapus dari Wishlist' : 'Tambah ke Wishlist'}
            </button>
          </div>
        </div>
      </div>

      {/* Customer Reviews Section */}
      <div className="card" style={{ marginTop: '48px', padding: '24px' }}>
        <h3 style={{ fontSize: '18px', fontWeight: 800, marginBottom: '20px' }}>
          Ulasan Pembeli ({reviews.length})
        </h3>
        {reviews.length === 0 ? (
          <div style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>
            Belum ada ulasan untuk produk ini.
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            {reviews.map((rev) => (
              <div key={rev.id} style={{ borderBottom: '1px solid var(--border-light)', paddingBottom: '16px' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '6px' }}>
                  <div style={{ display: 'flex' }}>
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        size={14}
                        fill={i < rev.rating ? '#FFC400' : '#E0E0E0'}
                        color={i < rev.rating ? '#FFC400' : '#E0E0E0'}
                      />
                    ))}
                  </div>
                  <span style={{ fontWeight: 700, fontSize: '13px' }}>{rev.user?.name || 'Pembeli'}</span>
                  <span style={{ fontSize: '11px', color: 'var(--text-muted)' }}>
                    {new Date(rev.created_at).toLocaleDateString('id-ID')}
                  </span>
                </div>
                <p style={{ fontSize: '13px', color: 'var(--text-secondary)', lineHeight: '1.5' }}>
                  {rev.comment}
                </p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
