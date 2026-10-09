# Git & Change Management Rules

## 1. Objetivo

Estas reglas establecen cómo los agentes deben trabajar con cambios del repositorio.

El objetivo es:

- evitar modificaciones accidentales;
- mantener cambios pequeños y trazables;
- evitar conflictos entre agentes;
- proteger ramas principales;
- facilitar revisión y rollback.

## 2. Git No Es Automático

Los agentes NO deben automáticamente:

- crear commits;
- crear ramas;
- hacer merge;
- hacer rebase;
- hacer push;
- crear tags;
- eliminar ramas;
- modificar historial Git.

Estas acciones requieren una solicitud explícita del usuario o autorización clara dentro de una tarea.

## 3. Inspeccionar Antes de Modificar

Antes de modificar archivos existentes, revisar cuando corresponda:

- estado actual del archivo;
- cambios relacionados;
- dependencias;
- contratos;
- contexto suficiente del repositorio.

No sobrescribir trabajo existente basándose solamente en una descripción histórica.

## 4. Cambios Mínimos

Cada implementación debe modificar únicamente los archivos necesarios.

Evitar mezclar en una misma tarea:

- nueva funcionalidad;
- refactorizaciones no relacionadas;
- cambios cosméticos;
- limpieza general;
- actualización masiva de dependencias.

## 5. Trabajo Existente

Nunca descartar automáticamente cambios existentes que no pertenezcan al agente actual.

Si aparecen modificaciones previas o concurrentes:

- preservarlas;
- analizar si afectan la tarea;
- informar conflictos cuando corresponda.

No utilizar comandos destructivos para restaurar el repositorio sin autorización.

## 6. Comandos Destructivos

NO ejecutar automáticamente:

git reset --hard
git clean -fd
git checkout -- .
git restore .
git push --force
git push --force-with-lease

ni equivalentes destructivos.

Solo pueden considerarse con autorización explícita y comprensión clara del impacto.

## 7. Ramas Principales

No asumir permiso para modificar directamente ramas como:

main
master
production
release

La estrategia de ramas debe respetar la política real del repositorio.

Si no existe una política definida, no inventarla automáticamente.

## 8. Commits

Si el usuario solicita crear un commit:

el cambio debe ser coherente y revisable.

Evitar commits que mezclen múltiples objetivos no relacionados.

Antes del commit verificar cuando corresponda:

- archivos modificados;
- pruebas relevantes;
- artefactos accidentales;
- secretos;
- archivos temporales.

## 9. Mensajes de Commit

Cuando el usuario solicite commits, preferir mensajes claros y descriptivos.

Puede proponerse Conventional Commits:

feat:
fix:
refactor:
test:
docs:
chore:
perf:
build:
ci:

pero no imponerlo si el repositorio utiliza otra convención existente.

La convención actual del proyecto tiene prioridad.

## 10. Archivos Generados

No agregar accidentalmente al repositorio:

- binarios;
- builds;
- caches;
- archivos temporales;
- logs;
- dumps;
- secretos;
- credenciales;
- archivos locales de IDE;
- node_modules;
- .next;
- artefactos de pruebas;

salvo que exista una razón explícita.

Respetar .gitignore.

## 11. Secretos

Antes de preparar un commit revisar que no se incluyan:

- passwords;
- API keys;
- tokens;
- certificados privados;
- claves privadas;
- credenciales de base de datos;
- secretos de proveedores.

Nunca mover un secreto a Git solamente para facilitar el desarrollo.

## 12. Dependencias

Cambios en:

go.mod
go.sum
package.json
package-lock.json
pnpm-lock.yaml
yarn.lock

deben corresponder a una necesidad real de la tarea.

No actualizar dependencias de forma masiva durante una funcionalidad no relacionada.

## 13. Migraciones

Las migraciones de base de datos deben ser cambios explícitos y revisables.

No ejecutar automáticamente migraciones productivas como consecuencia de crear o modificar el archivo de migración.

Creación de migración y ejecución de migración son acciones distintas.

## 14. Configuración

No sobrescribir automáticamente archivos locales de configuración.

Distinguir entre:

- configuración versionada;
- configuración local;
- secretos;
- ejemplos de configuración.

Cuando corresponda utilizar archivos tipo:

.env.example

sin incluir valores secretos reales.

## 15. Cambios entre Capas

Cuando una funcionalidad modifique varias capas:

Database
→ Backend
→ Frontend/Mobile
→ QA

mantener coherencia de contratos.

No considerar una capa terminada si depende de un contrato todavía no acordado.

## 16. Cambios Concurrentes

Si dos agentes necesitan modificar el mismo archivo:

NO hacerlo simultáneamente de forma independiente.

fleet-tech-lead debe coordinar:

- orden;
- responsabilidad;
- contrato;
- integración.

## 17. Conflictos

Ante un conflicto:

no resolver automáticamente descartando uno de los lados.

Determinar:

- intención de cada cambio;
- responsable;
- contrato esperado;
- riesgo de regresión.

Escalar a fleet-tech-lead cuando afecte múltiples dominios.

## 18. Refactoring

No utilizar una solicitud pequeña como excusa para refactorizar grandes áreas.

Un refactor amplio debe tener:

- objetivo;
- alcance;
- beneficio;
- riesgos;
- pruebas;
- estrategia de rollback cuando corresponda.

## 19. Rollback

Los cambios importantes deben permitir identificar cómo volver a un estado funcional anterior.

Rollback no significa necesariamente Git reset.

Puede requerir:

- revertir código;
- revertir configuración;
- rollback de deployment;
- migración compensatoria;
- feature flag;

según el caso.

## 20. Estado del Repositorio

Antes de declarar terminado un trabajo que modificó archivos, revisar cuando sea posible:

- qué archivos cambiaron;
- qué archivos fueron creados;
- qué archivos fueron eliminados;
- si existen modificaciones inesperadas.

Informar claramente cualquier cambio fuera del alcance esperado.

## 21. Regla de Solo Lectura

Si el usuario solicita:

- auditoría;
- análisis;
- diagnóstico;
- revisión;
- explicación;

no modificar el repositorio.

Una tarea de solo lectura debe permanecer de solo lectura.

## 22. Prioridad de la Política Existente

Si el repositorio ya tiene una política Git documentada:

esa política tiene prioridad sobre convenciones genéricas propuestas por esta regla, salvo conflicto con una instrucción explícita del usuario.

No inventar Git Flow, trunk-based development u otra estrategia sin necesidad.

## 23. Cierre

Antes de finalizar una tarea con modificaciones:

confirmar:

- alcance respetado;
- cambios inesperados identificados;
- secretos no agregados;
- artefactos accidentales no agregados;
- dependencias justificadas;
- pruebas correspondientes consideradas;
- conflictos pendientes informados.

NO realizar commit, push, merge o deployment automáticamente.
