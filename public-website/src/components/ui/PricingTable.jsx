'use client';

import { motion } from 'framer-motion';
import { Check, X, Star } from 'lucide-react';
import Link from 'next/link';

const PricingCard = ({ plan, index, isPopular = false }) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.6, delay: index * 0.1 }}
      viewport={{ once: true }}
      className={`relative rounded-2xl p-8 ${isPopular 
        ? 'bg-blue-50 border-2 border-blue-500 shadow-xl scale-105' 
        : 'bg-white border border-gray-200 shadow-lg'
      }`}
    >
      {isPopular && (
        <div className="absolute -top-4 left-1/2 transform -translate-x-1/2">
          <div className="bg-blue-600 text-white px-4 py-2 rounded-full text-sm font-medium flex items-center">
            <Star className="w-4 h-4 mr-1" />
            Most Popular
          </div>
        </div>
      )}
      
      <div className="text-center mb-8">
        <h3 className="text-2xl font-bold text-gray-900 mb-2">{plan.name}</h3>
        <p className="text-gray-600 mb-4">{plan.description}</p>
        
        <div className="mb-4">
          <span className="text-4xl font-bold text-gray-900">${plan.price}</span>
          {plan.period && (
            <span className="text-gray-600 ml-1">/{plan.period}</span>
          )}
        </div>
        
        {plan.originalPrice && (
          <div className="text-sm text-gray-500">
            <span className="line-through">${plan.originalPrice}</span>
            <span className="ml-2 text-green-600 font-medium">
              Save ${plan.originalPrice - plan.price}
            </span>
          </div>
        )}
      </div>
      
      <div className="space-y-4 mb-8">
        {plan.features.map((feature, idx) => (
          <div key={idx} className="flex items-start">
            {feature.included ? (
              <Check className="w-5 h-5 text-green-500 mr-3 mt-0.5 flex-shrink-0" />
            ) : (
              <X className="w-5 h-5 text-gray-400 mr-3 mt-0.5 flex-shrink-0" />
            )}
            <span className={`text-sm ${feature.included ? 'text-gray-700' : 'text-gray-400'}`}>
              {feature.text}
            </span>
          </div>
        ))}
      </div>
      
      <Link
        href={plan.href || '/contact'}
        className={`block w-full text-center py-3 px-6 rounded-lg font-medium transition-all duration-300 ${
          isPopular
            ? 'bg-blue-600 text-white hover:bg-blue-700 shadow-lg hover:shadow-xl'
            : 'bg-gray-900 text-white hover:bg-gray-800'
        }`}
      >
        {plan.cta || 'Get Started'}
      </Link>
      
      {plan.note && (
        <p className="text-xs text-gray-500 text-center mt-4">{plan.note}</p>
      )}
    </motion.div>
  );
};

const PricingTable = ({ 
  plans, 
  title = "Choose Your Plan", 
  subtitle = "Select the perfect plan for your business needs",
  showBillingToggle = false,
  billingCycle = 'monthly',
  onBillingChange
}) => {
  return (
    <section className="py-16 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          viewport={{ once: true }}
          className="text-center mb-12"
        >
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            {title}
          </h2>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto mb-8">
            {subtitle}
          </p>
          
          {showBillingToggle && (
            <div className="flex items-center justify-center space-x-4">
              <span className={`text-sm ${billingCycle === 'monthly' ? 'text-gray-900 font-medium' : 'text-gray-500'}`}>
                Monthly
              </span>
              <button
                onClick={() => onBillingChange && onBillingChange(billingCycle === 'monthly' ? 'yearly' : 'monthly')}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  billingCycle === 'yearly' ? 'bg-blue-600' : 'bg-gray-200'
                }`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    billingCycle === 'yearly' ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
              <span className={`text-sm ${billingCycle === 'yearly' ? 'text-gray-900 font-medium' : 'text-gray-500'}`}>
                Yearly
                <span className="ml-1 text-green-600 font-medium">(Save 20%)</span>
              </span>
            </div>
          )}
        </motion.div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 max-w-5xl mx-auto">
          {plans.map((plan, index) => (
            <PricingCard
              key={plan.id || index}
              plan={plan}
              index={index}
              isPopular={plan.popular}
            />
          ))}
        </div>
        
        <motion.div
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.4 }}
          viewport={{ once: true }}
          className="text-center mt-12"
        >
          <p className="text-gray-600 mb-4">
            Need a custom solution? We can create a plan tailored to your specific needs.
          </p>
          <Link
            href="/contact"
            className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
          >
            Contact us for enterprise pricing
          </Link>
        </motion.div>
      </div>
    </section>
  );
};

export { PricingCard, PricingTable };
export default PricingTable;