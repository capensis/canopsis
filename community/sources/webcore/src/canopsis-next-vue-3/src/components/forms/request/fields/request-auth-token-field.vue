<template>
  <c-layout gap="2">
    <div :style="flexSizeStyle">
      <c-select
        :model-value="value.type"
        :label="$t('request.whereToAdd')"
        :name="`${name}.type`"
        :disabled="disabled"
        :required="required"
        :items="types"
        item-title="text"
        item-value="value"
        @update:model-value="updateField('type', $event)"
      />
    </div>
    <div v-if="hasParameterField" :style="flexSizeStyle">
      <c-text-field
        :model-value="value.parameter"
        :label="$t('request.parameterName')"
        :name="`${name}.parameter`"
        :disabled="disabled"
        :required="required"
        @update:model-value="updateField('parameter', $event)"
      />
    </div>
    <div :style="flexSizeStyle">
      <external-auth-token-field
        :value="value.rule"
        :label="$t('request.tokenValue')"
        :name="`${name}.rule`"
        :disabled="disabled"
        :required="required"
        @update:value="updateField('rule', $event)"
      />
    </div>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { REQUEST_AUTH_TOKEN_TYPES } from '@/constants';

import { gridFlexStyle } from '@/helpers/layout';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import ExternalAuthTokenField from '@/components/other/external-auth-token/form/fields/external-auth-token-field.vue';

export default {
  emits: ['update:value'],
  components: { ExternalAuthTokenField },
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
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const types = computed(() => Object.values(REQUEST_AUTH_TOKEN_TYPES).map(type => ({
        value: type,
        text: t(`request.authTokenTypes.${type}`),
      })),
    );

    const isAnyHeaderAuthorizationType = computed(() => [
        REQUEST_AUTH_TOKEN_TYPES.headerAuthorization,
        REQUEST_AUTH_TOKEN_TYPES.headerAuthorizationBearer,
      ].includes(props.value.type),
    );

    const isHeaderCustomParameterType = computed(
      () => props.value.type === REQUEST_AUTH_TOKEN_TYPES.headerCustomParameter,
    );

    const isPayloadType = computed(() => props.value.type === REQUEST_AUTH_TOKEN_TYPES.payload);
    const isUrlType = computed(() => props.value.type === REQUEST_AUTH_TOKEN_TYPES.url);

    const hasParameterField = computed(
      () => isHeaderCustomParameterType.value || isPayloadType.value || isUrlType.value,
    );

    const flexSizeStyle = computed(() => gridFlexStyle(hasParameterField.value ? 4 : 6));

    return {
      types,
      isAnyHeaderAuthorizationType,
      isHeaderCustomParameterType,
      isPayloadType,
      isUrlType,
      hasParameterField,
      flexSizeStyle,
      updateField,
    };
  },
};
</script>
