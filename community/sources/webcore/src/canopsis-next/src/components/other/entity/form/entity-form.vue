<template>
  <v-layout class="gap-3" column>
    <c-form-general-patterns-tabs
      :form="form"
      :additional-label="$t('entity.manageInfos')"
    >
      <template #general="{ setRef }">
        <entity-general-form
          v-field="form"
          :ref="setRef"
          :prepare-state-setting-form="prepareStateSettingForm"
        />
      </template>

      <template #additional="{ setRef }">
        <entity-manage-infos-form
          v-field="form.infos"
          :ref="setRef"
        />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import EntityGeneralForm from './entity-general-form.vue';
import EntityManageInfosForm from './entity-manage-infos-form.vue';

export default {
  components: {
    EntityGeneralForm,
    EntityManageInfosForm,
  },
  model: {
    prop: 'form',
    event: 'input',
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
};
</script>
