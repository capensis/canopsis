<template>
  <c-advanced-data-table
    :headers="headers"
    :items="maps"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable && hasAnyDeletableMap"
    :is-disabled-item="isDisabledMap"
    advanced-pagination
    expand
    search
    @update:options="$emit('update:options', $event)"
  >
    <template v-if="removable && hasAnyDeletableMap" #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        map
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #type="{ item }">
      <span>{{ $t(`map.types.${item.type}`) }}</span>
    </template>
    <template #updated="{ item }">
      {{ item.updated | date }}
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
          :tooltip="item.deletable ? $t('common.delete') : $t('map.usingMap')"
          :disabled="!item.deletable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <maps-list-expand-item :map="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import MapsListExpandItem from './partials/maps-list-expand-item.vue';

export default {
  components: {
    MapsListExpandItem,
  },
  props: {
    maps: {
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
  },
  setup(props) {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'name',
      },
      {
        text: t('common.type'),
        value: 'type',
      },
      {
        text: t('common.lastModifiedOn'),
        value: 'updated',
      },
      {
        text: t('common.lastModifiedBy'),
        value: 'author.display_name',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const hasAnyDeletableMap = computed(() => props.maps.some(({ deletable }) => deletable));

    const isDisabledMap = ({ deletable }) => !deletable;

    return {
      headers,
      hasAnyDeletableMap,

      isDisabledMap,
    };
  },
};
</script>
