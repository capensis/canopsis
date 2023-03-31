import { ACTION_TYPES } from '@/constants';

/**
 * Check if action type is associate ticket
 *
 * @param {ActionType} type
 * @returns {boolean}
 */
export const isAssociateTicketActionType = type => type === ACTION_TYPES.assocticket;
