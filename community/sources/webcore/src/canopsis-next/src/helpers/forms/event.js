import {
  ENTITIES_STATES,
  EVENT_DEFAULT_ORIGIN,
  EVENT_ENTITY_TYPES,
  EVENT_INITIATORS,
  MANUAL_META_ALARM_EVENT_DEFAULT_FIELDS,
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
 *   'done' |
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
 * @typedef {Object} ManualMetaAlarmEvent
 *
 * @property {EventType} event_type
 * @property {string} connector
 * @property {string} connector_name
 * @property {string} component
 * @property {string} source_type
 * @property {number|string} state
 * @property {string} [display_name]
 * @property {string[]} ma_children
 * @property {string[]} [ma_parent]
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
 * Prepare manual meta alarm event by: type, alarms and already prepared data
 *
 * @param {EventType} type
 * @param {Alarm[]} alarms
 * @param {Object|ManualMetaAlarmEvent} [data = {}]
 * @return {ManualMetaAlarmEvent[]}
 */
export const prepareManualMetaAlarmEventByAlarms = (type, alarms, data = {}) => [{
  ...MANUAL_META_ALARM_EVENT_DEFAULT_FIELDS,

  event_type: type,
  ma_children: alarms.map(({ entity }) => entity._id),
  state: ENTITIES_STATES.minor,

  ...data,
}];

/**
 * Prepare event by: type, alarms and already prepared data
 *
 * @param {EventType} type
 * @param {Alarm[]} alarms
 * @param {Object|Event} [data]
 * @return {Event[]|ManualMetaAlarmEvent[]}
 */
export const prepareEventsByAlarms = (type, alarms, data) => (
  [EVENT_ENTITY_TYPES.manualMetaAlarmGroup, EVENT_ENTITY_TYPES.manualMetaAlarmUpdate].includes(type)
    ? prepareManualMetaAlarmEventByAlarms(type, alarms, data)
    : alarms.map(alarm => prepareEventByAlarm(type, alarm, data))
);
