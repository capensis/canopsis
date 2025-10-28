<template>
  <c-advanced-data-table
    :headers="headers"
    :items="items"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    search
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar="">
      <v-layout align-center>
        <c-quick-date-interval-field
          :interval="options.interval"
          with-hours
          @input="updateInterval"
        />
      </v-layout>
    </template>
    <template #time="{ item }">
      {{ item.time | date }}
    </template>
  </c-advanced-data-table>
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

    /**
     * Updates the table options and emits the changes to the parent component
     *
     * @param {Object} options - The updated options object
     */
    const updateOptions = options => emit('update:options', options);

    /**
     * Updates the interval in the table options and emits the changes to the parent component
     *
     * @param {Object} interval - The new interval object to apply to the options
     */
    const updateInterval = interval => updateOptions({ ...props.options, interval });

    return {
      headers,
      updateOptions,
      updateInterval,
    };
  },
};
</script>
