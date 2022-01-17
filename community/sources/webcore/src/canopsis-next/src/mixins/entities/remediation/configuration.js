import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('remediationConfiguration');

/**
 * @mixin
 */
export const entitiesRemediationConfigurationMixin = {
  computed: {
    ...mapGetters({
      remediationConfigurations: 'items',
      remediationConfigurationsPending: 'pending',
      remediationConfigurationsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchRemediationConfigurationsList: 'fetchList',
      fetchRemediationConfigurationsListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchRemediationConfigurationsListWithoutStore: 'fetchListWithoutStore',
      createRemediationConfiguration: 'create',
      updateRemediationConfiguration: 'update',
      removeRemediationConfiguration: 'remove',
    }),
  },
};
