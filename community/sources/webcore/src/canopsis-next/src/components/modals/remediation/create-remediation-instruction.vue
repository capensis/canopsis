<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        remediation-instruction-form(v-model="form", :disabled="disabled")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import {
  formToRemediationInstruction,
  remediationInstructionErrorsToForm,
  remediationInstructionToForm,
} from '@/helpers/forms/remediation-instruction';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { validationErrorsMixinCreator } from '@/mixins/form/validation-errors';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import RemediationInstructionForm from '@/components/other/remediation/instructions/form/remediation-instruction-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRemediationInstruction,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ModalWrapper,
    RemediationInstructionForm,
  },
  mixins: [
    modalInnerMixin,
    validationErrorsMixinCreator(),
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
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

    disabled() {
      return this.config.disabled;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToRemediationInstruction(this.form));
          }

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(remediationInstructionErrorsToForm(err));
        }
      }
    },
  },
};
</script>
