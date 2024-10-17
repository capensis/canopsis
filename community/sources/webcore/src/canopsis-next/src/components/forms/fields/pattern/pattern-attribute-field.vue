<template>
  <c-select-field
    v-field="value"
    v-validate="rules"
    :items="items"
    :disabled="disabled"
    :error-messages="errors.collect(name)"
    :label="label || $tc('common.searchBy')"
    :name="name"
    item-disabled="options.disabled"
    ellipsis
    autocomplete
  />
</template>

<script>
export default {
  inject: ['$validator'],
  props: {
    value: {
      type: [String, Object],
      required: false,
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
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
};
</script>
