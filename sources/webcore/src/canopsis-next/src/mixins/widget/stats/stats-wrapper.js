import { CANOPSIS_STACK } from '@/constants';

import entitiesInfoMixin from '@/mixins/entities/info';

export default {
  mixins: [entitiesInfoMixin],
  data() {
    return {
      serverErrorMessage: null,
    };
  },
  computed: {
    stackError() {
      return this.stack === CANOPSIS_STACK.python;
    },

    hasError() {
      return this.stackError || this.serverErrorMessage;
    },

    errorMessage() {
      if (this.stackError) {
        return this.$t('errors.statsWrongStackError');
      }

      if (this.serverErrorMessage) {
        return this.$t('errors.statsRequestProblem');
      }

      return '';
    },
  },
};
