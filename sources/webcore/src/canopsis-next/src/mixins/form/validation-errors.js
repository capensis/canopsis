import { has, keyBy } from 'lodash';

export default ({ formField = 'form' } = {}) => ({
  computed: {
    fieldByName() {
      return keyBy(this.$validator.fields.items, 'name');
    },
  },
  methods: {
    setFormErrors(err = {}) {
      const existFieldErrors = Object.entries(err)
        .filter(([field]) => has(this[formField], field) || this.fieldByName[field]);

      if (existFieldErrors.length) {
        this.errors.add(existFieldErrors.map(([field, msg]) => ({ field, msg })));
      }
    },
  },
});
