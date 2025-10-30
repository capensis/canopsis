<template>
  <top-bar-menu :title="userName" :links="links" without-sort />
</template>

<script>
import { computed, inject } from 'vue';
import { useRouter } from 'vue-router/composables';

import { MODALS, ROUTES_NAMES, USER_PERMISSIONS } from '@/constants';

import { useAuth } from '@/hooks/auth';
import { useUser } from '@/hooks/store/modules/user';
import { useInfo } from '@/hooks/store/modules/info';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';

import { useTopBarMenu } from './hooks/top-bar-menu';
import TopBarMenu from './top-bar-menu.vue';

export default {
  components: { TopBarMenu },
  setup() {
    const $system = inject('$system');
    const router = useRouter();

    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const { currentUser, logout, fetchCurrentUser } = useAuth();
    const { updateCurrentUser } = useUser();
    const { defaultColorTheme } = useInfo();
    const { filterLinks } = useTopBarMenu();

    const userName = computed(() => currentUser.value.display_name || currentUser.value._id);

    /**
     * Opens the user profile modal for editing user preferences.
     * After successful update, shows a success popup, refreshes the current user data,
     * and applies the updated theme.
     */
    const showEditUserModal = () => {
      modals.show({
        name: MODALS.createUser,
        config: {
          title: t('common.profile'),
          user: currentUser.value,
          onlyUserPrefs: true,
          action: async (data) => {
            await updateCurrentUser({ data });

            popups.success({ text: t('success.default', data.ui_language) });

            await fetchCurrentUser();

            $system.setTheme(currentUser.value.ui_theme_colors);
          },
        },
      });
    };

    /**
     * Handles user logout.
     * Logs out the current user, redirects to the login page,
     * and resets the theme to the default color theme.
     */
    const logoutHandler = async () => {
      await logout({
        redirect: () => router.replaceAsync({ name: ROUTES_NAMES.login }),
      });

      $system.setTheme(defaultColorTheme.value);
    };

    const links = computed(() => {
      const rawLinks = [
        {
          icon: 'person',
          title: t('user.seeProfile'),
          handler: showEditUserModal,
        },
        {
          icon: 'filter_list',
          title: t('pattern.patterns'),
          route: { name: ROUTES_NAMES.profilePatterns },
        },
        {
          icon: 'palette',
          title: t('theme.themes'),
          route: { name: ROUTES_NAMES.profileThemes },
          permission: USER_PERMISSIONS.technical.profile.theme,
        },
        {
          icon: 'exit_to_app',
          color: 'error',
          class: 'top-bar-user-menu__logout-btn',
          title: t('common.logout'),
          handler: logoutHandler,
        },
      ].map(link => ({
        ...link,

        class: [link.class, 'text-uppercase text-body-2'].filter(Boolean).join(' '),
      }));

      return filterLinks(rawLinks);
    });

    return {
      userName,
      links,
    };
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
