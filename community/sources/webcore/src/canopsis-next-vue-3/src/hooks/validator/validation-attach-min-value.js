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
