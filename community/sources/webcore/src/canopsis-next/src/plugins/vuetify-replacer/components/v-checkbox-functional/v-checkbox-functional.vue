<template lang="pug">
  div.v-input.v-input--selection-controls.v-input--checkbox.theme--light(
    :class="{ 'v-input--is-disabled': disabled, 'v-input--hide-details': hideDetails }"
  )
    div.v-input__control
      div.v-input__slot
        div.v-input--selection-controls__input(@click="change")
          input(
            class="hidden",
            :aria-checked="String(inputValue)",
            :checked="inputValue",
            :disabled="disabled",
            role="checkbox",
            type="checkbox"
          )
          div.v-input--selection-controls__ripple.primary--text(v-ripple="{ center: true }")
          i.v-icon.material-icons.theme--light(
            :class="{ 'primary--text': inputValue }"
          ) {{ inputValue ? 'check_box' : 'check_box_outline_blank' }}
        label(
          v-show="label !== ''",
          @click="change"
        ) {{ label }}
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
    hideDetails: {
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
  methods: {
    change() {
      this.$emit('change', !this.inputValue);
    },
  },
};
</script>

<style scoped>
  label {
    cursor: pointer;
  }
</style>
