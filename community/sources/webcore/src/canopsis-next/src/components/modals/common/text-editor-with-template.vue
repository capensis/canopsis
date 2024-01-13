<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <text-editor-with-template-field
          v-model="form"
          :templates="config.templates"
          :variables="config.variables"
          :rules="config.rules"
          :label="config.label"
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
import { CUSTOM_WIDGET_TEMPLATE, MODALS, VALIDATION_DELAY } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import TextEditorWithTemplateField from '@/components/common/text-editor/text-editor-with-template.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.textEditorWithTemplate,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { TextEditorWithTemplateField, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { text = '', template = CUSTOM_WIDGET_TEMPLATE } = this.modal.config;

    return {
      form: {
        text,
        template,
      },
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.textEditor.title');
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
