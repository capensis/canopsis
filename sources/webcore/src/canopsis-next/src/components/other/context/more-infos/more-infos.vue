<template lang="pug">
  v-tabs.visible(v-model="activeTab", color="secondary lighten-1", dark, slider-color="primary", centered)
    v-tab(v-for="(tab, index) in tabs", :key="index") {{ tab }}
    v-tab-item
      pbehaviors-list(:itemId="item._id", :tabId="tabId")
    v-tab-item
      impact-depends(:impact="impact", :depends="depends")
    v-tab-item
      infos(:infos="item.infos")
</template>

<script>
import PbehaviorsList from './tabs/pbehaviors.vue';
import ImpactDepends from './tabs/impact-depends.vue';
import Infos from './tabs/infos.vue';

export default {
  components: {
    PbehaviorsList,
    ImpactDepends,
    Infos,
  },
  props: {
    item: {
      type: Object,
      required: true,
    },
    tabId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      activeTab: 0,
      tabs: [
        this.$t('context.moreInfos.tabs.pbehaviors'),
        this.$t('context.moreInfos.tabs.impactDepends'),
        this.$t('context.moreInfos.tabs.infos'),
      ],
    };
  },
  computed: {
    impact() {
      return this.item.impact.map((id, index) => ({ name: this.item.impact_name[index], id }));
    },
    depends() {
      return this.item.depends.map((id, index) => ({ name: this.item.depends_name[index], id }));
    },
  },
};
</script>

<style lang="scss" scoped>
  .v-tabs.visible {
    & /deep/ > .v-tabs__bar {
      display: block;
    }
  }
</style>
