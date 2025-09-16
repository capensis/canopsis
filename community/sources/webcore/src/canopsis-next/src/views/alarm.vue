<template>
  <alarms-list-table
    :widget="widget"
    :alarms="alarmItems"
    :columns="columns"
    :loading="pending"
    :total-items="alarmItems.length"
    expandable
    hide-pagination
    has-columns
  />
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router/composables';

import { WIDGET_TYPES } from '@/constants';

import { prepareAlarmListWidget } from '@/helpers/entities/widget/forms/alarm';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { useAlarm } from '@/hooks/store/modules/alarm';
import { useWidget } from '@/hooks/store/modules/widget';
import { usePopups } from '@/hooks/popups';
import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';

import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';

export default {
  components: { AlarmsListTable },
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  setup(props) {
    const widget = ref(generatePreparedDefaultAlarmListWidget());

    const { t } = useI18n();
    const popups = usePopups();
    const { getAlarmItem, fetchAlarmItem } = useAlarm();
    const { fetchWidgetWithoutStore } = useWidget();
    const route = useRoute();

    const widgetId = computed(() => route.query.widgetId);

    const alarmItems = computed(() => {
      const alarm = getAlarmItem.value(props.id);

      return alarm ? [alarm] : [];
    });

    const columns = computed(() => widget.value.parameters.widgetColumns.map(column => ({
      ...column,
      sortable: false,
    })));

    /**
     * Fetches alarm data and optional widget configuration for the alarm view.
     */
    const fetchAlarmAndWidgetHandler = async () => {
      try {
        const requests = [fetchAlarmItem({ id: props.id })];

        if (widgetId.value) {
          requests.push(fetchWidgetWithoutStore({ id: widgetId.value }));
        }

        const [, widgetResponse] = await Promise.all(requests);

        if (widgetResponse?.type === WIDGET_TYPES.alarmList) {
          widget.value = prepareAlarmListWidget(widgetResponse);
        }
      } catch (err) {
        console.error(err);

        popups.error({ text: err.description || t('errors.default') });
      }
    };

    const { pending, handler: fetchAlarmAndWidget } = usePendingHandler(fetchAlarmAndWidgetHandler);

    onMounted(fetchAlarmAndWidget);

    return {
      pending,
      widget,
      alarmItems,
      columns,
    };
  },
};
</script>
