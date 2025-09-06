import React, { forwardRef } from 'react';

const InvoicePDF = forwardRef(({ order, companyInfo }, ref) => {
  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount || 0);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const formatDateTime = (dateString) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (!order) {
    return <div>No order data available</div>;
  }

  const defaultCompanyInfo = {
    name: "ShopVendor",
    subtitle: "Single Vendor E-commerce",
    email: "admin@shopvendor.com",
    phone: "(555) 123-4567",
    website: "www.shopvendor.com",
    address: {
      street: "123 Business Street",
      city: "Commerce City",
      state: "CC",
      postal_code: "12345",
      country: "United States"
    },
    ...companyInfo
  };

  return (
    <div ref={ref} className="bg-white p-8 max-w-4xl mx-auto" style={{ fontFamily: 'Arial, sans-serif' }}>
      {/* Header */}
      <div className="flex justify-between items-start mb-8 pb-6 border-b-2 border-gray-800">
        <div className="flex-1">
          <h1 className="text-3xl font-bold text-gray-800 mb-2">{defaultCompanyInfo.name}</h1>
          <p className="text-lg text-gray-600 mb-3">{defaultCompanyInfo.subtitle}</p>
          <div className="text-sm text-gray-600 space-y-1">
            <p>📧 {defaultCompanyInfo.email} | 📞 {defaultCompanyInfo.phone}</p>
            <p>🌐 {defaultCompanyInfo.website}</p>
            <p>
              {defaultCompanyInfo.address.street}, {defaultCompanyInfo.address.city}, {defaultCompanyInfo.address.state} {defaultCompanyInfo.address.postal_code}
            </p>
          </div>
        </div>
        <div className="text-right">
          <h2 className="text-2xl font-bold text-gray-800 mb-3">INVOICE</h2>
          <div className="text-sm space-y-1">
            <p><span className="font-semibold">Order ID:</span> {order.id}</p>
            <p><span className="font-semibold">Invoice Date:</span> {formatDate(order.date_created || order.created_at)}</p>
            <p><span className="font-semibold">Status:</span> <span className="capitalize font-medium">{order.status}</span></p>
            {order.payment_status && (
              <p><span className="font-semibold">Payment:</span> <span className="capitalize">{order.payment_status}</span></p>
            )}
          </div>
        </div>
      </div>

      {/* Customer & Billing Information */}
      <div className="grid grid-cols-2 gap-8 mb-8">
        <div>
          <h3 className="text-lg font-bold text-gray-800 mb-3">Bill To:</h3>
          <div className="space-y-1 text-sm">
            <p className="font-semibold">{order.customer?.name || order.user?.name || 'Customer Name'}</p>
            <p>{order.customer?.email || order.user?.email || 'customer@email.com'}</p>
            {order.customer?.phone && <p>{order.customer.phone}</p>}
            {order.billing?.address && (
              <>
                <p className="mt-2 font-medium">Billing Address:</p>
                <p>{order.billing.address.street}</p>
                {order.billing.address.apartment && <p>{order.billing.address.apartment}</p>}
                <p>{order.billing.address.city}, {order.billing.address.state} {order.billing.address.postal_code}</p>
                <p>{order.billing.address.country}</p>
              </>
            )}
          </div>
        </div>
        <div>
          <h3 className="text-lg font-bold text-gray-800 mb-3">Ship To:</h3>
          <div className="space-y-1 text-sm">
            {order.shipping?.address ? (
              <>
                <p>{order.shipping.address.street}</p>
                {order.shipping.address.apartment && <p>{order.shipping.address.apartment}</p>}
                <p>{order.shipping.address.city}, {order.shipping.address.state} {order.shipping.address.postal_code}</p>
                <p>{order.shipping.address.country}</p>
                {order.shipping.method && (
                  <>
                    <p className="mt-2 font-medium">Shipping Method:</p>
                    <p>{order.shipping.method} {order.shipping.carrier && `- ${order.shipping.carrier}`}</p>
                  </>
                )}
                {order.shipping.tracking_number && (
                  <p><span className="font-medium">Tracking:</span> {order.shipping.tracking_number}</p>
                )}
              </>
            ) : order.shipping_address ? (
              <p>{order.shipping_address}</p>
            ) : (
              <p>Same as billing address</p>
            )}
          </div>
          
          {order.billing?.payment_method && (
            <div className="mt-4">
              <h4 className="font-medium text-gray-800 mb-2">Payment Method:</h4>
              <div className="text-sm space-y-1">
                <p>{order.billing.payment_method}</p>
                {order.billing.transaction_id && (
                  <p><span className="font-medium">Transaction:</span> {order.billing.transaction_id}</p>
                )}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Items Table */}
      <div className="mb-8">
        <table className="w-full border-collapse border border-gray-300">
          <thead>
            <tr className="bg-gray-100">
              <th className="border border-gray-300 px-4 py-3 text-left font-bold">Item</th>
              <th className="border border-gray-300 px-4 py-3 text-left font-bold">SKU</th>
              <th className="border border-gray-300 px-4 py-3 text-left font-bold">Variant</th>
              <th className="border border-gray-300 px-4 py-3 text-center font-bold">Qty</th>
              <th className="border border-gray-300 px-4 py-3 text-right font-bold">Unit Price</th>
              <th className="border border-gray-300 px-4 py-3 text-right font-bold">Total</th>
            </tr>
          </thead>
          <tbody>
            {(order.items || order.order_items || []).map((item, index) => (
              <tr key={item.id || index}>
                <td className="border border-gray-300 px-4 py-3">{item.name || item.product_name || 'Product'}</td>
                <td className="border border-gray-300 px-4 py-3 text-sm text-gray-600">
                  {item.sku || item.product_sku || '-'}
                </td>
                <td className="border border-gray-300 px-4 py-3 text-sm text-gray-600">
                  {item.variant || item.product_variant || '-'}
                </td>
                <td className="border border-gray-300 px-4 py-3 text-center">
                  {item.qty || item.quantity || 1}
                </td>
                <td className="border border-gray-300 px-4 py-3 text-right">
                  {formatCurrency(item.price || item.unit_price || 0)}
                </td>
                <td className="border border-gray-300 px-4 py-3 text-right font-medium">
                  {formatCurrency((item.price || item.unit_price || 0) * (item.qty || item.quantity || 1))}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Totals Section */}
      <div className="flex justify-end mb-8">
        <div className="w-1/2 space-y-2">
          <div className="flex justify-between py-1">
            <span>Subtotal:</span>
            <span>{formatCurrency(order.subtotal || order.total_price || 0)}</span>
          </div>
          {order.discount > 0 && (
            <div className="flex justify-between py-1 text-green-600">
              <span>Discount:</span>
              <span>-{formatCurrency(order.discount || order.discount_amount || 0)}</span>
            </div>
          )}
          {order.tax > 0 && (
            <div className="flex justify-between py-1">
              <span>Tax:</span>
              <span>{formatCurrency(order.tax || 0)}</span>
            </div>
          )}
          {order.shipping_cost > 0 && (
            <div className="flex justify-between py-1">
              <span>Shipping:</span>
              <span>{formatCurrency(order.shipping_cost || 0)}</span>
            </div>
          )}
          <div className="border-t border-gray-300 pt-2">
            <div className="flex justify-between py-2 text-lg font-bold">
              <span>Total:</span>
              <span>{formatCurrency(order.total || order.final_price || order.total_amount || 0)}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Notes Section */}
      {(order.notes || order.internal_notes) && (
        <div className="mb-8">
          {order.notes && (
            <div className="mb-4">
              <h4 className="font-bold text-gray-800 mb-2">Customer Notes:</h4>
              <p className="text-sm text-gray-600 bg-gray-50 p-3 rounded border">{order.notes}</p>
            </div>
          )}
          {order.internal_notes && (
            <div>
              <h4 className="font-bold text-gray-800 mb-2">Internal Notes:</h4>
              <p className="text-sm text-gray-600 bg-gray-50 p-3 rounded border">{order.internal_notes}</p>
            </div>
          )}
        </div>
      )}

      {/* Order Timeline */}
      {order.timeline && order.timeline.length > 0 && (
        <div className="mb-8">
          <h4 className="font-bold text-gray-800 mb-3">Order Timeline:</h4>
          <div className="space-y-2">
            {order.timeline.map((event, index) => (
              <div key={index} className="flex justify-between items-center py-1 text-sm">
                <div>
                  <span className="font-medium capitalize">{event.status}</span>
                  <span className="text-gray-600 ml-2">{event.description}</span>
                </div>
                <div className="text-gray-500 text-xs">
                  {formatDateTime(event.timestamp)} • {event.user}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Footer */}
      <div className="text-center text-gray-600 text-sm pt-8 border-t border-gray-300">
        <p className="font-medium mb-2">Thank you for your business!</p>
        <p>For support or questions, contact us at {defaultCompanyInfo.email}</p>
        <p className="mt-2 text-xs">This invoice was generated on {formatDateTime(new Date().toISOString())}</p>
      </div>
    </div>
  );
});

InvoicePDF.displayName = 'InvoicePDF';

export default InvoicePDF;