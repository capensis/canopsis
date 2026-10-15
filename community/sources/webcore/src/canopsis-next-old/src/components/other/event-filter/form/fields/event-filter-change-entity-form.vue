<template>
  <v-layout class="gap-3" column>
    <span class="text-subtitle-1">
      {{ $t('eventFilter.changeEntityOptions') }}
    </span>

    <c-form-block>
      <c-form-block-row :label="$t('externalData.title')">
        <div class="py-3">
          <external-data-form
            v-field="form.external_data"
            :variables="templateVars.external_data"
            optionally
          />
        </div>
      </c-form-block-row>

      <c-form-block-row :label="$t('common.resource')">
        <c-payload-text-field
          v-field="form.config.resource"
          :label="$t('eventFilter.resource')"
          :name="`${name}.resource`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.component')">
        <c-payload-text-field
          v-field="form.config.component"
          :label="$t('eventFilter.component')"
          :name="`${name}.component`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.connector')">
        <c-payload-text-field
          v-field="form.config.connector"
          :label="$t('eventFilter.connector')"
          :name="`${name}.connector`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.connectorName')">
        <c-payload-text-field
          v-field="form.config.connector_name"
          :label="$t('eventFilter.connectorName')"
          :name="`${name}.connector_name`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.upstream')">
        <c-payload-text-field
          v-field="form.config.upstream"
          :label="$t('eventFilter.upstream')"
          :name="`${name}.upstream`"
          :variables="templateVars.config"
          :required="someRequired"
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

export default {
  components: {
    ExternalDataForm,
  },
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
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const someRequired = computed(() => !(
      props.form.config.resource
       || props.form.config.component
       || props.form.config.connector
       || props.form.config.connector_name
       || props.form.config.upstream
    ));

    return {
      someRequired,
    };
  },
};
</script>
