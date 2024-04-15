<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <c-alarm-extra-details-chip :color="color" icon="center_focus_strong" v-on="on" />
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
import { computed } from 'vue';

import { COLORS } from '@/config';

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
  setup(props) {
    const ruleName = computed(() => props.rule?.name ?? '');
    const color = COLORS.alarmExtraDetails.children;

    return {
      ruleName,
      color,
    };
  },
};
</script>
