<template lang="pug">
  div
    v-list.pt-0(expand)
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
import FieldTitle from '@/components/side-bars/settings/partial/fields/title.vue';
import FieldWeatherDataSet from '@/components/side-bars/settings/partial/fields/weather-data-set.vue';
import FieldWeatherTemplate from '@/components/side-bars/settings/partial/fields/weather-template.vue';
import FieldGridSize from '@/components/side-bars/settings/partial/fields/grid-size.vue';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';

export default {
  name: SIDE_BARS.weatherSettings,
  components: {
    FieldTitle,
    FieldWeatherDataSet,
    FieldWeatherTemplate,
    FieldGridSize,
  },
  mixins: [widgetSettingsMixin],
  data() {
    return {
      settings: {
        widget: {
          title: '',
          block_template: '',
          modal_template: '',
          entity_template: '',
          columnSM: 0,
          columnMD: 0,
          columnLG: 0,
          columnHG: 0,
        },
      },
    };
  },
  mounted() {
    this.settings.widget = { ...this.widget };
  },
};
</script>

