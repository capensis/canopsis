<template>
  <c-advanced-data-table
    :headers="headers"
    :items="items"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :toolbar-props="{ class: 'mb-3' }"
    search
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar>
      <c-quick-date-interval-field
        :interval="options.interval"
        with-hours
        @input="updateInterval"
      />
    </template>
    <template #time="{ item }">
      {{ dateToString(item.time) }}
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:options'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.date'),
        value: 'time',
      },
      {
        text: t('entity.infosLog.eventFilterId'),
        value: 'rule._id',
        sortable: false,
      },
      {
        text: t('entity.infosLog.eventFilterDescription'),
        value: 'rule.description',
        sortable: false,
      },
      {
        text: t('common.name'),
        value: 'name',
        sortable: false,
      },
      {
        text: t('entity.infosLog.oldValue'),
        value: 'prev_value',
        sortable: false,
      },
      {
        text: t('entity.infosLog.newValue'),
        value: 'new_value',
        sortable: false,
      },
    ]);

    const updateOptions = options => emit('update:options', options);
    const updateInterval = interval => updateOptions({ ...props.options, interval });
    const dateToString = date => convertDateToString(date);

    return {
      dateToString,
      headers,
      updateInterval,
      updateOptions,
    };
  },
};
</script>
