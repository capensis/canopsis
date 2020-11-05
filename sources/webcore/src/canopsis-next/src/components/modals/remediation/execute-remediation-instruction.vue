<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        remediation-instruction-execute(:execution-instruction-id="config.executionInstructionId")
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import authMixin from '@/mixins/auth';

import RemediationInstructionExecute from '@/components/other/remediation/instruction-execute/remediation-instruction-execute.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.executeRemediationInstruction,
  components: {
    ModalWrapper,
    RemediationInstructionExecute,
  },
  mixins: [authMixin, modalInnerMixin],
};
</script>
