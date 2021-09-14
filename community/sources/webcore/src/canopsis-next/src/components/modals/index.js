import CModalLoaderOverlay from '@/components/common/overlay/modal-loader-overlay.vue';

export const CreateEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-event.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateAckEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-ack-event.vue'),
  loading: CModalLoaderOverlay,
});
export const ConfirmAckWithTicket = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/confirm-ack-with-ticket.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateAssociateTicketEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-associate-ticket-event.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateChangeStateEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-change-state-event.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateDeclareTicketEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-declare-ticket-event.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateSnoozeEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-snooze-event.vue'),
  loading: CModalLoaderOverlay,
});
export const VariablesHelp = () => ({
  component: import(/* webpackChunkName: "Events" */ './common/variables-help.vue'),
  loading: CModalLoaderOverlay,
});
export const InfoPopupSetting = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/info-popup-setting/info-popup-setting.vue'),
  loading: CModalLoaderOverlay,
});
export const AddInfoPopup = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/info-popup-setting/add-info-popup.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateManualMetaAlarm = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/create-manual-meta-alarm.vue'),
  loading: CModalLoaderOverlay,
});
export const PbehaviorList = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-list.vue'),
  loading: CModalLoaderOverlay,
});
export const EditLiveReporting = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/edit-live-reporting.vue'),
  loading: CModalLoaderOverlay,
});
export const Confirmation = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/confirmation.vue'),
  loading: CModalLoaderOverlay,
});
export const ClickOutsideConfirmation = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/click-outside-confirmation.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateWidget = () => ({
  component: import(/* webpackChunkName: "Widget" */ './view/create-widget.vue'),
  loading: CModalLoaderOverlay,
});
export const ColorPicker = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/color-picker.vue'),
  loading: CModalLoaderOverlay,
});
export const TextEditor = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/text-editor.vue'),
  loading: CModalLoaderOverlay,
});
export const TextFieldEditor = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/text-field-editor.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateService = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/create-service.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateEntity = () => ({
  component: import(/* webpackChunkName: "Context" */ './entity/create-entity.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateEntityInfo = () => ({
  component: import(/* webpackChunkName: "Context" */ './entity/create-entity-info.vue'),
  loading: CModalLoaderOverlay,
});
export const SelectView = () => ({
  component: import(/* webpackChunkName: "Views" */ './view/select-view.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateView = () => ({
  component: import(/* webpackChunkName: "Views" */ './view/create-view.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateGroup = () => ({
  component: import(/* webpackChunkName: "Views" */ './view/create-group.vue'),
  loading: CModalLoaderOverlay,
});
export const ImportGroupsAndViews = () => ({
  component: import(/* webpackChunkName: "Views" */ './view/import-groups-and-views.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateFilter = () => ({
  component: import(/* webpackChunkName: "Filters" */ './common/create-filter.vue'),
  loading: CModalLoaderOverlay,
});
export const ServiceEntities = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/service-entities.vue'),
  loading: CModalLoaderOverlay,
});
export const ServiceDependencies = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/service-dependencies.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateServicePauseEvent = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/create-service-pause-event.vue'),
  loading: CModalLoaderOverlay,
});
export const AddStat = () => ({
  component: import(/* webpackChunkName: "Stats" */ './stats/add-stat.vue'),
  loading: CModalLoaderOverlay,
});
export const StatsDateInterval = () => ({
  component: import(/* webpackChunkName: "Stats" */ './stats/stats-date-interval.vue'),
  loading: CModalLoaderOverlay,
});
export const StatsDisplayMode = () => ({
  component: import(/* webpackChunkName: "Stats" */ './stats/stats-display-mode.vue'),
  loading: CModalLoaderOverlay,
});
export const AlarmsList = () => ({
  component: import(/* webpackChunkName: "Alarms" */ './alarm/alarms-list.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateUser = () => ({
  component: import(/* webpackChunkName: "Users" */ './admin/create-user.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateRole = () => ({
  component: import(/* webpackChunkName: "Roles" */ './admin/create-role.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateBroadcastMessage = () => ({
  component: import(/* webpackChunkName: "BroadcastMessages" */ './admin/create-broadcast-message.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateEventFilterRule = () => ({
  component: import(/* webpackChunkName: "EventFilters" */ './event-filter/create-event-filter-rule.vue'),
  loading: CModalLoaderOverlay,
});
export const CreatePattern = () => ({
  component: import(/* webpackChunkName: "Patterns" */ './pattern/create-pattern.vue'),
  loading: CModalLoaderOverlay,
});
export const CreatePatternRule = () => ({
  component: import(/* webpackChunkName: "Patterns" */ './pattern/create-pattern-rule.vue'),
  loading: CModalLoaderOverlay,
});
export const EventFilterRuleActions = () => ({
  component: import(/* webpackChunkName: "EventFilters" */ './event-filter/event-filter-rule-actions.vue'),
  loading: CModalLoaderOverlay,
});
export const EventFilterRuleExternalData = () => ({
  component: import(/* webpackChunkName: "EventFilters" */ './event-filter/event-filter-rule-external-data.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateEventFilterRuleAction = () => ({
  component: import(/* webpackChunkName: "EventFilters" */ './event-filter/create-event-filter-rule-action.vue'),
  loading: CModalLoaderOverlay,
});
export const FiltersList = () => ({
  component: import(/* webpackChunkName: "Filters" */ './common/filters-list.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateSnmpRule = () => ({
  component: import(/* webpackChunkName: "Snmp" */ './snmp-rule/create-snmp-rule.vue'),
  loading: CModalLoaderOverlay,
});
export const SelectViewTab = () => ({
  component: import(/* webpackChunkName: "Views" */ './view/select-view-tab.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateDynamicInfo = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/create-dynamic-info.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateDynamicInfoInformation = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/create-dynamic-info-information.vue'),
  loading: CModalLoaderOverlay,
});
export const DynamicInfoTemplatesList = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/dynamic-info-templates-list.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateDynamicInfoTemplate = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/create-dynamic-info-template.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateMetaAlarmRule = () => ({
  component: import(/* webpackChunkName: "Alarms" */ './meta-alarm-rule/create-meta-alarm-rule.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateCommentEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './common/create-comment-event.vue'),
  loading: CModalLoaderOverlay,
});
export const CreatePlaylist = () => ({
  component: import(/* webpackChunkName: "Playlist" */ './admin/create-playlist.vue'),
  loading: CModalLoaderOverlay,
});
export const ManagePlaylistTabs = () => ({
  component: import(/* webpackChunkName: "Playlist" */ './admin/manage-playlist-tabs.vue'),
  loading: CModalLoaderOverlay,
});
export const PbehaviorPlanning = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-planning.vue'),
  loading: CModalLoaderOverlay,
});
export const SelectExceptionsLists = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/select-exceptions-lists.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateRRule = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-r-rule.vue'),
  loading: CModalLoaderOverlay,
});
export const PbehaviorRecurrentChangesConfirmation = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-recurrent-changes-confirmation.vue'),
  loading: CModalLoaderOverlay,
});
export const CreatePbehavior = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior.vue'),
  loading: CModalLoaderOverlay,
});
export const CreatePbehaviorType = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior-type.vue'),
  loading: CModalLoaderOverlay,
});
export const CreatePbehaviorReason = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior-reason.vue'),
  loading: CModalLoaderOverlay,
});
export const CreatePbehaviorException = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior-exception.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateRemediationInstruction = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-instruction.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateRemediationConfiguration = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-configuration.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateRemediationJob = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-job.vue'),
  loading: CModalLoaderOverlay,
});
export const ExecuteRemediationInstruction = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/execute-remediation-instruction.vue'),
  loading: CModalLoaderOverlay,
});
export const RemediationPatterns = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/remediation-patterns.vue'),
  loading: CModalLoaderOverlay,
});
export const RemediationInstructionApproval = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/remediation-instruction-approval.vue'),
  loading: CModalLoaderOverlay,
});
export const ImageViewer = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/image-viewer.vue'),
  loading: CModalLoaderOverlay,
});
export const ImagesViewer = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/images-viewer.vue'),
  loading: CModalLoaderOverlay,
});
export const Patterns = () => ({
  component: import(/* webpackChunkName: "Patterns" */ './pattern/patterns.vue'),
  loading: CModalLoaderOverlay,
});
export const Rate = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './common/rate.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateRemediationInstructionsFilter = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-instructions-filter.vue'),
  loading: CModalLoaderOverlay,
});
export const TestSuite = () => ({
  component: import(/* webpackChunkName: "Junit" */ './test-suite/test-suite.vue'),
  loading: CModalLoaderOverlay,
});
export const StateSetting = () => ({
  component: import(/* webpackChunkName: "Parameters" */ './state-setting/state-setting.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateIdleRule = () => ({
  component: import(/* webpackChunkName: "IdleRules" */ './idle-rule/create-idle-rule.vue'),
  loading: CModalLoaderOverlay,
});
export const HealthcheckEngine = () => ({
  component: import(/* webpackChunkName: "Healthcheck" */ './healthcheck/healthcheck-engine.vue'),
  loading: CModalLoaderOverlay,
});
export const HealthcheckEnginesChainReference = () => ({
  component: import(/* webpackChunkName: "Healthcheck" */ './healthcheck/healthcheck-engines-chain-reference.vue'),
  loading: CModalLoaderOverlay,
});
export const CreateScenario = () => ({
  component: import(/* webpackChunkName: "Scenarios" */ './scenario/create-scenario.vue'),
  loading: CModalLoaderOverlay,
});
