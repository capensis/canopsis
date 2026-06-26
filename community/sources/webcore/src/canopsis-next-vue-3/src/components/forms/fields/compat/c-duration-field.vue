<template>
  <c-number-field
    :value="duration.value"
    v-bind="$attrs"
    @update:value="updateValue"
  />
</template>

<script>
export default {
  props: {
    duration: {
      type: Object,
      default: () => ({ value: 0, unit: 's' }),
    },
  },

  emits: ['update:duration'],

  methods: {
    updateValue(value) {
      this.$emit('update:duration', {
        ...this.duration,
        value,
      });
    },
  },
};
</script>
