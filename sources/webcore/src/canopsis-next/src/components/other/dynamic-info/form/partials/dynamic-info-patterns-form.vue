<template lang="pug">
  div
    v-tabs(centered, slider-color="primary")
      v-tab {{ $t('modals.createDynamicInfo.steps.patterns.alarmPatterns') }}
      v-tab {{ $t('modals.createDynamicInfo.steps.patterns.entityPatterns') }}
      v-tab-item
        patterns-list(
          v-field="form.alarm_patterns",
          :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS"
        )
      v-tab-item
        patterns-list(
          v-field="form.entity_patterns",
          :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS"
        )
</template>

<script>
import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  components: {
    PatternsList,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
