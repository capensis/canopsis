import { defineAsyncComponent } from 'vue';

import { MODALS } from '@/constants';

export default {
  [MODALS.confirmation]: defineAsyncComponent(() => import('./common/confirmation.vue')),
  [MODALS.confirmationPhrase]: defineAsyncComponent(
    () => import('./common/confirmation-phrase.vue'),
  ),
  [MODALS.duration]: defineAsyncComponent(() => import('./common/duration.vue')),
  [MODALS.rate]: defineAsyncComponent(() => import('./common/rate.vue')),
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
  [MODALS.createCommentTemplate]: defineAsyncComponent(
    () => import('./comment-template/create-comment-template.vue'),
  ),
  [MODALS.createStateSetting]: defineAsyncComponent(
    () => import('./state-setting/create-state-setting.vue'),
  ),
  [MODALS.createJunitStateSetting]: defineAsyncComponent(
    () => import('./state-setting/create-junit-state-setting.vue'),
  ),
  [MODALS.userInterface]: defineAsyncComponent(() => import('./admin/user-interface.vue')),
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
  [MODALS.textFieldEditor]: defineAsyncComponent(() => import('./common/text-field-editor.vue')),
  [MODALS.createView]: defineAsyncComponent(() => import('./view/create-view.vue')),
  [MODALS.selectView]: defineAsyncComponent(() => import('./view/select-view.vue')),
  [MODALS.shareView]: defineAsyncComponent(() => import('./view/share-view.vue')),
  [MODALS.createWidget]: defineAsyncComponent(() => import('./view/create-widget.vue')),
  [MODALS.createShareToken]: defineAsyncComponent(
    () => import('./share-token/create-share-token.vue'),
  ),
  [MODALS.createMaintenance]: defineAsyncComponent(
    () => import('./maintenance/create-maintenance.vue'),
  ),
  [MODALS.createPattern]: defineAsyncComponent(() => import('./pattern/create-pattern.vue')),
  [MODALS.createTheme]: defineAsyncComponent(() => import('./theme/create-theme.vue')),
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
  [MODALS.pbehaviorPlanning]: defineAsyncComponent(
    () => import('./pbehavior/pbehavior-planning.vue'),
  ),
  [MODALS.pbehaviorRecurrentChangesConfirmation]: defineAsyncComponent(
    () => import('./pbehavior/pbehavior-recurrent-changes-confirmation.vue'),
  ),
  [MODALS.createAlarmStatusRule]: defineAsyncComponent(
    () => import('./alarm-status-rule/create-alarm-status-rule.vue'),
  ),
  [MODALS.createIdleRule]: defineAsyncComponent(() => import('./idle-rule/create-idle-rule.vue')),
  [MODALS.createLinkRule]: defineAsyncComponent(() => import('./link-rule/create-link-rule.vue')),
  [MODALS.createDeclareTicketRule]: defineAsyncComponent(
    () => import('./declare-ticket/create-declare-ticket-rule.vue'),
  ),
  [MODALS.createDeclareTicketEvent]: defineAsyncComponent(
    () => import('./declare-ticket/create-declare-ticket-event.vue'),
  ),
  [MODALS.createAssociateTicketEvent]: defineAsyncComponent(
    () => import('./declare-ticket/create-associate-ticket-event.vue'),
  ),
  [MODALS.executeDeclareTickets]: defineAsyncComponent(
    () => import('./declare-ticket/execute-declare-tickets-rule.vue'),
  ),
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
  [MODALS.createAnomalyMonitoredConnector]: defineAsyncComponent(
    () => import('./anomaly-monitored-connector/create-anomaly-monitored-connector.vue'),
  ),
  [MODALS.createSnmpRule]: defineAsyncComponent(() => import('./snmp-rule/create-snmp-rule.vue')),
  [MODALS.createTemplateTestingData]: defineAsyncComponent(
    () => import('./template-testing/create-template-testing-data.vue'),
  ),
  [MODALS.createTemplateTestingTest]: defineAsyncComponent(
    () => import('./template-testing/create-template-testing-test.vue'),
  ),
  [MODALS.createExternalAuthToken]: defineAsyncComponent(
    () => import('./external-auth-token/create-external-auth-token.vue'),
  ),
  [MODALS.createExternalDataTable]: defineAsyncComponent(
    () => import('./external-data-table/create-external-data-table.vue'),
  ),
  [MODALS.createEntityInfoProperty]: defineAsyncComponent(
    () => import('./entity-info-property/create-entity-info-property.vue'),
  ),
  [MODALS.createIcon]: defineAsyncComponent(() => import('./icon/create-icon.vue')),
  [MODALS.createMap]: defineAsyncComponent(() => import('./map/create-map.vue')),
  [MODALS.createGeoMap]: defineAsyncComponent(() => import('./map/create-geo-map.vue')),
  [MODALS.createFlowchartMap]: defineAsyncComponent(() => import('./map/create-flowchart-map.vue')),
  [MODALS.createMermaidMap]: defineAsyncComponent(() => import('./map/create-mermaid-map.vue')),
  [MODALS.createTreeOfDependenciesMap]: defineAsyncComponent(
    () => import('./map/create-tree-of-dependencies-map.vue'),
  ),
  [MODALS.createTag]: defineAsyncComponent(() => import('./tag/create-tag.vue')),
  [MODALS.createLlm]: defineAsyncComponent(() => import('./llm/create-llm.vue')),
  [MODALS.importExportViews]: defineAsyncComponent(
    () => import('./view/import-groups-and-views.vue'),
  ),
  [MODALS.createExternalDataTableRecord]: defineAsyncComponent(
    () => import('./external-data-table/create-external-data-table-record.vue'),
  ),
  [MODALS.importExternalDataTableRecords]: defineAsyncComponent(
    () => import('./external-data-table/import-external-data-table-records.vue'),
  ),
  [MODALS.createRemediationConfiguration]: defineAsyncComponent(
    () => import('./remediation/create-remediation-configuration.vue'),
  ),
  [MODALS.createRemediationInstruction]: defineAsyncComponent(
    () => import('./remediation/create-remediation-instruction.vue'),
  ),
  [MODALS.createRemediationJob]: defineAsyncComponent(
    () => import('./remediation/create-remediation-job.vue'),
  ),
};
