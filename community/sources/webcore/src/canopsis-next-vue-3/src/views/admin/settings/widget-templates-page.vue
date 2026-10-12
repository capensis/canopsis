<template>
  <c-page
    :creatable="hasCreateAnyWidgetTemplateAccess"
    :create-tooltip="t('modals.createWidgetTemplate.create.title')"
    @refresh="fetchList"
    @create="showSelectWidgetTemplateTypeModal"
  >
    <widget-templates-list
      :widget-templates="widgetTemplates"
      :pending="pending"
      :options="options"
      :total-items="meta.total_count"
      :updatable="hasUpdateAnyWidgetTemplateAccess"
      :removable="hasDeleteAnyWidgetTemplateAccess"
      @update:options="options = $event"
      @edit="showEditWidgetTemplateModal"
      @remove="showRemoveWidgetTemplateModal"
    />
  </c-page>
</template>

<script>
import { onMounted, ref } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useCallActionWithPopup } from '@/hooks/actions/call';
import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useQueryOptions } from '@/hooks/query/options';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useWidgetTemplate } from '@/hooks/store/modules/widget-template';

import WidgetTemplatesList from '@/components/other/widget-template/widget-templates-list.vue';

export default {
  components: { WidgetTemplatesList },

  setup() {
    const widgetTemplates = ref([]);
    const meta = ref({});

    const { t } = useI18n();
    const modals = useModals();
    const { callActionWithPopup } = useCallActionWithPopup();

    const {
      hasCreateAccess: hasCreateAnyWidgetTemplateAccess,
      hasUpdateAccess: hasUpdateAnyWidgetTemplateAccess,
      hasDeleteAccess: hasDeleteAnyWidgetTemplateAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.widgetTemplate);

    const {
      createWidgetTemplate,
      updateWidgetTemplate,
      removeWidgetTemplate,
      fetchWidgetTemplatesListWithoutStore,
    } = useWidgetTemplate();

    const {
      query,
      pending,
      updateQuery,
      handler: fetchList,
    } = usePendingWithLocalQuery({
      initialQuery: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
      },
      fetchHandler: async (fetchQuery) => {
        const response = await fetchWidgetTemplatesListWithoutStore({
          params: convertQueryToRequest(fetchQuery),
        });

        widgetTemplates.value = response.data;
        meta.value = response.meta;
      },
    });

    const { options } = useQueryOptions(query, updateQuery);

    const showSelectWidgetTemplateTypeModal = () => {
      modals.show({
        name: MODALS.selectWidgetTemplateType,
        config: {
          action: newWidgetTemplate => callActionWithPopup(
              () => createWidgetTemplate({ data: newWidgetTemplate }),
              fetchList,
            ),
        },
      });
    };

    const showEditWidgetTemplateModal = (widgetTemplate) => {
      modals.show({
        name: MODALS.createWidgetTemplate,
        config: {
          widgetTemplate,
          action: newWidgetTemplate => callActionWithPopup(
              () => updateWidgetTemplate({ id: widgetTemplate._id, data: newWidgetTemplate }),
              fetchList,
            ),
        },
      });
    };

    const showRemoveWidgetTemplateModal = (id) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => callActionWithPopup(() => removeWidgetTemplate({ id }), fetchList),
        },
      });
    };

    onMounted(() => fetchList());

    return {
      t,
      widgetTemplates,
      meta,
      pending,
      options,
      hasCreateAnyWidgetTemplateAccess,
      hasUpdateAnyWidgetTemplateAccess,
      hasDeleteAnyWidgetTemplateAccess,
      fetchList,
      showSelectWidgetTemplateTypeModal,
      showEditWidgetTemplateModal,
      showRemoveWidgetTemplateModal,
    };
  },
};
</script>
