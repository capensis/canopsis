import { isBoolean, cloneDeep } from 'lodash';

import { TIME_UNITS } from '@/constants';

import { durationToForm, formToDuration } from '@/helpers/date/duration';

/**
 * @typedef {Array<string>} TabsList
 */

/**
 * @typedef {Object} Playlist
 * @property {string} name
 * @property {number} created
 * @property {number} updated
 * @property {boolean} enabled
 * @property {boolean} fullscreen
 * @property {Duration} interval
 * @property {TabsList} tabs_list
 * @property {string} _id
 */

/**
 * @typedef {Playlist} PlaylistForm
 * @property {DurationForm} interval
 */

/**
 * Convert playlist to playlist form
 *
 * @param {Playlist} [playlist = {}]
 * @returns {PlaylistForm}
 */
export const playlistToForm = (playlist = {}) => ({
  interval: playlist.interval
    ? durationToForm(playlist.interval)
    : { value: 10, unit: TIME_UNITS.second },
  name: playlist.name || '',
  fullscreen: isBoolean(playlist.fullscreen) ? playlist.fullscreen : true,
  enabled: isBoolean(playlist.enabled) ? playlist.enabled : true,
  tabs_list: playlist.tabs_list ? cloneDeep(playlist.tabs_list) : [],
});

/**
 * Convert playlist form to playlist
 *
 * @param {PlaylistForm} [form = {}]
 * @returns {Playlist}
 */
export const formToPlaylist = (form = {}) => ({
  ...form,

  interval: formToDuration(form.interval),
  tabs_list: form.tabs_list.map(tab => tab._id),
});
