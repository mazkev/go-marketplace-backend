import React from 'react';
import { ShieldCheck, Truck, Headphones, RotateCcw } from 'lucide-react';

export const Footer = () => {
  return (
    <footer style={{ background: '#FFFFFF', borderTop: '1px solid var(--border-light)', marginTop: '48px', paddingTop: '32px' }}>
      {/* Value Proposition */}
      <div className="container" style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '24px', paddingBottom: '32px', borderBottom: '1px solid var(--border-light)' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <div style={{ background: 'var(--color-primary-light)', padding: '10px', borderRadius: '50%', color: 'var(--color-primary-dark)' }}>
            <ShieldCheck size={24} />
          </div>
          <div>
            <div style={{ fontWeight: 700, fontSize: '14px' }}>Transaksi Aman & Escrow</div>
            <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Dana diteruskan saat barang diterima</div>
          </div>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <div style={{ background: 'var(--color-secondary-orange-light)', padding: '10px', borderRadius: '50%', color: 'var(--color-secondary-orange)' }}>
            <Truck size={24} />
          </div>
          <div>
            <div style={{ fontWeight: 700, fontSize: '14px' }}>Pengiriman Seluruh Indonesia</div>
            <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Multi-kurir terpercaya & resi real-time</div>
          </div>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <div style={{ background: 'var(--color-primary-light)', padding: '10px', borderRadius: '50%', color: 'var(--color-primary-dark)' }}>
            <RotateCcw size={24} />
          </div>
          <div>
            <div style={{ fontWeight: 700, fontSize: '14px' }}>Garansi Kepuasan</div>
            <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Bebas komplain jika barang tidak sesuai</div>
          </div>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <div style={{ background: '#E0F2FE', padding: '10px', borderRadius: '50%', color: '#0284C7' }}>
            <Headphones size={24} />
          </div>
          <div>
            <div style={{ fontWeight: 700, fontSize: '14px' }}>Layanan Bantuan 24/7</div>
            <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Tim customer care siap membantu</div>
          </div>
        </div>
      </div>

      {/* Main Footer Links */}
      <div className="container" style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '32px', padding: '32px 16px', fontSize: '13px', color: 'var(--text-secondary)' }}>
        <div>
          <div style={{ fontWeight: 700, color: 'var(--text-primary)', marginBottom: '12px', fontSize: '14px' }}>TokoMarket</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <span>Tentang TokoMarket</span>
            <span>Hak Kekayaan Intelektual</span>
            <span>Karir</span>
            <span>Blog TokoMarket</span>
            <span>Mitra Blog</span>
          </div>
        </div>

        <div>
          <div style={{ fontWeight: 700, color: 'var(--text-primary)', marginBottom: '12px', fontSize: '14px' }}>Beli</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <span>Tagihan & Top Up</span>
            <span>TokoMarket COD</span>
            <span>Bebas Ongkir</span>
            <span>Promo Hari Ini</span>
          </div>
        </div>

        <div>
          <div style={{ fontWeight: 700, color: 'var(--text-primary)', marginBottom: '12px', fontSize: '14px' }}>Jual</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <span>Pusat Edukasi Seller</span>
            <span>Mitra Toppers</span>
            <span>Daftar Official Store</span>
          </div>
        </div>

        <div>
          <div style={{ fontWeight: 700, color: 'var(--text-primary)', marginBottom: '12px', fontSize: '14px' }}>Bantuan dan Panduan</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <span>TokoMarket Care</span>
            <span>Syarat dan Ketentuan</span>
            <span>Kebijakan Privasi</span>
          </div>
        </div>
      </div>

      <div style={{ background: '#F8F9FA', borderTop: '1px solid var(--border-light)', padding: '16px 0', textAlign: 'center', fontSize: '12px', color: 'var(--text-muted)' }}>
        © 2026 TokoMarket Indonesia. All rights reserved. Clean Architecture E-Commerce.
      </div>
    </footer>
  );
};
