import { nextTick, onBeforeUnmount, onMounted } from 'vue';

import { useComponentInstance } from '../vue';

import { useValidator } from './validator';

/**
 * Hook for add and remove min_value validation rule for field
 *
 * @param {string} name
 * @return {Object}
 */
export const useValidationAttachMinValue = (name) => {
  const validator = useValidator();
  const instance = useComponentInstance();

  /**
   * Attaches a min_value validation rule to the field if it doesn't already exist
   *
   * @param {Function} getter - Function that returns the field's current value for validation
   * @param {number} [minValue=1] - Minimum allowed value
   */
  const attachMinValueRule = (getter, minValue = 1) => {
    const oldField = validator?.fields?.find?.({ name });

    if (!oldField) {
      validator?.attach?.({
        name,
        rules: `min_value:${minValue}`,
        getter,
        vm: instance,
      });
    }
  };

  /**
   * Validates the min_value rule for the field
   *
   * @returns {Promise|undefined} Promise resolving to validation result, or undefined if validator not available
   */
  const validateMinValueRule = () => validator?.validate?.(name);

  /**
   * Resets the validation state for the field
   */
  const resetMinValueRule = () => validator?.reset?.({ name });

  /**
   * Detaches the validation rule from the field
   */
  const detachMinValueRule = () => validator?.detach?.(name);

  return {
    validator,
    attachMinValueRule,
    validateMinValueRule,
    resetMinValueRule,
    detachMinValueRule,
  };
};

/**
 * Hook for attaching, validating, and managing min_value validation rules for a specific field
 *
 * @param {string} name - The name of the field to validate
 * @param {Function} getter - Function that returns the field's value for validation
 * @param {boolean} [validateOnMount=true] - Whether to validate the field on mount
 * @param {number} [minValue=1] - Minimum allowed value
 * @returns {Object} An object containing validation control functions
 * @property {Function} attachMinValueRule - Attaches the min_value validation rule to the field
 * @property {Function} detachMinValueRule - Detaches the min_value validation rule from the field
 * @property {Function} validateMinValueRule - Synchronously validates the min_value rule
 * @property {Function} asyncValidateMinValueRule - Asynchronously validates the min_value rule on the next tick
 */
export const useValidationAttachMinValueForField = (name, getter, validateOnMount = true, minValue = 1) => {
  const { attachMinValueRule, detachMinValueRule, validateMinValueRule } = useValidationAttachMinValue(name);

  const asyncValidateMinValueRule = () => nextTick(validateMinValueRule);

  onMounted(() => {
    attachMinValueRule(getter, minValue);

    if (validateOnMount) {
      validateMinValueRule();
    }
  });

  onBeforeUnmount(detachMinValueRule);

  return {
    attachMinValueRule,
    detachMinValueRule,
    validateMinValueRule,
    asyncValidateMinValueRule,
  };
};
