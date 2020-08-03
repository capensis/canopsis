import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

/**
 * Convert type data to type form
 *
 * @param {Object} type
 * @return {Object}
 */
export function typeToForm(type = {}) {
  return {
    name: type.name || '',
    description: type.description || '',
    type: type.type || PBEHAVIOR_TYPE_TYPES.activeState,
    priority: type.priority || '',
    iconName: type.iconName || '',
  };
}

/**
 * Convert type form to type data
 *
 * @param {Object} typeForm
 * @return {Object}
 */
export function formToType(typeForm = {}) {
  const { iconName, ...form } = typeForm;

  return {
    icon_name: iconName,
    ...form,
  };
}
