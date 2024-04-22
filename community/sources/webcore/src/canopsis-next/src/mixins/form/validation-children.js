import { isParent } from '@/helpers/vue-base';

/**
 * @deprecated Should be used useValidationChildren
 */
export const validationChildrenMixin = {
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

    formValidated() {
      if (this.$validator) {
        return this.$validator.fields.items.some(({ flags }) => flags.validated);
      }

      return false;
    },
  },
  methods: {
    validateChildren() {
      return this.$validator.validateAll(this.childrenFields);
    },
  },
};

export default validationChildrenMixin;
