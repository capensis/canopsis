<template lang="pug">
  v-container.pa-0(fluid)
    v-layout(wrap)
      //v-flex(xs12)
        v-btn(@click.stop="click") Test
      v-flex.my-2(xs12)
        v-layout(justify-space-around, wrap)
          v-flex(xs4, md2)
            div Type : {{ item.props.type }}
          v-flex(xs4, md2)
            template(v-if="item.props.enabled")
              v-chip.green.white--text Enabled
            template(v-else)
              v-chip.red.white--text.title Disabled
          v-flex(xs4, md2)
            div Last active date :
          v-flex(xs6, md2)
            v-menu(:value="isImpactExpanded", bottom, offset-y, fixed)
              v-btn(@click.stop="isImpactExpanded = !isImpactExpanded", slot="activator") Impact
              v-list(dense)
                template(v-for="item in item.props.impact")
                  v-list-tile
                    v-list-tile-content {{ item }}
          v-flex(xs6, md2)
            v-menu(:value="isDependsExpanded", bottom, offset-y, fixed)
              v-btn(@click.stop="isDependsExpanded = !isDependsExpanded", slot="activator") Depends
              v-list(dense)
                template(v-for="item in item.props.depends")
                  v-list-tile
                    v-list-tile-content {{ item }}
      v-flex.my-2(xs12)
        h3.text-xs-center.my-2 Pbehaviors
        pbehaviors-list(:itemId="item.props._id")
      v-flex.my-2(xs12)
        h3.text-xs-center Infos
</template>

<script>
import PbehaviorsList from '@/components/other/context-explorer/pbehaviors-list.vue';

export default {
  name: 'context-more-infos',
  components: {
    PbehaviorsList,
  },
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
  methods: {
    click() {
      console.log(this.item);
    },
  },
};
</script>

