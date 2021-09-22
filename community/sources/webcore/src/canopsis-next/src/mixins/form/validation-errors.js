import { has, get, keyBy } from 'lodash';

export const validationErrorsMixinCreator = ({ formField = 'form' } = {}) => ({
  computed: {
    fieldsByName() {
      return keyBy(this.$validator.fields.items, 'name');
    },
  },
  methods: {
    /**
     * Get errors for exists fields in current form
     *
     * @param {any} err
     * @returns {[string, string][]}
     */
    getExistsFieldsErrors(err) {
      return Object.entries(err)
        .filter(([field]) => this.fieldsByName[field] || has(get(this, formField), field));
    },

    /**
     * Add exists fields errors to validator errors
     *
     * @param {[string, string][]} existsFieldsErrors
     */
    addExistsFieldsErrors(existsFieldsErrors) {
      this.errors.add(existsFieldsErrors.map(([field, msg]) => ({ field, msg })));
    },

    /**
     * Set form errors from response error and returns true if form errors exists
     *
     * @param {any} [err = {}]
     * @return {boolean}
     */
    setFormErrors(err = {}) {
      if (!this.$validator) {
        return false;
      }

      const existFieldErrors = this.getExistsFieldsErrors(err);

      if (existFieldErrors.length) {
        this.addExistsFieldsErrors(existFieldErrors);

        return true;
      }

      return false;
    },
  },
});
