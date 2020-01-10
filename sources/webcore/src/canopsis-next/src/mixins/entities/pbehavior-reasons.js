import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('pbehaviorReasons');

export default {
  computed: {
    ...mapGetters({
      pbehaviorReasonsPending: 'pending',
      pbehaviorReasons: 'pbehaviorReasons',
    }),
  },
  methods: {
    ...mapActions(['fetchPbehaviorReasons']),
  },
};
