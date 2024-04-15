<template>
  <v-tabs
    v-model="activeTab"
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $tc('common.graph') }}</v-tab>
    <v-tab>{{ $t('context.activeAlarm') }}</v-tab>
    <v-tab>{{ $t('context.resolvedAlarms') }}</v-tab>

    <v-tabs-items
      v-model="activeTab"
      mandatory
    >
      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card>
              <v-card-text>
                <availability-history
                  :availability="availability"
                  :interval="interval"
                  :default-show-type="defaultShowType"
                  :display-parameter="displayParameter"
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>

      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card>
              <v-card-text>
                <entity-alarms-list-table
                  :entity="availability.entity"
                  :columns="activeAlarmsColumns"
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>

      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card>
              <v-card-text>
                <entity-alarms-list-table
                  :entity="availability.entity"
                  :columns="resolvedAlarmsColumns"
                  resolved
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import { ref } from 'vue';

import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

import EntityAlarmsListTable from '@/components/widgets/context/partials/expand-panel-tabs/entity-alarms-list-table.vue';
import AvailabilityHistory from '@/components/other/availability/partials/availability-graph.vue';

export default {
  components: {
    AvailabilityHistory,
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
    const activeTab = ref(0);

    return {
      activeTab,
    };
  },
};
</script>
