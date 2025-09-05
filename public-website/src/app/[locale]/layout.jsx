import Header from '@/components/layout/Header';
import Footer from '@/components/layout/Footer';
import { Providers } from '../providers';
import '../globals.css';

export default async function RootLayout({
  children,
  params,
}) {
  const { locale } = await params;

  return (
    <html lang={locale}>
      <body>
        <Providers>
          <div className="min-h-screen flex flex-col">
            <Header />
            <main className="flex-grow">
              {children}
            </main>
            <Footer />
          </div>
        </Providers>
      </body>
    </html>
  );
}