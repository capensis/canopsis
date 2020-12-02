<template lang="pug">
  v-layout(row)
    v-flex(xs12)
      v-tabs(v-model="activePatternTab", slider-color="primary", fixed-tabs)
        template(v-if="alarm")
          v-tab(:class="{ 'error--text': errors.has('alarm_patterns') }") {{ $t('common.alarmPatterns') }}
          v-tab-item
            patterns-list(
              v-field="value.alarm_patterns",
              :disabled="disabled",
              name="alarm_patterns",
              @input="errors.remove('alarm_patterns')"
            )
        template(v-if="event")
          v-tab(:class="{ 'error--text': errors.has('event_patterns') }") {{ $t('common.eventPatterns') }}
          v-tab-item
            patterns-list(
              v-field="value.event_patterns",
              :disabled="disabled",
              name="event_patterns",
              @input="errors.remove('event_patterns')"
            )
        template(v-if="entity")
          v-tab(:class="{ 'error--text': errors.has('entity_patterns') }") {{ $t('common.entityPatterns') }}
          v-tab-item
            patterns-list(
              v-field="value.entity_patterns",
              :disabled="disabled",
              name="entity_patterns",
              @input="errors.remove('entity_patterns')"
            )
        template(v-if="totalEntity")
          v-tab(:class="{ 'error--text': errors.has('total_entity_patterns') }") {{ $t('common.totalEntityPatterns') }}
          v-tab-item
            patterns-list(
              v-field="value.total_entity_patterns",
              :disabled="disabled",
              name="total_entity_patterns",
              @input="errors.remove('total_entity_patterns')"
            )
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
    totalEntity: {
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

