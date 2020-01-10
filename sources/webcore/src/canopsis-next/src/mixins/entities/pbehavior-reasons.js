import { createNamespacedHelpers } from 'vuex';

import { PAUSE_REASONS } from '@/constants';

const { mapActions, mapGetters } = createNamespacedHelpers('pbehaviorReasons');

export default {
  computed: {
    ...mapGetters({
      pbehaviorReasonsPending: 'pending',
      pbehaviorReasons: 'pbehaviorReasons',
    }),
    pbehaviorReasonsOrDefault() {
      return this.pbehaviorReasons.length ? this.pbehaviorReasons : Object.values(PAUSE_REASONS);
    },
  },
  methods: {
    ...mapActions(['fetchPbehaviorReasons']),
  },
};
