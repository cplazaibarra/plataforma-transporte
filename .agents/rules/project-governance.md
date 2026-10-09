# Project Governance Rules

## 1. Objetivo

Estas reglas aplican globalmente al desarrollo de la plataforma de logística de última milla.

Su objetivo es mantener:

- arquitectura coherente;
- separación de responsabilidades;
- cambios controlados;
- seguridad;
- calidad;
- trazabilidad;
- compatibilidad;
- mínima complejidad necesaria.

## 2. Principio de Cambio Mínimo

Ante cualquier requerimiento:

comprender primero el problema antes de modificar código.

Implementar el cambio mínimo correcto que resuelva el requerimiento.

Evitar:

- refactors innecesarios;
- reescrituras completas;
- nuevas dependencias sin justificación;
- cambios fuera del alcance;
- modificaciones cosméticas no solicitadas.

## 3. Preservación de Funcionalidad

Nunca romper funcionalidad existente para implementar una nueva característica.

Antes de modificar una funcionalidad existente considerar:

- comportamiento actual;
- dependencias;
- contratos;
- datos;
- usuarios afectados;
- regresiones potenciales.

## 4. No Diseño Prematuro

Si un requerimiento depende de decisiones funcionales todavía no confirmadas:

NO comenzar diseño técnico detallado.

Primero resolver las decisiones funcionales necesarias.

No definir prematuramente:

- tablas;
- columnas;
- relaciones;
- endpoints;
- contratos API;
- componentes;
- infraestructura;
- algoritmos;
- tecnologías.

## 5. Hechos, Supuestos y Propuestas

Cuando exista incertidumbre distinguir explícitamente:

CONFIRMADO

Información establecida por el usuario, negocio, sistema existente o evidencia verificable.

SUPUESTO

Hipótesis temporal necesaria para razonar.

PROPUESTA

Recomendación que todavía requiere aceptación.

PREGUNTA ABIERTA

Decisión que debe ser respondida antes de continuar cuando sea relevante.

Nunca presentar un supuesto o propuesta como hecho confirmado.

## 6. Separación Funcional y Técnica

Las decisiones funcionales deben responder principalmente:

QUÉ
POR QUÉ
QUIÉN
CUÁNDO
REGLAS
EXCEPCIONES

Las decisiones técnicas deben responder principalmente:

CÓMO

No utilizar arquitectura técnica para resolver ambigüedad de negocio.

## 7. Fuente de Verdad

Antes de modificar comportamiento existente:

inspeccionar el código, configuración, modelo o documentación real correspondiente.

No asumir que una descripción histórica representa necesariamente el estado actual del proyecto.

El repositorio actual es la fuente técnica de verdad.

## 8. Especialización

Cuando exista un agente especializado en una materia, utilizarlo para el análisis profundo correspondiente.

El Tech Lead coordina.

Los especialistas analizan e implementan dentro de su dominio.

Evitar que un agente invada innecesariamente responsabilidades de otro.

## 9. Responsabilidad Final

fleet-tech-lead es responsable de:

- coordinación;
- arquitectura global;
- resolución de conflictos;
- integración entre dominios;
- decisión técnica final.

Los agentes especializados no deben redefinir unilateralmente la arquitectura global.

## 10. Seguridad

Aplicar siempre:

- mínimo privilegio;
- aislamiento multi-tenant;
- validación de entradas;
- protección de información sensible;
- gestión segura de credenciales;
- autorización en backend;
- minimización de datos.

Nunca confiar exclusivamente en controles del frontend.

## 11. Multi-Tenant

Toda funcionalidad que manipule información perteneciente a clientes debe considerar aislamiento por tenant.

Nunca mezclar datos entre tenants sin autorización explícita para una vista administrativa global.

## 12. Datos

No eliminar, transformar o migrar datos productivos sin analizar:

- impacto;
- compatibilidad;
- rollback;
- integridad;
- respaldo cuando corresponda.

Las migraciones deben ser controladas y revisables.

## 13. APIs y Contratos

Evitar romper contratos existentes.

Cuando sea necesario cambiar un contrato evaluar:

- compatibilidad;
- consumidores;
- versionado;
- migración;
- rollback.

## 14. Dependencias

No agregar librerías, servicios o tecnologías nuevas únicamente por conveniencia.

Toda dependencia nueva debe justificar:

- necesidad;
- beneficio;
- mantenimiento;
- seguridad;
- impacto operacional.

Preferir capacidades ya existentes cuando sean suficientes.

## 15. Complejidad

No sobreingenierizar.

No introducir:

- microservicios;
- Kubernetes;
- colas;
- event sourcing;
- nuevas bases de datos;
- nuevos frameworks;
- servicios externos;

sin una necesidad real y demostrable.

## 16. Rendimiento

Primero medir.

Luego optimizar.

No realizar optimizaciones prematuras basadas solamente en intuición.

## 17. Observabilidad

Las funcionalidades críticas deben permitir diagnóstico operacional suficiente mediante:

- logs;
- métricas;
- health checks;
- trazabilidad;

cuando corresponda.

Nunca registrar secretos o información sensible innecesaria.

## 18. Calidad

Todo cambio importante debe considerar:

- happy path;
- errores;
- casos límite;
- regresión;
- concurrencia cuando aplique;
- conectividad inestable cuando aplique;
- permisos;
- multi-tenant.

## 19. QA

No considerar una funcionalidad terminada únicamente porque compila.

Debe verificarse su comportamiento.

Los cambios relevantes deben pasar por fleet-qa cuando corresponda.

## 20. Code Review

Los cambios relevantes deben ser revisados por:

fleet-code-reviewer

antes del cierre técnico cuando corresponda.

El objetivo es detectar:

- errores;
- regresiones;
- problemas de arquitectura;
- problemas de seguridad;
- complejidad innecesaria;
- problemas de mantenibilidad.

## 21. DevOps

Cambios relacionados con:

- Docker;
- Docker Compose;
- despliegues;
- TLS;
- backups;
- observabilidad;
- variables de entorno;
- almacenamiento;
- reverse proxy;
- operación;

deben considerar:

fleet-devops-sre

## 22. Seguridad Especializada

Cambios sensibles relacionados con:

- autenticación;
- autorización;
- RBAC;
- multi-tenant;
- tracking público;
- POD;
- archivos;
- PII;
- tokens;
- sesiones;
- integraciones externas;

deben considerar:

fleet-security

## 23. Analítica

No construir indicadores ambiguos.

Antes de implementar un KPI nuevo definir claramente:

- significado;
- fórmula;
- fuente;
- filtros;
- exclusiones;
- periodo;
- dimensiones;
- interpretación.

Utilizar fleet-data-analytics cuando corresponda.

## 24. Estado Actual vs Futuro

Distinguir siempre entre:

IMPLEMENTADO ACTUALMENTE

REQUERIMIENTO ACTUAL

PROPUESTA

MEJORA FUTURA

No presentar una capacidad futura como si ya estuviera implementada.

## 25. SaaS Futuro

La plataforma debe permitir evolución futura hacia SaaS multi-tenant.

Sin embargo:

no implementar complejidad SaaS anticipadamente cuando todavía no sea necesaria.

Diseñar evitando bloqueos arquitectónicos evidentes, pero implementar solamente lo requerido.

## 26. Cierre de Trabajo

Antes de declarar una tarea terminada verificar:

- requerimiento cumplido;
- funcionalidad existente preservada;
- pruebas realizadas;
- errores conocidos;
- seguridad cuando corresponda;
- migraciones cuando corresponda;
- impacto operacional;
- documentación necesaria;
- revisión de código cuando corresponda.

## 27. Regla de No Modificación Accidental

Si una solicitud es únicamente:

- análisis;
- auditoría;
- revisión;
- diagnóstico;
- propuesta;
- explicación;

NO modificar archivos automáticamente.

Modificar archivos solamente cuando la solicitud autorice claramente implementación o cambios.

## 28. Regla de Alcance

No modificar archivos fuera del alcance necesario.

Si durante una tarea aparece un problema no relacionado:

reportarlo separadamente.

No corregirlo automáticamente salvo que sea necesario para completar de forma segura el requerimiento actual.
