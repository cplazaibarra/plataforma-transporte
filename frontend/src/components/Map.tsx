'use client';

import React, { useEffect, useRef } from 'react';
import L from 'leaflet';
import { VehiclePosition, RouteStop } from '../types';

interface MapProps {
  vehicles?: VehiclePosition[];
  stops?: RouteStop[];
  center?: [number, number];
  zoom?: number;
  selectedVehicleId?: string | null;
  onSelectStop?: (stop: RouteStop) => void;
}

export default function Map({
  vehicles = [],
  stops = [],
  center = [-33.435, -70.62], // Santiago por defecto
  zoom = 13,
  selectedVehicleId,
  onSelectStop,
}: MapProps) {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<L.Map | null>(null);
  const vehicleMarkersRef = useRef<{ [id: string]: L.Marker }>({});
  const stopMarkersRef = useRef<L.Marker[]>([]);

  useEffect(() => {
    if (!mapContainerRef.current) return;

    // Inicializar mapa de Leaflet solo una vez
    if (!mapRef.current) {
      const map = L.map(mapContainerRef.current, {
        zoomControl: true,
      }).setView(center, zoom);

      // Capa base de OpenStreetMap
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 19,
      }).addTo(map);

      mapRef.current = map;
    }

    return () => {
      // Cleanup al desmontar si es necesario
    };
  }, []);

  // Actualizar marcadores de vehículos
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;

    // Remover marcadores de vehículos que ya no existan
    const currentVehicleIds = new Set(vehicles.map((v) => v.vehicle_id));
    Object.keys(vehicleMarkersRef.current).forEach((id) => {
      if (!currentVehicleIds.has(id)) {
        vehicleMarkersRef.current[id].remove();
        delete vehicleMarkersRef.current[id];
      }
    });

    // Actualizar o crear marcadores para cada vehículo
    vehicles.forEach((v) => {
      const isSelected = selectedVehicleId === v.vehicle_id;
      const truckHtml = `
        <div class="relative flex items-center justify-center">
          <div class="w-10 h-10 rounded-full ${isSelected ? 'bg-sky-500 ring-4 ring-sky-300' : 'bg-slate-900'} text-white flex items-center justify-center shadow-lg transform transition-all duration-300">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8a1 1 0 011-1h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1m-6-1a1 1 0 001 1h1M5 17a2 2 0 104 0m-4 0a2 2 0 114 0m6 0a2 2 0 104 0m-4 0a2 2 0 114 0" />
            </svg>
          </div>
          <div class="absolute -top-7 bg-slate-900/90 text-white text-[11px] font-bold px-2 py-0.5 rounded shadow whitespace-nowrap">
            ${v.plate} • ${v.speed.toFixed(0)} km/h
          </div>
        </div>
      `;

      const truckIcon = L.divIcon({
        className: 'truck-marker',
        html: truckHtml,
        iconSize: [40, 40],
        iconAnchor: [20, 20],
      });

      if (vehicleMarkersRef.current[v.vehicle_id]) {
        // Mover marcador existente con animación suave
        const marker = vehicleMarkersRef.current[v.vehicle_id];
        marker.setLatLng([v.latitude, v.longitude]);
        marker.setIcon(truckIcon);
      } else {
        // Crear nuevo marcador
        const marker = L.marker([v.latitude, v.longitude], { icon: truckIcon }).addTo(map);
        marker.bindPopup(`
          <div class="p-2 text-sm font-sans">
            <p class="font-bold text-slate-800 text-base mb-1">Móvil ${v.plate}</p>
            <p class="text-slate-600"><strong>Modelo:</strong> ${v.model || 'N/A'}</p>
            <p class="text-slate-600"><strong>Chofer:</strong> ${v.driver_name || 'Sin asignar'}</p>
            <p class="text-slate-600"><strong>Velocidad:</strong> ${v.speed.toFixed(1)} km/h</p>
            ${v.route_code ? `<p class="text-sky-600 font-semibold mt-1">Ruta: ${v.route_code}</p>` : ''}
          </div>
        `);
        vehicleMarkersRef.current[v.vehicle_id] = marker;
      }
    });
  }, [vehicles, selectedVehicleId]);

  // Actualizar marcadores de paradas / pedidos
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;

    // Limpiar paradas anteriores
    stopMarkersRef.current.forEach((m) => m.remove());
    stopMarkersRef.current = [];

    stops.forEach((stop) => {
      const order = stop.order;
      if (!order || !order.latitude || !order.longitude) return;

      // Color del pin según estado
      let colorClass = 'bg-amber-500'; // pending
      let statusLabel = 'Pendiente';
      if (stop.status === 'en_route') {
        colorClass = 'bg-sky-600 ring-2 ring-sky-300 animate-bounce';
        statusLabel = 'En camino';
      } else if (stop.status === 'delivered') {
        colorClass = 'bg-emerald-600';
        statusLabel = 'Entregado';
      } else if (stop.status === 'failed') {
        colorClass = 'bg-rose-600';
        statusLabel = 'Fallido';
      }

      const stopHtml = `
        <div class="flex items-center justify-center">
          <div class="w-8 h-8 rounded-full ${colorClass} text-white font-bold text-xs flex items-center justify-center shadow-md border-2 border-white">
            ${stop.stop_number}
          </div>
        </div>
      `;

      const stopIcon = L.divIcon({
        className: 'stop-marker',
        html: stopHtml,
        iconSize: [32, 32],
        iconAnchor: [16, 16],
      });

      const marker = L.marker([order.latitude, order.longitude], { icon: stopIcon }).addTo(map);
      marker.bindPopup(`
        <div class="p-2 text-sm font-sans max-w-xs">
          <div class="flex items-center justify-between gap-2 mb-1">
            <span class="font-bold text-slate-900">#${stop.stop_number} ${order.customer_name}</span>
            <span class="text-[10px] font-bold uppercase px-1.5 py-0.5 rounded text-white ${colorClass}">
              ${statusLabel}
            </span>
          </div>
          <p class="text-slate-600 text-xs mb-1">📍 ${order.delivery_address}</p>
          <p class="text-slate-500 text-xs">📦 ${order.packages_count} bulto(s) • ${order.tracking_number}</p>
          <p class="text-slate-500 text-xs">⏰ Ventana: ${order.window_start} - ${order.window_end}</p>
        </div>
      `);

      if (onSelectStop) {
        marker.on('click', () => onSelectStop(stop));
      }

      stopMarkersRef.current.push(marker);
    });
  }, [stops, onSelectStop]);

  return (
    <div className="relative w-full h-full min-h-[450px]">
      <div ref={mapContainerRef} className="w-full h-full rounded-xl shadow-inner" />
    </div>
  );
}
