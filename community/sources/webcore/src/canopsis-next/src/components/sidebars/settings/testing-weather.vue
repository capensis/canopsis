<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="form.title")
      v-divider
      field-periodic-refresh(v-model="form.parameters.periodic_refresh")
      v-divider
      field-storage(
        v-model="form.parameters.directory",
        :title="$t('settings.resultDirectory')",
        :disabled="form.parameters.is_api",
        @add="editResultDirectory",
        @edit="editResultDirectory",
        @remove="removeResultDirectory"
      )
      v-divider
      v-list-group
        template(#activator="")
          v-list-tile {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-switcher(
            v-model="form.parameters.is_api",
            :title="$t('settings.receiveByApi')"
          )
          v-divider
          field-storages(
            v-model="form.parameters.screenshot_directories",
            :disabled="form.parameters.is_api",
            :help-text="$t('settings.screenshotDirectories.helpText')",
            :title="$t('settings.screenshotDirectories.title')",
            @add="editScreenshotStorage",
            @edit="editScreenshotStorage"
          )
          v-divider
          field-storages(
            v-model="form.parameters.video_directories",
            :disabled="form.parameters.is_api",
            :help-text="$t('settings.videoDirectories.helpText')",
            :title="$t('settings.videoDirectories.title')",
            @add="editVideoStorage",
            @edit="editVideoStorage"
          )
          v-divider
          field-file-name-masks(v-model="form.parameters")
          v-divider
    v-btn.primary(
      :loading="submitting",
      :disabled="submitting",
      @click="submit"
    ) {{ $t('common.save') }}
</template>

<script>
import { MODALS, SIDE_BARS } from '@/constants';

import uid from '@/helpers/uid';

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
