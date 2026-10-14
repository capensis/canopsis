# Architecture

## Project Context

This repository is a new Vue 3 rewrite of the Canopsis `canopsis-next` frontend.

- New project: `/Users/verokhov/Workspace/canopsis-pro-vue-3/canopsis`
- Source project: `/Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next`
- Migration inventory: `docs/migration-inventory.md`

The current strategy is not an in-place Vue 2 to Vue 3 upgrade. It is a new application that ports the old product feature by feature while preserving proven business logic and UI behavior.

## Current Stack

- Vue 3.5
- Vuetify 4.1 (`package.json`: `^4.1.2`)
- Pinia 3
- Vue Router 5
- vee-validate 4
- vue-i18n 11
- Vite 8
- Vitest and Playwright
- ESLint 10, Oxlint, Oxfmt
- npm

Prefer JavaScript during the active migration. Do not convert files to TypeScript unless a task explicitly asks for it or the file already needs TypeScript-specific behavior.

## Main Runtime Shape

- `src/main.js`: application bootstrap, global component registration, plugins, modals, socket, Vuetify aliases.
- `src/app.vue`: application shell, app-level state, navigation visibility, theme/system logic.
- `src/router`: route records, access flags, redirects and placeholder pages.
- `src/plugins/vuetify`: Vuetify 4 theme, defaults, custom icons and aliases.
- `src/plugins/modals`: modal registry and runtime modal wrapper.
- `src/plugins/socket`: socket plugin and services.

## Source Layout

- `src/components/common`: reusable UI primitives and compatibility components.
- `src/components/forms`: form blocks and field components.
- `src/components/modals`: concrete modal implementations.
- `src/components/other`: feature-level components ported from the old app.
- `src/components/layout`: navigation, top bar, side bars and app layout.
- `src/views`: route page components. Page component filenames should use the `-page` suffix.
- `src/stores`: Pinia modules and store plugins.
- `src/hooks/store/modules`: public store-facing hooks. Components should use these, not stores directly.
- `src/hooks/form`: the new form layer over vee-validate.
- `src/hooks/query`: query state, query options and query mixins.
- `src/helpers/entities`: entity-specific mappers and form/query helpers.
- `src/helpers/date`: date and duration helpers.
- `src/i18n/messages`: translation messages.
- `docs`: migration tracking and AI artifacts.

## Migration Architecture

The old source of truth is `canopsis-next`. When porting a feature:

1. Locate the old page/component/store/helper in the source project.
2. Copy the old logic structure as much as possible.
3. Adapt only what is required by Vue 3, Vuetify 4, Pinia, vee-validate 4 and the new app shell.
4. Port the whole dependency slice used by the migrated code: child components, mixins, hooks, helpers, constants, typedefs, i18n keys, styles and related utility logic.
5. Skip only dependencies that are obsolete or incompatible with the new stack, and explicitly report what was skipped and why.
6. Avoid simplified interim components unless the user explicitly asks for a placeholder.
7. When finishing any migration stage, update migration documentation in the same change: `docs/migration-inventory.md` for page/widget/shared-component state, the `Ready For Testing` section at the top of that file when the area becomes runtime-testable, and the related `docs/ai/*-migration-plan.md` checklist when one exists.

For non-UI infrastructure files, prefer full-file ports instead of extracting
small fragments. Helpers, services, constants, store facades, domain hooks,
form/query mappers and utility modules from `canopsis-next` should be copied as
complete files first, preserving JSDoc and typedefs, then adapted only for
target-stack API differences. Do not cherry-pick individual helper functions
from a source file unless the rest of the file is obsolete/incompatible or
belongs to an explicitly deferred feature; document that reason when it happens.

## Deferred AI Chat Migration

The AI chat/LLM UI layer is intentionally deferred until the end of the migration.

Until that final pass:

- Do not port `src/hooks/ai/ai-chat.js`, full `src/hooks/ai/ai-chat-form.js`, `src/components/other/llm/chat/**`, or `src/helpers/entities/llm/chat/**` as part of ordinary page/form migrations.
- Keep only explicit compatibility placeholders when already migrated forms import AI chat entry points.
- The placeholders must keep forms compilable and keep `$aiChat` injections safe, but they must not fetch LLMs, open sidebars, join sockets or apply generated patterns.
- When the AI chat migration starts, replace the placeholders with a full dependency-slice port from `../canopsis-next` and adapt it as one coherent feature.

Current placeholders:

- `src/hooks/ai/ai-chat-form.js`
- `src/components/other/llm/chat/ai-chat-sidebar.vue`

## Component Strategy

Do not rewrite components just because they are old.

- If the old component uses Options API and it can still work, keep Options API.
- If the old component uses Composition API, keep Composition API.
- Use `<script setup>` only when the existing target component already uses it or when creating a small wrapper where that is consistent with nearby code.
- Do not set the component `name` option unless recursive component behavior requires it.

All `c-*` components must be registered globally from index barrels and used in templates without local imports or local `components` registration. When a migrated `c-*` component is missing, add or update its barrel export and global registration in `src/main.js` instead of importing it locally.

## UI Architecture

Vuetify 4 is the UI base, but the target app intentionally carries Vuetify 2 compatibility where old screens depend on it. Do not call the target stack Vuetify 3; check existing Vuetify 4 wrappers and nearby migrated components before adapting old markup.

- Use `c-layout` instead of raw `d-flex flex-column`.
- Vertical field spacing is usually `gap="2"`.
- Horizontal field spacing is usually `gap="3"`.
- Avoid nested card-in-card layouts.
- Preserve old UI behavior and visual density unless the user asks to redesign.
- Tables should preserve Vuetify 2-like borders, row height, hover and sticky header behavior.
- Top bar height is 48px.
- Tooltips and menus should use fade transitions by default.
- Button and tab text should be uppercase.

## Form Architecture

`src/hooks/form/index.js` is the project form layer.

The form layer separates:

- local form values as the only source of truth;
- vee-validate as validation state and submit engine only.

Pattern:

```js
const form = useForm({
  initialValues: entityToForm(entity),
});

const titleField = form.model('title');

const submit = form.handleSubmit(() => {
  save(formToEntity(form.values.value));
});
```

Subforms (`*-form.vue` components inside a `useForm()` tree) should use form hooks, not prop/event bridging:

```js
const {
  externalDataField,
  configField,
} = useFormModels({
  externalDataField: 'external_data',
  configField: 'config',
});
```

```vue
<external-data-form v-model:form="externalDataField" />
<c-payload-text-field v-model:value="resourceField" />
```

Do not pass `:form="values"` and `@update:form="updateModel"` from the parent. The subform participates in the shared form context provided by the root `useForm()`.

For array fields:
```js
const form = useFormContext();
const itemsField = form.array('items');
```

For single fields:
```js
const form = useFormContext();
const titleField = form.model('title');
```

Field components under `src/components/**/fields/` that are leaf or composite inputs (not `*-form` subforms) must not call `useForm()`, `useFormContext()`, `useFormModels()` or `useFormModel()`. The parent subform binds models with `v-model` / `v-model:<prop>` and the field emits matching `update:*` events.

Example:

```vue
<pbehavior-recurrence-rule-field
  v-model:rrule="rruleField"
  v-model:exdates="exdatesField"
  v-model:exceptions="exceptionsField"
  :start="startField"
/>
```

Leaf field wrappers can use `useField()` for validation and errors, but their value should still come through `v-model`, `v-model:<prop>`, or explicit `:prop` + `@update:<prop>`. Do not use `v-field` or component `@input`.

Do not use `reactive`, `cloneDeep` or `useFieldArray` inside the form layer.

Do not pass ref or reactive value of form to `initialValue` of useField (to avoid mutation by `vee-validate` problem)

## Store Architecture

Pinia stores live under `src/stores`. Components should not import stores directly unless a task explicitly requires low-level store work. Use `src/hooks/store/modules/*` as the public interface.

Request flow:

```text
.vue component / field / page
  -> useEntity() in src/hooks/store/modules/entity.js
    -> useEntityStore() in src/stores/modules/entity.js
      -> request.get/post/... in src/services/request
```

Rules:

- Hooks do not call `request`/API directly.
- Store modules own server requests.
- Hooks wrap stores and expose domain-specific names such as `usersMeta`, `usersPending`, `createUser`, `bulkEnableUser`.
- CRUD modules should follow the current Pinia CRUD plugin shape unless the old Vuex module has important custom behavior.
- For associative table entities, preserve the old API shape: `{ name, content }` payloads and `content` response extraction.
- When migrating an entity, add `src/stores/modules/<name>.js` first, then `src/hooks/store/modules/<name>.js`, then wire components to the hook.

## Routing Architecture

Router migration is tracked in `docs/migration-inventory.md`.

- Missing pages should use placeholders until migrated.
- Preserve old redirect and permission rules.
- Route title rendering in top bar should use translated metadata, not raw route names.
- When routes are added, update inventory and page blocks.

## Known Active Migration Areas

- Keep `normalizr` out of the target migration path; use Pinia stores and explicit mapper helpers instead of restoring the old Vuex/normalizr entity plugin.
- Continue porting pages and modal forms.
- Continue deep parity for `c-advanced-data-table`, advanced search and `c-patterns-field`.
- Migrate the deferred AI chat/LLM UI layer at the end, after core pages/forms are stable.
- Keep date/time picker behavior aligned with old UX while using Vuetify 4 APIs.
- Avoid touching groups sidebar drag-and-drop unless the user explicitly asks.
