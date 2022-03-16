<template lang="pug">
  c-select-field(
    v-field="value",
    :items="items",
    :disabled="disabled",
    :label="label || $tc('common.searchBy')",
    :name="name",
    ellipsis
  )
</template>

<script>
export default {
  props: {
    value: {
      type: [String, Object],
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'attribute',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
