<template>
  <div>
    <c-page-header />
    <v-card class="ma-4 mt-0">
      <maps-list
        :maps="maps"
        :pending="mapsPending"
        :options.sync="options"
        :total-items="mapsMeta.total_count"
        :updatable="hasUpdateAnyMapAccess"
        :removable="hasDeleteAnyMapAccess"
        :duplicable="hasCreateAnyMapAccess"
        @edit="showEditMapModal"
        @remove="showRemoveMapModal"
        @duplicate="showDuplicateMapModal"
        @remove-selected="showDeleteSelectedMapsModal"
      />
    </v-card>
    <c-fab-btn
      :has-access="hasCreateAnyMapAccess"
      @refresh="fetchList"
      @create="showCreateMapModal"
    >
      <span>{{ $t('modals.createMap.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { omit } from 'lodash';

import { MODALS, MAP_TYPES, CREATE_MAP_MODAL_NAMES_BY_TYPE } from '@/constants';

import { pickIds } from '@/helpers/array';

import { authMixin } from '@/mixins/auth';
import { permissionsTechnicalMapMixin } from '@/mixins/permissions/technical/map';
import { entitiesMapMixin } from '@/mixins/entities/map';
import { localQueryMixin } from '@/mixins/query/query';

import MapsList from '@/components/other/map/maps-list.vue';

export default {
  components: {
    MapsList,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    permissionsTechnicalMapMixin,
    entitiesMapMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    showCreateMapModal() {
      this.$modals.show({
        name: MODALS.createMap,
        config: {
          action: async (newMap) => {
            await this.createMap({ data: newMap });

            return this.fetchList();
          },
        },
      });
    },

    async showEditMapModal({ _id: id }) {
      const map = await this.fetchMapWithoutStore({ id });

      const title = {
        [MAP_TYPES.geo]: this.$t('modals.createGeoMap.edit.title'),
        [MAP_TYPES.flowchart]: this.$t('modals.createFlowchartMap.edit.title'),
        [MAP_TYPES.mermaid]: this.$t('modals.createMermaidMap.edit.title'),
        [MAP_TYPES.treeOfDependencies]: this.$t('modals.createTreeOfDependenciesMap.edit.title'),
      }[map.type];

      this.$modals.show({
        name: CREATE_MAP_MODAL_NAMES_BY_TYPE[map.type],
        config: {
          map,
          title,
          action: async (newMap) => {
            await this.updateMap({ id: map._id, data: newMap });

            return this.fetchList();
          },
        },
      });
    },

    async showDuplicateMapModal({ _id: id }) {
      const map = await this.fetchMapWithoutStore({ id });

      const title = {
        [MAP_TYPES.geo]: this.$t('modals.createGeoMap.duplicate.title'),
        [MAP_TYPES.flowchart]: this.$t('modals.createFlowchartMap.duplicate.title'),
        [MAP_TYPES.mermaid]: this.$t('modals.createMermaidMap.duplicate.title'),
        [MAP_TYPES.treeOfDependencies]: this.$t('modals.createTreeOfDependenciesMap.duplicate.title'),
      }[map.type];

      this.$modals.show({
        name: CREATE_MAP_MODAL_NAMES_BY_TYPE[map.type],
        config: {
          map: omit(map, ['_id']),
          title,
          action: async (newMap) => {
            await this.createMap({ data: newMap });

            return this.fetchList();
          },
        },
      });
    },

    showRemoveMapModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeMap({ id });

            return this.fetchList();
          },
        },
      });
    },

    showDeleteSelectedMapsModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.bulkRemoveMaps({
              data: pickIds(selected),
            });

            return this.fetchList();
          },
        },
      });
    },

    fetchList() {
      const params = this.getQuery();
      params.with_flags = true;

      return this.fetchMapsList({ params });
    },
  },
};
</script>
