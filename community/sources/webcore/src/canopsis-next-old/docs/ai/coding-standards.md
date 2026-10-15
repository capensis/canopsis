# Coding Standards

## Baseline

Follow the repository linters:

- `npm run lint`
- `npx eslint <files>` for targeted checks
- `npm run build` for compile verification

During migration, targeted checks are often preferred because the full lint may fail on unrelated existing files.

ESLint is the final formatting authority. `oxfmt` may run before lint for plain JS/TS files, but do not run it on Vue SFC files because it can rewrite templates in ways that conflict with migrated Vuetify markup.

Do not run `oxlint --fix` on Vue SFC files. `oxlint` may check `.vue` files,
but its fixer can collapse template attribute indentation. For `.vue` files,
use ESLint as the auto-fixer.

The installed UI dependency is Vuetify `^4.1.2`. When migrating old Vuetify 2 screens, use the current Vuetify 4 component API plus project wrappers/compat components from this repository; do not assume Vuetify 3 behavior.

## JavaScript

- Prefer JavaScript for new migration code unless the user asks for TypeScript.
- Keep semicolons.
- Do not use nested ternary expressions; use `if`/`else` or a `switch`/lookup map for multiple branches.
- Use single quotes for JavaScript strings. If a string contains a single quote, escape it instead of switching to double quotes.
- Omit parentheses for single-argument expression-body arrow functions, but keep parentheses when the arrow function uses a block body. `oxfmt` uses `arrowParens: avoid`, and ESLint with `requireForBlockBody` is the final fixer.
- Keep an empty line between import groups, enforced by `perfectionist/sort-imports`: builtin, external, internal `@/`, parent, sibling and index imports.
- Use named exports for shared `.js` modules.
- Avoid default exports in `.js` helpers.
- Avoid broad utility abstractions unless they remove real duplication.
- Prefer existing helpers over ad hoc parsing.
- For new form/API mapping helpers, use directional names such as `somethingToForm` and `formToSomething`; these names do not need a redundant `convert*` or `prepare*` prefix. Use `prepare*` or `convert*` when the function prepares/converts data without a directional `ToForm`/`FormTo...` name. Do not introduce new `normalize*`, `clamp*`, or similar names for payload preparation/bounds handling.
- Do not remove existing `console.log` calls unless the user explicitly asks for it. Treat them as user debugging markers, even when targeted lint reports `no-console`.
- When migrating helper files from `canopsis-next`, migrate their original JSDoc with the code, including function comments and all `@typedef` blocks, and update it if target signatures, types or behavior change.

## Prefer Simple Optional Access Over Verbose Guards

Avoid long guard chains when optional chaining and a fallback are enough.

```js
// Bad
if (item && typeof item === 'object' && keyField in item && item[keyField] !== undefined) {
  return item[keyField];
}

return index;

// Good
return item?.[keyField] ?? index;
```

```js
// Bad
if ('header' in item && item.header) {
  return item.header;
}

// Good
if (item?.header) {
  return item.header;
}
```

Use `??` instead of `||` when `0`, `''`, or `false` are valid values.

Only keep explicit type or shape checks when they encode real business rules, not when they only restate optional access.

- Preserve component API style from old source.
- Do not simplify migrated source code as a temporary compatibility shortcut. A port must keep the old component/helper behavior, state, events, slots, validation timing, menu/tooltip behavior, selection handling, and edge cases unless the user explicitly approves changing them.
- For widget settings migrations, port every original setting and its complete dependency slice, including child components, forms, modals, helpers, stores, constants, translations and styles. Do not replace missing functionality with TODOs, stubs, disabled controls or visual placeholders.
- Adapt incompatible Vue 2/Vuetify 2 APIs to Vue 3/Vuetify 4, but do not replace a complex old component with a reduced approximation just because it compiles.
- When porting code from `canopsis-next`, port all accompanying dependencies used by that code: mixins, hooks, helpers, constants, typedefs, child components, styles and i18n keys. Do not leave required support logic behind just because the top-level component compiles.
- For non-UI infrastructure files, port whole files first instead of copying individual functions. This applies to helpers, services, constants, store facades, domain hooks, form/query mappers and utility modules. After the full-file port, adapt only target-stack API differences and remove only obsolete/incompatible or explicitly deferred code with a documented reason.
- Skip only support code that is obsolete or incompatible with the new stack. In that case, document or report exactly what was not ported and the reason.
- If the exact behavior depends on missing infrastructure, port or identify that infrastructure first; otherwise document the blocker and ask before reducing the behavior.
- Do not rewrite Options API components to Composition API unless required by a library migration.
- Do not rewrite Composition API components to Options API.
- Do not add component `name` unless recursion requires it.
- In Vue SFC files, keep top-level blocks in this order: `<template>`, then `<script>`, then `<style>`. The `<template>` block must be first whenever it exists.
- Keep prop contracts compatible with old usage when porting old components.
- Keep Vue `props` declarations multi-line, matching the old project style. Do not compress prop definitions into one-line objects like `modelValue: { type: String, default: undefined }`.
- Use `modelValue` / `update:modelValue` for standard leaf fields; use named `update:<prop>` for subform/object/array components.
- Do not add new `v-field` or `@input` component bindings.

## Replace `v-field` And `@input`

Cross-cutting migration rule for form bindings ported from `canopsis-next`:

| Old | Target |
| --- | --- |
| `v-field="form.field"` on a leaf/native field | `:model-value="form.field"` + `@update:model-value="updateField('field', $event)"`, or `:value` + `@update:value` for compat wrappers |
| `v-field="form"` on a child subform (`*-form.vue`) | `useFormModels()` / `useFormContext()` + `v-model` / `v-model:<prop>`; do not pass `:form` or emit `@update:form` |
| `v-field="items[index]"` on an array item | `:item="items[index]"` + `@update:item="updateItemInArray(index, $event)"` |
| `v-field="items"` on `c-form-block-array-field` | `:form="items"` + `@update:form="updateModel"` |
| `model: { prop: 'form', event: 'input' }` | `emits: ['update:form']` |
| `@input="handler"` on a Vue component | `@update:<prop>="handler"` matching the child model prop |
| `emit('input', value)` | `emit('update:<prop>', value)` |

Implementation helpers:

- Subforms (`*-form.vue` inside a `useForm()` tree): `useFormContext()`, `useFormModels()` or `useFormModel(path)`; bind fields with `v-model` / `v-model:<prop>`. Do not accept a `form` prop or emit `update:form`.
- Object subforms (legacy, outside form context): `useModelField(props, emit)` → `updateModel`, `updateField`. Migrate to form hooks when the parent uses `useForm()`.
- Array subforms: `useArrayModelField(props, emit)` → `updateModel`, `updateItemInArray`, `removeItemFromArray`.
- vee-validate form roots: `useForm()` + `form.model()` / `useFormModels()` on bound fields.
- Field components under `src/components/**/fields/` (leaf/composite inputs, not `*-form` subforms): props + `update:*` emits only; no `useForm*` hooks.

Exceptions:

- Native DOM `@input` on raw `<input>` elements stays as-is.
- Do not restore the old compile-time `v-field` directive.

## Use `c-*` Field Wrappers

In migrated app code, do not use raw Vuetify form controls when a project wrapper exists.

All `c-*` components are global components. If a `c-*` component is missing in a template, export it from the closest `index.js` barrel and register that barrel in `src/main.js`; do not import or register `c-*` components locally in feature components.

Available wrappers in `src/components/common/vuetify/`:

- `c-text-field`
- `c-select`
- `c-combobox`
- `c-textarea`
- `c-radio-group`
- `c-checkbox`
- `c-autocomplete`

Compat/domain wrappers may compose them (`c-name-field`, `c-payload-text-field`, etc.).

| Avoid | Use |
| --- | --- |
| `v-text-field` | `c-text-field` |
| `v-select` | `c-select` |
| `v-combobox` | `c-combobox` |
| `v-textarea` | `c-textarea` |
| `v-radio-group` | `c-radio-group` |

Prefer wrapper validation props (`required`, `name`, …) over `v-validate` + manual `errors.collect()` on the same control.

Use `item-title` / `item-value` with Vuetify 4 selects.

Project wrappers must normalize Vuetify 4 slot item wrappers before exposing
slots to feature components. In feature/business components that use `c-select`,
`c-combobox`, `c-autocomplete`, `c-lazy-search-field`, `c-advanced-data-table`,
or another project wrapper, treat slot `item` as the original item and do not
read `item.raw`. Accessing `item.raw` is allowed only inside wrapper source
files or in intentionally raw Vuetify components that have no project wrapper.

`c-select` binding:

```vue
<c-select
  :model-value="value"
  :items="items"
  item-title="text"
  item-value="value"
  @update:model-value="updateModel"
/>
```

Do not use the removed `c-select-field` compat wrapper.

Exceptions: wrapper source files themselves, and specialized low-level components (`time-picker-field`, `c-search-field`) that intentionally bind Vuetify APIs directly.

## Use `c-tooltip`

In application templates, always use the global `c-tooltip` wrapper. Do not use
raw `v-tooltip` and do not restore the removed `c-simple-tooltip` component.
Only `src/components/common/tooltip/c-tooltip.vue` may render `v-tooltip`
directly.

`c-tooltip` forwards Vuetify tooltip attributes and slots, preserves the legacy
`content`, `html`, `top`, `right`, `bottom` and `left` convenience props, and
forces `eager: false`. Keeping that lazy behavior centralized prevents every
tooltip from creating an overlay subtree before it is opened.

```vue
<c-tooltip location="top">
  <template #activator="{ props }">
    <v-btn v-bind="props" icon="info" />
  </template>

  Tooltip text
</c-tooltip>
```

`c-tooltip` is globally registered; do not import or register it locally in
feature components. In component tests, stub `c-tooltip` rather than
`v-tooltip`; stub `v-tooltip` only when testing the wrapper itself.

## Payload Textarea Parity

When the old project used `c-payload-textarea-field`, the migrated Vue 3 code must use `c-payload-textarea-field` too. Do not replace it with `c-textarea`, `text-editor-field`, `c-payload-text-field`, or a plain `c-text-field`: the payload textarea owns multiline template variables, `{{ ... }}` variable menu behavior, line-level payload errors and the simple-template editing UX.

Old-project `c-payload-textarea-field` usage inventory:

| Old source | Target status |
| --- | --- |
| `components/other/action/form/action-note-form.vue` | Migrated: `src/components/other/action/form/action-note-form.vue` |
| `components/forms/fields/c-change-state-field.vue` | Migrated: `src/components/forms/fields/c-change-state-field.vue`; the Note/output field always uses `c-payload-textarea-field` |
| `components/forms/request/request-form.vue` | Migrated: `src/components/forms/request/request-form.vue` |
| `components/forms/fields/c-json-field.vue` | Migrated: `src/components/forms/fields/compat/c-json-field.vue` |
| `components/other/remediation/jobs/form/remediation-job-general-form.vue` | Migrated in the new job form: `src/components/other/remediation/jobs/form/remediation-job-form.vue` |
| `components/other/template-testing/test-variables/partials/template-testing-test-variables-template-item.vue` | Migrated: `src/components/other/template-testing/test-variables/partials/template-testing-test-variables-template-item.vue` |
| `components/modals/common/payload-textarea-editor.vue` | Migrated: `src/components/modals/common/payload-textarea-editor.vue`; keep `MODALS.payloadTextareaEditor` for simple column/template editors |
| `components/other/meta-alarm-rule/form/meta-alarm-rule-general-form.vue` | Migrated: `src/components/other/meta-alarm-rule/form/meta-alarm-rule-general-form.vue`; `output_template` uses `c-payload-textarea-field` with `templateVars.output` |

## Vuetify 4 Migration

Common mappings:

- `value` -> `model-value`
- `@input` -> `@update:model-value` for native Vuetify components
- component `@input` -> `@update:<prop>` for custom/migrated components
- `v-tabs-items` / `v-tab-item` -> `v-window` / `v-window-item`
- `append-icon` events may become `append-inner-icon` and `click:append-inner`
- `v-subheader` is not available; use `v-list-subheader`
- Do not use `c-layout` in migrated form/layout templates. Use Vuetify `v-row`/`v-col` instead. When an existing `c-layout` is encountered in code being edited, rewrite that layout to the equivalent `v-row`/`v-col` structure while preserving its alignment, spacing and responsive behavior.
- Do not use Vuetify `v-layout` in migrated code. Old `v-layout` behavior must be migrated to `v-row`/`v-col` or plain `<div>` where appropriate.
- Old `v-flex` behavior is emulated with `v-flex` compat, `v-row`/`v-col`, or plain `<div>` where appropriate.

### Layout Gaps

Use one gap mechanism per layout element.

- For Vuetify layouts, use `v-row`/`v-col` and their spacing props/classes.
- Do not use old `gap-*` classes in migrated templates. Convert the layout to `v-row`/`v-col` and preserve the intended spacing.
- Do not introduce `c-layout` in new or migrated form/layout code.

### Field Icon And Affix Slots

Vuetify field affix slots changed meaning between Vuetify 2 and Vuetify 4. Do not copy old slot names mechanically.

| Vuetify 2 source | Vuetify 4 target | Notes |
| --- | --- | --- |
| `#append` | `#append-inner` | Use for help icons, menu activators and icons that must live inside the input field. |
| `#append-outer` | `#append` | Use only when the old icon/control was intentionally outside the field outline/underline. |
| `#prepend-inner` | `#prepend-inner` | Still inside the input field. |
| `#prepend` | `#prepend` | Outside the input field. |
| `append-icon` / `@click:append` | `append-inner-icon` / `@click:append-inner` | Use for native Vuetify fields when the icon belongs inside the field. |

Before migrating any field with a help icon, compare its visual placement in the old project. If the icon was inside the field, it must remain inside the field in the target wrapper (`c-text-field`, `c-combobox`, `c-select`, `c-textarea`, payload fields, etc.).

### Dotted `v-slot` Names

Vue parses `#item.foo` as slot `item` with modifier `foo`, not as slot name `item.foo`. `vue/valid-v-slot` allows simple modifiers (`#item.actions`) via `allowModifiers: true`, but hyphenated names such as `#item.row-icons` fail with `'v-slot' directive doesn't support any modifier`.

For Vuetify data-table / `c-advanced-data-table` slots whose names contain a dot, always use dynamic slot syntax. Do not disable `vue/valid-v-slot` to work around this.

| Avoid | Use |
| --- | --- |
| `#item.row-icons` | `` #[`item.row-icons`] `` |
| `#header.row-icons` | `` #[`header.row-icons`] `` |
| `#item.foo-bar` | `` #[`item.${column.value}`] `` |

```vue
<!-- Bad -->
<template #item.row-icons="{ item }">

<!-- Good -->
<template #[`item.row-icons`]="{ item }">

<!-- Good: dynamic column keys -->
<template
  v-for="column in columns"
  :key="`item.${column.value}`"
  #[`item.${column.value}`]="{ item }"
>
```

Simple unhyphenated Vuetify slots such as `#item.actions` may keep the shorthand. Prefer the dynamic form whenever the slot name includes a hyphen or is built from a column key.

### `v-alert` visibility

In Vuetify 2, `v-alert` used `:value` / `v-model` to show and hide the alert (often with a transition).

In Vuetify 4, do not bind `:value`, `:model-value`, or `v-model` on `v-alert` or `c-alert` for conditional display. Put `v-if` on the alert itself instead.

| Old | Target |
| --- | --- |
| `<v-alert :value="hasError" type="error">` | `<v-alert v-if="hasError" type="error">` |
| `<c-alert :value="!!message" type="error">` | `<c-alert v-if="message" type="error">` |
| `<v-alert v-model="visible" closable>` | `<v-alert v-if="visible" closable @click:close="visible = false">` |

For closable alerts that need a leave transition, wrap with `v-fade-transition` (or another transition) around the `v-if` alert.

Do not pass a `value` prop to `c-alert`; the wrapper renders the alert as-is and callers control visibility with `v-if`.

Do not assume Vuetify 2 labels, slots or menu behavior still work. Verify component APIs when porting.

## `v-bind` And Listener Objects For Project Components

When binding **project components** (`c-*`, field wrappers, dynamic `<component>` configs), do not split props and listeners into paired `v-bind` + `v-on`.

Merge listeners into one bind object and use Vue 3 listener prop names.

| Old | Target |
| --- | --- |
| `v-bind="config.props" v-on="config.on"` | `v-bind="config.props"` only |
| `{ bind, on }` config objects | one object, usually still named `bind` or `props` |
| `'update:value': handler` in listener objects | `'onUpdate:value': handler` |
| `'update:modelValue': handler` | `'onUpdate:modelValue': handler` |
| `'update:duration': handler` | `'onUpdate:duration': handler` |
| `input: handler` on project components | `'onUpdate:modelValue': handler` and/or `'onUpdate:value': handler` matching the child emit |
| `click: handler` | `onClick: handler` |
| `close: handler` | `onClose: handler` |
| `'click:close': handler` | `'onClick:close': handler` |
| `refresh: handler` on `c-fab-btn` | `onRefresh: handler` |

```vue
<!-- Bad -->
<component
  :is="valueComponent.is"
  v-bind="valueComponent.props"
  v-on="valueComponent.on"
/>

<!-- Good -->
<component
  :is="valueComponent.is"
  v-bind="valueComponent.props"
/>
```

```js
// Bad
return {
  is: 'c-number-field',
  props: { modelValue: value, name },
  on: { 'update:modelValue': updateValue },
};

// Good
return {
  is: 'c-number-field',
  props: {
    modelValue: value,
    name,
    'onUpdate:modelValue': updateValue,
  },
};
```

Rules:

- Prefer explicit `@update:<prop>` in templates when the binding is static.
- Use merged `v-bind` objects only for dynamic `<component>` configs, chip/menu wrappers and similar listener maps.
- Do not keep separate `on` keys in component config objects.
- Native DOM listeners (`@input` on `<input>`, `@keydown`, …) stay in templates; this rule is for **Vue project components** only.

## Forms And Validation

Use the local form layer:

```js
const form = useForm({
  initialValues: entityToForm(entity),
});

const name = form.model('name');

const submit = form.handleSubmit(() => {
  save(formToEntity(form.values.value));
});
```

Rules:

- Form values live in `form.values`.
- vee-validate stores validation state only.
- `field wrappers -> useField() -> validation/meta/errors`.
- Form roots use `useForm()`.
- Components using `useForm()`, `useFormContext()`, `useFormModels()` or `useFormModel()` are Composition API components. Do not mix an Options API component body with `setup()` only to install a form context.
- Subforms (`*-form.vue`) use `useFormContext()`, `useFormModels()` or `useFormModel()`. Bind child fields with `v-model` / `v-model:<prop>`. Do not pass `form` by props or return changes through `update:form` events.
- Field components under `src/components/**/fields/` (leaf/composite inputs, not `*-form` subforms) receive model props and emit `update:*`; do not use `useForm()`, `useFormContext()`, `useFormModels()` or `useFormModel()` inside them.
- Do not use `useFieldArray`.
- Do not use `cloneDeep` in the form layer.
- Do not call `setFieldValue` to mirror form values into vee-validate.

## Date And Time Fields

Date/time fields were ported from old `canopsis-next` and adapted to Vuetify 4.

- `date-time-picker-field` works with `Date` values.
- `date-time-picker-text-field` works with formatted date strings.
- `time-picker-field` uses a separate combobox search state.
- `c-duration-field` keeps the old value/unit shape and Options API surface, but uses `c-number-field` and `c-select`; replace old `$tc` unit labels with Vue I18n `t(key, count)` while preserving pluralization.
- Invalid time search text must not update the model.
- Date values should be converted with helpers from `src/helpers/date/date.js`.
- Time/date picker-specific helpers live in `src/helpers/date/date-time-picker.js`.

## Stores

Use Pinia modules and store hooks:

```
Component / composable
  -> src/hooks/store/modules/<entity>.js
    -> src/stores/modules/<entity>.js
      -> src/services/request
```

Rules:

- **Store modules own all server requests.** Only files under `src/stores/modules/**` (and shared store plugins under `src/stores/plugins/**`) may call `request` or API routes directly.
- **Hooks wrap stores.** Files under `src/hooks/store/modules/**` expose domain-specific refs and methods (`usersMeta`, `fetchAlarmTagsListWithoutStore`, `createUser`). They must not import `request`.
- **Components use hooks only.** Do not import Pinia stores or `request` from components, field wrappers, or page logic. Call methods from `useUser()`, `useAlarmTag()`, etc.
- Preserve old API payload shapes when porting stores.
- Prefer `createCRUDStore()` for list/CRUD entities unless the old module has important custom behavior.

Example:

```js
// src/stores/modules/alarm-tag.js
export const useAlarmTagStore = createCRUDStore({
  id: 'alarmTag',
  route: API_ROUTES.alarmTag.list,
  withWithoutStore: true,
});

// src/hooks/store/modules/alarm-tag.js
export const useAlarmTag = () => {
  const alarmTagStore = useAlarmTagStore();

  return {
    fetchAlarmTagsListWithoutStore: alarmTagStore.fetchListWithoutStore,
  };
};

// component setup()
const { fetchAlarmTagsListWithoutStore } = useAlarmTag();
```

When adding a new entity:

1. Create `src/stores/modules/<entity>.js` first.
2. Create `src/hooks/store/modules/<entity>.js` that wraps the store.
3. Use the hook from components.

See also `docs/ai/architecture.md` (Store Architecture) and `docs/ai/feature-template.md` (Store Checklist).

## Tables

`c-advanced-data-table` is a compatibility layer over Vuetify 4 tables.

- Preserve old slot names and old header shapes where possible.
- Keep pagination/search/sort behavior compatible with old consumers.
- Table rows and headers should have Vuetify 2-like bottom borders.
- Sticky headers and nested table spacing should be checked for each migrated page.

## CSS And Layout

- Prefer scoped styles unless global compatibility styles are intentional.
- Use existing Vuetify 2 compatibility SCSS for migrated classes.
- Do not create one-off utility classes if a compatibility class already exists.
- Keep fields dense but readable.
- Avoid UI text overlap across mobile and desktop.

## Verification

Before finishing a coding task:

- Run a targeted eslint command for changed files.
- Run `npm run build` for compile checks when changes affect app code.
- Mention if full `npm run lint` was not run or fails on unrelated files.
- Do not start the dev server unless the user explicitly asks.
