<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <external-data-table-form
      v-model="form"
      :widget-id="widget._id"
    />
  </widget-settings>
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import WidgetSettings from '../partials/widget-settings.vue';

import ExternalDataTableForm from './form/external-data-table-form.vue';

export default {
  name: SIDE_BARS.externalDataTableSettings,
  components: {
    WidgetSettings,
    ExternalDataTableForm,
  },
  mixins: [
    widgetSettingsMixin,
  ],
};
</script>
