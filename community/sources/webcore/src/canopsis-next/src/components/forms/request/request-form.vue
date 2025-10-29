<template>
  <v-layout class="gap-2" column>
    <request-url-field
      v-if="!hideUrl"
      v-field="form"
      :help-text="urlHelpText ||$t('common.request.urlHelp')"
      :name="name"
      :disabled="disabled"
      :url-variables="urlVariables"
    />
    <v-layout v-if="withMultipleUrls">
      <v-flex
        offset-xs6
        xs6
      >
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
      </v-flex>
    </v-layout>
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
    <c-information-block :title="$t('user.auth')" class="mb-2 mt-2">
      <request-auth-with-token-field
        v-field="form.auth"
        :auth-token="authToken"
        :name="`${name}.auth`"
        :disabled="disabled"
        :only-credentials="!withAuthToken"
        @update:auth-token="updateAuthToken"
      />
    </c-information-block>
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
    urlHelpText: {
      type: String,
      default: '',
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
