import validationChildrenMixin from './validation-children';

export default {
  mixins: [validationChildrenMixin],
  computed: {
    hasAnyError() {
      return this.hasChildrenError;
    },

    validationHeaderClass() {
      return {
        'validation-header': true,
        'error--text': this.hasAnyError,
      };
    },
  },
};
