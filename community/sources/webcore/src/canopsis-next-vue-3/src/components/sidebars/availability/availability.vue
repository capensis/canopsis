<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <availability-form
      v-model="form"
      :widget-id="widget._id"
      :entity-columns-widget-templates="entityColumnsWidgetTemplates"
      :entity-columns-widget-templates-pending="widgetTemplatesPending"
      :alarm-columns-widget-templates="alarmColumnsWidgetTemplates"
      :alarm-columns-widget-templates-pending="widgetTemplatesPending"
      :show-interval="hasAccessToInterval"
      :show-export="hasAccessToExportAsCsv"
      :show-filter="hasAccessToFilter"
      @update:widget-columns-template="updateWidgetColumnsTemplate"
      @update:active-alarms-columns-template="updateActiveAlarmsColumnsTemplate"
      @update:resolved-alarms-columns-template="updateResolvedAlarmsColumnsTemplate"
    />
  </widget-settings>
</template>

<script>
import { computed, onMounted, ref, toRef } from 'vue';

import { MAX_LIMIT, SIDE_BARS, WIDGET_TEMPLATES_TYPES } from '@/constants';

import { prepareWidgetTemplates } from '@/helpers/entities/widget/template/prepare';

import {
  useAvailabilityExportPermissions,
  useAvailabilityFiltersPermissions,
  useAvailabilityIntervalPermissions,
} from '@/hooks/permissions/widgets/availability/filters';
import { useWidgetTemplate } from '@/hooks/store/modules/widget-template';
import { useWidgetSettings } from '@/hooks/widget/settings';

import AvailabilityForm from '@/components/sidebars/availability/form/availability-form.vue';
import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

export default {
  name: SIDE_BARS.availabilitySettings,

  components: {
    WidgetSettings,
    AvailabilityForm,
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
    const { hasAccessToFilter } = useAvailabilityFiltersPermissions();
    const { hasAccessToInterval } = useAvailabilityIntervalPermissions();
    const { hasAccessToExportAsCsv } = useAvailabilityExportPermissions();
    const { fetchWidgetTemplatesListWithoutStore } = useWidgetTemplate();

    const widgetTemplates = ref([]);
    const widgetTemplatesPending = ref(false);

    const preparedWidgetTemplates = computed(() => prepareWidgetTemplates(widgetTemplates.value));

    const alarmColumnsWidgetTemplates = computed(() =>
      preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.alarmColumns),
    );
    const entityColumnsWidgetTemplates = computed(() =>
      preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.entityColumns),
    );

    const updateWidgetColumnsTemplate = (template, columns) => {
      form.value.parameters.widget_columns_template = template;
      form.value.parameters.widget_columns = columns;
    };

    const updateActiveAlarmsColumnsTemplate = (template, columns) => {
      form.value.parameters.active_alarms_columns_template = template;
      form.value.parameters.active_alarms_columns = columns;
    };

    const updateResolvedAlarmsColumnsTemplate = (template, columns) => {
      form.value.parameters.resolved_alarms_columns_template = template;
      form.value.parameters.resolved_alarms_columns = columns;
    };

    const fetchWidgetTemplates = async () => {
      try {
        widgetTemplatesPending.value = true;

        const { data } = await fetchWidgetTemplatesListWithoutStore({ params: { limit: MAX_LIMIT } });

        widgetTemplates.value = data;
      } catch (err) {
        console.warn(err);
      } finally {
        widgetTemplatesPending.value = false;
      }
    };

    onMounted(() => {
      fetchWidgetTemplates();
    });

    return {
      form,
      widget,
      hasChanges,
      submitting,
      hasAccessToFilter,
      hasAccessToInterval,
      hasAccessToExportAsCsv,
      widgetTemplatesPending,
      alarmColumnsWidgetTemplates,
      entityColumnsWidgetTemplates,

      close,
      submit,
      updateWidgetColumnsTemplate,
      updateActiveAlarmsColumnsTemplate,
      updateResolvedAlarmsColumnsTemplate,
    };
  },
};
</script>
