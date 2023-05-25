import { isNull, uniq } from 'lodash';

import {
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  ENTITY_EVENT_BY_ACTION_TYPE,
  WEATHER_ACK_EVENT_OUTPUT,
  WEATHER_ACTIONS_TYPES,
} from '@/constants';

import {
  createAckEventByEntity,
  createAssociateTicketEventByEntity,
  createCancelEventByEntity,
  createCommentEventByEntity,
  createDeclareTicketEventByEntity,
  createInvalidateEventByEntity,
  createRemoveAckEventByEntity,
  createValidateEventByEntity,
} from '@/helpers/entities/event/form';

import { getEntityEventIcon } from './icons';
import { getEntityEventColor } from './color';

/**
 * @typedef {Object} EntityAction
 * @property {string} actionType
 * @property {Entity[]} entities
 * @property {Object} payload
 */

/**
 * Check is action available for the entity
 *
 * @param {String} actionType
 * @param {Entity} entity
 * @return {boolean}
 */
export const isActionTypeAvailableForEntity = (actionType, entity) => {
  const {
    state,
    ack,
    status,
    alarm_display_name: alarmDisplayName,
    assigned_instructions: assignedInstructions,
    pbh_origin_icon: pbhOriginIcon,
  } = entity;

  const paused = pbhOriginIcon !== '';
  const stateIsOk = state?.val === ENTITIES_STATES.ok;
  const statusIsCancelled = status?.val === ENTITIES_STATUSES.cancelled;

  if (stateIsOk || statusIsCancelled) {
    return false;
  }

  switch (actionType) {
    case WEATHER_ACTIONS_TYPES.entityAck:
      return !stateIsOk && isNull(ack);
    case WEATHER_ACTIONS_TYPES.entityAckRemove:
      return !stateIsOk && !isNull(ack);

    case WEATHER_ACTIONS_TYPES.entityValidate:
    case WEATHER_ACTIONS_TYPES.entityInvalidate:
      return state?.val === ENTITIES_STATES.major;

    case WEATHER_ACTIONS_TYPES.entityCancel:
      return alarmDisplayName && (!status || !statusIsCancelled);

    case WEATHER_ACTIONS_TYPES.entityPlay:
      return paused;

    case WEATHER_ACTIONS_TYPES.entityPause:
      return !paused;

    case WEATHER_ACTIONS_TYPES.executeInstruction:
      return !!assignedInstructions?.length;

    case WEATHER_ACTIONS_TYPES.declareTicket:
    case WEATHER_ACTIONS_TYPES.entityAssocTicket:
    default:
      return true;
  }
};

/**
 * Get all available action types for entity
 *
 * @param {Entity} entity
 * @param {string[]} actionTypes
 * @returns {string[]}
 */
export const getAvailableEntityActionsTypes = (
  entity,
  actionTypes = [
    WEATHER_ACTIONS_TYPES.entityComment,
    WEATHER_ACTIONS_TYPES.executeInstruction,
    WEATHER_ACTIONS_TYPES.entityAck,
    WEATHER_ACTIONS_TYPES.entityAckRemove,
    WEATHER_ACTIONS_TYPES.entityAssocTicket,
    WEATHER_ACTIONS_TYPES.declareTicket,
    WEATHER_ACTIONS_TYPES.entityValidate,
    WEATHER_ACTIONS_TYPES.entityInvalidate,
    WEATHER_ACTIONS_TYPES.entityPlay,
    WEATHER_ACTIONS_TYPES.entityPause,
    WEATHER_ACTIONS_TYPES.entityCancel,
  ],
) => actionTypes.filter(actionType => isActionTypeAvailableForEntity(actionType, entity));

/**
 * Convert entity action type to action object
 *
 * @param {string} type
 * @return {Object}
 */
const convertEntityActionTypeToAction = (type) => {
  const eventType = ENTITY_EVENT_BY_ACTION_TYPE[type];

  return {
    type,
    icon: getEntityEventIcon(eventType),
    color: getEntityEventColor(eventType),
  };
};

/**
 * Get all available actions for entity
 *
 * @param {Entity} entity
 * @returns {Object[]}
 */
export const getAvailableActionsByEntity = entity => getAvailableEntityActionsTypes(entity)
  .map(convertEntityActionTypeToAction);

/**
 * Get all available actions for entities
 *
 * @param {Entity[]} entities
 * @param {string[]} actionTypes
 * @returns {Object[]}
 */
export const getAvailableActionsByEntities = (
  entities = [],
  actionTypes = [
    WEATHER_ACTIONS_TYPES.entityComment,
    WEATHER_ACTIONS_TYPES.entityAck,
    WEATHER_ACTIONS_TYPES.entityAckRemove,
    WEATHER_ACTIONS_TYPES.entityAssocTicket,
    WEATHER_ACTIONS_TYPES.declareTicket,
    WEATHER_ACTIONS_TYPES.entityValidate,
    WEATHER_ACTIONS_TYPES.entityInvalidate,
    WEATHER_ACTIONS_TYPES.entityPlay,
    WEATHER_ACTIONS_TYPES.entityPause,
    WEATHER_ACTIONS_TYPES.entityCancel,
  ],
) => {
  const types = entities.reduce(
    (acc, entity) => acc.concat(getAvailableEntityActionsTypes(entity, actionTypes)),
    [],
  );

  return uniq(types).map(convertEntityActionTypeToAction);
};

/**
 * Convert action to events by type
 *
 * @param {string} actionType
 * @param {Entity} entity
 * @param {Object} payload
 * @return {Event[]}
 */
export const convertActionToEvents = ({ actionType, entity, payload }) => {
  switch (actionType) {
    case WEATHER_ACTIONS_TYPES.entityAck:
      return [
        createAckEventByEntity({ entity, output: WEATHER_ACK_EVENT_OUTPUT.ack }),
      ];
    case WEATHER_ACTIONS_TYPES.entityAckRemove:
      return [
        createRemoveAckEventByEntity({ entity, output: payload.output }),
      ];
    case WEATHER_ACTIONS_TYPES.entityComment:
      return [
        createCommentEventByEntity({ entity, output: payload.output }),
      ];
    case WEATHER_ACTIONS_TYPES.entityCancel:
      return [
        createCancelEventByEntity({ entity, output: payload.output }),
      ];
    case WEATHER_ACTIONS_TYPES.entityAssocTicket:
      return [
        createAckEventByEntity({ entity, output: WEATHER_ACK_EVENT_OUTPUT.ack }),
        createAssociateTicketEventByEntity({ entity, ...payload }),
      ];
    case WEATHER_ACTIONS_TYPES.entityValidate:
      return [
        createAckEventByEntity({ entity, output: WEATHER_ACK_EVENT_OUTPUT.validateOk }),
        createValidateEventByEntity({ entity }),
      ];
    case WEATHER_ACTIONS_TYPES.entityInvalidate:
      return [
        createAckEventByEntity({ entity, output: WEATHER_ACK_EVENT_OUTPUT.ack }),
        createInvalidateEventByEntity({ entity }),
      ];
    case WEATHER_ACTIONS_TYPES.declareTicket:
      return [
        createDeclareTicketEventByEntity({ entity }),
      ];
  }

  return [];
};
