import './globals.css'
import { Providers } from './providers'
import Header from '@/components/layout/Header'
import Footer from '@/components/layout/Footer'

export const metadata = {
  title: {
    template: '%s | Hexagonal Ecommerce',
    default: 'Hexagonal Ecommerce - Launch Your Dream Store in Minutes'
  },
  description: 'The most powerful SaaS ecommerce platform with beautiful themes, advanced features, and enterprise-grade hexagonal architecture. Start your free trial today.'
}

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body className="min-h-screen flex flex-col">
        <Providers>
          <Header />
          <main className="flex-grow">
            {children}
          </main>
          <Footer />
        </Providers>
      </body>
    </html>
  )
}