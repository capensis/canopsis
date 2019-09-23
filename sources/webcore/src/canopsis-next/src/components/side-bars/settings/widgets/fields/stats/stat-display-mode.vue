<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ $t('settings.statsNumbers.displayMode') }}
        .font-italic.caption.ml-1 ({{ $t('common.optional') }})
      div
        v-btn.primary(
          small,
          @click="showStatDisplayModeModal"
        ) {{ $t('common.edit') }}
</template>

<script>
import { MODALS } from '@/constants';


export default {
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  methods: {
    showStatDisplayModeModal() {
      this.$modals.show({
        name: MODALS.statsDisplayMode,
        config: {
          displayMode: this.value,
          action: data => this.$emit('input', data),
        },
      });
    },
  },
};
</script>
