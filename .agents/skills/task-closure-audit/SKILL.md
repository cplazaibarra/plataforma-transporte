---
name: task-closure-audit
description: Ejecuta una auditoría técnica final antes de declarar completada una tarea que haya modificado código, configuración, infraestructura, migraciones u otros archivos del proyecto.
---

# Task Closure Audit

## 1. Objetivo

Este Skill proporciona un procedimiento común y repetible para verificar una tarea antes de declararla técnicamente terminada.

Su propósito es detectar:

- cambios accidentales;
- archivos fuera de alcance;
- secretos;
- artefactos temporales;
- errores de compilación;
- pruebas faltantes;
- regresiones evidentes;
- dependencias inesperadas;
- acciones Git no autorizadas.

Este Skill verifica.

No corrige automáticamente problemas encontrados.

## 2. Cuándo Utilizarlo

Ejecutar este Skill antes del cierre de una tarea que haya modificado:

- código;
- configuración;
- migraciones;
- infraestructura;
- contratos;
- dependencias;
- scripts;
- archivos versionados relevantes.

No es necesario para una tarea puramente:

- informativa;
- explicativa;
- de análisis;
- de auditoría de solo lectura;

si no hubo modificaciones.

## 3. Responsabilidad

fleet-tech-lead debe utilizar este Skill antes del cierre final de cambios relevantes.

fleet-code-reviewer puede utilizarlo como parte de la revisión final.

Los agentes implementadores pueden utilizarlo antes de entregar su trabajo al Tech Lead.

## 4. Principio de Solo Verificación

Este Skill NO autoriza automáticamente:

- modificar archivos;
- corregir código;
- instalar dependencias;
- ejecutar migraciones;
- hacer commit;
- hacer push;
- hacer merge;
- desplegar;
- borrar archivos.

Si encuentra un problema:

reportarlo.

La corrección debe realizarse mediante la tarea o agente correspondiente.

## 5. Contexto Inicial

Antes de comenzar identificar:

TASK:

objetivo original de la tarea.

EXPECTED_SCOPE:

archivos, módulos o componentes que razonablemente deberían haber cambiado.

ACTUAL_CHANGES:

cambios reales encontrados.

VALIDATION_AVAILABLE:

pruebas, build, linters u otras validaciones disponibles.

## 6. Estado del Repositorio

Inspeccionar el estado actual del repositorio cuando Git esté disponible.

Preferir comandos de solo lectura como:

git status --short

git diff --stat

git diff

No utilizar comandos destructivos.

## 7. Archivos Modificados

Identificar:

- archivos modificados;
- archivos nuevos;
- archivos eliminados;
- archivos renombrados.

Comparar contra EXPECTED_SCOPE.

Clasificar cualquier cambio inesperado como:

OUT_OF_SCOPE_CHANGE

No eliminarlo automáticamente.

## 8. Cambios Concurrentes

No asumir que todos los cambios encontrados pertenecen a la tarea actual.

Puede existir trabajo previo del usuario u otro agente.

Cuando no sea posible atribuir un cambio con seguridad:

marcarlo como:

PRE_EXISTING_OR_CONCURRENT_CHANGE

No sobrescribirlo ni revertirlo.

## 9. Secretos

Revisar los cambios buscando exposición accidental de:

- passwords;
- tokens;
- API keys;
- private keys;
- credenciales;
- connection strings con secretos;
- cookies;
- certificados privados;
- secretos de proveedores.

No imprimir el valor completo de un secreto encontrado.

Reportar:

SECRET_EXPOSURE_RISK

indicando únicamente archivo y ubicación aproximada.

## 10. Información Sensible

Revisar exposición innecesaria de:

- RUT;
- DNI;
- teléfono;
- correo;
- ubicación GPS;
- POD;
- firmas;
- fotografías;
- datos de clientes;
- información entre tenants.

Distinguir datos legítimos del dominio de filtraciones accidentales.

## 11. Archivos Temporales y Artefactos

Buscar cambios accidentales relacionados con:

- logs;
- dumps;
- caches;
- binarios;
- coverage temporal;
- archivos de IDE;
- archivos de prueba temporales;
- .next;
- node_modules;
- builds;
- archivos generados localmente.

Respetar las necesidades reales del repositorio.

No asumir que todo archivo generado es incorrecto.

## 12. Dependencias

Si cambiaron archivos de dependencias:

- go.mod;
- go.sum;
- package.json;
- package-lock.json;
- pnpm-lock.yaml;
- yarn.lock;

verificar que el cambio corresponda a la tarea.

Clasificar dependencias inesperadas como:

UNEXPECTED_DEPENDENCY_CHANGE

## 13. Configuración

Si cambiaron:

- .env.example;
- Docker;
- Compose;
- configuración;
- variables de entorno;
- proxy;
- CI/CD;

verificar que estén dentro del alcance.

No exigir cambios de infraestructura si la tarea no los necesita.

## 14. Migraciones

Si existen migraciones nuevas o modificadas:

confirmar que fueron identificadas explícitamente.

No ejecutar automáticamente la migración.

Indicar que requiere:

database-migration-verifier

cuando dicho Skill esté disponible y corresponda.

## 15. Multi-Tenant

Si el cambio afecta:

- queries;
- endpoints;
- autorización;
- WebSockets;
- recursos de clientes;
- tracking;
- POD;

verificar si existe impacto multi-tenant.

Cuando corresponda indicar:

REQUIRES_MULTI_TENANT_CHECK

No declarar aislamiento correcto sin haberlo validado.

## 16. Build

Cuando exista un mecanismo de build seguro y conocido:

verificar si fue ejecutado.

No inventar comandos.

Utilizar los comandos reales del proyecto.

Reportar:

BUILD_PASS

BUILD_FAIL

BUILD_NOT_RUN

BUILD_NOT_APPLICABLE

## 17. Tests

Identificar las pruebas pertinentes al cambio.

No afirmar que una prueba pasó si no fue ejecutada.

Reportar:

TESTS_PASS

TESTS_FAIL

TESTS_NOT_RUN

TESTS_NOT_AVAILABLE

TESTS_NOT_APPLICABLE

## 18. Lint y Validación Estática

Cuando el proyecto disponga de validaciones existentes, considerar:

- gofmt;
- go vet;
- linters Go;
- TypeScript;
- ESLint;
- validadores SQL;
- validadores de configuración.

No instalar nuevas herramientas únicamente para completar este Skill.

## 19. Regresión

Evaluar si el cambio puede afectar funcionalidad existente.

Considerar especialmente:

- contratos API;
- estados;
- autenticación;
- autorización;
- tenant;
- sincronización offline;
- WebSockets;
- GPS;
- POD;
- rutas;
- tracking público.

Si existe riesgo no validado:

REGRESSION_RISK

## 20. Logging

Revisar que el cambio no introduzca:

- console.log permanente innecesario;
- fmt.Println como logging productivo;
- secretos en logs;
- PII innecesaria;
- payloads completos sensibles;
- overlogging evidente.

Aplicar las reglas globales de logging.

## 21. Código de Debug

Buscar señales evidentes de código temporal:

- TODO crítico introducido por la tarea;
- hacks temporales;
- mocks accidentales;
- endpoints de debug;
- bypass de autenticación;
- flags temporales inseguros.

No considerar automáticamente todo TODO como defecto.

Evaluar contexto.

## 22. Acciones Git

Verificar que no se hayan realizado como consecuencia automática de la tarea:

- commit;
- push;
- merge;
- rebase;
- tag;
- modificación destructiva del historial.

Si fueron explícitamente solicitadas por el usuario, indicarlo como autorizado.

## 23. Acciones Operacionales

Verificar que el cierre no dependa silenciosamente de:

- migration pendiente;
- variable de entorno no documentada;
- servicio externo;
- certificado;
- cambio de DNS;
- reinicio;
- deployment;
- configuración manual.

Reportar dependencias pendientes.

## 24. Documentación

Si el cambio modifica:

- arquitectura;
- contrato;
- configuración;
- procedimiento operacional;

verificar si la documentación correspondiente necesita actualización.

No exigir documentación nueva para cambios triviales.

## 25. Clasificación de Hallazgos

Utilizar:

BLOCKER

Impide cerrar la tarea.

HIGH

Riesgo importante que debe resolverse antes del cierre salvo aceptación explícita.

MEDIUM

Debe evaluarse o corregirse, pero puede no bloquear.

LOW

Mejora menor.

INFO

Información relevante sin defecto.

## 26. Checklist Final

Emitir exactamente esta estructura:

TASK CLOSURE AUDIT

1. SCOPE
Status:
Evidence:

2. UNEXPECTED FILES
Status:
Evidence:

3. SECRETS / SENSITIVE DATA
Status:
Evidence:

4. DEPENDENCIES
Status:
Evidence:

5. BUILD
Status:
Evidence:

6. TESTS
Status:
Evidence:

7. STATIC VALIDATION
Status:
Evidence:

8. REGRESSION RISK
Status:
Evidence:

9. MULTI-TENANT IMPACT
Status:
Evidence:

10. MIGRATION IMPACT
Status:
Evidence:

11. LOGGING / DEBUG CODE
Status:
Evidence:

12. GIT ACTIONS
Status:
Evidence:

13. OPERATIONAL DEPENDENCIES
Status:
Evidence:

14. DOCUMENTATION
Status:
Evidence:

## 27. Veredicto

Finalizar siempre con uno de:

APPROVED

La tarea puede cerrarse.

APPROVED WITH OBSERVATIONS

La tarea puede cerrarse, pero existen observaciones no bloqueantes.

REQUIRES CORRECTION

Existen problemas que deben corregirse antes del cierre.

BLOCKED

No existe información suficiente o existe un bloqueo crítico que impide validar el cambio.

## 28. Evidencia

No declarar:

PASS

sin evidencia suficiente.

La evidencia puede provenir de:

- comandos ejecutados;
- tests;
- build;
- diff;
- inspección del código;
- configuración;
- resultados de QA.

Si algo no fue comprobado:

indicarlo explícitamente.

## 29. Economía

No ejecutar suites enormes innecesariamente.

Seleccionar validaciones proporcionales al cambio.

Una modificación pequeña debe recibir una validación pequeña pero suficiente.

Una modificación transversal requiere validación más amplia.

## 30. Regla Final

Este Skill existe para aumentar confianza en el cierre.

No debe convertirse en una ceremonia que genere cambios innecesarios.

Verificar primero.

Reportar después.

Corregir únicamente mediante autorización o flujo de implementación correspondiente.
