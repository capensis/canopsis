<template>
  <div>
    <c-page-header />
    <v-card class="ma-4 mt-0">
      <maps-list
        :maps="maps"
        :pending="mapsPending"
        :options="options"
        :total-items="mapsMeta.total_count"
        :updatable="hasUpdateAccess"
        :removable="hasDeleteAccess"
        :duplicable="hasCreateAccess"
        @edit="showEditMapModal"
        @remove="showRemoveMapModal"
        @duplicate="showDuplicateMapModal"
        @refresh="fetchList"
        @update:options="updateOptions"
      />
    </v-card>
    <c-fab-btn
      :has-access="hasCreateAccess"
      @refresh="fetchList"
      @create="showCreateMapModal"
    >
      <span>{{ $t('modals.createMap.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { omit } from 'lodash';
import { onMounted } from 'vue';

import { MODALS, MAP_TYPES, CREATE_MAP_MODAL_NAMES_BY_TYPE, USER_PERMISSIONS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useMaps } from '@/hooks/store/modules/maps';

import MapsList from '@/components/other/map/maps-list.vue';

export default {
  components: {
    MapsList,
  },
  setup() {
    const {
      items: maps,
      pending: mapsPending,
      meta: mapsMeta,
      fetchList: fetchMapsList,
      fetchItemWithoutStore: fetchMapWithoutStore,
      createMap,
      updateMap,
      removeMap,
    } = useMaps();

    const {
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.map);

    const modals = useModals();
    const { t } = useI18n();

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useLocalQueryWithOptions({
      onUpdate: (fetchQuery) => {
        const params = convertQueryToRequest(fetchQuery);
        params.with_flags = true;
        return fetchMapsList({ params });
      },
    });

    /**
     * Shows the modal for creating a new map.
     * After successful creation, refreshes the maps list.
     */
    const showCreateMapModal = () => {
      modals.show({
        name: MODALS.createMap,
        config: {
          action: async (newMap) => {
            await createMap({ data: newMap });

            return fetchList();
          },
        },
      });
    };

    /**
     * Shows the modal for editing an existing map.
     * After successful update, refreshes the maps list.
     *
     * @param {Object} params - The parameters object.
     * @param {string} params._id - The unique identifier of the map to edit.
     */
    const showEditMapModal = async ({ _id: id }) => {
      const map = await fetchMapWithoutStore({ id });

      const title = {
        [MAP_TYPES.geo]: t('modals.createGeoMap.edit.title'),
        [MAP_TYPES.flowchart]: t('modals.createFlowchartMap.edit.title'),
        [MAP_TYPES.mermaid]: t('modals.createMermaidMap.edit.title'),
        [MAP_TYPES.treeOfDependencies]: t('modals.createTreeOfDependenciesMap.edit.title'),
      }[map.type];

      modals.show({
        name: CREATE_MAP_MODAL_NAMES_BY_TYPE[map.type],
        config: {
          map,
          title,
          action: async (newMap) => {
            await updateMap({ id: map._id, data: newMap });

            return fetchList();
          },
        },
      });
    };

    /**
     * Shows the modal for duplicating an existing map.
     * After successful duplication, refreshes the maps list.
     *
     * @param {Object} params - The parameters object.
     * @param {string} params._id - The unique identifier of the map to duplicate.
     */
    const showDuplicateMapModal = async ({ _id: id }) => {
      const map = await fetchMapWithoutStore({ id });

      const title = {
        [MAP_TYPES.geo]: t('modals.createGeoMap.duplicate.title'),
        [MAP_TYPES.flowchart]: t('modals.createFlowchartMap.duplicate.title'),
        [MAP_TYPES.mermaid]: t('modals.createMermaidMap.duplicate.title'),
        [MAP_TYPES.treeOfDependencies]: t('modals.createTreeOfDependenciesMap.duplicate.title'),
      }[map.type];

      modals.show({
        name: CREATE_MAP_MODAL_NAMES_BY_TYPE[map.type],
        config: {
          map: omit(map, ['_id']),
          title,
          action: async (newMap) => {
            await createMap({ data: newMap });

            return fetchList();
          },
        },
      });
    };

    /**
     * Shows the confirmation modal for deleting a map.
     * After successful deletion, refreshes the maps list.
     *
     * @param {string} id - The unique identifier of the map to delete.
     */
    const showRemoveMapModal = (id) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removeMap({ id });

            return fetchList();
          },
        },
      });
    };

    onMounted(fetchList);

    return {
      maps,
      mapsPending,
      mapsMeta,
      options,
      updateOptions,
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,
      fetchList,
      showCreateMapModal,
      showEditMapModal,
      showDuplicateMapModal,
      showRemoveMapModal,
    };
  },
};
</script>
