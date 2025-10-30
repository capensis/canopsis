<template>
  <c-page
    :creatable="hasCreateAnyWidgetTemplateAccess"
    :create-tooltip="$t('modals.createWidgetTemplate.create.title')"
    @create="showSelectWidgetTemplateTypeModal"
    @refresh="fetchWidgetTemplatesListWithPreviousParams"
  >
    <template #header>
      {{ $tc('common.widgetTemplate', 2) }}
    </template>
    <widget-templates />
  </c-page>
</template>

<script>
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';
import { useWidgetTemplate } from '@/hooks/store/modules/widget-template';

import WidgetTemplates from '@/components/other/widget-template/widget-templates.vue';

export default {
  components: { WidgetTemplates },
  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const { fetchWidgetTemplatesListWithPreviousParams, createWidgetTemplate } = useWidgetTemplate();
    const {
      hasCreateAccess: hasCreateAnyWidgetTemplateAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.widgetTemplate);

    /**
     * Shows modal for selecting widget template type and creating a new widget template.
     * After successful creation, refreshes the widget templates list.
     */
    const showSelectWidgetTemplateTypeModal = () => modals.show({
      name: MODALS.selectWidgetTemplateType,
      config: {
        title: t('modals.createWidgetTemplate.create.title'),
        action: async (newWidgetTemplate) => {
          await createWidgetTemplate({ data: newWidgetTemplate });
          fetchWidgetTemplatesListWithPreviousParams();
        },
      },
    });

    return {
      hasCreateAnyWidgetTemplateAccess,
      showSelectWidgetTemplateTypeModal,
      fetchWidgetTemplatesListWithPreviousParams,
    };
  },
};
</script>
