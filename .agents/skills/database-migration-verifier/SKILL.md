---
name: database-migration-verifier
description: Audita migraciones PostgreSQL y PostGIS antes de su ejecución, evaluando seguridad, compatibilidad, locks, datos, rollback, rendimiento y riesgo operacional.
---

# Database Migration Verifier

## 1. Objetivo

Este Skill proporciona un procedimiento repetible para revisar migraciones PostgreSQL/PostGIS antes de considerarlas listas para ejecución.

Debe detectar riesgos relacionados con:

- pérdida de datos;
- incompatibilidad;
- locks;
- downtime;
- constraints;
- índices;
- tipos de datos;
- PostGIS;
- grandes volúmenes;
- rollback;
- despliegues;
- multi-tenant;
- compatibilidad entre versiones de aplicación.

Este Skill VERIFICA.

No ejecuta migraciones.

## 2. Cuándo Utilizarlo

Utilizar cuando una tarea:

- cree una migración;
- modifique una migración todavía no aplicada;
- cambie tablas;
- cambie columnas;
- cambie constraints;
- cambie índices;
- cambie tipos;
- cambie relaciones;
- agregue extensiones;
- modifique estructuras PostGIS;
- realice backfills;
- transforme datos.

## 3. Agentes Principales

Especialmente:

fleet-database-postgis
fleet-backend-go
fleet-devops-sre
fleet-tech-lead
fleet-code-reviewer

fleet-security debe participar cuando existan implicancias de seguridad o aislamiento.

## 4. Regla Fundamental

CREAR una migración

NO significa

EJECUTAR una migración.

Este Skill nunca debe ejecutar automáticamente:

- migration up;
- migration down;
- ALTER productivo;
- DROP;
- TRUNCATE;
- DELETE;
- backfill;
- reindex;
- cambios de extensiones.

## 5. Contexto

Antes de revisar identificar:

MIGRATION:
objetivo de la migración.

DATABASE:
PostgreSQL/PostGIS y versión conocida cuando esté disponible.

ENVIRONMENT:
development / test / staging / production / unknown.

DATA_VOLUME:
volumen relevante conocido o UNKNOWN.

DEPLOYMENT_CONTEXT:
cómo convivirá con versiones actuales de la aplicación.

## 6. Fuente de Verdad

Inspeccionar cuando corresponda:

- esquema actual;
- migraciones anteriores;
- modelos;
- repositories;
- queries;
- configuración real.

No asumir que el esquema descrito históricamente sigue siendo idéntico.

## 7. Sintaxis

Validar sintaxis SQL cuando sea posible sin ejecutar cambios destructivos.

No afirmar que una migración es válida solo porque parece correcta visualmente.

## 8. Transacciones

Evaluar si la migración debe ejecutarse dentro de una transacción.

NO imponer BEGIN/COMMIT universalmente.

Algunas operaciones PostgreSQL pueden tener restricciones transaccionales o estrategias operacionales diferentes.

Clasificar:

TRANSACTION_SAFE

TRANSACTION_NOT_APPROPRIATE

TRANSACTION_REQUIRES_REVIEW

## 9. Atomicidad

Cuando una migración pueda ejecutarse atómicamente de forma segura:

preferir atomicidad.

Cuando no pueda:

identificar claramente puntos de fallo parcial y estrategia de recuperación.

## 10. Destructividad

Identificar operaciones como:

DROP TABLE
DROP COLUMN
TRUNCATE
DELETE masivo
ALTER TYPE destructivo
reemplazo de datos
cambio irreversible.

Clasificarlas explícitamente como:

DESTRUCTIVE_CHANGE

No aprobarlas silenciosamente.

## 11. Preservación de Datos

Preguntar:

¿puede perderse información?

¿puede truncarse?

¿puede cambiar significado?

¿puede quedar inaccesible?

¿puede convertirse incorrectamente?

Si existe riesgo:

DATA_LOSS_RISK

## 12. Columnas Nuevas

Para nuevas columnas evaluar:

- NULL / NOT NULL;
- default;
- filas existentes;
- volumen;
- impacto;
- compatibilidad con código anterior.

No asumir que NOT NULL + DEFAULT siempre es gratuito.

Considerar versión real de PostgreSQL y contexto.

## 13. Renombres

Un rename puede romper:

- backend;
- frontend indirectamente;
- queries;
- reportes;
- integraciones;
- analytics.

Evaluar compatibilidad antes de aprobar.

## 14. Cambio de Tipo

Para ALTER COLUMN TYPE evaluar:

- conversión;
- USING;
- pérdida de precisión;
- rewrite;
- locks;
- tamaño;
- compatibilidad.

## 15. Constraints

Para:

PRIMARY KEY
FOREIGN KEY
UNIQUE
CHECK
NOT NULL

evaluar datos existentes antes de asumir que pueden aplicarse.

En tablas grandes considerar estrategias seguras de validación cuando corresponda.

No imponer una estrategia única.

## 16. Foreign Keys

Verificar:

- integridad existente;
- comportamiento ON DELETE;
- comportamiento ON UPDATE;
- índices relevantes;
- impacto operacional.

No agregar CASCADE automáticamente.

## 17. Índices

Para índices evaluar:

- necesidad;
- columnas;
- selectividad;
- queries reales;
- tamaño;
- costo de escritura;
- duplicación con índices existentes.

No agregar índices únicamente por intuición.

## 18. Índices en Producción

En tablas grandes evaluar impacto de creación de índices.

Considerar mecanismos PostgreSQL apropiados para reducir bloqueo cuando corresponda.

No imponer CREATE INDEX CONCURRENTLY sin revisar compatibilidad con el framework de migraciones y contexto transaccional.

## 19. Locks

Identificar operaciones que puedan adquirir locks significativos.

Evaluar:

LOCK_LEVEL

LOCK_DURATION_RISK

BUSINESS_IMPACT

Si no existe información suficiente:

LOCK_RISK_UNKNOWN

## 20. Tablas Grandes

Si una tabla puede ser grande:

no asumir que un ALTER o UPDATE masivo será rápido.

Evaluar:

- filas;
- tamaño;
- rewrite;
- WAL;
- I/O;
- locks;
- duración;
- replicación si existe.

## 21. Backfills

Separar conceptualmente cuando convenga:

SCHEMA CHANGE

de:

DATA BACKFILL.

Un backfill grande puede necesitar estrategia operacional propia.

No ejecutar automáticamente.

## 22. Compatibilidad Hacia Atrás

Evaluar si durante un deployment pueden coexistir:

OLD APPLICATION

y:

NEW DATABASE SCHEMA.

Preferir cambios compatibles durante despliegues progresivos cuando corresponda.

## 23. Expand / Migrate / Contract

Para cambios complejos considerar el patrón conceptual:

EXPAND
→ MIGRATE
→ CONTRACT.

No imponerlo para cambios triviales.

## 24. Rollback

No asumir que toda migración debe tener un DOWN automático.

Evaluar:

ROLLBACK_SQL

ROLL_FORWARD

RESTORE_FROM_BACKUP

COMPENSATING_MIGRATION

MANUAL_RECOVERY

según el cambio.

## 25. Rollback Destructivo

Un DOWN puede ser más peligroso que mantener el cambio.

Ejemplo:

eliminar una columna nueva que ya contiene información.

No aprobar rollback destructivo solo por simetría.

## 26. Idempotencia

No exigir IF EXISTS o IF NOT EXISTS universalmente.

A veces pueden ocultar un estado inesperado del esquema.

Evaluar si la herramienta de migraciones ya garantiza ejecución única.

## 27. Orden

Verificar dependencia con migraciones anteriores.

No asumir que una migración puede ejecutarse de forma independiente.

## 28. PostGIS

Para cambios geoespaciales revisar cuando corresponda:

- geometry/geography;
- SRID;
- tipo geométrico;
- dimensiones;
- índices espaciales;
- conversiones;
- funciones PostGIS;
- volumen.

## 29. SRID

Un cambio de SRID requiere distinguir:

cambiar metadata

de:

transformar coordenadas.

No tratarlos como equivalentes.

## 30. Índices Espaciales

Evaluar necesidad y compatibilidad de índices espaciales según las queries reales.

No crear índices espaciales indiscriminadamente.

## 31. Extensiones

No agregar, eliminar o actualizar extensiones PostgreSQL/PostGIS automáticamente.

Verificar:

- disponibilidad;
- permisos;
- versión;
- compatibilidad;
- impacto operacional.

## 32. Multi-Tenant

Si la migración afecta entidades tenantizadas:

evaluar:

- tenant_id;
- relaciones;
- constraints;
- unicidad;
- índices;
- datos existentes.

Cuando corresponda requerir:

multi-tenant-isolation-check.

## 33. Unicidad Multi-Tenant

Revisar si una regla UNIQUE debe ser:

global

o:

por tenant.

No asumir ninguna de las dos sin definición funcional.

## 34. Datos Existentes

Antes de agregar reglas más estrictas:

considerar si los datos actuales las cumplen.

No asumir base de datos limpia.

## 35. Defaults

Evaluar defaults según:

- semántica;
- compatibilidad;
- filas existentes;
- impacto;
- comportamiento de aplicación.

No utilizar defaults para ocultar ausencia de definición de negocio.

## 36. Timestamps

Revisar coherencia temporal.

Distinguir:

timestamp

de:

timestamptz

según significado real.

No cambiar tipos temporales indiscriminadamente.

## 37. Estados y Enums

Cambios en estados operacionales pueden afectar:

- código;
- reportes;
- analytics;
- filtros;
- integraciones.

Evaluar consumidores antes de modificar.

## 38. Datos Productivos

Nunca utilizar datos productivos reales como campo de prueba destructiva.

No ejecutar operaciones de validación que alteren producción.

## 39. Backup

Para cambios de riesgo significativo determinar si existe requisito de backup previo.

fleet-devops-sre y fleet-database-postgis deben coordinar la estrategia.

Este Skill no crea el backup automáticamente.

## 40. RPO / RTO

No inventar RPO o RTO.

Si la estrategia de migración depende de ellos y no están definidos:

marcar:

BUSINESS_CONTINUITY_REQUIREMENT_UNKNOWN.

## 41. Rendimiento

Evaluar impacto sobre queries y escrituras existentes.

No declarar mejora de rendimiento sin evidencia.

## 42. Plan de Ejecución

Para migraciones relevantes definir conceptualmente:

PRECHECK

EXECUTION

VALIDATION

RECOVERY.

Sin ejecutar ninguno automáticamente.

## 43. Precheck

Puede incluir:

- versión;
- estado de esquema;
- volumen;
- datos incompatibles;
- locks activos;
- espacio;
- dependencias.

Solo utilizar comprobaciones seguras.

## 44. Validación Posterior

Definir qué debería comprobarse después de una eventual ejecución:

- esquema esperado;
- integridad;
- aplicación;
- queries críticas;
- errores;
- rendimiento;
- datos.

Definir no significa ejecutar.

## 45. Migraciones Ya Aplicadas

No modificar silenciosamente una migración que ya haya sido aplicada en otros ambientes.

Preferir una nueva migración correctiva cuando corresponda.

Primero verificar la política real del proyecto.

## 46. Framework de Migraciones

No asumir una herramienta específica.

Primero identificar si el proyecto utiliza:

- golang-migrate;
- goose;
- Atlas;
- scripts SQL;
- otra herramienta.

Adaptar la revisión al mecanismo real.

## 47. Severidad

Clasificar hallazgos:

BLOCKER

Puede causar pérdida de datos, indisponibilidad grave o inconsistencia crítica.

HIGH

Riesgo importante antes de producción.

MEDIUM

Debe revisarse.

LOW

Mejora recomendada.

INFO

Observación.

## 48. Formato de Salida

Emitir:

DATABASE MIGRATION VERIFICATION

MIGRATION:
ENVIRONMENT:
DATABASE VERSION:
MIGRATION TOOL:
DATA VOLUME:

1. PURPOSE
Status:
Evidence:

2. SYNTAX
Status:
Evidence:

3. TRANSACTION STRATEGY
Status:
Evidence:

4. DESTRUCTIVE OPERATIONS
Status:
Evidence:

5. DATA PRESERVATION
Status:
Evidence:

6. LOCK RISK
Status:
Evidence:

7. LARGE TABLE IMPACT
Status:
Evidence:

8. INDEX / CONSTRAINT IMPACT
Status:
Evidence:

9. BACKWARD COMPATIBILITY
Status:
Evidence:

10. POSTGIS IMPACT
Status:
Evidence:

11. MULTI-TENANT IMPACT
Status:
Evidence:

12. ROLLBACK / RECOVERY
Status:
Evidence:

13. DEPLOYMENT DEPENDENCIES
Status:
Evidence:

14. POST-MIGRATION VALIDATION
Status:
Evidence:

15. FINDINGS
Severity:
Evidence:
Impact:

## 49. Veredicto

Finalizar con uno de:

READY FOR CONTROLLED EXECUTION

La migración está técnicamente preparada para una ejecución controlada.

READY WITH OBSERVATIONS

Puede avanzar con observaciones no bloqueantes.

REQUIRES CORRECTION

Debe corregirse antes de ejecución.

REQUIRES OPERATIONAL PLAN

La migración puede ser técnicamente válida pero requiere planificación operacional adicional.

INSUFFICIENT EVIDENCE

No existe información suficiente.

## 50. No Autorización de Ejecución

Un veredicto:

READY FOR CONTROLLED EXECUTION

NO significa:

ejecutar ahora.

La ejecución requiere una acción separada y el contexto correspondiente.

## 51. Escalamiento

Involucrar:

fleet-database-postgis

para diseño y revisión de base de datos.

fleet-devops-sre

para ejecución operacional, backups y deployment.

fleet-security

si afecta aislamiento o datos sensibles.

fleet-tech-lead

para coordinación.

## 52. Regla Final

La migración más elegante no es necesariamente la más segura en producción.

Evaluar siempre:

datos + volumen + locks + compatibilidad + recuperación.
