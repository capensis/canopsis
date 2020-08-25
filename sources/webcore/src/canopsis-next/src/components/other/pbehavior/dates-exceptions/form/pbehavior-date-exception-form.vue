<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('modals.createPbehaviorDateException.fields.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(row)
      v-text-field(
        v-field="form.description",
        v-validate="'required'",
        :label="$t('modals.createPbehaviorDateException.fields.description')",
        :error-messages="errors.collect('description')",
        name="description"
      )
    v-layout(row)
      pbehavior-exceptions-dates-field(v-model="form.exdates")
</template>

<script>
import formMixin from '@/mixins/form';
import formArrayMixin from '@/mixins/form/array';

import PbehaviorExceptionsDatesField from '@/components/other/pbehavior/dates-exceptions/form/pbehavior-exceptions-dates-field.vue';

export default {
  components: { PbehaviorExceptionsDatesField },
  inject: ['$validator'],
  mixins: [formMixin, formArrayMixin],
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
