<template lang="pug">
  v-container(fluid)
    v-data-table(
      :items="statsList",
      :headers="headers",
    )
      template(slot="items", slot-scope="props")
          tr
            td {{ props.item.entity.name }}
            td {{ props.item["Durée d'indisponibilité"].value }}
</template>

<script>
import entitiesStatsMixin from '@/mixins/entities/stats';

export default {
  mixins: [entitiesStatsMixin],
  data() {
    return {
      headers: [
        {
          text: 'Entity',
          value: 'entity',
        },
        {
          text: 'Durée d\'indisponibilité',
          value: 'Durée d\'indisponibilité',
        },
      ],
    };
  },
  mounted() {
    this.fetchStats();
  },
};
</script>
