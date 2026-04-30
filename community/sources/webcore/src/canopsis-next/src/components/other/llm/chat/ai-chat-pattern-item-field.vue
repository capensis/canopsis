<template>
  <div>
    <c-select-field
      v-field="value"
      v-validate="'required'"
      :label="label"
      :items="items"
      :error-messages="errors.collect('patternItem')"
      name="patternItem"
      item-text="text"
      item-value="value"
      hide-details
      return-object
    />
  </div>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
  },
};
</script>
