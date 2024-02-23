<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ config.title }}</span>
      </template>
      <template #text="">
        <pattern-form v-model="form" />
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

import { patternToForm, formToPattern } from '@/helpers/entities/pattern/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import PatternForm from '@/components/forms/pattern-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPattern,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    PatternForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const form = patternToForm(this.modal.config.pattern);

    if (this.modal.config.type) {
      form.type = this.modal.config.type;
    }

    return { form };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToPattern(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
