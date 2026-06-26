<template>
  <advanced-search-chip
    :model-value="modelValue"
    :active="active"
    :always-active="alwaysActive"
    :closable="closable"
    :disabled="disabled"
    :input-placeholder="$t('advancedSearch.duration')"
    allow-text
    @update:model-value="$emit('update:modelValue', $event)"
    @click="$emit('click')"
    @close="$emit('close')"
  />
</template>

<script>
import AdvancedSearchChip from './advanced-search-chip.vue';

export default {
  components: { AdvancedSearchChip },

  props: {
    modelValue: {
      type: [String, Number, Object],
      default: '',
    },
    active: {
      type: Boolean,
      default: false,
    },
    alwaysActive: {
      type: Boolean,
      default: false,
    },
    closable: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'click', 'close'],
};
</script>
