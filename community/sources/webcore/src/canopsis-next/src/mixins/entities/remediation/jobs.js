import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('remediationJob');

/**
 * @mixin
 */
export const entitiesRemediationJobsMixin = {
  computed: {
    ...mapGetters({
      remediationJobs: 'items',
      remediationJobsPending: 'pending',
      remediationJobsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchRemediationJobsList: 'fetchList',
      fetchRemediationJobsListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchRemediationJobsListWithoutStore: 'fetchListWithoutStore',
      createRemediationJob: 'create',
      updateRemediationJob: 'update',
      removeRemediationJob: 'remove',
    }),
  },
};

export default entitiesRemediationJobsMixin;
