<template>
  <div class="map-widget pa-2">
    <v-row
      v-if="hasAccessToCategory || hasAccessToUserFilter"
      class="map-widget__filters mx-1 mb-2 flex-grow-0"
      align="center"
      no-gutters
    >
      <v-col
        v-if="hasAccessToCategory"
        class="map-widget__filter-cell"
        cols="3"
      >
        <c-entity-category-field
          :category="query.category"
          class="mr-3"
          @update:category="updateCategory"
        />
      </v-col>

      <v-col
        v-if="hasAccessToUserFilter"
        class="map-widget__filter-cell d-flex align-center"
        cols="4"
      >
        <filter-selector
          :model-value="query.filter"
          :locked-value="query.lockedFilter"
          :label="$t('settings.selectAFilter')"
          :filters="userPreferenceFilters"
          :locked-filters="widget.filters"
          @update:model-value="updateSelectedFilter"
        />

        <filters-list-btn
          :widget-id="widget._id"
          addable
          editable
          with-entity
          private
        />
      </v-col>
    </v-row>

    <v-fade-transition v-if="mapState" mode="out-in">
      <v-progress-linear
        v-if="pending"
        key="progress"
        class="progress-linear-absolute--top"
        height="2"
        indeterminate
      />
    </v-fade-transition>

    <c-layout v-else-if="pending" class="pa-4" justify="center">
      <v-progress-circular color="primary" indeterminate />
    </c-layout>

    <map-breadcrumbs
      v-if="previousMaps.length"
      :previous-maps="previousMaps"
      :active-map="mapState"
      :pending="pending"
      class="mb-2"
      @click="backToBreadcrumb"
    />

    <component
      v-if="component && mapState"
      :is="component"
      :map="mapState"
      :columns="preparedWidget.parameters.entitiesColumns"
      :popup-template="preparedWidget.parameters.entity_info_template"
      :color-indicator="preparedWidget.parameters.color_indicator"
      :pbehavior-enabled="preparedWidget.parameters.entities_under_pbehavior_enabled"
      popup-actions
      @show:map="showMap"
      @show:alarms="showAlarmListModal"
    />
  </div>
</template>

<script>
import { pick } from 'lodash';
import { mapActions } from 'pinia';

import { ENTITY_TYPES, MAP_TYPES, MODALS } from '@/constants';

import { useAlarmStore } from '@/stores/modules/alarm';
import { useMapStore } from '@/stores/modules/map';
import { useServiceStore } from '@/stores/modules/service';

import { durationToSeconds } from '@/helpers/date/duration';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';
import { prepareMapWidget } from '@/helpers/entities/widget/forms/map';

import { activeViewMixin } from '@/mixins/active-view';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';
import { permissionsWidgetsMapCategory } from '@/mixins/permissions/widgets/map/category';
import { permissionsWidgetsMapFilters } from '@/mixins/permissions/widgets/map/filters';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';
import FlowchartPreview from '@/components/other/map/partials/flowchart-preview.vue';
import GeomapPreview from '@/components/other/map/partials/geomap-preview.vue';
import MermaidPreview from '@/components/other/map/partials/mermaid-preview.vue';
import TreeOfDependenciesPreview from '@/components/other/map/partials/tree-of-dependencies-preview.vue';

import MapBreadcrumbs from './partials/map-breadcrumbs.vue';

const MAP_COMPONENTS_BY_TYPES = {
  [MAP_TYPES.geo]: GeomapPreview,
  [MAP_TYPES.flowchart]: FlowchartPreview,
  [MAP_TYPES.mermaid]: MermaidPreview,
  [MAP_TYPES.treeOfDependencies]: TreeOfDependenciesPreview,
};

export default {
  inject: ['$periodicRefresh'],

  components: {
    FilterSelector,
    FiltersListBtn,
    FlowchartPreview,
    GeomapPreview,
    MapBreadcrumbs,
    MermaidPreview,
    TreeOfDependenciesPreview,
  },

  mixins: [
    activeViewMixin,
    entitiesUserPreferenceMixin,
    permissionsWidgetsMapCategory,
    permissionsWidgetsMapFilters,
  ],

  props: {
    widget: {
      type: Object,
      required: true,
    },
    visible: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      pending: false,
      mapState: undefined,
      previousMaps: [],
      activeMapId: this.widget.parameters.map,
      query: {
        category: undefined,
        filter: null,
        lockedFilter: this.widget.parameters.mainFilter ?? null,
      },
      periodicRefreshInterval: undefined,
      intervalStartedAt: undefined,
      intervalDelay: undefined,
    };
  },

  computed: {
    component() {
      return MAP_COMPONENTS_BY_TYPES[this.mapState?.type];
    },

    preparedWidget() {
      return prepareMapWidget(this.widget);
    },

    userPreferenceFilters() {
      return this.userPreference?.filters ?? [];
    },

    queryFilters() {
      return [this.query.filter, this.query.lockedFilter].flat().filter(Boolean);
    },

    periodicRefreshEnabled() {
      return this.widget.parameters.periodic_refresh?.enabled;
    },

    periodicRefreshSeconds() {
      const { value, unit } = this.widget.parameters.periodic_refresh ?? {};

      return durationToSeconds({ value, unit });
    },
  },

  watch: {
    'widget.parameters.periodic_refresh': {
      immediate: true,
      handler(value, oldValue) {
        const periodicRefresh = value;
        const oldPeriodicRefresh = oldValue ?? {};

        if (this.activeViewPeriodicRefreshPaused || !this.visible) {
          return;
        }

        if (periodicRefresh?.enabled && periodicRefresh?.value) {
          const valueIsChanged = periodicRefresh.value !== oldPeriodicRefresh.value;
          const enabledIsChanged = periodicRefresh.enabled !== oldPeriodicRefresh.enabled;

          if (valueIsChanged || enabledIsChanged) {
            if (this.periodicRefreshInterval) {
              this.stopPeriodicRefresh();
            }

            if (!this.periodicRefreshSeconds) {
              return;
            }

            this.startPeriodicRefresh();
          }
        } else {
          this.stopPeriodicRefresh();
        }
      },
    },

    activeViewPeriodicRefreshPaused(value) {
      if (!this.visible) {
        return;
      }

      if (value) {
        this.pausePeriodicRefresh();

        return;
      }

      this.resumePeriodicRefresh();
    },

    activeMapId(id, oldId) {
      if (id !== oldId) {
        this.fetchList();
      }
    },

    visible: {
      immediate: true,
      handler(visible) {
        if (visible) {
          this.prepareQueryAndFetch();

          if (!this.activeViewPeriodicRefreshPaused) {
            this.startPeriodicRefresh();
          }
        } else {
          this.stopPeriodicRefresh();
        }
      },
    },

    'widget.parameters.map': function mapParameter(mapId) {
      this.activeMapId = mapId;
    },
  },

  created() {
    this.registerEditingOffHandler(this.clearPreviousMaps);
    this.$periodicRefresh.registerChild(this.refresh);
  },

  beforeUnmount() {
    this.unregisterEditingOffHandler(this.clearPreviousMaps);
    this.$periodicRefresh.unregisterChild(this.refresh);
    this.stopPeriodicRefresh();
  },

  methods: {
    ...mapActions(useMapStore, {
      fetchMapStateWithoutStore: 'fetchItemStateWithoutStore',
    }),
    ...mapActions(useAlarmStore, {
      fetchOpenAlarmsListWithoutStore: 'fetchOpenAlarmsListWithoutStore',
    }),
    ...mapActions(useServiceStore, {
      fetchServiceAlarmsWithoutStore: 'fetchAlarmsWithoutStore',
    }),

    getQuery() {
      return {
        filters: this.queryFilters,
        category: this.query.category,
      };
    },

    refresh() {
      if (this.visible) {
        return this.fetchList();
      }

      return Promise.resolve();
    },

    async prepareQueryAndFetch() {
      if (this.widget._id) {
        await this.fetchUserPreference({ id: this.widget._id });
      }

      this.setQueryFromWidgetAndPreference();
      this.fetchList();
    },

    setQueryFromWidgetAndPreference() {
      const content = this.userPreference?.content ?? {};

      this.query = {
        ...this.query,
        category: content.category,
        filter: content.mainFilter ?? null,
        lockedFilter: this.widget.parameters.mainFilter ?? null,
      };
    },

    updateSelectedFilter(filter) {
      this.updateContentInUserPreference({
        mainFilter: filter,
      });

      this.query = {
        ...this.query,
        filter,
      };

      this.fetchList();
    },

    updateCategory(category) {
      const categoryId = category?._id;

      this.updateContentInUserPreference({
        category: categoryId,
      });

      this.query = {
        ...this.query,
        category: categoryId,
      };

      this.fetchList();
    },

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

    async fetchList() {
      if (!this.activeMapId) {
        return;
      }

      try {
        this.pending = true;

        const params = this.getQuery();

        this.mapState = await this.fetchMapStateWithoutStore({
          id: this.activeMapId,
          params: pick(params, ['filters', 'category']),
        });
      } finally {
        this.pending = false;
      }
    },

    showAlarmListModal(point) {
      const widget = generatePreparedDefaultAlarmListWidget();

      widget.parameters = widget.parameters ?? {};
      widget.parameters.widgetColumns = this.preparedWidget.parameters.alarmsColumns;

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

    startPeriodicRefresh() {
      this.stopPeriodicRefresh();

      if (!this.periodicRefreshEnabled || !this.periodicRefreshSeconds) {
        return;
      }

      const timeout = this.periodicRefreshSeconds * 1000;

      this.periodicRefreshInterval = setInterval(this.fetchList, timeout);
      this.intervalStartedAt = Date.now();
    },

    resumePeriodicRefresh() {
      setTimeout(() => {
        this.fetchList();
        this.startPeriodicRefresh();
      }, this.intervalDelay || 0);
    },

    pausePeriodicRefresh() {
      if (!this.periodicRefreshInterval || !this.intervalStartedAt) {
        return;
      }

      clearInterval(this.periodicRefreshInterval);

      this.intervalDelay =
        this.periodicRefreshSeconds * 1000 - (Date.now() - this.intervalStartedAt);
    },

    stopPeriodicRefresh() {
      if (this.periodicRefreshInterval) {
        clearInterval(this.periodicRefreshInterval);
        this.periodicRefreshInterval = undefined;
      }

      this.intervalStartedAt = undefined;
      this.intervalDelay = undefined;
    },
  },
};
</script>

<style lang="scss" scoped>
.map-widget {
  &__filter-cell {
    min-width: 0;

    :deep(.filter-selector) {
      flex: 1 1 auto;
      min-width: 0;
    }
  }
}
</style>
