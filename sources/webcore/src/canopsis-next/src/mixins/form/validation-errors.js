import { has } from 'lodash';

export default ({ formField = 'form' } = {}) => ({
  computed: {
    isFormChanged() {
      return Object.keys(this.fields).some(key => this.fields[key].changed);
    },
  },
  methods: {
    setFormErrors(err = {}) {
      const existFieldErrors = Object.entries(err).filter(([field]) => has(this[formField], field));

      if (existFieldErrors.length) {
        this.errors.add(existFieldErrors.map(([field, msg]) => ({ field, msg })));
      }
    },
  },
});
