---
name: multi-tenant-isolation-check
description: Audita de forma sistemática el aislamiento de datos y autorización entre tenants en APIs, consultas, WebSockets, procesos y funcionalidades de la plataforma.
---

# Multi-Tenant Isolation Check

## 1. Objetivo

Este Skill verifica que una funcionalidad no permita acceso accidental o intencional a información perteneciente a otro tenant.

Su objetivo es detectar riesgos de:

- acceso cruzado;
- IDOR;
- queries sin filtro de tenant;
- autorización insuficiente;
- WebSockets compartidos incorrectamente;
- recursos públicos con exposición excesiva;
- jobs o procesos sin contexto de tenant;
- relaciones de base de datos ambiguas.

Este Skill verifica.

No corrige automáticamente los problemas encontrados.

## 2. Cuándo Utilizarlo

Utilizar cuando una tarea cree o modifique:

- endpoints;
- handlers;
- services;
- repositories;
- queries SQL;
- tablas;
- relaciones;
- WebSockets;
- tracking;
- POD;
- archivos;
- reportes;
- dashboards;
- procesos batch;
- integraciones;
- caches;
- recursos pertenecientes a clientes.

No ejecutarlo automáticamente en cambios sin relación con datos o autorización.

## 3. Agentes Principales

Puede ser utilizado especialmente por:

fleet-security
fleet-backend-go
fleet-database-postgis
fleet-qa
fleet-tech-lead

## 4. Principio Fundamental

Nunca confiar en un tenant_id enviado libremente por el cliente como única prueba de autorización.

El contexto autorizado del tenant debe provenir de un mecanismo confiable definido por la arquitectura.

No asumir que ese mecanismo es necesariamente JWT si el proyecto actual utiliza otro.

Primero inspeccionar la implementación real.

## 5. Fuente de Verdad

Antes de evaluar aislamiento:

identificar cómo funciona actualmente:

- autenticación;
- sesión;
- identidad;
- tenant;
- roles;
- autorización.

No inventar mecanismos que todavía no existan.

## 6. Identificación del Tenant

Determinar:

TENANT_SOURCE

Ejemplos posibles:

- sesión autenticada;
- claims;
- contexto backend;
- relación usuario-organización;
- credencial de integración.

Clasificar:

TRUSTED

UNTRUSTED

UNKNOWN

## 7. Entrada Manipulable

Identificar cualquier tenant o identificador de recurso proveniente de:

- URL;
- query parameter;
- body;
- header;
- WebSocket message;
- formulario;
- almacenamiento local.

Considerarlo manipulable por defecto hasta demostrar lo contrario.

## 8. Autorización del Recurso

No basta con verificar que el recurso exista.

Debe verificarse que el actor autorizado pueda acceder a ese recurso.

Ejemplo conceptual:

order_id válido

NO implica:

order_id autorizado para este tenant.

## 9. APIs REST

Para endpoints que acceden a recursos de clientes verificar:

- autenticación;
- autorización;
- tenant;
- ownership cuando corresponda;
- filtros;
- errores seguros.

Intentar conceptualmente:

Tenant A solicita recurso de Tenant B.

El resultado esperado debe ser denegación segura.

## 10. Queries

Revisar queries relacionadas con recursos tenantizados.

Buscar ausencia de restricciones equivalentes a:

tenant context + resource criteria

No imponer una sintaxis SQL específica si el modelo utiliza otra estrategia segura.

## 11. Repositories

Cuando exista patrón Repository:

verificar que el aislamiento no dependa exclusivamente de que cada caller recuerde agregar tenant_id.

Identificar riesgos de métodos demasiado genéricos que permitan consultas sin contexto.

## 12. Relaciones

Revisar relaciones indirectas.

Ejemplo conceptual:

tenant
→ order
→ route
→ stop
→ POD.

Un recurso descendiente no debe quedar accesible solamente porque se conoce su ID.

## 13. WebSockets

Verificar:

- autenticación de conexión;
- tenant asociado;
- autorización de suscripciones;
- rooms/channels;
- broadcasts;
- reconexión;
- cambios de sesión.

Un cliente de Tenant A no debe recibir eventos de Tenant B.

## 14. Broadcast

Buscar broadcasts globales de información tenantizada.

Distinguir:

evento realmente global

de:

evento perteneciente a cliente.

No asumir que todo broadcast es inseguro.

## 15. Tracking Público

El tracking público requiere tratamiento especial.

Verificar:

- identificadores difíciles de enumerar cuando corresponda;
- información mínima;
- ausencia de datos internos;
- ausencia de datos de otros pedidos;
- ausencia de información administrativa;
- controles adicionales cuando correspondan.

No asumir que tracking público pertenece a una sesión autenticada.

## 16. POD

Revisar acceso a:

- fotografías;
- firmas;
- documentos;
- comprobantes;
- metadata;
- URLs de descarga.

Conocer una URL o ID no debe permitir acceso cruzado.

## 17. Archivos

Si existen archivos asociados a tenants:

verificar autorización tanto para:

UPLOAD

como para:

DOWNLOAD.

No confiar únicamente en nombres de archivo o paths.

## 18. Base de Datos

Revisar si las entidades tenantizadas tienen suficiente información para aplicar aislamiento de forma consistente.

No rediseñar el esquema durante este Skill.

Si el modelo impide aplicar aislamiento seguro:

reportarlo.

## 19. PostGIS

Las consultas espaciales también deben respetar tenant.

Ejemplos:

- vehículos cercanos;
- rutas;
- posiciones;
- geofences;
- telemetría.

Una consulta geoespacial eficiente pero sin aislamiento sigue siendo incorrecta.

## 20. Telemetría GPS

Verificar que posiciones de vehículos estén asociadas correctamente al contexto autorizado.

Tenant A no debe poder consultar:

- vehículos;
- posiciones;
- historial;
- rutas;

de Tenant B.

## 21. Cache

Si existe caching:

verificar que la clave considere tenant cuando los datos sean tenantizados.

Buscar riesgos como:

cache key = order_id

cuando pudiera requerirse conceptualmente:

tenant + order.

No imponer formato específico.

## 22. Jobs y Procesos Batch

Los procesos internos también requieren contexto correcto.

Verificar:

- imports;
- exports;
- reconciliaciones;
- notificaciones;
- generación de reportes;
- procesos programados.

No asumir que por ejecutarse en backend son automáticamente seguros.

## 23. Integraciones

Cuando una integración externa opere para múltiples clientes:

verificar:

- credenciales;
- contexto;
- configuración;
- datos enviados;
- datos recibidos.

Evitar reutilización incorrecta entre tenants.

## 24. Roles Administrativos

Puede existir un rol autorizado para acceso multi-tenant.

No marcarlo automáticamente como vulnerabilidad.

Debe existir autorización explícita y trazable.

Distinguir:

CROSS_TENANT_AUTHORIZED

de:

CROSS_TENANT_UNAUTHORIZED.

## 25. Errores

Los errores no deben revelar innecesariamente:

- existencia de recursos ajenos;
- nombres de tenants;
- información interna;
- identificadores sensibles.

## 26. Logs

Los logs internos pueden contener tenant_id para diagnóstico cuando sea apropiado.

Pero logs expuestos a clientes nunca deben revelar información de otros tenants.

Aplicar logging-and-telemetry-standards.md.

## 27. Pruebas Negativas

Cuando sea posible definir pruebas como:

Tenant A crea recurso A.

Tenant B intenta:

GET A
UPDATE A
DELETE A
SUBSCRIBE A
DOWNLOAD A

Todas las operaciones no autorizadas deben fallar de forma segura.

## 28. Pruebas Positivas

También verificar:

Tenant A puede utilizar correctamente sus propios recursos.

Seguridad no significa bloquear funcionalidad legítima.

## 29. IDs

UUID, números aleatorios o identificadores difíciles de adivinar NO reemplazan autorización.

Un ID no enumerable reduce descubrimiento.

No establece ownership.

## 30. Frontend

Ocultar información en frontend no constituye control de seguridad.

La autorización debe aplicarse en la capa confiable correspondiente.

## 31. PWA

La PWA puede almacenar temporalmente datos offline.

Verificar que:

- pertenezcan al usuario/contexto correcto;
- no sobrevivan incorrectamente a cambio de sesión;
- no mezclen información entre conductores o tenants;
- sincronización respete autorización.

## 32. Multi-Tenant Futuro

Si una funcionalidad todavía no es multi-tenant:

no inventar infraestructura SaaS completa.

Pero identificar decisiones que crearían un bloqueo evidente para aislamiento futuro.

Clasificarlo como riesgo arquitectónico, no necesariamente como vulnerabilidad actual.

## 33. Evidencia

No declarar aislamiento correcto únicamente porque existe tenant_id.

Buscar evidencia del flujo:

IDENTIDAD
→ CONTEXTO AUTORIZADO
→ RECURSO
→ QUERY/OPERACIÓN
→ RESPUESTA.

## 34. Clasificación

Utilizar:

CRITICAL

Existe posibilidad clara de acceso no autorizado entre tenants a información o acciones sensibles.

HIGH

Existe un bypass probable o control insuficiente.

MEDIUM

Existe riesgo dependiente de condiciones adicionales.

LOW

Mejora defensiva.

INFO

Observación sin vulnerabilidad demostrada.

## 35. Formato de Salida

Emitir:

MULTI-TENANT ISOLATION CHECK

COMPONENT:
SCOPE:
TENANT_SOURCE:

1. AUTHENTICATION CONTEXT
Status:
Evidence:

2. TENANT SOURCE TRUST
Status:
Evidence:

3. RESOURCE AUTHORIZATION
Status:
Evidence:

4. DATABASE / QUERY ISOLATION
Status:
Evidence:

5. INDIRECT RELATIONSHIPS
Status:
Evidence:

6. WEBSOCKET ISOLATION
Status:
Evidence:

7. FILE / POD ACCESS
Status:
Evidence:

8. CACHE / BACKGROUND PROCESSING
Status:
Evidence:

9. PUBLIC EXPOSURE
Status:
Evidence:

10. NEGATIVE TESTS
Status:
Evidence:

11. POSITIVE TESTS
Status:
Evidence:

12. FINDINGS
Severity:
Evidence:
Impact:

## 36. Veredicto

Finalizar con:

ISOLATION VERIFIED

No se identificaron problemas con la evidencia disponible.

ISOLATION VERIFIED WITH OBSERVATIONS

Existen observaciones no bloqueantes.

ISOLATION RISK FOUND

Se identificaron riesgos que requieren corrección.

INSUFFICIENT EVIDENCE

No existe información suficiente para afirmar aislamiento correcto.

## 37. No Corrección Automática

Este Skill no debe modificar:

- handlers;
- queries;
- tablas;
- políticas;
- WebSockets;
- frontend;
- configuración.

Debe entregar hallazgos al agente responsable.

## 38. Escalamiento

Hallazgos CRITICAL o HIGH deben ser enviados a:

fleet-security

y:

fleet-tech-lead.

Los cambios correctivos posteriores deben pasar por QA y Code Review cuando corresponda.

## 39. Regla Final

Nunca afirmar:

"el sistema es multi-tenant seguro"

basándose solamente en arquitectura esperada.

Verificar implementación y evidencia real.
