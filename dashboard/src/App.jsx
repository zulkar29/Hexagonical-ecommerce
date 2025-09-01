// TODO: Implement main App component during development

import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
// TODO: Add MUI theme provider
// import { ThemeProvider } from '@mui/material/styles'
// import CssBaseline from '@mui/material/CssBaseline'

// TODO: Import pages
// import Dashboard from './pages/Dashboard'
// import Products from './pages/Products'
// import Orders from './pages/Orders'
// import Customers from './pages/Customers'
// import Settings from './pages/Settings'
// import Login from './pages/Login'

function App() {
  // TODO: Implement during development:
  // - Add authentication check
  // - Add Jotai provider
  // - Add React Query provider
  // - Add MUI theme
  // - Add protected routes
  // - Add error boundary

  return (
    <div className="App">
      {/* TODO: Add providers */}
      {/* <ThemeProvider theme={theme}> */}
      {/* <CssBaseline /> */}
      {/* <QueryClient> */}
      {/* <JotaiProvider> */}
        <Router>
          <Routes>
            {/* TODO: Add actual routes */}
            <Route path="/" element={<div>TODO: Dashboard</div>} />
            <Route path="/products" element={<div>TODO: Products</div>} />
            <Route path="/orders" element={<div>TODO: Orders</div>} />
            <Route path="/customers" element={<div>TODO: Customers</div>} />
            <Route path="/settings" element={<div>TODO: Settings</div>} />
            <Route path="/login" element={<div>TODO: Login</div>} />
          </Routes>
        </Router>
      {/* </JotaiProvider> */}
      {/* </QueryClient> */}
      {/* </ThemeProvider> */}
    </div>
  )
}

export default App