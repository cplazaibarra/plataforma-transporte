'use client';

import React, { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import { 
  Truck, Navigation, CheckCircle2, Clock, MapPin, 
  Phone, Package, Radio, ArrowLeft, ChevronRight, ShieldCheck 
} from 'lucide-react';
import { Route, RouteStop } from '../../types';
import { getRoutes, updateStopStatus, submitPOD, sendTelemetry } from '../../lib/api';
import PODModal from '../../components/PODModal';

export default function DriverAppPage() {
  const [route, setRoute] = useState<Route | null>(null);
  const [selectedStop, setSelectedStop] = useState<RouteStop | null>(null);
  const [isPodModalOpen, setIsPodModalOpen] = useState(false);
  const [isGpsStreaming, setIsGpsStreaming] = useState(false);
  const [gpsCoords, setGpsCoords] = useState<{ lat: number; lng: number } | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const watchIdRef = useRef<number | null>(null);

  const loadDriverRoute = async () => {
    try {
      setIsLoading(true);
      const routes = await getRoutes();
      if (routes && routes.length > 0) {
        // En una app completa, se filtra por el ID del chofer logueado
        setRoute(routes[0]);
      }
    } catch (err) {
      console.error('Error cargando ruta del chofer:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadDriverRoute();
  }, []);

  // Control de transmisión de GPS en vivo desde el navegador del conductor
  const toggleGpsStreaming = () => {
    if (isGpsStreaming) {
      if (watchIdRef.current !== null) {
        navigator.geolocation.clearWatch(watchIdRef.current);
        watchIdRef.current = null;
      }
      setIsGpsStreaming(false);
    } else {
      if (!('geolocation' in navigator)) {
        alert('Geolocalización no soportada en este dispositivo.');
        return;
      }

      setIsGpsStreaming(true);
      watchIdRef.current = navigator.geolocation.watchPosition(
        (pos) => {
          const lat = pos.coords.latitude;
          const lng = pos.coords.longitude;
          const speed = pos.coords.speed ? pos.coords.speed * 3.6 : 30; // km/h
          const heading = pos.coords.heading || 0;

          setGpsCoords({ lat, lng });

          // Enviar telemetría al backend en Go
          if (route?.vehicle_id) {
            sendTelemetry({
              vehicle_id: route.vehicle_id,
              driver_id: route.driver_id,
              latitude: lat,
              longitude: lng,
              speed,
              heading,
            }).catch(console.error);
          }
        },
        (err) => {
          console.warn('Error de GPS:', err.message);
        },
        { enableHighAccuracy: true, maximumAge: 5000, timeout: 10000 }
      );
    }
  };

  const handleStartEnRoute = async (stopId: string) => {
    try {
      await updateStopStatus(stopId, 'en_route');
      loadDriverRoute();
    } catch (err: any) {
      alert('Error: ' + err.message);
    }
  };

  const handleSubmitPOD = async (data: any) => {
    await submitPOD({
      ...data,
      latitude: gpsCoords?.lat,
      longitude: gpsCoords?.lng,
    });
    loadDriverRoute();
  };

  const openWaze = (lat: number, lng: number) => {
    window.open(`https://waze.com/ul?ll=${lat},${lng}&navigate=yes`, '_blank');
  };

  const openGoogleMaps = (lat: number, lng: number) => {
    window.open(`https://www.google.com/maps/dir/?api=1&destination=${lat},${lng}`, '_blank');
  };

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100 flex flex-col max-w-md mx-auto shadow-2xl pb-10">
      
      {/* Cabecera Móvil */}
      <header className="bg-slate-950 p-4 border-b border-slate-800 sticky top-0 z-30 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link href="/" className="p-1.5 rounded-lg bg-slate-800 text-slate-400 hover:text-white">
            <ArrowLeft className="w-5 h-5" />
          </Link>
          <div>
            <h1 className="font-bold text-sm tracking-tight flex items-center gap-1.5">
              <Truck className="w-4 h-4 text-sky-400" /> Móvil {route?.vehicle_plate || 'KJ-DF-44'}
            </h1>
            <p className="text-xs text-slate-400">Chofer: {route?.driver_name || 'Carlos Morales'}</p>
          </div>
        </div>

        {/* Botón de Transmisión GPS */}
        <button
          onClick={toggleGpsStreaming}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold transition shadow-lg ${
            isGpsStreaming
              ? 'bg-emerald-500 text-slate-950 animate-pulse'
              : 'bg-slate-800 text-slate-300 border border-slate-700'
          }`}
        >
          <Radio className="w-3.5 h-3.5" />
          {isGpsStreaming ? 'GPS ACTIVO' : 'ACTIVAR GPS'}
        </button>
      </header>

      {/* Progreso de la Hoja de Ruta */}
      <div className="p-4 bg-slate-800/60 border-b border-slate-800">
        <div className="flex items-center justify-between text-xs mb-1.5">
          <span className="font-semibold text-slate-300">Ruta: {route?.code || 'RUTA-HOY'}</span>
          <span className="font-bold text-sky-400">
            {route?.completed_stops || 0} de {route?.total_stops || 0} Entregas
          </span>
        </div>
        <div className="w-full bg-slate-700 h-2.5 rounded-full overflow-hidden">
          <div 
            className="bg-gradient-to-r from-sky-500 to-emerald-400 h-full transition-all duration-500 rounded-full"
            style={{ 
              width: `${((route?.completed_stops || 0) / (route?.total_stops || 1)) * 100}%` 
            }}
          />
        </div>
      </div>

      {/* Lista de Paradas / Entregas */}
      <main className="flex-1 p-4 space-y-3">
        {route?.stops?.map((stop) => {
          const order = stop.order;
          if (!order) return null;

          const isDelivered = stop.status === 'delivered';
          const isFailed = stop.status === 'failed';
          const isEnRoute = stop.status === 'en_route';
          const isPending = stop.status === 'pending';

          return (
            <div 
              key={stop.id}
              className={`rounded-2xl p-4 border transition-all ${
                isDelivered
                  ? 'bg-slate-800/40 border-emerald-900/50 opacity-75'
                  : isEnRoute
                  ? 'bg-slate-800 border-sky-500 shadow-lg shadow-sky-950/50 ring-1 ring-sky-500'
                  : 'bg-slate-800/80 border-slate-700'
              }`}
            >
              {/* Header de la Parada */}
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2">
                  <span className={`w-7 h-7 rounded-full text-xs font-bold flex items-center justify-center ${
                    isDelivered ? 'bg-emerald-500 text-slate-950' : isEnRoute ? 'bg-sky-500 text-white animate-bounce' : 'bg-slate-700 text-slate-300'
                  }`}>
                    {stop.stop_number}
                  </span>
                  <div>
                    <h3 className="font-bold text-sm text-white">{order.customer_name}</h3>
                    <p className="text-[11px] text-slate-400">Guía: {order.tracking_number}</p>
                  </div>
                </div>

                <span className={`text-[10px] font-bold uppercase px-2 py-0.5 rounded-full ${
                  isDelivered 
                    ? 'bg-emerald-950 text-emerald-400 border border-emerald-800'
                    : isEnRoute
                    ? 'bg-sky-950 text-sky-400 border border-sky-700 animate-pulse'
                    : isFailed
                    ? 'bg-rose-950 text-rose-400 border border-rose-800'
                    : 'bg-slate-700 text-slate-400'
                }`}>
                  {isDelivered ? 'Entregado' : isEnRoute ? 'En Camino' : isFailed ? 'Fallido' : 'Pendiente'}
                </span>
              </div>

              {/* Dirección */}
              <div className="flex items-start gap-2 text-xs text-slate-300 mb-3 bg-slate-900/50 p-2.5 rounded-xl">
                <MapPin className="w-4 h-4 text-sky-400 shrink-0 mt-0.5" />
                <div>
                  <p className="font-semibold text-slate-200">{order.delivery_address}</p>
                  <p className="text-slate-400 text-[11px]">{order.commune} • {order.items_description || '1 bulto'}</p>
                </div>
              </div>

              {/* Acciones para el Chofer */}
              {!isDelivered && !isFailed && (
                <div className="space-y-2 pt-1">
                  {/* Botones de Navegación GPS (Waze / Google Maps) */}
                  <div className="grid grid-cols-2 gap-2">
                    <button
                      onClick={() => openWaze(order.latitude, order.longitude)}
                      className="py-2 px-3 bg-sky-600/20 hover:bg-sky-600/30 text-sky-300 rounded-xl text-xs font-semibold flex items-center justify-center gap-1.5 border border-sky-500/30 transition"
                    >
                      <Navigation className="w-3.5 h-3.5 text-sky-400" /> Abrir Waze
                    </button>
                    <button
                      onClick={() => openGoogleMaps(order.latitude, order.longitude)}
                      className="py-2 px-3 bg-slate-700/50 hover:bg-slate-700 text-slate-200 rounded-xl text-xs font-semibold flex items-center justify-center gap-1.5 border border-slate-600 transition"
                    >
                      <MapPin className="w-3.5 h-3.5 text-emerald-400" /> Google Maps
                    </button>
                  </div>

                  {/* Cambio de Estado o Registrar Entrega */}
                  {isPending ? (
                    <button
                      onClick={() => handleStartEnRoute(stop.id)}
                      className="w-full py-2.5 bg-sky-600 hover:bg-sky-500 text-white rounded-xl text-xs font-bold flex items-center justify-center gap-2 shadow-md transition"
                    >
                      <Truck className="w-4 h-4" /> Iniciar Camino a Destino
                    </button>
                  ) : (
                    <button
                      onClick={() => {
                        setSelectedStop(stop);
                        setIsPodModalOpen(true);
                      }}
                      className="w-full py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl text-xs font-bold flex items-center justify-center gap-2 shadow-md transition animate-pulse"
                    >
                      <ShieldCheck className="w-4 h-4" /> Registrar Entrega (POD / Foto / Firma)
                    </button>
                  )}
                </div>
              )}

              {/* Si ya está entregado, opción para ver el comprobante */}
              {(isDelivered || isFailed) && (
                <div className="pt-2 flex justify-end">
                  <button
                    onClick={() => {
                      setSelectedStop(stop);
                      setIsPodModalOpen(true);
                    }}
                    className="text-xs text-sky-400 hover:text-sky-300 flex items-center gap-1 font-semibold"
                  >
                    Ver comprobante POD <ChevronRight className="w-3.5 h-3.5" />
                  </button>
                </div>
              )}

            </div>
          );
        })}
      </main>

      {/* Modal de POD / Firma Táctil */}
      {selectedStop && (
        <PODModal
          stop={selectedStop}
          isOpen={isPodModalOpen}
          onClose={() => setIsPodModalOpen(false)}
          onSubmitPOD={handleSubmitPOD}
          readOnly={selectedStop.status === 'delivered' || selectedStop.status === 'failed'}
        />
      )}

    </div>
  );
}
