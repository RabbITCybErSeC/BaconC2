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
          {/* Redirect base path "/" to "/playbooks" */}
          <Route index element={<Navigate to="/playbooks" replace />} />

          {/* Playbook Inventory Route */}
          <Route path="playbooks" element={<AgentInventoryPage />} />

          {/* Catch-all for unknown routes within the layout, redirects to playbooks */}
          <Route path="*" element={<Navigate to="/playbooks" replace />} />
        </Route>
      </Route>
    </Routes>
  );
}

export default App;
