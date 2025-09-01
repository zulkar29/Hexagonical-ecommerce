// Individual Product Page
import ProductImages from '@/components/product/ProductImages'
import ProductInfo from '@/components/product/ProductInfo'
import ProductTabs from '@/components/product/ProductTabs'
import RelatedProducts from '@/components/product/RelatedProducts'
import Breadcrumbs from '@/components/ui/Breadcrumbs'
import { notFound } from 'next/navigation'

// Generate metadata for SEO
export async function generateMetadata({ params }) {
  // TODO: Fetch product data
  const product = await getProduct(params.slug)
  
  if (!product) {
    return {}
  }

  return {
    title: product.name,
    description: product.description,
    openGraph: {
      title: product.name,
      description: product.description,
      images: [product.image],
    },
  }
}

async function getProduct(slug) {
  // TODO: Implement API call
  // const response = await fetch(`${process.env.API_URL}/products/${slug}`)
  // return response.json()
  return null
}

export default async function ProductPage({ params }) {
  const product = await getProduct(params.slug)

  if (!product) {
    notFound()
  }

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Breadcrumbs */}
      <Breadcrumbs 
        items={[
          { label: 'Home', href: '/' },
          { label: 'Products', href: '/products' },
          { label: product.category, href: `/products?category=${product.category}` },
          { label: product.name }
        ]} 
      />

      {/* Product Details */}
      <div className="grid md:grid-cols-2 gap-8 mb-16">
        <ProductImages images={product.images} />
        <ProductInfo product={product} />
      </div>

      {/* Product Tabs */}
      <ProductTabs product={product} />

      {/* Related Products */}
      <section className="mt-16">
        <h2 className="text-2xl font-bold mb-8">You Might Also Like</h2>
        <RelatedProducts categoryId={product.categoryId} />
      </section>
    </div>
  )
}