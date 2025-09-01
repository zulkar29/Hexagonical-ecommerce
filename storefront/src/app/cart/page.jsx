// Shopping Cart Page
import CartItems from '@/components/cart/CartItems'
import CartSummary from '@/components/cart/CartSummary'
import EmptyCart from '@/components/cart/EmptyCart'
import { useAtom } from 'jotai'
import { cartItemsAtom } from '@/stores/cartStore'

export const metadata = {
  title: 'Shopping Cart',
  description: 'Review your items and proceed to checkout.'
}

export default function CartPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Shopping Cart</h1>
      
      <CartPageContent />
    </div>
  )
}

function CartPageContent() {
  // TODO: This should be a client component
  // const [cartItems] = useAtom(cartItemsAtom)
  const cartItems = [] // Placeholder

  if (cartItems.length === 0) {
    return <EmptyCart />
  }

  return (
    <div className="lg:grid lg:grid-cols-12 lg:gap-12">
      {/* Cart Items */}
      <div className="lg:col-span-8">
        <CartItems items={cartItems} />
      </div>

      {/* Order Summary */}
      <div className="lg:col-span-4 mt-8 lg:mt-0">
        <CartSummary />
      </div>
    </div>
  )
}