<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="form._id",
        v-validate="'required'",
        :label="$t('modals.createRight.fields.id')",
        :error-messages="errors.collect('id')",
        name="id"
      )
    v-layout(row)
      v-text-field(
        v-field="form.desc",
        :label="$t('modals.createRight.fields.description')"
      )
    v-layout(row)
      v-select(
        v-field="form.type",
        :label="$t('modals.createRight.fields.type')",
        :items="types"
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
