'use client';

import React, { useRef, useState, useEffect } from 'react';
import { RouteStop, DeliveryProof } from '../types';
import { CheckCircle, XCircle, FileText, Camera, Edit3, X } from 'lucide-react';

interface PODModalProps {
  stop: RouteStop;
  isOpen: boolean;
  onClose: () => void;
  onSubmitPOD?: (data: {
    route_stop_id: string;
    recipient_name: string;
    recipient_dni: string;
    recipient_relationship: string;
    signature_base64: string;
    photo_url: string;
    failure_reason?: string;
    notes: string;
    status: 'delivered' | 'failed';
  }) => Promise<void>;
  readOnly?: boolean;
}

export default function PODModal({
  stop,
  isOpen,
  onClose,
  onSubmitPOD,
  readOnly = false,
}: PODModalProps) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [recipientName, setRecipientName] = useState('');
  const [recipientDNI, setRecipientDNI] = useState('');
  const [relationship, setRelationship] = useState('titular');
  const [deliveryStatus, setDeliveryStatus] = useState<'delivered' | 'failed'>('delivered');
  const [failureReason, setFailureReason] = useState('cliente_ausente');
  const [notes, setNotes] = useState('');
  const [photoPreview, setPhotoPreview] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (stop.order) {
      setRecipientName(stop.order.customer_name || '');
    }
  }, [stop]);

  if (!isOpen) return null;

  // Manejadores del Canvas de Firma Táctil
  const startDrawing = (e: React.MouseEvent<HTMLCanvasElement> | React.TouchEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    setIsDrawing(true);
    const rect = canvas.getBoundingClientRect();
    const x = 'touches' in e ? e.touches[0].clientX - rect.left : e.clientX - rect.left;
    const y = 'touches' in e ? e.touches[0].clientY - rect.top : e.clientY - rect.top;

    ctx.beginPath();
    ctx.moveTo(x, y);
    ctx.lineWidth = 2.5;
    ctx.lineCap = 'round';
    ctx.strokeStyle = '#0f172a';
  };

  const draw = (e: React.MouseEvent<HTMLCanvasElement> | React.TouchEvent<HTMLCanvasElement>) => {
    if (!isDrawing) return;
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const rect = canvas.getBoundingClientRect();
    const x = 'touches' in e ? e.touches[0].clientX - rect.left : e.clientX - rect.left;
    const y = 'touches' in e ? e.touches[0].clientY - rect.top : e.clientY - rect.top;

    ctx.lineTo(x, y);
    ctx.stroke();
  };

  const stopDrawing = () => {
    setIsDrawing(false);
  };

  const clearCanvas = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);
  };

  // Captura de Foto simulada o via input file
  const handlePhotoUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setPhotoPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!onSubmitPOD) return;

    let signatureBase64 = '';
    if (deliveryStatus === 'delivered' && canvasRef.current) {
      signatureBase64 = canvasRef.current.toDataURL('image/png');
    }

    try {
      setIsSubmitting(true);
      await onSubmitPOD({
        route_stop_id: stop.id,
        recipient_name: recipientName,
        recipient_dni: recipientDNI,
        recipient_relationship: relationship,
        signature_base64: signatureBase64,
        photo_url: photoPreview,
        failure_reason: deliveryStatus === 'failed' ? failureReason : '',
        notes,
        status: deliveryStatus,
      });
      onClose();
    } catch (err: any) {
      alert('Error al guardar comprobante: ' + err.message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const proof = stop.proof;

  return (
    <div className="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4 overflow-y-auto">
      <div className="bg-white rounded-2xl max-w-lg w-full shadow-2xl overflow-hidden border border-slate-100 animate-in fade-in zoom-in-95 duration-200">
        
        {/* Cabecera */}
        <div className="bg-slate-900 text-white p-5 flex items-center justify-between">
          <div>
            <div className="flex items-center gap-2">
              <FileText className="w-5 h-5 text-sky-400" />
              <h3 className="font-bold text-lg">Prueba de Entrega (POD)</h3>
            </div>
            <p className="text-xs text-slate-300 mt-1">
              Guía #{stop.order?.tracking_number} • Parada #{stop.stop_number}
            </p>
          </div>
          <button onClick={onClose} className="p-1 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-white transition">
            <X className="w-6 h-6" />
          </button>
        </div>

        {/* Modal de Solo Lectura (Auditoría en Torre de Control) */}
        {readOnly || proof ? (
          <div className="p-6 space-y-4">
            <div className={`p-4 rounded-xl border flex items-center gap-3 ${
              stop.status === 'delivered' ? 'bg-emerald-50 border-emerald-200 text-emerald-900' : 'bg-rose-50 border-rose-200 text-rose-900'
            }`}>
              {stop.status === 'delivered' ? <CheckCircle className="w-6 h-6 text-emerald-600 shrink-0" /> : <XCircle className="w-6 h-6 text-rose-600 shrink-0" />}
              <div>
                <p className="font-bold text-sm">
                  {stop.status === 'delivered' ? 'Entrega Concretada Exitosamente' : 'Entrega No Concretada / Rechazada'}
                </p>
                {proof?.failure_reason && (
                  <p className="text-xs font-medium text-rose-700 mt-0.5">
                    Motivo: {proof.failure_reason.replace('_', ' ').toUpperCase()}
                  </p>
                )}
                {proof?.submitted_at && (
                  <p className="text-xs opacity-75 mt-0.5">
                    Fecha y Hora: {new Date(proof.submitted_at).toLocaleString('es-CL')}
                  </p>
                )}
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3 text-sm bg-slate-50 p-4 rounded-xl">
              <div>
                <p className="text-xs text-slate-400 font-semibold">RECEPTOR</p>
                <p className="font-medium text-slate-800">{proof?.recipient_name || 'N/A'}</p>
              </div>
              <div>
                <p className="text-xs text-slate-400 font-semibold">RUT / DNI</p>
                <p className="font-medium text-slate-800">{proof?.recipient_dni || 'N/A'}</p>
              </div>
              <div>
                <p className="text-xs text-slate-400 font-semibold">VÍNCULO</p>
                <p className="font-medium text-slate-800 capitalize">{proof?.recipient_relationship || 'N/A'}</p>
              </div>
              <div>
                <p className="text-xs text-slate-400 font-semibold">NOTAS</p>
                <p className="font-medium text-slate-800">{proof?.notes || 'Sin observaciones'}</p>
              </div>
            </div>

            {/* Firma Digital */}
            {proof?.signature_base64 && (
              <div>
                <p className="text-xs font-bold text-slate-700 uppercase tracking-wider mb-1 flex items-center gap-1.5">
                  <Edit3 className="w-3.5 h-3.5 text-sky-600" /> Firma Digital
                </p>
                <div className="border border-slate-200 rounded-xl p-2 bg-white flex items-center justify-center">
                  <img src={proof.signature_base64} alt="Firma del receptor" className="max-h-28 object-contain" />
                </div>
              </div>
            )}

            {/* Foto de Evidencia */}
            {proof?.photo_url && (
              <div>
                <p className="text-xs font-bold text-slate-700 uppercase tracking-wider mb-1 flex items-center gap-1.5">
                  <Camera className="w-3.5 h-3.5 text-sky-600" /> Foto de Evidencia
                </p>
                <div className="border border-slate-200 rounded-xl overflow-hidden max-h-48 flex items-center justify-center bg-slate-900">
                  <img src={proof.photo_url} alt="Foto de entrega" className="object-cover max-h-48 w-full" />
                </div>
              </div>
            )}

            <div className="pt-2">
              <button onClick={onClose} className="w-full py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-800 font-medium rounded-xl text-sm transition">
                Cerrar
              </button>
            </div>
          </div>
        ) : (
          /* Formulario de Cierre de Entrega (Para el Móvil / Conductor) */
          <form onSubmit={handleSubmit} className="p-5 space-y-4 max-h-[80vh] overflow-y-auto">
            {/* Toggle de Resultado */}
            <div className="grid grid-cols-2 gap-2 p-1 bg-slate-100 rounded-xl">
              <button
                type="button"
                onClick={() => setDeliveryStatus('delivered')}
                className={`py-2 text-sm font-semibold rounded-lg flex items-center justify-center gap-2 transition ${
                  deliveryStatus === 'delivered' ? 'bg-emerald-600 text-white shadow' : 'text-slate-600 hover:text-slate-900'
                }`}
              >
                <CheckCircle className="w-4 h-4" /> Entregado
              </button>
              <button
                type="button"
                onClick={() => setDeliveryStatus('failed')}
                className={`py-2 text-sm font-semibold rounded-lg flex items-center justify-center gap-2 transition ${
                  deliveryStatus === 'failed' ? 'bg-rose-600 text-white shadow' : 'text-slate-600 hover:text-slate-900'
                }`}
              >
                <XCircle className="w-4 h-4" /> No Entregado
              </button>
            </div>

            {deliveryStatus === 'delivered' ? (
              <>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-xs font-semibold text-slate-700 mb-1">Nombre Receptor *</label>
                    <input
                      type="text"
                      required
                      value={recipientName}
                      onChange={(e) => setRecipientName(e.target.value)}
                      className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-slate-700 mb-1">RUT / DNI Receptor *</label>
                    <input
                      type="text"
                      required
                      placeholder="12.345.678-9"
                      value={recipientDNI}
                      onChange={(e) => setRecipientDNI(e.target.value)}
                      className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-700 mb-1">Vínculo con el destinatario</label>
                  <select
                    value={relationship}
                    onChange={(e) => setRelationship(e.target.value)}
                    className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none bg-white"
                  >
                    <option value="titular">Titular</option>
                    <option value="familiar">Familiar / Pareja</option>
                    <option value="conserje">Conserjería / Recepción</option>
                    <option value="vecino">Vecino autorizado</option>
                    <option value="companero_trabajo">Compañero de Trabajo</option>
                  </select>
                </div>

                {/* Canvas de Firma Digital */}
                <div>
                  <div className="flex items-center justify-between mb-1">
                    <label className="text-xs font-semibold text-slate-700 flex items-center gap-1">
                      <Edit3 className="w-3.5 h-3.5 text-sky-600" /> Firma del Receptor *
                    </label>
                    <button type="button" onClick={clearCanvas} className="text-xs text-sky-600 hover:text-sky-800 font-medium">
                      Limpiar
                    </button>
                  </div>
                  <div className="border-2 border-dashed border-slate-300 rounded-xl bg-slate-50 overflow-hidden touch-none">
                    <canvas
                      ref={canvasRef}
                      width={400}
                      height={120}
                      className="w-full cursor-crosshair bg-white"
                      onMouseDown={startDrawing}
                      onMouseMove={draw}
                      onMouseUp={stopDrawing}
                      onMouseLeave={stopDrawing}
                      onTouchStart={startDrawing}
                      onTouchMove={draw}
                      onTouchEnd={stopDrawing}
                    />
                  </div>
                  <p className="text-[11px] text-slate-400 mt-0.5 text-center">Firme con el dedo o lápiz táctil dentro del recuadro</p>
                </div>
              </>
            ) : (
              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">Motivo de Rechazo / Fallido *</label>
                <select
                  value={failureReason}
                  onChange={(e) => setFailureReason(e.target.value)}
                  className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-rose-500 focus:outline-none bg-white"
                >
                  <option value="cliente_ausente">Cliente Ausente / No contesta</option>
                  <option value="direccion_incorrecta">Dirección Errónea / Inexistente</option>
                  <option value="rechazado_por_dano">Rechazado por Daño en el Producto</option>
                  <option value="pedido_cancelado">Pedido Cancelado por el Cliente</option>
                  <option value="zona_inaccesible">Zona Inaccesible / Peligrosa</option>
                </select>
              </div>
            )}

            {/* Captura de Foto de Evidencia */}
            <div>
              <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1">
                <Camera className="w-3.5 h-3.5 text-sky-600" /> Foto de Evidencia (Fachada / Paquete)
              </label>
              <input
                type="file"
                accept="image/*"
                capture="environment"
                onChange={handlePhotoUpload}
                className="block w-full text-xs text-slate-500 file:mr-3 file:py-1.5 file:px-3 file:rounded-lg file:border-0 file:text-xs file:font-semibold file:bg-sky-50 file:text-sky-700 hover:file:bg-sky-100"
              />
              {photoPreview && (
                <div className="mt-2 relative rounded-lg overflow-hidden border border-slate-200 max-h-32 bg-black">
                  <img src={photoPreview} alt="Preview" className="w-full h-32 object-cover" />
                  <button
                    type="button"
                    onClick={() => setPhotoPreview('')}
                    className="absolute top-1 right-1 bg-rose-600 text-white rounded-full p-1 shadow"
                  >
                    <X className="w-3 h-3" />
                  </button>
                </div>
              )}
            </div>

            {/* Observaciones */}
            <div>
              <label className="block text-xs font-semibold text-slate-700 mb-1">Observaciones</label>
              <textarea
                rows={2}
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Ej: Se entregó a conserje Juan en recepción edificio."
                className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-sky-500 focus:outline-none"
              />
            </div>

            <div className="pt-3 flex gap-2">
              <button
                type="button"
                onClick={onClose}
                className="flex-1 py-2.5 border border-slate-300 text-slate-700 font-semibold rounded-xl text-sm hover:bg-slate-50 transition"
              >
                Cancelar
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className={`flex-1 py-2.5 text-white font-bold rounded-xl text-sm shadow-md transition ${
                  deliveryStatus === 'delivered' ? 'bg-emerald-600 hover:bg-emerald-700' : 'bg-rose-600 hover:bg-rose-700'
                }`}
              >
                {isSubmitting ? 'Guardando...' : deliveryStatus === 'delivered' ? 'Confirmar Entrega' : 'Registrar Falla'}
              </button>
            </div>
          </form>
        )}

      </div>
    </div>
  );
}
