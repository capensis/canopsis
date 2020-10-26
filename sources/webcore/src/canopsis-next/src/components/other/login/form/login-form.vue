<template lang="pug">
  v-layout(column)
    v-flex.mt-1
      v-text-field(
        v-field="form.username",
        v-validate="'required'",
        :label="$t('common.username')",
        :error-messages="errors.collect('username')",
        color="primary",
        name="username",
        data-test="username",
        autofocus,
        clearable,
        outline
      )
    v-flex
      v-text-field(
        v-field="form.password",
        v-validate="'required'",
        :label="$t('common.password')",
        :error-messages="errors.collect('password')",
        color="primary",
        name="password",
        type="password",
        data-test="password",
        clearable,
        outline
      )
</template>

<script>
import authMixin from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';

export default {
  mixins: [authMixin, entitiesInfoMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
