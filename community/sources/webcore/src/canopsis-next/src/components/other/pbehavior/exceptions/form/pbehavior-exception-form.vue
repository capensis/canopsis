<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('modals.createPbehaviorException.fields.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(row)
      v-text-field(
        v-field="form.description",
        v-validate="'required'",
        :label="$t('modals.createPbehaviorException.fields.description')",
        :error-messages="errors.collect('description')",
        name="description"
      )
    v-layout(row)
      pbehavior-exceptions-field(v-model="form.exdates")
        v-alert(
          slot="no-data",
          :value="true",
          type="info"
        ) {{ $t('modals.createPbehaviorException.emptyExdates') }}
</template>

<script>
import formMixin from '@/mixins/form';
import formArrayMixin from '@/mixins/form/array';

import PbehaviorExceptionsField from '@/components/other/pbehavior/exceptions/form/pbehavior-exceptions-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorExceptionsField },
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
