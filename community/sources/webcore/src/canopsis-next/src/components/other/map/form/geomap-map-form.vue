<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", required)
    geomap-editor(v-field="form.parameters")
</template>

<script>
import GeomapEditor from './partials/geomap-editor.vue';

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
