<template lang="pug">
  v-card.text-editor-modal
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.textEditor.title') }}
    v-card-text
      text-editor-component(v-model="text")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/inner';

import TextEditorComponent from '@/components/other/text-editor/text-editor.vue';

export default {
  name: MODALS.textEditor,
  components: {
    TextEditorComponent,
  },
  mixins: [modalInnerMixin],
  data() {
    const text = this.modal.config.text || '';

    return {
      text,
    };
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action(this.text);
      }

      this.hideModal();
    },
  },
};
</script>

<style lang="scss">
  .text-editor-modal {
    .quill-editor {
      .ql-editor{
        min-height: 120px !important;
        max-height: 300px;
        overflow: hidden;
        overflow-y: auto;
      }
    }
  }
</style>

