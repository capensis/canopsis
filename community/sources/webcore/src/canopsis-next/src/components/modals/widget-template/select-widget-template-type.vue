<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('modals.selectWidgetTemplateType.title') }}</span>
    </template>
    <template #text="">
      <v-layout column>
        <v-card
          v-for="{ value, text, icon } in availableTypes"
          :key="value"
          class="my-1 cursor-pointer"
          @click="selectType(value)"
        >
          <v-card-title primary-title>
            <v-layout
              wrap
              justify-between
            >
              <v-flex xs11>
                <div class="text-subtitle-1">
                  {{ text }}
                </div>
              </v-flex>
              <v-flex>
                <v-icon>{{ icon }}</v-icon>
              </v-flex>
            </v-layout>
          </v-card-title>
        </v-card>
      </v-layout>
    </template>
  </modal-wrapper>
</template>

<script>
import { computed } from 'vue';

import { MODALS, WIDGET_TEMPLATES_TYPES, WIDGET_TEMPLATE_TYPES_TO_ICONS } from '@/constants';

import { useInnerModal, useModals } from '@/hooks/modals';
import { useI18n } from '@/hooks/i18n';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createWidgetTemplate,

  components: {
    ModalWrapper,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { config } = useInnerModal(props);
    const { t } = useI18n();
    const modals = useModals();

    const availableTypes = computed(() => Object.values(WIDGET_TEMPLATES_TYPES).map(type => ({
      value: type,
      icon: WIDGET_TEMPLATE_TYPES_TO_ICONS[type] ?? 'description',
      text: t(`widgetTemplate.types.${type}`),
    })));

    const selectType = (type) => {
      const TEMPLATE_TYPES_TO_DEFAULT_DATA = {
        [WIDGET_TEMPLATES_TYPES.alarmExportToPdf]: { content: ALARM_EXPORT_PDF_TEMPLATE },
      };
      const defaultData = TEMPLATE_TYPES_TO_DEFAULT_DATA[type];

      let widgetTemplate = { type };

      if (defaultData) {
        widgetTemplate = { ...widgetTemplate, ...defaultData };
      }

      modals.show({
        name: MODALS.createWidgetTemplate,
        config: {
          widgetTemplate,
          title: t('modals.createWidgetTemplate.create.title'),
          action: config.value.action,
        },
      });

      modals.hide();
    };

    return {
      availableTypes,
      selectType,
    };
  },
};
</script>
