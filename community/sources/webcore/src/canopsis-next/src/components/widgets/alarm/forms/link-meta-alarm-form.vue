<template>
  <v-layout column>
    <v-combobox
      v-field="form.metaAlarm"
      v-validate="'required'"
      :items="metaAlarms"
      :label="$t('modals.linkToMetaAlarm.fields.metaAlarm')"
      :error-messages="errors.collect('manualMetaAlarm')"
      :loading="pending"
      item-value="_id"
      item-text="name"
      name="manualMetaAlarm"
      return-object
      blur-on-create
    >
      <template #no-data="">
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title v-html="$t('modals.linkToMetaAlarm.noData')" />
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-combobox>
    <v-text-field
      v-field="form.comment"
      :label="$t('common.note')"
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
import { createNamespacedHelpers } from 'vuex';

import MetaAlarmRuleTagsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-tags-form.vue';
import MetaAlarmRuleInfosForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-infos-form.vue';

const { mapActions } = createNamespacedHelpers('metaAlarm');

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
  data() {
    return {
      pending: false,
      metaAlarms: [],
    };
  },
  mounted() {
    this.fetchMetaAlarms();
  },
  methods: {
    ...mapActions({
      fetchMetaAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchMetaAlarms() {
      this.pending = true;

      const alarms = await this.fetchMetaAlarmsListWithoutStore();

      this.metaAlarms = alarms ?? [];
      this.pending = false;
    },
  },
};
</script>
