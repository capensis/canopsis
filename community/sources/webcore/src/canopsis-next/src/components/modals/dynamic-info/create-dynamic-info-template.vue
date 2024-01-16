<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <dynamic-info-template-form v-model="form" />
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
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS } from '@/constants';

import { templateToForm, formToTemplate } from '@/helpers/entities/dynamic-info/template/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import DynamicInfoTemplateForm from '@/components/other/dynamic-info/form/dynamic-info-template-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDynamicInfoTemplate,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DynamicInfoTemplateForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { template = {} } = this.modal.config;

    return {
      form: templateToForm(template),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createDynamicInfoTemplate.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToTemplate(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
