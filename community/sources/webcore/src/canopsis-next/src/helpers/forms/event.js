import moment from 'moment';

import { EVENT_DEFAULT_ORIGIN, EVENT_ENTITY_TYPES, EVENT_INITIATORS } from '@/constants';

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
 * @property {string} author
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
 * @property {number} state_type
 */

/**
 *
 * @param {EventType} type
 * @param {Alarm} alarm
 * @param {Object|Event} [data]
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
    timestamp: moment().unix(),
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
 *
 * @param {EventType} type
 * @param {Alarm[]} alarms
 * @param {Object|Event} [data]
 * @return {Event[]}
 */
export const prepareEventsByAlarms = (type, alarms, data) =>
  alarms.map(alarm => prepareEventByAlarm(type, alarm, data));
