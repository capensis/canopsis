import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('remediationJobExecution');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      remediationJobsExecutions: 'items',
      getRemediationJobExecution: 'getItemById',
      remediationJobsExecutionsPending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchRemediationJobExecution: 'fetchItem',
      createRemediationJobExecution: 'create',
      cancelRemediationJobExecution: 'cancel',
    }),
  },
};
