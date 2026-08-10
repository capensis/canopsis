<template>
  <div class="secondary pa-3">
    <v-card>
      <v-card-text>
        <v-data-table :headers="preparedHeaders" :items="items" hide-default-footer>
          <template #item.failed="{ item }">
            {{ formatDate(item.failed) }}
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    token: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.requestDate'), value: 'failed', width: 200, sortable: false },
      { text: t('common.output'), value: 'fail_reason', sortable: false },
    ]);
    const preparedHeaders = computed(() =>
      headers.value.map(header => ({
        ...header,
        title: header.text,
        key: header.value,
      })),
    );
    const items = computed(() => [props.token]);
    const formatDate = date => convertDateToString(date);

    return {
      preparedHeaders,
      items,
      formatDate,
    };
  },
};
</script>
