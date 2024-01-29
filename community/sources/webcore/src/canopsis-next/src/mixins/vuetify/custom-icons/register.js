import { createNamespacedHelpers } from 'vuex';

import { SOCKET_ROOMS } from '@/config';
import { MAX_LIMIT } from '@/constants';

import { ICONS_RESPONSE_MESSAGES_TYPES } from '@/plugins/socket/constants';

import { vuetifyCustomIconsBaseMixin } from './base';

const { mapGetters, mapActions } = createNamespacedHelpers('icon');

export const vuetifyCustomIconsRegisterMixin = {
  mixins: [vuetifyCustomIconsBaseMixin],
  data() {
    return {
      customIconsById: {},
    };
  },
  computed: {
    ...mapGetters(['registeredIconsById']),
  },
  async mounted() {
    this.$socket
      .join(SOCKET_ROOMS.icons)
      .addListener(this.setIcon);
  },
  methods: {
    ...mapActions({
      fetchIconsListWithoutStore: 'fetchListWithoutStore',
      fetchIconWithoutStore: 'fetchItemWithoutStore',
      setRegisteredIcons: 'setRegisteredIcons',
      addRegisteredIcon: 'addRegisteredIcon',
      removeRegisteredIcon: 'removeRegisteredIcon',
    }),

    async setIcon({ _id: id, type }) {
      const prevIcon = this.registeredIconsById[id];

      if (prevIcon) {
        this.unregisterIconFromVuetify(prevIcon.title);
      }

      if (type === ICONS_RESPONSE_MESSAGES_TYPES.delete) {
        await this.$nextTick();
        this.removeRegisteredIcon({ id });

        return;
      }

      const icon = await this.fetchIconWithoutStore({ id });

      this.registerIconInVuetify(icon.title, icon.content);
      this.addRegisteredIcon({ id, icon });
    },

    async fetchIconsWithRegistering() {
      try {
        const { data: icons } = await this.fetchIconsListWithoutStore({ params: { limit: MAX_LIMIT } });

        icons.map(({ title, content }) => this.registerIconInVuetify(title, content));

        this.setRegisteredIcons({ icons });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
