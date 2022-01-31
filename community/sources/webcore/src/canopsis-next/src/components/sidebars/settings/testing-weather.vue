<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-periodic-refresh(v-model="settings.widget.parameters.periodic_refresh")
      v-divider
      field-storage(
        v-model="settings.widget.parameters.directory",
        :title="$t('settings.resultDirectory')",
        :disabled="settings.widget.parameters.is_api",
        @add="editResultDirectory",
        @edit="editResultDirectory",
        @remove="removeResultDirectory"
      )
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-switcher(
            v-model="settings.widget.parameters.is_api",
            :title="$t('settings.receiveByApi')"
          )
          v-divider
          field-storages(
            v-model="settings.widget.parameters.screenshot_directories",
            :disabled="settings.widget.parameters.is_api",
            :help-text="$t('settings.screenshotDirectories.helpText')",
            :title="$t('settings.screenshotDirectories.title')",
            @add="editScreenshotStorage",
            @edit="editScreenshotStorage"
          )
          v-divider
          field-storages(
            v-model="settings.widget.parameters.video_directories",
            :disabled="settings.widget.parameters.is_api",
            :help-text="$t('settings.videoDirectories.helpText')",
            :title="$t('settings.videoDirectories.title')",
            @add="editVideoStorage",
            @edit="editVideoStorage"
          )
          v-divider
          field-file-name-masks(v-model="settings.widget.parameters")
          v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, SIDE_BARS } from '@/constants';

import uid from '@/helpers/uid';
import { formToTestingWeatherWidget, testingWeatherWidgetToForm } from '@/helpers/forms/widgets/testing-weather';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesTestSuitesMixin } from '@/mixins/entities/test-suite';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldDuration from '@/components/sidebars/settings/fields/common/duration.vue';
import FieldPeriodicRefresh from '@/components/sidebars/settings/fields/common/periodic-refresh.vue';
import FieldSwitcher from '@/components/sidebars/settings/fields/common/switcher.vue';
import FieldStorages from '@/components/sidebars/settings/fields/testing-weather/storages.vue';
import FieldStorage from '@/components/sidebars/settings/fields/testing-weather/storage.vue';
import FieldFileNameMasks from '@/components/sidebars/settings/fields/testing-weather/file-name-masks.vue';

export default {
  name: SIDE_BARS.testingWeatherSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldTitle,
    FieldDuration,
    FieldSwitcher,
    FieldPeriodicRefresh,
    FieldStorages,
    FieldStorage,
    FieldFileNameMasks,
  },
  mixins: [
    widgetSettingsMixin,
    entitiesTestSuitesMixin,
  ],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: testingWeatherWidgetToForm(cloneDeep(widget)),
      },
    };
  },
  methods: {
    showDefineStorageModal({ title, field, action }) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title,
          field: {
            ...field,
            name: 'directory',
            validationRules: 'required',
          },
          action: async (directory) => {
            await this.validateTestSuitesDirectory({ data: { directory } });

            await action(directory);
          },
        },
      });
    },

    editResultDirectory(value) {
      this.showDefineStorageModal({
        title: this.$t('modals.defineStorage.title'),
        field: {
          value,
          placeholder: this.$t('modals.defineStorage.field.placeholder'),
        },
        action: (directory) => {
          this.settings.widget.parameters.directory = directory;
        },
      });
    },

    removeResultDirectory() {
      this.settings.widget.parameters.directory = '';
    },

    getStorageByString(directory) {
      return {
        directory: directory || '',
        key: uid(),
      };
    },

    editScreenshotStorage(storage, index) {
      this.showDefineStorageModal({
        title: this.$t('modals.defineScreenshotStorage.title'),
        field: {
          value: storage && storage.directory,
          placeholder: this.$t('modals.defineScreenshotStorage.field.placeholder'),
        },
        action: (directory) => {
          if (storage) {
            this.settings.widget.parameters.screenshot_directories.splice(index, 1, directory);
          } else {
            this.settings.widget.parameters.screenshot_directories.push(this.getStorageByString(directory));
          }
        },
      });
    },

    editVideoStorage(storage, index) {
      this.showDefineStorageModal({
        title: this.$t('modals.defineVideoStorage.title'),
        field: {
          value: storage && storage.directory,
          placeholder: this.$t('modals.defineVideoStorage.field.placeholder'),
        },
        action: (directory) => {
          if (storage) {
            this.settings.widget.parameters.video_directories.splice(index, 1, directory);
          } else {
            this.settings.widget.parameters.video_directories.push(this.getStorageByString(directory));
          }
        },
      });
    },

    prepareWidgetSettings() {
      const { widget } = this.settings;

      return formToTestingWeatherWidget(widget);
    },
  },
};
</script>
