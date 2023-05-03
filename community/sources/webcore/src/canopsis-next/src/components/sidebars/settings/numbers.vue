<template lang="pug">
  widget-settings(:submitting="submitting", divider, @submit="submit")
    numbers-widget-form(
      v-model="form",
      :with-filters="hasAccessToListFilters",
      with-periodic-refresh,
      with-preset
    )
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { permissionsWidgetsNumbersFilters } from '@/mixins/permissions/widgets/chart/numbers/filters';

import WidgetSettings from './partials/widget-settings.vue';
import NumbersWidgetForm from './forms/widgets/numbers-widget-form.vue';

export default {
  name: SIDE_BARS.numbersSettings,
  components: {
    WidgetSettings,
    NumbersWidgetForm,
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsNumbersFilters,
  ],
};
</script>
