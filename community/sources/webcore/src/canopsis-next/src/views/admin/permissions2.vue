<template>
  <v-container class="admin-rights">
    <c-page-header />
    <v-card class="position-relative">
      <v-tabs fixed-tabs>
        <v-tab>BUSINESS</v-tab>
        <v-tab-item>
          <permissions-table
            v-model="items"
            :headers="headers"
          />
        </v-tab-item>
      </v-tabs>
    </v-card>
  </v-container>
</template>

<script>
import { computed, ref } from 'vue';

import PermissionsTable from '@/components/other/permission2/permissions-table.vue';

export default {
  components: { PermissionsTable },
  setup() {
    const items = ref([
      {
        text: 'Common',
        value: { manager: true, pilot: true, admin: true },
      },
      {
        text: 'Alarms list',
        value: { manager: true, pilot: true, admin: true },
        children: [
          {
            text: 'Alarms list widget settings',
            value: { manager: true, pilot: true, admin: true },
            children: [
              { text: 'Set alarm filters', value: { manager: true, pilot: true, admin: true } },
              { text: 'Set filters by remediation instructions', value: { manager: true, pilot: true, admin: true } },
            ],
          },
        ],
      },
    ]);

    const headers = computed(() => [
      { text: '', sortable: false },
      { text: 'Manager', sortable: false, value: 'manager' },
      { text: 'Pilot', sortable: false, value: 'pilot' },
      { text: 'Admin', sortable: false, value: 'admin' },
    ]);

    return {
      headers,
      items,
    };
  },
};
</script>
