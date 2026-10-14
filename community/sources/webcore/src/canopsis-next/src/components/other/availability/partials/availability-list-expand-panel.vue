<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab value="graph">
        {{ $t('common.graph') }}
      </v-tab>
      <v-tab value="active">
        {{ $t('context.activeAlarm') }}
      </v-tab>
      <v-tab value="resolved">
        {{ $t('context.resolvedAlarms') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="pa-3 bg-secondary-lighten-2">
      <c-window-expand-item value="graph">
        <availability-graph
          :availability="availability"
          :interval="interval"
          :default-show-type="defaultShowType"
          :display-parameter="displayParameter"
        />
      </c-window-expand-item>

      <c-window-expand-item value="active">
        <entity-alarms-list-table
          :entity="availability.entity"
          :columns="activeAlarmsColumns"
        />
      </c-window-expand-item>

      <c-window-expand-item value="resolved">
        <entity-alarms-list-table
          :entity="availability.entity"
          :columns="resolvedAlarmsColumns"
          resolved
        />
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { ref } from 'vue';

import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

import AvailabilityGraph from '@/components/other/availability/partials/availability-graph.vue';
import EntityAlarmsListTable from '@/components/widgets/context/partials/expand-panel-tabs/entity-alarms-list-table.vue';

export default {
  components: {
    AvailabilityGraph,
    EntityAlarmsListTable,
  },

  props: {
    availability: {
      type: Object,
      required: true,
    },
    activeAlarmsColumns: {
      type: Array,
      required: true,
    },
    resolvedAlarmsColumns: {
      type: Array,
      required: true,
    },
    interval: {
      type: Object,
      required: true,
    },
    defaultShowType: {
      type: Number,
      required: false,
    },
    displayParameter: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
  },

  setup() {
    const activeTab = ref('graph');

    return {
      activeTab,
    };
  },
};
</script>
