import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

/**
 * Convert planning type data to type form
 *
 * @param {Object} type
 * @return {Object}
 */
export function planningTypeToForm(type = {}) {
  return {
    name: type.name || '',
    description: type.description || '',
    type: type.type || PBEHAVIOR_TYPE_TYPES.activeState,
    priority: type.priority || '',
    iconName: type.iconName || '',
  };
}

/**
 * Convert type form to planning type data
 *
 * @param {Object} typeForm
 * @return {Object}
 */
export function formToPlanningType(typeForm = {}) {
  const { iconName, ...form } = typeForm;

  return {
    icon_name: iconName,
    ...form,
  };
}
