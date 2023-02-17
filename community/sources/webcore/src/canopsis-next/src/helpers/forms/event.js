import {
  DECLARE_TICKET_OUTPUT,
  ENTITIES_STATES,
  EVENT_DEFAULT_ORIGIN,
  EVENT_ENTITY_TYPES,
  EVENT_INITIATORS,
  WEATHER_ACK_EVENT_OUTPUT,
  WEATHER_EVENT_DEFAULT_ENTITY,
} from '@/constants';

import { getNowTimestamp } from '@/helpers/date/date';

/**
 * @typedef {
 *   'ack' |
 *   'fastAck' |
 *   'ackremove' |
 *   'pbehaviorAdd' |
 *   'pbehaviorList' |
 *   'assocticket' |
 *   'cancel' |
 *   'delete' |
 *   'changestate' |
 *   'declareticket' |
 *   'snooze' |
 *   'validate' |
 *   'invalidate' |
 *   'pause' |
 *   'play' |
 *   'groupRequest' |
 *   'group' |
 *   'pbhenter' |
 *   'pbhleave' |
 *   'comment' |
 *   'manual_metaalarm_group' |
 *   'manual_metaalarm_ungroup' |
 *   'manual_metaalarm_update' |
 *   'stateinc' |
 *   'statedec' |
 *   'statusinc' |
 *   'statusdec' |
 *   'unsooze' |
 *   'metaalarmattach' |
 *   'executeInstruction' |
 *   'instructionstart' |
 *   'instructionpause' |
 *   'instructionresume' |
 *   'instructioncomplete' |
 *   'instructionabort' |
 *   'instructionfail' |
 *   'instructionjobstart' |
 *   'instructionjobcomplete' |
 *   'instructionjobabort' |
 *   'instructionjobfail' |
 *   'junittestsuiteupdate' |
 *   'junittestcaseupdate'
 * } EventType
 */

/**
 * @typedef {Object} Event
 *
 * @property {string} ref_rk
 * @property {string} resource
 * @property {string} origin
 * @property {string} source_type
 * @property {string} pbh_origin_icon
 * @property {string} connector_name
 * @property {string} component
 * @property {EventType} event_type
 * @property {EventType} crecord_type
 * @property {string} connector
 * @property {string} id
 * @property {number|string} state
 * @property {number} timestamp
 * @property {number} [state_type]
 */

/**
 * Prepare event by: type, alarm and already prepared data
 *
 * @param {EventType} type
 * @param {Alarm} alarm
 * @param {Object|Event} [data = {}]
 * @return {Event}
 */
export const prepareEventByAlarm = (type, alarm, data = {}) => {
  const event = {
    id: alarm._id,
    connector: alarm.v.connector,
    connector_name: alarm.v.connector_name,
    source_type: alarm.entity ? alarm.entity.type : null,
    component: alarm.v.component,
    state: alarm.v.state.val,
    event_type: type,
    crecord_type: type,
    timestamp: getNowTimestamp(),
    resource: alarm.v.resource,
    ref_rk: `${alarm.v.resource}/${alarm.v.component}`,
    origin: EVENT_DEFAULT_ORIGIN,
    initiator: EVENT_INITIATORS.user,
  };

  if (type !== EVENT_ENTITY_TYPES.snooze) {
    event.state_type = alarm.v.status.val;
  }

  return { ...event, ...data };
};

/**
 * Prepare event by: type, alarms and already prepared data
 *
 * @param {EventType} type
 * @param {Alarm[]} alarms
 * @param {Object|Event} [data]
 * @return {Event[]}
 */
export const prepareEventsByAlarms = (type, alarms, data) => (
  alarms.map(alarm => prepareEventByAlarm(type, alarm, data))
);

/**
 * Prepare event by: type, entity and already prepared data
 *
 * @param {Entity} entity
 * @param {EventType} type
 * @param {Event} [data = {}]
 * @return {Event}
 */
export const prepareEventByEntity = (entity, type, data = {}) => ({
  component: entity.component || WEATHER_EVENT_DEFAULT_ENTITY,
  connector: entity.connector || WEATHER_EVENT_DEFAULT_ENTITY,
  connector_name: entity.connector_name || WEATHER_EVENT_DEFAULT_ENTITY,
  crecord_type: type,
  event_type: type,
  ref_rk: `${entity.resource || WEATHER_EVENT_DEFAULT_ENTITY}/${entity.component || WEATHER_EVENT_DEFAULT_ENTITY}`,
  resource: entity.resource || WEATHER_EVENT_DEFAULT_ENTITY,
  source_type: entity.source_type,
  ...data,
});

/**
 * Create acknowledge event by entity data
 *
 * @param {Entity} entity
 * @param {string} output
 * @return {Event}
 */
export const createAckEventByEntity = ({ entity, output }) => prepareEventByEntity(
  entity,
  EVENT_ENTITY_TYPES.ack,
  { output },
);

/**
 * Create remove acknowledge event by entity data
 *
 * @param {Entity} entity
 * @param {string} output
 * @return {Event}
 */
export const createRemoveAckEventByEntity = ({ entity, output }) => prepareEventByEntity(
  entity,
  EVENT_ENTITY_TYPES.ackRemove,
  { output },
);

/**
 * Create associate ticket event by entity data
 *
 * @param {Entity} entity
 * @param {AssociateTicketEvent} payload
 * @return {Event}
 */
export const createAssociateTicketEventByEntity = ({ entity, ...payload }) => prepareEventByEntity(
  entity,
  EVENT_ENTITY_TYPES.assocTicket,
  payload,
);

/**
 * Create declare ticket event by entity data
 *
 * @param {Entity} entity
 * @return {Event}
 */
export const createDeclareTicketEventByEntity = ({ entity }) => prepareEventByEntity(
  entity,
  EVENT_ENTITY_TYPES.declareTicket,
  { output: DECLARE_TICKET_OUTPUT },
);

/**
 * Create validate event by entity data
 *
 * @param {Entity} entity
 * @return {Event}
 */
export const createValidateEventByEntity = ({ entity }) => prepareEventByEntity(
  entity,
  EVENT_ENTITY_TYPES.validate,
  {
    state: ENTITIES_STATES.critical,
    output: WEATHER_ACK_EVENT_OUTPUT.validateOk,
    keep_state: true,
  },
);

/**
 * Create invalidate event by entity data
 *
 * @param {Entity} entity
 * @return {Event}
 */
export const createInvalidateEventByEntity = ({ entity }) => prepareEventByEntity(
  entity,
  EVENT_ENTITY_TYPES.invalidate,
  {
    state: ENTITIES_STATES.major,
    output: WEATHER_ACK_EVENT_OUTPUT.validateCancel,
    keep_state: true,
  },
);

/**
 * Create comment event by entity data
 *
 * @param {Entity} entity
 * @param {string} output
 * @return {Event}
 */
export const createCommentEventByEntity = ({ entity, output }) => prepareEventByEntity(
  entity,
  EVENT_ENTITY_TYPES.comment,
  { output },
);

/**
 * Create cancel event by entity data
 *
 * @param {Entity} entity
 * @param {string} output
 * @return {Event}
 */
export const createCancelEventByEntity = ({ entity, output }) => prepareEventByEntity(
  entity,
  EVENT_ENTITY_TYPES.cancel,
  { output },
);
