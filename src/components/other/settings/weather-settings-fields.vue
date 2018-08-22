<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title")
      field-weather-data-set
      field-weather-template(:content.sync="settings.widget.block_template", :value="widget.block_template")
      field-weather-modal-template(:content.sync="settings.widget.modal_template", :value="widget.modal_template")
      field-weather-entity-template(:content.sync="settings.widget.entity_template", :value="widget.entity_template")
      field-grid-size(:value.sync="widget.columnSM")
      field-grid-size(:value.sync="widget.columnMD")
      field-grid-size(:value.sync="widget.columnLG")
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldWeatherDataSet from '@/components/other/settings/fields/weather-data-set.vue';
import FieldWeatherTemplate from '@/components/other/settings/fields/weather-template.vue';
import FieldWeatherModalTemplate from '@/components/other/settings/fields/weather-modal-template.vue';
import FieldWeatherEntityTemplate from '@/components/other/settings/fields/weather-entity-template.vue';
import FieldGridSize from '@/components/other/settings/fields/grid-size.vue';
import widgetSettingsMixin from '@/mixins/widget/settings';

export default {
  components: {
    FieldTitle,
    FieldWeatherDataSet,
    FieldWeatherTemplate,
    FieldWeatherModalTemplate,
    FieldWeatherEntityTemplate,
    FieldGridSize,
  },
  mixins: [widgetSettingsMixin],
  data() {
    return {
      settings: {
        widget: {
          title: this.widget.title,
          block_template: '',
          modal_template: '',
          entity_template: '',
          columnSM: 0,
          columnMD: 2,
          columnLG: 4,
        },
      },
    };
  },
  mounted() {
    this.settings.widget = { ...this.widget };
  },
};
</script>

