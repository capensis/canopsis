<template>
  <request-form
    v-field="form.request"
    v-bind="$attrs"
    :name="name"
    :auth-token="form.auth_token"
    :multiple="form.multiple_urls"
    :headers-variables="headersVariables"
    :payload-variables="payloadVariables"
    :url-variables="urlVariables"
    with-auth-token
    @update:auth-token="updateAuthToken"
    @update:multiple="updateMultiple"
  />
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';

import RequestForm from '@/components/forms/request/request-form.vue';

export default {
  components: { RequestForm },
  inheritAttrs: false,
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
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const updateAuthToken = authToken => updateField('auth_token', authToken);
    const updateMultiple = multiple => updateField('multiple_urls', multiple);

    return {
      updateAuthToken,
      updateMultiple,
    };
  },
};
</script>
