<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <v-layout>
      <v-flex
        class="pr-3"
        xs6
      >
        <v-select
          v-validate="'required'"
          :value="form.type"
          :items="remediationJobConfigTypes"
          :label="$t('common.type')"
          :error-messages="errors.collect('type')"
          name="type"
          item-text="name"
          item-value="name"
          return-object
          @input="updateType"
        />
      </v-flex>
      <v-flex xs6>
        <v-text-field
          v-field="form.host"
          v-validate="'required|url'"
          :label="$t('modals.createRemediationConfiguration.fields.host')"
          :error-messages="errors.collect('host')"
          name="host"
        />
      </v-flex>
    </v-layout>
    <v-text-field
      v-field="form.auth_token"
      v-validate="'required'"
      :label="$t('modals.createRemediationConfiguration.fields.token')"
      :error-messages="errors.collect('token')"
      name="token"
    />
    <c-name-field
      v-if="isShownUserNameField"
      v-field="form.auth_username"
      :label="$t('common.username')"
      name="username"
    />
    <c-enabled-field
      v-field="form.skip_verify"
      :label="$t('common.request.skipVerify')"
    />
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { isJobTypeIncludesUserName } from '@/helpers/entities/remediation/configuration/form';

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

    typeObject() {
      return this.remediationJobConfigTypes.find(({ name }) => name === this.form.type);
    },

    isShownUserNameField() {
      return isJobTypeIncludesUserName(this.typeObject);
    },
  },

  methods: {
    updateType(type) {
      const hasUserName = isJobTypeIncludesUserName(type);

      this.updateModel({
        ...this.form,
        type: type.name,
        auth_username: hasUserName ? this.form.auth_username : '',
      });
    },
  },
};
</script>
