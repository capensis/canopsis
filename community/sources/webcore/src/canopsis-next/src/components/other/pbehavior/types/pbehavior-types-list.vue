<template>
  <c-advanced-data-table
    :headers="headers"
    :items="pbehaviorTypes"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledType"
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
        pbehavior-type
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #icon_name="{ item }">
      <v-chip
        :color="item.color"
        class="pbehavior-type-icon"
      >
        <v-icon
          :color="getIconColor(item.color)"
          size="18"
        >
          {{ item.icon_name }}
        </v-icon>
      </v-chip>
    </template>
    <template #priority="{ item }">
      {{ item.priority || '-' }}
    </template>
    <template #visible="{ item }">
      <c-enabled :value="!item.hidden" />
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          :disabled="!item.deletable"
          :tooltip="item.deletable ? $t('common.delete') : $t('pbehavior.types.usingType')"
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <pbehavior-types-list-expand-panel :pbehavior-type="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { getMostReadableTextColor } from '@/helpers/color';

import { useI18n } from '@/hooks/i18n';

import PbehaviorTypesListExpandPanel from './partials/pbehavior-types-list-expand-panel.vue';

export default {
  components: {
    PbehaviorTypesListExpandPanel,
  },
  props: {
    pbehaviorTypes: {
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
    const { t, tc } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'name',
      },
      {
        text: tc('common.icon', 1),
        value: 'icon_name',
        sortable: false,
      },
      {
        text: t('common.priority'),
        value: 'priority',
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

    const isDisabledType = ({ deletable }) => !deletable;

    const getIconColor = color => getMostReadableTextColor(color, { level: 'AA', size: 'large' });

    return {
      headers,
      isDisabledType,
      getIconColor,
    };
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-type-icon {
    width: 42px;
    height: 24px;
  }
</style>
