import { isParent } from '@/helpers/vue-base';

export default {
  computed: {
    hasAnyError() {
      if (this.$validator && this.errors.any()) {
        const fields = this.$validator.fields.items.filter(({ vm }) => isParent(vm, this));

        return fields.some(field => this.errors.has(field.name));
      }

      return false;
    },

    validationHeaderClass() {
      return {
        'validation-header': true,
        'error--text': this.hasAnyError,
      };
    },
  },
};
