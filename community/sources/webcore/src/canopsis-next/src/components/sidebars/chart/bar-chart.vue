<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <bar-chart-widget-form
      v-model="form"
      :with-filters="hasAccessToFilter"
      with-periodic-refresh
      with-preset
    />
  </widget-settings>
</template>

<script>
import { toRef } from 'vue';

import { SIDE_BARS } from '@/constants';

import { useBarChartFiltersPermissions } from '@/hooks/permissions/widgets/chart/filters';
import { useWidgetSettings } from '@/hooks/widget/settings';

import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import BarChartWidgetForm from './form/bar-chart-widget-form.vue';

export default {
  name: SIDE_BARS.barChartSettings,

  components: {
    BarChartWidgetForm,
    WidgetSettings,
  },

  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const sidebar = toRef(props, 'sidebar');
    const { form, hasChanges, submitting, close, submit } = useWidgetSettings({ sidebar });
    const { hasAccessToFilter } = useBarChartFiltersPermissions();

    return {
      form,
      hasChanges,
      submitting,
      hasAccessToFilter,

      close,
      submit,
    };
  },
};
</script>
