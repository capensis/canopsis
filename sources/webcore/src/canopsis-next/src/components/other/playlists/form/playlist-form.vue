<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    time-interval-field(v-field="form.interval")
    v-layout(row)
      v-switch(
        v-field="form.enabled",
        :label="$t('common.enabled')",
        color="primary"
      )
      v-switch(
        v-field="form.fullscreen",
        :label="$t('common.fullscreen')",
        color="primary"
      )
    v-btn.secondary(@click="showManageTabsModal") {{ $t('modals.createPlaylist.manageTabs') }}
    v-layout.py-4(row)
      v-flex(xs12)
        v-flex.text-xs-center.mb-2 {{ $t('modals.createPlaylist.result') }}
        draggable-tabs(v-field="form.tabs_list", put, pull, @change="validateTabs")
          v-layout.tab-panel-content(xs12, slot="title", slot-scope="{ tab }")
            v-flex.tab-content-block.secondary.pa-2.white--text(xs4) {{ getGroupByTab(tab).name }}
            v-flex.tab-content-block.secondary.lighten-1.pa-2.white--text(xs4) {{ getViewByTab(tab).name }}
            v-flex.tab-content-block.pa-2.white--text(xs4) {{ tab.title }}
    v-layout
      v-alert(:value="errors.has('tabs')", type="error") {{ $t('modals.createPlaylist.errors.emptyTabs') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';
import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';
import DraggableTabs from '@/components/other/playlists/form/partials/draggable-tabs.vue';
import TimeIntervalField from '@/components/forms/fields/time-interval.vue';

import { ENTITIES_TYPES, MODALS } from '@/constants';
import { SCHEMA_EMBEDDED_KEY } from '@/config';

import formMixin from '@/mixins/form';

import TabPanelContent from './partials/tab-panel-content.vue';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

export default {
  inject: ['$validator'],
  components: {
    TimeIntervalField,
    DraggableTabs,
    TabPanelContent,
    GroupViewPanel,
    GroupPanel,
    GroupsSideBarGroup,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: false,
    },
    groups: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    ...entitiesMapGetters({
      getEntityItem: 'getItem',
    }),
  },
  created() {
    this.$validator.attach({
      name: 'tabs',
      rules: 'required:true',
      getter: () => this.form.tabs_list.length > 0,
      context: () => this,
      vm: this,
    });
  },
  methods: {
    validateTabs() {
      this.$nextTick(() => this.$validator.validate('tabs'));
    },
    showManageTabsModal() {
      this.$modals.show({
        name: MODALS.managePlaylistTabs,
        config: {
          groups: this.groups,
          selectedTabs: this.form.tabs_list,
          action: tabs => this.updateField('tabs_list', tabs),
        },
      });
    },
    getViewByTab(tab) {
      const tabWithEmbedded = this.getEntityItem(ENTITIES_TYPES.viewTab, tab._id, true);

      const { parents: [parent] } = tabWithEmbedded[SCHEMA_EMBEDDED_KEY];

      return this.getEntityItem(parent.type, parent.id);
    },
    getGroupByTab(tab) {
      const view = this.getViewByTab(tab);

      return this.getEntityItem(ENTITIES_TYPES.group, view.group_id);
    },
  },
};
</script>

<style lang="scss" scoped>
  .tab-panel-content {
    height: 100%;
  }
  .tab-content-block {
    display: flex;
    align-items: center;
  }
</style>
