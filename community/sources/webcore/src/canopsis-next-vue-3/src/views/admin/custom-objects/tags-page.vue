<template>
  <c-page
    :creatable="hasCreateAccess"
    :create-tooltip="t('modals.createTag.create.title')"
    @refresh="fetchList"
    @create="showCreateTagModal"
  >
    <tags-list
      :tags="alarmTags"
      :pending="alarmTagsPending"
      :options="options"
      :total-items="alarmTagsMeta.total_count"
      :updatable="hasUpdateAccess"
      :removable="hasDeleteAccess"
      :duplicable="hasCreateAccess"
      @edit="showEditTagModal"
      @duplicate="showDuplicateTagModal"
      @remove="showRemoveTagModal"
      @remove-selected="showRemoveSelectedTagsModal"
      @update:options="updateOptions"
    />
  </c-page>
</template>

<script>
import { omit } from 'lodash';
import { onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { pickIds } from '@/helpers/array';
import { convertQueryToRequest } from '@/helpers/query';
import { isImportedTag } from '@/helpers/entities/tag/entity';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useAlarmTag } from '@/hooks/store/modules/alarm-tag';

import TagsList from '@/components/other/tag/tags-list.vue';

export default {
  components: {
    TagsList,
  },

  setup() {
    const {
      alarmTags,
      alarmTagsPending,
      alarmTagsMeta,
      fetchAlarmTagsList,
      createAlarmTag,
      updateAlarmTag,
      removeAlarmTag,
      bulkRemoveAlarmTag,
    } = useAlarmTag();

    const { hasCreateAccess, hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.tag,
    );

    const { t } = useI18n();
    const modals = useModals();

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useLocalQueryWithOptions({
      onUpdate: (fetchQuery) => {
        const params = convertQueryToRequest(fetchQuery);
        params.with_flags = true;

        return fetchAlarmTagsList({ params });
      },
    });

    const showCreateTagModal = () =>
      modals.show({
        name: MODALS.createTag,
        config: {
          action: async (newTag) => {
            const result = await createAlarmTag({ data: newTag });

            await fetchList();

            return result;
          },
        },
      });

    const showEditTagModal = (tag) =>
      modals.show({
        name: MODALS.createTag,
        config: {
          title: t('modals.createTag.edit.title'),
          tag,
          isImported: isImportedTag(tag),
          action: async (newTag) => {
            const result = await updateAlarmTag({ id: tag._id, data: newTag });

            await fetchList();

            return result;
          },
        },
      });

    const showDuplicateTagModal = (tag) =>
      modals.show({
        name: MODALS.createTag,
        config: {
          title: t('modals.createTag.duplicate.title'),
          tag: omit(tag, ['_id']),
          action: async (newTag) => {
            const result = await createAlarmTag({ data: newTag });

            await fetchList();

            return result;
          },
        },
      });

    const showRemoveTagModal = (id) =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          text: t('tag.deleteConfirmation'),
          action: async () => {
            await removeAlarmTag({ id });

            return fetchList();
          },
        },
      });

    const showRemoveSelectedTagsModal = (selected = [], clearSelected) =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          text: t('tag.deleteConfirmation', selected.length),
          action: async () => {
            await bulkRemoveAlarmTag({ data: pickIds(selected) });

            clearSelected?.();

            return fetchList();
          },
        },
      });

    onMounted(fetchList);

    return {
      t,
      alarmTags,
      alarmTagsPending,
      alarmTagsMeta,
      options,
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,

      fetchList,
      updateOptions,
      showCreateTagModal,
      showEditTagModal,
      showDuplicateTagModal,
      showRemoveTagModal,
      showRemoveSelectedTagsModal,
    };
  },
};
</script>
