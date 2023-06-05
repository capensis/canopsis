<template lang="pug">
  div.v-input.v-input--selection-controls.v-input--checkbox(
    :class="[{ 'v-input--is-disabled': disabled, 'v-input--hide-details': hideDetails }, themeClasses]"
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
          i.v-icon.material-icons(
            :class="[{ 'primary--text': inputValue }, themeClasses]"
          ) {{ inputValue ? 'check_box' : 'check_box_outline_blank' }}
        label.v-label(
          v-show="label !== ''",
          :class="themeClasses",
          @click="change"
        ) {{ label }}
</template>

<script>
import Themeable from 'vuetify/es5/mixins/themeable';

export default {
  mixins: [Themeable],
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
