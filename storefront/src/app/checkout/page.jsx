// Checkout Page
import CheckoutForm from '@/components/checkout/CheckoutForm'
import OrderSummary from '@/components/checkout/OrderSummary'
import CheckoutSteps from '@/components/checkout/CheckoutSteps'

export const metadata = {
  title: 'Checkout',
  description: 'Complete your purchase securely.'
}

export default function CheckoutPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Checkout</h1>
      
      {/* Progress Steps */}
      <CheckoutSteps currentStep={1} />

      <div className="lg:grid lg:grid-cols-12 lg:gap-12 mt-8">
        {/* Checkout Form */}
        <div className="lg:col-span-7">
          <CheckoutForm />
        </div>

        {/* Order Summary */}
        <div className="lg:col-span-5 mt-8 lg:mt-0">
          <div className="bg-gray-50 rounded-lg p-6">
            <OrderSummary />
          </div>
        </div>
      </div>
    </div>
  )
}