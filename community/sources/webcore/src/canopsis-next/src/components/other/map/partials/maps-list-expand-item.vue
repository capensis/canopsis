<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $t(`map.types.${map.type}`) }}
    v-tab-item
      v-layout.pa-3
        component(:is="component", :map="map")
</template>

<script>
import { MAP_TYPES } from '@/constants';

import MapsListExpandMermaidItem from './maps-list-expand-mermaid-item.vue';

export default {
  components: { MapsListExpandMermaidItem },
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
  computed: {
    component() {
      return {
        [MAP_TYPES.mermaid]: 'maps-list-expand-mermaid-item',
      }[this.map.type];
    },
  },
};
</script>
