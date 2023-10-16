<template>
  <v-layout column>
    <c-description-field
      v-field="form.message"
      :label="$t('maintenance.message')"
      name="message"
      required
    />
    <c-color-picker-field v-field="form.color" />
  </v-layout>
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
