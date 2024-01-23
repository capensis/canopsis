<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <state-setting-form v-model="form" />
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
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { stateSettingToForm, formToStateSetting } from '@/helpers/entities/state-setting/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { authMixin } from '@/mixins/auth';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import StateSettingForm from '@/components/other/state-setting/form/state-setting-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createStateSetting,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    StateSettingForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    authMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: stateSettingToForm(this.modal.config.stateSetting),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createStateSetting.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToStateSetting(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
