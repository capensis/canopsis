<template lang="pug">
  div.v-input.v-textarea.v-text-field.v-text-field--box.v-text-field--enclosed.v-input--is-label-active(
    class="v-input--is-dirty theme--light"
  )
    div.v-input__control(@click="$emit('click', $event)")
      div.v-input__slot
        div.v-text-field__slot
          label.v-label.v-label--active.theme--light {{ label }}
          div(v-html="value")
      div.v-text-field__details
        div.v-messages.theme--light.error--text
          div.v-messages__wrapper
            div.v-messages__message(v-for="errorMessage in errorMessages") {{ errorMessage }}
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
  },
};
</script>

<style scoped>
.v-label {
  left: 0;
  right: auto;
  position: absolute;
}

.v-text-field__slot {
  padding-top: 28px!important;
}
</style>
