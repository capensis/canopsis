/**
 * Create submittable mixin for components
 *
 * @param {string} [method='submit'] - name of submit method which we will wrap into try catch block
 * @param {string} [property='submitting'] - property name for submitting flag functional
 * @param {string} [computedProperty='isDisabled'] - computed property name for buttons disabling
 * @returns {{data(): *, computed: {}, created(): void}}
 */
export const submittableMixinCreator = ({
  method = 'submit',
  property = 'submitting',
  computedProperty = 'isDisabled',
} = {}) => ({
  data() {
    return {
      [property]: false,
    };
  },
  created() {
    const sourceSubmit = this[method];

    if (sourceSubmit) {
      this[method] = async (...args) => {
        try {
          this[property] = true;

          await sourceSubmit.apply(this, args);
        } catch (err) {
          console.warn(err);
          const message = Object.values(err).join('\n');

          this.$popups.error({ text: message || err.details || this.$t('errors.default') });
        } finally {
          this[property] = false;
        }
      };
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
