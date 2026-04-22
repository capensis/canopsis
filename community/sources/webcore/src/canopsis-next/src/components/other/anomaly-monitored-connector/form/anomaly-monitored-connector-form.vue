<template>
  <v-layout column>
    <c-enabled-field v-field="form.enabled" with-background />
    <c-name-field
      v-field="form.name"
      :label="$t('common.connectorName')"
      :max-length="255"
      name="name"
      required
    />
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
      required: true,
    },
  },
};
</script>
