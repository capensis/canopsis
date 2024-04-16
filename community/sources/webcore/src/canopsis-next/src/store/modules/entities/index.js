import alarmModule from './alarm';
import entityModule from './entity';
import serviceModule from './service';
import pbehaviorModule from './pbehavior';
import pbehaviorReasonsModule from './pbehavior-reasons';
import pbehaviorTimespanModule from './pbehavior-timespan';
import pbehaviorExceptionsModule from './pbehavior-exceptions';
import pbehaviorTypesModule from './pbehavior-types';
import remediationInstructionModule from './remediation-instruction';
import remediationInstructionExecutionModule from './remediation-instruction-execution';
import remediationConfigurationModule from './remediation-configuration';
import remediationJobModule from './remediation-job';
import remediationJobExecutionModule from './remediation-job-execution';
import remediationInstructionStatsModule from './remediation-instruction-stats';
import remediationStatisticModule from './remediation-statistic';
import userPreferenceModule from './user-preference';
import viewModule from './view';
import roleModule from './role';
import userModule from './user';
import permissionModule from './permission';
import eventFilterModule from './event-filter';
import infoModule from './info';
import infosModule from './infos';
import snmpRuleModule from './snmp/rule';
import snmpMibModule from './snmp/mib';
import dynamicInfoModule from './dynamic-info';
import broadcastMessageModule from './broadcast-message';
import counterModule from './counter';
import playlistModule from './playlist';
import metaAlarmRuleModule from './meta-alarm-rule';
import scenarioModule from './scenario';
import entityCategoryModule from './entity-category';
import testSuiteModule from './test-suite';
import stateSettingModule from './state-setting';
import associativeTableModule from './associative-table';
import dataStorageModule from './data-storage';
import notificationSettingsModule from './notification-settings';
import idleRulesModule from './idle-rules';
import flappingRulesModule from './flapping-rules';
import resolveRulesModule from './resolve-rules';
import healthcheckModule from './healthcheck';
import healthcheckParametersModule from './healthcheck-parameters';
import messageRateStatsModule from './message-rate-stats';
import metricsModule from './metrics';
import filterModule from './filter';
import ratingSettingsModule from './rating-settings';
import patternModule from './pattern';
import mapModule from './map';
import alarmTagModule from './alarm-tag';
import shareTokenModule from './share-token';
import techMetricsModule from './tech-metrics';
import widgetTemplateModule from './widget-template';
import manualMetaAlarmModule from './manual-meta-alarm';
import metaAlarmModule from './meta-alarm';
import templateVarsModule from './template-vars';
import declareTicketRuleModule from './declare-ticket-rule';
import templateValidatorModule from './template-validator';
import LinkRuleRuleModule from './links-rule';
import metricsSettingsModule from './metrics-settings';
import aggregatedMetricsModule from './aggregated-metrics';
import vectorMetricsModule from './vector-metrics';
import groupMetricsModule from './group-metrics';
import tagModule from './tag';
import themeModule from './theme';
import iconModule from './icon';
import availabilityModule from './availability';

export default {
  alarm: alarmModule,
  entity: entityModule,
  service: serviceModule,
  pbehavior: pbehaviorModule,
  pbehaviorReasons: pbehaviorReasonsModule,
  pbehaviorTimespan: pbehaviorTimespanModule,
  pbehaviorExceptions: pbehaviorExceptionsModule,
  pbehaviorTypes: pbehaviorTypesModule,
  userPreference: userPreferenceModule,
  view: viewModule,
  role: roleModule,
  user: userModule,
  permission: permissionModule,
  eventFilter: eventFilterModule,
  info: infoModule,
  infos: infosModule,
  snmpRule: snmpRuleModule,
  snmpMib: snmpMibModule,
  dynamicInfo: dynamicInfoModule,
  broadcastMessage: broadcastMessageModule,
  counter: counterModule,
  playlist: playlistModule,
  metaAlarmRule: metaAlarmRuleModule,
  remediationInstruction: remediationInstructionModule,
  remediationJob: remediationJobModule,
  remediationConfiguration: remediationConfigurationModule,
  remediationInstructionExecution: remediationInstructionExecutionModule,
  remediationJobExecution: remediationJobExecutionModule,
  remediationInstructionStats: remediationInstructionStatsModule,
  remediationStatistic: remediationStatisticModule,
  scenario: scenarioModule,
  entityCategory: entityCategoryModule,
  testSuite: testSuiteModule,
  stateSetting: stateSettingModule,
  associativeTable: associativeTableModule,
  dataStorage: dataStorageModule,
  notificationSettings: notificationSettingsModule,
  idleRules: idleRulesModule,
  flappingRules: flappingRulesModule,
  resolveRules: resolveRulesModule,
  healthcheck: healthcheckModule,
  healthcheckParameters: healthcheckParametersModule,
  messageRateStats: messageRateStatsModule,
  metrics: metricsModule,
  filter: filterModule,
  ratingSettings: ratingSettingsModule,
  pattern: patternModule,
  map: mapModule,
  alarmTag: alarmTagModule,
  shareToken: shareTokenModule,
  techMetrics: techMetricsModule,
  widgetTemplate: widgetTemplateModule,
  manualMetaAlarm: manualMetaAlarmModule,
  metaAlarm: metaAlarmModule,
  templateVars: templateVarsModule,
  declareTicketRule: declareTicketRuleModule,
  templateValidator: templateValidatorModule,
  linkRule: LinkRuleRuleModule,
  metricsSettings: metricsSettingsModule,
  aggregatedMetrics: aggregatedMetricsModule,
  vectorMetrics: vectorMetricsModule,
  groupMetrics: groupMetricsModule,
  tag: tagModule,
  theme: themeModule,
  icon: iconModule,
  availability: availabilityModule,
};
