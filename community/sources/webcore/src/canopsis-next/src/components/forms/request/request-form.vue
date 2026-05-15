<template>
  <div>
    <c-form-block-row v-if="!hideUrl" :label="urlLabel || $t('common.url')" :depth="depth">
      <request-url-field
        v-field="form"
        :help-text="urlHelpText ||$t('common.request.urlHelp')"
        :name="name"
        :disabled="disabled"
        :url-variables="urlVariables"
      />
    </c-form-block-row>

    <c-form-block-row v-if="withMultipleUrls" :label="$t('scenario.allowMultipleUrls')" :depth="depth">
      <c-enabled-field
        :value="multiple"
        :label="$t('scenario.allowMultipleUrls')"
        :disabled="disabled"
        @input="updateMultiple"
      >
        <template #append>
          <c-help-icon
            :text="$t('scenario.allowMultipleUrlsTooltip')"
            icon="help"
            color="grey darken-1"
            left
          />
        </template>
      </c-enabled-field>
    </c-form-block-row>

    <c-form-block-row :label="$t('common.request.timeoutSettings')" :depth="depth">
      <c-duration-field
        v-field="form.timeout"
        :disabled="disabled"
        :units-label="$t('common.unit')"
        clearable
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('common.request.repeatRequest')" :depth="depth">
      <span v-if="hideRepeat" class="font-italic mt-4">
        {{ $t('common.request.repeatRequestInTomlFile') }}
      </span>
      <c-retry-field
        v-else
        v-field="form"
        :disabled="disabled"
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('common.request.skipVerify')" :depth="depth">
      <c-enabled-field
        v-field="form.skip_verify"
        :label="$t('common.request.skipVerify')"
        :disabled="disabled"
        hide-details
      />
    </c-form-block-row>

    <slot name="additional-fields" />

    <c-form-block-row v-if="!hideAuth" :label="$t('user.auth')" :depth="depth">
      <request-auth-with-token-field
        v-field="form.auth"
        :auth-token="authToken"
        :name="`${name}.auth`"
        :disabled="disabled"
        :only-credentials="!withAuthToken"
        @update:auth-token="updateAuthToken"
      />
    </c-form-block-row>

    <c-form-block-row
      v-if="!hideHeaders"
      :label="$tc('common.header', 2)"
      :depth="depth"
      indented
    >
      <request-headers-field
        v-field="form.headers"
        :name="`${name}.headers`"
        :disabled="disabled"
        :headers-variables="headersVariables"
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('common.payload')" :depth="depth">
      <c-payload-textarea-field
        v-field="form.payload"
        :label="$t('common.payload')"
        :line-height="16"
        :disabled="disabled"
        :variables="payloadVariables"
        :name="`${name}.payload`"
      />
    </c-form-block-row>
  </div>
</template>

<script>
import RequestUrlField from './fields/request-url-field.vue';
import RequestHeadersField from './fields/request-headers-field.vue';
import RequestAuthWithTokenField from './fields/request-auth-with-token-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RequestUrlField,
    RequestHeadersField,
    RequestAuthWithTokenField,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    authToken: {
      type: Object,
      default: () => ({}),
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
    hideRepeat: {
      type: Boolean,
      default: false,
    },
    hideAuth: {
      type: Boolean,
      default: false,
    },
    hideHeaders: {
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
    withMultipleUrls: {
      type: Boolean,
      default: false,
    },
    withAuthToken: {
      type: Boolean,
      default: false,
    },
    urlLabel: {
      type: String,
      default: '',
    },
    urlHelpText: {
      type: String,
      default: '',
    },
    depth: {
      type: [Number, String],
      default: 0,
    },
  },

  setup(props, { emit }) {
    const updateMultiple = multiple => emit('update:multiple', multiple);
    const updateAuthToken = authToken => emit('update:auth-token', authToken);

    return {
      updateMultiple,
      updateAuthToken,
    };
  },
};
</script>
