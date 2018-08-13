<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title")
      field-weather-data-set
      field-weather-template
      field-weather-modal-template
      field-weather-entity-template
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldWeatherDataSet from '@/components/other/settings/fields/weather-data-set.vue';
import FieldWeatherTemplate from '@/components/other/settings/fields/weather-template.vue';
import FieldWeatherModalTemplate from '@/components/other/settings/fields/weather-modal-template.vue';
import FieldWeatherEntityTemplate from '@/components/other/settings/fields/weather-entity-template.vue';
import widgetSettingsMixin from '@/mixins/widget/settings';

export default {
  components: {
    FieldTitle,
    FieldWeatherDataSet,
    FieldWeatherTemplate,
    FieldWeatherModalTemplate,
    FieldWeatherEntityTemplate,
  },
  mixins: [widgetSettingsMixin],
  data() {
    return {
      settings: {
        widget: {
          title: this.widget.title,
        },
      },
    };
  },
};
</script>

