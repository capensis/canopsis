import { useValidator } from '@/hooks/validator/validator';

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
  setup() {
    const validator = useValidator();

    return {
      $validator: validator,
      errors: validator.errors,
    };
  },
};

export default validatorMixin;
