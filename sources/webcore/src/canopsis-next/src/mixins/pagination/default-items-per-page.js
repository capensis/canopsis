import { get } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { PAGINATION_LIMIT } from '@/config';

const { mapGetters } = createNamespacedHelpers('auth');

export default {
  computed: {
    ...mapGetters(['currentUser']),

    defaultItemsPerPage() {
      return get(this.currentUser, 'itemsPerPage', PAGINATION_LIMIT);
    },
  },
};
