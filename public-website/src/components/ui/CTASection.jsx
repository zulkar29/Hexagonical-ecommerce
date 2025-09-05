'use client';

import { motion } from 'framer-motion';
import { ArrowRight, Phone, Mail } from 'lucide-react';
import Link from 'next/link';

const CTASection = ({ 
  title,
  subtitle,
  description,
  primaryButton,
  secondaryButton,
  variant = 'default',
  background = 'gradient',
  showContactInfo = false
}) => {
  const variants = {
    default: 'text-center',
    split: 'text-left lg:flex lg:items-center lg:justify-between',
    minimal: 'text-center py-12'
  };

  const backgrounds = {
    gradient: 'bg-blue-600',
    dark: 'bg-gray-900',
    light: 'bg-gray-50',
    white: 'bg-white border border-gray-200',
    blue: 'bg-blue-600'
  };

  const textColors = {
    gradient: 'text-white',
    dark: 'text-white',
    light: 'text-gray-900',
    white: 'text-gray-900',
    blue: 'text-white'
  };

  const isLightBackground = background === 'light' || background === 'white';

  return (
    <section className={`py-16 ${backgrounds[background]}`}>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className={variants[variant]}>
          {variant === 'split' ? (
            <>
              <div className="lg:w-2/3">
                <motion.div
                  initial={{ opacity: 0, x: -20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.6 }}
                  viewport={{ once: true }}
                >
                  {subtitle && (
                    <p className={`text-sm font-semibold uppercase tracking-wide mb-2 ${
                      isLightBackground ? 'text-blue-600' : 'text-blue-200'
                    }`}>
                      {subtitle}
                    </p>
                  )}
                  <h2 className={`text-3xl md:text-4xl font-bold mb-4 ${textColors[background]}`}>
                    {title}
                  </h2>
                  {description && (
                    <p className={`text-lg ${
                      isLightBackground ? 'text-gray-600' : 'text-gray-200'
                    }`}>
                      {description}
                    </p>
                  )}
                </motion.div>
              </div>
              
              <motion.div
                initial={{ opacity: 0, x: 20 }}
                whileInView={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.6, delay: 0.2 }}
                viewport={{ once: true }}
                className="mt-8 lg:mt-0 lg:w-1/3 lg:flex lg:justify-end"
              >
                <div className="space-y-4 lg:space-y-0 lg:space-x-4 lg:flex lg:flex-col xl:flex-row">
                  {primaryButton && (
                    <Link
                      href={primaryButton.href || '/contact'}
                      className={`inline-flex items-center justify-center px-8 py-3 rounded-lg font-medium transition-all duration-300 ${
                        isLightBackground
                          ? 'bg-blue-600 text-white hover:bg-blue-700 shadow-lg hover:shadow-xl'
                          : 'bg-white text-gray-900 hover:bg-gray-100 shadow-lg hover:shadow-xl'
                      }`}
                    >
                      {primaryButton.text}
                      <ArrowRight className="w-4 h-4 ml-2" />
                    </Link>
                  )}
                  {secondaryButton && (
                    <Link
                      href={secondaryButton.href || '/about'}
                      className={`inline-flex items-center justify-center px-8 py-3 rounded-lg font-medium border transition-all duration-300 ${
                        isLightBackground
                          ? 'border-gray-300 text-gray-700 hover:bg-gray-50'
                          : 'border-white/30 text-white hover:bg-white/10'
                      }`}
                    >
                      {secondaryButton.text}
                    </Link>
                  )}
                </div>
              </motion.div>
            </>
          ) : (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
              viewport={{ once: true }}
              className="max-w-4xl mx-auto"
            >
              {subtitle && (
                <p className={`text-sm font-semibold uppercase tracking-wide mb-2 ${
                  isLightBackground ? 'text-blue-600' : 'text-blue-200'
                }`}>
                  {subtitle}
                </p>
              )}
              <h2 className={`text-3xl md:text-4xl lg:text-5xl font-bold mb-6 ${textColors[background]}`}>
                {title}
              </h2>
              {description && (
                <p className={`text-lg md:text-xl mb-8 ${
                  isLightBackground ? 'text-gray-600' : 'text-gray-200'
                }`}>
                  {description}
                </p>
              )}
              
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                {primaryButton && (
                  <Link
                    href={primaryButton.href || '/contact'}
                    className={`inline-flex items-center justify-center px-8 py-4 rounded-lg font-medium text-lg transition-all duration-300 ${
                      isLightBackground
                        ? 'bg-blue-600 text-white hover:bg-blue-700 shadow-lg hover:shadow-xl hover:scale-105'
                        : 'bg-white text-gray-900 hover:bg-gray-100 shadow-lg hover:shadow-xl hover:scale-105'
                    }`}
                  >
                    {primaryButton.text}
                    <ArrowRight className="w-5 h-5 ml-2" />
                  </Link>
                )}
                {secondaryButton && (
                  <Link
                    href={secondaryButton.href || '/about'}
                    className={`inline-flex items-center justify-center px-8 py-4 rounded-lg font-medium text-lg border transition-all duration-300 ${
                      isLightBackground
                        ? 'border-gray-300 text-gray-700 hover:bg-gray-50'
                        : 'border-white/30 text-white hover:bg-white/10'
                    }`}
                  >
                    {secondaryButton.text}
                  </Link>
                )}
              </div>
              
              {showContactInfo && (
                <motion.div
                  initial={{ opacity: 0 }}
                  whileInView={{ opacity: 1 }}
                  transition={{ duration: 0.6, delay: 0.4 }}
                  viewport={{ once: true }}
                  className="mt-8 pt-8 border-t border-white/20"
                >
                  <p className={`text-sm mb-4 ${
                    isLightBackground ? 'text-gray-600' : 'text-gray-200'
                  }`}>
                    Or reach out directly:
                  </p>
                  <div className="flex flex-col sm:flex-row gap-4 justify-center">
                    <a
                      href="tel:1-800-HEXAGON"
                      className={`inline-flex items-center ${
                        isLightBackground ? 'text-gray-700 hover:text-blue-600' : 'text-gray-200 hover:text-white'
                      } transition-colors`}
                    >
                      <Phone className="w-4 h-4 mr-2" />
                      1-800-HEXAGON
                    </a>
                    <a
                      href="mailto:hello@storebuilder.com"
                      className={`inline-flex items-center ${
                        isLightBackground ? 'text-gray-700 hover:text-blue-600' : 'text-gray-200 hover:text-white'
                      } transition-colors`}
                    >
                      <Mail className="w-4 h-4 mr-2" />
                      hello@storebuilder.com
                    </a>
                  </div>
                </motion.div>
              )}
            </motion.div>
          )}
        </div>
      </div>
    </section>
  );
};

export default CTASection;