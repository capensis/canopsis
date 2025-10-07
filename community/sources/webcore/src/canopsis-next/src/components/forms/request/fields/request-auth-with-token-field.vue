<template>
  <v-layout class="gap-2" column>
    <v-radio-group
      v-field="value.type"
      :disabled="disabled"
      column
    >
      <v-radio
        v-for="type in types"
        :key="type.value"
        :label="type.label"
        :value="type.value"
      />
    </v-radio-group>
    <v-layout v-if="isCredentialsType" class="gap-2">
      <c-name-field
        v-field="value.username"
        :label="$t('common.username')"
        :name="`${name}.username`"
        :disabled="disabled"
      />
      <c-password-field
        v-field="value.password"
        :name="`${name}.password`"
        :disabled="disabled"
      />
    </v-layout>
    <request-auth-token-field
      v-if="isTokenType"
      :value="authToken"
      :name="`${name}.token`"
      :disabled="disabled"
      @input="updateAuthToken"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { REQUEST_AUTH_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import RequestAuthTokenField from '@/components/forms/request/fields/request-auth-token-field.vue';

export default {
  components: { RequestAuthTokenField },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    authToken: {
      type: Object,
      default: () => ({}),
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

    const types = computed(() => Object.values(REQUEST_AUTH_TYPES).map(type => ({
      value: type,
      label: t(`request.authTypes.${type}`),
    })));

    const isCredentialsType = computed(() => props.value.type === REQUEST_AUTH_TYPES.credentials);
    const isTokenType = computed(() => props.value.type === REQUEST_AUTH_TYPES.token);

    const updateAuthToken = authToken => emit('update:auth-token', authToken);

    return {
      types,
      isCredentialsType,
      isTokenType,
      updateAuthToken,
    };
  },
};
</script>
