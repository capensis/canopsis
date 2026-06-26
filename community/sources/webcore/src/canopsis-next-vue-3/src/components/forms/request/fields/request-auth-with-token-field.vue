<template>
  <c-layout gap="2" column>
    <c-radio-group
      :model-value="value.type"
      :disabled="disabled"
      column
      @update:model-value="updateField('type', $event)"
    >
      <v-radio
        v-for="type in types"
        :key="type.value"
        :label="type.label"
        :value="type.value"
      />
    </c-radio-group>
    <c-layout v-if="isCredentialsType" gap="2">
      <c-name-field
        :model-value="value.username"
        :label="$t('common.username')"
        :name="`${name}.username`"
        :disabled="disabled"
        @update:model-value="updateField('username', $event)"
      />
      <c-password-field
        :model-value="value.password"
        :name="`${name}.password`"
        :disabled="disabled"
        @update:model-value="updateField('password', $event)"
      />
    </c-layout>
    <request-auth-token-field
      v-if="isTokenType"
      :value="authToken"
      :name="`${name}.token`"
      :disabled="disabled"
      @update:value="updateAuthToken"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { REQUEST_AUTH_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

import RequestAuthTokenField from '@/components/forms/request/fields/request-auth-token-field.vue';

export default {
  emits: ['update:value'],
  components: { RequestAuthTokenField },
  props: {
    value: {
      type: Object,
      required: true,
    },
    authToken: {
      type: Object,
      default: () => ({}),
    },
    onlyCredentials: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'auth',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const types = computed(() => Object.values(REQUEST_AUTH_TYPES).reduce((acc, type) => {
      if (props.onlyCredentials && type === REQUEST_AUTH_TYPES.token) {
        return acc;
      }

      acc.push({
        value: type,
        label: t(`request.authTypes.${type}`),
      });

      return acc;
    }, []));

    const isCredentialsType = computed(() => props.value.type === REQUEST_AUTH_TYPES.credentials);
    const isTokenType = computed(() => props.value.type === REQUEST_AUTH_TYPES.token);

    /**
     * Updates the authentication token value and emits the change to parent component
     *
     * @param {Object} authToken - The authentication token object to update
     */
    const updateAuthToken = authToken => emit('update:auth-token', authToken);

    return {
      types,
      isCredentialsType,
      isTokenType,
      updateField,
      updateAuthToken,
    };
  },
};
</script>
