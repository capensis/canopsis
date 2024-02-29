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
          <v-icon color="white" small>center_focus_weak</v-icon>
        </span>
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
};
</script>

<style lang="scss" scoped>
.rule-name:nth-of-type(odd) {
  color: #b5b5b5;
}
</style>
