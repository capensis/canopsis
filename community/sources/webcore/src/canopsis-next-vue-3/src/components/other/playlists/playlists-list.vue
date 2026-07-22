<template>
  <c-advanced-data-table
    :headers="headers"
    :items="playlists"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable || updatable"
    advanced-pagination
    expand
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        :enablable="updatable"
        :disablable="updatable"
        playlist
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>

    <template #fullscreen="{ item }">
      <c-enabled :value="item.fullscreen" />
    </template>

    <template #interval="{ item }">
      {{ formatDuration(item.interval) }}
    </template>

    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>

    <template #actions="{ item }">
      <c-layout gap="1">
        <c-action-btn :tooltip="t('common.play')">
          <template #button="{ props: buttonProps }">
            <v-btn v-bind="buttonProps" :to="getPlaylistRouteById(item._id)" variant="text" icon>
              <c-icon>play_arrow</c-icon>
            </v-btn>
          </template>
        </c-action-btn>
        <c-copy-btn
          :value="getPlaylistRouteFullUrlById(item._id)"
          :tooltip="t('common.copyLink')"
          @success="onSuccessCopied"
          @error="onErrorCopied"
        />
        <c-action-btn v-if="duplicable" type="duplicate" @click="$emit('duplicate', item)" />
        <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn v-if="removable" type="delete" @click="$emit('remove', item._id)" />
      </c-layout>
    </template>

    <template #expand="{ item }">
      <playlists-list-expand-item :playlist="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { APP_HOST } from '@/config'
import { ROUTES_NAMES } from '@/constants'

import { convertDurationToString } from '@/helpers/date/duration'
import { removeTrailingSlashes } from '@/helpers/url'

import { useI18n } from '@/hooks/i18n'
import { usePopups } from '@/hooks/popups'

import PlaylistsListExpandItem from './partials/playlists-list-expand-item.vue'

export default {
  components: {
    PlaylistsListExpandItem,
  },

  props: {
    playlists: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:options', 'edit', 'duplicate', 'remove', 'refresh'],

  setup() {
    const router = useRouter()
    const popups = usePopups()
    const { t } = useI18n()

    const headers = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('common.fullscreen'), value: 'fullscreen', sortable: false },
      { text: t('common.enabled'), value: 'enabled' },
      { text: t('common.interval'), value: 'interval', sortable: false },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ])

    const getPlaylistRouteById = (id) => ({
      name: ROUTES_NAMES.playlist,
      params: { id },
      query: { autoplay: true },
    })

    const getPlaylistRouteFullUrlById = (id) => {
      const { href } = router.resolve(getPlaylistRouteById(id))

      return removeTrailingSlashes(`${APP_HOST}${href}`)
    }

    const onSuccessCopied = () => popups.success({ text: t('success.linkCopied') })
    const onErrorCopied = () => popups.error({ text: t('errors.default') })

    return {
      headers,
      t,
      formatDuration: convertDurationToString,
      getPlaylistRouteById,
      getPlaylistRouteFullUrlById,
      onSuccessCopied,
      onErrorCopied,
    }
  },
}
</script>
