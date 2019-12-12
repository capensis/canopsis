import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('filterHint');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters(['alarmFilterHints', 'entityFilterHints']),
  },
  methods: {
    ...mapActions(['fetchFilterHints']),
  },
};
