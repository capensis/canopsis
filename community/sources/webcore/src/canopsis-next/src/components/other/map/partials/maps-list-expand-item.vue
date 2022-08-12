<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $t(`map.types.${map.type}`) }}
    v-tab-item
      v-layout.pa-3
        component(:is="component", :map="map")
</template>

<script>
import { MAP_TYPES } from '@/constants';

import MapsListExpandGeomapItem from './maps-list-expand-geomap-item.vue';

export default {
  components: { MapsListExpandGeomapItem },
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
  computed: {
    component() {
      return {
        [MAP_TYPES.geo]: 'maps-list-expand-geomap-item',
      }[this.map.type];
    },
  },
};
</script>
