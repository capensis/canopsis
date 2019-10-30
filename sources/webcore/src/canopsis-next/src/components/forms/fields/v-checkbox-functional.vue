<template lang="pug" functional>
  .v-input.v-input--selection-controls.v-input--checkbox(:class="{ 'v-input--is-disabled': props.disabled }")
    .v-input__control
      .v-input__slot
        .v-input--selection-controls__input(@click="listeners.change(!props.inputValue)")
          .v-input--selection-controls__ripple.primary--text(v-ripple="{ center: true }")
          i.v-icon.material-icons(
            :class="{ 'primary--text': props.inputValue }"
          ) {{ props.inputValue ? 'check_box' : 'check_box_outline_blank' }}
        label(
          v-show="props.label !== ''",
          @click="listeners.change(!props.inputValue)"
        ) {{ props.label }}
</template>

<script>
export default {
  model: {
    prop: 'inputValue',
    event: 'change',
  },
  props: {
    inputValue: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style scoped>
  label {
    cursor: pointer;
  }
</style>
