<template>
  <c-layout gap="3" column>
    <c-name-field v-model="name" autofocus required />

    <c-layout gap="3">
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
      <c-text-field
        v-model="host"
        :label="t('modals.createRemediationConfiguration.fields.host')"
        name="host"
        rules="required|url"
      />
    </c-layout>

    <c-text-field
      v-model="authToken"
      :label="t('modals.createRemediationConfiguration.fields.token')"
      name="token"
      required
    />

    <c-name-field
      v-if="isShownUserNameField"
      v-model="authUsername"
      :label="t('common.username')"
      name="username"
    />

    <c-enabled-field v-model="skipVerify" :label="t('common.request.skipVerify')" />
  </c-layout>
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
