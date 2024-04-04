<template>
  <v-radio-group
    v-field="value"
    :label="$t('metaAlarmRule.selectType')"
    class="my-0"
    hide-details
  >
    <v-radio
      v-for="ruleTypeOption in ruleTypesOptions"
      :key="ruleTypeOption.value"
      :value="ruleTypeOption.value"
      :label="ruleTypeOption.label"
      color="primary"
    >
      <template #label>
        {{ ruleTypeOption.label }}
        <c-help-icon
          v-if="ruleTypeOption.helpText"
          :text="ruleTypeOption.helpText"
          icon="help"
          icon-class="ml-2"
          max-width="300"
          top
        />
      </template>
    </v-radio>
  </v-radio-group>
</template>

<script>
import { computed } from 'vue';

import { META_ALARMS_RULE_TYPES } from '@/constants';

import { isManualGroupMetaAlarmRuleType } from '@/helpers/entities/meta-alarm/rule/form';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    value: {
      type: String,
      required: false,
    },
  },
  setup() {
    const { t, te } = useI18n();

    const ruleTypesOptions = computed(() => Object.values(META_ALARMS_RULE_TYPES).reduce((acc, type) => {
      /**
       * We are filtered 'manualgroup' because we are using in only in the alarms list widget directly
       */
      if (!isManualGroupMetaAlarmRuleType(type)) {
        const messageKey = `metaAlarmRule.types.${type}`;

        const { text, helpText } = te(messageKey) ? t(messageKey) : {};

        acc.push({
          value: type,
          label: text,
          helpText,
        });
      }

      return acc;
    }, []));

    return {
      ruleTypesOptions,
    };
  },
};
</script>
