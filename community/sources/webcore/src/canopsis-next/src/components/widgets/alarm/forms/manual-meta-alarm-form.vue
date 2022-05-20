<template lang="pug">
  div
    v-layout(row)
      v-combobox(
        v-model="form.metaAlarm",
        v-validate="'required'",
        :items="manualMetaAlarms",
        :label="$t('modals.createManualMetaAlarm.fields.metaAlarm')",
        :error-messages="errors.collect('manualMetaAlarm')",
        :loading="pending",
        item-value="_id",
        item-text="name",
        name="manualMetaAlarm",
        return-object,
        blur-on-create
      )
        template(#no-data="")
          v-list-tile
            v-list-tile-content
              v-list-tile-title(v-html="$t('modals.createManualMetaAlarm.noData')")
    v-layout(row)
      v-text-field(
        v-model="form.output",
        :label="$t('modals.createManualMetaAlarm.fields.output')"
      )
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
      fetchManualMetaAlarmListWithoutStore: 'fetchManualMetaAlarmListWithoutStore',
    }),

    async fetchManualMetaAlarms() {
      this.pending = true;

      const alarms = await this.fetchManualMetaAlarmListWithoutStore();

      this.manualMetaAlarms = alarms ?? [];
      this.pending = false;
    },
  },
};
</script>
