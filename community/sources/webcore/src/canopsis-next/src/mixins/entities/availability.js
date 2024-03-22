import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('availability');

export const entitiesAvailabilityMixin = {
  computed: {
    ...mapGetters({
      getAvailabilityListByWidgetId: 'getListByWidgetId',
      getAvailabilityPendingByWidgetId: 'getPendingByWidgetId',
      getAvailabilityMetaByWidgetId: 'getMetaByWidgetId',
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
    ...mapActions({
      fetchAvailabilityList: 'fetchList',
      createAvailabilityExport: 'createAvailabilityExport',
      fetchAvailabilityExport: 'fetchAvailabilityExport',
    }),
  },
};
