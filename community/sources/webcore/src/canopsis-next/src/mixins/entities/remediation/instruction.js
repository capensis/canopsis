import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('remediationInstruction');

/**
 * @mixin
 */
export const entitiesRemediationInstructionMixin = {
  computed: {
    ...mapGetters({
      remediationInstructions: 'items',
      remediationInstructionsPending: 'pending',
      remediationInstructionsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchRemediationInstructionsList: 'fetchList',
      fetchRemediationInstructionWithoutStore: 'fetchItemWithoutStore',
      fetchRemediationInstructionsListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchRemediationInstructionsListWithoutStore: 'fetchListWithoutStore',
      createRemediationInstruction: 'create',
      updateRemediationInstruction: 'update',
      removeRemediationInstruction: 'remove',
      rateRemediationInstruction: 'rateInstruction',
    }),
  },
};
