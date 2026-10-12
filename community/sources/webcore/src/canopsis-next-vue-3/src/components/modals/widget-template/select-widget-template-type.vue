<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ t('modals.selectWidgetTemplateType.title') }}</span>
    </template>

    <template #text>
      <c-layout column>
        <v-card
          v-for="{ value, text, icon } in availableTypes"
          :key="value"
          class="my-1 cursor-pointer"
          @click="selectType(value)"
        >
          <v-card-title class="py-3">
            <c-layout align="center" justify="space-between">
              <span class="text-subtitle-1">
                {{ text }}
              </span>

              <c-icon>{{ icon }}</c-icon>
            </c-layout>
          </v-card-title>
        </v-card>
      </c-layout>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed } from 'vue';

import { MODALS, WIDGET_TEMPLATES_TYPES, WIDGET_TEMPLATE_TYPES_TO_ICONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html?raw';

export default {
  name: MODALS.selectWidgetTemplateType,

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, modals, close } = useInnerModal(props);

    const availableTypes = computed(() => Object.values(WIDGET_TEMPLATES_TYPES).map(type => ({
        value: type,
        icon: WIDGET_TEMPLATE_TYPES_TO_ICONS[type] ?? 'description',
        text: t(`widgetTemplate.types.${type}`),
      })),
    );

    const selectType = (type) => {
      const defaultDataByType = {
        [WIDGET_TEMPLATES_TYPES.alarmExportToPdf]: {
          content: ALARM_EXPORT_PDF_TEMPLATE,
        },
      };
      const widgetTemplate = {
        type,
        ...defaultDataByType[type],
      };

      modals.show({
        name: MODALS.createWidgetTemplate,
        config: {
          widgetTemplate,
          action: config.value.action,
        },
      });

      close();
    };

    return {
      t,
      availableTypes,
      selectType,
    };
  },
};
</script>
