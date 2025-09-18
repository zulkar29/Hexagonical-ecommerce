import React, { useState } from 'react';
import {
  ArrowLeft,
  Building2,
  User,
  Mail,
  Phone,
  MapPin,
  CreditCard,
  Package,
  Settings,
  Save,
  X,
  CheckCircle,
  Clock,
  AlertCircle,
  Upload,
  Eye,
  EyeOff
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Link } from 'react-router-dom';

const onboardingSteps = [
  { id: 'company', title: 'Company Information', icon: Building2 },
  { id: 'admin', title: 'Admin User', icon: User },
  { id: 'subscription', title: 'Subscription Plan', icon: Package },
  { id: 'settings', title: 'Initial Settings', icon: Settings },
  { id: 'review', title: 'Review & Launch', icon: CheckCircle }
];

const subscriptionPlans = [
  {
    id: 'starter',
    name: 'Starter Plan',
    description: 'Perfect for small businesses getting started',
    price: 29,
    features: ['Up to 100 products', 'Basic analytics', 'Email support', '5GB storage'],
    recommended: false
  },
  {
    id: 'professional',
    name: 'Professional Plan',
    description: 'Advanced features for growing businesses',
    price: 79,
    features: ['Up to 1000 products', 'Advanced analytics', 'Priority support', '50GB storage', 'API access'],
    recommended: true
  },
  {
    id: 'enterprise',
    name: 'Enterprise Plan',
    description: 'Full-featured solution for large organizations',
    price: 199,
    features: ['Unlimited products', 'Custom analytics', '24/7 support', 'Unlimited storage', 'Full API access', 'Custom integrations'],
    recommended: false
  }
];

const industries = [
  'Technology', 'Retail', 'Healthcare', 'Finance', 'Education', 'Manufacturing',
  'Real Estate', 'Food & Beverage', 'Fashion', 'Automotive', 'Other'
];

const countries = [
  'United States', 'Canada', 'United Kingdom', 'Germany', 'France', 'Australia',
  'Japan', 'Singapore', 'Netherlands', 'Sweden', 'Other'
];

export default function TenantOnboarding() {
  const [currentStep, setCurrentStep] = useState(0);
  const [showPassword, setShowPassword] = useState(false);
  const [formData, setFormData] = useState({
    // Company Information
    companyName: '',
    companyEmail: '',
    companyPhone: '',
    website: '',
    industry: '',
    companySize: '',
    description: '',
    
    // Address
    address: '',
    city: '',
    state: '',
    zipCode: '',
    country: '',
    
    // Admin User
    adminFirstName: '',
    adminLastName: '',
    adminEmail: '',
    adminPhone: '',
    adminPassword: '',
    adminConfirmPassword: '',
    
    // Subscription
    selectedPlan: '',
    billingCycle: 'monthly',
    trialDays: '14',
    
    // Settings
    timezone: '',
    currency: 'USD',
    language: 'en',
    enableNotifications: true,
    enableAnalytics: true,
    customDomain: '',
    
    // Additional
    notes: '',
    sendWelcomeEmail: true
  });

  const handleInputChange = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const getCurrentStepProgress = () => {
    return ((currentStep + 1) / onboardingSteps.length) * 100;
  };

  const validateCurrentStep = () => {
    switch (currentStep) {
      case 0: // Company Information
        return formData.companyName && formData.companyEmail && formData.industry;
      case 1: // Admin User
        return formData.adminFirstName && formData.adminLastName && formData.adminEmail && 
               formData.adminPassword && formData.adminPassword === formData.adminConfirmPassword;
      case 2: // Subscription
        return formData.selectedPlan;
      case 3: // Settings
        return formData.timezone && formData.currency;
      case 4: // Review
        return true;
      default:
        return false;
    }
  };

  const handleNext = () => {
    if (validateCurrentStep() && currentStep < onboardingSteps.length - 1) {
      setCurrentStep(currentStep + 1);
    }
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleSubmit = () => {
    console.log('Onboarding tenant:', formData);
    // Handle tenant creation
  };

  const renderStepContent = () => {
    switch (currentStep) {
      case 0:
        return (
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Building2 className="h-5 w-5" />
                  Company Details
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="companyName">Company Name *</Label>
                    <Input
                      id="companyName"
                      value={formData.companyName}
                      onChange={(e) => handleInputChange('companyName', e.target.value)}
                      placeholder="Enter company name"
                    />
                  </div>
                  <div>
                    <Label htmlFor="companyEmail">Company Email *</Label>
                    <Input
                      id="companyEmail"
                      type="email"
                      value={formData.companyEmail}
                      onChange={(e) => handleInputChange('companyEmail', e.target.value)}
                      placeholder="company@example.com"
                    />
                  </div>
                  <div>
                    <Label htmlFor="companyPhone">Phone Number</Label>
                    <Input
                      id="companyPhone"
                      value={formData.companyPhone}
                      onChange={(e) => handleInputChange('companyPhone', e.target.value)}
                      placeholder="+1 (555) 123-4567"
                    />
                  </div>
                  <div>
                    <Label htmlFor="website">Website</Label>
                    <Input
                      id="website"
                      value={formData.website}
                      onChange={(e) => handleInputChange('website', e.target.value)}
                      placeholder="https://company.com"
                    />
                  </div>
                  <div>
                    <Label htmlFor="industry">Industry *</Label>
                    <Select value={formData.industry} onValueChange={(value) => handleInputChange('industry', value)}>
                      <SelectTrigger>
                        <SelectValue placeholder="Select industry" />
                      </SelectTrigger>
                      <SelectContent>
                        {industries.map((industry) => (
                          <SelectItem key={industry} value={industry.toLowerCase()}>
                            {industry}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div>
                    <Label htmlFor="companySize">Company Size</Label>
                    <Select value={formData.companySize} onValueChange={(value) => handleInputChange('companySize', value)}>
                      <SelectTrigger>
                        <SelectValue placeholder="Select company size" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="1-10">1-10 employees</SelectItem>
                        <SelectItem value="11-50">11-50 employees</SelectItem>
                        <SelectItem value="51-200">51-200 employees</SelectItem>
                        <SelectItem value="201-1000">201-1000 employees</SelectItem>
                        <SelectItem value="1000+">1000+ employees</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div>
                  <Label htmlFor="description">Company Description</Label>
                  <Textarea
                    id="description"
                    value={formData.description}
                    onChange={(e) => handleInputChange('description', e.target.value)}
                    placeholder="Brief description of your company and business..."
                    rows={3}
                  />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  Company Address
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="address">Street Address</Label>
                  <Input
                    id="address"
                    value={formData.address}
                    onChange={(e) => handleInputChange('address', e.target.value)}
                    placeholder="123 Business Street"
                  />
                </div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <Label htmlFor="city">City</Label>
                    <Input
                      id="city"
                      value={formData.city}
                      onChange={(e) => handleInputChange('city', e.target.value)}
                      placeholder="City"
                    />
                  </div>
                  <div>
                    <Label htmlFor="state">State/Province</Label>
                    <Input
                      id="state"
                      value={formData.state}
                      onChange={(e) => handleInputChange('state', e.target.value)}
                      placeholder="State"
                    />
                  </div>
                  <div>
                    <Label htmlFor="zipCode">ZIP/Postal Code</Label>
                    <Input
                      id="zipCode"
                      value={formData.zipCode}
                      onChange={(e) => handleInputChange('zipCode', e.target.value)}
                      placeholder="12345"
                    />
                  </div>
                </div>
                <div>
                  <Label htmlFor="country">Country</Label>
                  <Select value={formData.country} onValueChange={(value) => handleInputChange('country', value)}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select country" />
                    </SelectTrigger>
                    <SelectContent>
                      {countries.map((country) => (
                        <SelectItem key={country} value={country.toLowerCase().replace(' ', '_')}>
                          {country}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </CardContent>
            </Card>
          </div>
        );

      case 1:
        return (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-5 w-5" />
                Admin User Setup
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>
                  This user will have full administrative access to the tenant account.
                </AlertDescription>
              </Alert>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="adminFirstName">First Name *</Label>
                  <Input
                    id="adminFirstName"
                    value={formData.adminFirstName}
                    onChange={(e) => handleInputChange('adminFirstName', e.target.value)}
                    placeholder="John"
                  />
                </div>
                <div>
                  <Label htmlFor="adminLastName">Last Name *</Label>
                  <Input
                    id="adminLastName"
                    value={formData.adminLastName}
                    onChange={(e) => handleInputChange('adminLastName', e.target.value)}
                    placeholder="Doe"
                  />
                </div>
                <div>
                  <Label htmlFor="adminEmail">Email Address *</Label>
                  <Input
                    id="adminEmail"
                    type="email"
                    value={formData.adminEmail}
                    onChange={(e) => handleInputChange('adminEmail', e.target.value)}
                    placeholder="admin@company.com"
                  />
                </div>
                <div>
                  <Label htmlFor="adminPhone">Phone Number</Label>
                  <Input
                    id="adminPhone"
                    value={formData.adminPhone}
                    onChange={(e) => handleInputChange('adminPhone', e.target.value)}
                    placeholder="+1 (555) 123-4567"
                  />
                </div>
                <div>
                  <Label htmlFor="adminPassword">Password *</Label>
                  <div className="relative">
                    <Input
                      id="adminPassword"
                      type={showPassword ? 'text' : 'password'}
                      value={formData.adminPassword}
                      onChange={(e) => handleInputChange('adminPassword', e.target.value)}
                      placeholder="Enter secure password"
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                    </Button>
                  </div>
                </div>
                <div>
                  <Label htmlFor="adminConfirmPassword">Confirm Password *</Label>
                  <Input
                    id="adminConfirmPassword"
                    type="password"
                    value={formData.adminConfirmPassword}
                    onChange={(e) => handleInputChange('adminConfirmPassword', e.target.value)}
                    placeholder="Confirm password"
                  />
                </div>
              </div>
              
              {formData.adminPassword && formData.adminConfirmPassword && 
               formData.adminPassword !== formData.adminConfirmPassword && (
                <Alert>
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    Passwords do not match.
                  </AlertDescription>
                </Alert>
              )}
            </CardContent>
          </Card>
        );

      case 2:
        return (
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Package className="h-5 w-5" />
                  Choose Subscription Plan
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {subscriptionPlans.map((plan) => (
                    <div
                      key={plan.id}
                      className={`border rounded-lg p-4 cursor-pointer transition-colors relative ${
                        formData.selectedPlan === plan.id
                          ? 'border-primary bg-primary/5'
                          : 'border-border hover:border-primary/50'
                      }`}
                      onClick={() => handleInputChange('selectedPlan', plan.id)}
                    >
                      {plan.recommended && (
                        <Badge className="absolute -top-2 left-4 bg-primary text-primary-foreground">
                          Recommended
                        </Badge>
                      )}
                      <div className="flex items-center justify-between mb-2">
                        <h3 className="font-semibold">{plan.name}</h3>
                        <div className="text-right">
                          <div className="text-2xl font-bold">${plan.price}</div>
                          <div className="text-sm text-muted-foreground">/month</div>
                        </div>
                      </div>
                      <p className="text-sm text-muted-foreground mb-3">{plan.description}</p>
                      <ul className="space-y-1">
                        {plan.features.map((feature, index) => (
                          <li key={index} className="text-sm flex items-center gap-2">
                            <CheckCircle className="h-3 w-3 text-green-500" />
                            {feature}
                          </li>
                        ))}
                      </ul>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Billing Configuration</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="billingCycle">Billing Cycle</Label>
                    <Select value={formData.billingCycle} onValueChange={(value) => handleInputChange('billingCycle', value)}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="monthly">Monthly</SelectItem>
                        <SelectItem value="yearly">Yearly (10% discount)</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div>
                    <Label htmlFor="trialDays">Trial Period (days)</Label>
                    <Select value={formData.trialDays} onValueChange={(value) => handleInputChange('trialDays', value)}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="0">No trial</SelectItem>
                        <SelectItem value="7">7 days</SelectItem>
                        <SelectItem value="14">14 days</SelectItem>
                        <SelectItem value="30">30 days</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        );

      case 3:
        return (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Settings className="h-5 w-5" />
                Initial Configuration
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <Label htmlFor="timezone">Timezone *</Label>
                  <Select value={formData.timezone} onValueChange={(value) => handleInputChange('timezone', value)}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select timezone" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="UTC">UTC</SelectItem>
                      <SelectItem value="America/New_York">Eastern Time</SelectItem>
                      <SelectItem value="America/Chicago">Central Time</SelectItem>
                      <SelectItem value="America/Denver">Mountain Time</SelectItem>
                      <SelectItem value="America/Los_Angeles">Pacific Time</SelectItem>
                      <SelectItem value="Europe/London">London</SelectItem>
                      <SelectItem value="Europe/Paris">Paris</SelectItem>
                      <SelectItem value="Asia/Tokyo">Tokyo</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label htmlFor="currency">Currency *</Label>
                  <Select value={formData.currency} onValueChange={(value) => handleInputChange('currency', value)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="USD">USD - US Dollar</SelectItem>
                      <SelectItem value="EUR">EUR - Euro</SelectItem>
                      <SelectItem value="GBP">GBP - British Pound</SelectItem>
                      <SelectItem value="CAD">CAD - Canadian Dollar</SelectItem>
                      <SelectItem value="AUD">AUD - Australian Dollar</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label htmlFor="language">Language</Label>
                  <Select value={formData.language} onValueChange={(value) => handleInputChange('language', value)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="en">English</SelectItem>
                      <SelectItem value="es">Spanish</SelectItem>
                      <SelectItem value="fr">French</SelectItem>
                      <SelectItem value="de">German</SelectItem>
                      <SelectItem value="ja">Japanese</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              
              <div>
                <Label htmlFor="customDomain">Custom Domain (Optional)</Label>
                <Input
                  id="customDomain"
                  value={formData.customDomain}
                  onChange={(e) => handleInputChange('customDomain', e.target.value)}
                  placeholder="shop.yourcompany.com"
                />
              </div>
              
              <div className="space-y-3">
                <div className="flex items-center space-x-2">
                  <Switch
                    checked={formData.enableNotifications}
                    onCheckedChange={(checked) => handleInputChange('enableNotifications', checked)}
                  />
                  <Label>Enable email notifications</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Switch
                    checked={formData.enableAnalytics}
                    onCheckedChange={(checked) => handleInputChange('enableAnalytics', checked)}
                  />
                  <Label>Enable analytics tracking</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Switch
                    checked={formData.sendWelcomeEmail}
                    onCheckedChange={(checked) => handleInputChange('sendWelcomeEmail', checked)}
                  />
                  <Label>Send welcome email to admin</Label>
                </div>
              </div>
              
              <div>
                <Label htmlFor="notes">Additional Notes</Label>
                <Textarea
                  id="notes"
                  value={formData.notes}
                  onChange={(e) => handleInputChange('notes', e.target.value)}
                  placeholder="Any special requirements or notes for this tenant..."
                  rows={3}
                />
              </div>
            </CardContent>
          </Card>
        );

      case 4:
        const selectedPlan = subscriptionPlans.find(p => p.id === formData.selectedPlan);
        return (
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5" />
                  Review & Launch Tenant
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                <Alert>
                  <CheckCircle className="h-4 w-4" />
                  <AlertDescription>
                    Please review all information before creating the tenant account.
                  </AlertDescription>
                </Alert>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <h4 className="font-semibold mb-3">Company Information</h4>
                    <div className="space-y-2 text-sm">
                      <div><span className="font-medium">Name:</span> {formData.companyName}</div>
                      <div><span className="font-medium">Email:</span> {formData.companyEmail}</div>
                      <div><span className="font-medium">Industry:</span> {formData.industry}</div>
                      {formData.companySize && <div><span className="font-medium">Size:</span> {formData.companySize}</div>}
                    </div>
                  </div>
                  
                  <div>
                    <h4 className="font-semibold mb-3">Admin User</h4>
                    <div className="space-y-2 text-sm">
                      <div><span className="font-medium">Name:</span> {formData.adminFirstName} {formData.adminLastName}</div>
                      <div><span className="font-medium">Email:</span> {formData.adminEmail}</div>
                      {formData.adminPhone && <div><span className="font-medium">Phone:</span> {formData.adminPhone}</div>}
                    </div>
                  </div>
                  
                  <div>
                    <h4 className="font-semibold mb-3">Subscription</h4>
                    <div className="space-y-2 text-sm">
                      <div><span className="font-medium">Plan:</span> {selectedPlan?.name}</div>
                      <div><span className="font-medium">Price:</span> ${selectedPlan?.price}/{formData.billingCycle}</div>
                      <div><span className="font-medium">Trial:</span> {formData.trialDays} days</div>
                    </div>
                  </div>
                  
                  <div>
                    <h4 className="font-semibold mb-3">Configuration</h4>
                    <div className="space-y-2 text-sm">
                      <div><span className="font-medium">Timezone:</span> {formData.timezone}</div>
                      <div><span className="font-medium">Currency:</span> {formData.currency}</div>
                      <div><span className="font-medium">Language:</span> {formData.language}</div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="p-6">
      {/* Page Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Link to="/tenants" className="flex items-center gap-2 text-muted-foreground hover:text-foreground">
              <ArrowLeft className="h-4 w-4" />
              Back to Tenants
            </Link>
            <Separator orientation="vertical" className="h-6" />
            <Building2 className="h-5 w-5 text-primary" />
            <h1 className="text-xl font-semibold">Onboard New Tenant</h1>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline">
              <X className="h-4 w-4 mr-2" />
              Cancel
            </Button>
            {currentStep === onboardingSteps.length - 1 ? (
              <Button onClick={handleSubmit}>
                <CheckCircle className="h-4 w-4 mr-2" />
                Create Tenant
              </Button>
            ) : (
              <Button onClick={handleNext} disabled={!validateCurrentStep()}>
                Next Step
              </Button>
            )}
          </div>
        </div>
      </div>

      <div className="space-y-6">
          {/* Progress */}
          <Card>
            <CardContent className="pt-6">
              <div className="space-y-4">
                <div className="flex justify-between text-sm">
                  <span>Step {currentStep + 1} of {onboardingSteps.length}</span>
                  <span>{Math.round(getCurrentStepProgress())}% Complete</span>
                </div>
                <Progress value={getCurrentStepProgress()} className="h-2" />
                
                <div className="flex justify-between">
                  {onboardingSteps.map((step, index) => {
                    const Icon = step.icon;
                    const isActive = index === currentStep;
                    const isCompleted = index < currentStep;
                    
                    return (
                      <div key={step.id} className="flex flex-col items-center gap-2">
                        <div className={`w-10 h-10 rounded-full flex items-center justify-center ${
                          isCompleted ? 'bg-green-100 text-green-600' :
                          isActive ? 'bg-primary text-primary-foreground' :
                          'bg-muted text-muted-foreground'
                        }`}>
                          {isCompleted ? <CheckCircle className="h-5 w-5" /> : <Icon className="h-5 w-5" />}
                        </div>
                        <span className={`text-xs text-center ${
                          isActive ? 'font-medium' : 'text-muted-foreground'
                        }`}>
                          {step.title}
                        </span>
                      </div>
                    );
                  })}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Step Content */}
          {renderStepContent()}

          {/* Navigation */}
          <div className="flex justify-between">
            <Button 
              variant="outline" 
              onClick={handlePrevious} 
              disabled={currentStep === 0}
            >
              Previous
            </Button>
            <div className="flex gap-2">
              {currentStep === onboardingSteps.length - 1 ? (
                <Button onClick={handleSubmit}>
                  <Save className="h-4 w-4 mr-2" />
                  Create Tenant
                </Button>
              ) : (
                <Button onClick={handleNext} disabled={!validateCurrentStep()}>
                  Next Step
                </Button>
              )}
            </div>
          </div>
        </div>
    </div>
  );
}