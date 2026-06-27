import { Link } from 'react-router-dom';
import { Building2, LogIn } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

export function Navbar() {
  const { isAuthenticated, isAdmin } = useAuth();

  return (
    <nav className="fixed top-0 left-0 right-0 z-40 glass border-b border-white/5">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="flex items-center gap-3 group">
            <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-emerald-500 to-blue-500 flex items-center justify-center group-hover:scale-110 transition-transform">
              <Building2 className="w-5 h-5 text-white" />
            </div>
            <span className="text-xl font-bold text-white tracking-tight">
              Q-TA
            </span>
          </Link>

          <div className="flex items-center gap-3">
            {isAuthenticated ? (
              <Link
                to={isAdmin ? '/admin' : '/tenant'}
                className="flex items-center gap-2 px-4 py-2 rounded-xl bg-primary/10 text-primary hover:bg-primary/20 transition-all font-medium text-sm"
              >
                Dashboard
              </Link>
            ) : (
              <Link
                to="/login"
                className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 text-white hover:from-emerald-400 hover:to-emerald-500 transition-all font-medium text-sm shadow-lg shadow-emerald-500/20 hover:shadow-emerald-500/30"
              >
                <LogIn className="w-4 h-4" />
                Masuk
              </Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
