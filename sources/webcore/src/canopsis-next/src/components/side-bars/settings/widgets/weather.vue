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
      field-filter-editor(v-model="settings.widget.parameters.mfilter", :hidden-fields="['title']")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.titles.alarmListSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-columns(v-model="settings.widget.parameters.alarmsList.widgetColumns")
          v-divider
          field-default-elements-per-page(v-model="settings.widget.parameters.alarmsList.itemsPerPage")
          v-divider
          field-info-popup(v-model="settings.widget.parameters.alarmsList.infoPopups")
          v-divider
          field-more-info(v-model="settings.widget.parameters.alarmsList.moreInfoTemplate")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-weather-template(
          v-model="settings.widget.parameters.blockTemplate",
          :title="$t('settings.weatherTemplate')"
          )
          v-divider
          field-weather-template(
          v-model="settings.widget.parameters.modalTemplate",
          :title="$t('settings.modalTemplate')"
          )
          v-divider
          field-weather-template(
          v-model="settings.widget.parameters.entityTemplate",
          :title="$t('settings.entityTemplate')"
          )
          v-divider
          field-grid-size(v-model="settings.widget.parameters.columnSM", :title="$t('settings.columnSM')")
          v-divider
          field-grid-size(v-model="settings.widget.parameters.columnMD", :title="$t('settings.columnMD')")
          v-divider
          field-grid-size(v-model="settings.widget.parameters.columnLG", :title="$t('settings.columnLG')")
          v-divider
          v-list-group
            v-list-tile(slot="activator") {{ $t('settings.margin.title') }}
            v-list.grey.lighten-4.px-2.py-0(expand)
              field-slider(
              v-model="settings.widget.parameters.margin.top",
              :title="$t('settings.margin.top')",
              :min="0",
              :max="5",
              )
              v-divider
              field-slider(
              v-model="settings.widget.parameters.margin.right",
              :title="$t('settings.margin.right')",
              :min="0",
              :max="5",
              )
              v-divider
              field-slider(
              v-model="settings.widget.parameters.margin.bottom",
              :title="$t('settings.margin.bottom')",
              :min="0",
              :max="5",
              )
              v-divider
              field-slider(
              v-model="settings.widget.parameters.margin.left",
              :title="$t('settings.margin.left')",
              :min="0",
              :max="5",
              )
          v-divider
          field-slider(
          v-model="settings.widget.parameters.heightFactor",
          :title="$t('settings.height')",
          :min="1",
          :max="20",
          )
          v-divider
          field-modal-type(v-model="settings.widget.parameters.modalType")
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';
import sideBarSettingsWidgetAlarmMixin from '@/mixins/side-bar/settings/widgets/alarm';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldWeatherTemplate from './fields/weather/weather-template.vue';
import FieldGridSize from './fields/common/grid-size.vue';
import FieldSlider from './fields/common/slider.vue';
import FieldModalType from './fields/weather/modal-type.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldInfoPopup from './fields/alarm/info-popup.vue';
import FieldMoreInfo from './fields/alarm/more-info.vue';

export default {
  name: SIDE_BARS.weatherSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldFilterEditor,
    FieldWeatherTemplate,
    FieldGridSize,
    FieldSlider,
    FieldModalType,
    FieldColumns,
    FieldDefaultElementsPerPage,
    FieldInfoPopup,
    FieldMoreInfo,
  },
  mixins: [widgetSettingsMixin, sideBarSettingsWidgetAlarmMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: this.prepareWidgetWithAlarmParametersSettings(cloneDeep(widget), true),
      },
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

