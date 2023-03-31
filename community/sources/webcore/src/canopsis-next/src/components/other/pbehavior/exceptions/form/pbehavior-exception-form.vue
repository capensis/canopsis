<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", required)
    v-text-field(
      v-field="form.description",
      v-validate="'required'",
      :label="$t('modals.createPbehaviorException.fields.description')",
      :error-messages="errors.collect('description')",
      name="description"
    )
    pbehavior-exceptions-field(v-model="form.exdates")
      template(#no-data="")
        v-alert(:value="true", type="info") {{ $t('modals.createPbehaviorException.emptyExdates') }}
</template>

<script>
import { formMixin, formArrayMixin } from '@/mixins/form';

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
