<template lang="pug">
  v-container(fluid)
    v-data-table(
      :items="statsList",
      :headers="columns",
    )
      template(slot="headers", slot-scope="props")
        tr
          th(v-for="header in props.headers", :key="header.value") {{ header.value }}
      template(slot="items", slot-scope="props")
          tr.text-xs-center
            td(v-for="(property, key) in props.item")
              p(v-if="key === 'entity'") {{ property.name }}
              template(v-else)
                p(v-if="property.value !== null && property.value !== undefined") {{ property.value }}
                p(v-else) No data
</template>

<script>
import entitiesStatsMixin from '@/mixins/entities/stats';

export default {
  mixins: [entitiesStatsMixin],
  computed: {
    columns() {
      const columnsList = [];

      Object.keys(this.statsList[0]).map(item => columnsList.push({ value: item }));

      return columnsList;
    },
  },
  mounted() {
    this.fetchStats();
  },
};
</script>
