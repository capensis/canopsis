# Profile Patterns Migration Plan

Temporary working plan for migrating `/filters` (`profile-patterns`) from `../canopsis-next`.

## Rules

- Read old source only from `../canopsis-next`; edit only `canopsis-next-vue-3`.
- Preserve old pattern list, corporate tab and create/edit/delete flows.
- Reuse the already migrated `c-patterns-field` stack; do not simplify pattern form behavior.
- Do not port the full AI chat stack; keep the existing deferred placeholder behavior.
- Update `docs/migration-inventory.md` and this file after completed migration stages.

## Old Source

- `../canopsis-next/src/views/profile/patterns.vue`
- `../canopsis-next/src/components/other/pattern/patterns.vue`
- `../canopsis-next/src/components/other/pattern/corporate-patterns.vue`
- `../canopsis-next/src/components/other/pattern/patterns-list.vue`
- `../canopsis-next/src/components/modals/pattern/create-pattern.vue`
- `../canopsis-next/src/components/forms/pattern-form.vue`
- `../canopsis-next/src/helpers/entities/pattern/form.js`
- `../canopsis-next/src/store/modules/entities/pattern/**`
- `../canopsis-next/src/hooks/store/modules/pattern*.js`

## Target

- `src/views/profile/patterns-page.vue`
- `src/components/other/pattern/**`
- `src/components/modals/pattern/create-pattern.vue`
- `src/components/forms/pattern-form.vue`
- `src/router/index.js`
- `src/components/modals/index.js`
- `docs/migration-inventory.md`

## Already Present

- `src/helpers/entities/pattern/form.js`
- `src/stores/modules/pattern.js`
- `src/stores/modules/pattern-corporate.js`
- `src/hooks/store/modules/pattern.js`
- `src/hooks/store/modules/pattern-corporate.js`
- `src/components/forms/fields/pattern/**`
- `src/components/common/buttons/c-fab-expand-btn.vue`
- `src/components/common/buttons/c-action-fab-btn.vue`

## Checklist

- `[x]` Read migration rules and old page/list/modal source.
- `[x]` Confirm target stores/hooks/constants/shared pattern field stack are already present.
- `[x]` Port page shell and corporate tabs.
- `[x]` Port pattern list components.
- `[x]` Port create/edit pattern modal with `useForm`.
- `[x]` Port `pattern-form.vue` using form context.
- `[x]` Register modal and route component.
- `[x]` Update `docs/migration-inventory.md`.
- `[x]` Run targeted eslint.
- `[x]` Run `npm run build`.

## Runtime Checklist

- `[x]` Open `/filters`.
- `[x]` Verify personal patterns list/search/pagination/delete/bulk delete.
- `[x]` Verify corporate tab visibility and corporate patterns list/search/pagination/delete/bulk delete.
- `[x]` Verify create/edit flows for alarm/entity/pbehavior/service-weather patterns.
- `[x]` Verify pattern advanced editor and counters in the modal.
