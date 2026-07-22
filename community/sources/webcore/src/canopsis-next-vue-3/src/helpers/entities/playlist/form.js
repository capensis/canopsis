import { cloneDeep, isBoolean } from 'lodash'

import { TIME_UNITS } from '@/constants'

import { durationToForm } from '@/helpers/date/duration'

export const playlistToForm = (playlist = {}) => ({
  interval: durationToForm(playlist.interval ?? { value: 10, unit: TIME_UNITS.second }),
  name: playlist.name ?? '',
  fullscreen: isBoolean(playlist.fullscreen) ? playlist.fullscreen : true,
  enabled: isBoolean(playlist.enabled) ? playlist.enabled : true,
  tabs_list: playlist.tabs_list ? cloneDeep(playlist.tabs_list) : [],
})

export const formToPlaylist = (form = {}) => ({
  ...form,

  tabs_list: form.tabs_list.map((tab) => tab._id),
})

export const getTabsFromGroupsByIds = (groups = [], tabIds = []) => {
  const tabIdsSet = new Set(tabIds)

  return groups.reduce((acc, group) => {
    group.views?.forEach((view) => {
      view.tabs?.forEach((tab) => {
        if (tabIdsSet.has(tab._id)) {
          acc.push(tab)
        }
      })
    })

    return acc
  }, [])
}
