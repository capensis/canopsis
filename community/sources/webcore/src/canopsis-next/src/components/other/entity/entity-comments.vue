<template>
  <v-layout class="gap-3" column>
    <v-flex v-if="addable">
      <v-btn
        color="primary"
        @click="addComment"
      >
        {{ $t('common.add') }}
      </v-btn>
    </v-flex>
    <entity-comments-list
      :comments="comments"
      :pending="pending"
      :editable-first="editable && isFirstPage"
      @edit="editComment"
    />
    <c-table-pagination
      :total-items="meta.total_count"
      :items-per-page="query.itemsPerPage"
      :page="query.page"
      @update:page="updateQueryPage"
      @update:items-per-page="updateQueryItemsPerPage"
      @input="updateQuery"
    />
  </v-layout>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useLocalQuery } from '@/hooks/query/local-query';
import { usePendingHandler } from '@/hooks/query/pending';
import { useEntityComments } from '@/hooks/store/modules/entity-comment';

import EntityCommentsList from './partials/entity-comments-list.vue';

export default {
  components: { EntityCommentsList },
  props: {
    entity: {
      type: Object,
      required: true,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const comments = ref([]);
    const meta = ref({});

    const { t } = useI18n();
    const modals = useModals();
    const {
      createEntityComment,
      updateEntityComment,
      fetchEntityCommentsListWithoutStore,
    } = useEntityComments();

    const {
      pending,
      handler: fetchList,
    } = usePendingHandler(async (fetchQuery) => {
      const response = await fetchEntityCommentsListWithoutStore({
        params: {
          limit: fetchQuery.itemsPerPage,
          page: fetchQuery.page,
          entity: props.entity._id,
        },
      });

      comments.value = response.data;
      meta.value = response.meta;
    });

    const {
      query,
      updateQuery,
      updateQueryPage,
      updateQueryItemsPerPage,
    } = useLocalQuery({
      initialQuery: { page: 1, itemsPerPage: PAGINATION_LIMIT },
      onUpdate: fetchList,
    });

    const isFirstPage = computed(() => query.value.page === 1);

    const addComment = async () => {
      modals.show({
        name: MODALS.textEditor,
        config: {
          textarea: true,
          rules: { required: true },
          label: t('common.message'),
          action: async (message) => {
            await createEntityComment({ data: { entity: props.entity._id, message } });
            return fetchList(query.value);
          },
        },
      });
    };

    const editComment = (comment) => {
      modals.show({
        name: MODALS.textEditor,
        config: {
          textarea: true,
          rules: { required: true },
          label: t('common.message'),
          text: comment.message,
          action: async (message) => {
            await updateEntityComment({ id: comment._id, data: { entity: props.entity._id, message } });
            return fetchList(query.value);
          },
        },
      });
    };

    onMounted(() => fetchList(query.value));

    return {
      pending,
      comments,
      meta,
      query,
      isFirstPage,
      addComment,
      editComment,
      updateQuery,
      updateQueryPage,
      updateQueryItemsPerPage,
    };
  },
};
</script>
