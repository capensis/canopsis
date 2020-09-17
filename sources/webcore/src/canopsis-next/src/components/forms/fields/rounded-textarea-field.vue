<template lang="pug">
  v-textarea.rounded-field(
    v-field="value",
    :label="label",
    :error-messages="errorMessages",
    :name="name",
    :type="type",
    box
  )
</template>

<script>
export default {
  $_veeValidate: {
    value() {
      return this.value;
    },

    name() {
      return this.name;
    },
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      default: 'text',
    },
    label: {
      type: String,
      default: '',
    },
  },
};
</script>
