'use client';

import { motion } from 'framer-motion';
import { ArrowRight } from 'lucide-react';
import Link from 'next/link';

const ServiceCard = ({ service, index, variant = 'default' }) => {
  const cardVariants = {
    default: 'bg-white border border-gray-200 hover:shadow-lg',
    featured: 'bg-blue-50 border border-blue-200 hover:shadow-xl',
    minimal: 'bg-transparent border-0 hover:bg-gray-50'
  };

  const iconVariants = {
    default: 'bg-blue-100 text-blue-600',
    featured: 'bg-blue-600 text-white',
    minimal: 'bg-gray-100 text-gray-600 group-hover:bg-blue-100 group-hover:text-blue-600'
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.6, delay: index * 0.1 }}
      viewport={{ once: true }}
      className={`group rounded-xl p-6 transition-all duration-300 ${cardVariants[variant]}`}
    >
      <div className={`w-12 h-12 rounded-lg flex items-center justify-center mb-4 transition-colors duration-300 ${iconVariants[variant]}`}>
        <service.icon className="w-6 h-6" />
      </div>
      
      <h3 className="text-xl font-semibold text-gray-900 mb-3 group-hover:text-blue-600 transition-colors">
        {service.title}
      </h3>
      
      <p className="text-gray-600 mb-4 leading-relaxed">
        {service.description}
      </p>
      
      {service.features && (
        <ul className="space-y-2 mb-6">
          {service.features.slice(0, 3).map((feature, idx) => (
            <li key={idx} className="flex items-center text-sm text-gray-600">
              <div className="w-1.5 h-1.5 bg-blue-600 rounded-full mr-3"></div>
              {feature}
            </li>
          ))}
        </ul>
      )}
      
      {service.price && (
        <div className="mb-4">
          <span className="text-2xl font-bold text-gray-900">{service.price}</span>
          {service.period && (
            <span className="text-gray-600 ml-1">/{service.period}</span>
          )}
        </div>
      )}
      
      <Link
        href={service.href || '/contact'}
        className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium group-hover:translate-x-1 transition-all duration-300"
      >
        {service.cta || 'Learn More'}
        <ArrowRight className="w-4 h-4 ml-2" />
      </Link>
    </motion.div>
  );
};

const ServiceGrid = ({ services, title, subtitle, variant = 'default', columns = 3 }) => {
  const gridCols = {
    2: 'grid-cols-1 md:grid-cols-2',
    3: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
    4: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4'
  };

  return (
    <section className="py-16">
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
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
                {title}
              </h2>
            )}
            {subtitle && (
              <p className="text-xl text-gray-600 max-w-3xl mx-auto">
                {subtitle}
              </p>
            )}
          </motion.div>
        )}
        
        <div className={`grid ${gridCols[columns]} gap-8`}>
          {services.map((service, index) => (
            <ServiceCard
              key={service.id || index}
              service={service}
              index={index}
              variant={variant}
            />
          ))}
        </div>
      </div>
    </section>
  );
};

export { ServiceCard, ServiceGrid };
export default ServiceCard;