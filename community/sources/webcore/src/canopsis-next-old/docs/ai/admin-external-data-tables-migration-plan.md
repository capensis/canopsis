# Admin External Data Tables Migration Plan

This file tracks the temporary migration plan for `/admin/external-data-tables`. Remove it after the page is fully migrated and verified.

## Source And Target

- Source page: `../canopsis-next/src/views/admin/custom-objects/external-data-tables.vue`
- Target page: `src/views/admin/custom-objects/external-data-tables-page.vue`
- Route: `ROUTES.adminCustomObjectsExternalDataTables`
- Access: `USER_PERMISSIONS.technical.externalDataTable`

## Migration Rules For This Page

- Preserve the source component style. The source page and most external data table components use Composition API, so keep Composition API for those files.
- Do not simplify the old behavior. Port the full dependency slice: page, list, expand panel, record list, import/export helpers, modals, forms, i18n, constants, hooks and store methods.
- Use `c-layout` as the default layout migration strategy for old `v-layout` / `v-flex` markup unless a real grid is needed.
- Keep migrated helper JSDoc and typedefs from the source project.
- Adapt only Vue 3, Vuetify 4, Pinia, router, i18n and form APIs.

## Implementation Plan

1. Store and hooks
   - `[x]` Wrap the external data table store with `createCRUDStore` for table CRUD and list loading.
   - `[x]` Preserve existing import and record actions.
   - `[x]` Add export actions from the source Vuex module: create export and fetch export status.
   - `[x]` Expose all table, import, export and record actions through `useExternalDataTable`.
   - `[x]` Add source-compatible import and record hook aliases.

2. Page and route
   - `[x]` Create `external-data-tables-page.vue` from the source page.
   - `[x]` Replace `.sync` with `:options` / `@update:options`.
   - `[x]` Wire the router route to the migrated page instead of `PlaceholderPage`.

3. Table list
   - `[x]` Port `external-data-tables-list.vue`.
   - `[x]` Preserve headers, linked-rules delete protection, removed-from-config warning tooltip and expand slot.
   - `[x]` Adapt tooltip activators and action layout to Vuetify 4/project wrappers.

4. Expand panel and records
   - `[x]` Port expand panel components used by the list.
   - `[x]` Port records table, data loading, search/options, row actions and bulk remove behavior.
   - `[x]` Preserve expanded row state across list refresh when the table still exists through `c-advanced-data-table`.
   - `[x]` Port column settings hooks and small density/grid controls used by the records table.

5. Modals and forms
   - `[x]` Port create/edit external data table modal.
   - `[x]` Port create/edit external data table record modal.
   - `[x]` Port import external data table records modal, upload section, preview and completion flow.
   - `[x]` Reuse already migrated external data table form fields where they are complete; port missing child components without reducing behavior.
   - `[x]` Add `c-array-text-field`, `c-csv-separator-field` and `file-drag-selector` dependencies.

6. Helpers, constants and i18n
   - `[x]` Port missing URL and record helpers with JSDoc and typedefs.
   - `[x]` Verify modal constants and modal registry entries.
   - `[x]` Compare English and French external data messages and port missing keys.
   - `[x]` Add file size/download helper support required by import/export.

7. Verification
   - `[x]` Run targeted lint for changed files.
   - `[x]` Run `npm run build`.
   - `[x]` Manually verify `/admin/external-data-tables`: list load, pagination, create/edit/delete, expand records, record create/edit/delete, import preview/complete and export/download.

## Out Of Scope

- `src/components/widgets/external-data-table/external-data-table-widget.vue` is tracked separately in the widget inventory and is not part of this page migration.
- Sidebars for widget configuration are out of scope unless the page directly requires them.
