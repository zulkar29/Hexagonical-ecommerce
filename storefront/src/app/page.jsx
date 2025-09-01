// E-commerce Homepage
import HeroSection from '@/components/layout/HeroSection'
import FeaturedProducts from '@/components/product/FeaturedProducts'
import CategoryGrid from '@/components/layout/CategoryGrid'
import Newsletter from '@/components/layout/Newsletter'

export default function HomePage() {
  return (
    <>
      {/* Hero Banner */}
      <HeroSection />
      
      {/* Featured Products */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12">
            Featured Products
          </h2>
          <FeaturedProducts />
        </div>
      </section>

      {/* Shop by Category */}
      <section className="py-16 bg-gray-50">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12">
            Shop by Category
          </h2>
          <CategoryGrid />
        </div>
      </section>

      {/* Newsletter */}
      <Newsletter />
    </>
  )
}