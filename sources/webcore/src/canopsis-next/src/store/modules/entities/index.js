import alarmModule from './alarm';
import entityModule from './entity';
import watcherModule from './watcher';
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
import rightModule from './right';
import eventFilterRuleModule from './event-filter-rule';
import webhookModule from './webhook';
import infoModule from './info';
import filterHintModule from './filter-hint';
import snmpRuleModule from './snmp/rule';
import snmpMibModule from './snmp/mib';
import actionModule from './action';
import heartbeatModule from './heartbeat';
import keepaliveModule from './keepalive';
import dynamicInfoModule from './dynamic-info';
import dynamicInfoTemplateModule from './dynamic-info-template';
import alarmColumnFiltersModule from './alarm-column-filters';
import sessionModule from './session';
import broadcastMessageModule from './broadcast-message';
import counterModule from './counter';
import playlistModule from './playlist';
import metaAlarmRuleModule from './meta-alarm-rule';
import engineRunInfoModule from './engine-run-info';

export default {
  alarm: alarmModule,
  entity: entityModule,
  watcher: watcherModule,
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
  right: rightModule,
  eventFilterRule: eventFilterRuleModule,
  webhook: webhookModule,
  info: infoModule,
  snmpRule: snmpRuleModule,
  snmpMib: snmpMibModule,
  action: actionModule,
  heartbeat: heartbeatModule,
  keepalive: keepaliveModule,
  dynamicInfo: dynamicInfoModule,
  dynamicInfoTemplate: dynamicInfoTemplateModule,
  filterHint: filterHintModule,
  alarmColumnFilters: alarmColumnFiltersModule,
  session: sessionModule,
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
};
