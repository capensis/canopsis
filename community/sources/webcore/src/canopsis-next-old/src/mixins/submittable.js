import { validationErrorsMixinCreator } from '@/mixins/form';

/**
 * Create submittable mixin for components
 *
 * @param {string} [formField = 'form'] - name of form field which we will use in the validationErrorsMixinCreator mixin
 * @param {string} [method = 'submit'] - name of submit method which we will wrap into try catch block
 * @param {string} [property = 'submitting'] - property name for submitting flag functional
 * @param {string} [computedProperty = 'isDisabled'] - computed property name for buttons disabling
 * @param {boolean} [withTimeout = true] - property for timeout enabling
 * @returns {{data(): *, computed: {}, created(): void}}
 */
export const submittableMixinCreator = ({
  formField = 'form',
  method = 'submit',
  property = 'submitting',
  computedProperty = 'isDisabled',
  withTimeout = true,
} = {}) => ({
  mixins: [validationErrorsMixinCreator({ formField })],
  data() {
    return {
      [property]: false,
    };
  },
  created() {
    const sourceSubmit = this[method];

    if (sourceSubmit) {
      const submit = async (...args) => {
        try {
          this[property] = true;

          await sourceSubmit.apply(this, args);
        } catch (err) {
          const wasSet = this.setFormErrors(err);

          if (!wasSet) {
            console.error(err);

            const message = Object.values(err).join('\n');

            this.$popups.error({ text: message || err.details || this.$t('errors.default') });
          }
        } finally {
          this[property] = false;
        }
      };

      /**
       * If `withTimeout` is true, a timeout is set to call `submitHandler` with the provided arguments after 0 ms.
       * Otherwise, `submitHandler` is called directly.
       */
      this[method] = withTimeout ? (...args) => setTimeout(() => submit.apply(this, args), 0) : submit;
    }
  },
  computed: {
    [computedProperty]() {
      if (this.errors) {
        return this.errors.any() || this[property];
      }

      return this[property];
    },
  },
});

export default submittableMixinCreator;
