import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

/**
 * Convert pbehavior type data to type form
 *
 * @param {Object} type
 * @return {Object}
 */
export function pbehaviorTypeToForm(type = {}) {
  return {
    name: type.name || '',
    description: type.description || '',
    type: type.type || PBEHAVIOR_TYPE_TYPES.active,
    priority: type.priority || '',
    iconName: type.icon_name || '',
  };
}

/**
 * Convert type form to pbehavior type data
 *
 * @param {Object} typeForm
 * @return {Object}
 */
export function formToPbehaviorType(typeForm = {}) {
  const { iconName, ...form } = typeForm;

  return {
    icon_name: iconName,
    ...form,
  };
}
