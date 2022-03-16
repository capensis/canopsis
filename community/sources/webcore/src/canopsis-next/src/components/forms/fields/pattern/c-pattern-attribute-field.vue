<template lang="pug">
  v-select.c-pattern-attribute-field(
    v-field="value",
    :items="items",
    :disabled="disabled",
    :label="label || $tc('common.searchBy')",
    :name="name"
  )
    template(#selection="{ item }")
      span.ellipsis {{ item.text }}
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

<style lang="scss">
$selectIconWidth: 24px;

.c-pattern-attribute-field {
  .v-select__selections {
    width: calc(100% - #{$selectIconWidth});
    flex-wrap: nowrap;
  }
}
</style>
