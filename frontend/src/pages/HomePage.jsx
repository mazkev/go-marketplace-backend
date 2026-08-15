import React, { useState, useEffect } from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import { Sparkles, Tag, Flame, CheckCircle, Smartphone, Laptop, Shirt, ShoppingBag, Coffee, ChevronRight } from 'lucide-react';
import { productService } from '../services/productService';
import { voucherService } from '../services/voucherService';
import { ProductCard } from '../components/ProductCard';

export const HomePage = () => {
  const [searchParams] = useSearchParams();
  const search = searchParams.get('search') || '';
  const categoryId = searchParams.get('category_id') || '';

  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [vouchers, setVouchers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [copiedCode, setCopiedCode] = useState('');

  useEffect(() => {
    fetchData();
  }, [search, categoryId]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [prodRes, catRes, vouchRes] = await Promise.all([
        productService.getProducts({ search, category_id: categoryId || undefined }),
        productService.getCategories(),
        voucherService.getAvailableVouchers(),
      ]);

      if (prodRes.data) setProducts(prodRes.data);
      if (catRes.data) setCategories(catRes.data);
      if (vouchRes.data) setVouchers(vouchRes.data);
    } catch (err) {
      console.error('Error fetching home data:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleCopyVoucher = (code) => {
    navigator.clipboard.writeText(code);
    setCopiedCode(code);
    setTimeout(() => setCopiedCode(''), 2000);
  };

  const categoryIcons = {
    'Elektronik': <Smartphone size={24} color="#03AC0E" />,
    'Fashion Pria': <Shirt size={24} color="#FA591D" />,
    'Fashion Wanita': <ShoppingBag size={24} color="#E02954" />,
    'Handphone & Tablet': <Smartphone size={24} color="#0284C7" />,
    'Komputer & Laptop': <Laptop size={24} color="#7C3AED" />,
    'Makanan & Minuman': <Coffee size={24} color="#D97706" />,
  };

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      {/* Hero Banner Section */}
      {!search && !categoryId && (
        <>
          <div
            style={{
              background: 'linear-gradient(135deg, #03AC0E 0%, #008844 100%)',
              borderRadius: '16px',
              padding: '36px 48px',
              color: 'white',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              marginBottom: '32px',
              boxShadow: 'var(--shadow-md)',
              position: 'relative',
              overflow: 'hidden',
            }}
          >
            <div style={{ maxWidth: '600px', zIndex: 1 }}>
              <div style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', background: 'rgba(255, 255, 255, 0.2)', padding: '4px 12px', borderRadius: '20px', fontSize: '12px', fontWeight: 700, marginBottom: '12px' }}>
                <Sparkles size={14} /> PROMO SPESIAL TOKOMARKET
              </div>
              <h1 style={{ fontSize: '32px', fontWeight: 800, lineHeight: '1.2', marginBottom: '12px' }}>
                Belanja Hemat & Transaksi Aman Bersama TokoMarket
              </h1>
              <p style={{ fontSize: '15px', opacity: 0.95, marginBottom: '24px' }}>
                Nikmati multi-store checkout, gratis ongkir per toko, dan proteksi saldo escrow 100%.
              </p>
              <div style={{ display: 'flex', gap: '12px' }}>
                <a href="#product-section" className="btn btn-orange btn-lg">
                  Mulai Belanja
                </a>
                <Link to="/seller" className="btn btn-outline btn-lg" style={{ color: 'white', borderColor: 'white' }}>
                  Buka Toko Gratis
                </Link>
              </div>
            </div>
            <div style={{ fontSize: '110px', opacity: 0.85, transform: 'rotate(10deg)' }}>
              🛍️
            </div>
          </div>

          {/* Category Icon Grid */}
          <div className="card" style={{ padding: '20px', marginBottom: '32px' }}>
            <div style={{ fontWeight: 800, fontSize: '16px', marginBottom: '16px', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <span>Kategori Pilihan</span>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))', gap: '16px' }}>
              {categories.map((cat) => (
                <Link
                  key={cat.id}
                  to={`/?category_id=${cat.id}`}
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    padding: '16px 8px',
                    borderRadius: '12px',
                    border: '1px solid var(--border-light)',
                    background: '#FAFAFA',
                    transition: 'all 0.2s',
                    textAlign: 'center',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.borderColor = 'var(--color-primary)';
                    e.currentTarget.style.transform = 'translateY(-2px)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.borderColor = 'var(--border-light)';
                    e.currentTarget.style.transform = 'translateY(0)';
                  }}
                >
                  <div style={{ marginBottom: '8px' }}>
                    {categoryIcons[cat.name] || <Tag size={24} color="#03AC0E" />}
                  </div>
                  <div style={{ fontSize: '13px', fontWeight: 600, color: 'var(--text-primary)' }}>
                    {cat.name}
                  </div>
                </Link>
              ))}
            </div>
          </div>

          {/* Promo & Voucher Banner Bar */}
          {vouchers.length > 0 && (
            <div className="card" style={{ padding: '20px', marginBottom: '32px', background: 'linear-gradient(to right, #FFF7ED, #FEF2F2)' }}>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '16px' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontWeight: 800, fontSize: '16px', color: 'var(--color-secondary-orange)' }}>
                  <Flame size={20} />
                  <span>Klaim Voucher Diskon Spesial Hari Ini</span>
                </div>
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: '16px' }}>
                {vouchers.map((v) => (
                  <div
                    key={v.id}
                    style={{
                      background: 'white',
                      borderRadius: '10px',
                      border: '1px dashed var(--color-secondary-orange)',
                      padding: '14px 16px',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                      boxShadow: 'var(--shadow-xs)',
                    }}
                  >
                    <div>
                      <div style={{ fontWeight: 800, fontSize: '15px', color: 'var(--color-secondary-orange)' }}>
                        {v.code}
                      </div>
                      <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>
                        {v.voucher_type === 'percentage'
                          ? `Diskon ${v.discount_percent}% (Maks. Rp ${v.max_discount?.toLocaleString('id-ID')})`
                          : `Potongan Rp ${v.discount_amount?.toLocaleString('id-ID')}`}
                      </div>
                      <div style={{ fontSize: '11px', color: 'var(--text-muted)' }}>
                        Min. Belanja: Rp {v.min_spend?.toLocaleString('id-ID')}
                      </div>
                    </div>
                    <button
                      onClick={() => handleCopyVoucher(v.code)}
                      className="btn btn-outline btn-sm"
                      style={{ borderColor: 'var(--color-secondary-orange)', color: 'var(--color-secondary-orange)' }}
                    >
                      {copiedCode === v.code ? <CheckCircle size={14} /> : 'Salin'}
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </>
      )}

      {/* Main Product Grid Section */}
      <div id="product-section">
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '20px' }}>
          <div>
            <h2 style={{ fontSize: '20px', fontWeight: 800 }}>
              {search ? `Hasil Pencarian: "${search}"` : categoryId ? 'Produk Kategori' : 'Rekomendasi Untukmu'}
            </h2>
            <p style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>
              {products.length} produk tersedia dari penjual terpercaya
            </p>
          </div>
          {(search || categoryId) && (
            <Link to="/" className="btn btn-secondary btn-sm">
              Reset Filter
            </Link>
          )}
        </div>

        {loading ? (
          <div style={{ padding: '60px 0', textAlign: 'center', color: 'var(--text-secondary)' }}>
            Memuat produk marketplace...
          </div>
        ) : products.length === 0 ? (
          <div className="card" style={{ padding: '60px 20px', textAlign: 'center' }}>
            <div style={{ fontSize: '48px', marginBottom: '16px' }}>🔍</div>
            <div style={{ fontWeight: 700, fontSize: '16px', marginBottom: '8px' }}>Tidak Ada Produk Ditemukan</div>
            <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '20px' }}>
              Coba cari dengan kata kunci lain atau periksa kategori lainnya.
            </p>
            <Link to="/" className="btn btn-primary">
              Lihat Semua Produk
            </Link>
          </div>
        ) : (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))', gap: '16px' }}>
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
