<template lang="pug">
  div.pa-2
    v-layout.mx-1(wrap)
      v-flex(v-if="hasAccessToCategory", xs3)
        c-entity-category-field.mr-3(:category="query.category", @input="updateCategory")
      v-flex(xs5)
        v-layout(row, align-center)
          template(v-if="hasAccessToUserFilter")
            filter-selector(
              :label="$t('settings.selectAFilter')",
              :filters="userPreference.filters",
              :locked-filters="widget.filters",
              :locked-value="lockedFilter",
              :value="mainFilter",
              :disabled="!hasAccessToListFilters",
              @input="updateSelectedFilter"
            )
            filters-list-btn(
              v-if="hasAccessToAddFilter || hasAccessToEditFilter",
              :widget-id="widget._id",
              :addable="hasAccessToAddFilter",
              :editable="hasAccessToEditFilter",
              :entity-types="[$constants.ENTITY_TYPES.service]",
              with-entity,
              with-service-weather,
              private
            )
          c-enabled-field.ml-3(
            v-if="isHideGrayEnabled"
            :value="query.hide_grey",
            :label="$t('serviceWeather.hideGrey')",
            @input="updateHideGray"
          )
    v-fade-transition(v-if="servicesPending", key="progress", mode="out-in")
      v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
    v-layout.fill-height(key="content", wrap)
      v-alert(v-if="hasNoData && servicesError", :value="true", type="error")
        v-layout(align-center)
          div.mr-4 {{ $t('errors.default') }}
          c-help-icon(icon="help", top)
            div(v-if="servicesError.name") {{ $t('common.name') }}: {{ servicesError.name }}
            div(v-if="servicesError.description") {{ $t('common.description') }}: {{ servicesError.description }}
      v-alert(v-else-if="hasNoData", :value="true", type="info") {{ $t('common.noData') }}
      template(v-else)
        v-flex(
          v-for="service in services",
          :key="service._id",
          :class="flexSize"
        )
          service-weather-item(:service="service", :widget="widget")
</template>

<script>
import { permissionsWidgetsServiceWeatherFilters } from '@/mixins/permissions/widgets/service-weather/filters';
import { permissionsWidgetsServiceWeatherCategory } from '@/mixins/permissions/widgets/service-weather/category';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import entitiesServiceMixin from '@/mixins/entities/service';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/filters-list-btn.vue';

import ServiceWeatherItem from './service-weather-item.vue';

export default {
  components: {
    FilterSelector,
    FiltersListBtn,
    ServiceWeatherItem,
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
      const columnsCount = {
        m: this.widget.parameters.columnMobile,
        t: this.widget.parameters.columnTablet,
        l: this.widget.parameters.columnDesktop,
        xl: this.widget.parameters.columnDesktop,
      }[this.$mq];

      return `xs${12 / columnsCount}`;
    },

    hasNoData() {
      return this.services.length === 0;
    },

    isHideGrayEnabled() {
      return this.widget.parameters.isHideGrayEnabled ?? true;
  },
  },
  methods: {
    updateHideGray(hideGrey) {
      this.updateContentInUserPreference({
        hide_grey: hideGrey,
      });

      this.query = {
        ...this.query,

        hide_grey: hideGrey,
      };
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

    fetchList() {
      this.fetchServicesList({
        params: this.getQuery(),
        widgetId: this.widget._id,
      });
    },
  },
};
</script>
