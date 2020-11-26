<template lang="pug">
  v-card-text
    remediation-instructions-list(
      :remediation-instructions="remediationInstructions",
      :pending="remediationInstructionsPending",
      :totalItems="remediationInstructionsMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedRemediationInstructionModal",
      @assign-patterns="showAssignPatternsModal",
      @remove="showRemoveRemediationInstructionModal",
      @edit="showEditRemediationInstructionModal"
    )
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import { remediationInstructionToForm, formToRemediationInstruction } from '@/helpers/forms/remediation-instruction';

import entitiesRemediationInstructionsMixin from '@/mixins/entities/remediation/instructions';
import localQueryMixin from '@/mixins/query-local/query';

import RemediationInstructionsList from './remediation-instructions-list.vue';

export default {
  components: { RemediationInstructionsList },
  mixins: [
    entitiesRemediationInstructionsMixin,
    localQueryMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchRemediationInstructionsList({ params: this.getQuery() });
    },

    showEditRemediationInstructionModal(remediationInstruction) {
      this.$modals.show({
        name: MODALS.createRemediationInstruction,
        config: {
          remediationInstruction,
          title: this.$t('modals.createRemediationInstruction.edit.title'),
          action: async (instruction) => {
            await this.updateRemediationInstructionWithConfirm(remediationInstruction, instruction);

            this.$popups.success({
              text: this.$t('modals.createRemediationInstruction.edit.popups.success', {
                instructionName: instruction.name,
              }),
            });

            await this.fetchList();
          },
        },
      });
    },

    showConfirmModalOnRunningRemediationInstruction(action) {
      return new Promise((resolve) => {
        this.$modals.show({
          name: MODALS.confirmation,
          dialogProps: { persistent: true },
          config: {
            text: this.$t('remediationInstructions.errors.runningInstruction'),
            action: async () => {
              await action();
              resolve();
            },
            cancel: resolve,
          },
        });
      });
    },

    async updateRemediationInstructionWithConfirm(remediationInstruction, data) {
      if (remediationInstruction.running) {
        await this.showConfirmModalOnRunningRemediationInstruction(() => {
          this.updateRemediationInstruction({ id: remediationInstruction._id, data });
        });
      } else {
        await this.updateRemediationInstruction({ id: remediationInstruction._id, data });
      }
    },

    showAssignPatternsModal(remediationInstruction) {
      this.$modals.show({
        name: MODALS.patterns,
        config: {
          alarm: true,
          entity: true,
          patterns: {
            alarm_patterns: remediationInstruction.alarm_patterns,
            entity_patterns: remediationInstruction.entity_patterns,
          },
          action: async ({ alarm_patterns: alarmPatterns, entity_patterns: entityPattens }) => {
            if (
              isEqual(remediationInstruction.alarm_patterns, alarmPatterns) &&
              isEqual(remediationInstruction.entity_patterns, entityPattens)
            ) {
              return;
            }

            const form = {
              ...remediationInstructionToForm(remediationInstruction),
              author: remediationInstruction.author,
              alarm_patterns: alarmPatterns,
              entity_patterns: entityPattens,
            };

            await this.updateRemediationInstructionWithConfirm(
              remediationInstruction,
              formToRemediationInstruction(form),
            );
            await this.fetchList();
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
