<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name")
    v-layout(row)
      v-flex
        v-select(
          v-field="form.type",
          v-validate="'required'",
          :items="remediationJobConfigTypes",
          :label="$t('common.type')",
          :error-messages="errors.collect('type')",
          name="type",
          item-text="name",
          item-value="name",
          return-object
        )
      v-flex
        v-text-field(
          v-field="form.host",
          v-validate="'required|url'",
          :label="$t('modals.createRemediationConfiguration.fields.host')",
          :error-messages="errors.collect('host')",
          name="host"
        )
    v-text-field(
      v-field="form.auth_token",
      v-validate="'required'",
      :label="$t('modals.createRemediationConfiguration.fields.token')",
      :error-messages="errors.collect('token')",
      name="token"
    )
    v-text-field(
      v-if="isShownUserNameField",
      v-field="form.auth_username",
      v-validate="'required'",
      :label="$t('users.username')",
      :error-messages="errors.collect('username')",
      name="username"
    )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { get, isString } from 'lodash';

import { REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME } from '@/constants';

import { formMixin } from '@/mixins/form';

const { mapGetters } = createNamespacedHelpers('info');

export default {
  inject: ['$validator'],
  mixins: [formMixin],
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
    ...mapGetters(['remediationJobConfigTypes']),

    isShownUserNameField() {
      return REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME.includes(get(this.form.type, 'auth_type'));
    },
  },

  mounted() {
    if (isString(this.form.type)) {
      const typeObject = this.remediationJobConfigTypes.find(({ name }) => name === this.form.type);

      if (typeObject) {
        this.updateField('type', typeObject);
      }
    }
  },
};
</script>
