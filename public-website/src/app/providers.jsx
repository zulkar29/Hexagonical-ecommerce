'use client';

import { Provider } from 'jotai';

export function Providers({ children }) {
  return (
    <Provider>
      {children}
    </Provider>
  );
}