import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Package, Truck, CheckCircle, Star, AlertCircle, Clock } from 'lucide-react';
import { orderService } from '../services/orderService';
import { useAuth } from '../context/AuthContext';

export const OrderHistoryPage = () => {
  const { isAuthenticated, openLoginModal } = useAuth();
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('ALL');

  // Review Modal State
  const [reviewModalOpen, setReviewModalOpen] = useState(false);
  const [reviewOrderItemId, setReviewOrderItemId] = useState(null);
  const [reviewRating, setReviewRating] = useState(5);
  const [reviewComment, setReviewComment] = useState('');
  const [submittingReview, setSubmittingReview] = useState(false);
  const [message, setMessage] = useState('');

  useEffect(() => {
    if (isAuthenticated) {
      fetchOrders();
    }
  }, [isAuthenticated]);

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const res = await orderService.getUserOrders();
      if (res.data) setOrders(res.data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleCompleteItem = async (orderItemId) => {
    try {
      await orderService.completeOrderItem(orderItemId);
      setMessage('Pesanan berhasil diselesaikan! Dana escrow telah diteruskan ke toko penjual.');
      await fetchOrders();
      setTimeout(() => setMessage(''), 4000);
    } catch (err) {
      alert(err.response?.data?.message || 'Gagal menyelesaikan pesanan');
    }
  };

  const handleOpenReview = (orderItemId) => {
    setReviewOrderItemId(orderItemId);
    setReviewRating(5);
    setReviewComment('');
    setReviewModalOpen(true);
  };

  const handleSubmitReview = async (e) => {
    e.preventDefault();
    setSubmittingReview(true);
    try {
      await orderService.createReview({
        order_item_id: reviewOrderItemId,
        rating: reviewRating,
        comment: reviewComment,
      });
      setReviewModalOpen(false);
      setMessage('Ulasan berhasil dikirim! Terima kasih atas ulasan Anda.');
      await fetchOrders();
      setTimeout(() => setMessage(''), 4000);
    } catch (err) {
      alert(err.response?.data?.message || 'Gagal mengirim ulasan');
    } finally {
      setSubmittingReview(false);
    }
  };

  if (!isAuthenticated) {
    return (
      <div className="container" style={{ padding: '60px 16px', textAlign: 'center' }}>
        <button onClick={openLoginModal} className="btn btn-primary btn-lg">
          Silakan Masuk Untuk Melihat Pesanan
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

  const getStatusBadge = (status) => {
    switch (status) {
      case 'PENDING':
        return <span className="badge badge-orange">Menunggu Pembayaran</span>;
      case 'PROCESSING':
        return <span className="badge badge-green">Sedang Diproses Penjual</span>;
      case 'SHIPPED':
        return <span className="badge badge-green" style={{ background: '#E0F2FE', color: '#0284C7' }}>Sedang Dikirim</span>;
      case 'COMPLETED':
        return <span className="badge badge-green">Pesanan Selesai</span>;
      case 'CANCELLED':
        return <span className="badge badge-red">Dibatalkan</span>;
      default:
        return <span className="badge">{status}</span>;
    }
  };

  return (
    <div className="container" style={{ padding: '24px 16px' }}>
      <h1 style={{ fontSize: '22px', fontWeight: 800, marginBottom: '20px' }}>
        Daftar Pesanan Saya
      </h1>

      {message && (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '12px 16px', background: 'var(--color-primary-light)', color: 'var(--color-primary-dark)', borderRadius: '8px', marginBottom: '20px', fontWeight: 600 }}>
          <CheckCircle size={18} />
          <span>{message}</span>
        </div>
      )}

      {loading ? (
        <div style={{ textAlign: 'center', padding: '40px 0' }}>Memuat riwayat pesanan...</div>
      ) : orders.length === 0 ? (
        <div className="card" style={{ padding: '60px 20px', textAlign: 'center' }}>
          <div style={{ fontSize: '64px', marginBottom: '16px' }}>📦</div>
          <h3 style={{ fontSize: '18px', fontWeight: 800, marginBottom: '8px' }}>
            Belum Ada Transaksi
          </h3>
          <p style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '24px' }}>
            Yuk, mulai belanja dan temukan berbagai produk menarik di TokoMarket!
          </p>
          <Link to="/" className="btn btn-primary btn-lg">
            Belanja Sekarang
          </Link>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          {orders.map((order) => (
            <div key={order.id} className="card" style={{ padding: '20px' }}>
              {/* Order Header */}
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid var(--border-light)', paddingBottom: '12px', marginBottom: '16px' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px', fontSize: '13px' }}>
                  <Package size={18} color="var(--color-primary)" />
                  <span style={{ fontWeight: 700 }}>Belanja</span>
                  <span style={{ color: 'var(--text-muted)' }}>{new Date(order.created_at).toLocaleDateString('id-ID')}</span>
                  <span style={{ fontWeight: 600, color: 'var(--text-secondary)' }}>{order.invoice_number}</span>
                </div>
                <div>{getStatusBadge(order.payment_status)}</div>
              </div>

              {/* Items Breakdown */}
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                {order.order_items?.map((item) => (
                  <div
                    key={item.id}
                    style={{
                      display: 'grid',
                      gridTemplateColumns: '64px 1fr auto',
                      gap: '16px',
                      alignItems: 'center',
                    }}
                  >
                    <div style={{ width: '64px', height: '64px', borderRadius: '8px', background: '#F0FDF4', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '24px' }}>
                      🛍️
                    </div>

                    <div>
                      <div style={{ fontSize: '12px', color: 'var(--color-primary)', fontWeight: 700, marginBottom: '2px' }}>
                        {item.store?.store_name || 'Toko Penjual'}
                      </div>
                      <div style={{ fontWeight: 600, fontSize: '14px', marginBottom: '4px' }}>
                        {item.product?.name || 'Produk'}
                      </div>
                      <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>
                        {item.quantity} barang x {formatRupiah(item.price)}
                      </div>
                      {item.tracking_number && (
                        <div style={{ display: 'flex', alignItems: 'center', gap: '4px', fontSize: '12px', color: '#0284C7', marginTop: '4px', fontWeight: 600 }}>
                          <Truck size={14} /> Resi: {item.tracking_number} ({item.courier_name})
                        </div>
                      )}
                    </div>

                    {/* Actions per item */}
                    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: '8px' }}>
                      <div style={{ fontSize: '12px', fontWeight: 600 }}>
                        Status Item: {getStatusBadge(item.status)}
                      </div>

                      {/* Complete order button */}
                      {item.status === 'SHIPPED' && (
                        <button
                          onClick={() => handleCompleteItem(item.id)}
                          className="btn btn-primary btn-sm"
                        >
                          <CheckCircle size={14} /> Selesaikan Pesanan
                        </button>
                      )}

                      {/* Review Button */}
                      {item.status === 'COMPLETED' && !item.review && (
                        <button
                          onClick={() => handleOpenReview(item.id)}
                          className="btn btn-orange btn-sm"
                        >
                          <Star size={14} /> Beri Ulasan
                        </button>
                      )}

                      {item.review && (
                        <div style={{ display: 'flex', alignItems: 'center', gap: '4px', fontSize: '12px', color: 'var(--color-primary)', fontWeight: 600 }}>
                          <Star size={14} fill="#FFC400" color="#FFC400" /> Sudah Diulas ({item.review.rating}★)
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>

              {/* Total Order Amount */}
              <div style={{ display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: '16px', borderTop: '1px solid var(--border-light)', paddingTop: '16px', marginTop: '16px' }}>
                <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>Total Tagihan:</span>
                <span style={{ fontSize: '18px', fontWeight: 900, color: 'var(--color-primary)' }}>
                  {formatRupiah(order.final_amount || order.total_amount)}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Review Modal */}
      {reviewModalOpen && (
        <div className="modal-overlay" onClick={() => setReviewModalOpen(false)}>
          <div className="modal-card" onClick={(e) => e.stopPropagation()} style={{ padding: '24px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 800, marginBottom: '16px' }}>
              Beri Ulasan Produk
            </h3>

            <form onSubmit={handleSubmitReview}>
              {/* Interactive Stars */}
              <div className="form-group" style={{ textAlign: 'center', marginBottom: '20px' }}>
                <div style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '8px' }}>
                  Pilih Rating Bintang:
                </div>
                <div style={{ display: 'flex', justifyContent: 'center', gap: '8px' }}>
                  {[1, 2, 3, 4, 5].map((star) => (
                    <button
                      key={star}
                      type="button"
                      onClick={() => setReviewRating(star)}
                      style={{ padding: '4px', transition: 'transform 0.1s' }}
                    >
                      <Star
                        size={32}
                        fill={star <= reviewRating ? '#FFC400' : '#E0E0E0'}
                        color={star <= reviewRating ? '#FFC400' : '#E0E0E0'}
                      />
                    </button>
                  ))}
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Tulis Ulasan Anda</label>
                <textarea
                  rows="4"
                  required
                  placeholder="Ceritakan kepuasan Anda tentang kualitas barang, kecepatan pengiriman, atau pelayanan penjual..."
                  value={reviewComment}
                  onChange={(e) => setReviewComment(e.target.value)}
                  className="form-textarea"
                />
              </div>

              <div style={{ display: 'flex', gap: '8px', marginTop: '20px' }}>
                <button
                  type="submit"
                  disabled={submittingReview}
                  className="btn btn-primary btn-block btn-lg"
                >
                  {submittingReview ? 'Mengirim...' : 'Kirim Ulasan'}
                </button>
                <button
                  type="button"
                  onClick={() => setReviewModalOpen(false)}
                  className="btn btn-secondary"
                >
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
