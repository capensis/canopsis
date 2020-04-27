<template lang="pug">
  v-layout(row)
    v-flex(xs12)
      v-tabs(v-model="activePatternTab", fixed-tabs, slider-color="primary")
        v-tab {{ $t('metaAlarmRule.fields.alarmPatterns') }}
        v-tab {{ $t('metaAlarmRule.fields.eventPatterns') }}
        v-tab {{ $t('metaAlarmRule.fields.entityPatterns') }}
        v-tab-item
          patterns-list(v-field="patterns.alarm_patterns", :disabled="disabled")
        v-tab-item
          patterns-list(v-field="patterns.event_patterns", :disabled="disabled")
        v-tab-item
          patterns-list(v-field="patterns.entity_patterns", :disabled="disabled")
</template>

<script>
import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  inject: ['$validator'],
  components: { PatternsList },
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      activePatternTab: 0,
    };
  },
};
</script>

