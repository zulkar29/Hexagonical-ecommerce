'use client';

import { usePathname } from 'next/navigation';
import { useState, useEffect } from 'react';

export function useTranslations() {
  const pathname = usePathname();
  const [translations, setTranslations] = useState({});

  // Extract locale from pathname
  const locale = pathname.split('/')[1] || 'en';

  useEffect(() => {
    const loadTranslations = async () => {
      try {
        const messages = await import(`@/messages/${locale}.json`);
        setTranslations(messages.default);
      } catch (error) {
        console.error('Failed to load translations:', error);
        // Fallback to English
        const fallback = await import('@/messages/en.json');
        setTranslations(fallback.default);
      }
    };

    loadTranslations();
  }, [locale]);

  const t = (key) => {
    const keys = key.split('.');
    let value = translations;

    for (const k of keys) {
      if (value && typeof value === 'object' && k in value) {
        value = value[k];
      } else {
        return key; // Return key if translation not found
      }
    }

    return value || key;
  };

  return { t, locale };
}