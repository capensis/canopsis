<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        text-editor-with-template-field(
          v-model="form",
          :templates="config.templates",
          :variables="config.variables",
          :rules="config.rules",
          :label="config.label"
        )
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

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
    const { text = '', template = '' } = this.modal.config;

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
