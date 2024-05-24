<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <remediation-instruction-approval-alert
          v-if="hasApproval && isChangesByCurrentUser"
          :user-name="alertUserName"
          :comment="alertComment"
          :dismissed="isChangesDismissed"
          class="mb-3"
        />
        <remediation-instruction-form
          v-model="form"
          :disabled="disabled"
          :is-new="isNew"
          :required-approve="requiredInstructionApprove"
          :no-pattern="noPattern"
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
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
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

import { MODALS, VALIDATION_DELAY, PATTERNS_FIELDS } from '@/constants';

import {
  formToRemediationInstruction,
  remediationInstructionErrorsToForm,
  remediationInstructionToForm,
} from '@/helpers/entities/remediation/instruction/form';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { validationErrorsMixinCreator } from '@/mixins/form/validation-errors';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { authMixin } from '@/mixins/auth';

import RemediationInstructionForm from '@/components/other/remediation/instructions/form/remediation-instruction-form.vue';
import RemediationInstructionApprovalAlert from '@/components/other/remediation/instructions/partials/approval-alert.vue';

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
    RemediationInstructionApprovalAlert,
  },
  mixins: [
    authMixin,
    modalInnerMixin,
    validationErrorsMixinCreator(),
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: { ...remediationInstructionToForm(this.modal.config.remediationInstruction),
        patterns: {
          ...filterPatternsToForm(
            this.modal.config.remediationInstruction,
            [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
          ),
          active_on_pbh: this.modal.config.remediationInstruction?.active_on_pbh ?? [],
          disabled_on_pbh: this.modal.config.remediationInstruction?.disabled_on_pbh ?? [],
        } },
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

    approval() {
      return this.modal.config.remediationInstruction?.approval;
    },

    hasApproval() {
      return !!this.approval;
    },

    isChangesDismissed() {
      return !!this.approval?.dismissed_by;
    },

    isChangesByCurrentUser() {
      return this.approval?.requested_by?._id === this.currentUser._id;
    },

    alertUserName() {
      const { dismissed_by: dismissedBy, requested_by: requestedBy } = this.approval ?? {};

      return dismissedBy?.display_name ?? requestedBy?.display_name;
    },

    alertComment() {
      return this.approval?.dismiss_comment ?? this.approval?.comment;
    },

    noPattern() {
      return !!this.config.noPattern;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action({
              ...formToRemediationInstruction(this.form),
              ...formFilterToPatterns(this.form.patterns, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
              active_on_pbh: this.form.patterns.active_on_pbh,
              disabled_on_pbh: this.form.patterns.disabled_on_pbh,
            });
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
