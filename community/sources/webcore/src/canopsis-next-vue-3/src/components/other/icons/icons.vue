<template>
  <icons-list
    :options="options"
    :icons="icons"
    :total-items="meta.total_count"
    :pending="pending"
    :addable="hasCreateAnyIconAccess"
    :updatable="hasUpdateAnyIconAccess"
    :removable="hasDeleteAnyIconAccess"
    @edit="showEditIconModal"
    @remove="showRemoveIconModal"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted, watch } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useIcon } from '@/hooks/store/modules/icon';
import { useCallActionWithPopup } from '@/hooks/actions/call';

import {
  registerCustomIcon,
  customIconsRegistry,
  unregisterCustomIcon,
} from '@/plugins/vuetify/icons';

import IconsList from './icons-list.vue';

export default {
  components: { IconsList },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const {
      hasCreateAccess: hasCreateAnyIconAccess,
      hasUpdateAccess: hasUpdateAnyIconAccess,
      hasDeleteAccess: hasDeleteAnyIconAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.icon);
    const { fetchIconsListWithoutStore, updateIcon, removeIcon, setRegisteredIcons } = useIcon();
    const { callActionWithPopup } = useCallActionWithPopup();

    const {
      data: icons,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: fetchIconsListWithoutStore,
    });

    watch(icons, (newIcons) => {
      newIcons.forEach(({ title, content, fill_border: fillBorder }) => {
        if (!customIconsRegistry[title]) {
          registerCustomIcon(title, content, { fillBorder });
        }
      });

      setRegisteredIcons({ icons: newIcons });
    });

    const showEditIconModal = (icon) => {
      modals.show({
        name: MODALS.createIcon,
        config: {
          icon,
          title: t('modals.createIcon.create.title'),
          action: (newIcon) =>
            callActionWithPopup(
              () => updateIcon({ data: newIcon, id: icon._id }),
              fetchList,
              t('modals.createIcon.create.success'),
            ),
        },
      });
    };

    const showRemoveIconModal = (icon) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: () =>
            callActionWithPopup(
              async () => {
                await removeIcon({ id: icon._id });
                unregisterCustomIcon(icon.title);
              },
              fetchList,
              t('modals.createIcon.remove.success'),
            ),
        },
      });
    };

    onMounted(fetchList);

    return {
      icons,
      meta,
      pending,
      options,
      hasCreateAnyIconAccess,
      hasUpdateAnyIconAccess,
      hasDeleteAnyIconAccess,
      showEditIconModal,
      showRemoveIconModal,
      updateOptions,
      fetchList,
    };
  },
};
</script>
