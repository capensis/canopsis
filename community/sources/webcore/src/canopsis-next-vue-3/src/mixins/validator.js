import { validate as validateValue, defineRule } from 'vee-validate';

import { FORM_CONTEXT_KEY } from '@/hooks/form';

const getFieldMessage = (result, name) => result.errors[0] || `${name} is not valid.`;

/**
 * Provides a vee-validate 2-like validator facade for legacy Options API components.
 *
 * This mixin exists only as a migration bridge for old components that still use
 * `this.$validator` or `this.errors`. New code should use `useFormContext`,
 * `useValidator` or field-level validation hooks directly.
 *
 * @deprecated Replace with explicit vee-validate 4 form/field APIs during form migration.
 */
export const validatorMixin = {
  inject: {
    formContext: {
      from: FORM_CONTEXT_KEY,
      default: null,
    },
  },

  data() {
    return {
      validatorFields: {},
    };
  },

  computed: {
    errors() {
      return this.validator.errors;
    },

    validator() {
      const getFieldError = name => this.formContext?.errors?.value?.[name];
      const setFieldError = (name, message) => this.formContext?.setFieldError?.(name, message);
      const clear = name =>
        name ? setFieldError(name, undefined) : this.formContext?.setErrors?.({});

      return {
        errors: {
          any: () => Object.keys(this.formContext?.errors?.value ?? {}).length > 0,
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
          find: ({ name } = {}) => this.validatorFields[name] ?? null,
          items: Object.values(this.validatorFields),
        },
        attach: field => this.attachValidatorField(field),
        detach: name => this.detachValidatorField(name),
        extend: (name, rule) => defineRule(name, rule),
        reset: ({ name } = {}) => clear(name),
        validate: name => this.validateValidatorField(name),
        validateAll: () => this.validateAllValidatorFields(),
      };
    },
  },

  beforeCreate() {
    Object.defineProperty(this, '$validator', {
      configurable: true,
      get: () => this.validator,
    });
  },

  methods: {
    attachValidatorField(field) {
      if (!field?.name) {
        return;
      }

      this.validatorFields = {
        ...this.validatorFields,
        [field.name]: field,
      };
    },

    detachValidatorField(name) {
      const fields = { ...this.validatorFields };

      delete fields[name];

      this.validatorFields = fields;
      this.errors.clear(name);
    },

    async validateValidatorField(name) {
      const field = this.validatorFields[name];

      if (!field) {
        const { valid } = (await this.formContext?.validateField?.(name)) ?? { valid: true };

        return valid;
      }

      const result = await validateValue(field.getter?.() ?? undefined, field.rules, { name });

      if (result.valid) {
        this.errors.clear(name);
      } else {
        this.errors.add({
          field: name,
          msg: getFieldMessage(result, name),
        });
      }

      return result.valid;
    },

    async validateAllValidatorFields() {
      const results = await Promise.all(
        Object.keys(this.validatorFields).map(name => this.validateValidatorField(name)),
      );
      const { valid } = (await this.formContext?.validate?.()) ?? { valid: true };

      return valid && results.every(Boolean);
    },
  },
};

export default validatorMixin;
