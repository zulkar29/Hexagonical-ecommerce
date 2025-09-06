import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';

export const generatePDFFromHTML = async (element, filename = 'invoice.pdf') => {
  try {
    // Configure html2canvas options for better quality
    const canvas = await html2canvas(element, {
      scale: 2, // Higher scale for better quality
      useCORS: true,
      allowTaint: false,
      backgroundColor: '#ffffff',
      logging: false,
      onclone: (clonedDoc) => {
        // Ensure all styles are properly applied in the cloned document
        const clonedElement = clonedDoc.querySelector('[data-pdf-content]') || clonedDoc.body;
        if (clonedElement) {
          clonedElement.style.width = '210mm'; // A4 width
          clonedElement.style.padding = '20mm';
          clonedElement.style.fontSize = '12px';
          clonedElement.style.lineHeight = '1.4';
        }
      }
    });

    const imgData = canvas.toDataURL('image/png');
    const pdf = new jsPDF('p', 'mm', 'a4');
    
    const pdfWidth = pdf.internal.pageSize.getWidth();
    const pdfHeight = pdf.internal.pageSize.getHeight();
    const imgWidth = canvas.width;
    const imgHeight = canvas.height;
    
    const ratio = Math.min(pdfWidth / imgWidth, pdfHeight / imgHeight);
    const imgX = (pdfWidth - imgWidth * ratio) / 2;
    const imgY = 0;
    
    pdf.addImage(imgData, 'PNG', imgX, imgY, imgWidth * ratio, imgHeight * ratio);
    
    // If content is too long, add new pages
    const totalPDFPages = Math.ceil((imgHeight * ratio) / pdfHeight);
    for (let i = 1; i < totalPDFPages; i++) {
      pdf.addPage();
      pdf.addImage(
        imgData,
        'PNG',
        imgX,
        -(pdfHeight * i) + imgY,
        imgWidth * ratio,
        imgHeight * ratio
      );
    }
    
    pdf.save(filename);
    return { success: true, pdf };
  } catch (error) {
    console.error('PDF generation failed:', error);
    throw new Error(`Failed to generate PDF: ${error.message}`);
  }
};

export const generateInvoicePDF = async (order, companyInfo = {}, filename) => {
  // Create a temporary container for the PDF content
  const tempContainer = document.createElement('div');
  tempContainer.style.position = 'absolute';
  tempContainer.style.left = '-9999px';
  tempContainer.style.width = '210mm'; // A4 width
  tempContainer.style.backgroundColor = 'white';
  tempContainer.setAttribute('data-pdf-content', 'true');
  
  const defaultFilename = `invoice-${order.id || 'unknown'}-${new Date().toISOString().split('T')[0]}.pdf`;
  
  // Generate invoice HTML content
  tempContainer.innerHTML = generateInvoiceHTML(order, companyInfo);
  
  document.body.appendChild(tempContainer);
  
  try {
    await generatePDFFromHTML(tempContainer, filename || defaultFilename);
    return { success: true };
  } finally {
    // Clean up
    document.body.removeChild(tempContainer);
  }
};

const generateInvoiceHTML = (order, companyInfo = {}) => {
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

  return `
    <div style="font-family: Arial, sans-serif; padding: 20px; max-width: 800px; margin: 0 auto; background: white; color: #333;">
      <!-- Header -->
      <div style="display: flex; justify-content: space-between; margin-bottom: 30px; padding-bottom: 20px; border-bottom: 2px solid #333;">
        <div style="flex: 1;">
          <h1 style="margin: 0 0 8px 0; font-size: 28px; font-weight: bold; color: #333;">${defaultCompanyInfo.name}</h1>
          <p style="margin: 0 0 12px 0; font-size: 16px; color: #666;">${defaultCompanyInfo.subtitle}</p>
          <div style="font-size: 12px; color: #666; line-height: 1.4;">
            <p style="margin: 0;">📧 ${defaultCompanyInfo.email} | 📞 ${defaultCompanyInfo.phone}</p>
            <p style="margin: 4px 0 0 0;">🌐 ${defaultCompanyInfo.website}</p>
            <p style="margin: 4px 0 0 0;">
              ${defaultCompanyInfo.address.street}, ${defaultCompanyInfo.address.city}, ${defaultCompanyInfo.address.state} ${defaultCompanyInfo.address.postal_code}
            </p>
          </div>
        </div>
        <div style="text-align: right;">
          <h2 style="margin: 0 0 12px 0; font-size: 24px; font-weight: bold; color: #333;">INVOICE</h2>
          <div style="font-size: 12px; line-height: 1.5;">
            <p style="margin: 0;"><strong>Order ID:</strong> ${order.id}</p>
            <p style="margin: 4px 0;"><strong>Invoice Date:</strong> ${formatDate(order.date_created || order.created_at)}</p>
            <p style="margin: 4px 0;"><strong>Status:</strong> <span style="text-transform: capitalize; font-weight: 600;">${order.status}</span></p>
            ${order.payment_status ? `<p style="margin: 4px 0;"><strong>Payment:</strong> <span style="text-transform: capitalize;">${order.payment_status}</span></p>` : ''}
          </div>
        </div>
      </div>

      <!-- Customer & Billing Information -->
      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-bottom: 30px;">
        <div>
          <h3 style="margin: 0 0 12px 0; font-size: 16px; font-weight: bold; color: #333;">Bill To:</h3>
          <div style="font-size: 12px; line-height: 1.4; color: #333;">
            <p style="margin: 0 0 4px 0; font-weight: bold;">${order.customer?.name || order.user?.name || 'Customer Name'}</p>
            <p style="margin: 0 0 4px 0;">${order.customer?.email || order.user?.email || 'customer@email.com'}</p>
            ${order.customer?.phone ? `<p style="margin: 0 0 4px 0;">${order.customer.phone}</p>` : ''}
            ${order.billing?.address ? `
              <p style="margin: 8px 0 4px 0; font-weight: 600;">Billing Address:</p>
              <p style="margin: 0 0 2px 0;">${order.billing.address.street}</p>
              ${order.billing.address.apartment ? `<p style="margin: 0 0 2px 0;">${order.billing.address.apartment}</p>` : ''}
              <p style="margin: 0 0 2px 0;">${order.billing.address.city}, ${order.billing.address.state} ${order.billing.address.postal_code}</p>
              <p style="margin: 0 0 2px 0;">${order.billing.address.country}</p>
            ` : ''}
          </div>
        </div>
        <div>
          <h3 style="margin: 0 0 12px 0; font-size: 16px; font-weight: bold; color: #333;">Ship To:</h3>
          <div style="font-size: 12px; line-height: 1.4; color: #333;">
            ${order.shipping?.address ? `
              <p style="margin: 0 0 2px 0;">${order.shipping.address.street}</p>
              ${order.shipping.address.apartment ? `<p style="margin: 0 0 2px 0;">${order.shipping.address.apartment}</p>` : ''}
              <p style="margin: 0 0 2px 0;">${order.shipping.address.city}, ${order.shipping.address.state} ${order.shipping.address.postal_code}</p>
              <p style="margin: 0 0 4px 0;">${order.shipping.address.country}</p>
              ${order.shipping.method ? `
                <p style="margin: 8px 0 4px 0; font-weight: 600;">Shipping Method:</p>
                <p style="margin: 0 0 2px 0;">${order.shipping.method}${order.shipping.carrier ? ` - ${order.shipping.carrier}` : ''}</p>
              ` : ''}
              ${order.shipping.tracking_number ? `<p style="margin: 0 0 2px 0;"><strong>Tracking:</strong> ${order.shipping.tracking_number}</p>` : ''}
            ` : order.shipping_address ? `
              <p style="margin: 0 0 4px 0;">${order.shipping_address}</p>
            ` : `
              <p style="margin: 0 0 4px 0;">Same as billing address</p>
            `}
          </div>
          
          ${order.billing?.payment_method ? `
            <div style="margin-top: 16px;">
              <h4 style="margin: 0 0 8px 0; font-weight: 600; color: #333;">Payment Method:</h4>
              <div style="font-size: 12px; line-height: 1.4;">
                <p style="margin: 0 0 4px 0;">${order.billing.payment_method}</p>
                ${order.billing.transaction_id ? `<p style="margin: 0 0 4px 0;"><strong>Transaction:</strong> ${order.billing.transaction_id}</p>` : ''}
              </div>
            </div>
          ` : ''}
        </div>
      </div>

      <!-- Items Table -->
      <div style="margin-bottom: 30px;">
        <table style="width: 100%; border-collapse: collapse; border: 1px solid #ddd;">
          <thead>
            <tr style="background-color: #f8f9fa;">
              <th style="border: 1px solid #ddd; padding: 12px; text-align: left; font-weight: bold;">Item</th>
              <th style="border: 1px solid #ddd; padding: 12px; text-align: left; font-weight: bold;">SKU</th>
              <th style="border: 1px solid #ddd; padding: 12px; text-align: left; font-weight: bold;">Variant</th>
              <th style="border: 1px solid #ddd; padding: 12px; text-align: center; font-weight: bold;">Qty</th>
              <th style="border: 1px solid #ddd; padding: 12px; text-align: right; font-weight: bold;">Unit Price</th>
              <th style="border: 1px solid #ddd; padding: 12px; text-align: right; font-weight: bold;">Total</th>
            </tr>
          </thead>
          <tbody>
            ${(order.items || order.order_items || []).map(item => `
              <tr>
                <td style="border: 1px solid #ddd; padding: 12px;">${item.name || item.product_name || 'Product'}</td>
                <td style="border: 1px solid #ddd; padding: 12px; font-size: 11px; color: #666;">
                  ${item.sku || item.product_sku || '-'}
                </td>
                <td style="border: 1px solid #ddd; padding: 12px; font-size: 11px; color: #666;">
                  ${item.variant || item.product_variant || '-'}
                </td>
                <td style="border: 1px solid #ddd; padding: 12px; text-align: center;">
                  ${item.qty || item.quantity || 1}
                </td>
                <td style="border: 1px solid #ddd; padding: 12px; text-align: right;">
                  ${formatCurrency(item.price || item.unit_price || 0)}
                </td>
                <td style="border: 1px solid #ddd; padding: 12px; text-align: right; font-weight: 600;">
                  ${formatCurrency((item.price || item.unit_price || 0) * (item.qty || item.quantity || 1))}
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>

      <!-- Totals Section -->
      <div style="display: flex; justify-content: flex-end; margin-bottom: 30px;">
        <div style="width: 300px;">
          <div style="display: flex; justify-content: space-between; padding: 4px 0;">
            <span>Subtotal:</span>
            <span>${formatCurrency(order.subtotal || order.total_price || 0)}</span>
          </div>
          ${(order.discount || order.discount_amount) > 0 ? `
            <div style="display: flex; justify-content: space-between; padding: 4px 0; color: #059669;">
              <span>Discount:</span>
              <span>-${formatCurrency(order.discount || order.discount_amount || 0)}</span>
            </div>
          ` : ''}
          ${order.tax > 0 ? `
            <div style="display: flex; justify-content: space-between; padding: 4px 0;">
              <span>Tax:</span>
              <span>${formatCurrency(order.tax || 0)}</span>
            </div>
          ` : ''}
          ${order.shipping_cost > 0 ? `
            <div style="display: flex; justify-content: space-between; padding: 4px 0;">
              <span>Shipping:</span>
              <span>${formatCurrency(order.shipping_cost || 0)}</span>
            </div>
          ` : ''}
          <div style="border-top: 1px solid #333; padding-top: 8px; margin-top: 8px;">
            <div style="display: flex; justify-content: space-between; padding: 8px 0; font-size: 18px; font-weight: bold;">
              <span>Total:</span>
              <span>${formatCurrency(order.total || order.final_price || order.total_amount || 0)}</span>
            </div>
          </div>
        </div>
      </div>

      ${(order.notes || order.internal_notes) ? `
        <!-- Notes Section -->
        <div style="margin-bottom: 30px;">
          ${order.notes ? `
            <div style="margin-bottom: 16px;">
              <h4 style="margin: 0 0 8px 0; font-weight: bold; color: #333;">Customer Notes:</h4>
              <p style="margin: 0; font-size: 12px; color: #666; background-color: #f8f9fa; padding: 12px; border-radius: 4px; border: 1px solid #e9ecef;">${order.notes}</p>
            </div>
          ` : ''}
          ${order.internal_notes ? `
            <div>
              <h4 style="margin: 0 0 8px 0; font-weight: bold; color: #333;">Internal Notes:</h4>
              <p style="margin: 0; font-size: 12px; color: #666; background-color: #f8f9fa; padding: 12px; border-radius: 4px; border: 1px solid #e9ecef;">${order.internal_notes}</p>
            </div>
          ` : ''}
        </div>
      ` : ''}

      ${order.timeline && order.timeline.length > 0 ? `
        <!-- Order Timeline -->
        <div style="margin-bottom: 30px;">
          <h4 style="margin: 0 0 12px 0; font-weight: bold; color: #333;">Order Timeline:</h4>
          <div>
            ${order.timeline.map(event => `
              <div style="display: flex; justify-content: space-between; align-items: center; padding: 4px 0; font-size: 12px;">
                <div>
                  <span style="font-weight: 600; text-transform: capitalize;">${event.status}</span>
                  <span style="color: #666; margin-left: 8px;">${event.description}</span>
                </div>
                <div style="color: #999; font-size: 10px;">
                  ${formatDateTime(event.timestamp)} • ${event.user}
                </div>
              </div>
            `).join('')}
          </div>
        </div>
      ` : ''}

      <!-- Footer -->
      <div style="text-align: center; color: #666; font-size: 12px; padding-top: 30px; border-top: 1px solid #ddd;">
        <p style="margin: 0 0 8px 0; font-weight: 600;">Thank you for your business!</p>
        <p style="margin: 0 0 8px 0;">For support or questions, contact us at ${defaultCompanyInfo.email}</p>
        <p style="margin: 0; font-size: 10px;">This invoice was generated on ${formatDateTime(new Date().toISOString())}</p>
      </div>
    </div>
  `;
};

export default { generatePDFFromHTML, generateInvoicePDF };