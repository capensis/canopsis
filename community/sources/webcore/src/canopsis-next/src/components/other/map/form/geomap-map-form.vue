<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <geomap-editor v-field="form.parameters" />
  </v-layout>
</template>

<script>
import GeomapEditor from './fields/geomap-editor.vue';

export default {
  inject: ['$validator'],
  components: { GeomapEditor },
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
