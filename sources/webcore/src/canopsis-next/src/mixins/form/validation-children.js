import { isParent } from '@/helpers/vue-base';

export default {
  computed: {
    childrenFields() {
      return this.$validator ? this.$validator.fields.items.filter(({ vm }) => isParent(vm, this)) : [];
    },

    hasChildrenError() {
      if (this.$validator && this.errors.any()) {
        return this.childrenFields.some(field => this.errors.has(field.name));
      }

      return false;
    },
  },
};
