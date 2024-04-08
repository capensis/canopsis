import { useComponentInstance } from '../vue';

import { useInjectValidator } from './useValidationChildren';

/**
 * Hook for add and remove validation rule for field
 *
 * @param {string} name
 * @return {Object}
 */
export const useValidationAttachRequired = (name) => {
  const validator = useInjectValidator();
  const instance = useComponentInstance();

  const attachRequiredRule = (getter) => {
    const oldField = validator.fields.find({ name });

    if (!oldField) {
      validator.attach({
        name,
        rules: 'required:true',
        getter,
        vm: instance,
      });
    }
  };
  const validateRequiredRule = () => validator.validate(name);
  const detachRequiredRule = () => validator.detach(name);

  return {
    validator,
    attachRequiredRule,
    detachRequiredRule,
    validateRequiredRule,
  };
};
