<template>
  <modal-wrapper close>
    <template #title="">
      {{ $t('modals.remediationInstructionApproval.title') }}
    </template>
    <template #text="">
      <v-fade-transition>
        <v-layout
          v-if="!remediationInstructionApproval"
          justify-center
        >
          <v-progress-circular
            color="primary"
            indeterminate
          />
        </v-layout>
        <v-layout
          v-else
          column
        >
          <remediation-instruction-approval-alert
            :user-name="remediationInstructionApproval.approval.requested_by.name"
            :comment="remediationInstructionApproval.approval.comment"
          >
            <remediation-instruction-approval-form
              :disabled="!remediationInstructionApproval"
              :submitting="submitting"
              @approve="approve"
              @dismiss="dismiss"
            />
          </remediation-instruction-approval-alert>
          <remediation-instruction-approval-tabs
            :original="remediationInstructionApproval.original"
            :updated="remediationInstructionApproval.updated"
          />
        </v-layout>
      </v-fade-transition>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.close') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';

import RemediationInstructionApprovalAlert from '@/components/other/remediation/instructions/partials/approval-alert.vue';
import RemediationInstructionApprovalForm from '@/components/other/remediation/instructions/partials/approval-form.vue';
import RemediationInstructionApprovalTabs from '@/components/other/remediation/instructions/partials/approval-tabs.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions } = createNamespacedHelpers('remediationInstruction');

export default {
  name: MODALS.remediationInstructionApproval,
  components: {
    RemediationInstructionApprovalAlert,
    RemediationInstructionApprovalForm,
    RemediationInstructionApprovalTabs,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
  ],
  data() {
    return {
      remediationInstructionApproval: null,
    };
  },
  mounted() {
    this.fetchItem();
  },
  methods: {
    ...mapActions({
      fetchRemediationInstructionApprovalWithoutStore: 'fetchItemApprovalWithoutStore',
      updateRemediationInstructionApproval: 'updateApproval',
    }),

    async fetchItem() {
      this.remediationInstructionApproval = await this.fetchRemediationInstructionApprovalWithoutStore({
        id: this.config.remediationInstructionId,
      });
    },

    approve() {
      return this.submit({ approve: true });
    },

    dismiss(comment) {
      return this.submit({
        approve: false,
        comment,
      });
    },

    async submit(data) {
      await this.updateRemediationInstructionApproval({
        id: this.config.remediationInstructionId,
        data,
      });

      if (this.config.afterSubmit) {
        await this.config.afterSubmit();
      }

      this.$modals.hide();
    },
  },
};
</script>
