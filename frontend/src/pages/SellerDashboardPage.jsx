import React, { useState, useEffect } from 'react';
import {
  Store,
  Wallet,
  Package,
  Truck,
  Tag,
  Plus,
  ArrowDownRight,
  ArrowUpRight,
  CheckCircle,
  AlertCircle,
  Clock,
} from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { sellerService } from '../services/sellerService';
import { productService } from '../services/productService';

export const SellerDashboardPage = () => {
  const { user, isAuthenticated, openLoginModal, openRegisterModal, logout } = useAuth();

  const [store, setStore] = useState(null);
  const [balance, setBalance] = useState(0);
  const [mutations, setMutations] = useState([]);
  const [orders, setOrders] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('overview'); // 'overview' | 'orders' | 'add_product' | 'vouchers'

  // Register Store State
  const [newStoreName, setNewStoreName] = useState('');
  const [newDomainSlug, setNewDomainSlug] = useState('');
  const [creatingStore, setCreatingStore] = useState(false);

  // Withdrawal State
  const [withdrawModalOpen, setWithdrawModalOpen] = useState(false);
  const [withdrawBank, setWithdrawBank] = useState('BCA');
  const [withdrawAccNo, setWithdrawAccNo] = useState('');
  const [withdrawAccHolder, setWithdrawAccHolder] = useState(user?.name || '');
  const [withdrawAmount, setWithdrawAmount] = useState('');
  const [withdrawing, setWithdrawing] = useState(false);

  // Ship Order State
  const [shipModalOpen, setShipModalOpen] = useState(false);
  const [selectedOrderItemId, setSelectedOrderItemId] = useState(null);
  const [trackingNumber, setTrackingNumber] = useState('');
  const [shippingCourier, setShippingCourier] = useState('JNE Reguler');

  // Add Product State
  const [productName, setProductName] = useState('');
  const [productDesc, setProductDesc] = useState('');
  const [productCategoryId, setProductCategoryId] = useState('');
  const [productPrice, setProductPrice] = useState('');
  const [productStock, setProductStock] = useState('');
  const [productWeight, setProductWeight] = useState('500');
  const [variants, setVariants] = useState([]);
  const [newVariantName, setNewVariantName] = useState('');
  const [newVariantStock, setNewVariantStock] = useState('');

  // Store Voucher State
  const [voucherCode, setVoucherCode] = useState('');
  const [voucherPercent, setVoucherPercent] = useState('10');
  const [voucherMaxDiscount, setVoucherMaxDiscount] = useState('50000');
  const [voucherMinSpend, setVoucherMinSpend] = useState('100000');
  const [voucherQuota, setVoucherQuota] = useState('50');

  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    if (isAuthenticated && user?.role === 'seller') {
      fetchSellerData();
    } else {
      setLoading(false);
    }
  }, [isAuthenticated, user]);

  const fetchSellerData = async () => {
    setLoading(true);
    try {
      const [storeRes, balRes, mutRes, ordRes, catRes] = await Promise.allSettled([
        sellerService.getMyStore(),
        sellerService.getStoreBalance(),
        sellerService.getMutations(),
        sellerService.getStoreOrders(),
        productService.getCategories(),
      ]);

      if (storeRes.status === 'fulfilled' && storeRes.value?.data) {
        setStore(storeRes.value.data);
      }
      if (balRes.status === 'fulfilled' && balRes.value?.data) {
        setBalance(balRes.value.data.balance || 0);
      }
      if (mutRes.status === 'fulfilled' && mutRes.value?.data) {
        setMutations(mutRes.value.data);
      }
      if (ordRes.status === 'fulfilled' && ordRes.value?.data) {
        setOrders(ordRes.value.data);
      }
      if (catRes.status === 'fulfilled' && catRes.value?.data) {
        setCategories(catRes.value.data);
        if (catRes.value.data.length > 0) {
          setProductCategoryId(catRes.value.data[0].id);
        }
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateStore = async (e) => {
    e.preventDefault();
    setCreatingStore(true);
    setError('');
    try {
      const res = await sellerService.registerStore({
        store_name: newStoreName,
        domain_slug: newDomainSlug,
        city_id: 1,
      });
      if (res.data) {
        setStore(res.data);
        setMessage('Selamat! Toko Anda berhasil dibuka di TokoMarket!');
        await fetchSellerData();
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Gagal membuka toko');
    } finally {
      setCreatingStore(false);
    }
  };

  const handleWithdrawal = async (e) => {
    e.preventDefault();
    setWithdrawing(true);
    setError('');
    try {
      await sellerService.requestWithdrawal({
        bank_name: withdrawBank,
        account_number: withdrawAccNo,
        account_holder: withdrawAccHolder,
        amount: parseFloat(withdrawAmount),
      });
      setWithdrawModalOpen(false);
      setMessage('Pengajuan penarikan dana berhasil diproses!');
      setWithdrawAmount('');
      setWithdrawAccNo('');
      await fetchSellerData();
    } catch (err) {
      setError(err.response?.data?.message || 'Gagal melakukan penarikan dana');
    } finally {
      setWithdrawing(false);
    }
  };

  const handleShipOrder = async (e) => {
    e.preventDefault();
    try {
      await sellerService.shipOrderItem(selectedOrderItemId, {
        tracking_number: trackingNumber,
        courier_name: shippingCourier,
      });
      setShipModalOpen(false);
      setTrackingNumber('');
      setMessage('Resi pengiriman berhasil diinput! Status pesanan kini DIKIRIM.');
      await fetchSellerData();
    } catch (err) {
      alert(err.response?.data?.message || 'Gagal menginput resi');
    }
  };

  const handleAddVariant = () => {
    if (!newVariantName || !newVariantStock) return;
    setVariants([
      ...variants,
      {
        variant_name: newVariantName,
        stock: parseInt(newVariantStock, 10),
      },
    ]);
    setNewVariantName('');
    setNewVariantStock('');
  };

  const handleCreateProduct = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await sellerService.createProduct({
        category_id: parseInt(productCategoryId, 10),
        name: productName,
        description: productDesc,
        price: parseFloat(productPrice),
        stock: parseInt(productStock, 10) || 10,
        weight: parseInt(productWeight, 10) || 500,
        variants: variants.length > 0 ? variants : undefined,
      });
      setMessage('Produk berhasil ditambahkan ke etalase toko!');
      setProductName('');
      setProductDesc('');
      setProductPrice('');
      setProductStock('');
      setVariants([]);
      setActiveTab('overview');
    } catch (err) {
      setError(err.response?.data?.message || 'Gagal menambahkan produk');
    }
  };

  const handleCreateStoreVoucher = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const now = new Date();
      const nextMonth = new Date();
      nextMonth.setMonth(nextMonth.getMonth() + 1);

      await sellerService.createStoreVoucher({
        code: voucherCode.toUpperCase(),
        voucher_type: 'percentage',
        discount_percent: parseFloat(voucherPercent),
        max_discount: parseFloat(voucherMaxDiscount),
        min_spend: parseFloat(voucherMinSpend),
        quota: parseInt(voucherQuota, 10),
        start_date: now.toISOString(),
        end_date: nextMonth.toISOString(),
      });
      setMessage(`Voucher toko ${voucherCode} berhasil dibuat!`);
      setVoucherCode('');
      setActiveTab('overview');
    } catch (err) {
      setError(err.response?.data?.message || 'Gagal membuat voucher');
    }
  };

  if (!isAuthenticated) {
    return (
      <div className="container" style={{ padding: '60px 16px', textAlign: 'center' }}>
        <button onClick={openLoginModal} className="btn btn-primary btn-lg">
          Silakan Masuk Untuk Mengakses Seller Center
        </button>
      </div>
    );
  }

  if (user?.role === 'buyer') {
    return (
      <div className="container" style={{ padding: '60px 16px', textAlign: 'center' }}>
        <div className="card" style={{ maxWidth: '520px', margin: '0 auto', padding: '40px 28px' }}>
          <div style={{ width: '64px', height: '64px', borderRadius: '50%', background: '#FEF2F2', color: '#E02954', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px' }}>
            <AlertCircle size={32} />
          </div>
          <h2 style={{ fontSize: '20px', fontWeight: 800, marginBottom: '8px' }}>Akses Seller Center Dibatasi</h2>
          <p style={{ fontSize: '14px', color: 'var(--text-secondary)', lineHeight: '1.6', marginBottom: '24px' }}>
            Akun Anda (<strong>{user?.email}</strong>) terdaftar sebagai <strong>Pembeli (Buyer)</strong>. Akun Pembeli tidak diizinkan membuka toko atau mengakses Seller Center.
          </p>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <button
              onClick={() => {
                logout();
                openRegisterModal();
              }}
              className="btn btn-primary btn-lg btn-block"
            >
              Daftar Sebagai Penjual (Seller)
            </button>
            <a href="/" className="btn btn-secondary btn-block">
              Kembali ke Beranda Belanja
            </a>
          </div>
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

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      {message && (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '12px 16px', background: 'var(--color-primary-light)', color: 'var(--color-primary-dark)', borderRadius: '8px', marginBottom: '20px', fontWeight: 600 }}>
          <CheckCircle size={18} />
          <span>{message}</span>
        </div>
      )}
      {error && (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '12px 16px', background: 'var(--color-accent-red-light)', color: 'var(--color-accent-red)', borderRadius: '8px', marginBottom: '20px', fontWeight: 600 }}>
          <AlertCircle size={18} />
          <span>{error}</span>
        </div>
      )}

      {/* If User has NO Store -> Prompt Registration */}
      {!store && !loading ? (
        <div className="card" style={{ maxWidth: '540px', margin: '40px auto', padding: '36px 28px', textAlign: 'center' }}>
          <div style={{ width: '64px', height: '64px', borderRadius: '50%', background: 'var(--color-primary-light)', color: 'var(--color-primary)', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px' }}>
            <Store size={32} />
          </div>
          <h2 style={{ fontSize: '22px', fontWeight: 800, marginBottom: '8px' }}>
            Buka Toko Gratis di TokoMarket
          </h2>
          <p style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '24px' }}>
            Mulai jualan online, jangkau jutaan pembeli, dan dapatkan pencairan saldo instan!
          </p>

          <form onSubmit={handleCreateStore} style={{ textAlign: 'left' }}>
            <div className="form-group">
              <label className="form-label">Nama Toko</label>
              <input
                type="text"
                required
                placeholder="Contoh: Official Gadget Store"
                value={newStoreName}
                onChange={(e) => setNewStoreName(e.target.value)}
                className="form-input"
              />
            </div>
            <div className="form-group">
              <label className="form-label">Domain Slug Toko (URL)</label>
              <input
                type="text"
                required
                placeholder="contoh: gadget-store"
                value={newDomainSlug}
                onChange={(e) => setNewDomainSlug(e.target.value.toLowerCase().replace(/[^a-z0-9]/g, ''))}
                className="form-input"
              />
              <div style={{ fontSize: '11px', color: 'var(--text-muted)', marginTop: '4px' }}>
                URL Toko: tokomarket.com/store/{newDomainSlug || 'nama-toko'}
              </div>
            </div>
            <button
              type="submit"
              disabled={creatingStore}
              className="btn btn-primary btn-block btn-lg"
              style={{ marginTop: '16px' }}
            >
              {creatingStore ? 'Membuat Toko...' : 'Buka Toko Sekarang'}
            </button>
          </form>
        </div>
      ) : (
        /* Seller Dashboard */
        <div>
          {/* Store Header Banner */}
          <div className="card" style={{ padding: '24px', marginBottom: '24px', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
              <div style={{ width: '56px', height: '56px', borderRadius: '12px', background: 'var(--color-primary-light)', color: 'var(--color-primary-dark)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 800, fontSize: '24px' }}>
                🏬
              </div>
              <div>
                <h1 style={{ fontSize: '20px', fontWeight: 800 }}>{store?.store_name}</h1>
                <div style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>
                  Domain: <strong>tokomarket.com/store/{store?.domain_slug}</strong> • Kota Jakarta
                </div>
              </div>
            </div>

            {/* Balance Card */}
            <div style={{ background: '#F8F9FA', padding: '12px 20px', borderRadius: '10px', border: '1px solid var(--border-light)', display: 'flex', alignItems: 'center', gap: '20px' }}>
              <div>
                <div style={{ fontSize: '11px', color: 'var(--text-secondary)', fontWeight: 600 }}>Saldo Penghasilan Toko</div>
                <div style={{ fontSize: '22px', fontWeight: 900, color: 'var(--color-primary)' }}>
                  {formatRupiah(balance)}
                </div>
              </div>
              <button
                disabled={balance <= 0}
                onClick={() => setWithdrawModalOpen(true)}
                className="btn btn-orange btn-sm"
              >
                Tarik Dana
              </button>
            </div>
          </div>

          {/* Navigation Tabs */}
          <div style={{ display: 'flex', gap: '8px', borderBottom: '1px solid var(--border-light)', paddingBottom: '12px', marginBottom: '24px' }}>
            {[
              { id: 'overview', label: 'Ringkasan & Mutasi Saldo' },
              { id: 'orders', label: `Pesanan Masuk (${orders.length})` },
              { id: 'add_product', label: '+ Tambah Produk' },
              { id: 'vouchers', label: '+ Buat Voucher Toko' },
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                style={{
                  padding: '8px 16px',
                  borderRadius: '8px',
                  fontWeight: 700,
                  fontSize: '14px',
                  background: activeTab === tab.id ? 'var(--color-primary)' : 'transparent',
                  color: activeTab === tab.id ? 'white' : 'var(--text-secondary)',
                }}
              >
                {tab.label}
              </button>
            ))}
          </div>

          {/* Tab 1: Overview & Mutation Ledger */}
          {activeTab === 'overview' && (
            <div style={{ display: 'grid', gridTemplateColumns: '1fr', gap: '24px' }}>
              <div className="card" style={{ padding: '24px' }}>
                <h3 style={{ fontSize: '16px', fontWeight: 800, marginBottom: '16px' }}>
                  Buku Mutasi Saldo (Audit Ledger)
                </h3>
                {mutations.length === 0 ? (
                  <div style={{ color: 'var(--text-secondary)', fontSize: '13px', textAlign: 'center', padding: '24px 0' }}>
                    Belum ada riwayat transaksi saldo.
                  </div>
                ) : (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    {mutations.map((m) => (
                      <div
                        key={m.id}
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'space-between',
                          padding: '12px 16px',
                          borderRadius: '8px',
                          border: '1px solid var(--border-light)',
                          background: m.type === 'CREDIT' ? '#F0FDF4' : '#FEF2F2',
                        }}
                      >
                        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                          <div style={{ padding: '8px', borderRadius: '50%', background: m.type === 'CREDIT' ? 'var(--color-primary-light)' : 'var(--color-accent-red-light)', color: m.type === 'CREDIT' ? 'var(--color-primary-dark)' : 'var(--color-accent-red)' }}>
                            {m.type === 'CREDIT' ? <ArrowDownRight size={18} /> : <ArrowUpRight size={18} />}
                          </div>
                          <div>
                            <div style={{ fontWeight: 700, fontSize: '14px' }}>{m.description}</div>
                            <div style={{ fontSize: '11px', color: 'var(--text-muted)' }}>
                              {new Date(m.created_at).toLocaleString('id-ID')}
                            </div>
                          </div>
                        </div>

                        <div style={{ textAlign: 'right' }}>
                          <div style={{ fontWeight: 800, fontSize: '15px', color: m.type === 'CREDIT' ? 'var(--color-primary)' : 'var(--color-accent-red)' }}>
                            {m.type === 'CREDIT' ? '+' : '-'}{formatRupiah(m.amount)}
                          </div>
                          <div style={{ fontSize: '11px', color: 'var(--text-secondary)' }}>
                            Saldo Akhir: {formatRupiah(m.balance_after)}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Tab 2: Orders & Shipping Input */}
          {activeTab === 'orders' && (
            <div className="card" style={{ padding: '24px' }}>
              <h3 style={{ fontSize: '16px', fontWeight: 800, marginBottom: '16px' }}>
                Daftar Pesanan Toko Masuk
              </h3>
              {orders.length === 0 ? (
                <div style={{ textAlign: 'center', padding: '32px 0', color: 'var(--text-secondary)', fontSize: '14px' }}>
                  Belum ada pesanan masuk.
                </div>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  {orders.map((ord) => (
                    <div key={ord.id} style={{ border: '1px solid var(--border-light)', borderRadius: '8px', padding: '16px' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '12px' }}>
                        <div>
                          <div style={{ fontWeight: 700, fontSize: '14px' }}>{ord.product?.name}</div>
                          <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>
                            Jumlah: {ord.quantity}x • Total: {formatRupiah(ord.price * ord.quantity + ord.shipping_cost)}
                          </div>
                        </div>
                        <div>
                          <span className="badge badge-green">{ord.status}</span>
                        </div>
                      </div>

                      {ord.status === 'PROCESSING' && (
                        <button
                          onClick={() => {
                            setSelectedOrderItemId(ord.id);
                            setShipModalOpen(true);
                          }}
                          className="btn btn-primary btn-sm"
                        >
                          <Truck size={14} /> Input Nomor Resi
                        </button>
                      )}

                      {ord.tracking_number && (
                        <div style={{ fontSize: '12px', color: '#0284C7', fontWeight: 600 }}>
                          Resi: {ord.tracking_number} ({ord.courier_name})
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Tab 3: Add Product Form */}
          {activeTab === 'add_product' && (
            <div className="card" style={{ maxWidth: '640px', padding: '24px' }}>
              <h3 style={{ fontSize: '16px', fontWeight: 800, marginBottom: '16px' }}>
                Tambah Produk Baru ke Katalog Toko
              </h3>
              <form onSubmit={handleCreateProduct}>
                <div className="form-group">
                  <label className="form-label">Nama Produk</label>
                  <input
                    type="text"
                    required
                    placeholder="Contoh: Laptop Asus ROG Strix"
                    value={productName}
                    onChange={(e) => setProductName(e.target.value)}
                    className="form-input"
                  />
                </div>

                <div className="form-group">
                  <label className="form-label">Kategori</label>
                  <select
                    value={productCategoryId}
                    onChange={(e) => setProductCategoryId(e.target.value)}
                    className="form-select"
                  >
                    {categories.map((c) => (
                      <option key={c.id} value={c.id}>{c.name}</option>
                    ))}
                  </select>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                  <div className="form-group">
                    <label className="form-label">Harga (Rp)</label>
                    <input
                      type="number"
                      required
                      placeholder="1500000"
                      value={productPrice}
                      onChange={(e) => setProductPrice(e.target.value)}
                      className="form-input"
                    />
                  </div>
                  <div className="form-group">
                    <label className="form-label">Stok</label>
                    <input
                      type="number"
                      required
                      placeholder="10"
                      value={productStock}
                      onChange={(e) => setProductStock(e.target.value)}
                      className="form-input"
                    />
                  </div>
                </div>

                <div className="form-group">
                  <label className="form-label">Deskripsi Lengkap</label>
                  <textarea
                    rows="3"
                    placeholder="Spesifikasi, keunggulan, garansi..."
                    value={productDesc}
                    onChange={(e) => setProductDesc(e.target.value)}
                    className="form-textarea"
                  />
                </div>

                {/* Varian Builder */}
                <div style={{ background: '#F8F9FA', padding: '16px', borderRadius: '8px', marginBottom: '20px' }}>
                  <div style={{ fontWeight: 700, fontSize: '13px', marginBottom: '8px' }}>
                    Varian Produk (Opsional - Size / Warna / RAM)
                  </div>
                  <div style={{ display: 'flex', gap: '8px', marginBottom: '8px' }}>
                    <input
                      type="text"
                      placeholder="Nama Varian (e.g. Size M)"
                      value={newVariantName}
                      onChange={(e) => setNewVariantName(e.target.value)}
                      className="form-input"
                      style={{ fontSize: '13px' }}
                    />
                    <input
                      type="number"
                      placeholder="Stok Varian"
                      value={newVariantStock}
                      onChange={(e) => setNewVariantStock(e.target.value)}
                      className="form-input"
                      style={{ width: '120px', fontSize: '13px' }}
                    />
                    <button type="button" onClick={handleAddVariant} className="btn btn-outline btn-sm">
                      + Tambah
                    </button>
                  </div>

                  {variants.map((v, i) => (
                    <div key={i} style={{ fontSize: '12px', color: 'var(--text-secondary)', padding: '2px 0' }}>
                      • {v.variant_name} (Stok: {v.stock})
                    </div>
                  ))}
                </div>

                <button type="submit" className="btn btn-primary btn-block btn-lg">
                  Simpan & Terbitkan Produk
                </button>
              </form>
            </div>
          )}

          {/* Tab 4: Create Store Voucher */}
          {activeTab === 'vouchers' && (
            <div className="card" style={{ maxWidth: '540px', padding: '24px' }}>
              <h3 style={{ fontSize: '16px', fontWeight: 800, marginBottom: '16px' }}>
                Buat Kupon Diskon Promo Toko
              </h3>
              <form onSubmit={handleCreateStoreVoucher}>
                <div className="form-group">
                  <label className="form-label">Kode Voucher</label>
                  <input
                    type="text"
                    required
                    placeholder="Contoh: DISKONHEMAT"
                    value={voucherCode}
                    onChange={(e) => setVoucherCode(e.target.value.toUpperCase())}
                    className="form-input"
                  />
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                  <div className="form-group">
                    <label className="form-label">Diskon (%)</label>
                    <input
                      type="number"
                      required
                      value={voucherPercent}
                      onChange={(e) => setVoucherPercent(e.target.value)}
                      className="form-input"
                    />
                  </div>
                  <div className="form-group">
                    <label className="form-label">Maksimal Potongan (Rp)</label>
                    <input
                      type="number"
                      required
                      value={voucherMaxDiscount}
                      onChange={(e) => setVoucherMaxDiscount(e.target.value)}
                      className="form-input"
                    />
                  </div>
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                  <div className="form-group">
                    <label className="form-label">Minimal Belanja (Rp)</label>
                    <input
                      type="number"
                      required
                      value={voucherMinSpend}
                      onChange={(e) => setVoucherMinSpend(e.target.value)}
                      className="form-input"
                    />
                  </div>
                  <div className="form-group">
                    <label className="form-label">Kuota Penggunaan</label>
                    <input
                      type="number"
                      required
                      value={voucherQuota}
                      onChange={(e) => setVoucherQuota(e.target.value)}
                      className="form-input"
                    />
                  </div>
                </div>
                <button type="submit" className="btn btn-primary btn-block btn-lg">
                  Terbitkan Voucher Toko
                </button>
              </form>
            </div>
          )}
        </div>
      )}

      {/* Withdrawal Modal */}
      {withdrawModalOpen && (
        <div className="modal-overlay" onClick={() => setWithdrawModalOpen(false)}>
          <div className="modal-card" onClick={(e) => e.stopPropagation()} style={{ padding: '24px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 800, marginBottom: '16px' }}>
              Tarik Dana Saldo Toko
            </h3>
            <form onSubmit={handleWithdrawal}>
              <div className="form-group">
                <label className="form-label">Bank Tujuan</label>
                <select
                  value={withdrawBank}
                  onChange={(e) => setWithdrawBank(e.target.value)}
                  className="form-select"
                >
                  <option value="BCA">Bank BCA</option>
                  <option value="Mandiri">Bank Mandiri</option>
                  <option value="BRI">Bank BRI</option>
                  <option value="BNI">Bank BNI</option>
                </select>
              </div>
              <div className="form-group">
                <label className="form-label">Nomor Rekening</label>
                <input
                  type="text"
                  required
                  placeholder="0123456789"
                  value={withdrawAccNo}
                  onChange={(e) => setWithdrawAccNo(e.target.value)}
                  className="form-input"
                />
              </div>
              <div className="form-group">
                <label className="form-label">Nama Pemilik Rekening</label>
                <input
                  type="text"
                  required
                  value={withdrawAccHolder}
                  onChange={(e) => setWithdrawAccHolder(e.target.value)}
                  className="form-input"
                />
              </div>
              <div className="form-group">
                <label className="form-label">Nominal Penarikan (Rp)</label>
                <input
                  type="number"
                  required
                  max={balance}
                  placeholder={`Maksimal ${formatRupiah(balance)}`}
                  value={withdrawAmount}
                  onChange={(e) => setWithdrawAmount(e.target.value)}
                  className="form-input"
                />
              </div>
              <div style={{ display: 'flex', gap: '8px', marginTop: '20px' }}>
                <button type="submit" disabled={withdrawing} className="btn btn-primary btn-block btn-lg">
                  {withdrawing ? 'Memproses...' : 'Ajukan Penarikan Dana'}
                </button>
                <button type="button" onClick={() => setWithdrawModalOpen(false)} className="btn btn-secondary">
                  Batal
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Input Resi Modal */}
      {shipModalOpen && (
        <div className="modal-overlay" onClick={() => setShipModalOpen(false)}>
          <div className="modal-card" onClick={(e) => e.stopPropagation()} style={{ padding: '24px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 800, marginBottom: '16px' }}>
              Input Nomor Resi Pengiriman
            </h3>
            <form onSubmit={handleShipOrder}>
              <div className="form-group">
                <label className="form-label">Kurir Pengiriman</label>
                <input
                  type="text"
                  value={shippingCourier}
                  onChange={(e) => setShippingCourier(e.target.value)}
                  className="form-input"
                />
              </div>
              <div className="form-group">
                <label className="form-label">Nomor Resi / AWB</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: JNE-TRACK-998877"
                  value={trackingNumber}
                  onChange={(e) => setTrackingNumber(e.target.value)}
                  className="form-input"
                />
              </div>
              <div style={{ display: 'flex', gap: '8px', marginTop: '20px' }}>
                <button type="submit" className="btn btn-primary btn-block btn-lg">
                  Kirim Pesanan
                </button>
                <button type="button" onClick={() => setShipModalOpen(false)} className="btn btn-secondary">
                  Batal
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
