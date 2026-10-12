<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <statistics-form
      v-model="form"
      :widget="widget"
      :type="$constants.KPI_RATING_SETTINGS_TYPES.user"
      :show-filter="hasAccessToFilter"
      :show-interval="hasAccessToInterval"
    />
  </widget-settings>
</template>

<script>
import { toRef } from 'vue';

import { SIDE_BARS } from '@/constants';

import {
  useUserStatisticsFiltersPermissions,
  useUserStatisticsIntervalPermissions,
} from '@/hooks/permissions/widgets/statistics/filters';
import { useWidgetSettings } from '@/hooks/widget/settings';

import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import StatisticsForm from './form/statistics-form.vue';

export default {
  name: SIDE_BARS.userStatisticsSettings,

  components: {
    WidgetSettings,
    StatisticsForm,
  },

  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const sidebar = toRef(props, 'sidebar');
    const { form, widget, hasChanges, submitting, close, submit } = useWidgetSettings({ sidebar });
    const { hasAccessToFilter } = useUserStatisticsFiltersPermissions();
    const { hasAccessToInterval } = useUserStatisticsIntervalPermissions();

    return {
      form,
      widget,
      hasChanges,
      submitting,
      hasAccessToFilter,
      hasAccessToInterval,

      close,
      submit,
    };
  },
};
</script>
