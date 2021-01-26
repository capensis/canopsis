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
        item-value="entity._id",
        item-text="v.display_name",
        name="manualMetaAlarm",
        return-object,
        blur-on-create
      )
        template(slot="no-data")
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
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchManualMetaAlarms() {
      this.pending = true;

      const params = {
        manual: true,
        correlation: true,
        page: 1,

        /**
         * We need this option for fetching of every items
         */
        limit: 10000,
      };

      const { data = [] } = await this.fetchAlarmsListWithoutStore({ params });

      this.manualMetaAlarms = data;
      this.pending = false;
    },
  },
};
</script>
