<template lang="pug">
  v-layout(column)
    v-layout(row)
      slot(name="prepend-url")
      request-url-field(
        v-field="form",
        :help-text="$t('common.request.urlHelp')",
        :name="`${name}.request`",
        :disabled="disabled"
      )
      slot(name="append-url")
    v-layout(row)
      v-flex.mr-3(xs6)
        c-information-block(:title="$t('common.request.timeoutSettings')")
          c-duration-field(v-field="form.timeout", :disabled="disabled", clearable)
      v-flex(xs6)
        c-information-block(:title="$t('common.request.repeatRequest')")
          c-retry-field(v-field="form", :disabled="disabled")
    c-enabled-field(
      v-field="form.skip_verify",
      :label="$t('common.request.skipVerify')",
      :disabled="disabled",
      hide-details
    )
    request-auth-field(v-field="form.auth", :name="`${name}.auth`", :disabled="disabled")

    c-information-block.mb-2(
      :title="$tc('common.header', 2)",
      :help-text="$t('common.request.headersHelpText')",
      help-icon="help",
      help-icon-color="grey darken-1"
    )
      v-flex(v-if="!form.headers.length", xs12)
        v-alert(:value="true", type="info") {{ $t('common.request.emptyHeaders') }}
      request-headers-field(v-field="form.headers", :name="`${name}.headers`", :disabled="disabled")

    c-payload-field(
      v-field="form.payload",
      :label="$t('common.payload')",
      :line-height="16",
      :disabled="disabled"
    )
</template>

<script>
import { formMixin } from '@/mixins/form';

import RequestAuthField from './fields/request-auth-field.vue';
import RequestHeadersField from './fields/request-headers-field.vue';
import RequestUrlField from './fields/request-url-field.vue';

export default {
  inject: ['$validator'],
  components: { RequestUrlField, RequestHeadersField, RequestAuthField },
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
      default: 'request',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
