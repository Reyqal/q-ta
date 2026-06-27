import { Routes, Route } from 'react-router-dom'
import { ProtectedRoute } from './components/ProtectedRoute'
import { AdminRoute } from './components/AdminRoute'
import { AdminLayout } from './layouts/AdminLayout'
import { TenantLayout } from './layouts/TenantLayout'
import { LandingPage } from './pages/LandingPage'
import { LoginPage } from './pages/LoginPage'
import { DashboardPage } from './pages/admin/DashboardPage'
import { RoomsPage } from './pages/admin/RoomsPage'
import { TenantsPage } from './pages/admin/TenantsPage'
import { InvoicesPage as AdminInvoicesPage } from './pages/admin/InvoicesPage'
import { NotificationsPage } from './pages/admin/NotificationsPage'
import { InvoicesPage as TenantInvoicesPage } from './pages/tenant/InvoicesPage'
import { PaymentPage } from './pages/tenant/PaymentPage'

function App() {
  return (
    <Routes>
      {/* Public Routes */}
      <Route path="/" element={<LandingPage />} />
      <Route path="/login" element={<LoginPage />} />

      {/* Admin Routes */}
      <Route element={<ProtectedRoute />}>
        <Route element={<AdminRoute />}>
          <Route element={<AdminLayout />}>
            <Route path="/admin/dashboard" element={<DashboardPage />} />
            <Route path="/admin/rooms" element={<RoomsPage />} />
            <Route path="/admin/tenants" element={<TenantsPage />} />
            <Route path="/admin/invoices" element={<AdminInvoicesPage />} />
            <Route path="/admin/notifications" element={<NotificationsPage />} />
          </Route>
        </Route>

        {/* Tenant Routes */}
        <Route element={<TenantLayout />}>
          <Route path="/tenant/invoices" element={<TenantInvoicesPage />} />
          <Route path="/tenant/payment/:invoiceId" element={<PaymentPage />} />
        </Route>
      </Route>
    </Routes>
  )
}

export default App
