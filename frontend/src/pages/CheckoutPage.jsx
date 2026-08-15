import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import {
  MapPin,
  Truck,
  Tag,
  CreditCard,
  CheckCircle,
  AlertCircle,
  Store,
  ChevronRight,
  ShieldCheck,
} from 'lucide-react';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';
import { authService } from '../services/authService';
import { orderService } from '../services/orderService';
import { voucherService } from '../services/voucherService';

export const CheckoutPage = () => {
  const { cart, fetchCart } = useCart();
  const { user, isAuthenticated, openLoginModal } = useAuth();
  const navigate = useNavigate();

  const [addresses, setAddresses] = useState([]);
  const [selectedAddressId, setSelectedAddressId] = useState(null);
  const [isAddressModalOpen, setIsAddressModalOpen] = useState(false);
  const [newAddress, setNewAddress] = useState({
    receiver_name: user?.name || '',
    phone: user?.phone || '',
    full_address: '',
    city_id: 1,
    is_primary: true,
  });

  const [storeCouriers, setStoreCouriers] = useState({});
  const [voucherCode, setVoucherCode] = useState('');
  const [voucherDiscount, setVoucherDiscount] = useState(0);
  const [voucherMessage, setVoucherMessage] = useState('');
  const [paymentMethod, setPaymentMethod] = useState('BCA_VA');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [orderResult, setOrderResult] = useState(null);

  useEffect(() => {
    if (isAuthenticated) {
      authService.getAddresses()
        .then((res) => {
          if (res.data && res.data.length > 0) {
            setAddresses(res.data);
            const primary = res.data.find((a) => a.is_primary) || res.data[0];
            setSelectedAddressId(primary.id);
          }
        })
        .catch(console.error);
    }
  }, [isAuthenticated]);

  if (!isAuthenticated) {
    return (
      <div className="container" style={{ padding: '60px 16px', textAlign: 'center' }}>
        <button onClick={openLoginModal} className="btn btn-primary btn-lg">
          Silakan Masuk Untuk Melanjutkan Checkout
        </button>
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

  const storeShippingFee = 10000;
  const totalShipping = (cart.stores?.length || 0) * storeShippingFee;
  const totalOrderBeforeDiscount = (cart.total_price || 0) + totalShipping;
  const finalTotal = Math.max(0, totalOrderBeforeDiscount - voucherDiscount);

  const handleApplyVoucher = async () => {
    if (!voucherCode.trim()) return;
    try {
      const res = await voucherService.applyVoucher(voucherCode.trim(), totalOrderBeforeDiscount);
      if (res.data && res.data.valid) {
        setVoucherDiscount(res.data.discount_amount);
        setVoucherMessage(`Voucher ${res.data.voucher_code} berhasil diterapkan! Hemat ${formatRupiah(res.data.discount_amount)}`);
      }
    } catch (err) {
      setVoucherDiscount(0);
      setVoucherMessage(err.response?.data?.message || 'Kode voucher tidak valid');
    }
  };

  const handleAddAddress = async (e) => {
    e.preventDefault();
    try {
      const res = await authService.addAddress(newAddress);
      if (res.data) {
        const updated = [...addresses, res.data];
        setAddresses(updated);
        setSelectedAddressId(res.data.id);
        setIsAddressModalOpen(false);
      }
    } catch (err) {
      alert('Gagal menambahkan alamat');
    }
  };

  const handlePlaceOrder = async () => {
    if (!selectedAddressId) {
      setError('Silakan pilih alamat pengiriman terlebih dahulu');
      return;
    }
    setLoading(true);
    setError('');

    const storeOptions = (cart.stores || []).map((s) => ({
      store_id: s.store_id,
      courier_name: storeCouriers[s.store_id] || 'JNE Reguler',
    }));

    try {
      const payload = {
        address_id: selectedAddressId,
        payment_method: paymentMethod,
        voucher_code: voucherCode || undefined,
        stores: storeOptions,
      };

      const res = await orderService.checkout(payload);
      if (res.data) {
        setOrderResult(res.data);
        await fetchCart();
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Gagal memproses pesanan');
    } finally {
      setLoading(false);
    }
  };

  const handleSimulatePayment = async () => {
    if (!orderResult) return;
    try {
      await orderService.simulatePaymentWebhook(orderResult.invoice_number, orderResult.final_amount);
      alert('Pembayaran Berhasil Dikonfirmasi via Webhook!');
      navigate('/orders');
    } catch (err) {
      alert('Simulasi pembayaran gagal');
    }
  };

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      <h1 style={{ fontSize: '22px', fontWeight: 800, marginBottom: '20px' }}>
        Pengiriman & Checkout
      </h1>

      {error && (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '12px 16px', background: 'var(--color-accent-red-light)', color: 'var(--color-accent-red)', borderRadius: '8px', marginBottom: '20px', fontWeight: 600 }}>
          <AlertCircle size={18} />
          <span>{error}</span>
        </div>
      )}

      {/* Order Success Modal */}
      {orderResult && (
        <div className="modal-overlay">
          <div className="modal-card" style={{ padding: '32px 24px', textAlign: 'center' }}>
            <div style={{ width: '64px', height: '64px', borderRadius: '50%', background: 'var(--color-primary-light)', color: 'var(--color-primary-dark)', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px' }}>
              <CheckCircle size={36} />
            </div>
            <h2 style={{ fontSize: '20px', fontWeight: 800, marginBottom: '8px' }}>Pesanan Berhasil Dibuat!</h2>
            <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '20px' }}>
              Nomor Invoice: <strong style={{ color: 'var(--text-primary)' }}>{orderResult.invoice_number}</strong>
            </p>

            <div className="card" style={{ background: '#F8F9FA', padding: '16px', marginBottom: '24px', textAlign: 'left' }}>
              <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Metode Pembayaran</div>
              <div style={{ fontSize: '14px', fontWeight: 700, marginBottom: '8px' }}>{orderResult.payment_method}</div>
              <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Nomor Virtual Account</div>
              <div style={{ fontSize: '18px', fontWeight: 900, color: 'var(--color-primary)', letterSpacing: '1px' }}>
                {orderResult.va_number || '8808123456789'}
              </div>
              <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '8px' }}>
                Total Bayar: <strong style={{ color: 'var(--text-primary)' }}>{formatRupiah(orderResult.final_amount)}</strong>
              </div>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <button onClick={handleSimulatePayment} className="btn btn-primary btn-block btn-lg">
                Simulasi Bayar Sekarang (Webhook)
              </button>
              <Link to="/orders" className="btn btn-secondary btn-block">
                Lihat Daftar Pesanan
              </Link>
            </div>
          </div>
        </div>
      )}

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 360px', gap: '32px', alignItems: 'start' }}>
        {/* Left Column: Address, Order Items per Store, Payment */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          {/* Address Section */}
          <div className="card" style={{ padding: '20px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontWeight: 800, fontSize: '15px' }}>
                <MapPin size={18} color="var(--color-primary)" />
                <span>Alamat Pengiriman</span>
              </div>
              <button
                onClick={() => setIsAddressModalOpen(true)}
                className="btn btn-outline btn-sm"
              >
                + Tambah Alamat
              </button>
            </div>

            {addresses.length === 0 ? (
              <div style={{ padding: '16px', background: '#F8F9FA', borderRadius: '8px', fontSize: '13px', color: 'var(--text-secondary)' }}>
                Belum ada alamat tersimpan. Silakan klik tombol "Tambah Alamat" di atas.
              </div>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                {addresses.map((addr) => (
                  <label
                    key={addr.id}
                    style={{
                      display: 'flex',
                      alignItems: 'flex-start',
                      gap: '12px',
                      padding: '12px',
                      borderRadius: '8px',
                      border: selectedAddressId === addr.id ? '2px solid var(--color-primary)' : '1px solid var(--border-light)',
                      background: selectedAddressId === addr.id ? 'var(--color-primary-ultra-light)' : 'white',
                      cursor: 'pointer',
                    }}
                  >
                    <input
                      type="radio"
                      name="address"
                      checked={selectedAddressId === addr.id}
                      onChange={() => setSelectedAddressId(addr.id)}
                      style={{ marginTop: '3px' }}
                    />
                    <div>
                      <div style={{ fontWeight: 700, fontSize: '14px' }}>
                        {addr.receiver_name}{' '}
                        <span style={{ fontSize: '12px', fontWeight: 400, color: 'var(--text-muted)' }}>({addr.phone})</span>
                      </div>
                      <div style={{ fontSize: '13px', color: 'var(--text-secondary)', marginTop: '2px' }}>
                        {addr.full_address}
                      </div>
                    </div>
                  </label>
                ))}
              </div>
            )}
          </div>

          {/* Store Sub-Orders Breakdown */}
          {cart.stores?.map((store) => (
            <div key={store.store_id} className="card" style={{ padding: '20px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontWeight: 700, fontSize: '14px', marginBottom: '16px' }}>
                <Store size={16} color="var(--color-primary)" />
                <span>Pesanan dari {store.store_name}</span>
              </div>

              {/* Items */}
              <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', marginBottom: '16px' }}>
                {store.items.map((item) => (
                  <div key={item.id} style={{ display: 'flex', justifyContent: 'space-between', fontSize: '13px' }}>
                    <div>
                      <span style={{ fontWeight: 600 }}>{item.quantity}x</span> {item.product_name}
                      {item.variant_name && <span style={{ color: 'var(--text-muted)' }}> ({item.variant_name})</span>}
                    </div>
                    <div style={{ fontWeight: 700 }}>{formatRupiah(item.item_total)}</div>
                  </div>
                ))}
              </div>

              {/* Courier Selector per Store */}
              <div style={{ background: '#F8F9FA', padding: '12px', borderRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '13px', fontWeight: 600 }}>
                  <Truck size={16} color="var(--color-primary)" />
                  <span>Pilih Kurir Pengiriman</span>
                </div>
                <select
                  value={storeCouriers[store.store_id] || 'JNE Reguler'}
                  onChange={(e) => setStoreCouriers({ ...storeCouriers, [store.store_id]: e.target.value })}
                  className="form-select"
                  style={{ width: 'auto', padding: '6px 12px', fontSize: '13px' }}
                >
                  <option value="JNE Reguler">JNE Reguler (Rp 10.000)</option>
                  <option value="SiCepat BEST">SiCepat BEST (Rp 10.000)</option>
                  <option value="J&T Express">J&T Express (Rp 10.000)</option>
                  <option value="Anteraja">Anteraja (Rp 10.000)</option>
                </select>
              </div>
            </div>
          ))}

          {/* Payment Method */}
          <div className="card" style={{ padding: '20px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontWeight: 800, fontSize: '15px', marginBottom: '16px' }}>
              <CreditCard size={18} color="var(--color-primary)" />
              <span>Metode Pembayaran</span>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '12px' }}>
              {[
                { id: 'BCA_VA', name: 'BCA Virtual Account' },
                { id: 'MANDIRI_VA', name: 'Mandiri Virtual Account' },
                { id: 'QRIS', name: 'QRIS / GoPay' },
              ].map((m) => (
                <button
                  key={m.id}
                  onClick={() => setPaymentMethod(m.id)}
                  style={{
                    padding: '12px',
                    borderRadius: '8px',
                    border: paymentMethod === m.id ? '2px solid var(--color-primary)' : '1px solid var(--border-medium)',
                    background: paymentMethod === m.id ? 'var(--color-primary-light)' : 'white',
                    fontWeight: 700,
                    fontSize: '13px',
                    textAlign: 'center',
                    color: paymentMethod === m.id ? 'var(--color-primary-dark)' : 'var(--text-primary)',
                  }}
                >
                  {m.name}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Right Column: Voucher & Final Summary */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', position: 'sticky', top: '100px' }}>
          {/* Voucher Box */}
          <div className="card" style={{ padding: '20px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '6px', fontWeight: 800, fontSize: '14px', marginBottom: '12px' }}>
              <Tag size={16} color="var(--color-secondary-orange)" />
              <span>Makin Hemat Pakai Promo</span>
            </div>
            <div style={{ display: 'flex', gap: '8px' }}>
              <input
                type="text"
                placeholder="Masukkan kode promo (e.g. DISKON10)"
                value={voucherCode}
                onChange={(e) => setVoucherCode(e.target.value.toUpperCase())}
                className="form-input"
                style={{ fontSize: '13px' }}
              />
              <button
                onClick={handleApplyVoucher}
                className="btn btn-outline"
                style={{ borderColor: 'var(--color-secondary-orange)', color: 'var(--color-secondary-orange)', whiteSpace: 'nowrap' }}
              >
                Gunakan
              </button>
            </div>
            {voucherMessage && (
              <div style={{ fontSize: '12px', color: voucherDiscount > 0 ? 'var(--color-primary)' : 'var(--color-accent-red)', marginTop: '8px', fontWeight: 600 }}>
                {voucherMessage}
              </div>
            )}
          </div>

          {/* Summary Box */}
          <div className="card" style={{ padding: '20px' }}>
            <div style={{ fontWeight: 800, fontSize: '16px', marginBottom: '16px' }}>
              Ringkasan Pembayaran
            </div>

            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '13px', marginBottom: '8px', color: 'var(--text-secondary)' }}>
              <span>Total Harga Barang</span>
              <span style={{ fontWeight: 600, color: 'var(--text-primary)' }}>{formatRupiah(cart.total_price)}</span>
            </div>

            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '13px', marginBottom: '8px', color: 'var(--text-secondary)' }}>
              <span>Total Ongkos Kirim ({cart.stores?.length} Toko)</span>
              <span style={{ fontWeight: 600, color: 'var(--text-primary)' }}>{formatRupiah(totalShipping)}</span>
            </div>

            {voucherDiscount > 0 && (
              <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '13px', marginBottom: '8px', color: 'var(--color-primary)' }}>
                <span>Diskon Promo</span>
                <span style={{ fontWeight: 700 }}>-{formatRupiah(voucherDiscount)}</span>
              </div>
            )}

            <div style={{ borderTop: '1px solid var(--border-light)', paddingTop: '16px', marginTop: '12px', marginBottom: '20px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span style={{ fontSize: '15px', fontWeight: 800 }}>Total Tagihan</span>
              <span style={{ fontSize: '20px', fontWeight: 900, color: 'var(--color-primary)' }}>
                {formatRupiah(finalTotal)}
              </span>
            </div>

            <button
              onClick={handlePlaceOrder}
              disabled={loading || !selectedAddressId}
              className="btn btn-primary btn-block btn-lg"
            >
              {loading ? 'Memproses Pesanan...' : 'Bayar Sekarang'}
            </button>
          </div>
        </div>
      </div>

      {/* Add Address Modal */}
      {isAddressModalOpen && (
        <div className="modal-overlay" onClick={() => setIsAddressModalOpen(false)}>
          <div className="modal-card" onClick={(e) => e.stopPropagation()} style={{ padding: '24px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 800, marginBottom: '16px' }}>Tambah Alamat Baru</h3>
            <form onSubmit={handleAddAddress}>
              <div className="form-group">
                <label className="form-label">Nama Penerima</label>
                <input
                  type="text"
                  required
                  value={newAddress.receiver_name}
                  onChange={(e) => setNewAddress({ ...newAddress, receiver_name: e.target.value })}
                  className="form-input"
                />
              </div>
              <div className="form-group">
                <label className="form-label">Nomor Telepon</label>
                <input
                  type="tel"
                  required
                  value={newAddress.phone}
                  onChange={(e) => setNewAddress({ ...newAddress, phone: e.target.value })}
                  className="form-input"
                />
              </div>
              <div className="form-group">
                <label className="form-label">Alamat Lengkap</label>
                <textarea
                  rows="3"
                  required
                  placeholder="Nama jalan, gedung, nomor rumah, RT/RW..."
                  value={newAddress.full_address}
                  onChange={(e) => setNewAddress({ ...newAddress, full_address: e.target.value })}
                  className="form-textarea"
                />
              </div>
              <div style={{ display: 'flex', gap: '8px', marginTop: '16px' }}>
                <button type="submit" className="btn btn-primary btn-block">Simpan Alamat</button>
                <button type="button" onClick={() => setIsAddressModalOpen(false)} className="btn btn-secondary">Batal</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
