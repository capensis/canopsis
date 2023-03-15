<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", required)
    c-duration-field(
      v-field="form.delay",
      :label="$t('common.delay')",
      :units-label="$t('common.unit')",
      name="delay",
      clearable
    )
    c-enabled-field(v-field="form.enabled")
    c-triggers-field(:value="form.triggers", @input="updateField('triggers', $event)")
    c-disable-during-periods-field(v-field="form.disable_during_periods")
    c-priority-field(v-model="form.priority", :min="1", required)
    scenario-actions-form.mt-2(v-field="form.actions", name="actions")
</template>

<script>
import { formMixin } from '@/mixins/form';

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
