import { createNamespacedHelpers } from 'vuex';

import { entitiesInfoMixin } from '@/mixins/entities/info';

const { mapGetters, mapActions } = createNamespacedHelpers('service');

export default {
  mixins: [entitiesInfoMixin],
  computed: {
    ...mapGetters({
      getServicesListByWidgetId: 'getListByWidgetId',
      getServicesPendingByWidgetId: 'getPendingByWidgetId',
      getServicesErrorByWidgetId: 'getErrorByWidgetId',
      getService: 'getItem',
    }),

    services() {
      return this.getServicesListByWidgetId(this.widget._id);
    },

    servicesPending() {
      return this.getServicesPendingByWidgetId(this.widget._id);
    },

    servicesError() {
      return this.getServicesErrorByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchServiceItemWithoutStore: 'fetchItemWithoutStore',
      fetchServicesList: 'fetchList',
      createService: 'create',
      editService: 'update',
      removeService: 'remove',
    }),
  },
};
