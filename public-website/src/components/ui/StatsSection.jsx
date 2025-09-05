'use client';

import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';

const CountUpNumber = ({ end, duration = 2000, suffix = '', prefix = '' }) => {
  const [count, setCount] = useState(0);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (!isVisible) return;

    let startTime;
    const startValue = 0;
    const endValue = end;

    const animate = (currentTime) => {
      if (!startTime) startTime = currentTime;
      const progress = Math.min((currentTime - startTime) / duration, 1);
      
      const easeOutQuart = 1 - Math.pow(1 - progress, 4);
      const currentCount = Math.floor(easeOutQuart * (endValue - startValue) + startValue);
      
      setCount(currentCount);
      
      if (progress < 1) {
        requestAnimationFrame(animate);
      }
    };

    requestAnimationFrame(animate);
  }, [isVisible, end, duration]);

  return (
    <motion.div
      initial={{ opacity: 0 }}
      whileInView={{ opacity: 1 }}
      onViewportEnter={() => setIsVisible(true)}
      viewport={{ once: true }}
      className="text-3xl md:text-4xl font-bold text-white"
    >
      {prefix}{count.toLocaleString()}{suffix}
    </motion.div>
  );
};

const StatCard = ({ stat, index, variant = 'default' }) => {
  const variants = {
    default: 'bg-white text-gray-900 shadow-lg',
    dark: 'bg-gray-800 text-white',
    gradient: 'bg-blue-600 text-white',
    minimal: 'bg-transparent text-gray-900 border border-gray-200'
  };

  const iconVariants = {
    default: 'bg-blue-100 text-blue-600',
    dark: 'bg-gray-700 text-blue-400',
    gradient: 'bg-white/20 text-white',
    minimal: 'bg-blue-50 text-blue-600'
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.6, delay: index * 0.1 }}
      viewport={{ once: true }}
      className={`rounded-xl p-6 text-center ${variants[variant]}`}
    >
      {stat.icon && (
        <div className={`w-12 h-12 rounded-lg flex items-center justify-center mx-auto mb-4 ${iconVariants[variant]}`}>
          <stat.icon className="w-6 h-6" />
        </div>
      )}
      
      <div className="mb-2">
        {variant === 'gradient' ? (
          <CountUpNumber 
            end={stat.value} 
            suffix={stat.suffix || ''} 
            prefix={stat.prefix || ''}
          />
        ) : (
          <div className="text-3xl md:text-4xl font-bold">
            {stat.prefix || ''}{stat.value.toLocaleString()}{stat.suffix || ''}
          </div>
        )}
      </div>
      
      <h3 className="text-lg font-semibold mb-2">{stat.label}</h3>
      
      {stat.description && (
        <p className={`text-sm ${variant === 'gradient' || variant === 'dark' ? 'text-gray-200' : 'text-gray-600'}`}>
          {stat.description}
        </p>
      )}
    </motion.div>
  );
};

const StatsSection = ({ 
  stats, 
  title, 
  subtitle, 
  variant = 'default',
  columns = 4,
  background = 'default'
}) => {
  const gridCols = {
    2: 'grid-cols-1 md:grid-cols-2',
    3: 'grid-cols-1 md:grid-cols-3',
    4: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4'
  };

  const backgrounds = {
    default: 'bg-white',
    gray: 'bg-gray-50',
    dark: 'bg-gray-900',
    gradient: 'bg-blue-600'
  };

  const textColors = {
    default: 'text-gray-900',
    gray: 'text-gray-900',
    dark: 'text-white',
    gradient: 'text-white'
  };

  return (
    <section className={`py-16 ${backgrounds[background]}`}>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {(title || subtitle) && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            viewport={{ once: true }}
            className="text-center mb-12"
          >
            {title && (
              <h2 className={`text-3xl md:text-4xl font-bold mb-4 ${textColors[background]}`}>
                {title}
              </h2>
            )}
            {subtitle && (
              <p className={`text-xl max-w-3xl mx-auto ${
                background === 'dark' || background === 'gradient' ? 'text-gray-200' : 'text-gray-600'
              }`}>
                {subtitle}
              </p>
            )}
          </motion.div>
        )}
        
        <div className={`grid ${gridCols[columns]} gap-8`}>
          {stats.map((stat, index) => (
            <StatCard
              key={stat.id || index}
              stat={stat}
              index={index}
              variant={variant}
            />
          ))}
        </div>
      </div>
    </section>
  );
};

export { StatCard, StatsSection, CountUpNumber };
export default StatsSection;