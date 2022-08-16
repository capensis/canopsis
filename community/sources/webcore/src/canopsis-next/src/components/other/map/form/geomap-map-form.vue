<template lang="pug">
  v-layout.geomap-map-form(column)
    c-name-field.mx-4(v-field="form.name")
    div.geomap-map-form__editor
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

<style lang="scss">
.geomap-map-form {
  &__editor {
    height: 500px;
  }
}
</style>
