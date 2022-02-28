<template lang="pug">
  v-select(
    v-field="value",
    :items="items",
    :disabled="disabled",
    :label="label || $tc('common.searchBy')",
    :name="name",
    hide-details,
    return-object
  )
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
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
