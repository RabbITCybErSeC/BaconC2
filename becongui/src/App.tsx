import { Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import DashboardLayout from './layouts/DashboardLayout';
import AgentInventoryPage from './pages/AgentInventoryPage.tsx';
import ProtectedRoute from './services/routedProtector.tsx';

function App() {
  return (
    <Routes>
      {/* Public route for login */}
      <Route path="/login" element={<LoginPage />} />

      {/* Protected routes under DashboardLayout */}
      <Route element={<ProtectedRoute />}>
        <Route path="/" element={<DashboardLayout />}>
          {/* Redirect base path "/" to "/agents" */}
          <Route index element={<Navigate to="/agents" replace />} />

          {/* Agent Inventory Route */}
          <Route path="agents" element={<AgentInventoryPage />} />

          {/* Catch-all for unknown routes within the layout, redirects to agents */}
          <Route path="*" element={<Navigate to="/agents" replace />} />
        </Route>
      </Route> 
    </Routes >
  );
}

export default App;
