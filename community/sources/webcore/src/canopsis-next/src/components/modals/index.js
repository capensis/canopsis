import LoaderOverlay from '@/components/common/overlay/loader-overlay.vue';

export const CreateEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-event.vue'),
  loading: LoaderOverlay,
});
export const CreateAckEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-ack-event.vue'),
  loading: LoaderOverlay,
});
export const ConfirmAckWithTicket = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/confirm-ack-with-ticket.vue'),
  loading: LoaderOverlay,
});
export const CreateAssociateTicketEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-associate-ticket-event.vue'),
  loading: LoaderOverlay,
});
export const CreateChangeStateEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-change-state-event.vue'),
  loading: LoaderOverlay,
});
export const CreateDeclareTicketEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-declare-ticket-event.vue'),
  loading: LoaderOverlay,
});
export const CreateSnoozeEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/create-snooze-event.vue'),
  loading: LoaderOverlay,
});
export const VariablesHelp = () => ({
  component: import(/* webpackChunkName: "Events" */ './common/variables-help.vue'),
  loading: LoaderOverlay,
});
export const EditLiveReporting = () => ({
  component: import(/* webpackChunkName: "Events" */ './alarm/edit-live-reporting.vue'),
  loading: LoaderOverlay,
});
export const CreateCommentEvent = () => ({
  component: import(/* webpackChunkName: "Events" */ './common/create-comment-event.vue'),
  loading: LoaderOverlay,
});
export const InfoPopupSetting = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/info-popup-setting/info-popup-setting.vue'),
  loading: LoaderOverlay,
});
export const AddInfoPopup = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/info-popup-setting/add-info-popup.vue'),
  loading: LoaderOverlay,
});
export const CreateManualMetaAlarm = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/create-manual-meta-alarm.vue'),
  loading: LoaderOverlay,
});
export const AlarmsList = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './alarm/alarms-list.vue'),
  loading: LoaderOverlay,
});
export const CreateMetaAlarmRule = () => ({
  component: import(/* webpackChunkName: "Alarm" */ './meta-alarm-rule/create-meta-alarm-rule.vue'),
  loading: LoaderOverlay,
});

export const PbehaviorList = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-list.vue'),
  loading: LoaderOverlay,
});
export const PbehaviorPlanning = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-planning.vue'),
  loading: LoaderOverlay,
});
export const PbehaviorsCalendar = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehaviors-calendar.vue'),
  loading: LoaderOverlay,
});
export const PbehaviorRecurrenceRule = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-recurrence-rule.vue'),
  loading: LoaderOverlay,
});
export const SelectExceptionsLists = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/select-exceptions-lists.vue'),
  loading: LoaderOverlay,
});
export const CreateRecurrenceRule = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-recurrence-rule.vue'),
  loading: LoaderOverlay,
});
export const PbehaviorRecurrentChangesConfirmation = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/pbehavior-recurrent-changes-confirmation.vue'),
  loading: LoaderOverlay,
});
export const CreatePbehavior = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior.vue'),
  loading: LoaderOverlay,
});
export const CreatePbehaviorType = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior-type.vue'),
  loading: LoaderOverlay,
});
export const CreatePbehaviorReason = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior-reason.vue'),
  loading: LoaderOverlay,
});
export const CreatePbehaviorException = () => ({
  component: import(/* webpackChunkName: "Pbehavior" */ './pbehavior/create-pbehavior-exception.vue'),
  loading: LoaderOverlay,
});
export const Confirmation = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/confirmation.vue'),
  loading: LoaderOverlay,
});
export const ClickOutsideConfirmation = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/click-outside-confirmation.vue'),
  loading: LoaderOverlay,
});
export const ColorPicker = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/color-picker.vue'),
  loading: LoaderOverlay,
});
export const TextEditor = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/text-editor.vue'),
  loading: LoaderOverlay,
});
export const TextFieldEditor = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/text-field-editor.vue'),
  loading: LoaderOverlay,
});
export const ConfirmationPhrase = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/confirmation-phrase.vue'),
  loading: LoaderOverlay,
});
export const ImageViewer = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/image-viewer.vue'),
  loading: LoaderOverlay,
});
export const ImagesViewer = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/images-viewer.vue'),
  loading: LoaderOverlay,
});
export const Info = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/info.vue'),
  loading: LoaderOverlay,
});
export const CreateWidget = () => ({
  component: import(/* webpackChunkName: "Widget" */ './view/create-widget.vue'),
  loading: LoaderOverlay,
});
export const CreateService = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/create-service.vue'),
  loading: LoaderOverlay,
});
export const CreateEntity = () => ({
  component: import(/* webpackChunkName: "Context" */ './entity/create-entity.vue'),
  loading: LoaderOverlay,
});
export const CreateEntityInfo = () => ({
  component: import(/* webpackChunkName: "Context" */ './entity/create-entity-info.vue'),
  loading: LoaderOverlay,
});
export const ServiceEntities = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/service-entities.vue'),
  loading: LoaderOverlay,
});
export const ServiceDependencies = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/service-dependencies.vue'),
  loading: LoaderOverlay,
});
export const CreateServicePauseEvent = () => ({
  component: import(/* webpackChunkName: "Context" */ './service/create-service-pause-event.vue'),
  loading: LoaderOverlay,
});
export const SelectView = () => ({
  component: import(/* webpackChunkName: "View" */ './view/select-view.vue'),
  loading: LoaderOverlay,
});
export const SelectViewTab = () => ({
  component: import(/* webpackChunkName: "View" */ './view/select-view-tab.vue'),
  loading: LoaderOverlay,
});
export const CreateView = () => ({
  component: import(/* webpackChunkName: "View" */ './view/create-view.vue'),
  loading: LoaderOverlay,
});
export const CreateGroup = () => ({
  component: import(/* webpackChunkName: "View" */ './view/create-group.vue'),
  loading: LoaderOverlay,
});
export const ImportGroupsAndViews = () => ({
  component: import(/* webpackChunkName: "View" */ './view/import-groups-and-views.vue'),
  loading: LoaderOverlay,
});
export const CreateFilter = () => ({
  component: import(/* webpackChunkName: "Filters" */ './common/create-filter.vue'),
  loading: LoaderOverlay,
});
export const FiltersList = () => ({
  component: import(/* webpackChunkName: "Filters" */ './common/filters-list.vue'),
  loading: LoaderOverlay,
});
export const CreateUser = () => ({
  component: import(/* webpackChunkName: "User" */ './admin/create-user.vue'),
  loading: LoaderOverlay,
});
export const CreateRole = () => ({
  component: import(/* webpackChunkName: "Role" */ './admin/create-role.vue'),
  loading: LoaderOverlay,
});
export const CreateBroadcastMessage = () => ({
  component: import(/* webpackChunkName: "BroadcastMessage" */ './admin/create-broadcast-message.vue'),
  loading: LoaderOverlay,
});
export const CreateEventFilterRule = () => ({
  component: import(/* webpackChunkName: "EventFilter" */ './event-filter/create-event-filter-rule.vue'),
  loading: LoaderOverlay,
});
export const EventFilterRuleActions = () => ({
  component: import(/* webpackChunkName: "EventFilter" */ './event-filter/event-filter-rule-actions.vue'),
  loading: LoaderOverlay,
});
export const EventFilterRuleExternalData = () => ({
  component: import(/* webpackChunkName: "EventFilter" */ './event-filter/event-filter-rule-external-data.vue'),
  loading: LoaderOverlay,
});
export const CreateEventFilterRuleAction = () => ({
  component: import(/* webpackChunkName: "EventFilter" */ './event-filter/create-event-filter-rule-action.vue'),
  loading: LoaderOverlay,
});
export const CreatePattern = () => ({
  component: import(/* webpackChunkName: "Pattern" */ './pattern/create-pattern.vue'),
  loading: LoaderOverlay,
});
export const CreatePatternRule = () => ({
  component: import(/* webpackChunkName: "Pattern" */ './pattern/create-pattern-rule.vue'),
  loading: LoaderOverlay,
});
export const CreateSnmpRule = () => ({
  component: import(/* webpackChunkName: "SnmpRule" */ './snmp-rule/create-snmp-rule.vue'),
  loading: LoaderOverlay,
});
export const CreateDynamicInfo = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/create-dynamic-info.vue'),
  loading: LoaderOverlay,
});
export const CreateDynamicInfoInformation = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/create-dynamic-info-information.vue'),
  loading: LoaderOverlay,
});
export const DynamicInfoTemplatesList = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/dynamic-info-templates-list.vue'),
  loading: LoaderOverlay,
});
export const CreateDynamicInfoTemplate = () => ({
  component: import(/* webpackChunkName: "DynamicInfo" */ './dynamic-info/create-dynamic-info-template.vue'),
  loading: LoaderOverlay,
});
export const CreatePlaylist = () => ({
  component: import(/* webpackChunkName: "Playlist" */ './admin/create-playlist.vue'),
  loading: LoaderOverlay,
});
export const ManagePlaylistTabs = () => ({
  component: import(/* webpackChunkName: "Playlist" */ './admin/manage-playlist-tabs.vue'),
  loading: LoaderOverlay,
});
export const CreateRemediationInstruction = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-instruction.vue'),
  loading: LoaderOverlay,
});
export const CreateRemediationConfiguration = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-configuration.vue'),
  loading: LoaderOverlay,
});
export const CreateRemediationJob = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-job.vue'),
  loading: LoaderOverlay,
});
export const ExecuteRemediationInstruction = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/execute-remediation-instruction.vue'),
  loading: LoaderOverlay,
});
export const RemediationPatterns = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/remediation-patterns.vue'),
  loading: LoaderOverlay,
});
export const RemediationInstructionApproval = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/remediation-instruction-approval.vue'),
  loading: LoaderOverlay,
});
export const Rate = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './common/rate.vue'),
  loading: LoaderOverlay,
});
export const Patterns = () => ({
  component: import(/* webpackChunkName: "Common" */ './common/patterns.vue'),
  loading: LoaderOverlay,
});
export const CreateRemediationInstructionsFilter = () => ({
  component: import(/* webpackChunkName: "Remediation" */ './remediation/create-remediation-instructions-filter.vue'),
  loading: LoaderOverlay,
});
export const TestSuite = () => ({
  component: import(/* webpackChunkName: "Junit" */ './test-suite/test-suite.vue'),
  loading: LoaderOverlay,
});
export const StateSetting = () => ({
  component: import(/* webpackChunkName: "Parameters" */ './state-setting/state-setting.vue'),
  loading: LoaderOverlay,
});
export const CreateIdleRule = () => ({
  component: import(/* webpackChunkName: "IdleRule" */ './idle-rule/create-idle-rule.vue'),
  loading: LoaderOverlay,
});
export const HealthcheckEngine = () => ({
  component: import(/* webpackChunkName: "Healthcheck" */ './healthcheck/healthcheck-engine.vue'),
  loading: LoaderOverlay,
});
export const HealthcheckEnginesChainReference = () => ({
  component: import(/* webpackChunkName: "Healthcheck" */ './healthcheck/healthcheck-engines-chain-reference.vue'),
  loading: LoaderOverlay,
});
export const CreateScenario = () => ({
  component: import(/* webpackChunkName: "Scenario" */ './scenario/create-scenario.vue'),
  loading: LoaderOverlay,
});
export const CreateAlarmStatusRule = () => ({
  component: import(/* webpackChunkName: "AlarmStatusRule" */ './alarm-status-rule/create-alarm-status-rule.vue'),
  loading: LoaderOverlay,
});
