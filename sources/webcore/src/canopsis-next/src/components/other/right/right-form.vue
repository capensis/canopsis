<template lang="pug">
  v-form
    v-layout(row)
      v-text-field(
        :value="form._id",
        v-validate="'required'",
        :label="$t('modals.createRight.fields.id')",
        :error-messages="errors.collect('id')",
        name="id",
        @input="updateField('_id', $event)"
      )
    v-layout(row)
      v-text-field(
        :value="form.desc",
        :label="$t('modals.createRight.fields.description')",
        @input="updateField('desc', $event)"
      )
    v-layout(row)
      v-select(
        :value="form.type",
        :label="$t('modals.createRight.fields.type')",
        :items="types",
        @input="updateField('type', $event)"
      )
</template>

<script>
import { USERS_RIGHTS_TYPES } from '@/constants';

import formMixin from '@/mixins/form';

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
  computed: {
    types() {
      return [
        { value: '', text: 'Default' },
        { value: USERS_RIGHTS_TYPES.rw, text: USERS_RIGHTS_TYPES.rw },
        { value: USERS_RIGHTS_TYPES.crud, text: USERS_RIGHTS_TYPES.crud },
      ];
    },
  },
};
</script>
