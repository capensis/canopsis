<template>
  <top-bar-menu
    :title="userName"
    :links="links"
    content-class="top-bar-user-menu__content"
    without-sort
  />
</template>

<script>
import { computed, inject } from 'vue';
import { useRouter } from 'vue-router';

import { MODALS, ROUTES_NAMES, USER_PERMISSIONS } from '@/constants';

import { useAuth } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useInfo } from '@/hooks/store/modules/info';
import { useUser } from '@/hooks/store/modules/user';

import { useTopBarMenu } from './hooks/top-bar-menu';
import TopBarMenu from './top-bar-menu.vue';

export default {
  components: { TopBarMenu },

  setup() {
    const $system = inject('$system', null);
    const router = useRouter();

    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const { currentUser, logout } = useAuth();
    const { updateCurrentUser } = useUser();
    const { defaultColorTheme } = useInfo();
    const { filterLinks } = useTopBarMenu();

    const userName = computed(() => currentUser.value.display_name || currentUser.value._id || '');

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

            $system?.setTheme(currentUser.value.ui_theme_colors);
          },
        },
      });
    };

    const logoutHandler = async () => {
      await logout({
        redirect: () => router.replace({ name: ROUTES_NAMES.login }),
      });

      $system?.setTheme(defaultColorTheme.value);
    };

    const links = computed(() => filterLinks([
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
      ]),
    );

    return {
      userName,
      links,
    };
  },
};
</script>

<style lang="scss">
.top-bar-user-menu {
  --btn-error-color: rgba(255, 82, 82, 0.1);

  &__content {
    .v-list-item-title {
      text-transform: uppercase;
    }
  }

  &__logout-btn {
    &:hover,
    &:active {
      background: var(--btn-error-color);
    }
  }
}
</style>
