<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <remediation-instruction-form
          v-model="form"
          :disabled="disabled"
          :is-new="isNew"
          :required-approve="requiredInstructionApprove"
        />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          :disabled="isDisabled"
          :loading="submitting"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import {
  formToRemediationInstruction,
  remediationInstructionErrorsToForm,
  remediationInstructionToForm,
} from '@/helpers/entities/remediation/instruction/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { validationErrorsMixinCreator } from '@/mixins/form/validation-errors';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import RemediationInstructionForm from '@/components/other/remediation/instructions/form/remediation-instruction-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapGetters } = createNamespacedHelpers('info');

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
    ...mapGetters({
      requiredInstructionApprove: 'requiredInstructionApprove',
    }),

    title() {
      return this.config.title || this.$t('modals.createRemediationInstruction.create.title');
    },

    disabled() {
      return this.config.disabled;
    },

    isNew() {
      return !this.modal.config.remediationInstruction?._id;
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
          this.setFormErrors(remediationInstructionErrorsToForm(err, this.form));
        }
      }
    },
  },
};
</script>
