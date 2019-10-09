import { CANOPSIS_EDITION } from '@/constants';

import entitiesInfoMixin from '@/mixins/entities/info';

export default {
  mixins: [entitiesInfoMixin],
  data() {
    return {
      serverErrorMessage: null,
    };
  },
  computed: {
    editionError() {
      return this.edition === CANOPSIS_EDITION.core;
    },

    hasError() {
      return this.editionError || this.serverErrorMessage;
    },

    errorMessage() {
      if (this.editionError) {
        return this.$t('errors.statsWrongEditionError');
      }

      if (this.serverErrorMessage) {
        return this.$t('errors.statsRequestProblem');
      }

      return '';
    },
  },
};
