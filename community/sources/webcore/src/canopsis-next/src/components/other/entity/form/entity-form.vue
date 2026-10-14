<template>
  <c-layout gap="3" column>
    <c-form-general-patterns-tabs
      :form="form"
      :additional-label="$t('entity.manageInfos')"
    >
      <template #general="{ setRef }">
        <entity-general-form
          :ref="setRef"
          :form="form"
          :prepare-state-setting-form="prepareStateSettingForm"
          @update:form="updateModel"
        />
      </template>

      <template #additional="{ setRef }">
        <entity-manage-infos-form
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

import EntityGeneralForm from './entity-general-form.vue';
import EntityManageInfosForm from './entity-manage-infos-form.vue';

export default {
  components: {
    EntityGeneralForm,
    EntityManageInfosForm,
  },

  props: {
    form: {
      type: Object,
      required: true,
    },
    prepareStateSettingForm: {
      type: Function,
      default: data => data,
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
