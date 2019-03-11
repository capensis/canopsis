import { createNamespacedHelpers } from 'vuex';

import authMixin from './auth';

const { mapActions } = createNamespacedHelpers('event');

export default {
  mixins: [authMixin],
  methods: {
    ...mapActions({
      createEventAction: 'create',
    }),
  },
};
