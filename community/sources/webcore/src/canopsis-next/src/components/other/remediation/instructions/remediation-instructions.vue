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
      @assign-patterns="showAssignPatternsModal"
      @duplicate="showDuplicateRemediationInstructionModal"
      @remove="showRemoveRemediationInstructionModal"
      @approve="showApproveRemediationInstructionModal"
      @edit="showEditRemediationInstructionModal"
    />
  </v-card-text>
</template>

<script>
import { isEqual, omit } from 'lodash';

import { MODALS } from '@/constants';

import {
  remediationInstructionToForm,
  formToRemediationInstruction,
} from '@/helpers/entities/remediation/instruction/form';
import { isSeveralEqual } from '@/helpers/collection';

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

    showAssignPatternsModal(instruction) {
      this.$modals.show({
        name: MODALS.remediationPatterns,
        config: {
          instruction,

          action: async (data) => {
            const isPbehaviorsEqual = isSeveralEqual(instruction, data, [
              'active_on_pbh',
              'disabled_on_pbh',
            ]);

            const isAlarmPatternEqual = instruction.corporate_alarm_pattern === data.corporate_alarm_pattern
              || isEqual(instruction.alarm_pattern, data.alarm_pattern);

            const isEntityPatternEqual = instruction.corporate_entity_pattern === data.corporate_entity_pattern
              || isEqual(instruction.entity_pattern, data.entity_pattern);

            if (isPbehaviorsEqual && isAlarmPatternEqual && isEntityPatternEqual) {
              return;
            }

            const form = {
              ...remediationInstructionToForm(instruction),
              ...data,
            };

            await this.updateRemediationInstruction({
              id: instruction._id,
              data: formToRemediationInstruction(form),
            });

            await this.fetchList();
          },
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
