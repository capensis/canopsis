<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title>
        {{ title }}
      </template>
      <template #text>
        <c-payload-textarea-field
          v-model="form.text"
          :required="isRequired"
          :label="config.label"
          :name="config.name"
          :variables="config.variables"
        />
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
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.payloadTextareaEditor,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  inject: ['$system'],
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const text = this.modal.config.text ?? '';

    return {
      form: {
        text,
      },
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.textEditor.title');
    },

    variables() {
      return this.config.variables;
    },

    isRequired() {
      return this.config.rules?.required;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form.text);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
