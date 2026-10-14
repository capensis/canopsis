<template>
  <c-advanced-data-table
    :headers="headers"
    :items="remediationConfigurations"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledConfiguration"
    :select-all="removable"
    :hide-mass-actions="!active"
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        remediation-configuration
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>

    <template #actions="{ item, disabled }">
      <c-layout>
        <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn v-if="duplicable" type="duplicate" @click="$emit('duplicate', item)" />
        <c-db-export-btn :id="item._id" job-config />
        <c-action-btn
          v-if="removable"
          :tooltip="
            disabled ? t('remediation.configuration.usingConfiguration') : t('common.delete')
          "
          :disabled-button="disabled"
          type="delete"
          @click="$emit('remove', item)"
        />
      </c-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    remediationConfigurations: {
      type: Array,
      required: true,
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
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
    active: {
      type: Boolean,
      default: true,
    },
  },

  emits: ['update:options', 'refresh', 'remove', 'duplicate', 'edit'],

  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'name',
      },
      {
        text: t('common.author'),
        value: 'author.display_name',
      },
      {
        text: t('common.type'),
        value: 'type',
      },
      {
        text: t('remediation.configuration.host'),
        value: 'host',
      },
      {
        text: t('remediation.configuration.clientId'),
        value: 'client_id',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const isDisabledConfiguration = ({ deletable }) => !deletable;

    return {
      t,
      headers,
      isDisabledConfiguration,
    };
  },
};
</script>
