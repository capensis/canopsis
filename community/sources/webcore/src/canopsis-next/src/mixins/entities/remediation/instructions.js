import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('remediationInstruction');

/**
 * @mixin
 */
export default {
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
      fetchRemediationInstructionsListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchRemediationInstructionsListWithoutStore: 'fetchListWithoutStore',
      createRemediationInstruction: 'create',
      updateRemediationInstruction: 'update',
      removeRemediationInstruction: 'remove',
      rateRemediationInstruction: 'rateInstruction',
    }),
  },
};
