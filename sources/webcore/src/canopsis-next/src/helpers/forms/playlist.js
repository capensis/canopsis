/**
 * Convert playlist to playlist form
 *
 * @param {Object} [playlist = {}]
 * @returns {Object}
 */
export function playlistToForm(playlist = {}) {
  return {
    ...playlist,

    interval: { ...playlist.interval },
    tabs_list: [],
  };
}

/**
 * Convert playlist form to playlist
 *
 * @param {Object} [form = {}]
 * @returns {Object}
 */
export function formToPlaylist(form = {}) {
  return {
    ...form,

    tabs_list: form.tabs_list.map(tab => tab._id),
  };
}
