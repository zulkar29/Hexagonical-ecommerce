// Mock implementation of referral service
// This file provides a temporary implementation while the actual service is being developed

const referralService = {
  // Statistics and analytics
  getReferralStatistics: async (timeRange = '30d') => {
    console.log('Mock getReferralStatistics called with:', timeRange);
    return {
      totalReferrals: 248,
      activeReferrals: 183,
      conversionRate: 21.5,
      totalCommissions: 12480,
      pendingCommissions: 3250,
    };
  },
  
  getReferralAnalytics: async (timeRange = '30d') => {
    console.log('Mock getReferralAnalytics called with:', timeRange);
    return {
      referralTrends: Array(7).fill().map((_, i) => ({
        date: new Date(Date.now() - (6-i) * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        count: Math.floor(Math.random() * 20) + 5
      })),
      conversionTrends: Array(7).fill().map((_, i) => ({
        date: new Date(Date.now() - (6-i) * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        rate: Math.floor(Math.random() * 15) + 10
      }))
    };
  },
  
  exportReferralData: async (options) => {
    console.log('Mock exportReferralData called with:', options);
    return { 
      success: true,
      downloadUrl: '#'
    };
  },
  
  // Commission management
  getCommissions: async (filters = {}) => {
    console.log('Mock getCommissions called with:', filters);
    return {
      commissions: Array(10).fill().map((_, i) => ({
        id: `com_${100+i}`,
        referralId: `ref_${200+i}`,
        status: ['pending', 'approved', 'paid', 'rejected'][Math.floor(Math.random() * 4)],
        amount: Math.floor(Math.random() * 200) + 50,
        date: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
        customer: `Customer ${i+1}`,
        plan: ['Basic', 'Professional', 'Enterprise'][Math.floor(Math.random() * 3)],
        paymentMethod: ['bank_transfer', 'paypal', 'stripe'][Math.floor(Math.random() * 3)],
        notes: `Commission notes ${i+1}`
      })),
      total: 56,
      page: 1,
      totalPages: 6
    };
  },
  
  updateCommissionStatus: async (id, status) => {
    console.log('Mock updateCommissionStatus called with:', id, status);
    return { success: true };
  },
  
  processCommissionPayment: async (id, paymentDetails) => {
    console.log('Mock processCommissionPayment called with:', id, paymentDetails);
    return { 
      success: true,
      transactionId: 'txn_' + Math.random().toString(36).substr(2, 9)
    };
  },
  
  bulkUpdateCommissions: async (ids, updateData) => {
    console.log('Mock bulkUpdateCommissions called with:', ids, updateData);
    return { success: true, updated: ids.length };
  },
  
  exportCommissions: async (filters) => {
    console.log('Mock exportCommissions called with:', filters);
    return { 
      success: true,
      downloadUrl: '#'
    };
  },
  
  // Referral code management
  isReferralCodeValid: (code) => {
    return code && code.length >= 4 && /^[A-Za-z0-9]+$/.test(code);
  },
  
  generateReferralCode: async (options) => {
    console.log('Mock generateReferralCode called with:', options);
    const code = options.code || ('REF' + Math.random().toString(36).substr(2, 6).toUpperCase());
    return {
      success: true,
      code,
      expiresAt: options.expiration ? new Date(options.expiration).toISOString() : null,
      maxUses: options.maxUses || null,
      discountPercentage: options.discountPercentage || 10
    };
  },
  
  generateQRCode: async (code, options = {}) => {
    console.log('Mock generateQRCode called with:', code, options);
    // Return a placeholder QR code data URL
    return {
      success: true,
      qrCodeDataUrl: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=='
    };
  },
  
  generateReferralLink: (code) => {
    return `https://example.com/referral/${code}`;
  },
  
  // This would be used by backend code
  applyReferralCode: async (code, customerId) => {
    console.log('Mock applyReferralCode called with:', code, customerId);
    return { success: true, discountApplied: true, discountAmount: 10 };
  }
};

export default referralService;
