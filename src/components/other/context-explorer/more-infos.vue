<template lang="pug">
  v-container(fluid)
    v-layout(wrap)
      //v-flex(xs12)
        v-btn(@click.stop="click") Test
      v-flex.my-2(xs12)
        v-layout(justify-space-around)
          v-flex(xs2)
            div Type : {{ item.props.type }}
          v-flex(xs2)
            template(v-if="item.props.enabled")
              v-chip.green.white--text Enabled
            template(v-else)
              v-chip.red.white--text.title Disabled
          v-flex(xs2)
            div Last active date :
          v-flex(xs3)
            div(@click.stop="isImpactExpanded = true")
              v-expansion-panel(:value="isImpactExpanded")
                v-expansion-panel-content
                  div(slot="header") Impact
                  v-list(dense)
                    template(v-for="item in item.props.impact")
                      v-list-tile
                        v-list-tile-content {{ item }}
          v-flex(xs3)
            div(@click.stop="isDependsExpanded = true")
              v-expansion-panel(:value="isDependsExpanded")
                v-expansion-panel-content
                  div(slot="header") Depends
                  v-list(dense)
                    template(v-for="item in item.props.depends")
                      v-list-tile
                        v-list-tile-content {{ item }}
      v-flex.my-2(xs12)
        h3.text-xs-center.my-2 Pbehaviors
        v-data-table(v-if="pbehaviorsList.length > 0", :items="pbehaviorsList", :headers="pbehaviorsTableHeaders")
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
        div.red.darken-2.white--text.py-3.text-xs-center(v-else) No pbehaviors
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions: pbehaviorMapAction, mapGetters: pbehaviorMapGetters } = createNamespacedHelpers('pbehavior');

export default {
  name: 'context-more-infos',
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      isImpactExpanded: false,
      isDependsExpanded: false,
      pbehaviors: [],
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
    ...pbehaviorMapGetters(['error', 'pbehaviorsList']),
  },
  mounted() {
    this.fetchPbehaviorsList({ id: this.item.props._id });
  },
  methods: {
    ...pbehaviorMapAction({
      fetchPbehaviorsList: 'fetchById',
    }),
    click() {
      console.log(this.item);
    },
  },
};
</script>

