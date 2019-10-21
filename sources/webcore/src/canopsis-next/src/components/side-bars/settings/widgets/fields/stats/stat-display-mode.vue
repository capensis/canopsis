<template lang="pug">
  v-container.pa-3(data-test="statDisplayMode", fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ $t('settings.statsNumbers.displayMode') }}
        .font-italic.caption.ml-1 ({{ $t('common.optional') }})
      div
        v-btn.primary(
          data-test="editButton",
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
