<template lang="pug">
  v-layout(column)
    c-pattern-panel.mb-2(v-if="withAlarm", :title="$t('common.alarmPatterns')")
      c-alarm-patterns-field(
        v-field="value.alarm_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="alarmPatternsFieldName",
        :check-count-name="$constants.PATTERNS_FIELDS.alarm",
        :attributes="alarmAttributes",
        with-type,
        @input="errors.remove(alarmPatternsFieldName)"
      )

    c-pattern-panel.mb-2(v-if="withEntity", :title="$t('common.entityPatterns')")
      c-entity-patterns-field(
        v-field="value.entity_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="entityPatternsFieldName",
        :check-count-name="$constants.PATTERNS_FIELDS.entity",
        :attributes="entityAttributes",
        with-type,
        @input="errors.remove(entityPatternsFieldName)"
      )

    c-pattern-panel.mb-2(v-if="withPbehavior", :title="$t('common.pbehaviorPatterns')")
      c-pbehavior-patterns-field(
        v-field="value.pbehavior_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="entityPbehaviorFieldName",
        :check-count-name="$constants.PATTERNS_FIELDS.pbehavior",
        with-type,
        @input="errors.remove(entityPbehaviorFieldName)"
      )

    c-pattern-panel(v-if="withEvent", :title="$t('common.eventPatterns')")
      c-event-filter-patterns-field(
        v-field="value.event_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="eventPatternsFieldName",
        @input="errors.remove(eventPatternsFieldName)"
      )
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

export default {
  inject: ['$validator'],
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
    alarmAttributes: {
      type: Array,
      required: false,
    },
    entityAttributes: {
      type: Array,
      required: false,
    },
    withAlarm: {
      type: Boolean,
      default: false,
    },
    withEvent: {
      type: Boolean,
      default: false,
    },
    withEntity: {
      type: Boolean,
      default: false,
    },
    withPbehavior: {
      type: Boolean,
      default: false,
    },
    totalEntity: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    someRequired: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      activePatternTab: 0,
    };
  },
  computed: {
    hasPatterns() {
      const {
        alarm_pattern: alarmPattern,
        entity_pattern: entityPattern,
        pbehavior_pattern: pbehaviorPattern,
        event_pattern: eventPattern,
      } = this.value;

      return alarmPattern?.groups?.length
        || entityPattern?.groups?.length
        || pbehaviorPattern?.groups?.length
        || eventPattern?.groups?.length;
    },

    isPatternRequired() {
      return this.someRequired ? !this.hasPatterns : this.required;
    },

    alarmPatternsFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.alarm);
    },

    eventPatternsFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.event);
    },

    entityPatternsFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.entity);
    },

    entityPbehaviorFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.pbehavior);
    },
  },
  methods: {
    preparePatternsFieldName(name) {
      return [this.name, name].filter(Boolean).join('.');
    },
  },
};
</script>
