import { useValidator } from './validator';

/**
 * Compatibility helper for old dynamic min_value validation.
 *
 * @param {string} name Field name.
 * @returns {Object}
 */
export const useValidationAttachMinValue = (name) => {
  const validator = useValidator();

  return {
    validator,
    attachMinValueRule: () => {},
    detachMinValueRule: () => {},
    validateMinValueRule: () => validator.validate(name),
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
  const validator = useValidator();

  return {
    validator,
    asyncValidateMinValueRule: () => validator.validate(name, getValue?.()),
  };
};
