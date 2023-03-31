<template lang="pug">
  div
    c-name-field(v-field="form.name", required)
    c-duration-field(v-field="form.interval")
    v-layout(row)
      c-enabled-field(v-field="form.enabled")
      c-enabled-field(v-field="form.fullscreen", :label="$t('common.fullscreen')")
    v-btn.secondary.ml-0(@click="showManageTabsModal") {{ $t('modals.createPlaylist.manageTabs') }}
    v-layout.py-4(row)
      v-layout(v-if="tabsPending", justify-center)
        v-progress-circular(color="primary", indeterminate)
      v-flex(v-else, xs12)
        v-flex.text-xs-center.mb-2 {{ $t('common.result') }}
        draggable-playlist-tabs(v-field="form.tabs_list")
    v-layout
      v-alert(:value="errors.has('tabs')", type="error") {{ $t('modals.createPlaylist.errors.emptyTabs') }}
</template>

<script>
import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';
import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';
import DraggablePlaylistTabs from '@/components/other/playlists/form/partials/draggable-playlist-tabs.vue';

import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

import TabPanelContent from './partials/tab-panel-content.vue';

export default {
  inject: ['$validator'],
  components: {
    DraggablePlaylistTabs,
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
    tabsPending: {
      type: Boolean,
      default: false,
    },
  },
  created() {
    this.$validator.attach({
      name: 'tabs',
      rules: 'required:true',
      getter: () => this.form.tabs_list.length > 0,
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
