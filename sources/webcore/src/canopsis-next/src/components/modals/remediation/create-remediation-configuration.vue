<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        remediation-configuration-form(v-model="form")
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

import { formToRemediationConfiguration, remediationConfigurationToForm } from '@/helpers/forms/remediation-configuration';

import authMixin from '@/mixins/auth';
import validationErrorsMixin from '@/mixins/form/validation-errors';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import RemediationConfigurationForm from '@/components/other/remediation/configurations/form/remediation-configuration-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRemediationConfiguration,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    RemediationConfigurationForm,
    ModalWrapper,
  },
  mixins: [
    authMixin,
    validationErrorsMixin(),
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: remediationConfigurationToForm(this.modal.config.remediationConfiguration),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createRemediationConfiguration.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            const configuration = formToRemediationConfiguration(this.form);
            configuration.author = this.currentUser._id;

            await this.config.action(configuration);
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
