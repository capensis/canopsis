<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $t(`map.types.${map.type}`) }}
    v-tab-item
      v-layout.pa-3
        v-layout(v-if="!mapDetails", justify-center)
          v-progress-circular.pa-4(color="white", indeterminate)
        component(v-else, :is="component", :map="mapDetails")
</template>

<script>
import { MAP_TYPES } from '@/constants';

import { entitiesMapMixin } from '@/mixins/entities/map';

import MapsListExpandMermaidItem from './maps-list-expand-mermaid-item.vue';

export default {
  components: { MapsListExpandMermaidItem },
  mixins: [entitiesMapMixin],
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: true,
      mapDetails: undefined,
    };
  },
  computed: {
    component() {
      return {
        [MAP_TYPES.mermaid]: 'maps-list-expand-mermaid-item',
      }[this.map.type];
    },
  },
  mounted() {
    this.fetchMapDetails();
  },
  methods: {
    async fetchMapDetails() {
      this.pending = true;

      this.mapDetails = await this.fetchMapWithoutStore({ id: this.map._id });

      this.pending = false;
    },
  },
};
</script>
