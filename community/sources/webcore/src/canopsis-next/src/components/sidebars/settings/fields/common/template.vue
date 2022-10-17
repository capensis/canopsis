<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ title }}
      v-btn.primary(
        small,
        @click="showTextEditorModal"
      ) {{ $t('common.show') }}/{{ $t('common.edit') }}
</template>

<script>
import { MODALS } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    variables: {
      type: Array,
      required: false,
    },
  },
  methods: {
    showTextEditorModal() {
      this.$modals.show({
        name: MODALS.textEditor,
        config: {
          title: this.title,
          text: this.value,
          variables: this.variables,
          action: value => this.$emit('input', value),
        },
      });
    },
  },
};
</script>
