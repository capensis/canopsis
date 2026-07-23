import { nextTick, onBeforeUnmount, onMounted } from 'vue';

import { useValidator } from './validator';

/**
 * Hook for add and remove validation rule for field
 *
 * @param {string} name
 * @return {Object}
 */
export const useValidationAttachRequired = (name) => {
  const validator = useValidator();

  const attachRequiredRule = () => {};

  const validateRequiredRule = () => validator?.validate?.(name);

  const resetRequiredRule = () => validator?.reset?.({ name });

  const detachRequiredRule = () => {};

  return {
    validator,
    attachRequiredRule,
    validateRequiredRule,
    resetRequiredRule,
    detachRequiredRule,
  };
};

/**
 * Hook for attaching, validating, and managing required validation rules for a specific field
 *
 * @param {string} name - The name of the field to validate
 * @param {Function} getter - Function that returns the field's value for validation
 * @param {boolean} validateOnMount - Whether to validate the field on mount
 * @returns {Object} An object containing validation control functions
 * @property {Function} attachRequiredRule - Attaches the required validation rule to the field
 * @property {Function} detachRequiredRule - Detaches the required validation rule
 * @property {Function} validateRequiredRule - Synchronously validates the required rule
 * @property {Function} asyncValidateRequiredRule - Asynchronously validates the required rule on the next tick
 */
export const useValidationAttachRequiredForField = (name, getter, validateOnMount = true) => {
  const { attachRequiredRule, detachRequiredRule, validateRequiredRule } =
    useValidationAttachRequired(name);

  const validate = () => nextTick(validateRequiredRule);

  onMounted(() => {
    attachRequiredRule(getter);

    if (validateOnMount) {
      validateRequiredRule();
    }
  });

  onBeforeUnmount(detachRequiredRule);

  return {
    attachRequiredRule,
    detachRequiredRule,
    validateRequiredRule,
    asyncValidateRequiredRule: validate,
  };
};
