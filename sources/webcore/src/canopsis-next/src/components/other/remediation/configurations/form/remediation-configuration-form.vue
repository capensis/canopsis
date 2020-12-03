<template lang="pug">
  v-layout(column)
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(row)
      v-flex
        v-select(
          v-field="form.type",
          v-validate="'required'",
          :items="availableTypes",
          :label="$t('common.type')",
          :error-messages="errors.collect('type')",
          name="type",
          item-text="text",
          item-value="value"
        )
      v-flex
        v-text-field(
          v-field="form.host",
          v-validate="'required|url'",
          :label="$t('modals.createRemediationConfiguration.fields.host')",
          :error-messages="errors.collect('host')",
          name="host"
        )
    v-layout(row)
      v-text-field(
        v-field="form.auth_token",
        v-validate="'required'",
        :label="$t('modals.createRemediationConfiguration.fields.token')",
        :error-messages="errors.collect('token')",
        name="token"
      )
</template>

<script>
import { REMEDIATION_CONFIGURATION_TYPES } from '@/constants';

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
  },
  computed: {
    availableTypes() {
      return Object.values(REMEDIATION_CONFIGURATION_TYPES).map(type => ({
        value: type,
        text: this.$t(`modals.createRemediationConfiguration.types.${type}`),
      }));
    },
  },
};
</script>
