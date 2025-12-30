<template>
  <c-advanced-data-table
    :headers="headers"
    :items="pbehaviorReasons"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledReason"
    :select-all="removable || updatable"
    expand
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        :hideable="updatable"
        :unhideable="updatable"
        pbehavior-reason
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #visible="{ item }">
      <c-enabled :value="!item.hidden" />
    </template>
    <template #actions="{ item }">
      <c-action-btn
        v-if="updatable"
        type="edit"
        @click="$emit('edit', item)"
      />
      <c-action-btn
        v-if="removable"
        :tooltip="item.deletable ? $t('common.delete') : $t('pbehavior.reasons.usingReason')"
        :disabled="!item.deletable"
        type="delete"
        @click="$emit('remove', item._id)"
      />
    </template>
    <template #expand="{ item }">
      <pbehavior-reasons-list-expand-panel :pbehavior-reason="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import PbehaviorReasonsListExpandPanel from './partials/pbehavior-reasons-list-expand-panel.vue';

export default {
  components: {
    PbehaviorReasonsListExpandPanel,
  },
  props: {
    pbehaviorReasons: {
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
      required: true,
    },
    updatable: {
      type: Boolean,
      required: true,
    },
  },
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'name',
      },
      {
        text: t('pbehavior.visible'),
        value: 'visible',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const isDisabledReason = ({ deletable }) => !deletable;

    return {
      headers,

      isDisabledReason,
    };
  },
};
</script>
