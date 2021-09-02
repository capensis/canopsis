import { validationChildrenMixin } from './validation-children';

export const formValidationHeaderMixin = {
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

export default formValidationHeaderMixin;
