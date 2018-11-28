<template lang="pug">
  v-container.pa-0(fluid)
    v-layout(wrap)
      v-flex.my-2(xs12)
        v-layout(justify-space-around, align-center, wrap)
          v-flex(xs4, md2)
            h4.text-xs-center {{ $t('context.moreInfos.type') }} :
            p.text-xs-center {{ item.type }}
          v-flex(xs4, md2)
            template(v-if="item.enabled")
              v-chip.primary.darken-1.white--text {{ $t('common.enabled') }}
            template(v-else)
              v-chip.red.white--text.title {{ $t('common.disabled') }}
          v-flex(xs4, md2, v-if="lastActiveDate")
            h4.text-xs-center {{ $t('context.moreInfos.lastActiveDate') }} :
            p.text-xs-center {{ lastActiveDate | date('long') }}
          v-flex(xs6, md2)
            v-menu(:value="isImpactExpanded", bottom, offset-y, fixed)
              v-btn(@click.stop="isImpactExpanded = !isImpactExpanded", slot="activator") {{ $t('context.impacts') }}
              v-list(dense)
                template(v-for="item in item.impact")
                  v-list-tile
                    v-list-tile-content {{ item }}
          v-flex(xs6, md2)
            v-menu(:value="isDependsExpanded", bottom, offset-y, fixed)
              v-btn(@click.stop="isDependsExpanded = !isDependsExpanded",
                    slot="activator") {{ $t('context.dependencies') }}
              v-list(dense)
                template(v-for="item in item.depends")
                  v-list-tile
                    v-list-tile-content {{ item }}
      v-flex.my-2(xs12)
        h3.text-xs-center.my-2 Pbehaviors
        pbehaviors-list(:itemId="item._id")
      v-flex.my-2(xs12)
        h3.text-xs-center Infos
        v-container(fluid, grid-list-sm)
          v-layout(row, wrap)
            v-flex(v-for="(value, key) in item.infos", :key="key", xs4)
              h4.text-xs-center {{ key }}
              p.text-xs-center {{ $t('common.description') }} : {{ value.description }}
              p.text-xs-center {{ $t('common.value') }} : {{ value.value }}


</template>

<script>
import PbehaviorsList from './pbehaviors-list.vue';

export default {
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
    };
  },
  computed: {
    lastActiveDate() {
      if (this.item.enable_history) {
        return Math.max(this.item.enable_history) * 1000;
      }

      return null;
    },
  },
  methods: {
  },
};
</script>

