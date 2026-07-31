import { defineRule, useFormErrors } from 'vee-validate';
import { inject } from 'vue';

import { FORM_CONTEXT_KEY } from '@/hooks/form';

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
  const form = inject(FORM_CONTEXT_KEY, null);

  const getFieldError = name => formErrors.value[name];
  const setFieldError = (name, message) => form?.setFieldError?.(name, message);
  const clear = name => (name ? setFieldError(name, undefined) : form?.setErrors?.({}));

  return {
    errors: {
      any: () => Object.keys(formErrors.value).length > 0,
      has: name => !!getFieldError(name),
      collect: (name) => {
        const error = getFieldError(name);

        return error ? [error] : [];
      },
      add: (errors) => {
        const normalizedErrors = Array.isArray(errors) ? errors : [errors];

        normalizedErrors.forEach(({ field, msg }) => setFieldError(field, msg));
      },
      clear,
      remove: name => setFieldError(name, undefined),
    },
    fields: {
      find: ({ name } = {}) => (getFieldError(name) ? { name } : null),
      items: [],
    },
    attach: () => {},
    detach: name => setFieldError(name, undefined),
    extend: (name, rule) => defineRule(name, adaptRule(name, rule)),
    reset: ({ name } = {}) => clear(name),
    validate: async (name) => {
      const { valid } = (await form?.validateField?.(name)) ?? { valid: true };

      return valid;
    },
    validateAll: async () => {
      const { valid } = (await form?.validate?.()) ?? { valid: true };

      return valid;
    },
  };
};
