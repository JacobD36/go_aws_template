// Header component

'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { clearAuthData } from '@/lib/auth';
import { Button } from '@/components/ui/Button';

export const Header: React.FC = () => {
  const router = useRouter();

  const handleLogout = () => {
    clearAuthData();
    router.push('/login');
  };

  return (
    <header className="bg-white border-b border-gray-200 px-6 py-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-semibold text-gray-800">
            Bienvenido
          </h2>
          <p className="text-sm text-gray-600">
            Gestiona tu sistema desde aquí
          </p>
        </div>
        
        <div className="flex items-center space-x-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={handleLogout}
          >
            <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Cerrar Sesión
          </Button>
        </div>
      </div>
    </header>
  );
};
