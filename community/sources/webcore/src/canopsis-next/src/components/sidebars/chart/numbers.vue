<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <numbers-widget-form
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

import { useNumbersChartFiltersPermissions } from '@/hooks/permissions/widgets/chart/filters';
import { useWidgetSettings } from '@/hooks/widget/settings';

import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import NumbersWidgetForm from './form/numbers-widget-form.vue';

export default {
  name: SIDE_BARS.numbersSettings,

  components: {
    NumbersWidgetForm,
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
    const { hasAccessToFilter } = useNumbersChartFiltersPermissions();

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
