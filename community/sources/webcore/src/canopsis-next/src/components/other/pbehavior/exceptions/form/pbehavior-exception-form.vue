<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <v-text-field
      v-field="form.description"
      v-validate="'required'"
      :label="$t('modals.createPbehaviorException.fields.description')"
      :error-messages="errors.collect('description')"
      name="description"
    />
    <pbehavior-exceptions-field
      v-field="form.exdates"
      with-exdate-type
    >
      <template #no-data="">
        <c-alert type="info">
          {{ $t('modals.createPbehaviorException.emptyExdates') }}
        </c-alert>
      </template>
    </pbehavior-exceptions-field>
  </v-layout>
</template>

<script>
import { formMixin, formArrayMixin } from '@/mixins/form';

import PbehaviorExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exceptions-field.vue';

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
