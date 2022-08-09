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
      getAlarmDetailsQuery: 'getQuery',
    }),
  },
  methods: {
    ...mapActions({
      fetchAlarmDetails: 'fetchItem',
      fetchAlarmsDetailsList: 'fetchList',
      updateAlarmDetailsQuery: 'updateQuery',
      removeAlarmDetailsQuery: 'removeQuery',
    }),
  },
};
