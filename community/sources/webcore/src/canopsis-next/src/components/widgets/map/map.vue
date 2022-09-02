<template lang="pug">
  div.pa-2
    v-layout.mx-1(wrap)
      v-flex(v-if="hasAccessToCategory", xs3)
        c-entity-category-field.mr-3(:category="query.category", @input="updateCategory")
      v-flex(v-if="hasAccessToUserFilter", xs4)
        v-layout(row, wrap, align-center)
          filter-selector(
            :label="$t('settings.selectAFilter')",
            :filters="userPreference.filters",
            :locked-filters="widget.filters",
            :value="mainFilter",
            :disabled="!hasAccessToListFilters && !hasAccessToUserFilter",
            @input="updateSelectedFilter"
          )
          filters-list-btn(
            :widget-id="widget._id",
            :addable="hasAccessToAddFilter",
            :editable="hasAccessToEditFilter",
            with-entity,
            with-service-weather,
            private
          )

    template(v-if="mapState")
      v-fade-transition(v-if="pending", key="progress", mode="out-in")
        v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
    template(v-else)
      v-layout.pa-4(v-if="pending", justify-center)
        v-progress-circular(color="primary", indeterminate)

    map-breadcrumbs.mb-2(
      v-if="previousMaps.length",
      :previous-maps="previousMaps",
      :active-map="mapState",
      :pending="pending",
      @click="backToBreadcrumb"
    )
    component(
      v-if="mapState",
      :is="component",
      :map="mapState",
      :popup-template="widget.parameters.entity_info_template",
      :color-indicator="widget.parameters.color_indicator",
      :pbehavior-enabled="widget.parameters.entities_under_pbehavior_enabled",
      popup-actions,
      @show:map="showMap",
      @show:alarms="showAlarmListModal"
    )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { pick } from 'lodash';

import { ENTITY_TYPES, MAP_TYPES, MODALS } from '@/constants';

import { generateDefaultAlarmListWidget } from '@/helpers/entities';

import { permissionsWidgetsMapCategory } from '@/mixins/permissions/widgets/map/category';
import { permissionsWidgetsMapFilters } from '@/mixins/permissions/widgets/map/filters';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/filters-list-btn.vue';
import MermaidPreview from '@/components/other/map/partials/mermaid-preview.vue';
import GeomapPreview from '@/components/other/map/partials/geomap-preview.vue';
import FlowchartPreview from '@/components/other/map/partials/flowchart-preview.vue';

import MapBreadcrumbs from './partials/map-breadcrumbs.vue';

const { mapActions: mapMapActions } = createNamespacedHelpers('map');
const { mapActions: mapActiveViewActions } = createNamespacedHelpers('activeView');
const { mapActions: mapAlarmActions } = createNamespacedHelpers('alarm');
const { mapActions: mapServiceActions } = createNamespacedHelpers('service');

export default {
  components: {
    MapBreadcrumbs,
    FilterSelector,
    FiltersListBtn,
    MermaidPreview,
    GeomapPreview,
    FlowchartPreview,
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
        [MAP_TYPES.treeOfDependencies]: 'span',
      }[this.mapState.type];
    },
  },
  watch: {
    activeMapId: 'fetchList',
  },
  created() {
    this.registerEditingOffHandler(this.clearPreviousMaps);
  },
  beforeDestroy() {
    this.unregisterEditingOffHandler(this.clearPreviousMaps);
  },
  mounted() {
    if (this.editing) {
      this.fetchList();
    }
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

      this.mapState = await this.fetchMapStateWithoutStore({
        id: this.activeMapId,
        params: pick(this.query, ['filter', 'category']),
      });

      this.pending = false;
    },

    showAlarmListModal(point) {
      try {
        const widget = generateDefaultAlarmListWidget();

        widget.parameters.widgetColumns = this.widget.parameters.alarms_columns;

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
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
