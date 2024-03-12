<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <icon-form v-model="form" />
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

import { iconToForm, formToRequest } from '@/helpers/entities/icon/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { authMixin } from '@/mixins/auth';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import IconForm from '@/components/other/icons/form/icon-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createIcon,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    IconForm,
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
      form: iconToForm(this.modal.config.icon),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createIcon.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToRequest(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
