<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        remediation-job-form(v-model="form")
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

import { formToRemediationJob, remediationJobToForm } from '@/helpers/forms/remediation-job';

import { validationErrorsMixin } from '@/mixins/form/validation-errors';
import { submittableMixin } from '@/mixins/submittable';
import { confirmableModalMixin } from '@/mixins/confirmable-modal';

import RemediationJobForm from '@/components/other/remediation/jobs/form/remediation-job-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRemediationJob,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    ModalWrapper,
    RemediationJobForm,
  },
  mixins: [
    validationErrorsMixin(),
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: remediationJobToForm(this.modal.config.remediationJob),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createRemediationJob.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToRemediationJob(this.form));
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