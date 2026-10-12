# AI Chat Migration Plan

Source project: `../canopsis-next`
Target project: `src`

This plan tracks the deferred AI chat / LLM chat migration. The feature must be
ported as one coherent dependency slice, not as isolated fixes inside individual
forms.

## Migration Rules

- Port the whole old AI chat stack before enabling it in consumers:
  - `hooks/ai/ai-chat.js`;
  - `hooks/ai/ai-chat-form.js`;
  - `helpers/entities/llm/chat/**`;
  - `components/other/llm/chat/**`;
  - `components/modals/llm/ai-chat-history.vue`.
- Keep existing migrated forms working while the stack is deferred. Temporary
  placeholders may keep imports/injections safe, but must not fetch LLMs, open
  sidebars, join sockets or apply generated patterns.
- Preserve old behavior first, then adapt only Vue 3, Vuetify 4, Pinia, router,
  socket and form APIs.
- AI chat sidebar belongs to the host form modal lifecycle. Opening nested or
  parallel modals, confirmations, help popups, selectors, histories or other UI
  layers must not close, hide, reset or unmount the chat; it fully closes only
  when the form modal that created it is closed/unmounted.
- Do not mix Options API and Composition API in one component during this pass.
  Old Composition components should remain Composition; old Options components
  should remain Options unless there is a clear local reason to rewrite.
- Keep helper typedefs/JSDoc from the old project.
- Do not migrate this feature piecemeal through consumers. Consumers should be
  re-enabled after the shared AI chat stack is complete.

## Source Inventory

- `../canopsis-next/src/hooks/ai/ai-chat.js`
- `../canopsis-next/src/hooks/ai/ai-chat-form.js`
- `../canopsis-next/src/helpers/entities/llm/chat/form.js`
- `../canopsis-next/src/helpers/entities/llm/chat/message/entity.js`
- `../canopsis-next/src/components/other/llm/chat/ai-chat.vue`
- `../canopsis-next/src/components/other/llm/chat/ai-chat-greeting.vue`
- `../canopsis-next/src/components/other/llm/chat/ai-chat-info-alert.vue`
- `../canopsis-next/src/components/other/llm/chat/ai-chat-messages-list.vue`
- `../canopsis-next/src/components/other/llm/chat/ai-chat-pattern-item-field.vue`
- `../canopsis-next/src/components/other/llm/chat/ai-chat-sidebar-wrapper.vue`
- `../canopsis-next/src/components/other/llm/chat/ai-chat-sidebar.vue`
- `../canopsis-next/src/components/other/llm/chat/ai-chat-suggestions.vue`
- `../canopsis-next/src/components/other/llm/chat/ai-chat-textarea.vue`
- `../canopsis-next/src/components/other/llm/chat/ai-chat-warning-alert.vue`
- `../canopsis-next/src/components/other/llm/chat/partials/ai-chat-llm-field.vue`
- `../canopsis-next/src/components/other/llm/chat/partials/ai-chat-message.vue`
- `../canopsis-next/src/components/other/llm/chat/partials/ai-chat-pattern.vue`
- `../canopsis-next/src/components/modals/llm/ai-chat-history.vue`

## Current Target State

- `src/hooks/ai/ai-chat-form.js` contains the migrated AI chat form
  integration.
- `src/hooks/ai/ai-chat.js` contains the migrated chat/session/socket logic.
- `src/components/other/llm/chat/**` contains the migrated chat sidebar,
  messages, model selector, suggestions and pattern version UI.
- `src/components/modals/llm/ai-chat-history.vue` is registered and reachable
  from `/admin/llms` prompt history.
- LLM administration is already migrated:
  - `/admin/llms`;
  - `components/modals/llm/create-llm.vue`;
  - `components/other/llm/form/**`;
  - `stores/modules/llm.js`;
  - `hooks/store/modules/llm.js`.
- `MODALS.aiChatHistory`, base LLM constants and i18n keys already exist.

## Consumers To Re-Enable

These files currently import or depend on AI chat entry points and must be
checked after the shared stack is restored:

- `src/components/modals/remediation/create-remediation-instruction.vue`
- `src/components/modals/link-rule/create-link-rule.vue`
- `src/components/modals/event-filter/create-event-filter.vue`
- `src/components/modals/pbehavior/create-pbehavior.vue`
- `src/components/modals/pattern/create-pattern.vue`
- `src/components/modals/state-setting/create-state-setting.vue`
- `src/components/modals/meta-alarm-rule/create-meta-alarm-rule.vue`
- `src/components/modals/declare-ticket/create-declare-ticket-rule.vue`
- `src/components/other/pbehavior/calendar/partials/pbehavior-create-event.vue`
- `src/components/common/block/c-form-general-patterns-tabs.vue`
- `src/components/forms/fields/pattern/pattern-fix-with-ai-btn.vue`
- `src/components/forms/fields/pattern/hooks/patterns-field-expand.js`
- `src/components/other/link-rule/form/link-rule-form.vue`
- `/admin/llms` AI chat history replay modal/session behavior.

## Implementation Plan

### 1. Constants, API And Store Surface

- `[x]` Compare old and target `constants/llm.js`, `constants/modal.js`,
  `config.js` socket/API route constants and permission constants.
- `[x]` Add or align `SOCKET_ROOMS.llmChat`, LLM chat API routes and
  `MODALS.aiChatHistory` registration.
- `[x]` Extend `stores/modules/llm.js` and `hooks/store/modules/llm.js` with the
  old chat-specific actions:
  - fetch enabled LLMs for chat;
  - create/send chat request;
  - fetch history list;
  - fetch history item/session;
  - bulk link chat history to rule ids.
- `[x]` Preserve the already migrated LLM CRUD/admin behavior.

### 2. Chat Helpers

- `[x]` Port `helpers/entities/llm/chat/form.js` completely.
- `[x]` Port `helpers/entities/llm/chat/message/entity.js` completely.
- `[x]` Adapt imports to target helper paths without changing behavior.
- `[x]` Preserve pattern conversion behavior for all supported contexts and
  fields.

### 3. Core Hooks

- `[x]` Replace placeholder `hooks/ai/ai-chat.js` with the full old hook adapted
  to Vue 3/Pinia/socket plugin.
- `[x]` Replace placeholder `hooks/ai/ai-chat-form.js` with the full old form
  integration hook.
- `[x]` Preserve:
  - expand callback registration;
  - pattern item mapping;
  - minimized state in local storage;
  - pending/cancel state;
  - JSON push into chat;
  - form pattern sync;
  - LLM availability check;
  - chat id collection and bulk link after submit;
  - `$aiChat` provide contract.
- `[x]` Keep existing consumer API shape so migrated forms do not need local
  rewrites beyond necessary bindings.

### 4. Chat UI Components

- `[x]` Port `components/other/llm/chat/**` completely.
- `[x]` Replace old Vuetify 2 layout/field syntax with target wrappers and
  Vuetify 4 APIs.
- `[x]` Preserve:
  - greeting;
  - info/warning alerts;
  - suggestions;
  - model selector;
  - textarea ask/cancel behavior;
  - messages list;
  - thinking/error states;
  - pattern versions;
  - restore version;
  - see pattern JSON;
  - restart confirmation;
  - minimized sidebar mode.

### 5. History Modal

- `[x]` Port `components/modals/llm/ai-chat-history.vue`.
- `[x]` Register `MODALS.aiChatHistory` in `components/modals/index.js`.
- `[x]` Wire `/admin/llms` prompts history actions to open the migrated history
  modal.
- `[x]` Preserve grouped chat/session display and replay behavior from the old
  project.

### 6. Consumer Re-Enable Pass

- `[x]` Re-check every consumer listed above against the restored hook/sidebar
  API.
- `[x]` Ensure forms with pattern tabs expand to the correct tab/panel when AI
  chat requests focus.
- `[x]` Ensure generated patterns update the correct scoped form fields.
- `[x]` Ensure modal close/submit flows link collected chat histories to created
  or updated rule ids where the old project did.
- `[x]` Ensure disabled/without-link modes still work where consumers pass those
  options.

### 7. Storage Settings

- `[x]` Verify `storage-settings-llm-chat-form.vue` and
  `helpers/entities/data-storage/form.js` still match old AI chat history
  retention settings.
- `[x]` Confirm storage settings save payload includes `llm_chat.delete_after`
  exactly as before.

### 8. Verification

- `[x]` Run targeted eslint for all migrated AI chat files and touched
  consumers.
- `[x]` Run `npm run build`.
- `[x]` Runtime check with enabled LLM:
  - AI chat appears in every consumer that should show it;
  - model selection works;
  - suggestions populate the prompt;
  - send/cancel/restart works;
  - socket room joins/leaves cleanly;
  - generated patterns are applied to the correct tab;
  - `Fix with AI` sends current JSON to chat;
  - submit links chat histories to saved rules;
  - `/admin/llms` history modal opens and renders sessions.
- `[x]` Runtime check without enabled LLM:
  - chat stays hidden or shows old no-model behavior;
  - forms remain usable;
  - no console warnings/errors from `$aiChat` injections.

## Completion Criteria

- Placeholders are removed or replaced by full implementations.
- All consumers listed in this plan are re-enabled.
- `Deferred AI Chat Migration` in `docs/migration-inventory.md` is replaced with
  a completed entry pointing to this plan.
- `/admin/llms` history replay/session behavior is ready for QA.
