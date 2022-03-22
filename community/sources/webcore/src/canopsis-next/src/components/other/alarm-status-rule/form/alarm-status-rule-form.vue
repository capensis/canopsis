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
    v-textarea(
      v-field="form.description",
      v-validate="'required'",
      :label="$t('common.description')",
      :error-messages="errors.collect('description')",
      name="description"
    )
    c-patterns-field(
      v-field="form.patterns",
      :some-required="flapping",
      with-alarm,
      with-entity
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
