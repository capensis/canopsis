<template>
  <c-layout gap="3" column>
    <span class="text-subtitle-1">
      {{ $t('eventFilter.changeEntityOptions') }}
    </span>

    <c-form-block>
      <c-form-block-row :label="$t('externalData.title')" indented>
        <external-data-field
          v-model:external-data="externalDataField"
          :variables="templateVars.external_data"
          optionally
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.resource')">
        <c-payload-text-field
          v-model:value="configField.resource"
          :label="$t('eventFilter.resource')"
          :name="`${name}.resource`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.component')">
        <c-payload-text-field
          v-model:value="configField.component"
          :label="$t('eventFilter.component')"
          :name="`${name}.component`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.connector')">
        <c-payload-text-field
          v-model:value="configField.connector"
          :label="$t('eventFilter.connector')"
          :name="`${name}.connector`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.connectorName')">
        <c-payload-text-field
          v-model:value="configField.connector_name"
          :label="$t('eventFilter.connectorName')"
          :name="`${name}.connector_name`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.upstream')">
        <c-payload-text-field
          v-model:value="configField.upstream"
          :label="$t('eventFilter.upstream')"
          :name="`${name}.upstream`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import ExternalDataField from '@/components/forms/external-data/external-data-field.vue';
import { useFormModels } from '@/hooks/form';

export default {
  components: {
    ExternalDataField,
  },
  props: {
    name: {
      type: String,
      default: 'config',
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup() {
    const { externalDataField, configField } = useFormModels({
      externalDataField: 'external_data',
      configField: 'config',
    });

    const someRequired = computed(
      () =>
        !(
          configField.value?.resource ||
          configField.value?.component ||
          configField.value?.connector ||
          configField.value?.connector_name ||
          configField.value?.upstream
        ),
    );

    return {
      someRequired,

      externalDataField,
      configField,
    };
  },
};
</script>
