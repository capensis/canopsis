<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $t(`map.types.${map.type}`) }}</v-tab>
    <v-tab-item>
      <v-layout class="pa-3">
        <v-layout
          v-if="!mapDetails"
          justify-center
        >
          <v-progress-circular
            class="pa-4"
            color="white"
            indeterminate
          />
        </v-layout>
        <component
          v-else
          :is="component"
          :map="mapDetails"
        />
      </v-layout>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { MAP_TYPES } from '@/constants';

import { entitiesMapMixin } from '@/mixins/entities/map';

import MapsListExpandMermaidItem from './maps-list-expand-mermaid-item.vue';
import MapsListExpandGeomapItem from './maps-list-expand-geomap-item.vue';
import MapsListExpandFlowchartItem from './maps-list-expand-flowchart-item.vue';
import MapsListExpandTreeOfDependenciesItem from './maps-list-expand-tree-of-dependencies-item.vue';

export default {
  components: {
    MapsListExpandMermaidItem,
    MapsListExpandGeomapItem,
    MapsListExpandFlowchartItem,
    MapsListExpandTreeOfDependenciesItem,
  },
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
        [MAP_TYPES.geo]: 'maps-list-expand-geomap-item',
        [MAP_TYPES.mermaid]: 'maps-list-expand-mermaid-item',
        [MAP_TYPES.flowchart]: 'maps-list-expand-flowchart-item',
        [MAP_TYPES.treeOfDependencies]: 'maps-list-expand-tree-of-dependencies-item',
      }[this.map.type];
    },
  },
  watch: {
    map: 'fetchMapDetails',
  },
  mounted() {
    this.fetchMapDetails();
  },
  methods: {
    async fetchMapDetails() {
      this.pending = true;

      if (this.map.type === MAP_TYPES.treeOfDependencies) {
        this.mapDetails = await this.fetchMapStateWithoutStore({ id: this.map._id });
      } else {
        this.mapDetails = await this.fetchMapWithoutStore({ id: this.map._id });
      }

      this.pending = false;
    },
  },
};
</script>
