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
      field-filters(
        :entitiesType="$constants.ENTITIES_TYPES.alarm",
        :filters.sync="settings.widget.parameters.viewFilters",
        hideSelect
      )
      v-divider
      v-list-group()
        v-list-tile(slot="activator") {{ $t('settings.titles.alarmListSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-columns(v-model="settings.widget.parameters.alarmsList.widgetColumns", withHtml)
          v-divider
          field-default-elements-per-page(v-model="settings.widget.parameters.alarmsList.itemsPerPage")
          v-divider
          field-info-popup(
            v-model="settings.widget.parameters.alarmsList.infoPopups",
            :columns="settings.widget.parameters.alarmsList.widgetColumns"
          )
          v-divider
          field-text-editor(
            v-model="settings.widget.parameters.alarmsList.moreInfoTemplate",
            :title="$t('settings.moreInfosModal')"
          )
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-template(
            v-model="settings.widget.parameters.blockTemplate",
            :title="$t('settings.tileTemplate')"
          )
          v-divider
          field-grid-size(
            v-model="settings.widget.parameters.columnSM",
            :title="$t('settings.columnSM')"
          )
          v-divider
          field-grid-size(
            v-model="settings.widget.parameters.columnMD",
            :title="$t('settings.columnMD')"
          )
          v-divider
          field-grid-size(
            v-model="settings.widget.parameters.columnLG",
            :title="$t('settings.columnLG')"
          )
          v-divider
          v-list-group
            v-list-tile(slot="activator") {{ $t('settings.margin.title') }}
            v-list.grey.lighten-4.px-2.py-0(expand)
              field-slider(
                v-model="settings.widget.parameters.margin.top",
                :title="$t('settings.margin.top')",
                :min="0",
                :max="5"
              )
              v-divider
              field-slider(
                v-model="settings.widget.parameters.margin.right",
                :title="$t('settings.margin.right')",
                :min="0",
                :max="5"
              )
              v-divider
              field-slider(
                v-model="settings.widget.parameters.margin.bottom",
                :title="$t('settings.margin.bottom')",
                :min="0",
                :max="5"
              )
              v-divider
              field-slider(
                v-model="settings.widget.parameters.margin.left",
                :title="$t('settings.margin.left')",
                :min="0",
                :max="5"
              )
          v-divider
          counter-levels-form(v-model="settings.widget.parameters.levels")
      v-divider
    v-btn.primary(data-test="submitWeather", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';
import sideBarSettingsWidgetAlarmMixin from '@/mixins/side-bar/settings/widgets/alarm';

import FieldSortColumn from './fields/weather/sort-column.vue';
import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldTemplate from './fields/common/template.vue';
import FieldGridSize from './fields/common/grid-size.vue';
import FieldSlider from './fields/common/slider.vue';
import FieldModalType from './fields/weather/modal-type.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldInfoPopup from './fields/alarm/info-popup.vue';
import FieldTextEditor from './fields/common/text-editor.vue';
import FieldNumber from './fields/common/number.vue';
import FieldFilters from './fields/common/filters.vue';
import CounterLevelsForm from './forms/counter-levels.vue';

export default {
  name: SIDE_BARS.weatherSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldSortColumn,
    FieldRowGridSize,
    FieldTitle,
    FieldPeriodicRefresh,
    FieldFilterEditor,
    FieldDefaultSortColumn,
    FieldTemplate,
    FieldGridSize,
    FieldSlider,
    FieldModalType,
    FieldColumns,
    FieldDefaultElementsPerPage,
    FieldInfoPopup,
    FieldTextEditor,
    FieldNumber,
    FieldFilters,
    CounterLevelsForm,
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

