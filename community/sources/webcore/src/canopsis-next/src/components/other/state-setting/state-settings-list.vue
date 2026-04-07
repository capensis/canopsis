<template>
  <c-advanced-data-table
    :items="stateSettings"
    :headers="headers"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    select-all
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        state-setting
        @clear:items="clearSelected"
      />
    </template>
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>
    <template #type="{ item }">
      {{ item.type || '-' }}
    </template>
    <template #priority="{ item }">
      {{ item.priority || '-' }}
    </template>
    <template #method="{ item }">
      {{ getMethodLabel(item.method) }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          :disabled="!item.editable"
          type="edit"
          @click.stop="$emit('edit', item)"
        />
        <c-action-btn
          v-if="addable"
          :disabled="!isDuplicable(item)"
          type="duplicate"
          @click.stop="$emit('duplicate', item)"
        />
        <c-action-btn
          v-if="removable"
          :disabled="!item.deletable"
          type="delete"
          @click.stop="$emit('remove', item)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { JUNIT_STATE_SETTING_ID, SERVICE_STATE_SETTING_ID } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    stateSettings: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t, te } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.title'),
        value: 'title',
      },
      {
        text: t('common.enabled'),
        value: 'enabled',
      },
      {
        text: t('common.priority'),
        value: 'priority',
      },
      {
        text: t('stateSetting.appliedFor'),
        value: 'type',
      },
      {
        text: t('common.method'),
        value: 'method',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const isDuplicable = item => ![JUNIT_STATE_SETTING_ID, SERVICE_STATE_SETTING_ID].includes(item._id);

    const getMethodLabel = method => (te(`stateSetting.methods.${method}.label`)
      ? t(`stateSetting.methods.${method}.label`)
      : t(`stateSetting.junit.methods.${method}`));

    return {
      headers,
      isDuplicable,
      getMethodLabel,
    };
  },
};
</script>
