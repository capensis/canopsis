<template lang="pug">
  div
    v-layout(row)
      v-switch(
        v-field="form.enabled",
        :label="$t('common.enabled')"
      )
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
        v-field="form.fullscreen",
        :label="$t('common.fullscreen')"
      )
    v-layout.py-4(row)
      v-flex.export-views-block.mr-2(xs6)
        v-flex.text-xs-center.mb-2 {{ $t('modals.createPlaylist.groups') }}
        v-expansion-panel(readonly, hide-actions, expand, dark, focusable, :value="openedPanels")
          group-panel(
            v-for="group in groups",
            :group="group",
            :key="group._id",
            hideActions
          )
            v-expansion-panel.tabs-panel(
              v-for="view in group.views",
              :key="view._id",
              :value="getPanelValueFromArray(view.tabs)",
              readonly,
              hide-actions,
              expand,
              dark,
              focusable
            )
              v-expansion-panel-content(hide-actions)
                group-view-panel(slot="header", :view="view")
                draggable-tabs(:tabs="view.tabs", pull="clone")
      v-flex.export-views-block.ml-2(xs6)
        v-flex.text-xs-center.mb-2 {{ $t('modals.createPlaylist.result') }}
        draggable-tabs(v-field="form.tabs_list", put, pull, @change="validateTabs")
    v-layout
      v-alert(:value="errors.has('tabs')", type="error") {{ $t('modals.createPlaylist.errors.emptyTabs') }}
</template>

<script>
import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';
import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';
import DraggableTabs from '@/components/other/playlists/form/partials/draggable-tabs.vue';
import TimeIntervalField from '@/components/forms/fields/time-interval.vue';

import TabPanelContent from './partials/tab-panel-content.vue';

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
    openedPanels() {
      return this.getPanelValueFromArray(this.groups);
    },
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
    getPanelValueFromArray(values = []) {
      return new Array(values.length).fill(true);
    },

    validateTabs() {
      this.$nextTick(() => this.$validator.validate('tabs'));
    },
  },
};
</script>

<style lang="scss" scoped>
  .tabs-panel {
    /deep/ .v-expansion-panel__header {
      padding: 0;
      margin: 0;
    }
  }
  .tab-panel-item {
    display: flex;
    align-items: center;
    height: 48px;
  }
</style>
