<template lang="pug">
v-data-table(v-if="!pending", :items="items")
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
    };
  },
  computed: {
    ...pbehaviorMapgetters(['pbehaviorsList', 'pending']),
  },
  mounted() {
    this.items = this.fetchPbehaviorsList({ id: this.itemId });
  },
  methods: {
    ...pbehaviorMapAction({
      fetchPbehaviorsList: 'fetchById',
    }),
  },
};
</script>

