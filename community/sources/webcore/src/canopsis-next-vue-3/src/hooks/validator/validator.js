import { isFunction, isPlainObject } from 'lodash';
import { defineRule, useFormErrors } from 'vee-validate';

import { toLiteralVeeFieldPath, useFormContext } from '@/hooks/form';

/**
 * Converts a legacy validator rule definition to a vee-validate compatible rule function.
 *
 * @param {string} name Rule name used as a fallback field label in validation messages.
 * @param {Function|Object} rule Rule function or legacy rule object with validate/getMessage methods.
 * @returns {Function} Vee-validate compatible validation rule.
 */
const adaptRule = (name, rule) => {
  if (isFunction(rule)) {
    return rule;
  }

  return async (value, params, ctx) => {
    const result = await rule.validate(value, params, ctx);
    const valid = isPlainObject(result) ? result.valid !== false : !!result;

    return valid || rule.getMessage?.(ctx?.field || name, params, ctx) || false;
  };
};

/**
 * Returns a small compatibility facade backed by vee-validate form state.
 *
 * @returns {Object}
 */
export const useValidator = () => {
  const form = useFormContext();
  const formErrors = useFormErrors();

  const getFieldError = (name) => {
    const errors = formErrors.value;

    return errors?.[name] ?? errors?.[toLiteralVeeFieldPath(name)];
  };
  const setFieldError = (name, message) => form.setFieldError(name, message);
  const clear = name => (name ? setFieldError(name, undefined) : form.setErrors({}));

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
