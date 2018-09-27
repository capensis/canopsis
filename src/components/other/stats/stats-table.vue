<template lang="pug">
  v-container(fluid)
    v-layout(justify-end)
      v-btn(icon, @click="showSettings")
        v-icon settings
    v-data-table(
      :items="stats",
      :headers="columns",
      :loading="statsPending",
    )
      v-progress-linear(slot="progress", color="blue", indeterminate)
      template(slot="headers", slot-scope="{ headers }")
        th Entity
        th(v-for="header in headers", :key="header.value") {{ header.value }}
      template(slot="items", slot-scope="{ item }")
        td {{ item.entity.name }}
        td(v-for="(property, key) in widget.parameters.stats")
          template(
          v-if="item[key].value !== undefined && item[key].value !== null"
          )
            td
              stats-number(:item="item[key]")
          div(v-else) No data
</template>

<script>
import { SIDE_BARS } from '@/constants';
import entitiesStatsMixin from '@/mixins/entities/stats';
import sideBarMixin from '@/mixins/side-bar/side-bar';

import StatsNumber from './stats-number.vue';

export default {
  components: {
    StatsNumber,
  },
  mixins: [entitiesStatsMixin, sideBarMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    columns() {
      return Object.keys(this.widget.parameters.stats).map(item => ({ value: item }));
    },

  },
  mounted() {
    this.fetchStats({ params: this.widget.parameters, widgetId: this.widget._id });
  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.statsTableSettings,
        config: {
          widget: this.widget,
          rowId: this.rowId,
        },
      });
    },
  },
};
</script>
