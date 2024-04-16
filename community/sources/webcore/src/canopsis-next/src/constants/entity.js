import { COLORS } from '@/config';

// eslint-disable-next-line import/no-webpack-loader-syntax
import engineeringIcon from '!!svg-inline-loader?modules!@/assets/images/engineering.svg';

import { PBEHAVIOR_TYPE_TYPES } from './pbehavior';

export const ENTITY_FIELDS = {
  id: '_id',
  name: 'name',
  categoryName: 'category.name',
  type: 'type',
  component: 'component',
  connector: 'connector',
  connectorName: 'connector_name',
  resource: 'resource',
  impactLevel: 'impact_level',
  lastEventDate: 'last_event_date',
  lastPbehaviorDate: 'last_pbehavior_date',
  lastUpdateDate: 'last_update_date',
  koEvents: 'ko_events',
  okEvents: 'ok_events',
  statsOk: 'stats.ok',
  statsKo: 'stats.ko',
  pbehaviorInfo: 'pbehavior_info',
  state: 'state',
  impactState: 'impact_state',
  status: 'status',
  idleSince: 'idle_since',
  enabled: 'enabled',
  infos: 'infos',
  componentInfos: 'component_infos',
  links: 'links',
  alarmDisplayName: 'alarm_display_name',
  alarmCreationDate: 'alarm_creation_date',
  importSource: 'import_source',
  imported: 'imported',
  alarmLastComment: 'alarm_last_comment',

  /**
   * OBJECTS
   */
  ack: 'ack',
  category: 'category',
  ticket: 'ticket',
  snooze: 'snooze',
};

export const EVENT_ENTITY_TYPES = { // TODO: remove it in the future
  ack: 'ack', // TODO: ServiceWeather
  fastAck: 'fastAck', // TODO: ServiceWeather
  ackRemove: 'ackremove', // TODO: ServiceWeather, EventPattern
  pbehaviorList: 'pbehaviorList', // TODO: ServiceWeather
  cancel: 'cancel', // TODO: ServiceWeather, EventPattern
  delete: 'delete', // TODO: ServiceWeather
  assocTicket: 'assocticket', // TODO: DeclareTicket, EventPattern
  declareTicket: 'declareticket', // TODO: DeclareTicket, EventPattern
  declareTicketFail: 'declareticketfail', // TODO: DeclareTicket, EventPattern
  webhookStart: 'webhookstart', // TODO: DeclareTicket, EventPattern
  webhookComplete: 'webhookcomplete', // TODO: DeclareTicket, EventPattern
  webhookFail: 'webhookfail', // TODO: DeclareTicket, EventPattern
  snooze: 'snooze', // TODO: EventPattern
  unsooze: 'unsooze', // TODO: remove
  validate: 'validate', // TODO: ServiceWeather
  invalidate: 'invalidate', // TODO: ServiceWeather
  pause: 'pause', // TODO: ServiceWeather
  play: 'play', // TODO: ServiceWeather
  pbhenter: 'pbhenter', // TODO: remove
  pbhleave: 'pbhleave', // TODO: remove
  comment: 'comment', // TODO: ServiceWeather, EventPattern
  metaalarmattach: 'metaalarmattach', // TODO: remove
  executeInstruction: 'executeInstruction', // TODO: ServiceWeather
};

export const ENTITY_INFOS_TYPE = { // TODO: move to alarm
  state: 'state',
  status: 'status',
  action: 'action',
};

export const SERVICE_STATES = {
  ok: 'ok',
  minor: 'minor',
  major: 'major',
  critical: 'critical',
  pause: 'pause',
};

export const WEATHER_ICONS = {
  [SERVICE_STATES.ok]: 'wb_sunny',
  [SERVICE_STATES.minor]: 'person',
  [SERVICE_STATES.major]: 'person',
  [SERVICE_STATES.critical]: 'wb_cloudy',

  [PBEHAVIOR_TYPE_TYPES.maintenance]: 'build',
  [PBEHAVIOR_TYPE_TYPES.inactive]: 'brightness_3',
  [PBEHAVIOR_TYPE_TYPES.pause]: 'pause',
};

export const WEATHER_ACTIONS_TYPES = {
  entityAck: 'entityAck',
  entityAckRemove: 'entityAckRemove',
  entityAssocTicket: 'entityAssocTicket',
  entityValidate: 'entityValidate',
  entityInvalidate: 'entityInvalidate',
  entityPause: 'entityPause',
  entityPlay: 'entityPlay',
  entityCancel: 'entityCancel',
  entityComment: 'entityComment',
  declareTicket: 'declareTicket',

  entityLinks: 'entityLinks',

  moreInfos: 'moreInfos',
  alarmsList: 'alarmsList',
  pbehaviorList: 'pbehaviorList',

  executeInstruction: 'executeInstruction',
};

export const EVENT_ENTITY_ICONS_BY_TYPE = { // TODO: remove it
  [EVENT_ENTITY_TYPES.ack]: 'check',
  [EVENT_ENTITY_TYPES.fastAck]: 'done_all',
  [EVENT_ENTITY_TYPES.pbehaviorList]: 'list',
  [EVENT_ENTITY_TYPES.ackRemove]: 'remove_done',
  [EVENT_ENTITY_TYPES.declareTicket]: 'note_add',
  [EVENT_ENTITY_TYPES.declareTicketFail]: 'report_problem',
  [EVENT_ENTITY_TYPES.webhookStart]: 'report_problem',
  [EVENT_ENTITY_TYPES.webhookComplete]: 'report_problem',
  [EVENT_ENTITY_TYPES.webhookFail]: 'report_problem',
  [EVENT_ENTITY_TYPES.assocTicket]: '$vuetify.icons.sticky_note_2',
  [EVENT_ENTITY_TYPES.delete]: 'delete',
  [EVENT_ENTITY_TYPES.snooze]: 'alarm',
  [EVENT_ENTITY_TYPES.validate]: 'thumb_up',
  [EVENT_ENTITY_TYPES.invalidate]: 'thumb_down',
  [EVENT_ENTITY_TYPES.pause]: 'pause',
  [EVENT_ENTITY_TYPES.play]: 'play_arrow',
  [EVENT_ENTITY_TYPES.pbhenter]: 'pause',
  [EVENT_ENTITY_TYPES.pbhleave]: 'play_arrow',
  [EVENT_ENTITY_TYPES.comment]: 'comment',
  [EVENT_ENTITY_TYPES.metaalarmattach]: 'center_focus_weak',
  [EVENT_ENTITY_TYPES.executeInstruction]: 'assignment',
  [EVENT_ENTITY_TYPES.cancel]: 'block',
  groupChildren: 'center_focus_strong',
  groupParents: 'center_focus_weak',
};

export const EVENT_ENTITY_COLORS_BY_TYPE = {
  [EVENT_ENTITY_TYPES.ack]: COLORS.entitiesEvents.ack,
  [EVENT_ENTITY_TYPES.ackRemove]: COLORS.entitiesEvents.ackRemove,
  [EVENT_ENTITY_TYPES.snooze]: COLORS.entitiesEvents.snooze,
  [EVENT_ENTITY_TYPES.pbhenter]: COLORS.entitiesEvents.pbhenter,
  [EVENT_ENTITY_TYPES.pbhleave]: COLORS.entitiesEvents.pbhleave,
  [EVENT_ENTITY_TYPES.comment]: COLORS.entitiesEvents.comment,
  [EVENT_ENTITY_TYPES.metaalarmattach]: COLORS.entitiesEvents.metaalarmattach,
};

export const ENTITY_EVENT_BY_ACTION_TYPE = {
  [WEATHER_ACTIONS_TYPES.entityAck]: EVENT_ENTITY_TYPES.ack,
  [WEATHER_ACTIONS_TYPES.entityAssocTicket]: EVENT_ENTITY_TYPES.assocTicket,
  [WEATHER_ACTIONS_TYPES.entityValidate]: EVENT_ENTITY_TYPES.validate,
  [WEATHER_ACTIONS_TYPES.entityInvalidate]: EVENT_ENTITY_TYPES.invalidate,
  [WEATHER_ACTIONS_TYPES.entityPause]: EVENT_ENTITY_TYPES.pause,
  [WEATHER_ACTIONS_TYPES.entityPlay]: EVENT_ENTITY_TYPES.play,
  [WEATHER_ACTIONS_TYPES.entityCancel]: EVENT_ENTITY_TYPES.cancel,
  [WEATHER_ACTIONS_TYPES.entityComment]: EVENT_ENTITY_TYPES.comment,
  [WEATHER_ACTIONS_TYPES.pbehaviorList]: EVENT_ENTITY_TYPES.pbehaviorList,
  [WEATHER_ACTIONS_TYPES.executeInstruction]: EVENT_ENTITY_TYPES.executeInstruction,
  [WEATHER_ACTIONS_TYPES.declareTicket]: EVENT_ENTITY_TYPES.declareTicket,
  [WEATHER_ACTIONS_TYPES.entityAckRemove]: EVENT_ENTITY_TYPES.ackRemove,
};

export const SERVICE_WEATHER_WIDGET_MODAL_TYPES = {
  moreInfo: 'more-info',
  alarmList: 'alarm-list',
  both: 'both',
};

export const WEATHER_ACK_EVENT_OUTPUT = {
  ack: 'MDS_ACKNOWLEDGE',
  validateOk: 'MDS_VALIDATEOK',
  validateCancel: 'MDS_VALIDATECANCEL',
};

export const DEFAULT_CONTEXT_WIDGET_COLUMNS = [
  { value: ENTITY_FIELDS.name },
  { value: ENTITY_FIELDS.type },
];

export const DEFAULT_SERVICE_DEPENDENCIES_COLUMNS = [...DEFAULT_CONTEXT_WIDGET_COLUMNS];

export const AVAILABLE_COUNTERS = {
  total: 'total',
  total_active: 'total_active',
  snooze: 'snooze',
  ack: 'ack',
  ticket: 'ticket',
  pbehavior_active: 'pbehavior_active',
};

export const DEFAULT_COUNTER_BLOCK_TEMPLATE = `<h2 style="text-align: justify;">
  <strong>{{ counter.filter.title }}</strong></h2><br>
  <center><strong><span style="font-size: 18px;">{{ counter.total_active }} alarmes actives</span></strong></center>
  <br>Seuil mineur à {{ levels.values.minor }}, seuil critique à {{ levels.values.critical }}
  <p style="text-align: justify;">{{ counter.ack }} acquittées, {{ counter.ticket}} avec ticket</p>`;

export const PBEHAVIOR_COUNTERS_LIMIT = 3;

export const BASIC_ENTITY_TYPES = {
  connector: 'connector',
  component: 'component',
  resource: 'resource',
};

export const ENTITY_TYPES = {
  service: 'service',

  ...BASIC_ENTITY_TYPES,
};

export const COLOR_INDICATOR_TYPES = {
  state: 'state',
  impactState: 'impact_state',
};

export const STATE_SETTING_ENTITY_TYPES = [
  ENTITY_TYPES.component,
  ENTITY_TYPES.service,
];

export const STATE_SETTING_METHODS = {
  inherited: 'inherited',
  dependencies: 'dependencies',
};

export const STATE_SETTING_THRESHOLDS_METHODS = {
  share: 'share',
  number: 'number',
};

export const JUNIT_STATE_SETTING_METHODS = {
  worst: 'worst',
  worstOfShare: 'worst_of_share',
};

export const JUNIT_STATE_SETTING_THRESHOLDS_TYPES = {
  number: 0,
  percent: 1,
};

export const STATE_SETTING_THRESHOLDS_CONDITIONS = {
  greater: 'gt',
  less: 'lt',
};

export const JUNIT_STATE_SETTING_ID = 'junit';

export const SERVICE_STATE_SETTING_ID = 'service';

export const CONTEXT_ACTIONS_TYPES = {
  createEntity: 'createEntity',
  editEntity: 'editEntity',
  duplicateEntity: 'duplicateEntity',
  deleteEntity: 'deleteEntity',
  pbehaviorAdd: 'pbehaviorAdd',
  pbehaviorList: 'pbehaviorList',
  pbehaviorDelete: 'pbehaviorDelete',
  variablesHelp: 'variablesHelp',
  massEnable: 'massEnable',
  massDisable: 'massDisable',

  listFilters: 'listFilters',
  editFilter: 'editFilter',
  addFilter: 'addFilter',
};

export const COUNTER_ACTIONS_TYPES = {
  alarmsList: 'alarmsList',
  variablesHelp: 'variablesHelp',
};

export const ENTITY_INFOS_FIELDS = [
  ENTITY_FIELDS.infos,
  ENTITY_FIELDS.componentInfos,
];

const {
  ack,
  category,
  ticket,
  snooze,
  connectorName,
  resource,
  statsOk,
  statsKo,
  alarmCreationDate,
  alarmDisplayName,
  ...contextWidgetColumns
} = ENTITY_FIELDS;

export const CONTEXT_WIDGET_COLUMNS = contextWidgetColumns;

export const ENTITY_PATTERN_FIELDS = {
  id: ENTITY_FIELDS.id,
  name: ENTITY_FIELDS.name,
  type: ENTITY_FIELDS.type,
  component: ENTITY_FIELDS.component,
  connector: ENTITY_FIELDS.connector,
  infos: ENTITY_FIELDS.infos,
  componentInfos: ENTITY_FIELDS.componentInfos,
  category: ENTITY_FIELDS.category,
  impactLevel: ENTITY_FIELDS.impactLevel,
  lastEventDate: ENTITY_FIELDS.lastEventDate,
};

export const ENTITY_TEMPLATE_FIELDS = {
  id: `entity.${ENTITY_FIELDS.id}`,
  name: `entity.${ENTITY_FIELDS.name}`,
  infos: `entity.${ENTITY_FIELDS.infos}`,
  connector: `entity.${ENTITY_FIELDS.connector}`,
  connectorName: `entity.${ENTITY_FIELDS.connectorName}`,
  component: `entity.${ENTITY_FIELDS.component}`,
  resource: `entity.${ENTITY_FIELDS.resource}`,
  state: `entity.${ENTITY_FIELDS.state}`,
  status: `entity.${ENTITY_FIELDS.status}`,
  snooze: `entity.${ENTITY_FIELDS.snooze}`,
  ack: `entity.${ENTITY_FIELDS.ack}`,
  lastUpdateDate: `entity.${ENTITY_FIELDS.lastUpdateDate}`,
  impactLevel: `entity.${ENTITY_FIELDS.impactLevel}`,
  impactState: `entity.${ENTITY_FIELDS.impactState}`,
  categoryName: `entity.${ENTITY_FIELDS.categoryName}`,
  alarmDisplayName: `entity.${ENTITY_FIELDS.alarmDisplayName}`,
  pbehaviorInfo: `entity.${ENTITY_FIELDS.pbehaviorInfo}`,
  alarmCreationDate: `entity.${ENTITY_FIELDS.alarmCreationDate}`,
  ticket: `entity.${ENTITY_FIELDS.ticket}`,
  statsOk: `entity.${ENTITY_FIELDS.statsOk}`,
  statsKo: `entity.${ENTITY_FIELDS.statsKo}`,
  links: `entity.${ENTITY_FIELDS.links}`,
  alarmLastComment: `entity.${ENTITY_FIELDS.alarmLastComment}`,
};

export const ENTITY_FIELDS_TO_LABELS_KEYS = {
  [ENTITY_FIELDS.id]: 'common.id',
  [ENTITY_FIELDS.name]: 'common.name',
  [ENTITY_FIELDS.categoryName]: 'entity.fields.categoryName',
  [ENTITY_FIELDS.type]: 'common.type',
  [ENTITY_FIELDS.component]: 'common.component',
  [ENTITY_FIELDS.connector]: 'common.connector',
  [ENTITY_FIELDS.connectorName]: 'common.connectorName',
  [ENTITY_FIELDS.resource]: 'common.resource',
  [ENTITY_FIELDS.impactLevel]: 'common.impactLevel',
  [ENTITY_FIELDS.lastEventDate]: 'common.lastEventDate',
  [ENTITY_FIELDS.lastPbehaviorDate]: 'common.lastPbehaviorDate',
  [ENTITY_FIELDS.lastUpdateDate]: 'common.updated',
  [ENTITY_FIELDS.koEvents]: 'entity.fields.koEvents',
  [ENTITY_FIELDS.okEvents]: 'entity.fields.okEvents',
  [ENTITY_FIELDS.statsOk]: 'entity.fields.statsOk',
  [ENTITY_FIELDS.statsKo]: 'entity.fields.statsKo',
  [ENTITY_FIELDS.pbehaviorInfo]: 'pbehavior.pbehaviorInfo',
  [ENTITY_FIELDS.state]: 'common.state',
  [ENTITY_FIELDS.impactState]: 'common.impactState',
  [ENTITY_FIELDS.status]: 'common.status',
  [ENTITY_FIELDS.idleSince]: 'entity.fields.idleSince',
  [ENTITY_FIELDS.enabled]: 'common.enabled',
  [ENTITY_FIELDS.infos]: 'common.infos',
  [ENTITY_FIELDS.componentInfos]: 'entity.fields.componentInfos',
  [ENTITY_FIELDS.links]: 'common.link',
  [ENTITY_FIELDS.alarmDisplayName]: 'entity.fields.alarmDisplayName',
  [ENTITY_FIELDS.alarmCreationDate]: 'entity.fields.alarmCreationDate',
  [ENTITY_FIELDS.importSource]: 'entity.fields.importSource',
  [ENTITY_FIELDS.imported]: 'entity.fields.imported',
  [ENTITY_FIELDS.alarmLastComment]: 'entity.fields.alarmLastComment',

  /**
   * OBJECTS
   */
  [ENTITY_FIELDS.ack]: 'common.ack',
  [ENTITY_FIELDS.category]: 'common.category',
  [ENTITY_FIELDS.ticket]: 'common.ticket',
  [ENTITY_FIELDS.snooze]: 'common.snooze',
  [ENTITY_FIELDS.pbehaviorInfo]: 'pbehavior.pbehaviorInfo',
};

export const ENTITY_UNSORTABLE_FIELDS = [
  ENTITY_FIELDS.links,
  ENTITY_FIELDS.pbehaviorInfo,
];

export const ENTITY_PAYLOADS_VARIABLES = {
  entity: '.Entity',
  entities: '.Entities',
  name: '.Name',
  infosValue: '(index .Infos "%infos_name%").Value',
};

export const SERVICE_WEATHER_DEFAULT_EM_HEIGHT = 4;

export const ENTITY_EXPORT_FILE_NAME_PREFIX = 'entity';

export const TREE_OF_DEPENDENCIES_SHOW_TYPES = {
  allDependencies: 0,
  dependenciesDefiningTheState: 1,
  custom: 2,
};

export const ROOT_CAUSE_DIAGRAM_NODE_SIZE = 50;

export const ROOT_CAUSE_DIAGRAM_EVENTS_NODE_SIZE = 30;

export const ROOT_CAUSE_DIAGRAM_OPTIONS = {
  fitPadding: 40,
  wheelSensitivity: 0.3,
  minZoom: 0.05,
  maxZoom: 1.5,
  nodeSize: ROOT_CAUSE_DIAGRAM_NODE_SIZE,
};

export const ROOT_CAUSE_DIAGRAM_TOOLTIP_OFFSET = 5;

export const ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS = {
  name: 'dagre',
  animate: 'end',
  fit: true,
  rankDir: 'LR',
  padding: 40,
  minLen: 5,
};

export const ENTITY_TYPES_ICONS = {
  [ENTITY_TYPES.component]: 'developer_board',
  [ENTITY_TYPES.connector]: 'perm_identity',
  [ENTITY_TYPES.resource]: 'perm_identity',
  [ENTITY_TYPES.service]: '$vuetify.icons.engineering',
};

export const ENTITY_TYPES_ICONS_FOR_CYTOSCAPE = {
  ...ENTITY_TYPES_ICONS,

  [ENTITY_TYPES.service]: engineeringIcon,
};
