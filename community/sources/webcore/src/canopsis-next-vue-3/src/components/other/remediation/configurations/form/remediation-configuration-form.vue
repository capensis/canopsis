<template>
  <v-row>
    <v-col cols="12">
      <c-name-field v-model="name" autofocus required />
    </v-col>

    <v-col cols="12">
      <c-form-block>
        <c-form-block-row :label="t('common.type')">
          <c-select
            v-model="type"
            :items="remediationJobConfigTypes"
            :label="t('common.type')"
            name="type"
            item-title="name"
            item-value="name"
            return-object
            required
            @update:model-value="updateType"
          />
        </c-form-block-row>

        <c-form-block-row :label="t('modals.createRemediationConfiguration.fields.host')">
          <c-text-field
            v-model="host"
            :label="t('modals.createRemediationConfiguration.fields.host')"
            name="host"
            required
            url
          />
        </c-form-block-row>

        <c-form-block-row :label="t('modals.createRemediationConfiguration.fields.token')">
          <c-text-field
            v-model="authToken"
            :label="t('modals.createRemediationConfiguration.fields.token')"
            name="token"
            required
          />
        </c-form-block-row>

        <c-form-block-row
          v-if="isShownUserNameField"
          :label="t('common.username')"
        >
          <c-name-field
            v-model="authUsername"
            :label="t('common.username')"
            name="username"
          />
        </c-form-block-row>

        <c-form-block-row
          :label="t('common.request.skipVerify')"
          align-center
        >
          <c-enabled-field
            v-model="skipVerify"
            :label="t('common.request.skipVerify')"
            hide-details
            no-margin
          />
        </c-form-block-row>
      </c-form-block>
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { isJobTypeIncludesUserName } from '@/helpers/entities/remediation/configuration/form';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';

export default {
  setup() {
    const { t } = useI18n();
    const { remediationJobConfigTypes } = useInfo();
    const {
      name,
      type,
      host,
      auth_token: authToken,
      auth_username: authUsername,
      skip_verify: skipVerify,
    } = useFormModels(['name', 'type', 'host', 'auth_token', 'auth_username', 'skip_verify']);

    const typeObject = computed(
      () => remediationJobConfigTypes.value.find(({ name: typeName }) => typeName === type.value),
    );
    const isShownUserNameField = computed(() => isJobTypeIncludesUserName(typeObject.value));

    const updateType = (nextType) => {
      const hasUserName = isJobTypeIncludesUserName(nextType);

      type.value = nextType?.name ?? '';
      authUsername.value = hasUserName ? authUsername.value : '';
    };

    return {
      t,
      remediationJobConfigTypes,
      name,
      type,
      host,
      authToken,
      authUsername,
      skipVerify,
      isShownUserNameField,
      updateType,
    };
  },
};
</script>
