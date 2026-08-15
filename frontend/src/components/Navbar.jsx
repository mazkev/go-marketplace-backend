import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Search,
  ShoppingCart,
  Heart,
  Store,
  User as UserIcon,
  LogOut,
  Package,
  Menu,
  ChevronDown,
  Sparkles,
} from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';
import { useWishlist } from '../context/WishlistContext';
import { productService } from '../services/productService';

export const Navbar = () => {
  const { user, isAuthenticated, logout, openLoginModal, openRegisterModal } = useAuth();
  const { totalItems } = useCart();
  const { wishlistCount } = useWishlist();
  const navigate = useNavigate();

  const [searchQuery, setSearchQuery] = useState('');
  const [categories, setCategories] = useState([]);
  const [isCategoryOpen, setIsCategoryOpen] = useState(false);
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

  useEffect(() => {
    productService.getCategories()
      .then((res) => {
        if (res.data) setCategories(res.data);
      })
      .catch((err) => console.error(err));
  }, []);

  const handleSearch = (e) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/?search=${encodeURIComponent(searchQuery.trim())}`);
    } else {
      navigate('/');
    }
  };

  return (
    <header style={{ position: 'sticky', top: 0, zIndex: 100, background: '#ffffff', borderBottom: '1px solid var(--border-light)', boxShadow: 'var(--shadow-xs)' }}>
      {/* Top Header Bar */}
      <div style={{ background: '#F8F9FA', borderBottom: '1px solid #ECEEEF', fontSize: '12px', color: '#6D7588', padding: '4px 0' }}>
        <div className="container" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div style={{ display: 'flex', gap: '16px' }}>
            <span>Download TokoMarket App</span>
            <span style={{ color: '#E0E0E0' }}>|</span>
            <span>Mitra TokoMarket</span>
            <span style={{ color: '#E0E0E0' }}>|</span>
            <span>Pusat Edukasi Seller</span>
          </div>
          <div style={{ display: 'flex', gap: '16px' }}>
            <Link to="/seller" style={{ color: 'var(--color-primary)', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '4px' }}>
              <Store size={14} /> Toko Saya / Seller Center
            </Link>
            <span style={{ color: '#E0E0E0' }}>|</span>
            <span>Tentang TokoMarket</span>
            <span style={{ color: '#E0E0E0' }}>|</span>
            <span>Bantuan</span>
          </div>
        </div>
      </div>

      {/* Main Navigation Bar */}
      <div className="container" style={{ display: 'flex', alignItems: 'center', gap: '24px', padding: '12px 16px' }}>
        {/* Brand Logo */}
        <Link to="/" style={{ display: 'flex', alignItems: 'center', gap: '6px', textDecoration: 'none' }}>
          <div style={{ background: 'var(--color-primary)', color: 'white', padding: '6px 10px', borderRadius: '8px', fontWeight: 900, fontSize: '20px', letterSpacing: '-0.5px' }}>
            TM
          </div>
          <div>
            <div style={{ color: 'var(--color-primary)', fontWeight: 800, fontSize: '20px', lineHeight: '1.1' }}>
              tokomarket
            </div>
            <div style={{ fontSize: '10px', color: 'var(--color-secondary-orange)', fontWeight: 700 }}>
              Official C2C Store
            </div>
          </div>
        </Link>

        {/* Category Dropdown */}
        <div style={{ position: 'relative' }}>
          <button
            onClick={() => setIsCategoryOpen(!isCategoryOpen)}
            style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '14px', fontWeight: 600, color: 'var(--text-secondary)', padding: '8px 12px', borderRadius: '8px' }}
          >
            <Menu size={18} />
            <span>Kategori</span>
            <ChevronDown size={14} />
          </button>

          {isCategoryOpen && (
            <div
              style={{
                position: 'absolute',
                top: '100%',
                left: 0,
                width: '240px',
                background: 'white',
                boxShadow: 'var(--shadow-lg)',
                borderRadius: '8px',
                border: '1px solid var(--border-light)',
                padding: '8px 0',
                zIndex: 200,
              }}
            >
              {categories.map((cat) => (
                <Link
                  key={cat.id}
                  to={`/?category_id=${cat.id}`}
                  onClick={() => setIsCategoryOpen(false)}
                  style={{
                    display: 'block',
                    padding: '8px 16px',
                    fontSize: '13px',
                    color: 'var(--text-primary)',
                    transition: 'background 0.2s',
                  }}
                  onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--color-primary-light)'}
                  onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                >
                  {cat.name}
                </Link>
              ))}
            </div>
          )}
        </div>

        {/* Search Bar */}
        <form onSubmit={handleSearch} style={{ flex: 1, position: 'relative' }}>
          <input
            type="text"
            placeholder="Cari laptop gaming, hoodie unisex, aksesoris..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            style={{
              width: '100%',
              padding: '10px 42px 10px 16px',
              fontSize: '14px',
              border: '1px solid var(--border-medium)',
              borderRadius: '8px',
              outline: 'none',
            }}
          />
          <button
            type="submit"
            style={{
              position: 'absolute',
              right: '6px',
              top: '50%',
              transform: 'translateY(-50%)',
              background: 'var(--color-primary)',
              color: 'white',
              borderRadius: '6px',
              padding: '6px 10px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <Search size={16} />
          </button>
        </form>

        {/* Action Icons */}
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          {/* Wishlist */}
          <Link
            to="/wishlist"
            title="Wishlist"
            style={{ position: 'relative', padding: '8px', color: 'var(--text-primary)' }}
          >
            <Heart size={22} />
            {wishlistCount > 0 && (
              <span
                style={{
                  position: 'absolute',
                  top: '2px',
                  right: '2px',
                  background: 'var(--color-accent-red)',
                  color: 'white',
                  fontSize: '10px',
                  fontWeight: 700,
                  borderRadius: '10px',
                  padding: '1px 5px',
                  lineHeight: '1',
                }}
              >
                {wishlistCount}
              </span>
            )}
          </Link>

          {/* Cart */}
          <Link
            to="/cart"
            title="Keranjang Belanja"
            style={{ position: 'relative', padding: '8px', color: 'var(--text-primary)' }}
          >
            <ShoppingCart size={22} />
            {totalItems > 0 && (
              <span
                style={{
                  position: 'absolute',
                  top: '2px',
                  right: '2px',
                  background: 'var(--color-accent-red)',
                  color: 'white',
                  fontSize: '10px',
                  fontWeight: 700,
                  borderRadius: '10px',
                  padding: '1px 5px',
                  lineHeight: '1',
                }}
              >
                {totalItems}
              </span>
            )}
          </Link>

          {/* User Auth Section */}
          <div style={{ borderLeft: '1px solid var(--border-light)', paddingLeft: '16px' }}>
            {isAuthenticated ? (
              <div style={{ position: 'relative' }}>
                <button
                  onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '8px',
                    padding: '6px 12px',
                    borderRadius: '8px',
                    border: '1px solid var(--border-light)',
                    background: 'white',
                  }}
                >
                  <div style={{ width: '28px', height: '28px', borderRadius: '50%', background: 'var(--color-primary-light)', color: 'var(--color-primary-dark)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 700, fontSize: '13px' }}>
                    {user?.name ? user.name[0].toUpperCase() : 'U'}
                  </div>
                  <span style={{ fontSize: '13px', fontWeight: 600, maxWidth: '100px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {user?.name || 'User'}
                  </span>
                  <ChevronDown size={14} color="#6D7588" />
                </button>

                {isUserMenuOpen && (
                  <div
                    style={{
                      position: 'absolute',
                      right: 0,
                      top: '110%',
                      width: '200px',
                      background: 'white',
                      boxShadow: 'var(--shadow-lg)',
                      borderRadius: '8px',
                      border: '1px solid var(--border-light)',
                      padding: '8px 0',
                      zIndex: 200,
                    }}
                  >
                    <div style={{ padding: '8px 16px', borderBottom: '1px solid var(--border-light)' }}>
                      <div style={{ fontWeight: 700, fontSize: '13px' }}>{user?.name}</div>
                      <div style={{ fontSize: '11px', color: 'var(--text-secondary)' }}>{user?.email}</div>
                    </div>
                    <Link
                      to="/orders"
                      onClick={() => setIsUserMenuOpen(false)}
                      style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '8px 16px', fontSize: '13px', color: 'var(--text-primary)' }}
                    >
                      <Package size={16} /> Pesanan Saya
                    </Link>
                    <Link
                      to="/seller"
                      onClick={() => setIsUserMenuOpen(false)}
                      style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '8px 16px', fontSize: '13px', color: 'var(--color-primary)', fontWeight: 600 }}
                    >
                      <Store size={16} /> Dashboard Seller
                    </Link>
                    <div style={{ borderTop: '1px solid var(--border-light)', marginTop: '4px' }}>
                      <button
                        onClick={() => {
                          logout();
                          setIsUserMenuOpen(false);
                        }}
                        style={{ display: 'flex', alignItems: 'center', gap: '8px', width: '100%', padding: '8px 16px', fontSize: '13px', color: 'var(--color-accent-red)', textAlign: 'left' }}
                      >
                        <LogOut size={16} /> Keluar
                      </button>
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div style={{ display: 'flex', gap: '8px' }}>
                <button
                  onClick={openLoginModal}
                  className="btn btn-outline btn-sm"
                >
                  Masuk
                </button>
                <button
                  onClick={openRegisterModal}
                  className="btn btn-primary btn-sm"
                >
                  Daftar
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};
