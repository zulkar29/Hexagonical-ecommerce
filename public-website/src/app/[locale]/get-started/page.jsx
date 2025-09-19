'use client';

import { useState } from 'react';
import { ArrowRight, ArrowLeft, Check, Globe, ShoppingCart, CreditCard, Crown, Zap, Star, Clock, User, Palette, Mail, Phone, Lock, Eye, EyeOff } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { useTranslations } from '@/hooks/useTranslations';

export default function GetStartedPage() {
  const { t } = useTranslations();

  // Helper function to safely get features array
  const getFeatures = (key) => {
    try {
      const features = t(key);
      return Array.isArray(features) ? features : [];
    } catch {
      return [];
    }
  };
  const [currentStep, setCurrentStep] = useState(1);

  // User registration data
  const [userInfo, setUserInfo] = useState({
    name: '',
    email: '',
    phone: '',
    password: '',
    businessName: ''
  });
  const [showPassword, setShowPassword] = useState(false);

  // Store setup data
  const [selectedDomain, setSelectedDomain] = useState('');
  const [domainType, setDomainType] = useState('subdomain');
  const [selectedTemplate, setSelectedTemplate] = useState(null);
  const [selectedPlan, setSelectedPlan] = useState(null);

  const templateKeys = ['fashion', 'electronics', 'food', 'cosmetics', 'books', 'jewelry'];
  const paidPlanKeys = ['starter', 'professional', 'pro'];
  const freePlanKey = 'free';
  
  const planPrices = {
    free: 0,
    starter: 990,
    professional: 2990,
    pro: 4990
  };

  const planColors = {
    free: 'border-emerald-500',
    starter: 'border-gray-200',
    professional: 'border-orange-500',
    pro: 'border-purple-500'
  };

  const templateColors = {
    fashion: 'bg-rose-600',
    electronics: 'bg-indigo-600',
    food: 'bg-emerald-600',
    cosmetics: 'bg-pink-600',
    books: 'bg-amber-600',
    jewelry: 'bg-purple-600'
  };

  const steps = [
    { id: 1, titleKey: 'getStarted.steps.account', icon: User, descriptionKey: 'getStarted.stepDescriptions.account' },
    { id: 2, titleKey: 'getStarted.steps.store', icon: Globe, descriptionKey: 'getStarted.stepDescriptions.store' },
    { id: 3, titleKey: 'getStarted.steps.template', icon: Palette, descriptionKey: 'getStarted.stepDescriptions.template' },
    { id: 4, titleKey: 'getStarted.steps.plan', icon: CreditCard, descriptionKey: 'getStarted.stepDescriptions.plan' }
  ];

  const nextStep = () => {
    if (currentStep < 4) {
      setCurrentStep(currentStep + 1);
    }
  };

  const prevStep = () => {
    if (currentStep > 1) {
      setCurrentStep(currentStep - 1);
    }
  };

  const canProceed = () => {
    if (currentStep === 1) {
      return userInfo.name.length > 0 &&
             userInfo.email.length > 0 &&
             userInfo.phone.length > 0 &&
             userInfo.password.length >= 6 &&
             userInfo.businessName.length > 0;
    }
    if (currentStep === 2) return selectedDomain.length > 0;
    if (currentStep === 3) return selectedTemplate !== null;
    if (currentStep === 4) return selectedPlan !== null;
    return false;
  };

  const handleUserInfoChange = (field, value) => {
    setUserInfo(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const handleComplete = () => {
    // Here you would typically integrate with your backend
    const storeData = {
      user: userInfo,
      domain: {
        name: selectedDomain,
        type: domainType
      },
      template: selectedTemplate,
      plan: selectedPlan
    };

    const successMessage = t('getStarted.successMessage')
      .replace('{{businessName}}', userInfo.businessName)
      .replace('{{domain}}', `${selectedDomain}${domainType === 'subdomain' ? '.storebuilder.com' : '.com'}`)
      .replace('{{template}}', t(`getStarted.templates.${selectedTemplate}.name`))
      .replace('{{plan}}', t(`getStarted.plans.${selectedPlan}.name`))
      .replace('{{userName}}', userInfo.name);

    console.log('Store Data:', storeData);
    alert(successMessage);
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            {t('getStarted.title')}
          </h1>
          <p className="text-lg text-gray-600 mb-2">
            {t('getStarted.subtitle')}
          </p>
          <div className="flex items-center justify-center gap-2 text-sm text-gray-500">
            <Clock className="w-4 h-4" />
            <span>{t('getStarted.estimatedTime')}</span>
          </div>
          <div className="mt-4 flex items-center justify-center gap-4 text-sm text-gray-500">
            <div className="flex items-center gap-2">
              <Check className="w-4 h-4 text-emerald-500" />
              <span>Free 14-day trial</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="w-4 h-4 text-emerald-500" />
              <span>No credit card required</span>
            </div>
          </div>
        </div>

        {/* Progress Bar */}
        <div className="mb-12">
          <div className="flex items-start justify-between mb-4 relative">
            {steps.map((step, index) => {
              const IconComponent = step.icon;
              const isActive = currentStep === step.id;
              const isCompleted = currentStep > step.id;

              return (
                <div key={step.id} className="flex flex-col items-center flex-1 relative">
                  <div className={`flex items-center justify-center w-12 h-12 rounded-full border-2 transition-all duration-300 mb-3 relative z-10 ${
                    isCompleted
                      ? 'bg-gray-900 border-gray-900 text-white'
                      : isActive
                        ? 'border-orange-500 bg-orange-50 text-orange-600'
                        : 'border-gray-300 bg-white text-gray-400'
                  }`}>
                    {isCompleted ? (
                      <Check className="w-6 h-6" />
                    ) : (
                      <IconComponent className="w-6 h-6" />
                    )}
                  </div>

                  <div className="text-center">
                    <div className={`text-sm font-medium mb-1 ${isActive ? 'text-orange-600' : isCompleted ? 'text-gray-900' : 'text-gray-500'}`}>
                      {t(step.titleKey)}
                    </div>
                    <div className="text-xs text-gray-400 hidden sm:block">
                      {t(step.descriptionKey)}
                    </div>
                  </div>
                </div>
              );
            })}

            {/* Connecting Lines */}
            <div className="absolute top-6 left-0 right-0 flex justify-between px-6 -z-10">
              <div className={`h-0.5 flex-1 ${currentStep > 1 ? 'bg-gray-900' : 'bg-gray-300'}`} />
              <div className="w-12" />
              <div className={`h-0.5 flex-1 ${currentStep > 2 ? 'bg-gray-900' : 'bg-gray-300'}`} />
              <div className="w-12" />
              <div className={`h-0.5 flex-1 ${currentStep > 3 ? 'bg-gray-900' : 'bg-gray-300'}`} />
            </div>
          </div>
        </div>

        {/* Step Content */}
        <div className="bg-white rounded-2xl shadow-sm border p-6 md:p-8">
          <AnimatePresence mode="wait">
            {/* Step 1: Account Creation */}
            {currentStep === 1 && (
              <motion.div
                key="step1"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3 }}
              >
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 bg-orange-100 rounded-xl flex items-center justify-center">
                    <User className="w-6 h-6 text-orange-600" />
                  </div>
                  <h2 className="text-2xl md:text-3xl font-bold text-gray-900">
                    {t('getStarted.account.title')}
                  </h2>
                </div>
                <p className="text-gray-600 mb-8">
                  {t('getStarted.account.subtitle')}
                </p>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* Full Name */}
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                      {t('getStarted.account.fullName')} *
                    </label>
                    <div className="relative">
                      <User className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
                      <input
                        type="text"
                        value={userInfo.name}
                        onChange={(e) => handleUserInfoChange('name', e.target.value)}
                        placeholder={t('getStarted.account.fullNamePlaceholder')}
                        className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none transition-colors"
                      />
                    </div>
                  </div>

                  {/* Email */}
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                      {t('getStarted.account.email')} *
                    </label>
                    <div className="relative">
                      <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
                      <input
                        type="email"
                        value={userInfo.email}
                        onChange={(e) => handleUserInfoChange('email', e.target.value)}
                        placeholder={t('getStarted.account.emailPlaceholder')}
                        className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none transition-colors"
                      />
                    </div>
                  </div>

                  {/* Phone */}
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                      {t('getStarted.account.phone')} *
                    </label>
                    <div className="relative">
                      <Phone className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
                      <input
                        type="tel"
                        value={userInfo.phone}
                        onChange={(e) => handleUserInfoChange('phone', e.target.value)}
                        placeholder={t('getStarted.account.phonePlaceholder')}
                        className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none transition-colors"
                      />
                    </div>
                  </div>

                  {/* Password */}
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                      {t('getStarted.account.password')} *
                    </label>
                    <div className="relative">
                      <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
                      <input
                        type={showPassword ? "text" : "password"}
                        value={userInfo.password}
                        onChange={(e) => handleUserInfoChange('password', e.target.value)}
                        placeholder={t('getStarted.account.passwordPlaceholder')}
                        className="w-full pl-10 pr-12 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none transition-colors"
                      />
                      <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                      >
                        {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                      </button>
                    </div>
                    <p className="text-xs text-gray-500 mt-1">{t('getStarted.account.passwordHint')}</p>
                  </div>
                </div>

                {/* Business Name - Full Width */}
                <div className="mt-6">
                  <label className="block text-sm font-semibold text-gray-700 mb-2">
                    {t('getStarted.account.businessName')} *
                  </label>
                  <div className="relative">
                    <ShoppingCart className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
                    <input
                      type="text"
                      value={userInfo.businessName}
                      onChange={(e) => handleUserInfoChange('businessName', e.target.value)}
                      placeholder={t('getStarted.account.businessNamePlaceholder')}
                      className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none transition-colors"
                    />
                  </div>
                  <p className="text-xs text-gray-500 mt-1">{t('getStarted.account.businessNameHint')}</p>
                </div>

                {/* Terms Agreement */}
                <div className="mt-6 p-4 bg-gray-50 rounded-xl">
                  <div className="flex items-start gap-3">
                    <Check className="w-5 h-5 text-emerald-500 mt-0.5" />
                    <div className="text-sm text-gray-600">
                      <p>{t('getStarted.account.agreement')}</p>
                    </div>
                  </div>
                </div>
              </motion.div>
            )}

            {/* Step 2: Store Setup (Previously Step 1) */}
            {currentStep === 2 && (
              <motion.div
                key="step1"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3 }}
              >
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 bg-orange-100 rounded-xl flex items-center justify-center">
                    <Globe className="w-6 h-6 text-orange-600" />
                  </div>
                  <h2 className="text-2xl md:text-3xl font-bold text-gray-900">
                    {t('getStarted.domain.title')}
                  </h2>
                </div>
                <p className="text-gray-600 mb-8">
                  {t('getStarted.domain.subtitle')}
                </p>

                <div className="space-y-6">
                  {/* Domain Type Selection */}
                  <div>
                    <label className="text-lg font-semibold text-gray-900 mb-4 block">
                      {t('getStarted.domain.typeLabel')}
                    </label>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div
                        className={`p-6 border-2 rounded-xl cursor-pointer transition-all duration-200 ${
                          domainType === 'subdomain'
                            ? 'border-orange-500 bg-orange-50'
                            : 'border-gray-200 hover:border-gray-300'
                        }`}
                        onClick={() => setDomainType('subdomain')}
                      >
                        <div className="flex items-center justify-between mb-3">
                          <h3 className="font-semibold text-gray-900">{t('getStarted.domain.freeSubdomain')}</h3>
                          <div className="bg-emerald-100 text-emerald-700 px-3 py-1 rounded-full text-xs font-medium">
                            {t('getStarted.domain.free')}
                          </div>
                        </div>
                        <p className="text-gray-600 text-sm mb-2">yourstore.storebuilder.com</p>
                        <p className="text-gray-500 text-xs">{t('getStarted.domain.freeDescription')}</p>
                      </div>

                      <div
                        className={`p-6 border-2 rounded-xl cursor-pointer transition-all duration-200 ${
                          domainType === 'custom'
                            ? 'border-orange-500 bg-orange-50'
                            : 'border-gray-200 hover:border-gray-300'
                        }`}
                        onClick={() => setDomainType('custom')}
                      >
                        <div className="flex items-center justify-between mb-3">
                          <h3 className="font-semibold text-gray-900">{t('getStarted.domain.customDomain')}</h3>
                          <div className="bg-purple-100 text-purple-700 px-3 py-1 rounded-full text-xs font-medium">
                            {t('getStarted.domain.pro')}
                          </div>
                        </div>
                        <p className="text-gray-600 text-sm mb-2">yourstore.com</p>
                        <p className="text-gray-500 text-xs">{t('getStarted.domain.customDescription')}</p>
                      </div>
                    </div>
                  </div>

                  {/* Domain Input */}
                  <div>
                    <label className="text-lg font-semibold text-gray-900 mb-4 block">
                      {t('getStarted.domain.storeNameLabel')}
                    </label>
                    <div className="flex items-center">
                      <input
                        type="text"
                        placeholder={t('getStarted.domain.storeNamePlaceholder')}
                        value={selectedDomain}
                        onChange={(e) => setSelectedDomain(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ''))}
                        className="flex-1 px-4 py-3 border border-gray-300 rounded-l-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none text-lg"
                      />
                      <div className="px-4 py-3 bg-gray-100 border border-l-0 border-gray-300 rounded-r-lg text-gray-600 font-mono">
                        {domainType === 'subdomain' ? '.storebuilder.com' : '.com'}
                      </div>
                    </div>
                    <p className="text-sm text-gray-500 mt-2">
                      {t('getStarted.domain.instructions')}
                    </p>
                  </div>

                  {/* Preview */}
                  {selectedDomain && (
                    <div className="bg-gray-50 rounded-xl p-6 border">
                      <h4 className="font-semibold text-gray-900 mb-3">{t('getStarted.domain.preview')}</h4>
                      <div className="flex items-center gap-3 p-4 bg-white rounded-lg border">
                        <Globe className="w-5 h-5 text-orange-600" />
                        <span className="font-mono text-lg text-gray-900">
                          {selectedDomain}{domainType === 'subdomain' ? '.storebuilder.com' : '.com'}
                        </span>
                      </div>
                    </div>
                  )}
                </div>
              </motion.div>
            )}

            {/* Step 3: Template Selection */}
            {currentStep === 3 && (
              <motion.div
                key="step2"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3 }}
              >
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 bg-orange-100 rounded-xl flex items-center justify-center">
                    <Palette className="w-6 h-6 text-orange-600" />
                  </div>
                  <h2 className="text-2xl md:text-3xl font-bold text-gray-900">
                    Choose Your Store Design
                  </h2>
                </div>
                <p className="text-gray-600 mb-8">
                  Select a beautiful template that matches your business style. You can customize it later.
                </p>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {templateKeys.map((templateKey) => (
                    <div
                      key={templateKey}
                      className={`cursor-pointer transition-all duration-200 ${
                        selectedTemplate === templateKey
                          ? 'transform scale-[1.02]'
                          : 'hover:transform hover:scale-[1.01]'
                      }`}
                      onClick={() => setSelectedTemplate(templateKey)}
                    >
                      <div className={`bg-white rounded-xl border-2 overflow-hidden transition-all duration-200 ${
                        selectedTemplate === templateKey
                          ? 'border-orange-500 ring-2 ring-orange-100'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}>
                        {/* Template Preview */}
                        <div className={`h-32 ${templateColors[templateKey]} flex items-center justify-center relative`}>
                          <div className="text-white text-center">
                            <div className="w-12 h-12 bg-white/20 rounded-xl flex items-center justify-center mx-auto mb-2">
                              <ShoppingCart className="w-6 h-6" />
                            </div>
                            <div className="font-semibold text-sm">{t(`getStarted.templates.${templateKey}.name`)}</div>
                          </div>
                          {selectedTemplate === templateKey && (
                            <div className="absolute top-2 right-2 w-6 h-6 bg-orange-500 rounded-full flex items-center justify-center">
                              <Check className="w-4 h-4 text-white" />
                            </div>
                          )}
                        </div>

                        {/* Template Details */}
                        <div className="p-4">
                          <div className="flex items-center justify-between mb-2">
                            <h3 className="font-bold text-gray-900">{t(`getStarted.templates.${templateKey}.name`)}</h3>
                            <span className="text-xs text-orange-600 bg-orange-50 px-2 py-1 rounded-full font-medium">
                              {t(`getStarted.templates.${templateKey}.category`)}
                            </span>
                          </div>
                          <p className="text-gray-600 text-sm mb-3">{t(`getStarted.templates.${templateKey}.description`)}</p>
                          <div className="space-y-1">
                            <div className="flex items-center gap-2">
                              <Check className="w-3 h-3 text-emerald-500" />
                              <span className="text-xs text-gray-600">Mobile Responsive</span>
                            </div>
                            <div className="flex items-center gap-2">
                              <Check className="w-3 h-3 text-emerald-500" />
                              <span className="text-xs text-gray-600">SEO Optimized</span>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </motion.div>
            )}


            {/* Step 4: Plan Selection */}
            {currentStep === 4 && (
              <motion.div
                key="step3"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3 }}
              >
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 bg-orange-100 rounded-xl flex items-center justify-center">
                    <CreditCard className="w-6 h-6 text-orange-600" />
                  </div>
                  <h2 className="text-2xl md:text-3xl font-bold text-gray-900">
                    {t('getStarted.plan.title')}
                  </h2>
                </div>
                <p className="text-gray-600 mb-8">
                  {t('getStarted.plan.subtitle')}
                </p>

                {/* Paid Plans Grid */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                  {paidPlanKeys.map((planKey) => (
                    <div
                      key={planKey}
                      className={`cursor-pointer transition-all duration-300 ${
                        selectedPlan === planKey
                          ? 'transform scale-105'
                          : 'hover:transform hover:scale-102'
                      }`}
                      onClick={() => setSelectedPlan(planKey)}
                    >
                      <div className={`bg-white rounded-xl shadow-sm hover:shadow-lg transition-all duration-300 p-6 border-2 relative ${
                        selectedPlan === planKey
                          ? 'border-orange-500 ring-2 ring-orange-100'
                          : planColors[planKey]
                      }`}>
                        {planKey === 'professional' && (
                          <div className="absolute -top-3 left-1/2 transform -translate-x-1/2">
                            <div className="bg-orange-600 text-white px-4 py-1 rounded-full text-xs font-bold flex items-center gap-1">
                              <Crown className="w-3 h-3" />
                              {t('getStarted.plan.recommended')}
                            </div>
                          </div>
                        )}

                        {selectedPlan === planKey && (
                          <div className="absolute top-4 right-4 w-6 h-6 bg-orange-500 rounded-full flex items-center justify-center">
                            <Check className="w-4 h-4 text-white" />
                          </div>
                        )}

                        <div className="mb-6">
                          <h3 className="text-xl font-bold text-gray-900 mb-2">{t(`getStarted.plans.${planKey}.name`)}</h3>
                          <p className="text-gray-600 text-sm mb-4">{t(`getStarted.plans.${planKey}.description`)}</p>
                          <div className="flex items-baseline gap-1 mb-2">
                            <span className="text-3xl font-bold text-gray-900">৳{planPrices[planKey].toLocaleString()}</span>
                            <span className="text-gray-600">{t('getStarted.plan.monthly')}</span>
                          </div>
                          <p className="text-sm text-emerald-600 font-medium">✓ {t('getStarted.plan.freeTrial')}</p>
                        </div>

                        <div className="space-y-2 mb-4">
                          {getFeatures(`getStarted.plans.${planKey}.features`).slice(0, 4).map((feature, index) => (
                            <div key={index} className="flex items-start gap-2">
                              <Check className="w-4 h-4 text-emerald-500 mt-0.5 flex-shrink-0" />
                              <span className="text-sm text-gray-700">{feature}</span>
                            </div>
                          ))}
                          {getFeatures(`getStarted.plans.${planKey}.features`).length > 4 && (
                            <div className="text-sm text-gray-500">
                              +{getFeatures(`getStarted.plans.${planKey}.features`).length - 4} {t('getStarted.plan.moreFeatures')}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Free Trial Option - Full Width */}
                <div className="border-t pt-6">
                  <div className="text-center mb-4">
                    <h4 className="text-lg font-semibold text-gray-900 mb-2">Or start with our free trial</h4>
                    <p className="text-gray-600">No credit card required • Cancel anytime</p>
                  </div>

                  <div
                    className={`cursor-pointer transition-all duration-300 ${
                      selectedPlan === freePlanKey
                        ? 'transform scale-[1.02]'
                        : 'hover:transform hover:scale-[1.01]'
                    }`}
                    onClick={() => setSelectedPlan(freePlanKey)}
                  >
                    <div className={`bg-white rounded-xl shadow-sm hover:shadow-lg transition-all duration-300 p-6 border-2 relative ${
                      selectedPlan === freePlanKey
                        ? 'border-emerald-500 ring-2 ring-emerald-100 bg-emerald-50'
                        : 'border-emerald-200 hover:border-emerald-300'
                    }`}>
                      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                        <div className="flex items-center gap-4 flex-1">
                          <div className="flex items-center gap-3">
                            <div className="w-12 h-12 bg-emerald-100 rounded-xl flex items-center justify-center flex-shrink-0">
                              <Star className="w-6 h-6 text-emerald-600" />
                            </div>
                            <div>
                              <h3 className="text-xl font-bold text-gray-900">{t(`getStarted.plans.${freePlanKey}.name`)}</h3>
                              <p className="text-gray-600 text-sm">{t(`getStarted.plans.${freePlanKey}.description`)}</p>
                            </div>
                          </div>
                          <div className="text-right md:ml-auto">
                            <div className="text-2xl font-bold text-emerald-600">FREE</div>
                            <div className="text-sm text-emerald-600 font-medium">14 days trial</div>
                          </div>
                        </div>

                        <div className="flex items-center justify-between md:justify-end gap-6">
                          <div className="flex md:hidden items-center gap-4 text-sm text-gray-600">
                            {getFeatures(`getStarted.plans.${freePlanKey}.features`).slice(0, 2).map((feature, index) => (
                              <div key={index} className="flex items-center gap-1">
                                <Check className="w-3 h-3 text-emerald-500" />
                                <span>{feature}</span>
                              </div>
                            ))}
                          </div>
                          <div className="hidden md:flex items-center gap-4 text-sm text-gray-600">
                            {getFeatures(`getStarted.plans.${freePlanKey}.features`).slice(0, 3).map((feature, index) => (
                              <div key={index} className="flex items-center gap-1">
                                <Check className="w-3 h-3 text-emerald-500" />
                                <span>{feature}</span>
                              </div>
                            ))}
                          </div>

                          {selectedPlan === freePlanKey && (
                            <div className="w-8 h-8 bg-emerald-500 rounded-full flex items-center justify-center flex-shrink-0">
                              <Check className="w-5 h-5 text-white" />
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="mt-6 p-4 bg-amber-50 rounded-lg border border-amber-200">
                  <div className="flex items-start gap-3">
                    <Star className="w-5 h-5 text-amber-600 mt-0.5" />
                    <div>
                      <h4 className="font-semibold text-amber-800 mb-1">{t('getStarted.plan.specialOffer.title')}</h4>
                      <p className="text-sm text-amber-700">{t('getStarted.plan.specialOffer.description')}</p>
                    </div>
                  </div>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        {/* Navigation Buttons */}
        <div className="flex items-center justify-between mt-8">
          <button
            onClick={prevStep}
            disabled={currentStep === 1}
            className="flex items-center gap-2 px-6 py-3 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
          >
            <ArrowLeft className="w-5 h-5" />
            {t('getStarted.navigation.previous')}
          </button>

          {currentStep < 4 ? (
            <button
              onClick={nextStep}
              disabled={!canProceed()}
              className="flex items-center gap-2 px-8 py-3 bg-gray-900 text-white rounded-lg hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              {t('getStarted.navigation.next')}
              <ArrowRight className="w-5 h-5" />
            </button>
          ) : (
            <button
              onClick={handleComplete}
              disabled={!canProceed()}
              className="flex items-center gap-2 px-8 py-3 bg-orange-600 text-white rounded-lg hover:bg-orange-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              <Zap className="w-5 h-5" />
              {t('getStarted.navigation.createStore')}
            </button>
          )}
        </div>

        {/* Summary Sidebar */}
        {currentStep > 1 && (
          <div className="mt-8 bg-gray-50 rounded-xl p-6 border">
            <h4 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
              <User className="w-5 h-5" />
              {t('getStarted.summary.title')}
            </h4>
            <div className="space-y-3 text-sm">
              {userInfo.name && (
                <div className="flex items-center gap-2">
                  <User className="w-4 h-4 text-blue-600" />
                  <span><strong>{t('getStarted.summary.account')}</strong> {userInfo.name} ({userInfo.businessName})</span>
                </div>
              )}
              {selectedDomain && (
                <div className="flex items-center gap-2">
                  <Globe className="w-4 h-4 text-orange-600" />
                  <span><strong>{t('getStarted.summary.domain')}</strong> {selectedDomain}{domainType === 'subdomain' ? '.storebuilder.com' : '.com'}</span>
                </div>
              )}
              {selectedTemplate && (
                <div className="flex items-center gap-2">
                  <Palette className="w-4 h-4 text-emerald-600" />
                  <span><strong>{t('getStarted.summary.template')}</strong> {t(`getStarted.templates.${selectedTemplate}.name`)}</span>
                </div>
              )}
              {selectedPlan && (
                <div className="flex items-center gap-2">
                  <CreditCard className="w-4 h-4 text-purple-600" />
                  <span><strong>{t('getStarted.summary.plan')}</strong> {t(`getStarted.plans.${selectedPlan}.name`)} - {selectedPlan === 'free' ? 'FREE' : `৳${planPrices[selectedPlan].toLocaleString()}${t('getStarted.summary.monthly')}`}</span>
                </div>
              )}
            </div>
            {currentStep === 4 && selectedPlan && (
              <div className="mt-4 pt-4 border-t border-gray-200">
                {selectedPlan === 'free' ? (
                  <div className="text-lg font-bold text-emerald-600">
                    Free for 14 days - No credit card required
                  </div>
                ) : (
                  <>
                    <div className="text-lg font-bold text-gray-900">
                      {t('getStarted.summary.firstMonthPrice')} <span className="text-orange-600">৳০ ({t('getStarted.summary.freeTrial')})</span>
                    </div>
                    <div className="text-sm text-gray-500 mt-1">
                      {t('getStarted.summary.thenMonthly')} ৳{planPrices[selectedPlan].toLocaleString()} {t('getStarted.summary.perMonth')}
                    </div>
                  </>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}