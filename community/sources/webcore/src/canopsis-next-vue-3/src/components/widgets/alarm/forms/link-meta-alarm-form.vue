<template>
  <c-layout gap="2" column>
    <c-meta-alarm-field
      :value="form.metaAlarm"
      :label="t('modals.linkToMetaAlarm.fields.metaAlarm')"
      :autocomplete="false"
      name="metaAlarm"
      required
      addable
      @update:value="updateField('metaAlarm', $event)"
    >
      <template #no-data>
        <v-list-item>
          <v-list-item-title>{{ t('modals.linkToMetaAlarm.noData') }}</v-list-item-title>
        </v-list-item>
      </template>
    </c-meta-alarm-field>
    <c-name-field
      :model-value="form.comment"
      :label="t('common.note')"
      :max-length="255"
      name="comment"
      @update:model-value="updateField('comment', $event)"
    />
    <c-text-field
      :model-value="form.component"
      :label="t('common.component')"
      @update:model-value="updateField('component', $event)"
    />
    <c-text-field
      :model-value="form.resource"
      :label="t('common.resource')"
      @update:model-value="updateField('resource', $event)"
    />
    <meta-alarm-rule-tags-form
      :form="form.tags"
      @update:form="updateField('tags', $event)"
    />
    <meta-alarm-rule-infos-form
      :infos="form.infos"
      @update:infos="updateField('infos', $event)"
    />
    <c-enabled-field
      :model-value="form.auto_resolve"
      :label="t('metaAlarmRule.autoResolve')"
      @update:model-value="updateField('auto_resolve', $event)"
    />
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import MetaAlarmRuleTagsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-tags-form.vue';
import MetaAlarmRuleInfosForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-infos-form.vue';

export default {
  components: { MetaAlarmRuleTagsForm, MetaAlarmRuleInfosForm },

  props: {
    form: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    return {
      t,
      updateField,
    };
  },
};
</script>
