<template>
  <c-layout gap="3" column>
    <c-layout gap="2">
      <div class="flex-grow-1">
        <c-id-field
          v-model="idField"
          :disabled="disabledIdField"
          :help-text="t('metaAlarmRule.idHelp')"
          autofocus
        />
      </div>

      <div class="flex-grow-1">
        <c-name-field
          v-model="nameField"
          :autofocus="disabledIdField"
          required
        />
      </div>
    </c-layout>

    <c-form-block>
      <c-form-block-row :label="t('metaAlarmRule.outputTemplate')" indented>
        <c-payload-textarea-field
          v-model:value="outputTemplateField"
          :label="t('metaAlarmRule.outputTemplate')"
          :help-text="t('metaAlarmRule.outputTemplateHelp')"
          :variables="templateVars.output"
          name="output_template"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.tag', 2)" indented>
        <meta-alarm-rule-tags-form
          v-model:form="tagsField"
          :variables="templateVars.output"
          class="mb-4"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.infos')" indented>
        <meta-alarm-rule-infos-form v-model:infos="infosField" />
      </c-form-block-row>

      <c-form-block-row :label="t('metaAlarmRule.autoResolve')" align-center>
        <c-enabled-field
          v-model="autoResolveField"
          :label="t('metaAlarmRule.autoResolve')"
          hide-details
        />
      </c-form-block-row>

      <c-form-block-row :label="t('metaAlarmRule.grouping')" indented>
        <meta-alarm-rule-type-field v-model="typeField" />
      </c-form-block-row>

      <template v-if="hasTemplateFields">
        <c-form-block-row :label="t('metaAlarmRule.componentTemplate')" align-center>
          <c-payload-text-field
            v-model:value="componentTemplateField"
            :label="t('metaAlarmRule.componentTemplate')"
            :variables="templateVars.entity"
            name="component_template"
            hide-details
          />
        </c-form-block-row>

        <c-form-block-row :label="t('metaAlarmRule.resourceTemplate')" align-center>
          <c-payload-text-field
            v-model:value="resourceTemplateField"
            :label="t('metaAlarmRule.resourceTemplate')"
            :variables="templateVars.entity"
            name="resource_template"
            hide-details
          />
        </c-form-block-row>
      </template>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { hasTemplateConfigFields } from '@/helpers/entities/meta-alarm/rule/form';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import MetaAlarmRuleTypeField from './fields/meta-alarm-rule-type-field.vue';
import MetaAlarmRuleInfosForm from './meta-alarm-rule-infos-form.vue';
import MetaAlarmRuleTagsForm from './meta-alarm-rule-tags-form.vue';

export default {
  components: {
    MetaAlarmRuleInfosForm,
    MetaAlarmRuleTagsForm,
    MetaAlarmRuleTypeField,
  },

  props: {
    disabledIdField: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();
    const {
      idField,
      nameField,
      outputTemplateField,
      tagsField,
      infosField,
      autoResolveField,
      typeField,
      componentTemplateField,
      resourceTemplateField,
    } = useFormModels({
      idField: '_id',
      nameField: 'name',
      outputTemplateField: 'output_template',
      tagsField: 'tags',
      infosField: 'infos',
      autoResolveField: 'auto_resolve',
      typeField: 'type',
      componentTemplateField: 'config.component_template',
      resourceTemplateField: 'config.resource_template',
    });

    const hasTemplateFields = computed(() => hasTemplateConfigFields(typeField.value));

    return {
      t,
      idField,
      nameField,
      outputTemplateField,
      tagsField,
      infosField,
      autoResolveField,
      typeField,
      componentTemplateField,
      resourceTemplateField,
      hasTemplateFields,
    };
  },
};
</script>
