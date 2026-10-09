'use client';

import React, { useState, useEffect } from 'react';
import dynamic from 'next/dynamic';
import { useParams } from 'next/navigation';
import { 
  Package, Truck, CheckCircle2, Clock, MapPin, 
  ShieldCheck, AlertCircle, FileText 
} from 'lucide-react';
import { PublicTracking, VehiclePosition, RouteStop } from '../../../types';
import { getPublicTracking } from '../../../lib/api';

const DynamicMap = dynamic(() => import('../../../components/Map'), {
  ssr: false,
  loading: () => (
    <div className="w-full h-72 bg-slate-100 animate-pulse rounded-2xl flex items-center justify-center text-slate-400 text-xs font-semibold">
      Cargando mapa en vivo...
    </div>
  ),
});

export default function CustomerTrackingPage() {
  const params = useParams();
  const trackingNumber = params?.id as string;
  const [data, setData] = useState<PublicTracking | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchTracking = async () => {
    try {
      if (!trackingNumber) return;
      const res = await getPublicTracking(trackingNumber);
      setData(res);
    } catch (err: any) {
      setError(err.message || 'No se encontró la guía de despacho');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTracking();

    // Polling cada 6 segundos para actualizar la posición del móvil del chofer
    const interval = setInterval(fetchTracking, 6000);
    return () => clearInterval(interval);
  }, [trackingNumber]);

  if (loading) {
    return (
      <div className="min-h-screen bg-slate-100 flex items-center justify-center p-4">
        <div className="text-center space-y-3">
          <div className="w-12 h-12 border-4 border-sky-600 border-t-transparent rounded-full animate-spin mx-auto" />
          <p className="text-slate-600 font-semibold text-sm">Consultando estado del envío...</p>
        </div>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="min-h-screen bg-slate-100 flex items-center justify-center p-4">
        <div className="bg-white p-6 rounded-2xl shadow-xl max-w-sm w-full text-center space-y-3 border border-slate-200">
          <AlertCircle className="w-12 h-12 text-rose-500 mx-auto" />
          <h2 className="font-bold text-slate-900 text-lg">Guía No Encontrada</h2>
          <p className="text-xs text-slate-500">
            No pudimos encontrar información para el código <span className="font-bold text-slate-700">{trackingNumber}</span>.
          </p>
        </div>
      </div>
    );
  }

  const { order, stop_status, driver_name, vehicle_plate, vehicle_lat, vehicle_lng, proof } = data;

  const isDelivered = order.status === 'delivered';
  const isEnRoute = order.status === 'in_transit' || stop_status === 'en_route';
  const isFailed = order.status === 'failed';

  // Preparar vehículos y paradas para el mapa
  const vehiclesForMap: VehiclePosition[] = vehicle_lat && vehicle_lng ? [{
    vehicle_id: 'v1',
    plate: vehicle_plate || 'MÓVIL',
    model: 'Furgón de Reparto',
    driver_name: driver_name || 'Conductor asignado',
    latitude: vehicle_lat,
    longitude: vehicle_lng,
    speed: 32,
    heading: 0,
    updated_at: new Date().toISOString(),
  }] : [];

  const stopsForMap: RouteStop[] = [{
    id: 's1',
    route_id: 'r1',
    order_id: order.id,
    stop_number: data.stop_number || 1,
    status: stop_status as any,
    notes: '',
    order: order,
    createdAt: new Date(),
  }];

  return (
    <div className="min-h-screen bg-slate-100 text-slate-900 pb-12">
      
      {/* Barra superior de marca */}
      <div className="bg-slate-900 text-white py-4 px-6 shadow-sm">
        <div className="max-w-xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Package className="w-5 h-5 text-sky-400" />
            <span className="font-extrabold tracking-tight">LastMile Tracking</span>
          </div>
          <span className="text-xs text-slate-400 font-mono">#{order.tracking_number}</span>
        </div>
      </div>

      <div className="max-w-xl mx-auto px-4 mt-6 space-y-4">
        
        {/* Tarjeta de Estado del Envío */}
        <div className="bg-white rounded-2xl p-6 shadow-sm border border-slate-200/80">
          <div className="flex items-center justify-between mb-4">
            <div>
              <p className="text-xs font-semibold uppercase text-slate-400">Hola, {order.customer_name}</p>
              <h2 className="text-xl font-black text-slate-900 tracking-tight mt-0.5">
                {isDelivered 
                  ? '¡Tu pedido ha sido entregado!' 
                  : isEnRoute 
                  ? 'Tu pedido está en camino' 
                  : isFailed
                  ? 'No pudimos concretar la entrega'
                  : 'Estamos preparando tu despacho'}
              </h2>
            </div>
            <div className={`p-3 rounded-2xl ${
              isDelivered ? 'bg-emerald-50 text-emerald-600' : isEnRoute ? 'bg-sky-50 text-sky-600 animate-pulse' : 'bg-slate-100 text-slate-600'
            }`}>
              {isDelivered ? <CheckCircle2 className="w-8 h-8" /> : <Truck className="w-8 h-8" />}
            </div>
          </div>

          {/* Stepper de Progreso */}
          <div className="relative flex items-center justify-between pt-2 pb-4">
            <div className="absolute left-0 top-1/2 -translate-y-1/2 h-1 bg-slate-100 w-full z-0" />
            <div 
              className="absolute left-0 top-1/2 -translate-y-1/2 h-1 bg-sky-600 z-0 transition-all duration-500"
              style={{ width: isDelivered ? '100%' : isEnRoute ? '66%' : '33%' }}
            />

            <div className="relative z-10 flex flex-col items-center">
              <div className="w-8 h-8 rounded-full bg-sky-600 text-white text-xs font-bold flex items-center justify-center shadow">
                ✓
              </div>
              <span className="text-[11px] font-semibold text-slate-600 mt-1.5">Bodega</span>
            </div>

            <div className="relative z-10 flex flex-col items-center">
              <div className={`w-8 h-8 rounded-full text-xs font-bold flex items-center justify-center shadow ${
                isEnRoute || isDelivered ? 'bg-sky-600 text-white' : 'bg-white border-2 border-slate-300 text-slate-400'
              }`}>
                2
              </div>
              <span className="text-[11px] font-semibold text-slate-600 mt-1.5">En Reparto</span>
            </div>

            <div className="relative z-10 flex flex-col items-center">
              <div className={`w-8 h-8 rounded-full text-xs font-bold flex items-center justify-center shadow ${
                isDelivered ? 'bg-emerald-600 text-white' : 'bg-white border-2 border-slate-300 text-slate-400'
              }`}>
                3
              </div>
              <span className="text-[11px] font-semibold text-slate-600 mt-1.5">Entregado</span>
            </div>
          </div>

          {/* Datos del conductor y móvil */}
          {driver_name && (
            <div className="mt-2 bg-slate-50 rounded-xl p-3.5 flex items-center justify-between text-xs border border-slate-100">
              <div className="flex items-center gap-3">
                <div className="w-9 h-9 rounded-full bg-slate-200 flex items-center justify-center font-bold text-slate-700">
                  {driver_name.charAt(0)}
                </div>
                <div>
                  <p className="font-bold text-slate-800">{driver_name}</p>
                  <p className="text-slate-500 text-[11px]">Móvil {vehicle_plate}</p>
                </div>
              </div>
              <span className="bg-sky-100 text-sky-800 font-bold px-2.5 py-1 rounded-lg text-[11px]">
                Ventana {order.window_start} - {order.window_end}
              </span>
            </div>
          )}
        </div>

        {/* Mapa en Vivo con OpenStreetMap y Leaflet */}
        <div className="bg-white rounded-2xl p-4 shadow-sm border border-slate-200/80">
          <div className="flex items-center justify-between mb-3">
            <h3 className="font-bold text-slate-900 text-sm flex items-center gap-1.5">
              <MapPin className="w-4 h-4 text-sky-600" /> Ubicación en Tiempo Real
            </h3>
            <span className="text-[11px] text-slate-400">Actualización automática</span>
          </div>

          <div className="h-72 w-full rounded-xl overflow-hidden border border-slate-200">
            <DynamicMap
              vehicles={vehiclesForMap}
              stops={stopsForMap}
              center={[order.latitude, order.longitude]}
              zoom={14}
            />
          </div>

          <div className="mt-3 flex items-start gap-2 text-xs text-slate-600 bg-slate-50 p-2.5 rounded-xl">
            <MapPin className="w-4 h-4 text-slate-400 shrink-0 mt-0.5" />
            <p><strong>Destino:</strong> {order.delivery_address}, {order.commune}</p>
          </div>
        </div>

        {/* Comprobante de Entrega Digital (Si ya fue entregado) */}
        {proof && (
          <div className="bg-white rounded-2xl p-6 shadow-sm border border-emerald-200 space-y-4">
            <div className="flex items-center gap-2 text-emerald-800">
              <ShieldCheck className="w-5 h-5 text-emerald-600" />
              <h3 className="font-bold text-sm">Comprobante de Entrega Digital</h3>
            </div>

            <div className="grid grid-cols-2 gap-3 text-xs bg-emerald-50/50 p-3.5 rounded-xl border border-emerald-100">
              <div>
                <p className="text-slate-400 font-semibold">RECIBIDO POR</p>
                <p className="font-bold text-slate-800 mt-0.5">{proof.recipient_name}</p>
              </div>
              <div>
                <p className="text-slate-400 font-semibold">RUT / DNI</p>
                <p className="font-bold text-slate-800 mt-0.5">{proof.recipient_dni}</p>
              </div>
              <div>
                <p className="text-slate-400 font-semibold">VÍNCULO</p>
                <p className="font-bold text-slate-800 mt-0.5 capitalize">{proof.recipient_relationship}</p>
              </div>
              <div>
                <p className="text-slate-400 font-semibold">FECHA Y HORA</p>
                <p className="font-bold text-slate-800 mt-0.5">
                  {new Date(proof.submitted_at).toLocaleTimeString('es-CL', { hour: '2-digit', minute: '2-digit' })} hrs
                </p>
              </div>
            </div>

            {proof.signature_base64 && (
              <div>
                <p className="text-xs font-bold text-slate-700 mb-1">Firma Registrada:</p>
                <div className="border border-slate-200 rounded-xl p-2 bg-slate-50 flex items-center justify-center">
                  <img src={proof.signature_base64} alt="Firma del receptor" className="max-h-24 object-contain" />
                </div>
              </div>
            )}
          </div>
        )}

      </div>
    </div>
  );
}
