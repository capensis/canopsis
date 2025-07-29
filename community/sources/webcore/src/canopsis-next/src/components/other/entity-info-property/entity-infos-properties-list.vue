<template>
  <c-advanced-data-table
    :headers="headers"
    :items="entityInfosProperties"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        :tooltip="$t('common.bulkDelete')"
        icon="delete"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #type="{ item }">
      <span>{{ $t(ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS[item.type]) }}</span>
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="duplicable"
          type="duplicate"
          @click="$emit('duplicate', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    entityInfosProperties: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    duplicable: {
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
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('entityInfoProperties.infosKey'),
        value: 'name',
      },
      {
        text: t('common.description'),
        value: 'description',
      },
      {
        text: t('common.type'),
        value: 'type',
      },
      {
        text: t('common.alias'),
        value: 'alias',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    return {
      ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS,

      headers,
    };
  },
};
</script>
