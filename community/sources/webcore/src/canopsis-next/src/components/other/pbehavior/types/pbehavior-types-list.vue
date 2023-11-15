<template>
  <c-advanced-data-table
    :headers="headers"
    :items="pbehaviorTypes"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledType"
    :select-all="removable"
    expand
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #icon_name="{ item }">
      <v-chip
        class="pbehavior-type-icon"
        :color="item.color"
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
    <template #hidden="{ item }">
      <c-enabled :value="item.hidden" />
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
import { getMostReadableTextColor } from '@/helpers/color';

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
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.icon'),
          value: 'icon_name',
          sortable: false,
        },
        {
          text: this.$t('common.priority'),
          value: 'priority',
        },
        {
          text: this.$t('common.hidden'),
          value: 'hidden',
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
  methods: {
    isDisabledType({ deletable }) {
      return !deletable;
    },

    getIconColor(color) {
      return getMostReadableTextColor(color, { level: 'AA', size: 'large' });
    },
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-type-icon {
    width: 42px;
    height: 24px;
  }
</style>
