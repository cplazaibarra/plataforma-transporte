# Multi-Agent Coordination Rules

## 1. Objetivo

Estas reglas definen cómo deben colaborar los agentes especializados del proyecto.

El objetivo es evitar:

- modificaciones contradictorias;
- cambios simultáneos incompatibles;
- contratos desalineados;
- duplicación de trabajo;
- decisiones técnicas conflictivas;
- pérdida de contexto entre especialistas.

## 2. Coordinación Central

fleet-tech-lead es el coordinador principal del sistema de agentes.

Cuando una tarea involucre múltiples dominios:

fleet-tech-lead debe:

- identificar áreas afectadas;
- seleccionar agentes;
- definir orden de participación;
- establecer responsabilidades;
- consolidar decisiones;
- resolver conflictos;
- validar integración final.

## 3. Especialización

Cada agente debe actuar principalmente dentro de su dominio.

Ejemplos:

fleet-product-analyst
→ reglas de negocio y requerimientos.

fleet-data-analytics
→ KPIs y definiciones analíticas.

fleet-database-postgis
→ modelo de datos y rendimiento PostgreSQL/PostGIS.

fleet-backend-go
→ servicios, APIs, WebSockets y lógica backend.

fleet-frontend-next
→ interfaz web.

fleet-mobile-pwa
→ experiencia móvil y offline.

fleet-geo-routing
→ geolocalización, ETA y optimización de rutas.

fleet-security
→ seguridad.

fleet-devops-sre
→ operación e infraestructura.

fleet-qa
→ validación.

fleet-code-reviewer
→ revisión final del código.

Un agente no debe apropiarse innecesariamente del dominio de otro.

## 4. Orden Funcional Antes de Técnico

Si existe ambigüedad funcional:

fleet-product-analyst debe participar antes del diseño técnico.

Cuando corresponda:

USUARIO
→ PRODUCT ANALYST
→ TECH LEAD
→ ESPECIALISTAS TÉCNICOS.

No diseñar implementación para resolver preguntas de negocio no respondidas.

## 5. Orden Analítico

Cuando una métrica dependa de una definición funcional:

PRODUCT ANALYST
→ DATA ANALYTICS
→ TECH LEAD
→ IMPLEMENTACIÓN.

Data Analytics no debe inventar significado de negocio.

## 6. Contratos Compartidos

Cuando dos o más capas compartan un contrato:

definir primero el contrato.

Ejemplos:

- request/response API;
- DTO;
- evento;
- estado;
- formato temporal;
- identificador;
- error;
- payload WebSocket;
- esquema de integración.

No permitir que cada agente invente su propia versión.

## 7. Contratos Backend / Frontend

Para cambios que involucren backend y frontend:

preferir:

1. definición funcional;
2. contrato;
3. backend;
4. frontend;
5. QA.

Frontend no debe asumir campos que Backend todavía no haya confirmado.

Backend no debe cambiar contratos existentes sin considerar consumidores.

## 8. Contratos Backend / Mobile

Para flujos de conductor:

coordinar especialmente:

- estados;
- sincronización;
- idempotencia;
- errores;
- operación offline;
- timestamps;
- identificadores;
- reintentos.

No diseñar la PWA y el backend de forma independiente.

## 9. Contratos Database / Backend

Cuando un cambio requiera persistencia:

fleet-database-postgis y fleet-backend-go deben coordinar:

- significado de datos;
- restricciones;
- integridad;
- consultas;
- migraciones;
- compatibilidad.

Backend no debe imponer estructuras físicas sin análisis de Database.

Database no debe cambiar estructuras sin considerar uso desde Backend.

## 10. Geo y Backend

fleet-geo-routing define:

- modelos;
- reglas espaciales;
- criterios;
- cálculos;
- semántica.

fleet-backend-go implementa e integra cuando corresponda.

No duplicar lógica geoespacial divergente.

## 11. Geo y Database

Cuando se utilicen funciones PostGIS:

fleet-geo-routing define la necesidad geoespacial.

fleet-database-postgis define la implementación eficiente y segura.

## 12. Security

fleet-security debe participar en cambios sensibles.

Los demás agentes no deben ignorar una observación de seguridad de severidad alta o crítica sin resolución o aceptación explícita del riesgo.

## 13. DevOps

fleet-devops-sre debe participar cuando un cambio afecte:

- deployment;
- Docker;
- networking;
- TLS;
- storage;
- variables de entorno;
- observabilidad;
- backups;
- disponibilidad.

Los agentes de aplicación no deben modificar infraestructura crítica unilateralmente.

## 14. QA

fleet-qa debe validar el comportamiento implementado.

QA no redefine el requerimiento.

Si detecta ambigüedad funcional:

debe escalarla al Tech Lead.

## 15. Code Review

fleet-code-reviewer revisa el resultado técnico.

No debe cambiar automáticamente el diseño funcional aprobado.

Si detecta un problema estructural:

debe reportarlo al Tech Lead.

## 16. Responsabilidad Única por Archivo

Cuando sea posible:

un archivo debe tener un responsable principal durante una tarea.

Si varios agentes necesitan cambios sobre el mismo archivo:

fleet-tech-lead debe definir el orden.

Evitar modificaciones simultáneas independientes.

## 17. Cambios Concurrentes

Antes de modificar un archivo que otro agente podría estar trabajando:

verificar estado actual.

No asumir que el contenido visto previamente sigue siendo idéntico.

## 18. Cambios Sobre el Trabajo de Otro Agente

No sobrescribir automáticamente trabajo generado por otro agente.

Primero determinar:

- qué cambió;
- por qué cambió;
- si el cambio es compatible;
- si existe conflicto.

## 19. Handoff

Cuando un agente entrega trabajo a otro debe proporcionar contexto suficiente.

Un handoff técnico debería incluir cuando corresponda:

- objetivo;
- decisiones confirmadas;
- archivos afectados;
- contratos;
- dependencias;
- restricciones;
- riesgos;
- casos límite;
- pendientes;
- criterios de aceptación.

## 20. No Repetir Análisis

Si un especialista ya resolvió una cuestión de su dominio:

otros agentes deben reutilizar esa decisión cuando siga vigente.

No repetir análisis profundo innecesariamente.

## 21. Verificación de Vigencia

Antes de reutilizar una decisión anterior:

verificar que:

- el requerimiento no haya cambiado;
- el código no haya cambiado significativamente;
- la decisión siga siendo aplicable.

## 22. Conflictos Técnicos

Si dos agentes proponen soluciones incompatibles:

no seleccionar arbitrariamente.

fleet-tech-lead debe comparar:

- requerimiento;
- complejidad;
- seguridad;
- mantenibilidad;
- rendimiento;
- impacto;
- compatibilidad;
- costo operacional.

## 23. Conflictos Funcionales

Los agentes técnicos no deben decidir unilateralmente conflictos de negocio.

Escalar a:

fleet-product-analyst

y cuando corresponda al usuario.

## 24. Cambios de Alcance

Si durante una implementación aparece una necesidad nueva fuera del alcance:

no incorporarla automáticamente.

Reportarla como:

HALLAZGO

RIESGO

MEJORA FUTURA

o

BLOQUEANTE

según corresponda.

## 25. Dependencias entre Tareas

Si una tarea depende de otra:

respetar el orden.

Ejemplo:

migración
→ backend
→ frontend
→ QA.

No simular que una dependencia está terminada cuando todavía no lo está.

## 26. Paralelización

Se puede trabajar en paralelo solo cuando las tareas sean suficientemente independientes.

Ejemplos válidos:

- Security revisando diseño mientras QA prepara casos de prueba.
- Frontend trabajando en una vista estática mientras el contrato API ya está confirmado.

Ejemplos riesgosos:

- dos agentes modificando el mismo archivo;
- Backend y Frontend inventando simultáneamente un contrato;
- Database y Backend creando estructuras incompatibles.

## 27. Branch Workspaces

Si Antigravity utiliza workspaces o ramas separadas para subagentes:

usar aislamiento cuando ayude a evitar conflictos.

El uso de branch workspace no elimina la necesidad de coordinar contratos.

No crear ramas Git reales automáticamente solo por esta regla.

## 28. Fuente de Verdad

Para decisiones técnicas actuales:

el repositorio es la fuente principal de verdad.

Para decisiones funcionales confirmadas:

utilizar la definición aprobada vigente.

No confiar exclusivamente en conversaciones antiguas si el proyecto actual demuestra otra realidad.

## 29. Estados de Trabajo

Cuando sea útil, una tarea puede expresarse mediante estados:

ANALISIS

BLOQUEADA

LISTA PARA IMPLEMENTACION

EN IMPLEMENTACION

LISTA PARA QA

CAMBIOS REQUERIDOS

LISTA PARA REVIEW

APROBADA

No es obligatorio implementar estos estados en software.

Son estados conceptuales de coordinación.

## 30. Bloqueos

Un agente debe detener su parte cuando exista un bloqueo real.

Ejemplos:

- requisito funcional crítico no definido;
- contrato pendiente;
- migración incompatible;
- riesgo crítico de seguridad;
- dependencia inexistente.

No ocultar bloqueos con supuestos silenciosos.

## 31. Supuestos

Si es necesario continuar con un supuesto:

debe marcarse explícitamente como:

SUPUESTO

y no como decisión confirmada.

## 32. Integración

Una tarea multiagente no se considera terminada únicamente porque cada agente completó su parte.

Debe validarse la integración entre capas.

## 33. Flujo Recomendado General

Para requerimientos funcionales importantes:

USUARIO
→ PRODUCT ANALYST cuando corresponda
→ TECH LEAD
→ ESPECIALISTAS
→ SECURITY cuando corresponda
→ DEVOPS/SRE cuando corresponda
→ QA
→ CODE REVIEW
→ TECH LEAD
→ CIERRE.

## 34. Flujo de Cambio de Datos

Cuando exista cambio de modelo:

PRODUCT / TECH LEAD
→ DATABASE
→ BACKEND
→ FRONTEND/MOBILE cuando corresponda
→ QA
→ CODE REVIEW.

## 35. Flujo de Funcionalidad Geográfica

Para funcionalidades geográficas:

PRODUCT cuando corresponda
→ TECH LEAD
→ GEO-ROUTING
→ DATABASE cuando corresponda
→ BACKEND
→ FRONTEND/MOBILE
→ QA
→ CODE REVIEW.

## 36. Flujo Analítico

Para KPIs y dashboards:

PRODUCT ANALYST si falta definición
→ DATA ANALYTICS
→ TECH LEAD
→ DATABASE/BACKEND cuando corresponda
→ FRONTEND cuando corresponda
→ QA
→ CODE REVIEW.

## 37. Flujo de Seguridad

Para cambios sensibles:

TECH LEAD
→ ESPECIALISTA IMPLEMENTADOR
→ SECURITY
→ QA
→ CODE REVIEW
→ TECH LEAD.

Security puede participar antes de implementación si el diseño lo requiere.

## 38. Flujo DevOps

Para cambios operacionales:

TECH LEAD
→ DEVOPS/SRE
→ SECURITY cuando corresponda
→ DATABASE cuando corresponda
→ QA o validación operacional
→ CODE REVIEW cuando exista código/configuración versionada
→ TECH LEAD.

## 39. Cierre del Tech Lead

Antes de cerrar una tarea multiagente, fleet-tech-lead debe verificar:

- requerimiento cumplido;
- decisiones consistentes;
- contratos sincronizados;
- integración correcta;
- pruebas realizadas;
- seguridad atendida;
- observaciones de review resueltas;
- cambios fuera de alcance identificados;
- riesgos pendientes documentados.

## 40. Regla de Economía

No convocar todos los agentes para todas las tareas.

Utilizar únicamente los especialistas necesarios.

Una modificación pequeña y localizada no necesita una ceremonia multiagente completa si no existe impacto transversal.

## 41. Regla de No Delegación Circular

Evitar ciclos como:

Tech Lead
→ Backend
→ Tech Lead
→ Backend

sin nueva información.

Cada delegación debe tener un objetivo concreto.

## 42. Regla de Autoridad

La jerarquía operativa es:

USUARIO
→ fleet-tech-lead
→ agentes especializados.

Las reglas globales del proyecto aplican transversalmente.

Ningún subagente puede ignorar una instrucción explícita del usuario que esté dentro del alcance permitido y sea coherente con las reglas del proyecto.

## 43. Cierre

Esta regla no autoriza por sí misma:

- modificar archivos;
- ejecutar código;
- crear branches;
- hacer commits;
- desplegar;
- ejecutar migraciones.

Solo define coordinación.
