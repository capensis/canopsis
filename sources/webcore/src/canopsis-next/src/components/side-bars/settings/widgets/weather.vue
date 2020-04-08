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
      field-filter-editor(
        data-test="widgetFilterEditor",
        v-model="settings.widget.parameters.mfilter",
        :hidden-fields="['title']",
        :entitiesType="$constants.ENTITIES_TYPES.entity"
      )
      v-divider
      alarms-list-modal-form(v-model="settings.widget.parameters.alarmsList")
      v-divider
      field-number(
        data-test="widgetLimit",
        v-model="settings.widget.parameters.limit",
        :title="$t('settings.limit')"
      )
      v-divider
      v-list-group(data-test="advancedSettings")
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-sort-column(
            v-model="settings.widget.parameters.sort",
            :columns="sortColumns",
            :columnsLabel="$t('settings.orderBy')"
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
          field-modal-type(v-model="settings.widget.parameters.modalType")
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
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldNumber from './fields/common/number.vue';
import AlarmsListModalForm from './forms/alarms-list-modal.vue';
import MarginsForm from './forms/margins.vue';

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
    FieldDefaultElementsPerPage,
    FieldNumber,
    AlarmsListModalForm,
    MarginsForm,
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

