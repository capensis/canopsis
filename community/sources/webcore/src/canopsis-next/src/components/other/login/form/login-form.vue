<template>
  <v-layout column>
    <v-flex class="mt-1">
      <v-text-field
        v-field.mutate="form.username"
        v-validate="'required'"
        :label="$t('common.username')"
        :error-messages="errors.collect('username')"
        color="primary"
        name="username"
        autofocus
        clearable
        outlined
      />
    </v-flex>
    <v-flex>
      <c-password-field
        v-field.mutate="form.password"
        required
        clearable
        outlined
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';

export default {
  inject: ['$validator'],
  mixins: [authMixin, entitiesInfoMixin],
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
