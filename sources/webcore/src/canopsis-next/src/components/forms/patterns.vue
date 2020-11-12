<template lang="pug">
  v-layout(row)
    v-flex(xs12)
      v-tabs(v-model="activePatternTab", fixed-tabs, slider-color="primary")
        template(v-if="alarm")
          v-tab {{ $t('common.alarmPatterns') }}
          v-tab-item
            patterns-list(v-field="value.alarm_patterns", :disabled="disabled")
        template(v-if="event")
          v-tab {{ $t('common.eventPatterns') }}
          v-tab-item
            patterns-list(v-field="value.event_patterns", :disabled="disabled")
        template(v-if="entity")
          v-tab {{ $t('common.entityPatterns') }}
          v-tab-item
            patterns-list(v-field="value.entity_patterns", :disabled="disabled")
</template>

<script>
import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  inject: ['$validator'],
  components: { PatternsList },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    alarm: {
      type: Boolean,
      default: false,
    },
    event: {
      type: Boolean,
      default: false,
    },
    entity: {
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

