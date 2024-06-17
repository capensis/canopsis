<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-textarea
          v-if="textarea"
          v-model="form.text"
          v-validate="config.rules"
          :label="config.label"
          :error-messages="errors.collect('text')"
          name="text"
        />
        <text-editor-field
          v-else
          v-model="form.text"
          v-validate="config.rules"
          :label="config.label"
          :sanitize-options="config.sanitizeOptions"
          :error-messages="errors.collect('text')"
          :variables="variables"
          :dark="$system.dark"
          name="text"
        />
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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

const TextEditorField = () => import(/* webpackChunkName: "TextEditor" */ '@/components/common/text-editor/text-editor.vue');

export default {
  name: MODALS.textEditor,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  inject: ['$system'],
  components: { TextEditorField, ModalWrapper },
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

    textarea() {
      return this.config.textarea;
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
