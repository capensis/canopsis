<template lang="pug">
  v-layout.geomap-map-form(column)
    c-name-field.mx-4(v-field="form.name")
    div.geomap-map-form__editor
      geomap-editor(v-field="form.parameters")
</template>

<script>
import MermaidEditor from '@/components/common/mermaid/mermaid-editor.vue';

import GeomapEditor from './partials/geomap-editor.vue';

export default {
  inject: ['$validator'],
  components: { GeomapEditor, MermaidEditor },
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
    height: 800px;
  }
}
</style>
