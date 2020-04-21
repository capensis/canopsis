<template lang="pug">
  v-container
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(fixed-tabs, slider-color="primary")
            v-tab {{ $t('parameters.tabs.parameters') }}
            v-tab-item
              v-card-text
                user-interface-form(:disabled="!hasUpdateParametersAccess")
            v-tab {{ $t('parameters.tabs.importExportViews') }}
            v-tab-item
              views-import-export
            v-tab {{ $t('parameters.tabs.createPlaylist') }}
            v-tab-item
              create-playlist
</template>

<script>
import rightsTechnicalParametersMixin from '@/mixins/rights/technical/parameters';

import UserInterfaceForm from '@/components/forms/user-interface.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';
import ViewsImportExport from '@/components/other/view/views-import-export.vue';
import CreatePlaylist from '@/components/other/playlist/create-playlist.vue';

export default {
  components: {
    ViewsImportExport,
    UserInterfaceForm,
    GroupsSideBarGroup,
    CreatePlaylist,
  },
  mixins: [rightsTechnicalParametersMixin],
};
</script>
