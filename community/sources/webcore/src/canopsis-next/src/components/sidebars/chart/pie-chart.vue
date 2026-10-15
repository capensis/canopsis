<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <pie-chart-widget-form
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

import { usePieChartFiltersPermissions } from '@/hooks/permissions/widgets/chart/filters';
import { useWidgetSettings } from '@/hooks/widget/settings';

import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import PieChartWidgetForm from './form/pie-chart-widget-form.vue';

export default {
  name: SIDE_BARS.pieChartSettings,

  components: {
    PieChartWidgetForm,
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
    const { hasAccessToFilter } = usePieChartFiltersPermissions();

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
