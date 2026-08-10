<template>
  <c-layout gap="2" column>
    <c-enabled-field v-model="enabledField" hide-details with-background />

    <v-tabs v-model="activeTab" align-tabs="center" color="primary">
      <v-tab value="general">
        {{ t('common.general') }}
      </v-tab>
      <v-tab value="patterns">
        {{ t('common.pattern', 2) }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item value="general">
        <idle-rule-general-form :is-entity-type="isEntityType" />
      </v-window-item>

      <v-window-item value="patterns">
        <idle-rule-patterns-form
          :form="patternsField"
          :is-entity-type="isEntityType"
          class="mt-2"
          @update:model-value="patternsField = $event"
        />
      </v-window-item>
    </v-window>
  </c-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { isIdleRuleEntityType } from '@/helpers/entities/idle-rule/form';
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import IdleRuleGeneralForm from './idle-rule-general-form.vue';
import IdleRulePatternsForm from './idle-rule-patterns-form.vue';

export default {
  components: {
    IdleRuleGeneralForm,
    IdleRulePatternsForm,
  },

  setup() {
    const { t } = useI18n();
    const activeTab = ref('general');

    const { enabledField, typeField, patternsField } = useFormModels({
      enabledField: 'enabled',
      typeField: 'type',
      patternsField: 'patterns',
    });

    const isEntityType = computed(() => isIdleRuleEntityType(typeField.value));

    return {
      t,
      activeTab,
      enabledField,
      patternsField,
      isEntityType,
    };
  },
};
</script>
