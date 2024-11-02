<template>
  <v-layout class="gap-4" column>
    <v-radio-group
      v-field="form.type"
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
    <v-expand-transition>
      <v-layout v-if="hasTemplateFields" column>
        <c-payload-text-field
          v-field="form.config.component_template"
          :label="$t('metaAlarmRule.componentTemplate')"
          :variables="variables"
        />
        <c-payload-text-field
          v-field="form.config.resource_template"
          :label="$t('metaAlarmRule.resourceTemplate')"
          :variables="variables"
        />
      </v-layout>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { META_ALARMS_RULE_TYPES } from '@/constants';

import { isManualGroupMetaAlarmRuleType } from '@/helpers/entities/meta-alarm/rule/form';

import { useI18n } from '@/hooks/i18n';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const { t, te } = useI18n();

    const ruleTypesOptions = computed(() => Object.values(META_ALARMS_RULE_TYPES).reduce((acc, type) => {
      /**
       * We are filtered 'manualgroup' because we are using it only in the alarms list widget directly
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

    const hasTemplateFields = computed(() => [
      META_ALARMS_RULE_TYPES.timebased,
      META_ALARMS_RULE_TYPES.attribute,
      META_ALARMS_RULE_TYPES.complex,
      META_ALARMS_RULE_TYPES.valuegroup,
    ].includes(props.form.type));

    return {
      ruleTypesOptions,
      hasTemplateFields,
    };
  },
};
</script>
