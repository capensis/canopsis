<template lang="pug">
  geomap-preview.geomap-expand-panel(:points="map.properties.points")
</template>

<script>
import GeomapPreview from '@/components/other/map/form/partials/geomap-preview.vue';

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

<style lang="scss">
.geomap-expand-panel {
  min-height: 600px;
}
</style>
