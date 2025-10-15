<template>
  <v-layout column>
    <c-meta-alarm-field
      v-field="form.metaAlarm"
      :label="$t('modals.linkToMetaAlarm.fields.metaAlarm')"
      :autocomplete="false"
      name="metaAlarm"
      required
      blur-on-create
      addable
    >
      <template #no-data="">
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title v-html="$t('modals.linkToMetaAlarm.noData')" />
          </v-list-item-content>
        </v-list-item>
      </template>
    </c-meta-alarm-field>
    <c-name-field
      v-field="form.comment"
      :label="$tc('common.note')"
      :max-length="255"
      name="comment"
    />
    <v-text-field
      v-field="form.component"
      :label="$t('common.component')"
    />
    <v-text-field
      v-field="form.resource"
      :label="$t('common.resource')"
    />
    <meta-alarm-rule-tags-form v-field="form.tags" />
    <meta-alarm-rule-infos-form v-field="form.infos" />
    <c-enabled-field
      v-field="form.auto_resolve"
      :label="$t('metaAlarmRule.autoResolve')"
    />
  </v-layout>
</template>

<script>
import MetaAlarmRuleTagsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-tags-form.vue';
import MetaAlarmRuleInfosForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-infos-form.vue';

export default {
  inject: ['$validator'],
  components: { MetaAlarmRuleTagsForm, MetaAlarmRuleInfosForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
