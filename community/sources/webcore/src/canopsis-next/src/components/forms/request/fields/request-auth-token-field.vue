<template>
  <v-layout class="gap-2">
    <v-flex v-bind="flexSizeProp">
      <c-select-field
        v-field="value.type"
        :label="$t('request.whereToAdd')"
        :name="`${name}.type`"
        :disabled="disabled"
        :required="required"
        :items="types"
      />
    </v-flex>
    <v-flex v-if="hasParameterField" v-bind="flexSizeProp">
      <v-text-field
        v-field="value.parameter"
        :label="$t('request.parameterName')"
        :name="`${name}.parameter`"
        :disabled="disabled"
        :required="required"
      />
    </v-flex>
    <v-flex v-bind="flexSizeProp">
      <external-auth-token-field
        v-field="value.rule"
        :label="$t('request.tokenValue')"
        :name="`${name}.rule`"
        :disabled="disabled"
        :required="required"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { REQUEST_AUTH_TOKEN_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import ExternalAuthTokenField from '@/components/other/external-auth-token/form/fields/external-auth-token-field.vue';

export default {
  components: { ExternalAuthTokenField },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'token',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const types = computed(() => Object.values(REQUEST_AUTH_TOKEN_TYPES).map(type => ({
      value: type,
      text: t(`request.authTokenTypes.${type}`),
    })));

    const isAnyHeaderAuthorizationType = computed(() => [
      REQUEST_AUTH_TOKEN_TYPES.headerAuthorization,
      REQUEST_AUTH_TOKEN_TYPES.headerAuthorizationBearer,
    ].includes(props.value.type));

    const isHeaderCustomParameterType = computed(() => (
      props.value.type === REQUEST_AUTH_TOKEN_TYPES.headerCustomParameter
    ));

    const isPayloadType = computed(() => props.value.type === REQUEST_AUTH_TOKEN_TYPES.payload);
    const isUrlType = computed(() => props.value.type === REQUEST_AUTH_TOKEN_TYPES.url);

    const hasParameterField = computed(() => (
      isHeaderCustomParameterType.value || isPayloadType.value || isUrlType.value
    ));

    const flexSizeProp = computed(() => {
      if (hasParameterField.value) {
        return { xs4: true };
      }
      return { xs6: true };
    });

    return {
      types,
      isAnyHeaderAuthorizationType,
      isHeaderCustomParameterType,
      isPayloadType,
      isUrlType,
      hasParameterField,
      flexSizeProp,
    };
  },
};
</script>
