<template lang="pug">
  div
    v-layout(justify-space-between, align-center)
      v-flex.pa-1(xs6)
        v-select(
          v-field="request.method",
          v-validate="'required'",
          :items="availableMethods",
          :disabled="disabled",
          :label="$t('webhook.tabs.request.fields.method')",
          :error-messages="errors.collect('request.method')",
          name="request.method"
        )
      v-flex.pa-1(xs6)
        v-text-field(
          v-field="request.url",
          v-validate="'required'",
          :disabled="disabled",
          :label="$t('webhook.tabs.request.fields.url')",
          :error-messages="errors.collect('request.url')",
          name="request.url"
        )
    v-layout(row, wrap)
      v-flex(xs12)
        v-switch(v-model="withAuth", :label="$t('webhook.tabs.request.fields.authSwitch')", :disabled="disabled")
      template(v-if="withAuth")
        v-flex(xs12)
          h4.ml-1 {{ $t('webhook.tabs.request.fields.auth') }}
        v-flex.pa-1(xs6)
          v-text-field(
            v-field="request.auth.username",
            :disabled="disabled",
            :label="$t('webhook.tabs.request.fields.username')"
          )
        v-flex.pa-1(xs6)
          v-text-field(
            v-field="request.auth.password",
            :disabled="disabled",
            :label="$t('webhook.tabs.request.fields.password')"
          )
    text-pairs(
      v-field="request.headers",
      :disabled="disabled",
      :title="$t('webhook.tabs.request.fields.headers')",
      :textLabel="$t('webhook.tabs.request.fields.headerKey')",
      :valueLabel="$t('webhook.tabs.request.fields.headerValue')"
    )
    v-layout
      v-flex
        v-textarea(
          v-field="request.payload",
          v-validate="'required'",
          :disabled="disabled",
          :read-only="disabled",
          :label="$t('webhook.tabs.request.fields.payload')",
          :error-messages="errors.collect('request.payload')",
          name="request.payload"
        )
</template>

<script>
import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import TextPairs from '@/components/forms/fields/text-pairs.vue';

export default {
  inject: ['$validator'],
  components: { TextPairs },
  mixins: [formMixin, formValidationHeaderMixin],
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
        return Boolean(this.request.auth);
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
