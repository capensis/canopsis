<template>
  <geomap-preview :map="map" />
</template>

<script>
import GeomapPreview from './geomap-preview.vue';

export default {
  components: { GeomapPreview },

  props: {
    map: {
      type: Object,
      required: true,
    },
  },
};
</script>
