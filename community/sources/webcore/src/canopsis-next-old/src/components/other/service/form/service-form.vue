<template>
  <v-layout class="gap-3" column>
    <c-enabled-field v-field="form.enabled" with-background />
    <c-form-general-patterns-tabs
      :form="form"
      :additional-label="$t('common.entityPatterns')"
      :patterns-label="$t('entity.manageInfos')"
    >
      <template #general="{ setRef }">
        <service-general-form
          v-field="form"
          :ref="setRef"
          :prepare-state-setting-form="prepareStateSettingForm"
          :template-vars="templateVars"
        />
      </template>

      <template #additional="{ setRef }">
        <service-entity-patterns-form
          v-field="form.patterns"
          :ref="setRef"
        />
      </template>

      <template #patterns="{ setRef }">
        <service-manage-infos-form
          v-field="form.infos"
          :ref="setRef"
        />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import ServiceGeneralForm from './service-general-form.vue';
import ServiceEntityPatternsForm from './service-entity-patterns-form.vue';
import ServiceManageInfosForm from './service-manage-infos-form.vue';

export default {
  components: {
    ServiceGeneralForm,
    ServiceEntityPatternsForm,
    ServiceManageInfosForm,
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
    prepareStateSettingForm: {
      type: Function,
      default: data => data,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
