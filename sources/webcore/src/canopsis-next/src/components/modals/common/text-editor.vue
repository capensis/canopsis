<template lang="pug">
  modal-wrapper(data-test="textEditorModal")
    template(slot="title")
      span {{ title }}
    template(slot="text")
      text-editor-component(v-model="text", data-test="jodit")
    template(slot="actions")
      v-btn(data-test="textEditorCancelButton", @click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(data-test="textEditorSubmitButton", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/inner';

import TextEditorComponent from '@/components/other/text-editor/text-editor.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.textEditor,
  components: {
    TextEditorComponent,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    const text = this.modal.config.text || '';

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
      if (this.config.action) {
        await this.config.action(this.text);
      }

      this.$modals.hide();
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

