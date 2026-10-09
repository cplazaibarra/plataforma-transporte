import './globals.css';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'LastMile Delivery Platform | Monitoreo y POD',
  description: 'Plataforma integral de gestión de entregas de última milla con backend en Go, monitoreo GPS y prueba de entrega digital.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="es">
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
      </head>
      <body className="min-h-screen bg-slate-50">{children}</body>
    </html>
  );
}
