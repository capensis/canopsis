import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('healthcheck');

export const entitiesHealthcheckMixin = {
  computed: {
    ...mapGetters({
      healthcheckPending: 'pending',
      services: 'services',
      engines: 'engines',
      maxQueueLength: 'maxQueueLength',
      hasInvalidEnginesOrder: 'hasInvalidEnginesOrder',
      healthcheckError: 'error',
    }),
  },
  methods: {
    ...mapActions({
      fetchHealthcheckStatus: 'fetchStatus',
    }),
  },
};
