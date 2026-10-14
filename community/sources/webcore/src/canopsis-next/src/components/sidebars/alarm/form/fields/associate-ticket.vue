<template>
  <widget-settings-group :title="t('settings.associateTicket.title')">
    <v-container>
      <div class="d-flex flex-column gap-2">
        <div class="d-flex align-center">
          <span class="text-body-2">
            {{ t('settings.associateTicket.ticketFormats') }}
          </span>
          <c-help-icon
            :text="t('settings.associateTicket.formatsHelp')"
            icon-class="ml-2"
            max-width="400"
            top
          />
        </div>
        <span>{{ t('settings.associateTicket.description') }}</span>
        <c-widget-template-field
          :model-value="template"
          :templates="templates"
          :pending="templatesPending"
          @update:model-value="updateTemplate"
        />
        <span v-if="isSelectedTemplateEmpty">
          {{ t('settings.associateTicket.emptyTemplate') }}
        </span>
        <ticket-formats-form
          :model-value="modelValue"
          name="associateTicket.formats"
          @update:model-value="updateValue"
        />
      </div>
    </v-container>
  </widget-settings-group>
</template>

<script>
import { computed } from 'vue';

import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useWidgetTemplateField } from '@/hooks/widget/widget-template';

import TicketFormatsForm from '@/components/forms/ticket-formats-form.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

export default {
  components: { TicketFormatsForm, WidgetSettingsGroup },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    template: {
      type: [String, Symbol],
      default: '',
    },
    templates: {
      type: Array,
      default: () => [],
    },
    templatesPending: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'update:template'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateTemplate, updateValue } = useWidgetTemplateField(props, 'formats', emit);
    const isSelectedTemplateEmpty = computed(
      () => props.template !== CUSTOM_WIDGET_TEMPLATE && !props.modelValue.length,
    );

    return {
      t,
      isSelectedTemplateEmpty,
      updateTemplate,
      updateValue,
    };
  },
};
</script>
