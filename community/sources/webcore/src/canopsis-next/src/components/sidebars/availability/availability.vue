<template>
  <widget-settings
    :submitting="submitting"
    divider
    @submit="submit"
  >
    <availability-form v-model="form" />
  </widget-settings>
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import WidgetSettings from '../partials/widget-settings.vue';

import AvailabilityForm from './form/availability-form.vue';

export default {
  name: SIDE_BARS.availabilitySettings,
  components: {
    WidgetSettings,
    AvailabilityForm,
  },
  mixins: [
    widgetSettingsMixin,
  ],
};
</script>
