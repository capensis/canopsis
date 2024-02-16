<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('modals.selectWidgetTemplateType.title') }}</span>
    </template>
    <template #text="">
      <v-layout column>
        <v-card
          class="my-1 cursor-pointer"
          v-for="{ value, text, icon } in availableTypes"
          :key="value"
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
import { MODALS, WIDGET_TEMPLATES_TYPES, COLUMNS_WIDGET_TEMPLATES_TYPES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWidgetTemplate,
  components: { ModalWrapper },
  mixins: [modalInnerMixin],
  computed: {
    availableTypes() {
      return Object.values(WIDGET_TEMPLATES_TYPES).map(type => ({
        value: type,
        icon: COLUMNS_WIDGET_TEMPLATES_TYPES.includes(type) ? 'view_week' : 'description',
        text: this.$t(`widgetTemplate.types.${type}`),
      }));
    },
  },
  methods: {
    selectType(type) {
      const TEMPLATE_TYPES_TO_DEFAULT_DATA = {
        [WIDGET_TEMPLATES_TYPES.alarmExportToPdf]: { content: ALARM_EXPORT_PDF_TEMPLATE },
      };
      const defaultData = TEMPLATE_TYPES_TO_DEFAULT_DATA[type];

      let widgetTemplate = { type };

      if (defaultData) {
        widgetTemplate = { ...widgetTemplate, ...defaultData };
      }

      this.$modals.show({
        name: MODALS.createWidgetTemplate,
        config: {
          widgetTemplate,

          title: this.$t('modals.createWidgetTemplate.create.title'),
          action: this.config.action,
        },
      });

      this.$modals.hide();
    },
  },
};
</script>
