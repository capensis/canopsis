import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('alarm/details');

/**
 * @mixin
 */
export const entitiesAlarmDetailsMixin = {
  computed: {
    ...mapGetters({
      getAlarmDetailsItem: 'getItem',
      getAlarmDetailsPending: 'getPending',
    }),
  },
  methods: {
    ...mapActions({
      fetchAlarmItemDetails: 'fetchItem',
    }),
  },
};
