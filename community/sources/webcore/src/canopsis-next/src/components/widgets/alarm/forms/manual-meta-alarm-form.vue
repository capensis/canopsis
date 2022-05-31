<template lang="pug">
  v-layout(column)
    v-combobox(
      v-field="form.metaAlarm",
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
      template(#no-data="")
        v-list-tile
          v-list-tile-content
            v-list-tile-title(v-html="$t('modals.createManualMetaAlarm.noData')")
    v-text-field(
      v-field="form.output",
      :label="$t('common.note')"
    )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

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
        limit: MAX_LIMIT,
      };

      const { data = [] } = await this.fetchAlarmsListWithoutStore({ params });

      this.manualMetaAlarms = data;
      this.pending = false;
    },
  },
};
</script>
