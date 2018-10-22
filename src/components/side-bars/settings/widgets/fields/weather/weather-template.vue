<template lang="pug">
  v-container(fluid).pa-3
    v-layout(align-center, justify-space-between)
      div.subheading {{ title }}
      v-btn.green.darken-4.white--text(
      small,
      @click="showtextEditorModal"
      ) {{ $t('common.show') }}/{{ $t('common.edit') }}
</template>

<script>
import { MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';
import TextEditor from '@/components/other/text-editor/text-editor.vue';

export default {
  components: {
    TextEditor,
  },
  mixins: [modalMixin],
  props: {
    value: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
  },
  methods: {
    showtextEditorModal() {
      this.showModal({
        name: MODALS.textEditor,
        config: {
          text: this.value,
          action: value => this.$emit('input', value),
        },
      });
    },
  },
};
</script>

