import {
  defineRule,
  useFormErrors,
  useSetFieldError,
  useSetFormErrors,
  useValidateField,
  useValidateForm,
} from 'vee-validate';

const adaptRule = (name, rule) => {
  if (typeof rule === 'function') {
    return rule;
  }

  return async (value, params, ctx) => {
    const result = await rule.validate(value, params, ctx);
    const valid = typeof result === 'object' ? result.valid !== false : !!result;

    return valid || rule.getMessage?.(ctx?.field || name, params, ctx) || false;
  };
};

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
  const clear = (name) => (name ? setFieldError(name, undefined) : setFormErrors({}));

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
      clear,
      remove: (name) => setFieldError(name, undefined),
    },
    fields: {
      find: ({ name } = {}) => (getFieldError(name) ? { name } : null),
      items: [],
    },
    attach: () => {},
    detach: (name) => setFieldError(name, undefined),
    extend: (name, rule) => defineRule(name, adaptRule(name, rule)),
    reset: ({ name } = {}) => clear(name),
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
