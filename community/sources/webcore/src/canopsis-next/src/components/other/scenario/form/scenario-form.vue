<template lang="pug">
  v-layout(column)
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(row)
      c-duration-field(
        v-field="form.delay",
        :label="$t('common.delay')",
        :units-label="$t('common.unit')",
        name="delay",
        clearable
      )
    v-layout(row)
      c-enabled-field(v-field="form.enabled")
    v-layout(row)
      c-triggers-field(
        :value="form.triggers",
        name="triggers",
        @input="updateField('triggers', $event)"
      )
    v-layout(row)
      c-disable-during-periods-field(
        v-field="form.disable_during_periods",
        name="disable_during_periods"
      )
    v-layout(row)
      v-text-field(
        v-field.number="form.priority",
        v-validate="'required'",
        :label="$t('common.priority')",
        :error-messages="errors.collect('priority')",
        :min="1",
        name="priority",
        type="number"
      )
    v-layout(column)
      scenario-actions-form(v-field="form.actions", name="actions")
</template>

<script>
import formMixin from '@/mixins/form/object';

import ScenarioActionsForm from './scenario-actions-form.vue';

export default {
  inject: ['$validator'],
  components: { ScenarioActionsForm },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
