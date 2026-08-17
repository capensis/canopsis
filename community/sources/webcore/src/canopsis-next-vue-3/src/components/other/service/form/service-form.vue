<template>
  <c-layout gap="3" column>
    <c-enabled-field
      :model-value="form.enabled"
      with-background
      @update:model-value="updateField('enabled', $event)"
    />
    <c-form-general-patterns-tabs
      :form="form"
      :additional-label="$t('common.entityPatterns')"
      :patterns-label="$t('entity.manageInfos')"
    >
      <template #general="{ setRef }">
        <service-general-form
          :ref="setRef"
          :form="form"
          :prepare-state-setting-form="prepareStateSettingForm"
          :template-vars="templateVars"
          @update:form="updateModel"
        />
      </template>

      <template #additional="{ setRef }">
        <service-entity-patterns-form
          :ref="setRef"
          :patterns="form.patterns"
          @update:patterns="updateField('patterns', $event)"
        />
      </template>

      <template #patterns="{ setRef }">
        <service-manage-infos-form
          :ref="setRef"
          :infos="form.infos"
          @update:infos="updateField('infos', $event)"
        />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';

import ServiceEntityPatternsForm from './service-entity-patterns-form.vue';
import ServiceGeneralForm from './service-general-form.vue';
import ServiceManageInfosForm from './service-manage-infos-form.vue';

export default {
  components: {
    ServiceGeneralForm,
    ServiceEntityPatternsForm,
    ServiceManageInfosForm,
  },

  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    prepareStateSettingForm: {
      type: Function,
      default: data => data,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);

    return {
      updateModel,
      updateField,
    };
  },
};
</script>
