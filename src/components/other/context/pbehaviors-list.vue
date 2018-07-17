<template lang="pug">
v-data-table(:items="items", :headers="pbehaviorsTableHeaders")
  template(slot="items" slot-scope="props")
    td {{ props.item.name }}
    td {{ props.item.author }}
    td {{ props.item.connector }}
    td {{ props.item.connector_name }}
    td {{ props.item.enabled }}
    td {{ props.item.tstart }}
    td {{ props.item.tstop }}
    td {{ props.item.type }}
    td {{ props.item.reason }}
    td {{ props.item.rrule }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions: pbehaviorMapAction, mapGetters: pbehaviorMapgetters } = createNamespacedHelpers('pbehavior');

export default {
  props: {
    itemId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      items: [],
      pbehaviorsTableHeaders: [
        {
          text: this.$t('common.name'),
          sortable: false,
        },
        {
          text: this.$t('common.author'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviorsList.connector'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviorsList.connectorName'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviorsList.isEnabled'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviorsList.begins'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviorsList.ends'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviorsList.type'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviorsList.reason'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviorsList.rrule'),
          sortable: false,
        },
      ],
    };
  },
  computed: {
    ...pbehaviorMapgetters(['pbehaviorsList']),
  },
  mounted() {
    this.fetchItems();
  },
  methods: {
    ...pbehaviorMapAction({
      fetchPbehaviorsList: 'fetchById',
    }),
    async fetchItems() {
      await this.fetchPbehaviorsList({ id: this.itemId });
      this.items = [...this.pbehaviorsList];
    },
  },
};
</script>

