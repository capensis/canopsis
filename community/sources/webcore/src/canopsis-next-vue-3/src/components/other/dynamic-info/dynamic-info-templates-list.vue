<template>
  <div class="dynamic-info-templates-list-root">
    <c-advanced-data-table
      :options="options"
      :items="paginatedTemplates"
      :loading="pending"
      :headers="headers"
      :total-items="totalItems"
      :items-per-page-items="itemsPerPageItems"
      select-all
      advanced-pagination
      expand
      hide-actions
      @update:options="updateOptions"
    >
      <template #mass-actions="{ selected, clearSelected }">
        <c-table-mass-actions-panel
          :items="selected"
          :removable="removable"
          dynamic-info-template
          @clear:items="clearSelected"
          @refresh="$emit('refresh')"
        />
      </template>
      <template #expand="{ item }">
        <v-container>
          <v-card>
            <v-card-text>
              <dynamic-info-template-names-expand :names="item.names" />
            </v-card-text>
          </v-card>
        </v-container>
      </template>

      <template #actions="{ item }">
        <c-layout gap="1">
          <c-action-btn
            v-if="usable"
            :tooltip="$t('dynamicInfo.createFromTemplateTooltip')"
            icon="assignment"
            @click="$emit('use', item)"
          />
          <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
          <c-action-btn v-if="removable" type="delete" @click="$emit('remove', item._id)" />
        </c-layout>
      </template>
    </c-advanced-data-table>
  </div>
</template>

<script>
import { get, orderBy } from 'lodash';
import { computed, watch } from 'vue';

import { getMaxPageForItemsCount } from '@/helpers/pagination';

import { useI18n } from '@/hooks/i18n';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';

import DynamicInfoTemplateNamesExpand from './partials/dynamic-info-template-names-expand.vue';

export default {
  components: {
    DynamicInfoTemplateNamesExpand,
  },
  props: {
    templates: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    usable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['use', 'edit', 'remove', 'refresh'],

  setup(props) {
    const { t } = useI18n();

    const itemsPerPageItems = [5, 10, 25, 50, 100];

    const { options, updateOptions } = useLocalQueryWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: 10,
        search: '',
        sortBy: [],
        sortDesc: [],
      },
      onUpdate: () => {},
    });

    watch(
      () => props.templates.length,
      (length) => {
        const opts = options.value;
        const maxPage = getMaxPageForItemsCount(length, opts.itemsPerPage);

        if (opts.page > maxPage) {
          updateOptions({ ...opts, page: maxPage });
        }
      },
    );

    const headers = computed(() => [
      {
        text: t('common.title'),
        sortable: true,
        value: 'title',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
        width: '1%',
      },
    ]);

    const sortedTemplates = computed(() => {
      const { sortBy = [], sortDesc = [] } = options.value;

      if (!sortBy.length) {
        return props.templates;
      }

      const [sortItem] = sortBy;
      const sortKey = sortItem.key ?? sortItem;
      const isDesc = sortItem.order === 'desc' || sortDesc[0];

      return orderBy(
        props.templates,
        [item => String(get(item, sortKey, '') ?? '')],
        [isDesc ? 'desc' : 'asc'],
      );
    });

    const totalItems = computed(() => props.templates.length);

    const paginatedTemplates = computed(() => {
      const { page = 1, itemsPerPage = 10 } = options.value;
      const start = (page - 1) * itemsPerPage;

      return sortedTemplates.value.slice(start, start + itemsPerPage);
    });

    return {
      options,
      updateOptions,
      headers,
      itemsPerPageItems,
      totalItems,
      paginatedTemplates,
    };
  },
};
</script>

<style lang="scss" scoped>
.dynamic-info-templates-list-root {
  max-width: 720px;
  width: 100%;
}
</style>
