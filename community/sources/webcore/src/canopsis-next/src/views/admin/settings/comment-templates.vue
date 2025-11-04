<template>
  <c-page
    :creatable="hasCreateAnyCommentTemplateAccess"
    :create-tooltip="$t('modals.createCommentTemplate.create.title')"
    @refresh="fetchList"
    @create="showCreateCommentTemplateModal"
  >
    <comment-templates-list
      :comment-templates="commentTemplates"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      :editable="hasUpdateAnyCommentTemplateAccess"
      :deletable="hasDeleteAnyCommentTemplateAccess"
      @edit="showEditCommentTemplateModal"
      @remove="showRemoveCommentTemplateModal"
    />
  </c-page>
</template>

<script>
import { ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCommentTemplates } from '@/hooks/store/modules/comment-template';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useCallActionWithPopup } from '@/hooks/actions/call';
import { useQueryOptions } from '@/hooks/query/options';
import { useCRUDPermissions } from '@/hooks/auth';

import CommentTemplatesList from '@/components/other/comment-template/comment-templates-list.vue';

export default {
  components: { CommentTemplatesList },
  setup() {
    const commentTemplates = ref([]);
    const meta = ref({});

    const { t } = useI18n();
    const modals = useModals();

    /**
     * PERMISSIONS
     */
    const {
      hasCreateAccess: hasCreateAnyCommentTemplateAccess,
      hasUpdateAccess: hasUpdateAnyCommentTemplateAccess,
      hasDeleteAccess: hasDeleteAnyCommentTemplateAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.commentTemplate);

    /**
     * STORE
     */
    const {
      createCommentTemplate,
      updateCommentTemplate,
      removeCommentTemplate,
      fetchCommentTemplatesListWithoutStore,
    } = useCommentTemplates();
    const { callActionWithPopup } = useCallActionWithPopup();

    /**
     * QUERY
     */
    const {
      query,
      pending,
      updateQuery,
      handler: fetchList,
    } = usePendingWithLocalQuery({
      initialQuery: { page: 1, itemsPerPage: PAGINATION_LIMIT },
      fetchHandler: async (fetchQuery) => {
        const response = await fetchCommentTemplatesListWithoutStore({
          params: {
            limit: fetchQuery.itemsPerPage,
            page: fetchQuery.page,
          },
        });

        commentTemplates.value = response.data;
        meta.value = response.meta;
      },
    });

    const { options } = useQueryOptions(query, updateQuery);

    /**
     * Show modal for creating a new comment template
     */
    const showCreateCommentTemplateModal = () => {
      modals.show({
        name: MODALS.createCommentTemplate,
        config: {
          action: newTemplate => callActionWithPopup(
            () => createCommentTemplate({ data: newTemplate }),
            fetchList,
          ),
        },
      });
    };

    /**
     * Show modal for editing an existing comment template
     *
     * @param {Object} template - The comment template object to edit
     */
    const showEditCommentTemplateModal = (template) => {
      modals.show({
        name: MODALS.createCommentTemplate,
        config: {
          template,
          title: t('modals.createCommentTemplate.edit.title'),

          action: newTemplate => callActionWithPopup(
            () => updateCommentTemplate({ id: template._id, data: newTemplate }),
            fetchList,
          ),
        },
      });
    };

    /**
     * Show confirmation modal for removing a comment template
     *
     * @param {string} id - The ID of the comment template to remove
     */
    const showRemoveCommentTemplateModal = (id) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => callActionWithPopup(
            () => removeCommentTemplate({ id }),
            fetchList,
          ),
        },
      });
    };

    onMounted(() => fetchList());

    return {
      hasCreateAnyCommentTemplateAccess,
      hasUpdateAnyCommentTemplateAccess,
      hasDeleteAnyCommentTemplateAccess,
      commentTemplates,
      meta,
      pending,
      options,

      updateQuery,
      showCreateCommentTemplateModal,
      showEditCommentTemplateModal,
      showRemoveCommentTemplateModal,
      fetchList,
    };
  },
};
</script>
