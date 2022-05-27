import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('alarm/details');

/**
 * @mixin
 */
export const entitiesAlarmDetailsMixin = {
  computed: {
    ...mapGetters({
      getAlarmDetailsItem: 'getItem',
    }),
  },
  methods: {
    ...mapActions({
      fetchAlarmItemDetails: 'fetchItem',
    }),
  },
};
