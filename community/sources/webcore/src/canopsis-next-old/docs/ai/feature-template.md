# Feature Migration Template

Use this template when starting a new migrated page, modal, store module or complex shared component.

## Request

```text
Migrate: <feature name>
Old source:
- /Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next/src/<path>

Target scope:
- src/<path>/**

Do not touch:
- src/<path>/**
```

## Source Inventory

- Old route/component:
- Old child components:
- Old modals:
- Old store modules:
- Old hooks/mixins:
- Old helpers:
- Old i18n keys:
- Old tests:

## Target Ownership

- Page files:
- Feature components:
- Modal files:
- Store files:
- Hook files:
- Helper files:
- Common components needed:
- Docs to update:

## Migration Checklist

- `[ ]` Read old source component tree.
- `[ ]` Read old store/API logic.
- `[ ]` Check existing target patterns.
- `[ ]` Preserve Options API or Composition API style.
- `[ ]` Port every used dependency: child components, mixins, hooks, helpers, constants, typedefs, styles and i18n keys.
- `[ ]` Document any obsolete or incompatible dependency that is intentionally not ported.
- `[ ]` Port helpers before UI when possible.
- `[ ]` Port store module and hook before page uses it.
- `[ ]` Port leaf components before parent components.
- `[ ]` Port modal wrapper and register modal.
- `[ ]` Replace Vuetify 2-only APIs with Vuetify 4 equivalents.
- `[ ]` Use `c-layout` and standard gaps.
- `[ ]` Use `useForm` for form roots.
- `[ ]` Use store hooks from components.
- `[ ]` Update `docs/migration-inventory.md` and this page-specific plan after each completed migration stage.
- `[ ]` Run targeted eslint.
- `[ ]` Run `npm run build`.

## Form Checklist

- `[ ]` `entityToForm` helper exists.
- `[ ]` `formToEntity` helper exists.
- `[ ]` Root modal/page uses `useForm`.
- `[ ]` Subforms (`*-form.vue`) use `useFormContext`, `useFormModels` or `useFormModel`; no `:form` prop or `update:form` emit.
- `[ ]` Field components use props and `update:*` emits; no `useForm*` hooks inside leaf/composite `**/fields/**` components.
- `[ ]` Leaf fields use wrappers.
- `[ ]` Validation rules are explicit props or wrapper rules.
- `[ ]` Submit uses `form.handleSubmit`.
- `[ ]` Dirty close uses confirmable modal where needed.

## Store Checklist

- `[ ]` Pinia module exists under `src/stores/modules`.
- `[ ]` Store uses `request` / API routes (not the hook).
- `[ ]` Hook exists under `src/hooks/store/modules`.
- `[ ]` Hook delegates to the store; hook file does not import `request`.
- `[ ]` Hook exposes domain-specific refs and methods.
- `[ ]` Components do not import store or `request` directly.
- `[ ]` Old payload shape is preserved.

## Verification Notes

```text
Targeted eslint:

Build:

Known unrelated failures:

Runtime checks still needed:
```

## Final Response Template

```text
Migrated <feature>.

Changed:
- <file>: <short reason>
- <file>: <short reason>

Verification:
- <command> passes
- <command> passes

Not done / pending:
- <item>
```
