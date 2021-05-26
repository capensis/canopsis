<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.remediationInstructionApproval.title') }}
      template(slot="text")
        v-layout(v-if="remediationInstructionApproval", column)
          remediation-instruction-approval-alert(
            v-if="remediationInstructionApproval",
            :approval="remediationInstructionApproval.approval"
          )
          v-tabs.mt-3(slider-color="primary", fixed-tabs)
            v-tab {{ $t('modals.remediationInstructionApproval.tabs.updated') }}
            v-tab-item
              span Updated
            v-tab {{ $t('modals.remediationInstructionApproval.tabs.original') }}
            v-tab-item
              remediation-instruction-form(
                :form="remediationInstructionApproval.original",
                disabled-common,
                disabled
              )
        v-layout(v-else, justify-center)
          v-progress-circular(indeterminate, color="primary")
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
import RemediationInstructionForm from
  '@/components/other/remediation/instructions/form/remediation-instruction-form.vue';


import ModalWrapper from '../modal-wrapper.vue';

const { mapActions } = createNamespacedHelpers('remediationInstruction');

export default {
  name: MODALS.remediationInstructionApproval,
  components: {
    RemediationInstructionApprovalAlert,
    RemediationInstructionForm,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      pending: false,
      remediationInstructionApproval: null,
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

      this.remediationInstructionApproval = await this.fetchRemediationInstructionApprovalWithoutStore({
        id: this.config.remediationInstructionId,
      });

      this.pending = false;
    },

    dismiss() {},

    approve() {},
  },
};
</script>
