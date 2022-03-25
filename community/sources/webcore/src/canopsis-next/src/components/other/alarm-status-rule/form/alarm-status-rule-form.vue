<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    c-duration-field(v-field="form.duration", required)
    c-priority-field(v-field="form.priority", required)
    c-number-field(
      v-if="flapping",
      v-field="form.freq_limit",
      :label="$t('alarmStatusRules.frequencyLimit')",
      :min="1",
      name="freq_limit"
    )
    c-description-field(v-field="form.description", required)
    c-patterns-field(
      v-field="form.patterns",
      :some-required="flapping",
      alarm,
      entity
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
      default: () => ({}),
    },
    flapping: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
