<template>
  <advanced-search-chip
    :model-value="preparedValue"
    :active="active"
    :always-active="alwaysActive"
    :closable="closable"
    :disabled="disabled"
    :input-placeholder="$t('advancedSearch.period')"
    allow-text
    @update:model-value="update"
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
      type: Object,
      default: () => ({ from: '', to: '' }),
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

  computed: {
    preparedValue() {
      return [this.modelValue.from, this.modelValue.to].filter(Boolean).join(' - ');
    },
  },

  methods: {
    update(value) {
      const prepared = { from: value, to: value };

      this.$emit('update:modelValue', prepared);
    },
  },
};
</script>
