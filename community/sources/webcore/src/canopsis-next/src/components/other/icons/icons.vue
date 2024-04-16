<template>
  <icons-list
    :options.sync="options"
    :icons="icons"
    :total-items="iconsMeta.total_count"
    :pending="iconsPending"
    :addable="hasCreateAnyIconAccess"
    :updatable="hasUpdateAnyIconAccess"
    :removable="hasDeleteAnyIconAccess"
    @edit="showEditIconModal"
    @remove="showRemoveIconModal"
  />
</template>

<script>
import { MODALS } from '@/constants';

import { localQueryMixin } from '@/mixins/query/query';
import { entitiesIconMixin } from '@/mixins/entities/icon';
import { permissionsTechnicalIconMixin } from '@/mixins/permissions/technical/icon';
import { vuetifyCustomIconsBaseMixin } from '@/mixins/vuetify/custom-icons/base';

import IconsList from '@/components/other/icons/icons-list.vue';

export default {
  components: { IconsList },
  mixins: [
    localQueryMixin,
    entitiesIconMixin,
    permissionsTechnicalIconMixin,
    vuetifyCustomIconsBaseMixin,
  ],
  watch: {
    /**
     * We need to use this watcher (even if we have sockets) to displaying icons in table without blink
     */
    icons(icons) {
      icons.forEach(({ title, content }) => (
        !this.hasIconInVuetify(title) && this.registerIconInVuetify(title, content)
      ));
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    showEditIconModal(icon) {
      this.$modals.show({
        name: MODALS.createIcon,
        config: {
          icon,
          title: this.$t('modals.createIcon.create.title'),
          action: async (newIcon) => {
            await this.updateIcon({ data: newIcon, id: icon._id });

            this.$popups.success({ text: this.$t('modals.createIcon.create.success') });
            this.fetchList();
          },
        },
      });
    },

    showRemoveIconModal(icon) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeIcon({ id: icon._id });

            this.$popups.success({ text: this.$t('modals.createIcon.remove.success') });
            this.fetchList();
          },
        },
      });
    },

    fetchList() {
      return this.fetchIconsList({ params: this.getQuery() });
    },
  },
};
</script>
