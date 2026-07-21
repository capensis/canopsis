import { mapActions, mapState, mapStores } from 'pinia'

import { SOCKET_ROOMS } from '@/config'
import { MAX_LIMIT } from '@/constants'

import { ICONS_RESPONSE_MESSAGES_TYPES } from '@/plugins/socket/constants'

import { useIconStore } from '@/stores/modules/icon'

import { vuetifyCustomIconsBaseMixin } from './base'

export const vuetifyCustomIconsRegisterMixin = {
  mixins: [vuetifyCustomIconsBaseMixin],
  data() {
    return {
      customIconsById: {},
    }
  },
  computed: {
    ...mapStores(useIconStore),
    ...mapState(useIconStore, ['registeredIconsById']),
  },
  async mounted() {
    this.$socket?.join?.(SOCKET_ROOMS.icons)?.addListener?.(this.setIcon)
  },
  methods: {
    ...mapActions(useIconStore, [
      'fetchListWithoutStore',
      'fetchItemWithoutStore',
      'setRegisteredIcons',
      'addRegisteredIcon',
      'removeRegisteredIcon',
    ]),

    async setIcon({ _id: id, type }) {
      const prevIcon = this.registeredIconsById[id]

      if (prevIcon) {
        this.unregisterIconFromVuetify(prevIcon.title)
      }

      if (type === ICONS_RESPONSE_MESSAGES_TYPES.delete) {
        await this.$nextTick()
        this.removeRegisteredIcon({ id })

        return
      }

      const icon = await this.fetchItemWithoutStore({ id })

      this.registerIconInVuetify(icon.title, icon.content, { fillBorder: icon.fill_border })
      this.addRegisteredIcon({ id, icon })
    },

    async fetchIconsWithRegistering() {
      try {
        const response = await this.fetchListWithoutStore({ params: { limit: MAX_LIMIT } })
        const icons = Array.isArray(response?.data) ? response.data : (response?.data?.data ?? [])

        icons.forEach(({ title, content, fill_border: fillBorder }) =>
          this.registerIconInVuetify(title, content, { fillBorder }),
        )

        this.setRegisteredIcons({ icons })
      } catch (err) {
        console.error(err)
      }
    },
  },
}

export default vuetifyCustomIconsRegisterMixin
