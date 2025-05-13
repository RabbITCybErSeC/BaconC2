import { Navigate, Outlet } from 'react-router-dom';
import { isAuthenticated } from './authService.tsx';

const ProtectedRoute = () => {
  return isAuthenticated() ? <Outlet /> : <Navigate to="/login" replace />;
};

export default ProtectedRoute;
