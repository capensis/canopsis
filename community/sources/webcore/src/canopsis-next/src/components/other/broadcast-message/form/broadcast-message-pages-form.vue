<template>
  <v-treeview
    v-model="selectedItems"
    :items="treeItems"
    :open="openedNodes"
    selected-color="primary"
    item-key="id"
    item-children="children"
    selectable
    open-on-click
    transition
    @input="updateSelection"
  >
    <template #prepend="{ item, open }">
      <v-icon v-if="item.children">
        {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
      </v-icon>
      <v-icon v-else>
        {{ item.icon || 'mdi-file-document' }}
      </v-icon>
    </template>
    <template #label="{ item }">
      <span>{{ item.name }}</span>
    </template>
  </v-treeview>
</template>

<script>
import { ref, computed, onMounted } from 'vue';

import { BROADCAST_MESSAGE_VIEWS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useViewGroup } from '@/hooks/store/modules/view';
import { usePlaylist } from '@/hooks/store/modules/playlist';

export default {
  model: {
    prop: 'views',
    event: 'input',
  },
  props: {
    views: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const { t, tc } = useI18n();

    const selectedItems = ref([]);
    const openedNodes = ref([]);

    const { groups, fetchAllGroupsListWithWidgets } = useViewGroup();
    const { items: playlistItems, fetchList: fetchPlaylistList } = usePlaylist();

    const viewGroups = computed(() => {
      if (!groups.value || !groups.value.length) {
        return [];
      }

      return groups.value.map((group) => {
        const groupId = `${BROADCAST_MESSAGE_VIEWS.allViews}_${group._id}`;

        const child = {
          id: groupId,
          name: group.title || group.name,
          children: (group.views || []).map(view => ({
            id: `${groupId}_${view._id}`,
            name: view.title || view.name,
          })),
        };

        return child;
      });
    });

    const playlistsTree = computed(() => {
      if (!playlistItems.value || !playlistItems.value.length) {
        return [];
      }

      return playlistItems.value.map(playlist => ({
        id: `${BROADCAST_MESSAGE_VIEWS.allPlaylists}_${playlist._id}`,
        name: playlist.name,
      }));
    });

    const treeItems = computed(() => [
      {
        id: BROADCAST_MESSAGE_VIEWS.login,
        name: t('common.login'),
      },
      {
        id: BROADCAST_MESSAGE_VIEWS.exploitation,
        name: t('common.exploitation'),
      },
      {
        id: BROADCAST_MESSAGE_VIEWS.administration,
        name: t('common.administration'),
      },
      {
        id: BROADCAST_MESSAGE_VIEWS.notifications,
        name: tc('common.notification', 2),
      },
      {
        id: BROADCAST_MESSAGE_VIEWS.profile,
        name: t('common.profile'),
      },
      {
        id: BROADCAST_MESSAGE_VIEWS.allViews,
        name: t('broadcastMessage.allViews'),
        children: viewGroups.value,
      },
      {
        id: BROADCAST_MESSAGE_VIEWS.allPlaylists,
        name: t('broadcastMessage.allPlaylists'),
        children: playlistsTree.value,
      },
    ]);

    const updateSelection = (value) => {
      emit('input', value);
    };

    onMounted(async () => {
      try {
        await Promise.all([
          fetchAllGroupsListWithWidgets(),
          fetchPlaylistList(),
        ]);
      } catch (err) {
        console.error(err);
      }

      openedNodes.value = [BROADCAST_MESSAGE_VIEWS.allViews, BROADCAST_MESSAGE_VIEWS.allPlaylists];
    });

    return {
      selectedItems,
      openedNodes,
      treeItems,
      updateSelection,
    };
  },
};
</script>
