<template lang="pug">
  v-layout(column)
    c-request-url-field(
      v-field="webhhook.request",
      :help-text="$t('scenario.urlHelp')",
      :name="`${name}.request`"
    )
    c-retry-field(v-field="webhhook.retry")
    c-enabled-field(v-model="webhhook.request.skip_verify", :label="$t('scenario.skipVerify')")
    c-enabled-field(v-model="withAuth", :label="$t('scenario.withAuth')")
    v-layout(v-if="withAuth", row)
      v-flex.pa-1(xs6)
        v-text-field(
          v-field="webhhook.request.auth.username",
          :label="$t('common.username')",
          :name="`${name}.username`"
        )
      v-flex.pa-1(xs6)
        v-text-field(
          v-field="webhhook.request.auth.password",
          :label="$t('common.password')",
          :name="`${name}.password`"
        )
    c-text-pairs-field(
      v-field="webhhook.request.headers",
      :title="$t('scenario.headers')",
      :text-label="$t('scenario.headerKey')",
      :value-label="$t('scenario.headerValue')",
      :name="name"
    )
    v-textarea(
      v-field="webhhook.request.payload",
      :label="$t('common.payload')"
    )
      v-tooltip(slot="append", left)
        v-icon(slot="activator") help
        div(v-html="$t('scenario.payloadHelp')")
    h4.ml-1 {{ $t('scenario.declareTicket') }}
    c-enabled-field(v-model="webhhook.empty_response", :label="$t('scenario.emptyResponse')")
    c-enabled-field(v-model="webhhook.is_regexp", :label="$t('scenario.isRegexp')")
    c-text-pairs-field(
      v-field="webhhook.declare_ticket",
      :text-label="$t('scenario.key')",
      :name="name"
    )
</template>

<script>
import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'webhhook',
    event: 'input',
  },
  props: {
    webhhook: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: 'webhook',
    },
  },
  computed: {
    withAuth: {
      set(value) {
        this.updateField('request.auth', value ? { username: '', password: '' } : undefined);
      },
      get() {
        return !!this.webhhook.request.auth;
      },
    },
  },
};
</script>