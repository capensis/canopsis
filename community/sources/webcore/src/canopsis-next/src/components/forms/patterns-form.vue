<template lang="pug">
  v-layout(column)
    v-text-field(
      v-if="withTitle",
      v-field="form.title",
      v-validate="titleRules",
      :label="$t('common.title')",
      :error-messages="errors.collect('title')",
      name="title"
    )
    c-patterns-field.mt-2(
      v-field="form",
      :with-alarm="withAlarm",
      :with-entity="withEntity",
      :with-pbehavior="withPbehavior",
      :with-event="withEvent",
      :alarm-attributes="alarmAttributes",
      :alarm-excluded-attributes="alarmExcludedAttributes",
      :entity-attributes="entityAttributes",
      :entity-excluded-items="entityExcludedItems"
    )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    withTitle: {
      type: Boolean,
      default: false,
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
  },
  computed: {
    titleRules() {
      return {
        required: this.withTitle,
      };
    },
  },
};
</script>
