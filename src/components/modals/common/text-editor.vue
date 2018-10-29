<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.textEditor.title') }}
    v-card-text
      quill-editor(v-model="text")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';
import { quillEditor as QuillEditor } from 'vue-quill-editor';
import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.textEditor,
  components: {
    QuillEditor,
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

