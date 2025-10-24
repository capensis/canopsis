<template>
  <v-layout column>
    <request-url-field
      v-if="!hideUrl"
      v-field="form"
      :help-text="$t('common.request.urlHelp')"
      :name="name"
      :disabled="disabled"
      :url-variables="urlVariables"
    />
    <v-layout>
      <v-flex
        class="mr-3"
        xs6
      >
        <c-information-block :title="$t('common.request.timeoutSettings')">
          <c-duration-field
            v-field="form.timeout"
            :disabled="disabled"
            :units-label="$t('common.unit')"
            clearable
          />
        </c-information-block>
      </v-flex>
      <v-flex xs6>
        <c-information-block :title="$t('common.request.repeatRequest')">
          <c-retry-field
            v-field="form"
            :disabled="disabled"
          />
        </c-information-block>
      </v-flex>
    </v-layout>
    <c-enabled-field
      v-field="form.skip_verify"
      :label="$t('common.request.skipVerify')"
      :disabled="disabled"
      hide-details
    />
    <request-auth-field
      v-field="form.auth"
      :name="`${name}.auth`"
      :disabled="disabled"
    />
    <c-information-block
      :title="$tc('common.header', 2)"
      :help-text="$t('common.request.headersHelpText')"
      class="mb-2"
      help-icon="help"
      help-icon-color="grey darken-1"
    >
      <request-headers-field
        v-field="form.headers"
        :name="`${name}.headers`"
        :disabled="disabled"
        :headers-variables="headersVariables"
      />
    </c-information-block>
    <c-payload-textarea-field
      v-field="form.payload"
      :label="$t('common.payload')"
      :line-height="16"
      :disabled="disabled"
      :variables="payloadVariables"
      :name="`${name}.payload`"
    />
  </v-layout>
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
    hideUrl: {
      type: Boolean,
      default: false,
    },
    urlVariables: {
      type: Array,
      default: () => [],
    },
    headersVariables: {
      type: Array,
      default: () => [],
    },
    payloadVariables: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
