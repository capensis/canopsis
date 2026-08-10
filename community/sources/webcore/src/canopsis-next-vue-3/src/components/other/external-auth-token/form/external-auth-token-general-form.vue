<template>
  <c-layout column>
    <c-name-field
      v-model="name"
      :label="t('externalAuthToken.tokenName')"
      name="name"
      autofocus
      required
    />
    <c-description-field v-model="description" />

    <c-form-block>
      <request-form
        :form="request"
        :url-variables="templateVars.request"
        :headers-variables="templateVars.request"
        :payload-variables="templateVars.request"
        :url-help-text="t('externalAuthToken.urlHelp')"
        name="request"
        @update:form="request = $event"
      />

      <c-form-block-row :label="t('common.token')" top-border>
        <c-payload-text-field
          v-if="allowVariables"
          :value="template"
          :label="t('common.token')"
          :variables="templateVars.template"
          name="response_field"
          required
          @update:value="template = $event"
        >
          <template #append-inner>
            <c-help-icon
              :text="t('externalAuthToken.tokenExpirationHelpText')"
              icon="help"
              color="grey darken-1"
              top
            />
          </template>
        </c-payload-text-field>

        <c-text-field
          v-else
          v-model="responseField"
          :label="t('common.token')"
          name="response_field"
          required
        >
          <template #append-inner>
            <c-help-icon
              :text="t('externalAuthToken.tokenExpirationHelpText')"
              icon="help"
              color="grey darken-1"
              top
            />
          </template>
        </c-text-field>

        <c-enabled-field v-model="allowVariables" :label="t('externalAuthToken.allowVariables')" />
      </c-form-block-row>

      <c-form-block-row :label="t('externalAuthToken.tokenExpirationTime')">
        <c-duration-field v-model:duration="expirationDuration" required same-width />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import RequestForm from '@/components/forms/request/request-form.vue';
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  components: { RequestForm },

  props: {
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();
    const {
      name,
      description,
      request,
      template,
      responseField,
      allowVariables,
      expirationDuration,
    } = useFormModels({
      name: 'name',
      description: 'description',
      request: 'request',
      template: 'template',
      responseField: 'response_field',
      allowVariables: 'allow_variables',
      expirationDuration: 'expiration_duration',
    });

    return {
      t,
      name,
      description,
      request,
      template,
      responseField,
      allowVariables,
      expirationDuration,
    };
  },
};
</script>
