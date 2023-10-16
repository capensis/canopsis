<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <v-text-field
      v-field="form.description"
      v-validate="'required'"
      :label="$t('modals.createPbehaviorReason.fields.description')"
      :error-messages="errors.collect('description')"
      name="description"
    />
    <c-enabled-field
      v-field="form.hidden"
      :label="$t('pbehavior.reasons.hidden')"
    />
  </v-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

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
  },
};
</script>
