'use client';

import Link from 'next/link';
import { ArrowRight, Check, ShoppingCart, Play, Star, Filter } from 'lucide-react';
import { motion } from 'framer-motion';
import { useTranslations } from '@/hooks/useTranslations';
import { useState } from 'react';

export default function TemplatesPage() {
  const { t } = useTranslations();
  const [selectedCategory, setSelectedCategory] = useState('all');

  const templateKeys = ['fashion', 'electronics', 'food', 'cosmetics', 'books', 'jewelry'];
  
  const templateColors = {
    fashion: 'bg-rose-600',
    electronics: 'bg-indigo-600',
    food: 'bg-emerald-600',
    cosmetics: 'bg-pink-600',
    books: 'bg-amber-600',
    jewelry: 'bg-purple-600'
  };

  const categories = [
    { key: 'all', name: t('templatesPage.categories.all') },
    { key: 'fashion', name: t('templatesPage.categories.fashion') },
    { key: 'electronics', name: t('templatesPage.categories.electronics') },
    { key: 'food', name: t('templatesPage.categories.food') },
    { key: 'health', name: t('templatesPage.categories.health') },
    { key: 'books', name: t('templatesPage.categories.books') }
  ];

  const additionalTemplates = [
    {
      key: 'health',
      name: t('templatesPage.additionalTemplates.health.name'),
      category: t('templatesPage.additionalTemplates.health.category'),
      color: 'bg-teal-600',
      description: t('templatesPage.additionalTemplates.health.description'),
      features: Array.isArray(t('templatesPage.additionalTemplates.health.features')) ? 
        t('templatesPage.additionalTemplates.health.features') : 
        [t('templatesPage.additionalTemplates.health.features')]
    },
    {
      key: 'home',
      name: t('templatesPage.additionalTemplates.home.name'),
      category: t('templatesPage.additionalTemplates.home.category'),
      color: 'bg-lime-600',
      description: t('templatesPage.additionalTemplates.home.description'),
      features: Array.isArray(t('templatesPage.additionalTemplates.home.features')) ? 
        t('templatesPage.additionalTemplates.home.features') : 
        [t('templatesPage.additionalTemplates.home.features')]
    },
    {
      key: 'toys',
      name: t('templatesPage.additionalTemplates.toys.name'),
      category: t('templatesPage.additionalTemplates.toys.category'),
      color: 'bg-orange-600',
      description: t('templatesPage.additionalTemplates.toys.description'),
      features: Array.isArray(t('templatesPage.additionalTemplates.toys.features')) ? 
        t('templatesPage.additionalTemplates.toys.features') : 
        [t('templatesPage.additionalTemplates.toys.features')]
    }
  ];

  const filteredTemplates = selectedCategory === 'all' 
    ? templateKeys 
    : templateKeys.filter(key => key === selectedCategory);

  const filteredAdditionalTemplates = selectedCategory === 'all' 
    ? additionalTemplates 
    : additionalTemplates.filter(template => template.key === selectedCategory);

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center mb-16">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <h1 className="text-5xl font-bold text-gray-900 mb-6">
              {t('themes.title')}
            </h1>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto mb-8">
              {t('templatesPage.description')} {t('templatesPage.descriptionExtended')}
            </p>
            
            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12">
              <Link
                href="/get-started"
                className="bg-gray-900 text-white px-8 py-4 rounded-lg font-semibold hover:bg-gray-800 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105 inline-flex items-center justify-center"
              >
                {t('common.getStarted')}
                <ArrowRight className="w-5 h-5 ml-2" />
              </Link>
              <Link
                href="/pricing"
                className="border-2 border-gray-300 text-gray-700 px-8 py-4 rounded-lg font-semibold hover:border-gray-400 hover:bg-gray-50 transition-all duration-300 inline-flex items-center justify-center"
              >
                {t('pricing.title')}
              </Link>
            </div>
          </motion.div>
        </div>

        {/* Filter Section */}
        <div className="mb-12">
          <div className="text-center mb-8">
            <div className="flex items-center justify-center gap-2 text-gray-600 mb-6">
              <Filter className="w-5 h-5" />
              <span className="font-medium text-lg">{t('templatesPage.categoryFilter')}</span>
            </div>
            <div className="flex flex-wrap gap-3 justify-center">
              {categories.map((category) => (
                <button
                  key={category.key}
                  onClick={() => setSelectedCategory(category.key)}
                  className={`px-6 py-3 rounded-lg font-medium transition-all duration-200 ${
                    selectedCategory === category.key
                      ? 'bg-gray-900 text-white shadow-lg'
                      : 'bg-white text-gray-700 hover:bg-gray-100 hover:text-gray-900 border border-gray-200 shadow-sm'
                  }`}
                >
                  {category.name}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Popular Templates Section */}
        <div className="mb-16">
          <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">{t('templatesPage.popularTemplates')}</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {filteredTemplates.map((templateKey, index) => (
              <motion.div
                key={templateKey}
                className="group cursor-pointer"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                whileHover={{ y: -5 }}
              >
                <div className="bg-white rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100">
                  {/* Template Preview */}
                  <div className={`h-48 ${templateColors[templateKey]} flex items-center justify-center relative`}>
                    <div className="text-white text-center">
                      <div className="w-16 h-16 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-4">
                        <ShoppingCart className="w-8 h-8" />
                      </div>
                      <div className="text-lg font-semibold">{t(`getStarted.templates.${templateKey}.name`)}</div>
                    </div>
                    <div className="absolute top-4 right-4 bg-white/20 backdrop-blur-sm rounded-full p-2">
                      <Play className="w-4 h-4 text-white" />
                    </div>
                    {(templateKey === 'fashion' || templateKey === 'electronics') && (
                      <div className="absolute top-4 left-4 bg-orange-600 text-white px-3 py-1 rounded-full text-sm font-semibold flex items-center gap-1">
                        <Star className="w-3 h-3 fill-current" />
                        {t('templatesPage.popular')}
                      </div>
                    )}
                  </div>

                  {/* Template Details */}
                  <div className="p-6">
                    <div className="flex items-center justify-between mb-4">
                      <h3 className="text-xl font-bold text-gray-900">{t(`getStarted.templates.${templateKey}.name`)}</h3>
                      <span className="text-xs text-orange-600 bg-orange-50 px-2 py-1 rounded-full font-medium">
                        {t(`getStarted.templates.${templateKey}.category`)}
                      </span>
                    </div>
                    <p className="text-gray-600 text-sm mb-4">{t(`getStarted.templates.${templateKey}.description`)}</p>
                    <div className="space-y-2 mb-6">
                      {Array.isArray(t(`getStarted.templates.${templateKey}.features`)) ? 
                        t(`getStarted.templates.${templateKey}.features`).map((feature, featureIndex) => (
                          <div key={featureIndex} className="flex items-center gap-2">
                            <Check className="w-4 h-4 text-emerald-500" />
                            <span className="text-sm text-gray-600">{feature}</span>
                          </div>
                        )) : (
                          <div className="flex items-center gap-2">
                            <Check className="w-4 h-4 text-emerald-500" />
                            <span className="text-sm text-gray-600">{t(`getStarted.templates.${templateKey}.features`)}</span>
                          </div>
                        )
                      }
                    </div>
                    

                    
                    <div className="flex gap-3">
                      <Link
                        href="/get-started"
                        className="flex-1 bg-gray-900 text-white text-center py-3 rounded-lg font-medium hover:bg-gray-800 transition-colors duration-200"
                      >
                        {t('templatesPage.useTemplate')}
                      </Link>
                      <button className="px-4 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">
                        <Play className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        {/* Additional Templates Section */}
        {filteredAdditionalTemplates.length > 0 && (
          <div className="mb-16">
            <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">{t('templatesPage.moreTemplates')}</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {filteredAdditionalTemplates.map((template, index) => (
                <motion.div
                  key={template.key}
                  className="group cursor-pointer"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: (index + filteredTemplates.length) * 0.1 }}
                  whileHover={{ y: -5 }}
                >
                  <div className="bg-white rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100">
                    {/* Template Preview */}
                    <div className={`h-48 ${template.color} flex items-center justify-center relative`}>
                      <div className="text-white text-center">
                        <div className="w-16 h-16 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-4">
                          <ShoppingCart className="w-8 h-8" />
                        </div>
                        <div className="text-lg font-semibold">{template.name}</div>
                      </div>
                      <div className="absolute top-4 right-4 bg-white/20 backdrop-blur-sm rounded-full p-2">
                        <Play className="w-4 h-4 text-white" />
                      </div>
                    </div>

                    {/* Template Details */}
                    <div className="p-6">
                      <div className="flex items-center justify-between mb-4">
                        <h3 className="text-xl font-bold text-gray-900">{template.name}</h3>
                        <span className="text-xs text-orange-600 bg-orange-50 px-2 py-1 rounded-full font-medium">
                          {template.category}
                        </span>
                      </div>
                      <p className="text-gray-600 text-sm mb-4">{template.description}</p>
                      <div className="space-y-2 mb-6">
                        {Array.isArray(template.features) ? 
                          template.features.map((feature, featureIndex) => (
                            <div key={featureIndex} className="flex items-center gap-2">
                              <Check className="w-4 h-4 text-emerald-500" />
                              <span className="text-sm text-gray-600">{feature}</span>
                            </div>
                          )) : (
                            <div className="flex items-center gap-2">
                              <Check className="w-4 h-4 text-emerald-500" />
                              <span className="text-sm text-gray-600">{template.features}</span>
                            </div>
                          )
                        }
                      </div>
                      

                      
                      <div className="flex gap-3">
                        <Link
                          href="/get-started"
                          className="flex-1 bg-gray-900 text-white text-center py-3 rounded-lg font-medium hover:bg-gray-800 transition-colors duration-200"
                        >
                          {t('templatesPage.useTemplate')}
                        </Link>
                        <button className="px-4 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">
                          <Play className="w-4 h-4" />
                        </button>
                      </div>
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        )}

        {/* Features Section */}
        <div className="bg-white rounded-2xl p-12 shadow-sm border mb-16">
          <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">{t('templatesPage.allTemplatesFeatures')}</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {[
              {
                title: t('templatesPage.templateFeatures.mobileFriendly.title'),
                description: t('templatesPage.templateFeatures.mobileFriendly.description')
              },
              {
                title: t('templatesPage.templateFeatures.fastLoad.title'),
                description: t('templatesPage.templateFeatures.fastLoad.description')
              },
              {
                title: t('templatesPage.templateFeatures.seoOptimized.title'),
                description: t('templatesPage.templateFeatures.seoOptimized.description')
              },
              {
                title: t('templatesPage.templateFeatures.customization.title'),
                description: t('templatesPage.templateFeatures.customization.description')
              }
            ].map((feature, index) => (
              <div key={index} className="text-center">
                <div className="w-16 h-16 bg-orange-50 rounded-2xl flex items-center justify-center mx-auto mb-4">
                  <Check className="w-8 h-8 text-orange-600" />
                </div>
                <h3 className="text-lg font-bold text-gray-900 mb-2">{feature.title}</h3>
                <p className="text-gray-600">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>

        {/* CTA Section */}
        <div className="text-center">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="bg-gray-900 rounded-2xl p-12 text-white"
          >
            <h2 className="text-3xl font-bold mb-6">
              {t('templatesPage.startWithTemplate')}
            </h2>
            <p className="text-xl text-gray-300 mb-8">
              {t('templatesPage.quickStart')}
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link
                href="/get-started"
                className="inline-flex items-center gap-2 bg-orange-600 text-white px-8 py-4 rounded-lg hover:bg-orange-700 transition-all duration-300 font-semibold shadow-lg hover:shadow-xl transform hover:scale-105"
              >
                {t('templatesPage.startNowFree')}
                <ArrowRight className="w-5 h-5" />
              </Link>
              <Link
                href="/contact"
                className="inline-flex items-center gap-2 border-2 border-gray-600 text-gray-300 px-8 py-4 rounded-lg hover:bg-gray-800 hover:border-gray-500 transition-all duration-200"
              >
                {t('templatesPage.customDesign')}
              </Link>
            </div>
          </motion.div>
        </div>
      </div>
    </div>
  );
}