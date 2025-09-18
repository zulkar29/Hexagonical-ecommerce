'use client'

import { motion } from 'framer-motion'
import Link from 'next/link'
import { Check, X, Star, ArrowRight, Crown } from 'lucide-react'
import { useState } from 'react'
import { useTranslations } from '@/hooks/useTranslations'

export default function PricingPage() {
  const { t, locale } = useTranslations()
  const [billingCycle, setBillingCycle] = useState('monthly')

  const pricingPlans = [
    {
      name: 'free',
      price: 0,
      isFree: true
    },
    {
      name: 'starter',
      price: 990
    },
    {
      name: 'professional',
      price: 2990,
      popular: true
    },
    {
      name: 'pro',
      price: 4990
    }
  ]

  // Helper function to safely get features array
  const getFeatures = (key) => {
    try {
      const features = t(key);
      return Array.isArray(features) ? features : [];
    } catch {
      return [];
    }
  };

  const addOns = [
    {
      name: 'অতিরিক্ত স্টোরেজ',
      description: 'পণ্যের ছবি এবং ফাইলের জন্য ১০০জিবি',
      price: '৳৪৯০/মাস'
    },
    {
      name: 'উন্নত মার্কেটিং টুলস',
      description: 'ইমেইল ক্যাম্পেইন ও সোশ্যাল মিডিয়া',
      price: '৳১৯৯০/মাস'
    },
    {
      name: 'কাস্টম থিম ডিজাইন',
      description: 'পেশাদার কাস্টম থিম তৈরি',
      price: '৳৯৯৯০/একবার'
    },
    {
      name: 'মাইগ্রেশন সেবা',
      description: 'অন্য প্ল্যাটফর্ম থেকে সম্পূর্ণ স্থানান্তর',
      price: '৳৪৯৯০/একবার'
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <motion.h1 
              className="text-4xl md:text-5xl font-bold text-gray-900 mb-6"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
            >
              {t('pricing.title')}
            </motion.h1>
            <motion.p
              className="text-xl text-gray-600 max-w-3xl mx-auto mb-12"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              {t('pricing.description')}
            </motion.p>

            {/* Billing Toggle */}
            <motion.div 
              className="flex items-center justify-center mb-12"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.4 }}
            >
              <span className={`mr-3 ${billingCycle === 'monthly' ? 'text-gray-900 font-semibold' : 'text-gray-500'}`}>
                মাসিক
              </span>
              <button
                onClick={() => setBillingCycle(billingCycle === 'monthly' ? 'yearly' : 'monthly')}
                className="relative inline-flex h-6 w-11 items-center rounded-full bg-gray-200 transition-colors focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    billingCycle === 'yearly' ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
              <span className={`ml-3 ${billingCycle === 'yearly' ? 'text-gray-900 font-semibold' : 'text-gray-500'}`}>
                বার্ষিক
                <span className="ml-1 text-emerald-600 text-sm font-medium">(২০% সাশ্রয়)</span>
              </span>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Pricing Plans */}
      <section className="py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Paid Plans Grid */}
          <div className="grid md:grid-cols-3 gap-8 mb-16">
            {pricingPlans.filter(plan => !plan.isFree).map((plan, index) => {
              const planFeatures = getFeatures(`getStarted.plans.${plan.name}.features`)
              const isPopular = plan.popular

              return (
                <motion.div
                  key={plan.name}
                  className={`relative bg-white rounded-2xl shadow-sm border-2 p-6 ${
                    isPopular ? 'border-orange-500 scale-105' : 'border-gray-100'
                  }`}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  {isPopular && (
                    <div className="absolute -top-4 left-1/2 transform -translate-x-1/2">
                      <span className="bg-orange-600 text-white px-4 py-2 rounded-full text-sm font-semibold flex items-center">
                        <Crown className="w-4 h-4 mr-1" />
                        {t('getStarted.plan.recommended')}
                      </span>
                    </div>
                  )}

                  <div className="text-center mb-6">
                    <h3 className="text-2xl font-bold text-gray-900 mb-4">{t(`getStarted.plans.${plan.name}.name`)}</h3>
                    <div className="mb-4">
                      <span className="text-4xl font-bold text-gray-900">
                        {billingCycle === 'yearly'
                          ? `৳${Math.round(plan.price * 0.8)}`
                          : `৳${plan.price}`
                        }
                      </span>
                      <span className="text-gray-600 ml-1">
                        /{billingCycle === 'yearly' ? t('getStarted.plan.yearly') : t('getStarted.plan.monthly')}
                      </span>
                    </div>
                    <p className="text-gray-600">{t(`getStarted.plans.${plan.name}.description`)}</p>
                  </div>

                  <ul className="space-y-3 mb-8">
                    {planFeatures.slice(0, 8).map((feature, idx) => (
                      <li key={idx} className="flex items-start">
                        <Check className="w-4 h-4 text-emerald-500 mr-3 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700 text-sm">{feature}</span>
                      </li>
                    ))}
                  </ul>

                  <Link
                    href={`/${locale}/get-started`}
                    className={`block w-full py-3 px-6 rounded-lg font-semibold text-center transition-colors duration-200 ${
                      isPopular
                        ? 'bg-orange-600 text-white hover:bg-orange-700'
                        : 'bg-gray-100 text-gray-900 hover:bg-gray-200'
                    }`}
                  >
                    {t('common.getStarted')}
                  </Link>
                </motion.div>
              )
            })}
          </div>

          {/* Free Trial Option - Full Width */}
          <div className="border-t pt-12">
            <div className="text-center mb-8">
              <h3 className="text-2xl font-bold text-gray-900 mb-2">Or start with our free trial</h3>
              <p className="text-gray-600">No credit card required • Cancel anytime</p>
            </div>

            {pricingPlans.filter(plan => plan.isFree).map((plan) => {
              const planFeatures = getFeatures(`getStarted.plans.${plan.name}.features`)

              return (
                <motion.div
                  key={plan.name}
                  className="bg-white rounded-2xl shadow-sm border-2 border-emerald-200 p-8"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.3 }}
                >
                  <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-6">
                    <div className="flex items-center gap-6 flex-1">
                      <div className="flex items-center gap-4">
                        <div className="w-16 h-16 bg-emerald-100 rounded-2xl flex items-center justify-center">
                          <Star className="w-8 h-8 text-emerald-600" />
                        </div>
                        <div>
                          <h3 className="text-3xl font-bold text-gray-900 mb-2">{t(`getStarted.plans.${plan.name}.name`)}</h3>
                          <p className="text-gray-600">{t(`getStarted.plans.${plan.name}.description`)}</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="text-4xl font-bold text-emerald-600 mb-1">FREE</div>
                        <div className="text-emerald-600 font-medium">14 days trial</div>
                      </div>
                    </div>

                    <div className="flex items-center gap-8">
                      <div className="hidden lg:flex items-center gap-6 text-sm text-gray-600">
                        {planFeatures.slice(0, 4).map((feature, index) => (
                          <div key={index} className="flex items-center gap-2">
                            <Check className="w-4 h-4 text-emerald-500" />
                            <span>{feature}</span>
                          </div>
                        ))}
                      </div>

                      <Link
                        href={`/${locale}/get-started`}
                        className="bg-emerald-600 text-white px-8 py-4 rounded-lg font-semibold hover:bg-emerald-700 transition-colors duration-200 whitespace-nowrap"
                      >
                        Start Free Trial
                      </Link>
                    </div>
                  </div>
                </motion.div>
              )
            })}
          </div>
        </div>
      </section>

      {/* Feature Comparison */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              প্ল্যান তুলনা
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              প্রতিটি প্ল্যানে কী কী আছে দেখুন এবং আপনার জন্য উপযুক্ত প্ল্যান বেছে নিন।
            </p>
          </div>

          <div className="bg-white rounded-2xl shadow-sm border overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">Feature</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">{t('getStarted.plans.free.name')}</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">{t('getStarted.plans.starter.name')}</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">{t('getStarted.plans.professional.name')}</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">{t('getStarted.plans.pro.name')}</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {[
                    { feature: 'Products', free: '50', starter: '500', professional: '2000', pro: '10000' },
                    { feature: 'Storage', free: '1 GB', starter: '5 GB', professional: '20 GB', pro: '100 GB' },
                    { feature: 'Staff Accounts', free: '1', starter: '1', professional: '3', pro: '10' },
                    { feature: 'Support', free: 'Email', starter: 'Email', professional: 'Priority', pro: 'Priority' },
                    { feature: 'Custom Domain', free: false, starter: true, professional: true, pro: true },
                    { feature: 'Advanced Analytics', free: false, starter: false, professional: true, pro: true },
                    { feature: 'Email Marketing', free: false, starter: false, professional: false, pro: true },
                    { feature: 'API Access', free: false, starter: false, professional: false, pro: true }
                  ].map((row, index) => (
                    <tr key={index} className="hover:bg-gray-50">
                      <td className="px-6 py-4 text-sm font-medium text-gray-900">{row.feature}</td>
                      <td className="px-6 py-4 text-center">
                        {typeof row.free === 'boolean' ? (
                          row.free ? (
                            <Check className="w-5 h-5 text-emerald-500 mx-auto" />
                          ) : (
                            <X className="w-5 h-5 text-gray-300 mx-auto" />
                          )
                        ) : (
                          <span className="text-sm text-gray-700">{row.free}</span>
                        )}
                      </td>
                      <td className="px-6 py-4 text-center">
                        {typeof row.starter === 'boolean' ? (
                          row.starter ? (
                            <Check className="w-5 h-5 text-emerald-500 mx-auto" />
                          ) : (
                            <X className="w-5 h-5 text-gray-300 mx-auto" />
                          )
                        ) : (
                          <span className="text-sm text-gray-700">{row.starter}</span>
                        )}
                      </td>
                      <td className="px-6 py-4 text-center">
                        {typeof row.professional === 'boolean' ? (
                          row.professional ? (
                            <Check className="w-5 h-5 text-emerald-500 mx-auto" />
                          ) : (
                            <X className="w-5 h-5 text-gray-300 mx-auto" />
                          )
                        ) : (
                          <span className="text-sm text-gray-700">{row.professional}</span>
                        )}
                      </td>
                      <td className="px-6 py-4 text-center">
                        {typeof row.pro === 'boolean' ? (
                          row.pro ? (
                            <Check className="w-5 h-5 text-emerald-500 mx-auto" />
                          ) : (
                            <X className="w-5 h-5 text-gray-300 mx-auto" />
                          )
                        ) : (
                          <span className="text-sm text-gray-700">{row.pro}</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </section>

      {/* Add-ons */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              অতিরিক্ত সেবা
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              আপনার প্রয়োজন অনুযায়ী অতিরিক্ত ফিচার এবং সেবা যোগ করুন।
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
            {addOns.map((addon, index) => (
              <motion.div
                key={addon.name}
                className="bg-white rounded-xl p-6 shadow-sm hover:shadow-lg transition-shadow duration-300 border border-gray-100"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                <h3 className="text-lg font-semibold text-gray-900 mb-3">{addon.name}</h3>
                <p className="text-gray-600 text-sm mb-4">{addon.description}</p>
                <div className="flex items-center justify-between">
                  <span className="text-xl font-bold text-orange-600">{addon.price}</span>
                  <button className="text-orange-600 hover:text-orange-700 font-medium text-sm flex items-center">
                    যোগ করুন
                    <ArrowRight className="w-4 h-4 ml-1" />
                  </button>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* FAQ */}
      <section className="py-20 bg-white">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              সাধারণ প্রশ্ন
            </h2>
          </div>

          <div className="space-y-6">
            {[
              {
                question: 'আমি কি যেকোনো সময় প্ল্যান পরিবর্তন করতে পারি?',
                answer: 'হ্যাঁ, আপনি যেকোনো সময় প্ল্যান আপগ্রেড বা ডাউনগ্রেড করতে পারেন। পরিবর্তনটি পরবর্তী বিলিং চক্রে প্রতিফলিত হবে।'
              },
              {
                question: 'কি ফ্রি ট্রায়াল আছে?',
                answer: 'সব প্ল্যানে ১৪ দিনের ফ্রি ট্রায়াল রয়েছে। আপনি সব ফিচার পরীক্ষা করতে পারেন। কোনো ক্রেডিট কার্ডের প্রয়োজন নেই।'
              },
              {
                question: 'কি কি পেমেন্ট পদ্ধতি গ্রহণযোগ্য?',
                answer: 'আমরা বিকাশ, নগদ, ব্যাংক ট্রান্সফার এবং সব প্রধান ক্রেডিট/ডেবিট কার্ড গ্রহণ করি।'
              },
              {
                question: 'রিফান্ড নীতি কি?',
                answer: 'সব পেইড প্ল্যানের জন্য ৩০ দিনের মানি-ব্যাক গ্যারান্টি রয়েছে। সন্তুষ্ট না হলে পূর্ণ রিফান্ড পাবেন।'
              },
              {
                question: 'বাতিল করলে আমার ডেটার কি হবে?',
                answer: 'আপনি যেকোনো সময় বাতিল করতে পারেন। আপনার স্টোর বিলিং পিরিয়ডের শেষ পর্যন্ত চালু থাকবে এবং সব ডেটা এক্সপোর্ট করতে পারবেন।'
              }
            ].map((faq, index) => (
              <motion.div
                key={index}
                className="bg-gray-50 rounded-xl p-6 border border-gray-100"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                <h3 className="text-lg font-semibold text-gray-900 mb-3">{faq.question}</h3>
                <p className="text-gray-600">{faq.answer}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>
    </div>
  )
}