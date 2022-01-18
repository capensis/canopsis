import ModalLoaderOverlay from '@/components/common/overlay/modal-loader-overlay.vue';

export const CreateEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-event.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateAckEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-ack-event.vue'),
  loading: ModalLoaderOverlay,
});
export const ConfirmAckWithTicket = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/confirm-ack-with-ticket.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateAssociateTicketEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-associate-ticket-event.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateChangeStateEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-change-state-event.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateDeclareTicketEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-declare-ticket-event.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateSnoozeEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-snooze-event.vue'),
  loading: ModalLoaderOverlay,
});
export const VariablesHelp = () => ({
  component: import(/* webpackChunkName: "Events" */ './common/variables-help.vue'),
  loading: ModalLoaderOverlay,
});
export const EditLiveReporting = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/edit-live-reporting.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateCommentEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './common/create-comment-event.vue'),
  loading: ModalLoaderOverlay,
});
export const InfoPopupSetting = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/info-popup-setting/info-popup-setting.vue'),
  loading: ModalLoaderOverlay,
});
export const AddInfoPopup = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/info-popup-setting/add-info-popup.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateManualMetaAlarm = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/create-manual-meta-alarm.vue'),
  loading: ModalLoaderOverlay,
});
export const AlarmsList = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/alarms-list.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateMetaAlarmRule = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './meta-alarm-rule/create-meta-alarm-rule.vue'),
  loading: ModalLoaderOverlay,
});

export const PbehaviorList = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-list.vue'),
  loading: ModalLoaderOverlay,
});
export const PbehaviorPlanning = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-planning.vue'),
  loading: ModalLoaderOverlay,
});
export const PbehaviorRecurrenceRule = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-recurrence-rule.vue'),
  loading: ModalLoaderOverlay,
});
export const SelectExceptionsLists = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/select-exceptions-lists.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateRecurrenceRule = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-recurrence-rule.vue'),
  loading: ModalLoaderOverlay,
});
export const PbehaviorRecurrentChangesConfirmation = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-recurrent-changes-confirmation.vue'),
  loading: ModalLoaderOverlay,
});
export const CreatePbehavior = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior.vue'),
  loading: ModalLoaderOverlay,
});
export const CreatePbehaviorType = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior-type.vue'),
  loading: ModalLoaderOverlay,
});
export const CreatePbehaviorReason = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior-reason.vue'),
  loading: ModalLoaderOverlay,
});
export const CreatePbehaviorException = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior-exception.vue'),
  loading: ModalLoaderOverlay,
});
export const Confirmation = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/confirmation.vue'),
  loading: ModalLoaderOverlay,
});
export const ClickOutsideConfirmation = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/click-outside-confirmation.vue'),
  loading: ModalLoaderOverlay,
});
export const ColorPicker = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/color-picker.vue'),
  loading: ModalLoaderOverlay,
});
export const TextEditor = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/text-editor.vue'),
  loading: ModalLoaderOverlay,
});
export const TextFieldEditor = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/text-field-editor.vue'),
  loading: ModalLoaderOverlay,
});
export const ConfirmationPhrase = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/confirmation-phrase.vue'),
  loading: ModalLoaderOverlay,
});
export const ImageViewer = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/image-viewer.vue'),
  loading: ModalLoaderOverlay,
});
export const ImagesViewer = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/images-viewer.vue'),
  loading: ModalLoaderOverlay,
});
export const Info = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/info.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateWidget = () => ({
  component: import(/* webpackChunkName: "Widget" */ './view/create-widget.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateService = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/create-service.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateEntity = () => ({
  component: import(/* webpackChunkName: "Context" */ './entity/create-entity.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateEntityInfo = () => ({
  component: import(/* webpackChunkName: "Context" */ './entity/create-entity-info.vue'),
  loading: ModalLoaderOverlay,
});
export const ServiceEntities = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/service-entities.vue'),
  loading: ModalLoaderOverlay,
});
export const ServiceDependencies = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/service-dependencies.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateServicePauseEvent = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/create-service-pause-event.vue'),
  loading: ModalLoaderOverlay,
});
export const SelectView = () => ({
  component: import(/* webpackChunkName: "View" */ './view/select-view.vue'),
  loading: ModalLoaderOverlay,
});
export const SelectViewTab = () => ({
  component: import(/* webpackChunkName: "View" */ './view/select-view-tab.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateView = () => ({
  component: import(/* webpackChunkName: "View" */ './view/create-view.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateGroup = () => ({
  component: import(/* webpackChunkName: "View" */ './view/create-group.vue'),
  loading: ModalLoaderOverlay,
});
export const ImportGroupsAndViews = () => ({
  component: import(/* webpackChunkName: "View" */ './view/import-groups-and-views.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateFilter = () => ({
  component: import(/* webpackChunkName: "Filters" */ './common/create-filter.vue'),
  loading: ModalLoaderOverlay,
});
export const FiltersList = () => ({
  component: import(/* webpackChunkName: "Filters" */ './common/filters-list.vue'),
  loading: ModalLoaderOverlay,
});
export const AddStat = () => ({
  component: import(/* webpackChunkName: "Stats" */ './stats/add-stat.vue'),
  loading: ModalLoaderOverlay,
});
export const StatsDateInterval = () => ({
  component: import(/* webpackChunkName: "Stats" */ './stats/stats-date-interval.vue'),
  loading: ModalLoaderOverlay,
});
export const StatsDisplayMode = () => ({
  component: import(/* webpackChunkName: "Stats" */ './stats/stats-display-mode.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateUser = () => ({
  component: import(/* webpackChunkName: "User" */ './admin/create-user.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateRole = () => ({
  component: import(/* webpackChunkName: "Role" */ './admin/create-role.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateBroadcastMessage = () => ({
  component: import(/* webpackChunkName: "BroadcastMessage" */ './admin/create-broadcast-message.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateEventFilterRule = () => ({
  component: import(/* webpackChunkName: "EventFilter" */ './event-filter/create-event-filter-rule.vue'),
  loading: ModalLoaderOverlay,
});
export const EventFilterRuleActions = () => ({
  component: import(/* webpackChunkName: "EventFilter" */ './event-filter/event-filter-rule-actions.vue'),
  loading: ModalLoaderOverlay,
});
export const EventFilterRuleExternalData = () => ({
  component: import(/* webpackChunkName: "EventFilter" */ './event-filter/event-filter-rule-external-data.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateEventFilterRuleAction = () => ({
  component: import(/* webpackChunkName: "EventFilter" */ './event-filter/create-event-filter-rule-action.vue'),
  loading: ModalLoaderOverlay,
});
export const CreatePattern = () => ({
  component: import(/* webpackChunkName: "Pattern" */ './pattern/create-pattern.vue'),
  loading: ModalLoaderOverlay,
});
export const CreatePatternRule = () => ({
  component: import(/* webpackChunkName: "Pattern" */ './pattern/create-pattern-rule.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateSnmpRule = () => ({
  component: import(/* webpackChunkName: "SnmpRule" */ './snmp-rule/create-snmp-rule.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateDynamicInfo = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/create-dynamic-info.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateDynamicInfoInformation = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/create-dynamic-info-information.vue'),
  loading: ModalLoaderOverlay,
});
export const DynamicInfoTemplatesList = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/dynamic-info-templates-list.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateDynamicInfoTemplate = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/create-dynamic-info-template.vue'),
  loading: ModalLoaderOverlay,
});
export const CreatePlaylist = () => ({
  component: import(/* webpackChunkName: "Playlist" */ './admin/create-playlist.vue'),
  loading: ModalLoaderOverlay,
});
export const ManagePlaylistTabs = () => ({
  component: import(/* webpackChunkName: "Playlist" */ './admin/manage-playlist-tabs.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateRemediationInstruction = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-instruction.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateRemediationConfiguration = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-configuration.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateRemediationJob = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-job.vue'),
  loading: ModalLoaderOverlay,
});
export const ExecuteRemediationInstruction = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/execute-remediation-instruction.vue'),
  loading: ModalLoaderOverlay,
});
export const RemediationPatterns = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/remediation-patterns.vue'),
  loading: ModalLoaderOverlay,
});
export const RemediationInstructionApproval = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/remediation-instruction-approval.vue'),
  loading: ModalLoaderOverlay,
});
export const Rate = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './common/rate.vue'),
  loading: ModalLoaderOverlay,
});
export const Patterns = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/patterns.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateRemediationInstructionsFilter = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-instructions-filter.vue'),
  loading: ModalLoaderOverlay,
});
export const TestSuite = () => ({
  component: import(/* webpackChunkName: "Junit" */ './test-suite/test-suite.vue'),
  loading: ModalLoaderOverlay,
});
export const StateSetting = () => ({
  component: import(/* webpackChunkName: "Parameters" */ './state-setting/state-setting.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateIdleRule = () => ({
  component: import(/* webpackChunkName: "IdleRule" */ './idle-rule/create-idle-rule.vue'),
  loading: ModalLoaderOverlay,
});
export const HealthcheckEngine = () => ({
  component: import(/* webpackChunkName: "Healthcheck" */ './healthcheck/healthcheck-engine.vue'),
  loading: ModalLoaderOverlay,
});
export const HealthcheckEnginesChainReference = () => ({
  component: import(/* webpackChunkName: "Healthcheck" */ './healthcheck/healthcheck-engines-chain-reference.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateScenario = () => ({
  component: import(/* webpackChunkName: "Scenario" */ './scenario/create-scenario.vue'),
  loading: ModalLoaderOverlay,
});
export const CreateAlarmStatusRule = () => ({
  component: import(/* webpackChunkName: "AlarmStatusRule" */ './alarm-status-rule/create-alarm-status-rule.vue'),
  loading: ModalLoaderOverlay,
});
