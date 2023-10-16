<template>
  <v-rating
    class="layout"
    v-field="value"
    :readonly="readonly"
    :hover="hover"
    background-color="primary"
    medium
    half-increments
    @input="$emit('input', $event)"
  />
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Number, String],
      default: 0,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    hover: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
