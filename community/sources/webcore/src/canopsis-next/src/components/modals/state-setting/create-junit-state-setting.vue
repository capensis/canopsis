<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>
      <template #text>
        <junit-state-setting-form v-model="form" />
      </template>
      <template #actions>
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
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { junitStateSettingToForm, formToJunitStateSetting } from '@/helpers/entities/state-setting/junit/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { authMixin } from '@/mixins/auth';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import JunitStateSettingForm from '@/components/other/state-setting/junit/form/junit-state-setting-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createJunitStateSetting,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    JunitStateSettingForm,
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
      form: junitStateSettingToForm(this.modal.config.stateSetting),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createJunitStateSetting.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToJunitStateSetting(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
