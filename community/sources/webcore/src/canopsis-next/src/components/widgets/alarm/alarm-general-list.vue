<template>
  <div>
    <v-subheader>{{ $t('common.general') }}</v-subheader>
    <v-data-table
      :headers="headers"
      :items="items"
    >
      <template #item="{ item }">
        <tr>
          <td class="text-left">
            {{ item.v.state.a }}
          </td>
          <td class="text-left">
            {{ item.v.connector }}
          </td>
          <td class="text-left">
            {{ item.v.component }}
          </td>
          <td class="text-left">
            {{ item.v.resource }}
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.author'),
        sortable: false,
      },
      {
        text: t('common.connector'),
        sortable: false,
      },
      {
        text: t('common.component'),
        sortable: false,
      },
      {
        text: t('common.resource'),
        sortable: false,
      },
    ]);

    return {
      headers,
    };
  },
};
</script>
