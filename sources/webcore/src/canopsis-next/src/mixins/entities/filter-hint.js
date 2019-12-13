import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('filterHint');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      filterHintsPending: 'pending',
      alarmFilterHints: 'alarmFilterHints',
      entityFilterHints: 'entityFilterHints',
    }),
  },
  methods: {
    ...mapActions(['fetchFilterHints']),
  },
};
