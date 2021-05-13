<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.remediationInstructionApproval.title') }}
      template(slot="text")
        v-fade-transition
          c-progress-overlay(v-if="pending", pending)
          span(v-else-if="remediationInstructionApproval") something
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions } = createNamespacedHelpers('remediationInstruction');

export default {
  name: MODALS.remediationInstructionApproval,
  components: { ModalWrapper },
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
  },
};
</script>
