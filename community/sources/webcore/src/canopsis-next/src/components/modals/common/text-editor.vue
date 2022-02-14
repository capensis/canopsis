<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        text-editor-component(
          v-model="text",
          v-validate="config.rules",
          :label="config.label",
          :error-messages="errors.collect('text')",
          name="text"
        )
      template(slot="actions")
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
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import TextEditorComponent from '@/components/common/text-editor/text-editor.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.textEditor,
  $_veeValidate: {
    validator: 'new',
  },
  components: { TextEditorComponent, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator({ field: 'text' }),
  ],
  data() {
    const text = this.modal.config.text ?? '';

    return {
      text,
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.textEditor.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.text);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

<style lang="scss">
  .text-editor-modal {
    .quill-editor {
      .ql-editor {
        min-height: 120px !important;
        max-height: 300px;
        overflow: hidden;
        overflow-y: auto;
      }
    }
  }
</style>
