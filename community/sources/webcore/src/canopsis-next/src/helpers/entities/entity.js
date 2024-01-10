import { isNull, uniq } from 'lodash';

import {
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  ENTITY_EVENT_BY_ACTION_TYPE,
  ENTITY_TYPES,
  WEATHER_ACK_EVENT_OUTPUT,
  WEATHER_ACTIONS_TYPES,
} from '@/constants';

import { getEntityEventIcon } from '@/helpers/icon';
import { getEntityEventColor } from '@/helpers/color';
import {
  createAckEventByEntity,
  createAssociateTicketEventByEntity,
  createCancelEventByEntity,
  createCommentEventByEntity,
  createDeclareTicketEventByEntity,
  createInvalidateEventByEntity,
  createRemoveAckEventByEntity,
  createValidateEventByEntity,
} from '@/helpers/forms/event';

/**
 * @typedef {Object} EntityAction
 * @property {string} actionType
 * @property {Entity[]} entities
 * @property {Object} payload
 */

/**
 * Check if entity is component
 *
 * @param {EntityType} type
 * @returns {boolean}
 */
export const isEntityComponentType = type => type === ENTITY_TYPES.component;

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
    alarm_id: alarmId,
    alarm_display_name: alarmDisplayName,
    assigned_instructions: assignedInstructions,
    pbh_origin_icon: pbhOriginIcon,
  } = entity;

  const paused = pbhOriginIcon !== '';
  const stateIsOk = state?.val === ENTITIES_STATES.ok;
  const statusIsCancelled = status?.val === ENTITIES_STATUSES.cancelled;
  const hasAlarm = !!alarmId;

  if (statusIsCancelled) {
    return false;
  }

  switch (actionType) {
    case WEATHER_ACTIONS_TYPES.entityAck:
      return isNull(ack);
    case WEATHER_ACTIONS_TYPES.entityAckRemove:
      return !isNull(ack);

    case WEATHER_ACTIONS_TYPES.entityValidate:
    case WEATHER_ACTIONS_TYPES.entityInvalidate:
      return !stateIsOk && state?.val === ENTITIES_STATES.major;

    case WEATHER_ACTIONS_TYPES.entityCancel:
      return !stateIsOk && alarmDisplayName && (!status || !statusIsCancelled);

    case WEATHER_ACTIONS_TYPES.entityPlay:
      return paused;

    case WEATHER_ACTIONS_TYPES.entityPause:
      return !paused;

    case WEATHER_ACTIONS_TYPES.executeInstruction:
      return !stateIsOk && !!assignedInstructions?.length;

    case WEATHER_ACTIONS_TYPES.declareTicket:
    case WEATHER_ACTIONS_TYPES.entityAssocTicket:
      return hasAlarm;

    case WEATHER_ACTIONS_TYPES.entityComment:
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

/**
 * Get icon by entity type
 *
 * @param {EntityType} type
 * @returns {string}
 */
export const getIconByEntityType = type => (
  type === ENTITY_TYPES.service
    ? '$vuetify.icons.engineering'
    : 'person'
);

/**
* Is disabled action for entity by actions requests
*
* @param {string} entityId
* @param {Object[]} actionsRequests
* @param {string} actionType
* @returns {boolean}
*/
export const isDisabledActionForEntityByActionsRequests = (entityId, actionType, actionsRequests) => {
  switch (actionType) {
    case WEATHER_ACTIONS_TYPES.declareTicket:
    case WEATHER_ACTIONS_TYPES.entityAssocTicket:
    case WEATHER_ACTIONS_TYPES.entityComment:
    case WEATHER_ACTIONS_TYPES.executeInstruction:
      return false;
    default:
      return actionsRequests.some(({ actionType: requestActionType, entitiesIds: requestEntitiesIds }) => (
        requestActionType === actionType && requestEntitiesIds.includes(entityId)
      ));
  }
};
