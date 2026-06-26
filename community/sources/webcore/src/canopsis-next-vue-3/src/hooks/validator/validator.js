import {
  useFormErrors,
  useSetFieldError,
  useValidateField,
  useValidateForm,
} from 'vee-validate';

/**
 * Returns a small compatibility facade backed by vee-validate form state.
 *
 * @returns {Object}
 */
export const useValidator = () => {
  const formErrors = useFormErrors();
  const setFieldError = useSetFieldError();
  const validateField = useValidateField();
  const validateForm = useValidateForm();

  const getFieldError = name => formErrors.value[name];

  return {
    errors: {
      has: name => !!getFieldError(name),
      collect: name => {
        const error = getFieldError(name);

        return error ? [error] : [];
      },
      remove: name => setFieldError(name, undefined),
    },
    validate: async name => {
      const { valid } = await validateField(name);

      return valid;
    },
    validateAll: async () => {
      const { valid } = await validateForm();

      return valid;
    },
  };
};
