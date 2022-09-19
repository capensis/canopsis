<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        dynamic-info-form(v-model="form", :is-disabled-id-field="isDisabledIdField")
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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/forms/dynamic-info';

import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDynamicInfo,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { DynamicInfoForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
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
        if (this.config.action) {
          await this.config.action(formToDynamicInfo(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
