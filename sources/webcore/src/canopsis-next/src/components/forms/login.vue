<template lang="pug">
  v-form(@submit.prevent="submit")
    v-flex(px-3)
      v-text-field(
      :label="$t('common.username')",
      :error-messages="errors.collect('username')",
      :value="value.username",
      @input="updateField('username', $event)",
      v-validate="'required'",
      color="primary",
      name="username",
      autofocus,
      clearable,
      outline
      )
    v-flex(px-3)
      v-text-field(
      :label="$t('common.password')",
      :error-messages="errors.collect('password')",
      :value="value.password",
      @input="updateField('password', $event)",
      v-validate="'required'",
      color="primary",
      name="password",
      type="password",
      clearable,
      outline
      )
    v-flex.px-3.py-2
      v-alert(:value="hasServerError", type="error")
        span {{ $t('login.errors.incorrectEmailOrPassword') }}
</template>

<script>
import formMixin from '@/mixins/form';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      required: true,
    },
    hasServerError: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
