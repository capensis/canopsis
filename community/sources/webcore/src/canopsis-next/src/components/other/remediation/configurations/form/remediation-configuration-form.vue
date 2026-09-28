<template>
  <v-layout class="gap-3" column>
    <c-name-field
      v-field="form.name"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="$t('common.type')">
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
      </c-form-block-row>

      <c-form-block-row :label="$t('modals.createRemediationConfiguration.fields.host')">
        <v-text-field
          v-field="form.host"
          v-validate="'required|url'"
          :label="$t('modals.createRemediationConfiguration.fields.host')"
          :error-messages="errors.collect('host')"
          name="host"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('modals.createRemediationConfiguration.fields.token')">
        <v-text-field
          v-field="form.auth_token"
          v-validate="'required'"
          :label="$t('modals.createRemediationConfiguration.fields.token')"
          :error-messages="errors.collect('token')"
          name="token"
        />
      </c-form-block-row>

      <c-form-block-row
        v-if="isShownUserNameField"
        :label="$t('common.username')"
      >
        <c-name-field
          v-field="form.auth_username"
          :label="$t('common.username')"
          name="username"
        />
      </c-form-block-row>

      <c-form-block-row
        :label="$t('common.request.skipVerify')"
        align-center
      >
        <c-enabled-field
          v-field="form.skip_verify"
          :label="$t('common.request.skipVerify')"
          hide-details
          no-margin
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { isJobTypeIncludesUserName } from '@/helpers/entities/remediation/configuration/form';

import { useInfo } from '@/hooks/store/modules/info';
import { useModelField } from '@/hooks/form/model-field';

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
  setup(props, { emit }) {
    const { remediationJobConfigTypes } = useInfo();
    const { updateModel } = useModelField(props, emit);

    const typeObject = computed(() => remediationJobConfigTypes.value.find(
      ({ name }) => name === props.form.type,
    ));

    const isShownUserNameField = computed(() => isJobTypeIncludesUserName(typeObject.value));

    const updateType = (type) => {
      const hasUserName = isJobTypeIncludesUserName(type);

      updateModel({
        ...props.form,
        type: type.name,
        auth_username: hasUserName ? props.form.auth_username : '',
      });
    };

    return {
      remediationJobConfigTypes,
      isShownUserNameField,
      updateType,
    };
  },
};
</script>
