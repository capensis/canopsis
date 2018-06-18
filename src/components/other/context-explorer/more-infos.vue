<template lang="pug">
  v-container(fluid)
    v-layout(wrap)
      //v-flex(xs12)
        v-btn(@click.stop="click") Test
      v-flex(xs12)
        v-layout(justify-space-around, align-center)
          v-flex(xs3)
            div Type : {{ item.props.type }}
          v-flex(xs3)
            template(v-if="item.props.enabled")
              v-chip.green.white--text Enabled
            template(v-else)
              v-chip.red.white--text.title Disabled
          v-flex(xs3)
            div Last active date : {{  }}
      v-flex.my-2(xs12)
        v-layout(wrap, justify-space-around)
          v-flex(xs5)
            div(@click.stop="isImpactExpanded = true")
              v-expansion-panel(:value="isImpactExpanded")
                v-expansion-panel-content
                  div(slot="header") Impact
                  v-list(dense)
                    template(v-for="item in item.props.impact")
                      v-list-tile
                        v-list-tile-content {{ item }}
          v-flex(xs5)
            div(@click.stop="isDependsExpanded = true")
              v-expansion-panel(:value="isDependsExpanded")
                v-expansion-panel-content
                  div(slot="header") Depends
                  v-list(dense)
                    template(v-for="item in item.props.depends")
                      v-list-tile
                        v-list-tile-content {{ item }}
</template>

<script>
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
    };
  },
  methods: {
    click() {
      console.log(this.item);
    },
  },
};
</script>

