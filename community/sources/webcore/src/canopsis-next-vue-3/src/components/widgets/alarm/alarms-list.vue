<template>
  <div class="alarms-list-widget pa-2">
    <v-data-table
      :headers="headers"
      :items="items"
      class="alarms-list-widget__table"
      density="compact"
      hide-default-footer
    >
      <template #no-data>
        <tr>
          <td :colspan="headersColspan" class="text-center text-medium-emphasis py-6">
            {{ $t('common.noData') }}
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script>
import { prepareAlarmListWidget } from '@/helpers/entities/widget/forms/alarm';

/**
 * Alarm-list widget runtime foundation.
 *
 * The dedicated table, alarm query state and actions are migrated in the next
 * AlarmsList stages; this component already consumes prepared widget
 * parameters so create/edit/duplicate flows use the final settings payload.
 */
export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tabId: {
      type: String,
      default: '',
    },
    visible: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    hideMassSelection: {
      type: Boolean,
      default: false,
    },
    hideToolbar: {
      type: Boolean,
      default: false,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
  },
  computed: {
    preparedWidget() {
      return prepareAlarmListWidget(this.widget);
    },

    columns() {
      return this.preparedWidget.parameters?.widgetColumns ?? [];
    },

    headers() {
      const headers = this.columns.map(column => ({
        ...column,
        key: column.value,
        title: column.text,
      }));

      if (!this.hideActions) {
        headers.push({
          key: 'actions',
          title: this.$t('common.actionsLabel'),
          sortable: false,
        });
      }

      return headers;
    },

    headersColspan() {
      return Math.max(this.headers.length, 1);
    },

    items() {
      return [];
    },
  },
};
</script>

<style lang="scss" scoped>
.alarms-list-widget {
  height: 100%;
}

.alarms-list-widget__table {
  min-height: 100%;
}
</style>
