export interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'operator' | 'driver';
  vehicle_id?: string;
  vehicle_plate?: string;
  vehicle_model?: string;
  phone: string;
  status: string;
  created_at: string;
}

export interface Driver {
  id: string;
  name: string;
  email: string;
  phone: string;
  license_number: string;
  status: string;
}

export interface Vehicle {
  id: string;
  plate: string;
  model: string;
  brand: string;
  capacity_kg: number;
  status: string;
}

export interface Order {
  id: string;
  tracking_number: string;
  external_order_id: string;
  customer_name: string;
  customer_phone: string;
  customer_email: string;
  delivery_address: string;
  commune: string;
  latitude: number;
  longitude: number;
  window_start: string;
  window_end: string;
  packages_count: number;
  items_description: string;
  status: 'ready_for_dispatch' | 'assigned' | 'in_transit' | 'delivered' | 'failed';
  created_at: string;
}

export interface DeliveryProof {
  id: string;
  route_stop_id: string;
  recipient_name: string;
  recipient_dni: string;
  recipient_relationship: string;
  signature_base64: string;
  photo_url: string;
  failure_reason?: string;
  notes: string;
  latitude?: number;
  longitude?: number;
  submitted_at: string;
}

export interface RouteStop {
  id: string;
  route_id: string;
  order_id: string;
  stop_number: number;
  status: 'pending' | 'en_route' | 'arrived' | 'delivered' | 'failed';
  eta?: string;
  arrived_at?: string;
  completed_at?: string;
  notes: string;
  order?: Order;
  proof?: DeliveryProof;
}

export interface Route {
  id: string;
  code: string;
  date: string;
  driver_id?: string;
  driver_name?: string;
  vehicle_id?: string;
  vehicle_plate?: string;
  status: 'draft' | 'assigned' | 'in_progress' | 'completed' | 'cancelled';
  started_at?: string;
  completed_at?: string;
  total_stops: number;
  completed_stops: number;
  stops?: RouteStop[];
}

export interface VehiclePosition {
  vehicle_id: string;
  plate: string;
  model: string;
  driver_id?: string;
  driver_name?: string;
  latitude: number;
  longitude: number;
  speed: number;
  heading: number;
  route_id?: string;
  route_code?: string;
  current_stop?: number;
  total_stops?: number;
  updated_at: string;
}

export interface PublicTracking {
  order: Order;
  stop_status: string;
  stop_number: number;
  total_stops: number;
  driver_name?: string;
  vehicle_plate?: string;
  vehicle_lat?: number;
  vehicle_lng?: number;
  proof?: DeliveryProof;
}
