'use client'

import { motion } from 'framer-motion'
import Link from 'next/link'
import { Users, Target, Award, TrendingUp, Shield, Globe, ArrowRight, ShoppingCart, Zap, Heart } from 'lucide-react'
import { useTranslations } from '@/hooks/useTranslations'

export default function AboutPage() {
  const { t } = useTranslations()
  
  const statsArray = [
    { label: 'সফল প্রকল্প', value: '10', suffix: '+' },
    { label: 'টেমপ্লেট', value: '15', suffix: '+' },
    { label: 'জেলায় সেবা', value: '15', suffix: '+' },
    { label: 'ভাষা সাপোর্ট', value: '2', suffix: '' }
  ];

  const values = [
    {
      icon: Target,
      title: 'সহজ ব্যবহার',
      description: 'কোনো প্রযুক্তিগত জ্ঞান ছাড়াই যে কেউ সুন্দর অনলাইন দোকান তৈরি করতে পারেন।'
    },
    {
      icon: Users,
      title: 'গ্রাহক সেবা',
      description: 'আমরা আমাদের গ্রাহকদের সফলতাকে সর্বোচ্চ অগ্রাধিকার দিয়ে ২৪/৭ সাহায্য প্রদান করি।'
    },
    {
      icon: Award,
      title: 'মানসম্মত পণ্য',
      description: 'আমাদের সকল টেমপ্লেট এবং ফিচার উচ্চ মানের এবং আধুনিক প্রযুক্তিতে তৈরি।'
    },
    {
      icon: Shield,
      title: 'নিরাপত্তা',
      description: 'আপনার এবং আপনার গ্রাহকদের তথ্যের সম্পূর্ণ নিরাপত্তা নিশ্চিত করি।'
    }
  ]

  const team = [
    {
      name: 'মোহাম্মদ রহিম',
      position: 'প্রধান নির্বাহী',
      bio: 'ই-কমার্সে ১০+ বছরের অভিজ্ঞতা নিয়ে সবার জন্য অনলাইন ব্যবসাকে সহজ করার স্বপ্ন দেখেন।',
      initial: 'MR'
    },
    {
      name: 'ফাতেমা খান',
      position: 'প্রযুক্তি পরিচালক',
      bio: 'সফটওয়্যার ডেভেলপমেন্টে বিশেষজ্ঞ, আমাদের প্ল্যাটফর্মের প্রযুক্তিগত উন্নতির দায়িত্বে।',
      initial: 'FK'
    },
    {
      name: 'আহমেদ হাসান',
      position: 'ডিজাইন প্রধান',
      bio: 'ব্যবহারকারী অভিজ্ঞতা এবং সুন্দর ডিজাইনের মাধ্যমে গ্রাহক সন্তুষ্টি নিশ্চিত করেন।',
      initial: 'AH'
    },
    {
      name: 'নুসরাত জাহান',
      position: 'গ্রাহক সেবা প্রধান',
      bio: 'গ্রাহকদের সমস্যা সমাধান এবং তাদের সফলতা নিশ্চিত করার জন্য নিরলস কাজ করেন।',
      initial: 'NJ'
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
              আমাদের সম্পর্কে
            </motion.h1>
            <motion.p 
              className="text-xl text-gray-600 max-w-3xl mx-auto mb-12"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              আমরা বাংলাদেশী উদ্যোক্তাদের জন্য সহজ এবং কার্যকর অনলাইন দোকান তৈরির প্ল্যাটফর্ম নিয়ে কাজ করছি। 
              আমাদের লক্ষ্য প্রতিটি ব্যবসায়ীর স্বপ্নের দোকানকে বাস্তব করা।
            </motion.p>
          </div>

          {/* Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            {statsArray.map((stat, index) => (
              <motion.div
                key={stat.label}
                className="text-center bg-gray-50 rounded-2xl p-6 border"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 + 0.4 }}
              >
                <div className="text-3xl md:text-4xl font-bold text-orange-600 mb-2">
                  {stat.value}{stat.suffix}
                </div>
                <div className="text-gray-600 font-medium">{stat.label}</div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Mission & Vision */}
      <section className="py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6 }}
            >
              <div className="flex items-center gap-3 mb-6">
                <div className="w-12 h-12 bg-orange-100 rounded-xl flex items-center justify-center">
                  <Target className="w-6 h-6 text-orange-600" />
                </div>
                <h2 className="text-3xl font-bold text-gray-900">আমাদের লক্ষ্য</h2>
              </div>
              <p className="text-lg text-gray-600 mb-6 leading-relaxed">
                প্রতিটি বাংলাদেশী উদ্যোক্তার হাতে শক্তিশালী অনলাইন ব্যবসার টুলস তুলে দেওয়া। 
                আমরা চাই যেন কোনো প্রযুক্তিগত জ্ঞান ছাড়াই যে কেউ তার স্বপ্নের অনলাইন দোকান 
                তৈরি করতে পারে এবং সফল হতে পারে।
              </p>
              <p className="text-lg text-gray-600 leading-relaxed">
                আমাদের প্ল্যাটফর্ম ব্যবহার করে ব্যবসায়ীরা তাদের পণ্য বিক্রয়, গ্রাহক ব্যবস্থাপনা 
                এবং ব্যবসায়িক বৃদ্ধিতে মনোনিবেশ করতে পারেন। প্রযুক্তিগত জটিলতার ভার আমাদের উপর ছেড়ে দিন।
              </p>
            </motion.div>
            
            <motion.div
              className="bg-white rounded-2xl p-8 border shadow-sm"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-10 h-10 bg-gray-900 rounded-lg flex items-center justify-center">
                  <Heart className="w-5 h-5 text-white" />
                </div>
                <h3 className="text-2xl font-bold text-gray-900">আমাদের দৃষ্টিভঙ্গি</h3>
              </div>
              <p className="text-gray-700 leading-relaxed">
                বাংলাদেশের সবচেয়ে বিশ্বস্ত এবং জনপ্রিয় অনলাইন স্টোর তৈরির প্ল্যাটফর্ম হয়ে ওঠা। 
                আমরা স্বপ্ন দেখি এমন একটি ভবিষ্যতের যেখানে প্রতিটি উদ্যোক্তা মাত্র কয়েক মিনিটেই 
                তার পেশাদার অনলাইন দোকান চালু করতে পারবেন।
              </p>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Values */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              আমাদের মূল্যবোধ
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              যে নীতিমালা আমাদের কাজ এবং গ্রাহক সেবায় অনুপ্রাণিত করে
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {values.map((value, index) => (
              <motion.div
                key={value.title}
                className="bg-gray-50 rounded-2xl p-8 border hover:shadow-lg transition-all duration-300 text-center group"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                <div className="w-16 h-16 bg-orange-100 rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:bg-orange-200 transition-colors">
                  <value.icon className="w-8 h-8 text-orange-600" />
                </div>
                <h3 className="text-xl font-bold text-gray-900 mb-4">{value.title}</h3>
                <p className="text-gray-600 leading-relaxed">{value.description}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Team */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              আমাদের টিম
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              অভিজ্ঞ এবং দক্ষ পেশাদারদের নিয়ে গঠিত আমাদের টিম আপনার সফলতার জন্য নিবেদিত
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {team.map((member, index) => (
              <motion.div
                key={member.name}
                className="text-center bg-white rounded-2xl p-8 border shadow-sm hover:shadow-lg transition-all duration-300"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                <div className="w-20 h-20 bg-gray-900 rounded-2xl flex items-center justify-center text-white text-2xl font-bold mx-auto mb-6">
                  {member.initial}
                </div>
                <h3 className="text-xl font-bold text-gray-900 mb-2">{member.name}</h3>
                <p className="text-orange-600 font-medium mb-4">{member.position}</p>
                <p className="text-gray-600 text-sm leading-relaxed">{member.bio}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-20 bg-gray-900 text-white">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <div className="w-16 h-16 bg-orange-600 rounded-2xl flex items-center justify-center mx-auto mb-6">
              <Zap className="w-8 h-8 text-white" />
            </div>
            <h2 className="text-3xl md:text-4xl font-bold mb-6">
              আমাদের সাথে যাত্রা শুরু করুন
            </h2>
            <p className="text-xl text-gray-300 mb-8 max-w-2xl mx-auto">
              হাজারো সফল উদ্যোক্তার সাথে যোগ দিন যারা আমাদের প্ল্যাটফর্ম দিয়ে তাদের 
              স্বপ্নের ব্যবসা গড়ে তুলেছেন। আজই শুরু করুন আপনার যাত্রা।
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link
                href="/get-started"
                className="inline-flex items-center gap-2 bg-orange-600 text-white px-8 py-4 rounded-lg font-semibold hover:bg-orange-700 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105"
              >
                <ShoppingCart className="w-5 h-5" />
                ফ্রি ট্রায়াল শুরু করুন
                <ArrowRight className="w-5 h-5" />
              </Link>
              <Link
                href="/contact"
                className="inline-flex items-center gap-2 border-2 border-gray-600 text-gray-300 px-8 py-4 rounded-lg font-semibold hover:bg-gray-800 hover:border-gray-500 transition-all duration-300"
              >
                যোগাযোগ করুন
              </Link>
            </div>
          </motion.div>
        </div>
      </section>
    </div>
  )
}