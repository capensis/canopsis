<template lang="pug">
  div.pa-2
    v-layout.mx-1(wrap)
      v-flex(v-if="hasAccessToCategory", xs3)
        c-entity-category-field.mr-3(:category="query.category", @input="updateCategory")
      v-flex(v-if="hasAccessToUserFilter", xs4)
        filter-selector(
          :value="mainFilter",
          :filters="viewFilters",
          :locked-filters="widgetViewFilters",
          :condition="mainFilterCondition",
          :label="$t('settings.selectAFilter')",
          :has-access-to-user-filter="hasAccessToUserFilter",
          @input="updateSelectedFilter",
          @update:condition="updateSelectedCondition",
          @update:filters="updateFilters"
        )
    v-fade-transition(v-if="servicesPending", key="progress", mode="out-in")
      v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
    v-layout.fill-height(key="content", wrap)
      v-alert(v-if="hasNoData && servicesError", :value="true", type="error")
        v-layout(align-center)
          div.mr-4 {{ $t('errors.default') }}
          v-tooltip(top)
            v-icon(slot="activator") help
            div(v-if="servicesError.name") {{ $t('common.name') }}: {{ servicesError.name }}
            div(v-if="servicesError.description") {{ $t('common.description') }}: {{ servicesError.description }}
      v-alert(v-else-if="hasNoData", :value="true", type="info") {{ $t('tables.noData') }}
      template(v-else)
        v-flex(v-for="service in services", :key="service._id", :class="flexSize")
          service-weather-item.weather-item(
            :service="service",
            :widget="widget",
            :template="widget.parameters.blockTemplate"
          )
</template>

<script>
import { omit } from 'lodash';

import { permissionsWidgetsServiceWeatherFilters } from '@/mixins/permissions/widgets/service-weather/filters';
import { permissionsWidgetsServiceWeatherCategory } from '@/mixins/permissions/widgets/service-weather/category';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import widgetFilterSelectMixin from '@/mixins/widget/filter-select';
import entitiesServiceMixin from '@/mixins/entities/service';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

import FilterSelector from '@/components/other/filter/filter-selector.vue';

import ServiceWeatherItem from './service-weather-item.vue';

export default {
  components: {
    ServiceWeatherItem,
    FilterSelector,
  },
  mixins: [
    permissionsWidgetsServiceWeatherFilters,
    permissionsWidgetsServiceWeatherCategory,
    widgetPeriodicRefreshMixin,
    widgetFilterSelectMixin,
    entitiesServiceMixin,
    widgetFetchQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    flexSize() {
      return [
        `xs${this.widget.parameters.columnSM}`,
        `md${this.widget.parameters.columnMD}`,
        `lg${this.widget.parameters.columnLG}`,
      ];
    },
    hasNoData() {
      return this.services.length === 0;
    },
  },
  methods: {
    updateCategory(category) {
      const categoryId = category && category._id;

      this.updateWidgetPreferencesInUserPreference({
        ...this.userPreference.widget_preferences,

        category: categoryId,
      });

      this.query = {
        ...this.query,

        category: categoryId,
      };
    },

    getQuery() {
      const query = omit(this.query, [
        'page',
        'sortKey',
        'sortDir',
      ]);

      if (this.query.sortKey) {
        query.sort_by = this.query.sortKey;
        query.sort = this.query.sortDir.toLowerCase();
      }

      return query;
    },

    fetchList() {
      this.fetchServicesList({
        params: this.getQuery(),
        widgetId: this.widget._id,
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .weather-item {
    height: 100%;
  }
</style>