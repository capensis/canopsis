<template>
  <c-layout gap="2" column>
    <c-enabled-field v-model="enabledField" hide-details with-background />

    <c-form-general-patterns-tabs>
      <template #general="{ setRef }">
        <idle-rule-general-form :ref="setRef" :is-entity-type="isEntityType" />
      </template>

      <template #patterns="{ setRef }">
        <idle-rule-patterns-form :ref="setRef" :is-entity-type="isEntityType" />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { isIdleRuleEntityType } from '@/helpers/entities/idle-rule/form';

import { useFormModels } from '@/hooks/form';

import IdleRuleGeneralForm from './idle-rule-general-form.vue';
import IdleRulePatternsForm from './idle-rule-patterns-form.vue';

export default {
  components: {
    IdleRuleGeneralForm,
    IdleRulePatternsForm,
  },

  setup() {
    const { enabledField, typeField } = useFormModels({
      enabledField: 'enabled',
      typeField: 'type',
    });

    const isEntityType = computed(() => isIdleRuleEntityType(typeField.value));

    return {
      enabledField,
      isEntityType,
    };
  },
};
</script>
