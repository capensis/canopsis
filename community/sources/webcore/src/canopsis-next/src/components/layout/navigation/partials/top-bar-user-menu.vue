<template>
  <v-menu
    bottom
    offset-y
    offset-x
  >
    <template #activator="{ on }">
      <v-btn
        class="white--text"
        text
        v-on="on"
      >
        {{ userName }}
      </v-btn>
    </template>
    <v-list class="py-0">
      <top-bar-profile-menu-link
        v-for="link in links"
        :key="link.title"
        :link="link"
      />
    </v-list>
  </v-menu>
</template>

<script>
import { MODALS, ROUTES_NAMES, USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { entitiesUserMixin } from '@/mixins/entities/user';
import { entitiesInfoMixin } from '@/mixins/entities/info';
import { layoutNavigationTopBarMenuMixin } from '@/mixins/layout/navigation/top-bar-menu';

import TopBarProfileMenuLink from './top-bar-profile-menu-link.vue';

export default {
  inject: ['$system'],
  components: { TopBarProfileMenuLink },
  mixins: [
    authMixin,
    entitiesUserMixin,
    entitiesInfoMixin,
    layoutNavigationTopBarMenuMixin,
  ],
  computed: {
    userName() {
      return this.currentUser.display_name || this.currentUser._id;
    },

    links() {
      const links = [
        {
          icon: 'person',
          title: this.$t('user.seeProfile'),
          handler: this.showEditUserModal,
        },
        {
          icon: 'filter_list',
          title: this.$t('pattern.patterns'),
          route: { name: ROUTES_NAMES.profilePatterns },
        },
        {
          icon: 'palette',
          title: this.$t('theme.themes'),
          route: { name: ROUTES_NAMES.profileThemes },
          permission: USERS_PERMISSIONS.technical.profile.theme,
        },
        {
          icon: 'exit_to_app',
          color: 'error',
          class: 'top-bar-user-menu__logout-btn',
          title: this.$t('common.logout'),
          handler: this.logoutHandler,
        },
      ];

      return this.filterLinks(links);
    },
  },
  methods: {
    showEditUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          title: this.$t('common.profile'),
          user: this.currentUser,
          onlyUserPrefs: true,
          action: async (data) => {
            await this.updateCurrentUser({ data });

            this.$popups.success({ text: this.$t('success.default', data.ui_language) });

            await this.fetchCurrentUser();

            this.$system.setTheme(this.currentUser.ui_theme);
          },
        },
      });
    },

    async logoutHandler() {
      await this.logout({
        redirect: () => this.$router.replaceAsync({ name: ROUTES_NAMES.login }),
      });

      this.$system.setTheme(this.defaultColorTheme);
    },
  },
};
</script>

<style lang="scss">
.top-bar-user-menu {
  --btn-error-color: rgba(255, 82, 82, .1);

  &__logout-btn {
    &:hover, &:active {
      background: var(--btn-error-color);
    }
  }
}
</style>
