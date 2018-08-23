<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title")
      field-weather-data-set
      field-weather-template(:value.sync="settings.widget.block_template", :title="$t('settings.weatherTemplate')")
      field-weather-template(:value.sync="settings.widget.modal_template", :title="$t('settings.modalTemplate')")
      field-weather-template(:value.sync="settings.widget.entity_template", :title="$t('settings.entityTemplate')")
      field-grid-size(v-model="settings.widget.columnSM", :title="$t('settings.columnSM')")
      field-grid-size(:value.sync="settings.widget.columnMD", :title="$t('settings.columnMD')")
      field-grid-size(:value.sync="settings.widget.columnLG", :title="$t('settings.columnLG')")
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldWeatherDataSet from '@/components/other/settings/fields/weather-data-set.vue';
import FieldWeatherTemplate from '@/components/other/settings/fields/weather-template.vue';
import FieldGridSize from '@/components/other/settings/fields/grid-size.vue';
import widgetSettingsMixin from '@/mixins/widget/settings';

export default {
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
          title: this.widget.title || '',
          block_template: this.widget.block_template || '',
          modal_template: this.widget.modal_template || '',
          entity_template: this.widget.entity_template || '',
          columnSM: this.widget.columnSM || 0,
          columnMD: this.widget.columnMD || 0,
          columnLG: this.widget.columnLG || 0,
          columnHG: this.widget.columnHG || 0,
        },
      },
    };
  },
  mounted() {
    this.settings.widget = { ...this.widget };
  },
};
</script>

