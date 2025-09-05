'use client'

import { motion } from 'framer-motion'
import Link from 'next/link'
import { Check, X, Star, ArrowRight, Crown } from 'lucide-react'
import { useState } from 'react'
import { useTranslations } from '@/hooks/useTranslations'

export default function PricingPage() {
  const { t } = useTranslations()
  const [billingCycle, setBillingCycle] = useState('monthly')

  const pricingPlans = [
    {
      name: 'starter',
      price: 990,
      description: 'নতুন ব্যবসার জন্য উপযুক্ত'
    },
    {
      name: 'professional', 
      price: 2990,
      description: 'ক্রমবর্ধমান ব্যবসার জন্য',
      popular: true
    },
    {
      name: 'pro',
      price: 4990,
      description: 'প্রতিষ্ঠিত ব্যবসার জন্য'
    },
    {
      name: 'enterprise',
      price: null,
      description: 'বড় কোম্পানির জন্য'
    }
  ]

  const features = {
    starter: [
      '৫০০টি পর্যন্ত পণ্য',
      '৫জিবি স্টোরেজ',
      'বেসিক অ্যানালিটিক্স',
      'ইমেইল সাপোর্ট',
      'এসএসএল সার্টিফিকেট',
      'মোবাইল ফ্রেন্ডলি',
      'কাস্টম ডোমেইন',
      'বেসিক এসইও টুলস'
    ],
    professional: [
      '২০০০টি পর্যন্ত পণ্য',
      '২০জিবি স্টোরেজ',
      'অ্যাডভান্সড অ্যানালিটিক্স',
      'অগ্রাধিকার সাপোর্ট',
      'কাস্টম ডোমেইন',
      'অ্যাডভান্সড এসইও',
      'মাল্টি-কারেন্সি সাপোর্ট',
      'পরিত্যক্ত কার্ট রিকভারি',
      'পেমেন্ট গেটওয়ে',
      'ইনভেন্টরি ব্যবস্থাপনা'
    ],
    pro: [
      '১০০০০টি পর্যন্ত পণ্য',
      '১০০জিবি স্টোরেজ',
      'অ্যাডভান্সড রিপোর্ট',
      'অগ্রাধিকার সাপোর্ট',
      'মাল্টি-স্টোর ব্যবস্থাপনা',
      'ইন্টিগ্রেশন সাপোর্ট',
      'অ্যাডভান্সড ইন্টিগ্রেশন',
      'কাস্টম থিম',
      'মার্কেটিং অটোমেশন',
      'অ্যাডভান্সড নিরাপত্তা'
    ],
    enterprise: [
      'আনলিমিটেড পণ্য',
      'আনলিমিটেড স্টোরেজ',
      'কাস্টম ড্যাশবোর্ড',
      '২৪/৭ ডেডিকেটেড সাপোর্ট',
      'হোয়াইট-লেবেল সমাধান',
      'কাস্টম ইন্টিগ্রেশন',
      'অ্যাডভান্সড নিরাপত্তা',
      'এসএলএ গ্যারান্টি',
      'মাল্টি-স্টোর ব্যবস্থাপনা',
      'কাস্টম ডেভেলপমেন্ট'
    ]
  }

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
              আপনার ব্যবসার জন্য সেরা প্ল্যান বেছে নিন। সব প্ল্যানে ১৪ দিন ফ্রি ট্রায়াল এবং কোনো সেটআপ ফি নেই।
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
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {pricingPlans.map((plan, index) => {
              const planFeatures = features[plan.name] || []
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
                        জনপ্রিয়
                      </span>
                    </div>
                  )}

                  <div className="text-center mb-6">
                    <h3 className="text-2xl font-bold text-gray-900 mb-4">{t(`getStarted.plans.${plan.name}.name`)}</h3>
                    <div className="mb-4">
                      <span className="text-4xl font-bold text-gray-900">
                        {plan.price === null 
                          ? 'যোগাযোগ করুন'
                          : billingCycle === 'yearly' 
                            ? `৳${Math.round(plan.price * 0.8)}`
                            : `৳${plan.price}`
                        }
                      </span>
                      {plan.price !== null && (
                        <span className="text-gray-600 ml-1">
                          /{billingCycle === 'yearly' ? 'মাস (বার্ষিক বিলিং)' : 'মাস'}
                        </span>
                      )}
                    </div>
                    <p className="text-gray-600">{plan.description}</p>
                  </div>

                  <ul className="space-y-3 mb-8">
                    {planFeatures.map((feature, idx) => (
                      <li key={idx} className="flex items-start">
                        <Check className="w-4 h-4 text-emerald-500 mr-3 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700 text-sm">{feature}</span>
                      </li>
                    ))}
                  </ul>

                  <Link 
                    href="/get-started"
                    className={`block w-full py-3 px-6 rounded-lg font-semibold text-center transition-colors duration-200 ${
                      isPopular
                        ? 'bg-orange-600 text-white hover:bg-orange-700'
                        : 'bg-gray-100 text-gray-900 hover:bg-gray-200'
                    }`}
                  >
                    শুরু করুন
                  </Link>
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
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">ফিচার</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">স্টার্টার</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">প্রফেশনাল</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">প্রো</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">এন্টারপ্রাইজ</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {[
                    { feature: 'পণ্য সংখ্যা', basic: '৫০০', professional: '২০০০', pro: '১০০০০', enterprise: 'আনলিমিটেড' },
                    { feature: 'স্টোরেজ', basic: '৫ জিবি', professional: '২০ জিবি', pro: '১০০ জিবি', enterprise: 'আনলিমিটেড' },
                    { feature: 'সাপোর্ট', basic: 'ইমেইল', professional: 'অগ্রাধিকার', pro: 'অগ্রাধিকার', enterprise: '২৪/৭ ডেডিকেটেড' },
                    { feature: 'উন্নত অ্যানালিটিক্স', basic: false, professional: true, pro: true, enterprise: true },
                    { feature: 'মাল্টি-স্টোর ব্যবস্থাপনা', basic: false, professional: false, pro: true, enterprise: true },
                    { feature: 'ইন্টিগ্রেশন সাপোর্ট', basic: false, professional: false, pro: true, enterprise: true },
                    { feature: 'হোয়াইট-লেবেল', basic: false, professional: false, pro: false, enterprise: true },
                    { feature: 'কাস্টম ডেভেলপমেন্ট', basic: false, professional: false, pro: false, enterprise: true }
                  ].map((row, index) => (
                    <tr key={index} className="hover:bg-gray-50">
                      <td className="px-6 py-4 text-sm font-medium text-gray-900">{row.feature}</td>
                      <td className="px-6 py-4 text-center">
                        {typeof row.basic === 'boolean' ? (
                          row.basic ? (
                            <Check className="w-5 h-5 text-emerald-500 mx-auto" />
                          ) : (
                            <X className="w-5 h-5 text-gray-300 mx-auto" />
                          )
                        ) : (
                          <span className="text-sm text-gray-700">{row.basic}</span>
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
                      <td className="px-6 py-4 text-center">
                        {typeof row.enterprise === 'boolean' ? (
                          row.enterprise ? (
                            <Check className="w-5 h-5 text-emerald-500 mx-auto" />
                          ) : (
                            <X className="w-5 h-5 text-gray-300 mx-auto" />
                          )
                        ) : (
                          <span className="text-sm text-gray-700">{row.enterprise}</span>
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

      {/* CTA */}
      <section className="py-20 bg-gray-900">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">
            আপনার স্টোর চালু করতে প্রস্তুত?
          </h2>
          <p className="text-xl text-gray-300 mb-8 max-w-2xl mx-auto">
            হাজারো ব্যবসায়ীর সাথে যোগ দিন যারা আমাদের বিশ্বাস করেন। আজই ফ্রি ট্রায়াল শুরু করুন।
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/get-started"
              className="bg-orange-600 text-white px-8 py-4 rounded-lg font-semibold hover:bg-orange-700 transition-colors duration-200"
            >
              ফ্রি ট্রায়াল শুরু করুন
            </Link>
            <Link
              href="/contact"
              className="border-2 border-gray-600 text-gray-300 px-8 py-4 rounded-lg font-semibold hover:bg-gray-800 hover:border-gray-500 transition-colors duration-200"
            >
              যোগাযোগ করুন
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}