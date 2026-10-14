# Admin Widget Templates Migration Plan

Source page: `/Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next/src/views/admin/settings/widget-templates.vue`

Target page: `src/views/admin/settings/widget-templates-page.vue`

Route: `/admin/widget-templates`

Status: `[x]` fully migrated and ready for testing.

## Scope

- `[x]` Audit old route, page, list, modals, form subtree, helper payload mapper, store/hook, constants and i18n.
- `[x]` Restore the target route so `/admin/widget-templates` renders the migrated page instead of a placeholder/wrong component.
- `[x]` Migrate the page shell with `c-page`, create FAB, refresh action and technical widget-template CRUD permissions.
- `[x]` Migrate the list with `c-advanced-data-table`, search, advanced pagination, formatted updated date, translated type and edit/delete actions.
- `[x]` Migrate `select-widget-template-type` modal with all legacy template types and default PDF content.
- `[x]` Migrate `create-widget-template` modal with `useForm`, dirty close confirmation and create/edit submit labels.
- `[x]` Migrate the widget-template form subtree:
  - name field;
  - columns template editor;
  - alarm quick actions and mass quick actions editor;
  - alarm sort columns editor;
  - text editor template with alarm/entity variables.
- `[x]` Restore full `widget/template/form.js` payload parity, including quick actions.
- `[x]` Register modals in the Vue 3 modal registry and keep existing modal sizing.
- `[x]` Verify widget settings consumers against the old `widgetTemplatesMixin` usages:
  `alarm`, `alarms-list-modal`, `stats-calendar`, `counter`, `context`, `availability`, `map`
  and `service-weather`.
- `[x]` Update `docs/migration-inventory.md` and `Ready For Testing`.

## Adaptation Notes

- Preserve Composition API in source components that already use Composition API.
- Do not reintroduce Vue 2 `v-field`, `.sync`, `$tc`, `$validator` or raw Vuetify form fields when target wrappers exist.
- `c-*` components remain global; feature components should not locally register them.
- `c-advanced-data-table` normalizes item slots, so feature code should use `item` directly and not `item.raw`.
- QA owns runtime verification; this plan only tracks migration completion and Ready For Testing handoff.
