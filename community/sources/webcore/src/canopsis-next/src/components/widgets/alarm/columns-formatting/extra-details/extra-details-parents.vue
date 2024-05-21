<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
      disable-resize
    >
      <template #activator="{ on }">
        <c-alarm-extra-details-chip :color="color" icon="center_focus_weak" v-on="on" />
      </template>
      <div class="text-md-center">
        <strong>{{ $t('alarm.actions.iconsTitles.grouping') }}</strong>
        <v-layout column>
          <div>{{ $tc('common.rule', rules.length) }}&nbsp;:</div>
          <div
            v-for="rule in rules"
            :key="rule.id"
            class="rule-name"
          >
            &nbsp;{{ rule.name }}
          </div>
        </v-layout>
        <div>{{ $t('alarm.actions.iconsFields.parents') }} : {{ total }}</div>
      </div>
    </v-tooltip>
  </div>
</template>

<script>
import { COLORS } from '@/config';

export default {
  props: {
    rules: {
      type: Array,
      default: () => [],
    },
    total: {
      type: Number,
      default: 0,
    },
  },
  setup() {
    const color = COLORS.alarmExtraDetails.parent;

    return {
      color,
    };
  },
};
</script>

<style lang="scss" scoped>
.rule-name:nth-of-type(odd) {
  color: #b5b5b5;
}
</style>
