<template lang="pug">
  alarms-list-table(
    :widget="widget",
    :alarms="alarmItems",
    :columns="columns",
    expandable,
    hasColumns
  )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ENTITIES_TYPES, WIDGET_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

import AlarmsListTable from '@/components/other/alarm/partials/alarms-list-table.vue';
import { generateWidgetByType } from '@/helpers/entities';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

export default {
  components: { AlarmsListTable },
  mixins: [
    authMixin,
    entitiesAlarmMixin,
    entitiesViewGroupMixin,
  ],
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  computed: {
    ...entitiesMapGetters({ getEntityItem: 'getItem' }),

    widgetId() {
      return this.$route.query.widgetId;
    },

    alarmItems() {
      return this.getAlarmsList([this.id]);
    },

    widget() {
      const widget = this.getEntityItem(ENTITIES_TYPES.widget, this.widgetId);

      return !widget || widget.type !== WIDGET_TYPES.alarmList
        ? generateWidgetByType(WIDGET_TYPES.alarmList)
        : widget;
    },

    columns() {
      return this.widget.parameters.widgetColumns.map(({ label: text, value }) => ({
        sortable: false,
        value,
        text,
      }));
    },
  },

  mounted() {
    this.fetchAlarmItem({
      id: this.id,
      params: {
        opened: true,
        resolved: true,
      },
    });

    if (this.widgetId && !this.groupsPending) {
      this.fetchGroupsList();
    }
  },
};
</script>
