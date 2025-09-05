'use client'

import { motion } from 'framer-motion'
import { Mail, Phone, MapPin, Clock, Send, MessageSquare } from 'lucide-react'
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

  const handleSubmit = (e) => {
    e.preventDefault()
    // Handle form submission here
    console.log('Form submitted:', formData)
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
      description: t('contactPage.contactInfo.email.description')
    },
    {
      icon: Phone,
      title: t('contactPage.contactInfo.phone.title'),
      details: t('contactPage.contactInfo.phone.details'),
      description: t('contactPage.contactInfo.phone.description')
    },
    {
      icon: MapPin,
      title: t('contactPage.contactInfo.address.title'),
      details: t('contactPage.contactInfo.address.details'),
      description: t('contactPage.contactInfo.address.description')
    },
    {
      icon: Clock,
      title: t('contactPage.contactInfo.hours.title'),
      details: t('contactPage.contactInfo.hours.details'),
      description: t('contactPage.contactInfo.hours.description')
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <section className="py-20 bg-gradient-to-br from-orange-50 to-gray-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <motion.h1 
              className="text-4xl md:text-5xl font-bold text-gray-900 mb-6"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
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

      {/* Contact Info Cards */}
      <section className="py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8 mb-20">
            {contactInfo.map((info, index) => {
              const IconComponent = info.icon
              return (
                <motion.div
                  key={info.title}
                  className="text-center p-6 bg-white rounded-xl shadow-lg hover:shadow-xl transition-shadow duration-300 border border-gray-100"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <div className="w-12 h-12 bg-orange-100 rounded-lg flex items-center justify-center mx-auto mb-4">
                    <IconComponent className="w-6 h-6 text-orange-600" />
                  </div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">{info.title}</h3>
                  <p className="text-orange-600 font-medium mb-1">{info.details}</p>
                  <p className="text-sm text-gray-600">{info.description}</p>
                </motion.div>
              )
            })}
          </div>
        </div>
      </section>

      {/* Contact Form & Map */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid lg:grid-cols-2 gap-12">
            {/* Contact Form */}
            <motion.div
              className="bg-white rounded-2xl p-8 shadow-xl border border-gray-100"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6 }}
            >
              <div className="mb-8">
                <h2 className="text-2xl font-bold text-gray-900 mb-4">{t('contactPage.form.title')}</h2>
                <p className="text-gray-600">{t('contactPage.form.subtitle')}</p>
              </div>

              <form onSubmit={handleSubmit} className="space-y-6">
                <div className="grid md:grid-cols-2 gap-6">
                  <div>
                    <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('contactPage.form.fullName')}
                    </label>
                    <input
                      type="text"
                      id="name"
                      name="name"
                      required
                      value={formData.name}
                      onChange={handleChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-colors duration-200"
                      placeholder={t('contactPage.form.placeholders.fullName')}
                    />
                  </div>
                  <div>
                    <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('contactPage.form.email')}
                    </label>
                    <input
                      type="email"
                      id="email"
                      name="email"
                      required
                      value={formData.email}
                      onChange={handleChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-colors duration-200"
                      placeholder={t('contactPage.form.placeholders.email')}
                    />
                  </div>
                </div>

                <div className="grid md:grid-cols-2 gap-6">
                  <div>
                    <label htmlFor="company" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('contactPage.form.company')}
                    </label>
                    <input
                      type="text"
                      id="company"
                      name="company"
                      value={formData.company}
                      onChange={handleChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-colors duration-200"
                      placeholder={t('contactPage.form.placeholders.company')}
                    />
                  </div>
                  <div>
                    <label htmlFor="subject" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('contactPage.form.subject')}
                    </label>
                    <select
                      id="subject"
                      name="subject"
                      required
                      value={formData.subject}
                      onChange={handleChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-colors duration-200"
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
                  <label htmlFor="message" className="block text-sm font-medium text-gray-700 mb-2">
                    {t('contactPage.form.message')}
                  </label>
                  <textarea
                    id="message"
                    name="message"
                    required
                    rows={6}
                    value={formData.message}
                    onChange={handleChange}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-colors duration-200 resize-none"
                    placeholder={t('contactPage.form.placeholders.message')}
                  />
                </div>

                <button
                  type="submit"
                  className="w-full bg-orange-600 text-white py-4 px-6 rounded-lg font-semibold hover:bg-orange-700 transition-colors duration-200 flex items-center justify-center shadow-lg hover:shadow-xl transform hover:scale-105"
                >
                  <Send className="w-5 h-5 mr-2" />
                  {t('contactPage.form.sendMessage')}
                </button>
              </form>
            </motion.div>

            {/* Map & Additional Info */}
            <motion.div
              className="space-y-8"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              {/* Map Placeholder */}
              <div className="bg-white rounded-2xl p-8 shadow-xl border border-gray-100">
                <h3 className="text-xl font-bold text-gray-900 mb-6">{t('contactPage.location.title')}</h3>
                <div className="bg-gray-100 rounded-lg h-64 flex items-center justify-center">
                  <div className="text-center">
                    <MapPin className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                    <p className="text-gray-600">{t('contactPage.location.mapText')}</p>
                    <p className="text-sm text-gray-500">{t('contactPage.contactInfo.address.details')}</p>
                    <p className="text-sm text-gray-500">{t('contactPage.contactInfo.address.description')}</p>
                  </div>
                </div>
              </div>

              {/* FAQ */}
              <div className="bg-white rounded-2xl p-8 shadow-xl border border-gray-100">
                <h3 className="text-xl font-bold text-gray-900 mb-6">{t('contactPage.faq.title')}</h3>
                <div className="space-y-4">
                  <div>
                    <h4 className="font-semibold text-gray-900 mb-2">{t('contactPage.faq.questions.response.question')}</h4>
                    <p className="text-gray-600 text-sm">{t('contactPage.faq.questions.response.answer')}</p>
                  </div>
                  <div>
                    <h4 className="font-semibold text-gray-900 mb-2">{t('contactPage.faq.questions.onboarding.question')}</h4>
                <p className="text-gray-600 text-sm">{t('contactPage.faq.questions.onboarding.answer')}</p>
                  </div>
                  <div>
                    <h4 className="font-semibold text-gray-900 mb-2">{t('contactPage.faq.questions.timeline.question')}</h4>
                    <p className="text-gray-600 text-sm">{t('contactPage.faq.questions.timeline.answer')}</p>
                  </div>
                </div>
              </div>

              {/* Live Chat */}
              <div className="bg-orange-600 rounded-2xl p-8 text-white">
                <div className="flex items-center mb-4">
                  <MessageSquare className="w-8 h-8 mr-3" />
                  <h3 className="text-xl font-bold">{t('contactPage.liveChat.title')}</h3>
                </div>
                <p className="text-orange-100 mb-6">{t('contactPage.liveChat.description')}</p>
                <button className="bg-white text-orange-600 px-6 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors duration-200 shadow-lg hover:shadow-xl transform hover:scale-105">
                  {t('contactPage.liveChat.startChat')}
                </button>
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Office Hours */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl font-bold text-gray-900 mb-8">{t('contactPage.officeHours.title')}</h2>
          <div className="grid md:grid-cols-2 gap-8">
            <div className="bg-white rounded-xl p-6 shadow-lg border border-gray-100">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">{t('contactPage.officeHours.regular.title')}</h3>
              <div className="space-y-2 text-gray-600">
                <div className="flex justify-between">
                  <span>{t('contactPage.officeHours.regular.monday')}</span>
                  <span>{t('contactPage.officeHours.regular.mondayTime')}</span>
                </div>
                <div className="flex justify-between">
                  <span>{t('contactPage.officeHours.regular.saturday')}</span>
                  <span>{t('contactPage.officeHours.regular.saturdayTime')}</span>
                </div>
                <div className="flex justify-between">
                  <span>{t('contactPage.officeHours.regular.sunday')}</span>
                  <span>{t('contactPage.officeHours.regular.sundayTime')}</span>
                </div>
              </div>
            </div>
            <div className="bg-white rounded-xl p-6 shadow-lg border border-gray-100">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">{t('contactPage.officeHours.emergency.title')}</h3>
              <div className="space-y-2 text-gray-600">
                <div className="flex justify-between">
                  <span>{t('contactPage.officeHours.emergency.line')}</span>
                  <span>{t('contactPage.officeHours.emergency.phone')}</span>
                </div>
                <div className="flex justify-between">
                  <span>{t('contactPage.officeHours.emergency.responseTime')}</span>
                  <span>{t('contactPage.officeHours.emergency.responseValue')}</span>
                </div>
                <div className="flex justify-between">
                  <span>{t('contactPage.officeHours.emergency.available')}</span>
                  <span>{t('contactPage.officeHours.emergency.availableValue')}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}