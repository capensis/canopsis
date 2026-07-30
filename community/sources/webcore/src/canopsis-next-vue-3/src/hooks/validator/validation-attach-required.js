import { nextTick, onBeforeUnmount, onMounted, unref } from 'vue';

import { useValidationAttachRule } from './validation-attach-rule';

const REQUIRED_RULE = 'required';

/**
 * Hook for add and remove validation rule for field
 *
 * @param {string} name
 * @return {Object}
 */
export const useValidationAttachRequired = (name) => {
  const { validator, attachRule, detachRule, validateRule } = useValidationAttachRule(
    name,
    REQUIRED_RULE,
  );

  const resetRequiredRule = () => validator?.reset?.({ name: unref(name) });

  return {
    validator,
    attachRequiredRule: attachRule,
    validateRequiredRule: validateRule,
    resetRequiredRule,
    detachRequiredRule: detachRule,
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
