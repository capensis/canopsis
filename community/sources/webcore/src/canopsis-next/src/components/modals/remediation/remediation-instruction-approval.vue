<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.remediationInstructionApproval.title') }}
      template(slot="text")
        remediation-instruction-approval-alert(v-if="approval", :approval="approval.approval")
      template(slot="actions")
        v-btn(depressed, flat, @click="dismiss") {{ $t('common.dismiss') }}
        v-btn.primary(@click="approve") {{ $t('common.approve') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/plugins/modals/mixins/inner';

import RemediationInstructionApprovalAlert from
  '@/components/other/remediation/instructions/partials/approval-alert.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions } = createNamespacedHelpers('remediationInstruction');

export default {
  name: MODALS.remediationInstructionApproval,
  components: {
    RemediationInstructionApprovalAlert,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      pending: false,
      approval: null,
    };
  },
  mounted() {
    this.fetchItem();
  },
  methods: {
    ...mapActions({
      fetchRemediationInstructionApprovalWithoutStore: 'fetchItemApprovalWithoutStore',
    }),

    async fetchItem() {
      this.pending = true;

      this.approval = await this.fetchRemediationInstructionApprovalWithoutStore({
        id: this.config.remediationInstructionId,
      });

      this.pending = false;
    },

    dismiss() {},

    approve() {},
  },
};
</script>
