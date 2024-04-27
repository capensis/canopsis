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
import { WIDGET_TYPES } from '@/constants';

import { prepareAlarmListWidget } from '@/helpers/entities/widget/forms/alarm';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { authMixin } from '@/mixins/auth';
import { entitiesAlarmMixin } from '@/mixins/entities/alarm';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';

export default {
  components: { AlarmsListTable },
  mixins: [
    authMixin,
    entitiesAlarmMixin,
    entitiesWidgetMixin,
  ],
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      widget: generatePreparedDefaultAlarmListWidget(),
    };
  },
  computed: {
    widgetId() {
      return this.$route.query.widgetId;
    },

    alarmItems() {
      const alarm = this.getAlarmItem(this.id);

      return alarm ? [alarm] : [];
    },

    columns() {
      return this.widget.parameters.widgetColumns.map(column => ({
        ...column,

        sortable: false,
      }));
    },
  },

  async mounted() {
    this.fetchAlarmAndWidget();
  },
  methods: {
    async fetchAlarmAndWidget() {
      try {
        this.pending = true;

        const requests = [this.fetchAlarmItem({ id: this.id })];

        if (this.widgetId) {
          requests.push(this.fetchWidgetWithoutStore({ id: this.widgetId }));
        }

        const [, widget] = await Promise.all(requests);

        if (widget?.type === WIDGET_TYPES.alarmList) {
          this.widget = prepareAlarmListWidget(widget);
        }
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.description || this.$t('errors.default') });
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
