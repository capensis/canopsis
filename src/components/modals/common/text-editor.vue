<template lang="pug">
  v-card
    v-card-title.green.darken-4.white--text
      v-layout(justify-space-between, align-center)
        h2 {{ $t('modals.textEditor.title') }}
        v-btn(@click='hideModal', icon, small)
          v-icon.white--text close
    v-card-text
      quill-editor(v-model="text")
    v-btn(@click="submit") {{ $t('common.submit') }}
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

