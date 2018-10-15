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
      field-title(v-model="settings.widget.title")
      v-divider
      field-weather-data-set
      v-divider
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
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldRowGridSize from '../partial/fields/row-grid-size.vue';
import FieldTitle from '../partial/fields/title.vue';
import FieldWeatherDataSet from '../partial/fields/weather-data-set.vue';
import FieldWeatherTemplate from '../partial/fields/weather-template.vue';
import FieldGridSize from '../partial/fields/grid-size.vue';

export default {
  name: SIDE_BARS.weatherSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldWeatherDataSet,
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

