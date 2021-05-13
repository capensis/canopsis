import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('remediationInstructionStats');

export const entitiesRemediationInstructionStatsMixin = {
  computed: {
    ...mapGetters({
      remediationInstructionStats: 'items',
      remediationInstructionStatsPending: 'pending',
      remediationInstructionStatsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchRemediationInstructionStatsList: 'fetchList',
    }),
  },
};

