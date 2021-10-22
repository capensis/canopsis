<template lang="pug">
  v-card-text
    remediation-instructions-list(
      :remediation-instructions="remediationInstructions",
      :pending="remediationInstructionsPending",
      :total-items="remediationInstructionsMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedRemediationInstructionModal",
      @assign-patterns="showAssignPatternsModal",
      @remove="showRemoveRemediationInstructionModal",
      @approve="showApproveRemediationInstructionModal",
      @edit="showEditRemediationInstructionModal"
    )
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import { remediationInstructionToForm, formToRemediationInstruction } from '@/helpers/forms/remediation-instruction';

import { authMixin } from '@/mixins/auth';
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationInstructionMixin } from '@/mixins/entities/remediation/instruction';

import RemediationInstructionsList from './remediation-instructions-list.vue';

export default {
  components: { RemediationInstructionsList },
  mixins: [
    authMixin,
    localQueryMixin,
    entitiesRemediationInstructionMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      const params = this.getQuery();
      params.with_flags = true;
      params.with_month_executions = true;

      return this.fetchRemediationInstructionsList({ params });
    },

    showEditRemediationInstructionModal(remediationInstruction) {
      const wasRequestedByAnotherUser = !!remediationInstruction.approval
        && !(remediationInstruction.approval.requested_by === this.currentUser._id);

      this.$modals.show({
        name: MODALS.createRemediationInstruction,
        config: {
          remediationInstruction,
          disabled: wasRequestedByAnotherUser,
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

    showApproveRemediationInstructionModal(remediationInstruction) {
      this.$modals.show({
        name: MODALS.remediationInstructionApproval,
        config: {
          remediationInstructionId: remediationInstruction._id,
          afterSubmit: this.fetchList,
        },
      });
    },

    showConfirmModalOnRunningRemediationInstruction(action) {
      return new Promise((resolve, reject) => {
        this.$modals.show({
          name: MODALS.confirmation,
          dialogProps: { persistent: true },
          config: {
            text: this.$t('remediationInstructions.errors.runningInstruction'),
            action: async () => {
              try {
                await action();

                resolve();
              } catch (err) {
                reject(err);
              }
            },
            cancel: resolve,
          },
        });
      });
    },

    async updateRemediationInstructionWithConfirm(remediationInstruction, data) {
      if (remediationInstruction.running) {
        await this.showConfirmModalOnRunningRemediationInstruction(
          () => this.updateRemediationInstruction({ id: remediationInstruction._id, data }),
        );
      } else {
        await this.updateRemediationInstruction({ id: remediationInstruction._id, data });
      }
    },

    showAssignPatternsModal(remediationInstruction) {
      const patterns = {
        alarm_patterns: remediationInstruction.alarm_patterns || [],
        entity_patterns: remediationInstruction.entity_patterns || [],
        active_on_pbh: remediationInstruction.active_on_pbh || [],
        disabled_on_pbh: remediationInstruction.disabled_on_pbh || [],
      };

      this.$modals.show({
        name: MODALS.remediationPatterns,
        config: {
          patterns,

          action: async (newPatterns) => {
            if (isEqual(patterns, newPatterns)) {
              return;
            }

            const form = {
              ...remediationInstructionToForm(remediationInstruction),
              ...newPatterns,
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
