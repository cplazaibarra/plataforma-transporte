import { Order, Route, VehiclePosition, PublicTracking, DeliveryProof } from '../types';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8095/api/v1';

export async function getOrders(): Promise<Order[]> {
  const res = await fetch(`${API_BASE}/orders`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Error al cargar órdenes');
  return res.json();
}

export async function createOrder(order: Partial<Order>): Promise<Order> {
  const res = await fetch(`${API_BASE}/orders`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(order),
  });
  if (!res.ok) throw new Error('Error al crear orden');
  return res.json();
}

export async function getRoutes(): Promise<Route[]> {
  const res = await fetch(`${API_BASE}/routes`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Error al cargar rutas');
  return res.json();
}

export async function getRouteDetail(id: string): Promise<Route> {
  const res = await fetch(`${API_BASE}/routes/${id}`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Error al cargar detalle de ruta');
  return res.json();
}

export async function updateStopStatus(stopId: string, status: string): Promise<void> {
  const res = await fetch(`${API_BASE}/routes/stops/${stopId}/status`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status }),
  });
  if (!res.ok) throw new Error('Error al actualizar estado de la parada');
}

export async function submitPOD(data: {
  route_stop_id: string;
  recipient_name?: string;
  recipient_dni?: string;
  recipient_relationship?: string;
  signature_base64?: string;
  photo_url?: string;
  failure_reason?: string;
  notes?: string;
  latitude?: number;
  longitude?: number;
  status: 'delivered' | 'failed';
}): Promise<void> {
  const res = await fetch(`${API_BASE}/pod`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Error al enviar prueba de entrega');
}

export async function getVehiclePositions(): Promise<VehiclePosition[]> {
  const res = await fetch(`${API_BASE}/telemetry/positions`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Error al cargar posiciones de vehículos');
  return res.json();
}

export async function sendTelemetry(telemetry: {
  vehicle_id: string;
  driver_id?: string;
  latitude: number;
  longitude: number;
  speed: number;
  heading: number;
  battery_level?: number;
}): Promise<void> {
  const res = await fetch(`${API_BASE}/telemetry`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(telemetry),
  });
  if (!res.ok) throw new Error('Error al enviar telemetría');
}

export async function getPublicTracking(trackingNumber: string): Promise<PublicTracking> {
  const res = await fetch(`${API_BASE}/tracking/${trackingNumber}`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Número de seguimiento no encontrado');
  return res.json();
}

export async function getUsers(): Promise<User[]> {
  const res = await fetch(`${API_BASE}/users`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Error al cargar usuarios');
  return res.json();
}

export async function createUser(user: Partial<User>): Promise<User> {
  const res = await fetch(`${API_BASE}/users`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(user),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Error al crear usuario');
  }
  return res.json();
}

export async function deleteUser(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/users/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error('Error al eliminar usuario');
}

export async function getVehicles(): Promise<Vehicle[]> {
  const res = await fetch(`${API_BASE}/vehicles`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Error al cargar vehículos');
  return res.json();
}

