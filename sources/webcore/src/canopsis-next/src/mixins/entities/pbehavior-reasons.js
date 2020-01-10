import { createNamespacedHelpers } from 'vuex';

import { PAUSE_REASONS } from '@/constants';

const { mapActions, mapGetters } = createNamespacedHelpers('pbehaviorReasons');

export default {
  computed: {
    ...mapGetters({
      pbehaviorReasonsPending: 'pending',
      pbehaviorReasonsData: 'pbehaviorReasons',
    }),
    pbehaviorReasons() {
      return this.pbehaviorReasonsData.length ? this.pbehaviorReasonsData : Object.values(PAUSE_REASONS);
    },
  },
  methods: {
    ...mapActions(['fetchPbehaviorReasons']),
  },
};
