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
    component(v-if="mapState", :is="component", :map="mapState")
</template>

<script>
import { MAP_TYPES } from '@/constants';

import { permissionsWidgetsMapCategory } from '@/mixins/permissions/widgets/map/category';
import { permissionsWidgetsMapFilters } from '@/mixins/permissions/widgets/map/filters';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { entitiesMapMixin } from '@/mixins/entities/map';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/filters-list-btn.vue';

export default {
  components: {
    FilterSelector,
    FiltersListBtn,
  },
  mixins: [
    permissionsWidgetsMapCategory,
    permissionsWidgetsMapFilters,
    widgetPeriodicRefreshMixin,
    widgetFilterSelectMixin,
    widgetFetchQueryMixin,
    entitiesMapMixin,
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
    };
  },
  computed: {
    component() {
      return {
        [MAP_TYPES.geo]: 'span',
        [MAP_TYPES.flowchart]: 'span',
        [MAP_TYPES.mermaid]: 'span',
        [MAP_TYPES.treeOfDependencies]: 'span',
      }[this.mapState.type];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
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
        id: this.widget.parameters.map,
        params: this.query,
      });

      this.pending = false;
    },
  },
};
</script>
