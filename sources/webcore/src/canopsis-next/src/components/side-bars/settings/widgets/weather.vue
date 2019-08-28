<template lang="pug">
  div
    v-list.pt-0(expand)
      field-row-grid-size(
      :rowId.sync="settings.rowId",
      :size.sync="settings.widget.size",
      :availableRows="availableRows",
      @createRow="createRow"
      )
      v-divider
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-periodic-refresh(v-model="settings.widget.parameters.periodicRefresh")
      v-divider
      field-filter-editor(v-model="settings.widget.parameters.mfilter", :hidden-fields="['title']")
      v-divider
      v-list-group(data-test='alarmsList')
        v-list-tile(slot="activator") {{ $t('settings.titles.alarmListSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-columns(v-model="settings.widget.parameters.alarmsList.widgetColumns", withHtml)
          v-divider
          field-default-elements-per-page(v-model="settings.widget.parameters.alarmsList.itemsPerPage")
          v-divider
          field-info-popup(v-model="settings.widget.parameters.alarmsList.infoPopups")
          v-divider
          field-text-editor(
          v-model="settings.widget.parameters.alarmsList.moreInfoTemplate",
          :title="$t('settings.moreInfosModal')"
          )
      v-divider
      field-number(data-test='widgetLimit', v-model="settings.widget.parameters.limit", :title="$t('settings.limit')")
      v-divider
      v-list-group(data-test="advancedSettings")
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(
          v-model="settings.widget.parameters.sort",
          :columns="sortColumns",
          :columnsLabel="$t('settings.orderBy')"
          )
          v-divider
          field-weather-template(
          data-test="widgetTemplateWeatherItem",
          v-model="settings.widget.parameters.blockTemplate",
          :title="$t('settings.weatherTemplate')"
          )
          v-divider
          field-weather-template(
          data-test="widgetTemplateModal",
          v-model="settings.widget.parameters.modalTemplate",
          :title="$t('settings.modalTemplate')"
          )
          v-divider
          field-weather-template(
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
          v-list-group(data-test="widgetMargin")
            v-list-tile(slot="activator") {{ $t('settings.margin.title') }}
            v-list.grey.lighten-4.px-2.py-0(expand)
              field-slider(
              data-test="widget-margin-top",
              v-model="settings.widget.parameters.margin.top",
              :title="$t('settings.margin.top')",
              :min="0",
              :max="5",
              )
              v-divider
              field-slider(
              data-test="widget-margin-right",
              v-model="settings.widget.parameters.margin.right",
              :title="$t('settings.margin.right')",
              :min="0",
              :max="5",
              )
              v-divider
              field-slider(
              data-test="widget-margin-bottom",
              v-model="settings.widget.parameters.margin.bottom",
              :title="$t('settings.margin.bottom')",
              :min="0",
              :max="5",
              )
              v-divider
              field-slider(
              data-test="widget-margin-left",
              v-model="settings.widget.parameters.margin.left",
              :title="$t('settings.margin.left')",
              :min="0",
              :max="5",
              )
          v-divider
          field-slider(
          data-test="widgetHeightFactory"
          v-model="settings.widget.parameters.heightFactor",
          :title="$t('settings.height')",
          :min="1",
          :max="20",
          )
          v-divider
          field-modal-type(v-model="settings.widget.parameters.modalType")
    v-btn.primary(data-test="submitWeather", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';
import sideBarSettingsWidgetAlarmMixin from '@/mixins/side-bar/settings/widgets/alarm';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldWeatherTemplate from './fields/weather/weather-template.vue';
import FieldGridSize from './fields/common/grid-size.vue';
import FieldSlider from './fields/common/slider.vue';
import FieldModalType from './fields/weather/modal-type.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldInfoPopup from './fields/alarm/info-popup.vue';
import FieldTextEditor from './fields/common/text-editor.vue';
import FieldNumber from './fields/common/number.vue';

export default {
  name: SIDE_BARS.weatherSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldPeriodicRefresh,
    FieldFilterEditor,
    FieldDefaultSortColumn,
    FieldWeatherTemplate,
    FieldGridSize,
    FieldSlider,
    FieldModalType,
    FieldColumns,
    FieldDefaultElementsPerPage,
    FieldInfoPopup,
    FieldTextEditor,
    FieldNumber,
  },
  mixins: [widgetSettingsMixin, sideBarSettingsWidgetAlarmMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: this.prepareWidgetWithAlarmParametersSettings(cloneDeep(widget), true),
      },
      sortColumns: [
        { label: 'name', value: 'name' },
        { label: 'state', value: 'state' },
        { label: 'status', value: 'status' },
        { label: 'criticity', value: 'criticity' },
        { label: 'org', value: 'org' },
      ],
    };
  },
  methods: {
    prepareWidgetSettings() {
      const { widget } = this.settings;

      return this.prepareWidgetWithAlarmParametersSettings(widget);
    },
  },
};
</script>

