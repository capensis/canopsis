<template>
  <v-layout column>
    <c-meta-alarm-field
      v-field="form.metaAlarm"
      :label="$t('modals.linkToMetaAlarm.fields.metaAlarm')"
      :autocomplete="false"
      name="manualMetaAlarm"
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
      :label="$t('common.note')"
    />
    <c-enabled-field
      v-field="form.auto_resolve"
      :label="$t('metaAlarmRule.autoResolve')"
    />
  </v-layout>
</template>

<script>
export default {
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
