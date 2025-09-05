<template>
  <div style="width: 700px">
    <template-testing-test-variables />
  </div>
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { isImportedTag } from '@/helpers/entities/tag/entity';
import { pickIds } from '@/helpers/array';

import { authMixin } from '@/mixins/auth';
import { localQueryMixin } from '@/mixins/query/query';
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';
import { permissionsTechnicalTagMixin } from '@/mixins/permissions/technical/tag';

import TemplateTestingTestVariables from '@/components/other/template-testing/template-testing-test-variables.vue';

export default {
  components: {
    TemplateTestingTestVariables,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    entitiesAlarmTagMixin,
    permissionsTechnicalTagMixin,
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
            await this.createAlarmTag({ data: newTag });

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
            await this.updateAlarmTag({ id: tag._id, data: newTag });

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
            await this.createAlarmTag({ data: newTag });

            return this.fetchList();
          },
        },
      });
    },

    showRemoveTagModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          text: this.$tc('tag.deleteConfirmation'),
          action: async () => {
            await this.removeAlarmTag({ id });

            return this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedTagsModal(selected = []) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          text: this.$tc('tag.deleteConfirmation', selected.length),
          action: async () => {
            await this.bulkRemoveAlarmTag({ data: pickIds(selected) });

            return this.fetchList();
          },
        },
      });
    },

    fetchList() {
      const params = this.getQuery();

      params.with_flags = true;

      return this.fetchAlarmTagsList({ params });
    },
  },
};
</script>
