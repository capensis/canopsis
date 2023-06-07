<template lang="pug">
  c-page(
    :creatable="hasCreateAnyTagAccess",
    :create-tooltip="$t('modals.createTag.create.title')",
    @refresh="fetchList",
    @create="showCreateTagModal"
  )
    tags-list(
      :tags="tags",
      :pending="tagsPending",
      :pagination.sync="pagination",
      :total-items="tagsMeta.total_count",
      :updatable="hasUpdateAnyTagAccess",
      :removable="hasDeleteAnyTagAccess",
      :duplicable="hasCreateAnyTagAccess",
      @edit="showEditTagModal",
      @remove="showRemoveTagModal",
      @duplicate="showDuplicateTagModal",
      @remove-selected="showDeleteSelectedTagsModal"
    )
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { isImportedTag } from '@/helpers/entities/tag/entity';
import { mapIds } from '@/helpers/array';

import { authMixin } from '@/mixins/auth';
import { permissionsTechnicalTagMixin } from '@/mixins/permissions/technical/tag';
import { entitiesTagMixin } from '@/mixins/entities/tag';
import { localQueryMixin } from '@/mixins/query-local/query';

import TagsList from '@/components/other/tag/tags-list.vue';

export default {
  components: {
    TagsList,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    permissionsTechnicalTagMixin,
    entitiesTagMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    showCreateTagModal() {
      this.$modals.show({
        name: MODALS.createTag,
        config: {
          action: async (newTag) => {
            await this.createTag({ data: newTag });

            return this.fetchList();
          },
        },
      });
    },

    showEditTagModal(tag) {
      this.$modals.show({
        name: MODALS.createTag,
        config: {
          title: this.$t('modals.createTag.edit.title'),
          tag,
          isImported: isImportedTag(tag),
          action: async (newTag) => {
            await this.updateTag({ id: tag._id, data: newTag });

            return this.fetchList();
          },
        },
      });
    },

    showDuplicateTagModal(tag) {
      this.$modals.show({
        name: MODALS.createTag,
        config: {
          title: this.$t('modals.createTag.duplicate.title'),
          tag: omit(tag, ['_id']),
          action: async (newTag) => {
            await this.createTag({ data: newTag });

            return this.fetchList();
          },
        },
      });
    },

    showRemoveTagModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeTag({ id });

            return this.fetchList();
          },
        },
      });
    },

    showDeleteSelectedTagsModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.bulkRemoveTags({ data: mapIds(selected) });

            return this.fetchList();
          },
        },
      });
    },

    fetchList() {
      return this.fetchTagsList({ params: this.getQuery() });
    },
  },
};
</script>
