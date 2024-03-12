<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <c-duration-field
      v-field="form.duration"
      required
    />
    <c-priority-field v-field="form.priority" />
    <c-number-field
      v-if="flapping"
      v-field="form.freq_limit"
      :label="$t('common.frequencyLimit')"
      :min="1"
      name="freq_limit"
    />
    <c-description-field
      v-field="form.description"
      required
    />
    <alarm-status-rule-patterns-form
      v-field="form.patterns"
      :flapping="flapping"
      class="mt-2"
    />
  </v-layout>
</template>

<script>
import AlarmStatusRulePatternsForm from './alarm-status-rule-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: { AlarmStatusRulePatternsForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    flapping: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
