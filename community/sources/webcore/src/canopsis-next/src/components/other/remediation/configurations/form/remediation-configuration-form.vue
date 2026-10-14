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
          <c-password-field
            v-model="authToken"
            :label="t('modals.createRemediationConfiguration.fields.token')"
            :required="isNew"
            :masked="!isNew"
            name="token"
            autocomplete="new-password"
            visibility
          />
        </c-form-block-row>

        <c-form-block-row v-if="isShownUserNameField" :label="t('common.username')">
          <c-name-field
            v-model="authUsername"
            :label="t('common.username')"
            name="username"
            required
          />
        </c-form-block-row>

        <template v-if="isShownAwxFields">
          <c-form-block-row
            :label="t('modals.createRemediationConfiguration.fields.enableRefreshToken')"
            align-center
          >
            <c-enabled-field
              v-model="enableRefreshToken"
              :label="t('modals.createRemediationConfiguration.fields.enableRefreshToken')"
              hide-details
              no-margin
            />
          </c-form-block-row>

          <v-expand-transition>
            <div v-if="enableRefreshToken">
              <c-form-block-row :label="t('modals.createRemediationConfiguration.fields.refreshToken')">
                <c-password-field
                  v-model="authRefreshToken"
                  :label="t('modals.createRemediationConfiguration.fields.refreshToken')"
                  :required="areAwxAuthFieldsRequired"
                  :masked="hasExistingClientId"
                  name="refresh_token"
                  autocomplete="new-password"
                  visibility
                />
              </c-form-block-row>

              <c-form-block-row :label="t('modals.createRemediationConfiguration.fields.clientId')">
                <c-text-field
                  v-model="clientId"
                  :label="t('modals.createRemediationConfiguration.fields.clientId')"
                  :required="enableRefreshToken"
                  name="client_id"
                />
              </c-form-block-row>

              <c-form-block-row :label="t('modals.createRemediationConfiguration.fields.clientSecret')" bottom-border>
                <c-password-field
                  v-model="clientSecret"
                  :label="t('modals.createRemediationConfiguration.fields.clientSecret')"
                  :required="areAwxAuthFieldsRequired"
                  :masked="hasExistingClientId"
                  name="client_secret"
                  autocomplete="new-password"
                  visibility
                />
              </c-form-block-row>
            </div>
          </v-expand-transition>
        </template>

        <c-form-block-row :label="t('common.request.skipVerify')" align-center>
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

import {
  isAwxConfigurationType,
  isJobTypeIncludesUserName,
} from '@/helpers/entities/remediation/configuration/form';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';

export default {
  props: {
    isNew: {
      type: Boolean,
      default: true,
    },
  },

  setup() {
    const { t } = useI18n();
    const { remediationJobConfigTypes } = useInfo();
    const {
      name,
      type,
      host,
      auth_token: authToken,
      enable_refresh_token: enableRefreshToken,
      auth_refresh_token: authRefreshToken,
      client_id: clientId,
      client_secret: clientSecret,
      auth_username: authUsername,
      skip_verify: skipVerify,
    } = useFormModels([
      'name',
      'type',
      'host',
      'auth_token',
      'enable_refresh_token',
      'auth_refresh_token',
      'client_id',
      'client_secret',
      'auth_username',
      'skip_verify',
    ]);

    const typeObject = computed(() =>
      remediationJobConfigTypes.value.find(({ name: typeName }) => typeName === type.value),
    );
    const isShownUserNameField = computed(() => isJobTypeIncludesUserName(typeObject.value));
    const isShownAwxFields = computed(() => isAwxConfigurationType(type.value));
    const hasExistingClientId = Boolean(clientId.value);
    const wasRefreshTokenEnabled = enableRefreshToken.value;
    const areAwxAuthFieldsRequired = computed(() =>
      enableRefreshToken.value && !wasRefreshTokenEnabled,
    );

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
      enableRefreshToken,
      authRefreshToken,
      clientId,
      clientSecret,
      authUsername,
      skipVerify,
      isShownUserNameField,
      isShownAwxFields,
      hasExistingClientId,
      areAwxAuthFieldsRequired,
      updateType,
    };
  },
};
</script>
