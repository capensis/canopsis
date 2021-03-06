<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        remediation-instruction-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { formToRemediationInstruction, remediationInstructionToForm } from '@/helpers/forms/remediation-instruction';

import authMixin from '@/mixins/auth';
import validationErrorsMixin from '@/mixins/form/validation-errors';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import RemediationInstructionForm from '@/components/other/remediation/instructions/form/remediation-instruction-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRemediationInstruction,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    ModalWrapper,
    RemediationInstructionForm,
  },
  mixins: [
    authMixin,
    validationErrorsMixin(),
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: remediationInstructionToForm(this.modal.config.remediationInstruction),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createRemediationInstruction.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            const instruction = formToRemediationInstruction(this.form);
            instruction.author = this.currentUser._id;

            await this.config.action(instruction);
          }

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>
