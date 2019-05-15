<template lang="pug">
  div
    v-layout(justify-space-between, align-center)
      v-flex(xs6).pa-1
        v-select(
        :value="request.method",
        :items="availableMethods",
        :disabled="disabled",
        :label="$t('webhook.tabs.request.fields.method')",
        v-validate="'required'",
        name="request.method",
        :error-messages="errors.collect('request.method')",
        @input="updateField('method', $event)"
        )
      v-flex(xs6).pa-1
        v-text-field(
        :value="request.url",
        :disabled="disabled",
        :label="$t('webhook.tabs.request.fields.url')",
        v-validate="'required|url'",
        name="request.url",
        :error-messages="errors.collect('request.url')",
        @input="updateField('url', $event)"
        )
    v-layout(row, wrap)
      v-flex(xs12)
        v-switch(v-model="withAuth", :label="$t('webhook.tabs.request.fields.authSwitch')", :disabled="disabled")
      template(v-if="withAuth")
        v-flex(xs12)
          h4.ml-1 {{ $t('webhook.tabs.request.fields.auth') }}
        v-flex(xs6).pa-1
          v-text-field(
          :value="request | get('auth.username')",
          :disabled="disabled",
          :label="$t('webhook.tabs.request.fields.username')",
          @input="updateField('auth.username', $event)"
          )
        v-flex(xs6).pa-1
          v-text-field(
          :value="request | get('auth.password')",
          :disabled="disabled",
          :label="$t('webhook.tabs.request.fields.password')",
          @input="updateField('auth.password', $event)"
          )
    text-pairs(
    :items="request.headers",
    :disabled="disabled",
    :title="$t('webhook.tabs.request.fields.headers')",
    :textLabel="$t('webhook.tabs.request.fields.headerKey')",
    :valueLabel="$t('webhook.tabs.request.fields.headerValue')",
    @input="updateField('headers', $event)"
    )
    v-layout
      v-flex
        v-textarea(
        :value="request.payload",
        :disabled="disabled",
        :read-only="disabled",
        :label="$t('webhook.tabs.request.fields.payload')",
        v-validate="'required'",
        name="request.payload",
        :error-messages="errors.collect('request.payload')",
        @input="updateField('payload', $event)"
        )
</template>

<script>
import formMixin from '@/mixins/form';

import TextPairs from '@/components/forms/fields/text-pairs.vue';

export default {
  inject: ['$validator'],
  components: { TextPairs },
  mixins: [formMixin],
  model: {
    prop: 'request',
    event: 'input',
  },
  props: {
    request: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      headers: [],
    };
  },
  computed: {
    withAuth: {
      get() {
        return !!this.request.auth;
      },
      set(value) {
        if (value) {
          this.updateField('auth', { username: '', password: '' });
        } else {
          this.removeField('auth');
        }
      },
    },

    availableMethods() {
      return [
        'POST', 'GET', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'CONNECT', 'OPTIONS', 'TRACE',
      ];
    },
  },
};
</script>
