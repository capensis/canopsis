<template>
  <span class="service-dependencies-expand">
    <c-tooltip v-if="item.loadMore" :content="$t('common.loadMore')" right>
      <template #activator="{ props: tooltipProps }">
        <v-btn
          v-bind="tooltipProps"
          :loading="pending"
          icon="more_horiz"
          variant="text"
          size="small"
          @click="$emit('load', item)"
        />
      </template>
    </c-tooltip>
    <v-badge
      v-else
      :model-value="hasStatusIcon"
      color="grey"
      location="top end"
    >
      <v-btn
        :class="entityColorClass"
        variant="flat"
        size="small"
        icon
        @click="$emit('show', item)"
      >
        <v-icon>{{ entityIcon }}</v-icon>
      </v-btn>
      <template #badge>
        <v-icon class="pa-0" size="12">
          {{ statusIcon }}
        </v-icon>
      </template>
    </v-badge>
    <c-tooltip v-if="item.cycle" :content="$t('common.cycleDependency')" top>
      <template #activator="{ props: tooltipProps }">
        <v-icon
          v-bind="tooltipProps"
          class="ml-1"
          color="error"
          size="14"
        >
          autorenew
        </v-icon>
      </template>
    </c-tooltip>
  </span>
</template>

<script>
import {
  COLOR_INDICATOR_TYPES,
  ENTITY_TYPES_ICONS,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import { getEntityColorClass } from '@/helpers/entities/entity/color';

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

  emits: ['load', 'show'],

  computed: {
    entity() {
      return this.item?.entity ?? {};
    },

    entityColorClass() {
      return getEntityColorClass(this.entity, COLOR_INDICATOR_TYPES.impactState);
    },

    entityIcon() {
      return ENTITY_TYPES_ICONS[this.entity.type] ?? 'perm_identity';
    },

    hasStatusIcon() {
      return Boolean(this.entity.pause || this.entity.pbehavior_info?.id);
    },

    statusIcon() {
      const pbehavior = this.entity.pbehavior_info;
      const hasNonPausePbehavior = pbehavior?.id &&
        pbehavior.canonical_type !== PBEHAVIOR_TYPE_TYPES.pause;

      if (hasNonPausePbehavior) {
        return pbehavior.icon_name ?? 'build';
      }

      if (this.entity.pause) {
        return 'pause';
      }

      return '$vuetify.icons.pbehavior';
    },
  },
};
</script>

<style lang="scss" scoped>
.service-dependencies-expand {
  display: inline-flex;
  align-items: center;
  flex-wrap: nowrap;
}
</style>
