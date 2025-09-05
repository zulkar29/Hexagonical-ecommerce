'use client';

import { useRouter, usePathname } from 'next/navigation';
import { useState } from 'react';

export default function LanguageSwitcher() {
  const router = useRouter();
  const pathname = usePathname();

  // Extract locale from pathname
  const locale = pathname.split('/')[1] || 'en';
  const isEnglish = locale === 'en';

  const toggleLanguage = () => {
    const newLocale = isEnglish ? 'bn' : 'en';
    // Remove the current locale from the pathname
    const pathWithoutLocale = pathname.replace(`/${locale}`, '') || '/';
    // Navigate to the new locale
    router.push(`/${newLocale}${pathWithoutLocale}`);
  };

  return (
    <div className="flex items-center">
      <button
        onClick={toggleLanguage}
        className="flex items-center space-x-2 px-3 py-2 rounded-lg border border-gray-200 hover:border-gray-300 hover:bg-gray-50 transition-all duration-200"
        aria-label="Switch language"
      >
        <span className="text-sm font-medium text-gray-700">
          {isEnglish ? '🇺🇸 EN' : '🇧🇩 বাং'}
        </span>
      </button>
    </div>
  );
}