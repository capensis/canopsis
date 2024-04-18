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
      <v-flex xs5>
        <v-layout align-center>
          <template v-if="hasAccessToUserFilter">
            <filter-selector
              :value="query.filter"
              :locked-value="query.lockedFilter"
              :label="$t('settings.selectAFilter')"
              :filters="userPreference.filters"
              :locked-filters="widget.filters"
              :disabled="!hasAccessToListFilters"
              @input="updateSelectedFilter"
            />
            <filters-list-btn
              v-if="hasAccessToAddFilter || hasAccessToEditFilter"
              :widget-id="widget._id"
              :addable="hasAccessToAddFilter"
              :editable="hasAccessToEditFilter"
              :entity-types="[$constants.ENTITY_TYPES.service]"
              with-entity
              with-service-weather
              private
              entity-counters-type
            />
          </template>
          <c-enabled-field
            v-if="isHideGrayEnabled"
            :value="query.hide_grey"
            :label="$t('serviceWeather.hideGrey')"
            class="ml-3"
            @input="updateHideGray"
          />
        </v-layout>
      </v-flex>
    </v-layout>
    <v-fade-transition
      v-if="servicesPending"
      key="progress"
      mode="out-in"
    >
      <v-progress-linear
        class="progress-linear-absolute--top"
        height="2"
        indeterminate
      />
    </v-fade-transition>
    <v-layout
      key="content"
      class="fill-height"
      wrap
    >
      <v-layout
        v-if="hasNoData"
        justify-center
      >
        <v-alert
          v-if="servicesError"
          type="error"
        >
          <v-layout align-center>
            <div class="mr-4">
              {{ $t('errors.default') }}
            </div>
            <c-help-icon
              icon="help"
              top
            >
              <div v-if="servicesError.name">
                {{ $t('common.name') }}: {{ servicesError.name }}
              </div>
              <div v-if="servicesError.description">
                {{ $t('common.description') }}: {{ servicesError.description }}
              </div>
            </c-help-icon>
          </v-layout>
        </v-alert>
        <v-alert type="info">
          {{ $t('common.noData') }}
        </v-alert>
      </v-layout>
      <template v-else>
        <v-flex
          v-for="service in services"
          :key="service._id"
          :class="flexSize"
        >
          <service-weather-item
            :service="service"
            :action-required-blinking="actionRequiredSettings.is_blinking"
            :action-required-color="actionRequiredSettings.color"
            :action-required-icon="actionRequiredSettings.icon_name"
            :show-alarms-button="isBothModalType && hasAlarmsListAccess"
            :show-variables-help-button="hasVariablesHelpAccess"
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
        </v-flex>
      </template>
    </v-layout>
  </div>
</template>

<script>
import { MODALS, SERVICE_WEATHER_WIDGET_MODAL_TYPES, USERS_PERMISSIONS } from '@/constants';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';
import { getEntityColor } from '@/helpers/entities/entity/color';

import { permissionsWidgetsServiceWeatherFilters } from '@/mixins/permissions/widgets/service-weather/filters';
import { permissionsWidgetsServiceWeatherCategory } from '@/mixins/permissions/widgets/service-weather/category';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { entitiesServiceMixin } from '@/mixins/entities/service';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { authMixin } from '@/mixins/auth';
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';

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
    widgetFilterSelectMixin,
    entitiesAlarmTagMixin,
    entitiesServiceMixin,
    widgetFetchQueryMixin,
    authMixin,
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

    hasMoreInfosAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.moreInfos);
    },

    hasAlarmsListAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList);
    },

    hasVariablesHelpAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.variablesHelp);
    },

    actionRequiredSettings() {
      return this.widget.parameters.actionRequiredSettings ?? {};
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
    showAdditionalInfoModal(service) {
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
          color: getEntityColor(service, this.widget.parameters.colorIndicator),
          service,
          widgetParameters: this.widget.parameters,
        },
      });
    },

    showAlarmListModal(service) {
      try {
        const widget = generatePreparedDefaultAlarmListWidget();

        widget.parameters = {
          ...widget.parameters,
          ...this.widget.parameters.alarmsList,

          serviceDependenciesColumns: this.widget.parameters.serviceDependenciesColumns,
        };

        this.$modals.show({
          name: MODALS.alarmsList,
          config: {
            widget,
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

      if (!this.alarmTagsPending) {
        this.fetchAlarmTagsList({ params: { paginate: false } });
      }
    },
  },
};
</script>
