<template>
  <v-menu>
    <template #activator="{ on }">
      <input
        v-field="value"
        type="text"
        autocomplete="off"
        v-on="on"
      >
    </template>
  </v-menu>
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
      default: '',
    },
  },
};
</script>
