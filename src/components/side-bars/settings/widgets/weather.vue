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
      field-filter-editor(v-model="settings.widget.parameters.mfilter")
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
    v-btn(@click="submit", color="green darken-4 white--text") {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldWeatherTemplate from './fields/weather/weather-template.vue';
import FieldGridSize from './fields/common/grid-size.vue';

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
  },
  mixins: [widgetSettingsMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: cloneDeep(widget),
      },
    };
  },
};
</script>

