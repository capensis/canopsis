<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        remediation-instruction-execute(:execution-instruction-id="config.executionInstructionId")
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import authMixin from '@/mixins/auth';
import validationErrorsMixin from '@/mixins/form/validation-errors';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import RemediationInstructionExecute from '@/components/other/remediation/instruction-execute/remediation-instruction-execute.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.executeRemediationInstruction,
  components: {
    ModalWrapper,
    RemediationInstructionExecute,
  },
  mixins: [
    authMixin,
    modalInnerMixin,
    validationErrorsMixin(),
    submittableMixin(),
    confirmableModalMixin(),
  ],
  computed: {
    title() {
      return this.config.executionInstruction.name;
    },
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action();
      }
    },
  },
};
</script>
