<template>
  <widget-settings
    :submitting="submitting"
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <v-divider />
    <field-periodic-refresh v-model="form.parameters" />
    <v-divider />
    <field-storage
      v-model="form.parameters.directory"
      :title="$t('settings.resultDirectory')"
      :disabled="form.parameters.is_api"
      @add="editResultDirectory"
      @edit="editResultDirectory"
      @remove="removeResultDirectory"
    />
    <v-divider />
    <widget-settings-group :title="$t('settings.advancedSettings')">
      <field-switcher
        v-model="form.parameters.is_api"
        :title="$t('settings.receiveByApi')"
      />
      <v-divider />
      <field-storages
        v-model="form.parameters.screenshot_directories"
        :disabled="form.parameters.is_api"
        :help-text="$t('settings.screenshotDirectories.helpText')"
        :title="$t('settings.screenshotDirectories.title')"
        @add="editScreenshotStorage"
        @edit="editScreenshotStorage"
      />
      <v-divider />
      <field-storages
        v-model="form.parameters.video_directories"
        :disabled="form.parameters.is_api"
        :help-text="$t('settings.videoDirectories.helpText')"
        :title="$t('settings.videoDirectories.title')"
        @add="editVideoStorage"
        @edit="editVideoStorage"
      />
      <v-divider />
      <field-file-name-masks v-model="form.parameters" />
    </widget-settings-group>
    <v-divider />
  </widget-settings>
</template>

<script>
import { MODALS, SIDE_BARS } from '@/constants';

import { uid } from '@/helpers/uid';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesTestSuitesMixin } from '@/mixins/entities/test-suite';

import FieldTitle from '../form/fields/title.vue';
import FieldPeriodicRefresh from '../form/fields/periodic-refresh.vue';
import FieldSwitcher from '../form/fields/switcher.vue';
import WidgetSettings from '../partials/widget-settings.vue';
import WidgetSettingsGroup from '../partials/widget-settings-group.vue';

import FieldStorages from './form/fields/storages.vue';
import FieldStorage from './form/fields/storage.vue';
import FieldFileNameMasks from './form/fields/file-name-masks.vue';

export default {
  name: SIDE_BARS.testingWeatherSettings,
  components: {
    FieldTitle,
    FieldSwitcher,
    FieldPeriodicRefresh,
    FieldStorages,
    FieldStorage,
    FieldFileNameMasks,
    WidgetSettings,
    WidgetSettingsGroup,
  },
  mixins: [
    widgetSettingsMixin,
    entitiesTestSuitesMixin,
  ],
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
          this.form.parameters.directory = directory;
        },
      });
    },

    removeResultDirectory() {
      this.form.parameters.directory = '';
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
            this.form.parameters.screenshot_directories.splice(index, 1, directory);
          } else {
            this.form.parameters.screenshot_directories.push(this.getStorageByString(directory));
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
            this.form.parameters.video_directories.splice(index, 1, directory);
          } else {
            this.form.parameters.video_directories.push(this.getStorageByString(directory));
          }
        },
      });
    },
  },
};
</script>
