'use client';

import { useState } from 'react';
import Link from 'next/link';
import { ArrowRight, ArrowLeft, Check, Globe, ShoppingCart, CreditCard, Crown, Zap, Star, Clock, User, Store, Palette } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { useTranslations } from '@/hooks/useTranslations';

export default function GetStartedPage() {
  const { t } = useTranslations();
  const [currentStep, setCurrentStep] = useState(1);
  const [selectedDomain, setSelectedDomain] = useState('');
  const [domainType, setDomainType] = useState('subdomain');
  const [selectedTemplate, setSelectedTemplate] = useState(null);
  const [selectedPlan, setSelectedPlan] = useState(null);
  const [businessInfo, setBusinessInfo] = useState({
    name: '',
    description: '',
    category: ''
  });

  const templateKeys = ['fashion', 'electronics', 'food', 'cosmetics', 'books', 'jewelry'];
  const planKeys = ['starter', 'professional', 'pro', 'enterprise'];
  
  const planPrices = {
    starter: 990,
    professional: 2990,
    pro: 4990,
    enterprise: 8990
  };

  const planColors = {
    starter: 'border-gray-200',
    professional: 'border-orange-500',
    pro: 'border-gray-200',
    enterprise: 'border-purple-500'
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
    { id: 1, titleKey: 'getStarted.steps.businessInfo', icon: Store, descriptionKey: 'getStarted.stepDescriptions.businessInfo' },
    { id: 2, titleKey: 'getStarted.steps.domain', icon: Globe, descriptionKey: 'getStarted.stepDescriptions.domain' },
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
    if (currentStep === 1) return businessInfo.name.length > 0 && businessInfo.category.length > 0;
    if (currentStep === 2) return selectedDomain.length > 0;
    if (currentStep === 3) return selectedTemplate !== null;
    if (currentStep === 4) return selectedPlan !== null;
    return false;
  };

  const handleComplete = () => {
    const storeData = {
      business: businessInfo,
      domain: {
        name: selectedDomain,
        type: domainType
      },
      template: selectedTemplate,
      plan: selectedPlan
    };
    
    // Here you would typically integrate with your backend
    const successMessage = t('getStarted.successMessage')
      .replace('{{businessName}}', businessInfo.name)
      .replace('{{domain}}', `${selectedDomain}${domainType === 'subdomain' ? '.ourplatform.com' : '.com'}`)
      .replace('{{template}}', t(`getStarted.templates.${selectedTemplate}.name`))
      .replace('{{plan}}', t(`getStarted.plans.${selectedPlan}.name`));
    
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
        </div>

        {/* Progress Bar */}
        <div className="mb-12">
          <div className="flex items-center justify-between mb-4">
            {steps.map((step, index) => {
              const IconComponent = step.icon;
              const isActive = currentStep === step.id;
              const isCompleted = currentStep > step.id;
              
              return (
                <div key={step.id} className="flex flex-col items-center flex-1">
                  <div className="flex items-center w-full">
                    <div className={`flex items-center justify-center w-10 h-10 rounded-full border-2 transition-all duration-300 ${
                      isCompleted 
                        ? 'bg-gray-900 border-gray-900 text-white' 
                        : isActive 
                          ? 'border-orange-500 bg-orange-50 text-orange-600' 
                          : 'border-gray-300 text-gray-400'
                    }`}>
                      {isCompleted ? (
                        <Check className="w-5 h-5" />
                      ) : (
                        <IconComponent className="w-5 h-5" />
                      )}
                    </div>
                    {index < steps.length - 1 && (
                      <div className={`flex-1 h-0.5 mx-4 ${currentStep > step.id ? 'bg-gray-900' : 'bg-gray-300'}`} />
                    )}
                  </div>
                  <div className="mt-2 text-center">
                    <div className={`text-sm font-medium ${isActive ? 'text-orange-600' : isCompleted ? 'text-gray-900' : 'text-gray-500'}`}>
                      {t(step.titleKey)}
                    </div>
                    <div className="text-xs text-gray-400 mt-1 hidden sm:block">
                      {t(step.descriptionKey)}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Step Content */}
        <div className="bg-white rounded-2xl shadow-sm border p-6 md:p-8">
          <AnimatePresence mode="wait">
            {/* Step 1: Business Information */}
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
                    <Store className="w-6 h-6 text-orange-600" />
                  </div>
                  <h2 className="text-2xl md:text-3xl font-bold text-gray-900">
                    {t('getStarted.businessInfo.title')}
                  </h2>
                </div>
                <p className="text-gray-600 mb-8">
                  {t('getStarted.businessInfo.subtitle')}
                </p>

                <div className="space-y-6">
                  {/* Business Name */}
                  <div>
                    <label className="block text-lg font-semibold text-gray-900 mb-3">
                      {t('getStarted.businessInfo.nameLabel')} <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      placeholder={t('getStarted.businessInfo.namePlaceholder')}
                      value={businessInfo.name}
                      onChange={(e) => setBusinessInfo({...businessInfo, name: e.target.value})}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none text-lg"
                    />
                  </div>

                  {/* Business Category */}
                  <div>
                    <label className="block text-lg font-semibold text-gray-900 mb-3">
                      {t('getStarted.businessInfo.categoryLabel')} <span className="text-red-500">*</span>
                    </label>
                    <select
                      value={businessInfo.category}
                      onChange={(e) => setBusinessInfo({...businessInfo, category: e.target.value})}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none text-lg"
                    >
                      <option value="">{t('getStarted.businessInfo.categoryPlaceholder')}</option>
                      <option value="fashion">{t('getStarted.businessInfo.categories.fashion')}</option>
                      <option value="electronics">{t('getStarted.businessInfo.categories.electronics')}</option>
                      <option value="food">{t('getStarted.businessInfo.categories.food')}</option>
                      <option value="cosmetics">{t('getStarted.businessInfo.categories.cosmetics')}</option>
                      <option value="books">{t('getStarted.businessInfo.categories.books')}</option>
                      <option value="jewelry">{t('getStarted.businessInfo.categories.jewelry')}</option>
                      <option value="health">{t('getStarted.businessInfo.categories.health')}</option>
                      <option value="home">{t('getStarted.businessInfo.categories.home')}</option>
                      <option value="other">{t('getStarted.businessInfo.categories.other')}</option>
                    </select>
                  </div>

                  {/* Business Description */}
                  <div>
                    <label className="block text-lg font-semibold text-gray-900 mb-3">
                      {t('getStarted.businessInfo.descriptionLabel')}
                    </label>
                    <textarea
                      placeholder={t('getStarted.businessInfo.descriptionPlaceholder')}
                      value={businessInfo.description}
                      onChange={(e) => setBusinessInfo({...businessInfo, description: e.target.value})}
                      rows={3}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none resize-none"
                    />
                  </div>
                </div>
              </motion.div>
            )}

            {/* Step 2: Domain Selection */}
            {currentStep === 2 && (
              <motion.div
                key="step2"
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
                key="step3"
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
                    {t('getStarted.template.title')}
                  </h2>
                </div>
                <p className="text-gray-600 mb-8">
                  {t('getStarted.template.subtitle')}
                </p>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {templateKeys.map((templateKey) => (
                    <div
                      key={templateKey}
                      className={`cursor-pointer transition-all duration-300 ${
                        selectedTemplate === templateKey
                          ? 'transform scale-105'
                          : 'hover:transform hover:scale-102'
                      }`}
                      onClick={() => setSelectedTemplate(templateKey)}
                    >
                      <div className={`bg-white rounded-xl shadow-sm hover:shadow-lg transition-all duration-300 overflow-hidden border-2 ${
                        selectedTemplate === templateKey 
                          ? 'border-orange-500 ring-2 ring-orange-100' 
                          : 'border-gray-100'
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
                            <h3 className="font-bold text-gray-900 text-sm">{t(`getStarted.templates.${templateKey}.name`)}</h3>
                            <span className="text-xs text-orange-600 bg-orange-50 px-2 py-1 rounded-full font-medium">
                              {t(`getStarted.templates.${templateKey}.category`)}
                            </span>
                          </div>
                          <p className="text-gray-600 text-xs mb-3 line-clamp-2">{t(`getStarted.templates.${templateKey}.description`)}</p>
                          <div className="space-y-1">
                            {t(`getStarted.templates.${templateKey}.features`).slice(0, 2).map((feature, index) => (
                              <div key={index} className="flex items-center gap-1">
                                <Check className="w-3 h-3 text-emerald-500" />
                                <span className="text-xs text-gray-600">{feature}</span>
                              </div>
                            ))}
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
                key="step4"
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

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {planKeys.slice(0, 2).map((planKey) => (
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
                          {t(`getStarted.plans.${planKey}.features`).slice(0, 4).map((feature, index) => (
                            <div key={index} className="flex items-start gap-2">
                              <Check className="w-4 h-4 text-emerald-500 mt-0.5 flex-shrink-0" />
                              <span className="text-sm text-gray-700">{feature}</span>
                            </div>
                          ))}
                          {t(`getStarted.plans.${planKey}.features`).length > 4 && (
                            <div className="text-sm text-gray-500">
                              +{t(`getStarted.plans.${planKey}.features`).length - 4} {t('getStarted.plan.moreFeatures')}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
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
              {businessInfo.name && (
                <div className="flex items-start gap-2">
                  <Store className="w-4 h-4 text-orange-600 mt-0.5" />
                  <div>
                    <span className="font-medium">{t('getStarted.summary.business')}</span> {businessInfo.name}
                    {businessInfo.category && <div className="text-gray-500">{t('getStarted.summary.businessType')} {businessInfo.category}</div>}
                  </div>
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
                  <span><strong>{t('getStarted.summary.plan')}</strong> {t(`getStarted.plans.${selectedPlan}.name`)} - ৳{planPrices[selectedPlan].toLocaleString()}{t('getStarted.summary.monthly')}</span>
                </div>
              )}
            </div>
            {currentStep === 4 && (
              <div className="mt-4 pt-4 border-t border-gray-200">
                <div className="text-lg font-bold text-gray-900">
                  {t('getStarted.summary.firstMonthPrice')} <span className="text-orange-600">৳০ ({t('getStarted.summary.freeTrial')})</span>
                </div>
                <div className="text-sm text-gray-500 mt-1">
                  {t('getStarted.summary.thenMonthly')} ৳{selectedPlan ? planPrices[selectedPlan].toLocaleString() : '0'} {t('getStarted.summary.perMonth')}
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}