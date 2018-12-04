<template lang="pug">
  v-layout
    v-flex
      v-switch(
      label="Mix filters",
      :value="multiple",
      @change="$emit('update:multiple', $event)"
      )
    v-flex
      v-select(
      v-bind="$props"
      return-object,
      clearable,
      @input="updateFilter"
      )
</template>

<script>


export default {
  props: {
    value: {
      type: [Object, Array],
      default: () => null,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      default: 'title',
    },
    itemValue: {
      type: String,
      default: 'filter',
    },
    multiple: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    updateFilter(value) {
      this.$emit('input', value);
    },
  },
};
</script>
