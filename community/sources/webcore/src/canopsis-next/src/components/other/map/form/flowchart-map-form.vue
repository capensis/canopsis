<template>
  <v-layout
    class="geomap-map-form"
    column
  >
    <c-name-field
      v-field="form.name"
      required
    />
    <flowchart-editor v-field="form.parameters" />
  </v-layout>
</template>

<script>
import FlowchartEditor from './fields/flowchart-editor.vue';

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
