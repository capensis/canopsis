<template>
  <geomap-preview :map="map" />
</template>

<script>
import GeomapPreview from '@/components/other/map/partials/geomap-preview.vue';

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
