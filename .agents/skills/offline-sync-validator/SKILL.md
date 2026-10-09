---
name: offline-sync-validator
description: Audita flujos offline y sincronización diferida entre la PWA del conductor y el backend, verificando persistencia local, idempotencia, reintentos, evidencias POD, GPS, conflictos y recuperación ante conectividad inestable.
---

# Offline Sync Validator

## 1. Objetivo

Este Skill proporciona un procedimiento sistemático para verificar funcionalidades que deben continuar operando cuando el conductor pierde conectividad.

Debe evaluar:

- persistencia local;
- cierre y reapertura de la PWA;
- cola de sincronización;
- idempotencia;
- reintentos;
- duplicados;
- orden de eventos;
- conflictos;
- POD;
- fotografías;
- firmas;
- GPS;
- recuperación de conectividad;
- confirmación real del servidor.

Este Skill VERIFICA.

No implementa ni modifica automáticamente el flujo offline.

## 2. Cuándo Utilizarlo

Utilizar cuando una funcionalidad involucre:

- operación sin Internet;
- sincronización diferida;
- POD;
- fotografías;
- firma;
- llegada a parada;
- entrega;
- entrega fallida;
- actualización de estados;
- telemetría GPS almacenada temporalmente;
- reintentos;
- recuperación de conectividad.

## 3. Agentes Principales

Especialmente:

fleet-mobile-pwa
fleet-backend-go
fleet-qa
fleet-tech-lead

Cuando corresponda:

fleet-security
fleet-database-postgis.

## 4. Principio Fundamental

Una operación offline crítica no debe depender únicamente de memoria React o estado temporal del navegador.

Debe sobrevivir, cuando el requerimiento lo exija, a:

- pérdida de Internet;
- cambio de vista;
- cierre de la PWA;
- reapertura;
- reinicio del dispositivo cuando la plataforma lo permita.

## 5. Fuente de Verdad

Antes de evaluar:

inspeccionar la implementación real.

No asumir automáticamente:

- IndexedDB;
- Service Worker;
- Background Sync;
- UUID;
- Idempotency-Key;
- una librería específica.

Evaluar primero el mecanismo existente.

## 6. Flujo Base

Para una operación offline crítica evaluar conceptualmente:

USER ACTION
→ LOCAL VALIDATION
→ LOCAL PERSISTENCE
→ LOCAL CONFIRMATION
→ SYNC QUEUE
→ CONNECTIVITY AVAILABLE
→ SERVER SUBMISSION
→ SERVER CONFIRMATION
→ LOCAL RECONCILIATION.

No exigir exactamente esta implementación si existe otra equivalente y segura.

## 7. Persistencia Local

Determinar:

qué información debe sobrevivir a pérdida de conectividad.

Verificar que los datos críticos no existan únicamente en:

- memoria;
- estado React;
- variables temporales.

## 8. IndexedDB

Para datos estructurados, colas, fotografías o blobs:

IndexedDB puede ser apropiado.

No imponer IndexedDB si la implementación actual utiliza otro mecanismo robusto compatible.

No utilizar localStorage para grandes evidencias binarias.

## 9. Estados de Sincronización

La aplicación debe poder distinguir conceptualmente entre:

LOCAL_ONLY

PENDING

SYNCING

SYNCED

FAILED

o estados equivalentes.

No asumir que estos nombres deben existir literalmente en código.

## 10. Confirmación Visual

El conductor debe poder distinguir:

operación guardada localmente

de:

operación confirmada por servidor.

No mostrar:

"Sincronizado"

si únicamente se guardó localmente.

## 11. Idempotencia

Una operación reintentada no debe generar efectos duplicados.

Evaluar idempotencia especialmente en:

- entrega;
- POD;
- entrega fallida;
- cambio de estado;
- evidencias;
- eventos operacionales.

## 12. Clave de Idempotencia

Puede utilizarse:

- UUID;
- operation_id;
- client_event_id;
- Idempotency-Key;
- identificador equivalente.

No imponer UUID v4 ni una cabecera específica si existe otra estrategia segura.

Lo importante es la propiedad de idempotencia.

## 13. Backend

El backend debe poder reconocer operaciones repetidas cuando el flujo lo requiera.

No confiar exclusivamente en que el cliente nunca enviará dos veces.

## 14. Reintentos

Los reintentos deben ser controlados.

Evaluar:

- cantidad;
- intervalo;
- backoff;
- errores recuperables;
- errores permanentes;
- conectividad;
- respuesta del servidor.

No imponer una cantidad fija de intentos sin definición.

## 15. Backoff

Puede utilizarse backoff cuando sea apropiado.

Evitar loops agresivos que:

- consuman batería;
- saturen backend;
- consuman datos móviles;
- generen múltiples operaciones simultáneas.

## 16. Errores Recuperables

Ejemplos posibles:

- timeout;
- sin conectividad;
- HTTP 5xx temporal;
- dependencia temporalmente caída.

Deben distinguirse de errores como:

- operación inválida;
- autorización denegada;
- recurso incompatible.

No reintentar indefinidamente errores permanentes.

## 17. Orden de Eventos

Evaluar si el orden importa.

Ejemplo conceptual:

ARRIVAL
→ DELIVERY
→ POD.

Si el backend recibe eventos fuera de orden:

debe existir comportamiento definido.

No inventar automáticamente la regla de negocio.

## 18. Timestamps

Distinguir cuando corresponda:

EVENT_TIME

momento real de la acción del conductor.

SYNC_TIME

momento en que el servidor recibió la operación.

No reemplazar automáticamente EVENT_TIME por la hora de sincronización.

## 19. Zona Horaria

Mantener representación temporal coherente.

Preferir formatos no ambiguos.

No inferir la hora real del evento únicamente desde la hora del servidor cuando la operación ocurrió offline.

## 20. POD

Para Proof of Delivery offline verificar:

- estado;
- timestamp;
- posición cuando corresponda;
- firma;
- fotografía;
- identificador de operación;
- persistencia local;
- estado de sincronización.

No declarar POD remoto confirmado hasta recibir confirmación del servidor.

## 21. Fotografías

Evaluar:

- persistencia;
- tamaño;
- compresión cuando corresponda;
- orientación;
- reintentos;
- upload parcial;
- duplicados;
- eliminación local.

No eliminar la evidencia local antes de confirmación segura del servidor.

## 22. Firma

La firma debe conservarse de manera suficientemente robusta durante el periodo offline.

No asumir que un Canvas en memoria constituye persistencia.

## 23. GPS

Si la operación requiere georreferencia:

guardar cuando corresponda:

- posición;
- timestamp;
- precisión;
- fuente disponible.

No reemplazar una posición capturada durante la acción por una posición posterior durante sincronización.

## 24. Telemetría GPS

Para telemetría continua:

evaluar buffering cuando sea necesario.

No exigir que cada posición se conserve indefinidamente.

Considerar:

- volumen;
- batería;
- almacenamiento;
- orden;
- utilidad operacional.

## 25. Background Sync

No asumir que Background Sync está disponible de manera idéntica en:

- Chrome Android;
- Safari iOS;
- PWA instalada;
- navegador normal.

El flujo debe considerar las capacidades reales de la plataforma.

## 26. Reapertura Manual

Si sincronización automática en background no está disponible:

la PWA debe poder continuar el proceso cuando el usuario vuelva a abrirla, según el requerimiento.

## 27. Cierre de Aplicación

Probar conceptualmente:

acción offline
→ persistencia
→ cierre
→ reapertura
→ recuperación del pendiente.

Si el dato desaparece:

OFFLINE_DATA_LOSS_RISK.

## 28. Pérdida Durante Upload

Probar:

upload comienza
→ conectividad se pierde
→ estado incierto.

El sistema no debe asumir:

SUCCESS

ni crear duplicados ciegamente.

## 29. Confirmación Perdida

Caso crítico:

servidor procesa correctamente

pero:

cliente pierde la respuesta.

Al reintentar:

la operación debe poder reconciliarse sin duplicar efectos.

## 30. Conflictos

Evaluar qué ocurre cuando:

estado local

y:

estado servidor

han cambiado de forma incompatible.

No inventar automáticamente estrategia:

last-write-wins

server-wins

client-wins.

La política debe depender del dominio.

## 31. Reconciliación

Cuando exista conflicto:

identificar si puede resolverse automáticamente o requiere intervención.

Operaciones críticas pueden requerir tratamiento explícito.

## 32. Sesión Expirada

Evaluar:

usuario realiza operaciones offline

→ sesión expira
→ vuelve conectividad.

No descartar automáticamente datos pendientes.

Tampoco sincronizarlos sin autorización válida.

Debe existir comportamiento seguro.

## 33. Cambio de Usuario

Caso crítico:

Conductor A tiene pendientes locales.

Se cierra sesión.

Conductor B inicia sesión.

Los pendientes de A no deben mezclarse con B.

## 34. Cambio de Tenant

Si la plataforma permite contextos tenantizados:

los datos offline deben mantener asociación correcta.

Nunca sincronizar pendientes bajo un tenant diferente.

Cuando corresponda utilizar:

multi-tenant-isolation-check.

## 35. Duplicados

Evaluar duplicados en:

- eventos;
- fotografías;
- POD;
- estados;
- telemetría.

No deduplicar únicamente por timestamp si no es suficientemente confiable.

## 36. Cola

La cola debe permitir identificar conceptualmente:

- operación;
- estado;
- intentos;
- última falla;
- siguiente intento cuando corresponda;
- identidad/contexto;
- fecha del evento.

No imponer un esquema físico específico.

## 37. Tamaño de Cola

Considerar escenarios de desconexión prolongada.

Evaluar:

- crecimiento;
- almacenamiento;
- fotografías;
- telemetría;
- límites del navegador.

No asumir almacenamiento infinito.

## 38. Presión de Almacenamiento

Si el dispositivo tiene poco espacio:

la aplicación debe evitar corrupción silenciosa o pérdida inadvertida.

Si no existe estrategia:

reportar riesgo.

## 39. Batería

No utilizar sincronización o GPS agresivo sin considerar batería.

No optimizar prematuramente, pero identificar loops o polling evidentemente excesivos.

## 40. Conectividad

No confiar únicamente en:

navigator.onLine

como prueba de acceso real al backend.

Puede existir:

red local

sin:

conectividad efectiva al servicio.

## 41. Seguridad Local

Los datos almacenados offline pueden incluir información sensible.

Evaluar:

- minimización;
- duración;
- logout;
- cambio de usuario;
- POD;
- ubicación;
- credenciales.

fleet-security debe participar cuando corresponda.

## 42. Tokens

No almacenar tokens de forma insegura solo para facilitar sincronización offline.

La estrategia de autenticación debe respetar la arquitectura real.

## 43. Logs

No registrar fotografías, firmas, tokens ni payloads POD completos para diagnosticar sincronización.

Aplicar:

logging-and-telemetry-standards.md.

## 44. Pruebas Mínimas

Cuando aplique, verificar escenarios:

ONLINE
→ operación
→ confirmación.

OFFLINE
→ operación
→ persistencia.

OFFLINE
→ cerrar PWA
→ abrir PWA
→ pendiente continúa.

OFFLINE
→ ONLINE
→ sincronización.

SYNC
→ timeout
→ retry.

SERVER SUCCESS
→ RESPONSE LOST
→ retry.

USER A pending
→ logout
→ USER B login.

## 45. Pruebas Adicionales

Cuando corresponda:

- múltiples pendientes;
- orden alterado;
- fotografía grande;
- poco almacenamiento;
- GPS sin precisión;
- sesión expirada;
- HTTP 401;
- HTTP 409;
- HTTP 429;
- HTTP 500;
- conexión intermitente;
- cierre durante upload.

## 46. Evidencia

No declarar que el modo offline funciona únicamente porque existe un Service Worker.

Buscar evidencia del flujo completo.

## 47. Severidad

Clasificar:

BLOCKER

Puede provocar pérdida de evidencia, duplicación crítica o corrupción operacional.

HIGH

Falla importante en sincronización o seguridad.

MEDIUM

Problema recuperable con impacto operacional.

LOW

Mejora recomendada.

INFO

Observación.

## 48. Formato de Salida

Emitir:

OFFLINE SYNC VALIDATION

FEATURE:
CLIENT:
BACKEND:
OFFLINE STORAGE:

1. LOCAL PERSISTENCE
Status:
Evidence:

2. APP RESTART RECOVERY
Status:
Evidence:

3. SYNC QUEUE
Status:
Evidence:

4. IDEMPOTENCY
Status:
Evidence:

5. RETRY STRATEGY
Status:
Evidence:

6. EVENT ORDER
Status:
Evidence:

7. TIMESTAMPS
Status:
Evidence:

8. POD / EVIDENCE
Status:
Evidence:

9. GPS
Status:
Evidence:

10. SERVER CONFIRMATION
Status:
Evidence:

11. LOST RESPONSE RECOVERY
Status:
Evidence:

12. SESSION / USER CHANGE
Status:
Evidence:

13. MULTI-TENANT
Status:
Evidence:

14. STORAGE / BATTERY
Status:
Evidence:

15. SECURITY
Status:
Evidence:

16. TEST COVERAGE
Status:
Evidence:

17. FINDINGS
Severity:
Evidence:
Impact:

## 49. Veredicto

Finalizar con:

OFFLINE FLOW VERIFIED

Existe evidencia suficiente de funcionamiento correcto.

OFFLINE FLOW VERIFIED WITH OBSERVATIONS

Funciona, pero existen observaciones no bloqueantes.

OFFLINE FLOW RISK FOUND

Existen problemas que requieren corrección.

INSUFFICIENT EVIDENCE

No existe evidencia suficiente para validar el flujo.

## 50. No Corrección Automática

Este Skill no modifica:

- PWA;
- backend;
- IndexedDB;
- Service Worker;
- endpoints;
- base de datos;
- configuración.

Entrega hallazgos a los agentes correspondientes.

## 51. Coordinación

Problemas de PWA:

fleet-mobile-pwa.

Problemas de backend/idempotencia:

fleet-backend-go.

Problemas de persistencia:

fleet-database-postgis cuando corresponda.

Problemas de seguridad:

fleet-security.

Pruebas:

fleet-qa.

Coordinación:

fleet-tech-lead.

## 52. Regla Final

Offline-first no significa simplemente:

"la pantalla abre sin Internet".

Significa que las operaciones críticas pueden conservarse, recuperarse y sincronizarse de forma segura y verificable según el requerimiento funcional.
