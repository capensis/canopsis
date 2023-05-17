import {
  DECLARE_TICKET_OUTPUT,
  ENTITIES_STATES,
  EVENT_ENTITY_TYPES,
  TIME_UNITS,
  WEATHER_ACK_EVENT_OUTPUT,
  WEATHER_EVENT_DEFAULT_ENTITY,
} from '@/constants';

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
 * @typedef {Object} SnoozeAction
 * @property {number} duration
 * @property {string} comment
 */

/**
 * @typedef {SnoozeAction} SnoozeActionForm
 * @property {Duration} duration
 */

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

/**
 * Convert snooze object to form snooze
 *
 * @param {SnoozeAction} snooze
 * @returns {SnoozeActionForm}
 */
export const snoozeToForm = (snooze = {}) => ({
  duration: {
    unit: snooze.duration?.unit ?? TIME_UNITS.minute,
    value: snooze.duration?.seconds ?? 1,
  },
  comment: snooze.comment ?? '',
});
