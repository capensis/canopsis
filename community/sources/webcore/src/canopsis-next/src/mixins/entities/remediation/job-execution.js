import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('remediationJobExecution');

/**
 * @mixin
 */
export const entitiesRemediationJobExecutionMixin = {
  methods: {
    ...mapActions({
      createRemediationJobExecution: 'create',
      cancelRemediationJobExecution: 'cancel',
    }),
  },
};
