import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

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
    /* this.$socket
      .join(SOCKET_ROOMS.icons)
      .addListener(this.setIcon); */
  },
  methods: {
    ...mapActions({
      fetchIconsListWithoutStore: 'fetchListWithoutStore',
      fetchIconWithoutStore: 'fetchItemWithoutStore',
      setRegisteredIcons: 'setRegisteredIcons',
      addRegisteredIcon: 'addRegisteredIcon',
      removeRegisteredIcon: 'removeRegisteredIcon',
    }),

    async setIcon({ id }) {
      const prevIcon = this.registeredIconsById[id];
      const icon = await this.fetchIconWithoutStore({ id });

      if (prevIcon) {
        this.unregisterIconFromVuetify(prevIcon.title);
      }

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
