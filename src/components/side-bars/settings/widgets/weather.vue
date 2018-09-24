<template lang="pug">
  div
    v-list.pt-0(expand)
      field-row-grid-size(
      :rowId.sync="settings.rowId",
      :size.sync="settings.widget.size",
      :availableRows="availableRows",
      @createRow="createRow"
      )
      field-title(v-model="settings.widget.title")
      field-weather-data-set
      field-weather-template(v-model="settings.widget.block_template", :title="$t('settings.weatherTemplate')")
      field-weather-template(v-model="settings.widget.modal_template", :title="$t('settings.modalTemplate')")
      field-weather-template(v-model="settings.widget.entity_template", :title="$t('settings.entityTemplate')")
      field-grid-size(v-model="settings.widget.columnSM", :title="$t('settings.columnSM')")
      field-grid-size(v-model="settings.widget.columnMD", :title="$t('settings.columnMD')")
      field-grid-size(v-model="settings.widget.columnLG", :title="$t('settings.columnLG')")
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
  mounted() {
    this.settings.widget = { ...this.widget };
  },
};
</script>

