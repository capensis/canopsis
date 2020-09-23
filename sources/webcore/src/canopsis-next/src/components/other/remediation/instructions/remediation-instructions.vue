<template lang="pug">
  v-card-text
    remediation-instructions-list(
      :remediationInstructions="remediationInstructions",
      :pending="remediationInstructionsPending",
      :totalItems="remediationInstructionsMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedRemediationInstructionModal",
      @assign-filter="showCreateFilterModal",
      @remove="showRemoveRemediationInstructionModal",
      @edit="showEditRemediationInstructionModal"
    )
</template>

<script>
import { isEqual } from 'lodash';

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

    showEditRemediationInstructionModal() {
    },

    showConfirmModalOnRunningRemediationInstruction(action) {
      return new Promise((resolve) => {
        this.$modals.show({
          name: MODALS.confirmation,
          config: {
            text: this.$t('remediationInstructions.errors.runningInstruction'),
            action: async (...args) => {
              await action(...args);
              resolve();
            },
            cancel: resolve,
          },
        });
      });
    },

    async updateRemediationInstructionFilter(remediationInstruction, filter) {
      if (isEqual(remediationInstruction.filter, filter)) {
        return;
      }

      const id = remediationInstruction._id;
      const data = { ...remediationInstruction, filter };

      if (remediationInstruction.running) {
        await this.showConfirmModalOnRunningRemediationInstruction(() => {
          this.updateRemediationInstruction({ id, data });
        });
      } else {
        await this.updateRemediationInstruction({ id, data });
      }
    },

    showCreateFilterModal(remediationInstruction) {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          filter: { filter: remediationInstruction.filter },
          hiddenFields: ['title'],
          action: async ({ filter }) => {
            await this.updateRemediationInstructionFilter(remediationInstruction, filter);
          },
        },
      });
    },

    showRemoveRemediationInstructionModal(remediationInstruction) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          text: remediationInstruction.running
            ? this.$t('remediationInstructions.errors.runningInstruction')
            : undefined,
          action: async () => {
            await this.removeRemediationInstruction({ id: remediationInstruction._id });
            await this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedRemediationInstructionModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id: id }) => this.removeRemediationInstruction({ id })));

            await this.fetchList();
          },
        },
      });
    },
  },
};
</script>
