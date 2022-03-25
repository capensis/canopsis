<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.title",
      v-validate="titleRules",
      :label="$t('common.title')",
      :error-messages="errors.collect('title')",
      name="title"
    )
    c-alarm-patterns-field(v-if="isAlarmPattern", v-field="form")
    c-entity-patterns-field(v-else-if="isEntityPattern", v-field="form")
    c-pbehavior-patterns-field(v-else-if="isPbehaviorPattern", v-field="form")
</template>

<script>
import { PATTERN_TYPES } from '@/constants';

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
  },
  computed: {
    isAlarmPattern() {
      return this.form.type === PATTERN_TYPES.alarm;
    },

    isEntityPattern() {
      return this.form.type === PATTERN_TYPES.entity;
    },

    isPbehaviorPattern() {
      return this.form.type === PATTERN_TYPES.pbehavior;
    },

    titleRules() {
      return {
        required: true,
      };
    },
  },
};
</script>
