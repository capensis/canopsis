import { createNamespacedHelpers } from 'vuex';

import { CANOPSIS_EDITION } from '@/constants';

const { mapGetters: mapInfoGetters } = createNamespacedHelpers('info');


export default {
  data() {
    return {
      serverErrorMessage: null,
    };
  },
  computed: {
    ...mapInfoGetters(['edition']),

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
