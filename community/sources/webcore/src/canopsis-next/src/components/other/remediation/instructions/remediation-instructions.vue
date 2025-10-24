<template>
  <v-card-text>
    <remediation-instructions-list
      :remediation-instructions="remediationInstructions"
      :pending="remediationInstructionsPending"
      :total-items="remediationInstructionsMeta.total_count"
      :options.sync="options"
      :updatable="hasUpdateAnyRemediationInstructionAccess"
      :removable="hasDeleteAnyRemediationInstructionAccess"
      :duplicable="hasCreateAnyRemediationInstructionAccess"
      @remove-selected="showRemoveSelectedRemediationInstructionModal"
      @duplicate="showDuplicateRemediationInstructionModal"
      @remove="showRemoveRemediationInstructionModal"
      @approve="showApproveRemediationInstructionModal"
      @edit="showEditRemediationInstructionModal"
    />
  </v-card-text>
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { localQueryMixin } from '@/mixins/query/query';
import { entitiesRemediationInstructionMixin } from '@/mixins/entities/remediation/instruction';
import {
  permissionsTechnicalRemediationInstructionMixin,
} from '@/mixins/permissions/technical/remediation-instruction';

import RemediationInstructionsList from './remediation-instructions-list.vue';

export default {
  components: { RemediationInstructionsList },
  mixins: [
    authMixin,
    localQueryMixin,
    entitiesRemediationInstructionMixin,
    permissionsTechnicalRemediationInstructionMixin,
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
        && !(remediationInstruction.approval.requested_by?._id === this.currentUser._id);

      this.$modals.show({
        name: MODALS.createRemediationInstruction,
        config: {
          remediationInstruction,
          disabled: wasRequestedByAnotherUser,
          title: this.$t('modals.createRemediationInstruction.edit.title'),
          action: async (instruction) => {
            await this.updateRemediationInstruction({ id: remediationInstruction._id, data: instruction });

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

    showRemoveRemediationInstructionModal(remediationInstruction) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeRemediationInstruction({ id: remediationInstruction._id });
            await this.fetchList();
          },
        },
      });
    },

    showDuplicateRemediationInstructionModal(remediationInstruction) {
      this.$modals.show({
        name: MODALS.createRemediationInstruction,
        config: {
          remediationInstruction: omit(remediationInstruction, ['_id']),
          title: this.$t('modals.createRemediationInstruction.duplicate.title'),
          action: async (instruction) => {
            await this.createRemediationInstruction({ data: instruction });

            this.$popups.success({
              text: this.$t('modals.createRemediationInstruction.duplicate.popups.success', {
                instructionName: remediationInstruction.name,
              }),
            });

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
