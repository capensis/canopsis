<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createDynamicInfo.create.title') }}
      template(slot="text")
        dynamic-info-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/forms/dynamic-info';

import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDynamicInfo,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DynamicInfoForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { dynamicInfo = {} } = this.modal.config;

    return {
      form: dynamicInfoToForm(dynamicInfo),
    };
  },
  computed: {
    patterns() {
      return this.form.patterns;
    },
  },
  watch: {
    patterns() {
      this.$validator.validate('patterns');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const preparedData = formToDynamicInfo(this.form);

        if (this.config.action) {
          await this.config.action(preparedData);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
