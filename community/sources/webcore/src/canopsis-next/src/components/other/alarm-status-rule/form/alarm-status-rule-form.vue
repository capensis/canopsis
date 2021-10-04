<template lang="pug">
  v-layout(column)
    c-duration-field(v-field="form.duration", required)
    c-priority-field(v-field="form.priority", required)
    v-text-field(
      v-if="flapping",
      v-field.number="form.freq_limit",
      v-validate="'required|numeric|min_value:1'",
      :label="$t('alarmStatusRules.frequencyLimit')",
      :error-messages="errors.collect('freq_limit')",
      :min="1",
      name="freq_limit",
      type="number"
    )
    v-textarea(
      v-field="form.description",
      v-validate="'required'",
      :label="$t('common.description')",
      :error-messages="errors.collect('description')",
      name="description"
    )
    c-patterns-field(
      v-field="form.patterns",
      alarm,
      entity,
      some-required
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
