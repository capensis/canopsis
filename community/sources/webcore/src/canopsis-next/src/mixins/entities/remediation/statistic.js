import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('remediationStatistic');

/**
 * @mixin
 */
export const entitiesRemediationStatisticMixin = {
  computed: {
    ...mapGetters({
      remediationMetrics: 'metrics',
      remediationMetricsPending: 'pending',
      remediationMetricsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchRemediationMetricsList: 'fetchList',
      fetchRemediationMetricsListWithPreviousParams: 'fetchListWithPreviousParams',
    }),
  },
};
