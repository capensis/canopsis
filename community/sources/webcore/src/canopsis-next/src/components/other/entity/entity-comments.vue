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
import {
  computed,
  ref,
  inject,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import Observer from '@/services/observer';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useEntityComments } from '@/hooks/store/modules/entity-comment';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';

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

    const periodicRefresh = inject('$periodicRefresh', new Observer());

    const {
      pending,
      query,
      updateQuery,
      updateQueryPage,
      updateQueryItemsPerPage,
      fetchHandlerWithQuery: fetchList,
    } = usePendingWithLocalQuery({
      initialQuery: { page: 1, itemsPerPage: PAGINATION_LIMIT },
      fetchHandler: async (fetchQuery) => {
        const response = await fetchEntityCommentsListWithoutStore({
          params: {
            limit: fetchQuery.itemsPerPage,
            page: fetchQuery.page,
            entity: props.entity._id,
          },
        });

        comments.value = response.data;
        meta.value = response.meta;
      },
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
            return fetchList();
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
            return fetchList();
          },
        },
      });
    };

    onMounted(() => {
      fetchList();

      periodicRefresh.register(fetchList);
    });
    onBeforeUnmount(() => periodicRefresh.unregister(fetchList));

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
