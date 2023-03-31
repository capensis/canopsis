<template lang="pug">
  v-layout.geomap-map-form(column)
    c-name-field(v-field="form.name", required)
    flowchart-editor(v-field="form.parameters")
</template>

<script>
import FlowchartEditor from './partials/flowchart-editor.vue';

export default {
  inject: ['$validator'],
  components: { FlowchartEditor },
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

<style lang="scss">
.geomap-map-form {
  &__editor {
    height: 500px;
  }
}
</style>
