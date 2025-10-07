<template>
  <request-form
    v-field="form.request"
    v-bind="$attrs"
    :name="name"
    :auth-token="form.auth_token"
    :headers-variables="variables"
    :payload-variables="variables"
    with-auth-token
    @update:auth-token="updateAuthToken"
  />
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';

import RequestForm from '@/components/forms/request/request-form.vue';

export default {
  components: { RequestForm },
  inheritAttrs: false, // TODO: check it
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
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const updateAuthToken = authToken => updateField('auth_token', authToken);

    return {
      updateAuthToken,
    };
  },
};
</script>
