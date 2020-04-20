<template lang="pug">
  div
    h1 Playlist
    v-fade-transition(v-if="activeTab", mode="out-in")
      view-tab-rows(:tab="activeTab", :key="activeTab._id")
    v-btn(@click="changeActiveTab") Change active tab
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ENTITIES_TYPES } from '@/constants';

import ViewTabRows from '@/components/other/view/view-tab-rows.vue';

const { mapGetters } = createNamespacedHelpers('entities');

export default {
  components: { ViewTabRows },
  data() {
    return {
      activeTabIndex: 0,
      playlist: {
        _id: 'id123',
        name: 'Playlist #1',
        interval: {
          value: 10,
          unit: 'm',
        },
        tabs: [
          'view-tab_5a339b3a-0611-4d4c-b307-dc1b92aeb27d', // Meteo technic
          '875df4c2-027b-4549-8add-e20ed7ff7d4f', // Alarm default
          'view-tab_c02ae48e-7f0a-4ba4-9215-ba5662e1550c', // Meteo correct
        ],
      },
    };
  },
  computed: {
    ...mapGetters(['getList']),

    tabs() {
      return this.getList(ENTITIES_TYPES.viewTab, this.playlist.tabs);
    },

    activeTab() {
      return this.tabs[this.activeTabIndex];
    },
  },
  methods: {
    changeActiveTab() {
      this.activeTabIndex = this.activeTabIndex >= this.tabs.length - 1 ? 0 : this.activeTabIndex + 1;
    },
  },
};
</script>
