<template lang="pug">
  v-card-text
    remediation-instructions-list(
      :remediationInstructions="remediationInstructions",
      :pending="remediationInstructionsPending",
      :totalItems="remediationInstructionsMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedRemediationInstructionModal",
      @assign="showEditPatternModal",
      @remove="showRemoveRemediationInstructionModal",
      @edit="showEditRemediationInstructionModal"
    )
</template>

<script>
import { isEmpty, isEqual } from 'lodash';

import { MODALS } from '@/constants';

import entitiesRemediationInstructionMixin from '@/mixins/entities/remediation/instruction';
import remediationQueryMixin from '@/mixins/remediation/query';

import RemediationInstructionsList from './remediation-instructions-list.vue';

export default {
  components: { RemediationInstructionsList },
  inject: ['$validator'],
  mixins: [
    entitiesRemediationInstructionMixin,
    remediationQueryMixin,
  ],
  watch: {
    query(query, oldQuery) {
      if (!isEqual(query, oldQuery)) {
        this.fetchList();
      }
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchRemediationInstructionsList({ params: this.getQuery() });
    },

    async tryRemoveRemediationInstruction(remediationInstructionId) {
      try {
        await this.removeRemediationInstruction({ id: remediationInstructionId });
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    showEditRemediationInstructionModal() {
    },

    showConfirmEditRunningRemediationInstructionModal(action) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          text: 'There are instruction in progress assigned to this pattern. Would you like to cancel them?',
          action,
        },
      });
    },

    showEditPatternModal(remediationInstruction) {
      this.$modals.show({
        name: MODALS.createEventFilterRulePattern,
        config: {
          isSimplePattern: true,
          pattern: remediationInstruction.pattern,
          action: (pattern) => {
            const data = { id: remediationInstruction._id, data: { ...remediationInstruction, pattern } };

            if (remediationInstruction.running) {
              this.showConfirmEditRunningRemediationInstructionModal(() => this.updateRemediationInstruction(data));
            } else if (isEmpty(remediationInstruction.pattern)) {
              this.updateRemediationInstruction(data);
            }
          },
        },
      });
    },

    showRemoveRemediationInstructionModal(remediationInstruction) {
      const action = async () => {
        await this.tryRemoveRemediationInstruction(remediationInstruction._id);
        await this.fetchList();
      };

      if (remediationInstruction.running) {
        this.showConfirmEditRunningRemediationInstructionModal(action);
      } else {
        this.$modals.show({
          name: MODALS.confirmation,
          config: {
            action,
          },
        });
      }
    },

    showRemoveSelectedRemediationInstructionModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id: id }) => this.tryRemoveRemediationInstruction(id)));

            await this.fetchList();
          },
        },
      });
    },
  },
};
</script>
