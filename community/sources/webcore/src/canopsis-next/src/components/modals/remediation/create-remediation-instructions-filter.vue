<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('common.filters') }}
      template(#text="")
        remediation-instructions-filter-form(
          v-model="form",
          :filters="config.filters"
        )
      template(#actions="")
        v-btn(
          :disabled="submitting",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { remediationInstructionFilterToForm } from '@/helpers/forms/remediation-instruction-filter';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import RemediationInstructionsFilterForm
  from '@/components/other/remediation/instructions-filter/remediation-instructions-filter-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRemediationInstructionsFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { RemediationInstructionsFilterForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: remediationInstructionFilterToForm(this.modal.config.filter),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config?.action?.(this.form);

        this.$modals.hide();
      }
    },
  },
};
</script>
