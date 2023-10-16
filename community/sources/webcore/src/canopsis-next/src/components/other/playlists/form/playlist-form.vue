<template>
  <div>
    <c-name-field
      v-field="form.name"
      required="required"
    />
    <c-duration-field v-field="form.interval" />
    <v-layout>
      <c-enabled-field v-field="form.enabled" />
      <c-enabled-field
        v-field="form.fullscreen"
        :label="$t('common.fullscreen')"
      />
    </v-layout>
    <v-btn
      class="secondary ml-0"
      @click="showManageTabsModal"
    >
      {{ $t('modals.createPlaylist.manageTabs') }}
    </v-btn>
    <v-layout
      class="py-4"
    >
      <v-layout
        v-if="tabsPending"
        justify-center="justify-center"
      >
        <v-progress-circular
          color="primary"
          indeterminate="indeterminate"
        />
      </v-layout>
      <v-flex
        v-else
        xs12="xs12"
      >
        <v-flex class="text-center mb-2">
          {{ $t('common.result') }}
        </v-flex>
        <draggable-playlist-tabs v-field="form.tabs_list" />
      </v-flex>
    </v-layout>
    <v-layout>
      <v-alert
        :value="errors.has('tabs')"
        type="error"
      >
        {{ $t('modals.createPlaylist.errors.emptyTabs') }}
      </v-alert>
    </v-layout>
  </div>
</template>

<script>
import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

import GroupViewPanel from '@/components/layout/navigation/partials/groups-side-bar/group-view-panel.vue';
import GroupPanel from '@/components/layout/navigation/partials/groups-side-bar/group-panel.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partials/groups-side-bar/groups-side-bar-group.vue';
import DraggablePlaylistTabs from '@/components/other/playlists/form/fields/draggable-playlist-tabs.vue';

import TabPanelContent from '../partials/tab-panel-content.vue';

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
