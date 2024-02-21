<template>
  <v-layout>
    <v-tooltip
      v-if="item.loadMore"
      right
    >
      <template #activator="{ on }">
        <v-btn
          v-on="on"
          :loading="pending"
          icon
          @click="$emit('load', item)"
        >
          <v-icon>more_horiz</v-icon>
        </v-btn>
      </template>
      <span>{{ $t('common.loadMore') }}</span>
    </v-tooltip>
    <v-badge
      v-else
      :value="!!item.entity?.pbehavior_info"
      color="grey"
      overlap
    >
      <v-btn
        :color="entityColor"
        fab
        small
        depressed
        dark
        @click="$emit('show', item)"
      >
        <v-icon>{{ entityIcon }}</v-icon>
      </v-btn>
      <template #badge>
        <v-icon
          class="pa-0"
          dark
        >
          {{ pbehaviorIcon }}
        </v-icon>
      </template>
    </v-badge>
    <v-tooltip
      v-if="item.cycle"
      top
    >
      <template #activator="{ on }">
        <v-icon
          class="ml-1"
          v-on="on"
          color="error"
          size="14"
        >
          autorenew
        </v-icon>
      </template>
      <span>{{ $t('common.cycleDependency') }}</span>
    </v-tooltip>
  </v-layout>
</template>

<script>
import { COLOR_INDICATOR_TYPES } from '@/constants';

import { getEntityColor } from '@/helpers/entities/entity/color';
import { getIconByEntityType } from '@/helpers/entities/entity/icons';

export default {
  props: {
    item: {
      type: Object,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    entity() {
      return this.item?.entity ?? {};
    },

    entityColor() {
      return getEntityColor(this.entity, COLOR_INDICATOR_TYPES.impactState);
    },

    entityIcon() {
      return getIconByEntityType(this.entity.type);
    },

    pbehaviorIcon() {
      return this.entity?.pbehavior_info?.icon_name ?? 'pause';
    },
  },
};
</script>
