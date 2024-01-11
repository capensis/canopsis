<template>
  <v-select
    v-field="value"
    :items="languages"
    :disabled="disabled"
    :label="label"
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
      type: String,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'language',
    },
  },
  computed: {
    languages() {
      return Object.keys(this.$i18n.messages);
    },
  },
};
</script>
