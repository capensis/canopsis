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
// div.red.darken-2.white--text.py-3.text-xs-center(v-else) No pbehaviors
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions: pbehaviorMapAction, mapGetters: pbehaviorMapgetters } = createNamespacedHelpers('pbehavior');

export default {
  name: 'pbehaviors-list',
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
          text: 'Name',
          sortable: false,
        },
        {
          text: 'Author',
          sortable: false,
        },
        {
          text: 'Connector',
          sortable: false,
        },
        {
          text: 'Connector name',
          sortable: false,
        },
        {
          text: 'Is Enabled',
          sortable: false,
        },
        {
          text: 'Begins',
          sortable: false,
        },
        {
          text: 'Ends',
          sortable: false,
        },
        {
          text: 'Type',
          sortable: false,
        },
        {
          text: 'Reason',
          sortable: false,
        },
        {
          text: 'Rrule',
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

