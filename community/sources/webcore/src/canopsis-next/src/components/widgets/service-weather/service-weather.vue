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
              private,
              entity-counters-type
            )
          c-enabled-field.ml-3(
            v-if="isHideGrayEnabled",
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
          service-weather-item(
            :service="service",
            :action-required-blinking="actionRequiredSettings.is_blinking",
            :action-required-color="actionRequiredSettings.color",
            :action-required-icon="actionRequiredSettings.icon_name",
            :show-alarms-button="isBothModalType && hasAlarmsListAccess",
            :show-variables-help-button="hasVariablesHelpAccess",
            :template="widget.parameters.blockTemplate",
            :height-factor="widget.parameters.heightFactor",
            :color-indicator="widget.parameters.colorIndicator",
            :priority-enabled="widget.parameters.isPriorityEnabled",
            :secondary-icon-enabled="widget.parameters.isSecondaryIconEnabled",
            :counters-settings="widget.parameters.counters",
            :margin="widget.parameters.margin",
            @show:service="showAdditionalInfoModal(service)",
            @show:alarms="showAlarmListModal(service)"
          )
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
