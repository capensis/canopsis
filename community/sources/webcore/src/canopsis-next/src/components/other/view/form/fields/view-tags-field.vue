<template>
  <c-combobox
    :model-value="value"
    :label="$t('view.groupTags')"
    append-icon=""
    clearable
    multiple
    chips
    closable-chips
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
export default {
  props: {
    value: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:value'],
};
</script>
