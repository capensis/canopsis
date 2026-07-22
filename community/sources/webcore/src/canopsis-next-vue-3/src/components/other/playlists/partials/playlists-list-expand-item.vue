<template>
  <div>
    <v-tabs v-model="activeTab" class="secondary lighten-1" color="primary" align-tabs="center">
      <v-tab>{{ t('common.information', 1) }}</v-tab>
    </v-tabs>

    <v-window v-model="activeTab" mandatory>
      <v-window-item>
        <div class="pa-3">
          <v-card>
            <v-card-text>
              <c-progress-overlay :pending="pending" />
              <draggable-playlist-tabs :tabs="availableTabs" :groups="groups" disabled />
            </v-card-text>
          </v-card>
        </div>
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { computed, onMounted, ref } from 'vue'

import { MAX_LIMIT } from '@/constants'

import { getTabsFromGroupsByIds } from '@/helpers/entities/playlist/form'

import { useI18n } from '@/hooks/i18n'
import { useView } from '@/hooks/store/modules/view'

import DraggablePlaylistTabs from '@/components/other/playlists/form/fields/draggable-playlist-tabs.vue'

export default {
  components: { DraggablePlaylistTabs },

  props: {
    playlist: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n()
    const { fetchViewGroupsListWithoutStore } = useView()

    const activeTab = ref(0)
    const groups = ref([])
    const pending = ref(false)

    const availableTabs = computed(() =>
      getTabsFromGroupsByIds(groups.value, props.playlist.tabs_list ?? []),
    )

    const fetchGroups = async () => {
      try {
        pending.value = true

        const { data = [] } = await fetchViewGroupsListWithoutStore({
          params: {
            limit: MAX_LIMIT,
            with_views: true,
            with_tabs: true,
          },
        })

        groups.value = data
      } finally {
        pending.value = false
      }
    }

    onMounted(fetchGroups)

    return {
      t,
      activeTab,
      groups,
      pending,
      availableTabs,
    }
  },
}
</script>
