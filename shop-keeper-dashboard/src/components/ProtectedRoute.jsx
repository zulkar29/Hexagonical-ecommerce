// ProtectedRoute component - Authentication disabled for development
const ProtectedRoute = ({ children }) => {
  // Return children directly without authentication checks
  return children;
};

export default ProtectedRoute;