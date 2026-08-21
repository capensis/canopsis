<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />

    <alarms-list-modal-form
      v-model:value="form.parameters.alarmsList"
      :templates="preparedWidgetTemplates"
      :templates-pending="widgetTemplatesPending"
    />

    <widget-settings-group :title="t('settings.advancedSettings')">
      <field-filters
        v-model:filters="form.filters"
        :widget-id="widget._id"
        :widget-type="form.type"
        addable
        editable
        with-alarm
        with-entity
        with-pbehavior
        hide-selector
      />

      <field-opened-resolved-filter v-model="form.parameters.opened" />

      <field-switcher
        v-model="form.parameters.considerPbehaviors"
        :title="t('settings.considerPbehaviors.title')"
      />

      <field-criticity-levels v-model="form.parameters.criticityLevels" />

      <field-levels-colors-selector
        v-model="form.parameters.criticityLevelsColors"
        color-type="hex"
        hide-suffix
      />
    </widget-settings-group>
  </widget-settings>
</template>

<script>
import { computed, onMounted, ref, toRef } from 'vue';

import { MAX_LIMIT, SIDE_BARS } from '@/constants';

import { prepareWidgetTemplates } from '@/helpers/entities/widget/template/prepare';

import { useI18n } from '@/hooks/i18n';
import { useInfos } from '@/hooks/store/modules/infos';
import { useWidgetTemplate } from '@/hooks/store/modules/widget-template';
import { useWidgetSettings } from '@/hooks/widget/settings';

import AlarmsListModalForm from '@/components/sidebars/alarm/form/alarms-list-modal.vue';
import FieldOpenedResolvedFilter from '@/components/sidebars/alarm/form/fields/opened-resolved-filter.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import FieldCriticityLevels from './form/fields/criticity-levels.vue';
import FieldLevelsColorsSelector from './form/fields/levels-colors-selector.vue';

export default {
  name: SIDE_BARS.statsCalendarSettings,

  components: {
    AlarmsListModalForm,
    FieldCriticityLevels,
    FieldFilters,
    FieldLevelsColorsSelector,
    FieldOpenedResolvedFilter,
    FieldSwitcher,
    FieldTitle,
    WidgetSettings,
    WidgetSettingsGroup,
  },

  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const sidebar = toRef(props, 'sidebar');
    const widgetTemplates = ref([]);
    const widgetTemplatesPending = ref(false);
    const { fetchInfos } = useInfos();
    const { fetchWidgetTemplatesListWithoutStore } = useWidgetTemplate();

    const { form, widget, hasChanges, submitting, close, submit } = useWidgetSettings({ sidebar });

    const preparedWidgetTemplates = computed(() => prepareWidgetTemplates(widgetTemplates.value));

    const fetchWidgetTemplates = async () => {
      try {
        widgetTemplatesPending.value = true;

        const { data } = await fetchWidgetTemplatesListWithoutStore({ params: { limit: MAX_LIMIT } });

        widgetTemplates.value = data;
      } catch (err) {
        console.warn(err);
      } finally {
        widgetTemplatesPending.value = false;
      }
    };

    onMounted(() => {
      fetchInfos();
      fetchWidgetTemplates();
    });

    return {
      form,
      hasChanges,
      preparedWidgetTemplates,
      submitting,
      t,
      widget,
      widgetTemplatesPending,

      close,
      submit,
    };
  },
};
</script>
