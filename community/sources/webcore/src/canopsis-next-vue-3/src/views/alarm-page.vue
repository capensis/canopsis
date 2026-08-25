<template>
  <alarms-list-table
    :widget="widget"
    :alarms="alarmItems"
    :columns="columns"
    :meta="meta"
    :query="query"
    :loading="pending"
    expandable
    hide-pagination
    eager
  />
</template>

<script>
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { WIDGET_TYPES } from '@/constants';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';
import { prepareAlarmListWidget } from '@/helpers/entities/widget/forms/alarm';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';
import { useAlarm } from '@/hooks/store/modules/alarm';
import { useWidget } from '@/hooks/store/modules/widget';

import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';

export default {
  components: {
    AlarmsListTable,
  },

  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },

  setup(props) {
    const route = useRoute();
    const widget = ref(generatePreparedDefaultAlarmListWidget());
    const query = ref({
      page: 1,
      itemsPerPage: 1,
      sortBy: [],
      sortDesc: [],
    });

    const { t } = useI18n();
    const popups = usePopups();
    const { getItem: getAlarmItem, fetchItem: fetchAlarmItem } = useAlarm();
    const { fetchWidgetWithoutStore } = useWidget();

    const widgetId = computed(() => route.query.widgetId);

    const alarmItems = computed(() => {
      const alarm = getAlarmItem.value(props.id);

      return alarm ? [alarm] : [];
    });

    const meta = computed(() => ({
      total_count: alarmItems.value.length,
    }));

    const columns = computed(() =>
      widget.value.parameters.widgetColumns.map(column => ({
        ...column,
        sortable: false,
      })),
    );

    const fetchAlarmAndWidgetHandler = async () => {
      try {
        const requests = [fetchAlarmItem({ id: props.id })];

        if (widgetId.value) {
          requests.push(fetchWidgetWithoutStore({ id: widgetId.value }));
        }

        const [, widgetResponse] = await Promise.all(requests);

        if (widgetResponse?.type === WIDGET_TYPES.alarmList) {
          widget.value = prepareAlarmListWidget(widgetResponse);
        } else {
          widget.value = generatePreparedDefaultAlarmListWidget();
        }
      } catch (err) {
        console.error(err);

        popups.error({ text: err.description || t('errors.default') });
      }
    };

    const { pending, handler: fetchAlarmAndWidget } = usePendingHandler(
      fetchAlarmAndWidgetHandler,
    );

    onMounted(fetchAlarmAndWidget);
    watch([() => props.id, widgetId], fetchAlarmAndWidget);

    return {
      alarmItems,
      columns,
      meta,
      pending,
      query,
      widget,
    };
  },
};
</script>
