# Logging, Errors & Telemetry Standards

## 1. Objetivo

Estas reglas establecen un estándar transversal para:

- logging;
- manejo de errores;
- trazabilidad;
- correlación;
- telemetría;
- diagnóstico operacional.

Aplican a backend, frontend, PWA, infraestructura y servicios relacionados.

## 2. Principio General

Los logs deben servir para diagnosticar problemas reales.

Evitar:

- logs excesivos;
- logs sin contexto;
- mensajes ambiguos;
- información sensible;
- ruido permanente en producción.

## 3. Logging Estructurado

Preferir logging estructurado cuando el componente lo permita.

Cada evento relevante debería poder incluir campos como:

- timestamp;
- level;
- service;
- environment;
- request_id;
- correlation_id;
- tenant_id cuando corresponda;
- user_id anonimizado o identificador interno cuando sea necesario;
- operation;
- duration;
- status;
- error_code.

No todos los campos son obligatorios en todos los eventos.

## 4. Correlation ID

Las operaciones distribuidas deben poder correlacionarse cuando corresponda.

Preferir un identificador como:

correlation_id

o:

request_id

Debe propagarse entre capas cuando sea técnicamente posible.

Ejemplo:

PWA
→ Backend
→ integración externa
→ logs
→ respuesta

No generar identificadores incompatibles entre capas sin necesidad.

## 5. Request ID

Las solicitudes HTTP importantes deberían disponer de un identificador trazable.

Cuando exista un identificador entrante válido, evaluar reutilizarlo.

Si no existe, el backend puede generarlo.

No confiar en valores externos sin validación.

## 6. Niveles de Log

Utilizar niveles coherentes:

DEBUG

Información detallada para diagnóstico técnico.

INFO

Eventos normales relevantes.

WARN

Condiciones anómalas recuperables.

ERROR

Errores que impiden completar una operación.

FATAL o equivalente

Solo para condiciones que impiden que el proceso continúe.

No utilizar ERROR para situaciones normales del negocio.

Ejemplo:

"cliente ausente"

puede ser un evento operacional esperado, no necesariamente un error técnico.

## 7. Errores Técnicos vs Eventos de Negocio

Distinguir:

ERROR TÉCNICO

de:

EVENTO DE NEGOCIO

Ejemplos técnicos:

- timeout;
- fallo de base de datos;
- dependencia externa caída;
- panic;
- error de autenticación interno;
- archivo corrupto.

Ejemplos de negocio:

- entrega fallida;
- cliente ausente;
- dirección incorrecta;
- pedido cancelado;
- rechazo de recepción.

No mezclar ambos conceptos.

## 8. PII y Datos Sensibles

Nunca registrar innecesariamente:

- passwords;
- tokens;
- API keys;
- cookies de sesión;
- claves privadas;
- credenciales;
- firmas completas;
- fotografías POD;
- documentos de identidad completos;
- números de tarjeta;
- headers sensibles.

Minimizar datos personales.

## 9. Identificadores Personales

Cuando sea necesario registrar una referencia a una persona:

preferir identificadores internos.

Evitar registrar:

- nombre completo;
- RUT completo;
- DNI completo;
- teléfono;
- correo;

salvo necesidad operacional justificada.

## 10. GPS y Ubicación

La ubicación GPS puede constituir información sensible.

No registrar posiciones completas de forma indiscriminada en logs generales.

La telemetría operacional debe almacenarse en el mecanismo diseñado para telemetría, no como log de aplicación.

Los logs pueden referenciar eventos de GPS mediante identificadores o metadatos mínimos cuando sea suficiente.

## 11. Payloads

No registrar automáticamente bodies completos de:

- requests;
- responses;
- webhooks;
- POD;
- formularios;
- cargas de archivos.

Si se necesita diagnosticar un problema:

registrar solo campos necesarios o versiones sanitizadas.

## 12. WebSockets

Para WebSocket registrar cuando corresponda:

- conexión;
- desconexión;
- autenticación fallida;
- errores de protocolo;
- reconexión;
- cierre inesperado;
- backpressure;
- mensajes rechazados.

Evitar registrar cada actualización GPS individual en nivel INFO en producción.

## 13. Telemetría GPS

La telemetría masiva no debe saturar el sistema de logging.

Distinguir:

datos operacionales GPS

de:

logs técnicos del pipeline GPS.

Ejemplos de logs útiles:

- pérdida prolongada de telemetría;
- mensajes inválidos;
- retraso excesivo;
- cola saturada;
- error de persistencia;
- dispositivo desconocido.

## 14. Frontend

En frontend de producción evitar:

console.log

dispersos y permanentes.

Los errores importantes deben manejarse mediante mecanismos controlados.

No mostrar detalles internos sensibles al usuario final.

## 15. PWA

La PWA puede requerir diagnóstico local.

Los logs locales deben ser:

- limitados;
- rotables o eliminables;
- no sensibles;
- compatibles con operación offline.

No usar almacenamiento local como repositorio indefinido de logs.

## 16. Backend Go

En backend preferir un logger centralizado y estructurado.

Evitar:

fmt.Println

para logging de producción.

Los handlers y services deberían entregar contexto suficiente para diagnóstico sin duplicar logs innecesariamente.

## 17. Base de Datos

No registrar consultas SQL completas si contienen datos sensibles.

Los errores de base de datos deben permitir identificar:

- operación;
- duración;
- error;
- contexto mínimo;

sin exponer credenciales o información sensible.

## 18. Integraciones Externas

Registrar:

- proveedor;
- operación;
- resultado;
- latencia;
- código de respuesta;
- timeout;
- retry;

cuando corresponda.

No registrar secretos ni payloads completos por defecto.

## 19. Retries

Cuando exista retry:

distinguir entre:

intento inicial

y:

reintento.

Evitar generar múltiples errores idénticos sin contexto.

Cuando sea útil incluir:

attempt_number

max_attempts

retry_reason

## 20. Códigos de Error

Cuando el sistema utilice códigos de error, preferir códigos estables y documentables.

Ejemplo conceptual:

AUTH_INVALID_TOKEN

ORDER_NOT_FOUND

TENANT_ACCESS_DENIED

GPS_INVALID_POSITION

No depender exclusivamente de texto libre para identificar errores.

## 21. Mensajes al Usuario

Los mensajes mostrados al usuario deben ser comprensibles y seguros.

No exponer:

- stack traces;
- SQL;
- rutas internas;
- secretos;
- nombres internos de servicios;
- detalles innecesarios de infraestructura.

## 22. Stack Traces

Stack traces pueden registrarse para errores técnicos relevantes.

No mostrarlos directamente al usuario.

Evitar duplicar el mismo stack trace múltiples veces en la misma cadena de error.

## 23. Panic y Excepciones No Controladas

Los procesos críticos deben tener mecanismos para registrar adecuadamente fallos inesperados.

En Go:

un panic no debe convertirse silenciosamente en comportamiento normal.

Cuando corresponda:

- capturar;
- registrar;
- devolver error seguro;
- permitir diagnóstico.

## 24. Latencia

Para operaciones críticas puede registrarse duración.

Ejemplos:

- API;
- consulta DB;
- integración externa;
- cálculo de ruta;
- procesamiento de POD.

No medir absolutamente todo si genera ruido innecesario.

## 25. Métricas vs Logs

No utilizar logs como sustituto permanente de métricas.

Ejemplo:

cantidad de requests por minuto

debe ser una métrica si se requiere monitoreo continuo.

Logs y métricas tienen objetivos distintos.

## 26. Eventos Operacionales

Cuando el negocio requiera historial:

no asumir que un log técnico reemplaza un evento persistente de negocio.

Ejemplo:

entrega realizada

debe formar parte del modelo operacional correspondiente, no existir únicamente en logs.

## 27. Retención

La retención de logs debe definirse según:

- necesidad operacional;
- volumen;
- privacidad;
- costo;
- cumplimiento.

No asumir retención indefinida.

## 28. Rotación

Los logs no deben crecer indefinidamente.

En ambientes productivos debe existir una estrategia de:

- rotación;
- retención;
- eliminación;
- centralización cuando corresponda.

Coordinar con fleet-devops-sre.

## 29. Tiempo

Los timestamps deben ser consistentes.

Preferir almacenar y transmitir tiempos técnicos en UTC cuando corresponda.

La presentación al usuario puede convertirse a zona horaria local.

No mezclar zonas horarias sin claridad.

## 30. Tenant

Cuando un evento esté asociado a un tenant:

incluir contexto de tenant cuando sea necesario para diagnóstico.

Nunca permitir que logs expuestos a clientes revelen información de otros tenants.

## 31. Ambientes

Los logs deben identificar el ambiente cuando sea relevante:

development
test
staging
production

Evitar confundir errores de pruebas con producción.

## 32. Sanitización

Antes de registrar estructuras provenientes del usuario:

sanitizar o seleccionar campos.

No serializar automáticamente objetos completos solo por conveniencia.

## 33. Sampling

Para eventos de alta frecuencia puede utilizarse sampling cuando sea necesario.

Ejemplo:

telemetría GPS.

El sampling no debe ocultar errores críticos.

## 34. Alertas

No toda línea ERROR debe generar una alerta.

Una alerta operacional debe considerar:

- frecuencia;
- persistencia;
- impacto;
- severidad;
- acción esperada.

Coordinar con fleet-devops-sre.

## 35. Auditoría

Eventos sensibles pueden requerir audit trail específico.

Ejemplos:

- cambios de permisos;
- login administrativo;
- cambios críticos de configuración;
- descarga de POD sensible;
- cambios en tenant.

Audit logs y application logs no deben confundirse.

## 36. Seguridad

fleet-security debe participar cuando el logging implique:

- PII;
- autenticación;
- autorización;
- sesiones;
- tokens;
- POD;
- tracking público;
- integraciones sensibles.

## 37. Observabilidad

fleet-devops-sre es responsable de la infraestructura necesaria para:

- centralización;
- retención;
- métricas;
- alertas;
- dashboards operacionales.

Los agentes de aplicación deben generar señales utilizables.

## 38. Calidad

Los logs deben permitir responder preguntas como:

- qué ocurrió;
- cuándo;
- dónde;
- en qué servicio;
- para qué operación;
- con qué resultado;
- cuánto demoró;
- qué error ocurrió.

No agregar campos sin utilidad real.

## 39. Desarrollo

En desarrollo puede utilizarse mayor nivel de detalle.

Eso no significa que el mismo volumen o sensibilidad de logs sea aceptable en producción.

## 40. Regla de No Overlogging

No registrar cada función ejecutada.

No registrar cada query.

No registrar cada render de componente.

No registrar cada posición GPS.

Logging excesivo:

- aumenta costos;
- dificulta diagnóstico;
- puede filtrar información;
- afecta rendimiento.

## 41. Cambio de Estándar

No introducir una nueva librería de logging o plataforma de observabilidad sin:

- necesidad;
- evaluación;
- compatibilidad;
- impacto;
- aprobación del Tech Lead cuando corresponda.

## 42. Diagnóstico de Incidentes

Durante un incidente:

se puede aumentar temporalmente el nivel de logging.

Debe existir cuidado especial con:

- PII;
- volumen;
- rendimiento.

Después del incidente, restaurar el nivel adecuado.

## 43. Cierre

Antes de aprobar cambios relevantes verificar cuando corresponda:

- logs útiles;
- errores diferenciados;
- secretos protegidos;
- PII minimizada;
- correlation/request ID coherentes;
- ausencia de overlogging;
- mensajes seguros para usuario;
- señales suficientes para diagnóstico.

No modificar automáticamente infraestructura de logging por crear esta regla.
