<template lang="pug">
  v-container(fluid)
    v-btn(icon, @click="showSettings")
      v-icon settings
    v-data-table(
      :items="statsList",
      :headers="columns",
    )
      template(slot="headers", slot-scope="props")
        tr
          th Entity
          th(v-for="header in props.headers", :key="header.value") {{ header.value }}
      template(slot="items", slot-scope="props")
          tr.text-xs-center
            td {{ props.item.entity.name }}
            td(v-for="(property, key) in widget.parameters.stats")
              template(
              v-if="props.item[key].value !== undefined && props.item[key].value !== null"
              )
                div {{ props.item[key].value }}
                  sub {{ props.item[key].trend }}
              div(v-else) No data
</template>

<script>
import { SIDE_BARS } from '@/constants';
import entitiesStatsMixin from '@/mixins/entities/stats';
import sideBarMixin from '@/mixins/side-bar/side-bar';

export default {
  mixins: [entitiesStatsMixin, sideBarMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    columns() {
      const columnsList = [];
      Object.keys(this.widget.parameters.stats).map(item => columnsList.push({ value: item }));

      return columnsList;
    },

  },
  mounted() {
    this.fetchStats({ params: this.widget.parameters });
  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.statsTableSettings,
        config: {
          widget: this.widget,
        },
      });
    },
  },
};
</script>
