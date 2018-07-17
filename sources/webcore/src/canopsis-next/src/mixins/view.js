import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the alarms view
 */
export default {
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
    }),
  },
  computed: {
    ...mapGetters({
      getView: 'getItem',
      viewPending: 'pending',
    }),
  },
};
