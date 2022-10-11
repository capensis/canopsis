<template lang="pug">
  v-layout(column)
    c-request-url-field(
      v-field="form.request",
      :help-text="$t('scenario.urlHelp')",
      :name="`${name}.request`"
    )
    c-retry-field(v-field="form.retry")
    c-enabled-field(v-model="form.request.skip_verify", :label="$t('scenario.skipVerify')")
    c-enabled-field(v-model="withAuth", :label="$t('scenario.withAuth')")
    v-layout(v-if="withAuth", row)
      v-flex.pa-1(xs6)
        v-text-field(
          v-field="form.request.auth.username",
          :label="$t('common.username')",
          :name="`${name}.username`"
        )
      v-flex.pa-1(xs6)
        v-text-field(
          v-field="form.request.auth.password",
          :label="$t('common.password')",
          :name="`${name}.password`"
        )
    c-text-pairs-field(
      v-field="form.request.headers",
      :title="$t('scenario.headers')",
      :text-label="$t('scenario.headerKey')",
      :value-label="$t('scenario.headerValue')",
      :name="name"
    )
    v-textarea(
      v-field="form.request.payload",
      :label="$t('common.payload')"
    )
      v-tooltip(slot="append", left)
        v-icon(slot="activator") help
        div(v-html="$t('scenario.payloadHelp')")
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
