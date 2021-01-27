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
        :units-label="$t('common.unit')"
      )
    v-layout(row)
      c-enabled-field(v-field="form.enabled")
    v-layout(row)
      c-triggers-field(
        :value="form.triggers",
        @input="updateField('triggers', $event)"
      )
    v-layout(row)
      c-disable-during-periods-field(v-field="form.disable_during_periods")
    v-layout(row)
      v-text-field(
        v-field.number="form.priority",
        :label="$t('common.priority')",
        :min="0",
        type="number"
      )
</template>

<script>
import formMixin from '@/mixins/form/object';

export default {
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
