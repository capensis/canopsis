<template lang="pug">
  v-card.text-editor-modal(data-test="textEditorModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text(data-test="jodit")
      text-editor-component(v-model="text")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(data-test="textEditorCancelButton", @click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(data-test="textEditorSubmitButton", @click="submit") {{ $t('common.submit') }}
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

