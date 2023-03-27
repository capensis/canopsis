<template lang="pug">
  v-chip.declare-ticket-event-tickets-chip-field(
    :selected="value",
    :color="value ? 'primary' : undefined",
    :text-color="value ? 'white' : undefined",
    :disabled="disabled",
    small,
    @click="updateValue"
  )
    slot
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Boolean,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    updateValue() {
      if (!this.disabled) {
        this.$emit('input', !this.value);
      }
    },
  },
};
</script>

<style lang="scss">
.declare-ticket-event-tickets-chip-field {
  margin: 0;

  &:not(.v-chip--disabled) .v-chip__content {
    cursor: pointer;
  }
}
</style>
