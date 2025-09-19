import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  Copy,
  QrCode,
  Share2,
  RefreshCw,
  Calendar,
  DollarSign,
  Users,
  Link as LinkIcon,
  Download
} from 'lucide-react';
import { toast } from 'sonner';
import referralService from '@/services/referralService';

const ReferralCodeGenerator = () => {
  const [formData, setFormData] = useState({
    code: '',
    description: '',
    commissionRate: '',
    commissionType: 'percentage',
    maxUses: '',
    expiresAt: '',
    isActive: true
  });
  
  const [generatedCode, setGeneratedCode] = useState(null);
  const [qrCodeUrl, setQrCodeUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [showQRDialog, setShowQRDialog] = useState(false);

  const generateRandomCode = () => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    let result = 'REF';
    for (let i = 0; i < 6; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    setFormData(prev => ({ ...prev, code: result }));
  };

  const handleInputChange = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.code.trim()) {
      toast.error('Please enter a referral code');
      return;
    }

    if (!referralService.isReferralCodeValid(formData.code)) {
      toast.error('Invalid referral code format. Use only letters and numbers.');
      return;
    }

    setLoading(true);
    try {
      const response = await referralService.generateReferralCode({
        code: formData.code.toUpperCase(),
        description: formData.description,
        commissionRate: parseFloat(formData.commissionRate) || 0,
        commissionType: formData.commissionType,
        maxUses: formData.maxUses ? parseInt(formData.maxUses) : null,
        expiresAt: formData.expiresAt || null,
        isActive: formData.isActive
      });

      setGeneratedCode(response.data);
      toast.success('Referral code generated successfully!');
      
      // Reset form
      setFormData({
        code: '',
        description: '',
        commissionRate: '',
        commissionType: 'percentage',
        maxUses: '',
        expiresAt: '',
        isActive: true
      });
    } catch (error) {
      toast.error(error.message || 'Failed to generate referral code');
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = async (text, type = 'code') => {
    try {
      await navigator.clipboard.writeText(text);
      toast.success(`${type} copied to clipboard!`);
    } catch (error) {
      toast.error('Failed to copy to clipboard');
    }
  };

  const generateQRCode = async (code) => {
    try {
      setLoading(true);
      const response = await referralService.generateQRCode(code, {
        size: 256,
        format: 'png'
      });
      setQrCodeUrl(response.qrCodeUrl);
      setShowQRDialog(true);
    } catch (error) {
      toast.error('Failed to generate QR code');
    } finally {
      setLoading(false);
    }
  };

  const shareReferralLink = async (code) => {
    const link = referralService.generateReferralLink(code);
    
    if (navigator.share) {
      try {
        await navigator.share({
          title: 'Join with my referral code',
          text: `Use my referral code ${code} to get started!`,
          url: link
        });
      } catch (error) {
        // Fallback to clipboard
        copyToClipboard(link, 'Referral link');
      }
    } else {
      copyToClipboard(link, 'Referral link');
    }
  };

  const downloadQRCode = () => {
    if (qrCodeUrl) {
      const link = document.createElement('a');
      link.href = qrCodeUrl;
      link.download = `referral-qr-${generatedCode?.code}.png`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  };

  return (
    <div className="space-y-6">
      {/* Generator Form */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Share2 className="h-5 w-5" />
            Generate Referral Code
          </CardTitle>
          <CardDescription>
            Create new referral codes with custom settings and commission rates
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Referral Code */}
              <div className="space-y-2">
                <Label htmlFor="code">Referral Code</Label>
                <div className="flex gap-2">
                  <Input
                    id="code"
                    value={formData.code}
                    onChange={(e) => handleInputChange('code', e.target.value.toUpperCase())}
                    placeholder="Enter custom code or generate"
                    className="font-mono"
                  />
                  <Button
                    type="button"
                    variant="outline"
                    onClick={generateRandomCode}
                    className="shrink-0"
                  >
                    <RefreshCw className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              {/* Commission Rate */}
              <div className="space-y-2">
                <Label htmlFor="commissionRate">Commission Rate</Label>
                <div className="flex gap-2">
                  <Input
                    id="commissionRate"
                    type="number"
                    step="0.01"
                    min="0"
                    value={formData.commissionRate}
                    onChange={(e) => handleInputChange('commissionRate', e.target.value)}
                    placeholder="10"
                  />
                  <Select
                    value={formData.commissionType}
                    onValueChange={(value) => handleInputChange('commissionType', value)}
                  >
                    <SelectTrigger className="w-32">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="percentage">%</SelectItem>
                      <SelectItem value="fixed">৳</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              {/* Max Uses */}
              <div className="space-y-2">
                <Label htmlFor="maxUses">Max Uses (Optional)</Label>
                <Input
                  id="maxUses"
                  type="number"
                  min="1"
                  value={formData.maxUses}
                  onChange={(e) => handleInputChange('maxUses', e.target.value)}
                  placeholder="Unlimited"
                />
              </div>

              {/* Expiry Date */}
              <div className="space-y-2">
                <Label htmlFor="expiresAt">Expiry Date (Optional)</Label>
                <Input
                  id="expiresAt"
                  type="datetime-local"
                  value={formData.expiresAt}
                  onChange={(e) => handleInputChange('expiresAt', e.target.value)}
                />
              </div>
            </div>

            {/* Description */}
            <div className="space-y-2">
              <Label htmlFor="description">Description (Optional)</Label>
              <Textarea
                id="description"
                value={formData.description}
                onChange={(e) => handleInputChange('description', e.target.value)}
                placeholder="Describe the purpose of this referral code..."
                rows={3}
              />
            </div>

            <Button type="submit" disabled={loading} className="w-full">
              {loading ? (
                <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
              ) : (
                <Share2 className="h-4 w-4 mr-2" />
              )}
              Generate Referral Code
            </Button>
          </form>
        </CardContent>
      </Card>

      {/* Generated Code Display */}
      {generatedCode && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Badge variant="default" className="text-sm">
                Generated
              </Badge>
              Referral Code: {generatedCode.code}
            </CardTitle>
            <CardDescription>
              Your referral code has been created successfully
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Code Details */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="flex items-center gap-2 p-3 bg-muted rounded-lg">
                <DollarSign className="h-4 w-4 text-green-500" />
                <div>
                  <p className="text-sm font-medium">
                    {generatedCode.commissionRate}{
                      generatedCode.commissionType === 'percentage' ? '%' : '৳'
                    } Commission
                  </p>
                  <p className="text-xs text-muted-foreground">Rate</p>
                </div>
              </div>
              
              <div className="flex items-center gap-2 p-3 bg-muted rounded-lg">
                <Users className="h-4 w-4 text-blue-500" />
                <div>
                  <p className="text-sm font-medium">
                    {generatedCode.maxUses || 'Unlimited'}
                  </p>
                  <p className="text-xs text-muted-foreground">Max Uses</p>
                </div>
              </div>
              
              <div className="flex items-center gap-2 p-3 bg-muted rounded-lg">
                <Calendar className="h-4 w-4 text-orange-500" />
                <div>
                  <p className="text-sm font-medium">
                    {generatedCode.expiresAt 
                      ? new Date(generatedCode.expiresAt).toLocaleDateString()
                      : 'No Expiry'
                    }
                  </p>
                  <p className="text-xs text-muted-foreground">Expires</p>
                </div>
              </div>
            </div>

            {/* Referral Link */}
            <div className="space-y-2">
              <Label>Referral Link</Label>
              <div className="flex gap-2">
                <Input
                  value={referralService.generateReferralLink(generatedCode.code)}
                  readOnly
                  className="font-mono text-sm"
                />
                <Button
                  variant="outline"
                  onClick={() => copyToClipboard(
                    referralService.generateReferralLink(generatedCode.code),
                    'Referral link'
                  )}
                >
                  <Copy className="h-4 w-4" />
                </Button>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-wrap gap-2">
              <Button
                variant="outline"
                onClick={() => copyToClipboard(generatedCode.code, 'Referral code')}
              >
                <Copy className="h-4 w-4 mr-2" />
                Copy Code
              </Button>
              
              <Button
                variant="outline"
                onClick={() => generateQRCode(generatedCode.code)}
              >
                <QrCode className="h-4 w-4 mr-2" />
                Generate QR
              </Button>
              
              <Button
                variant="outline"
                onClick={() => shareReferralLink(generatedCode.code)}
              >
                <LinkIcon className="h-4 w-4 mr-2" />
                Share Link
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* QR Code Dialog */}
      <Dialog open={showQRDialog} onOpenChange={setShowQRDialog}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>QR Code for {generatedCode?.code}</DialogTitle>
            <DialogDescription>
              Scan this QR code to access the referral link
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col items-center space-y-4">
            {qrCodeUrl && (
              <img
                src={qrCodeUrl}
                alt="Referral QR Code"
                className="w-64 h-64 border rounded-lg"
              />
            )}
            <div className="flex gap-2">
              <Button onClick={downloadQRCode} variant="outline">
                <Download className="h-4 w-4 mr-2" />
                Download
              </Button>
              <Button
                onClick={() => copyToClipboard(qrCodeUrl, 'QR code URL')}
                variant="outline"
              >
                <Copy className="h-4 w-4 mr-2" />
                Copy URL
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default ReferralCodeGenerator;