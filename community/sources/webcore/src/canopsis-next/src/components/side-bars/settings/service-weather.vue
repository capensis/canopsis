<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-periodic-refresh(v-model="settings.widget.parameters.periodicRefresh")
      v-divider
      template(v-if="hasAccessToListFilters")
        field-filters(
          v-model="settings.widget.parameters.mainFilter",
          :entities-type="$constants.ENTITIES_TYPES.entity",
          :filters.sync="settings.widget.parameters.viewFilters",
          :condition.sync="settings.widget.parameters.mainFilterCondition",
          :addable="hasAccessToAddFilter",
          :editable="hasAccessToEditFilter",
          @input="updateMainFilterUpdatedAt"
        )
        v-divider
      alarms-list-modal-form(v-model="settings.widget.parameters.alarmsList")
      v-divider
      field-number(
        v-model="settings.widget.parameters.limit",
        :title="$t('settings.limit')"
      )
      v-divider
      field-color-indicator(v-model="settings.widget.parameters.colorIndicator")
      v-divider
      field-columns(
        v-model="settings.widget.parameters.serviceDependenciesColumns",
        :label="$t('settings.treeOfDependenciesColumnNames')",
        with-color-indicator
      )
      v-divider
      v-list-group(data-test="advancedSettings")
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-sort-column(
            v-model="settings.widget.parameters.sort",
            :columns="sortColumns",
            :columns-label="$t('settings.orderBy')"
          )
          v-divider
          field-default-elements-per-page(v-model="settings.widget.parameters.modalItemsPerPage")
            span(slot="title") {{ $t('settings.defaultNumberOfElementsPerPage') }}
              span.font-italic.caption.ml-1 (Modal)
          v-divider
          field-template(
            data-test="widgetTemplateWeatherItem",
            v-model="settings.widget.parameters.blockTemplate",
            :title="$t('settings.weatherTemplate')"
          )
          v-divider
          field-template(
            data-test="widgetTemplateModal",
            v-model="settings.widget.parameters.modalTemplate",
            :title="$t('settings.modalTemplate')"
          )
          v-divider
          field-template(
            data-test="widgetTemplateEntities",
            v-model="settings.widget.parameters.entityTemplate",
            :title="$t('settings.entityTemplate')"
          )
          v-divider
          field-grid-size(
            data-test="columnSM",
            v-model="settings.widget.parameters.columnSM",
            :title="$t('settings.columnSM')"
          )
          v-divider
          field-grid-size(
            data-test="columnMD",
            v-model="settings.widget.parameters.columnMD",
            :title="$t('settings.columnMD')"
          )
          v-divider
          field-grid-size(
            data-test="columnLG",
            v-model="settings.widget.parameters.columnLG",
            :title="$t('settings.columnLG')"
          )
          v-divider
          margins-form(v-model="settings.widget.parameters.margin")
          v-divider
          field-slider(
            data-test="widgetHeightFactory",
            v-model="settings.widget.parameters.heightFactor",
            :title="$t('settings.height')",
            :min="1",
            :max="20"
          )
          v-divider
          field-counters-selector(
            v-model="settings.widget.parameters.counters",
            :title="$t('settings.counters')"
          )
          v-divider
          field-modal-type(v-model="settings.widget.parameters.modalType")
    v-btn.primary(data-test="submitWeather", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { sideBarSettingsWidgetAlarmMixin } from '@/mixins/side-bar/settings/widgets/alarm';
import { rightsWidgetsServiceWeatherListFilters } from '@/mixins/rights/widgets/service-weather/filters';

import FieldTitle from './fields/common/title.vue';
import FieldSortColumn from './fields/service-weather/sort-column.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldTemplate from './fields/common/template.vue';
import FieldGridSize from './fields/common/grid-size.vue';
import FieldSlider from './fields/common/slider.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import FieldModalType from './fields/service-weather/modal-type.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldNumber from './fields/common/number.vue';
import FieldCountersSelector from './fields/common/counters-selector.vue';
import FieldColorIndicator from './fields/common/color-indicator.vue';
import AlarmsListModalForm from './forms/alarms-list-modal.vue';
import MarginsForm from './forms/margins.vue';

export default {
  name: SIDE_BARS.serviceWeatherSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldTitle,
    FieldSortColumn,
    FieldPeriodicRefresh,
    FieldFilters,
    FieldColumns,
    FieldDefaultSortColumn,
    FieldTemplate,
    FieldGridSize,
    FieldSlider,
    FieldSwitcher,
    FieldModalType,
    FieldDefaultElementsPerPage,
    FieldNumber,
    FieldCountersSelector,
    FieldColorIndicator,
    AlarmsListModalForm,
    MarginsForm,
  },
  mixins: [
    widgetSettingsMixin,
    sideBarSettingsWidgetAlarmMixin,
    rightsWidgetsServiceWeatherListFilters,
  ],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: this.prepareWidgetWithAlarmParametersSettings(cloneDeep(widget), true),
      },
      sortColumns: [
        { label: 'name', value: 'name' },
        { label: 'state', value: 'state' },
      ],
    };
  },
  methods: {
    prepareWidgetSettings() {
      const { widget } = this.settings;

      return this.prepareWidgetWithAlarmParametersSettings(widget);
    },

    updateMainFilterUpdatedAt() {
      this.settings.widget.parameters.mainFilterUpdatedAt = Date.now();
    },
  },
};
</script>
