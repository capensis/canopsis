<template>
  <div class="theme-form">
    <c-name-field v-model="nameField" autofocus required />

    <v-tabs v-model="activeTab" color="primary" align-tabs="center">
      <v-tab value="main">{{ t('theme.tabs.mainColors') }}</v-tab>
      <v-tab value="widget-table">{{ t('theme.tabs.widgetTable') }}</v-tab>
      <v-tab value="severity">{{ t('theme.tabs.severity') }}</v-tab>
      <v-tab value="font-size">{{ t('theme.tabs.fontSize') }}</v-tab>
    </v-tabs>

    <v-window v-model="activeTab" :touchless="true" class="theme-form__window">
      <v-window-item value="main">
        <theme-main-colors-form />
      </v-window-item>
      <v-window-item value="widget-table">
        <theme-widget-table-colors-form />
      </v-window-item>
      <v-window-item value="severity">
        <theme-severity-colors-form />
      </v-window-item>
      <v-window-item value="font-size">
        <theme-font-size-form />
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { ref } from 'vue';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ThemeFontSizeForm from './theme-font-size-form.vue';
import ThemeMainColorsForm from './theme-main-colors-form.vue';
import ThemeSeverityColorsForm from './theme-severity-colors-form.vue';
import ThemeWidgetTableColorsForm from './theme-widget-table-colors-form.vue';

export default {
  components: {
    ThemeFontSizeForm,
    ThemeMainColorsForm,
    ThemeSeverityColorsForm,
    ThemeWidgetTableColorsForm,
  },

  setup() {
    const { t } = useI18n();
    const { nameField } = useFormModels({ nameField: 'name' });
    const activeTab = ref('main');

    return { t, nameField, activeTab };
  },
};
</script>

<style lang="scss">
.theme-form {
  min-width: 0;

  &__window {
    overflow: visible;
  }

  &__columns {
    margin: 0 -8px;
    padding: 16px 0 24px;
  }

  &__colors {
    width: 100%;
    padding: 16px 0 24px;

    .theme-color-picker-field,
    .theme-enabled-color-picker-field {
      align-items: center;
      width: 100%;
      flex: 0 0 100%;
      padding: 12px 0 12px 16px;

      &:not(:last-of-type) {
        border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
      }
    }
  }

  &__column {
    padding: 0 8px;

    & > * {
      padding: 10px 0 10px 16px;

      &:not(:last-child) {
        border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
      }
    }
  }

}
</style>
