<template>
  <div class="service-weather-widget pa-2">
    <v-row
      v-if="hasAccessToCategory || hasAccessToUserFilter"
      class="service-weather-widget__filters mx-1 mb-2 flex-grow-0"
      align="center"
      no-gutters
    >
      <v-col
        v-if="hasAccessToCategory"
        class="service-weather-widget__filter-cell"
        cols="3"
      >
        <c-entity-category-field
          :category="query.category"
          class="mr-3"
          @update:category="updateCategory"
        />
      </v-col>

      <v-col
        v-if="hasAccessToUserFilter || isHideGrayEnabled"
        class="service-weather-widget__filter-cell d-flex align-center"
        cols="5"
      >
        <template v-if="hasAccessToUserFilter">
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
            :entity-types="[ENTITY_TYPES.service]"
            addable
            editable
            with-entity
            with-service-weather
            private
            entity-counters-type
          />
        </template>

        <c-enabled-field
          v-if="isHideGrayEnabled"
          :model-value="query.hide_grey"
          :label="$t('serviceWeather.hideGrey')"
          class="ml-3 mt-0"
          color="primary"
          hide-details
          inline
          no-margin
          @update:model-value="updateHideGray"
        />
      </v-col>
    </v-row>

    <v-fade-transition v-if="servicesPending" mode="out-in">
      <v-progress-linear
        key="progress"
        class="progress-linear-absolute--top"
        height="2"
        indeterminate
      />
    </v-fade-transition>

    <v-row class="fill-height" density="compact">
      <v-col v-if="hasNoData" cols="12">
        <c-layout justify="center" column>
          <v-alert v-if="servicesError" type="error">
            <c-layout align="center">
              <div class="mr-4">
                {{ $t('errors.default') }}
              </div>
              <c-help-icon icon="help" location="top">
                <div v-if="servicesError.name">
                  {{ $t('common.name') }}: {{ servicesError.name }}
                </div>
                <div v-if="servicesError.description">
                  {{ $t('common.description') }}: {{ servicesError.description }}
                </div>
              </c-help-icon>
            </c-layout>
          </v-alert>

          <v-alert type="info">
            {{ $t('common.noData') }}
          </v-alert>
        </c-layout>
      </v-col>

      <template v-else>
        <v-col
          v-for="service in services"
          :key="service._id"
          :cols="columnSize"
        >
          <service-weather-item
            :service="service"
            :action-required-blinking="actionRequiredSettings.is_blinking"
            :action-required-color="actionRequiredSettings.color"
            :action-required-icon="actionRequiredSettings.icon_name"
            :no-action-required-color="noActionRequiredSettings.color"
            :no-action-required-icon="noActionRequiredSettings.icon_name"
            :show-alarms-button="!hideActions && isBothModalType && hasAlarmsListAccess"
            :show-variables-help-button="!hideActions && hasVariablesHelpAccess"
            :template="widget.parameters.blockTemplate"
            :height-factor="widget.parameters.heightFactor"
            :color-indicator="widget.parameters.colorIndicator"
            :priority-enabled="widget.parameters.isPriorityEnabled"
            :secondary-icon-enabled="widget.parameters.isSecondaryIconEnabled"
            :counters-settings="widget.parameters.counters"
            :show-root-cause-by-state-click="showRootCauseByStateClick"
            :margin="widget.parameters.margin"
            @show:service="showAdditionalInfoModal(service)"
            @show:alarms="showAlarmListModal(service)"
            @show:root-cause="openRootCauseDiagram(service)"
          />
        </v-col>
      </template>
    </v-row>
  </div>
</template>

<script>
import { mapActions, mapState } from 'pinia';

import {
  ENTITY_TYPES,
  MODALS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  USER_PERMISSIONS,
} from '@/constants';

import { useWeatherServiceWidgetStore } from '@/stores/modules/weather-service-widget';

import { getEntityColorClass } from '@/helpers/entities/entity/color';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { authMixin } from '@/mixins/auth';
import { entitiesServiceMixin } from '@/mixins/entities/service';
import { permissionsWidgetsServiceWeatherCategory } from '@/mixins/permissions/widgets/service-weather/category';
import { permissionsWidgetsServiceWeatherFilters } from '@/mixins/permissions/widgets/service-weather/filters';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';

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
    entitiesServiceMixin,
    widgetFetchQueryMixin,
    authMixin,
  ],

  props: {
    widget: {
      type: Object,
      required: true,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    visible: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      ENTITY_TYPES,
    };
  },

  computed: {
    ...mapState(useWeatherServiceWidgetStore, {
      getWeatherServicesListByWidgetId: 'getListByWidgetId',
      getWeatherServicesPendingByWidgetId: 'getPendingByWidgetId',
      getWeatherServicesErrorByWidgetId: 'getErrorByWidgetId',
    }),

    columnSize() {
      return 12 / this.columnsCount;
    },

    columnsCount() {
      const displayName = this.$vuetify.display.name;

      if (['xs', 'sm'].includes(displayName)) {
        return this.widget.parameters.columnMobile;
      }

      if (displayName === 'md') {
        return this.widget.parameters.columnTablet;
      }

      return this.widget.parameters.columnDesktop;
    },

    services() {
      return this.getWeatherServicesListByWidgetId(this.widget._id);
    },

    servicesPending() {
      return this.getWeatherServicesPendingByWidgetId(this.widget._id);
    },

    servicesError() {
      return this.getWeatherServicesErrorByWidgetId(this.widget._id);
    },

    hasNoData() {
      return this.services.length === 0;
    },

    userPreferenceFilters() {
      return this.userPreference?.filters ?? [];
    },

    hasMoreInfosAccess() {
      return this.checkAccess(USER_PERMISSIONS.business.serviceWeather.actions.moreInfos);
    },

    hasAlarmsListAccess() {
      return this.checkAccess(USER_PERMISSIONS.business.serviceWeather.actions.alarmsList);
    },

    hasVariablesHelpAccess() {
      return this.checkAccess(USER_PERMISSIONS.business.serviceWeather.actions.variablesHelp);
    },

    actionRequiredSettings() {
      return this.widget.parameters.actionRequiredSettings ?? {};
    },

    noActionRequiredSettings() {
      return this.widget.parameters.noActionRequiredSettings ?? {};
    },

    isBothModalType() {
      return this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.both;
    },

    isAlarmListModalType() {
      return this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList;
    },

    isHideGrayEnabled() {
      return this.widget.parameters.isHideGrayEnabled ?? true;
    },

    showRootCauseByStateClick() {
      return this.widget.parameters.showRootCauseByStateClick ?? true;
    },
  },

  methods: {
    ...mapActions(useWeatherServiceWidgetStore, {
      fetchWeatherServicesList: 'fetchList',
    }),

    showAdditionalInfoModal(service) {
      if (this.hideActions) {
        return;
      }

      if (this.isAlarmListModalType && this.hasAlarmsListAccess) {
        this.showAlarmListModal(service);
      } else if (!this.isAlarmListModalType && this.hasMoreInfosAccess) {
        this.showMainInfoModal(service);
      }
    },

    showMainInfoModal(service) {
      this.$modals.show({
        name: MODALS.serviceEntities,
        config: {
          service,
          titleClass: getEntityColorClass(service, this.widget.parameters.colorIndicator),
          widgetParameters: this.widget.parameters,
        },
      });
    },

    showAlarmListModal(service) {
      try {
        const alarmListWidget = generatePreparedDefaultAlarmListWidget();

        alarmListWidget._id = this.widget._id;
        alarmListWidget.parameters = {
          ...alarmListWidget.parameters,
          ...this.widget.parameters.alarmsList,
          serviceDependenciesColumns: this.widget.parameters.serviceDependenciesColumns,
        };

        this.$modals.show({
          name: MODALS.alarmsList,
          config: {
            widget: alarmListWidget,
            title: this.$t('modals.alarmsList.prefixTitle', { prefix: service.name }),
            fetchList: params => this.fetchServiceAlarmsWithoutStore({ id: service._id, params }),
          },
        });
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    openRootCauseDiagram(service) {
      this.$modals.show({
        name: MODALS.entitiesRootCauseDiagram,
        config: {
          entity: service,
          colorIndicator: this.widget.parameters.rootCauseColorIndicator,
        },
      });
    },

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
      const categoryId = category?._id;

      this.updateContentInUserPreference({
        category: categoryId,
      });

      this.query = {
        ...this.query,
        category: categoryId,
      };
    },

    updateSelectedFilter(filter) {
      this.updateContentInUserPreference({ mainFilter: filter });

      this.query = {
        ...this.query,
        filter,
      };
    },

    fetchList() {
      return this.fetchWeatherServicesList({
        params: this.getQuery(),
        widgetId: this.widget._id,
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.service-weather-widget {
  &__filter-cell {
    min-width: 0;

    :deep(.filter-selector) {
      flex: 1 1 auto;
      min-width: 0;
    }
  }
}
</style>
