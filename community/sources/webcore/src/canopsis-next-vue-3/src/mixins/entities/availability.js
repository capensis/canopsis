import { mapActions, mapState } from 'pinia';

import { useMetricsStore } from '@/stores/modules/metrics';

export const entitiesAvailabilityMixin = {
  computed: {
    ...mapState(useMetricsStore, {
      getAvailabilityListByWidgetId: 'getAvailabilityListByWidgetId',
      getAvailabilityPendingByWidgetId: 'getAvailabilityPendingByWidgetId',
      getAvailabilityMetaByWidgetId: 'getAvailabilityMetaByWidgetId',
    }),

    availabilities() {
      return this.getAvailabilityListByWidgetId(this.widget._id);
    },

    availabilitiesPending() {
      return this.getAvailabilityPendingByWidgetId(this.widget._id);
    },

    availabilitiesMeta() {
      return this.getAvailabilityMetaByWidgetId(this.widget._id);
    },
  },

  methods: {
    ...mapActions(useMetricsStore, {
      fetchAvailabilityList: 'fetchAvailabilityList',
      createAvailabilityExport: 'createAvailabilityExport',
      fetchAvailabilityExport: 'fetchAvailabilityExport',
    }),
  },
};
