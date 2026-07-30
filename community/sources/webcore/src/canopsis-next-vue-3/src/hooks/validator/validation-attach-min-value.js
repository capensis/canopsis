import { useValidationAttachRule } from './validation-attach-rule';

const MIN_VALUE_RULE = 'min_value:1';

/**
 * Compatibility helper for old dynamic min_value validation.
 *
 * @param {string} name Field name.
 * @returns {Object}
 */
export const useValidationAttachMinValue = (name) => {
  const { validator, attachRule, detachRule, validateRule } = useValidationAttachRule(
    name,
    MIN_VALUE_RULE,
    () => 0,
  );

  return {
    validator,
    attachMinValueRule: attachRule,
    detachMinValueRule: detachRule,
    validateMinValueRule: validateRule,
  };
};

/**
 * Compatibility helper for old dynamic min_value validation by field value getter.
 *
 * @param {string} name Field name.
 * @param {Function} getValue Field value getter.
 * @returns {Object}
 */
export const useValidationAttachMinValueForField = (name, getValue) => {
  const { validator, attachRule, validateRule } = useValidationAttachRule(
    name,
    MIN_VALUE_RULE,
    () => getValue?.() ?? 0,
  );

  attachRule();

  return {
    validator,
    asyncValidateMinValueRule: validateRule,
  };
};
