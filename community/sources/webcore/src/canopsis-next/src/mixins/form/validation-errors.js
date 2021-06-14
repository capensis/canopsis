import { has, keyBy } from 'lodash';

export default ({ formField = 'form' } = {}) => ({
  computed: {
    fieldsByName() {
      return keyBy(this.$validator.fields.items, 'name');
    },
  },
  methods: {
    setFormErrors(err = {}) {
      const existFieldErrors = Object.entries(err)
        .filter(([field]) => this.fieldsByName[field] || has(this[formField], field));

      if (existFieldErrors.length) {
        this.errors.add(existFieldErrors.map(([field, msg]) => ({ field, msg })));
      } else {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
});
