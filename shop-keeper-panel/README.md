# Shop Owner Dashboard - Business Operating System

## Project Overview

A comprehensive business dashboard for Bangladesh small businesses featuring inventory management, POS system, customer management, and sales analytics. Built with React and modern web technologies for optimal performance and user experience.

## 🎯 Core Features

### 1. Dashboard Overview
- Real-time business metrics and KPIs
- Sales, revenue, and inventory widgets
- Quick action buttons
- Mobile-responsive design

### 2. Inventory Management
- Product catalog with SKU generation
- Stock tracking and alerts
- Barcode scanning support
- Supplier management
- Low stock notifications

### 3. Point of Sale (POS)
- Quick product search
- Multiple payment methods (Cash, bKash, Nagad)
- Receipt generation
- Customer selection
- Discount management

### 4. Customer Management
- Customer database
- Purchase history tracking
- Loyalty program
- Communication history

### 5. Sales & Orders
- Order management
- Invoice generation
- Payment tracking
- Sales reporting

### 6. Analytics & Reports
- Sales performance analytics
- Inventory reports
- Customer insights
- Export capabilities (PDF, Excel)

## 🛠️ Tech Stack

- **React 18+** - Frontend framework
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Styling framework
- **shadcn/ui** - UI component library
- **Jotai** - State management (minimal usage)
- **TanStack Query** - Server state management
- **React Hook Form** - Form handling
- **React Router v6** - Routing
- **Lucide React** - Icons
- **Recharts** - Data visualization
- **Axios** - HTTP client

## 📁 Project Structure

```
src/
├── components/
│   ├── ui/              # shadcn/ui components
│   ├── layout/          # Layout components (Sidebar, Header)
│   ├── forms/           # Reusable form components
│   └── charts/          # Chart components
├── pages/
│   ├── Dashboard.jsx    # Main dashboard
│   ├── Inventory/       # Inventory pages
│   ├── POS.jsx          # Point of sale
│   ├── Customers.jsx    # Customer management
│   ├── Sales.jsx        # Sales & orders
│   └── Reports.jsx      # Analytics & reports
├── services/
│   └── api.js           # API service functions
├── hooks/
│   └── useAuth.js       # Authentication hook
├── utils/
│   ├── constants.js     # App constants
│   └── helpers.js       # Utility functions
└── lib/
    ├── axios.js         # Axios config
    └── utils.js         # shadcn/ui utils
```

## 🎨 Design System

### Colors (Bangladesh Theme)
```css
--primary: hsl(142, 76%, 36%)        /* Bangladesh Green */
--secondary: hsl(210, 40%, 96%)      /* Light Gray */
--accent: hsl(45, 93%, 47%)          /* Golden Yellow */
--success: hsl(142, 76%, 36%)        /* Success */
--warning: hsl(38, 92%, 50%)         /* Warning */
--destructive: hsl(0, 72%, 51%)      /* Error */
```

### Typography
- **Font**: Inter (English), Noto Sans Bengali (Bengali)
- **Sizes**: text-sm (14px), text-base (16px), text-lg (18px), text-xl (20px)

## 🏗️ Component Patterns

### Layout Structure
```jsx
<div className="min-h-screen bg-background">
  <Sidebar />
  <div className="ml-64">
    <Header />
    <main className="p-6">
      <Outlet />
    </main>
  </div>
</div>
```

### Page Layout
```jsx
<div className="space-y-6">
  <div className="flex items-center justify-between">
    <h1 className="text-2xl font-bold">Page Title</h1>
    <Button>Primary Action</Button>
  </div>
  <Card>
    <CardContent>
      {/* Content */}
    </CardContent>
  </Card>
</div>
```

### Data Table
```jsx
<Card>
  <CardHeader>
    <div className="flex items-center justify-between">
      <CardTitle>Products</CardTitle>
      <Input placeholder="Search..." className="w-64" />
    </div>
  </CardHeader>
  <CardContent>
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Stock</TableHead>
          <TableHead>Price</TableHead>
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {products.map((product) => (
          <TableRow key={product.id}>
            <TableCell>{product.name}</TableCell>
            <TableCell>
              <Badge variant={product.stock < 10 ? "destructive" : "default"}>
                {product.stock}
              </Badge>
            </TableCell>
            <TableCell>৳{product.price}</TableCell>
            <TableCell>
              <Button variant="ghost" size="sm">Edit</Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  </CardContent>
</Card>
```

### Stats Cards
```jsx
<div className="grid grid-cols-1 md:grid-cols-4 gap-6">
  {stats.map((stat) => (
    <Card key={stat.id}>
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm text-muted-foreground">{stat.label}</p>
            <p className="text-2xl font-bold">{stat.value}</p>
          </div>
          <stat.icon className="h-8 w-8 text-primary" />
        </div>
      </CardContent>
    </Card>
  ))}
</div>
```

## 🔌 API Integration

### API Service
```javascript
// services/api.js
import axios from 'axios'

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:3001/api',
})

// Products
export const productsApi = {
  getAll: () => api.get('/products'),
  create: (data) => api.post('/products', data),
  update: (id, data) => api.put(`/products/${id}`, data),
  delete: (id) => api.delete(`/products/${id}`),
}

// Customers
export const customersApi = {
  getAll: () => api.get('/customers'),
  create: (data) => api.post('/customers', data),
  update: (id, data) => api.put(`/customers/${id}`, data),
}

// Sales
export const salesApi = {
  getAll: () => api.get('/sales'),
  create: (data) => api.post('/sales', data),
  getStats: () => api.get('/sales/stats'),
}
```

### TanStack Query Usage
```javascript
// In components - direct usage
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { productsApi } from '../services/api'

const ProductList = () => {
  const queryClient = useQueryClient()
  
  // Fetch products
  const { data: products, isLoading } = useQuery({
    queryKey: ['products'],
    queryFn: () => productsApi.getAll().then(res => res.data),
  })

  // Create product mutation
  const createProduct = useMutation({
    mutationFn: productsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries(['products'])
      toast.success('Product created!')
    },
  })

  // Delete product mutation
  const deleteProduct = useMutation({
    mutationFn: productsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries(['products'])
      toast.success('Product deleted!')
    },
  })

  if (isLoading) return <div>Loading...</div>

  return (
    <div>
      {/* Product list UI */}
    </div>
  )
}
```

## 🗂️ State Management

### Minimal Jotai Usage
```javascript
// Only for UI state that needs global access
import { atom } from 'jotai'

// UI state atoms
export const sidebarOpenAtom = atom(false)
export const themeAtom = atom('light')

// Usage in components
const [sidebarOpen, setSidebarOpen] = useAtom(sidebarOpenAtom)
```

### Local State for Most Cases
```javascript
// Use React useState for component-specific state
const [selectedProducts, setSelectedProducts] = useState([])
const [filters, setFilters] = useState({ search: '', category: '' })
```

## 📱 Responsive Design

### Mobile-First Grid
```jsx
<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
  {/* Auto-responsive cards */}
</div>
```

### Mobile Navigation
```jsx
<Sheet>
  <SheetTrigger asChild>
    <Button variant="ghost" className="lg:hidden">
      <Menu className="h-6 w-6" />
    </Button>
  </SheetTrigger>
  <SheetContent side="left">
    <nav className="space-y-4">
      {/* Mobile menu items */}
    </nav>
  </SheetContent>
</Sheet>
```

## 📊 Charts & Data Visualization

### Simple Chart Component
```jsx
import { LineChart, Line, XAxis, YAxis, ResponsiveContainer } from 'recharts'

const SalesChart = ({ data }) => (
  <ResponsiveContainer width="100%" height={300}>
    <LineChart data={data}>
      <XAxis dataKey="date" />
      <YAxis />
      <Line type="monotone" dataKey="sales" stroke="#8884d8" />
    </LineChart>
  </ResponsiveContainer>
)
```

## 🎯 Key Implementation Notes

### Keep It Simple
- Use React's built-in state management for most cases
- Only use Jotai for truly global UI state
- Direct TanStack Query usage in components
- Avoid over-abstraction with custom hooks

### Performance Tips
- Use React.memo for expensive components
- Implement pagination for large data sets
- Use debounced search inputs
- Lazy load route components

### Form Handling
```jsx
import { useForm } from 'react-hook-form'

const ProductForm = ({ onSubmit }) => {
  const { register, handleSubmit, formState: { errors } } = useForm()

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="name">Product Name</Label>
        <Input 
          id="name"
          {...register('name', { required: 'Name is required' })}
        />
        {errors.name && <p className="text-red-500 text-sm">{errors.name.message}</p>}
      </div>
      
      <div>
        <Label htmlFor="price">Price</Label>
        <Input 
          id="price"
          type="number"
          {...register('price', { required: 'Price is required' })}
        />
        {errors.price && <p className="text-red-500 text-sm">{errors.price.message}</p>}
      </div>

      <Button type="submit">Save Product</Button>
    </form>
  )
}
```

### Error Handling
```jsx
// Simple error boundary or use react-error-boundary
const { data, error, isLoading } = useQuery({
  queryKey: ['products'],
  queryFn: productsApi.getAll,
})

if (error) return <div>Error: {error.message}</div>
if (isLoading) return <div>Loading...</div>
```
