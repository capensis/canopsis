<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $tc('common.graph') }}</v-tab>
    <v-tab-item>
      <v-layout class="pa-3">
        <v-flex xs12>
          <v-card>
            <v-card-text>
              GRAPH
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>

    <v-tab>{{ $t('context.activeAlarm') }}</v-tab>
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

    <v-tab>{{ $t('context.resolvedAlarms') }}</v-tab>
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
  </v-tabs>
</template>

<script>
import EntityAlarmsListTable from '@/components/widgets/context/partials/expand-panel-tabs/entity-alarms-list-table.vue';

export default {
  components: {
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
  },
};
</script>
