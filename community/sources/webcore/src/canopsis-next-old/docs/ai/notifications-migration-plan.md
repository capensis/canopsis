# Notifications Page Migration Plan

Source page: `../canopsis-next/src/views/notification/notifications.vue`
Target page: `src/views/notification/notifications-page.vue`

## Completed

- `[x]` Wire `/notifications/:tabId` to the migrated page instead of the placeholder.
- `[x]` Port the route page shell with page header, tabs, tab badges, refresh FAB and notifications socket refresh.
- `[x]` Port the three notification tabs:
  - instructions to approve;
  - instructions to rate;
  - event filter failures.
- `[x]` Reuse existing remediation instruction, remediation instruction stats and event filter Pinia hooks; no direct request calls in components.
- `[x]` Port `useNotificationActiveId` for opening the referenced notification item from the top-bar notification query.
- `[x]` Port instruction rating dependencies used by the notifications page: `rate` modal, `rate` form and `rating-field`.
- `[x]` Port the remediation instruction stats list enough for notifications: main stats list, interval toolbar, summary/statistics/rating expand tabs, affect-state chips and state-count chips.
- `[x]` Port remediation instruction stats alarms timeline expand tab for notifications, including executions list, horizontal alarm timeline, failed-executions switch and per-execution steps expand panel.

## Runtime Checks

- `[x]` Verify `/notifications/instructions-to-approve` list load, approval/edit action from active notification and pagination.
- `[x]` Verify `/notifications/instructions-to-rate` list load, quick interval changes, rating modal submit and expand tabs.
- `[x]` Verify `/notifications/event-filter-failures` list load, active notification row expansion and errors tab selection.
- `[x]` Verify top-bar “see all” navigation and socket refresh update all three counters.
