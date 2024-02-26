<template>
  <div class="pa-2">
    <v-layout
      class="mx-1"
      wrap
    >
      <v-flex
        v-if="hasAccessToCategory"
        xs3
      >
        <c-entity-category-field
          :category="query.category"
          class="mr-3"
          @input="updateCategory"
        />
      </v-flex>
      <v-flex
        v-if="hasAccessToUserFilter"
        xs4
      >
        <v-layout align-center>
          <filter-selector
            :label="$t('settings.selectAFilter')"
            :filters="userPreference.filters"
            :locked-filters="widget.filters"
            :value="mainFilter"
            :locked-value="lockedFilter"
            :disabled="!hasAccessToListFilters && !hasAccessToUserFilter"
            @input="updateSelectedFilter"
          />
          <filters-list-btn
            v-if="hasAccessToAddFilter || hasAccessToEditFilter"
            :widget-id="widget._id"
            :addable="hasAccessToAddFilter"
            :editable="hasAccessToEditFilter"
            with-entity
            with-service-weather
            private
          />
        </v-layout>
      </v-flex>
    </v-layout>
    <template v-if="mapState">
      <v-fade-transition
        v-if="pending"
        key="progress"
        mode="out-in"
      >
        <v-progress-linear
          class="progress-linear-absolute--top"
          height="2"
          indeterminate
        />
      </v-fade-transition>
    </template>
    <template v-else>
      <v-layout
        v-if="pending"
        class="pa-4"
        justify-center
      >
        <v-progress-circular
          color="primary"
          indeterminate
        />
      </v-layout>
    </template>
    <map-breadcrumbs
      v-if="previousMaps.length"
      :previous-maps="previousMaps"
      :active-map="mapState"
      :pending="pending"
      class="mb-2"
      @click="backToBreadcrumb"
    />
    <component
      v-if="mapState"
      :is="component"
      :map="mapState"
      :columns="widget.parameters.entitiesColumns"
      :popup-template="widget.parameters.entity_info_template"
      :color-indicator="widget.parameters.color_indicator"
      :pbehavior-enabled="widget.parameters.entities_under_pbehavior_enabled"
      popup-actions
      @show:map="showMap"
      @show:alarms="showAlarmListModal"
    />
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { pick } from 'lodash';

import { ENTITY_TYPES, MAP_TYPES, MODALS } from '@/constants';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { permissionsWidgetsMapCategory } from '@/mixins/permissions/widgets/map/category';
import { permissionsWidgetsMapFilters } from '@/mixins/permissions/widgets/map/filters';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';

import MapBreadcrumbs from './partials/map-breadcrumbs.vue';

const { mapActions: mapMapActions } = createNamespacedHelpers('map');
const { mapActions: mapActiveViewActions } = createNamespacedHelpers('activeView');
const { mapActions: mapAlarmActions } = createNamespacedHelpers('alarm');
const { mapActions: mapServiceActions } = createNamespacedHelpers('service');

const MermaidPreview = () => import(/* webpackChunkName: "Maps" */ '@/components/other/map/partials/mermaid-preview.vue');
const GeomapPreview = () => import(/* webpackChunkName: "Maps" */ '@/components/other/map/partials/geomap-preview.vue');
const FlowchartPreview = () => import(/* webpackChunkName: "Maps" */ '@/components/other/map/partials/flowchart-preview.vue');
const TreeOfDependenciesPreview = () => import(/* webpackChunkName: "Maps" */ '@/components/other/map/partials/tree-of-dependencies-preview.vue');

export default {
  components: {
    MapBreadcrumbs,
    FilterSelector,
    FiltersListBtn,
    MermaidPreview,
    GeomapPreview,
    FlowchartPreview,
    TreeOfDependenciesPreview,
  },
  mixins: [
    permissionsWidgetsMapCategory,
    permissionsWidgetsMapFilters,
    widgetPeriodicRefreshMixin,
    widgetFilterSelectMixin,
    widgetFetchQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      mapState: undefined,
      previousMaps: [],
      activeMapId: this.widget.parameters.map,
    };
  },
  computed: {
    component() {
      return {
        [MAP_TYPES.geo]: 'geomap-preview',
        [MAP_TYPES.flowchart]: 'flowchart-preview',
        [MAP_TYPES.mermaid]: 'mermaid-preview',
        [MAP_TYPES.treeOfDependencies]: 'tree-of-dependencies-preview',
      }[this.mapState.type];
    },
  },
  watch: {
    activeMapId(id, oldId) {
      if (id !== oldId) {
        this.fetchList();
      }
    },
  },
  created() {
    this.registerEditingOffHandler(this.clearPreviousMaps);
  },
  beforeDestroy() {
    this.unregisterEditingOffHandler(this.clearPreviousMaps);
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapMapActions({
      fetchMapStateWithoutStore: 'fetchItemStateWithoutStore',
    }),
    ...mapActiveViewActions({
      registerEditingOffHandler: 'registerEditingOffHandler',
      unregisterEditingOffHandler: 'unregisterEditingOffHandler',
    }),
    ...mapAlarmActions({
      fetchOpenAlarmsListWithoutStore: 'fetchOpenAlarmsListWithoutStore',
    }),
    ...mapServiceActions({
      fetchServiceAlarmsWithoutStore: 'fetchAlarmsWithoutStore',
    }),

    showMap(map) {
      this.activeMapId = map._id;
      this.previousMaps.push({ ...this.mapState });
    },

    clearPreviousMaps() {
      this.activeMapId = this.widget.parameters.map;
      this.previousMaps = [];
    },

    backToBreadcrumb({ index }) {
      this.mapState = this.previousMaps[index];
      this.activeMapId = this.mapState._id;

      this.previousMaps = this.previousMaps.slice(0, index);
    },

    updateCategory(category) {
      const categoryId = category && category._id;

      this.updateContentInUserPreference({
        category: categoryId,
      });

      this.query = {
        ...this.query,

        category: categoryId,
      };
    },

    async fetchList() {
      this.pending = true;

      const params = this.getQuery();

      this.mapState = await this.fetchMapStateWithoutStore({
        id: this.activeMapId,
        params: pick(params, ['filters', 'category']),
      });

      this.pending = false;
    },

    showAlarmListModal(point) {
      const widget = generatePreparedDefaultAlarmListWidget();

      widget.parameters.widgetColumns = this.widget.parameters.alarmsColumns;

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
          fetchList: async (params) => {
            const { entity } = point;

            if (entity.type === ENTITY_TYPES.service) {
              return this.fetchServiceAlarmsWithoutStore({ id: entity._id, params });
            }

            const alarm = await this.fetchOpenAlarmsListWithoutStore({
              params: { ...params, _id: point.entity._id },
            });
            const data = alarm ? [alarm] : [];

            return {
              meta: { total_count: data.length },
              data,
            };
          },
        },
      });
    },
  },
};
</script>
