// Products Catalog Page
import ProductGrid from '@/components/product/ProductGrid'
import ProductFilters from '@/components/product/ProductFilters'
import Breadcrumbs from '@/components/ui/Breadcrumbs'

export const metadata = {
  title: 'All Products',
  description: 'Browse our complete product catalog with filters and search.'
}

export default function ProductsPage({ searchParams }) {
  const { category, search, sort, price } = searchParams

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Breadcrumbs */}
      <Breadcrumbs 
        items={[
          { label: 'Home', href: '/' },
          { label: 'Products', href: '/products' }
        ]} 
      />

      {/* Page Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-4">All Products</h1>
        {search && (
          <p className="text-gray-600">
            Showing results for "{search}"
          </p>
        )}
      </div>

      <div className="lg:grid lg:grid-cols-4 lg:gap-8">
        {/* Filters Sidebar */}
        <div className="lg:col-span-1">
          <ProductFilters />
        </div>

        {/* Products Grid */}
        <div className="lg:col-span-3">
          <ProductGrid 
            category={category}
            search={search}
            sort={sort}
            price={price}
          />
        </div>
      </div>
    </div>
  )
}