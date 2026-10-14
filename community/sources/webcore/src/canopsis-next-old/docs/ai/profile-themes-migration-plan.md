# Profile Themes Migration Plan

Source page: `../canopsis-next/src/views/profile/themes.vue`
Target page: `src/views/profile/themes-page.vue`

## Completed

- `[x]` Wire `/themes` to the migrated profile themes page instead of the placeholder.
- `[x]` Port the themes list with search, server pagination, select-all, mass delete and row actions.
- `[x]` Port create/edit/duplicate/delete/bulk-delete flows through the existing theme Pinia CRUD facade.
- `[x]` Port the create/edit/duplicate theme modal root with `useForm`, `useSubmittableForm` and confirmable close behavior.
- `[x]` Port the theme form and color preview components with Vue 3 `v-model` bindings and form context models.
- `[x]` Port the small color-picker menu field as a global `c-*` form field so the theme form keeps the old compact color-square control.
- `[x]` Add `formToTheme` conversion helpers, preserving old table color payload shape.

## Remaining Runtime Checks

- `[x]` Verify create/edit/duplicate saves against the backend.
- `[x]` Verify editing the active user theme updates the runtime theme immediately.
- `[x]` Verify default themes are disabled in edit/delete actions and row selection.
- `[x]` Verify color preview readability messages and font-size preview parity.
- `[x]` Verify the compact color picker menu visual parity against the old Vue 2 field.
