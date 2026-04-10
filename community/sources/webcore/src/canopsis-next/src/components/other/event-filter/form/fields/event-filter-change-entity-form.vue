<template>
  <div>
    <c-payload-text-field
      v-field="form.resource"
      :label="$t('eventFilter.resource')"
      :name="`${name}.resource`"
      :variables="variables"
    />
    <c-payload-text-field
      v-field="form.component"
      :label="$t('eventFilter.component')"
      :name="`${name}.component`"
      :variables="variables"
    />
    <c-payload-text-field
      v-field="form.connector"
      :label="$t('eventFilter.connector')"
      :name="`${name}.connector`"
      :variables="variables"
    />
    <c-payload-text-field
      v-field="form.connector_name"
      :label="$t('eventFilter.connectorName')"
      :name="`${name}.connector_name`"
      :variables="variables"
    />
    <c-payload-text-field
      v-field="form.upstream"
      :label="$t('eventFilter.upstream')"
      :name="`${name}.upstream`"
      :variables="variables"
    />
    <v-alert
      :value="errors.has(name)"
      type="error"
    >
      {{ $t('eventFilter.configRequired') }}
    </v-alert>
  </div>
</template>

<script>
import { onBeforeUnmount, watch } from 'vue';

import { useModelField } from '@/hooks/form/model-field';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'config',
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    useModelField(props, emit);

    const {
      attachRequiredRule,
      detachRequiredRule,
      validateRequiredRule,
    } = useValidationAttachRequired(props.name);

    const getRequiredValue = () => (
      props.form.resource
       || props.form.component
       || props.form.connector
       || props.form.connector_name
       || props.form.upstream
    );

    watch(() => props.form, validateRequiredRule);

    attachRequiredRule(getRequiredValue);

    onBeforeUnmount(detachRequiredRule);
  },
};
</script>
