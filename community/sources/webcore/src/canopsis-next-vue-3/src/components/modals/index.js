import { defineAsyncComponent } from 'vue'

import { MODALS } from '@/constants'

export default {
  [MODALS.confirmation]: defineAsyncComponent(() => import('./common/confirmation.vue')),
  [MODALS.clickOutsideConfirmation]: defineAsyncComponent(
    () => import('./common/click-outside-confirmation.vue'),
  ),
  [MODALS.createUser]: defineAsyncComponent(() => import('./admin/create-user.vue')),
  [MODALS.createRole]: defineAsyncComponent(() => import('./admin/create-role.vue')),
  [MODALS.createPlaylist]: defineAsyncComponent(() => import('./admin/create-playlist.vue')),
  [MODALS.managePlaylistTabs]: defineAsyncComponent(
    () => import('./admin/manage-playlist-tabs.vue'),
  ),
  [MODALS.createBroadcastMessage]: defineAsyncComponent(
    () => import('./admin/create-broadcast-message.vue'),
  ),
  [MODALS.createDynamicInfo]: defineAsyncComponent(
    () => import('./dynamic-info/create-dynamic-info.vue'),
  ),
  [MODALS.createEventFilter]: defineAsyncComponent(
    () => import('./event-filter/create-event-filter.vue'),
  ),
  [MODALS.createDynamicInfoTemplate]: defineAsyncComponent(
    () => import('./dynamic-info/create-dynamic-info-template.vue'),
  ),
  [MODALS.addDynamicInfoInfosFromTemplate]: defineAsyncComponent(
    () => import('./dynamic-info/add-dynamic-info-infos-from-template.vue'),
  ),
  [MODALS.imageViewer]: defineAsyncComponent(() => import('./common/image-viewer.vue')),
  [MODALS.imagesViewer]: defineAsyncComponent(() => import('./common/images-viewer.vue')),
  [MODALS.textEditor]: defineAsyncComponent(() => import('./common/text-editor.vue')),
  [MODALS.textEditorWithTemplate]: defineAsyncComponent(
    () => import('./common/text-editor-with-template.vue'),
  ),
  [MODALS.createMaintenance]: defineAsyncComponent(
    () => import('./maintenance/create-maintenance.vue'),
  ),
  [MODALS.createPbehaviorType]: defineAsyncComponent(
    () => import('./pbehavior/create-pbehavior-type.vue'),
  ),
  [MODALS.createPbehaviorReason]: defineAsyncComponent(
    () => import('./pbehavior/create-pbehavior-reason.vue'),
  ),
  [MODALS.createPbehaviorException]: defineAsyncComponent(
    () => import('./pbehavior/create-pbehavior-exception.vue'),
  ),
  [MODALS.importPbehaviorException]: defineAsyncComponent(
    () => import('./pbehavior/import-pbehavior-exception.vue'),
  ),
  [MODALS.createAlarmStatusRule]: defineAsyncComponent(
    () => import('./alarm-status-rule/create-alarm-status-rule.vue'),
  ),
  [MODALS.createIdleRule]: defineAsyncComponent(() => import('./idle-rule/create-idle-rule.vue')),
  [MODALS.healthcheckEngine]: defineAsyncComponent(
    () => import('./healthcheck/healthcheck-engine.vue'),
  ),
  [MODALS.healthcheckEnginesChainReference]: defineAsyncComponent(
    () => import('./healthcheck/healthcheck-engines-chain-reference.vue'),
  ),
  [MODALS.createKpiFilter]: defineAsyncComponent(() => import('./kpi/create-kpi-filter.vue')),
  [MODALS.eventsRecord]: defineAsyncComponent(() => import('./events-record/events-record.vue')),
  [MODALS.applyEventFilter]: defineAsyncComponent(
    () => import('./events-record/apply-event-filter.vue'),
  ),
  [MODALS.createTemplateTestingData]: defineAsyncComponent(
    () => import('./template-testing/create-template-testing-data.vue'),
  ),
  [MODALS.createTemplateTestingTest]: defineAsyncComponent(
    () => import('./template-testing/create-template-testing-test.vue'),
  ),
}
