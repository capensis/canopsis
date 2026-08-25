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
  [MODALS.stateSettingInheritedEntityPattern]: defineAsyncComponent(
    () => import('./state-setting/state-setting-inherited-entity-pattern.vue'),
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
  [MODALS.payloadTextareaEditor]: defineAsyncComponent(
    () => import('./common/payload-textarea-editor.vue'),
  ),
  [MODALS.textEditorWithTemplate]: defineAsyncComponent(
    () => import('./common/text-editor-with-template.vue'),
  ),
  [MODALS.infoPopupSetting]: defineAsyncComponent(
    () => import('./alarm/info-popup-setting/info-popup-setting.vue'),
  ),
  [MODALS.addInfoPopup]: defineAsyncComponent(
    () => import('./alarm/info-popup-setting/add-info-popup.vue'),
  ),
  [MODALS.textFieldEditor]: defineAsyncComponent(() => import('./common/text-field-editor.vue')),
  [MODALS.createFilter]: defineAsyncComponent(() => import('./common/create-filter.vue')),
  [MODALS.filtersList]: defineAsyncComponent(() => import('./common/filters-list.vue')),
  [MODALS.variablesHelp]: defineAsyncComponent(() => import('./common/variables-help.vue')),
  [MODALS.editLiveReporting]: defineAsyncComponent(() => import('./alarm/edit-live-reporting.vue')),
  [MODALS.createAlarmChart]: defineAsyncComponent(() => import('./alarm/create-alarm-chart.vue')),
  [MODALS.createView]: defineAsyncComponent(() => import('./view/create-view.vue')),
  [MODALS.createGroup]: defineAsyncComponent(() => import('./view/create-group.vue')),
  [MODALS.selectView]: defineAsyncComponent(() => import('./view/select-view.vue')),
  [MODALS.selectViewTab]: defineAsyncComponent(() => import('./view/select-view-tab.vue')),
  [MODALS.shareView]: defineAsyncComponent(() => import('./view/share-view.vue')),
  [MODALS.createWidget]: defineAsyncComponent(() => import('./view/create-widget.vue')),
  [MODALS.selectWidgetTemplateType]: defineAsyncComponent(
    () => import('./widget-template/select-widget-template-type.vue'),
  ),
  [MODALS.createWidgetTemplate]: defineAsyncComponent(
    () => import('./widget-template/create-widget-template.vue'),
  ),
  [MODALS.createShareToken]: defineAsyncComponent(
    () => import('./share-token/create-share-token.vue'),
  ),
  [MODALS.serviceDependencies]: defineAsyncComponent(
    () => import('./service/service-dependencies.vue'),
  ),
  [MODALS.serviceEntities]: defineAsyncComponent(() => import('./service/service-entities.vue')),
  [MODALS.createServicePauseEvent]: defineAsyncComponent(
    () => import('./service/create-service-pause-event.vue'),
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
  [MODALS.pbehaviorsCalendar]: defineAsyncComponent(
    () => import('./pbehavior/pbehaviors-calendar.vue'),
  ),
  [MODALS.pbehaviorRecurrentChangesConfirmation]: defineAsyncComponent(
    () => import('./pbehavior/pbehavior-recurrent-changes-confirmation.vue'),
  ),
  [MODALS.createAlarmStatusRule]: defineAsyncComponent(
    () => import('./alarm-status-rule/create-alarm-status-rule.vue'),
  ),
  [MODALS.createIdleRule]: defineAsyncComponent(() => import('./idle-rule/create-idle-rule.vue')),
  [MODALS.createLinkRule]: defineAsyncComponent(() => import('./link-rule/create-link-rule.vue')),
  [MODALS.createScenario]: defineAsyncComponent(() => import('./scenario/create-scenario.vue')),
  [MODALS.createDeclareTicketRule]: defineAsyncComponent(
    () => import('./declare-ticket/create-declare-ticket-rule.vue'),
  ),
  [MODALS.createDeclareTicketEvent]: defineAsyncComponent(
    () => import('./declare-ticket/create-declare-ticket-event.vue'),
  ),
  [MODALS.createAssociateTicketEvent]: defineAsyncComponent(
    () => import('./declare-ticket/create-associate-ticket-event.vue'),
  ),
  [MODALS.removeAssociatedTicketEvent]: defineAsyncComponent(
    () => import('./alarm/create-remove-associated-ticket-event.vue'),
  ),
  [MODALS.executeDeclareTickets]: defineAsyncComponent(
    () => import('./declare-ticket/execute-declare-tickets-rule.vue'),
  ),
  [MODALS.createEvent]: defineAsyncComponent(() => import('./alarm/create-event.vue')),
  [MODALS.createAckEvent]: defineAsyncComponent(() => import('./alarm/create-ack-event.vue')),
  [MODALS.createSnoozeEvent]: defineAsyncComponent(() => import('./alarm/create-snooze-event.vue')),
  [MODALS.createChangeStateEvent]: defineAsyncComponent(
    () => import('./alarm/create-change-state-event.vue'),
  ),
  [MODALS.createCommentEvent]: defineAsyncComponent(
    () => import('./common/create-comment-event.vue'),
  ),
  [MODALS.linkToMetaAlarm]: defineAsyncComponent(() => import('./alarm/link-to-meta-alarm.vue')),
  [MODALS.removeAlarmsFromMetaAlarm]: defineAsyncComponent(
    () => import('./alarm/remove-alarms-from-meta-alarm.vue'),
  ),
  [MODALS.alarmsList]: defineAsyncComponent(() => import('./alarm/alarms-list.vue')),
  [MODALS.executeRemediationInstruction]: defineAsyncComponent(
    () => import('./remediation/execute-remediation-instruction.vue'),
  ),
  [MODALS.executeRemediationSimpleInstruction]: defineAsyncComponent(
    () => import('./remediation/execute-remediation-simple-instruction.vue'),
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
  [MODALS.anomalyMonitoredConnectorHistory]: defineAsyncComponent(
    () => import('./anomaly-monitored-connector/anomaly-monitored-connector-history.vue'),
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
  [MODALS.createEntity]: defineAsyncComponent(() => import('./entity/create-entity.vue')),
  [MODALS.createEntityInfo]: defineAsyncComponent(() => import('./entity/create-entity-info.vue')),
  [MODALS.entitiesList]: defineAsyncComponent(() => import('./entity/entities-list.vue')),
  [MODALS.entitiesRootCauseDiagram]: defineAsyncComponent(
    () => import('./entity/entities-root-cause-diagram.vue'),
  ),
  [MODALS.createService]: defineAsyncComponent(() => import('./service/create-service.vue')),
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
  [MODALS.remediationInstructionApproval]: defineAsyncComponent(
    () => import('./remediation/remediation-instruction-approval.vue'),
  ),
  [MODALS.createRemediationJob]: defineAsyncComponent(
    () => import('./remediation/create-remediation-job.vue'),
  ),
  [MODALS.testSuite]: defineAsyncComponent(() => import('./test-suite/test-suite.vue')),
};
