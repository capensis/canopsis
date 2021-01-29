<template lang="pug">
  v-layout(column)
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        :disabled="disabled",
        name="name"
      )
    v-layout(row)
      c-duration-field(
        v-field="form.delay",
        :label="$t('common.delay')",
        :units-label="$t('common.unit')",
        :disabled="disabled",
        name="delay"
      )
    v-layout(row)
      c-enabled-field(v-field="form.enabled")
    v-layout(row)
      c-triggers-field(
        :value="form.triggers",
        :disabled="disabled",
        name="triggers",
        @input="updateField('triggers', $event)"
      )
    v-layout(row)
      c-disable-during-periods-field(
        v-field="form.disable_during_periods",
        name="disable_during_periods",
        :disabled="disabled"
      )
    v-layout(row)
      v-text-field(
        v-field.number="form.priority",
        v-validate="'required'",
        :label="$t('common.priority')",
        :error-messages="errors.collect('priority')",
        :disabled="disabled",
        :min="1",
        name="priority",
        type="number"
      )
    v-layout(column)
      scenario-action-field(
        v-for="(action, index) in form.actions",
        v-field="form.actions[index]",
        :name="`actions.${index}`",
        :key="action.key"
      )
</template>

<script>
import formMixin from '@/mixins/form/object';

import ScenarioActionField from '@/components/other/scenarios/form/fields/scenario-action-field.vue';

export default {
  components: { ScenarioActionField },
  inject: ['$validator'],
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
    name: {
      type: String,
      default: 'scenario',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
