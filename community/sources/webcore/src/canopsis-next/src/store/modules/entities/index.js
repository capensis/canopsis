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
import userPreferenceModule from './user-preference';
import viewModule from './view';
import statsModule from './stats';
import roleModule from './role';
import userModule from './user';
import permissionModule from './permission';
import eventFilterRuleModule from './event-filter-rule';
import infoModule from './info';
import snmpRuleModule from './snmp/rule';
import snmpMibModule from './snmp/mib';
import heartbeatModule from './heartbeat';
import keepaliveModule from './keepalive';
import dynamicInfoModule from './dynamic-info';
import sessionsCountModule from './sessions-count';
import broadcastMessageModule from './broadcast-message';
import counterModule from './counter';
import playlistModule from './playlist';
import metaAlarmRuleModule from './meta-alarm-rule';
import engineRunInfoModule from './engine-run-info';
import scenarioModule from './scenario';
import entityCategoryModule from './entity-category';
import testSuiteModule from './test-suite';
import stateSettingModule from './state-setting';
import associativeTableModule from './associative-table';
import dataStorageModule from './data-storage';
import notificationSettingsModule from './notification-settings';

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
  stats: statsModule,
  role: roleModule,
  user: userModule,
  permission: permissionModule,
  eventFilterRule: eventFilterRuleModule,
  info: infoModule,
  snmpRule: snmpRuleModule,
  snmpMib: snmpMibModule,
  heartbeat: heartbeatModule,
  keepalive: keepaliveModule,
  dynamicInfo: dynamicInfoModule,
  sessionsCount: sessionsCountModule,
  broadcastMessage: broadcastMessageModule,
  counter: counterModule,
  playlist: playlistModule,
  metaAlarmRule: metaAlarmRuleModule,
  engineRunInfo: engineRunInfoModule,
  remediationInstruction: remediationInstructionModule,
  remediationJob: remediationJobModule,
  remediationConfiguration: remediationConfigurationModule,
  remediationInstructionExecution: remediationInstructionExecutionModule,
  remediationJobExecution: remediationJobExecutionModule,
  scenario: scenarioModule,
  entityCategory: entityCategoryModule,
  testSuite: testSuiteModule,
  stateSetting: stateSettingModule,
  associativeTable: associativeTableModule,
  dataStorage: dataStorageModule,
  notificationSettings: notificationSettingsModule,
};