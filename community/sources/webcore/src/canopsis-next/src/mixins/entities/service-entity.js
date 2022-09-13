import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('service/entity');

export const entitiesServiceEntityMixin = {
  computed: {
    ...mapGetters({
      getServiceEntitiesListByServiceId: 'getListByServiceId',
      getServiceEntitiesPendingByServiceId: 'getPendingByServiceId',
      getServiceEntitiesMetaByServiceId: 'getMetaByServiceId',
    }),

    serviceEntities() {
      return this.getServiceEntitiesListByServiceId(this.service._id);
    },

    serviceEntitiesPending() {
      return this.getServiceEntitiesPendingByServiceId(this.service._id);
    },

    serviceEntitiesMeta() {
      return this.getServiceEntitiesMetaByServiceId(this.service._id);
    },
  },
  methods: {
    ...mapActions({
      fetchServiceEntitiesList: 'fetchList',
    }),
  },
};
