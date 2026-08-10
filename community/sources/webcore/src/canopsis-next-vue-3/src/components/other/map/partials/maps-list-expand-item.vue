<template>
  <div>
    <v-tabs
      :model-value="tab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab :value="tab">
        {{ $t(`map.types.${map.type}`) }}
      </v-tab>
    </v-tabs>

    <v-window :model-value="tab" class="bg-secondary-lighten-2">
      <v-window-item :value="tab">
        <c-layout class="pa-3">
          <c-layout v-if="!mapDetails" justify="center">
            <v-progress-circular class="pa-4" color="white" indeterminate />
          </c-layout>
          <component v-else-if="component" :is="component" :map="mapDetails" />
          <span v-else />
        </c-layout>
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { MAP_TYPES } from '@/constants';
import { useMap } from '@/hooks/store/modules/map';

import MapsListExpandFlowchartItem from './maps-list-expand-flowchart-item.vue';
import MapsListExpandGeomapItem from './maps-list-expand-geomap-item.vue';
import MapsListExpandMermaidItem from './maps-list-expand-mermaid-item.vue';
import MapsListExpandTreeOfDependenciesItem from './maps-list-expand-tree-of-dependencies-item.vue';

export default {
  components: {
    MapsListExpandFlowchartItem,
    MapsListExpandGeomapItem,
    MapsListExpandMermaidItem,
    MapsListExpandTreeOfDependenciesItem,
  },

  props: {
    map: {
      type: Object,
      required: true,
    },
  },

  setup() {
    const { fetchMapWithoutStore, fetchMapStateWithoutStore } = useMap();

    return {
      fetchMapWithoutStore,
      fetchMapStateWithoutStore,
    };
  },

  data() {
    return {
      tab: 'preview',
      pending: true,
      mapDetails: undefined,
    };
  },

  computed: {
    component() {
      return {
        [MAP_TYPES.geo]: 'maps-list-expand-geomap-item',
        [MAP_TYPES.flowchart]: 'maps-list-expand-flowchart-item',
        [MAP_TYPES.mermaid]: 'maps-list-expand-mermaid-item',
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
