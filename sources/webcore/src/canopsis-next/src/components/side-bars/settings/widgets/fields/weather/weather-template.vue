<template lang="pug">
  v-container(fluid).pa-3
    v-layout(align-center, justify-space-between)
      div.subheading {{ title }}
      v-btn.primary(
      data-test="showEditButton",
      small,
      @click="showTextEditorModal"
      ) {{ $t('common.show') }}/{{ $t('common.edit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

export default {
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
    showTextEditorModal() {
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

