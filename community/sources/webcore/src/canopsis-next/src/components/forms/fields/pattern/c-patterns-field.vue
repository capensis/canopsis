<template lang="pug">
  v-layout(column)
    c-pattern-panel.mb-2(v-if="withAlarm", :title="$t('common.alarmPatterns')")
      c-alarm-patterns-field(
        v-field="value.alarm_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="alarmFieldName",
        :excluded="alarmExcludedAttributes",
        :attributes="alarmAttributes",
        with-type,
        @input="errors.remove(alarmFieldName)"
      )

    c-pattern-panel.mb-2(v-if="withEntity", :title="$t('common.entityPatterns')")
      c-entity-patterns-field(
        v-field="value.entity_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="entityFieldName",
        :attributes="entityAttributes",
        :excluded="entityExcludedItems",
        with-type,
        @input="errors.remove(entityFieldName)"
      )

    c-pattern-panel.mb-2(v-if="withPbehavior", :title="$t('common.pbehaviorPatterns')")
      c-pbehavior-patterns-field(
        v-field="value.pbehavior_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="pbehaviorFieldName",
        with-type,
        @input="errors.remove(pbehaviorFieldName)"
      )

    c-pattern-panel.mb-2(v-if="withEvent", :title="$t('common.eventPatterns')")
      c-event-filter-patterns-field(
        v-field="value.event_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="eventFieldName",
        @input="errors.remove(eventFieldName)"
      )

    c-pattern-panel(v-if="withTotalEntity", :title="$t('common.totalEntityPatterns')")
      c-entity-patterns-field(
        v-field="value.total_entity_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="totalEntityFieldName",
        with-type,
        @input="errors.remove(totalEntityFieldName)"
      )
</template>

<script>
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
    alarmExcludedAttributes: {
      type: Array,
      required: false,
    },
    entityAttributes: {
      type: Array,
      required: false,
    },
    entityExcludedItems: {
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
    withTotalEntity: {
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

    alarmFieldName() {
      return this.preparePatternsFieldName('alarm_pattern');
    },

    eventFieldName() {
      return this.preparePatternsFieldName('event_pattern');
    },

    entityFieldName() {
      return this.preparePatternsFieldName('entity_pattern');
    },

    pbehaviorFieldName() {
      return this.preparePatternsFieldName('pbehavior_pattern');
    },

    totalEntityFieldName() {
      return this.preparePatternsFieldName('total_entity_pattern');
    },
  },
  methods: {
    preparePatternsFieldName(name) {
      return [this.name, name].filter(Boolean).join('.');
    },
  },
};
</script>
