'use client';

import { motion } from 'framer-motion';
import { ArrowRight, Play, CheckCircle } from 'lucide-react';
import Link from 'next/link';
import Image from 'next/image';

const HeroSection = ({ 
  title,
  subtitle,
  description,
  primaryButton,
  secondaryButton,
  features = [],
  image,
  video,
  variant = 'default',
  background = 'gradient'
}) => {
  const variants = {
    default: 'text-center',
    split: 'lg:grid lg:grid-cols-2 lg:gap-12 lg:items-center',
    minimal: 'text-center py-20'
  };

  const backgrounds = {
    gradient: 'bg-blue-600',
    dark: 'bg-gray-900',
    light: 'bg-gray-50',
    white: 'bg-white'
  };

  const textColors = {
    gradient: 'text-white',
    dark: 'text-white',
    light: 'text-gray-900',
    white: 'text-gray-900'
  };

  const isLightBackground = background === 'light' || background === 'white';

  return (
    <section className={`relative py-20 lg:py-32 overflow-hidden ${backgrounds[background]}`}>
      {/* Background Pattern */}
      {background === 'gradient' && (
        <div className="absolute inset-0 bg-black/10">
          <div className="absolute inset-0" style={{
            backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Ccircle cx='30' cy='30' r='2'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
          }} />
        </div>
      )}
      
      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className={variants[variant]}>
          {/* Content */}
          <div className={variant === 'split' ? 'lg:pr-8' : 'max-w-4xl mx-auto'}>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
            >
              {subtitle && (
                <motion.p
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.6, delay: 0.1 }}
                  className={`text-sm font-semibold uppercase tracking-wide mb-4 ${
                    isLightBackground ? 'text-blue-600' : 'text-blue-200'
                  }`}
                >
                  {subtitle}
                </motion.p>
              )}
              
              <motion.h1
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: 0.2 }}
                className={`text-4xl md:text-5xl lg:text-6xl font-bold mb-6 leading-tight ${textColors[background]}`}
              >
                {title}
              </motion.h1>
              
              {description && (
                <motion.p
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.6, delay: 0.3 }}
                  className={`text-lg md:text-xl mb-8 leading-relaxed ${
                    isLightBackground ? 'text-gray-600' : 'text-gray-200'
                  }`}
                >
                  {description}
                </motion.p>
              )}
              
              {/* Features List */}
              {features.length > 0 && (
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.6, delay: 0.4 }}
                  className="mb-8"
                >
                  <ul className={`space-y-2 ${variant === 'default' ? 'max-w-md mx-auto' : ''}`}>
                    {features.map((feature, index) => (
                      <li key={index} className={`flex items-center ${
                        variant === 'default' ? 'justify-center' : 'justify-start'
                      }`}>
                        <CheckCircle className={`w-5 h-5 mr-3 ${
                          isLightBackground ? 'text-green-600' : 'text-green-400'
                        }`} />
                        <span className={isLightBackground ? 'text-gray-700' : 'text-gray-200'}>
                          {feature}
                        </span>
                      </li>
                    ))}
                  </ul>
                </motion.div>
              )}
              
              {/* Buttons */}
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: 0.5 }}
                className={`flex flex-col sm:flex-row gap-4 ${
                  variant === 'default' ? 'justify-center' : 'justify-start'
                }`}
              >
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
                    href={secondaryButton.href || '#'}
                    className={`inline-flex items-center justify-center px-8 py-4 rounded-lg font-medium text-lg border transition-all duration-300 ${
                      isLightBackground
                        ? 'border-gray-300 text-gray-700 hover:bg-gray-50'
                        : 'border-white/30 text-white hover:bg-white/10'
                    }`}
                  >
                    {secondaryButton.icon === 'play' && <Play className="w-5 h-5 mr-2" />}
                    {secondaryButton.text}
                  </Link>
                )}
              </motion.div>
            </motion.div>
          </div>
          
          {/* Media */}
          {variant === 'split' && (image || video) && (
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.8, delay: 0.3 }}
              className="mt-12 lg:mt-0"
            >
              {video ? (
                <div className="relative rounded-xl overflow-hidden shadow-2xl">
                  <video
                    autoPlay
                    muted
                    loop
                    playsInline
                    className="w-full h-auto"
                  >
                    <source src={video} type="video/mp4" />
                  </video>
                </div>
              ) : image ? (
                <div className="relative rounded-xl overflow-hidden shadow-2xl">
                  <Image
                    src={image.src}
                    alt={image.alt || 'Hero image'}
                    width={600}
                    height={400}
                    className="w-full h-auto object-cover"
                    priority
                  />
                </div>
              ) : null}
            </motion.div>
          )}
        </div>
        
        {/* Stats or Trust Indicators */}
        {variant === 'default' && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.6 }}
            className="mt-16 pt-8 border-t border-white/20"
          >
            <p className={`text-sm mb-6 ${
              isLightBackground ? 'text-gray-600' : 'text-gray-300'
            }`}>
              Trusted by leading companies worldwide
            </p>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-8 items-center opacity-60">
              {/* Placeholder for company logos */}
              <div className={`h-8 rounded ${
                isLightBackground ? 'bg-gray-300' : 'bg-white/20'
              }`} />
              <div className={`h-8 rounded ${
                isLightBackground ? 'bg-gray-300' : 'bg-white/20'
              }`} />
              <div className={`h-8 rounded ${
                isLightBackground ? 'bg-gray-300' : 'bg-white/20'
              }`} />
              <div className={`h-8 rounded ${
                isLightBackground ? 'bg-gray-300' : 'bg-white/20'
              }`} />
            </div>
          </motion.div>
        )}
      </div>
    </section>
  );
};

export default HeroSection;