<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <span
          class="c-extra-details__badge brown darken-1"
          v-on="on"
        >
          <v-icon
            color="white"
            small
          >
            {{ icon }}
          </v-icon>
        </span>
      </template>
      <div class="text-md-center">
        <strong>{{ $t('alarm.actions.iconsTitles.grouping') }}</strong>
        <div>{{ $t('common.title') }} : {{ ruleName }}</div>
        <div>{{ $t('alarm.actions.iconsFields.children') }} : {{ total }}</div>
        <div>{{ $t('alarm.fields.openedChildren') }} : {{ opened }}</div>
        <div>{{ $t('alarm.fields.closedChildren') }} : {{ closed }}</div>
      </div>
    </v-tooltip>
  </div>
</template>

<script>
import { EVENT_ENTITY_ICONS_BY_TYPE } from '@/constants';

export default {
  props: {
    rule: {
      type: Object,
      required: true,
    },
    total: {
      type: Number,
      default: 0,
    },
    opened: {
      type: Number,
      default: 0,
    },
    closed: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    ruleName() {
      return this.rule?.name ?? '';
    },

    icon() {
      return EVENT_ENTITY_ICONS_BY_TYPE.groupChildren;
    },
  },
};
</script>
