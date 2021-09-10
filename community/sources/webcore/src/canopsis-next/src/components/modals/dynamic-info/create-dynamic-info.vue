<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        dynamic-info-form(v-model="form", :is-disabled-id-field="isDisabledIdField")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { createSubmittableMixin } from '@/mixins/submittable';
import { createConfirmableModalMixin } from '@/mixins/confirmable-modal';
import { createValidationErrorsMixin } from '@/mixins/form/validation-errors';

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/forms/dynamic-info';

import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDynamicInfo,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DynamicInfoForm, ModalWrapper },
  mixins: [
    createSubmittableMixin(),
    createConfirmableModalMixin(),
    createValidationErrorsMixin(),
  ],
  data() {
    return {
      form: dynamicInfoToForm(this.modal.config.dynamicInfo),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createDynamicInfo.create.title');
    },

    isDisabledIdField() {
      return this.config.isDisabledIdField;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToDynamicInfo(this.form));
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
