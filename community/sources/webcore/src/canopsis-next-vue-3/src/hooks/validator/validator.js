import {
  useFormErrors,
  useSetFieldError,
  useSetFormErrors,
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
  const setFormErrors = useSetFormErrors();
  const validateField = useValidateField();
  const validateForm = useValidateForm();

  const getFieldError = (name) => formErrors.value[name];

  return {
    errors: {
      any: () => Object.keys(formErrors.value).length > 0,
      has: (name) => !!getFieldError(name),
      collect: (name) => {
        const error = getFieldError(name);

        return error ? [error] : [];
      },
      add: (errors) => {
        const normalizedErrors = Array.isArray(errors) ? errors : [errors];

        normalizedErrors.forEach(({ field, msg }) => setFieldError(field, msg));
      },
      clear: () => setFormErrors({}),
      remove: (name) => setFieldError(name, undefined),
    },
    fields: {
      find: () => null,
      items: [],
    },
    attach: () => {},
    detach: () => {},
    extend: () => {},
    reset: ({ name } = {}) => (name ? setFieldError(name, undefined) : setFormErrors({})),
    validate: async (name) => {
      const { valid } = await validateField(name);

      return valid;
    },
    validateAll: async () => {
      const { valid } = await validateForm();

      return valid;
    },
  };
};
