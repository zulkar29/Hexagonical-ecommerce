'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { ArrowRight, Play, Check, Star, Zap, ShoppingCart, Settings, BarChart3, Users, Shield, Globe, Smartphone, CreditCard, Package, TrendingUp, Clock, Award } from 'lucide-react';
import { motion } from 'framer-motion';
import { useTranslations } from '@/hooks/useTranslations';

export default function HomePage() {
  const { t } = useTranslations();

  const [animatedStats, setAnimatedStats] = useState({
    customers: 0,
    stores: 0,
    revenue: 0,
    countries: 0
  });

  const stats = {
    customers: 150,
    stores: 75,
    revenue: 500000,
    countries: 15
  };

  // Animate stats on mount
  useEffect(() => {
    const animateStats = () => {
      const duration = 2000;
      const steps = 60;
      const stepDuration = duration / steps;
      
      let step = 0;
      const timer = setInterval(() => {
        step++;
        const progress = step / steps;
        
        setAnimatedStats({
          customers: Math.floor(stats.customers * progress),
          stores: Math.floor(stats.stores * progress),
          revenue: Math.floor(stats.revenue * progress),
          countries: Math.floor(stats.countries * progress)
        });
        
        if (step >= steps) clearInterval(timer);
      }, stepDuration);
    };
    
    animateStats();
  }, []);

  return (
    <div className="min-h-screen">
      {/* Hero Section - Professional & Clean */}
      <section className="relative py-20 lg:py-32 bg-white overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            {/* Content */}
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.8 }}
            >
              <div className="inline-flex items-center gap-2 bg-amber-50 text-amber-700 px-4 py-2 rounded-full text-sm font-medium mb-6">
                <Award className="w-4 h-4" />
                {t('hero.badge')}
              </div>
              
              <h1 className="text-5xl lg:text-6xl font-bold text-gray-900 mb-6 leading-tight">
                {t('hero.title')}
                <span className="text-orange-600"> {t('hero.titleHighlight')}</span>
              </h1>
              
              <p className="text-xl text-gray-600 mb-8 leading-relaxed">
                {t('hero.subtitle')}
              </p>

              {/* Key Benefits */}
              <div className="space-y-3 mb-8">
                {[
                  t('common.noCodingRequired'),
                  t('common.support247'),
                  t('common.freeDomainSSL')
                ].map((benefit, index) => (
                  <div key={index} className="flex items-center gap-3">
                    <div className="w-2 h-2 bg-orange-600 rounded-full"></div>
                    <span className="text-gray-700">{benefit}</span>
                  </div>
                ))}
              </div>
              
              <div className="flex flex-col sm:flex-row gap-4 mb-8">
                <Link
                  href="/get-started"
                  className="bg-gray-900 text-white px-8 py-4 rounded-lg font-semibold hover:bg-gray-800 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105 inline-flex items-center justify-center"
                >
                  {t('common.getStarted')}
                  <ArrowRight className="w-5 h-5 ml-2" />
                </Link>
                <button className="border-2 border-gray-300 text-gray-700 px-8 py-4 rounded-lg font-semibold hover:border-gray-400 hover:bg-gray-50 transition-all duration-300 inline-flex items-center justify-center">
                  <Play className="w-5 h-5 mr-2" />
                  {t('common.watchDemo')}
                </button>
              </div>

              {/* Trust indicators */}
              <div className="flex items-center gap-6 text-sm text-gray-500">
                <div className="flex items-center gap-2">
                  <div className="flex">
                    {[...Array(5)].map((_, i) => (
                      <Star key={i} className="w-4 h-4 text-amber-400 fill-current" />
                    ))}
                  </div>
                  <span>{t('common.rating48')}</span>
                </div>
                <div className="w-px h-4 bg-gray-300"></div>
                <span>{t('common.projectCount')} {t('common.successfulProjects')}</span>
              </div>
            </motion.div>

            {/* Visual */}
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.8, delay: 0.2 }}
              className="relative"
            >
              <div className="bg-gray-50 rounded-2xl p-8 shadow-2xl border">
                <div className="space-y-6">
                  {/* Mock dashboard header */}
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-orange-600 rounded-lg flex items-center justify-center">
                        <ShoppingCart className="w-6 h-6 text-white" />
                      </div>
                      <div>
                        <div className="h-3 bg-gray-300 rounded w-24 mb-1"></div>
                        <div className="h-2 bg-gray-200 rounded w-16"></div>
                      </div>
                    </div>
                    <div className="w-8 h-8 bg-gray-300 rounded-full"></div>
                  </div>

                  {/* Mock stats */}
                  <div className="grid grid-cols-3 gap-4">
                    {[
                      { label: t('common.orders'), value: t('common.ordersCount'), color: 'bg-emerald-100 text-emerald-700' },
                      { label: t('common.sales'), value: t('common.salesAmount'), color: 'bg-amber-100 text-amber-700' },
                      { label: t('common.customers'), value: t('common.customersCount'), color: 'bg-purple-100 text-purple-700' }
                    ].map((stat, index) => (
                      <div key={index} className="bg-white rounded-lg p-3 border">
                        <div className="text-lg font-bold text-gray-900">{stat.value}</div>
                        <div className="text-xs text-gray-500">{stat.label}</div>
                      </div>
                    ))}
                  </div>

                  {/* Mock chart */}
                  <div className="bg-white rounded-lg p-4 border">
                    <div className="flex items-end gap-2 h-24">
                      {[40, 60, 35, 80, 45, 70, 55].map((height, index) => (
                        <div
                          key={index}
                          className="bg-orange-200 rounded-t flex-1"
                          style={{ height: `${height}%` }}
                        ></div>
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              {/* Floating elements */}
              <div className="absolute -top-4 -right-4 bg-white rounded-lg shadow-lg p-3 border">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></div>
                  <span className="text-xs font-medium">{t('common.liveStore')}</span>
                </div>
              </div>
              
              <div className="absolute -bottom-4 -left-4 bg-white rounded-lg shadow-lg p-3 border">
                <div className="flex items-center gap-2">
                  <TrendingUp className="w-4 h-4 text-orange-600" />
                  <span className="text-xs font-medium">{t('common.growth85')}</span>
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-6">
              {t('homepage.features.title')}
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              {t('homepage.features.description')}
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {[
              {
                icon: ShoppingCart,
                title: t('homepage.featureList.easyStore.title'),
                description: t('homepage.featureList.easyStore.description'),
                features: t('homepage.featureList.easyStore.features'),
                color: "text-orange-600",
                bgColor: "bg-orange-50"
              },
              {
                icon: CreditCard,
                title: t('homepage.featureList.securePayment.title'),
                description: t('homepage.featureList.securePayment.description'),
                features: t('homepage.featureList.securePayment.features'),
                color: "text-emerald-600",
                bgColor: "bg-emerald-50"
              },
              {
                icon: BarChart3,
                title: t('homepage.featureList.analytics.title'),
                description: t('homepage.featureList.analytics.description'),
                features: t('homepage.featureList.analytics.features'),
                color: "text-purple-600",
                bgColor: "bg-purple-50"
              },
              {
                icon: Package,
                title: t('homepage.featureList.inventory.title'),
                description: t('homepage.featureList.inventory.description'),
                features: t('homepage.featureList.inventory.features'),
                color: "text-amber-600",
                bgColor: "bg-amber-50"
              },
              {
                icon: Smartphone,
                title: t('homepage.featureList.mobileApp.title'),
                description: t('homepage.featureList.mobileApp.description'),
                features: t('homepage.featureList.mobileApp.features'),
                color: "text-cyan-600",
                bgColor: "bg-cyan-50"
              },
              {
                icon: Users,
                title: t('homepage.featureList.support247.title'),
                description: t('homepage.featureList.support247.description'),
                features: t('homepage.featureList.support247.features'),
                color: "text-rose-600",
                bgColor: "bg-rose-50"
              }
            ].map((feature, index) => {
              const IconComponent = feature.icon;
              return (
                <motion.div
                  key={index}
                  className="group"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <div className="bg-white rounded-2xl p-8 shadow-sm hover:shadow-lg transition-all duration-300 border border-gray-100 group-hover:border-gray-200 h-full">
                    <div className={`w-16 h-16 ${feature.bgColor} rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300`}>
                      <IconComponent className={`w-8 h-8 ${feature.color}`} />
                    </div>
                    <h3 className="text-2xl font-bold text-gray-900 mb-4">{feature.title}</h3>
                    <p className="text-gray-600 leading-relaxed mb-6">{feature.description}</p>
                    <div className="space-y-2">
                      {Array.isArray(feature.features) ? feature.features.map((item, featureIndex) => (
                        <div key={featureIndex} className="flex items-center gap-2">
                          <Check className="w-4 h-4 text-emerald-500" />
                          <span className="text-sm text-gray-600">{item}</span>
                        </div>
                      )) : null}
                    </div>
                  </div>
                </motion.div>
              );
            })}
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-6">
              {t('homepage.stats.title')}
            </h2>
            <p className="text-xl text-gray-600">
              {t('homepage.stats.description')}
            </p>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            {[
              { key: 'customers', label: t('homepage.stats.happyCustomers'), suffix: '+', icon: Users, color: 'text-orange-600' },
              { key: 'stores', label: t('homepage.stats.activeStores'), suffix: '+', icon: ShoppingCart, color: 'text-emerald-600' },
              { key: 'revenue', label: t('homepage.stats.totalSales'), suffix: '', icon: TrendingUp, color: 'text-purple-600' },
              { key: 'countries', label: t('homepage.stats.districts'), suffix: '+', icon: Globe, color: 'text-amber-600' }
            ].map((stat, index) => {
              const IconComponent = stat.icon;
              return (
                <motion.div
                  key={index}
                  className="text-center"
                  initial={{ opacity: 0, scale: 0.5 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ delay: index * 0.1, duration: 0.6 }}
                >
                  <div className="bg-gray-50 rounded-2xl p-8 border">
                    <IconComponent className={`w-12 h-12 ${stat.color} mx-auto mb-4`} />
                    <div className="text-4xl font-bold text-gray-900 mb-2">
                      {stat.key === 'revenue' 
                        ? `৳${(animatedStats[stat.key] / 100000).toFixed(1)}লাখ` 
                        : animatedStats[stat.key]}{stat.suffix}
                    </div>
                    <div className="text-gray-600 font-medium">{stat.label}</div>
                  </div>
                </motion.div>
              );
            })}
          </div>
        </div>
      </section>

      {/* Template Showcase */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-6">
              {t('homepage.templates.title')}
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              {t('homepage.templates.description')}
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 mb-12">
            {[
              {
                name: t('getStarted.templates.fashion.name'),
                category: t('getStarted.templates.fashion.category'),
                color: "bg-rose-600",
                features: t('getStarted.templates.fashion.features'),
                description: t('getStarted.templates.fashion.description')
              },
              {
                name: t('getStarted.templates.electronics.name'),
                category: t('getStarted.templates.electronics.category'),
                color: "bg-indigo-600",
                features: t('getStarted.templates.electronics.features'),
                description: t('getStarted.templates.electronics.description')
              },
              {
                name: t('getStarted.templates.food.name'),
                category: t('getStarted.templates.food.category'),
                color: "bg-emerald-600",
                features: t('getStarted.templates.food.features'),
                description: t('getStarted.templates.food.description')
              }
            ].map((theme, index) => (
              <motion.div
                key={index}
                className="group cursor-pointer"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                whileHover={{ y: -5 }}
              >
                <div className="bg-white rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100">
                  <div className={`h-48 ${theme.color} flex items-center justify-center relative`}>
                    <div className="text-white text-center">
                      <div className="w-16 h-16 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-4">
                        <ShoppingCart className="w-8 h-8" />
                      </div>
                      <div className="text-lg font-semibold">{theme.name}</div>
                    </div>
                    <div className="absolute top-4 right-4 bg-white/20 backdrop-blur-sm rounded-full p-2">
                      <Play className="w-4 h-4 text-white" />
                    </div>
                  </div>
                  <div className="p-6">
                    <div className="flex items-center justify-between mb-4">
                      <h3 className="text-xl font-bold text-gray-900">{theme.name}</h3>
                      <span className="text-xs text-orange-600 bg-orange-50 px-2 py-1 rounded-full font-medium">{theme.category}</span>
                    </div>
                    <p className="text-gray-600 text-sm mb-4">{theme.description}</p>
                    <div className="space-y-2">
                      {Array.isArray(theme.features) ? theme.features.map((feature, featureIndex) => (
                        <div key={featureIndex} className="flex items-center gap-2">
                          <Check className="w-4 h-4 text-emerald-500" />
                          <span className="text-sm text-gray-600">{feature}</span>
                        </div>
                      )) : null}
                    </div>
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
          
          <div className="text-center">
            <Link
              href="/templates"
              className="inline-flex items-center gap-2 bg-gray-900 text-white px-8 py-4 rounded-lg hover:bg-gray-800 transition-all duration-300 font-semibold shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              {t('homepage.templates.viewAllTemplates')}
              <ArrowRight className="w-5 h-5" />
            </Link>
          </div>
        </div>
      </section>

      {/* Get Started Flow */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-6">
              {t('homepage.stepsSection.title')}
            </h2>
            <p className="text-xl text-gray-600">
              {t('homepage.stepsSection.description')}
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-12">
            {[
              {
                step: "১",
                titleKey: "homepage.stepsDetails.step1.title",
                descriptionKey: "homepage.stepsDetails.step1.description",
                icon: Globe,
                color: "bg-orange-600"
              },
              {
                step: "২", 
                titleKey: "homepage.stepsDetails.step2.title",
                descriptionKey: "homepage.stepsDetails.step2.description",
                icon: ShoppingCart,
                color: "bg-emerald-600"
              },
              {
                step: "৩",
                titleKey: "homepage.stepsDetails.step3.title",
                descriptionKey: "homepage.stepsDetails.step3.description",
                icon: CreditCard,
                color: "bg-purple-600"
              }
            ].map((step, index) => {
              const IconComponent = step.icon;
              return (
                <motion.div
                  key={index}
                  className="text-center"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.2 }}
                >
                  <div className="relative">
                    <div className={`w-20 h-20 ${step.color} rounded-2xl flex items-center justify-center mx-auto mb-6 shadow-lg`}>
                      <IconComponent className="w-10 h-10 text-white" />
                    </div>
                    <div className="absolute -top-2 -right-2 w-8 h-8 bg-gray-900 text-white rounded-full flex items-center justify-center font-bold text-lg">
                      {step.step}
                    </div>
                  </div>
                  <h3 className="text-2xl font-bold text-gray-900 mb-4">{t(step.titleKey)}</h3>
                  <p className="text-gray-600 leading-relaxed">{t(step.descriptionKey)}</p>
                </motion.div>
              );
            })}
          </div>

          <div className="text-center">
            <Link
              href="/get-started"
              className="inline-flex items-center gap-2 bg-gray-900 text-white px-12 py-5 rounded-lg hover:bg-gray-800 transition-all duration-300 font-bold text-lg shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              <Clock className="w-6 h-6" />
              {t('homepage.stepsSection.startNowFree')}
              <ArrowRight className="w-6 h-6" />
            </Link>
            <p className="text-sm text-gray-500 mt-4">{t('homepage.stepsSection.noCreditCard')}</p>
          </div>
        </div>
      </section>
    </div>
  );
}