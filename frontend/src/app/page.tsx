'use client';

import React, { useState, useEffect, useRef } from 'react';
import dynamic from 'next/dynamic';
import Link from 'next/link';
import { 
  Truck, Package, CheckCircle, AlertTriangle, Clock, MapPin, 
  Search, Plus, Eye, RefreshCw, Smartphone, ExternalLink, ShieldCheck,
  Users, UserCheck, Trash2, Shield, Lock, ChevronRight
} from 'lucide-react';
import { Order, Route, VehiclePosition, RouteStop, User, Vehicle } from '../types';
import { 
  getOrders, getRoutes, getVehiclePositions, createOrder, 
  getUsers, createUser, deleteUser, getVehicles 
} from '../lib/api';
import PODModal from '../components/PODModal';

// Importación dinámica de Leaflet para evitar errores de SSR
const DynamicMap = dynamic(() => import('../components/Map'), {
  ssr: false,
  loading: () => (
    <div className="w-full h-full min-h-[450px] bg-slate-100 animate-pulse flex items-center justify-center text-slate-400 font-medium">
      Cargando mapa interactivo...
    </div>
  ),
});

export default function ControlTowerPage() {
  // Pestaña activa: 'operations' (Despachos y Mapa) o 'users' (Gestión de Usuarios y Roles)
  const [activeTab, setActiveTab] = useState<'operations' | 'users'>('operations');

  useEffect(() => {
    if (typeof window !== 'undefined') {
      const params = new URLSearchParams(window.location.search);
      if (params.get('tab') === 'users') {
        setActiveTab('users');
      }
    }
  }, []);

  // Simulación de rol activo en la interfaz
  const [currentRole, setCurrentRole] = useState<'admin' | 'operator'>('admin');

  // Estados de datos
  const [orders, setOrders] = useState<Order[]>([]);
  const [routes, setRoutes] = useState<Route[]>([]);
  const [vehicles, setVehicles] = useState<VehiclePosition[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [availableVehicles, setAvailableVehicles] = useState<Vehicle[]>([]);

  const [selectedRoute, setSelectedRoute] = useState<Route | null>(null);
  const [selectedStop, setSelectedStop] = useState<RouteStop | null>(null);
  const [isPodModalOpen, setIsPodModalOpen] = useState(false);
  const [isNewOrderModalOpen, setIsNewOrderModalOpen] = useState(false);
  const [isNewUserModalOpen, setIsNewUserModalOpen] = useState(false);
  const [wsConnected, setWsConnected] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [isLoading, setIsLoading] = useState(true);

  // Nuevo pedido state
  const [newCustomerName, setNewCustomerName] = useState('');
  const [newAddress, setNewAddress] = useState('');
  const [newCommune, setNewCommune] = useState('Santiago');
  const [newLat, setNewLat] = useState(-33.4372);
  const [newLng, setNewLng] = useState(-70.6350);
  const [newPackages, setNewPackages] = useState(1);

  // Nuevo usuario state
  const [newUserName, setNewUserName] = useState('');
  const [newUserEmail, setNewUserEmail] = useState('');
  const [newUserPhone, setNewUserPhone] = useState('');
  const [newUserRole, setNewUserRole] = useState<'admin' | 'operator' | 'driver'>('driver');
  const [newUserVehicleId, setNewUserVehicleId] = useState('');

  const wsRef = useRef<WebSocket | null>(null);

  // Cargar datos iniciales
  const loadInitialData = async () => {
    try {
      setIsLoading(true);
      const [fetchedOrders, fetchedRoutes, fetchedVehicles, fetchedUsers, fetchedVehiclesList] = await Promise.all([
        getOrders().catch(() => []),
        getRoutes().catch(() => []),
        getVehiclePositions().catch(() => []),
        getUsers().catch(() => []),
        getVehicles().catch(() => []),
      ]);
      setOrders(fetchedOrders);
      setRoutes(fetchedRoutes);
      setVehicles(fetchedVehicles);
      setUsers(fetchedUsers);
      setAvailableVehicles(fetchedVehiclesList);

      if (fetchedVehiclesList.length > 0 && !newUserVehicleId) {
        setNewUserVehicleId(fetchedVehiclesList[0].id);
      }

      if (fetchedRoutes.length > 0) {
        setSelectedRoute(fetchedRoutes[0]);
      }
    } catch (err) {
      console.error('Error cargando datos:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadInitialData();

    // Conexión WebSocket al backend en Go
    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8095/ws';
    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      setWsConnected(true);
      console.log('✅ Conectado al WebSocket de Torre de Control');
    };

    ws.onclose = () => {
      setWsConnected(false);
      console.log('❌ Desconectado de WebSocket');
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.type === 'telemetry_update') {
          const t = msg.payload;
          setVehicles((prev) => {
            const index = prev.findIndex((v) => v.vehicle_id === t.vehicle_id);
            if (index !== -1) {
              const updated = [...prev];
              updated[index] = {
                ...updated[index],
                latitude: t.latitude,
                longitude: t.longitude,
                speed: t.speed,
                heading: t.heading,
                updated_at: new Date().toISOString(),
              };
              return updated;
            }
            return prev;
          });
        } else if (
          msg.type === 'stop_status_changed' || 
          msg.type === 'pod_submitted' || 
          msg.type === 'order_created' ||
          msg.type === 'user_created' ||
          msg.type === 'user_updated' ||
          msg.type === 'user_deleted'
        ) {
          loadInitialData();
        }
      } catch (e) {
        console.error('Error procesando mensaje WS:', e);
      }
    };

    return () => {
      ws.close();
    };
  }, []);

  // Métricas de despachos
  const totalOrders = orders.length;
  const inTransitCount = orders.filter((o) => o.status === 'in_transit').length;
  const deliveredCount = orders.filter((o) => o.status === 'delivered').length;
  const failedCount = orders.filter((o) => o.status === 'failed').length;

  // Métricas de usuarios
  const adminUsersCount = users.filter((u) => u.role === 'admin').length;
  const operatorUsersCount = users.filter((u) => u.role === 'operator').length;
  const driverUsersCount = users.filter((u) => u.role === 'driver').length;

  const filteredOrders = orders.filter((o) => {
    const matchesSearch = 
      o.customer_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      o.tracking_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
      o.delivery_address.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || o.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const handleCreateOrder = async (e: React.FormEvent) => {
    e.preventDefault();
    if (currentRole === 'operator') {
      alert('Modo Lectura: Los operadores no tienen permisos para crear despachos.');
      return;
    }
    try {
      await createOrder({
        customer_name: newCustomerName,
        delivery_address: newAddress,
        commune: newCommune,
        latitude: Number(newLat),
        longitude: Number(newLng),
        packages_count: Number(newPackages),
      });
      setIsNewOrderModalOpen(false);
      setNewCustomerName('');
      setNewAddress('');
      loadInitialData();
    } catch (err: any) {
      alert('Error creando pedido: ' + err.message);
    }
  };

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    if (currentRole === 'operator') {
      alert('Modo Lectura: Solo los Administradores pueden registrar nuevos usuarios.');
      return;
    }
    try {
      await createUser({
        name: newUserName,
        email: newUserEmail,
        phone: newUserPhone,
        role: newUserRole,
        vehicle_id: newUserRole === 'driver' ? newUserVehicleId : undefined,
      });
      setIsNewUserModalOpen(false);
      setNewUserName('');
      setNewUserEmail('');
      setNewUserPhone('');
      setNewUserRole('driver');
      loadInitialData();
    } catch (err: any) {
      alert('Error al crear usuario: ' + err.message);
    }
  };

  const handleDeleteUser = async (id: string, name: string) => {
    if (currentRole === 'operator') {
      alert('Modo Lectura: Los operadores no pueden eliminar usuarios.');
      return;
    }
    if (confirm(`¿Estás seguro de eliminar al usuario "${name}"?`)) {
      try {
        await deleteUser(id);
        loadInitialData();
      } catch (err: any) {
        alert('Error eliminando usuario: ' + err.message);
      }
    }
  };

  return (
    <div className="flex flex-col min-h-screen bg-slate-50">
      
      {/* Barra Superior */}
      <header className="bg-slate-900 text-white border-b border-slate-800 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-sky-600 to-sky-400 flex items-center justify-center shadow-md shadow-sky-500/20">
                <Truck className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="font-extrabold text-lg tracking-tight flex items-center gap-2">
                  LastMile <span className="text-sky-400 font-medium text-xs bg-sky-950 px-2 py-0.5 rounded-full border border-sky-800/60">PLATAFORMA CENTRAL</span>
                </h1>
                <p className="text-xs text-slate-400">Torre de Control & Gestión de Flota</p>
              </div>
            </div>

            {/* Pestañas Principales */}
            <div className="hidden md:flex items-center gap-1 bg-slate-800/80 p-1 rounded-xl border border-slate-700 ml-4">
              <button
                onClick={() => setActiveTab('operations')}
                className={`px-3 py-1.5 rounded-lg text-xs font-bold transition flex items-center gap-1.5 ${
                  activeTab === 'operations' ? 'bg-sky-600 text-white shadow' : 'text-slate-300 hover:text-white'
                }`}
              >
                <Truck className="w-3.5 h-3.5" /> Monitoreo y Despachos
              </button>
              <button
                onClick={() => setActiveTab('users')}
                className={`px-3 py-1.5 rounded-lg text-xs font-bold transition flex items-center gap-1.5 ${
                  activeTab === 'users' ? 'bg-sky-600 text-white shadow' : 'text-slate-300 hover:text-white'
                }`}
              >
                <Users className="w-3.5 h-3.5" /> Usuarios y Roles ({users.length})
              </button>
            </div>
          </div>

          <div className="flex items-center gap-3">
            {/* Selector de Rol de Sesión (Simulador RBAC) */}
            <div className="flex items-center gap-1.5 bg-slate-800 border border-slate-700 px-2.5 py-1 rounded-xl text-xs">
              <span className="text-slate-400 text-[11px] font-medium hidden sm:inline">Modo:</span>
              <button
                onClick={() => setCurrentRole('admin')}
                className={`px-2 py-0.5 rounded-md font-bold text-[11px] transition flex items-center gap-1 ${
                  currentRole === 'admin' ? 'bg-purple-600 text-white shadow' : 'text-slate-400 hover:text-white'
                }`}
                title="Acceso total para crear, editar y asignar camiones"
              >
                <Shield className="w-3 h-3" /> Administrador
              </button>
              <button
                onClick={() => setCurrentRole('operator')}
                className={`px-2 py-0.5 rounded-md font-bold text-[11px] transition flex items-center gap-1 ${
                  currentRole === 'operator' ? 'bg-amber-600 text-white shadow' : 'text-slate-400 hover:text-white'
                }`}
                title="Acceso restringido a solo lectura y monitoreo"
              >
                <Eye className="w-3 h-3" /> Operador (Lectura)
              </button>
            </div>

            {/* Estado WebSocket */}
            <div className={`hidden lg:flex items-center gap-2 px-2.5 py-1 rounded-full text-xs font-semibold ${
              wsConnected ? 'bg-emerald-950/80 text-emerald-400 border border-emerald-800' : 'bg-rose-950/80 text-rose-400 border border-rose-800'
            }`}>
              <span className={`w-2 h-2 rounded-full ${wsConnected ? 'bg-emerald-400 animate-pulse' : 'bg-rose-400'}`} />
              {wsConnected ? 'EN VIVO' : 'OFFLINE'}
            </div>

            {/* Accesos directos a Móvil y Tracking */}
            <Link
              href="/driver"
              target="_blank"
              className="flex items-center gap-1.5 px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 hover:text-white rounded-lg text-xs font-semibold border border-slate-700 transition"
            >
              <Smartphone className="w-3.5 h-3.5 text-sky-400" /> App Móvil
            </Link>

            <button
              onClick={loadInitialData}
              className="p-2 text-slate-400 hover:text-white hover:bg-slate-800 rounded-lg transition"
              title="Refrescar datos"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
          </div>
        </div>
      </header>

      {/* Banner informativo cuando se está en modo Operador (Lectura) */}
      {currentRole === 'operator' && (
        <div className="bg-amber-500/10 border-b border-amber-500/20 text-amber-800 px-4 py-2 text-xs font-semibold flex items-center justify-between">
          <div className="max-w-7xl mx-auto flex items-center gap-2 w-full">
            <Lock className="w-4 h-4 text-amber-600 shrink-0" />
            <span>
              <strong>Modo Operador Activo (Capacidad de Solo Lectura):</strong> Tienes acceso a monitorear la flota, ver pedidos y auditar entregas. La creación y edición de datos está restringida al rol Administrador.
            </span>
          </div>
        </div>
      )}

      {/* Selector móvil de pestañas */}
      <div className="md:hidden bg-white border-b border-slate-200 px-4 py-2 flex gap-2">
        <button
          onClick={() => setActiveTab('operations')}
          className={`flex-1 py-1.5 text-xs font-bold rounded-lg ${
            activeTab === 'operations' ? 'bg-sky-600 text-white' : 'bg-slate-100 text-slate-700'
          }`}
        >
          Despachos y Mapa
        </button>
        <button
          onClick={() => setActiveTab('users')}
          className={`flex-1 py-1.5 text-xs font-bold rounded-lg ${
            activeTab === 'users' ? 'bg-sky-600 text-white' : 'bg-slate-100 text-slate-700'
          }`}
        >
          Usuarios y Roles ({users.length})
        </button>
      </div>

      {/* VISTA 1: OPERACIONES Y MONITOREO EN VIVO */}
      {activeTab === 'operations' && (
        <>
          {/* Métricas Rápidas */}
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-5 w-full">
            <div className="grid grid-cols-2 md:grid-cols-5 gap-3 sm:gap-4">
              <div className="bg-white p-4 rounded-xl border border-slate-200/80 shadow-sm flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-slate-100 text-slate-700 flex items-center justify-center font-bold">
                  <Package className="w-5 h-5" />
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium">Total Pedidos</p>
                  <p className="text-xl font-bold text-slate-900">{totalOrders}</p>
                </div>
              </div>

              <div className="bg-white p-4 rounded-xl border border-slate-200/80 shadow-sm flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-sky-50 text-sky-600 flex items-center justify-center font-bold">
                  <Clock className="w-5 h-5" />
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium">En Reparto</p>
                  <p className="text-xl font-bold text-sky-600">{inTransitCount}</p>
                </div>
              </div>

              <div className="bg-white p-4 rounded-xl border border-slate-200/80 shadow-sm flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center font-bold">
                  <CheckCircle className="w-5 h-5" />
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium">Entregados (POD)</p>
                  <p className="text-xl font-bold text-emerald-600">{deliveredCount}</p>
                </div>
              </div>

              <div className="bg-white p-4 rounded-xl border border-slate-200/80 shadow-sm flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-rose-50 text-rose-600 flex items-center justify-center font-bold">
                  <AlertTriangle className="w-5 h-5" />
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium">No Entregados</p>
                  <p className="text-xl font-bold text-rose-600">{failedCount}</p>
                </div>
              </div>

              <div className="bg-white p-4 rounded-xl border border-slate-200/80 shadow-sm flex items-center gap-3 col-span-2 md:col-span-1">
                <div className="w-10 h-10 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
                  <Truck className="w-5 h-5" />
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium">Móviles GPS</p>
                  <p className="text-xl font-bold text-indigo-600">{vehicles.length}</p>
                </div>
              </div>
            </div>
          </div>

          {/* Mapa + Despachos */}
          <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-10 w-full flex-1">
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 h-[720px]">
              
              {/* Mapa OpenStreetMap + Leaflet (7 cols) */}
              <div className="lg:col-span-7 bg-white rounded-2xl p-4 border border-slate-200 shadow-sm flex flex-col h-full">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <MapPin className="w-5 h-5 text-sky-600" />
                    <h2 className="font-bold text-slate-900 text-sm">Monitoreo Satelital OpenStreetMap</h2>
                  </div>
                  <div className="flex items-center gap-3 text-xs text-slate-500">
                    <span className="flex items-center gap-1"><span className="w-2.5 h-2.5 rounded-full bg-slate-900" /> Móvil</span>
                    <span className="flex items-center gap-1"><span className="w-2.5 h-2.5 rounded-full bg-sky-600" /> En camino</span>
                    <span className="flex items-center gap-1"><span className="w-2.5 h-2.5 rounded-full bg-emerald-600" /> Entregado</span>
                    <span className="flex items-center gap-1"><span className="w-2.5 h-2.5 rounded-full bg-amber-500" /> Pendiente</span>
                  </div>
                </div>

                <div className="flex-1 w-full rounded-xl overflow-hidden border border-slate-200">
                  <DynamicMap
                    vehicles={vehicles}
                    stops={selectedRoute?.stops || []}
                    onSelectStop={(stop) => {
                      setSelectedStop(stop);
                      setIsPodModalOpen(true);
                    }}
                  />
                </div>

                <div className="mt-3 pt-3 border-t border-slate-100 flex items-center justify-between text-xs">
                  <span className="text-slate-500 font-medium">Flota en movimiento:</span>
                  <div className="flex items-center gap-2">
                    {vehicles.map((v) => (
                      <span key={v.vehicle_id} className="bg-slate-100 border border-slate-200 px-2.5 py-1 rounded-lg font-semibold text-slate-800 flex items-center gap-1.5">
                        <span className="w-2 h-2 rounded-full bg-emerald-500" />
                        {v.plate} ({v.speed.toFixed(0)} km/h)
                      </span>
                    ))}
                  </div>
                </div>
              </div>

              {/* Guías y Despachos (5 cols) */}
              <div className="lg:col-span-5 bg-white rounded-2xl border border-slate-200 shadow-sm flex flex-col h-full overflow-hidden">
                <div className="p-4 border-b border-slate-200 flex items-center justify-between">
                  <div>
                    <h2 className="font-bold text-slate-900 text-sm">Guías y Despachos</h2>
                    <p className="text-xs text-slate-500">{filteredOrders.length} pedidos encontrados</p>
                  </div>
                  <button
                    onClick={() => {
                      if (currentRole === 'operator') {
                        alert('Modo Lectura: La creación de pedidos está deshabilitada para Operadores.');
                        return;
                      }
                      setIsNewOrderModalOpen(true);
                    }}
                    className={`flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-bold shadow-sm transition ${
                      currentRole === 'admin' 
                        ? 'bg-sky-600 hover:bg-sky-700 text-white' 
                        : 'bg-slate-200 text-slate-400 cursor-not-allowed'
                    }`}
                    title={currentRole === 'operator' ? 'Acción restringida a Administrador' : 'Ingresar nuevo despacho'}
                  >
                    <Plus className="w-3.5 h-3.5" /> Nuevo Despacho
                  </button>
                </div>

                <div className="p-3 border-b border-slate-100 bg-slate-50 flex gap-2">
                  <div className="relative flex-1">
                    <Search className="w-4 h-4 absolute left-3 top-2.5 text-slate-400" />
                    <input
                      type="text"
                      placeholder="Buscar por cliente, tracking o dirección..."
                      value={searchTerm}
                      onChange={(e) => setSearchTerm(e.target.value)}
                      className="w-full pl-9 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
                    />
                  </div>
                  <select
                    value={statusFilter}
                    onChange={(e) => setStatusFilter(e.target.value)}
                    className="text-xs bg-white border border-slate-200 rounded-lg px-2 py-1.5 focus:outline-none focus:ring-2 focus:ring-sky-500"
                  >
                    <option value="all">Todos</option>
                    <option value="ready_for_dispatch">En Bodega</option>
                    <option value="in_transit">En Ruta</option>
                    <option value="delivered">Entregado</option>
                    <option value="failed">Fallido</option>
                  </select>
                </div>

                <div className="flex-1 overflow-y-auto divide-y divide-slate-100 p-2">
                  {filteredOrders.map((o) => {
                    let badgeColor = 'bg-slate-100 text-slate-700';
                    let statusText = 'En Bodega';
                    if (o.status === 'in_transit') {
                      badgeColor = 'bg-sky-100 text-sky-700';
                      statusText = 'En Ruta';
                    } else if (o.status === 'delivered') {
                      badgeColor = 'bg-emerald-100 text-emerald-700';
                      statusText = 'Entregado';
                    } else if (o.status === 'failed') {
                      badgeColor = 'bg-rose-100 text-rose-700';
                      statusText = 'Fallido';
                    } else if (o.status === 'assigned') {
                      badgeColor = 'bg-amber-100 text-amber-700';
                      statusText = 'Asignado a Ruta';
                    }

                    return (
                      <div key={o.id} className="p-3 hover:bg-slate-50 rounded-xl transition flex items-start justify-between gap-3">
                        <div className="space-y-1 flex-1">
                          <div className="flex items-center gap-2">
                            <span className="font-bold text-xs text-slate-900">{o.customer_name}</span>
                            <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${badgeColor}`}>
                              {statusText}
                            </span>
                          </div>
                          <p className="text-xs text-slate-500 truncate max-w-[280px]">📍 {o.delivery_address}</p>
                          <div className="flex items-center gap-3 text-[11px] text-slate-400">
                            <span>📦 {o.tracking_number}</span>
                            <span>•</span>
                            <span>⏰ {o.window_start} - {o.window_end}</span>
                          </div>
                        </div>

                        <div className="flex items-center gap-1.5 shrink-0 pt-1">
                          <Link
                            href={`/tracking/${o.tracking_number}`}
                            target="_blank"
                            className="p-1.5 text-slate-400 hover:text-sky-600 hover:bg-sky-50 rounded-lg transition"
                            title="Ver tracking de cliente"
                          >
                            <ExternalLink className="w-4 h-4" />
                          </Link>

                          {(o.status === 'delivered' || o.status === 'failed') && (
                            <button
                              onClick={() => {
                                const stopMatch = selectedRoute?.stops?.find((s) => s.order_id === o.id);
                                if (stopMatch) {
                                  setSelectedStop(stopMatch);
                                  setIsPodModalOpen(true);
                                } else {
                                  alert('Guía entregada con comprobante registrado.');
                                }
                              }}
                              className="p-1.5 text-emerald-600 hover:text-emerald-700 hover:bg-emerald-50 rounded-lg transition"
                              title="Auditar Prueba de Entrega (POD)"
                            >
                              <ShieldCheck className="w-4 h-4" />
                            </button>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>

              </div>

            </div>
          </main>
        </>
      )}

      {/* VISTA 2: GESTIÓN DE USUARIOS Y ROLES (RBAC & ASIGNACIÓN DE CAMIONES) */}
      {activeTab === 'users' && (
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 w-full flex-1">
          
          {/* Tarjetas resumen de usuarios */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-slate-100 text-slate-700 flex items-center justify-center font-bold">
                <Users className="w-5 h-5" />
              </div>
              <div>
                <p className="text-xs text-slate-500 font-medium">Total Usuarios</p>
                <p className="text-xl font-bold text-slate-900">{users.length}</p>
              </div>
            </div>

            <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-purple-50 text-purple-600 flex items-center justify-center font-bold">
                <Shield className="w-5 h-5" />
              </div>
              <div>
                <p className="text-xs text-slate-500 font-medium">Administradores</p>
                <p className="text-xl font-bold text-purple-700">{adminUsersCount}</p>
              </div>
            </div>

            <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-amber-50 text-amber-600 flex items-center justify-center font-bold">
                <Eye className="w-5 h-5" />
              </div>
              <div>
                <p className="text-xs text-slate-500 font-medium">Operadores (Lectura)</p>
                <p className="text-xl font-bold text-amber-700">{operatorUsersCount}</p>
              </div>
            </div>

            <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-sky-50 text-sky-600 flex items-center justify-center font-bold">
                <Truck className="w-5 h-5" />
              </div>
              <div>
                <p className="text-xs text-slate-500 font-medium">Transportistas con Patente</p>
                <p className="text-xl font-bold text-sky-700">{driverUsersCount}</p>
              </div>
            </div>
          </div>

          {/* Tabla de Usuarios */}
          <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
            <div className="p-5 border-b border-slate-200 flex items-center justify-between">
              <div>
                <h2 className="font-bold text-slate-900 text-base">Control de Accesos y Asignación de Móviles</h2>
                <p className="text-xs text-slate-500">
                  Administra los roles del sistema y vincula a los transportistas con la patente del camión en el que operan.
                </p>
              </div>
              <button
                onClick={() => {
                  if (currentRole === 'operator') {
                    alert('Modo Lectura: Solo los Administradores pueden registrar nuevos usuarios.');
                    return;
                  }
                  setIsNewUserModalOpen(true);
                }}
                className={`flex items-center gap-1.5 px-4 py-2 rounded-xl text-xs font-bold shadow transition ${
                  currentRole === 'admin' 
                    ? 'bg-sky-600 hover:bg-sky-700 text-white' 
                    : 'bg-slate-200 text-slate-400 cursor-not-allowed'
                }`}
              >
                <Plus className="w-4 h-4" /> Agregar Usuario
              </button>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs text-slate-600">
                <thead className="bg-slate-50 text-slate-400 uppercase font-semibold text-[11px] border-b border-slate-200">
                  <tr>
                    <th className="px-6 py-3.5">Usuario / Contacto</th>
                    <th className="px-6 py-3.5">Rol de Acceso</th>
                    <th className="px-6 py-3.5">Camión / Patente Asociada</th>
                    <th className="px-6 py-3.5">Estado</th>
                    <th className="px-6 py-3.5 text-right">Acciones</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {users.map((u) => {
                    return (
                      <tr key={u.id} className="hover:bg-slate-50/80 transition">
                        <td className="px-6 py-4">
                          <p className="font-bold text-slate-900 text-sm">{u.name}</p>
                          <p className="text-slate-500 text-xs">{u.email}</p>
                          {u.phone && <p className="text-slate-400 text-[11px]">{u.phone}</p>}
                        </td>

                        <td className="px-6 py-4">
                          {u.role === 'admin' && (
                            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-purple-100 text-purple-800 border border-purple-200">
                              <Shield className="w-3.5 h-3.5 text-purple-600" /> Administrador
                            </span>
                          )}
                          {u.role === 'operator' && (
                            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-amber-100 text-amber-800 border border-amber-200">
                              <Eye className="w-3.5 h-3.5 text-amber-600" /> Operador (Lectura)
                            </span>
                          )}
                          {u.role === 'driver' && (
                            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-sky-100 text-sky-800 border border-sky-200">
                              <Truck className="w-3.5 h-3.5 text-sky-600" /> Conductor / Transporte
                            </span>
                          )}
                        </td>

                        <td className="px-6 py-4">
                          {u.role === 'driver' ? (
                            u.vehicle_plate ? (
                              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-lg bg-slate-900 text-white font-mono text-xs font-bold shadow-sm">
                                <Truck className="w-3.5 h-3.5 text-sky-400" />
                                <span>{u.vehicle_plate}</span>
                                {u.vehicle_model && (
                                  <span className="text-[11px] text-slate-400 font-sans font-normal">
                                    • {u.vehicle_model}
                                  </span>
                                )}
                              </div>
                            ) : (
                              <span className="text-amber-600 font-medium bg-amber-50 px-2 py-0.5 rounded text-[11px] border border-amber-200">
                                Sin camión asignado
                              </span>
                            )
                          ) : (
                            <span className="text-slate-400 text-xs italic">
                              N/A (Plataforma Central)
                            </span>
                          )}
                        </td>

                        <td className="px-6 py-4">
                          <span className="inline-flex items-center gap-1 text-emerald-600 font-semibold">
                            <span className="w-2 h-2 rounded-full bg-emerald-500" />
                            Activo
                          </span>
                        </td>

                        <td className="px-6 py-4 text-right">
                          {currentRole === 'admin' ? (
                            <button
                              onClick={() => handleDeleteUser(u.id, u.name)}
                              className="p-1.5 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition"
                              title="Eliminar usuario"
                            >
                              <Trash2 className="w-4 h-4" />
                            </button>
                          ) : (
                            <span className="text-slate-300 text-xs">Lectura</span>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>

        </main>
      )}

      {/* Modal de Nuevo Usuario */}
      {isNewUserModalOpen && (
        <div className="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl max-w-md w-full shadow-2xl p-6 border border-slate-100">
            <h3 className="font-bold text-lg text-slate-900 mb-1 flex items-center gap-2">
              <Users className="w-5 h-5 text-sky-600" /> Registrar Nuevo Usuario
            </h3>
            <p className="text-xs text-slate-500 mb-4">
              Asigna el nivel de permisos y vincula al transportista con su patente.
            </p>

            <form onSubmit={handleCreateUser} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">Nombre Completo *</label>
                <input
                  type="text"
                  required
                  placeholder="Ej: Marcelo Gómez"
                  value={newUserName}
                  onChange={(e) => setNewUserName(e.target.value)}
                  className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">Correo Electrónico *</label>
                <input
                  type="email"
                  required
                  placeholder="mgomez@transporte.com"
                  value={newUserEmail}
                  onChange={(e) => setNewUserEmail(e.target.value)}
                  className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">Teléfono Móvil</label>
                <input
                  type="tel"
                  placeholder="+56912345678"
                  value={newUserPhone}
                  onChange={(e) => setNewUserPhone(e.target.value)}
                  className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
                />
              </div>

              {/* Selector de Rol */}
              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">Rol de Usuario *</label>
                <select
                  value={newUserRole}
                  onChange={(e) => setNewUserRole(e.target.value as any)}
                  className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none bg-white font-medium"
                >
                  <option value="driver">🚚 Usuario de Transporte / Conductor</option>
                  <option value="operator">👁️ Operador (Capacidad de Lectura)</option>
                  <option value="admin">👑 Administrador del Sistema (Control Total)</option>
                </select>
              </div>

              {/* Selector de Patente si es Rol Conductor */}
              {newUserRole === 'driver' && (
                <div className="bg-sky-50 p-3.5 rounded-xl border border-sky-200">
                  <label className="block text-xs font-bold text-sky-900 mb-1 flex items-center gap-1.5">
                    <Truck className="w-3.5 h-3.5 text-sky-600" /> Patente del Camión Asignado *
                  </label>
                  <select
                    value={newUserVehicleId}
                    onChange={(e) => setNewUserVehicleId(e.target.value)}
                    className="w-full border border-sky-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none bg-white font-bold text-slate-800"
                  >
                    {availableVehicles.map((v) => (
                      <option key={v.id} value={v.id}>
                        Patente {v.plate} - {v.brand} {v.model} ({v.capacity_kg} kg)
                      </option>
                    ))}
                  </select>
                  <p className="text-[11px] text-sky-700 mt-1">
                    El chofer reportará la telemetría y entregas bajo esta patente.
                  </p>
                </div>
              )}

              <div className="pt-3 flex gap-2">
                <button
                  type="button"
                  onClick={() => setIsNewUserModalOpen(false)}
                  className="flex-1 py-2 border border-slate-300 text-slate-700 font-semibold rounded-lg text-sm hover:bg-slate-50 transition"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="flex-1 py-2 bg-sky-600 text-white font-semibold rounded-lg text-sm hover:bg-sky-700 shadow transition"
                >
                  Guardar Usuario
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Modal de Auditoría de POD */}
      {selectedStop && (
        <PODModal
          stop={selectedStop}
          isOpen={isPodModalOpen}
          onClose={() => setIsPodModalOpen(false)}
          readOnly={true}
        />
      )}

      {/* Modal de Nuevo Despacho */}
      {isNewOrderModalOpen && (
        <div className="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl max-w-md w-full shadow-2xl p-6 border border-slate-100">
            <h3 className="font-bold text-lg text-slate-900 mb-4 flex items-center gap-2">
              <Package className="w-5 h-5 text-sky-600" /> Ingresar Nuevo Despacho
            </h3>
            <form onSubmit={handleCreateOrder} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">Nombre Destinatario *</label>
                <input
                  type="text"
                  required
                  placeholder="Ej: Laura Mendoza"
                  value={newCustomerName}
                  onChange={(e) => setNewCustomerName(e.target.value)}
                  className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">Dirección de Entrega *</label>
                <input
                  type="text"
                  required
                  placeholder="Ej: Av. Apoquindo 4500"
                  value={newAddress}
                  onChange={(e) => setNewAddress(e.target.value)}
                  className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-700 mb-1">Latitud GPS</label>
                  <input
                    type="number"
                    step="0.0001"
                    value={newLat}
                    onChange={(e) => setNewLat(Number(e.target.value))}
                    className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-700 mb-1">Longitud GPS</label>
                  <input
                    type="number"
                    step="0.0001"
                    value={newLng}
                    onChange={(e) => setNewLng(Number(e.target.value))}
                    className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
                  />
                </div>
              </div>
              <div className="pt-2 flex gap-2">
                <button
                  type="button"
                  onClick={() => setIsNewOrderModalOpen(false)}
                  className="flex-1 py-2 border border-slate-300 text-slate-700 font-semibold rounded-lg text-sm hover:bg-slate-50 transition"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="flex-1 py-2 bg-sky-600 text-white font-semibold rounded-lg text-sm hover:bg-sky-700 shadow transition"
                >
                  Guardar Guía
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}
