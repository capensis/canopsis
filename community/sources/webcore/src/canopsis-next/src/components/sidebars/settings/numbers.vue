<template lang="pug">
  widget-settings(:submitting="submitting", divider, @submit="submit")
    numbers-widget-form(
      v-model="form",
      with-periodic-refresh,
      with-filters
    )
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import WidgetSettings from './partials/widget-settings.vue';
import NumbersWidgetForm from './forms/widgets/numbers-widget-form.vue';

export default {
  name: SIDE_BARS.numbersSettings,
  components: {
    WidgetSettings,
    NumbersWidgetForm,
  },
  mixins: [widgetSettingsMixin],
};
</script>
