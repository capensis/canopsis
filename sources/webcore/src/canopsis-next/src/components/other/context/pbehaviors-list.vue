<template lang="pug">
v-data-table(:items="items", :headers="pbehaviorsTableHeaders")
  template(slot="items" slot-scope="props")
    td {{ props.item.name }}
    td {{ props.item.author }}
    td {{ props.item.connector }}
    td {{ props.item.connector_name }}
    td {{ props.item.enabled }}
    td {{ props.item.tstart | date('long') }}
    td {{ props.item.tstop | date('long') }}
    td {{ props.item.type_ }}
    td {{ props.item.reason }}
    td {{ props.item.rrule }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions: pbehaviorMapAction, mapGetters: pbehaviorMapGetters } = createNamespacedHelpers('pbehavior');

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
          text: this.$t('pbehaviors.connector'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.connectorName'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.isEnabled'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.begins'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.ends'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.type'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.reason'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.rrule'),
          sortable: false,
        },
      ],
    };
  },
  computed: {
    ...pbehaviorMapGetters({
      pbehaviorsList: 'items',
    }),
  },
  mounted() {
    this.fetchItems();
  },
  methods: {
    ...pbehaviorMapAction({
      fetchPbehaviorsList: 'fetchListByEntityId',
    }),
    async fetchItems() {
      await this.fetchPbehaviorsList({ id: this.itemId });
      this.items = [...this.pbehaviorsList];
    },
  },
};
</script>

