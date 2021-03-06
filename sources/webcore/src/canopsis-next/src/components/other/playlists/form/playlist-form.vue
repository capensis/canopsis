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
      enabled-field(v-field="form.enabled")
      v-switch(
        v-field="form.fullscreen",
        :label="$t('common.fullscreen')",
        color="primary"
      )
    v-btn.secondary(@click="showManageTabsModal") {{ $t('modals.createPlaylist.manageTabs') }}
    v-layout.py-4(row)
      v-flex(xs12)
        v-flex.text-xs-center.mb-2 {{ $t('modals.createPlaylist.result') }}
        draggable-playlist-tabs(v-field="form.tabs_list")
    v-layout
      v-alert(:value="errors.has('tabs')", type="error") {{ $t('modals.createPlaylist.errors.emptyTabs') }}
</template>

<script>
import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';
import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';
import DraggablePlaylistTabs from '@/components/other/playlists/form/partials/draggable-playlist-tabs.vue';
import TimeIntervalField from '@/components/forms/fields/time-interval.vue';
import EnabledField from '@/components/forms/fields/enabled-field.vue';

import { MODALS } from '@/constants';

import formMixin from '@/mixins/form';

import TabPanelContent from './partials/tab-panel-content.vue';

export default {
  inject: ['$validator'],
  components: {
    EnabledField,
    DraggablePlaylistTabs,
    TimeIntervalField,
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
          action: (tabs) => {
            this.updateField('tabs_list', tabs);
            this.validateTabs();
          },
        },
      });
    },
  },
};
</script>
