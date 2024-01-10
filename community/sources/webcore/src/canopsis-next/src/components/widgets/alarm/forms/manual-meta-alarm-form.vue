<template>
  <v-layout column>
    <v-combobox
      v-field="form.metaAlarm"
      v-validate="'required'"
      :items="manualMetaAlarms"
      :label="$t('modals.createManualMetaAlarm.fields.metaAlarm')"
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
            <v-list-item-title v-html="$t('modals.createManualMetaAlarm.noData')" />
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-combobox>
    <v-text-field
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
import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('alarm');

export default {
  inject: ['$validator'],
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
      manualMetaAlarms: [],
    };
  },
  mounted() {
    this.fetchManualMetaAlarms();
  },
  methods: {
    ...mapActions({
      fetchManualMetaAlarmsListWithoutStore: 'fetchManualMetaAlarmsListWithoutStore',
    }),

    async fetchManualMetaAlarms() {
      this.pending = true;

      const alarms = await this.fetchManualMetaAlarmsListWithoutStore();

      this.manualMetaAlarms = alarms ?? [];
      this.pending = false;
    },
  },
};
</script>
