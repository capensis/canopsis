<template>
  <v-layout class="gap-3" column>
    <v-layout class="gap-2">
      <v-flex xs6>
        <c-id-field
          v-field="form._id"
          :disabled="disabledIdField"
          :help-text="$t('metaAlarmRule.idHelp')"
          autofocus
        />
      </v-flex>

      <v-flex xs6>
        <c-name-field
          v-field="form.name"
          :autofocus="disabledIdField"
          required
        />
      </v-flex>
    </v-layout>

    <c-form-block>
      <c-form-block-row :label="$t('metaAlarmRule.outputTemplate')">
        <c-payload-textarea-field
          v-field="form.output_template"
          :label="$t('metaAlarmRule.outputTemplate')"
          :help-text="$t('metaAlarmRule.outputTemplateHelp')"
          :variables="templateVars.output"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$tc('common.tag', 2)" indented>
        <meta-alarm-rule-tags-form
          v-field="form.tags"
          :variables="templateVars.output"
          class="mb-4"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.infos')" indented>
        <meta-alarm-rule-infos-form v-field="form.infos" />
      </c-form-block-row>

      <c-form-block-row :label="$t('metaAlarmRule.autoResolve')">
        <c-enabled-field
          v-field="form.auto_resolve"
          :label="$t('metaAlarmRule.autoResolve')"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('metaAlarmRule.grouping')" indented>
        <meta-alarm-rule-type-field v-field="form.type" />
      </c-form-block-row>

      <template v-if="hasTemplateFields">
        <c-form-block-row :label="$t('metaAlarmRule.componentTemplate')">
          <c-payload-text-field
            v-field="form.config.component_template"
            :label="$t('metaAlarmRule.componentTemplate')"
            :variables="templateVars.entity"
          />
        </c-form-block-row>

        <c-form-block-row :label="$t('metaAlarmRule.resourceTemplate')">
          <c-payload-text-field
            v-field="form.config.resource_template"
            :label="$t('metaAlarmRule.resourceTemplate')"
            :variables="templateVars.entity"
          />
        </c-form-block-row>
      </template>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { META_ALARMS_RULE_TYPES } from '@/constants';

import MetaAlarmRuleTagsForm from './meta-alarm-rule-tags-form.vue';
import MetaAlarmRuleInfosForm from './meta-alarm-rule-infos-form.vue';
import MetaAlarmRuleTypeField from './fields/meta-alarm-rule-type-field.vue';

export default {
  inject: ['$validator'],
  components: { MetaAlarmRuleTagsForm, MetaAlarmRuleInfosForm, MetaAlarmRuleTypeField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    disabledIdField: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const hasTemplateFields = computed(() => [
      META_ALARMS_RULE_TYPES.timebased,
      META_ALARMS_RULE_TYPES.attribute,
      META_ALARMS_RULE_TYPES.complex,
      META_ALARMS_RULE_TYPES.valuegroup,
    ].includes(props.form.type));

    return {
      hasTemplateFields,
    };
  },
};
</script>
