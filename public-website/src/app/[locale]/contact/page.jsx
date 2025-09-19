'use client'

import { motion } from 'framer-motion'
import { Mail, Phone, Clock, Send, MessageSquare, Headphones, Shield, Zap, Globe } from 'lucide-react'
import { useState } from 'react'
import { useTranslations } from '@/hooks/useTranslations'

export default function ContactPage() {
  const { t } = useTranslations()
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    company: '',
    subject: '',
    message: ''
  })
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitStatus, setSubmitStatus] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setIsSubmitting(true)

    // Simulate form submission
    setTimeout(() => {
      setIsSubmitting(false)
      setSubmitStatus('success')
      setFormData({
        name: '',
        email: '',
        company: '',
        subject: '',
        message: ''
      })

      // Reset status after 3 seconds
      setTimeout(() => setSubmitStatus(null), 3000)
    }, 1000)
  }

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
  }

  const contactInfo = [
    {
      icon: Mail,
      title: t('contactPage.contactInfo.email.title'),
      details: t('contactPage.contactInfo.email.details'),
      description: t('contactPage.contactInfo.email.description'),
      action: 'mailto:hello@storebuilder.com'
    },
    {
      icon: Phone,
      title: t('contactPage.contactInfo.phone.title'),
      details: t('contactPage.contactInfo.phone.details'),
      description: t('contactPage.contactInfo.phone.description'),
      action: 'tel:+8801700000000'
    },
    {
      icon: Clock,
      title: t('contactPage.contactInfo.hours.title'),
      details: t('contactPage.contactInfo.hours.details'),
      description: t('contactPage.contactInfo.hours.description')
    }
  ]

  const supportFeatures = [
    {
      icon: Headphones,
      title: t('contactPage.supportFeatures.support.title'),
      description: t('contactPage.supportFeatures.support.description')
    },
    {
      icon: Shield,
      title: t('contactPage.supportFeatures.security.title'),
      description: t('contactPage.supportFeatures.security.description')
    },
    {
      icon: Zap,
      title: t('contactPage.supportFeatures.response.title'),
      description: t('contactPage.supportFeatures.response.description')
    },
    {
      icon: Globe,
      title: t('contactPage.supportFeatures.global.title'),
      description: t('contactPage.supportFeatures.global.description')
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <section className="py-20 bg-gradient-to-br from-orange-50 via-white to-blue-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <motion.div
              className="inline-flex items-center px-4 py-2 rounded-full bg-orange-100 text-orange-800 text-sm font-medium mb-6"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
            >
              <MessageSquare className="w-4 h-4 mr-2" />
              {t('contactPage.badge')}
            </motion.div>
            <motion.h1
              className="text-4xl md:text-6xl font-bold text-gray-900 mb-6"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.1 }}
            >
              {t('contactPage.title')}
            </motion.h1>
            <motion.p
              className="text-xl text-gray-600 max-w-3xl mx-auto mb-12"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              {t('contactPage.description')}
            </motion.p>
          </div>
        </div>
      </section>

      {/* Support Features */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">{t('contactPage.supportFeaturesTitle')}</h2>
            <p className="text-lg text-gray-600 max-w-2xl mx-auto">{t('contactPage.supportFeaturesSubtitle')}</p>
          </div>
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {supportFeatures.map((feature, index) => {
              const IconComponent = feature.icon
              return (
                <motion.div
                  key={feature.title}
                  className="text-center p-6 rounded-xl hover:shadow-lg transition-all duration-300 border border-gray-100 hover:border-orange-200"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <div className="w-12 h-12 bg-gradient-to-br from-orange-100 to-blue-100 rounded-xl flex items-center justify-center mx-auto mb-4">
                    <IconComponent className="w-6 h-6 text-orange-600" />
                  </div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">{feature.title}</h3>
                  <p className="text-sm text-gray-600">{feature.description}</p>
                </motion.div>
              )
            })}
          </div>
        </div>
      </section>

      {/* Contact Info Cards */}
      <section className="py-16 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">{t('contactPage.contactInfoTitle')}</h2>
            <p className="text-lg text-gray-600">{t('contactPage.contactInfoSubtitle')}</p>
          </div>
          <div className="grid md:grid-cols-3 gap-8 mb-20">
            {contactInfo.map((info, index) => {
              const IconComponent = info.icon
              return (
                <motion.div
                  key={info.title}
                  className="text-center p-8 bg-white rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 border border-gray-100 hover:border-orange-200"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <div className="w-16 h-16 bg-gradient-to-br from-orange-100 to-blue-100 rounded-2xl flex items-center justify-center mx-auto mb-6">
                    <IconComponent className="w-8 h-8 text-orange-600" />
                  </div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-3">{info.title}</h3>
                  {info.action ? (
                    <a
                      href={info.action}
                      className="text-orange-600 font-medium text-lg mb-2 block hover:text-orange-700 transition-colors"
                    >
                      {info.details}
                    </a>
                  ) : (
                    <p className="text-orange-600 font-medium text-lg mb-2">{info.details}</p>
                  )}
                  <p className="text-gray-600">{info.description}</p>
                </motion.div>
              )
            })}
          </div>
        </div>
      </section>

      {/* Contact Form */}
      <section className="py-20 bg-white">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">{t('contactPage.form.title')}</h2>
            <p className="text-lg text-gray-600">{t('contactPage.form.subtitle')}</p>
          </div>

          <motion.div
            className="bg-white rounded-3xl p-8 shadow-2xl border border-gray-100"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <label htmlFor="name" className="block text-sm font-semibold text-gray-700 mb-2">
                    {t('contactPage.form.fullName')} *
                  </label>
                  <input
                    type="text"
                    id="name"
                    name="name"
                    required
                    value={formData.name}
                    onChange={handleChange}
                    className="w-full px-4 py-4 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-all duration-200 text-gray-900"
                    placeholder={t('contactPage.form.placeholders.fullName')}
                  />
                </div>
                <div>
                  <label htmlFor="email" className="block text-sm font-semibold text-gray-700 mb-2">
                    {t('contactPage.form.email')} *
                  </label>
                  <input
                    type="email"
                    id="email"
                    name="email"
                    required
                    value={formData.email}
                    onChange={handleChange}
                    className="w-full px-4 py-4 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-all duration-200 text-gray-900"
                    placeholder={t('contactPage.form.placeholders.email')}
                  />
                </div>
              </div>

              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <label htmlFor="company" className="block text-sm font-semibold text-gray-700 mb-2">
                    {t('contactPage.form.company')}
                  </label>
                  <input
                    type="text"
                    id="company"
                    name="company"
                    value={formData.company}
                    onChange={handleChange}
                    className="w-full px-4 py-4 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-all duration-200 text-gray-900"
                    placeholder={t('contactPage.form.placeholders.company')}
                  />
                </div>
                <div>
                  <label htmlFor="subject" className="block text-sm font-semibold text-gray-700 mb-2">
                    {t('contactPage.form.subject')} *
                  </label>
                  <select
                    id="subject"
                    name="subject"
                    required
                    value={formData.subject}
                    onChange={handleChange}
                    className="w-full px-4 py-4 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-all duration-200 text-gray-900"
                  >
                    <option value="">{t('contactPage.form.placeholders.subject')}</option>
                    <option value="store">{t('contactPage.form.subjects.store')}</option>
                    <option value="support">{t('contactPage.form.subjects.support')}</option>
                    <option value="technical">{t('contactPage.form.subjects.technical')}</option>
                    <option value="migration">{t('contactPage.form.subjects.migration')}</option>
                    <option value="partnership">{t('contactPage.form.subjects.partnership')}</option>
                    <option value="other">{t('contactPage.form.subjects.other')}</option>
                  </select>
                </div>
              </div>

              <div>
                <label htmlFor="message" className="block text-sm font-semibold text-gray-700 mb-2">
                  {t('contactPage.form.message')} *
                </label>
                <textarea
                  id="message"
                  name="message"
                  required
                  rows={6}
                  value={formData.message}
                  onChange={handleChange}
                  className="w-full px-4 py-4 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-all duration-200 resize-none text-gray-900"
                  placeholder={t('contactPage.form.placeholders.message')}
                />
              </div>

              <div className="pt-4">
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="w-full bg-gradient-to-r from-orange-600 to-orange-700 text-white py-4 px-8 rounded-xl font-semibold hover:from-orange-700 hover:to-orange-800 transition-all duration-200 flex items-center justify-center shadow-lg hover:shadow-xl transform hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                >
                  {isSubmitting ? (
                    <div className="w-6 h-6 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  ) : (
                    <>
                      <Send className="w-5 h-5 mr-2" />
                      {t('contactPage.form.sendMessage')}
                    </>
                  )}
                </button>
              </div>

              {/* Success Message */}
              {submitStatus === 'success' && (
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="p-6 bg-green-50 border border-green-200 rounded-xl"
                >
                  <p className="text-green-800 font-medium text-center">
                    {t('contactPage.form.successMessage')}
                  </p>
                </motion.div>
              )}
            </form>
          </motion.div>
        </div>
      </section>

      {/* FAQ & Live Chat */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid lg:grid-cols-2 gap-12">
            {/* FAQ */}
            <motion.div
              className="bg-white rounded-2xl p-8 shadow-xl border border-gray-100"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6 }}
            >
              <h3 className="text-2xl font-bold text-gray-900 mb-8">{t('contactPage.faq.title')}</h3>
              <div className="space-y-6">
                <div className="p-4 bg-gray-50 rounded-xl">
                  <h4 className="font-semibold text-gray-900 mb-3">{t('contactPage.faq.questions.response.question')}</h4>
                  <p className="text-gray-600">{t('contactPage.faq.questions.response.answer')}</p>
                </div>
                <div className="p-4 bg-gray-50 rounded-xl">
                  <h4 className="font-semibold text-gray-900 mb-3">{t('contactPage.faq.questions.onboarding.question')}</h4>
                  <p className="text-gray-600">{t('contactPage.faq.questions.onboarding.answer')}</p>
                </div>
                <div className="p-4 bg-gray-50 rounded-xl">
                  <h4 className="font-semibold text-gray-900 mb-3">{t('contactPage.faq.questions.timeline.question')}</h4>
                  <p className="text-gray-600">{t('contactPage.faq.questions.timeline.answer')}</p>
                </div>
              </div>
            </motion.div>

            {/* Live Chat */}
            <motion.div
              className="space-y-8"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              <div className="bg-gradient-to-br from-orange-600 to-orange-700 rounded-2xl p-8 text-white shadow-xl">
                <div className="flex items-center mb-6">
                  <MessageSquare className="w-10 h-10 mr-4" />
                  <h3 className="text-2xl font-bold">{t('contactPage.liveChat.title')}</h3>
                </div>
                <p className="text-orange-100 mb-8 text-lg leading-relaxed">{t('contactPage.liveChat.description')}</p>
                <button className="bg-white text-orange-600 px-8 py-4 rounded-xl font-semibold hover:bg-gray-100 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105">
                  {t('contactPage.liveChat.startChat')}
                </button>
              </div>

              {/* Support Hours */}
              <div className="bg-white rounded-2xl p-8 shadow-xl border border-gray-100">
                <h3 className="text-xl font-bold text-gray-900 mb-6">{t('contactPage.supportHours.title')}</h3>
                <div className="space-y-4">
                  <div className="flex justify-between items-center py-2 border-b border-gray-100">
                    <span className="font-medium text-gray-700">{t('contactPage.supportHours.weekdays')}</span>
                    <span className="text-orange-600 font-semibold">{t('contactPage.supportHours.weekdaysTime')}</span>
                  </div>
                  <div className="flex justify-between items-center py-2 border-b border-gray-100">
                    <span className="font-medium text-gray-700">{t('contactPage.supportHours.saturday')}</span>
                    <span className="text-orange-600 font-semibold">{t('contactPage.supportHours.saturdayTime')}</span>
                  </div>
                  <div className="flex justify-between items-center py-2">
                    <span className="font-medium text-gray-700">{t('contactPage.supportHours.sunday')}</span>
                    <span className="text-gray-500">{t('contactPage.supportHours.sundayTime')}</span>
                  </div>
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </section>
    </div>
  )
}