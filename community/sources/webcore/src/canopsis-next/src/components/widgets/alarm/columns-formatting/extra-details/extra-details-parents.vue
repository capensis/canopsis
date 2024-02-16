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
        <v-layout column>
          <div>{{ $tc('common.rule', rules.length) }}&nbsp;:</div>
          <div
            class="rule-name"
            v-for="rule in rules"
            :key="rule.id"
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
import { EVENT_ENTITY_ICONS_BY_TYPE } from '@/constants';

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
  computed: {
    icon() {
      return EVENT_ENTITY_ICONS_BY_TYPE.groupParents;
    },
  },
};
</script>

<style lang="scss" scoped>
.rule-name:nth-of-type(odd) {
  color: #b5b5b5;
}
</style>
