<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-layout align-center>
          <v-icon>lock</v-icon>
          <c-name-field
            v-field="form.title"
            :label="$t('common.title')"
            :error-messages="errors.collect('title')"
            name="title"
            required
          />
        </v-layout>
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
    const { icon = {} } = this.modal.config;

    return {
      form: {
        title: icon.title ?? '',
        file: icon.file ?? null,
      },
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.editIcon.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
