<template>
  <c-advanced-data-table
    :headers="headers"
    :items="pbehaviorExceptions"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledException"
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
        pbehavior-exception
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #visible="{ item }">
      <c-enabled :value="!item.hidden" />
    </template>
    <template #actions="{ item: actionsItem }">
      <c-action-btn
        v-if="updatable"
        type="edit"
        @click="$emit('edit', actionsItem)"
      />
      <c-action-btn
        v-if="removable"
        :tooltip="actionsItem.deletable ? $t('common.delete') : $t('pbehavior.exceptions.usingException')"
        :disabled="!actionsItem.deletable"
        type="delete"
        @click="$emit('remove', actionsItem._id)"
      />
    </template>
    <template #expand="{ item: expandItem }">
      <pbehavior-exceptions-list-expand-panel :pbehavior-exception="expandItem" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import PbehaviorExceptionsListExpandPanel from './partials/pbehavior-exceptions-list-expand-panel.vue';

export default {
  components: {
    PbehaviorExceptionsListExpandPanel,
  },
  props: {
    pbehaviorExceptions: {
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

    const isDisabledException = ({ deletable }) => !deletable;

    return {
      headers,

      isDisabledException,
    };
  },
};
</script>
