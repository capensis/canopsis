<template lang="pug">
  v-layout(column)
    c-request-url-field(
      v-field="form.request",
      :help-text="$t('scenario.urlHelp')",
      :name="`${name}.request`",
      :disabled="disabled"
    )
    c-retry-field(v-field="form.retry", :disabled="disabled")
    c-enabled-field(
      v-field="form.request.skip_verify",
      :label="$t('scenario.skipVerify')",
      :disabled="disabled"
    )
    c-enabled-field(
      v-model="withAuth",
      :label="$t('scenario.withAuth')",
      :disabled="disabled"
    )
    v-layout(v-if="withAuth", row)
      v-flex.pa-1(xs6)
        v-text-field(
          v-field="form.request.auth.username",
          :label="$t('common.username')",
          :name="`${name}.username`",
          :disabled="disabled"
        )
      v-flex.pa-1(xs6)
        v-text-field(
          v-field="form.request.auth.password",
          :label="$t('common.password')",
          :name="`${name}.password`",
          :disabled="disabled"
        )
    c-text-pairs-field(
      v-field="form.request.headers",
      :title="$t('scenario.headers')",
      :text-label="$t('scenario.headerKey')",
      :value-label="$t('scenario.headerValue')",
      :name="name",
      :disabled="disabled"
    )
    v-textarea(
      v-field="form.request.payload",
      :label="$t('common.payload')",
      :disabled="disabled"
    )
      template(#append="")
        c-help-icon(icon="help", :text="$t('scenario.payloadHelp')", left)
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
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    withAuth: {
      set(value) {
        this.updateField('request.auth', value ? { username: '', password: '' } : undefined);
      },
      get() {
        return !!this.form?.request?.auth;
      },
    },
  },
};
</script>
